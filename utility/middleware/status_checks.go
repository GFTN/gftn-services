// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package middleware

import (
	"net/http"
	"os"

	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"
	pr_client "github.com/GFTN/gftn-services/participant-registry-client/pr-client"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/response"
)

type MiddlewareHandler struct {
	prc *pr_client.RestPRServiceClient
}

func CreateMiddlewareHandler() *MiddlewareHandler {

	prServiceURL := os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL)
	LOGGER.Info("Using REST Participant Registry Service Client")
	prc, prcErr := pr_client.CreateRestPRServiceClient(prServiceURL)

	if prcErr != nil {
		LOGGER.Error(prcErr)
		LOGGER.Error("Can not create connection to PR client service")
		return nil
	}

	return &MiddlewareHandler{&prc}
}

func (mwh *MiddlewareHandler) ParticipantStatusCheck(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// // Ignore all non POST and DELETE requests
	// if (r.Method != http.MethodPost) && (r.Method != http.MethodDelete) {
	// 	next.ServeHTTP(w, r)
	// 	return
	// }

	// Retrieve participant id
	homeDomain := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)

	// Get participant id from request if global service
	if homeDomain == "ww" {
		participantID, err := middlewares.GetIdentity(r)
		if err != nil {
			msg := "Cannot retrieve participant ID"
			LOGGER.Error(err)
			LOGGER.Error(msg)
			response.NotifyFailure(w, r, http.StatusInternalServerError, msg)
			return
		}
		homeDomain = participantID
	}

	participant, prcGetErr := mwh.prc.GetParticipantForDomain(homeDomain)

	if prcGetErr != nil {
		msg := "Participant not found from PR service"
		LOGGER.Error(prcGetErr)
		LOGGER.Error(msg)
		response.NotifyFailure(w, r, http.StatusInternalServerError, msg)
		return
	}

	if participant.Status != "active" {
		msg := "Participant status not active"
		response.NotifyFailure(w, r, http.StatusBadRequest, msg)
		return
	}

	next.ServeHTTP(w, r)
}
