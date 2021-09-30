// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	authtesting "github.com/GFTN/gftn-services/utility/testing"
)

// type postPayoutResponse struct {
// 	ID string
// }

// var payoutResp []postPayoutResponse

// func executeRequest(req *http.Request) *httptest.ResponseRecorder {
// 	rr := httptest.NewRecorder()
// 	APP.Router.ServeHTTP(rr, req)
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
// 	os.Setenv(global_environment.ENV_KEY_SERVICE_PORT, "11111")
// 	os.Setenv(global_environment.ENV_KEY_SERVICE_ERROR_CODES_FILE, common.Abs("error-codes/payoutErrorCodes.toml"))
// 	os.Setenv(environment.ENV_KEY_PAYOUT_DB_NAME, "payout")
// 	os.Setenv(environment.ENV_KEY_LOCATION_DB_TABLE, "locations")
// 	os.Setenv(environment.ENV_KEY_DB_USER, "wwUserAdmin")
// 	os.Setenv(environment.ENV_KEY_DB_PWD, "wwgftn!*@")
// 	os.Setenv(environment.ENV_KEY_DB_TIMEOUT, "10")

// }

func TestAuthForExternalEndpoint(t *testing.T) {
	a := App{}
	a.initializeRoutes()
	Convey("Testing authorization for external endpoints...", t, func() {
		authtesting.InitAuthTesting()
		err := a.Router.Walk(authtesting.AuthWalker)
		So(err, ShouldBeNil)
	})
}

// // TestServiceCheck - Test service check endpoint exists
// func TestServiceCheck(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/service_check", nil)
// 	response := executeRequest(req)
// 	checkResponseCode(t, http.StatusOK, response.Code)
// }

// // tests for post endpoint
// func TestPostPayoutPoint(t *testing.T) {
// 	//Test posting invalid payout point to receive failure (3 length country code)
// 	Convey("Error code is returned as expected for post", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/invalid-payout-point-payload.json"))
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/client/payout", bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "PAYOUT-1002")
// 	})

// 	//Test posting valid payout point to receive success
// 	Convey("Successful creation of payout point", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-payout-point-payload.json"))
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/client/payout", bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		payoutResp = append(payoutResp, postPayoutResponse{})
// 		json.Unmarshal(response.Body.Bytes(), &payoutResp[0])
// 		So(response.Code, ShouldEqual, http.StatusOK)

// 	})

// 	//Test posting same valid payout point to receive failure
// 	Convey("Successful creation of payout point", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-payout-point-payload.json"))
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/client/payout", bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(msg.Code, ShouldEqual, "PAYOUT-1009")

// 	})

// 	//Test posting second valid payout point to receive success
// 	Convey("Successful creation of payout point", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-payout-point-payload2.json"))
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/client/payout", bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		payoutResp = append(payoutResp, postPayoutResponse{})
// 		json.Unmarshal(response.Body.Bytes(), &payoutResp[1])
// 		So(response.Code, ShouldEqual, http.StatusOK)

// 	})
// }

// // tests for update endpoint
// func TestUpdatePayoutPoint(t *testing.T) {

// 	// invalid update payload to receive failure (id not found)
// 	Convey("Error code is returned as expected due to target id not found", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/invalid-payout-point-update-payload.json"))
// 		req, _ := http.NewRequest("PATCH", "/"+serviceVersion+"/client/payout", bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(msg.Code, ShouldEqual, "PAYOUT-1003")
// 	})

// 	// invalid update payload to receive failure (try to concat a non-array attribute)
// 	Convey("Error code is returned as expected due to url is not an array", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/invalid-payout-point-update-payload.json"))

// 		var updatePayload model.PayoutPointUpdateRequest
// 		json.Unmarshal(payload, &updatePayload)
// 		*updatePayload.ID = payoutResp[0].ID
// 		b, _ := json.Marshal(updatePayload)

// 		req, _ := http.NewRequest("PATCH", "/"+serviceVersion+"/client/payout", bytes.NewBuffer(b))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(msg.Code, ShouldEqual, "PAYOUT-1010")
// 	})

