// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package crypto_client

import (
	"encoding/json"
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/go-resty/resty"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/utility/nodeconfig"
)

var LOGGER = logging.MustGetLogger("crypto-client")

type RestCryptoServiceClient struct {
	Internal_URL                  string
	CreateAccountURL              string
	SignPayloadURL                string
	SignXdrURL                    string
	GetIBMAccountURL              string
	GetIBMSigURL                  string
	ParticipantSignXdrURL         string
	SignPayloadByMasterAccountURL string
}

func CreateRestCryptoServiceClient(internalUrl string) (RestCryptoServiceClient, error) {

	client := RestCryptoServiceClient{}

	if internalUrl == "" {
		return client, errors.New("environment variables are not correctly set for signing service to work!")
	}
	client.Internal_URL = internalUrl
	client.CreateAccountURL = client.Internal_URL + "/internal/account/"
	client.SignPayloadURL = client.Internal_URL + "/internal/request/sign"
	client.SignPayloadByMasterAccountURL = client.Internal_URL + "/internal/payload/wwsign"
	client.SignXdrURL = client.Internal_URL + "/internal/sign"
	client.ParticipantSignXdrURL = client.Internal_URL + "/internal/participant/sign"
	client.GetIBMAccountURL = client.Internal_URL + "/internal/admin/account"
	client.GetIBMSigURL = client.Internal_URL + "/internal/admin/sign"
	return client, nil

}

func (client RestCryptoServiceClient) CreateAccount(accountName string) (nodeconfig.Account, error, int, string) {

	createAccountURL := client.CreateAccountURL + accountName
	LOGGER.Debug("CYPTO-URL: ", createAccountURL)
	resp, err := resty.R().Post(createAccountURL)
	if err != nil {
		LOGGER.Debugf("There was an error while querying the Crypto service for domain %v", err.Error())
		return nodeconfig.Account{}, err, resp.StatusCode(), ""
	}

	if resp.StatusCode() != http.StatusCreated {
		LOGGER.Debugf("The response from the Crypto service was not 201.  Instead, it was %v - %v", resp.StatusCode(), resp.Status())
		msg := model.WorldWireError{}
		err = json.Unmarshal(resp.Body(), &msg)
		if err != nil {
			LOGGER.Debugf("Error marshalling err response:%v", err.Error())
			return nodeconfig.Account{}, errors.New("Error marshalling err response in Crypto service CreateAccount"), resp.StatusCode(), ""
		}
		return nodeconfig.Account{}, errors.New(*msg.Details), resp.StatusCode(), msg.Code

	}

	var account nodeconfig.Account
	responseBodyBytes := resp.Body()
	err = json.Unmarshal(responseBodyBytes, &account)

	if err != nil {
		LOGGER.Debugf("In crypto-client:rest_crypto_service_client:CreateAccount: Error while marshalling response data:  %v", err)
		return nodeconfig.Account{}, err, 409, ""
	}

	//success!
	return account, nil, resp.StatusCode(), ""

}
func (client RestCryptoServiceClient) SignPayload(accountName string, payload []byte) (signedPayload []byte, err error, statusCode int, errorCode string) {
	var payloadReq model.RequestPayload
	stPayload := strfmt.Base64(payload[:])

	payloadReq.AccountName = &accountName
	payloadReq.Payload = &stPayload

	LOGGER.Debug("CYPTO-URL: ", client.SignPayloadURL)
	LOGGER.Debugf("Signing with account %v", *payloadReq.AccountName)
	resp, err := resty.R().SetBody(payloadReq).Post(client.SignPayloadURL)
	if err != nil {
		LOGGER.Debugf("There was an error while querying the Crypto service :  %v", err.Error())
		return nil, err, resp.StatusCode(), ""
	}

	if resp.StatusCode() != http.StatusOK {
		LOGGER.Debugf("The response from the Crypto service was not 200.  Instead, it was %v - %v", resp.StatusCode(), resp.Status())
		msg := model.WorldWireError{}
		err = json.Unmarshal(resp.Body(), &msg)
		if err != nil {
			LOGGER.Debugf("Error marshalling err response:%v", err.Error())
			return nil, errors.New("Error marshalling err response in Crypto service SignPayload"), resp.StatusCode(), ""
		}
		return nil, errors.New(*msg.Details), resp.StatusCode(), msg.Code

	}

	var sig model.Signature
	responseBodyBytes := resp.Body()
	err = json.Unmarshal(responseBodyBytes, &sig)

	if err != nil {
		LOGGER.Debugf("In crypto-client:rest_crypto_service_client:CreateAccount: Error while marshalling response data:  %v", err)
		return nil, err, 409, ""
	}

	//success!
	return *sig.TransactionSigned, nil, resp.StatusCode(), ""
}

