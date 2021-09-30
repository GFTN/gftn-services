// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package quoteshandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/stellar/go/clients/horizon"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/quotes-service/handler/helper"
	modeladaptor "github.com/GFTN/gftn-services/quotes-service/utility/modelAdaptor"
	"github.com/GFTN/gftn-services/quotes-service/utility/nqsdbclient"
	"github.com/GFTN/gftn-services/quotes-service/utility/nqsmodel"
	"github.com/GFTN/gftn-services/quotes-service/utility/participantregistry"
	"github.com/GFTN/gftn-services/quotes-service/utility/whitelistservice"
	comn "github.com/GFTN/gftn-services/utility/common"
	"github.com/GFTN/gftn-services/utility/kafka"
	"github.com/GFTN/gftn-services/utility/response"
)

type QuoteHandler struct {
	PRClient         participantregistry.InterfaceClient
	HTTP             *http.Client
	DBClient         nqsdbclient.DatabaseClient
	WLSClient        whitelistservice.InterfaceClient
	HorizonClient    *horizon.Client
	GatewayOperation *kafka.KafkaOpreations
}

// RequestQuote recieve request from OFI and send corresponding quote request to RFIs
func (qh QuoteHandler) RequestQuote(w http.ResponseWriter, request *http.Request) {
	apiQuoteRequest := model.QuoteRequest{}
	err := json.NewDecoder(request.Body).Decode(&apiQuoteRequest)
	if err != nil {
		msg := "Unable to parse body of REST call for quote transaction: " + err.Error()
		LOGGER.Warningf(msg)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1002", err)
		return
	}
	err = apiQuoteRequest.Validate(strfmt.Default)
	if err != nil {
		msg := "Error validating structure of quote request " + err.Error()
		LOGGER.Warningf(msg)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1001", err)
		return
	}
	quoteRequestJsonb, _ := json.Marshal(apiQuoteRequest)
	LOGGER.Info("QuoteRequest: ", string(quoteRequestJsonb))

	//adapter
	nqsQuoteRequest := modeladaptor.QuoteRequestToNqs(&apiQuoteRequest)

	ofiDomain, err := helper.GetIdentity(request)
	if err != nil {
		LOGGER.Error("Error extracting identity from jwt token")
		response.NotifyWWError(w, request, http.StatusBadRequest, "COMMON-1006", err)
		return
	}
	nqsQuoteRequest.OfiId = &ofiDomain

	whiteListParticipants, err := qh.WLSClient.GetMutualWhiteListParticipants(*nqsQuoteRequest.OfiId)

	if err != nil {
		response.NotifyWWError(w, request, http.StatusNotFound, "COMMON-1007",
			err)
		return
	}
	// check if there is mutual whiteListParticipant available
	if len(whiteListParticipants) == 0 {
		response.NotifyWWError(w, request, http.StatusNotFound, "COMMON-1007",
			nil)
		return
	}
	AddressSendOfi, err := qh.PRClient.GetParticipantAccount(*nqsQuoteRequest.OfiId, comn.ISSUING)
	if err != nil {
		LOGGER.Error("Error geting AddressSendOfi from PR")
		response.NotifyWWError(w, request, http.StatusInternalServerError, "QUOTE-1006", err)
		return
	}
	AddressReceiveOfi := AddressSendOfi

	ofiReceiveAccount, err := qh.HorizonClient.LoadAccount(AddressReceiveOfi)
	if err != nil {
		LOGGER.Error("Error loading OFIRecieveAccount for horizon")
		response.NotifyWWError(w, request, http.StatusInternalServerError, "QUOTE-1002", err)
		return
	}
	for idx, _ := range whiteListParticipants {
		LOGGER.Info(*whiteListParticipants[idx].ID)
	}

	// generate uuid and timestamp
	requestID := uuid.Must(uuid.NewV4()).String()
	timeOfRequest := time.Now().Unix()
	timeOfRequestStr := timeOfRequest
	maxLimit := *nqsQuoteRequest.LimitMaxOfi
	minLimit := *nqsQuoteRequest.LimitMinOfi
	sourceAsset := *nqsQuoteRequest.SourceAsset
	targetAsset := *nqsQuoteRequest.TargetAsset
	timeExpireOfi := *nqsQuoteRequest.TimeExpireOfi
	sourceAssetJSON, _ := json.Marshal(sourceAsset)
	targetAssetJSON, _ := json.Marshal(targetAsset)
	// create quote request
	err = qh.DBClient.CreateRequest(requestID, ofiDomain, maxLimit, minLimit, sourceAssetJSON, targetAssetJSON, timeOfRequest, timeExpireOfi)
	if err != nil {
		LOGGER.Error(err)
		LOGGER.Error("Create Quote Request failed")
	}

	requestIDResponse := model.QuoteRequestReceipt{&requestID, &timeOfRequestStr}

	// reponse with quoteID
	requestIDResponseBytes, _ := json.Marshal(requestIDResponse)
	response.Respond(w, http.StatusOK, requestIDResponseBytes)
	LOGGER.Debug(*nqsQuoteRequest.OfiId)

	for idx, _ := range whiteListParticipants {
		go func(idx int) {
			LOGGER.Info("whitelisted participate undergoes asset pair screening: " + *whiteListParticipants[idx].ID)
			AddressSendRfi, err := qh.PRClient.GetParticipantAccount(*whiteListParticipants[idx].ID, comn.ISSUING)
			if err != nil {
				LOGGER.Error("Error geting AddressSendOfi from PR")
				response.NotifyWWError(w, request, http.StatusInternalServerError, "QUOTE-1006", err)
				return
			}
			AddressReceiveRfi := AddressSendRfi

			rfiReceiveAccount, err := qh.HorizonClient.LoadAccount(AddressReceiveRfi)
			if err != nil {
				LOGGER.Error(err)
				LOGGER.Error("Error loading RFIReceiveAccount for horizon")
				return
			}

			issuerAddressSourceAsset, err := qh.PRClient.GetParticipantAccount(*nqsQuoteRequest.SourceAsset.IssuerID, comn.ISSUING)
			if err != nil {
				LOGGER.Error(err)
				LOGGER.Error("Error geting Source Asset issuing address")
				return
			}
			nqsQuoteRequest.IssuerAddressSourceAsset = &issuerAddressSourceAsset

			//get source asset issuer address
			issuerAddressTargetAsset, err := qh.PRClient.GetParticipantAccount(*nqsQuoteRequest.TargetAsset.IssuerID, comn.ISSUING)
			if err != nil {
				LOGGER.Error(err)
				LOGGER.Error("Error geting Target Asset issuing address")
				return
			}
			nqsQuoteRequest.IssuerAddressTargetAsset = &issuerAddressTargetAsset

			// if source is DO, source DO issuer has to be either reciever or sender ,
			if *nqsQuoteRequest.SourceAsset.AssetType == nqsmodel.AssetAssetTypeDO {
				if AddressSendOfi != *nqsQuoteRequest.IssuerAddressSourceAsset &&
					AddressReceiveRfi != *nqsQuoteRequest.IssuerAddressSourceAsset {
					LOGGER.Info("invalid DO quote operation: source DO issuer is neither reciever or sender ")
					return
				}
			}
			// if target is DO,target DO issuer has to be either reciever or sender
			if *nqsQuoteRequest.TargetAsset.AssetType == nqsmodel.AssetAssetTypeDO {
				if AddressSendRfi != *nqsQuoteRequest.IssuerAddressTargetAsset &&
					AddressReceiveOfi != *nqsQuoteRequest.IssuerAddressTargetAsset {
					LOGGER.Info("invalid DO quote operation: target DO issuer is neither reciever or sender")
					return
				}
			}

			// if source is DA, RFI has to trust it
			if *nqsQuoteRequest.SourceAsset.AssetType == nqsmodel.AssetAssetTypeDA {
				flag := false
				for _, balance := range rfiReceiveAccount.Balances {
					if balance.Asset.Code == *nqsQuoteRequest.SourceAsset.AssetCode &&
						balance.Asset.Issuer == *nqsQuoteRequest.IssuerAddressSourceAsset {
						flag = true
						break
					}
				}
				if flag == false {
					LOGGER.Info("invalid quote operation : RFI do not trust target asset")
					return
				}
			}
			// if target is DA, OFI has to trust it
			if *nqsQuoteRequest.TargetAsset.AssetType == nqsmodel.AssetAssetTypeDA {
				flag := false
				for _, balance := range ofiReceiveAccount.Balances {
					if balance.Asset.Code == *nqsQuoteRequest.SourceAsset.AssetCode &&
						balance.Asset.Issuer == *nqsQuoteRequest.IssuerAddressTargetAsset {
						flag = true
						break
					}
				}
				if flag == false {
					LOGGER.Info("invalid quote operation : OFI do not trust target asset")
				}
			} //asset pair valid

			// insert empty quote for each quote_id-RFIdomain
			quoteStatus := 1
			rfiDomain := *whiteListParticipants[idx].ID
			quoteID := requestID + "-" + rfiDomain
			LOGGER.Info("Creating Quote for quoteID:", requestID, "RFIDomain:", rfiDomain)
			err = qh.DBClient.CreateQuote(requestID, quoteID, rfiDomain, ofiDomain, maxLimit, minLimit, sourceAssetJSON, targetAssetJSON, timeOfRequest, quoteStatus, timeExpireOfi)
			if err != nil {
				LOGGER.Error(err)
				LOGGER.Error("Create Quote failed")
			}
			// submit http quote request to RFI
			quoteRequestToRFI := model.QuoteRequestNotification{
				QuoteID:      &quoteID,
				QuoteRequest: &apiQuoteRequest,
			}
			jsonValue, err := json.Marshal(quoteRequestToRFI)
			if err != nil {
				LOGGER.Error(err)
				LOGGER.Error("Error Marshaling quoteRequestToRFI")
			}
			LOGGER.Info("submitting quote request to RFI: ", *whiteListParticipants[idx].ID)

			/*
				sending message to Kafka
			*/

			topicName := *whiteListParticipants[idx].ID + "_" + kafka.QUOTES_TOPIC
			LOGGER.Info(*whiteListParticipants[idx].ID, " topic name: ", topicName)

			err = qh.GatewayOperation.Produce(topicName, jsonValue)
			if err != nil {
				newError := errors.New("Kafka producer failed: " + err.Error())
				LOGGER.Error(newError.Error())
				LOGGER.Error("Error sending quote request to Kafka topic: " + topicName)
				return
			}

		}(idx)
	}
	return
}

