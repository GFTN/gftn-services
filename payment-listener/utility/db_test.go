// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package utility

import (
	"fmt"
	"os"
	"testing"
)

var c *DynamoClient

func TestInit(t *testing.T) {
	os.Setenv("AWS_DYNAMO_REGION", "ap-southeast-1")
	os.Setenv("DYNAMO_DB_TABLE_NAME", "payment-test")
	os.Setenv("ENV_VERSION", "dev")
	os.Setenv("HOME_DOMAIN_NAME", "participant3")

}

func TestCreate(t *testing.T) {
	var err error
	c, err = CreateConnection()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("connection")
}

func TestDynamoClient_AddTransactionData(t *testing.T) {

	err := c.AddAccountCursor("issuing", "now")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("add")
}

func TestDynamoClient_GetTransactionData(t *testing.T) {
	cursor, err := c.GetAccountCursor("issuing")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("********* cursor is :%v \n", *cursor)
}

func TestDynamoClient_UpdateTransactionData(t *testing.T) {
	err := c.UpdateAccountCursor("issuing", "1234")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("update")
}
