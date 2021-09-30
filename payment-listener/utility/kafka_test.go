// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package utility

import (
	"log"
	"os"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func TestExport(t *testing.T) {
	os.Setenv("AWS_DYNAMO_REGION", "ap-southeast-1")
	os.Setenv("DYNAMO_DB_TABLE_NAME", "payment-test")
	os.Setenv("ENV_VERSION", "dev")
	os.Setenv("HOME_DOMAIN_NAME", "participant3")
	os.Setenv("KAFKA_BROKER_URL", "kafka-1:9091,kafka-2:9092,kafka-3:9093")
}

func TestInitialize(t *testing.T) {

	brokerUrl := os.Getenv("KAFKA_BROKER_URL")
	if len(brokerUrl) == 0 {
		log.Fatalf("Kafka broker URL is empty")
	}

	log.Println("initialize producer")

	p := &kafka.Producer{}
	var err error

	p, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokerUrl,
	})

	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Println("------ start sending message -------")

	var payload []byte
	err = Produce("test", payload, p)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = Produce("test", payload, p)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = Produce("test", payload, p)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = Produce("test", payload, p)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func Produce(topic string, msg []byte, p *kafka.Producer) error {

	LOGGER.Infof("~~~~~~ operation %+v", p)
	produceResult := make(chan bool)
	// Delivery report kafka for produced messages
	go func() {
		LOGGER.Debugf("---------------------- Entering go routine ----------------------")
		for e := range p.Events() {
			LOGGER.Debugf("---------------------- for loop %v ----------------------", e)

			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Produce delivery failed: %v\n", ev.TopicPartition)
					produceResult <- false
				} else {
					log.Printf("Produce message to %v\n", ev.TopicPartition)
					produceResult <- true
					log.Printf("sucess\n")
				}
				return
			case kafka.Error:
				log.Printf("Encounter error while producing message to Kafka: %v", ev)
				produceResult <- false
				return
			default:
				log.Printf("Exception occured: %v", ev)
				produceResult <- false
				return
			}
		}
	}()

	LOGGER.Infof("**************** Producing Message ****************")

	produceErr := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, nil)
	if produceErr != nil {
		LOGGER.Errorf("Kafka broker produce error: %s", produceErr.Error())
		return produceErr
	}

	for {
		log.Println("for loop")
		r := <-produceResult
		log.Println("for loop2")

		if !r {
			log.Println("error")
			log.Fatalf("kafka producer failed")
		} else {
			break
		}
	}
	LOGGER.Infof("**************** Message Produced ****************")
	// Wait for message deliveries
	//p.Flush(15 * 1000)

	return nil
}
