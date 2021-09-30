// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package pr_client

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMockGetAllParticipants(t *testing.T) {
	prc, _ := CreateRestPRServiceClient("http://localhost:8081/v1")
	participants, _ := prc.GetAllParticipants()
	fmt.Println(participants)
	Convey("Successful get caller identity", t, func() {
	})
}
