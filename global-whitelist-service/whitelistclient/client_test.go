// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package whitelistclient

import (
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateWhiteListParticipants(t *testing.T) {
	wlc := Client{
		HTTPClient: &http.Client{Timeout: time.Second * 10},
		WLURL:      "http://localhost:11234/v1",
	}
	err := wlc.CreateWhiteListParticipants("hk.one.payments.worldwire.io", "test.com")
	Convey("Successful get caller identity", t, func() {
		So(err, ShouldBeNil)
	})
}
func TestGetWhiteListParticipantDomains(t *testing.T) {
	wlc := Client{
		HTTPClient: &http.Client{Timeout: time.Second * 10},
		WLURL:      "http://localhost:11234/v1",
	}
	wlp, _ := wlc.GetWhiteListParticipantDomains("hk.one.payments.worldwire.io")
	Convey("Successful get caller identity", t, func() {
		So(wlp, ShouldContain, "test.com")
	})
}

func TestGetWhiteListParticipants(t *testing.T) {
	wlc := Client{
		HTTPClient: &http.Client{Timeout: time.Second * 10},
		WLURL:      "http://localhost:11234/v1",
	}
	wlp, _ := wlc.GetWhiteListParticipants("hk.one.payments.worldwire.io")
	Convey("Successful get caller identity", t, func() {
		So(wlp, ShouldNotContain, "test.com")
	})
}

func TestDeleteWhiteListParticipants(t *testing.T) {
	wlc := Client{
		HTTPClient: &http.Client{Timeout: time.Second * 10},
		WLURL:      "http://localhost:11234/v1",
	}
	err := wlc.DeleteWhiteListParticipants("hk.one.payments.worldwire.io", "test.com")
	Convey("Successful get caller identity", t, func() {
		So(err, ShouldBeNil)
	})
}

func TestGetMutualParicipantDomain(t *testing.T) {
	wlc := Client{
		HTTPClient: &http.Client{Timeout: time.Second * 10},
		WLURL:      "http://localhost:11234/v1",
	}
	_ = wlc.CreateWhiteListParticipants("test1.com", "test2.com")
	_ = wlc.CreateWhiteListParticipants("test2.com", "test1.com")
	_ = wlc.CreateWhiteListParticipants("test1.com", "test3.com")
	_ = wlc.CreateWhiteListParticipants("test3.com", "test1.com")
	wl, _ := wlc.GetMutualWhiteListParticipantDomains("test1.com")
	Convey("Successful get caller identity", t, func() {
		So(wl, ShouldContain, "test2.com")
		So(wl, ShouldContain, "test3.com")
	})
	wl2, _ := wlc.GetMutualWhiteListParticipantDomains("test2.com")
	Convey("Successful get caller identity", t, func() {
		So(wl2, ShouldContain, "test1.com")
		So(wl2, ShouldNotContain, "test3.com")
	})
	_ = wlc.DeleteWhiteListParticipants("test1.com", "test2.com")
	_ = wlc.DeleteWhiteListParticipants("test1.com", "test3.com")
	_ = wlc.DeleteWhiteListParticipants("test3.com", "test1.com")
	_ = wlc.DeleteWhiteListParticipants("test2.com", "test1.com")
}

// TODO:
func TestGetMutualParicipants(t *testing.T) {
}
