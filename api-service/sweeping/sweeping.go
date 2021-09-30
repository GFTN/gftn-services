// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package sweeping

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/stellar/go/clients/horizon"

	"github.com/GFTN/gftn-services/gftn-models/model"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/GFTN/gftn-services/api-service/environment"
	crypto_client "github.com/GFTN/gftn-services/crypto-service-client/crypto-client"
	gasserviceclient "github.com/GFTN/gftn-services/gas-service-client"
	pr_client "github.com/GFTN/gftn-services/participant-registry-client/pr-client"
	"github.com/GFTN/gftn-services/utility"
	"github.com/GFTN/gftn-services/utility/blockchain-adaptor/ww_stellar"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/participant"
	"github.com/GFTN/gftn-services/utility/response"
)

type Operations struct {
	ParticipantRegistryClient pr_client.PRServiceClient
	CryptoServiceClient       crypto_client.CryptoServiceClient
	GasServiceClient          gasserviceclient.GasServiceClient
	HorizonClient             *horizon.Client
}

func CreateSweepingOperations() (Operations, error) {

	op := Operations{}
	var prClient pr_client.PRServiceClient
	if os.Getenv(environment.ENV_KEY_PARTICIPANT_REGISTRY_SERVICE_CLIENT) == "mock" {
		LOGGER.Warningf("USING MOCK PARTICIPANT REGISTRY SERVICE CLIENT")
		prClient = pr_client.MockPRServiceClient{}
	} else {
		LOGGER.Infof("Using REST Participant Registry Service Client")
		cl, err := pr_client.CreateRestPRServiceClient(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL))
		utility.ExitOnErr(LOGGER, err, "Unable to create REST Participant Registry Service Client")
		prClient = cl
	}
	op.ParticipantRegistryClient = prClient

	LOGGER.Info("creating cryptoservice client")
	var cClient crypto_client.CryptoServiceClient
	if os.Getenv(environment.ENV_KEY_PARTICIPANT_REGISTRY_SERVICE_CLIENT) == "mock" {
		cClient, _ = crypto_client.CreateMockCryptoServiceClient()
	} else {
		err := errors.New("")
		cServiceInternalUrl, err := participant.GetServiceUrl(os.Getenv(global_environment.ENV_KEY_CRYPTO_SVC_INTERNAL_URL), os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME))
		if err != nil {
			return op, err
		}
		cClient, err = crypto_client.CreateRestCryptoServiceClient(cServiceInternalUrl)
		if err != nil {
			return op, err
		}
	}
	op.CryptoServiceClient = cClient

	gasServiceClient := gasserviceclient.Client{
		HTTP: &http.Client{Timeout: time.Second * 20},
		URL:  os.Getenv(global_environment.ENV_KEY_GAS_SVC_URL),
	}
	op.GasServiceClient = &gasServiceClient

	op.HorizonClient = &horizon.Client{
		URL:  os.Getenv(global_environment.ENV_KEY_HORIZON_CLIENT_URL),
		HTTP: &http.Client{Timeout: time.Second * 10},
	}

	return op, nil

}

