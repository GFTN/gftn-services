// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package database

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/op/go-logging"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
)

var DC DynamoClient
var LOGGER = logging.MustGetLogger("database-client")

type DynamoClient struct {
	svc         *dynamodb.DynamoDB
	ProfileFile string
	ProfileName string
	Region      string
	TableName   string
}

type Item struct {
	Category string `json:"category"`
	Details  Detail `json:"details"`
}

type Detail struct {
	VersionName    string `json:"version_name"`
	DeploymentFile string `json:"deployment_file"`
	Timestamp      string `json:"create_time"`
}

type Key struct {
	Category string `json:"category"`
}

type DeploymentUpdate struct {
	VersionName    string `json:":n"`
	DeploymentFile string `json:":d"`
	Timestamp      string `json:":t"`
}

func (dc *DynamoClient) CreateConnection() error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(dc.Region),
		Credentials: credentials.NewStaticCredentials(os.Getenv(global_environment.ENV_KEY_AWS_ACCESS_KEY_ID), os.Getenv(global_environment.ENV_KEY_AWS_SECRET_ACCESS_KEY), ""),
		//Credentials: credentials.NewSharedCredentials(dc.ProfileFile, dc.ProfileName),
	})
	if err != nil {
		return err
	}

	dc.svc = dynamodb.New(sess)
	return nil
}

func (dc *DynamoClient) DeleteData(tableName, tableKey string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"category": {
				S: aws.String(tableKey),
			},
		},
	}

	_, err := dc.svc.DeleteItem(input)

	if err != nil {
		return err
	}
	return nil
}

func (dc *DynamoClient) AddData(tableName, c, versionName, file string) error {
	item := Item{
		Category: c,
		Details: Detail{
			VersionName:    versionName,
			DeploymentFile: file,
			Timestamp:      time.Now().String(),
		},
	}
	av, err := dynamodbattribute.MarshalMap(item)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dc.svc.PutItem(input)

	if err != nil {
		return err
	}

	return nil
}

func (dc *DynamoClient) GetData(tableName, tableKey string) (*string, *string, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"category": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(tableKey),
					},
				},
			},
		},
	}
	results, err := dc.svc.Query(queryInput)
	if err != nil {
		LOGGER.Errorf("Database query error: %s", err.Error())
		return nil, nil, err
	}
	var versionName, file string
	for _, result := range results.Items {
		item := Item{}
		err = dynamodbattribute.UnmarshalMap(result, &item)
		if err != nil {
			LOGGER.Errorf("Failed to unmarshal Record, %s", err.Error())
			return nil, nil, err
		}
		versionName = item.Details.VersionName
		file = item.Details.DeploymentFile
		break
	}
	return &versionName, &file, nil
}

func (dc *DynamoClient) UpdateData(tableName, c, versionName, file string) error {
	txID, err := dynamodbattribute.MarshalMap(Key{
		Category: c,
	})

	updateItem, err := dynamodbattribute.MarshalMap(DeploymentUpdate{
		VersionName:    versionName,
		DeploymentFile: file,
		Timestamp:      time.Now().String(),
	})

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: updateItem,
		TableName:                 aws.String(tableName),
		Key:                       txID,
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String("set details.version_name = :n, details.deployment_file = :d, details.create_time = :t"),
	}

	_, err = dc.svc.UpdateItem(input)

	if err != nil {
		LOGGER.Errorf("Database update error: %s", err.Error())
		return err
	}

	LOGGER.Debug("Successfully updated")
	return nil
}
