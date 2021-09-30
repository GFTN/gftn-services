// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package gasserviceclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
)

type HTTPClient interface {
	HTTPPoster
	HTTPGetter
}

type HTTPPoster interface {
	Post(url, contentType string, body io.Reader) (*http.Response, error)
}

type HTTPGetter interface {
	Get(url string) (*http.Response, error)
}

type Client struct {
	HTTP HTTPClient
	URL  string
}

type AccountSeqResponse struct {
	Pkey           string
	SequenceNumber string
}
type SubmitTexResponse struct {
	Title  string
	Ledger uint64
	Hash   string
}

func (gs *Client) GetAccountAndSequence() (string, uint64, error) {

	attempts, _ := strconv.Atoi(os.Getenv(global_environment.ENV_KEY_GAS_ACCOUNT_ATTEMPTS))
	waitDuration, _ := strconv.ParseInt(os.Getenv(global_environment.ENV_KEY_WAIT_UNLOCK_DURATION), 10, 64)
	accountSeqResponse, err := RetryGetAccountAndSequence(gs, attempts, time.Duration(waitDuration)*time.Second)
	if err != nil {
		return "", 0, errors.New("GasService: Unable to obtain IBMAccount and Seq Num")
	}

	// seq := accountSeqResponse.SequenceNumber
	seq, err := strconv.ParseUint(accountSeqResponse.SequenceNumber, 10, 64)
	LOGGER.Info("seq number:", accountSeqResponse.SequenceNumber)
	return accountSeqResponse.Pkey, seq, nil
}

func (gs *Client) SubmitTxe(txeOfiRfiSignedB64 string) (string, uint64, error) {
	jsonValue := `{
		"oneSignedXDR": "` + txeOfiRfiSignedB64 + `" }`
	resp, err := gs.HTTP.Post(gs.URL+"/signXDRAndExecuteXDR", "application/json", bytes.NewBuffer([]byte(jsonValue)))
	if err != nil {
		return "", 0, err
	}
	if resp.StatusCode == http.StatusOK {
		submitTexResponse := SubmitTexResponse{}
		readbuff, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			LOGGER.Error(err)
			return "", 0, err
		}
		err = json.Unmarshal(readbuff, &submitTexResponse)
		if err != nil {
			return "", 0, err
		}
		LOGGER.Info("Submiting Txe to Gas Service Success. Status: " + resp.Status)
		return submitTexResponse.Hash, submitTexResponse.Ledger, nil
	}
	if resp.StatusCode == http.StatusBadRequest {
		submitTexResponse := make(map[string]interface{})
		readbuff, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			LOGGER.Error(err)
			return "", 0, err
		}
		err = json.Unmarshal(readbuff, &submitTexResponse)
		if err != nil {
			return "", 0, err
		}
		LOGGER.Error(submitTexResponse)
		return "", 0, errors.New("Submiting Txe to Gas Service Failed. Status: " + resp.Status)
	}
	if resp.StatusCode == http.StatusForbidden {
		submitTexResponse := make(map[string]interface{})
		readbuff, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			LOGGER.Error(err)
			return "", 0, err
		}
		err = json.Unmarshal(readbuff, &submitTexResponse)
		if err != nil {
			return "", 0, err
		}
		LOGGER.Error(submitTexResponse)
		return "", 0, errors.New("Submiting Txe to Gas Service Failed. Status: " + resp.Status)
	}

	return "", 0, errors.New("GasService: Unforseen Status:" + resp.Status)
}

func RetryGetAccountAndSequence(gs *Client, attempts int, sleep time.Duration) (AccountSeqResponse, error) {
	var err error
	for i := 0; ; i++ {
		LOGGER.Infof("The %v time attempt getting account", i+1)
		accountSeqResponse := AccountSeqResponse{}
		resp, err := gs.HTTP.Get(gs.URL + "/lockaccount")
		if err != nil {
			LOGGER.Error(err)
			return accountSeqResponse, err
		}
		readbuff, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			LOGGER.Error(err)
			return accountSeqResponse, err
		}
		_ = json.Unmarshal(readbuff, &accountSeqResponse)

		if accountSeqResponse.Pkey != "" {
			return accountSeqResponse, nil
		} else if err == nil {
			err = errors.New("No available account at the moment")
		}

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(sleep)

		LOGGER.Infof("retrying after error: %s", err.Error())
	}
	LOGGER.Errorf("after %d attempts, last error: %s", attempts, err)
	return AccountSeqResponse{}, err
}
