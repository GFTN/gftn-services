// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package message_handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stellar/go/xdr"
	camt029struct "github.com/GFTN/iso20022/camt02900109"
	camt056struct "github.com/GFTN/iso20022/camt05600108"
	pacs002struct "github.com/GFTN/iso20022/pacs00200109"
	pacs004struct "github.com/GFTN/iso20022/pacs00400109"

	camt029pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/camt02900109"
	camt056pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/camt05600108"
	blocklist_client "github.com/GFTN/gftn-services/administration-service/blocklist-client"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/utility/payment/constant"
	"github.com/GFTN/gftn-services/utility/payment/environment"
	message_converter "github.com/GFTN/gftn-services/utility/payment/message-converter"
	"github.com/GFTN/gftn-services/utility/payment/utils"
	"github.com/GFTN/gftn-services/utility/payment/utils/database"
	"github.com/GFTN/gftn-services/utility/payment/utils/horizon"
	"github.com/GFTN/gftn-services/utility/xmldsig"

	"github.com/GFTN/gftn-services/utility/payment/utils/parse"
	"github.com/GFTN/gftn-services/utility/payment/utils/transaction"

	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"

	"github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/kafka"
	whitelist_handler "github.com/GFTN/gftn-services/utility/payment/utils/whitelist-handler"
)

/*
 	camt.056.001.08 FIToFIPaymentCancellationRequest
 	The FIToFIPaymentCancellationRequest message supports both the request for cancellation (the
	instructed agent - or assignee - has not yet processed and forwarded the payment instruction) as well
	as the request for refund (payment has been fully processed already by the instructed agent - or
	assignee).

	A FIToFIPaymentCancellationRequest message concerns one and only one original payment
	instruction at a time.

	[Mandatory fields]

	Assignment
	- Identifies the assignment of an investigation case from an assigner to an assignee.
	Usage: The assigner must be the sender of this confirmation and the assignee must be the
	receiver.

	Underlying
	- Identifies the payment instruction to be cancelled.
*/
func (op *PaymentOperations) Camt056(camt056 message_converter.Camt056) ([]byte, error) {
	// Validate content inside the camt056 message
	structData := camt056.Message
	xmlMsgType := constant.CAMT056
	xmlData, statusData, getInfoErr := getCriticalInfoFromCAMT056(structData.Body, op.homeDomain)
	originalReqMsgId := xmlData.OriginalMsgId
	errCode := xmlData.ErrorCode
	BIC := os.Getenv(environment.ENV_KEY_PARTICIPANT_BIC)
	ofiId := xmlData.OFIId
	rfiId := xmlData.RFIId

	originalGrpInf := &pacs002struct.OriginalGroupInformation29{
		OrgnlMsgId:   getReportMax35Text(originalReqMsgId),
		OrgnlMsgNmId: getReportMax35Text(xmlMsgType),
	}

	target, _, err := parse.KafkaErrorRouter(xmlMsgType, xmlData.MessageId, ofiId, rfiId, 0, false, originalGrpInf)

	if getInfoErr != nil {
		LOGGER.Error(getInfoErr.Error())
		report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, errCode)
		return report, getInfoErr
	}

	/*
		blocklist check
	*/

	// Create admin-service client for connecting to admin-service
	blockListClient := blocklist_client.Client{
		HTTPClient: &http.Client{Timeout: time.Second * 80},
		AdminUrl:   os.Getenv(global_environment.ENV_KEY_ADMIN_SVC_URL),
	}

	// preparing the data that need to be verify against the block-list
	var countries []string
	countries = append(countries, strconv.Quote(xmlData.OfiCountry))
	countries = append(countries, strconv.Quote(xmlData.RfiCountry))

	var currencies []string
	currencies = append(currencies, strconv.Quote(xmlData.CurrencyCode))

	var participants []string
	participants = append(participants, strconv.Quote(xmlData.OFIId))
	participants = append(participants, strconv.Quote(xmlData.RFIId))

	// validate block-list
	res, err := blockListClient.ValidateFromBlocklist(countries, currencies, participants)
	if err != nil {
		LOGGER.Errorf("%v", err)
		report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, constant.STATUS_CODE_INTERNAL_ERROR)
		return report, err
	}
	if res == common.BlocklistDeniedString {
		LOGGER.Errorf("The transaction currency/country/institution is within the blocklist, transaction forbidden!")
		report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, constant.STATUS_CODE_BLOCKLIST)
		return report, errors.New("the transaction currency/country/institution is within the blocklist, transaction forbidden")
	}

	camt056LogHandler := transaction.InitiatePaymentLogOperation()
	camt056LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_INITIAL)
	msgType := constant.PAYMENT_TYPE_CANCELLATION
	msgName := constant.CAMT056

	err = database.SyncWithDynamo(constant.DATABASE_INIT, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_PENDING, constant.DATABASE_STATUS_NONE, camt056LogHandler)
	if err != nil {
		LOGGER.Errorf(err.Error())
		report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, constant.STATUS_CODE_DUP_ID)
		return report, err
	}

	// Check mutual whitelist
	LOGGER.Infof("Check whether RFI is in OFI's whitelist and vice versa")
	pKey, whiteListErr := op.whitelistHandler.CheckWhiteListParticipant(xmlData.OFIId, xmlData.RFIId, constant.EMPTY_STRING)
	if whiteListErr != nil {
		LOGGER.Errorf(whiteListErr.Error())
		camt056LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, camt056LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", camt056LogHandler, &op.fundHandler, statusData)
		report := parse.CreateCamt030(BIC, originalReqMsgId, xmlMsgType, target, constant.STATUS_CODE_INTERNAL_ERROR)
		return report, whiteListErr
	}

	if pKey == "" {
		errMsg := "OFI can not find RFI in whitelist and vice versa"
		LOGGER.Errorf(errMsg)
		camt056LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, camt056LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", camt056LogHandler, &op.fundHandler, statusData)
		report := parse.CreateCamt030(BIC, originalReqMsgId, xmlMsgType, target, constant.STATUS_CODE_RFI_OR_OFI_NOT_IN_WL)
		return report, whiteListErr
	}
	LOGGER.Infof("Yes, RFI is in OFI's whitelist and vice versa")

	// Parse the pacs008 message with signature into ProtoBuffer
	LOGGER.Infof("Start parsing Go struct to ProtoBuffer")
	protoBufData, parseErr := proto.Marshal(&camt056.SendPayload)
	if parseErr != nil {
		errMsg := "Parse data to ProtoBuf error: " + parseErr.Error()
		LOGGER.Errorf(errMsg)
		camt056LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, camt056LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", camt056LogHandler, &op.fundHandler, statusData)
		report := parse.CreateCamt030(BIC, originalReqMsgId, xmlMsgType, target, constant.STATUS_CODE_INTERNAL_ERROR)
		return report, parseErr
	}
	LOGGER.Infof("Finished parsing Go struct to ProtoBuffer")

	camt056LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_SUCCESS)
	//save the instruction id of camt056 for pacs004/camt029 msg to use
	dbData := sendmodel.DBData{
		MessageId: string(*structData.Body.Assgnmt.Id),
	}

	dbDataByte, _ := json.Marshal(dbData)
	base64DBData := parse.EncodeBase64(dbDataByte)

	// Add the transaction status into the FireBase
	database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, base64DBData, constant.DATABASE_STATUS_PENDING, constant.DATABASE_STATUS_NONE, camt056LogHandler)
	database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", camt056LogHandler, &op.fundHandler, statusData)

	// Send the ProtoBuffer to the request topic of RFI on Kafka broker
	LOGGER.Infof("Start to send request to Kafka broker")
	kafkaErr := op.KafkaActor.Produce(xmlData.RFIId+kafka.REQUEST_TOPIC, protoBufData)
	if kafkaErr != nil {
		errMsg := "Error while submit message to Kafka broker: " + kafkaErr.Error()
		LOGGER.Errorf(errMsg)
		camt056LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, camt056LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", camt056LogHandler, &op.fundHandler, statusData)
		report := parse.CreateCamt030(BIC, originalReqMsgId, xmlMsgType, target, constant.STATUS_CODE_INTERNAL_ERROR)
		return report, kafkaErr
	}
	LOGGER.Infof("Successfully produce message to Kafka broker")
	LOGGER.Debug("-----------------------------------------------------------------")

	// Send status back to OFI
	report := parse.CreateCamt030(BIC, originalReqMsgId, xmlMsgType, target, constant.STATUS_CODE_OFI_SEND_TO_KAFKA)

	return report, nil
}

