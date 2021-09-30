// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package response

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/message"
)

var buildVersion = "Local"
var envBuildVersion = "ENV_BUILD_VERSION"

func Respond(w http.ResponseWriter, statusCode int, response []byte) {
	// add for API Security Vulnerability
	common.SetContentType(w)
	common.SetXXSS(w)
	common.SetHSTS(w)
	common.SetXCTO(w)
	common.SetCSP(w)
	common.SetCacheControl(w)
	common.SetPragma(w)

	w.WriteHeader(statusCode)
	w.Write(response)
}

func RespondXML(w http.ResponseWriter, statusCode int, response []byte) {
	// add for API Security Vulnerability
	w.Header().Set("Content-Type", "application/xml")
	common.SetXXSS(w)
	common.SetHSTS(w)
	common.SetXCTO(w)
	common.SetCSP(w)
	common.SetCacheControl(w)
	common.SetPragma(w)

	w.WriteHeader(statusCode)
	w.Write(response)
}

func NotFound(w http.ResponseWriter, req *http.Request) {

	// add for API Security Vulnerability
	common.SetContentType(w)
	common.SetXXSS(w)
	common.SetHSTS(w)
	common.SetXCTO(w)
	common.SetCSP(w)
	common.SetCacheControl(w)
	common.SetPragma(w)

	// create new error
	err := errors.New("URL not found")
	msg := message.Translate(req, "API-NotFound", err)
	msg.Type = "NotFound"
	// update error msg to include 404 url
	/* This line got flagged highly unsafe in checkmarx scan
	so removing req.RequestURI
	msg.Details = "Request URL not found: " + req.RequestURI
	*/
	*msg.Details = "Request URL not found"

	errBytes, _ := json.Marshal(msg)
	w.WriteHeader(http.StatusNotFound)
	w.Write(errBytes)
	return
}

func NotifyWWError(w http.ResponseWriter, r *http.Request, statusCode int, errorCode string, err error) {
	w.Header().Set("Content-Type", "application/json")

	// add for API Security Vulnerability
	common.SetContentType(w)
	common.SetXXSS(w)
	common.SetHSTS(w)
	common.SetXCTO(w)
	common.SetCSP(w)
	common.SetCacheControl(w)
	common.SetPragma(w)

	msg := message.Translate(r, errorCode, err)
	msg.Type = "NotifyWWError"
	service := os.Getenv(global_environment.ENV_KEY_SERVICE_NAME)

	// Default to local version if build version env variable is not set
	// This version connects to env varible pushed into docker container at build time
	build := os.Getenv(envBuildVersion)
	if build == "" {
		build = buildVersion
	}

	msg.BuildVersion = build
	msg.Service = service

	// write to fb logs node as stream of logs information
	// TODO: by chase, need to consider for production
	// if this log should be made available to the user

	errBytes, _ := json.Marshal(msg)

	//Handle the case when statusCode is not set as a safe practice
	if statusCode == 0 {
		statusCode = http.StatusNotFound
	}

	w.WriteHeader(statusCode)
	w.Write(errBytes)
	return
}

// This is a generic common function to generate a standard message response for WW service
func createMessage(msg string, r *http.Request) []byte {
	service := os.Getenv(global_environment.ENV_KEY_SERVICE_NAME)
	timeNow := time.Now().Unix()
	ParticipantID := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	URL := r.RequestURI

	// Default to local version if build version env variable is not set
	// This version connects to env varible pushed into docker container at build time
	build := os.Getenv(envBuildVersion)
	if build == "" {
		build = buildVersion
	}

	successBytes, _ := json.Marshal(model.WorldWireMessage{Msg: msg, Service: service, TimeStamp: timeNow,
		ParticipantID: ParticipantID, URL: URL, BuildVersion: build})

	return successBytes
}

func NotifySuccess(w http.ResponseWriter, r *http.Request, msg string) {
	//construct ww reponse
	successBytes := createMessage(msg, r)
	//send status OK
	Respond(w, http.StatusOK, successBytes)
	return
}

func NotifyFailure(w http.ResponseWriter, r *http.Request, statusCode int, msg string) {
	errBytes := createMessage(msg, r)
	//send failure status code
	Respond(w, statusCode, errBytes)
	return
}