// 	// check no payout point provide a mobile receive mode
// 	Convey("We don't have a payout point providing mobile receive mode", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/payout?receive_mode=mobile", nil)
// 		response := executeRequest(req)
// 		responsePp := []model.PayoutPoint{}
// 		json.Unmarshal(response.Body.Bytes(), &responsePp)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(len(responsePp), ShouldEqual, 0)

// 	})

// 	//update a payout point to make him provide mobile receive mode
// 	Convey("Update all array attribute", t, func() {
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-payout-point-update-payload.json"))

// 		var updatePayload model.PayoutPointUpdateRequest
// 		json.Unmarshal(payload, &updatePayload)
// 		*updatePayload.ID = payoutResp[0].ID
// 		b, _ := json.Marshal(updatePayload)

// 		req, _ := http.NewRequest("PATCH", "/"+serviceVersion+"/client/payout", bytes.NewBuffer(b))
// 		response := executeRequest(req)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})

// 	//Test if anyone has provided mobile receive mode now
// 	Convey("Now we should have a payout point provide mobile receive mode", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/payout?receive_mode=mobile", nil)
// 		response := executeRequest(req)
// 		responsePp := []model.PayoutPoint{}
// 		json.Unmarshal(response.Body.Bytes(), &responsePp)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		So(responsePp[0].ID, ShouldEqual, payoutResp[0].ID)
// 		So(len(responsePp), ShouldEqual, 1)

// 	})

// }

// // tests for get endpoint
// func TestGetPayoutPoint(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	//Test invalid receive mode
// 	Convey("Error code is returned as expected for invalid receive mode", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/payout?receive_mode=test", nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(msg.Code, ShouldEqual, "PAYOUT-1008")
// 	})

// 	//Test agency_pickup receive mode & check if there are two results & id is correct
// 	Convey("Successful getting a payout point from service", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/payout?receive_mode=agency_pickup", nil)
// 		response := executeRequest(req)
// 		responsePp := []model.PayoutPoint{}
// 		json.Unmarshal(response.Body.Bytes(), &responsePp)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		So(responsePp[0].ID, ShouldEqual, payoutResp[0].ID)
// 		So(responsePp[1].ID, ShouldEqual, payoutResp[1].ID)

// 		So(len(responsePp), ShouldEqual, 2)

// 	})

// 	//Test delivery receive mode & check if there is only 1 result & id is correct
// 	Convey("Successful getting a payout point from service", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/payout?receive_mode=delivery", nil)
// 		response := executeRequest(req)
// 		responsePp := []model.PayoutPoint{}
// 		json.Unmarshal(response.Body.Bytes(), &responsePp)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		So(responsePp[0].ID, ShouldEqual, payoutResp[0].ID)
// 		So(len(responsePp), ShouldEqual, 1)

// 	})

// 	//Test get all payout point with no spcified criteria
// 	Convey("Successful getting a payout point from service", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/payout", nil)
// 		response := executeRequest(req)
// 		responsePp := []model.PayoutPoint{}
// 		json.Unmarshal(response.Body.Bytes(), &responsePp)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		So(len(responsePp), ShouldEqual, 2)

// 	})

// }

// // tests for delete endpoint
// func TestDeletePayoutPoint(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()

// 	//specify no id
// 	Convey("Try deleting with no id specified, should failed", t, func() {
// 		req, _ := http.NewRequest("DELETE", "/"+serviceVersion+"/client/payout", nil)
// 		response := executeRequest(req)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)

// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(msg.Code, ShouldEqual, "PAYOUT-1011")
// 	})

// 	// delete first payout point
// 	Convey("Successful delete payout point", t, func() {
// 		req, _ := http.NewRequest("DELETE", "/"+serviceVersion+"/client/payout?id="+payoutResp[0].ID, nil)

// 		response := executeRequest(req)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})

// 	// delete second payout point
// 	Convey("Successful delete payout point", t, func() {
// 		req, _ := http.NewRequest("DELETE", "/"+serviceVersion+"/client/payout?id="+payoutResp[1].ID, nil)

// 		response := executeRequest(req)
// 		So(response.Code, ShouldEqual, http.StatusOK)

// 	})

// }
