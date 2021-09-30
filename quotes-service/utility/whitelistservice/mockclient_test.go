// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package whitelistservice

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func (client *MockClient) TestGetWhiteListParticipants(t *testing.T) {
	wlClient := MockClient{
		HTTP:  &http.Client{},
		WLURL: "test",
	}
	participantOFIID := "hk.one.payments.worldwire.io"
	wlParticipants, _ := wlClient.GetWhiteListParticipants(participantOFIID)
	Convey("Get All participants", t, func(c C) {
		So(*wlParticipants[1].ID, ShouldEqual, "ie.one.payments.worldwire.io")
	})

	participantOFIID = "ie.one.payments.worldwire.io"
	wlParticipants, _ = wlClient.GetWhiteListParticipants(participantOFIID)
	Convey("Get All participants", t, func(c C) {
		So(*wlParticipants[0].ID, ShouldEqual, "hk.one.payments.worldwire.io")
	})
}
