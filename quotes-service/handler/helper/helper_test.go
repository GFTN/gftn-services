// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package helper

import (
	"errors"
	"log"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/GFTN/gftn-services/utility/wwfirebase"
)

func TestExtractCallerIdentity1(t *testing.T) {

	jwtStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ii1MZDMwM0V4a29vTkRRTVBoWWxLLjEtMTgifQ.eyJhY2MiOlsiaXNzdWluZyIsImRlZmF1bHQiLCJhZG1pbiIsIm1hbmFnZXIiLCJ2aWV3ZXIiXSwidmVyIjoiMi45LjMuN19SQzEiLCJpcHMiOlsiMjAyLjEzNS4yNDUuMzkiLCIyMDIuMTM1LjI0NS40IiwiMjAyLjEzNS4yNDUuMiJdLCJlbnYiOiJxYSIsImVucCI6WyIvdjEvYWRtaW4vcHIiLCIvdjEvb25ib2FyZGluZy9hY2NvdW50cyIsIi92MS9vbmJvYXJkaW5nL2lzc3VpbmdhY2NvdW50IiwiL3YxL29uYm9hcmRpbmcvb3BlcmF0aW5nYWNjb3VudCIsIi92MS9jbGllbnQvYXNzZXRzIiwiL3YxL2NsaWVudC9hc3NldHMvaXNzdWVkIiwiL3YxL2NsaWVudC9wYXJ0aWNpcGFudHMvd2hpdGVsaXN0IiwiL3YxL2NsaWVudC90cnVzdCIsIi92MS9jbGllbnQvYXNzZXRzL2FjY291bnRzIiwiL3YxL2NsaWVudC9hY2NvdW50cyIsIi92MS9jbGllbnQvdHJhbnNhY3Rpb25zL3NlbmQiLCIvdjEvY2xpZW50L3RyYW5zYWN0aW9ucy9yZWNlaXZlIiwiL3YxL2NsaWVudC90cmFuc2FjdGlvbnMiLCIvdjEvY2xpZW50L3BhcnRpY2lwYW50cyIsIi92MS9jbGllbnQvZmVlcyIsIi92MS9jbGllbnQvYmFsYW5jZXMvYWNjb3VudHMiLCIvdjEvY2xpZW50L2JhbGFuY2VzL2RvIiwiL3YxL2NsaWVudC9zaWduIiwiL3YxL2NsaWVudC9xdW90ZXMvcmVxdWVzdCIsIi92MS9jbGllbnQvcXVvdGVzIiwiL3YxL2NsaWVudC9leGNoYW5nZSIsIi92MS9jbGllbnQvdHJhbnNhY3Rpb25zL3NldHRsZS9kbyIsIi92MS9jbGllbnQvdHJhbnNhY3Rpb25zL3NldHRsZS9kYSIsIi92MS9jbGllbnQvcGF5b3V0IiwiL3YxL2NsaWVudC9hY2NvdW50cyIsIi92MS9jbGllbnQvYXNzZXRzL3BhcnRpY2lwYW50cyJdLCJuIjowLCJpYXQiOjE1NTU5MTkyOTEsIm5iZiI6MTU1NTkxOTI5NiwiZXhwIjoxNTU2MTc4NDkxLCJhdWQiOiJwYXJ0aWNpcGFudDIiLCJzdWIiOiItTFBtckYwUWdnNmZzaHpMYzQ4aSIsImp0aSI6Ii1MZDMwMW9Dc0ptRlBsQkViUlR5In0.123467890oiuytre"
	request, _ := http.NewRequest("GET", "", nil)
	request.Header.Set("Authorization", jwtStr)
	Convey("Successful get caller identity", t, func() {
		participantDomain, _ := GetIdentity(request)
		So(participantDomain, ShouldEqual, "participant2")
	})
}

