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
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/utility/payment/environment"
	"github.com/GFTN/gftn-services/utility/payment/utils"
	"github.com/GFTN/gftn-services/utility/xmldsig"

	"github.com/GFTN/gftn-services/utility/payment/utils/transaction"
	whitelist_handler "github.com/GFTN/gftn-services/utility/payment/utils/whitelist-handler"

	message_converter "github.com/GFTN/gftn-services/utility/payment/message-converter"

	camt026struct "github.com/GFTN/iso20022/camt02600107"
	camt087struct "github.com/GFTN/iso20022/camt08700106"
	pacs002struct "github.com/GFTN/iso20022/pacs00200109"

	camt026pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/camt02600107"
	camt087pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/camt08700106"

	blocklist_client "github.com/GFTN/gftn-services/administration-service/blocklist-client"
	"github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/kafka"
	"github.com/GFTN/gftn-services/utility/payment/constant"
	"github.com/GFTN/gftn-services/utility/payment/utils/database"
	"github.com/GFTN/gftn-services/utility/payment/utils/parse"
	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"
)

// camt.26 message handler at RFI side
func (op *PaymentOperations) Camt026(camt026 message_converter.Camt026) ([]byte, error) {

	structData := camt026.Message
	BIC := os.Getenv(environment.ENV_KEY_PARTICIPANT_BIC)
	xmlMsgType := constant.CAMT026

	/*
		payload check
	*/

	// validate content in the camt026 message and get all the necessary data from it
	xmlData, statsData, err := getCriticalInfoFromCamt026(structData, op.homeDomain)
	ofiId := xmlData.OFIId
	rfiId := xmlData.RFIId
	statusCode := xmlData.ErrorCode

	originalGrpInf := &pacs002struct.OriginalGroupInformation29{
		OrgnlMsgId:   getReportMax35Text(xmlData.OriginalMsgId),
		OrgnlMsgNmId: getReportMax35Text(xmlMsgType),
	}

	target, _, err := parse.KafkaErrorRouter(xmlMsgType, xmlData.MessageId, ofiId, rfiId, 0, false, originalGrpInf)
	if statusCode != constant.STATUS_CODE_DEFAULT || err != nil {
		LOGGER.Errorf("Something wrong with the transaction information")
		report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, xmlData.ErrorCode)
		return report, errors.New("something wrong with the transaction information")
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

	/*
		start operations that will rely on the other services
	*/

	// Message type for payment status log : credit_transfer
	msgType := constant.PAYMENT_TYPE_EXCEPTION
	msgName := constant.CAMT026
	logHandler := transaction.InitiatePaymentLogOperation()

	/*
		write into dynamo
	*/
	// Initialize log handler and set the payment status to `INITIAL`
	logHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_INITIAL)
	err = database.SyncWithDynamo(constant.DATABASE_INIT, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_PENDING, constant.DATABASE_STATUS_NONE, logHandler)
	if err != nil {
		LOGGER.Errorf(err.Error())
		report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, constant.STATUS_CODE_DUP_ID)
		return report, err
	}

	// Check mutual whitelist
	LOGGER.Infof("Check whether RFI is in OFI's whitelist and vice versa")
	pkey, whiteListErr := op.whitelistHandler.CheckWhiteListParticipant(xmlData.OFIId, xmlData.RFIId, constant.EMPTY_STRING)
	if whiteListErr != nil {
		LOGGER.Errorf(whiteListErr.Error())
		logHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.MessageId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", logHandler, &op.fundHandler, statsData)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, logHandler)
		report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, constant.STATUS_CODE_INTERNAL_ERROR)
		return report, whiteListErr
	}

	if pkey == "" {
		errMsg := "OFI can not find RFI in whitelist and vice versa"
		LOGGER.Errorf(errMsg)
		logHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_FAIL)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.MessageId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", logHandler, &op.fundHandler, statsData)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, logHandler)
		report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, constant.STATUS_CODE_RFI_OR_OFI_NOT_IN_WL)
		return report, nil
	}
	LOGGER.Infof("Yes, RFI is in OFI's whitelist and vice versa")

	/*
		Parse the camt026 message with signature into ProtoBuffer
	*/

	LOGGER.Infof("Start parsing Go struct to ProtoBuffer")
	protoBufData, parseErr := proto.Marshal(&camt026.SendPayload)
	if parseErr != nil {
		errMsg := "Parse data to ProtoBuf error: " + parseErr.Error()
		LOGGER.Errorf(errMsg)
		logHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_FAIL)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.MessageId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", logHandler, &op.fundHandler, statsData)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, logHandler)
		report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, constant.STATUS_CODE_INTERNAL_ERROR)
		return report, parseErr
	}
	LOGGER.Infof("Finished parsing Go struct to ProtoBuffer")

	// Process done with OFI side, now update the payment status
	logHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_SUCCESS)

	// dbData to be written into Dynamo DB

	//save the instruction id of camt087 for pacs004/camt029 msg to use
	reportData := sendmodel.DBData{
		InstrId: xmlData.InstructionId,
	}
	dbData, _ := json.Marshal(reportData)
	base64DBData := parse.EncodeBase64(dbData)

	database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, base64DBData, constant.DATABASE_STATUS_PENDING, constant.DATABASE_STATUS_NONE, logHandler)
	database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.MessageId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", logHandler, &op.fundHandler, statsData)

	/*
		Send the ProtoBuffer to the request topic of RFI on Kafka broker
	*/
	LOGGER.Infof("Start to send request to Kafka broker")
	kafkaErr := op.KafkaActor.Produce(xmlData.OFIId+kafka.REQUEST_TOPIC, protoBufData)
	if kafkaErr != nil {
		errMsg := "Error while submit message to Kafka broker: " + kafkaErr.Error()
		LOGGER.Errorf(errMsg)
		logHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.MessageId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", logHandler, &op.fundHandler, statsData)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, logHandler)
		report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, constant.STATUS_CODE_INTERNAL_ERROR)
		return report, kafkaErr
	}

	LOGGER.Infof("Successfully produce message to Kafka broker")
	LOGGER.Debug("-----------------------------------------------------------------")

	// Send status back to OFI
	report := parse.CreateCamt030(BIC, xmlData.InstructionId, xmlMsgType, target, constant.STATUS_CODE_OFI_SEND_TO_KAFKA)

	return report, nil
}

