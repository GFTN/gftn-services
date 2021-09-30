// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package message_handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	pacs002pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/pacs00200109"

	pacs002struct "github.com/GFTN/iso20022/pacs00200109"
	blocklist_client "github.com/GFTN/gftn-services/administration-service/blocklist-client"
	"github.com/GFTN/gftn-services/anchor-service/handlers"
	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/utility/common"
	"github.com/GFTN/gftn-services/utility/kafka"
	"github.com/GFTN/gftn-services/utility/payment/client"
	"github.com/GFTN/gftn-services/utility/payment/environment"
	message_converter "github.com/GFTN/gftn-services/utility/payment/message-converter"
	"github.com/GFTN/gftn-services/utility/payment/utils"
	"github.com/GFTN/gftn-services/utility/payment/utils/signing"
	"github.com/GFTN/gftn-services/utility/xmldsig"

	"os"

	"github.com/GFTN/gftn-services/utility/payment/constant"

	"github.com/go-openapi/strfmt"
	"github.com/golang/protobuf/proto"
	"github.com/lestrrat-go/libxml2/xsd"
	"github.com/op/go-logging"
	whitelist_handler "github.com/GFTN/gftn-services/utility/payment/utils/whitelist-handler"

	"github.com/GFTN/gftn-services/utility/payment/utils/database"

	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/payment/utils/parse"
	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"
	"github.com/GFTN/gftn-services/utility/payment/utils/transaction"
)

var LOGGER = logging.MustGetLogger("message-handler")

type PaymentOperations struct {
	xsdPath          string
	XsdSchemas       []*xsd.Schema
	homeDomain       string
	KafkaActor       *kafka.KafkaOpreations
	whitelistHandler whitelist_handler.ParticipantWhiteList
	signHandler      signing.CreateSignOperations
	statusFileName   string
	input            *parse.EndPointInput
	sendVars         sendmodel.SendVariables
	fundHandler      transaction.CreateFundingOpereations
	prServiceURL     string
}

func InitiatePaymentOperations() (PaymentOperations, error) {
	op := PaymentOperations{}
	op.homeDomain = os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	op.statusFileName = os.Getenv(environment.ENV_KEY_SERVICE_FILE)
	op.prServiceURL = os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL)
	if _, exists := os.LookupEnv(environment.ENV_KEY_WW_ID); !exists {
		panic("environment variable WW_ID is empty")
		return PaymentOperations{}, errors.New("environment variable WW_ID is empty")
	} else if _, exists := os.LookupEnv(environment.ENV_KEY_WW_BIC); !exists {
		panic("environment variable WW_BIC is empty")
		return PaymentOperations{}, errors.New("environment variable WW_BIC is empty")
	}

	// Initialize funding operations such as gas-service-client, rdo-service-client, admin-service-client, pr-service-client
	op.fundHandler = transaction.InitiateFundingOperations(op.prServiceURL, op.homeDomain)

	// Initialize whitelist-service-client
	op.whitelistHandler = whitelist_handler.CreateWhiteListServiceOperations()

	// Initialize crypto-service-client
	op.signHandler = signing.InitiateSignOperations(op.prServiceURL)

	// Get participant's BIC code
	if !utils.StringsEqual(op.homeDomain, os.Getenv(environment.ENV_KEY_WW_ID)) {
		participantBIC := client.GetParticipantAccount(op.prServiceURL, op.homeDomain, constant.BIC_STRING)
		os.Setenv(environment.ENV_KEY_PARTICIPANT_BIC, *participantBIC)
	}
	// Initializing Kafka consumer which will subscribe to two specific topic
	// [participant_id]_req for incoming request from other participants
	// [participant_id]_res for incoming response from other participants
	LOGGER.Infof("Initiate Kafka producer")
	//kafkaActor, initKafkaErr := kafkaHandler.InitiateKafkaOperation(op.homeDomain, *op.ParticipantBIC, op.fundHandler, op.signHandler)
	kafkaActor, initKafkaErr := kafka.Initialize()
	if initKafkaErr != nil {
		panic("Initialize Kafka producer failed: " + initKafkaErr.Error())
		return PaymentOperations{}, initKafkaErr
	}

	op.KafkaActor = kafkaActor

	// Initialize DynamoDB connection
	LOGGER.Infof("Initiate DynamoDB")
	database.DC = database.DynamoClient{
		Region: os.Getenv(environment.ENV_KEY_DYNAMO_DB_REGION),
	}
	database.DC.CreateConnection()

	// Read the variables files which contained the response status code and xml message type etc.
	LOGGER.Infof("Loading kafka status file")
	op.input = &parse.V
	jsonFile, openErr := os.Open(op.statusFileName)
	if openErr != nil {
		panic("Kafka status file not found")
		return PaymentOperations{}, openErr
	}
	byteData, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteData, &op.sendVars)
	parse.V.Vars = op.sendVars

	// Initialize the random number seed for message ID used by response xml
	parse.Init()

	// Setup all necessary XML schema for future validation
	LOGGER.Infof("Setting up XML schema at %v", op.sendVars.XSDPath)
	xsdPaths := op.sendVars.XSDPath
	schemas, err := parse.SchemaInitiate(xsdPaths)
	if err != nil {
		panic("Unable to initiate xsds")
		return PaymentOperations{}, err
	}
	op.XsdSchemas = schemas

	return op, nil
}

