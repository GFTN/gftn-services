// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package handlers

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	b "github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
	"github.com/GFTN/gftn-services/anchor-service/environment"
	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"
	gasserviceclient "github.com/GFTN/gftn-services/gas-service-client"
	"github.com/GFTN/gftn-services/gftn-models/model"
	pr_client "github.com/GFTN/gftn-services/participant-registry-client/pr-client"
	util "github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	participant_checks "github.com/GFTN/gftn-services/utility/participant"
	"github.com/GFTN/gftn-services/utility/response"
)

type FundHandler struct {
	prClient         pr_client.PRServiceClient
	GasServiceClient gasserviceclient.GasServiceClient
}

type GasAccountAndSequence struct {
	Pkey           string `json:"pkey,omitempty"`
	SequenceNumber int    `json:"SequenceNumber,omitempty"`
}

type SignedXDR struct {
	OneSignedXDR string `json:"oneSignedXDR,omitempty"`
}

type Success struct {
	Hash   string `json:"hash"`
	Ledger string `json:"ledger"`
}

type Failure struct {
	Title         string `json:"title"`
	FailureReason string `json:"failure_reason"`
}
type Failure403 struct {
	Title         string           `json:"title"`
	FailureReason FailureReason403 `json:"failure_reason"`
}

type FailureReason403 struct {
	EnvelopeXdr string      `json:"envelope_xdr"`
	ResultCodes ResultCodes `json:"result_codes"`
	ResultXdr   string      `json:"result_xdr"`
}

type ResultCodes struct {
	Transaction string `json:"transaction"`
}

func CreateFundHandler() (FundHandler, error) {
	fh := FundHandler{}
	var prClient pr_client.PRServiceClient
	var err error
	if os.Getenv(environment.ENV_KEY_PARTICIPANT_REGISTRY_SERVICE_CLIENT) == "mock" {
		prClient, err = pr_client.CreateMockPRServiceClient(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL))
	} else {
		prClient, err = pr_client.CreateRestPRServiceClient(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL))
	}

	if err != nil {
		LOGGER.Errorf(" Error CreateRestPRServiceClient failed  %v", err)
		return fh, err
	}
	fh.prClient = prClient
	gasServiceClient := gasserviceclient.Client{
		HTTP: &http.Client{Timeout: time.Second * 20},
		URL:  os.Getenv(global_environment.ENV_KEY_GAS_SVC_URL),
	}
	fh.GasServiceClient = &gasServiceClient
	return fh, nil
}

// this endpoint is to construct a funding transaction (xdr) and return it to the anchor
func (fh FundHandler) FundRequest(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	fundRequest := model.Funding{}
	err := json.NewDecoder(req.Body).Decode(&fundRequest)
	if err != nil {
		msg := "Unable to parse body of REST call to Fund Stable Coin Request: " + err.Error()
		LOGGER.Debugf(msg)
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0070]", err)
		return
	}

	err = fundRequest.Validate(strfmt.Default)

	if err != nil {
		msg := "Unable to validate fund request: " + err.Error()
		LOGGER.Debugf(msg)
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0070", err)
		return
	}

	//Check JWT token
	if os.Getenv(global_environment.ENV_KEY_ENABLE_JWT) != "false" {
		participantID, err := middlewares.GetIdentity(req)
		//Check if requesting anchor id is same as participant id in the token
		if participantID != *fundRequest.AnchorID {
			response.NotifyWWError(w, req, http.StatusUnauthorized, "ANCHOR-0067",
				err)
			return
		}
	}
	//get participant issuing account

	LOGGER.Infof("Participant Domain: %v:", *fundRequest.ParticipantID)
	var participantAddress string
	accountName := fundRequest.AccountName
	if accountName == "" {
		accountName = util.ISSUING
	}
	participantAddress, err = fh.prClient.GetParticipantAccount(*fundRequest.ParticipantID, accountName)
	if err != nil {
		msg := "Unable to get participant account from participant registry: " + err.Error()
		LOGGER.Warningf(msg)
		response.NotifyWWError(w, req, http.StatusNotFound, "ANCHOR-0076", err)
		return
	}

	LOGGER.Infof("Participant Address: %v", participantAddress)

	// get anchor account
	LOGGER.Infof("Anchor Domain: %v:", *fundRequest.AnchorID)
	anchor, err := fh.prClient.GetParticipantForDomain(*fundRequest.AnchorID)
	if err != nil {
		msg := "Unable to get anchor from participant registry: " + err.Error()
		LOGGER.Warningf(msg)
		response.NotifyWWError(w, req, http.StatusNotFound, "ANCHOR-0020", err)
		return
	}
	anchorAddress := anchor.IssuingAccount
	LOGGER.Infof("Anchor Address: %v", anchorAddress)

	participant, err := fh.prClient.GetParticipantForDomain(*fundRequest.ParticipantID)
	if err != nil {
		msg := "Unable to get participant from participant registry: " + err.Error()
		LOGGER.Warningf(msg)
		response.NotifyWWError(w, req, http.StatusNotFound, "ANCHOR-0006", err)
		return
	}

	// Check if participant is active
	LOGGER.Info("Check participant active")
	err = participant_checks.CheckStatusActive(participant)
	if err != nil {
		msg := err.Error()
		LOGGER.Error(msg)
		response.NotifyFailure(w, req, http.StatusBadRequest, msg)
		return
	}

	// VERIFY IF ASSET ISSUED BY ANCHOR
	// verify if the asset is registered in world wire
	// verify if the participant has trust relationship with this asset
	// assuming all of the above are true, create a funding transaction

	rawTx, err := fh.buildFundTransactionXDR(fundRequest, anchorAddress, participantAddress)
	if err != nil {
		msg := "Unable to build funding transaction: " + err.Error()
		LOGGER.Warningf(msg)
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0020", err)
		return
	}

	anchorfundresp := model.FundingInstruction{}
	anchorfundresp.DetailsFunding = &fundRequest
	anchorfundresp.InstructionUnsigned = rawTx
	acrJSON, _ := json.Marshal(anchorfundresp)
	response.Respond(w, http.StatusOK, acrJSON)

}

