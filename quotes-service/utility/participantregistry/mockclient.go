// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package participantregistry

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	comn "github.com/GFTN/gftn-services/utility/common"
	"github.com/GFTN/gftn-services/gftn-models/model"
)

type MockClient struct {
	HTTP *http.Client
	URL  string
}

func CreateMockClient() (MockClient, error) {

	client := MockClient{}
	return client, nil

}

func (client *MockClient) GetAllParticipants() ([]model.Participant, error) {

	responseBody, err := ioutil.ReadFile(os.Getenv("GOPATH") + "/src/github.com/GFTN/gftn-services/quotes-service/unit-test-data/quote/participant_sandbox.json")

	if err != nil {
		LOGGER.Debugf("There was an error while querying the pr service for domain (%v):  %v", err)
		return []model.Participant{}, err
	}
	var participants []model.Participant
	// LOGGER.Debugf(string(responseBody))
	var result map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	response, _ := json.Marshal(result["participants"])
	err = json.Unmarshal(response, &participants)

	if err != nil {
		LOGGER.Debugf("In participant-registry-client:pr-client:rest_pr_service_client:GetAllParticipants: Error while marshalling response data:  %v", err)
		return []model.Participant{}, err
	} else {
		return participants, nil
	}
}

func (client *MockClient) GetParticipantForDomain(participantID string) (model.Participant, error) {

	responseBody, err := ioutil.ReadFile(os.Getenv("GOPATH") + "/src/github.com/GFTN/gftn-services/quotes-service/unit-test-data/quote/participant_sandbox.json")
	var participants []model.Participant
	// LOGGER.Debugf(string(responseBody))
	var result map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	response, _ := json.Marshal(result["participants"])
	err = json.Unmarshal(response, &participants)

	if err != nil {
		LOGGER.Debugf("In participant-registry-client:pr-client:rest_pr_service_client:GetParticipantForDomain: Error while marshalling response data:  %v", err)
		return model.Participant{}, err
	}
	for _, participant := range participants {
		if *participant.ID == participantID {
			return participant, nil
		}
	}
	return model.Participant{}, errors.New("In participant-registry-client: No participant returned")
}

func (client *MockClient) GetParticipantAccount(domain string, account string) (string, error) {
	participant, _ := client.GetParticipantForDomain(domain)
	if account == comn.ISSUING {
		return participant.IssuingAccount, nil
	} else {
		for _, operatingAccount := range participant.OperatingAccounts {
			if operatingAccount.Name == account {
				return *operatingAccount.Address, nil
			}
		}
	}
	return "", nil
}
