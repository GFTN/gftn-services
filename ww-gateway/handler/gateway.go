// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	go_kafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/GFTN/gftn-services/gftn-models/model"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	kafka_utils "github.com/GFTN/gftn-services/utility/kafka"
	"github.com/GFTN/gftn-services/utility/response"
	"github.com/GFTN/gftn-services/ww-gateway/environment"
	"github.com/GFTN/gftn-services/ww-gateway/kafka"
	kafkaHandler "github.com/GFTN/gftn-services/ww-gateway/kafka"
	"github.com/GFTN/gftn-services/ww-gateway/utility"
)

type GatewayOperations struct {
	Consumer         map[string]*go_kafka.Consumer
	homeDomain       string
	CurrentPartition *int32
}

func InitGatewayOperation() (GatewayOperations, error) {

	var initPartition = int32(0)
	operation := GatewayOperations{
		homeDomain:       os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME),
		CurrentPartition: &initPartition,
	}

	LOGGER.Infof("* Initiate Kafka consumer")

	operation.homeDomain = os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)

	/*
		Init consumer settings
	*/
	operation.Consumer = make(map[string]*go_kafka.Consumer, len(kafka_utils.SUPPORT_MESSAGE_TYPES))

	for _, msg := range kafka_utils.SUPPORT_MESSAGE_TYPES {
		err := operation.InitializeConsumer(msg)
		if err != nil {
			LOGGER.Errorf("Error creating the Kafka consumer: %s", err.Error())
			return GatewayOperations{}, err
		}
		LOGGER.Infof("* Kafka consumer successfully subscribed to topic %v", msg)

	}

	LOGGER.Infof("* InitGatewayOperations finished")
	return operation, nil
}

func (operation GatewayOperations) ServiceCheck(w http.ResponseWriter, req *http.Request) {
	LOGGER.Infof("Performing service check")
	response.Respond(w, http.StatusOK, []byte(`{"status":"Alive"}`))
	return
}

func (operation GatewayOperations) GetMessage(w http.ResponseWriter, request *http.Request) {
	LOGGER.Infof("gateway-service:Gateway Operations :get message")

	/*
		parsing URL parameters
	*/

	queryParams := request.URL.Query()
	var topic string

	if len(queryParams["type"]) > 0 {
		topic = strings.ToUpper(strings.TrimSpace(queryParams["type"][0]))
	}

	if !utility.Contains(kafka_utils.SUPPORT_MESSAGE_TYPES, topic) {
		LOGGER.Errorf("Error while parsing message parameter : Specified message type is not supported")
		response.NotifyWWError(w, request, http.StatusBadRequest, "GATEWAY-1004", nil)
		return
	}

	topicName := operation.homeDomain + "_" + topic
	var totalPartition int32

	maxPartitionString := os.Getenv(environment.ENV_KEY_KAFKA_PARTITION_NUMBER)
	if maxPartitionString == "" {
		LOGGER.Errorf("Maximum Kafka partition number not set")
		response.NotifyWWError(w, request, http.StatusInternalServerError, "GATEWAY-1005", nil)
		return
	}

	temp, _ := strconv.Atoi(maxPartitionString)
	totalPartition = int32(temp)

	*operation.CurrentPartition = *operation.CurrentPartition % totalPartition
	/*
		reading message from Kafka
	*/
	newStart, messages, err := kafkaHandler.ReqConsumer(operation.Consumer[topic], topicName, *operation.CurrentPartition)
	if err != nil && err.Error() != utility.STATUS_QUEUE_EMPTY {
		LOGGER.Errorf("Encounter error while consuming message: %v", err)
		err := operation.InitializeConsumer(topic)
		if err != nil {
			LOGGER.Errorf("Failed re-intializing Kafka consumer: %v", err)
			response.NotifyWWError(w, request, http.StatusInternalServerError, "GATEWAY-1005", err)
			return
		}
		newStart, messages, err = kafkaHandler.ReqConsumer(operation.Consumer[topic], topicName, *operation.CurrentPartition)
		if err != nil {
			LOGGER.Errorf("Failed re-intializing Kafka consumer: %v", err)
			response.NotifyWWError(w, request, http.StatusInternalServerError, "GATEWAY-1005", err)
			return
		}
	}

	var results = model.GatewayResponse{}
	for _, elem := range messages {
		//LOGGER.Debugf("msg key %+v", string(elem.Value))
		var temp map[string]interface{}
		_ = json.Unmarshal(elem.Value, &temp)
		results.Data = append(results.Data, temp)
	}

	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	results.Timestamp = &timestamp
	b, _ := json.Marshal(&results)

	LOGGER.Infof("Get message from Kafka successfully")
	response.Respond(w, http.StatusOK, b)

	/*
		update offset after successfully return the message result
	*/
	if len(newStart) == 0 {
		LOGGER.Infof("No offset await commit, end of the handler")
		*operation.CurrentPartition++
		return
	}

	LOGGER.Infof("Now committing offset to %v[%v]@%v", topicName, newStart[0].Partition, newStart[0].Offset)

	part, err := operation.Consumer[topic].CommitOffsets(newStart)
	if err != nil {
		LOGGER.Warningf("error committing %s", err.Error())
	}
	part, err = operation.Consumer[topic].Committed(newStart, -1)
	if err != nil {
		LOGGER.Errorf("error fetching committed offset: %v", err)
	}

	if len(part) == 0 {
		LOGGER.Infof("Nothing to commit, move to next partition")
	} else {
		LOGGER.Infof("Offset %+v commited to topic: %v at partition %v", part[0].Offset, topicName, *operation.CurrentPartition)
	}
	*operation.CurrentPartition++
	return

}

func (operation GatewayOperations) InitializeConsumer(topic string) error {
	var err error
	LOGGER.Debugf("****** Initializing Kafka consumer ******")
	if operation.Consumer[topic] != nil {
		operation.Consumer[topic].Unsubscribe()
		operation.Consumer[topic].Unassign()
		operation.Consumer[topic].Close()
	}
	operation.Consumer[topic], err = kafka.Initialize()
	if err != nil {
		LOGGER.Errorf("Failed intializing Kafka consumer: %v", err)
		return err
	}

	LOGGER.Debugf("****** Kafka consumer successfully initialized ******")
	return nil
}

/*
func (operation GatewayOperations) ResetOffset(w http.ResponseWriter, request *http.Request) {
	LOGGER.Infof("gateway-service:Gateway Operations :reset offset")

	queryParams := request.URL.Query()
	var topic string

	if len(queryParams["type"]) > 0 {
		topic = strings.ToUpper(queryParams["type"][0])
	}

	if !utility.Contains(kafka_utils.SUPPORT_MESSAGE_TYPES, topic) {
		LOGGER.Errorf("Error while parsing message parameter : Specified message type is not supported")
		response.NotifyWWError(w, request, http.StatusBadRequest, "GATEWAY-1004", nil)
		return
	}

	topicName := operation.homeDomain + "_" + topic

	LOGGER.Infof("Resetting offset of topic %v to 0", topicName)
	err := kafkaHandler.ResetOffset(operation.Consumer[topic], topicName)
	if err != nil {
		LOGGER.Errorf("Error while resetting offset of Kafka :  %v", err)
		return
	}

	LOGGER.Infof("Reset offset of Kafka successfully")
	response.Respond(w, http.StatusOK, []byte(`{"status":"Success"}`))
	return

}
*/
