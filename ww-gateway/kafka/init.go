// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package kafka

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/GFTN/gftn-services/ww-gateway/environment"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	kafka_utils "github.com/GFTN/gftn-services/utility/kafka"
)

func Initialize() (*kafka.Consumer, error) {

	var kafkaActor *kafka_utils.KafkaOpreations

	/*
		Retrieve kafka settings
	*/

	if len(os.Getenv(kafka_utils.ENV_KEY_KAFKA_BROKER_URL)) == 0 {
		LOGGER.Errorf("Kafka broker URL is empty")
		return &kafka.Consumer{}, errors.New("kafka broker url is empty")
	}

	if len(os.Getenv(environment.ENV_KEY_KAFKA_PARTITION_NUMBER)) == 0 {
		LOGGER.Errorf("Kafka partition number not specified")
		return &kafka.Consumer{}, errors.New("Kafka partition number not specified")
	}

	kafkaActor = &kafka_utils.KafkaOpreations{}
	kafkaActor.BrokerURL = os.Getenv(kafka_utils.ENV_KEY_KAFKA_BROKER_URL)
	// Check if the environment variable KAFKA_ENABLE_SSL was set to true. If it's true, setting up certificate that
	// will be use by the Kafka producer and consumer.
	if os.Getenv(kafka_utils.ENV_KEY_KAFKA_ENABLE_SSL) == "true" {
		kafkaActor.SecurityProtocol = kafka_utils.KAFKA_SSL
		kafkaActor.SslCaLocation = os.Getenv(kafka_utils.ENV_KEY_KAFKA_CA_LOCATION)
		kafkaActor.SslCertificateLocation = os.Getenv(kafka_utils.ENV_KEY_KAFKA_CERTIFICATE_LOCATION)
		kafkaActor.SslKeyLocation = os.Getenv(kafka_utils.ENV_KEY_KAFKA_KEY_LOCATION)
		pw, _ := ioutil.ReadFile(os.Getenv(kafka_utils.ENV_KEY_KAFKA_KEY_PASSWORD))
		kafkaActor.SslKeyPassword = string(pw)
	} else {
		kafkaActor.SecurityProtocol = "false"
	}

	/*
		Init consumer settings
	*/

	consumer, err := InitConsumer(
		kafkaActor.BrokerURL,
		kafkaActor.SecurityProtocol,
		kafkaActor.SslCaLocation,
		kafkaActor.SslCertificateLocation,
		kafkaActor.SslKeyLocation,
		kafkaActor.SslKeyPassword)
	if err != nil {
		LOGGER.Errorf("Error creating the Kafka consumer: %s", err.Error())
		return &kafka.Consumer{}, err
	}
	return consumer, nil
}

func InitConsumer(brokerURL, securityProtocol, caLocation, certLocation, keyLocation, keyPassword string) (*kafka.Consumer, error) {

	c := &kafka.Consumer{}
	var err error
	if securityProtocol == kafka_utils.KAFKA_SSL {
		c, err = kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers":         brokerURL,
			"group.id":                  "G2",
			"auto.offset.reset":         "latest",
			"security.protocol":         securityProtocol,
			"ssl.ca.location":           caLocation,
			"ssl.certificate.location":  certLocation,
			"ssl.key.location":          keyLocation,
			"ssl.key.password":          keyPassword,
			"session.timeout.ms":        60000,
			"max.partition.fetch.bytes": 3000000,
			"enable.auto.commit":        false,
			"enable.auto.offset.store":  false,
			"max.poll.interval.ms":      86400000,
		})
	} else {
		c, err = kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers":         brokerURL,
			"group.id":                  "G2",
			"auto.offset.reset":         "latest",
			"session.timeout.ms":        60000,
			"max.partition.fetch.bytes": 3000000,
			"enable.auto.commit":        false,
			"enable.auto.offset.store":  false,
			"max.poll.interval.ms":      86400000,
		})
	}

	return c, err
}
