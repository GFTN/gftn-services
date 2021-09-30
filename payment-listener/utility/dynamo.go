// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package utility

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/op/go-logging"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/payment-listener/environment"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
)

var LOGGER = logging.MustGetLogger("database-client")

type DynamoClient struct {
	Session   *dynamodb.DynamoDB
	Region    *string
	TableName *string
}

type CursorData struct {
	Account string `json:"account_name"`
	Cursor  string `json:"account_cursor"`
}

type CursorUpdateData struct {
	Cursor string `json:":c"`
}

type CursorKey struct {
	Account string `json:"account_name"`
}

func CreateConnection() (*DynamoClient, error) {
	if os.Getenv(environment.ENV_KEY_AWS_DYNAMO_REGION) == "" {
		return &DynamoClient{}, errors.New("Dynamo DB region/table is empty")
	}
	region := os.Getenv(environment.ENV_KEY_AWS_DYNAMO_REGION)
	env := os.Getenv(global_environment.ENV_KEY_ENVIRONMENT_VERSION)
	participantId := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	table := env + "_" + participantId + "_cursor"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return &DynamoClient{}, err
	}
	client := &DynamoClient{}
	client.Session = dynamodb.New(sess)
	client.TableName = &table
	client.Region = &region
	return client, nil
}

func (client *DynamoClient) AddAccountCursor(accountName, cursor string) error {

	LOGGER.Infof("Adding cursor %v to %v account at table %v", cursor, accountName, *client.TableName)
	if err := client.CheckTableExists(); err != nil {
		LOGGER.Errorf("Failed updating account cursor: %v", err.Error())
		return err
	}

	item := CursorData{
		Account: accountName,
		Cursor:  cursor,
	}

	itemMap, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return nil
	}
	input := &dynamodb.PutItemInput{
		Item:      itemMap,
		TableName: client.TableName,
	}

	_, err = client.Session.PutItem(input)

	if err != nil {
		return err
	}
	LOGGER.Info("Cursor successfully added!")
	return nil
}

func (client *DynamoClient) GetAccountCursor(accountName string) (*model.Cursor, error) {

	LOGGER.Infof("Retrieving cursor of %v account at table %v", accountName, *client.TableName)
	var cursorResult = &model.Cursor{}
	if err := client.CheckTableExists(); err != nil {
		LOGGER.Errorf("Failed updating account cursor: %v", err.Error())
		return nil, err
	}

	var queryInput = &dynamodb.QueryInput{
		TableName: client.TableName,
		KeyConditions: map[string]*dynamodb.Condition{
			"account_name": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(accountName),
					},
				},
			},
		},
	}
	results, err := client.Session.Query(queryInput)
	if err != nil {
		LOGGER.Errorf("Database query error: %s", err.Error())
		return nil, err
	}
	var cursor string
	for _, result := range results.Items {
		item := CursorData{}
		err = dynamodbattribute.UnmarshalMap(result, &item)
		if err != nil {
			LOGGER.Errorf("Failed to unmarshal Record, %s", err.Error())
			return nil, err
		}
		cursor = item.Cursor
		break
	}
	cursorResult.Cursor = cursor
	LOGGER.Info("Cursor successfully retrieved!")
	return cursorResult, nil
}

func (client *DynamoClient) UpdateAccountCursor(accountName, cursor string) error {

	LOGGER.Infof("Updating cursor %v for %v account at table %v", cursor, accountName, *client.TableName)
	if err := client.CheckTableExists(); err != nil {
		LOGGER.Errorf("Failed updating account cursor: %v", err.Error())
		return err
	}
	cursorID, err := dynamodbattribute.MarshalMap(CursorKey{
		Account: accountName,
	})

	updateItem, err := dynamodbattribute.MarshalMap(CursorUpdateData{
		Cursor: cursor,
	})
	LOGGER.Infof("updateItem %v %+v", accountName, cursorID)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: updateItem,
		TableName:                 client.TableName,
		Key:                       cursorID,
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String("set account_cursor = :c"),
	}

	_, err = client.Session.UpdateItem(input)

	if err != nil {
		LOGGER.Errorf("Database update error: %s", err.Error())
		return err
	}

	LOGGER.Info("Cursor successfully updated!")
	return nil
}

func (client *DynamoClient) CheckTableExists() error {

	LOGGER.Infof("Querying table %v", *client.TableName)

	input := &dynamodb.DescribeTableInput{
		TableName: client.TableName,
	}

	result, err := client.Session.DescribeTable(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				LOGGER.Errorf("%v %s", dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
				LOGGER.Warningf("Cannot find the table: %v, creating now", *client.TableName)
				if err := client.CreateTable(); err != nil {
					LOGGER.Errorf("Failed creating table: %v", *client.TableName)
					return err
				}
			case dynamodb.ErrCodeInternalServerError:
				LOGGER.Errorf("%v %s", dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				LOGGER.Errorf("%s", aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			LOGGER.Errorf("%s", err.Error())
		}
		return err
	}
	LOGGER.Infof("Accessing table %v", *result.Table.TableName)
	return nil
}

func (client *DynamoClient) CreateTable() error {

	LOGGER.Infof("Creating table %v", *client.TableName)

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("account_name"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("account_name"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: client.TableName,
	}

	result, err := client.Session.CreateTable(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceInUseException:
				LOGGER.Errorf("%v %s", dynamodb.ErrCodeResourceInUseException, aerr.Error())
			case dynamodb.ErrCodeLimitExceededException:
				LOGGER.Errorf("%v %s", dynamodb.ErrCodeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				LOGGER.Errorf("%v %s", dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				LOGGER.Errorf("%s", aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			LOGGER.Errorf("%s", err.Error())
		}
		return err
	}

	LOGGER.Infof("Table %v created successfully", *result.TableDescription.TableName)
	return nil

}
