// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package crypto_handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/GFTN/gftn-services/utility/participant"
	"github.com/GFTN/gftn-services/utility/response"
)

//add next
func (op *CryptoOperations) CreateAccount(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	accountName := vars["account_name"]

	if accountName == "" {
		err := errors.New("Missing required parameter: account_name")
		LOGGER.Errorf("Error: %v", err.Error())
		response.NotifyWWError(w, req, http.StatusBadRequest, "CRYPTO-0004", err)
		return
	}

	account, err := participant.GenericGetAccount(op.VaultSession, accountName)
	if err != nil {
		LOGGER.Errorf("account: %v, %v, %v", account.NodeAddress, account.PrivateKeyLabel, account.PublicKeyLabel)
	}
	if account.NodeAddress != "" {
		//account already exists
		LOGGER.Errorf("Account already exists: %s", account.NodeAddress)
		response.NotifyWWError(w, req, http.StatusAlreadyReported, "CYRPTO-0004", errors.New("account already exists: "+account.NodeAddress))
		return
	}

	accountHSM, err := op.HSMInstance.GenericGenerateAccount()
	if err != nil {
		LOGGER.Errorf("Error: %v", err.Error())
		response.NotifyWWError(w, req, http.StatusFailedDependency, "CYRPTO-0004", err)
		return
	}

	LOGGER.Debugf("Account Generated: %+v", accountHSM)

	responseData, marshalErr := json.Marshal(accountHSM)
	if marshalErr != nil {
		response.NotifyWWError(w, req, http.StatusNotFound, "CYRPTO-0004", err)
		return
	}
	response.Respond(w, http.StatusCreated, responseData)
	return
}
