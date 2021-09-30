// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	authtesting "github.com/GFTN/gftn-services/utility/testing"
)

// var adminUrl = "http://localhost:11111/v1"

/*
	This is a helper function to return the pointer for a given string literal
*/
// func Ptr(v string) *string {
// 	return &v
// }

/*
Custom function to execute the unit tests
*/
// func executeRequest(req *http.Request) *httptest.ResponseRecorder {
// 	rr := httptest.NewRecorder()
// 	a.Router.ServeHTTP(rr, req)

// 	return rr
// }

// /*
// 	Checks the response code
// */

// func checkResponseCode(t *testing.T, expected, actual int) {
// 	if expected != actual {
// 		t.Errorf("Expected response code %d. Got %d", expected, actual)
// 	}
// }

// func setEnvVariables() {
// 	os.Setenv(global_environment.ENV_KEY_SERVICE_LOG_FILE, common.Abs("log.txt"))
// 	os.Setenv(global_environment.ENV_KEY_SERVICE_VERSION, "v1")
// 	os.Setenv(global_environment.ENV_KEY_SERVICE_PORT, "11111")
// 	os.Setenv(global_environment.ENV_KEY_SERVICE_ERROR_CODES_FILE, common.Abs("error-codes/adminErrorCodes.toml"))

// 	os.Setenv(environment.ENV_KEY_ADMIN_DB_NAME, "GFTN")
// 	os.Setenv(environment.ENV_KEY_BLOCKLIST_DB_TABLE, "blocklist")
// 	os.Setenv(environment.ENV_KEY_PAYMENTS_DB_TABLE, "payments")

// 	os.Setenv(environment.ENV_KEY_DB_USER, "wwUserAdmin")
// 	os.Setenv(environment.ENV_KEY_DB_PWD, "wwgftn!*@")
// }

func TestAuthForExternalEndpoint(t *testing.T) {
	a := App{}
	a.initRoutes()
	Convey("Testing authorization for external endpoints...", t, func() {
		authtesting.InitAuthTesting()
		err := a.Router.Walk(authtesting.AuthWalker)
		So(err, ShouldBeNil)
		err = a.AdminRouter.Walk(authtesting.AuthWalker)
		So(err, ShouldBeNil)
	})
}

// func TestCreateNativePayment(t *testing.T) {
// 	setEnvVariables()

// 	Convey("Payment Should be successful", t, func() {
// 		sourceAccount := nodeconfig.Account{}
// 		sourceAccount.NodeAddress = "GCSB6625EY6SYBJR6NJZXNJAWO4TZGMZTYT7H644OTHXVW25TPNPSM2Q"
// 		sourceAccount.NodeSeed = "(seed value)"
// 		destinationAddress := "GD5D6PKAYCKQPYDTOB2OFIKEFMJZFF6G3Q3WE4B4OSZSTBNVV5LCXY5I"
// 		//destinationSeed := "(seed value)"
// 		balStrBefore, _ := asset.GetStellarAccount(destinationAddress).GetNativeBalance()
// 		amount := "0.0001"
// 		result, err := asset.CreateNativePayment(sourceAccount, destinationAddress, amount)
// 		So(err == nil, ShouldEqual, true)
// 		balStr, _ := asset.GetStellarAccount(destinationAddress).GetNativeBalance()
// 		fmt.Println("Balance before : ", balStrBefore)
// 		fmt.Println("Amount paid    : ", amount)
// 		fmt.Println("Balance After  : ", balStr)
// 		So(result, ShouldNotBeEmpty)
// 	})
// }

// TestServiceCheck - Test service check endpoint exists
// func TestServiceCheck(t *testing.T) {
// 	setEnvVariables()
// 	a.InitApp()
// 	req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/service_check", nil)
// 	response := executeRequest(req)
// 	checkResponseCode(t, http.StatusOK, response.Code)
// }

// func TestAddBlocklistClient(t *testing.T) {
// 	client := blocklist_client.Client{
// 		HTTPClient: &http.Client{Timeout: time.Second * 10},
// 		AdminUrl:   adminUrl,
// 	}
// 	_, err := client.AddBlocklist("{\"type\":\"coUntrY\",\"value\":[\"Sa\" , \"us\"]}")
// 	Convey("Successful get caller identity", t, func() {
// 		So(err, ShouldBeNil)
// 	})
// }

// func TestGetBlocklistClient(t *testing.T) {
// 	client := blocklist_client.Client{
// 		HTTPClient: &http.Client{Timeout: time.Second * 10},
// 		AdminUrl:   adminUrl,
// 	}
// 	res, err := client.GetBlocklist("CountrY")
// 	fmt.Printf("%+v", res)

