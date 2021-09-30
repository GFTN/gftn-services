// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	authtesting "github.com/GFTN/gftn-services/utility/testing"
)

/*
Global environmental variables for all unit tests
*/

// func setEnvVariables() {
// 	os.Setenv(global_environment.ENV_KEY_SERVICE_LOG_FILE, common.Abs("log.txt"))
// 	os.Setenv(global_environment.ENV_KEY_SERVICE_VERSION, "v1")
// 	os.Setenv(global_environment.ENV_KEY_SERVICE_PORT, "9080")
// 	os.Setenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL, "http://localhost:10080/v1")
// 	os.Setenv(global_environment.ENV_KEY_SERVICE_ERROR_CODES_FILE, common.Abs("error-codes/apiErrorCodes.toml"))
// 	os.Setenv(global_environment.ENV_KEY_HORIZON_CLIENT_URL, "https://horizon-testnet.stellar.org")
// 	os.Setenv(global_environment.ENV_KEY_STELLAR_NETWORK, "Test SDF Network ; September 2015")
// 	os.Setenv(global_environment.ENV_KEY_ADMIN_SVC_URL, "http://pr.gftn.io:10080/v1")
// 	os.Setenv(environment.ENV_KEY_STELLAR_DB_HOST, "18.219.148.107")
// 	os.Setenv(environment.ENV_KEY_STELLAR_DB_PORT, "5432")
// 	os.Setenv(environment.ENV_KEY_STELLAR_DB_NAME, "core_bac")
// 	os.Setenv(environment.ENV_KEY_STELLAR_DB_USER, "stellar")
// 	os.Setenv(environment.ENV_KEY_STELLAR_DB_PASSWORD, "gftn123")
// 	os.Setenv(global_environment.ENV_KEY_ADMIN_SVC_URL, "http://localhost:8083/v1")
// 	os.Setenv(environment.ENV_KEY_ANCHOR_SH_PASS, "seema.sandbox.worldwire")
// 	os.Setenv(environment.ENV_KEY_ANCHOR_SH_SEC, "4MryDmUI9do5NOZaif28UoMraZX1O8P/")
// 	os.Setenv(environment.ENV_KEY_ANCHOR_SH_CRED, "9b13e830-702e-4a60-912f-f7197b674b12")
// 	os.Setenv(environment.ENV_KEY_ANCHOR_SH_VENEU, "f7e35a5a-7f73-4c66-8b90-9d347730ddd9")
// 	os.Setenv(global_environment.ENV_KEY_ENABLE_JWT, "false")
// 	os.Setenv(global_environment.ENV_KEY_FIREBASE_CREDENTIALS, common.Abs("./firebase/_next-gftn-firebase-adminsdk-wvjz8-67ea263932.json"))
// 	os.Setenv(environment.ENV_KEY_STRONGHOLD_ANCHOR_ID, "http://localhost:8888/v1")

// }

/*
Custom function to execute the unit tests
*/
// func executeRequest(req *http.Request) *httptest.ResponseRecorder {
// 	rr := httptest.NewRecorder()
// 	/*
// 		//These 2 lines were flagged by Checkmarx for SSRF attack. So commenting them out.
// 		req.Header.Set("Origin", "http://worldwire.com")
// 		req.Header.Set("Access-Control-Request-Method", req.Method)
// 	*/
// 	if APP.HTTPHandler == nil {
// 		APP.Router.ServeHTTP(rr, req)
// 	} else {
// 		APP.HTTPHandler(APP.Router).ServeHTTP(rr, req)
// 	}
// 	return rr
// }

/*
	Checks the response code
*/

// func checkResponseCode(t *testing.T, expected, actual int) {
// 	if expected != actual {
// 		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
// 	}
// }

func TestAuthForExternalEndpoint(t *testing.T) {
	a := App{}
	a.initializeRoutes()
	Convey("Testing authorization for external endpoints...", t, func() {
		authtesting.InitAuthTesting()
		err := a.Router.Walk(authtesting.AuthWalker)
		So(err, ShouldBeNil)
		err = a.InternalRouter.Walk(authtesting.AuthWalker)
		So(err, ShouldBeNil)
	})
}

/*
	Checks if the api service is running
*/
// func TestServiceCheck(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/service_check", nil)
// 	response := executeRequest(req)
// 	checkResponseCode(t, http.StatusOK, response.Code)
// }

/*
	This checks if the user accesses a non existing endpoint, they receive the error message saying that
	it doesn't exist
*/
// func TestNotFound(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()

// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/a/url/that/doesnt/exist", nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(msg.Code, ShouldEqual, "WW-001")
// 	})
// }

/*
	This endpoint communicates with the callback service to find out if a given account exists
*/
// func TestAccountFind(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()

// 	temp := os.Getenv(global_environment.ENV_KEY_NODE_CONFIG)
// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/account_find/", nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(msg.Code, ShouldEqual, "WW-001")
// 	})
// 	/*
// 		Convey("Successful response received as expected", t, func() {
// 			req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/account_find/account_number1/account_type1", nil)
// 			response := executeRequest(req)
// 			var respMsg model.AccountFoundResponse
// 			json.Unmarshal(response.Body.Bytes(), &respMsg)
// 			So(response.Code, ShouldEqual, http.StatusOK)
// 			So(respMsg.StellarNetworkAddress, ShouldNotBeEmpty)
// 			So(respMsg.OtherInfo, ShouldNotBeEmpty)
// 		})*/

// 	httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
// 	defer httpmock.DeactivateAndReset()

// 	pl := model.VerificationReceipt{AccountStatus: Ptr("OK"), AccountMemo: "Test Memo from Receiver"}
// 	payload, _ := json.Marshal(&pl)
// 	var mockResponse model.VerificationReceipt
// 	json.Unmarshal(payload, &mockResponse)
// 	//simulate response
// 	responder, _ := httpmock.NewJsonResponder(http.StatusOK, mockResponse)
// 	httpmock.RegisterResponder("POST", "http://localhost:8081/v1/callback/verifications/account", responder)

// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("./unit-test-data/nodeconfig/DA.toml"))

// 	Convey("HTTP status code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/account_find/account_number1/account_type1", nil)
// 		response := executeRequest(req)
// 		var actual model.AccountFoundResponse
// 		json.Unmarshal(response.Body.Bytes(), &actual)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		So("FakeOperatingAccountDefaultNodeAddr", ShouldEqual, *actual.StellarNetworkAddress)
// 		So("Test Memo from Receiver", ShouldEqual, actual.OtherInfo)
// 	})

// 	pl2 := model.VerificationReceipt{AccountStatus: Ptr("DENIED"), AccountMemo: "Beneficiary Account is invalid"}
// 	payload2, _ := json.Marshal(&pl2)
// 	var mockResponse2 model.VerificationReceipt
// 	json.Unmarshal(payload2, &mockResponse2)
// 	//simulate response
// 	responder2, _ := httpmock.NewJsonResponder(http.StatusNotFound, mockResponse2)
// 	httpmock.RegisterResponder("POST", "http://localhost:8081/v1/callback/verifications/account", responder2)

// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("./unit-test-data/nodeconfig/DA.toml"))

// 	Convey("HTTP status code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/account_find/account_number1/account_type1", nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1070")
// 	})

// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs(temp))
// }

