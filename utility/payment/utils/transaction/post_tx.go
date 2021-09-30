// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package transaction

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/GFTN/gftn-services/utility/payment/constant"

	"github.com/GFTN/gftn-services/gftn-models/model"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/wwfirebase"
)

// FbTrxLog : Logs Non-PII transaction data for the portal
type FbTrxLog struct {
	// Don't log "TransferRequest" remove this as it may include PII
	// The only reason we might want to consider keeping
	// this info is for addresses for GEO locations
	//  and payout point for tracking endpoint financial institutions
	// TransferRequest interface{}

	// HomeDomain
	ParticipantID *string `json:"participant_id"`

	// Transaction details included in transaction memo
	// debtor/creditor amounts, currencies, status
	TransactionMemo interface{} `json:"transaction_memo"`
}

type FbTrxUpdateLog struct {
	ParticipantID   string                 `json:"participant_id"`
	TransactionMemo map[string]interface{} `json:"transaction_memo"`
}

// SendToAdm - logs transaction in admin service, firebase, etc.
func (op *CreateFundingOpereations) SendToAdm(paymentInfo []model.TransactionReceipt, method, instructionId string, txMemo model.FitoFICCTMemoData) {
	timeStamp := time.Now().Unix()
	LOGGER.Debug("Store transactionMemo to WW admin-service and FireBase DB")

	op.sendTransactionToWWAdministrator(paymentInfo, timeStamp, method, instructionId, txMemo)
}

