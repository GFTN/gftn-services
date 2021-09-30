// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package kafka

import (
	"errors"
	"io/ioutil"
	"os"

	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	whitelist_handler "github.com/GFTN/gftn-services/utility/payment/utils/whitelist-handler"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/GFTN/gftn-services/utility/payment/constant"
	"github.com/GFTN/gftn-services/utility/payment/utils/signing"
	"github.com/GFTN/gftn-services/utility/payment/utils/transaction"
)

type KafkaOpreations struct {
	BrokerURL              string
	SecurityProtocol       string
	SslCaLocation          string
	SslCertificateLocation string
	SslKeyLocation         string
	SslKeyPassword         string
	Producer               *kafka.Producer
	Consumers              []*kafka.Consumer
	GroupId                string
	//used only by send-service
	FundHandler      transaction.CreateFundingOpereations
	SignHandler      signing.CreateSignOperations
	WhitelistHandler whitelist_handler.ParticipantWhiteList
}

func Initialize() (*KafkaOpreations, error) {

	var ops = &KafkaOpreations{}
	ops.BrokerURL = os.Getenv(ENV_KEY_KAFKA_BROKER_URL)
	if len(ops.BrokerURL) == 0 {
		LOGGER.Errorf("Kafka broker URL is empty")
		return &KafkaOpreations{}, errors.New("kafka broker url is empty")
	}

	// Check if the environment variable KAFKA_ENABLE_SSL was set to true. If it's true, setting up certificate that
	// will be use by the Kafka producer and consumer.
	if os.Getenv(ENV_KEY_KAFKA_ENABLE_SSL) == "true" {
		ops.SecurityProtocol = KAFKA_SSL
		ops.SslCaLocation = os.Getenv(ENV_KEY_KAFKA_CA_LOCATION)
		ops.SslCertificateLocation = os.Getenv(ENV_KEY_KAFKA_CERTIFICATE_LOCATION)
		ops.SslKeyLocation = os.Getenv(ENV_KEY_KAFKA_KEY_LOCATION)
		pw, _ := ioutil.ReadFile(os.Getenv(ENV_KEY_KAFKA_KEY_PASSWORD))
		ops.SslKeyPassword = string(pw)
	} else {
		ops.SecurityProtocol = "false"
	}
	err := ops.InitProducer()
	if err != nil {
		LOGGER.Errorf("Error initializing Kafka producer: %v", err.Error())
		return &KafkaOpreations{}, err
	}
	return ops, nil
}

func (ops *KafkaOpreations) InitProducer() error {

	p := &kafka.Producer{}
	var err error

	if ops.SecurityProtocol == KAFKA_SSL {
		p, err = kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers":        ops.BrokerURL,
			"security.protocol":        ops.SecurityProtocol,
			"ssl.ca.location":          ops.SslCaLocation,
			"ssl.certificate.location": ops.SslCertificateLocation,
			"ssl.key.location":         ops.SslKeyLocation,
			"ssl.key.password":         ops.SslKeyPassword,
		})
	} else {
		p, err = kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": ops.BrokerURL,
		})
	}

	if err != nil {
		return err
	}

	ops.Producer = p
	return nil

}

func (ops *KafkaOpreations) InitPaymentConsumer(groupId string, router func(string, []byte, *KafkaOpreations)) error {

	listeningTopics := []string{REQUEST_TOPIC, RESPONSE_TOPIC}

	for consumerIndex, topicType := range listeningTopics {

		c := &kafka.Consumer{}
		var err error

		if groupId != "" {
			ops.GroupId = groupId
		} else {
			return errors.New("Group ID is empty")
		}
		LOGGER.Infof("Initiating Kafka consumer at %v, consumer index: %v", ops.GroupId, consumerIndex)
		if ops.SecurityProtocol == constant.KAFKA_SSL {
			c, err = kafka.NewConsumer(&kafka.ConfigMap{
				"bootstrap.servers":        ops.BrokerURL,
				"group.id":                 ops.GroupId,
				"auto.offset.reset":        "latest",
				"security.protocol":        ops.SecurityProtocol,
				"ssl.ca.location":          ops.SslCaLocation,
				"ssl.certificate.location": ops.SslCertificateLocation,
				"ssl.key.location":         ops.SslKeyLocation,
				"ssl.key.password":         ops.SslKeyPassword,
			})
		} else {
			c, err = kafka.NewConsumer(&kafka.ConfigMap{
				"bootstrap.servers": ops.BrokerURL,
				"group.id":          ops.GroupId,
				"auto.offset.reset": "latest",
			})
		}

		if err != nil {
			LOGGER.Errorf("Error creating the Kafka consumer: %s", err.Error())
			return err
		}

		ops.Consumers = append(ops.Consumers, c)

		homeDomainName := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
		kafkaReqTopic := homeDomainName + topicType
		go ops.consumerStartListening(kafkaReqTopic, topicType, consumerIndex, router)

	}

	// initialize admin client

	homeDomain := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	prServiceURL := os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL)
	ops.FundHandler = transaction.InitiateFundingOperations(prServiceURL, homeDomain)

	// initailize sign handler
	ops.SignHandler = signing.InitiateSignOperations(prServiceURL)

	// initialize whitelist handler
	ops.WhitelistHandler = whitelist_handler.CreateWhiteListServiceOperations()

	return nil
}