/*
	pacs.004.001.08 PaymentReturn
	The PaymentReturn message is sent by an agent to the previous agent in the payment chain to undo a
	payment previously settled.

	The PaymentReturn message is exchanged between agents to return funds after settlement of credit
	transfer instructions (i.e. FIToFICustomerCreditTransfer message and FinancialInstitutionCreditTransfer
	message) or direct debit instructions (FIToFICustomerDirectDebit message).
*/
func (op *PaymentOperations) Pacs004_Cancellation(pacs004 message_converter.Pacs004) ([]byte, error) {
	// Validate content inside the pacs004 message
	structData := pacs004.Message
	xmlMsgType := constant.PACS004
	BIC := os.Getenv(environment.ENV_KEY_PARTICIPANT_BIC)

	xmlData, statusData, pacs008PaymentInfo, getDataErr := getCriticalInfoFromPacs004(structData.Body, op.homeDomain)
	errCode := xmlData.ErrorCode
	ofiId := xmlData.OFIId
	rfiId := xmlData.RFIId
	msgId := xmlData.MessageId

	originalGrpInf := &pacs002struct.OriginalGroupInformation29{
		OrgnlMsgId:   getReportMax35Text(msgId),
		OrgnlMsgNmId: getReportMax35Text(xmlMsgType),
	}

	target, _, err := parse.KafkaErrorRouter(xmlMsgType, msgId, ofiId, rfiId, 0, false, originalGrpInf)

	if getDataErr != nil {
		LOGGER.Error(getDataErr.Error())
		report := parse.CreatePacs002(BIC, xmlData.InstructionId, target, errCode, originalGrpInf)
		return report, getDataErr
	}

	/*
		blocklist check
	*/

	// Create admin-service client for connecting to admin-service
	blockListClient := blocklist_client.Client{
		HTTPClient: &http.Client{Timeout: time.Second * 80},
		AdminUrl:   os.Getenv(global_environment.ENV_KEY_ADMIN_SVC_URL),
	}

	// preparing the data that need to be verify against the block-list
	var countries []string
	countries = append(countries, strconv.Quote(xmlData.OfiCountry))
	countries = append(countries, strconv.Quote(xmlData.RfiCountry))

	var currencies []string
	currencies = append(currencies, strconv.Quote(xmlData.CurrencyCode))

	var participants []string
	participants = append(participants, strconv.Quote(xmlData.OFIId))
	participants = append(participants, strconv.Quote(xmlData.RFIId))

	// validate block-list
	res, err := blockListClient.ValidateFromBlocklist(countries, currencies, participants)
	if err != nil {
		LOGGER.Errorf("%v", err)
		report := parse.CreatePacs002(BIC, xmlData.InstructionId, target, constant.STATUS_CODE_INTERNAL_ERROR, originalGrpInf)
		return report, err
	}
	if res == common.BlocklistDeniedString {
		LOGGER.Errorf("The transaction currency/country/institution is within the blocklist, transaction forbidden!")
		report := parse.CreatePacs002(BIC, xmlData.InstructionId, target, constant.STATUS_CODE_BLOCKLIST, originalGrpInf)
		return report, errors.New("the transaction currency/country/institution is within the blocklist, transaction forbidden")
	}

	rfiAccountName := xmlData.RFISettlementAccountName

	pacs008LogHandler := transaction.InitiatePaymentLogOperation()
	pacs004LogHandler := transaction.InitiatePaymentLogOperation()
	pacs008LogHandler.PaymentStatuses = pacs008PaymentInfo
	// Store the transaction information into the administration service
	msgType := constant.PAYMENT_TYPE_CANCELLATION
	msgName := constant.PACS004
	pacs004LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_INITIAL)

	err = database.SyncWithDynamo(constant.DATABASE_INIT, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_PENDING, constant.DATABASE_STATUS_NONE, pacs004LogHandler)
	if err != nil {
		LOGGER.Errorf(err.Error())
		report := parse.CreatePacs002(BIC, xmlData.InstructionId, target, constant.STATUS_CODE_DUP_ID, originalGrpInf)
		return report, err
	}

	LOGGER.Infof("Check whether OFI is in RFI's whitelist and vice versa")
	rfiAccount, whiteListErr := op.whitelistHandler.CheckWhiteListParticipant(ofiId, rfiId, rfiAccountName)
	if whiteListErr != nil {
		LOGGER.Error(whiteListErr.Error())
		report := parse.CreatePacs002(BIC, xmlData.InstructionId, target, constant.STATUS_CODE_INTERNAL_ERROR, originalGrpInf)
		pacs004LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, pacs004LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", pacs004LogHandler, &op.fundHandler, statusData)
		return report, whiteListErr
	}
	if rfiAccount == "" {
		errMsg := "RFI can not find OFI in whitelist and vice versa"
		LOGGER.Error(errMsg)
		report := parse.CreatePacs002(BIC, xmlData.InstructionId, target, constant.STATUS_CODE_RFI_OR_OFI_NOT_IN_WL, originalGrpInf)
		pacs004LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, pacs004LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", pacs004LogHandler, &op.fundHandler, statusData)
		return report, whiteListErr
	}
	LOGGER.Infof("Yes, OFI is in RFI's whitelist and vice versa")

	pacs004LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_SUCCESS)

	signData := &sendmodel.SignData{
		OFIId: ofiId,
		// get OFI settlement account name from original transaction information
		SettlementAccountName: xmlData.OFISettlementAccountName,
		// use return interbank settlement amount as settlement amount
		SettlementAmount: structData.Body.TxInf[0].RtrdIntrBkSttlmAmt.Value,
		CurrencyCode:     structData.Body.TxInf[0].RtrdIntrBkSttlmAmt.Currency,
		// get asset issuer ID from original transaction information
		AssetIssuerId: string(*structData.Body.TxInf[0].OrgnlTxRef.PmtTpInf.SvcLvl[0].Prtry),
	}

	submitResult, txHash, ofiAccount := op.fundHandler.FundAndSubmitPaymentTransaction(rfiAccount, xmlData.InstructionId, msgId, xmlMsgType, rfiAccountName, *signData, xdr.Memo{})
	report := parse.CreatePacs002(BIC, xmlData.InstructionId, target, submitResult, originalGrpInf)

	if submitResult != constant.STATUS_CODE_TX_SEND_TO_STELLAR {
		pacs004LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		pacs008LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, pacs004LogHandler)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.OriginalInstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_DONE, constant.DATABASE_STATUS_NONE, pacs008LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", pacs004LogHandler, &op.fundHandler, statusData)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, constant.PACS008, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.OriginalInstructionId, "", "", pacs008LogHandler, &op.fundHandler, statusData)
		return report, nil
	} else {
		if string(*structData.Body.TxInf[0].OrgnlTxRef.SttlmInf.SttlmMtd) == constant.DO_SETTLEMENT {
			// record the payment status "cleared"
			pacs004LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RETURNED, txHash)
			pacs008LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RETURNED, txHash)
		} else {
			// record the payment status "settled"
			pacs004LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_SETTLED, txHash)
			pacs008LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_SETTLED, txHash)
		}
	}

	// update status in firebase

	go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.OriginalInstructionId, txHash, constant.DATABASE_STATUS_CANCELED, constant.DATABASE_STATUS_NONE, pacs008LogHandler)
	go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, txHash, constant.DATABASE_STATUS_DONE, constant.DATABASE_STATUS_NONE, pacs004LogHandler)
	go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, txHash, ofiAccount, pacs004LogHandler, &op.fundHandler, statusData)
	go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, msgType, constant.PACS008, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.OriginalInstructionId, txHash, "", pacs008LogHandler, &op.fundHandler, statusData)

	LOGGER.Debug("---------------------------------------------------------------------")

	return report, nil
}

