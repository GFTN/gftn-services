// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package participantregistry

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/GFTN/gftn-services/gftn-models/model"
	comn "github.com/GFTN/gftn-services/utility/common"
)

type Client struct {
	HTTP *http.Client
	URL  string
}

func (client Client) GetAllParticipants() ([]model.Participant, error) {
	participantRegistryURL := client.URL + "/internal/pr"
	LOGGER.Debug("PR-URL: ", participantRegistryURL)

	req, _ := http.NewRequest("GET", participantRegistryURL, nil)
	res, err := client.HTTP.Do(req)
	if err != nil {
		LOGGER.Debugf("There was an error while querying the pr service for domain (%v):  %v", err)
		return []model.Participant{}, err
	}

	if res.StatusCode != http.StatusOK {
		LOGGER.Debugf("The response from the PR service was not 200.  Instead, it was %v - %v", res.StatusCode, res.Status, participantRegistryURL)
		return []model.Participant{}, errors.New("Bad status response from remote PR service")
	}

	var participants []model.Participant
	responseBodyBytes, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(responseBodyBytes, &participants)

	if err != nil {
		LOGGER.Debugf("In participant-registry-client:pr-client:rest_pr_service_client:GetAllParticipants: Error while marshalling response data:  %v", err)
		return []model.Participant{}, err
	}
	return participants, nil
}

func (client Client) GetParticipantForDomain(domain string) (model.Participant, error) {
	participantRegistryURL := client.URL + "/internal/pr/domain/" + domain

	LOGGER.Debugf("PR-URL: ", participantRegistryURL)
	req, _ := http.NewRequest("GET", participantRegistryURL, nil)
	res, err := client.HTTP.Do(req)
	if err != nil {
		LOGGER.Errorf("There was an error while querying the pr service for domain (%v):  %v", err)
		return model.Participant{}, err
	}

	if res.StatusCode != http.StatusOK {
		LOGGER.Warningf("The response from the PR service was not 200.  Instead, it was %v - %v", res.StatusCode, res.Status, participantRegistryURL)
		return model.Participant{}, errors.New("Bad status response from remote PR service")
	}

	var participant model.Participant
	responseBodyBytes, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(responseBodyBytes, &participant)

	if err != nil {
		LOGGER.Errorf("In participant-registry-client:pr-client:rest_pr_service_client:GetParticipantForDomain: Error while marshalling response data:  %v", err)
		return model.Participant{}, err
	} else {
		return participant, nil

	}
}

func (client Client) GetParticipantIssuingAccount(domain string) (string, error) {

	participantRegistryURL := client.URL + "/internal/pr/domain/" + domain

	req, _ := http.NewRequest("GET", participantRegistryURL, nil)
	res, err := client.HTTP.Do(req)
	if err != nil {
		LOGGER.Errorf("There was an error while querying the pr service for domain (%v):  %v", err)
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		LOGGER.Warningf("The response from the participant service was not 200.  Instead, it was %v - %v", res.StatusCode, res.Status, participantRegistryURL)
		return "", errors.New("Bad status response from remote participant service")
	}

	var participant model.Participant
	responseBodyBytes, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(responseBodyBytes, &participant)

	if err != nil {
		LOGGER.Debugf("error parsing PR data %v", err.Error())
		return "", errors.New("error parsing PR data")
	}
	accountKey := string(participant.IssuingAccount)
	LOGGER.Debugf(domain, "issuing account %v", accountKey)
	return accountKey, nil
}

func (client Client) GetParticipantDistAccount(domain string, account string) (string, error) {

	participantRegistryURL := client.URL + "/internal/pr/account/" + domain + "/" + account

	req, _ := http.NewRequest("GET", participantRegistryURL, nil)
	res, err := client.HTTP.Do(req)
	if err != nil {
		LOGGER.Errorf("There was an error while querying the pr service for domain (%v):  %v", err)
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		LOGGER.Warningf("The response from the participant service was not 200.  Instead, it was %v - %v", res.StatusCode, res.Status, participantRegistryURL)
		return "", errors.New("Bad status response from remote participant service")
	}

	accountKey := ""
	responseBodyBytes, _ := ioutil.ReadAll(res.Body)
	//LOGGER.Infof("Key received: %v", responseBodyBytes)
	accountKey = string(responseBodyBytes)

	return accountKey, nil
}

func (client Client) GetParticipantAccount(domain string, account string) (string, error) {
	var accountKey string
	var err error
	if account == comn.ISSUING {
		accountKey, err = client.GetParticipantIssuingAccount(domain)
	} else {
		accountKey, err = client.GetParticipantDistAccount(domain, account)
	}
	return accountKey, err
}