/*
	This function tests the successful response code for the FundAccount function which does the following
	1. Get operating account & issuing account from nodeconfig.toml
	2. Verify if operating account has trust relationship for the token by querying stellar
	3. If no trust, If no trust, return error (no auto change trust anymore... all trust must be explicit)
*/
/*
func TestFundAccount(t *testing.T) {
   setEnvVariables()
   APP.Initialize()
   APP.initializeRoutes()
   temp := os.Getenv(global_environment.ENV_KEY_NODE_CONFIG)
   //run test of trust internal
   os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, "./unit-test-data/nodeconfig/fund.toml")
   payload := []byte(`{"account_name":"OperatingAccX","asset_code":"FJD","amount":3.0}`)

   Convey("Response code is returned as expected", t, func() {
	   req, _ := http.NewRequest("POST", "/"+serviceVersion+"/client/account/asset/fund", bytes.NewBuffer(payload))
	   response := executeRequest(req)
	   So(response.Code, ShouldEqual, http.StatusOK)
   })

   payloadDO := []byte(`{"account_name":"OperatingAccX","asset_code":"FJDDO","amount":3.0}`)
   Convey("Error code is returned as expected", t, func() {
	   req, _ := http.NewRequest("POST", "/"+serviceVersion+"/client/account/asset/fund", bytes.NewBuffer(payloadDO))
	   response := executeRequest(req)
	   msg := model.WorldWireError{}
	   json.Unmarshal(response.Body.Bytes(), &msg)
	   So(response.Code, ShouldEqual, http.StatusNotFound)
	   So(msg.Code, ShouldEqual, "API-1244")
   })

   payloadmga := []byte(`{"account_name":"OperatingAccX","asset_code":"MGA","amount":3.0}`)
   Convey("Error code is returned as expected", t, func() {
	   req, _ := http.NewRequest("POST", "/"+serviceVersion+"/client/account/asset/fund", bytes.NewBuffer(payloadmga))
	   response := executeRequest(req)
	   msg := model.WorldWireError{}
	   json.Unmarshal(response.Body.Bytes(), &msg)
	   So(response.Code, ShouldEqual, http.StatusNotFound)
	   So(msg.Code, ShouldEqual, "API-1245")
   })

   os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, temp)
}*/

// func TestGetParticipantsForBeneficiary(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()

// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/participant/beneficiary", nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1060")
// 	})

// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/participant/beneficiary", nil)
// 		q := req.URL.Query()
// 		q.Add("account_identifier", "xyz")
// 		req.URL.RawQuery = q.Encode()
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1060")
// 	})

// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/participant/beneficiary", nil)
// 		q := req.URL.Query()
// 		q.Add("country_code", "")
// 		req.URL.RawQuery = q.Encode()
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1061")
// 	})

// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/participant/beneficiary", nil)
// 		q := req.URL.Query()
// 		q.Add("country_code", "xyz")
// 		req.URL.RawQuery = q.Encode()
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1062")
// 	})

// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/participant/beneficiary", nil)
// 		q := req.URL.Query()
// 		q.Add("country_code", "xyz")
// 		q.Add("account_identifier", "")
// 		req.URL.RawQuery = q.Encode()
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1063")
// 	})

// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/participant/beneficiary", nil)
// 		q := req.URL.Query()
// 		q.Add("country_code", "xyz")
// 		q.Add("account_identifier", "abc")
// 		req.URL.RawQuery = q.Encode()
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1064")
// 	})

// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/participant/beneficiary", nil)
// 		q := req.URL.Query()
// 		q.Add("country_code", "xyz")
// 		q.Add("account_identifier", "abc")
// 		q.Add("account_type", "")
// 		req.URL.RawQuery = q.Encode()
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1065")
// 	})
// }

// //TestGetParticipants - unit test for get participants handler gets list of participants from PR
// func TestGetParticipants(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()

// 	Convey("Error code is returned as expected", t, func() {
// 		httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
// 		defer httpmock.DeactivateAndReset()
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-participant-response1.json"))
// 		mockResponse := []model.Participant{}
// 		json.Unmarshal(payload, &mockResponse)
// 		//simulate response
// 		responder, _ := httpmock.NewJsonResponder(http.StatusOK, mockResponse)
// 		httpmock.RegisterResponder("GET", "http://localhost:10080/v1/internal/pr", responder)

// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/participants", nil)
// 		response := executeRequest(req)
// 		expectedPrs := []model.Participant{}
// 		json.Unmarshal(response.Body.Bytes(), &expectedPrs)

// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		So(*expectedPrs[0].ID, ShouldEqual, *mockResponse[0].ID)
// 	})
// }

// func TestGetParticipantByAssetCountry(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()

// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/participants", nil)
// 		q := req.URL.Query()

// 		q.Add("asset_code", "xyz")
// 		q.Add("asset_issuer", "xyz")
// 		req.URL.RawQuery = q.Encode()
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(msg.Code, ShouldEqual, "API-1119")
// 	})

// 	Convey("Testing invalid query parameters", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/participants", nil)
// 		q := req.URL.Query()
// 		q.Add("asset_issuer", "xyz")
// 		q.Add("country_code", "abc")
// 		req.URL.RawQuery = q.Encode()
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotAcceptable)
// 		So(msg.Code, ShouldEqual, "API-1042")

// 		req, _ = http.NewRequest("GET", "/"+serviceVersion+"/client/participants", nil)
// 		q = req.URL.Query()
// 		q.Add("asset_code", "xyz")
// 		q.Add("country_code", "abc")
// 		req.URL.RawQuery = q.Encode()
// 		response = executeRequest(req)
// 		msg = model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotAcceptable)
// 		So(msg.Code, ShouldEqual, "API-1042")

// 		req, _ = http.NewRequest("GET", "/"+serviceVersion+"/client/participants", nil)
// 		q = req.URL.Query()
// 		q.Add("asset_code", "")
// 		q.Add("asset_issuer", "xyz")
// 		q.Add("country_code", "abc")
// 		req.URL.RawQuery = q.Encode()
// 		response = executeRequest(req)
// 		msg = model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotAcceptable)
// 		So(msg.Code, ShouldEqual, "API-1042")

// 		req, _ = http.NewRequest("GET", "/"+serviceVersion+"/client/participants", nil)
// 		q = req.URL.Query()
// 		q.Add("asset_code", "xyz")
// 		q.Add("asset_issuer", "")
// 		q.Add("country_code", "abc")
// 		req.URL.RawQuery = q.Encode()
// 		response = executeRequest(req)
// 		msg = model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotAcceptable)
// 		So(msg.Code, ShouldEqual, "API-1042")
// 	})
// }