func (op *PaymentOperations) Camt029(camt029 message_converter.Camt029) ([]byte, error) {
	// Validate content inside the camt029 message
	structData := camt029.Message
	xmlMsgType := constant.CAMT029
	BIC := os.Getenv(environment.ENV_KEY_PARTICIPANT_BIC)
	xmlData, statusData, getDataErr := getCriticalInfoFromCamt029(structData.Body, op.homeDomain)
	originalReqMsgId := xmlData.OriginalMsgId
	errCode := xmlData.ErrorCode
	rfiId := xmlData.RFIId
	ofiId := xmlData.OFIId
	msgId := xmlData.MessageId

	originalGrpInf := &pacs002struct.OriginalGroupInformation29{
		OrgnlMsgId:   getReportMax35Text(msgId),
		OrgnlMsgNmId: getReportMax35Text(xmlMsgType),
	}

	target, _, err := parse.KafkaErrorRouter(xmlMsgType, msgId, ofiId, rfiId, 0, false, originalGrpInf)

	if getDataErr != nil {
		LOGGER.Error(getDataErr.Error())
		report := parse.CreateCamt030(BIC, originalReqMsgId, xmlMsgType, target, errCode)
		return report, getDataErr
	}

	/*
		blocklist check
	*/

	// Create admin-service client for connecting to admin-service
	blockListClient := blocklist_client.Client{
		HTTPClient: &http.Client{Timeout: time.Second * 80},
		AdminUrl:   os.Getenv(global_environment.ENV_KEY_ADMIN_SVC_URL),
	}

	// preparing the data that need to be verify against the block-list
	var countries []string
	countries = append(countries, strconv.Quote(xmlData.OfiCountry))
	countries = append(countries, strconv.Quote(xmlData.RfiCountry))

	var currencies []string
	currencies = append(currencies, strconv.Quote(xmlData.CurrencyCode))

	var participants []string
	participants = append(participants, strconv.Quote(xmlData.OFIId))
	participants = append(participants, strconv.Quote(xmlData.RFIId))

	// validate block-list
	res, err := blockListClient.ValidateFromBlocklist(countries, currencies, participants)
	if err != nil {
		LOGGER.Errorf("%v", err)
		report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, constant.STATUS_CODE_INTERNAL_ERROR)
		return report, err
	}
	if res == common.BlocklistDeniedString {
		LOGGER.Errorf("The transaction currency/country/institution is within the blocklist, transaction forbidden!")
		report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, constant.STATUS_CODE_BLOCKLIST)
		return report, errors.New("the transaction currency/country/institution is within the blocklist, transaction forbidden")
	}

	camt029LogHandler := transaction.InitiatePaymentLogOperation()
	// Store the transaction information into the administration service
	msgType := constant.PAYMENT_TYPE_CANCELLATION
	msgName := constant.CAMT029
	camt029LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_INITIAL)
	err = database.SyncWithDynamo(constant.DATABASE_INIT, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_PENDING, constant.DATABASE_STATUS_NONE, camt029LogHandler)
	if err != nil {
		LOGGER.Errorf(err.Error())
		report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, constant.STATUS_CODE_DUP_ID)
		return report, err
	}

	LOGGER.Infof("Check whether OFI is in RFI's whitelist and vice versa")
	pKey, whiteListErr := op.whitelistHandler.CheckWhiteListParticipant(rfiId, ofiId, constant.EMPTY_STRING)
	if whiteListErr != nil {
		LOGGER.Errorf(whiteListErr.Error())
		camt029LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, camt029LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", camt029LogHandler, &op.fundHandler, statusData)
		report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, constant.STATUS_CODE_INTERNAL_ERROR)
		return report, whiteListErr
	}
	if pKey == "" {
		LOGGER.Errorf("RFI can not find OFI in whitelist and vice versa")
		camt029LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, camt029LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", camt029LogHandler, &op.fundHandler, statusData)
		report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, constant.STATUS_CODE_RFI_OR_OFI_NOT_IN_WL)
		return report, whiteListErr
	}
	LOGGER.Infof("Yes, OFI is in RFI's whitelist and vice versa")

	// Parse the camt029 message with signature into ProtoBuffer
	LOGGER.Infof("Start parsing Go struct to ProtoBuffer")

	protoBufData, parseErr := proto.Marshal(&camt029.SendPayload)
	if parseErr != nil {
		LOGGER.Errorf("Parse struct to ProtoBuf error %s", parseErr.Error())
		camt029LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_DONE, constant.DATABASE_STATUS_NONE, camt029LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", camt029LogHandler, &op.fundHandler, statusData)
		report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, constant.STATUS_CODE_INTERNAL_ERROR)
		return report, parseErr
	}
	LOGGER.Infof("Finished parsing Go struct to ProtoBuffer")

	camt029LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_SUCCESS)
	database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_PENDING, constant.DATABASE_STATUS_NONE, camt029LogHandler)
	database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", camt029LogHandler, &op.fundHandler, statusData)

	// Send the ProtoBuffer to the response topic of OFI on Kafka broker
	LOGGER.Infof("Start to send response to Kafka broker")
	kafkaErr := op.KafkaActor.Produce(ofiId+kafka.RESPONSE_TOPIC, protoBufData)
	if kafkaErr != nil {
		LOGGER.Errorf("Error while submit message to Kafka broker: %s", kafkaErr.Error())
		camt029LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_DONE, constant.DATABASE_STATUS_NONE, camt029LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", camt029LogHandler, &op.fundHandler, statusData)
		report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, constant.STATUS_CODE_INTERNAL_ERROR)
		return report, kafkaErr
	}
	LOGGER.Infof("Successfully produce message to Kafka broker")
	LOGGER.Debug("-----------------------------------------------------------")

	// Send status back to RFI
	report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, constant.STATUS_CODE_CXL_RES_SEND_TO_KAFKA)

	return report, nil
}

