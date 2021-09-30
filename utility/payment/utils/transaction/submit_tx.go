// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package transaction

import (
	"errors"

	"github.com/stellar/go/xdr"
	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"
)

func (op *CreateFundingOpereations) createStellarTransaction(ibmAccount, sendingAccount, receivingAccount, settlementAccountName string, signData sendmodel.SignData, seqNum uint64, memoHash xdr.Memo) (string, error) {
	signedTx, signErr := op.signHandler.SignTx(ibmAccount, sendingAccount, receivingAccount, settlementAccountName, seqNum, &signData, memoHash)
	if signErr != nil {
		LOGGER.Warningf("Fail to sign transaction")
		return "", signErr
	}

	return signedTx, nil
}

func (op *CreateFundingOpereations) submitToStellar(reqMsgId, msgId, signedTx string) (string, error) {
	//call gas service to sign and submit the transaction
	txHash, _, restErr := op.GasClient.SubmitTxe(signedTx)
	if restErr != nil {
		LOGGER.Errorf("Failed to submit transaction to Stellar network.")
		return "", restErr
	}

	if txHash == "" {
		return "", errors.New("transaction failed")
	}

	return txHash, nil
}
