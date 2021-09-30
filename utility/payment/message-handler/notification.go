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
	"github.com/GFTN/gftn-services/utility/payment/environment"
	"github.com/GFTN/gftn-services/utility/payment/utils"
	"github.com/GFTN/gftn-services/utility/xmldsig"

	blocklist_client "github.com/GFTN/gftn-services/administration-service/blocklist-client"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/utility/payment/client"
	"github.com/GFTN/gftn-services/utility/payment/utils/horizon"
	"github.com/GFTN/gftn-services/utility/payment/utils/transaction"

	ibwf002struct "github.com/GFTN/iso20022/ibwf00200101"
	pacs002struct "github.com/GFTN/iso20022/pacs00200109"
	pacs004struct "github.com/GFTN/iso20022/pacs00400109"

	ibwfPbStruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/ibwf00200101"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/kafka"
	"github.com/GFTN/gftn-services/utility/payment/constant"
	message_converter "github.com/GFTN/gftn-services/utility/payment/message-converter"
	"github.com/GFTN/gftn-services/utility/payment/utils/database"

	"github.com/GFTN/gftn-services/utility/common"
	"github.com/GFTN/gftn-services/utility/payment/utils/parse"
	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"
)

type FbTrxLog struct {
	ParticipantID   string                 `json:"participant_id"`
	TransactionMemo map[string]interface{} `json:"transaction_memo"`
}