// /*
// 	This tests issue asset functionality
// 	The endpoint is called with the asset code.
// 	Asset is represented by the asset code and the account of the issuer of that asset.
// 	1. Issuing account is retrieved from the node config toml. If it doesn't exist return error.
// 	2. IBM Token Account is loaded from the node config toml.
// 	3. Invoke change trust operation on stellar so that IBM account will trust the asset (issuing account & asset code)
// 	4. If the operation is successful return the asset information otherwise return error
// */
// func TestIssueAsset(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	temp := os.Getenv(global_environment.ENV_KEY_NODE_CONFIG)
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, "")
// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/client/assets?asset_code=NZD&asset_type=NONDO", nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1124")
// 	})
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("./unit-test-data/nodeconfig/common/nodeconfig.toml"))
// 	Convey("Status code and response message are returned as expected", t, func() {
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/client/assets?asset_code=AUDDO&asset_type=DO", nil)
// 		response := executeRequest(req)
// 		var msg model.Asset
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		So(*msg.AssetCode, ShouldEqual, "AUDDO")
// 	})
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("./unit-test-data/nodeconfig/sg.nodeconfig.toml"))
// 	Convey("Status code and response message are returned as expected", t, func() {
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/client/assets?asset_code=SGDDO&asset_type=DO", nil)
// 		response := executeRequest(req)
// 		var msg model.Asset
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		So(*msg.AssetCode, ShouldEqual, "SGDDO")
// 	})
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs(temp))
// }

// func TestCORS(t *testing.T) {
// 	setEnvVariables()
// 	temp := os.Getenv(global_environment.ENV_KEY_ORIGIN_ALLOWED)
// 	os.Setenv(global_environment.ENV_KEY_ORIGIN_ALLOWED, "true")
// 	temp = os.Getenv(global_environment.ENV_KEY_ORIGIN_ALLOWED)
// 	APP.Initialize()
// 	APP.initializeRoutes()

// 	Convey("CORS allow origin is set to * when cors env is set to true", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/service_check", nil)
// 		req.Header.Set("Origin", "localhost")
// 		response := executeRequest(req)
// 		header := response.HeaderMap.Get("Access-Control-Allow-Origin")
// 		So(header, ShouldEqual, "*")
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})

// 	/*Convey("CORS allow origin is set to * when cors env is not set", t, func() {
// 		r, _ := http.NewRequest("OPTIONS", "/"+serviceVersion+"/client/service_check", nil)
// 		r.Header.Set("Origin", "localhost")
// 		response := executeRequest(r)
// 		header := response.HeaderMap.Get("Access-Control-Allow-Methods")
// 		So(header, ShouldEqual, "OPTIONS")
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})*/

// 	os.Setenv(global_environment.ENV_KEY_ORIGIN_ALLOWED, "")
// 	APP.Initialize()
// 	APP.initializeRoutes()

// 	Convey("CORS allow origin is set to * when cors env is not set", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/service_check", nil)
// 		response := executeRequest(req)
// 		header := response.HeaderMap.Get("Access-Control-Allow-Origin")
// 		So(header, ShouldEqual, "")
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})

// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs(temp))
// }

// func TestAssetBalance(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	temp := os.Getenv(global_environment.ENV_KEY_NODE_CONFIG)
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, "")

// 	Convey("Testing invalid query parameters", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/balances/test", nil)
// 		q := req.URL.Query()
// 		q.Add("asset_code", "NZD")
// 		req.URL.RawQuery = q.Encode()
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1100")
// 	})
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("./unit-test-data/nodeconfig/common/nodeconfig.toml"))
// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/balances/sudu", nil)
// 		q := req.URL.Query()
// 		q.Add("asset_code", "NZD")
// 		q.Add("asset_issuer", "aa")
// 		req.URL.RawQuery = q.Encode()
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(msg.Code, ShouldEqual, "API-1056")
// 	})

// 	Convey("Returns native asset balance  ", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/balances/sudu", nil)
// 		q := req.URL.Query()
// 		q.Add("asset_code", "xlm")
// 		req.URL.RawQuery = q.Encode()
// 		response := executeRequest(req)
// 		var msg model.AssetBalance
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		f, _ := strconv.ParseFloat(*msg.Balance, 64)
// 		So(f, ShouldBeGreaterThan, 0.0)
// 	})

// 	Convey("Returns asset balance for the requested asset ", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/balances/sudu", nil)
// 		q := req.URL.Query()
// 		q.Add("asset_code", "USD")
// 		q.Add("asset_issuer", "GBOQ4U4OUIRNM5WLEL2XTHRU2OXHFLNW5JF63HIFXGVASGED6L7G4FS3")
// 		req.URL.RawQuery = q.Encode()
// 		response := executeRequest(req)
// 		var msg model.AssetBalance
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		/*
// 		 If the asset exists with zero balance, system will return Asset & balance info. That's why we are checking for the issuer.
// 		*/
// 		So(msg.IssuerID, ShouldNotBeEmpty)
// 	})

// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs(temp))
// }

// func TestDOAssetBalances(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	temp := os.Getenv(global_environment.ENV_KEY_NODE_CONFIG)

// 	Convey("Testing invalid query parameters", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/balances/do", nil)
// 		q := req.URL.Query()
// 		q.Add("asset_code", "NZD")
// 		req.URL.RawQuery = q.Encode()
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		LOGGER.Info("Response code : %d", response.Code)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1100")
// 	})

// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/balances/do", nil)
// 		q := req.URL.Query()
// 		q.Add("asset_code", "FakeAssetDO")
// 		q.Add("asset_issuer", "FakeIssuer")
// 		req.URL.RawQuery = q.Encode()
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(msg.Code, ShouldEqual, "API-1056")
// 	})

// 	/*os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, "./unit-test-data/nodeconfig/asset-balance.toml")
// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/asset/do/balance", nil)
// 		q := req.URL.Query()
// 		q.Add("asset_code", "PLN")
// 		q.Add("issuer_id", "GAGDTJA5R2UURAYTYWRSN6PFIIIVOBFIPDZUQN7XLVRIMNOCBAXC4JYV")
// 		req.URL.RawQuery = q.Encode()
// 		response := executeRequest(req)
// 		var list []*model.AssetBalance
// 		json.Unmarshal(response.Body.Bytes(), &list)
// 		So(len(list), ShouldBeGreaterThan, 0)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})*/

// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs(temp))
// }

// /*
// 	This function returns all the assets trusted by the Issuing account.
// 	Asset is represented by the asset code and the account of the issuer of that asset.
// 	1. Issuing Account is loaded from the node config toml.
// 	2. Stellar is queried for Issuing Account node address for the assets it trusts
// 	4. Return the slice of assets
// */
// // removing these test cases as duplicate endpoint is removed
// /*func TestTrustedAssetsForIA(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	temp := os.Getenv(global_environment.ENV_KEY_NODE_CONFIG)
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, "")

// 	Convey("Error code is returned as expected ", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/issuingaccount/assets", nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(msg.Code, ShouldEqual, "API-1067")
// 	})
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, "./unit-test-data/nodeconfig/common/nodeconfig.toml")
// 	Convey("An array of Assets is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/issuingaccount/assets", nil)
// 		response := executeRequest(req)
// 		var assets []*model.Asset
// 		json.Unmarshal(response.Body.Bytes(), &assets)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		LOGGER.Infof("Length :%d", len(assets))
// 		So(len(assets), ShouldBeGreaterThan, 0)
// 	})
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, temp)
// }*/

// /*
// 	This function returns all the assets trusted by the Operating account.
// 	Asset is represented by the asset code and the account of the issuer of that asset.
// 	1. Operating Account is loaded from the node config toml.
// 	2. Stellar is queried for Operating Account node address for the assets it trusts
// 	4. Return the slice of assets
// */
// func TestTrustedAssetsForDA(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	temp := os.Getenv(global_environment.ENV_KEY_NODE_CONFIG)
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, "")

