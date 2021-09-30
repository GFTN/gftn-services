// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package prclient

import (
	"encoding/json"
	"net/http"

	"github.com/go-resty/resty"
	logging "github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/GFTN/gftn-services/gftn-models/model"
)

var LOGGER = logging.MustGetLogger("pr-client")

type Client struct {
	HTTPClient *http.Client
	URL        string
}

func (client Client) GetAllParticipants() ([]model.Participant, error) {
	participantRegistryURL := client.URL + "/internal/pr"
	LOGGER.Debug("PR-URL: ", participantRegistryURL)
	response, err := resty.R().Get(participantRegistryURL)
	if err != nil {
		LOGGER.Debugf("There was an error while querying the pr service for domain (%v):  %v", err)
		return []model.Participant{}, err
	}

	if response.StatusCode() != http.StatusOK {
		LOGGER.Debugf("The response from the PR service was not 200.  Instead, it was %v - %v", response.StatusCode(), response.Status())
		return []model.Participant{}, errors.New("Bad status response from remote PR service")
	}

	var participants []model.Participant
	responseBodyBytes := response.Body()
	err = json.Unmarshal(responseBodyBytes, &participants)

	if err != nil {
		LOGGER.Debugf("In participant-registry-client:pr-client:rest_pr_service_client:GetAllParticipants: Error while marshalling response data:  %v", err)
		return []model.Participant{}, err
	}
	return participants, nil
}