// this endpoint is for submitting a funding transaction that has been signed by the anchor
func (fh FundHandler) SignedFundRequest(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	queryParams := req.URL.Query()
	fs := queryParams["funding_signed"]

	if fs == nil || len(fs) != 1 {

		LOGGER.Warningf("Anchor Fund request had a missing param")
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0074", nil)
		return

	}

	signedFundingRequest, err := base64.StdEncoding.DecodeString(fs[0])
	if err != nil {
		msg := "funding_signed is not valid base64 encode " + err.Error()
		LOGGER.Debugf(msg)
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0070", errors.New("funding_signed is not valid base64 encode"))
		return
	}

	is := queryParams["instruction_signed"]
	if is == nil || len(is) != 1 {
		LOGGER.Warningf("Anchor Fund request had a missing param")
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0075", nil)
		return
	}
	instructionSignedXDR := is[0]

	fundRequest := model.Funding{}
	err = json.NewDecoder(req.Body).Decode(&fundRequest)
	if err != nil {
		msg := "Unable to parse body of REST call to Signed Fund Stable Coin Request: " + err.Error()
		LOGGER.Debugf(msg)
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0070", err)
		return
	}

	err = fundRequest.Validate(strfmt.Default)

	if err != nil {
		msg := "Unable to validate signed fund request: " + err.Error()
		LOGGER.Debugf(msg)
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0070", err)
		return
	}

	//Check JWT token
	if os.Getenv(global_environment.ENV_KEY_ENABLE_JWT) != "false" {
		participantID, err := middlewares.GetIdentity(req)
		//Check if requesting anchor id is same as participant id in the token
		if participantID != *fundRequest.AnchorID {
			response.NotifyWWError(w, req, http.StatusUnauthorized, "ANCHOR-0067",
				err)
			return
		}
	}

	// get anchor signing key
	LOGGER.Infof("Anchor Domain: %v:", *fundRequest.AnchorID)
	// get participant config got domain
	signingKey, err := fh.prClient.GetParticipantIssuingAccount(*fundRequest.AnchorID)
	if err != nil {
		msg := "Unable to get anchor signing key from participant registry: " + err.Error()
		LOGGER.Warningf(msg)
		response.NotifyWWError(w, req, http.StatusNotFound, "ANCHOR-0020", err)
		return
	}

	kp, err := keypair.Parse(signingKey)
	LOGGER.Infof("Signing Key: %v", signingKey)
	if err != nil {
		LOGGER.Warningf("Unable to parse stellar signing key from anchor (%v):  %v", signingKey, err)
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0071", err)
		return
	}

	//unsigned request
	afrString, err := json.Marshal(fundRequest)
	base64afrString := base64.StdEncoding.EncodeToString(afrString)
	LOGGER.Infof("AnchorFundRequest: base64 %s", (base64afrString))
	LOGGER.Infof("Signed request: %v", signedFundingRequest)

	LOGGER.Infof("Unsigned request: %v", string(afrString))
	/*signereq := strfmt.Base64{}
	//signereq, err = kp.Sign([]byte(base64afrString))

	signed := base64.StdEncoding.EncodeToString(signereq[:])
	LOGGER.Debug("signed: ", signed)
	decode, err := base64.StdEncoding.DecodeString(signed)
	if err != nil {
		LOGGER.Error(err)
	}
	LOGGER.Debug("Decode success:", decode)
	*/

	err = kp.Verify([]byte(base64afrString), []byte(signedFundingRequest))

	if err != nil {
		LOGGER.Warningf("The Fund Transaction was not signed with the correct key:  %v", err)
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0072", err)
		return
	}

	// **************      SUBMIT TRANSACTION IN Gas Service ****************

	hash, _, err := fh.GasServiceClient.SubmitTxe(instructionSignedXDR)
	if err != nil {
		LOGGER.Warningf("The Fund Transaction submission to GAS failed:  %v", err)
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0073", err)
		return
	}

	status := "cleared"
	timeStamp := time.Now().Unix()
	txnReceipt := model.TransactionReceipt{Transactionid: &hash, Transactionstatus: &status, Timestamp: &timeStamp}

	signedFundResponse := model.FundingReceipt{}
	signedFundResponse.ReceiptFunding = &txnReceipt
	signedFundResponse.DetailsFunding = &fundRequest

	sfresp, _ := json.Marshal(signedFundResponse)
	response.Respond(w, http.StatusOK, sfresp)
}

