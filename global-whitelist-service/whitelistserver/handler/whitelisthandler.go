// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/global-whitelist-service/whitelistserver/database"
	"github.com/GFTN/gftn-services/global-whitelist-service/whitelistserver/utility/prclient"
	"github.com/GFTN/gftn-services/utility/response"
)

type RequestBody struct {
	WlParticipantID string `json:"participant_id"`
}

type RequestBodyMutual struct {
	WlParticipantIDs []string `json:"participant_id"`
}

type WhitelistHandler struct {
	DBClient database.InterfaceClient
	PRClient prclient.InterfaceClient
}

func (wlh WhitelistHandler) GetWLParticipants(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	participantID := vars["participant_id"]
	var whiteListParticipants []model.Participant
	LOGGER.Info("Participant", participantID, "getting wl_participants objects")
	wlparticipantIDs, err := wlh.DBClient.GetWhiteListParicipants(participantID)
	if err != nil {
		LOGGER.Error("Error geting participantIDs")
		response.NotifyWWError(w, request, http.StatusInternalServerError, "WL-1004", err)
		return
	}
	LOGGER.Info("whitelist service fetching participant registry")
	participants, err := wlh.PRClient.GetAllParticipants()
	if err != nil {
		LOGGER.Error("err")
		response.NotifyWWError(w, request, http.StatusInternalServerError, "WL-1005", err)
		return
	}
	for idx, participant := range participants {
		for _, wlparticipantID := range wlparticipantIDs {
			if *participant.ID == wlparticipantID {
				whiteListParticipants = append(whiteListParticipants, participants[idx])
			}
		}
	}
	LOGGER.Info("whitelist service fetching participant sucessfully")
	resBody, _ := json.Marshal(whiteListParticipants)
	response.Respond(w, http.StatusOK, resBody)
	return
}

func (wlh WhitelistHandler) GetWLParticipantIDs(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	participantID := vars["participant_id"]
	wlparticipantIDs, err := wlh.DBClient.GetWhiteListParicipants(participantID)
	LOGGER.Info("Participant", participantID, "getting wl_participants")
	if err != nil {
		LOGGER.Error("Error geting participantIDs")
		response.NotifyWWError(w, request, http.StatusInternalServerError, "WL-1004", err)
		return
	}
	LOGGER.Info("Participant", participantID, "got wl_participant: ", wlparticipantIDs)
	resBody, _ := json.Marshal(wlparticipantIDs)
	response.Respond(w, http.StatusOK, resBody)
	return
}

func (wlh WhitelistHandler) CreateWLParticipant(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	participantID := vars["participant_id"]
	reqBody := RequestBody{}
	err := json.NewDecoder(request.Body).Decode(&reqBody)
	wlparticipantID := reqBody.WlParticipantID

	LOGGER.Info("Participant", participantID, "creating wl_participant: ", wlparticipantID)
	LOGGER.Info("Checking if wl_participant is in Participant Registry")
	participants, err := wlh.PRClient.GetAllParticipants()
	if err != nil {
		LOGGER.Error("Connetion to PR failed")
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusInternalServerError, "WL-1005", nil)
		return
	}
	participantExist := false
	for _, participant := range participants {
		if *participant.ID == wlparticipantID {
			LOGGER.Info("Participant", wlparticipantID, "is in PR")
			participantExist = true
			break
		}
	}
	if participantExist == false {
		err = errors.New("Participant for whitelist does not exist")
		LOGGER.Info(err)
		LOGGER.Info("WL Participant ID:", wlparticipantID)
		response.NotifyWWError(w, request, http.StatusBadRequest, "WL-1006", err)
		return
	}
	LOGGER.Info("Checking Passed")
	LOGGER.Info("Participant", participantID, "creating wl_participant: ", wlparticipantID)
	err = wlh.DBClient.AddWhitelistParticipant(participantID, wlparticipantID)
	if err != nil {
		LOGGER.Error("Error geting participantIDs")
		response.NotifyWWError(w, request, http.StatusInternalServerError, "WL-1006", err)
		return
	}
	LOGGER.Info("Participant", participantID, "created wl_participant: ", wlparticipantID)

	response.Respond(w, http.StatusOK, []byte("Successfully added "+wlparticipantID+" to the whitelist."))
	return
}

