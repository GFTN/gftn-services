// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package cryptoservice

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/GFTN/gftn-services/gftn-models/model"
	participantutil "github.com/GFTN/gftn-services/utility/participant"
)

type Client struct {
	HTTP        *http.Client
	URLTemplate string
}
type SigningRequest struct {
	AccountName string `json:"account_name"`

	// This will be signed reference envelope to verify against partcipant's signature for authenticity.
	IdentificationSigned string `json:"id_signed,omitempty"`

	// This will be unsigned reference envelope to verify against partcipant's signature for authenticity.
	IdentificationUnsigned string `json:"id_unsigned,omitempty"`

	// reference transaction id, will be
	TransactionID string `json:"transaction_id,omitempty"`

	// unsigned transaction envelope to be signed by the participant
	TransactionUnsigned string `json:"transaction_unsigned,omitempty"`
}

func (client *Client) RequestSigning(txeBase64 string, requestBase64 string, signedRequestBase64 string, accountName string, participant model.Participant) (string, error) {

	// var txeSFB64, requestSFB64, signedRequestSFB64 strfmt.Base64
	// err := txeSFB64.UnmarshalText([]byte("AAAAACIKcSda2GY1UmuKYyRF2uvTJPI6uhi1tYQ/MzZktwAQAAAAyAANhhcAAAAeAAAAAQAAAAAAAAAAAAAAAGIFUsIAAAAAAAAAAgAAAAEAAAAAN56OXjAeFHiGiWIiJUAocnJK3tU6wm3JxUfakiTHSNkAAAABAAAAAF49wJinSedEzsd5aWgwQWSQs2akIOI9+A8HnQemh+B6AAAAAlNHRERPAAAAAAAAAAAAAABePcCYp0nnRM7HeWloMEFkkLNmpCDiPfgPB50HpofgegAAAAAAmJaAAAAAAQAAAABePcCYp0nnRM7HeWloMEFkkLNmpCDiPfgPB50HpofgegAAAAEAAAAAN56OXjAeFHiGiWIiJUAocnJK3tU6wm3JxUfakiTHSNkAAAACVEhCRE8AAAAAAAAAAAAAADeejl4wHhR4holiIiVAKHJySt7VOsJtycVH2pIkx0jZAAAAAABcQ4oAAAAAAAAAAA=="))
	// requestSFB64.UnmarshalText([]byte(requestBase64))
	// signedRequestSFB64.UnmarshalText([]byte(signedRequestBase64))
	requestBody := &SigningRequest{
		AccountName:            accountName,
		IdentificationSigned:   signedRequestBase64, //signature
		IdentificationUnsigned: requestBase64,       //unsigned data
		TransactionID:          "",
		TransactionUnsigned:    txeBase64,
	}
	requestBodyByte, _ := json.Marshal(requestBody)
	url, err := participantutil.GetServiceUrl(client.URLTemplate, *participant.ID)
	if err != nil {
		return "", err
	}
	req, _ := http.NewRequest("POST", url+"/internal/sign", bytes.NewBuffer(requestBodyByte))
	// req, _ := http.NewRequest("POST", "http://"+*participant.ID+"-cryptoservice:10042"+"/v1/internal/sign", bytes.NewBuffer(requestBodyByte))
	res, err := client.HTTP.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusCreated {
		LOGGER.Debugf("The response from the Crypto service was not 200.  Instead, it was %v - %v", res.StatusCode, res.Status)
		msg := model.WorldWireError{}
		readbuff, _ := ioutil.ReadAll(res.Body)
		err = json.Unmarshal(readbuff, &msg)
		if err != nil {
			LOGGER.Debugf("Error marshalling err response:%v", err.Error())
			return "", errors.New("Error marshalling err response in Crypto service SignXdr")
		}
		return "", errors.New(*msg.Details)
	}
	var signedTxObject model.SignedTransaction
	readbuff, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(readbuff, &signedTxObject)
	if err != nil {
		return "", err
	}
	signedTxe := *signedTxObject.TransactionSigned

	signedTxeBase64 := base64.StdEncoding.EncodeToString(signedTxe)
	return signedTxeBase64, nil
}