func (op *PaymentOperations) Ibwf002(ibwf002 message_converter.Ibwf002) ([]byte, error) {
	structData := ibwf002.Message

	/*
		validate content in the ibwf002 message and get all the necessary data from it
	*/
	msgName := constant.IBWF002
	BIC := os.Getenv(environment.ENV_KEY_PARTICIPANT_BIC)
	xmlData, statsData, err := getCriticalInfoFromIbwf002(structData, op.prServiceURL, op.homeDomain)
	statusCode := xmlData.ErrorCode
	rfiId := xmlData.RFIId
	ofiId := xmlData.OFIId
	target := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)

	originalGrpInf := &pacs002struct.OriginalGroupInformation29{
		OrgnlMsgId:   getReportMax35Text(string(*structData.DigOblSetNotif.GrpHdr.MsgId)),
		OrgnlMsgNmId: getReportMax35Text(msgName),
	}

	if statusCode != constant.STATUS_CODE_DEFAULT {
		LOGGER.Errorf("Something wrong with the transaction information: %v", err)
		report := parse.CreatePacs002(BIC, xmlData.InstructionId, target, statusCode, originalGrpInf)
		return report, errors.New("something wrong with the transaction information")
	}

	if err != nil {
		LOGGER.Error(err.Error())
		report := parse.CreatePacs002(BIC, xmlData.InstructionId, target, statusCode, originalGrpInf)
		return report, err
	}

	// Create admin-service client for connecting to admin-service
	blockListClient := blocklist_client.Client{
		HTTPClient: &http.Client{Timeout: time.Second * 80},
		AdminUrl:   os.Getenv(global_environment.ENV_KEY_ADMIN_SVC_URL),
	}

	/*
		preparing the data that need to be verify against the block-list
	*/
	var countries []string
	countries = append(countries, strconv.Quote(string(*structData.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.Cdtr.PstlAdr.Ctry)))
	countries = append(countries, strconv.Quote(string(*structData.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.Dbtr.PstlAdr.Ctry)))

	var currencies []string
	currencies = append(currencies, strconv.Quote(string(*structData.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.IntrBkSttlmAmt.Ccy)))

	var participants []string
	participants = append(participants, strconv.Quote(string(*structData.DigOblSetNotif.GrpHdr.InstgAgt.FinInstnId.Othr.Id)))
	participants = append(participants, strconv.Quote(string(*structData.DigOblSetNotif.GrpHdr.InstdAgt.FinInstnId.Othr.Id)))

	/*
		validate block-list
	*/
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

	ibwf002LogHandler := transaction.InitiatePaymentLogOperation()
	// Message type for payment status log : digital_obligation_notification
	msgType := constant.PAYMENT_TYPE_RDO
	// Initialize log handler and set the payment status to `INITIAL`
	ibwf002LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_INITIAL)
	err = database.SyncWithDynamo(constant.DATABASE_INIT, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_PENDING, constant.DATABASE_STATUS_NONE, ibwf002LogHandler)
	if err != nil {
		LOGGER.Errorf(err.Error())
		report := parse.CreatePacs002(BIC, xmlData.InstructionId, target, constant.STATUS_CODE_DUP_ID, originalGrpInf)
		return report, err
	}

	/*
		Check mutual whitelist
	*/

	LOGGER.Infof("Check whether RFI is in OFI's whitelist and vice versa")
	pkey, whiteListErr := op.whitelistHandler.CheckWhiteListParticipant(ofiId, rfiId, constant.EMPTY_STRING)
	if whiteListErr != nil {
		LOGGER.Errorf(whiteListErr.Error())
		ibwf002LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", ibwf002LogHandler, &op.fundHandler, statsData)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, ibwf002LogHandler)
		report := parse.CreatePacs002(BIC, xmlData.InstructionId, target, constant.STATUS_CODE_INTERNAL_ERROR, originalGrpInf)
		return report, whiteListErr
	}
	if pkey == "" {
		errMsg := "OFI can not find RFI in whitelist and vice versa"
		LOGGER.Errorf(errMsg)
		ibwf002LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_FAIL)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", ibwf002LogHandler, &op.fundHandler, statsData)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, ibwf002LogHandler)
		report := parse.CreatePacs002(BIC, xmlData.InstructionId, target, constant.STATUS_CODE_RFI_OR_OFI_NOT_IN_WL, originalGrpInf)
		return report, whiteListErr
	}
	LOGGER.Infof("Yes, RFI is in OFI's whitelist and vice versa")

	/*
		Parse the ibwf002 message with signature into ProtoBuffer
	*/

	LOGGER.Infof("Start parsing Go struct to ProtoBuffer")
	protoBufData, parseErr := proto.Marshal(&ibwf002.SendPayload)
	if parseErr != nil {
		errMsg := "Parse data to ProtoBuf error: " + parseErr.Error()
		LOGGER.Errorf(errMsg)
		ibwf002LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", ibwf002LogHandler, &op.fundHandler, statsData)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, ibwf002LogHandler)
		report := parse.CreatePacs002(BIC, xmlData.InstructionId, target, constant.STATUS_CODE_INTERNAL_ERROR, originalGrpInf)
		return report, parseErr
	}
	LOGGER.Infof("Finished parsing Go struct to ProtoBuffer")

	ibwf002LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_SUCCESS)

	dbData := sendmodel.DBData{
		MessageId: string(*structData.DigOblSetNotif.GrpHdr.MsgId),
	}

	dbDataByte, _ := json.Marshal(dbData)
	base64DBData := parse.EncodeBase64(dbDataByte)

	database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, base64DBData, constant.DATABASE_STATUS_PENDING, constant.DATABASE_STATUS_NONE, ibwf002LogHandler)
	database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", ibwf002LogHandler, &op.fundHandler, statsData)

	/*
		Send the ProtoBuffer to the request topic of RFI on Kafka broker
	*/

	LOGGER.Infof("Start to send request to Kafka broker")
	kafkaErr := op.KafkaActor.Produce(rfiId+kafka.REQUEST_TOPIC, protoBufData)
	if kafkaErr != nil {
		errMsg := "Error while submit message to Kafka broker: " + kafkaErr.Error()
		LOGGER.Errorf(errMsg)
		ibwf002LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", ibwf002LogHandler, &op.fundHandler, statsData)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, ibwf002LogHandler)
		report := parse.CreatePacs002(BIC, xmlData.InstructionId, target, constant.STATUS_CODE_INTERNAL_ERROR, originalGrpInf)
		return report, kafkaErr
	}

	LOGGER.Infof("Successfully produce message to Kafka broker")
	LOGGER.Debug("-----------------------------------------------------------------")
	report := parse.CreatePacs002(BIC, xmlData.InstructionId, target, constant.STATUS_CODE_RDO_REQ_SEND_TO_KAFKA, originalGrpInf)

	return report, nil
}