func (qh QuoteHandler) GetQuotes(w http.ResponseWriter, request *http.Request) {
	participantDomain, err := helper.GetIdentity(request)
	if err != nil {
		LOGGER.Error(err)
	}
	LOGGER.Debug(participantDomain)
	vars := mux.Vars(request)
	requestID := vars["request_id"]
	if requestID == "" {
		err := errors.New("requestID empty")
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1011", err)
		return
	}
	LOGGER.Info("Get quotes for request_id: ", requestID)

	ofiDomain, err := helper.GetIdentity(request)
	if err != nil {
		LOGGER.Error("Error extracting identity from jwt token")
		response.NotifyWWError(w, request, http.StatusBadRequest, "COMMON-1006", err)
		return
	}

	quotes, err := qh.DBClient.GetQuotes(requestID, ofiDomain)

	if err != nil {
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusNotFound, "QUOTE-1011", err)
	}

	quoteStatuses := []model.QuoteStatus{}
	// LOGGER.Debug(*quotes[0].RequestID)
	for _, quote := range quotes {
		quoteStatus := modeladaptor.QuoteDBToQuoteStatus(quote)
		quoteStatuses = append(quoteStatuses, quoteStatus)
	}
	quoteStatusesBytes, _ := json.Marshal(quoteStatuses)

	response.Respond(w, http.StatusOK, quoteStatusesBytes)
}

