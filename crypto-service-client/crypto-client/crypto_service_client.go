// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package crypto_client

import (
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/utility/nodeconfig"
)

type CryptoServiceClient interface {
	CreateAccount(accountName string) (account nodeconfig.Account, err error, statusCode int, errorCode string)
	SignPayload(accountName string, payload []byte) (signedPayload []byte, err error, statusCode int, errorCode string)
	SignXdr(accountName string, idUnsigned []byte, idSigned []byte, transactionUnsigned []byte) (transactionSigned []byte,
		err error, statusCode int, errorCode string)
	ParticipantSignXdr(accountName string, transactionUnsigned []byte) (transactionSigned []byte,
		err error, statusCode int, errorCode string)
	AddIBMSign(transactionUnsigned []byte) (transactionSigned []byte,
		err error, statusCode int, errorCode string)
	GetIBMAccount() (account model.Account, err error, statusCode int, errorCode string)
}