/*
	pacs.004.001.08 PaymentReturn
*/
func (op *PaymentOperations) Pacs004_Rdo(pacs004 message_converter.Pacs004) ([]byte, error) {
	// Validate content inside the pacs004 message
	structData := pacs004.Message
	BIC := os.Getenv(environment.ENV_KEY_PARTICIPANT_BIC)
	xmlMsgType := constant.PACS004
	xmlData, statusData, pacs008PaymentInfo, getDataErr := getCriticalInfoFromPacs004Rdo(structData.Body, op.homeDomain)
	errCode := xmlData.ErrorCode
	rfiId := xmlData.RFIId
	ofiId := xmlData.OFIId
	msgId := xmlData.MessageId
	target := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	msgName := constant.PACS004

	originalGrpInf := &pacs002struct.OriginalGroupInformation29{
		OrgnlMsgId:   getReportMax35Text(xmlData.MessageId),
		OrgnlMsgNmId: getReportMax35Text(msgName),
	}

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
	msgType := constant.PAYMENT_TYPE_RDO
	pacs008LogHandler := transaction.InitiatePaymentLogOperation()
	pacs004LogHandler := transaction.InitiatePaymentLogOperation()
	pacs008LogHandler.PaymentStatuses = pacs008PaymentInfo

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
		pacs004LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, pacs004LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", pacs004LogHandler, &op.fundHandler, statusData)
		report := parse.CreatePacs002(BIC, xmlData.InstructionId, target, constant.STATUS_CODE_INTERNAL_ERROR, originalGrpInf)
		return report, whiteListErr
	}
	if rfiAccount == "" {
		errMsg := "RFI can not find OFI in whitelist and vice versa"
		LOGGER.Error(errMsg)
		pacs004LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, pacs004LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", pacs004LogHandler, &op.fundHandler, statusData)
		report := parse.CreatePacs002(BIC, xmlData.InstructionId, target, constant.STATUS_CODE_RFI_OR_OFI_NOT_IN_WL, originalGrpInf)
		return report, whiteListErr
	}
	LOGGER.Infof("Yes, OFI is in RFI's whitelist and vice versa")

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

	pacs004LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_SUCCESS)

	submitResult, txHash, _ := op.fundHandler.FundAndSubmitPaymentTransaction(rfiAccount, xmlData.InstructionId, msgId, xmlMsgType, rfiAccountName, *signData, xdr.Memo{})
	report := parse.CreatePacs002(BIC, xmlData.InstructionId, target, submitResult, originalGrpInf)

	if submitResult != constant.STATUS_CODE_TX_SEND_TO_STELLAR {
		pacs004LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_SUBMIT_FAIL)
		pacs008LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_SUBMIT_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, pacs004LogHandler)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.OriginalInstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_DONE, constant.DATABASE_STATUS_NONE, pacs008LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", rfiAccount, pacs004LogHandler, &op.fundHandler, statusData)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, msgType, constant.PACS008, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.OriginalInstructionId, "", "", pacs008LogHandler, &op.fundHandler, statusData)
		return report, nil
	} else {
		pacs004LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_SETTLED, txHash)
		pacs008LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_SETTLED, txHash)
	}

	go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.OriginalInstructionId, txHash, constant.DATABASE_STATUS_SETTLED, constant.DATABASE_STATUS_NONE, pacs008LogHandler)
	go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, txHash, constant.DATABASE_STATUS_DONE, constant.DATABASE_STATUS_NONE, pacs004LogHandler)

	go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, txHash, rfiAccount, pacs004LogHandler, &op.fundHandler, statusData)
	go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, msgType, constant.PACS008, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.OriginalInstructionId, txHash, "", pacs008LogHandler, &op.fundHandler, statusData)

	LOGGER.Debug("---------------------------------------------------------------------")

	return report, nil
}

