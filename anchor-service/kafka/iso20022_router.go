// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package kafka

import (
	"errors"

	message_handler "github.com/GFTN/gftn-services/utility/payment/message-handler"

	"github.com/GFTN/gftn-services/utility/payment/constant"
	message_converter "github.com/GFTN/gftn-services/utility/payment/message-converter"
	"github.com/GFTN/gftn-services/utility/payment/utils/parse"
)

// Route to different ISO20022 message handler base on the message type
func iso20022Router(data []byte, BIC, messageType, target string, op message_handler.PaymentOperations) ([]byte, error) {

	var err error
	var report []byte
	report, err = message_handler.Iso20022Validator(data, BIC, messageType, target)
	if err != nil {
		return report, err
	}

	LOGGER.Infof("Receiving message type: %v", messageType)

	switch messageType {
	case constant.PACS002:
		var messageInstance message_converter.MessageInterface = &message_converter.Pacs002{Raw: data}
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
		pacs002Instance := messageInstance.(*message_converter.Pacs002)
		LOGGER.Infof("DB data successfully constructed")
		report, err = op.Pacs002(*pacs002Instance, target, BIC)

	default:
		errMsg := "unknown message format"
		LOGGER.Error(errMsg)
		err = errors.New(errMsg)
	}

	return report, err
}