// 	Convey("Error code is returned as expected ", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/assets/sudu", nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(msg.Code, ShouldEqual, "API-1067")
// 	})
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("./unit-test-data/nodeconfig/common/nodeconfig.toml"))
// 	Convey("An array of Assets is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/assets/sudu", nil)
// 		response := executeRequest(req)
// 		var assets []*model.Asset
// 		json.Unmarshal(response.Body.Bytes(), &assets)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		LOGGER.Infof("Length :%d", len(assets))
// 		So(len(assets), ShouldBeGreaterThan, 0)
// 	})
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs(temp))
// }

// /*
// 	This function returns all the assets trusted by the IBM account.
// 	Asset is represented by the asset code and the account of the issuer of that asset.
// 	1. IBM Token Account is loaded from the node config toml.
// 	2. Stellar is queried for IBM Account node address for the assets it trusts
// 	4. Return the slice of assets
// */
// func TestWorldWireAssets(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	temp := os.Getenv(global_environment.ENV_KEY_NODE_CONFIG)
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, "")

// 	Convey("Error code is returned as expected ", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/assets", nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(msg.Code, ShouldEqual, "API-1068")
// 	})
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("./unit-test-data/nodeconfig/common/nodeconfig.toml"))
// 	Convey("An array of Assets is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/assets", nil)
// 		response := executeRequest(req)
// 		var assets []*model.Asset
// 		json.Unmarshal(response.Body.Bytes(), &assets)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		LOGGER.Infof("Length :%d", len(assets))
// 		So(len(assets), ShouldBeGreaterThan, 0)
// 	})
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs(temp))
// }

// /*
// 	In order to test this..Prerequisite for this to create a participant in Participant registry
// 	with an an issuing account that has one or more valid assets
// 	This tests if the endpoint is able to
// 		1. find the issuing account for a given domain.
// 		2. find all the assets trusted by this issuing account
// 		3. find all the assets (which were issued by this issued account) trusted by IBM account
// 		4. Combine 2 & 3 and return all these assets
// 	Explanation: All the assets trusted by a given account are appended to stellar
// 	network and this information is available in stellar. In addition,
// 	an Issuing account trusts all the assets issued by itself.
// */
// func TestGetAssetsForParticipant(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	/*Convey("No of elements should be at least 1", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/participant/au.one.payments.gftn.io/asset", nil)
// 		response := executeRequest(req)
// 		s11 := ""
// 		json.Unmarshal(response.Body.Bytes(), &s11)
// 		LOGGER.Infof(s11)
// 		var assets []*model.Asset
// 		json.Unmarshal(response.Body.Bytes(), &assets)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		LOGGER.Infof("Length :%d", len(assets))
// 		So(len(assets), ShouldBeGreaterThan, 0)
// 	})*/

// 	Convey("HTTP Status code is returned as expected", t, func() {
// 		httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
// 		defer httpmock.DeactivateAndReset()
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/participant/domain.info.json"))
// 		mockResponse := model.Participant{}
// 		json.Unmarshal(payload, &mockResponse)
// 		//simulate response
// 		responder, _ := httpmock.NewJsonResponder(http.StatusOK, mockResponse)
// 		httpmock.RegisterResponder("GET", "http://localhost:10080/v1/internal/pr/domain/au.one.payments.gftn.io", responder)

// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/assets/participant/au.one.payments.gftn.io", nil)
// 		q := req.URL.Query()
// 		q.Add("type", "issued")
// 		req.URL.RawQuery = q.Encode()
// 		response := executeRequest(req)
// 		var expectedPrs []*model.Asset
// 		json.Unmarshal(response.Body.Bytes(), &expectedPrs)
// 		LOGGER.Infof("%d", len(expectedPrs))
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})
// }

// /*
// 	This is a helper function to return the pointer for a given string literal
// */
// func Ptr(v string) *string {
// 	return &v
// }

// //TestGetParticipantRegistryByCountry - unit test for getting participants by passing country code
// func TestGetParticipantRegistryByCountry(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()

// 	Convey("HTTP Status code is returned as expected", t, func() {
// 		httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
// 		defer httpmock.DeactivateAndReset()
// 		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/participant/country-code.json"))
// 		mockResponse := []model.Participant{}
// 		json.Unmarshal(payload, &mockResponse)
// 		//simulate response
// 		responder, _ := httpmock.NewJsonResponder(http.StatusOK, mockResponse)
// 		httpmock.RegisterResponder("GET", "http://localhost:10080/v1/internal/pr/country/NZ", responder)

// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/participants?country_code=NZ", nil)
// 		response := executeRequest(req)
// 		expectedPrs := []model.Participant{}
// 		json.Unmarshal(response.Body.Bytes(), &expectedPrs)

// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		So(*expectedPrs[0].ID, ShouldEqual, *mockResponse[0].ID)
// 	})
// }

// /*
// 1: Feed an empty nodeconfig file URI and see if it returns Error code & HTTP status code properly
// 2: Feed a proper nodeconfig file URI with an existing issuing account and see if it returns the issuing account correctly
// 3. Feed a proper nodeconfig file URI without an issuing account and see if it creates an issuing account and
//    updates the node config file with the newly created issuing account. Mocking: This flow calls the participant registry to
//    update the newly created issuing account. That call is mocked to return HTTP OK Status code. The result can be
//    verified by checking for a newly created issuing account in (./unit-test-data/nodeconfig/no-IA.toml). Once you are
//    done checking feel free to empty out (./unit-test-data/nodeconfig/no-IA.toml) or just delete the issuing account
//    from it. Because if you rerun this test, you should be able to see if a new issuing account gets added there after
//    executing this unit test.
// */
// func TestCreateIssuingAccount(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	temp := os.Getenv(global_environment.ENV_KEY_NODE_CONFIG)

// 	//Feed an empty nodeconfig file URI and see if it returns Error code & HTTP status code properly
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, "")
// 	/*Convey("Error code & HTTP status code are returned as expected", t, func() {
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/onboarding/issuingaccount", nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(msg.Code, ShouldEqual, "API-1090")
// 	})*/

// 	//Feed a proper nodeconfig file URI with an existing issuing account and see if it returns the issuing account correctly
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("./unit-test-data/nodeconfig/IA.toml"))
// 	Convey("HTTP status code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/onboarding/issuingaccount", nil)
// 		response := executeRequest(req)
// 		actual := model.Account{}
// 		json.Unmarshal(response.Body.Bytes(), &actual)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		expected := "FakeIssuingAccount"
// 		So(expected, ShouldEqual, *actual.Address)
// 	})

// 	/*
// 			Feed a proper nodeconfig file URI without an issuing account and see if it creates an issuing account and
// 		   updates the node config file with the newly created issuing account. Mocking: This flow calls the participant registry to
// 		   update the newly created issuing account. That call is mocked to return HTTP OK Status code. The result can be
// 		   verified by checking for a newly created issuing account in (./unit-test-data/nodeconfig/no-IA.toml). Once you are
// 		   done checking feel free to empty out (./unit-test-data/nodeconfig/no-IA.toml) or just delete the issuing account
// 		   from it. Because if you rerun this test, you should be able to see if a new issuing account gets added there after
// 		   executing this unit test.
// 	*/
// 	/*
// 		 httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
// 		defer httpmock.DeactivateAndReset()
// 		var payload []byte
// 		mockResponse := ""
// 		json.Unmarshal(payload, &mockResponse)
// 		//simulate response
// 		responder, _ := httpmock.NewJsonResponder(http.StatusOK, mockResponse)
// 		httpmock.RegisterResponder("POST", "http://localhost:10080/v1/internal/pr/issuingaccount/nz.one.payments.gftn.io", responder)

