// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package exchangehandler

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/stellar/go/clients/horizon"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/quotes-service/handler/helper"
	"github.com/GFTN/gftn-services/quotes-service/utility/cryptoservice"
	"github.com/GFTN/gftn-services/quotes-service/utility/gasservice"
	modeladaptor "github.com/GFTN/gftn-services/quotes-service/utility/modelAdaptor"
	"github.com/GFTN/gftn-services/quotes-service/utility/nqsdbclient"
	"github.com/GFTN/gftn-services/quotes-service/utility/participantregistry"
	"github.com/GFTN/gftn-services/quotes-service/utility/whitelistservice"
	"github.com/GFTN/gftn-services/utility/blockchain-adaptor/ww_stellar"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/participant"
	"github.com/GFTN/gftn-services/utility/response"
)

// ExchangeHandler : Holds necessary parameters used to create an atomic exchange
type ExchangeHandler struct {
	HTTP             *http.Client
	GasServiceClient gasservice.GasServiceClient
	HorizonClient    *horizon.Client
	CSClient         cryptoservice.InterfaceClient
	WLSClient        whitelistservice.InterfaceClient
	PRClient         participantregistry.InterfaceClient
	DBClient         nqsdbclient.DatabaseClient
}

// CreateAtomicExchange : creates atomic exchange between 2 different assets based on selected quotes
func (eh *ExchangeHandler) CreateAtomicExchange(w http.ResponseWriter, request *http.Request) {

	exchangeRequestEnvelope := model.ExchangeEnvelope{}
	err := json.NewDecoder(request.Body).Decode(&exchangeRequestEnvelope)
	if err != nil {
		msg := "Unable to parse body of REST call for exchange transaction: " + err.Error()
		LOGGER.Warningf(msg)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1002", err)
		return
	}
	err = exchangeRequestEnvelope.Validate(strfmt.Default)
	if err != nil {
		msg := "Error validating structure of exchange request " + err.Error()
		LOGGER.Warningf(msg)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1001", err)
		return
	}
	// Log exchange request
	LOGGER.Info("ExchangeRequest:", *exchangeRequestEnvelope.Exchange)
	LOGGER.Info("ExchangeRequestSignature:", *exchangeRequestEnvelope.Signature)

	//map model.exchange to nqsmodel.exchange
	nqsExchangeRequest, err := modeladaptor.ExchangeRequestEnvelopeToNqs(exchangeRequestEnvelope)
	if err != nil {
		msg := "Error validating structure of exchange request " + err.Error()
		LOGGER.Warningf(msg)
		response.NotifyWWError(w, request, http.StatusBadRequest, "QUOTE-1001", err)
		return
	}
	//verify OFI signature?

	//get caller identity
	ofiDomain, err := helper.GetIdentity(request)
	if err != nil {
		LOGGER.Error("Error extracting identity from jwt token")
		response.NotifyWWError(w, request, http.StatusBadRequest, "COMMON-1006", err)
		return
	}

	// reject if ofiID != jwt caller id
	if *nqsExchangeRequest.OfiId != ofiDomain {
		LOGGER.Error("ofiID is not equal caller identity from jwt token")
		response.NotifyWWError(w, request, http.StatusForbidden, "EXCHANGE-1021", err)
		return
	}

	// TODO: go routine
	//get OFI details from PR
	paricipantOfi, err := eh.PRClient.GetParticipantForDomain(*nqsExchangeRequest.OfiId)
	if err != nil {
		LOGGER.Error("Error getting Participant")
		response.NotifyWWError(w, request, http.StatusBadRequest, "EXCHANGE-1007", err)
		return
	}

	// get RFI details for PR
	paricipantRfi, err := eh.PRClient.GetParticipantForDomain(*nqsExchangeRequest.RfiId)
	if err != nil {
		LOGGER.Error("Error getting Participant")
		response.NotifyWWError(w, request, http.StatusBadRequest, "EXCHANGE-1007", err)
		return
	}

	// Check all participants are active before quote execution
	LOGGER.Info("Check all participants active")
	err = participant.CheckStatusActive(paricipantOfi, paricipantRfi)
	if err != nil {
		msg := err.Error()
		LOGGER.Error(msg)
		response.NotifyFailure(w, request, http.StatusBadRequest, msg)
		return
	}

	//Executing quote: update quotes status to quote_executing
	LOGGER.Info("Executing Quote")
	isExecuted := false
	var timeExecuted int64
	quoteResponse, err := json.Marshal(*nqsExchangeRequest.ExchangeRequestDecode.Quote)
	timeExecuting := time.Now().Unix()
	err = eh.DBClient.ExecutingQuote(*nqsExchangeRequest.QuoteID, *nqsExchangeRequest.OfiId, quoteResponse, timeExecuting, *nqsExchangeRequest.Amount)
	if err != nil {
		LOGGER.Error(err)
		LOGGER.Warning("Error initiate Executing Quote: " + *nqsExchangeRequest.QuoteID)
		response.NotifyWWError(w, request, http.StatusBadRequest, "EXCHANGE-1013", err)
		return
	}
	defer func() {
		if isExecuted == false {
			err := eh.DBClient.FailedQuote(*nqsExchangeRequest.QuoteID, *nqsExchangeRequest.RfiId)
			if err != nil {
				LOGGER.Error(err)
				LOGGER.Error("Error marking failed Quote: " + *nqsExchangeRequest.QuoteID)
				// response.NotifyWWError(w, request, http.StatusBadRequest, "EXCHANGE-1017", err)
				return
			}
			LOGGER.Warning("Quote execution failed: " + *nqsExchangeRequest.QuoteID)
			// response.NotifyWWError(w, request, http.StatusBadRequest, "EXCHANGE-1018", err)
			return

		} else {
			err := eh.DBClient.ExecutedQuote(*nqsExchangeRequest.QuoteID, *nqsExchangeRequest.RfiId, timeExecuted)
			if err != nil {
				LOGGER.Error(err)
				LOGGER.Error("Error marking Executed Quote: " + *nqsExchangeRequest.QuoteID)
				// response.NotifyWWError(w, request, http.StatusBadRequest, "EXCHANGE-1019", err)
				return
			}
		}
	}()

	//TODO: go
	//check if whitelisted
	wlparticipantIDs, err := eh.WLSClient.GetMutualWhiteListParticipantDomains(*nqsExchangeRequest.OfiId)
	if err != nil {
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusInternalServerError, "COMMON-1002", err)
		return
	}
	for _, wlparticipantID := range wlparticipantIDs {
		if wlparticipantID == *nqsExchangeRequest.RfiId {
			break
		} else {
			continue
		}
		msg := "RFI not in OFI's whitelist " + err.Error()
		LOGGER.Warningf(msg)
		response.NotifyWWError(w, request, http.StatusBadRequest, "EXCHANGE-1006", err)
		return
	}
	//calculate target amount = source amount * exchange rate
	amountTargetAsset := (*nqsExchangeRequest.Amount).Mul(*nqsExchangeRequest.ExchangeRate)
	// Log the exchange details
	LOGGER.Info("Exchange Details: Quote ID: ", *nqsExchangeRequest.QuoteID)
	targetAsset, _ := json.Marshal(*nqsExchangeRequest.TargetAsset)
	sourceAsset, _ := json.Marshal(*nqsExchangeRequest.SourceAsset)
	LOGGER.Info("Exchange Details: Target Asset: ", string(targetAsset))
	LOGGER.Info("Exchange Details: Source Asset: ", string(sourceAsset))
	LOGGER.Info("Exchange Details: Source Asset amount:", *nqsExchangeRequest.Amount)
	LOGGER.Info("Exchange Details: Exchange Rate:", *nqsExchangeRequest.ExchangeRate)
	LOGGER.Info("Exchange Details: Target Asset amount:", amountTargetAsset)

	// Construct payment structs
	payments := make([]ww_stellar.Payment, 2)
	payments[0] = ww_stellar.Payment{
		ww_stellar.Account{*nqsExchangeRequest.OfiId, *nqsExchangeRequest.AccountSendOfi},
		ww_stellar.Account{*nqsExchangeRequest.RfiId, *nqsExchangeRequest.AccountReceiveRfi},
		model.Asset{
			nqsExchangeRequest.SourceAsset.AssetCode,
			nqsExchangeRequest.SourceAsset.AssetType,
			*nqsExchangeRequest.SourceAsset.IssuerID,
		},
		*nqsExchangeRequest.Amount,
	}
	payments[1] = ww_stellar.Payment{
		ww_stellar.Account{*nqsExchangeRequest.RfiId, *nqsExchangeRequest.AccountSendRfi},
		ww_stellar.Account{*nqsExchangeRequest.OfiId, *nqsExchangeRequest.AccountReceiveOfi},
		model.Asset{
			nqsExchangeRequest.TargetAsset.AssetCode,
			nqsExchangeRequest.TargetAsset.AssetType,
			*nqsExchangeRequest.TargetAsset.IssuerID,
		},
		amountTargetAsset, //source amount * exchange rate
	}
	// Sanity Check
	err = ww_stellar.SanityCheck(payments, eh.PRClient, eh.HorizonClient)
	if err != nil {
		LOGGER.Error(err)
		msg := "Sanity Check: " + err.Error()
		if wwErr, ok := err.(ww_stellar.Error); ok {
			// return status code 400 if err related to account does not exist || does not trust asset|| insufficient balance || trust limit exceeded
			switch wwErr.Code() {
			case ww_stellar.ERROR_ACCOUNT_DOES_NOT_TRUST_ASSET:
				response.NotifyWWError(w, request, http.StatusBadRequest, "COMMON-1301", wwErr)
			case ww_stellar.ERROR_ACCOUNT_TRUST_LIMIT_EXCEEDED:
				response.NotifyWWError(w, request, http.StatusBadRequest, "COMMON-1302", wwErr)
			case ww_stellar.ERROR_INSUFFICIENT_BALANCE:
				response.NotifyWWError(w, request, http.StatusBadRequest, "COMMON-1303", wwErr)
			case ww_stellar.ERROR_MAP_STELLAR_ADDRESS:
				response.NotifyWWError(w, request, http.StatusBadRequest, "COMMON-1304", wwErr)
			case ww_stellar.ERROR_AMOUNT_LESSTHEN_EQUAL_ZERO:
				response.NotifyWWError(w, request, http.StatusBadRequest, "COMMON-1308", wwErr)
			case ww_stellar.ERROR_NO_ISSUER_INVOLVED:
				response.NotifyWWError(w, request, http.StatusBadRequest, "COMMON-1309", wwErr)
			default:
				response.NotifyFailure(w, request, http.StatusInternalServerError, msg)
			}
			return
		}
		response.NotifyFailure(w, request, http.StatusInternalServerError, msg)
		return
	}

	// construct transcation envelop
	tc := ww_stellar.StellarTransactionConstructor{}
	err = tc.InitTransaction(eh.GasServiceClient)
	if err != nil {
		msg := "Error init transaction"
		LOGGER.Error(msg)
		response.NotifyFailure(w, request, http.StatusInternalServerError, msg)
		return
	}
	type PaymentBuilderResult struct {
		Error error
	}
	ch := make(chan PaymentBuilderResult)
	// get payment concurrently and add to transaction
	for idx := range payments {
		go func(ch chan PaymentBuilderResult, idx int) {
			err := tc.AddPayment(payments[idx], eh.PRClient)
			ch <- PaymentBuilderResult{err}
			return
		}(ch, idx)
	}
	for _ = range payments {
		select {
		case paymentBuilderResult := <-ch:
			if paymentBuilderResult.Error != nil {
				err := paymentBuilderResult.Error
				LOGGER.Error(err)
				msg := "Error: Error when constructing transaction: " + err.Error()
				LOGGER.Info(msg)
				// return status code 400 if err related to account does not exist
				if wwErr, ok := err.(ww_stellar.Error); ok {
					switch wwErr.Code() {
					case ww_stellar.ERROR_MAP_STELLAR_ADDRESS:
						response.NotifyWWError(w, request, http.StatusBadRequest, "COMMON-1304", wwErr)
					default:
						response.NotifyFailure(w, request, http.StatusInternalServerError, msg)
					}
					return
				}
				response.NotifyFailure(w, request, http.StatusInternalServerError, msg)
				return

			} else {
				LOGGER.Info("Transaction construction: Payment Added")
			}
		case <-time.After(time.Second * 10):
			msg := "Error: Timeout when constructing transaction"
			LOGGER.Info(msg)
			response.NotifyFailure(w, request, http.StatusInternalServerError, msg)
			LOGGER.Error(err)
			return
			break
			// call timed out
		}
	}
	close(ch)
	//set fee
	tc.SetFee()
	//set expiry time, allow 60 seconds for acquiring signatures and submitting to stellar
	tc.SetTimeBound(0, time.Now().Unix()+60)

	LOGGER.Info("Result atomic swap transaction envelop of ", *nqsExchangeRequest.OfiId, "and", *nqsExchangeRequest.RfiId, ":", tc.Base64())
	// sign transcation and concat signatures concurrently
	type SignTransactionResult struct {
		Error error
	}
	chOFI := make(chan SignTransactionResult)
	chRFI := make(chan SignTransactionResult)
	// get OFI's signature
	go func(ch chan SignTransactionResult) {
		err = tc.SignTransactionGlobalAndAppend(*nqsExchangeRequest.AccountSendOfi, eh.CSClient, paricipantOfi, *nqsExchangeRequest.ExchangeRequestBase64, *nqsExchangeRequest.ExchangeRequestSignBase64)
		ch <- SignTransactionResult{err}
		return
	}(chOFI)
	// get RFI's signature
	go func(ch chan SignTransactionResult) {
		// retrieve signed quote response from DB
		LOGGER.Debug(*nqsExchangeRequest.QuoteID)
		nqsQuoteResponse, err := eh.DBClient.GetQuoteByQuoteID(*nqsExchangeRequest.QuoteID, *nqsExchangeRequest.RfiId)
		if err != nil {
			LOGGER.Error(err)
			response.NotifyWWError(w, request, http.StatusNotFound, "EXCHANGE-1013", err)
			return
		}
		QuoteResponseBase64 := *nqsQuoteResponse[0].QuoteResponseBase64
		QuoteResponseSignatureBase64 := *nqsQuoteResponse[0].QuoteResponseSignature
		err = tc.SignTransactionGlobalAndAppend(*nqsExchangeRequest.AccountSendRfi, eh.CSClient, paricipantRfi, QuoteResponseBase64, QuoteResponseSignatureBase64)
		ch <- SignTransactionResult{err}
		return
	}(chRFI)
	for _ = range payments {
		select {
		case SignTransactionResult := <-chOFI:
			if SignTransactionResult.Error != nil {
				err := SignTransactionResult.Error
				LOGGER.Error(err)
				msg := "Error: Error when acquiring OFI signature in Signing transaction"
				LOGGER.Error(msg)
				response.NotifyFailure(w, request, http.StatusInternalServerError, msg)
				return
			} else {
				LOGGER.Info("Transaction construction: OFI Signature Added")
			}
		case SignTransactionResult := <-chRFI:
			if SignTransactionResult.Error != nil {
				err := SignTransactionResult.Error
				LOGGER.Error(err)
				msg := "Error: Error when acquiring RFI signature in Signing transaction"
				LOGGER.Error(msg)
				response.NotifyFailure(w, request, http.StatusInternalServerError, msg)
				return
			} else {
				LOGGER.Info("Transaction construction: RFI Signature Added")
			}
		case <-time.After(time.Second * 10):
			msg := "Error: Timeout when signing transaction"
			LOGGER.Error(msg)
			response.NotifyFailure(w, request, http.StatusInternalServerError, msg)
			return
			// call timed out
		}
	}
	close(chOFI)
	close(chRFI)
	LOGGER.Info("Signed result transaction of ", nqsExchangeRequest.OfiId, "and", nqsExchangeRequest.RfiId, ":", tc.Base64())

	// submit transaction
	hash, ledger, err := tc.SendTransaction(eh.GasServiceClient)
	if err != nil {
		msg := "Error: Sending transaction to Gas Service Failed"
		LOGGER.Error(msg)
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusInternalServerError, "API-1306", err)
		return
	}
	LOGGER.Info("Transaction submitted successfully, hash: " + hash + " ledger: " + strconv.FormatUint(ledger, 10))

	isExecuted = true //executed successfully
	timeExecuted = time.Now().Unix()
	exchangeReceipt := model.ExchangeReceipt{}
	exchangeReceipt.Exchange = nqsExchangeRequest.ExchangeRequestDecode
	exchangeReceipt.TransactedAmountTarget = &amountTargetAsset
	exchangeReceipt.TransactedAmountSource = nqsExchangeRequest.ExchangeRequestDecode.Amount
	exchangeReceipt.TransactionHash = &hash
	statusExchange := "OK"
	exchangeReceipt.StatusExchange = &statusExchange
	exchangeReceipt.TimeExecuted = timeExecuted

	// Constructing Exchange Log for Firebase
	ExchangeLog := ExchangeLog{
		ExchangeReceipt: exchangeReceipt,
		OFIID:           ofiDomain,
		RFIID:           *nqsExchangeRequest.RfiId,
	}

	// Sending log to Firebase
	LOGGER.Info("Sending exchange log to Firebase")
	err = SendExchangeLogToFirebase(ExchangeLog)
	if err != nil {
		LOGGER.Error(err)
	}

	resBody, err := json.Marshal(exchangeReceipt)
	if err != nil {
		LOGGER.Error(err)
	}
	response.Respond(w, http.StatusOK, resBody)
	return
}

// SendExchangeLogToFirebase uses helper function to
// send exchange log to firebase
func SendExchangeLogToFirebase(log ExchangeLog) error {

	// separate transactions by release and participantId of the OFI
	release := os.Getenv(global_environment.ENV_KEY_SERVICE_VERSION)
	ofiID := log.OFIID
	rfiID := log.RFIID

	// constructing db ref for firebase
	fbRefOfi := release + "/txn/exchange/" + ofiID
	fbRefRfi := release + "/txn/exchange/" + rfiID

	// using generic helper to log to firebase
	err := helper.SendLogToFirebase(log, fbRefOfi, fbRefRfi)
	if err != nil {
		return err
	}

	return nil
}

// ExchangeLog holds relevant non-PII Exchange
// transaction data for view by the client
// This is currently logged to Firebase,
// but can be reused in other places in the future
type ExchangeLog struct {
	ExchangeReceipt model.ExchangeReceipt
	OFIID           string
	RFIID           string
}
