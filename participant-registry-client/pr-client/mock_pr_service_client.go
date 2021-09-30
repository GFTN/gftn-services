// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package pr_client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/GFTN/gftn-services/gftn-models/model"
	comn "github.com/GFTN/gftn-services/utility/common"
)

type MockPRServiceClient struct {
	URL                                                 string
	ParticipantRegistryGetAllParticipantsURL            string
	ParticipantRegistryGetParticipantsByCountryURL      string
	ParticipantRegistryDomainURL                        string
	ParticipantRegistryQuoteURL                         string
	ParticipantRegistryDistAccountURL                   string
	ParticipantRegistryIssuingAccountURL                string
	ParticipantRegistryGetParticipantsByAssetCountryURL string
	ParticipantRegistryGetParticipantByAddress          string
}

func CreateMockPRServiceClient(url string) (MockPRServiceClient, error) {

	client := MockPRServiceClient{}
	client.URL = url
	client.ParticipantRegistryDomainURL = client.URL + "/internal/pr/domain/"
	client.ParticipantRegistryQuoteURL = client.URL + "/internal/pr/assetpair/"
	client.ParticipantRegistryDistAccountURL = client.URL + "/internal/pr/account/"
	client.ParticipantRegistryIssuingAccountURL = client.URL + "/internal/pr/issuingaccount/"
	client.ParticipantRegistryGetAllParticipantsURL = client.URL + "/internal/pr"
	client.ParticipantRegistryGetParticipantsByCountryURL = client.URL + "/internal/pr/country/"
	client.ParticipantRegistryGetParticipantsByAssetCountryURL = client.URL + "/internal/pr/asset/country"
	client.ParticipantRegistryGetParticipantByAddress = client.URL + "/internal/pr/account"
	return client, nil

}

func (client MockPRServiceClient) GetParticipantByAddress(address string) (model.Participant, error) {
	return model.Participant{}, nil
}

func (client MockPRServiceClient) GetParticipantForDomain(participantID string) (model.Participant, error) {

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

func (client MockPRServiceClient) GetParticipantForIssuingAddress(address string) (model.Participant, error) {
	return model.Participant{}, nil
}

func (client MockPRServiceClient) GetParticipantDistAccount(domain string, account string) (string, error) {
	return "GCVK2PWR2ZVQBL3UUQP4Z2RO7M5NJWNPRIPQW3CNB5XBCPB4QK7V4AHX", nil
}

func (client MockPRServiceClient) GetParticipantIssuingAccount(domain string) (string, error) {
	return "GCVK2PWR2ZVQBL3UUQP4Z2RO7M5NJWNPRIPQW3CNB5XBCPB4QK7V4AHX", nil
}

func (client MockPRServiceClient) PostParticipantDistAccount(domain string, account model.Account) error {
	return nil
}

func (client MockPRServiceClient) GetAllParticipants() ([]model.Participant, error) {

	responseBody, err := ioutil.ReadFile(os.Getenv("GOPATH") + "/src/github.com/GFTN/gftn-services/new-quotes-service/unit-test-data/quote/participant_sandbox.json")

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

func (client MockPRServiceClient) GetParticipantsByCountry(countryCode string) ([]model.Participant, error) {
	return []model.Participant{}, nil
}

func (client MockPRServiceClient) PostParticipantIssuingAccount(domain string, account model.Account) error {
	return nil
}

func (client MockPRServiceClient) GetParticipantAccount(domain string, account string) (string, error) {
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
