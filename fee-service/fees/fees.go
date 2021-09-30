// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package fees

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/GFTN/gftn-services/fee-service/environment"
	"github.com/GFTN/gftn-services/gftn-models/model"
	pr_client "github.com/GFTN/gftn-services/participant-registry-client/pr-client"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/kafka"
	participant_checks "github.com/GFTN/gftn-services/utility/participant"
	"github.com/GFTN/gftn-services/utility/response"
)

type FeeOperations struct {
	prClient         pr_client.PRServiceClient
	HTTPClient       *http.Client
	GatewayOperation *kafka.KafkaOpreations
}

func CreateFeeOperations() (FeeOperations, error) {

	op := FeeOperations{}
	var prClient pr_client.PRServiceClient
	var err error
	if os.Getenv(environment.ENV_KEY_PARTICIPANT_REGISTRY_SERVICE_CLIENT) == "mock" {
		prClient, err = pr_client.CreateMockPRServiceClient(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL))
	} else {
		prClient, err = pr_client.CreateRestPRServiceClient(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL))
	}

	if err != nil {
		LOGGER.Errorf(" Error CreatePRServiceClient failed  %v", err)
		return op, err
	}
	op.prClient = prClient
	hc := &http.Client{Timeout: time.Second * 10}
	op.HTTPClient = hc

	LOGGER.Infof("Initiate Kafka producer for ww-gateway")
	op.GatewayOperation, err = kafka.Initialize()
	if err != nil {
		LOGGER.Errorf("Initialize Kafka producer for ww-gateway failed: %s", err.Error())
		return FeeOperations{}, err
	}
	return op, nil
}

