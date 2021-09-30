// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package client

import "github.com/GFTN/gftn-services/gftn-models/model"

type AdministrationServiceClient interface {
	GetTxnDetails(txnDetailsRequest model.FItoFITransactionRequest) ([]byte, int, error)
}