// if message type is ibwf.002
func RFI_Ibwf002(data ibwfPbStruct.SendPayload, op *kafka.KafkaOpreations) {
	// Parse the ProtoBuffer into Go struct and reconstruct it into pacs008 message
	LOGGER.Infof("Parsing ProtoBuffer to XML")
	standardType := constant.ISO20022
	// Message type for payment status log : digital_obligatoin_notification
	paymentStatusMsgType := constant.PAYMENT_TYPE_RDO
	msgName := constant.IBWF002
	// Initialize the payment status
	pacs008LogHandler := transaction.InitiatePaymentLogOperation()
	ibwf002LogHandler := transaction.InitiatePaymentLogOperation()

	pacs008InstructionId := data.OriginalInstructionId
	instructionId := data.InstructionId
	reqMsgType := data.MsgType
	ofiId := data.OfiId
	rfiId := data.RfiId
	originalMsgId := data.OriginalMsgId
	msgId := data.MsgId
	topicName := rfiId + "_" + kafka.TRANSACTION_TOPIC

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
		find ibwf002 record from DB
	*/
	ibwf002DynamoData, ibwf002PaymentInfo := parse.GetDBData(instructionId)
	if ibwf002DynamoData == nil || ibwf002PaymentInfo == nil {
		LOGGER.Errorf("The original message ID %v does not exist in DB", instructionId)
		op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_WRONG_INSTRUCTION_ID, originalGrpInf)
		return
	}

	ibwf002LogHandler.PaymentStatuses = ibwf002PaymentInfo
	ibwf002LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RFI_PROCESSING)

	/*
		verify signature
	*/
	LOGGER.Infof("Verifying the signature from OFI")
	result := xmldsig.VerifySignature(string(data.Message))
	if !result {
		LOGGER.Errorf("signature verification failed")
		ibwf002LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RFI_VALIDATION_FAIL)
		//ibwf002 status
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, ibwf002LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, pacs008InstructionId, instructionId, "", "", ibwf002LogHandler, &op.FundHandler, &sendmodel.StatusData{})
		op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_OFI_SIGNATURE_FAIL, originalGrpInf)
		return
	}
	LOGGER.Infof("OFI signature verified!")

	/*
		constructing protobuffer to go struct
	*/
	ibwf002 := &message_converter.Ibwf002{SendPayload: data}
	xmlData, err := ibwf002.ProtobuftoStruct()
	if xmlData == nil {
		LOGGER.Errorf("Encounter error while construncting proto buffer to go struct")
		ibwf002LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RFI_VALIDATION_FAIL)
		//ibwf002 status
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, ibwf002LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, pacs008InstructionId, instructionId, "", "", ibwf002LogHandler, &op.FundHandler, &sendmodel.StatusData{})
		op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR_PARSE, originalGrpInf)
		return
	} else if err != nil {
		LOGGER.Errorf("Parse request from kafka failed: %s", err.Error())
		ibwf002LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RFI_VALIDATION_FAIL)
		//ibwf002 status
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, ibwf002LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, pacs008InstructionId, instructionId, "", "", ibwf002LogHandler, &op.FundHandler, &sendmodel.StatusData{})
		op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR_PARSE, originalGrpInf)
		return
	}
	// Get important data from the XML data
	reqData := xmlData.RequestXMLMsg
	settlementAccountName := xmlData.OFISettlementAccountName
	originalInstructionId := xmlData.OriginalInstructionId

	LOGGER.Infof("Finished paring ProtoBuffer to XML")

	// Generate payment status data
	// Aggregate necessary data for transaction memo
	settlementAmount := ibwf002.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.IntrBkSttlmAmt.Value

	statusData := &sendmodel.StatusData{
		CityCdtr:             string(*ibwf002.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.Cdtr.PstlAdr.TwnNm),
		CountryCdtr:          string(*ibwf002.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.Cdtr.PstlAdr.Ctry),
		NameCdtr:             string(*ibwf002.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.Cdtr.Nm),
		IdCdtr:               string(*ibwf002.Message.DigOblSetNotif.GrpHdr.InstdAgt.FinInstnId.Othr.Id),
		CityDbtr:             string(*ibwf002.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.Dbtr.PstlAdr.TwnNm),
		CountryDbtr:          string(*ibwf002.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.Dbtr.PstlAdr.Ctry),
		NameDbtr:             string(*ibwf002.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.Dbtr.Nm),
		IdDbtr:               string(*ibwf002.Message.DigOblSetNotif.GrpHdr.InstgAgt.FinInstnId.Othr.Id),
		CurrencyCode:         string(*ibwf002.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.IntrBkSttlmAmt.Ccy),
		AssetType:            string(*ibwf002.Message.DigOblSetNotif.GrpHdr.SttlmInf.SttlmMtd),
		CreditorStreet:       string(*ibwf002.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.Cdtr.PstlAdr.StrtNm),
		CreditorBuildingNo:   string(*ibwf002.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.Cdtr.PstlAdr.BldgNb),
		CreditorPostalCode:   string(*ibwf002.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.Cdtr.PstlAdr.PstCd),
		AssetCodeBeneficiary: string(*ibwf002.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.IntrBkSttlmAmt.Ccy)[:3],
		CrtyCcy:              string(*ibwf002.Message.DigOblSetNotif.GrpHdr.InstdAgt.FinInstnId.BICFI)[:3],
		CustomerStreet:       string(*ibwf002.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.Cdtr.PstlAdr.StrtNm),
		CustomerBuildingNo:   string(*ibwf002.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.Cdtr.PstlAdr.BldgNb),
		CustomerCountry:      string(*ibwf002.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.Cdtr.PstlAdr.PstCd),
		AccountNameSend:      string(*ibwf002.Message.DigOblSetNotif.GrpHdr.SttlmInf.SttlmAcct.Nm),
		EndToEndID:           string(*ibwf002.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlEndToEndId),
		InstructionID:        string(*ibwf002.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].NtfId),
		AmountSettlement:     settlementAmount,
	}

	rfiVerifyRequestAndSendToKafka(topicName, msgId, msgName, originalMsgId, ofiId, settlementAccountName, standardType, msgName, instructionId, originalInstructionId, paymentStatusMsgType, pacs008LogHandler, reqData, statusData, pacs008DynamoData, op, originalGrpInf)

	rawMsg, _ := json.Marshal(ibwf002DynamoData)
	ibwf002LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RDO_INIT)
	go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, parse.EncodeBase64(rawMsg), constant.DATABASE_STATUS_DONE, constant.DATABASE_STATUS_NONE, ibwf002LogHandler)
	go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, paymentStatusMsgType, msgName, originalMsgId, originalInstructionId, instructionId, "", "", ibwf002LogHandler, &op.FundHandler, statusData)

	return
}