// this function will construct the funding transaction XDR
func (fh FundHandler) buildFundTransactionXDR(fundRequest model.Funding, anchorAddress string, participantAddress string) (rawTx string, err error) {
	var tx *b.TransactionBuilder
	horizonClient := util.GetHorizonClient(os.Getenv(global_environment.ENV_KEY_HORIZON_CLIENT_URL))
	stellarNetwork := util.GetStellarNetwork(os.Getenv(global_environment.ENV_KEY_STELLAR_NETWORK))
	LOGGER.Infof("Anchor Address: %v", anchorAddress)
	LOGGER.Infof("Participant Address: %v", participantAddress)

	// since IBM will py for transaction gas, get the source and sequence nbr from IBM)
	gaspkey, sequence, err := fh.GasServiceClient.GetAccountAndSequence()
	if err != nil {
		LOGGER.Warningf(" Error getting Gas account details: %s", err)
		return "", err
	}
	LOGGER.Infof("Gas Account Number: %v", gaspkey)
	LOGGER.Infof("Gas Sequence Number: %v", sequence)

	tx, err = b.Transaction(
		b.SourceAccount{AddressOrSeed: gaspkey},
		stellarNetwork,
		b.AutoSequence{SequenceProvider: &horizonClient},
		b.Payment(
			b.SourceAccount{AddressOrSeed: anchorAddress},
			b.Destination{AddressOrSeed: participantAddress},
			b.CreditAmount{Code: *fundRequest.AssetCodeIssued, Issuer: anchorAddress, Amount: util.FloatToString(*fundRequest.AmountFunding)},
		))
	if err != nil {
		LOGGER.Warningf(" Error creating xdr %s", err)
		return "", err
	}
	txnMemoBytes, _ := json.Marshal(fundRequest)
	x := sha512.Sum512_256(txnMemoBytes)
	memo := xdr.Hash(x)
	memoHash, err := xdr.NewMemo(xdr.MemoTypeMemoHash, memo)
	if err != nil {
		LOGGER.Warningf(" Error creating transaction memo: %s", err)
		return "", err
	}
	tx.TX.Memo = memoHash
	var txe b.TransactionEnvelopeBuilder
	err = txe.Mutate(tx)

	if err != nil {
		LOGGER.Warningf(" Error while mutating transaction: %s", err)
		return "", err
	}
	txeB64, err := txe.Base64()

	if err != nil {
		LOGGER.Warningf(" Error getting base64 for transaction: %s", err)
		return "", err
	}

	return txeB64, nil

}