// 	Convey("Successful get caller identity", t, func() {
// 		So(err, ShouldBeNil)
// 		So(len(res[0].Value), ShouldEqual, 2)
// 	})
// }

// func TestValidateFromBlocklistClient(t *testing.T) {
// 	client := blocklist_client.Client{
// 		HTTPClient: &http.Client{Timeout: time.Second * 10},
// 		AdminUrl:   adminUrl,
// 	}
// 	res, err := client.ValidateFromBlocklist("[{\"type\":\"coUntrY\",\"value\":[\"sa\"]},{\"type\":\"currENcy\",\"value\":[\"usd\"]}]")
// 	fmt.Println(res)
// 	if err != nil {
// 		fmt.Printf("%s", err)
// 	}

// 	Convey("SA should be in the list, therefore denied expected", t, func() {
// 		So(err, ShouldBeNil)
// 		So(res, ShouldEqual, "denied")
// 	})
// }

// func TestDeleteBlocklistClient(t *testing.T) {
// 	client := blocklist_client.Client{
// 		HTTPClient: &http.Client{Timeout: time.Second * 10},
// 		AdminUrl:   adminUrl,
// 	}
// 	res, err := client.RemoveBlocklist("{\"type\":\"country\",\"value\":[\"sa\"]}")
// 	fmt.Println(res)
// 	Convey("Successful get caller identity", t, func() {
// 		So(err, ShouldBeNil)
// 	})
// }

// func TestRevalidateFromBlocklistClient(t *testing.T) {
// 	client := blocklist_client.Client{
// 		HTTPClient: &http.Client{Timeout: time.Second * 10},
// 		AdminUrl:   adminUrl,
// 	}
// 	res, err := client.ValidateFromBlocklist("[{\"type\":\"coUntrY\",\"value\":[\"sa\"]},{\"type\":\"currENcy\",\"value\":[\"usd\"]}]")
// 	fmt.Println(res)
// 	Convey("SA should not be in the list anymore, therefore approved expected", t, func() {
// 		So(err, ShouldBeNil)
// 		So(res, ShouldEqual, "approved")
// 	})
// }

// func TestDeleteAllBlocklistClient(t *testing.T) {
// 	client := blocklist_client.Client{
// 		HTTPClient: &http.Client{Timeout: time.Second * 10},
// 		AdminUrl:   adminUrl,
// 	}
// 	res, err := client.RemoveBlocklist("{\"type\":\"country\",\"value\":[\"us\"]}")
// 	fmt.Println(res)
// 	Convey("Successful get caller identity", t, func() {
// 		So(err, ShouldBeNil)
// 	})
// }