// func (op *CreateSendOperations) RecordTimeLogsToKafka(action, task string, timeStamp time.Time) {
// 	logs := fmt.Sprintf("%s-[%s]:[%s][%s][Time:%s]", action, op.homeDomain, op.KafkaActor.ServiceName, task, time.Since(timeStamp).String())
// 	op.KafkaActor.SendServiceLogs(constant.SERVICE_LOG_TOPIC, []byte(logs))
// }

func rfiVerifyRequestAndSendToKafka(
	topicName, msgId, msgName, originalMsgId, ofiId, settlementAccountName, standardType, reqMsgType, instructionId, originalInstructionId, paymentStatusMsgType string,
	originalLogHandler transaction.Payment,
	reqData []byte,
	statusData *sendmodel.StatusData,
	dynamoData interface{},
	op *kafka.KafkaOpreations,
	originalGrpInf *pacs002struct.OriginalGroupInformation29) {

	rfiId := statusData.IdCdtr

	if len(reqMsgType) < 2 {
		LOGGER.Errorf("Error message type")
		originalLogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RFI_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, originalLogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, originalInstructionId, instructionId, "", "", originalLogHandler, &op.FundHandler, statusData)
		op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR_PARSE, originalGrpInf)
		return
	}
	// Aggregate necessary data for transaction memo
	commonStatusData := &sendmodel.StatusData{
		IdCdtr:                statusData.IdCdtr,
		IdDbtr:                statusData.IdDbtr,
		EndToEndID:            statusData.EndToEndID,
		InstructionID:         statusData.InstructionID,
		OriginalInstructionID: statusData.OriginalInstructionID,
	}

	LOGGER.Infof("Sending %v message to Kafka", reqMsgType)

	var successStatus string
	var txStatus string
	switch reqMsgType {
	case constant.PACS008:
		successStatus = constant.PAYMENT_STATUS_RFI_VALIDATION_SUCCESS
		originalLogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RFI_PROCESSING)
		txStatus = constant.DATABASE_STATUS_DONE
	case constant.CAMT056:
		successStatus = constant.PAYMENT_STATUS_CANCELLATION_INIT
		txStatus = constant.DATABASE_STATUS_CANCEL_INIT
	case constant.IBWF002:
		successStatus = constant.PAYMENT_STATUS_RDO_INIT
		txStatus = constant.DATABASE_STATUS_RDO_INIT
	case constant.PACS009:
		successStatus = constant.PAYMENT_STATUS_ASSET_REDEMPTION_INIT
		originalLogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RFI_PROCESSING)
		txStatus = constant.DATABASE_STATUS_ASSET_REDEMPTION_INIT
	case constant.CAMT087:
		successStatus = constant.PAYMENT_STATUS_MODIFY_PAYMENT
		txStatus = constant.DATABASE_STATUS_MODIFY_PAYMENT
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
	countries = append(countries, strconv.Quote(statusData.CountryCdtr))
	countries = append(countries, strconv.Quote(statusData.CountryDbtr))

	var currencies []string
	currencies = append(currencies, strconv.Quote(statusData.CurrencyCode))

	var participants []string
	participants = append(participants, strconv.Quote(statusData.IdCdtr))
	participants = append(participants, strconv.Quote(statusData.IdDbtr))

	// validate block-list
	res, err := blockListClient.ValidateFromBlocklist(countries, currencies, participants)
	if err != nil {
		LOGGER.Errorf("%v", err)
		originalLogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, originalLogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, paymentStatusMsgType, msgName, originalMsgId, originalInstructionId, instructionId, "", "", originalLogHandler, &op.FundHandler, commonStatusData)
		op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR, originalGrpInf)
		return
	}
	if res == common.BlocklistDeniedString {
		LOGGER.Errorf("The transaction currency/country/institution is within the blocklist, transaction forbidden!")
		originalLogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RFI_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, originalLogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_BOTH, paymentStatusMsgType, msgName, originalMsgId, originalInstructionId, instructionId, "", "", originalLogHandler, &op.FundHandler, commonStatusData)
		op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_BLOCKLIST, originalGrpInf)
		return
	}

	// Check if OFI was whitelisted by RFI and vice versa, if not, reject the payment request
	whitelistHandler := whitelist_handler.CreateWhiteListServiceOperations()
	LOGGER.Infof("Check whether OFI is in RFI's whitelist and vice versa.")
	pkey, whiteListErr := whitelistHandler.CheckWhiteListParticipant(rfiId, ofiId, settlementAccountName)
	if whiteListErr != nil {
		LOGGER.Errorf(whiteListErr.Error())
		originalLogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, originalLogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, originalInstructionId, instructionId, "", "", originalLogHandler, &op.FundHandler, statusData)
		op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR, originalGrpInf)
		return
	}
	if pkey == "" {
		LOGGER.Errorf("Can not find OFI or RFI in whitelist")
		originalLogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RFI_VALIDATION_FAIL)
		go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, originalLogHandler)
		go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, originalInstructionId, instructionId, "", "", originalLogHandler, &op.FundHandler, statusData)
		op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_RFI_OR_OFI_NOT_IN_WL, originalGrpInf)
		return
	}
	LOGGER.Infof("Yes, OFI is in RFI's whitelist.")

	/*
		Signing message with IBM master account
	*/
	//stronghold exception handling
	if utils.StringsEqual(reqMsgType, constant.PACS009) &&
		utils.StringsEqual(statusData.IdCdtr, os.Getenv(global_environment.ENV_KEY_STRONGHOLD_ANCHOR_ID)) {

		anchorOps, err := handlers.CreateAnchorOperations()

		bankName := statusData.SupplementaryData[constant.PACS009_SUPPLEMENTARY_DATA_BANK_NAME]
		bankRoutingNumber := statusData.SupplementaryData[constant.PACS009_SUPPLEMENTARY_DATA_BRANCH]
		initiatorIp := "127.0.0.1"
		bankAccountNumber := statusData.SupplementaryData[constant.PACS009_SUPPLEMENTARY_DATA_ACCOUNT_NUMBER]
		bankAccountType := statusData.SupplementaryData[constant.PACS009_SUPPLEMENTARY_DATA_ACCOUNT_TYPE]
		paymentMethod := statusData.SupplementaryData[constant.PACS009_SUPPLEMENTARY_DATA_PAYMENT_METHOD]
		withdrawAssetMap := os.Getenv(environment.ENV_KEY_ANCHOR_SH_ASSET_ID)
		customerReference := statusData.InstructionID
		amount := common.FloatToFixedPrecisionString(statusData.AmountSettlement, 2)

		err = model.IsValidDACode(statusData.CurrencyCode)

		if err != nil {
			LOGGER.Debug("WithDraw: asset code is invalid:", err)
			originalLogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_RFI_VALIDATION_FAIL)
			go database.SyncWithDynamo(constant.DATABASE_UPDATE, originalInstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, originalLogHandler)
			go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, originalInstructionId, originalInstructionId, "", "", originalLogHandler, &op.FundHandler, statusData)
			op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_REQ_VALIDATE_FAIL, originalGrpInf)
			return
		}

		assetID := statusData.CurrencyCode + withdrawAssetMap

		strongHoldInput := model.StrongholdWithdrawRequest{
			AssetID:       &assetID,
			PaymentMethod: &paymentMethod,
			PaymentMethodDetails: &model.StrongholdPaymentMethodDetails{
				Amount:            &amount,
				BankName:          &bankName,
				BankRoutingNumber: &bankRoutingNumber,
				InitiatorIP:       &initiatorIp,
				BankAccountNumber: &bankAccountNumber,
				BankAccountType:   &bankAccountType,
			},
			CustomerReference: &customerReference,
		}

		anchorRes, err := anchorOps.WithDraw(strongHoldInput)
		if anchorRes == nil {
			LOGGER.Errorf("Error response from stronghold API: %v", err)
			originalLogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
			go database.SyncWithDynamo(constant.DATABASE_UPDATE, originalInstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, originalLogHandler)
			go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, originalInstructionId, originalInstructionId, "", "", originalLogHandler, &op.FundHandler, statusData)
			op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR, originalGrpInf)
			return
		}
		msg, instructionId, err := constructPacsMessageForSH(statusData, anchorRes)
		if err != nil {
			LOGGER.Errorf("Encounter error while constructing pacs002 msg for stronghold: %v", err.Error())
			originalLogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
			go database.SyncWithDynamo(constant.DATABASE_UPDATE, originalInstructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, originalLogHandler)
			go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, originalInstructionId, originalInstructionId, "", "", originalLogHandler, &op.FundHandler, statusData)
			op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR, originalGrpInf)
			return
		}
		/*
			write into dynamo
		*/
		logHandler := transaction.InitiatePaymentLogOperation()
		// Initialize log handler and set the payment status to `INITIAL`
		logHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_INITIAL)
		err = database.SyncWithDynamo(constant.DATABASE_INIT, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_PENDING, constant.DATABASE_STATUS_NONE, logHandler)
		if err != nil {
			LOGGER.Errorf(err.Error())
			originalLogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
			op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_DUP_ID, originalGrpInf)
			return
		}

		var sendPayload pacs002pbstruct.SendPayload

		sendPayload.Message = msg
		sendPayload.MsgType = constant.ISO20022 + ":" + constant.PACS002
		sendPayload.OfiId = statusData.IdDbtr
		sendPayload.RfiId = statusData.IdCdtr
		sendPayload.InstructionId = instructionId
		sendPayload.OriginalInstructionId = statusData.OriginalInstructionID

		/*
			sanity check
		*/
		var messageInstance message_converter.MessageInterface = &message_converter.Pacs002{Raw: msg}
		err = messageInstance.RequestToStruct()
		if err != nil {
			LOGGER.Errorf("Constructing to go struct failed: %v", err.Error())
			originalLogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_FAIL)
			op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_PARSE_FAIL, originalGrpInf)
			return
		}
		pacs002Instance := messageInstance.(*message_converter.Pacs002)

		xmlData, _ := getCriticalInfoFromPacs002(pacs002Instance.Message, statusData.IdCdtr)
		statusCode := xmlData.ErrorCode
		if statusCode != constant.STATUS_CODE_DEFAULT {
			LOGGER.Errorf("Something wrong with the transaction information")
			originalLogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_FAIL)
			op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_PARSE_FAIL, originalGrpInf)
			return
		}
		/*
			prepare to send to Kafka
		*/
		protoBufData, parseErr := proto.Marshal(&sendPayload)
		if parseErr != nil {
			errMsg := "Parse data to ProtoBuf error: " + parseErr.Error()
			LOGGER.Errorf(errMsg)
			originalLogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
			go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, originalLogHandler)
			go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, originalInstructionId, instructionId, "", "", originalLogHandler, &op.FundHandler, statusData)
			op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_PARSE_FAIL, originalGrpInf)
			return
		}
		LOGGER.Infof("Finished parsing Go struct to ProtoBuffer")

		msgType := constant.PAYMENT_TYPE_ASSET_REDEMPTION
		msgName := constant.PACS009
		logHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_VALIDATION_SUCCESS)
		database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, "", constant.DATABASE_STATUS_PENDING, constant.DATABASE_STATUS_NONE, logHandler)
		database.SyncWithFirebase(constant.FIREBASE_INIT, msgType, msgName, instructionId, originalInstructionId, instructionId, "", "", logHandler, &op.FundHandler, statusData)
		op.Produce(ofiId+kafka.RESPONSE_TOPIC, protoBufData)

		// not stronghold
	} else {
		var gatewayMsg string
		if utils.StringsEqual(reqMsgType, constant.PACS009) {
			signedMessage, signErr := signing.SignPayloadByMasterAccount(string(reqData))
			if signErr != nil {
				LOGGER.Errorf("Failed to sign payload: %v", signErr.Error())
				originalLogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
				go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, originalLogHandler)
				go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, originalInstructionId, instructionId, "", "", originalLogHandler, &op.FundHandler, statusData)
				op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_SIGN_PAYLOAD_FAIL, originalGrpInf)
				return
			}
			gatewayMsg = parse.EncodeBase64([]byte(signedMessage))

		} else {
			signedMessage, signErr := op.SignHandler.SignPayloadByMasterAccount(reqData)
			if signErr != nil {
				LOGGER.Errorf("Failed to sign payload: %v", signErr.Error())
				originalLogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
				go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, originalLogHandler)
				go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, originalInstructionId, instructionId, "", "", originalLogHandler, &op.FundHandler, statusData)
				op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_SIGN_PAYLOAD_FAIL, originalGrpInf)
				return
			}
			gatewayMsg = parse.EncodeBase64(signedMessage)
		}

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

		err = op.Produce(topicName, msg)
		if err != nil {
			LOGGER.Errorf("Encounter error while producing message to Kafka topic: %v", rfiId+"_"+kafka.TRANSACTION_TOPIC)
			originalLogHandler.RecordPaymentStatus(constant.PAYMENT_STATUS_FAILED)
			go database.SyncWithDynamo(constant.DATABASE_UPDATE, instructionId, constant.DATABASE_STATUS_EMPTY, constant.DATABASE_STATUS_FAILED, constant.DATABASE_STATUS_FAILED, originalLogHandler)
			go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, msgName, originalMsgId, originalInstructionId, instructionId, "", "", originalLogHandler, &op.FundHandler, statusData)
			op.SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId, constant.STATUS_CODE_INTERNAL_ERROR, originalGrpInf)
			return
		}
	}

	// Create transaction memo
	//sendmodel.TxMemo = originalLogHandler.BuildTXMemo(statusData, "", msgId, ofiId, paymentStatusMsgType)
	//go kafkaOperations.fundHandler.SendToAdm(originalLogHandler.PaymentStatuses, constant.UPDATE_STATUS, instructionId)

	rawMsg, _ := json.Marshal(dynamoData)
	originalLogHandler.RecordPaymentStatus(successStatus)
	go database.SyncWithDynamo(constant.DATABASE_UPDATE, originalInstructionId, parse.EncodeBase64(rawMsg), txStatus, constant.DATABASE_STATUS_NONE, originalLogHandler)
	go database.SyncWithFirebase(constant.FIREBASE_UPDATE_PARTIAL, paymentStatusMsgType, reqMsgType, originalMsgId, originalInstructionId, originalInstructionId, "", "", originalLogHandler, &op.FundHandler, commonStatusData)
	LOGGER.Debug("--------------------------------------------------------------------")
	return
}