// camt.26 message handler at OFI side
func OFI_Camt026(sendPayload camt026pbstruct.SendPayload, op *kafka.KafkaOpreations) {
	// Parse the ProtoBuffer into Go struct and reconstruct it into camt026 message
	LOGGER.Infof("Parsing ProtoBuffer to XML")

	camt026InstructionId := sendPayload.InstructionId
	reqMsgType := sendPayload.MsgType
	rfiId := sendPayload.RfiId
	ofiId := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	originalInstrId := sendPayload.OriginalInstructionId
	originalMsgId := sendPayload.OriginalInstructionId

	standardType := constant.ISO20022
	xmlMsgType := constant.CAMT026
	msgType := constant.PAYMENT_TYPE_EXCEPTION
	topicName := ofiId + "_" + kafka.TRANSACTION_TOPIC

	originalGrpInf := &pacs002struct.OriginalGroupInformation29{
		OrgnlMsgId:   getReportMax35Text(originalMsgId),
		OrgnlMsgNmId: getReportMax35Text(xmlMsgType),
	}
	/*
		get camt026 data from DB
	*/
	camt026LogHandler := transaction.InitiatePaymentLogOperation()
	pacs008LogHandler := transaction.InitiatePaymentLogOperation()

	/*
	 Get pacs008 from database
	*/
	LOGGER.Infof("Get transaction related information from database")
	pacs008DbData, pacs008PaymentInfo := parse.GetDBData(originalInstrId)
	if pacs008DbData == nil || pacs008PaymentInfo == nil {
		LOGGER.Error("Can not get original pacs008 message from database")
		op.SendErrMsg(camt026InstructionId, standardType, xmlMsgType, ofiId, rfiId, constant.STATUS_CODE_WRONG_ORIGINAL_ID, originalGrpInf)
		return
	}
	pacs008LogHandler.PaymentStatuses = pacs008PaymentInfo

	/*
	 Get camt026 from database
	*/
	dynamoData, paymentInfo := parse.GetDBData(camt026InstructionId)
	if dynamoData == nil || paymentInfo == nil {
		LOGGER.Errorf("The original message ID %v does not exist in DB", camt026InstructionId)
		op.SendErrMsg(camt026InstructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_WRONG_ORIGINAL_ID, originalGrpInf)
		return
	}

	pacs008LogHandler.PaymentStatuses = pacs008PaymentInfo
	camt026LogHandler.PaymentStatuses = paymentInfo
	camt026LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_OFI_PROCESSING)

	/*
		verify signature
	*/
	LOGGER.Infof("Verifying the signature from OFI")
	result := xmldsig.VerifySignature(string(sendPayload.Message))
	if !result {
		LOGGER.Errorf("signature verification failed")
		camt026LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_OFI_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, camt026InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt026LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, msgType, xmlMsgType, originalMsgId, camt026InstructionId, camt026InstructionId, "", "", camt026LogHandler, &op.FundHandler, &sendmodel.StatusData{})
		op.SendErrMsg(camt026InstructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_OFI_SIGNATURE_FAIL, originalGrpInf)
		return
	}
	LOGGER.Infof("OFI signature verified!")

	/*
		constructing protobuffer to go struct
	*/
	camt026 := &message_converter.Camt026{SendPayload: sendPayload}
	xmlData, err := camt026.ProtobuftoStruct()
	if xmlData == nil {
		LOGGER.Errorf("Encounter error while construncting proto buffer to go struct")
		op.SendErrMsg(camt026InstructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR_PARSE, originalGrpInf)
		camt026LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_OFI_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, camt026InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt026LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, msgType, xmlMsgType, originalMsgId, camt026InstructionId, camt026InstructionId, "", "", camt026LogHandler, &op.FundHandler, &sendmodel.StatusData{})
		return
	} else if err != nil {
		LOGGER.Errorf("Parse request from kafka failed: %s", err.Error())
		op.SendErrMsg(camt026InstructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR_PARSE, originalGrpInf)
		camt026LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_OFI_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, camt026InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt026LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, msgType, xmlMsgType, originalMsgId, camt026InstructionId, camt026InstructionId, "", "", camt026LogHandler, &op.FundHandler, &sendmodel.StatusData{})
		return
	}

	LOGGER.Infof("Finished paring ProtoBuffer to Go struct")

	// Aggregate necessary data for transaction memo
	statusData := &sendmodel.StatusData{
		IdCdtr:                xmlData.RFIId,
		IdDbtr:                xmlData.OFIId,
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
		camt026LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, camt026InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt026LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, msgType, xmlMsgType, originalMsgId, originalInstrId, camt026InstructionId, "", "", camt026LogHandler, &op.FundHandler, statusData)
		op.SendErrMsg(camt026InstructionId, standardType, xmlMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR, originalGrpInf)
		return
	}
	if res == common.BlocklistDeniedString {
		LOGGER.Errorf("The transaction currency/country/institution is within the blocklist, transaction forbidden!")
		camt026LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_OFI_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, camt026InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt026LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, msgType, xmlMsgType, originalMsgId, originalInstrId, camt026InstructionId, "", "", camt026LogHandler, &op.FundHandler, statusData)
		op.SendErrMsg(camt026InstructionId, standardType, xmlMsgType, ofiId, rfiId, constant.STATUS_CODE_BLOCKLIST, originalGrpInf)
		return
	}

	// Check if RFI was whitelisted by OFI and vice versa, if not, reject the payment request
	whitelistHandler := whitelist_handler.CreateWhiteListServiceOperations()
	LOGGER.Infof("Check whether RFI is in OFI's whitelist and vice versa.")
	pKey, whiteListErr := whitelistHandler.CheckWhiteListParticipant(ofiId, rfiId, constant.EMPTY_STRING)
	if whiteListErr != nil {
		LOGGER.Errorf(whiteListErr.Error())
		op.SendErrMsg(camt026InstructionId, standardType, xmlMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR, originalGrpInf)
		camt026LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, camt026InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt026LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, msgType, xmlMsgType, originalMsgId, originalInstrId, camt026InstructionId, "", "", camt026LogHandler, &op.FundHandler, statusData)
		return
	}
	if pKey == "" {
		LOGGER.Errorf("Can not find RFI or OFI in whitelist and vice versa")
		op.SendErrMsg(camt026InstructionId, standardType, xmlMsgType, ofiId, rfiId, constant.STATUS_CODE_RFI_OR_OFI_NOT_IN_WL, originalGrpInf)
		camt026LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_OFI_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, camt026InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt026LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, msgType, xmlMsgType, originalMsgId, originalInstrId, camt026InstructionId, "", "", camt026LogHandler, &op.FundHandler, statusData)
		return
	}
	LOGGER.Infof("Yes, RFI is in OFI's whitelist and vice versa.")
	camt026LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_OFI_VALIDATION_SUCCESS)

	/*
		signing message with IBM master account
	*/
	signedMessage, signErr := op.SignHandler.SignPayloadByMasterAccount(xmlData.RequestXMLMsg)
	if signErr != nil {
		LOGGER.Errorf("Failed to sign payload: %v", signErr.Error())
		camt026LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, camt026InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt026LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, msgType, xmlMsgType, originalMsgId, originalInstrId, camt026InstructionId, "", "", camt026LogHandler, &op.FundHandler, statusData)
		op.SendErrMsg(camt026InstructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_SIGN_PAYLOAD_FAIL, originalGrpInf)
		return
	}
	gatewayMsg := parse.EncodeBase64(signedMessage)
	callBackMsg := &model.SendPacs{
		MessageType: &reqMsgType,
		Message:     &gatewayMsg,
	}

	// Send the encoded xml message to the callback service of RFI
	LOGGER.Infof("Send encoded message back to Kafka topic: %v", topicName)
	msg, _ := json.Marshal(callBackMsg)

	/*
		sending message to Kafka
	*/

	err = op.SendRequestToKafka(ofiId+"_"+kafka.TRANSACTION_TOPIC, msg)
	if err != nil {
		LOGGER.Errorf("Encounter error while producing message to Kafka topic: %v", ofiId+"_"+kafka.TRANSACTION_TOPIC)
		op.SendErrMsg(camt026InstructionId, standardType, xmlMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR, originalGrpInf)
		camt026LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, camt026InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt026LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, msgType, xmlMsgType, originalMsgId, originalInstrId, camt026InstructionId, "", "", camt026LogHandler, &op.FundHandler, statusData)
		return
	}

	pacs008LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_UNABLE_TO_APPLY)
	camt026LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_UNABLE_TO_APPLY)

	// Update transaction related information inside the DynamoDB base on message ID
	// (request ID, transaction hash, done, response ID, done)
	go database.SyncWithDynamo(constant.DATABASE_UPDATE, camt026InstructionId, "", constant.TX_DONE, xmlData.InstructionId, camt026LogHandler)
	go database.SyncWithDynamo(constant.DATABASE_UPDATE, originalInstrId, "", "", xmlData.InstructionId, pacs008LogHandler)

	// Store the transaction information into the administration service and FireBase
	go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, msgType, xmlMsgType, originalMsgId, originalInstrId, camt026InstructionId, "", "", camt026LogHandler, &op.FundHandler, statusData)
	go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, msgType, constant.PACS008, originalMsgId, originalInstrId, originalInstrId, "", "", pacs008LogHandler, &op.FundHandler, statusData)

	LOGGER.Debug("---------------------------------------------------------------------")
	return

}