// 		ResetContents("./unit-test-data/nodeconfig/without-IA.toml", "./unit-test-data/nodeconfig/no-IA.toml")
// 		os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, "./unit-test-data/nodeconfig/no-IA.toml")
// 		Convey("HTTP status code is returned as expected", t, func() {
// 			req, _ := http.NewRequest("POST", "/"+serviceVersion+"/onboarding/issuingaccount", nil)
// 			response := executeRequest(req)
// 			actual := model.IssuingAccount{}
// 			json.Unmarshal(response.Body.Bytes(), &actual)
// 			So(response.Code, ShouldEqual, http.StatusOK)
// 			So(*actual.StellarNetworkAddress, ShouldNotBeEmpty)
// 		}) */
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs(temp))
// }

// func ResetContents(sourceFilePath, destinationFilePath string) {
// 	from, _ := os.Open(common.Abs(sourceFilePath))
// 	defer from.Close()
// 	to, _ := os.Create(common.Abs(destinationFilePath))
// 	defer to.Close()
// 	io.Copy(to, from)
// }

// /*
// 1: Feed an empty nodeconfig file URI and see if it returns Error code & HTTP status code properly
// 2: Feed a proper nodeconfig file URI with an existing operating account and see if it returns the operating account correctly
// 3. Feed a proper nodeconfig file URI without an operating account and see if it creates an operating account and
//    updates the node config file with the newly created operating account. Mocking: This flow calls the participant registry to
//    update the newly created operating account. That call is mocked to return HTTP OK Status code. The result can be
//    verified by checking for a newly created operating account in (./unit-test-data/nodeconfig/no-DA.toml). Once you are
//    done checking feel free to delete the operating account from it. Because if you rerun this test, you should be
// 	able to see if a new operating account gets added there after executing this unit test.
// */
// func TestCreateOperatingAccount(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	temp := os.Getenv(global_environment.ENV_KEY_NODE_CONFIG)

// 	//Feed an empty nodeconfig file URI and see if it returns Error code & HTTP status code properly
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, "")
// 	Convey("Error code & HTTP status code are returned as expected", t, func() {
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/onboarding/operatingaccount/DA1/10000", nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 		So(msg.Code, ShouldEqual, "API-1095")
// 	})

// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("./unit-test-data/nodeconfig/DA.toml"))

// 	Convey("HTTP status code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/onboarding/operatingaccount/OperatingAcc1/10000", nil)
// 		response := executeRequest(req)
// 		actual := model.Account{}
// 		json.Unmarshal(response.Body.Bytes(), &actual)
// 		So(response.Code, ShouldEqual, http.StatusAlreadyReported)
// 		expected := "FakeOperatingAccount1NodeAddr"
// 		So(expected, ShouldEqual, *actual.Address)
// 	})
// 	ResetContents("./unit-test-data/nodeconfig/without-DA.toml", "./unit-test-data/nodeconfig/no-DA.toml")
// 	httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
// 	defer httpmock.DeactivateAndReset()
// 	var payload []byte
// 	mockResponse := ""
// 	json.Unmarshal(payload, &mockResponse)
// 	//simulate response
// 	responder, _ := httpmock.NewJsonResponder(http.StatusOK, mockResponse)
// 	httpmock.RegisterResponder("POST", "http://localhost:10080/v1/internal/pr/account/nz.one.payments.gftn.io", responder)

// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("./unit-test-data/nodeconfig/no-DA.toml"))
// 	/*Convey("HTTP status code is returned as expected", t, func() {
// 		// Note if you large amount like 10000 instead of 3, the transaction might
// 		// fail because it might exeed the native balance of the issuing account
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/onboarding/operatingaccount/OperatingAccX/3", nil)
// 		response := executeRequest(req)
// 		actual := model.OperatingAccount{}
// 		json.Unmarshal(response.Body.Bytes(), &actual)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		So(*actual.StellarNetworkAddress, ShouldNotBeEmpty)
// 	})*/
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs(temp))

// }

// /*
// 1: Pass a non existant operating account name and see if it returns Error code properly
// 2: Feed a proper nodeconfig file URI with an existing operating account and see if it returns the operating account correctly
// */
// func TestGetOperatingAccount(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	temp := os.Getenv(global_environment.ENV_KEY_NODE_CONFIG)

// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("./unit-test-data/nodeconfig/DA.toml"))

// 	Convey("HTTP status code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/onboarding/accounts/JunkData", nil)
// 		response := executeRequest(req)
// 		actual := model.Account{}
// 		json.Unmarshal(response.Body.Bytes(), &actual)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 	})

// 	Convey("HTTP status code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/onboarding/accounts/OperatingAcc1", nil)
// 		response := executeRequest(req)
// 		actual := model.Account{}
// 		json.Unmarshal(response.Body.Bytes(), &actual)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		expected := "FakeOperatingAccount1NodeAddr"
// 		So(expected, ShouldEqual, *actual.Address)
// 	})
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs(temp))
// }

// /*
// 	This tests if the response is coming back from compliance
// */
// func TestVerifyCompliance(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()

// 	Convey("Error code & HTTP status code are returned as expected", t, func() {
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/compliance", nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1122")
// 	})

// 	/*httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
// 	defer httpmock.DeactivateAndReset()
// 	mockResponse := model.CallbackComplianceCheckResponse{}
// 	mockResponse.ComplianceStatus = Ptr("OK")
// 	mockResponse.SanctionsStatus = Ptr("OK")
// 	mockResponse.ComplianceIdentification = Ptr("OK")
// 	mockResponse.AmlKycInfo = Ptr("OK")
// 	responder, _ := httpmock.NewJsonResponder(http.StatusOK, mockResponse)
// 	httpmock.RegisterResponder("POST", "http://localhost:11003/v1/compliance", responder)
// 	Convey("HTTP status code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/compliance", bytes.NewBuffer([]byte(`{"somefield":"somevalue"}`)))
// 		response := executeRequest(req)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})*/

// }

// /*
// 	This tests if the quote is coming back correctly
// */
// func TestGetQuote(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()

// 	Convey("Error code & HTTP status code are returned as expected", t, func() {
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/quote", bytes.NewBuffer([]byte("")))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1107")
// 	})

// 	httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
// 	defer httpmock.DeactivateAndReset()
// 	payloadResp, _ := ioutil.ReadFile(common.Abs("./unit-test-data/quote/quoteresponse1.json"))
// 	mockResponse := model.Quote{}
// 	json.Unmarshal(payloadResp, &mockResponse)
// 	//simulate response
// 	responder, _ := httpmock.NewJsonResponder(http.StatusOK, mockResponse)
// 	httpmock.RegisterResponder("POST", "http://localhost:8081/v1/callback/quote", responder)

