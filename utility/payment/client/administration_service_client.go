// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/go-resty/resty"
	"github.com/GFTN/gftn-services/gftn-models/model"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
)

type RestAdministrationServiceClient struct {
	AdministrationServiceURL string
}

func CreateRestAdministrationServiceClient() (*RestAdministrationServiceClient, error) {
	client := &RestAdministrationServiceClient{}

	url := os.Getenv(global_environment.ENV_KEY_ADMIN_SVC_URL)
	if url == "" {
		LOGGER.Errorf("You MUST set the %v environment variable to point to Administration Service URL", global_environment.ENV_KEY_ADMIN_SVC_URL)
		os.Exit(1)
	}

	client.AdministrationServiceURL = url

	return client, nil
}

func (client RestAdministrationServiceClient) StoreFITOFICCTMemo(fitoficctmemo model.FitoFICCTMemoData) error {
	url := client.AdministrationServiceURL + "/internal/fitoficct"
	LOGGER.Infof("Doing internal administration service call: %v", url)
	bodyBytes, err := json.Marshal(&fitoficctmemo)
	response, err := resty.R().SetBody(bodyBytes).SetHeader("Content-type", "application/json").Post(url)
	if err != nil {
		LOGGER.Errorf("Error with calling administration service API:  %v", err.Error())
		return err
	}

	if response.StatusCode() != http.StatusOK {
		LOGGER.Warningf("Internal Administration Service returned a non-200 status (%v)", response.StatusCode())
		return errors.New("Returned status code with error ")
	}

	return nil
}

func (client *RestAdministrationServiceClient) GetTxnDetails(txnDetailsRequest model.FItoFITransactionRequest) (model.TransactionReceipt, int, error) {

	url := client.AdministrationServiceURL + "/internal/transaction"
	LOGGER.Info("Doing internal administration service call")
	bodyBytes, err := json.Marshal(&txnDetailsRequest)
	response, err := resty.R().SetBody(bodyBytes).SetHeader("Content-type", "application/json").Post(url)
	if err != nil {
		LOGGER.Errorf("Error with calling administration service API:  %v", err.Error())
		return model.TransactionReceipt{}, http.StatusNotFound, err
	}

	if response.StatusCode() != http.StatusOK {
		LOGGER.Warningf("Internal Administration Service returned a non-200 status (%v)", response.StatusCode())
		return model.TransactionReceipt{}, response.StatusCode(), errors.New(string(response.Body()[:]))
	}

	var txnResponse model.TransactionReceipt
	responseBody := response.Body()
	err = json.Unmarshal(responseBody, &txnResponse)

	if err != nil {
		LOGGER.Warningf("Error unmarshalling compliance response: %v", err.Error())
		return model.TransactionReceipt{}, http.StatusNotFound, err
	}
	return txnResponse, http.StatusOK, nil

}