func (op *PaymentOperations) Camt087(camt087 message_converter.Camt087) ([]byte, error) {
	// Validate content inside the camt087 message
	structData := camt087.Message
	xmlMsgType := constant.CAMT087
	xmlData, statusData, getInfoErr := getCriticalInfoFromCAMT087(structData.Body, op.homeDomain)
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

	camt087LogHandler := transaction.InitiatePaymentLogOperation()
	camt087LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_INITIAL)
	msgType := constant.PAYMENT_TYPE_EXCEPTION
	msgName := constant.CAMT087

	err = database.SyncWithDynamo(constant.DATABASE_INIT, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_PENDING, constant.DATABASE_STATUS_NONE, camt087LogHandler)
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
		camt087LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, camt087LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", camt087LogHandler, &op.fundHandler, statusData)
		report := parse.CreateCamt030(BIC, originalReqMsgId, xmlMsgType, target, constant.STATUS_CODE_INTERNAL_ERROR)
		return report, whiteListErr
	}

	if pKey == "" {
		errMsg := "OFI can not find RFI in whitelist and vice versa"
		LOGGER.Errorf(errMsg)
		camt087LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, camt087LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", camt087LogHandler, &op.fundHandler, statusData)
		report := parse.CreateCamt030(BIC, originalReqMsgId, xmlMsgType, target, constant.STATUS_CODE_RFI_OR_OFI_NOT_IN_WL)
		return report, whiteListErr
	}
	LOGGER.Infof("Yes, RFI is in OFI's whitelist and vice versa")

	// Parse the pacs008 message with signature into ProtoBuffer
	LOGGER.Infof("Start parsing Go struct to ProtoBuffer")
	protoBufData, parseErr := proto.Marshal(&camt087.SendPayload)
	if parseErr != nil {
		errMsg := "Parse data to ProtoBuf error: " + parseErr.Error()
		LOGGER.Errorf(errMsg)
		camt087LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, camt087LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", camt087LogHandler, &op.fundHandler, statusData)
		report := parse.CreateCamt030(BIC, originalReqMsgId, xmlMsgType, target, constant.STATUS_CODE_INTERNAL_ERROR)
		return report, parseErr
	}
	LOGGER.Infof("Finished parsing Go struct to ProtoBuffer")

	camt087LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_SUCCESS)
	//save the instruction id of camt087 for pacs004/camt029 msg to use
	dbData := sendmodel.DBData{
		MessageId: string(*structData.Body.Assgnmt.Id),
	}

	dbDataByte, _ := json.Marshal(dbData)
	base64DBData := parse.EncodeBase64(dbDataByte)

	// Add the transaction status into the FireBase
	database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, base64DBData, constant.DATABASE_STATUS_PENDING, constant.DATABASE_STATUS_NONE, camt087LogHandler)
	database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", camt087LogHandler, &op.fundHandler, statusData)

	// Send the ProtoBuffer to the request topic of RFI on Kafka broker
	LOGGER.Infof("Start to send request to Kafka broker")
	kafkaErr := op.KafkaActor.Produce(xmlData.RFIId+kafka.REQUEST_TOPIC, protoBufData)
	if kafkaErr != nil {
		errMsg := "Error while submit message to Kafka broker: " + kafkaErr.Error()
		LOGGER.Errorf(errMsg)
		camt087LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, xmlData.InstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_NONE, camt087LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, xmlData.OriginalMsgId, xmlData.OriginalInstructionId, xmlData.InstructionId, "", "", camt087LogHandler, &op.fundHandler, statusData)
		report := parse.CreateCamt030(BIC, originalReqMsgId, xmlMsgType, target, constant.STATUS_CODE_INTERNAL_ERROR)
		return report, kafkaErr
	}
	LOGGER.Infof("Successfully produce message to Kafka broker")
	LOGGER.Debug("-----------------------------------------------------------------")

	// Send status back to OFI
	report := parse.CreateCamt030(BIC, originalReqMsgId, xmlMsgType, target, constant.STATUS_CODE_OFI_SEND_TO_KAFKA)

	return report, nil
}

