// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package kafka

import (
	"net/http"
	"os"
	"strings"

	global_environment "github.com/GFTN/gftn-services/utility/global-environment"

	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"

	pr_client "github.com/GFTN/gftn-services/participant-registry-client/pr-client"
	message_handler "github.com/GFTN/gftn-services/utility/payment/message-handler"

	"github.com/GFTN/gftn-services/utility/response"

	"github.com/GFTN/gftn-services/utility/payment/constant"
)

func Router(w http.ResponseWriter, req *http.Request, op message_handler.PaymentOperations) {

	var err error
	var report []byte
	var data []byte
	var payloadType string

	target, _ := middlewares.GetIdentity(req)
	prServiceURL := os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL)
	prc, prcErr := pr_client.CreateRestPRServiceClient(prServiceURL)
	if prcErr != nil {
		LOGGER.Error("Can not create connection to PR client service, please check if PR service is running")
		return
	}

	participant, prcGetErr := prc.GetParticipantForDomain(target)
	if prcGetErr != nil {
		LOGGER.Error("Could not found participant from PR service")
		return
	}
	BIC := *participant.Bic

	data, report, payloadType, err = message_handler.ValidateRequest(req, BIC, target)
	if err != nil {
		response.Respond(w, http.StatusBadRequest, report)
		return
	}

	standardType := strings.TrimSpace(strings.ToLower(strings.Split(payloadType, ":")[0]))
	messageType := strings.TrimSpace(strings.ToLower(strings.Split(payloadType, ":")[1]))

	LOGGER.Infof("Receiving standard type: %v", standardType)
	// Route to different messageType router base on the standardType
	switch standardType {
	case constant.ISO20022:
		report, err = iso20022Router(data, BIC, messageType, target, op)
	case constant.ISO8385:
		LOGGER.Warning("ISO8385 not support yet")
		response.Respond(w, http.StatusBadRequest, []byte("ISO8385 not support yet"))
		return
	case constant.MT:
		LOGGER.Warning("MT not support yet")
		response.Respond(w, http.StatusBadRequest, []byte("MT not support yet"))
		return
	case constant.JSON:
		LOGGER.Warning("JSON not support yet")
		response.Respond(w, http.StatusBadRequest, []byte("JSON not support yet"))
		return

		/*
			------------ New message standard type ------------
		*/

	default:
		LOGGER.Error("Undefined message standard type")
		response.Respond(w, http.StatusBadRequest, []byte("Undefined message standard type"))
		return
	}

	if err != nil {
		response.Respond(w, http.StatusBadRequest, report)
	} else {
		response.Respond(w, http.StatusOK, report)
	}

	return

}
