// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package prclient

import (
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetAllParticipants(t *testing.T) {
	prc := Client{
		HTTPClient: &http.Client{Timeout: time.Second * 10},
		URL:        "http://localhost:8081/v1",
	}
	participants, _ := prc.GetAllParticipants()
	Convey("Successful get caller identity", t, func() {
		LOGGER.Info(participants)
	})
}
