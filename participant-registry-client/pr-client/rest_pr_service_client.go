// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package pr_client

import (
	"encoding/json"
	"net/http"

	"github.com/go-resty/resty"
	logging "github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/GFTN/gftn-services/gftn-models/model"
	comn "github.com/GFTN/gftn-services/utility/common"
)

var LOGGER = logging.MustGetLogger("pr-client")

type RestPRServiceClient struct {
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

func CreateRestPRServiceClient(url string) (RestPRServiceClient, error) {

	client := RestPRServiceClient{}
	client.URL = url
	client.ParticipantRegistryDomainURL = client.URL + "/internal/pr/domain/"
	client.ParticipantRegistryQuoteURL = client.URL + "/internal/pr/assetpair/"
	client.ParticipantRegistryDistAccountURL = client.URL + "/internal/pr/account/"
	client.ParticipantRegistryIssuingAccountURL = client.URL + "/internal/pr/issuingaccount/"
	client.ParticipantRegistryGetAllParticipantsURL = client.URL + "/internal/pr"
	client.ParticipantRegistryGetParticipantsByCountryURL = client.URL + "/internal/pr/country/"
	client.ParticipantRegistryGetParticipantsByAssetCountryURL = client.URL + "/internal/pr/asset/country"
	client.ParticipantRegistryGetParticipantByAddress = client.URL + "/internal/pr/account/"
	return client, nil

}

func (client RestPRServiceClient) GetParticipantByAddress(address string) (model.Participant, error) {
	participantRegistryURL := client.ParticipantRegistryGetParticipantByAddress + address
	LOGGER.Debug("PR-URL: ", participantRegistryURL)
	response, err := resty.R().Get(participantRegistryURL)
	if err != nil {
		LOGGER.Debugf("There was an error while querying the pr service for domain (%v):  %v", err)
		return model.Participant{}, err
	}

	if response.StatusCode() != http.StatusOK {
		LOGGER.Debugf("The response from the PR service was not 200.  Instead, it was %v - %v", response.StatusCode(), response.Status())
		return model.Participant{}, err
	}

	var participant model.Participant
	responseBodyBytes := response.Body()
	err = json.Unmarshal(responseBodyBytes, &participant)

	if err != nil {
		LOGGER.Debugf("In participant-registry-client:pr-client:rest_pr_service_client:GetParticipantByAddress: Error while marshalling response data:  %v", err)
		return model.Participant{}, err
	}
	return participant, nil

}

func (client RestPRServiceClient) GetAllParticipants() ([]model.Participant, error) {

	participantRegistryURL := client.ParticipantRegistryGetAllParticipantsURL
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

func (client RestPRServiceClient) GetParticipantsByCountry(countryCode string) ([]model.Participant, error) {

	participantRegistryURL := client.ParticipantRegistryGetParticipantsByCountryURL + countryCode
	LOGGER.Debugf("PR-URL: ", participantRegistryURL)
	response, err := resty.R().Get(participantRegistryURL)
	if err != nil {
		LOGGER.Debugf("There was an error while querying the pr service for country (%v):  %v", err)
		return []model.Participant{}, err
	}

	if response.StatusCode() != 200 {
		LOGGER.Debugf("The response from the PR service was not 200.  Instead, it was %v - %v", response.StatusCode(), response.Status())
		return []model.Participant{}, errors.New("Bad status response from remote PR service")
	}

	var participants []model.Participant
	responseBodyBytes := response.Body()
	err = json.Unmarshal(responseBodyBytes, &participants)

	if err != nil {
		LOGGER.Debugf("In participant-registry-client:pr-client:rest_pr_service_client:GetParticipantsByCountry: Error while marshalling response data:  %v", err)
		return []model.Participant{}, err
	}
	return participants, nil

}

func (client RestPRServiceClient) GetParticipantForDomain(domain string) (model.Participant, error) {

	participantRegistryURL := client.ParticipantRegistryDomainURL + domain
	LOGGER.Debugf("PR-URL: ", participantRegistryURL)
	response, err := resty.R().Get(participantRegistryURL)
	if err != nil {
		LOGGER.Errorf("There was an error while querying the pr service for domain (%v):  %v", err)
		return model.Participant{}, err
	}

	if response.StatusCode() != 200 {
		LOGGER.Warningf("The response from the PR service was not 200.  Instead, it was %v - %v", response.StatusCode(), response.Status())
		return model.Participant{}, errors.New("Bad status response from remote PR service")
	}

	var participant model.Participant
	responseBodyBytes := response.Body()
	err = json.Unmarshal(responseBodyBytes, &participant)

	if err != nil {
		LOGGER.Errorf("In participant-registry-client:pr-client:rest_pr_service_client:GetParticipantForDomain: Error while marshalling response data:  %v", err)
		return model.Participant{}, err
	}
	return participant, nil

}

func (client RestPRServiceClient) GetParticipantForIssuingAddress(accountAddress string) (model.Participant, error) {

	participantRegistryURL := client.ParticipantRegistryIssuingAccountURL + accountAddress
	LOGGER.Debugf("PR-URL: ", participantRegistryURL)
	response, err := resty.R().Get(participantRegistryURL)
	if err != nil {
		LOGGER.Errorf("There was an error while querying the pr service for domain (%v):  %v", err)
		return model.Participant{}, err
	}

	if response.StatusCode() != 200 {
		LOGGER.Warningf("The response from the PR service was not 200.  Instead, it was %v - %v", response.StatusCode(), response.Status())
		return model.Participant{}, errors.New("Bad status response from remote PR service")
	}

	var participant model.Participant
	responseBodyBytes := response.Body()
	err = json.Unmarshal(responseBodyBytes, &participant)

	if err != nil {
		LOGGER.Errorf("In participant-registry-client:pr-client:rest_pr_service_client:GetParticipantForDomain: Error while marshalling response data:  %v", err)
		return model.Participant{}, err
	}
	return participant, nil

}

func (client RestPRServiceClient) GetParticipantsForAssetPair(sc string, tc string) ([]model.Participant, error) {

	participantRegistryURL := client.ParticipantRegistryQuoteURL + sc + "/" + tc
	response, err := resty.R().Get(participantRegistryURL)
	if err != nil {
		LOGGER.Errorf("In: rest_pr_service_client - Error while getting participants for assetPair (%v):  %v", err)
		return []model.Participant{}, err
	}

	if response.StatusCode() != 200 {
		LOGGER.Warningf("The response from the PR service was not 200.  Instead, it was %v - %v", response.StatusCode(), response.Status())
		return []model.Participant{}, errors.New("Bad status response from remote PR service")
	}

	var participants []model.Participant
	responseBodyBytes := response.Body()
	err = json.Unmarshal(responseBodyBytes, &participants)

	if err != nil {
		LOGGER.Errorf("Error while marshalling response data:  %v", err)
		return []model.Participant{}, err
	}
	return participants, nil

}

func (client RestPRServiceClient) GetParticipantIssuingAccount(domain string) (string, error) {

	participantRegistryURL := client.ParticipantRegistryDomainURL + domain
	response, err := resty.R().Get(participantRegistryURL)
	if err != nil {
		LOGGER.Errorf("There was an error while querying the pr service for domain (%v):  %v", err)
		return "", err
	}

	if response.StatusCode() != 200 {
		LOGGER.Warningf("The response from the participant service was not 200.  Instead, it was %v - %v", response.StatusCode(), response.Status())
		return "", errors.New("Bad status response from remote participant service")
	}

	var participant model.Participant
	responseBodyBytes := response.Body()
	err = json.Unmarshal(responseBodyBytes, &participant)

	if err != nil {
		LOGGER.Debugf("error parsing PR data %v", err.Error())
		return "", errors.New("error parsing PR data")
	}
	accountKey := string(participant.IssuingAccount)
	LOGGER.Debugf("issuing account %v", accountKey)
	return accountKey, nil
}

func (client RestPRServiceClient) GetParticipantDistAccount(domain string, account string) (string, error) {

	participantRegistryURL := client.ParticipantRegistryDistAccountURL + domain + "/" + account

	response, err := resty.R().Get(participantRegistryURL)
	if err != nil {
		LOGGER.Errorf("There was an error while querying the pr service for domain (%v):  %v", err)
		return "", err
	}

	if response.StatusCode() != 200 {
		LOGGER.Warningf("The response from the participant service was not 200.  Instead, it was %v - %v", response.StatusCode(), response.Status())
		return "", errors.New("Bad status response from remote participant service")
	}

	accountKey := ""
	responseBodyBytes := response.Body()
	//LOGGER.Infof("Key received: %v", responseBodyBytes)

	accountKey = string(responseBodyBytes)

	return accountKey, nil
}

func (client RestPRServiceClient) PostParticipantDistAccount(domain string, account model.Account) error {

	bodyBytes, err := json.Marshal(&account)
	if err != nil {
		LOGGER.Errorf("decoding Operating account requested  %v", err)
		return err
	}
	participantRegistryURL := client.ParticipantRegistryDistAccountURL + domain
	LOGGER.Infof("PostParticipant Operating: %v: %v", participantRegistryURL, account.Name)

	resp, err := resty.R().SetBody(bodyBytes).SetHeader("Content-type", "application/json").Post(
		participantRegistryURL)

	if err != nil {
		LOGGER.Errorf("Error while making post to save operating account:  %v", err)
		return err
	}

	if resp.StatusCode() != 200 {
		LOGGER.Infof("PostParticipant Operating: account a non-200 response")
		msg := model.WorldWireError{}
		json.Unmarshal(resp.Body(), &msg)
		return errors.New(*msg.Details)
	}
	return nil
}

func (client RestPRServiceClient) PostParticipantIssuingAccount(domain string, account model.Account) error {
	bodyBytes, err := json.Marshal(&account)
	if err != nil {
		LOGGER.Errorf("decoding issuing account requested  %v", err)
		return err
	}
	participantRegistryURL := client.ParticipantRegistryIssuingAccountURL + domain
	LOGGER.Infof("PostParticipantIssuingAccount: %v: %v", participantRegistryURL, account.Address)

	resp, err := resty.R().SetBody(bodyBytes).SetHeader("Content-type", "application/json").Post(
		participantRegistryURL)

	if err != nil {
		LOGGER.Errorf("Error while making post to save Issuing account:  %v", err)
		return errors.New("Unable to connect participant registry")
	}

	if resp.StatusCode() != 200 {
		LOGGER.Infof("Post Issuing account a non-200 response")
		msg := model.WorldWireError{}
		json.Unmarshal(resp.Body(), &msg)
		return errors.New(*msg.Details)
	}
	return nil
}

// return err if account does not exist
func (client RestPRServiceClient) GetParticipantAccount(domain string, account string) (string, error) {
	var accountKey string
	var err error
	if account == comn.ISSUING {
		accountKey, err = client.GetParticipantIssuingAccount(domain)
		if accountKey == "" {
			return "", errors.New("participant account does not exist: issuing")
		}
	} else {
		accountKey, err = client.GetParticipantDistAccount(domain, account)
		if accountKey == "" {
			return "", errors.New("participant account does not exist: " + account)
		}
	}

	return accountKey, err
}