/*
// tests for post endpoint
func TestPostBlocklist(t *testing.T) {
	//Test posting invalid blocklist record to receive failure
	Convey("Error code is returned as expected for post", t, func() {
		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/invalid-blocklist-payload.json"))
		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/blocklist", bytes.NewBuffer(payload))
		response := executeRequest(req)
		msg := model.WorldWireError{}
		json.Unmarshal(response.Body.Bytes(), &msg)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
		So(msg.Code, ShouldEqual, "ADMIN-0014")
	})

	//Test posting valid record to receive success
	Convey("Successful creation of blocklist record", t, func() {
		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-blocklist-payload-country.json"))
		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/blocklist", bytes.NewBuffer(payload))
		response := executeRequest(req)
		var resp blocklistResponse
		json.Unmarshal(response.Body.Bytes(), &resp)
		So(resp.Status, ShouldEqual, "Success")
		So(response.Code, ShouldEqual, http.StatusOK)
	})

	//Test posting same record to receive error
	Convey("Fail creating blocklist record", t, func() {
		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-blocklist-payload-country.json"))
		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/blocklist", bytes.NewBuffer(payload))
		response := executeRequest(req)
		msg := model.WorldWireError{}
		json.Unmarshal(response.Body.Bytes(), &msg)
		So(response.Code, ShouldEqual, http.StatusNotFound)
		So(msg.Code, ShouldEqual, "ADMIN-0017")
	})

	//Test posting second valid blocklist record to receive success
	Convey("Successful creation of blocklist record", t, func() {
		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-blocklist-payload-country2.json"))
		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/blocklist", bytes.NewBuffer(payload))
		response := executeRequest(req)
		var resp blocklistResponse
		json.Unmarshal(response.Body.Bytes(), &resp)
		So(resp.Status, ShouldEqual, "Success")
		So(response.Code, ShouldEqual, http.StatusOK)

	})

	//Test posting currency blocklist record to receive success
	Convey("Successful creation of currency record", t, func() {
		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-blocklist-payload-currency.json"))
		req, _ := http.NewRequest("POST", "/"+serviceVersion+"/internal/blocklist", bytes.NewBuffer(payload))
		response := executeRequest(req)
		var resp blocklistResponse
		json.Unmarshal(response.Body.Bytes(), &resp)
		So(resp.Status, ShouldEqual, "Success")
		So(response.Code, ShouldEqual, http.StatusOK)

	})
}

// tests for get endpoint
func TestGetBlocklist(t *testing.T) {
	//Test invalid type(should be institution)
	Convey("Error code is returned as expected for invalid query type", t, func() {
		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/blocklist?type=instItutioNs", nil)
		response := executeRequest(req)
		msg := model.WorldWireError{}
		json.Unmarshal(response.Body.Bytes(), &msg)
		So(response.Code, ShouldEqual, http.StatusNotFound)
		So(msg.Code, ShouldEqual, "ADMIN-0016")
	})

	//Test country type & check if there are 3 results
	Convey("Successful getting blocklist record from service", t, func() {
		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/blocklist?type=CoUntrY", nil)
		response := executeRequest(req)
		resp := []model.Blocklist{}
		json.Unmarshal(response.Body.Bytes(), &resp)
		So(response.Code, ShouldEqual, http.StatusOK)
		So(*resp[0].Type, ShouldEqual, "COUNTRY")
		So(len(resp[0].Value), ShouldEqual, 3)

	})

	//Test get all payout point with no spcified criteria
	Convey("Successful all blocklist records from the service", t, func() {
		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/blocklist", nil)
		response := executeRequest(req)
		resp := []model.Blocklist{}
		json.Unmarshal(response.Body.Bytes(), &resp)
		So(response.Code, ShouldEqual, http.StatusOK)
		So(len(resp), ShouldEqual, 2)

	})

}

// tests for validate endpoint
func TestValidateBlocklist(t *testing.T) {
	//Test country that is not on the blocklist
	Convey("Approved country SG", t, func() {
		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/blocklist/validate?type=countrY&value=sg", nil)
		response := executeRequest(req)
		var resp blocklistResponse
		json.Unmarshal(response.Body.Bytes(), &resp)
		So(resp.Status, ShouldEqual, "approved")
		So(response.Code, ShouldEqual, http.StatusOK)
	})

	//Test country that is on the blocklist
	Convey("Denied country TW", t, func() {
		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/blocklist/validate?type=CouNtry&value=tw", nil)
		response := executeRequest(req)
		var resp blocklistResponse
		json.Unmarshal(response.Body.Bytes(), &resp)
		So(resp.Status, ShouldEqual, "denied")
		So(response.Code, ShouldEqual, http.StatusOK)

	})

}

// tests for delete endpoint
func TestDeleteBlocklist(t *testing.T) {
	//find country that does not exists
	Convey("Try deleting country that does not exists, should failed", t, func() {
		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/invalid-delete-blocklist-payload.json"))
		req, _ := http.NewRequest("DELETE", "/"+serviceVersion+"/internal/blocklist", bytes.NewBuffer(payload))
		response := executeRequest(req)
		So(response.Code, ShouldEqual, http.StatusNotFound)

		msg := model.WorldWireError{}
		json.Unmarshal(response.Body.Bytes(), &msg)
		So(msg.Code, ShouldEqual, "ADMIN-0016")
	})

	// delete country records
	Convey("Successful delete country record", t, func() {

		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-delete-blocklist-payload.json"))
		req, _ := http.NewRequest("DELETE", "/"+serviceVersion+"/internal/blocklist", bytes.NewBuffer(payload))
		response := executeRequest(req)
		So(response.Code, ShouldEqual, http.StatusOK)
	})

	// delete currency records
	Convey("Successful delete currency record", t, func() {

		payload, _ := ioutil.ReadFile(common.Abs("./unit-test-data/valid-blocklist-payload-currency.json"))
		req, _ := http.NewRequest("DELETE", "/"+serviceVersion+"/internal/blocklist", bytes.NewBuffer(payload))
		response := executeRequest(req)
		So(response.Code, ShouldEqual, http.StatusOK)
	})

	//get country type and see if all records have been deleted
	Convey("Successful deleting all records on the blocklist", t, func() {
		req, _ := http.NewRequest("GET", "/"+serviceVersion+"/internal/blocklist?type=country", nil)
		response := executeRequest(req)
		resp := []model.Blocklist{}
		json.Unmarshal(response.Body.Bytes(), &resp)
		So(response.Code, ShouldEqual, http.StatusOK)
		So(*resp[0].Type, ShouldEqual, "COUNTRY")
		So(len(resp[0].Value), ShouldEqual, 0)

	})

}

*/
