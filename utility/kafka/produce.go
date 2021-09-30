// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package kafka

import (
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gogo/protobuf/proto"
	"github.com/GFTN/iso20022/pacs00200109"
	"github.com/GFTN/gftn-services/utility/payment/utils/parse"
	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"
)

func (ops *KafkaOpreations) Produce(topic string, msg []byte) error {

	produceResult := make(chan bool)
	// Delivery report kafka for produced messages
	go func() {
		for e := range ops.Producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					LOGGER.Errorf("Produce delivery failed: %v\n", ev.TopicPartition)
					produceResult <- false
				} else {
					LOGGER.Infof("Produce message to %v\n", ev.TopicPartition)
					produceResult <- true
				}
				return
			case kafka.Error:
				LOGGER.Errorf("Encounter error while producing message to Kafka: %v", ev)
				produceResult <- false
				return
			default:
				LOGGER.Errorf("Exception occured: %v", ev)
				produceResult <- false
				return
			}
		}
	}()

	LOGGER.Infof("**************** Producing Message ****************")

	now := time.Now().UTC()
	timestamp := strconv.FormatInt(now.UnixNano(), 10)

	produceErr := ops.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
		Key:            []byte(timestamp),
		Timestamp:      now,
	}, nil)
	if produceErr != nil {
		LOGGER.Errorf("Kafka broker produce error: %s", produceErr.Error())
		return produceErr
	}

	for {
		r := <-produceResult
		if !r {
			LOGGER.Errorf("kafka producer failed")
			break
		} else {
			LOGGER.Infof("**************** Message Produced ****************")
			// Wait for message deliveries
			ops.Producer.Flush(15 * 1000)
			return nil
		}
	}

	LOGGER.Warningf("Encounter error while producing message to Kafka, re-intializing producer...")
	err := ops.InitProducer()
	if err != nil {
		LOGGER.Errorf("Failed re-initializing producer: %v", err)
		return err
	}
	LOGGER.Warningf("Reproducing message to topic %v", topic)
	err = ops.Produce(topic, msg)
	if err != nil {
		LOGGER.Errorf("Failed re-producing message: %v", err)
		return err
	}
	return nil
}

func (op *KafkaOpreations) SendRequestToKafka(topic string, msg []byte) error {

	/*
		sending message to Kafka
	*/

	err := op.Produce(topic, msg)
	if err != nil {
		LOGGER.Errorf("Encounter error while producing message to Kafka topic: %v", topic)
		return err
	}

	return nil
}

// Send back errors happened on RFI site during request processing to OFI
func (op *KafkaOpreations) SendErrMsg(instructionId, standardType, reqMsgType, ofiId, rfiId string, errType int, orgnlgrpInf *pacs00200109.OriginalGroupInformation29) {
	errCode := strconv.Itoa(errType)

	res := &sendmodel.SendPayload{
		InstructionId: instructionId,
		MsgType:       standardType + ":" + reqMsgType + ":" + errCode,
		OfiId:         ofiId,
		RfiId:         rfiId,
	}

	dataBuffer, err := proto.Marshal(res)
	if err != nil {
		LOGGER.Error("Parse Error: " + err.Error())
		return
	}

	targetParticipant, _, err := parse.KafkaErrorRouter(reqMsgType, instructionId, ofiId, rfiId, 0, false, orgnlgrpInf)
	if err != nil {
		return
	}

	op.Produce(targetParticipant+RESPONSE_TOPIC, dataBuffer)

	return
}