// if message type is camt.087
func RFI_Camt087(data camt087pbstruct.SendPayload, op *kafka.KafkaOpreations) {
	// Parse the ProtoBuffer into Go struct and reconstruct it into pacs008 message
	LOGGER.Infof("Parsing ProtoBuffer to XML")
	standardType := constant.ISO20022
	paymentStatusMsgType := constant.PAYMENT_TYPE_EXCEPTION
	msgName := constant.CAMT087
	pacs008InstructionId := data.OriginalInstructionId
	instructionId := data.InstructionId
	reqMsgType := data.MsgType
	ofiId := data.OfiId
	rfiId := data.RfiId
	originalMsgId := data.OriginalMsgId
	msgId := data.MsgId

	pacs008LogHandler := transaction.InitiatePaymentLogOperation()
	camt087LogHandler := transaction.InitiatePaymentLogOperation()

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
		find camt087 record from DB
	*/
	dynamoData, paymentInfo := parse.GetDBData(instructionId)
	if dynamoData == nil || paymentInfo == nil {
		LOGGER.Errorf("The original message ID %v does not exist in DB", instructionId)
		op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_WRONG_INSTRUCTION_ID, originalGrpInf)
		return
	}
	camt087LogHandler.PaymentStatuses = paymentInfo
	camt087LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RFI_PROCESSING)

	/*
		verify signature
	*/
	LOGGER.Infof("Verifying the signature from OFI")
	result := xmldsig.VerifySignature(string(data.Message))
	if !result {
		LOGGER.Errorf("signature verification failed")
		camt087LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RFI_VALIDATION_FAIL)
		//camt087 status
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt087LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, pacs008InstructionId, instructionId, "", "", camt087LogHandler, &op.FundHandler, &sendmodel.StatusData{})
		op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_OFI_SIGNATURE_FAIL, originalGrpInf)
		return
	}
	LOGGER.Infof("OFI signature verified!")

	/*
		constructing protobuffer to go struct
	*/
	camt087 := &message_converter.Camt087{SendPayload: data}
	xmlData, err := camt087.ProtobuftoStruct()

	if xmlData == nil {
		LOGGER.Errorf("Encounter error while construncting proto buffer to go struct")
		camt087LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RFI_VALIDATION_FAIL)
		//camt087 status
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt087LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, pacs008InstructionId, instructionId, "", "", camt087LogHandler, &op.FundHandler, &sendmodel.StatusData{})
		op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR_PARSE, originalGrpInf)
		return
	} else if err != nil {
		LOGGER.Errorf("Parse request from kafka failed: %s", err.Error())
		camt087LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RFI_VALIDATION_FAIL)
		//camt087 status
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, camt087LogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, pacs008InstructionId, instructionId, "", "", camt087LogHandler, &op.FundHandler, &sendmodel.StatusData{})
		op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR_PARSE, originalGrpInf)
		return
	}
	// Get important data from the XML data
	reqData := xmlData.RequestXMLMsg
	originalMessageId := xmlData.OriginalMsgId
	originalInstructionId := xmlData.OriginalInstructionId

	LOGGER.Infof("Finished paring ProtoBuffer to XML")

	// Generate payment status data
	// Aggregate necessary data for transaction memo
	statusData := &sendmodel.StatusData{
		IdCdtr:                rfiId,
		IdDbtr:                ofiId,
		CurrencyCode:          camt087.Message.Body.Undrlyg.Initn.OrgnlInstdAmt.Currency,
		InstructionID:         xmlData.InstructionId,
		OriginalInstructionID: xmlData.OriginalInstructionId,
	}

	// update pacs008 record in DB
	rfiVerifyRequestAndSendToKafka(topicName, msgId, msgName, originalMessageId, ofiId, constant.EMPTY_STRING, standardType, msgName, instructionId, originalInstructionId, paymentStatusMsgType, pacs008LogHandler, reqData, statusData, pacs008DynamoData, op, originalGrpInf)

	//update camt087 record in DB
	rawMsg, _ := json.Marshal(dynamoData)
	camt087LogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_MODIFY_PAYMENT)
	go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, parse.EncodeBase64(rawMsg), constant.DATABASE_STATUS_DONE, constant.DATABASE_STATUS_NONE, camt087LogHandler)
	go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, paymentStatusMsgType, msgName, originalMessageId, originalInstructionId, instructionId, "", "", camt087LogHandler, &op.FundHandler, statusData)

	return
}

