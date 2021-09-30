// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package parse

import (
	"errors"
	"os"

	"github.com/GFTN/iso20022/pacs00200109"
	"github.com/GFTN/gftn-services/utility/payment/constant"
	"github.com/GFTN/gftn-services/utility/payment/environment"
	"github.com/GFTN/gftn-services/utility/payment/utils"
)

func KafkaErrorRouter(xmlMsgType, instructionId, ofiId, rfiId string, statusCode int, generateReport bool, originalGrpInf *pacs00200109.OriginalGroupInformation29) (string, []byte, error) {
	var report []byte
	var targetParticipant string
	BIC := os.Getenv(environment.ENV_KEY_PARTICIPANT_BIC)

	switch xmlMsgType {
	case constant.PACS008:
		targetParticipant = ofiId
	case constant.IBWF001:
		targetParticipant = rfiId
	case constant.CAMT056:
		targetParticipant = ofiId
	case constant.PACS004:
		targetParticipant = rfiId
	case constant.CAMT029:
		targetParticipant = rfiId
	case constant.IBWF002:
		targetParticipant = ofiId
	case constant.PACS009:
		targetParticipant = ofiId
	case constant.PACS002:
		targetParticipant = rfiId
	case constant.CAMT026:
		targetParticipant = rfiId
	case constant.CAMT087:
		targetParticipant = ofiId
	default:
		LOGGER.Errorf("No matching XML message type found")
		return "", nil, errors.New("No matching XML message type found")
	}

	if generateReport {
		if utils.Contains(constant.SUPPORT_CAMT_MESSAGES, xmlMsgType) {
			report = CreateCamt030(BIC, instructionId, xmlMsgType, targetParticipant, statusCode)
		} else {
			report = CreatePacs002(BIC, instructionId, targetParticipant, statusCode, originalGrpInf)
		}
	}
	return targetParticipant, report, nil
}