// if message type is camt.056
func RFI_Camt056(data camt056pbstruct.SendPayload, op *kafka.KafkaOpreations) {
	// Parse the ProtoBuffer into Go struct and reconstruct it into pacs008 message
	LOGGER.Infof("Parsing ProtoBuffer to XML")
	standardType := constant.ISO20022
	paymentStatusMsgType := constant.PAYMENT_TYPE_CANCELLATION
	msgName := constant.CAMT056
	pacs008InstructionId := data.OriginalInstructionId
	instructionId := data.InstructionId
	reqMsgType := data.MsgType
	ofiId := data.OfiId
	rfiId := data.RfiId
	originalMsgId := data.OriginalMsgId
	msgId := data.MsgId

	pacs008LogHandler := transaction.InitiatePaymentLogOperation()
	camt056LogHandler := transaction.InitiatePaymentLogOperation()

	participantId := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	topicName := participantId + "_" + kafka.TRANSACTION_TOPIC

	originalGrpInf := &pacs002struct.OriginalGroupInformation29{
		OrgnlMsgId:   getReportMax35Text(msgId),
		OrgnlMsgNmId: getReportMax35Text(msgName),
	}

	/*
		find pacs008 record from DB
	*/
	pacs008DynamoData, pacs008PaymentInfo := parse.GetDBData(pacs008InstructionId)
	if pacs008DynamoData == nil || pacs008PaymentInfo == nil {
		LOGGER.Errorf("The original message ID %v does not exist in DB", pacs008InstructionId)
		op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_WRONG_ORIGINAL_ID, originalGrpInf)
		return
	}
	pacs008LogHandler.PaymentStatuses = pacs008PaymentInfo

	/*
		find camt056 record from DB
	*/
	dynamoData, paymentInfo := parse.GetDBData(instructionId)
	if dynamoData == nil || paymentInfo == nil {
		LOGGER.Errorf("The original message ID %v does not exist in DB", instructionId)
		op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_WRONG_INSTRUCTION_ID, originalGrpInf)
		return
	}
	camt056LogHandler.PaymentStatuses = paymentInfo
	camt056LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RFI_PROCESSING)

	/*
		verify signature
	*/
	LOGGER.Infof("Verifying the signature from OFI")
	result := xmldsig.VerifySignature(string(data.Message))
	if !result {
		LOGGER.Errorf("signature verification failed")
		camt056LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RFI_VALIDATION_FAIL)
		//camt056 status
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt056LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, pacs008InstructionId, instructionId, "", "", camt056LogHandler, &op.FundHandler, &sendmodel.StatusData{})
		op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_OFI_SIGNATURE_FAIL, originalGrpInf)
		return
	}
	LOGGER.Infof("OFI signature verified!")

	/*
		constructing protobuffer to go struct
	*/
	camt056 := &message_converter.Camt056{SendPayload: data}
	xmlData, err := camt056.ProtobuftoStruct()

	if xmlData == nil {
		LOGGER.Errorf("Encounter error while construncting proto buffer to go struct")
		camt056LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RFI_VALIDATION_FAIL)
		//camt056 status
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt056LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, pacs008InstructionId, instructionId, "", "", camt056LogHandler, &op.FundHandler, &sendmodel.StatusData{})
		op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR_PARSE, originalGrpInf)
		return
	} else if err != nil {
		LOGGER.Errorf("Parse request from kafka failed: %s", err.Error())
		camt056LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RFI_VALIDATION_FAIL)
		//camt056 status
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt056LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, pacs008InstructionId, instructionId, "", "", camt056LogHandler, &op.FundHandler, &sendmodel.StatusData{})
		op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR_PARSE, originalGrpInf)
		return
	}
	// Get important data from the XML data
	reqData := xmlData.RequestXMLMsg
	originalMessageId := xmlData.OriginalMsgId
	settlementAccountName := xmlData.OFISettlementAccountName
	originalInstructionId := xmlData.OriginalInstructionId

	LOGGER.Infof("Finished paring ProtoBuffer to XML")

	// Generate payment status data
	// Aggregate necessary data for transaction memo
	statusData := &sendmodel.StatusData{
		IdCdtr:                rfiId,
		IdDbtr:                ofiId,
		CurrencyCode:          camt056.Message.Body.Undrlyg[0].TxInf[0].OrgnlIntrBkSttlmAmt.Currency,
		AssetType:             string(*camt056.Message.Body.Undrlyg[0].TxInf[0].OrgnlTxRef.SttlmInf.SttlmMtd),
		CrtyCcy:               string(*camt056.Message.Body.Assgnmt.Assgne.Agt.FinInstnId.BICFI)[:3],
		AccountNameSend:       string(*camt056.Message.Body.Undrlyg[0].TxInf[0].OrgnlTxRef.SttlmInf.SttlmAcct.Nm),
		EndToEndID:            string(*camt056.Message.Body.Undrlyg[0].TxInf[0].OrgnlEndToEndId),
		InstructionID:         xmlData.InstructionId,
		OriginalInstructionID: xmlData.OriginalInstructionId,
	}

	// update pacs008 record in DB
	rfiVerifyRequestAndSendToKafka(topicName, msgId, msgName, originalMessageId, ofiId, settlementAccountName, standardType, msgName, instructionId, originalInstructionId, paymentStatusMsgType, pacs008LogHandler, reqData, statusData, pacs008DynamoData, op, originalGrpInf)

	//update camt056 record in DB
	rawMsg, _ := json.Marshal(dynamoData)
	camt056LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_CANCELLATION_INIT)
	go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, parse.EncodeBase64(rawMsg), constant.DATABASE_STATUS_DONE, constant.DATABASE_STATUS_NONE, camt056LogHandler)
	go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, paymentStatusMsgType, msgName, originalMessageId, originalInstructionId, instructionId, "", "", camt056LogHandler, &op.FundHandler, statusData)

	return
}