func getCriticalInfoFromCAMT087(document *camt087struct.RequestToModifyPaymentV06, homeDomain string) (*sendmodel.XMLData, *sendmodel.StatusData, error) {
	ofiId := string(*document.Assgnmt.Assgnr.Agt.FinInstnId.Othr.Id)
	rfiId := string(*document.Assgnmt.Assgne.Agt.FinInstnId.Othr.Id)
	instrId := string(*document.Assgnmt.Id)
	originalPaymentMsgId := string(*document.Undrlyg.Initn.OrgnlGrpInf.OrgnlMsgId)
	originalPaymentType := string(*document.Undrlyg.Initn.OrgnlGrpInf.OrgnlMsgNmId)
	originalInstrId := string(*document.Undrlyg.Initn.OrgnlInstrId)

	checkData := &sendmodel.XMLData{
		OriginalMsgId:         originalPaymentMsgId,
		MessageId:             instrId,
		OFIId:                 ofiId,
		RFIId:                 rfiId,
		ErrorCode:             constant.STATUS_CODE_DEFAULT,
		InstructionId:         instrId,
		OriginalInstructionId: originalInstrId,
	}

	if !utils.StringsEqual(ofiId, homeDomain) {
		LOGGER.Error("Assigner is an incorrect participant")
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_FI
		return checkData, &sendmodel.StatusData{}, errors.New("instructing agent is an incorrect participant")
	}

	// Check if the original payment type is pacs.008
	if !utils.StringsEqual(originalPaymentType, constant.PACS008) {
		LOGGER.Error("Incorrect original message name ID")
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_MSG_NAME_ID
		return checkData, &sendmodel.StatusData{}, errors.New("incorrect original message name ID")
	}

	dbData, txStatus, _, paymentInfoBase64, dbErr := database.DC.GetTransactionData(originalInstrId)

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

	info, _ := parse.DecodeBase64(*paymentInfoBase64)

	var paymentInfo []model.TransactionReceipt
	json.Unmarshal(info, &paymentInfo)

	//check if camt.026 already happened, if it is not, then camt.087 should not happen
	unableToApply := false
	for _, elem := range paymentInfo {
		if elem.Transactionstatus == nil {
			continue
		}
		if *elem.Transactionstatus == constant.PAYMENT_STATUS_UNABLE_TO_APPLY {
			unableToApply = true
			break
		}
	}
	if !unableToApply {
		checkData.ErrorCode = constant.STATUS_CODE_UNABLE_TO_APPLY_NOT_INIT
		return checkData, &sendmodel.StatusData{}, errors.New("WW hasn't received unable to apply request yet")
	}

	/*
		Aggregate necessary data for transaction memo
	*/
	statusData := &sendmodel.StatusData{
		IdCdtr:        rfiId,
		IdDbtr:        ofiId,
		CurrencyCode:  document.Undrlyg.Initn.OrgnlInstdAmt.Currency,
		InstructionID: instrId,
	}

	return checkData, statusData, nil
}