func (wlh WhitelistHandler) DeleteWLParticipant(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	participantID := vars["participant_id"]
	reqBody := RequestBody{}
	err := json.NewDecoder(request.Body).Decode(&reqBody)
	wlparticipantID := reqBody.WlParticipantID
	LOGGER.Info("Participant", participantID, "deleting wl_participant: ", wlparticipantID)
	err = wlh.DBClient.DeleteWhitelistParticipant(participantID, wlparticipantID)
	if err != nil {
		LOGGER.Error("Error geting participantIDs")
		response.NotifyWWError(w, request, http.StatusInternalServerError, "WL-1007", err)
		return
	}
	LOGGER.Info("Participant", participantID, "deleted wl_participant: ", wlparticipantID)
	response.Respond(w, http.StatusOK, nil)
	return
}

type Result struct {
	ParticipantID  string
	Wlparticipants []string
	Error          error
}

func (wlh WhitelistHandler) GetMutualWLParticipantIDs(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	participantID := vars["participant_id"]
	// get the participant's whitelist
	mutualwl, err := wlh.getMutualWL(participantID)
	if err != nil {
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusInternalServerError, "WL-1004", err)
	}
	LOGGER.Info("whitelist service fetching participant sucessfully")
	resBody, _ := json.Marshal(mutualwl)
	response.Respond(w, http.StatusOK, resBody)
	return
}

func (wlh WhitelistHandler) GetMutualWLParticipants(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	participantID := vars["participant_id"]
	// get the participant's whitelist
	mutualwl, err := wlh.getMutualWL(participantID)
	if err != nil {
		LOGGER.Error(err)
		response.NotifyWWError(w, request, http.StatusInternalServerError, "WL-1004", err)
	}

	LOGGER.Info("whitelist service fetching participant registry")
	var whiteListParticipants []model.Participant
	participants, err := wlh.PRClient.GetAllParticipants()
	if err != nil {
		LOGGER.Error("err")
		response.NotifyWWError(w, request, http.StatusInternalServerError, "WL-1005", err)
		return
	}
	for idx, participant := range participants {
		for _, wlparticipantID := range mutualwl {
			if *participant.ID == wlparticipantID {
				whiteListParticipants = append(whiteListParticipants, participants[idx])
			}
		}
	}
	LOGGER.Info("whitelist service fetching participant sucessfully")
	resBody, _ := json.Marshal(whiteListParticipants)
	response.Respond(w, http.StatusOK, resBody)
	return
}

// Client Endpoints Handlers
func (wlh WhitelistHandler) GetWLParticipantsClient(w http.ResponseWriter, r *http.Request) {
	participantId, err := middlewares.GetIdentity(r)
	if err != nil {
		LOGGER.Error(err)
		response.NotifyWWError(w, r, http.StatusUnauthorized, "WL-1008", nil)
		return
	}
	mux.Vars(r)["participant_id"] = participantId
	wlh.GetWLParticipants(w, r)
}

func (wlh WhitelistHandler) GetWLParticipantIDsClient(w http.ResponseWriter, r *http.Request) {
	participantId, err := middlewares.GetIdentity(r)
	if err != nil {
		LOGGER.Error(err)
		response.NotifyWWError(w, r, http.StatusUnauthorized, "WL-1008", nil)
		return
	}
	mux.Vars(r)["participant_id"] = participantId
	wlh.GetWLParticipantIDs(w, r)
}

func (wlh WhitelistHandler) DeleteWLParticipantClient(w http.ResponseWriter, r *http.Request) {
	participantId, err := middlewares.GetIdentity(r)
	if err != nil {
		LOGGER.Error(err)
		response.NotifyWWError(w, r, http.StatusUnauthorized, "WL-1008", nil)
		return
	}
	mux.Vars(r)["participant_id"] = participantId
	wlh.DeleteWLParticipant(w, r)
}

func (wlh WhitelistHandler) CreateWLParticipantClient(w http.ResponseWriter, r *http.Request) {
	participantId, err := middlewares.GetIdentity(r)
	if err != nil {
		LOGGER.Error(err)
		response.NotifyWWError(w, r, http.StatusUnauthorized, "WL-1008", nil)
		return
	}
	mux.Vars(r)["participant_id"] = participantId
	wlh.CreateWLParticipant(w, r)
}
