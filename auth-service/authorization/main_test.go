// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	authutility "github.com/GFTN/gftn-services/auth-service/authorization/authutility"
	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/utility/wwfirebase"
)

/*
 * initAppForTest : Redundant code for all tests needs to be added here.
 */
func initAppForTest() {
	SetEnvVariables()
	APP.Initialize()
	FbClient, FbAuthClient, _ := wwfirebase.AuthenticateWithAdminPrivileges()
	wwfirebase.FbClient = FbClient
	wwfirebase.FbAuthClient = FbAuthClient
	wwfirebase.FbRef = wwfirebase.GetRootRef()
	APP.initializeRoutes()

}

/*
 *Custom function to execute the unit tests
 */
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	/*
		//These 2 lines were flagged by Checkmarx for SSRF attack. So commenting them out.
		req.Header.Set("Origin", "http://worldwire.com")
		req.Header.Set("Access-Control-Request-Method", req.Method)
	*/
	if APP.HTTPHandler == nil {
		APP.Router.ServeHTTP(rr, req)
	} else {
		APP.HTTPHandler(APP.Router).ServeHTTP(rr, req)
	}
	return rr
}

/*
 * Checks the response code
 */
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

/*
 * Service check for auth service
 */
func TestServiceCheck(t *testing.T) {
	initAppForTest()

	req, _ := http.NewRequest("GET", "/auth/service_check", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

}

/*
 *This checks if the user accesses a non existing endpoint, they receive the error message saying that
 *it doesn't exist
 */
func TestNotFound(t *testing.T) {
	initAppForTest()

	Convey("Error code is returned as expected", t, func() {
		req, _ := http.NewRequest("GET", "/nakul-is-awesome", nil)
		response := executeRequest(req)
		msg := model.WorldWireError{}
		json.Unmarshal(response.Body.Bytes(), &msg)
		So(response.Code, ShouldEqual, http.StatusNotFound)
		So(msg.Code, ShouldEqual, "WW-001")
	})
}

func TestParticipantAuthForMaker(t *testing.T) {
	initAppForTest()

	Convey("Participant should be authorized for requesting access to an endpoint", t, func() {
		authorized, _ := middlewares.CheckAccess("Participant_permissions", "admin", true, "POST", "/v1/anchor/trust")
		So(authorized, ShouldEqual, true)
	})
}

func TestSuperUserAuthForMaker(t *testing.T) {
	initAppForTest()

	Convey("SuperUser should be authorized for requesting access to an endpoint", t, func() {
		// authorizedGet, _ := middlewares.CheckAccess("Super_Permissions", "admin", true, "GET", "")
		authorizedPost, _ := middlewares.CheckAccess("Super_permissions", "admin", true, "POST", "/v1/admin/anchor")
		// So(authorizedGet, ShouldEqual, true)
		So(authorizedPost, ShouldEqual, true)
	})
}

// ExpiredFID
/*
 * This checks if an Expired FID responds correctly from token.DecodeFID
 */
func TestExpiredFID(t *testing.T) {
	initAppForTest()

	Convey("Expired FID should not be decoded", t, func() {
		user, err := authutility.DecodeFID(ExpiredFID)
		So(user, ShouldEqual, "")
		So(err, ShouldNotEqual, nil)
	})
}

func TestSuperDirectAccess(t *testing.T) {
	initAppForTest()

	Convey("Direct Access test for super User", t, func() {

		boolDirectSuperAdmin, _ := middlewares.CheckAccess("Super_permissions", "admin", false, "POST", "/v1/admin/payout")
		boolDirectParticipant, _ := middlewares.CheckAccess("Participant_permissions", "admin", false, "POST", "/v1/admin/payout")

		So(boolDirectSuperAdmin, ShouldEqual, true)
		So(boolDirectParticipant, ShouldEqual, false)

	})
}

func TestParticipantDirectAccess(t *testing.T) {
	initAppForTest()

	Convey("Direct Access test for participant", t, func() {
		boolDirectParticipant, _ := middlewares.CheckAccess("Participant_permissions", "admin", false, "POST", "/v1/client/participants")
		boolDirectParticipantFalse, _ := middlewares.CheckAccess("Participant_permissions", "admin", false, "POST", "/v1/client/assets")

		So(boolDirectParticipant, ShouldEqual, true)
		So(boolDirectParticipantFalse, ShouldEqual, false)

	})
}

// running maker checker will lead to hitting the request
// func TestMakerCheckerFlow(t *testing.T) {
// 	initAppForTest()

// 	Convey("Complete maker checker flow for 1 participant", t, func() {
// 		requestID, errRequest := token.MakerRequest("participant4", IIDApprover, "/", "POST", "participant", UserIDRequester)
// 		check, errApprove := token.CheckerApprove(requestID, IIDApprover, UserIDApprover, "participant")
// 		So(requestID, ShouldNotEqual, "")
// 		So(errRequest, ShouldEqual, nil)
// 		So(errApprove, ShouldEqual, nil)
// 		So(check, ShouldEqual, true)

// 	})

// }