func (op Operations) Sweep(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	targetAccount := vars["account_name"]
	if targetAccount == "" {
		response.NotifyFailure(w, req, http.StatusBadRequest, "account_name missing")
		return
	}
	participantID := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	// get account details
	sweepInstruction := model.SweepInstruction{}

	err := json.NewDecoder(req.Body).Decode(&sweepInstruction)
	if err != nil {
		msg := errors.Errorf("Error decoding request body")
		LOGGER.Error(msg)
		LOGGER.Error(err)
		response.NotifyWWError(w, req, http.StatusBadRequest, "API-1042", msg)
		return
	}

	err = sweepInstruction.Validate(strfmt.Default)
	if err != nil {
		msg := "Unable to validate sweepInstructions: " + err.Error()
		LOGGER.Debugf(msg)
		response.NotifyWWError(w, req, http.StatusBadRequest, "API-1041", err)
		return
	}

	if len(sweepInstruction) > 19 {
		msg := "Fobidden: SourceAccounts more than 19"
		LOGGER.Info(msg)
		response.NotifyWWError(w, req, http.StatusBadRequest, "API-1300", errors.New(msg))
		return
	}
	tc := ww_stellar.StellarTransactionConstructor{}
	//Sanity Check
	var payments []ww_stellar.Payment
	for _, sweep := range sweepInstruction {
		payment := ww_stellar.Payment{
			ww_stellar.Account{participantID, *sweep.AccountName},
			ww_stellar.Account{participantID, targetAccount},
			*sweep.Asset,
			*sweep.Amount,
		}
		payments = append(payments, payment)
	}
	err = ww_stellar.SanityCheck(payments, op.ParticipantRegistryClient, op.HorizonClient)
	if err != nil {
		LOGGER.Error(err)
		msg := "Sanity Check: " + err.Error()
		if wwErr, ok := err.(ww_stellar.Error); ok {
			// return status code 400 if err related to account does not exist || does not trust asset|| insufficient balance || trust limit exceeded
			switch wwErr.Code() {
			case ww_stellar.ERROR_ACCOUNT_DOES_NOT_TRUST_ASSET:
				response.NotifyWWError(w, req, http.StatusBadRequest, "API-1301", wwErr)
			case ww_stellar.ERROR_ACCOUNT_TRUST_LIMIT_EXCEEDED:
				response.NotifyWWError(w, req, http.StatusBadRequest, "API-1302", wwErr)
			case ww_stellar.ERROR_INSUFFICIENT_BALANCE:
				response.NotifyWWError(w, req, http.StatusBadRequest, "API-1303", wwErr)
			case ww_stellar.ERROR_MAP_STELLAR_ADDRESS:
				response.NotifyWWError(w, req, http.StatusBadRequest, "API-1304", wwErr)
			case ww_stellar.ERROR_AMOUNT_LESSTHEN_EQUAL_ZERO:
				response.NotifyWWError(w, req, http.StatusBadRequest, "API-1308", wwErr)
			default:
				response.NotifyFailure(w, req, http.StatusInternalServerError, msg)
			}
			return
		}
		response.NotifyFailure(w, req, http.StatusInternalServerError, msg)
		return
	}
	for idx, _ := range sweepInstruction {
		//check if DO is involved, if yes, reject.
		if *sweepInstruction[idx].Asset.AssetType == "DO" {
			msg := "DO sweeping is not allowed"
			LOGGER.Debug(msg)
			response.NotifyWWError(w, req, http.StatusForbidden, "API-1305", errors.New(msg))
			return
		}
	}
	// construct transcation envelop
	err = tc.InitTransaction(op.GasServiceClient)
	if err != nil {
		msg := "Error init transaction"
		LOGGER.Error(msg)
		response.NotifyFailure(w, req, http.StatusInternalServerError, msg)
		return
	}
	type PaymentBuilderResult struct {
		Error error
	}
	ch := make(chan PaymentBuilderResult)
	// get payment concurrently
	for idx := range sweepInstruction {
		payment := ww_stellar.Payment{
			ww_stellar.Account{participantID, *sweepInstruction[idx].AccountName},
			ww_stellar.Account{participantID, targetAccount},
			*sweepInstruction[idx].Asset,
			*sweepInstruction[idx].Amount}
		go func(ch chan PaymentBuilderResult, idx int) {
			err := tc.AddPayment(payment, op.ParticipantRegistryClient)
			ch <- PaymentBuilderResult{err}
			return
		}(ch, idx)
	}

	for _ = range sweepInstruction {
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
						response.NotifyWWError(w, req, http.StatusBadRequest, "API-1304", wwErr)
					default:
						response.NotifyFailure(w, req, http.StatusInternalServerError, msg)
					}
					return
				}
				response.NotifyFailure(w, req, http.StatusInternalServerError, msg)
				return

			} else {
				LOGGER.Info("Transaction construction: Payment Added")
			}
		case <-time.After(time.Second * 10):
			msg := "Error: Timeout when constructing transaction"
			LOGGER.Info(msg)
			response.NotifyFailure(w, req, http.StatusInternalServerError, msg)
			LOGGER.Error(err)
			break
			// call timed out
		}
	}
	close(ch)
	//set fee
	tc.SetFee()
	LOGGER.Info("result transaction of ", participantID, ":", tc.Base64())
	// sign transcation and concat signatures concurrently
	type SignTransactionResult struct {
		Error error
	}
	ch2 := make(chan SignTransactionResult)
	for idx, _ := range sweepInstruction {
		go func(ch chan SignTransactionResult, idx int) {
			err = tc.SignTransactionAndAppend(*sweepInstruction[idx].AccountName, op.CryptoServiceClient)
			ch2 <- SignTransactionResult{err}
			return
		}(ch2, idx)
	}
	for _ = range sweepInstruction {
		select {
		case SignTransactionResult := <-ch2:
			if SignTransactionResult.Error != nil {
				err := SignTransactionResult.Error
				LOGGER.Error(err)
				msg := "Error: Error when Signing transaction"
				LOGGER.Error(msg)
				response.NotifyFailure(w, req, http.StatusInternalServerError, msg)
				return
			} else {
				LOGGER.Info("Transaction construction: Signature Added")
			}
		case <-time.After(time.Second * 10):
			msg := "Error: Timeout when signing transaction"
			LOGGER.Error(msg)
			response.NotifyFailure(w, req, http.StatusInternalServerError, msg)
			return
			// call timed out
		}
	}
	close(ch2)
	LOGGER.Info("Signed result transaction of ", participantID, ":", tc.Base64())

	// submit to gas service
	hash, _, err := tc.SendTransaction(op.GasServiceClient)
	if err != nil {
		msg := "Error: Sending transaction to Gas Service Failed"
		LOGGER.Error(msg)
		LOGGER.Error(err)
		response.NotifyWWError(w, req, http.StatusInternalServerError, "API-1306", err)
		return
	}
	// fetch updated balance and return as receipt
	var resultBalances []*model.Sweep
	type BalanceResult struct {
		Error          error
		AccountBalance model.Sweep
	}
	ch3 := make(chan BalanceResult)
	account := ww_stellar.Account{
		participantID,
		targetAccount,
	}
	// step 0: filter the unique asset
	type Asset struct {
		AssetCode string
		AssetType string
		IssuerID  string
	}
	assethm := make(map[Asset]model.Asset)
	var assetfiltered []model.Asset
	for _, sc := range sweepInstruction {
		temp := Asset{}
		temp.AssetCode = *sc.Asset.AssetCode
		temp.AssetType = *sc.Asset.AssetType
		temp.IssuerID = sc.Asset.IssuerID
		assethm[temp] = *sc.Asset
	}
	for _, val := range assethm {
		assetfiltered = append(assetfiltered, val)
	}
	LOGGER.Debug("asset filtered for balance receipt: ", assetfiltered)
	// step 1: fetch filtered source asset balance of target account
	for idx, _ := range assetfiltered {
		go func(ch chan BalanceResult, idx int) {

			balance, err := ww_stellar.GetBalance(
				account,
				assetfiltered[idx],
				op.ParticipantRegistryClient,
				op.HorizonClient,
			)
			if err != nil {
				ch <- BalanceResult{err, model.Sweep{}}
				return
			}
			accountBalance := model.Sweep{
				&account.Account,
				&balance.Amount,
				&assetfiltered[idx],
			}
			ch <- BalanceResult{nil, accountBalance}
			return
		}(ch3, idx)
	}
	// Step2 : append results
	for _ = range assetfiltered {
		select {
		case balanceResult := <-ch3:
			if balanceResult.Error != nil {
				err := balanceResult.Error
				LOGGER.Error(err)
				msg := "Error: Error when fetching account balance after successful transaction"
				LOGGER.Error(msg)
				response.NotifyWWError(w, req, http.StatusInternalServerError, "API-1307", err)
				return
			} else {
				resultBalances = append(resultBalances, &balanceResult.AccountBalance)
				LOGGER.Info("Transaction Submit: Added updated account balance for receipt")
			}
		case <-time.After(time.Second * 10):
			msg := "Error: Timeout when fetching updated balance"
			LOGGER.Error(msg)
			response.NotifyFailure(w, req, http.StatusInternalServerError, msg)
			return
			// call timed out
		}
	}
	//return results
	resBody := model.SweepReceipt{resultBalances, time.Now().Unix(), &hash}
	resBodyByte, _ := json.Marshal(resBody)
	response.Respond(w, http.StatusOK, resBodyByte)

}