// TODO testing; define caller, this shall be OFI
func (qh QuoteHandler) GetQuotesByAttributes(w http.ResponseWriter, request *http.Request) {
	LOGGER.Info("GetQuotesByAttributes got called")

	apiQuery := model.QuoteFilter{}
	err := json.NewDecoder(request.Body).Decode(&apiQuery)
	if err != nil {
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1002", err)
		return
	}
	err = apiQuery.Validate(strfmt.Default)
	if err != nil {
		msg := "Error validating structure of query" + err.Error()
		LOGGER.Warningf(msg)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1001", err)
		return
	}
	query := modeladaptor.QueryToQueryDB(&apiQuery)
	// get ofi identity
	ofiDomain, err := helper.GetIdentity(request)
	if err != nil {
		LOGGER.Error("Error extracting identity from jwt token")
		response.NotifyWWError(w, request, http.StatusBadRequest, "COMMON-1006", err)
		return
	}
	query.OfiID = &ofiDomain
	LOGGER.Info("Get quotes for attributes ")
	quotes, err := qh.DBClient.GetQuotesByAttributes(&query)
	if err != nil {
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1011", err)
	}

	quoteStatuses := []model.QuoteStatus{}
	for _, quote := range quotes {
		quoteStatus := modeladaptor.QuoteDBToQuoteStatus(quote)
		quoteStatuses = append(quoteStatuses, quoteStatus)
	}
	quoteStatusesJson, _ := json.Marshal(quoteStatuses)
	response.Respond(w, http.StatusOK, quoteStatusesJson)
}