// Sending the error message back to Kafka when there's a error on the other side
func HandleErrMsg(pbsData sendmodel.SendPayload, op *kafka.KafkaOpreations) {
	responseMsg := strings.Split(pbsData.MsgType, ":")
	xmlMsgType := responseMsg[1]
	statusCode, _ := strconv.Atoi(responseMsg[2])
	instrId := pbsData.InstructionId

	originalGrpInf := &pacs002struct.OriginalGroupInformation29{
		OrgnlMsgId:   getReportMax35Text(pbsData.MsgId),
		OrgnlMsgNmId: getReportMax35Text(xmlMsgType),
	}

	LOGGER.Warningf("Receiving error message %s on OFI or RFI side during request handling: %d", pbsData.MsgType, statusCode)
	LOGGER.Infof("Incoming message type: %s", xmlMsgType)

	targetParticipant, report, err := parse.KafkaErrorRouter(xmlMsgType, instrId, pbsData.OfiId, pbsData.RfiId, statusCode, true, originalGrpInf)
	if err != nil {
		return
	}

	op.SendRequestToKafka(targetParticipant+"_"+kafka.TRANSACTION_TOPIC, report)

	return
}

//validate signature & xml payload
func Iso20022Validator(data []byte, bic, messageType, target string) ([]byte, error) {

	/*
		verify signature
	*/

	LOGGER.Infof("Verifying the XML signature")
	originalGrpInf := &pacs002struct.OriginalGroupInformation29{
		OrgnlMsgNmId: getReportMax35Text(messageType),
	}
	var report []byte
	result := xmldsig.VerifySignature(string(data))
	if !result {
		LOGGER.Errorf("signature verification failed")
		if utils.Contains(constant.SUPPORT_CAMT_MESSAGES, messageType) {
			report = parse.CreateCamt030(bic, "", messageType, target, constant.STATUS_CODE_OFI_SIGNATURE_FAIL)
		} else {
			report = parse.CreatePacs002(bic, "", target, constant.STATUS_CODE_OFI_SIGNATURE_FAIL, originalGrpInf)
		}
		return report, errors.New("Failed validating the signature in application header")
	}

	//Parsing the XML document
	LOGGER.Infof("Validating the XML against schema")
	err := parse.ValidateSchema(string(data))

	if err != nil {
		errMsg := "Schema failure error message: " + err.Error()
		LOGGER.Error(errMsg)
		if utils.Contains(constant.SUPPORT_CAMT_MESSAGES, messageType) {
			report = parse.CreateCamt030(bic, "", messageType, target, constant.STATUS_CODE_XML_VALIDATE_FAIL)
		} else {
			report = parse.CreatePacs002(bic, "", target, constant.STATUS_CODE_XML_VALIDATE_FAIL, originalGrpInf)
		}
		return report, err
	}

	return nil, nil
}

