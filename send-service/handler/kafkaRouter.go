// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package handler

import (
	"strings"

	message_handler "github.com/GFTN/gftn-services/utility/payment/message-handler"

	"github.com/golang/protobuf/proto"
	camt026pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/camt02600107"
	camt029pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/camt02900109"
	camt056pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/camt05600108"
	camt087pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/camt08700106"
	ibwf001pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/ibwf00100101"
	ibwf002pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/ibwf00200101"
	pacs002pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/pacs00200109"
	pacs008pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/pacs00800107"

	"github.com/GFTN/gftn-services/utility/kafka"
	"github.com/GFTN/gftn-services/utility/payment/constant"
	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"
)

/*
	Once the RFI consumed the Payment request ProtoBuffer from the Kafka broker, it will start processing the message
*/
func KafkaRouter(consumeType string, data []byte, op *kafka.KafkaOpreations) {
	var pbs sendmodel.SendPayload
	proto.Unmarshal(data, &pbs)

	LOGGER.Infof("Message Type: %s", pbs.MsgType)
	if len(pbs.MsgType) < 2 {
		LOGGER.Errorf("Error reading message type: %v from Kafka", pbs.MsgType)
		return
	}
	standardType := strings.TrimSpace(strings.ToLower(strings.Split(pbs.MsgType, ":")[0]))
	messageType := strings.TrimSpace(strings.ToLower(strings.Split(pbs.MsgType, ":")[1]))

	switch consumeType {
	case kafka.REQUEST_TOPIC:

		switch standardType {
		case constant.ISO20022:

			switch messageType {
			case constant.PACS008:
				// For requesting a payment transaction (pacs008 message)
				var pbs pacs008pbstruct.SendPayload
				proto.Unmarshal(data, &pbs)
				message_handler.RFI_Pacs008(pbs, op)
				return
			case constant.CAMT056:
				// For requesting a payment cancellation (camt056 message)
				var pbs camt056pbstruct.SendPayload
				proto.Unmarshal(data, &pbs)
				message_handler.RFI_Camt056(pbs, op)
				return
			case constant.IBWF002:
				// For requesting a DO settlement notification (ibwf002 message)
				var pbs ibwf002pbstruct.SendPayload
				proto.Unmarshal(data, &pbs)
				message_handler.RFI_Ibwf002(pbs, op)
				return
			case constant.CAMT026:
				// For requesting unable tp apply request (camt026 message)
				var pbs camt026pbstruct.SendPayload
				proto.Unmarshal(data, &pbs)
				message_handler.OFI_Camt026(pbs, op)
				return
			case constant.CAMT087:
				// For requesting unable tp apply request (camt026 message)
				var pbs camt087pbstruct.SendPayload
				proto.Unmarshal(data, &pbs)
				message_handler.RFI_Camt087(pbs, op)
				return
			default:
				LOGGER.Errorf("No matching XML message type found")
				return
			}

		case constant.ISO8385:
			//report, err = ISO8583_handler(op, messageType, data)
		case constant.MT:
			//report, err = MT_handler(op, messageType, data)
		case constant.JSON:
			//report, err = JSON_handler(op, messageType, data)
		default:
			LOGGER.Errorf("No matching standard message type found")
			return
		}
	case kafka.RESPONSE_TOPIC:
		responseMsgType := strings.Split(pbs.MsgType, ":")
		// There are two types of response
		// 1. The response XML from RFI backend, which include ibwf001, camt029
		// 2. The error response from OFI or RFI send-service. This happens when there is anything wrong during
		// the request or response processing on OFI or RFI end. The response message type will appended with the error code at the end
		// and separate with `:`.
		switch len(responseMsgType) {
		case 2:
			switch standardType {
			case constant.ISO20022:

				switch messageType {
				case constant.IBWF001:
					// For replying a payment transaction (ibwf001 message)
					var pbs ibwf001pbstruct.SendPayload
					proto.Unmarshal(data, &pbs)
					message_handler.OFI_Ibwf001(pbs, op)
					return
				case constant.CAMT029:
					// For replying a payment cancellation message (camt029 message)
					var pbs camt029pbstruct.SendPayload
					proto.Unmarshal(data, &pbs)
					message_handler.OFI_Camt029(pbs, op)
					return
				case constant.PACS002:
					// For replying a asset redemption message (pacs002 message)
					var pbs pacs002pbstruct.SendPayload
					proto.Unmarshal(data, &pbs)
					message_handler.OFI_Pacs002(pbs, op)
					return
				default:
					LOGGER.Errorf("No matching XML message type found")
					return
				}

			case constant.ISO8385:
				//report, err = ISO8583_handler(op, messageType, data)
			case constant.MT:
				//report, err = MT_handler(op, messageType, data)
			case constant.JSON:
				//report, err = JSON_handler(op, messageType, data)
			default:
				LOGGER.Errorf("No matching standard message type found")
				return
			}
		case 3:
			message_handler.HandleErrMsg(pbs, op)
			return
		}
	}
	return
}
