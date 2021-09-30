// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package automate

import (
	"errors"
	"github.com/GFTN/gftn-services/automation-service/constant"
	"net/http"
	"os"
	"time"

	"github.com/op/go-logging"
	"github.com/GFTN/gftn-services/automation-service/internal_model"
	"github.com/GFTN/gftn-services/automation-service/utility"
	"github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/response"
)

var LOGGER = logging.MustGetLogger("service-check")

type ServiceCheck struct {
}

func InitiateServiceCheck() (ServiceCheck, error) {
	op := ServiceCheck{}
	return op, nil
}

func (op *ServiceCheck) Check(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("X-XSS-Protection", "1")
	LOGGER.Debugf("Running K8s testing script")
	script := constant.K8sBasePath + "/script/change_cluster.sh"
	err := utility.RunBashCmd(script, os.Getenv(global_environment.ENV_KEY_ENVIRONMENT_VERSION))
	if err != nil {
		LOGGER.Errorf("Error: %s", err.Error())
		response.NotifyWWError(w, req, http.StatusBadRequest, "API-1100", errors.New("service check failed"))
		return
	}

	LOGGER.Debugf("PR-service connection testing")
	prServiceURL := os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL)
	prClient := &internal_model.Client{
		HTTPClient: &http.Client{Timeout: time.Second * 15},
		URL:        prServiceURL,
	}

	res, err := prClient.HTTPClient.Get(prClient.URL + "/internal/service_check")
	if err != nil || res.StatusCode != http.StatusOK {
		LOGGER.Error("Unable to hit pr-serivce")
		response.NotifyWWError(w, req, http.StatusBadRequest, "API-1100", errors.New("unable to hit pr-serivce"))
		return
	}

	utility.FetchDockerImageTag()

	LOGGER.Debugf("Success")
	response.NotifySuccess(w, req, "Service is good")
	return
}
