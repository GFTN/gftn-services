// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package blocklist_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/GFTN/gftn-services/gftn-models/model"
)

type Client struct {
	HTTPClient *http.Client
	AdminUrl   string
}

var endpoints map[string]string

func init() {
	endpoints = map[string]string{
		"add":      "/internal/blocklist",
		"remove":   "/internal/blocklist",
		"get":      "/internal/blocklist",
		"validate": "/internal/blocklist/validate",
	}

}

func (client *Client) GetBlocklist(blocklistType string) ([]model.Blocklist, error) {
	queryUrl := client.AdminUrl + endpoints["get"]
	if blocklistType != "" {
		queryUrl += "?type=" + blocklistType
	}
	LOGGER.Infof("GetBlocklist URL:" + queryUrl)
	res, err := client.HTTPClient.Get(queryUrl)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		return []model.Blocklist{}, errors.New("Get blocklist record failed")
	}

	var blocklists []model.Blocklist
	readbuff, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(readbuff, &blocklists)
	if err != nil {
		return nil, err
	}
	return blocklists, nil
}

func (client *Client) ValidateFromBlocklist(countries []string, currencies []string, participants []string) (string, error) {
	queryUrl := client.AdminUrl + endpoints["validate"]
	LOGGER.Infof("ValidateFromBlocklist URL: " + queryUrl)
	LOGGER.Infof("Validating Currency/Country/Participant ID in Blocklist...")

	if len(countries) == 0 && len(currencies) == 0 && len(participants) == 0 {
		LOGGER.Errorf("No payload attached")
		return "", errors.New("No payload attached")
	}

	payload := "[{\"type\":\"country\",\"value\":[" + strings.Join(countries, ", ") + "]},{\"type\":\"currency\",\"value\":[" + strings.Join(currencies, ", ") + "]},{\"type\":\"institution\",\"value\":[" + strings.Join(participants, ", ") + "]}]"

	res, err := client.HTTPClient.Post(queryUrl, "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		LOGGER.Errorf("Something is wrong when sending the validate request: %s", err)
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", errors.New("Validate blocklist record failed")
	}
	byteResult, _ := ioutil.ReadAll(res.Body)
	LOGGER.Infof("Blocklist validation success!")
	return string(byteResult), nil
}

func (client *Client) AddBlocklist(payload string) (string, error) {
	queryUrl := client.AdminUrl + endpoints["add"]
	LOGGER.Infof("AddBlocklist URL: " + queryUrl)
	if payload == "" {
		LOGGER.Errorf("No payload attached")
		return "", errors.New("No payload attached")
	}

	res, err := client.HTTPClient.Post(queryUrl, "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", err
	}
	LOGGER.Infof("res: %s", res)
	if res.StatusCode != http.StatusOK {
		return "", errors.New("Create blocklist record failed")
	}
	byteResult, _ := ioutil.ReadAll(res.Body)
	return string(byteResult), nil

}

func (client *Client) RemoveBlocklist(payload string) (string, error) {
	queryUrl := client.AdminUrl + endpoints["remove"]
	LOGGER.Infof("RemoveBlocklist URL: " + queryUrl)

	if payload == "" {
		LOGGER.Errorf("No payload attached")
		return "", errors.New("No payload attached")
	}

	req, _ := http.NewRequest("DELETE", queryUrl, bytes.NewBuffer([]byte(payload)))
	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", errors.New("Delete blocklist record failed")
	}
	byteResult, _ := ioutil.ReadAll(res.Body)
	return string(byteResult), nil

}