//validate incoming request
func ValidateRequest(raw *http.Request, bic, target string) ([]byte, []byte, string, error) {
	var rawMsg model.SendPacs
	var report []byte
	err := json.NewDecoder(raw.Body).Decode(&rawMsg)
	if err != nil {
		LOGGER.Errorf("Error  %v", err.Error())
		report = []byte("Unable to decode incoming request")
		return []byte{}, report, "", err
	}

	err = rawMsg.Validate(strfmt.Default)
	if err != nil {
		msg := "Unable to validate send request: " + err.Error()
		LOGGER.Error(msg)
		report = []byte("Unable to validate incoming request format")
		return []byte{}, report, "", err
	}

	messageType := *rawMsg.MessageType

	originalGrpInf := &pacs002struct.OriginalGroupInformation29{
		OrgnlMsgNmId: getReportMax35Text(messageType),
	}

	if len(strings.Split(*rawMsg.MessageType, ":")) != 2 {
		LOGGER.Error("Invalid message type")
		if utils.Contains(constant.SUPPORT_CAMT_MESSAGES, messageType) {
			report = parse.CreateCamt030(bic, "", messageType, target, constant.STATUS_CODE_MSG_TYPE_VALIDATE_FAIL)
		} else {
			report = parse.CreatePacs002(bic, "", target, constant.STATUS_CODE_MSG_TYPE_VALIDATE_FAIL, originalGrpInf)
		}
		return []byte{}, report, "", errors.New("invalid message type")
	}

	//Check if requested ofi account is the same as account token
	if os.Getenv(global_environment.ENV_KEY_ENABLE_JWT) != constant.FALSE_STRING {

		accountToken, err := middlewares.GetSessionContext(raw)
		if err != nil {
			if utils.Contains(constant.SUPPORT_CAMT_MESSAGES, messageType) {
				report = parse.CreateCamt030(bic, "", messageType, target, constant.STATUS_CODE_JWT_FAIL)
			} else {
				report = parse.CreatePacs002(bic, "", target, constant.STATUS_CODE_JWT_FAIL, originalGrpInf)
			}
			return []byte{}, report, "", errors.New("Get session context fail")
		}

		if accountToken.TimeTill <= 0 {
			// TODO which report type should we use here ?
			if utils.Contains(constant.SUPPORT_CAMT_MESSAGES, messageType) {
				report = parse.CreateCamt030(bic, "", messageType, target, constant.STATUS_CODE_JWT_FAIL)
			} else {
				report = parse.CreatePacs002(bic, "", target, constant.STATUS_CODE_JWT_FAIL, originalGrpInf)
			}
			return []byte{}, report, "", errors.New("JWT token timeout")
		}
	}

	// Decode base64 message
	data, decodeErr := parse.DecodeBase64(*rawMsg.Message)
	if decodeErr != nil {
		msg := "Unable to parse the encoded message: " + decodeErr.Error()
		LOGGER.Error(msg)
		if utils.Contains(constant.SUPPORT_CAMT_MESSAGES, messageType) {
			report = parse.CreateCamt030(bic, "", messageType, target, constant.STATUS_CODE_DECODE_FAIL)
		} else {
			report = parse.CreatePacs002(bic, "", target, constant.STATUS_CODE_DECODE_FAIL, originalGrpInf)
		}
		return []byte{}, report, "", decodeErr
	}

	return data, []byte{}, *rawMsg.MessageType, nil
}
