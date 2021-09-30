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

//SignXdr: verifies identification and then signs xdr if pass
func (op *CryptoOperations) SignXdr(w http.ResponseWriter, req *http.Request) {

	var signRequest model.Draft
	err := json.NewDecoder(req.Body).Decode(&signRequest)
	if err != nil {
		LOGGER.Errorf("Error: %v", err.Error())
		response.NotifyWWError(w, req, http.StatusBadRequest, "CRYPTO-0001", err)
		return
	}
	err = signRequest.Validate(strfmt.Default)
	if err != nil {
		msg := "Unable to validate sign xdr request: " + err.Error()
		LOGGER.Errorf("%v", msg)
		response.NotifyWWError(w, req, http.StatusBadRequest, "CRYPTO-0001", err)
		return
	}

	sEnc := base64.StdEncoding.EncodeToString([]byte(*signRequest.TransactionUnsigned))

	//LOGGER.Debugf("Unsigned Transaction: %v", sEnc)
	raw := strings.NewReader(sEnc)
	b64r := base64.NewDecoder(base64.StdEncoding, raw)

	var tx xdr.TransactionEnvelope
	bytesRead, err := xdr.Unmarshal(b64r, &tx)
	LOGGER.Infof("read %d bytes from Xdr \n", bytesRead)
	//LOGGER.Infof("This tx has %d operations, source Account %s\n", len(tx.Tx.Operations), tx.Tx.SourceAccount)
	LOGGER.Infof("This tx has paid %d fee\n", tx.Tx.Fee)
	LOGGER.Infof("This tx has sequence number %d \n", tx.Tx.SeqNum)
	stellarPassphrase := os.Getenv(global_environment.ENV_KEY_STELLAR_NETWORK)

	txeb := &b.TransactionEnvelopeBuilder{E: &tx}
	txeb.Init()
	stellarNetwork := util.GetStellarNetwork(stellarPassphrase)
	err = txeb.MutateTX(stellarNetwork)

	if err != nil {
		LOGGER.Errorf("Error: %v", err.Error())
		response.NotifyWWError(w, req, http.StatusFailedDependency, "CRYPTO-0007", err)
		return
	}

	account, err := participant.GenericGetAccount(op.VaultSession, *signRequest.AccountName)
	if err != nil {
		LOGGER.Errorf("Error: %v", err.Error())
		response.NotifyWWError(w, req, http.StatusFailedDependency, "CRYPTO-0005", err)
		return
	}
	//Verify signature

	//LOGGER.Debugf("Unsigned ID: %v", base64.StdEncoding.EncodeToString([]byte(*signRequest.IDUnsigned)))
	//LOGGER.Debugf("Signed ID: %v", base64.StdEncoding.EncodeToString([]byte(*signRequest.IDSigned)))
	verification, err := op.HSMInstance.GenericVerifySignatureIdentity(*signRequest.IDUnsigned, *signRequest.IDSigned, account)

	if err != nil {
		LOGGER.Errorf("Error: %v", err.Error())
		response.NotifyWWError(w, req, http.StatusFailedDependency, "CRYPTO-0002", err)
		return
	}
	if verification != true {
		if err != nil {
			LOGGER.Errorf("Error: %v", err.Error())
			response.NotifyWWError(w, req, http.StatusFailedDependency, "CRYPTO-0003", err)
			return
		}
	}

	signedXdr := model.SignedTransaction{}

	//LOGGER.Debugf("Transaction: %v, %v", len(txeb.E.Signatures), txeb.E.Signatures)
	signedXdrEnv, err := op.HSMInstance.GenericSign(txeb, account)

	if err != nil {
		LOGGER.Errorf("Error  %v", err.Error())
		response.NotifyWWError(w, req, http.StatusFailedDependency, "CRYPTO-0006", err)
		return
	}
	//LOGGER.Debugf("signed xdr: %v", signedXdrEnv)

	bytesData, err := signedXdrEnv.Bytes()
	stBytes := strfmt.Base64(bytesData[:])
	LOGGER.Debugf("bytes xdr: %v", len(bytesData))
	signedXdr.TransactionSigned = &stBytes
	signedXdr.TransactionID = signRequest.TransactionID

	responseData, marshalErr := json.Marshal(signedXdr)
	if marshalErr != nil {
		response.NotifyWWError(w, req, http.StatusNotFound, "CRYPTO-0006", err)
		return
	}
	//LOGGER.Debugf("success: %v", responseData)
	response.Respond(w, http.StatusCreated, responseData)
	return
}

//ParticipantSignXdr: This is participant only signing endpoint for signing xdr with out verification
func (op *CryptoOperations) ParticipantSignXdr(w http.ResponseWriter, req *http.Request) {

	var signRequest model.InternalDraft
	err := json.NewDecoder(req.Body).Decode(&signRequest)
	if err != nil {
		LOGGER.Debugf("Error  %v", err.Error())
		response.NotifyWWError(w, req, http.StatusBadRequest, "CRYPTO-0001", err)
		return
	}
	err = signRequest.Validate(strfmt.Default)
	if err != nil {
		msg := "Unable to validate sign xdr request: " + err.Error()
		LOGGER.Debugf(msg)
		response.NotifyWWError(w, req, http.StatusBadRequest, "CRYPTO-0001", err)
		return
	}

	sEnc := base64.StdEncoding.EncodeToString([]byte(*signRequest.TransactionUnsigned))

	raw := strings.NewReader(sEnc)
	b64r := base64.NewDecoder(base64.StdEncoding, raw)

	var tx xdr.TransactionEnvelope
	bytesRead, err := xdr.Unmarshal(b64r, &tx)
	LOGGER.Infof("read %d bytes from Xdr \n", bytesRead)
	//LOGGER.Infof("This tx has %d operations, source Account %s\n", len(tx.Tx.Operations), tx.Tx.SourceAccount)
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

	account, err := participant.GenericGetAccount(op.VaultSession, *signRequest.AccountName)
	if err != nil {
		LOGGER.Errorf("Error: %v", err.Error())
		response.NotifyWWError(w, req, http.StatusFailedDependency, "CRYPTO-0005", err)
		return
	}

	signedXdr := model.SignedTransaction{}

	signedXdrEnv, err := op.HSMInstance.GenericSign(txeb, account)

	if err != nil {
		LOGGER.Errorf("Error: %v", err.Error())
		response.NotifyWWError(w, req, http.StatusFailedDependency, "CRYPTO-0006", err)
		return
	}

	bytesData, err := signedXdrEnv.Bytes()
	stBytes := strfmt.Base64(bytesData[:])
	//LOGGER.Debugf("bytes xdr: %v", len(bytesData))
	signedXdr.TransactionSigned = &stBytes
	signedXdr.TransactionID = signRequest.TransactionID

	responseData, marshalErr := json.Marshal(signedXdr)
	if marshalErr != nil {
		LOGGER.Errorf("Error: %v", marshalErr.Error())
		response.NotifyWWError(w, req, http.StatusNotFound, "CRYPTO-0006", err)
		return
	}
	//LOGGER.Debugf("success: %v", responseData)
	response.Respond(w, http.StatusCreated, responseData)
	return
}
