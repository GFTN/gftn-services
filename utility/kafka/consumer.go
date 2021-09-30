// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package kafka

import (
	"github.com/GFTN/gftn-services/utility/payment/constant"

	"github.com/GFTN/gftn-services/utility/global-environment/services"
)

func (ops *KafkaOpreations) Consume(data chan<- []byte, topic string, consumerIndex int) {

	//data <- []byte(constant.KAFKA_INITIAL_SUCCESS)
	subscribeErr := ops.Consumers[consumerIndex].SubscribeTopics([]string{topic}, nil)
	if subscribeErr != nil {
		LOGGER.Error("Failed to subscribe topic: %s, %s", topic, subscribeErr.Error())
		return
	}

	LOGGER.Infof("Kafka req consumer start to subscribe on topic: %s", topic)

	for {
		msg, err := ops.Consumers[consumerIndex].ReadMessage(-1)
		if err == nil {
			data <- msg.Value
		} else {
			LOGGER.Errorf("Error reading message from the Kafka brokers: %s", err.Error())
			break
		}
	}

	LOGGER.Info("Try to reconnect to the Kafka brokers.")
	data <- []byte(constant.KAFKA_CONSUMER_RECONNECT)
	return
}

func (op *KafkaOpreations) consumerStartListening(topic, topicType string, consumerIndex int, kafkaRouter func(string, []byte, *KafkaOpreations)) {
	//start req consumer
	reqDataFromKafka := make(chan []byte)
	LOGGER.Debug("---------Start listening from Kafka---------")

	go op.Consume(reqDataFromKafka, topic, consumerIndex)

	for {
		reqData := <-reqDataFromKafka
		errorString := string(reqData)
		if errorString == constant.KAFKA_INITIAL_ERROR {
			LOGGER.Errorf("Failed to initiate the Kafka consumer client: %s, %s", topic, errorString)
			// Read the latest configuration from secret manager
			services.VariableCheck()
			services.InitEnv()
			op.InitPaymentConsumer(op.GroupId, kafkaRouter)
			go op.Consume(reqDataFromKafka, topic, consumerIndex)
			continue
		} else if errorString == constant.KAFKA_INITIAL_SUCCESS {
			continue
		} else if errorString == constant.KAFKA_CONSUMER_RECONNECT {
			LOGGER.Warning(errorString)
			go op.Consume(reqDataFromKafka, topic, consumerIndex)
			continue
		}

		LOGGER.Infof("Consume request from kafka topic: %s from consumer %v", topic, consumerIndex)
		go kafkaRouter(topicType, reqData, op)
		LOGGER.Debug("---------Message Consumed---------")
	}
}
