// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package crypto_handler

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/go-openapi/strfmt"
	b "github.com/stellar/go/build"
	"github.com/stellar/go/xdr"
	"github.com/GFTN/gftn-services/gftn-models/model"
	util "github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/participant"
	"github.com/GFTN/gftn-services/utility/response"
)

// Returns back address of IBM account
func (op *CryptoOperations) GetIBMAccount(w http.ResponseWriter, req *http.Request) {

	account, err := participant.GenericGetIBMTokenAccount(op.VaultSession)
	if err != nil {
		LOGGER.Debugf("IBM account: %v", account.NodeAddress)
		LOGGER.Errorf("Error: %v", err.Error())
		response.NotifyWWError(w, req, http.StatusNotFound, "CYRPTO-0005", err)
		return
	}

	accountResp := model.Account{}
	accountResp.Name = "IBM"
	accountResp.Address = &account.NodeAddress
	responseData, marshalErr := json.Marshal(accountResp)
	if marshalErr != nil {
		LOGGER.Errorf("Error: %v", marshalErr.Error())
		response.NotifyWWError(w, req, http.StatusNotFound, "CYRPTO-0005", err)
		return
	}
	response.Respond(w, http.StatusOK, responseData)
	return
}

// Takes XDR and signs and sends back XDR signed with IBM account
func (op *CryptoOperations) AddIBMSign(w http.ResponseWriter, req *http.Request) {

	var signRequest model.InternalAdminDraft
	err := json.NewDecoder(req.Body).Decode(&signRequest)
	if err != nil {
		LOGGER.Errorf("Error: %v", err.Error())
		response.NotifyWWError(w, req, http.StatusBadRequest, "CRYPTO-0001", err)
		return
	}
	err = signRequest.Validate(strfmt.Default)
	if err != nil {
		msg := "Unable to validate sign xdr request: " + err.Error()
		LOGGER.Errorf(msg)
		response.NotifyWWError(w, req, http.StatusBadRequest, "CRYPTO-0001", err)
		return
	}

	sEnc := base64.StdEncoding.EncodeToString([]byte(*signRequest.TransactionUnsigned))

	LOGGER.Debugf("Unsigned Transaction: %v", sEnc)
	raw := strings.NewReader(sEnc)
	b64r := base64.NewDecoder(base64.StdEncoding, raw)

	var tx xdr.TransactionEnvelope
	bytesRead, err := xdr.Unmarshal(b64r, &tx)
	LOGGER.Infof("read %d bytes from Xdr \n", bytesRead)
	LOGGER.Infof("This tx has %d operations, source Account %s\n", len(tx.Tx.Operations), tx.Tx.SourceAccount)
	LOGGER.Infof("This tx has paid %d fee\n", tx.Tx.Fee)
	LOGGER.Infof("This tx has sequence number %d \n", tx.Tx.SeqNum)
	stellarPassphrase := os.Getenv(global_environment.ENV_KEY_STELLAR_NETWORK)

	LOGGER.Infof("Stellar Network: %v", stellarPassphrase)

	txeb := &b.TransactionEnvelopeBuilder{E: &tx}
	txeb.Init()
	stellarNetwork := util.GetStellarNetwork(stellarPassphrase)
	err = txeb.MutateTX(stellarNetwork)

	if err != nil {
		LOGGER.Errorf("Error: %v", err.Error())
		response.NotifyWWError(w, req, http.StatusFailedDependency, "CRYPTO-0007", err)
		return
	}

	account, err := participant.GenericGetIBMTokenAccount(op.VaultSession)
	if err != nil {
		LOGGER.Debugf("IBM account: %v", account.NodeAddress)
		LOGGER.Errorf("Error: %v", err.Error())
		response.NotifyWWError(w, req, http.StatusNotFound, "CYRPTO-0005", err)
		return
	}

	// This should be signed as stellar SDK, as we have already getting the account information for IBM account where we have accessed
	// it from nodeconfig or vault depending on env variable
	// do not use generic sign here as we will not be storing key in HSM for IBM account
	sig := b.Sign{Seed: account.NodeSeed}
	err = sig.MutateTransactionEnvelope(txeb)

	signedXdr := model.SignedTransaction{}
	bytesData, err := txeb.Bytes()
	stBytes := strfmt.Base64(bytesData[:])
	LOGGER.Debugf("bytes xdr: %v", len(bytesData))
	signedXdr.TransactionSigned = &stBytes
	signedXdr.TransactionID = ""

	responseData, marshalErr := json.Marshal(signedXdr)
	if marshalErr != nil {
		LOGGER.Errorf("Error: %v", marshalErr.Error())
		response.NotifyWWError(w, req, http.StatusNotFound, "CRYPTO-0006", err)
		return
	}

	response.Respond(w, http.StatusOK, responseData)
	return
}
