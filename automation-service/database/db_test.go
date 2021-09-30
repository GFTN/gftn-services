// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package database

import (
	"fmt"
	"github.com/GFTN/gftn-services/automation-service/environment"
	"os"
	"testing"
)

var c = DynamoClient{
	//ProfileFile: os.Getenv(environment.ENV_KEY_DYNAMO_DB_PROFILEFILE),
	//ProfileName: os.Getenv(environment.ENV_KEY_DYNAMO_DB_PROFILENAME),
	Region:      os.Getenv(environment.ENV_KEY_DYNAMO_DB_REGION),
}

func TestDynamoClient_CreateConnection(t *testing.T) {
	c.CreateConnection()
	fmt.Println("connection")
}

func TestDynamoClient_AddTransactionData(t *testing.T) {
	c.AddTransactionData("dev-deployment","testid", "testdata", "PENDING")
	fmt.Println("add")
}

func TestDynamoClient_GetTransactionData(t *testing.T) {
	data, status, _ := c.GetTransactionData("dev-deployment", "testid")
	fmt.Println(*data)
	fmt.Println(*status)
}

func TestDynamoClient_UpdateTransactionData(t *testing.T) {
	c.UpdateTransactionData("dev-deployment", "testid", "updatedata", "DONE")
	fmt.Println("update")
}