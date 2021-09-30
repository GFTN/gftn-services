// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package whitelistservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/GFTN/gftn-services/gftn-models/model"
)

type Client struct {
	HTTP  *http.Client
	WLURL string
}

func (client *Client) GetWhiteListParticipantDomains(participantID string) ([]string, error) {
	resp, err := client.HTTP.Get(client.WLURL + "/gftn/whitelist/participants/" + participantID)
	LOGGER.Info(client.WLURL + "/gftn/whitelist/participants/" + participantID)
	if err != nil {
		return nil, err
	}
	var wlparitcipants []string
	readbuff, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(readbuff, &wlparitcipants)
	if err != nil {
		return nil, err
	}
	return wlparitcipants, nil
}

func (client *Client) GetWhiteListParticipants(participantID string) ([]model.Participant, error) {
	resp, err := client.HTTP.Get(client.WLURL + "/gftn/whitelist/participants/object/" + participantID)
	LOGGER.Info(client.WLURL + "/gftn/whitelist/participants/object/" + participantID)
	if err != nil {
		return nil, err
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
	requestBody["wl_participant_id"] = wlparitcipantID
	requestBodyByte, _ := json.Marshal(requestBody)
	res, err := client.HTTP.Post(client.WLURL+"/gftn/whitelist/participants/"+participantID, "application/json", bytes.NewBuffer(requestBodyByte))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("Create whitelist participant failed")
	}
	return nil
}

func (client *Client) DeleteWhiteListParticipants(participantID, wlparitcipantID string) error {
	requestBody := make(map[string]string)
	requestBody["wl_participant_id"] = wlparitcipantID
	requestBodyByte, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("DELETE", client.WLURL+"/gftn/whitelist/participants/"+participantID, bytes.NewBuffer(requestBodyByte))
	res, err := client.HTTP.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("Delete whitelist participant failed")
	}
	return nil
}
