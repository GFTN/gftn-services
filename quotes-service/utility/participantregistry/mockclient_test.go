// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package participantregistry

import (
	. "github.com/smartystreets/goconvey/convey"
	comn "github.com/GFTN/gftn-services/utility/common"
	"testing"
)

func TestGetAllParticipantsMock(t *testing.T) {
	prc, _ := CreateMockClient()
	participants, _ := prc.GetAllParticipants()
	LOGGER.Debug(participants)
	Convey("Get All participants", t, func(c C) {
		So(*participants[0].ID, ShouldEqual, "hk.one.payments.worldwire.io")
	})
}

func TestGetParticipantForDomain(t *testing.T) {
	prc, _ := CreateMockClient()
	participantID := "hk.one.payments.worldwire.io"
	participants, _ := prc.GetParticipantForDomain(participantID)
	LOGGER.Debug(participants)
	Convey("Get All participants", t, func(c C) {
		So(*participants.ID, ShouldEqual, "hk.one.payments.worldwire.io")
	})
}

func TestGetParticipantAccount(t *testing.T) {
	prc, _ := CreateMockClient()
	Convey("Get All participants", t, func(c C) {
		participantID := "hk.one.payments.worldwire.io"
		account := comn.ISSUING
		address, _ := prc.GetParticipantAccount(participantID, account)
		So(address, ShouldEqual, "GA3Z5DS6GAPBI6EGRFRCEJKAFBZHESW62U5ME3OJYVD5VEREY5ENTGIK")
	})
	Convey("Get All participants", t, func(c C) {
		participantID := "hk.one.payments.worldwire.io"
		account := "default"
		address, _ := prc.GetParticipantAccount(participantID, account)
		So(address, ShouldEqual, "GC42RQF6NC7ZQ4JM4EBZHYBDZRKBCVGKQTD5XY2WRKJSRF2ZJUKERV7T")
	})
}