// if message type is camt.029
func OFI_Camt029(data camt029pbstruct.SendPayload, op *kafka.KafkaOpreations) {
	// Parse the ProtoBuffer to Go struct and reconstruct it into camt029 message
	xmlMsgType := constant.CAMT029
	standardType := constant.ISO20022
	msgType := constant.PAYMENT_TYPE_CANCELLATION
	msgName := constant.CAMT029
	camt029LogHandler := transaction.InitiatePaymentLogOperation()
	pacs008LogHandler := transaction.InitiatePaymentLogOperation()
	originalInstrId := data.OriginalInstructionId
	instrId := data.InstructionId
	ofiId := data.OfiId
	rfiId := data.RfiId
	originalReqMsgId := data.OriginalMsgId

	originalGrpInf := &pacs002struct.OriginalGroupInformation29{
		OrgnlMsgId:   getReportMax35Text(originalReqMsgId),
		OrgnlMsgNmId: getReportMax35Text(msgName),
	}

	/*
	 Get camt029 data from database
	*/
	LOGGER.Infof("Get camt029 information from database")
	_, paymentInfo := parse.GetDBData(instrId)
	if paymentInfo == nil {
		LOGGER.Error("Can not get camt.029 information from database")
		op.SendErrMsg(instrId, standardType, xmlMsgType, ofiId, rfiId, constant.STATUS_CODE_WRONG_INSTRUCTION_ID, originalGrpInf)
		return
	}

	// Initialize the payment status
	camt029LogHandler.PaymentStatuses = paymentInfo
	camt029LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_OFI_PROCESSING)

	/*
	 Get pacs008 data from database
	*/
	LOGGER.Infof("Get pacs008 information from database")
	dbData, pacs008PaymentInfo := parse.GetDBData(originalInstrId)
	if dbData == nil || pacs008PaymentInfo == nil {
		LOGGER.Error("Can not get pacs.008 information from database")
		op.SendErrMsg(instrId, standardType, xmlMsgType, ofiId, rfiId, constant.STATUS_CODE_WRONG_ORIGINAL_ID, originalGrpInf)
		return
	}
	pacs008LogHandler.PaymentStatuses = pacs008PaymentInfo

	/*
		verify signature
	*/
	LOGGER.Infof("Verifying the signature from RFI")
	result := xmldsig.VerifySignature(string(data.Message))
	if !result {
		LOGGER.Errorf("signature verification failed")
		camt029LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_OFI_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instrId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt029LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, msgType, msgName, originalReqMsgId, originalInstrId, instrId, "", "", camt029LogHandler, &op.FundHandler, &sendmodel.StatusData{})
		op.SendErrMsg(instrId, standardType, xmlMsgType, ofiId, rfiId, constant.STATUS_CODE_RFI_SIGNATURE_FAIL, originalGrpInf)
		return
	}
	LOGGER.Infof("RFI signature verified!")

	LOGGER.Infof("Parsing ProtoBuffer to XML")

	var camt029 message_converter.MessageInterface = &message_converter.Camt029{SendPayload: data}
	xmlData, err := camt029.ProtobuftoStruct()
	status := xmlData.StatusCode
	statusCode, _ := strconv.Atoi(status)

	if xmlData == nil {
		LOGGER.Errorf("Encounter error while construncting proto buffer to go struct")
		camt029LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_OFI_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instrId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt029LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, msgType, msgName, originalReqMsgId, originalInstrId, instrId, "", "", camt029LogHandler, &op.FundHandler, &sendmodel.StatusData{})
		op.SendErrMsg(instrId, standardType, xmlMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR_PARSE, originalGrpInf)
		return
	} else if err != nil || statusCode != constant.STATUS_CODE_DEFAULT {
		camt029LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_OFI_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instrId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt029LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, msgType, msgName, originalReqMsgId, originalInstrId, instrId, "", "", camt029LogHandler, &op.FundHandler, &sendmodel.StatusData{})
		op.SendErrMsg(instrId, standardType, xmlMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR_PARSE, originalGrpInf)
		return
	}
	// Get important data from the XML data
	rfiAccountName := xmlData.RFISettlementAccountName
	resData := xmlData.RequestXMLMsg
	homeDomain := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)

	LOGGER.Infof("Finished paring ProtoBuffer to XML")

	// Aggregate necessary data for transaction memo
	commonStatusData := &sendmodel.StatusData{
		IdCdtr:                xmlData.RFIId,
		IdDbtr:                xmlData.OFIId,
		EndToEndID:            xmlData.OriginalEndtoEndId,
		InstructionID:         xmlData.InstructionId,
		OriginalInstructionID: xmlData.OriginalMsgId,
	}

	/*
		blocklist check
	*/

	// Create admin-service client for connecting to admin-service
	blockListClient := blocklist_client.Client{
		HTTPClient: &http.Client{Timeout: time.Second * 80},
		AdminUrl:   os.Getenv(global_environment.ENV_KEY_ADMIN_SVC_URL),
	}

	// preparing the data that need to be verify against the block-list
	var countries []string
	countries = append(countries, strconv.Quote(xmlData.OfiCountry))
	countries = append(countries, strconv.Quote(xmlData.RfiCountry))

	var currencies []string
	currencies = append(currencies, strconv.Quote(xmlData.CurrencyCode))

	var participants []string
	participants = append(participants, strconv.Quote(xmlData.RFIId))
	participants = append(participants, strconv.Quote(xmlData.OFIId))

	// validate block-list
	res, err := blockListClient.ValidateFromBlocklist(countries, currencies, participants)
	if err != nil {
		LOGGER.Errorf("%v", err)
		op.SendErrMsg(instrId, standardType, xmlMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR, originalGrpInf)
		camt029LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instrId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt029LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, msgType, msgName, originalReqMsgId, originalInstrId, instrId, "", "", camt029LogHandler, &op.FundHandler, commonStatusData)
		return
	}
	if res == common.BlocklistDeniedString {
		LOGGER.Errorf("The transaction currency/country/institution is within the blocklist, transaction forbidden!")
		op.SendErrMsg(instrId, standardType, xmlMsgType, ofiId, rfiId, constant.STATUS_CODE_BLOCKLIST, originalGrpInf)
		camt029LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_OFI_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instrId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt029LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, msgType, msgName, originalReqMsgId, originalInstrId, instrId, "", "", camt029LogHandler, &op.FundHandler, commonStatusData)
		return
	}

	// Check if RFI was whitelisted by OFI and vice versa, if not, reject it
	whitelistHandler := whitelist_handler.CreateWhiteListServiceOperations()
	LOGGER.Infof("Check whether RFI is in OFI's whitelist and vice versa.")
	rfiAccount, whiteListErr := whitelistHandler.CheckWhiteListParticipant(homeDomain, rfiId, rfiAccountName)
	if whiteListErr != nil {
		LOGGER.Errorf(whiteListErr.Error())
		op.SendErrMsg(instrId, standardType, xmlMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR, originalGrpInf)
		camt029LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instrId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt029LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, msgType, msgName, originalReqMsgId, originalInstrId, instrId, "", "", camt029LogHandler, &op.FundHandler, commonStatusData)
		return
	}
	if rfiAccount == "" {
		LOGGER.Errorf("Can not find RFI or OFI in whitelist and vice versa")
		op.SendErrMsg(instrId, standardType, xmlMsgType, ofiId, rfiId, constant.STATUS_CODE_RFI_OR_OFI_NOT_IN_WL, originalGrpInf)
		camt029LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_OFI_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instrId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt029LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, msgType, msgName, originalReqMsgId, originalInstrId, instrId, "", "", camt029LogHandler, &op.FundHandler, commonStatusData)
		return
	}
	LOGGER.Infof("Yes, RFI is in OFI's whitelist and vice versa.")

	/*
		Signing message with IBM master account
	*/
	signedMessage, signErr := op.SignHandler.SignPayloadByMasterAccount(resData)
	if signErr != nil {
		LOGGER.Errorf("Failed to sign payload: %v", signErr.Error())
		camt029LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instrId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt029LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, msgType, msgName, originalReqMsgId, originalInstrId, instrId, "", "", camt029LogHandler, &op.FundHandler, commonStatusData)
		op.SendErrMsg(instrId, standardType, xmlMsgType, ofiId, rfiId, constant.STATUS_CODE_SIGN_PAYLOAD_FAIL, originalGrpInf)
		return
	}
	gatewayMsg := parse.EncodeBase64(signedMessage)

	callBackMsg := &model.SendPacs{
		MessageType: &xmlMsgType,
		Message:     &gatewayMsg,
	}

	topicName := ofiId + "_" + kafka.TRANSACTION_TOPIC

	// Send the encoded xml message to the callback service of OFI
	LOGGER.Infof("Send encoded message back to Kafka topic: %v", topicName)
	msg, _ := json.Marshal(callBackMsg)

	/*
		sending message to Kafka
	*/

	err = op.SendRequestToKafka(topicName, msg)
	if err != nil {
		LOGGER.Errorf("Encounter error while producing message to Kafka topic: %v", ofiId+"_"+kafka.TRANSACTION_TOPIC)
		camt029LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instrId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt029LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, msgType, msgName, originalReqMsgId, originalInstrId, instrId, "", "", camt029LogHandler, &op.FundHandler, commonStatusData)
		op.SendErrMsg(instrId, standardType, xmlMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR, originalGrpInf)
		return
	}

	camt029LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_REJECTED)
	pacs008LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_REJECTED)

	go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.OriginalInstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_CLEARED, constant.DATABASE_STATUS_NONE, pacs008LogHandler)
	go database.SyncWithDynamo(constant.DATABASE_UPDATE, instrId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_DONE, constant.DATABASE_STATUS_NONE, camt029LogHandler)
	go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, msgType, msgName, originalReqMsgId, originalInstrId, instrId, "", rfiAccount, camt029LogHandler, &op.FundHandler, commonStatusData)
	go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, msgType, constant.PACS008, originalReqMsgId, originalInstrId, originalInstrId, "", "", pacs008LogHandler, &op.FundHandler, commonStatusData)

	LOGGER.Debug("---------------------------------------------------------------------")
	return
}

