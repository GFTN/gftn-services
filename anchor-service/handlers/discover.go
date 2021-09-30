// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package handlers

import (
	"net/http"
	"github.com/GFTN/gftn-services/utility/response"
	au "github.com/GFTN/gftn-services/anchor-service/anchor-util"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"encoding/json"
)

type DiscoverHandler struct {
}

func CreateDiscoverHandler() (DiscoverHandler, error) {
	return DiscoverHandler{}, nil
}

func (dh DiscoverHandler) DiscoverParticipant(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	LOGGER.Infof("Handling federation protocol request")
	queryParams := req.URL.Query()

	queryName := queryParams["name"]

	// Split the query name
	if queryName == nil || len(queryName) != 1 {

		LOGGER.Warningf("Discover Participant request had a missing param")
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0001", nil)
		return

	}

	// split the queryName
	accountIdentifier, participantDomain, err := au.ParseFederationName(queryName[0])
	if err != nil {
		LOGGER.Warningf("Discover Participant request had a missing param")
		response.NotifyWWError(w, req, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if accountIdentifier == "" {
		LOGGER.Warningf("Account Identifier cannot be empty")
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0004", nil)
		return
	}

	if participantDomain == "" {
		LOGGER.Warningf("Participant Domain cannot be empty")
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0008", nil)
		return
	}

	LOGGER.Infof("Asking Participant Registry about (%v) hosted at domain (%v)", accountIdentifier, participantDomain)
	participant, err := au.GetParticipantForDomain(participantDomain)
	if err != nil {
		LOGGER.Warningf(err.Error())
		response.NotifyWWError(w, req, http.StatusNotFound, "ANCHOR-0009", nil)
		return
	}
	stellarAddress := au.GetAccountAddressForParticipant(participant, accountIdentifier)
	if stellarAddress == "" {
		LOGGER.Warningf("Invalid Account: %v", accountIdentifier)
		response.NotifyWWError(w, req, http.StatusNotFound, "ANCHOR-0010", nil)
		return
	}

	discoverParticipantResponse := model.DiscoverParticipantResponse{Address: &stellarAddress, AccountName: &accountIdentifier}
	discoverParticipantResponseBytes, _ := json.Marshal(discoverParticipantResponse)
	response.Respond(w, http.StatusOK, discoverParticipantResponseBytes)

}