func getCriticalInfoFromIbwf002(document *ibwf002struct.Message, prServiceURL, homeDomain string) (*sendmodel.XMLData, *sendmodel.StatusData, error) {
	originalInstructionId := string(*document.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlInstrId)
	instructionID := string(*document.DigOblSetNotif.SttlOblInf[0].TxInf[0].NtfId)
	instructedAgent := string(*document.DigOblSetNotif.GrpHdr.InstdAgt.FinInstnId.Othr.Id)
	instructingAgent := string(*document.DigOblSetNotif.GrpHdr.InstgAgt.FinInstnId.Othr.Id)
	ofiSettlementAccountName := strings.ToLower(string(*document.DigOblSetNotif.GrpHdr.SttlmInf.SttlmAcct.Nm))
	currencyCode := string(*document.DigOblSetNotif.SttlOblInf[0].TxInf[0].ActSttldAmt.Ccy)
	settlementMethod := string(*document.DigOblSetNotif.GrpHdr.SttlmInf.SttlmMtd)
	assetIssuerId := string(*document.DigOblSetNotif.GrpHdr.PmtTpInf.SvcLvl.Prtry)
	originalMsgType := string(*document.DigOblSetNotif.SttlOblInf[0].OrgnlGrpInf.OrgnlMsgNmId)
	originalMsgID := string(*document.DigOblSetNotif.SttlOblInf[0].OrgnlGrpInf.OrgnlMsgId)

	checkData := &sendmodel.XMLData{}

	checkData.OriginalInstructionId = originalInstructionId
	checkData.OriginalMsgId = originalMsgID
	checkData.RequestMsgType = originalMsgType
	checkData.InstructionId = instructionID
	checkData.OFIId = instructingAgent
	checkData.RFIId = instructedAgent
	checkData.OFISettlementAccountName = ofiSettlementAccountName
	checkData.ErrorCode = constant.STATUS_CODE_DEFAULT

	if !utils.StringsEqual(instructingAgent, homeDomain) {
		LOGGER.Error("Instructing agent is an incorrect participant")
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_FI
		return checkData, &sendmodel.StatusData{}, errors.New("instructing agent is an incorrect participant")
	}

	account := client.GetParticipantAccount(prServiceURL, homeDomain, ofiSettlementAccountName)

	if account == nil {
		LOGGER.Error("No corresponding account for participant")
		checkData.ErrorCode = constant.STATUS_CODE_ACCOUNT_NOT_EXIST
		return checkData, &sendmodel.StatusData{}, errors.New("no corresponding account for participant")
	}

	dbData, txStatus, resId, _, dbErr := database.DC.GetTransactionData(originalInstructionId)

	if dbErr != nil {
		LOGGER.Errorf("database query error: %v", dbErr)
		checkData.ErrorCode = constant.STATUS_CODE_INTERNAL_ERROR
		return checkData, &sendmodel.StatusData{}, errors.New("database query error")
	}

	if *dbData == "" {
		LOGGER.Errorf("Incorrect instruction ID")
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_ORIGINAL_ID
		return checkData, &sendmodel.StatusData{}, errors.New("Incorrect instruction ID")
	}

	/*
		DO sanity check
	*/

	// if it is DO, check if they are using issuing account & if either OFI or RFI is the issuer
	if !utils.StringsEqual(settlementMethod, constant.DO_SETTLEMENT) {
		errMsg := "The currency code of this message type must be DO"
		LOGGER.Errorf(errMsg)
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_SETTLEMENT_METHOD
		return checkData, &sendmodel.StatusData{}, errors.New(errMsg)
	}

	// check if the settlement account name is "issuing"
	if !utils.StringsEqual(ofiSettlementAccountName, common.ISSUING) {
		LOGGER.Error("Account name should be \"issuing\", if settlement method is WWDO")
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_ACCOUNT_NAME
		return checkData, &sendmodel.StatusData{}, errors.New("wrong account name for DO")
	}

	if !horizon.IsIssuer(assetIssuerId, currencyCode) {
		errMsg := "The asset " + currencyCode + " is not issued by the asset issuer " + assetIssuerId
		LOGGER.Errorf(errMsg)
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_ASSET_ISSUER
		return checkData, &sendmodel.StatusData{}, errors.New(errMsg)
	}

	// check if settlement asset currency code is ended with "DO"
	if !strings.HasSuffix(currencyCode, constant.SETTLEMENT_METHOD_DIGITAL_OBLIGATION) {
		errMsg := "Settlement method is DO, please use DO as settlement currency"
		LOGGER.Error(errMsg)
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_SETTLEMENT_METHOD
		return checkData, &sendmodel.StatusData{}, errors.New(errMsg)
	}

	if *resId != constant.DATABASE_STATUS_NONE {
		if *resId == constant.DATABASE_STATUS_FAILED {
			checkData.ErrorCode = constant.STATUS_CODE_REQUEST_CLOSE
			return checkData, &sendmodel.StatusData{}, errors.New("request was closed due to internal errors")
		} else if *resId == constant.DATABASE_STATUS_SETTLED {
			LOGGER.Errorf("The DO is already settled")
			checkData.ErrorCode = constant.STATUS_CODE_ALREADY_SETTLED
			return checkData, &sendmodel.StatusData{}, errors.New("The DO is already settled")
		}
	}

	/*

		if *resId != constant.DATABASE_STATUS_NONE {
			if *resId == constant.DATABASE_STATUS_FAILED {
				checkData.ErrorCode = constant.STATUS_CODE_REQUEST_CLOSE
				return checkData, errors.New("request was closed due to internal errors")
			} else if *resId == constant.DATABASE_STATUS_ {
				LOGGER.Errorf("The DO is already st")
				checkData.ErrorCode = constant.STATUS_CODE_ALREADY_SETTLED
				return checkData, errors.New("message id was used")
			}else {
				checkData.ErrorCode = constant.STATUS_CODE_DUP_ID
				return checkData, errors.New("message id was used")
			}
		}

		if *resId != constant.DATABASE_STATUS_NONE {
			if *resId == constant.DATABASE_STATUS_FAILED {
				LOGGER.Errorf("request was closed due to internal errors")
				checkData.ErrorCode = constant.STATUS_CODE_REQUEST_CLOSE
				return checkData, errors.New("request was closed due to internal errors")
			} else if *resId == constant.DATABASE_STATUS_ {
				LOGGER.Errorf("The DO is already st")
				checkData.ErrorCode = constant.STATUS_CODE_ALREADY_SETTLED
				return checkData, errors.New("message id was used")
			}
		}
	*/

	if *txStatus != constant.DATABASE_STATUS_CLEARED && *txStatus != constant.DATABASE_STATUS_CANCELED {
		LOGGER.Error("Data not found in database")
		checkData.ErrorCode = constant.STATUS_CODE_ORIGINAL_REQUEST_NOT_DONE
		return checkData, &sendmodel.StatusData{}, dbErr
	}

	/*
		Aggregate necessary data for transaction memo
	*/

	statsData := &sendmodel.StatusData{
		NameCdtr:             string(*document.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.Cdtr.Nm),
		IdCdtr:               instructedAgent,
		NameDbtr:             string(*document.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.Dbtr.Nm),
		IdDbtr:               instructingAgent,
		CurrencyCode:         string(*document.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.IntrBkSttlmAmt.Ccy),
		AssetType:            string(*document.DigOblSetNotif.GrpHdr.SttlmInf.SttlmMtd),
		AssetCodeBeneficiary: string(*document.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.IntrBkSttlmAmt.Ccy)[:3],
		CrtyCcy:              string(*document.DigOblSetNotif.GrpHdr.InstdAgt.FinInstnId.BICFI)[:3],
		AccountNameSend:      string(*document.DigOblSetNotif.GrpHdr.SttlmInf.SttlmAcct.Nm),
		EndToEndID:           string(*document.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlEndToEndId),
		InstructionID:        string(*document.DigOblSetNotif.SttlOblInf[0].TxInf[0].NtfId),
		AmountSettlement:     document.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlTxRef.IntrBkSttlmAmt.Value,
		//IssuerID:             document.Body.GrpHdr.PmtTpInf.SvcLvl.Prtry,
	}

	return checkData, statsData, nil
}