func getCriticalInfoFromCAMT056(document *camt056struct.FIToFIPaymentCancellationRequestV08, homeDomain string) (*sendmodel.XMLData, *sendmodel.StatusData, error) {
	ofiId := string(*document.Assgnmt.Assgnr.Agt.FinInstnId.Othr.Id)
	rfiId := string(*document.Assgnmt.Assgne.Agt.FinInstnId.Othr.Id)
	msgId := string(*document.Assgnmt.Id)
	originalPaymentMsgId := string(*document.Undrlyg[0].OrgnlGrpInfAndCxl.OrgnlMsgId)
	originalPaymentType := string(*document.Undrlyg[0].OrgnlGrpInfAndCxl.OrgnlMsgNmId)
	ofiSettlementAccountName := string(*document.Undrlyg[0].TxInf[0].OrgnlTxRef.SttlmInf.SttlmAcct.Nm)
	caseId := string(*document.Case.Id)
	originalInstrId := string(*document.Undrlyg[0].TxInf[0].OrgnlInstrId)

	checkData := &sendmodel.XMLData{
		OriginalMsgId:            originalPaymentMsgId,
		MessageId:                msgId,
		OFIId:                    ofiId,
		RFIId:                    rfiId,
		OFISettlementAccountName: ofiSettlementAccountName,
		ErrorCode:                constant.STATUS_CODE_DEFAULT,
		InstructionId:            msgId,
		OriginalInstructionId:    originalInstrId,
	}

	if !utils.StringsEqual(ofiId, homeDomain) {
		LOGGER.Error("Instructing agent is an incorrect participant")
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_FI
		return checkData, &sendmodel.StatusData{}, errors.New("instructing agent is an incorrect participant")
	}

	// Check if the original payment type is pacs.008
	if !utils.StringsEqual(originalPaymentType, constant.PACS008) {
		LOGGER.Error("Incorrect original message name ID")
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_MSG_NAME_ID
		return checkData, &sendmodel.StatusData{}, errors.New("incorrect original message name ID")
	}

	if !utils.StringsEqual(msgId, caseId) {
		LOGGER.Error("Incompatible message ID in Assignment and Case")
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_MSG_ID
		return checkData, &sendmodel.StatusData{}, errors.New("incompatible message ID in Assignment and Case")
	}

	dbData, txStatus, _, _, dbErr := database.DC.GetTransactionData(originalInstrId)

	if dbErr != nil {
		checkData.ErrorCode = constant.STATUS_CODE_INTERNAL_ERROR
		return checkData, &sendmodel.StatusData{}, errors.New("database query error")
	}

	if *dbData == "" {
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_ORIGINAL_ID
		return checkData, &sendmodel.StatusData{}, errors.New("wrong original Instruction ID")
	}

	if *txStatus != constant.DATABASE_STATUS_CLEARED && *txStatus != constant.DATABASE_STATUS_SETTLED {
		checkData.ErrorCode = constant.STATUS_CODE_ORIGINAL_REQUEST_NOT_DONE
		return checkData, &sendmodel.StatusData{}, errors.New("original payment request is not CLEARED/SETTLED yet")
	}

	/*
		Aggregate necessary data for transaction memo
	*/
	statusData := &sendmodel.StatusData{
		IdCdtr:          rfiId,
		IdDbtr:          ofiId,
		CurrencyCode:    document.Undrlyg[0].TxInf[0].OrgnlIntrBkSttlmAmt.Currency,
		AssetType:       string(*document.Undrlyg[0].TxInf[0].OrgnlTxRef.SttlmInf.SttlmMtd),
		CrtyCcy:         string(*document.Assgnmt.Assgne.Agt.FinInstnId.BICFI)[:3],
		AccountNameSend: string(*document.Undrlyg[0].TxInf[0].OrgnlTxRef.SttlmInf.SttlmAcct.Nm),
		EndToEndID:      msgId,
		InstructionID:   msgId,
	}

	return checkData, statusData, nil
}