func (client RestCryptoServiceClient) SignPayloadByMasterAccount(payload []byte) (signedPayload []byte, err error) {
	var payloadReq model.RequestPayload
	stPayload := strfmt.Base64(payload[:])

	accountName := "ww"
	payloadReq.Payload = &stPayload
	payloadReq.AccountName = &accountName
	LOGGER.Debug("CYPTO-URL: ", client.SignPayloadByMasterAccountURL)
	resp, err := resty.R().SetBody(payloadReq).Post(client.SignPayloadByMasterAccountURL)
	if err != nil {
		LOGGER.Debugf("There was an error while querying the Crypto service :  %v", err.Error())
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		LOGGER.Debugf("The response from the Crypto service was not 200.  Instead, it was %v - %v", resp.StatusCode(), resp.Status())
		msg := model.WorldWireError{}
		err = json.Unmarshal(resp.Body(), &msg)
		if err != nil {
			LOGGER.Debugf("Error marshalling err response:%v", err.Error())
			return nil, errors.New("Error marshalling err response in Crypto service SignPayload")
		}
		return nil, errors.New(*msg.Details)

	}

	responseBodyBytes := resp.Body()

	//The returned response is JSON in format {"payload_with_signature": "Base64 encoded"}
	payloadWithSignature := model.PayloadWithSignature{}
	json.Unmarshal(responseBodyBytes, &payloadWithSignature)
	//LOGGER.Infof("Payload with signature: %s", payloadWithSignature.PayloadWithSignature)

	if err != nil {
		LOGGER.Debugf("In crypto-client:rest_crypto_service_client:CreateAccount: Error while marshalling response data:  %v", err)
		return nil, err
	}

	//success!
	return []byte(*payloadWithSignature.PayloadWithSignature), nil
}

func (client RestCryptoServiceClient) SignXdr(accountName string, idUnsigned []byte, idSigned []byte, transactionUnsigned []byte) (transactionSigned []byte,
	err error, statusCode int, errorCode string) {
	var draft model.Draft
	_idSigned := strfmt.Base64(idSigned[:])
	_idUnsigned := strfmt.Base64(idUnsigned[:])
	_transactionUnsigned := strfmt.Base64(transactionUnsigned[:])

	draft.AccountName = &accountName
	draft.IDSigned = &_idSigned
	draft.IDUnsigned = &_idUnsigned
	draft.TransactionID = ""
	draft.TransactionUnsigned = &_transactionUnsigned

	LOGGER.Debug("CYPTO-URL: ", client.SignXdrURL)
	resp, err := resty.R().SetBody(draft).Post(client.SignXdrURL)
	if err != nil {
		LOGGER.Debugf("There was an error while querying the Crypto service :  %v", err.Error())
		return nil, err, resp.StatusCode(), ""
	}

	if resp.StatusCode() != http.StatusCreated {
		LOGGER.Debugf("The response from the Crypto service was not 200.  Instead, it was %v - %v", resp.StatusCode(), resp.Status())
		msg := model.WorldWireError{}
		err = json.Unmarshal(resp.Body(), &msg)
		if err != nil {
			LOGGER.Debugf("Error marshalling err response:%v", err.Error())
			return nil, errors.New("Error marshalling err response in Crypto service SignXdr"), resp.StatusCode(), ""
		}
		return nil, errors.New(*msg.Details), resp.StatusCode(), msg.Code

	}

	var sig model.SignedTransaction
	responseBodyBytes := resp.Body()
	err = json.Unmarshal(responseBodyBytes, &sig)

	if err != nil {
		LOGGER.Debugf("In crypto-client:rest_crypto_service_client:sign XDR: Error while marshalling response data:  %v", err)
		return nil, err, 409, ""
	}

	//success!
	return *sig.TransactionSigned, nil, resp.StatusCode(), ""
}