func (qh QuoteHandler) GetQuotesByQuoteID(w http.ResponseWriter, request *http.Request) {
	LOGGER.Info("GetQuotesByQuoteID got called")
	vars := mux.Vars(request)
	quoteID := vars["quote_id"]
	rfiDomain, err := helper.GetIdentity(request)
	if err != nil {
		LOGGER.Error("Error extracting identity from jwt token")
		response.NotifyWWError(w, request, http.StatusBadRequest, "COMMON-1006", err)
		return
	}
	LOGGER.Info("quote_id :", quoteID)
	if quoteID == "" {
		err := errors.New("quoteID empty")
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1011", err)
		return
	}
	quoteDBs, err := qh.DBClient.GetQuoteByQuoteID(quoteID, rfiDomain)
	if err != nil {
		//TODO: quote error id
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusNotFound, "QUOTE-1011", err)
		return
	}

	quoteStatus := modeladaptor.QuoteDBToQuoteStatus(quoteDBs[0])
	quoteStatusJson, _ := json.Marshal(quoteStatus)
	response.Respond(w, http.StatusOK, []byte(quoteStatusJson))
	return
}

func (qh QuoteHandler) UpdateQuote(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	quoteID := vars["quote_id"]
	if quoteID == "" {
		err := errors.New("quoteID empty")
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1026", err)
		return
	}

	// quote := nqsmodel.NqsAssetPriceQuote{}
	quoteE := model.QuoteEnvelope{}
	err := json.NewDecoder(request.Body).Decode(&quoteE)
	if err != nil {
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1002", err)
		return
	}
	err = quoteE.Validate(strfmt.Default)
	if err != nil {
		msg := "Error validating structure of quote response" + err.Error()
		LOGGER.Warningf(msg)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1001", err)
		return
	}

	nqsquote, err := modeladaptor.QuoteResponseEnvelopeToNqs(&quoteE)
	if err != nil {
		msg := "Error validating structure of quote response " + err.Error()
		LOGGER.Warningf(msg)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1001", err)
		return
	}
	// get rfi identity
	rfiDomain, err := helper.GetIdentity(request)
	if err != nil {
		LOGGER.Error("Error extracting identity from jwt token")
		response.NotifyWWError(w, request, http.StatusBadRequest, "COMMON-1006", err)
		return
	}
	nqsquote.RfiId = &rfiDomain

	timeOfQuote := time.Now().Unix()
	nqsquote.QuoteID = &quoteID
	LOGGER.Info("Update quote for quote_id: ", quoteID)
	quoteDB := modeladaptor.NqsAssetPriceQuoteToQuoteDB(nqsquote)
	err = qh.DBClient.UpdateQuote(quoteDB, timeOfQuote)
	if err != nil {
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1026", err)
		return
	}
	response.Respond(w, http.StatusOK, []byte(strconv.Itoa(int(timeOfQuote))))
}