func TestExtractCallerIdentity2(t *testing.T) {

	jwtStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ii1MZDMwQVp2TjU3RDZ1YXpQcElxLjEtMCJ9.eyJhY2MiOlsiaXNzdWluZyIsImRlZmF1bHQiLCJhZG1pbiIsIm1hbmFnZXIiLCJ2aWV3ZXIiXSwidmVyIjoiMi45LjMuN19SQzEiLCJpcHMiOlsiMjAyLjEzNS4yNDUuMzkiLCIyMDIuMTM1LjI0NS40IiwiMjE5Ljc0LjE1LjIwNyIsIjIwMi4xMzUuMjQ1LjIiXSwiZW52IjoicWEiLCJlbnAiOlsiL3YxL2FkbWluL3ByIiwiL3YxL2FkbWluL3ByL2RvbWFpbiIsIi92MS9hZG1pbi9wci9jb3VudHJ5IiwiL3YxL29uYm9hcmRpbmcvYWNjb3VudHMiLCIvdjEvb25ib2FyZGluZy9pc3N1aW5nYWNjb3VudCIsIi92MS9vbmJvYXJkaW5nL29wZXJhdGluZ2FjY291bnQiLCIvdjEvY2xpZW50L2Fzc2V0cyIsIi92MS9jbGllbnQvcGFydGljaXBhbnRzL3doaXRlbGlzdCIsIi92MS9jbGllbnQvdHJ1c3QiLCIvdjEvY2xpZW50L2Fzc2V0cy9hY2NvdW50cyIsIi92MS9jbGllbnQvYWNjb3VudHMiLCIvdjEvY2xpZW50L3RyYW5zYWN0aW9ucy9zZW5kIiwiL3YxL2NsaWVudC90cmFuc2FjdGlvbnMvcmVjZWl2ZSIsIi92MS9jbGllbnQvdHJhbnNhY3Rpb25zIiwiL3YxL2NsaWVudC9wYXJ0aWNpcGFudHMiLCIvdjEvY2xpZW50L2ZlZXMiLCIvdjEvY2xpZW50L2JhbGFuY2VzL2FjY291bnRzIiwiL3YxL2NsaWVudC9iYWxhbmNlcy9kbyIsIi92MS9jbGllbnQvc2lnbiIsIi92MS9jbGllbnQvcXVvdGVzL3JlcXVlc3QiLCIvdjEvY2xpZW50L3F1b3RlcyIsIi92MS9jbGllbnQvZXhjaGFuZ2UiLCIvdjEvY2xpZW50L3RyYW5zYWN0aW9ucy9zZXR0bGUvZG8iLCIvdjEvY2xpZW50L3RyYW5zYWN0aW9ucy9zZXR0bGUvZGEiLCIvdjEvY2xpZW50L3BheW91dCIsIi92MS9jbGllbnQvc2VydmljZV9jaGVjayIsIi92MS9jbGllbnQvYWNjb3VudHMiLCIvdjEvY2xpZW50L2Fzc2V0cy9wYXJ0aWNpcGFudHMiXSwibiI6MCwiaWF0IjoxNTU1OTE5MzIxLCJuYmYiOjE1NTU5MTkzMjYsImV4cCI6MTU1NjE3ODUyMSwiYXVkIjoicGFydGljaXBhbnQxIiwic3ViIjoiLUxQbXJGMFFnZzZmc2h6TGM0OGkiLCJqdGkiOiItTGQzMDkxdWNwVXozSExtaHVZWiJ9.123467890oiuytre"
	request, _ := http.NewRequest("GET", "", nil)
	request.Header.Set("Authorization", jwtStr)
	Convey("Successful get caller identity", t, func() {
		participantDomain, _ := GetIdentity(request)
		So(participantDomain, ShouldEqual, "participant1")
	})
}

func TestSendLogToFirebase(t *testing.T) {
	err := errors.New("")
	wwfirebase.FbClient, _, err = wwfirebase.AuthenticateWithAdminPrivileges()
	if err != nil {
		log.Fatalln("Unable to authenticate firebase")
	}
	wwfirebase.FbRef = wwfirebase.GetRootRef()
	testObj := map[string]string{"abc": "efg"}
	err = SendLogToFirebase(testObj, "v1/txn/exchange/test1", "v1/txn/exchange/test2")
	Convey("Successful SendLogToFirebase", t, func() {
		So(err, ShouldBeNil)
	})
}

func TestSendLogToFirebase2(t *testing.T) {
	err := errors.New("")
	wwfirebase.FbClient, _, err = wwfirebase.AuthenticateWithAdminPrivileges()
	if err != nil {
		log.Fatalln("Unable to authenticate firebase")
	}
	wwfirebase.FbRef = wwfirebase.GetRootRef()
	testObj := map[string]string{"abc": "efg"}
	err = SendLogToFirebase(testObj)
	Convey("Fail SendLogToFirebase", t, func() {
		So(err, ShouldNotBeNil)
	})
}