func (fo FeeOperations) CalculateFees(w http.ResponseWriter, request *http.Request) {

	LOGGER.Infof("*******In CalculateFees function******")
	vars := mux.Vars(request)

	rfiId := strings.TrimSpace(vars["participant_id"])

	if rfiId == "" {
		err := errors.New("Participant ID is empty")
		LOGGER.Errorf(err.Error())
		response.NotifyWWError(w, request, http.StatusBadRequest, "FEES-1001", err)
		return
	}
	participant, err := fo.getParticipantForDomain(rfiId)
	if err != nil {
		LOGGER.Errorf("Error while getting Participant for Domain:", err.Error())
		response.NotifyWWError(w, request, http.StatusBadRequest, "FEES-1001", err)
		return
	}

	// Check if participant is active
	LOGGER.Info("Check participant active")
	err = participant_checks.CheckStatusActive(participant)
	if err != nil {
		msg := err.Error()
		LOGGER.Error(msg)
		response.NotifyFailure(w, request, http.StatusBadRequest, msg)
		return
	}

	var feesAndAmountRequest model.FeesRequest
	err = json.NewDecoder(request.Body).Decode(&feesAndAmountRequest)
	if err != nil {
		LOGGER.Warningf("Unable to parse body of REST call to /fees endpoint:  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "FEES-1002", err)
		return
	}
	err = feesAndAmountRequest.Validate(strfmt.Default)
	if err != nil {
		msg := "Unable to validate /fees request: " + err.Error()
		LOGGER.Debugf(msg)
		response.NotifyWWError(w, request, http.StatusBadRequest, "FEES-1002", err)
		return
	}
	if feesAndAmountRequest.AmountPayout == 0 && feesAndAmountRequest.AmountGross == 0 {
		msg := "either amount_payout or amount_gross should be filled in the request "
		LOGGER.Debugf(msg)
		response.NotifyWWError(w, request, http.StatusBadRequest, "FEES-1002", errors.New(msg))
		return
	}
	if feesAndAmountRequest.AmountPayout != 0 && feesAndAmountRequest.AmountGross != 0 {
		msg := "either amount_payout or amount_gross should be filled in the request "
		LOGGER.Debugf(msg)
		response.NotifyWWError(w, request, http.StatusBadRequest, "FEES-1002", errors.New(msg))
		return
	}
	pp := feesAndAmountRequest.DetailsPayoutLocation
	if pp == nil {
		LOGGER.Infof("There is no Payout Point in the request...")
	}

	LOGGER.Infof("Fee ID for this request is: %v", feesAndAmountRequest.RequestID)
	settlementAmount := feesAndAmountRequest.AmountGross
	payoutAmount := feesAndAmountRequest.AmountPayout
	if settlementAmount <= 0 && payoutAmount <= 0 {
		LOGGER.Warningf("Either gross amount or payout amount is required")
		response.NotifyWWError(w, request, http.StatusBadRequest, "FEES-1003", nil)
		return
	}

	if settlementAmount > 0 && payoutAmount > 0 {
		LOGGER.Warningf("Either gross amount or payout amount is required but not both")
		response.NotifyWWError(w, request, http.StatusBadRequest, "FEES-1004", nil)
		return
	}

	/*
		sending message to Kafka
	*/

	fee_req, _ := json.Marshal(feesAndAmountRequest)
	err = fo.GatewayOperation.Produce(rfiId+"_"+kafka.FEE_TOPIC, fee_req)
	if err != nil {
		newError := errors.New("Kafka producer failed")
		response.NotifyWWError(w, request, http.StatusBadRequest, "FEES-1009", newError)
		return
	}
	response.Respond(w, http.StatusOK, []byte(`{"status":"Success"}`))
}

func (fo FeeOperations) RespondFees(w http.ResponseWriter, request *http.Request) {

	LOGGER.Infof("*******Receiving fee response from RFI******")
	vars := mux.Vars(request)
	ofiId := strings.TrimSpace(vars["participant_id"])

	if ofiId == "" {
		err := errors.New("Participant ID is empty")
		LOGGER.Errorf(err.Error())
		response.NotifyWWError(w, request, http.StatusBadRequest, "FEES-1001", err)
		return
	}
	participant, err := fo.getParticipantForDomain(ofiId)
	if err != nil {
		LOGGER.Errorf("Error while getting Participant for Domain:", err.Error())
		response.NotifyWWError(w, request, http.StatusBadRequest, "FEES-1001", err)
		return
	}

	// Check if participant is active
	LOGGER.Info("Check participant active")
	err = participant_checks.CheckStatusActive(participant)
	if err != nil {
		msg := err.Error()
		LOGGER.Error(msg)
		response.NotifyFailure(w, request, http.StatusBadRequest, msg)
		return
	}

	var feesAndAmountresponse model.TransactionFees

	err = json.NewDecoder(request.Body).Decode(&feesAndAmountresponse)
	if err != nil {
		LOGGER.Warningf("Unable to parse body of REST call to /fees endpoint:  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "FEES-1002", err)
		return
	}

	err = feesAndAmountresponse.Validate(strfmt.Default)
	if err != nil {
		LOGGER.Warningf("Invalid TransactionFees received from RFI: %v", err.Error())
		response.NotifyWWError(w, request, http.StatusConflict, "FEES-1005", err)
		return
	}

	if *feesAndAmountresponse.DetailsAssetSettlement.AssetCode == *feesAndAmountresponse.AssetCodePayout {
		//enforce same price rate for same payout asset
		if *feesAndAmountresponse.AmountSettlement != *feesAndAmountresponse.AmountPayout {
			msg := "Warning: Settlement asset is same as payout asset, price rate should be the same, only fee should be charged, invalidating RFI's fee response"
			LOGGER.Warningf("%s", msg)
			response.NotifyWWError(w, request, http.StatusBadRequest, "FEES-1006", errors.New(msg))
			return
		}
	}
	farBytes, err := json.Marshal(feesAndAmountresponse)
	if err != nil {
		LOGGER.Warningf("Warning: Cannot marshall FeeAndAmountResponse from RFI: %v", err.Error())
		response.NotifyWWError(w, request, http.StatusBadRequest, "FEES-1006", err)
		return
	}

	/*
		sending message to Kafka
	*/

	err = fo.GatewayOperation.Produce(ofiId+"_"+kafka.FEE_TOPIC, farBytes)
	if err != nil {
		newError := errors.New("Kafka producer failed")
		response.NotifyWWError(w, request, http.StatusBadRequest, "FEES-1009", newError)
		return
	}

	response.Respond(w, http.StatusOK, []byte(`{"status":"Success"}`))
}

// getParticipantForDomain : Get participant for domain
func (fo FeeOperations) getParticipantForDomain(domain string) (model.Participant, error) {
	LOGGER.Debugf("getParticipantForDomain domain = %v", domain)
	var participant model.Participant

	participant, err := fo.prClient.GetParticipantForDomain(domain)
	if err != nil {
		LOGGER.Errorf(" Error getParticipantForDomain failed: %v", err)
		return participant, err
	}
	return participant, err
}
