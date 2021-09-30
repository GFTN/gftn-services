// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package gasserviceclient

type GasServiceClient interface {
	GetAccountAndSequence() (string, uint64, error)
	SubmitTxe(txeOfiRfiSignedB64 string) (string, uint64, error)
}
