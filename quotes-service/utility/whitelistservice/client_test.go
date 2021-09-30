// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package whitelistservice

import (
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateWhiteListParticipants(t *testing.T) {
	wlc := Client{
		HTTP:  &http.Client{Timeout: time.Second * 10},
		WLURL: "http://localhost:11234",
	}
	err := wlc.CreateWhiteListParticipants("hk.one.payments.worldwire.io", "test.com")
	Convey("Successful get caller identity", t, func() {
		So(err, ShouldBeNil)
	})
}
func TestGetWhiteListParticipantDomains(t *testing.T) {
	wlc := Client{
		HTTP:  &http.Client{Timeout: time.Second * 10},
		WLURL: "http://localhost:11234",
	}
	wlp, _ := wlc.GetWhiteListParticipantDomains("hk.one.payments.worldwire.io")
	Convey("Successful get caller identity", t, func() {
		So(wlp, ShouldContain, "test.com")
	})
}

func TestGetWhiteListParticipants(t *testing.T) {
	wlc := Client{
		HTTP:  &http.Client{Timeout: time.Second * 10},
		WLURL: "http://localhost:11234",
	}
	wlp, _ := wlc.GetWhiteListParticipants("hk.one.payments.worldwire.io")
	Convey("Successful get caller identity", t, func() {
		So(wlp, ShouldNotContain, "test.com")
	})
}

func TestDeleteWhiteListParticipants(t *testing.T) {
	wlc := Client{
		HTTP:  &http.Client{Timeout: time.Second * 10},
		WLURL: "http://localhost:11234",
	}
	err := wlc.DeleteWhiteListParticipants("hk.one.payments.worldwire.io", "test.com")
	Convey("Successful get caller identity", t, func() {
		So(err, ShouldBeNil)
	})
}