func getCriticalInfoFromPacs004Rdo(document *pacs004struct.PaymentReturnV09, homeDomain string) (*sendmodel.XMLData, *sendmodel.StatusData, []model.TransactionReceipt, error) {
	originalMsgID := string(*document.OrgnlGrpInf.OrgnlMsgId)
	settlementMethod := string(*document.GrpHdr.SttlmInf.SttlmMtd)
	instructedAgent := string(*document.GrpHdr.InstdAgt.FinInstnId.Othr.Id)
	instructingAgent := string(*document.GrpHdr.InstgAgt.FinInstnId.Othr.Id)
	ofiSettlementAccountName := strings.ToLower(string(*document.TxInf[0].OrgnlTxRef.SttlmInf.SttlmAcct.Nm))
	rfiSettlementAccountName := strings.ToLower(string(*document.GrpHdr.SttlmInf.SttlmAcct.Nm))
	msgId := string(*document.GrpHdr.MsgId)
	assetIssuerId := string(*document.TxInf[0].OrgnlTxRef.PmtTpInf.SvcLvl[0].Prtry)
	currencyCode := document.TxInf[0].RtrdIntrBkSttlmAmt.Currency
	originalMsgType := string(*document.OrgnlGrpInf.OrgnlMsgNmId)
	originalEndtoEndID := string(*document.TxInf[0].OrgnlEndToEndId)
	originalInstructionId := string(*document.TxInf[0].OrgnlInstrId)
	instructionId := string(*document.TxInf[0].RtrId)

	checkData := &sendmodel.XMLData{}

	checkData.OriginalMsgId = originalMsgID
	checkData.OriginalEndtoEndId = originalEndtoEndID
	checkData.RequestMsgType = originalMsgType
	checkData.OriginalInstructionId = originalInstructionId
	checkData.InstructionId = instructionId
	checkData.OFIId = instructedAgent
	checkData.RFIId = instructingAgent
	checkData.OFISettlementAccountName = ofiSettlementAccountName
	checkData.RFISettlementAccountName = rfiSettlementAccountName
	checkData.ErrorCode = constant.STATUS_CODE_DEFAULT
	checkData.MessageId = msgId
	checkData.AssetIssuer = assetIssuerId

	if !utils.StringsEqual(instructingAgent, homeDomain) {
		LOGGER.Error("Instructing agent is an incorrect participant")
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_FI
		return checkData, nil, nil, errors.New("instructing agent is an incorrect participant")
	}

	dbData, txStatus, _, paymentInfoBase64, dbErr := database.DC.GetTransactionData(checkData.OriginalInstructionId)

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

	// if it is DO, check if they are using issuing account & if either OFI or RFI is the issuer
	if utils.StringsEqual(settlementMethod, constant.DO_SETTLEMENT) {
		// check if this DO was issued by either OFI or RFI
		if !utils.StringsEqual(instructedAgent, assetIssuerId) && !utils.StringsEqual(instructingAgent, assetIssuerId) {
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

		if utils.StringsEqual(instructedAgent, instructingAgent) {
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

	if *txStatus == constant.DATABASE_STATUS_RDO_INIT {
		reqSettlementMethod := string(*document.TxInf[0].OrgnlTxRef.SttlmInf.SttlmMtd)
		if !utils.StringsEqual(reqSettlementMethod, settlementMethod) {
			checkData.ErrorCode = constant.STATUS_CODE_WRONG_SETTLEMENT_METHOD
			return checkData, nil, nil, errors.New("settlement method is not the same as rdo request")
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

	amountBeneficiary, _ := strconv.ParseFloat(document.TxInf[0].RtrdInstdAmt.Value, 64)
	settlementAmount, _ := strconv.ParseFloat(document.TxInf[0].RtrdIntrBkSttlmAmt.Value, 64)

	/*
		Aggregate necessary data for transaction memo
	*/
	statusData := &sendmodel.StatusData{
		IdCdtr:               instructingAgent,
		IdDbtr:               instructedAgent,
		CurrencyCode:         document.TxInf[0].RtrdIntrBkSttlmAmt.Currency,
		AssetType:            string(*document.GrpHdr.SttlmInf.SttlmMtd),
		AmountBeneficiary:    amountBeneficiary,
		AssetCodeBeneficiary: document.TxInf[0].RtrdInstdAmt.Currency,
		AccountNameSend:      string(*document.GrpHdr.SttlmInf.SttlmAcct.Nm),
		EndToEndID:           originalEndtoEndID,
		InstructionID:        string(*document.TxInf[0].RtrId),
		AmountSettlement:     settlementAmount,
		IssuerID:             string(*document.TxInf[0].OrgnlTxRef.PmtTpInf.SvcLvl[0].Prtry),
		ExchangeRate:         1.0,
	}

	return checkData, statusData, paymentInfo, nil
}