func (op *CreateFundingOpereations) sendTransactionToWWAdministrator(paymentInfo []model.TransactionReceipt, timeStamp int64, method, instructionId string, txMemo model.FitoFICCTMemoData) {

	LOGGER.Infof("Sending transaction to admin service")

	newMemoData := txMemo
	newMemoData.TimeStamp = &timeStamp

	for _, p := range paymentInfo {
		pByte, _ := json.Marshal(p)

		var pi *model.TransactionReceipt
		json.Unmarshal(pByte, &pi)
		newMemoData.TransactionStatus = append(newMemoData.TransactionStatus, pi)
	}

	// logs send transactions to firebase to be displayed in the portal
	// so that the results can be filtered by environment and participantId
	// for a given institution on the ww network
	release := os.Getenv(global_environment.ENV_KEY_SERVICE_VERSION)

	switch method {
	case constant.FIREBASE_INIT:

		ofiId := *txMemo.Fitoficctnonpiidata.Transactiondetails.OfiID
		rfiId := *txMemo.Fitoficctnonpiidata.Transactiondetails.RfiID
		LOGGER.Infof("Adding Firebase record for both OFI: %v and RFI: %v", ofiId, rfiId)
		participantIds := []string{ofiId, rfiId}

		for _, pID := range participantIds {
			err := wwfirebase.FbRef.Child(release+"/txn/transfer/"+pID+"/"+*txMemo.Fitoficctnonpiidata.InstructionID).Set(wwfirebase.AppContext, FbTrxLog{
				ParticipantID:   &pID,
				TransactionMemo: newMemoData,
			})
			if err != nil {
				LOGGER.Error("Error posting log to Firebase for %s: %s", pID, err.Error())
			}
		}

		// write to mongo DB
		op.admClient.StoreFITOFICCTMemo(newMemoData)

		if newMemoData.MessageName != nil {
			switch *newMemoData.MessageName {
			case constant.PACS004:
				if *newMemoData.MessageType == constant.PAYMENT_TYPE_CANCELLATION {
					updateOriginalTxDetails(newMemoData, release, constant.CAMT056)
				} else if *newMemoData.MessageType == constant.PAYMENT_TYPE_RDO {
					updateOriginalTxDetails(newMemoData, release, constant.IBWF002)
				}
			}
		}

		return
	case constant.FIREBASE_UPDATE_BOTH:
		var participantIds []string

		if txMemo.Fitoficctnonpiidata.Transactiondetails.OfiID != nil {
			participantIds = append(participantIds, *txMemo.Fitoficctnonpiidata.Transactiondetails.OfiID)
		}
		if txMemo.Fitoficctnonpiidata.Transactiondetails.RfiID != nil {
			participantIds = append(participantIds, *txMemo.Fitoficctnonpiidata.Transactiondetails.RfiID)
		}
		LOGGER.Infof("Updating Firebase record for participants %v", participantIds)
		for _, pID := range participantIds {
			var log interface{}

			LOGGER.Infof("Updating result to FireBase: participant: %s, instruction id: %s", pID, instructionId)
			// Get the tx record based on the instruction id
			ref := wwfirebase.FbClient.NewRef("/" + release + "/txn/transfer/" + pID + "/" + instructionId + "/transaction_memo")
			ref.Get(wwfirebase.AppContext, &log)
			if log == nil {
				LOGGER.Error("Unable to find instruction id %s for participant: %s", instructionId, pID)
				return
			}

			var updatedLog = make(map[string]interface{})

			byteData, _ := json.Marshal(log)
			var oldMemoData *model.FitoFICCTMemoData
			json.Unmarshal(byteData, &oldMemoData)

			if newMemoData.MessageName != nil {
				switch *newMemoData.MessageName {
				case constant.CAMT029:
					newMemoData, _ = getTxDetails(newMemoData, release, constant.CAMT056)
				case constant.IBWF001:
					newMemoData, _ = getTxDetails(newMemoData, release, constant.PACS008)
				}
			}

			// hash
			newHash := oldMemoData.TransactionIdentifier
			if len(newMemoData.TransactionIdentifier) > 0 {
				newHash = append(newHash, newMemoData.TransactionIdentifier[0])
			}

			updatedLog["transaction_identifier"] = newHash
			updatedLog["fitoficctnonpiidata"] = newMemoData.Fitoficctnonpiidata
			updatedLog["transaction_status"] = newMemoData.TransactionStatus

			// Update the log into FireBase
			err := ref.Update(wwfirebase.AppContext, updatedLog)
			if err != nil {
				LOGGER.Error("Error posting log to Firebase for participant %s: %s", pID, err.Error())
			}

			// write to mongo DB
			op.admClient.StoreFITOFICCTMemo(newMemoData)

		}

		return
	case constant.FIREBASE_UPDATE_PARTIAL:
		ofiId := *txMemo.Fitoficctnonpiidata.Transactiondetails.OfiID
		rfiId := *txMemo.Fitoficctnonpiidata.Transactiondetails.RfiID
		LOGGER.Infof("Updating Firebase record for both OFI: %v and RFI: %v", ofiId, rfiId)
		participantIds := []string{ofiId, rfiId}

		for _, pID := range participantIds {
			var log interface{}

			LOGGER.Infof("Updating result to FireBase: participant: %s, instruction id: %s", pID, instructionId)
			ref := wwfirebase.FbClient.NewRef("/" + release + "/txn/transfer/" + pID + "/" + instructionId + "/transaction_memo")
			ref.Get(wwfirebase.AppContext, &log)

			if log == nil {
				LOGGER.Error("Unable to find instruction id %s for participant: %s", instructionId, pID)
				return
			}
			updatedLog := make(map[string]interface{})

			byteData, _ := json.Marshal(log)
			var oldMemoData *model.FitoFICCTMemoData
			json.Unmarshal(byteData, &oldMemoData)

			newHash := oldMemoData.TransactionIdentifier
			if len(newMemoData.TransactionIdentifier) > 0 {
				newHash = append(newHash, newMemoData.TransactionIdentifier[0])
			}

			updatedLog["transaction_status"] = newMemoData.TransactionStatus
			updatedLog["transaction_identifier"] = newHash

			// Update the log into FireBase
			err := ref.Update(wwfirebase.AppContext, updatedLog)
			if err != nil {
				LOGGER.Error("Error posting log to Firebase for participant %s: %s", pID, err.Error())
			}

			// write to mongo DB
			op.admClient.StoreFITOFICCTMemo(*oldMemoData)

		}

		return

	}
}