func (client RestCryptoServiceClient) AddIBMSign(transactionUnsigned []byte) (transactionSigned []byte,
	err error, statusCode int, errorCode string) {
	var draft model.InternalAdminDraft
	_transactionUnsigned := strfmt.Base64(transactionUnsigned[:])
	draft.TransactionUnsigned = &_transactionUnsigned
	resp, err := resty.R().SetBody(draft).Post(client.GetIBMSigURL)
	if err != nil {
		LOGGER.Debugf("There was an error while signing with IBM account in the Crypto service  %v", err.Error())
		msg := model.WorldWireError{}
		err = json.Unmarshal(resp.Body(), &msg)
		if err != nil {
			LOGGER.Debugf("Error marshalling err response:%v", err.Error())
			return nil, errors.New("Error marshalling err response in Crypto service Get IBM account"), resp.StatusCode(), ""
		}
	}

	if resp.StatusCode() != http.StatusOK {
		LOGGER.Debugf("The response from the Crypto service was not 200.  Instead, it was %v - %v", resp.StatusCode(), resp.Status())
		return nil, errors.New(resp.String()), resp.StatusCode(), err.Error()
	}

	var sig model.SignedTransaction
	responseBodyBytes := resp.Body()
	err = json.Unmarshal(responseBodyBytes, &sig)

	if err != nil {
		LOGGER.Debugf("In crypto-client:rest_crypto_service_client: Error while marshalling response data:  %v", err)
		return nil, err, 409, ""
	}

	//success!
	return *sig.TransactionSigned, nil, resp.StatusCode(), ""
}

func (client RestCryptoServiceClient) GetIBMAccount() (account model.Account, err error, statusCode int, errorCode string) {
	resp, err := resty.R().Get(client.GetIBMAccountURL)
	if err != nil {
		LOGGER.Debugf("There was an error while getting IBM account in the Crypto service  %v", err.Error())
		msg := model.WorldWireError{}
		err = json.Unmarshal(resp.Body(), &msg)
		if err != nil {
			LOGGER.Debugf("Error marshalling err response:%v", err.Error())
			return model.Account{}, errors.New("Error marshalling err response in Crypto service Get IBM account"), resp.StatusCode(), ""
		}
	}

	if resp.StatusCode() != http.StatusOK {
		LOGGER.Debugf("The response from the Crypto service was not 200.  Instead, it was %v - %v", resp.StatusCode(), resp.Status())
		return model.Account{}, errors.New(resp.String()), resp.StatusCode(), ""
	}

	var ibmAccount model.Account
	responseBodyBytes := resp.Body()
	err = json.Unmarshal(responseBodyBytes, &ibmAccount)

	if err != nil {
		LOGGER.Debugf("In crypto-client:rest_crypto_service_client: Error while marshalling response data:  %v", err)
		return model.Account{}, err, 409, ""
	}

	//success!
	return ibmAccount, nil, resp.StatusCode(), ""
}

func (client RestCryptoServiceClient) ParticipantSignXdr(accountName string, transactionUnsigned []byte) (transactionSigned []byte,
	err error, statusCode int, errorCode string) {
	var draft model.InternalDraft
	_transactionUnsigned := strfmt.Base64(transactionUnsigned[:])
	draft.AccountName = &accountName
	draft.TransactionID = ""
	draft.TransactionUnsigned = &_transactionUnsigned

	LOGGER.Debug("CYPTO-URL: ", client.ParticipantSignXdrURL)
	resp, err := resty.R().SetBody(draft).Post(client.ParticipantSignXdrURL)
	if err != nil {
		LOGGER.Debugf("There was an error while querying the Crypto service :  %v", err.Error())
		return nil, err, resp.StatusCode(), ""
	}

	if resp.StatusCode() != http.StatusCreated {
		LOGGER.Debugf("The response from the Crypto service was not 200.  Instead, it was %v - %v", resp.StatusCode(), resp.Status())
		msg := model.WorldWireError{}
		err = json.Unmarshal(resp.Body(), &msg)
		if err != nil {
			LOGGER.Debugf("Error marshalling err response:%v", err.Error())
			return nil, errors.New("Error marshalling err response in Crypto service ParticipantSignXdr"), resp.StatusCode(), ""
		}
		return nil, errors.New(*msg.Details), resp.StatusCode(), msg.Code

	}

	var sig model.SignedTransaction
	responseBodyBytes := resp.Body()
	err = json.Unmarshal(responseBodyBytes, &sig)

	if err != nil {
		LOGGER.Debugf("In crypto-client:rest_crypto_service_client: ParticipantSignXdr : Error while marshalling response data:  %v", err)
		return nil, err, 409, ""
	}

	//success!
	return *sig.TransactionSigned, nil, resp.StatusCode(), ""
}