// Participant will use operating account to sign and verify the transaction, if the settlement method is DA
// Otherwise, use issuing account.
func getCriticalInfoFromCamt026(document *camt026struct.Message, homeDomain string) (*sendmodel.XMLData, *sendmodel.StatusData, error) {
	rfiId := string(*document.Body.Assgnmt.Assgnr.Agt.FinInstnId.Othr.Id)
	ofiId := string(*document.Body.Assgnmt.Assgne.Agt.FinInstnId.Othr.Id)
	rfiBic := string(*document.Body.Assgnmt.Assgnr.Agt.FinInstnId.BICFI)
	ofiBic := string(*document.Body.Assgnmt.Assgne.Agt.FinInstnId.BICFI)

	instrId := string(*document.Body.Assgnmt.Id)
	originalInstrId := string(*document.Body.Undrlyg.Initn.OrgnlInstrId)

	checkData := &sendmodel.XMLData{
		OFIBIC:                ofiBic,
		MessageId:             instrId,
		OFIId:                 ofiId,
		RFIId:                 rfiId,
		ErrorCode:             constant.STATUS_CODE_DEFAULT,
		InstructionId:         instrId,
		OriginalInstructionId: originalInstrId,
	}

	if !utils.StringsEqual(rfiId, homeDomain) {
		LOGGER.Error("Instructing agent is an incorrect participant")
		checkData.ErrorCode = constant.STATUS_CODE_WRONG_FI
		return checkData, &sendmodel.StatusData{}, errors.New("instructing agent is an incorrect participant")
	}

	_, txStatus, _, _, dbErr := database.DC.GetTransactionData(originalInstrId)

	if dbErr != nil {
		checkData.ErrorCode = constant.STATUS_CODE_INTERNAL_ERROR
		return checkData, &sendmodel.StatusData{}, errors.New("database query error")
	}

	if *txStatus != constant.DATABASE_STATUS_CLEARED && *txStatus != constant.DATABASE_STATUS_SETTLED {
		checkData.ErrorCode = constant.STATUS_CODE_ORIGINAL_REQUEST_NOT_DONE
		return checkData, &sendmodel.StatusData{}, errors.New("original payment request is not SETTLED or CLEARED yet")
	}

	/*
		Aggregate necessary data for transaction memo
	*/

	amount, _ := strconv.ParseFloat(document.Body.Undrlyg.Initn.OrgnlInstdAmt.Value, 64)
	statusData := &sendmodel.StatusData{
		IdCdtr:                rfiId,
		IdDbtr:                ofiId,
		BICCdtr:               rfiBic,
		BICDbtr:               ofiBic,
		CurrencyCode:          document.Body.Undrlyg.Initn.OrgnlInstdAmt.Currency,
		AmountSettlement:      amount,
		OriginalInstructionID: originalInstrId,
		InstructionID:         instrId,
	}

	return checkData, statusData, nil
}
