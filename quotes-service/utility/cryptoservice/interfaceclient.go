// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package cryptoservice

import "github.com/GFTN/gftn-services/gftn-models/model"

type InterfaceClient interface {
	RequestSigning(txeBase64 string, requestBase64 string, signedRequestBase64 string, accountName string, participant model.Participant) (string, error)
}
