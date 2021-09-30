// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package handler

import (
	"errors"
	"os"
	"strconv"

	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	message_handler "github.com/GFTN/gftn-services/utility/payment/message-handler"

	"github.com/GFTN/gftn-services/utility/payment/environment"

	"github.com/GFTN/gftn-services/utility/payment/constant"
	message_converter "github.com/GFTN/gftn-services/utility/payment/message-converter"
	"github.com/GFTN/gftn-services/utility/payment/utils/parse"
)

// Route to different ISO20022 message handler base on the message type
func iso20022Router(data []byte, messageType string, op message_handler.PaymentOperations, source string) ([]byte, error) {

	var err error
	var report []byte
	BIC := os.Getenv(environment.ENV_KEY_PARTICIPANT_BIC)
	target := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	report, err = message_handler.Iso20022Validator(data, BIC, messageType, target)
	if err != nil {
		return report, err
	}

	LOGGER.Infof("Receiving message type: %v", messageType)
	// request endpoint
	if source == constant.REQUEST {
		switch messageType {
		case constant.PACS008:
			// if the message type is pacs008, execute payment send flow
			var messageInstance message_converter.MessageInterface = &message_converter.Pacs008{Raw: data}
			err = messageInstance.RequestToStruct()
			if err != nil {
				LOGGER.Errorf("Constructing to go struct failed: %v", err.Error())
				report = parse.CreatePacs002(BIC, "", target, constant.STATUS_CODE_PARSE_FAIL, nil)
				return report, err
			}

			err = messageInstance.SanityCheck(BIC, target)
			if err != nil {
				LOGGER.Errorf("Message payload validation failed: %v", err.Error())
				report = parse.CreatePacs002(BIC, "", target, constant.STATUS_CODE_REQ_VALIDATE_FAIL, nil)
				return report, err
			}

			parseErr := messageInstance.StructToProto()
			if parseErr != nil {
				errMsg := "Parse XML to Go struct error: " + parseErr.Error()
				LOGGER.Error(errMsg)
				report = parse.CreatePacs002(BIC, "", target, constant.STATUS_CODE_PARSE_FAIL, nil)
				return report, parseErr
			}

			LOGGER.Infof("Constructing DB data model..")
			pacs008Instance := messageInstance.(*message_converter.Pacs008)
			LOGGER.Infof("DB data successfully constructed")
			report, err = op.Pacs008(*pacs008Instance)

		case constant.CAMT056:
			// if the message type is camt056, execute payment return flow
			var messageInstance message_converter.MessageInterface = &message_converter.Camt056{Raw: data}

			err = messageInstance.RequestToStruct()
			if err != nil {
				LOGGER.Errorf("Constructing to go struct failed: %v", err.Error())
				report = parse.CreateCamt030(BIC, "", constant.CAMT056, target, constant.STATUS_CODE_PARSE_FAIL)
				return report, err
			}

			err = messageInstance.SanityCheck(BIC, target)
			if err != nil {
				LOGGER.Errorf("Message payload validation failed: %v", err.Error())
				report = parse.CreateCamt030(BIC, "", constant.CAMT056, target, constant.STATUS_CODE_REQ_VALIDATE_FAIL)
				return report, err
			}

			parseErr := messageInstance.StructToProto()
			if parseErr != nil {
				errMsg := "Parse XML to Go struct error: " + parseErr.Error()
				LOGGER.Error(errMsg)
				report = parse.CreateCamt030(BIC, "", constant.CAMT056, target, constant.STATUS_CODE_PARSE_FAIL)
				return report, parseErr
			}

			camt056Instance := messageInstance.(*message_converter.Camt056)
			report, err = op.Camt056(*camt056Instance)
		case constant.IBWF002:
			// if the message type is camt056, execute return of digital obligation flow
			var messageInstance message_converter.MessageInterface = &message_converter.Ibwf002{Raw: data}
			err = messageInstance.RequestToStruct()
			if err != nil {
				LOGGER.Errorf("Constructing to go struct failed: %v", err.Error())
				report = parse.CreatePacs002(BIC, "", target, constant.STATUS_CODE_PARSE_FAIL, nil)
				return report, err
			}

			err = messageInstance.SanityCheck(BIC, target)
			if err != nil {
				LOGGER.Errorf("Message payload validation failed: %v", err.Error())
				report = parse.CreatePacs002(BIC, "", target, constant.STATUS_CODE_REQ_VALIDATE_FAIL, nil)
				return report, err
			}

			parseErr := messageInstance.StructToProto()
			if parseErr != nil {
				errMsg := "Parse XML to Go struct error: " + parseErr.Error()
				LOGGER.Error(errMsg)
				report = parse.CreatePacs002(BIC, "", target, constant.STATUS_CODE_PARSE_FAIL, nil)
				return report, parseErr
			}

			LOGGER.Infof("Constructing DB data model..")
			ibwf002Instance := messageInstance.(*message_converter.Ibwf002)
			report, err = op.Ibwf002(*ibwf002Instance)
		case constant.CAMT026:
			// if the message type is camt026, execute payment unable to apply flow
			var messageInstance message_converter.MessageInterface = &message_converter.Camt026{Raw: data}

			err = messageInstance.RequestToStruct()
			if err != nil {
				LOGGER.Errorf("Constructing to go struct failed: %v", err.Error())
				report = parse.CreateCamt030(BIC, "", target, constant.CAMT026, constant.STATUS_CODE_PARSE_FAIL)
				return report, err
			}
			err = messageInstance.SanityCheck(BIC, target)
			if err != nil {
				LOGGER.Errorf("Message payload validation failed: %v", err.Error())
				report = parse.CreateCamt030(BIC, "", target, constant.CAMT026, constant.STATUS_CODE_REQ_VALIDATE_FAIL)
				return report, err
			}

			parseErr := messageInstance.StructToProto()
			if parseErr != nil {
				errMsg := "Parse XML to Go struct error: " + parseErr.Error()
				LOGGER.Error(errMsg)
				report = parse.CreateCamt030(BIC, "", target, constant.CAMT026, constant.STATUS_CODE_PARSE_FAIL)
				return report, parseErr
			}
			camt026Instance := messageInstance.(*message_converter.Camt026)
			report, err = op.Camt026(*camt026Instance)
		case constant.CAMT087:
			// if the message type is camt087, execute request to modify payment flow
			var messageInstance message_converter.MessageInterface = &message_converter.Camt087{Raw: data}

			err = messageInstance.RequestToStruct()
			if err != nil {
				LOGGER.Errorf("Constructing to go struct failed: %v", err.Error())
				report = parse.CreateCamt030(BIC, "", target, constant.CAMT087, constant.STATUS_CODE_PARSE_FAIL)
				return report, err
			}
			err = messageInstance.SanityCheck(BIC, target)
			if err != nil {
				LOGGER.Errorf("Message payload validation failed: %v", err.Error())
				report = parse.CreateCamt030(BIC, "", target, constant.CAMT087, constant.STATUS_CODE_REQ_VALIDATE_FAIL)
				return report, err
			}

			parseErr := messageInstance.StructToProto()
			if parseErr != nil {
				errMsg := "Parse XML to Go struct error: " + parseErr.Error()
				LOGGER.Error(errMsg)
				report = parse.CreateCamt030(BIC, "", target, constant.CAMT087, constant.STATUS_CODE_PARSE_FAIL)
				return report, parseErr
			}
			camt087Instance := messageInstance.(*message_converter.Camt087)
			report, err = op.Camt087(*camt087Instance)
		default:
			errMsg := "unknown message format"
			LOGGER.Error(errMsg)
			err = errors.New(errMsg)
		}
		// response endpoint
	} else if source == constant.RESPONSE {
		switch messageType {
		case constant.IBWF001:
			// if message type is ibwf001, execute the payment reply flow
			var messageInstance message_converter.MessageInterface = &message_converter.Ibwf001{Raw: data}

			parseErr := messageInstance.RequestToStruct()
			if parseErr != nil {
				errMsg := "Parse XML to Go struct error: " + parseErr.Error()
				LOGGER.Error(errMsg)
				report = parse.CreatePacs002(BIC, "", target, constant.STATUS_CODE_PARSE_FAIL, nil)
				return report, parseErr
			}
			err = messageInstance.SanityCheck(BIC, target)
			if err != nil {
				LOGGER.Errorf("Message payload validation failed: %v", err.Error())
				report = parse.CreatePacs002(BIC, "", target, constant.STATUS_CODE_REQ_VALIDATE_FAIL, nil)
				return report, err
			}

			parseErr = messageInstance.StructToProto()
			if parseErr != nil {
				errMsg := "Parse XML to Go struct error: " + parseErr.Error()
				LOGGER.Error(errMsg)
				report = parse.CreatePacs002(BIC, "", target, constant.STATUS_CODE_PARSE_FAIL, nil)
				return report, parseErr
			}

			ibwf001Instance := messageInstance.(*message_converter.Ibwf001)
			report, err = op.Ibwf001(*ibwf001Instance)

		case constant.PACS004:
			// if the message type is pacs004, execute acknowledge payment return flow
			var messageInstance message_converter.MessageInterface = &message_converter.Pacs004{Raw: data}

			err = messageInstance.RequestToStruct()
			if err != nil {
				LOGGER.Errorf("Constructing to go struct failed: %v", err.Error())
				report = parse.CreatePacs002(BIC, "", target, constant.STATUS_CODE_PARSE_FAIL, nil)
				return report, err
			}
			err = messageInstance.SanityCheck(BIC, target)
			if err != nil {
				LOGGER.Errorf("Message payload validation failed: %v", err.Error())
				report = parse.CreatePacs002(BIC, "", target, constant.STATUS_CODE_REQ_VALIDATE_FAIL, nil)
				return report, err
			}

			pacs004Instance := messageInstance.(*message_converter.Pacs004)
			code, _ := strconv.Atoi(string(*pacs004Instance.Message.Body.OrgnlGrpInf.RtrRsnInf[0].Rsn.Cd))

			switch code {
			case constant.REASON_CODE_PAYMENT_CANCELLATION:
				report, err = op.Pacs004_Cancellation(*pacs004Instance)
			case constant.REASON_CODE_RDO:
				report, err = op.Pacs004_Rdo(*pacs004Instance)
			default:
				err = errors.New("Cannot identify the reason code inside the OrgnlGrpInf tag")
				LOGGER.Errorf(err.Error())
				report = parse.CreatePacs002(BIC, "", target, constant.STATUS_CODE_PARSE_FAIL, nil)
			}

		case constant.CAMT029:
			// if the message type is pacs004, execute rejection of payment return flow
			var messageInstance message_converter.MessageInterface = &message_converter.Camt029{Raw: data}

			err = messageInstance.RequestToStruct()
			if err != nil {
				LOGGER.Errorf("Constructing to go struct failed: %v", err.Error())
				report = parse.CreateCamt030(BIC, "", target, constant.CAMT029, constant.STATUS_CODE_PARSE_FAIL)
				return report, err
			}
			err = messageInstance.SanityCheck(BIC, target)
			if err != nil {
				LOGGER.Errorf("Message payload validation failed: %v", err.Error())
				report = parse.CreateCamt030(BIC, "", target, constant.CAMT029, constant.STATUS_CODE_REQ_VALIDATE_FAIL)
				return report, err
			}

			parseErr := messageInstance.StructToProto()
			if parseErr != nil {
				errMsg := "Parse XML to Go struct error: " + parseErr.Error()
				LOGGER.Error(errMsg)
				report = parse.CreateCamt030(BIC, "", target, constant.CAMT029, constant.STATUS_CODE_PARSE_FAIL)
				return report, parseErr
			}
			camt029Instance := messageInstance.(*message_converter.Camt029)
			report, err = op.Camt029(*camt029Instance)
		default:
			errMsg := "unknown message format"
			LOGGER.Error(errMsg)
			err = errors.New(errMsg)
		}
	} else if source == constant.REDEEM {
		switch messageType {
		case constant.PACS009:
			// if the message type is pacs009, execute digital asset redemption flow
			var messageInstance message_converter.MessageInterface = &message_converter.Pacs009{Raw: data}
			err = messageInstance.RequestToStruct()
			if err != nil {
				LOGGER.Errorf("Constructing to go struct failed: %v", err.Error())
				report = parse.CreatePacs002(BIC, "", target, constant.STATUS_CODE_PARSE_FAIL, nil)
				return report, err
			}

			err = messageInstance.SanityCheck(BIC, target)
			if err != nil {
				LOGGER.Errorf("Message payload validation failed: %v", err.Error())
				report = parse.CreatePacs002(BIC, "", target, constant.STATUS_CODE_REQ_VALIDATE_FAIL, nil)
				return report, err
			}

			parseErr := messageInstance.StructToProto()
			if parseErr != nil {
				errMsg := "Parse XML to Go struct error: " + parseErr.Error()
				LOGGER.Error(errMsg)
				report = parse.CreatePacs002(BIC, "", target, constant.STATUS_CODE_PARSE_FAIL, nil)
				return report, parseErr
			}

			LOGGER.Infof("Constructing DB data model..")
			pacs009Instance := messageInstance.(*message_converter.Pacs009)
			LOGGER.Infof("DB data successfully constructed")
			report, err = op.Pacs009(*pacs009Instance)

		default:
			errMsg := "unknown message format"
			LOGGER.Error(errMsg)
			err = errors.New(errMsg)
		}
		// response endpoint
	}

	return report, err
}
