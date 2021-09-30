// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package kafka

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/GFTN/gftn-services/ww-gateway/environment"
	"github.com/GFTN/gftn-services/ww-gateway/utility"
)

func ReqConsumer(c *kafka.Consumer, topic string, partition int32) (kafka.TopicPartitions, []*kafka.Message, error) {

	LOGGER.Infof("Requesting message from Kafka topic %v", topic)
	var startOffset int64

	LOGGER.Infof("Retrieving starting offset for partition %v", partition)

	/*
		retrieving lowest/highest offset in the partition
	*/
	lowOffset, highOffset, err := c.QueryWatermarkOffsets(topic, partition, -1)
	if err != nil {
		LOGGER.Errorf("Encounter error while querying the last commit offset %v", err)
		return kafka.TopicPartitions{}, []*kafka.Message{}, err
	}

	LOGGER.Infof("The highest offset: %v", highOffset)
	LOGGER.Infof("The lowest offset: %v", lowOffset)

	if lowOffset == highOffset {
		LOGGER.Warningf("Topic %v is empty", topic)
		return kafka.TopicPartitions{}, []*kafka.Message{}, errors.New(utility.STATUS_QUEUE_EMPTY)
	}

	/*
		retrieving the latest read offset in the partition
	*/
	part, err := c.Committed(kafka.TopicPartitions{{Topic: &topic, Partition: partition}}, -1)
	if err != nil {
		LOGGER.Errorf("error fetching commit %v", err)
		return kafka.TopicPartitions{}, []*kafka.Message{}, err
	}
	if len(part) == 0 {
		LOGGER.Infof("No committed offset detected, starting from low offset: %v", lowOffset)
		startOffset = lowOffset
	} else {
		commitedOffset := int64(part[0].Offset)
		if commitedOffset < lowOffset {
			startOffset = lowOffset
		} else {
			startOffset = commitedOffset
		}
	}

	// if starting offset is the same as the highest offset in this partition, it means no new message
	if startOffset == highOffset {
		LOGGER.Infof("No new message to consume from %v[%v]@%v", topic, partition, startOffset)
		return kafka.TopicPartitions{}, []*kafka.Message{}, errors.New(utility.STATUS_QUEUE_EMPTY)
	}

	LOGGER.Infof("Continue receiving message at offset: %v for partition %v", startOffset, partition)

	/*
		calculate which offset to stop for this round
	*/
	count, _ := strconv.Atoi(os.Getenv(environment.ENV_KEY_BATCH_LIMIT))
	count64 := int64(count)

	var upperBound int64
	if upperBound = startOffset + count64; upperBound >= highOffset {
		upperBound = highOffset
	}

	var partitions kafka.TopicPartitions

	/*
		assigning starting offset & topic/partition
	*/

	tempOffSet, err := kafka.NewOffset(fmt.Sprintf("%v", startOffset))
	if err != nil {
		return kafka.TopicPartitions{}, []*kafka.Message{}, err
	}
	partitions = append(partitions, kafka.TopicPartition{
		Topic:     &topic,
		Partition: partition,
		Offset:    tempOffSet,
		Error:     err,
	})

	LOGGER.Infof("Assgining partition & offset to consumer manually")

	err = c.Assign(partitions)
	if err != nil {
		LOGGER.Errorf("Assign failed: %s", err)
		return kafka.TopicPartitions{}, []*kafka.Message{}, err
	}
	defer c.Unassign()

	assignment, err := c.Assignment()
	if err != nil {
		LOGGER.Errorf("Assignment() failed: %s", err)
		return kafka.TopicPartitions{}, []*kafka.Message{}, err
	}
	LOGGER.Infof("Assignment %v\n", assignment)
	timeout, _ := strconv.Atoi(os.Getenv(environment.ENV_KEY_MESSAGE_RETRIEVE_TIMEOUT))
	// if not set, default value will be 3 seconds
	if timeout == 0 {
		timeout = 3000
	}
	LOGGER.Debugf("****** Start receiving message from Kafka offset ******")

	/*
		start getting message from starting offset to upperbound
	*/

	var actualVisitedIndex int64
	var messages []*kafka.Message
Loop:
	for actualVisitedIndex = startOffset; actualVisitedIndex < upperBound; actualVisitedIndex++ {
		ev := c.Poll(timeout)

		switch e := ev.(type) {
		case *kafka.Message:
			if e.TopicPartition.Error != nil {
				LOGGER.Infof("Encounter error: %+v", e)
				return kafka.TopicPartitions{}, []*kafka.Message{}, e.TopicPartition.Error
			}
			LOGGER.Infof("Received message from Kafka: %+v", e)
			messages = append(messages, e)
		case kafka.Error:
			LOGGER.Errorf("Error: %+v", e)
			return kafka.TopicPartitions{}, []*kafka.Message{}, e
		case kafka.PartitionEOF:
			LOGGER.Warningf("End of partition")
			break Loop
		default:
			LOGGER.Errorf("Message retrieval timeout")
			return kafka.TopicPartitions{}, []*kafka.Message{}, errors.New("Message retrieval timeout")
		}

	}
	LOGGER.Debug("****** Complete receiving message ******")

	/*
		Recording the last offset being read, and commit it manually to Kafka after respond to the client
	*/

	var finalPartition kafka.TopicPartitions

	LOGGER.Infof("Newly visited offset: %v at partition %v", actualVisitedIndex, partition)
	tempOffSet, err = kafka.NewOffset(fmt.Sprintf("%v", actualVisitedIndex))
	if err != nil {
		LOGGER.Errorf(err.Error())
		return kafka.TopicPartitions{}, []*kafka.Message{}, err
	}
	finalPartition = append(finalPartition, kafka.TopicPartition{
		Topic:     &topic,
		Partition: partition,
		Offset:    tempOffSet,
		Error:     err,
	})

	return finalPartition, messages, err

}

/* for testing purpose
func ResetOffset(c *kafka.Consumer, topic string) error {
	LOGGER.Infof("Resetting message offset of Kafka topic %v", topic)
	part, err := c.Committed(kafka.TopicPartitions{{Topic: &topic, Partition: 0}}, -1)
	if err != nil {
		LOGGER.Errorf("error fetching commit %v", err)
	}

	if len(part) == 0 {
		err := errors.New("Cannot find the specified kafka topic")
		LOGGER.Errorf(err.Error())
		return err
	}
	LOGGER.Infof("Current offset is %+v", part[0].Offset)

	var resetPartition kafka.TopicPartitions

	tempOffSet, err := kafka.NewOffset(fmt.Sprintf("%v", 0))
	if err != nil {
		return err
	}
	resetPartition = append(resetPartition, kafka.TopicPartition{
		Topic:     &topic,
		Partition: 0,
		Offset:    tempOffSet,
		Error:     err,
	})

	part, err = c.CommitOffsets(resetPartition)
	if err != nil {
		LOGGER.Errorf("error committing %v", err)
		return err
	}
	LOGGER.Infof("commiting offset %+v", part)

	part, err = c.Committed(kafka.TopicPartitions{{Topic: &topic, Partition: 0}}, -1)
	if err != nil {
		LOGGER.Errorf("error fetching commit %v", err)
		return err
	}
	LOGGER.Infof("The offset of %v topic is now %+v", topic, part)
	return nil
}
*/
