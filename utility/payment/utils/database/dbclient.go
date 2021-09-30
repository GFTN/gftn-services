// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package database

import (
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/op/go-logging"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/payment/environment"
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
	Transactionid string          `json:"transactionid"`
	Info          TransactionInfo `json:"info"`
}

type TransactionInfo struct {
	TxData      string `json:"tx_data"`
	TxStatus    string `json:"tx_status"`
	ResId       string `json:"res_id"`
	PaymentInfo string `json:"payment_info"`
}

type Key struct {
	Transactionid string `json:"transactionid"`
}

type TransactionInfoUpdate struct {
	TxData      string `json:":d,omitempty"`
	TxStatus    string `json:":s,omitempty"`
	ResId       string `json:":i,omitempty"`
	PaymentInfo string `json:":p,omitempty"`
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
	dc.TableName = os.Getenv(environment.ENV_KEY_DYNAMO_DB_TABLE_NAME)
	return nil
}

func (dc *DynamoClient) DeleteTransactionData(transactionID string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(dc.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"transactionid": {
				S: aws.String(transactionID),
			},
		},
	}

	_, err := dc.svc.DeleteItem(input)

	if err != nil {
		return err
	}
	return nil
}

func (dc *DynamoClient) AddTransactionData(transactionID, data, status, redId, paymentInfo string) error {
	tableName := os.Getenv(environment.ENV_KEY_DYNAMO_DB_TABLE_NAME)
	item := Item{
		Transactionid: transactionID,
		Info: TransactionInfo{
			TxData:      data,
			TxStatus:    status,
			ResId:       redId,
			PaymentInfo: paymentInfo,
		},
	}
	av, err := dynamodbattribute.MarshalMap(item)
	input := &dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(tableName),
		ConditionExpression: aws.String("attribute_not_exists(transactionid)"),
	}

	_, err = dc.svc.PutItem(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return errors.New("Instruction ID already been used, please use another one")
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				LOGGER.Errorf(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				LOGGER.Errorf(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				LOGGER.Errorf(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeTransactionConflictException:
				LOGGER.Errorf(dynamodb.ErrCodeTransactionConflictException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				LOGGER.Errorf(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				LOGGER.Errorf(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				LOGGER.Errorf(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			LOGGER.Errorf(err.Error())
		}
		return err
	}

	return nil
}

func (dc *DynamoClient) GetTransactionData(transactionID string) (*string, *string, *string, *string, error) {
	LOGGER.Infof("Retrieving transaction with id %v", transactionID)
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(dc.TableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"transactionid": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(transactionID),
					},
				},
			},
		},
	}
	results, err := dc.svc.Query(queryInput)
	if err != nil {
		LOGGER.Errorf("Database query error: %s", err.Error())
		return nil, nil, nil, nil, err
	}
	var data, status, resId, paymentInfo string
	for _, result := range results.Items {
		item := Item{}
		err = dynamodbattribute.UnmarshalMap(result, &item)
		if err != nil {
			LOGGER.Errorf("Failed to unmarshal Record, %s", err.Error())
			return nil, nil, nil, nil, err
		}
		data = item.Info.TxData
		status = item.Info.TxStatus
		resId = item.Info.ResId
		paymentInfo = item.Info.PaymentInfo
		break
	}
	return &data, &status, &resId, &paymentInfo, nil
}

func (dc *DynamoClient) UpdateTransactionData(transactionID, data, status, resId, paymentInfo string) error {
	txID, err := dynamodbattribute.MarshalMap(Key{
		Transactionid: transactionID,
	})

	var updateObj TransactionInfoUpdate
	queryString := "set "
	if data != "" {
		updateObj.TxData = data
		queryString += "info.tx_data = :d,"
	}
	if status != "" {
		updateObj.TxStatus = status
		queryString += "info.tx_status = :s,"
	}
	if resId != "" {
		updateObj.ResId = resId
		queryString += "info.res_id = :i,"
	}
	if paymentInfo != "" {
		updateObj.PaymentInfo = paymentInfo
		queryString += "info.payment_info = :p,"
	}

	queryString = strings.TrimSuffix(queryString, ",")

	updateItem, err := dynamodbattribute.MarshalMap(updateObj)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: updateItem,
		TableName:                 aws.String(dc.TableName),
		Key:                       txID,
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String(queryString),
	}

	_, err = dc.svc.UpdateItem(input)

	if err != nil {
		LOGGER.Errorf("Database update error: %s", err.Error())
		return err
	}

	LOGGER.Info("Successfully updated")
	return nil
}
