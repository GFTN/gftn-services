// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package crypto_client

import (
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/utility/nodeconfig"
)

type MockCryptoServiceClient struct {
	URL              string
	CreateAccountURL string
}

func CreateMockCryptoServiceClient() (MockCryptoServiceClient, error) {

	client := MockCryptoServiceClient{}
	return client, nil

}

func (client MockCryptoServiceClient) CreateAccount(accountName string) (nodeconfig.Account, error, int, string) {
	return nodeconfig.Account{}, nil, 200, ""
}

func (client MockCryptoServiceClient) SignPayload(accountName string, payload []byte) (signedPayload []byte, err error, statusCode int, errorCode string) {
	return nil, nil, 200, ""
}

func (client MockCryptoServiceClient) SignXdr(accountName string, idUnsigned []byte, idSigned []byte, transactionUnsigned []byte) (transactionSigned []byte,
	err error, statusCode int, errorCode string) {
	return nil, nil, 200, ""
}

func (client MockCryptoServiceClient) AddIBMSign(transactionUnsigned []byte) (transactionSigned []byte,
	err error, statusCode int, errorCode string) {
	return nil, nil, 200, ""
}

func (client MockCryptoServiceClient) GetIBMAccount() (account model.Account, err error, statusCode int, errorCode string) {
	return model.Account{}, nil, 200, ""
}

func (client MockCryptoServiceClient) ParticipantSignXdr(accountName string, transactionUnsigned []byte) (transactionSigned []byte,
	err error, statusCode int, errorCode string) {
	return nil, nil, 200, ""
}
