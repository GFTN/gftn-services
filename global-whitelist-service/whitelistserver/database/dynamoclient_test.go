// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package database

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func DynamoClientBuilder() DynamoClient {
	dc := DynamoClient{
		AWS_SECRET_KEY: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWS_KEY_ID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		TableName:      "Whitelist",
		// ProfileFile: "../awsprofile",
		// ProfileName: "test-account",
		Region: "ap-southeast-1",
	}
	dc.CreateConnection()
	return dc
}

// func TestCreateTable(t *testing.T) {
// 	dc := DynamoClientBuilder()
// 	dc.CreateTable()
// 	Convey("Successful get caller identity", t, func() {
// 	})
// }

// func TestDeleteTable(t *testing.T) {
// 	dc := DynamoClientBuilder()
// 	dc.DeleteTable()
// 	Convey("Successful get caller identity", t, func() {
// 	})
// }

func TestAddWhitelistParticipant(t *testing.T) {
	dc := DynamoClientBuilder()
	err := dc.AddWhitelistParticipant("hk", "test3")
	Convey("Successful get caller identity", t, func() {
		So(err, ShouldBeNil)
	})
}

func TestGetWhiteListParicipants(t *testing.T) {
	dc := DynamoClientBuilder()
	whitelist, _ := dc.GetWhiteListParicipants("hk")
	Convey("Successful get caller identity", t, func() {
		So(whitelist, ShouldContain, "test3")
	})
}

func TestDeleteWhitelistParticipant(t *testing.T) {
	dc := DynamoClientBuilder()
	err := dc.DeleteWhitelistParticipant("hk", "test3")
	Convey("Successful get caller identity", t, func() {
		So(err, ShouldBeNil)
		whitelist, _ := dc.GetWhiteListParicipants("hk")
		So(whitelist, ShouldNotContain, "test3")
	})
}
