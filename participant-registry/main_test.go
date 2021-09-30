// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	authtesting "github.com/GFTN/gftn-services/utility/testing"
)

// func executeRequest(req *http.Request) *httptest.ResponseRecorder {
// 	rr := httptest.NewRecorder()
// 	a.Router.ServeHTTP(rr, req)
// 	return rr
// }

// func checkResponseCode(t *testing.T, expected, actual int) {
// 	if expected != actual {
// 		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
// 	}
// }

// func setEnvVariables() {
// 	os.Setenv(global_environment.ENV_KEY_SERVICE_LOG_FILE, common.Abs("log.txt"))
// 	os.Setenv(global_environment.ENV_KEY_SERVICE_VERSION, "v1")
// 	os.Setenv(global_environment.ENV_KEY_SERVICE_PORT, "1080")
// 	os.Setenv(global_environment.ENV_KEY_SERVICE_INTERNAL_PORT, "9080")
// 	os.Setenv(global_environment.ENV_KEY_SERVICE_ERROR_CODES_FILE, common.Abs("error-codes/prErrorCodes.toml"))
// 	os.Setenv(environment.ENV_KEY_PR_DB_NAME, "test-registry")
// 	os.Setenv(environment.ENV_KEY_PARTICIPANTS_DB_TABLE, "test-participants")
// 	os.Setenv(environment.ENV_KEY_DB_USER, "wwTestUser")
// 	os.Setenv(environment.ENV_KEY_DB_PWD, "wwTestpwd")
// 	os.Setenv(environment.ENV_KEY_IS_UNIT_TEST, "true")
// 	os.Setenv(environment.ENV_KEY_DB_TIMEOUT, "10000")

// }

func TestAuthForExternalEndpoint(t *testing.T) {
	a := App{}
	a.initRoutes()
	Convey("Testing authorization for external endpoints...", t, func() {
		authtesting.InitAuthTesting()
		err := a.Router.Walk(authtesting.AuthWalker)
		So(err, ShouldBeNil)
		err = a.OnboardingRouter.Walk(authtesting.AuthWalker)
		So(err, ShouldBeNil)
	})
}

// TestServiceCheck - Test service check endpoint exists
// func TestServiceCheck(t *testing.T) {
// 	setEnvVariables()
// 	a.InitApp()
// 	req, _ := http.NewRequest("GET", "/"+serviceVersion+"/admin/service_check", nil)
// 	response := executeRequest(req)
// 	checkResponseCode(t, http.StatusOK, response.Code)
// }

// // TestPostPR tests for /internal/pr post endpoint
// func TestPostPR(t *testing.T) {
// 	//Test invalid Participant to receive empty response
// 	Convey("Error code is returned as expected for post", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/invalid-participant1.json"))
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/pr", bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "PR-1001")
// 	})

// 	//Test posting valid participant to receive success
// 	Convey("Successful creation of participant in PR", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-participant1.json"))
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/pr", bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		requestPr := model.Participant{}
// 		json.Unmarshal(response.Body.Bytes(), &requestPr)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		payload, _ = ioutil.ReadFile(common.Abs("./unit-test-data/valid-participant2.json"))
// 		req, _ = http.NewRequest("POST", "/"+serviceVersion+"/internal/pr", bytes.NewBuffer(payload))
// 		response = executeRequest(req)
// 		json.Unmarshal(response.Body.Bytes(), &requestPr)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})

// 	//Test invalid Participant to receive empty response
// 	Convey("Error code is returned as expected for post", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/invalid-participant2.json"))
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/pr", bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "PR-1001")
// 	})

// 	//Test invalid Participant to receive empty response, unique Participant
// 	Convey("Error code is returned as expected for post", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/invalid-participant3.json"))
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/pr", bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "PR-1001")
// 	})
// }

// // TestGetPRList tests for /internal/pr endpoint
// func TestGetPRList(t *testing.T) {

// 	//Test get PR list to retrieve expected participants
// 	Convey("", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/pr", nil)
// 		response := executeRequest(req)
// 		prList := []model.Participant{}
// 		expectedPr := model.Participant{}
// 		json.Unmarshal(response.Body.Bytes(), &prList)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-participant1.json"))
// 		json.Unmarshal(payload, &expectedPr)
// 		So(*prList[0].ID, ShouldEqual, *expectedPr.ID)
// 		payload, _ = ioutil.ReadFile(common.Abs("./unit-test-data/valid-participant2.json"))
// 		json.Unmarshal(payload, &expectedPr)
// 		So(*prList[1].ID, ShouldEqual, *expectedPr.ID)
// 	})

// }

// // TestGetPRDomain tests for /internal/pr/domain/{participant_id} endpoint
// func TestGetPRDomain(t *testing.T) {
// 	//Test invalid domain to receive empty response
// 	Convey("Error code is returned as expected for get domain", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/pr/domain/test", nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(msg.Code, ShouldEqual, "PR-1002")
// 	})

// 	//Test valid domain to receive success response
// 	Convey("Error code is returned as expected for get domain", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/pr/domain/xyz.xyz.payments.gftn.io", nil)
// 		response := executeRequest(req)
// 		requestPr := model.Participant{}
// 		expectedPr := model.Participant{}
// 		json.Unmarshal(response.Body.Bytes(), &requestPr)
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-participant1.json"))
// 		json.Unmarshal(payload, &expectedPr)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		So(*requestPr.ID, ShouldEqual, *expectedPr.ID)
// 		So(requestPr.IssuingAccount, ShouldEqual, expectedPr.IssuingAccount)
// 		So(requestPr.Status, ShouldEqual, "inactive")
// 	})
// }

