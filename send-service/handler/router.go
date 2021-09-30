// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package handler

import (
	"net/http"
	"os"
	"strings"

	global_environment "github.com/GFTN/gftn-services/utility/global-environment"

	"github.com/GFTN/gftn-services/utility/payment/environment"
	message_handler "github.com/GFTN/gftn-services/utility/payment/message-handler"

	"github.com/GFTN/gftn-services/utility/response"

	"github.com/GFTN/gftn-services/utility/payment/constant"
)

func Router(w http.ResponseWriter, req *http.Request, op message_handler.PaymentOperations, sourceType string) {

	var err error
	var report []byte
	var data []byte
	var payloadType string

	target := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	BIC := os.Getenv(environment.ENV_KEY_PARTICIPANT_BIC)
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
		report, err = iso20022Router(data, messageType, op, sourceType)
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
