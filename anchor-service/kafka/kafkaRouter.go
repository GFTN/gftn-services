// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package kafka

import (
	"strings"

	message_handler "github.com/GFTN/gftn-services/utility/payment/message-handler"

	"github.com/golang/protobuf/proto"

	pacs009pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/pacs00900108"
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
			case constant.PACS009:
				// For requesting a payment transaction (pacs008 message)
				var pbs pacs009pbstruct.SendPayload
				proto.Unmarshal(data, &pbs)
				message_handler.RFI_Pacs009(pbs, op)
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
				/*
					case constant.IBWF001:
						// For replying a payment transaction (ibwf001 message)
						var pbs ibwf001pbstruct.SendPayload
						proto.Unmarshal(data, &pbs)
						message_handler.OFI_Ibwf001(pbs, op)
						return
				*/
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
