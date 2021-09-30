// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package transaction

import (
	"net/http"
	"os"
	"time"

	"github.com/stellar/go/xdr"
	gasserviceclient "github.com/GFTN/gftn-services/gas-service-client"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/participant"
	"github.com/GFTN/gftn-services/utility/payment/client"
	"github.com/GFTN/gftn-services/utility/payment/constant"
	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"
	"github.com/GFTN/gftn-services/utility/payment/utils/signing"
	"github.com/GFTN/gftn-services/utility/vault/utils"
)

type CreateFundingOpereations struct {
	gasServiceURL string
	admClient     *client.RestAdministrationServiceClient
	homeDomain    string
	serviceName   string
	prServiceURL  string
	signHandler   signing.CreateSignOperations
	GasClient     gasserviceclient.Client
}

func InitiateFundingOperations(pr, domain string) (op CreateFundingOpereations) {
	op.gasServiceURL = os.Getenv(global_environment.ENV_KEY_GAS_SVC_URL)
	op.homeDomain = domain
	//op.serviceName = os.Getenv(global_environment.ENV_KEY_SERVICE_NAME)
	op.admClient, _ = client.CreateRestAdministrationServiceClient()
	op.signHandler = signing.InitiateSignOperations(pr)
	op.prServiceURL = pr

	op.GasClient = gasserviceclient.Client{
		HTTP: &http.Client{Timeout: time.Second * 80},
		URL:  op.gasServiceURL,
	}

	return op
}

func (op *CreateFundingOpereations) FundAndSubmitPaymentTransaction(rfiAccount, reqMsgId, msgId, xmlMsgType, rfiSettlementAccountName string, dbData sendmodel.SignData, memoHash xdr.Memo) (int, string, string) {
	var sendingAccount, receivingAccount, settlementAccountName, ofiAccount string

	switch xmlMsgType {
	case constant.IBWF001:
		// if OFI receive a ibwf001 message, the transaction sender will be OFI and receiver will be RFI and will use OFI settlement account to sign the transaction
		receivingAccount = rfiAccount
		settlementAccountName = dbData.SettlementAccountName
		account, getErr := participant.GenericGetAccount(utils.Session{}, dbData.SettlementAccountName)
		if getErr != nil {
			LOGGER.Error("Failed to get OFI account address from AWS secret manager")
			return constant.STATUS_CODE_INTERNAL_ERROR, "", ""
		}
		ofiAccount = account.NodeAddress
		sendingAccount = ofiAccount
	case constant.PACS004:
		// if RFI receive a pacs004 message, the transaction sender will be RFI and receiver will be OFI and will use RFI settlement account to sign the transaction
		sendingAccount = rfiAccount
		settlementAccountName = rfiSettlementAccountName
		account := client.GetParticipantAccount(op.prServiceURL, dbData.OFIId, dbData.SettlementAccountName)
		if account == nil {
			LOGGER.Error("Failed to get OFI account address from PR")
			return constant.STATUS_CODE_INTERNAL_ERROR, "", ""
		}
		receivingAccount = *account
	case constant.PACS002:
		account, getErr := participant.GenericGetAccount(utils.Session{}, dbData.SettlementAccountName)
		if getErr != nil {
			LOGGER.Error("Failed to get OFI account address from AWS secret manager")
			return constant.STATUS_CODE_INTERNAL_ERROR, "", ""
		}
		sendingAccount = account.NodeAddress
		receivingAccount = rfiSettlementAccountName
		settlementAccountName = dbData.SettlementAccountName
	}

	LOGGER.Infof("Get IBM account and sequence number from the gas service")
	ibmAccount, seqNum, gasErr := op.GasClient.GetAccountAndSequence()
	if gasErr != nil || ibmAccount == "" {
		LOGGER.Errorf("IBM account: %s, Seq: %d", ibmAccount, seqNum)
		if gasErr != nil {
			LOGGER.Errorf("Failed to get IBM account and tx sequence: %s", gasErr.Error())
		}
		return constant.STATUS_CODE_INTERNAL_ERROR, "", ""
	}

	LOGGER.Infof("Create Stellar transaction")
	signedTx, txErr := op.createStellarTransaction(ibmAccount, sendingAccount, receivingAccount, settlementAccountName, dbData, seqNum, memoHash)
	if txErr != nil {
		LOGGER.Errorf("Failed to create Stellar Transaction: %s", txErr.Error())
		return constant.STATUS_CODE_INTERNAL_ERROR, "", ""
	}

	LOGGER.Infof("Submit Stellar transaction")
	txHash, submitErr := op.submitToStellar(reqMsgId, msgId, signedTx)
	if submitErr != nil || txHash == "" {
		LOGGER.Errorf("Failed to submit Transaction to Stellar: %s", submitErr.Error())
		return constant.STATUS_CODE_INTERNAL_ERROR, "", ""
	}

	LOGGER.Infof("Successfully submit transaction to Stellar network.")

	return constant.STATUS_CODE_TX_SEND_TO_STELLAR, txHash, ofiAccount
}

//func (op *CreateFundingOpereations) RecordTimeLogsToKafka(action, task string, timeStamp time.Time) {
//	logs := fmt.Sprintf("%s-[%s]:[%s][%s][Time:%s]", action, op.homeDomain, op.serviceName, task, time.Since(timeStamp).String())
//	postgres.CreateServiceLogInDB(op.serviceName, logs)
//}
