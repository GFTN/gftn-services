// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package database

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoClient struct {
	svc *dynamodb.DynamoDB
	// ProfileFile    string
	// ProfileName    string
	AWS_SECRET_KEY string
	AWS_KEY_ID     string
	Region         string
	TableName      string
}
type Item struct {
	Participant string `json:"participant"`
	Whitelist   string `json:"whitelist"`
}

func (dc *DynamoClient) CreateConnection() error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(dc.Region),
		Credentials: credentials.NewStaticCredentials(dc.AWS_KEY_ID, dc.AWS_SECRET_KEY, ""),
		// Credentials: credentials.NewSharedCredentials(dc.ProfileFile, dc.ProfileName),
	})
	if err != nil {
		return err
	}
	dc.svc = dynamodb.New(sess)
	return nil
}

func (dc *DynamoClient) CreateTable() {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("participant"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("whitelist"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("participant"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("whitelist"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(dc.TableName),
	}

	_, err := dc.svc.CreateTable(input)

	if err != nil {
		LOGGER.Error("Got error calling CreateTable:")
		LOGGER.Error(err.Error())
		os.Exit(1)
	}
	LOGGER.Info("Created the table")

}

func (dc *DynamoClient) DeleteTable() {
	input := &dynamodb.DeleteTableInput{
		TableName: aws.String("ParticipantWhitelist"),
	}
	_, err := dc.svc.DeleteTable(input)
	if err != nil {
		LOGGER.Error("Got error calling DeleteItem")
		LOGGER.Error(err.Error())
		return
	}
	LOGGER.Info("Deleted table")
}

func (dc *DynamoClient) DeleteWhitelistParticipant(participantID, wlParticipant string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(dc.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"participant": {
				S: aws.String(participantID),
			},
			"whitelist": {
				S: aws.String(wlParticipant),
			},
		},
	}

	_, err := dc.svc.DeleteItem(input)

	if err != nil {
		return err
	}
	return nil
}

func (dc *DynamoClient) AddWhitelistParticipant(participant, wlparticipant string) error {
	item := Item{
		Participant: participant,
		Whitelist:   wlparticipant,
	}
	av, err := dynamodbattribute.MarshalMap(item)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(dc.TableName),
	}

	_, err = dc.svc.PutItem(input)

	if err != nil {
		return err
	}
	return nil
}

func (dc *DynamoClient) GetWhiteListParicipants(participantID string) ([]string, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(dc.TableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"participant": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(participantID),
					},
				},
			},
		},
	}
	results, err := dc.svc.Query(queryInput)
	if err != nil {
		return nil, err
	}
	var whitelist []string
	for _, result := range results.Items {
		item := Item{}
		err = dynamodbattribute.UnmarshalMap(result, &item)
		if err != nil {
			LOGGER.Error(fmt.Sprintf("Failed to unmarshal Record, %v", err))
			return nil, err
		}
		whitelist = append(whitelist, item.Whitelist)
	}
	return whitelist, nil
}
