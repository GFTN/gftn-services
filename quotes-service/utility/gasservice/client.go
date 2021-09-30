// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package gasservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Client struct {
	HTTP *http.Client
	URL  string
}

type AccountSeqResponse struct {
	Pkey           string
	SequenceNumber uint64
}
type SubmitTexResponse struct {
	Title  string
	Ledger uint64
	Hash   string
}

func (gs *Client) GetAccountAndSequence() (string, uint64, error) {
	resp, err := gs.HTTP.Get(gs.URL + "/lockaccount")
	if err != nil {
		return "", 0, err
	}
	accountSeqResponse := AccountSeqResponse{}
	readbuff, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(readbuff, &accountSeqResponse)
	if err != nil {
		return "", 0, err
	}
	if accountSeqResponse.Pkey == "" {
		return "", 0, errors.New("GasService: Unable to obtain IBMAccount and Seq Num")
	}
	seq := accountSeqResponse.SequenceNumber
	// seq, err := strconv.ParseUint(accountSeqResponse.SequenceNumber, 10, 64)

	return accountSeqResponse.Pkey, seq, nil
}

func (gs *Client) SubmitTxe(txeOfiRfiSignedB64 string) (string, uint64, error) {
	jsonValue := `{
		"oneSignedXDR": "` + txeOfiRfiSignedB64 + `" }`
	resp, err := gs.HTTP.Post(gs.URL+"/signXDRAndExecuteXDR", "application/json", bytes.NewBuffer([]byte(jsonValue)))
	if err != nil {
		return "", 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return "", 0, errors.New("GasService: Status Code != 200; Returned status code: " + strconv.Itoa(resp.StatusCode))
	}
	submitTexResponse := SubmitTexResponse{}
	readbuff, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(readbuff, &submitTexResponse)
	if err != nil {
		return "", 0, err
	}

	if submitTexResponse.Title == "Transaction successful" {
		return submitTexResponse.Hash, submitTexResponse.Ledger, nil
	}
	if submitTexResponse.Title == "Source Account Expire" {
		return "", 0, errors.New("GasService: Source Account Expire")
	}
	return "", 0, errors.New("GasService: " + submitTexResponse.Title)
}