func getCriticalInfoFromPacs004(document *pacs004struct.PaymentReturnV09, homeDomain string) (*sendmodel.XMLData, *sendmodel.StatusData, []model.TransactionReceipt, error) {
	reqMsgId := string(*document.OrgnlGrpInf.OrgnlMsgId)
	originalInstrId := string(*document.TxInf[0].OrgnlInstrId)
	settlementMethod := string(*document.GrpHdr.SttlmInf.SttlmMtd)
	ofiId := string(*document.GrpHdr.InstdAgt.FinInstnId.Othr.Id)
	rfiId := string(*document.GrpHdr.InstgAgt.FinInstnId.Othr.Id)
	ofiSettlementAccountName := strings.ToLower(string(*document.TxInf[0].OrgnlTxRef.SttlmInf.SttlmAcct.Nm))
	rfiSettlementAccountName := strings.ToLower(string(*document.GrpHdr.SttlmInf.SttlmAcct.Nm))
	msgId := string(*document.GrpHdr.MsgId)
	assetIssuerId := string(*document.TxInf[0].OrgnlTxRef.PmtTpInf.SvcLvl[0].Prtry)
	currencyCode := document.TxInf[0].RtrdIntrBkSttlmAmt.Currency
	charges, _ := strconv.ParseFloat(document.TxInf[0].ChrgsInf[0].Amt.Value, 64)
	originalAssetType := string(*document.TxInf[0].OrgnlTxRef.SttlmInf.SttlmMtd)

	checkData := &sendmodel.XMLData{
		OriginalMsgId:            reqMsgId,
		InstructionId:            string(*document.TxInf[0].RtrId),
		OriginalInstructionId:    originalInstrId,
		OFIId:                    ofiId,
		RFIId:                    rfiId,
		OFISettlementAccountName: ofiSettlementAccountName,
		RFISettlementAccountName: rfiSettlementAccountName,
		ErrorCode:                constant.STATUS_CODE_DEFAULT,
		MessageId:                msgId,
		AssetIssuer:              assetIssuerId,
	}

	if !utils.StringsEqual(rfiId, homeDomain) {
		LOGGER.Error("Instructing agent is an incorrect participant")
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_FI
		return checkData, nil, nil, errors.New("instructing agent is an incorrect participant")
	}

	dbData, txStatus, _, paymentInfoBase64, dbErr := database.DC.GetTransactionData(originalInstrId)

	if dbErr != nil {
		checkData.ErrorCode = constant.STATUS_CODE_INTERNAL_ERROR
		return checkData, nil, nil, errors.New("database query error")
	}

	if *dbData == "" {
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_ORIGINAL_ID
		return checkData, nil, nil, errors.New("wrong original instruction ID")
	}

	info, _ := parse.DecodeBase64(*paymentInfoBase64)

	var paymentInfo []model.TransactionReceipt
	json.Unmarshal(info, &paymentInfo)

	//check if camt.026 already happened, if it is, then no charges should be included
	//check if RDO already happened, if it is, then RFI should contact counterparty for reconciliation outside World Wire
	for _, elem := range paymentInfo {
		if elem.Transactionstatus == nil {
			continue
		}
		if *elem.Transactionstatus == constant.PAYMENT_STATUS_UNABLE_TO_APPLY && charges != 0 {
			return checkData, nil, nil, errors.New("The payment is not executed by RFI, so no charges should be included")
		}

		if *elem.Transactionstatus == constant.PAYMENT_STATUS_SETTLED && originalAssetType == constant.DO_SETTLEMENT {
			return checkData, nil, nil, errors.New("Please contact counterparty for reconciliation outside World Wire")
		}
	}

	// if it is DO, check if they are using issuing account & if either OFI or RFI is the issuer
	if utils.StringsEqual(settlementMethod, constant.DO_SETTLEMENT) {
		// check if this DO was issued by either OFI or RFI
		if !utils.StringsEqual(ofiId, assetIssuerId) && !utils.StringsEqual(rfiId, assetIssuerId) {
			errMsg := "Either OFI or RFI should be the asset issuer, if settlement method is WWDO"
			LOGGER.Error(errMsg)
			checkData.ErrorCode = constant.STATUS_CODE_WRONG_ASSET_ISSUER
			return checkData, nil, nil, errors.New(errMsg)
		}

		// check if the settlement account name is "issuing"
		if !utils.StringsEqual(rfiSettlementAccountName, common.ISSUING) || !utils.StringsEqual(ofiSettlementAccountName, common.ISSUING) {
			errMsg := "The settlement method is WWDO, so the account name of both OFI & RFI should be \"issuing\""
			LOGGER.Error(errMsg)
			checkData.ErrorCode = constant.STATUS_CODE_WRONG_ACCOUNT_NAME
			return checkData, nil, nil, errors.New(errMsg)
		}

		if utils.StringsEqual(ofiId, rfiId) {
			LOGGER.Error("Internal DO transfer is not allowed")
			checkData.ErrorCode = constant.STATUS_CODE_DO_INTERNAL_TRANSFER_ERROR
			return checkData, nil, nil, errors.New("Internal DO transfer is not allowed")
		}

		if !horizon.IsIssuer(assetIssuerId, currencyCode) {
			errMsg := "The asset " + currencyCode + " is not issued by the asset issuer " + assetIssuerId
			LOGGER.Errorf(errMsg)
			checkData.ErrorCode = constant.STATUS_CODE_WRONG_ASSET_ISSUER
			return checkData, nil, nil, errors.New(errMsg)
		}

		// check if settlement asset currency code is ended with "DO"
		if !strings.HasSuffix(currencyCode, constant.SETTLEMENT_METHOD_DIGITAL_OBLIGATION) {
			errMsg := "Settlement method is DO, please use DO as settlement currency"
			LOGGER.Error(errMsg)
			checkData.ErrorCode = constant.STATUS_CODE_WRONG_SETTLEMENT_METHOD
			return checkData, nil, nil, errors.New(errMsg)
		}
	}

	if *txStatus == constant.DATABASE_STATUS_CANCEL_INIT {
		reqSettlementMethod := string(*document.TxInf[0].OrgnlTxRef.SttlmInf.SttlmMtd)
		if !utils.StringsEqual(reqSettlementMethod, settlementMethod) {
			checkData.ErrorCode = constant.STATUS_CODE_WRONG_SETTLEMENT_METHOD
			return checkData, nil, nil, errors.New("settlement method is not the same as payment request")
		}

		if utils.StringsEqual(settlementMethod, constant.DO_SETTLEMENT) {
			if !utils.StringsEqual(rfiSettlementAccountName, common.ISSUING) {
				LOGGER.Error("Account name should be \"issuing\", if settlement method is WWDO")
				checkData.ErrorCode = constant.STATUS_CODE_WRONG_ACCOUNT_NAME
				return checkData, nil, nil, errors.New("wrong account name for DO")
			}
		}
	} else {
		LOGGER.Error("Data not found in database")
		checkData.ErrorCode = constant.STATUS_CODE_ORIGINAL_REQUEST_NOT_INIT
		return checkData, nil, nil, errors.New("Data not found in database")
	}

	/*
		get status data for firebase
	*/
	amountBeneficiary, _ := strconv.ParseFloat(document.TxInf[0].RtrdInstdAmt.Value, 64)
	settlementAmount, _ := strconv.ParseFloat(document.TxInf[0].RtrdIntrBkSttlmAmt.Value, 64)

	// Aggregate necessary data for transaction memo
	statusData := &sendmodel.StatusData{
		IdCdtr:               rfiId,
		IdDbtr:               ofiId,
		CurrencyCode:         document.TxInf[0].RtrdIntrBkSttlmAmt.Currency,
		AssetType:            string(*document.GrpHdr.SttlmInf.SttlmMtd),
		AmountBeneficiary:    amountBeneficiary,
		AssetCodeBeneficiary: document.TxInf[0].RtrdInstdAmt.Currency,
		AccountNameSend:      string(*document.GrpHdr.SttlmInf.SttlmAcct.Nm),
		EndToEndID:           string(*document.TxInf[0].OrgnlEndToEndId),
		InstructionID:        string(*document.TxInf[0].RtrId),
		AmountSettlement:     settlementAmount,
		IssuerID:             string(*document.TxInf[0].OrgnlTxRef.PmtTpInf.SvcLvl[0].Prtry),
		ExchangeRate:         1.0,
	}

	return checkData, statusData, paymentInfo, nil
}

