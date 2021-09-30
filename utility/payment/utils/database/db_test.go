// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package database

import (
	"fmt"
	"testing"
)

var c = DynamoClient{
	//ProfileFile: os.Getenv(environment.ENV_KEY_DYNAMO_DB_PROFILEFILE),
	//ProfileName: os.Getenv(environment.ENV_KEY_DYNAMO_DB_PROFILENAME),
	Region:      "us-east-1",
}

func TestDynamoClient_CreateConnection(t *testing.T) {
	c.CreateConnection()
	fmt.Println("connection")
}

func TestDynamoClient_AddTransactionData(t *testing.T) {
	err := c.AddTransactionData("testid1", "testdata", "PENDING", "resId", "paymentInfo")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
	fmt.Println("add")
}

func TestDynamoClient_GetTransactionData(t *testing.T) {
	data, status, _, _, err := c.GetTransactionData("testid1")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
	fmt.Println(*data)
	fmt.Println(*status)
}

func TestDynamoClient_UpdateTransactionData(t *testing.T) {
	err := c.UpdateTransactionData("testid1", "updatedata", "DONE", "resIdDone", "paymentInfoDone")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
	fmt.Println("update")
}