// 	Convey("HTTP status code is returned as expected", t, func() {
// 		payloadRequest, _ := ioutil.ReadFile(common.Abs("./unit-test-data/quote/quoterequest1.json"))
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/quote", bytes.NewBuffer(payloadRequest))
// 		response := executeRequest(req)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		expected := model.Quote{}
// 		json.Unmarshal(response.Body.Bytes(), &expected)
// 		So(*expected.ExchangeRate, ShouldEqual, 1.456)
// 	})
// }

// /*
// 	This tests if the fee is coming back correctly
// */
// func TestCalculateFees(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()

// 	Convey("Error code & HTTP status code are returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/fees", bytes.NewBuffer([]byte("")))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1073")
// 	})

// 	/*
// 		httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
// 		defer httpmock.DeactivateAndReset()
// 		payloadResp, _ := ioutil.ReadFile("./unit-test-data/fee/response.json")
// 		mockResponse := model.Fee{}
// 		json.Unmarshal(payloadResp, &mockResponse)
// 		//simulate response
// 		//responder1, _ := httpmock.NewJsonResponder(http.StatusOK, mockResponse)
// 		responder2, _ := httpmock.NewJsonResponder(http.StatusOK, mockResponse)

// 		//httpmock.RegisterResponder("POST", "http://localhost:8081/v1/callback/fees/USD/10/yes", responder1)
// 		httpmock.RegisterResponder("POST", "http://localhost:10101/v1/internal/fees/fj.one.payments.gftn.io/FJD/1685.15/false", responder2)

// 		Convey("HTTP status code is returned as expected", t, func() {
// 			urlWithQueryParams := "/client/fees?source_currency=NZD&target_currency=FJD&beneficiary_amount=1685.15&price=1.2&beneficiary_domain=fj.one.payments.gftn.io"
// 			req, _ := http.NewRequest("GET", "/"+serviceVersion + urlWithQueryParams, nil)
// 			response := executeRequest(req)
// 			So(response.Code, ShouldEqual, http.StatusOK)
// 			expected := model.FeesAndAmountResponse{}
// 			json.Unmarshal(response.Body.Bytes(), &expected)
// 			So(*expected.AdministratorFee.FeeAmount, ShouldNotBeEmpty)
// 		}) */

// }

// func TestGetWhiteListParticipants(t *testing.T) {
// 	LOGGER.Infof("TestGetWhiteListParticipants")
// 	setEnvVariables()
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("unit-test-data/nodeconfig/whitelist.nodeconfig.toml"))

// 	APP.Initialize()
// 	APP.initializeRoutes()

// 	//	execute unit test
// 	Convey("Good case: GET", t, func() {
// 		url := "/" + serviceVersion + "/client/participants/whitelist"
// 		req, _ := http.NewRequest(http.MethodGet, url, nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})
// }

// func TestIsParticipantWhiteListed(t *testing.T) {
// 	LOGGER.Infof("TestIsParticipantWhiteListed")
// 	setEnvVariables()
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("unit-test-data/nodeconfig/whitelist.nodeconfig.toml"))
// 	APP.Initialize()
// 	APP.initializeRoutes()

// 	var domain = "testdomain.com"
// 	//	execute unit test
// 	Convey("Good case: GET", t, func() {
// 		url := "/" + serviceVersion + "/client/whitelist/" + domain
// 		req, _ := http.NewRequest(http.MethodGet, url, nil)

// 		response := executeRequest(req)
// 		var msg resp.SuccessMessage
// 		json.Unmarshal(response.Body.Bytes(), &msg)

// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		So(msg.Msg, ShouldEqual, "true")

// 	})
// 	domain = "nodomain.com"
// 	Convey("Bad case: GET", t, func() {
// 		url := "/" + serviceVersion + "/client/whitelist/" + domain
// 		req, _ := http.NewRequest(http.MethodGet, url, nil)

// 		response := executeRequest(req)
// 		var msg resp.SuccessMessage
// 		json.Unmarshal(response.Body.Bytes(), &msg)

// 		So(response.Code, ShouldEqual, http.StatusOK)
// 		So(msg.Msg, ShouldEqual, "false")
// 	})

// }

// func TestAddWhiteListParticipant(t *testing.T) {
// 	LOGGER.Infof("TestAddWhiteListParticipant")
// 	setEnvVariables()
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("unit-test-data/nodeconfig/whitelist.nodeconfig.toml"))

// 	APP.Initialize()
// 	APP.initializeRoutes()

// 	var domain = "newdomain.com"
// 	//	execute unit test
// 	Convey("Good case: POST", t, func() {
// 		url := "/" + serviceVersion + "/client/whitelist/" + domain
// 		req, _ := http.NewRequest(http.MethodPost, url, nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})
// 	Convey("Bad case: POST", t, func() {
// 		url := "/" + serviceVersion + "/client/whitelist/" + domain
// 		req, _ := http.NewRequest(http.MethodPost, url, nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotAcceptable)
// 	})
// }

// func TestRemoveWhiteListParticipant(t *testing.T) {
// 	LOGGER.Infof("TestRemoveWhiteListParticipant")
// 	setEnvVariables()
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("unit-test-data/nodeconfig/whitelist.nodeconfig.toml"))

// 	APP.Initialize()
// 	APP.initializeRoutes()

// 	var domain = "newdomain.com"

// 	//	execute unit test
// 	Convey("Good case: DELETE", t, func() {
// 		url := "/" + serviceVersion + "/client/whitelist/" + domain
// 		req, _ := http.NewRequest(http.MethodDelete, url, nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})

// 	Convey("Bad case: DELETE", t, func() {
// 		url := "/" + serviceVersion + "/client/whitelist/" + domain
// 		req, _ := http.NewRequest(http.MethodDelete, url, nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		LOGGER.Infof("TestRemoveWhiteListParticipant %v", response.Code)
// 		So(response.Code, ShouldEqual, http.StatusNotAcceptable)
// 	})

// }

// ///////
// func TestCreateOrAllowTrust(t *testing.T) {
// 	LOGGER.Infof("TestCreateOrAllowTrust")
// 	setEnvVariables()

// 	os.Setenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME, "sg.one.payments.gftn.io")

// 	APP.Initialize()
// 	APP.initializeRoutes()

// 	//-------------------------
// 	//* I. Good case: change trust : POST
// 	//	1. mock to  "http://localhost:10080/v1/internal/pr/domain/au.one.payments.gftn.io"
// 	httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
// 	defer httpmock.DeactivateAndReset()

// 	payload, _ := ioutil.ReadFile(common.Abs("unit-test-data/asset/au.participant_registry.json"))
// 	mockParticipant := model.Participant{}

// 	json.Unmarshal(payload, &mockParticipant)
// 	LOGGER.Infof("Mock URLAPIService  %v", *mockParticipant.URLAPIService)

// 	//simulate response
// 	responder, _ := httpmock.NewJsonResponder(http.StatusOK, mockParticipant)
// 	endPointURL := "http://localhost:10080/v1/internal/pr/domain/au.one.payments.gftn.io"

// 	LOGGER.Infof("Mock url %v", endPointURL)
// 	httpmock.RegisterResponder(http.MethodGet, endPointURL, responder)

// 	//	execute unit test
// 	LOGGER.Infof("I. TestCreateOrAllowTrust: Good case: change trust : POST")
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("unit-test-data/nodeconfig/sg.nodeconfig.toml"))
// 	jsonParams := "unit-test-data/asset/changetrust.json"