func (qh QuoteHandler) CancelQuote(w http.ResponseWriter, request *http.Request) {
	LOGGER.Info("CancelQuote got called")
	vars := mux.Vars(request)
	quoteID := vars["quote_id"]
	rfiDomain, err := helper.GetIdentity(request)
	if err != nil {
		LOGGER.Error("Error extracting identity from jwt token")
		response.NotifyWWError(w, request, http.StatusBadRequest, "COMMON-1006", err)
		return
	}
	LOGGER.Info("quote_id :", quoteID)
	if quoteID == "" {
		err := errors.New("quoteID empty")
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1027", err)
		return
	}
	timeCancel := time.Now().Unix()
	err = qh.DBClient.CancelQuote(quoteID, rfiDomain, timeCancel)
	if err != nil {
		//TODO: quote error id
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusNotFound, "QUOTE-1027", err)
		return
	}
	response.Respond(w, http.StatusOK, []byte(strconv.Itoa(int(timeCancel))))
	return
}

func (qh QuoteHandler) CancelQuotesByAttributes(w http.ResponseWriter, request *http.Request) {
	LOGGER.Info("CancelQuotesByAttributes got called")
	apiQuery := model.QuoteFilter{}
	err := json.NewDecoder(request.Body).Decode(&apiQuery)
	if err != nil {
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1002", err)
		return
	}
	err = apiQuery.Validate(strfmt.Default)
	if err != nil {
		msg := "Error validating structure of query" + err.Error()
		LOGGER.Warningf(msg)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1001", err)
		return
	}
	query := modeladaptor.QueryToQueryDB(&apiQuery)
	// get rfi identity
	rfiDomain, err := helper.GetIdentity(request)
	if err != nil {
		LOGGER.Error("Error extracting identity from jwt token")
		response.NotifyWWError(w, request, http.StatusBadRequest, "COMMON-1006", err)
		return
	}
	query.RfiID = &rfiDomain
	LOGGER.Info("Cancel quotes by attributes ")
	TimeCancel := time.Now().Unix()
	quotes, err := qh.DBClient.CancelQuotesByAttributes(&query, TimeCancel)

	if err != nil {
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1027", err)
	}

	quoteStatuses := []model.QuoteStatus{}
	for _, quote := range quotes {
		quoteStatus := modeladaptor.QuoteDBToQuoteStatus(quote)
		quoteStatuses = append(quoteStatuses, quoteStatus)
	}
	quoteStatusesJson, _ := json.Marshal(quoteStatuses)
	response.Respond(w, http.StatusOK, quoteStatusesJson)
}

func (qh QuoteHandler) ExecutedQuote(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	quoteID := vars["quote_id"]
	rfiDomain := vars["rfi_id"]

	LOGGER.Info("ExecutedQuote...")
	if quoteID == "" {
		err := errors.New("quoteID empty")
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1028", err)
		return
	}
	timeExecuted := time.Now().Unix()
	err := qh.DBClient.ExecutedQuote(quoteID, rfiDomain, timeExecuted)
	if err != nil {
		//TODO: quote error id
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1028", err)
		return
	}
	response.Respond(w, http.StatusOK, []byte(strconv.Itoa(int(timeExecuted))))
	return
}