func updateOriginalTxDetails(txMemo model.FitoFICCTMemoData, release, originalMsgName string) error {
	ofiId := *txMemo.Fitoficctnonpiidata.Transactiondetails.OfiID
	rfiId := *txMemo.Fitoficctnonpiidata.Transactiondetails.RfiID
	instructionId := *txMemo.Fitoficctnonpiidata.OriginalInstructionID
	LOGGER.Infof("Updating Firebase record %v for both OFI: %v and RFI: %v", instructionId, ofiId, rfiId)
	participantIds := []string{ofiId, rfiId}

	for _, pID := range participantIds {
		var log interface{}

		LOGGER.Infof("Update result to FireBase for participant: %s", pID)
		// Get all the txn logs from FireBase
		ref := wwfirebase.FbClient.NewRef("/" + release + "/txn/transfer/" + pID + "/" + instructionId + "/transaction_memo")
		ref.Get(wwfirebase.AppContext, &log)
		if log == nil {
			LOGGER.Error("Unable to find instruction id %s for participant: %s", instructionId, pID)
			return errors.New("Unable to find original instruction id: " + instructionId)
		}

		updatedLog := make(map[string]interface{})
		updatedMemoMap := make(map[string]interface{})

		byteData, _ := json.Marshal(log)
		var oldMemoData *model.FitoFICCTMemoData
		json.Unmarshal(byteData, &oldMemoData)

		oldMemoData.Fitoficctnonpiidata.AccountNameSend = txMemo.Fitoficctnonpiidata.AccountNameSend
		oldMemoData.Fitoficctnonpiidata.ExchangeRate = txMemo.Fitoficctnonpiidata.ExchangeRate
		oldMemoData.Fitoficctnonpiidata.Transactiondetails = txMemo.Fitoficctnonpiidata.Transactiondetails
		byteMemo, _ := json.Marshal(oldMemoData)
		json.Unmarshal(byteMemo, &updatedMemoMap)

		updatedLog["fitoficctnonpiidata"] = txMemo.Fitoficctnonpiidata

		// Update the log into FireBase
		err := ref.Update(wwfirebase.AppContext, updatedLog)
		if err != nil {
			LOGGER.Error("Error posting log to Firebase for participant %s: %s", pID, err.Error())
		}

	}
	return nil
}

func getTxDetails(txMemo model.FitoFICCTMemoData, release, originalMsgName string) (model.FitoFICCTMemoData, error) {
	ofiId := *txMemo.Fitoficctnonpiidata.Transactiondetails.OfiID
	originalInstructionId := *txMemo.Fitoficctnonpiidata.OriginalInstructionID

	var log interface{}

	LOGGER.Infof("Update result to FireBase for participant: %s", ofiId)
	ref := wwfirebase.FbClient.NewRef("/" + release + "/txn/transfer/" + ofiId + "/" + originalInstructionId + "/transaction_memo/fitoficctnonpiidata")
	ref.Get(wwfirebase.AppContext, &log)
	if log == nil {
		LOGGER.Error("Unable to find instruction id %s for participant: %s", originalInstructionId, ofiId)
		return model.FitoFICCTMemoData{}, errors.New("Unable to find instruction id: " + originalInstructionId)
	}

	e := log.(map[string]interface{})
	byteData, _ := json.Marshal(e)
	var oldMemoData *model.FitoFICCTNonPiiData
	json.Unmarshal(byteData, &oldMemoData)

	txMemo.Fitoficctnonpiidata.AccountNameSend = oldMemoData.AccountNameSend
	txMemo.Fitoficctnonpiidata.CreditorPaymentAddress = oldMemoData.CreditorPaymentAddress
	txMemo.Fitoficctnonpiidata.EndToEndID = oldMemoData.EndToEndID
	txMemo.Fitoficctnonpiidata.ExchangeRate = oldMemoData.ExchangeRate
	txMemo.Fitoficctnonpiidata.Transactiondetails = oldMemoData.Transactiondetails

	return txMemo, nil
}