// 	Convey("I. Good case: change trust : POST", t, func() {
// 		payload, _ := ioutil.ReadFile(jsonParams)
// 		url := "/" + serviceVersion + "/client/trust/AUDDO?permission=request"
// 		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})

// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("unit-test-data/nodeconfig/au.nodeconfig.toml"))
// 	jsonParams = "unit-test-data/asset/changetrust_sgddo.json"
// 	os.Setenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME, "au.one.payments.gftn.io")

// 	httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
// 	defer httpmock.DeactivateAndReset()

// 	payload, _ = ioutil.ReadFile("unit-test-data/asset/sg.participant_registry.json")

// 	json.Unmarshal(payload, &mockParticipant)
// 	LOGGER.Infof("Mock URLAPIService  %v", *mockParticipant.URLAPIService)

// 	//simulate response
// 	responder, _ = httpmock.NewJsonResponder(http.StatusOK, mockParticipant)
// 	endPointURL = "http://localhost:10080/v1/internal/pr/domain/sg.one.payments.gftn.io"

// 	LOGGER.Infof("Mock url %v", endPointURL)
// 	httpmock.RegisterResponder(http.MethodGet, endPointURL, responder)

// 	Convey("II. Good case: change trust : POST", t, func() {
// 		payload, _ := ioutil.ReadFile(jsonParams)
// 		url := "/" + serviceVersion + "/client/trust/SGDDO?permission=request"
// 		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})
// 	//-------------------------
// 	//*	II. Bad case: change trust : POST:  domain is local domain
// 	//	1. mock to  "http://localhost:10080/v1/internal/pr/domain/au.one.payments.gftn.io"
// 	httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
// 	defer httpmock.DeactivateAndReset()
// 	os.Setenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME, "sg.one.payments.gftn.io")

// 	payload, _ = ioutil.ReadFile(common.Abs("unit-test-data/asset/sg.participant_registry.json"))
// 	mockParticipant = model.Participant{}

// 	json.Unmarshal(payload, &mockParticipant)
// 	LOGGER.Infof("Mock URLAPIService  %v", *mockParticipant.URLAPIService)

// 	//simulate response
// 	responder, _ = httpmock.NewJsonResponder(http.StatusOK, mockParticipant)
// 	endPointURL = "http://localhost:10080/v1/internal/pr/domain/sg.one.payments.gftn.io"

// 	LOGGER.Infof("Mock url %v", endPointURL)
// 	httpmock.RegisterResponder(http.MethodGet, endPointURL, responder)

// 	//	execute unit test
// 	jsonParams = common.Abs("unit-test-data/asset/changetrust_badcase1.json")
// 	LOGGER.Infof("II. Bad case: change trust : POST:  domain is local domain")
// 	Convey("II. Bad case: change trust : POST:  domain is local domain", t, func() {
// 		payload, _ := ioutil.ReadFile(jsonParams)
// 		url := "/" + serviceVersion + "/client/trust/AUDDO?permission=request"
// 		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 	})

// 	//-------------------------
// 	//* III. Bad case: change trust : POST : distAccount not exist in remote host
// 	//	1. mock to  "http://localhost:10080/v1/internal/pr/domain/au.one.payments.gftn.io"
// 	httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
// 	defer httpmock.DeactivateAndReset()

// 	payload, _ = ioutil.ReadFile(common.Abs("unit-test-data/asset/au.participant_registry.json"))
// 	mockParticipant = model.Participant{}

// 	json.Unmarshal(payload, &mockParticipant)
// 	LOGGER.Infof("Mock URLAPIService  %v", *mockParticipant.URLAPIService)

// 	//simulate response
// 	responder, _ = httpmock.NewJsonResponder(http.StatusOK, mockParticipant)
// 	endPointURL = "http://localhost:10080/v1/internal/pr/domain/au.one.payments.gftn.io"

// 	LOGGER.Infof("Mock url %v", endPointURL)
// 	httpmock.RegisterResponder(http.MethodGet, endPointURL, responder)

// 	//	execute unit test
// 	LOGGER.Infof("III. Bad case: change trust : POST : distAccount not exist in remote host")
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("unit-test-data/nodeconfig/sg.nodeconfig.toml"))
// 	jsonParams = common.Abs("unit-test-data/asset/changetrust_badcase2.json")

// 	Convey("III. Bad case: change trust : POST : distAccount not exist in remote host", t, func() {
// 		payload, _ := ioutil.ReadFile(jsonParams)
// 		url := "/" + serviceVersion + "/client/asset/AUDDO/trust?action=change"
// 		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 	})

// 	//-------------------------
// 	//* IV. Bad case: change trust : POST : AssetCode is not DA nor DO
// 	//	1. mock to  "http://localhost:10080/v1/internal/pr/domain/au.one.payments.gftn.io"
// 	httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
// 	defer httpmock.DeactivateAndReset()

// 	payload, _ = ioutil.ReadFile(common.Abs("unit-test-data/asset/au.participant_registry.json"))
// 	mockParticipant = model.Participant{}

// 	json.Unmarshal(payload, &mockParticipant)
// 	LOGGER.Infof("Mock URLAPIService  %v", *mockParticipant.URLAPIService)

// 	//simulate response
// 	responder, _ = httpmock.NewJsonResponder(http.StatusOK, mockParticipant)
// 	endPointURL = "http://localhost:10080/v1/internal/pr/domain/au.one.payments.gftn.io"

// 	LOGGER.Infof("Mock url %v", endPointURL)
// 	httpmock.RegisterResponder(http.MethodGet, endPointURL, responder)

// 	//	execute unit test
// 	LOGGER.Infof("IV. Bad case: change trust : POST : AssetCode is not DA nor DO")
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("unit-test-data/nodeconfig/sg.nodeconfig.toml"))
// 	jsonParams = common.Abs("unit-test-data/asset/changetrust_badcase2.json")

// 	Convey("IV. Bad case: change trust : POST : AssetCode is not DA nor DO", t, func() {
// 		payload, _ := ioutil.ReadFile(jsonParams)
// 		url := "/" + serviceVersion + "/client/asset/AUD/trust?action=change"
// 		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 	})

// 	//-------------------------
// 	//* V. Bad case: change trust : POST : Cannot do internal trust if asset type is DO
// 	//	1. mock to  "http://localhost:10080/v1/internal/pr/domain/au.one.payments.gftn.io"
// 	httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
// 	defer httpmock.DeactivateAndReset()

// 	payload, _ = ioutil.ReadFile(common.Abs("unit-test-data/asset/au.participant_registry.json"))
// 	mockParticipant = model.Participant{}

// 	json.Unmarshal(payload, &mockParticipant)
// 	LOGGER.Infof("Mock URLAPIService  %v", *mockParticipant.URLAPIService)

// 	//simulate response
// 	responder, _ = httpmock.NewJsonResponder(http.StatusOK, mockParticipant)
// 	endPointURL = "http://localhost:10080/v1/internal/pr/domain/au.one.payments.gftn.io"

// 	LOGGER.Infof("Mock url %v", endPointURL)
// 	httpmock.RegisterResponder(http.MethodGet, endPointURL, responder)

// 	//	execute unit test
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("unit-test-data/nodeconfig/sg.nodeconfig.toml"))
// 	jsonParams = common.Abs("unit-test-data/asset/changetrust_badcase2.json")

