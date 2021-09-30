// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package whitelistclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/utility/participant"
)

type Client struct {
	HTTPClient *http.Client
	WLURL      string
}

func (client *Client) GetWhiteListParticipantDomains(participantID string) ([]string, error) {
	resp, err := client.HTTPClient.Get(client.WLURL + "/internal/participants/whitelist/" + participantID)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		LOGGER.Error("The response from the whitelist service was not 200.  Instead, it was", resp.Status)
		return nil, errors.New("Get WhiteList Participant Domains failed")
	}
	var wlparitcipants []string
	readbuff, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(readbuff, &wlparitcipants)
	if err != nil {
		return nil, err
	}
	return wlparitcipants, nil
}

func (client *Client) IsParticipantWhiteListed(participantID string, targetDomain string) (bool, error) {

	participants, err := client.GetWhiteListParticipantDomains(participantID)

	if err != nil {
		return false, err
	}
	for _, participant := range participants {
		if participant == targetDomain {
			return true, nil
		}
	}
	return false, nil
}

func (client *Client) GetWhiteListParticipants(participantID string) ([]model.Participant, error) {
	resp, err := client.HTTPClient.Get(client.WLURL + "/internal/participants/whitelist/" + participantID + "/object")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		LOGGER.Error("The response from the whitelist service was not 200.  Instead, it was", resp.Status)
		return nil, errors.New("Get WhiteList Participant failed")
	}
	var wlparitcipants []model.Participant
	readbuff, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(readbuff, &wlparitcipants)
	if err != nil {
		return nil, err
	}
	return wlparitcipants, nil
}

func (client *Client) CreateWhiteListParticipants(participantID, wlparitcipantID string) error {
	requestBody := make(map[string]string)
	requestBody["participant_id"] = wlparitcipantID
	requestBodyByte, _ := json.Marshal(requestBody)
	resp, err := client.HTTPClient.Post(client.WLURL+"/internal/participants/whitelist/"+participantID, "application/json", bytes.NewBuffer(requestBodyByte))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		LOGGER.Error("The response from the whitelist service was not 200.  Instead, it was", resp.Status)
		return errors.New("Create whitelist participant failed")
	}
	return nil
}

func (client *Client) DeleteWhiteListParticipants(participantID, wlparitcipantID string) error {
	requestBody := make(map[string]string)
	requestBody["participant_id"] = wlparitcipantID
	requestBodyByte, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("DELETE", client.WLURL+"/internal/participants/whitelist/"+participantID, bytes.NewBuffer(requestBodyByte))
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		LOGGER.Error("The response from the whitelist service was not 200.  Instead, it was", resp.Status)
		return errors.New("Delete whitelist participant failed")
	}
	return nil
}

func (client *Client) GetMutualWhiteListParticipantDomains(participantID string) ([]string, error) {
	url := client.WLURL + "/internal/participants/whitelist/" + participantID + "/mutual"
	resp, err := client.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		LOGGER.Error("The response from the whitelist service was not 200.  Instead, it was", resp.Status)
		return nil, errors.New("Get Mutual WhiteList Participant Domains failed")
	}
	var wlparitcipants []string
	readbuff, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(readbuff, &wlparitcipants)
	if err != nil {
		return nil, err
	}
	return wlparitcipants, nil
}

func (client *Client) GetMutualWhiteListParticipants(participantID string) ([]model.Participant, error) {
	url := client.WLURL + "/internal/participants/whitelist/" + participantID + "/mutual/object"
	resp, err := client.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		LOGGER.Error("The response from the whitelist service was not 200.  Instead, it was", resp.Status)
		return nil, errors.New("Get Mutual WhiteList Participant failed")
	}
	var wlparitcipants []model.Participant
	readbuff, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(readbuff, &wlparitcipants)
	if err != nil {
		return nil, err
	}

	wlparitcipants = participant.ExtractActiveParticipants(wlparitcipants)

	return wlparitcipants, nil
}