// // TestPostIssuingAccount tests for /internal/pr/issuingaccount/{participant_id} endpoint
// func TestPostIssuingAccount(t *testing.T) {

// 	//Test invalid issuing account payload to receive error response
// 	Convey("Error code is returned as expected for invalid issue account post", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/invalid-account.json"))
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/pr/issuingaccount/xyz.xyz.payments.gftn.io", bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "PR-1001")
// 	})

// 	//Test posting valid issuing account payload to receive success
// 	Convey("success code is returned as expected for valid issue account post", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-account.json"))
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/pr/issuingaccount/xyz.xyz.payments.gftn.io", bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})

// 	Convey("Error code is returned as expected for repeated issue account post", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-account.json"))
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/pr/issuingaccount/xyz.xyz.payments.gftn.io", bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		So(response.Code, ShouldEqual, http.StatusConflict)
// 	})

// }

// // TestPostOperatingAccount tests for /internal/pr/account/{participant_id} endpoint
// func TestPostOperatingAccount(t *testing.T) {

// 	//Test invalid issuing account payload to receive error response
// 	Convey("Error code is returned as expected for invalid operating account post", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/invalid-dist-account.json"))
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/pr/account/xyz.xyz.payments.gftn.io", bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "PR-1001")
// 	})

// 	//Test posting valid issuing account payload to receive success
// 	Convey("Success is returned as expected for valid operating account post", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-dist-account.json"))
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/pr/account/xyz.xyz.payments.gftn.io", bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})
// 	//Test re-posting valid issuing account payload to receive failure
// 	Convey("Error code is returned as expected for invalid issue account post", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-dist-account.json"))
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/pr/account/xyz.xyz.payments.gftn.io", bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		So(response.Code, ShouldEqual, http.StatusConflict)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(msg.Code, ShouldEqual, "PR-1004")
// 	})

// }

// // TestGetOperatingAccount tests for /internal/pr/account/{participant_id}/{account_name} endpoint
// func TestGetOperatingAccount(t *testing.T) {

// 	//Test invalid issuing account payload to receive error response
// 	Convey("Error code is returned as expected for invalid operating account get", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/pr/account/xyz.xyz.payments.gftn.io/", nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(msg.Code, ShouldEqual, "WW-001")
// 	})

// 	//Test posting valid issuing account payload to receive success
// 	Convey("Error code is returned as expected for invalid issue account post", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/pr/account/xyz.xyz.payments.gftn.io/default", nil)
// 		response := executeRequest(req)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})

// 	//Test re-posting valid issuing account payload to receive failure
// 	Convey("Error code is returned as expected for invalid issue account post", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/pr/account/xyz.xyz.payments.gftn.io/test", nil)
// 		response := executeRequest(req)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(msg.Code, ShouldEqual, "PR-1005")
// 	})

// }

// // TestGetPRDomain tests for /internal/pr/country/{country_code} endpoint
// func TestGetPRCountry(t *testing.T) {
// 	//Test empty domain to receive empty response
// 	Convey("Error code is returned as expected for get pr by country", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/pr/country/", nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(msg.Code, ShouldEqual, "WW-001")
// 	})

// 	//Test invalid domain to receive empty response
// 	Convey("Error code is returned as expected for get pr by country", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/pr/country/test", nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(msg.Code, ShouldEqual, "PR-1002")
// 	})

// 	//Test valid domain to receive success response
// 	Convey("Error code is returned as expected for get domain", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/pr/country/NZL", nil)
// 		response := executeRequest(req)
// 		requestPr := []model.Participant{}
// 		expectedPr := model.Participant{}
// 		json.Unmarshal(response.Body.Bytes(), &requestPr)
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-participant1.json"))
// 		json.Unmarshal(payload, &expectedPr)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		So(*requestPr[0].ID, ShouldEqual, *expectedPr.ID)
// 		So(*requestPr[0].CountryCode, ShouldEqual, *expectedPr.CountryCode)
// 	})
// }

// // TestPutParticipantStatus tests for /internal/pr/{participant_id}/status endpoint
// func TestPutParticipantStatus(t *testing.T) {
// 	//Test invalid status payload to receive error response
// 	Convey("Error code is returned as expected for invalid status put", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/invalid-participant-status.json"))
// 		req, _ := http.NewRequest("PUT", "/"+serviceVersion+"/internal/pr/xyz.xyz.payments.gftn.io/status", bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "PR-1001")
// 	})

// 	//Test valid status payload to receive success response
// 	Convey("Success returned as expected for valid status put", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-participant-status.json"))
// 		prStatus := model.ParticipantStatus{}
// 		json.Unmarshal(payload, &prStatus)
// 		req, _ := http.NewRequest("PUT", "/"+serviceVersion+"/internal/pr/xyz.xyz.payments.gftn.io/status", bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		req, _ = http.NewRequest("GET", "/"+serviceVersion+"/internal/pr/domain/xyz.xyz.payments.gftn.io", nil)
// 		response = executeRequest(req)
// 		requestPr := model.Participant{}
// 		json.Unmarshal(response.Body.Bytes(), &requestPr)
// 		So(requestPr.Status, ShouldEqual, *prStatus.Status)

// 	})
// }