// 	LOGGER.Infof("V. Bad case: change trust : POST : Cannot do internal trust if asset type is DO")
// 	Convey("V. Bad case: change trust : POST : Cannot do internal trust if asset type is DO", t, func() {
// 		payload, _ := ioutil.ReadFile(jsonParams)
// 		url := "/" + serviceVersion + "/client/asset/AUDDO/trust?action=change"
// 		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 	})

// 	//-------------------------
// 	//*	VI.Good case: allow trust : POST : authorize=true
// 	//	1. mock to  "http://localhost:10080/v1/internal/pr/domain/au.one.payments.gftn.io"
// 	httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
// 	defer httpmock.DeactivateAndReset()

// 	payload, _ = ioutil.ReadFile(common.Abs("unit-test-data/asset/au.participant_registry.json"))
// 	mockParticipant = model.Participant{}

// 	json.Unmarshal(payload, &mockParticipant)
// 	LOGGER.Infof("Mock URLAPIService  %v", *mockParticipant.URLAPIService)

// 	//simulate response
// 	responder, _ = httpmock.NewJsonResponder(http.StatusOK, mockParticipant)
// 	endPointURL = "http://localhost:10080/v1/internal/pr/domain/au.one.payments.gftn.io"

// 	LOGGER.Infof("Mock url %v", endPointURL)
// 	httpmock.RegisterResponder(http.MethodGet, endPointURL, responder)

// 	//	execute unit test
// 	LOGGER.Infof("VI. Good case: allow trust : POST : authorize=true")
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("unit-test-data/nodeconfig/sg.nodeconfig.toml"))
// 	jsonParams = common.Abs("unit-test-data/asset/allowtrust.json")

// 	Convey("VI. Good case: allow trust : POST : authorize=true ", t, func() {
// 		payload, _ := ioutil.ReadFile(jsonParams)
// 		url := "/" + serviceVersion + "/client/trust/SGDDO?permission=allow"
// 		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})

// 	//-------------------------
// 	//*	VII. Good case: allow trust : POST : authorize=false
// 	//	execute unit test
// 	LOGGER.Infof("VII. Good case: allow trust : POST authorize=false")
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("unit-test-data/nodeconfig/sg.nodeconfig.toml"))
// 	jsonParams = common.Abs("unit-test-data/asset/allowtrust.json")

// 	Convey("VII. Good case: allow trust : POST authorize=false ", t, func() {
// 		payload, _ := ioutil.ReadFile(jsonParams)
// 		url := "/" + serviceVersion + "/client/trust/SGDDO?permission=allow"
// 		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})

// 	//-------------------------
// 	//*	VIII. Bad case: allow trust : POST: distAccount not exist
// 	//	1. mock to  "http://localhost:10080/v1/internal/pr/domain/au.one.payments.gftn.io"
// 	httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
// 	defer httpmock.DeactivateAndReset()

// 	payload, _ = ioutil.ReadFile(common.Abs("unit-test-data/asset/au.participant_registry.json"))
// 	mockParticipant = model.Participant{}

// 	json.Unmarshal(payload, &mockParticipant)
// 	LOGGER.Infof("Mock URLAPIService  %v", *mockParticipant.URLAPIService)

// 	//simulate response
// 	responder, _ = httpmock.NewJsonResponder(http.StatusOK, mockParticipant)
// 	endPointURL = "http://localhost:10080/v1/internal/pr/domain/au.one.payments.gftn.io"

// 	LOGGER.Infof("Mock url %v", endPointURL)
// 	httpmock.RegisterResponder(http.MethodGet, endPointURL, responder)

// 	//	execute unit test
// 	LOGGER.Infof("VIII. Bad case: allow trust : POST: distAccount not exist")
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("unit-test-data/nodeconfig/sg.nodeconfig.toml"))
// 	jsonParams = common.Abs("unit-test-data/asset/allowtrust_badcase1.json")

// 	Convey("VIII. Bad case: allow trust : POST: distAccount not exist", t, func() {
// 		payload, _ := ioutil.ReadFile(jsonParams)
// 		url := "/" + serviceVersion + "/client/trust/SGDDO?permission=allow"
// 		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 	})

// 	//-------------------------
// 	//*	IX. Bad case: allow trust : POST: domain name is local
// 	//	1. mock to  "http://localhost:10080/v1/internal/pr/domain/sg.one.payments.gftn.io"
// 	httpmock.ActivateNonDefault(resty.DefaultClient.GetClient())
// 	defer httpmock.DeactivateAndReset()

// 	payload, _ = ioutil.ReadFile(common.Abs("unit-test-data/asset/sg.participant_registry.json"))
// 	mockParticipant = model.Participant{}

// 	json.Unmarshal(payload, &mockParticipant)
// 	LOGGER.Infof("Mock URLAPIService  %v", *mockParticipant.URLAPIService)

// 	//simulate response
// 	responder, _ = httpmock.NewJsonResponder(http.StatusOK, mockParticipant)
// 	endPointURL = "http://localhost:10080/v1/internal/pr/domain/sg.one.payments.gftn.io"

// 	LOGGER.Infof("Mock url %v", endPointURL)
// 	httpmock.RegisterResponder(http.MethodGet, endPointURL, responder)

// 	//	execute unit test
// 	LOGGER.Infof("IX. Bad case: allow trust : POST: domain name is local")
// 	os.Setenv(global_environment.ENV_KEY_NODE_CONFIG, common.Abs("unit-test-data/nodeconfig/sg.nodeconfig.toml"))
// 	jsonParams = common.Abs("unit-test-data/asset/allowtrust_badcase2.json")

// 	Convey("IX. Bad case: allow trust : POST: domain name is local", t, func() {
// 		payload, _ := ioutil.ReadFile(jsonParams)
// 		url := "/" + serviceVersion + "/client/trust/SGDDO?permission=allow"
// 		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotFound)
// 	})
// }

// func TestGetMyFee(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/fees/FakeAsset", nil)
// 		q := req.URL.Query()
// 		q.Add("amount", "11s")
// 		q.Add("is_fee_included", "yes")
// 		req.URL.RawQuery = q.Encode()
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1059")
// 	})
// }

// func TestClearTransaction(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/client/transactions/send", bytes.NewBuffer([]byte("")))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1107")
// 	})
// }

// func TestProcessTxnDetails(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/client/transactions", nil)
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1100")
// 	})
// }

// func TestGetQuotes(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/client/quotes", bytes.NewBuffer([]byte("")))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1107")
// 	})
// }

// func TestVerifyExchange(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/verifications/exchange", bytes.NewBuffer([]byte("")))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusNotAcceptable)
// 		So(msg.Code, ShouldEqual, "API-1042")
// 	})
// }

// func TestExchange(t *testing.T) {
// 	setEnvVariables()
// 	APP.Initialize()
// 	APP.initializeRoutes()
// 	Convey("Error code is returned as expected", t, func() {
// 		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/client/transactions/exchange", bytes.NewBuffer([]byte("")))
// 		response := executeRequest(req)
// 		msg := model.WorldWireError{}
// 		json.Unmarshal(response.Body.Bytes(), &msg)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 		So(msg.Code, ShouldEqual, "API-1042")
// 	})
// }