func getCriticalInfoFromCamt029(document *camt029struct.ResolutionOfInvestigationV09, homeDomain string) (*sendmodel.XMLData, *sendmodel.StatusData, error) {
	rfiId := string(*document.Assgnmt.Assgnr.Agt.FinInstnId.Othr.Id)
	ofiId := string(*document.Assgnmt.Assgne.Agt.FinInstnId.Othr.Id)
	originalReqId := string(*document.CxlDtls[0].OrgnlGrpInfAndSts.OrgnlMsgId)
	originalInstrId := string(*document.CxlDtls[0].TxInfAndSts[0].OrgnlInstrId)
	rfiSettlementAccountName := strings.ToLower(string(*document.CxlDtls[0].TxInfAndSts[0].OrgnlTxRef.SttlmInf.SttlmAcct.Nm))
	rejectStatus := string(*document.Sts.Conf)

	checkData := &sendmodel.XMLData{
		OriginalMsgId:            originalReqId,
		OFIId:                    ofiId,
		RFIId:                    rfiId,
		RFISettlementAccountName: rfiSettlementAccountName,
		InstructionId:            string(*document.Assgnmt.Id),
		OriginalInstructionId:    originalInstrId,
	}

	// Check if the assignor is using the correct participant ID
	if !utils.StringsEqual(rfiId, homeDomain) {
		LOGGER.Error("Instructing agent is an incorrect participant")
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_FI
		return checkData, nil, errors.New("instructing agent is an incorrect participant")
	}

	// Get the status of the cancellation request data from DynamoDB
	dbData, txStatus, _, _, dbErr := database.DC.GetTransactionData(originalInstrId)

	// Query failed or data unmarshal failed
	if dbErr != nil {
		checkData.ErrorCode = constant.STATUS_CODE_INTERNAL_ERROR
		return checkData, nil, errors.New("database query error")
	}

	// Not corresponding data exist in the database
	if *dbData == "" {
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_ORIGINAL_ID
		return checkData, nil, errors.New("wrong original instruction ID")
	}

	if *txStatus != constant.DATABASE_STATUS_CANCEL_INIT {
		LOGGER.Error("Data not found in database")
		checkData.ErrorCode = constant.STATUS_CODE_ORIGINAL_REQUEST_NOT_INIT
		return checkData, nil, nil
	}

	// Check if the response status if RJCR
	if !utils.StringsEqual(rejectStatus, constant.PAYMENT_STATUS_RJCR) {
		LOGGER.Error("Unknown response status code")
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_RESPONSE_CODE
		return checkData, nil, errors.New("unknown response status code")
	}

	/*
		Aggregate necessary data for transaction memo
	*/
	statusData := &sendmodel.StatusData{
		IdCdtr:                rfiId,
		IdDbtr:                ofiId,
		InstructionID:         string(*document.Assgnmt.Id),
		OriginalInstructionID: originalInstrId,
	}

	return checkData, statusData, nil
}
