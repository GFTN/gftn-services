// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package nqsdbclient

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/jmoiron/sqlx/types"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/GFTN/gftn-services/quotes-service/utility/nqsmodel"
)

//config for testing purpose
const (
	host       = "localhost"
	port       = 5432
	dbuser     = "user"
	dbpassword = "password"
	dbname     = "user"
)

func defaultQuote() QuoteDB {
	requestID := uuid.Must(uuid.NewV4()).String()
	rfiID := "ie.one.payments.worldwire.io"
	ofiID := "hk.one.payments.worldwire.io"
	quoteID := uuid.Must(uuid.NewV4()).String() + rfiID
	limitMaxOfi := decimal.NewFromFloat(1000.0)
	limitMinOfi := decimal.NewFromFloat(100.0)
	sourceAsset := types.JSONText(`{
		"asset_code": "GBPDO",
		"asset_issuer": "GBPD3QEYU5E6ORGOY54WS2BQIFSJBM3GUQQOEPPYB4DZ2B5GQ7QHUVAV",
		"asset_type": "DO"
	}`)
	targetAsset := types.JSONText(`{
		"asset_code": "HKDDO",
		"asset_issuer": "GA3Z5DS6GAPBI6EGRFRCEJKAFBZHESW62U5ME3OJYVD5VEREY5ENTGIK",
		"asset_type": "DO"
	}`)
	timeRequest := time.Now().Unix()
	statusQuote := 1
	timeExpireOfi := int64(1796876024000)
	timeExpireRfi := time.Now().Unix() + 1000000
	timeStartRfi := int64(0)
	quote := QuoteDB{
		RequestID:     &requestID,
		QuoteID:       &quoteID,
		RfiId:         &rfiID,
		OfiId:         &ofiID,
		LimitMaxOfi:   &limitMaxOfi,
		LimitMinOfi:   &limitMinOfi,
		SourceAsset:   &sourceAsset,
		TargetAsset:   &targetAsset,
		TimeRequest:   &timeRequest,
		StatusQuote:   &statusQuote,
		TimeExpireOfi: &timeExpireOfi,
		TimeExpireRfi: &timeExpireRfi,
		TimeStartRfi:  &timeStartRfi,
	}
	return quote
}

func TestCreateConnection(t *testing.T) {
	pdg := PostgreDatabaseClient{
		Host:     host,
		Port:     port,
		User:     dbuser,
		Password: dbpassword,
		Dbname:   dbname,
	}
	err := pdg.CreateConnection()
	if err != nil {
		LOGGER.Error(err)
	}
	pdg.CloseConnection()

}

func TestCreateQuote(t *testing.T) {
	pdg := PostgreDatabaseClient{
		Host:     host,
		Port:     port,
		User:     dbuser,
		Password: dbpassword,
		Dbname:   dbname,
	}
	err := pdg.CreateConnection()
	if err != nil {
		LOGGER.Error(err)
		t.FailNow()
	}
	defer pdg.CloseConnection()
	Convey("CreateQuote shall create a new quote and gettable by GetQuotes ", t, func(c C) {
		payloadRequest, _ := ioutil.ReadFile("../../unit-test-data/quote/quoterequest1.json")
		assetPriceQuoteRequest := nqsmodel.NqsQuoteRequest{}
		err = json.Unmarshal(payloadRequest, &assetPriceQuoteRequest)
		if err != nil {
			LOGGER.Debug(err)
		}
		requestID := uuid.Must(uuid.NewV4()).String()
		LOGGER.Info("UUID generated: ", requestID)
		rfiDomain := "testRFI"
		quoteID := requestID + rfiDomain
		ofiDomain := *assetPriceQuoteRequest.OfiId
		maxLimit := *assetPriceQuoteRequest.LimitMaxOfi
		minLimit := *assetPriceQuoteRequest.LimitMinOfi
		sourceAsset := *assetPriceQuoteRequest.SourceAsset
		targetAsset := *assetPriceQuoteRequest.TargetAsset
		sourceAssetJson, _ := json.Marshal(sourceAsset)
		targetAssetJson, _ := json.Marshal(targetAsset)
		timeOfRequest := time.Now().Unix()
		timeExpireOfi := int64(1796876024000)
		quoteStatus := 1
		err := pdg.CreateQuote(requestID, quoteID, rfiDomain, ofiDomain, maxLimit, minLimit, sourceAssetJson, targetAssetJson, timeOfRequest, quoteStatus, timeExpireOfi)
		So(err, ShouldEqual, nil)
	})
}

func TestCreateGetQuote(t *testing.T) {
	pdg := PostgreDatabaseClient{
		Host:     host,
		Port:     port,
		User:     dbuser,
		Password: dbpassword,
		Dbname:   dbname,
	}
	err := pdg.CreateConnection()
	if err != nil {
		LOGGER.Error(err)
		t.FailNow()
	}
	defer pdg.CloseConnection()
	Convey("CreateQuote shall create a new quote and gettable by GetQuotes ", t, func(c C) {
		// payloadRequest, _ := ioutil.ReadFile("../../unit-test-data/quote/quoterequest1.json")
		payloadRequest := `{
			"time_expire_ofi": 1796876024000,
			"limit_max_ofi": 1000,
			"limit_min_ofi": 100,
			"source_asset": {"asset_code":"GBPDO",
			  "asset_issuer": "GBPD3QEYU5E6ORGOY54WS2BQIFSJBM3GUQQOEPPYB4DZ2B5GQ7QHUVAV",
			  "asset_type": "DO"
			},
			"target_asset": {"asset_code":"HKDDO",
			  "asset_issuer": "GA3Z5DS6GAPBI6EGRFRCEJKAFBZHESW62U5ME3OJYVD5VEREY5ENTGIK",
			  "asset_type": "DO"
			},
			"ofi_id": "testOFI"
		  }`
		assetPriceQuoteRequest := nqsmodel.NqsQuoteRequest{}
		err = json.Unmarshal([]byte(payloadRequest), &assetPriceQuoteRequest)
		if err != nil {
			LOGGER.Debug(err)
		}
		requestID := uuid.Must(uuid.NewV4()).String()
		LOGGER.Info("UUID generated: ", requestID)
		rfiDomain := "testRFI"
		quoteID := requestID + rfiDomain
		ofiDomain := *assetPriceQuoteRequest.OfiId
		maxLimit := *assetPriceQuoteRequest.LimitMaxOfi
		minLimit := *assetPriceQuoteRequest.LimitMinOfi
		sourceAsset := *assetPriceQuoteRequest.SourceAsset
		targetAsset := *assetPriceQuoteRequest.TargetAsset
		sourceAssetJson, _ := json.Marshal(sourceAsset)
		targetAssetJson, _ := json.Marshal(targetAsset)
		timeOfRequest := time.Now().Unix()
		timeExpireOfi := int64(1796876024000)

		quoteStatus := 1
		err := pdg.CreateQuote(requestID, quoteID, rfiDomain, ofiDomain, maxLimit, minLimit, sourceAssetJson, targetAssetJson, timeOfRequest, quoteStatus, timeExpireOfi)
		if err != nil {
			LOGGER.Error(err)
			t.FailNow()
		}
		quotes, err := pdg.GetQuotes(requestID, ofiDomain)
		if err != nil {
			LOGGER.Error(err)
			t.FailNow()
		}
		So(len(quotes), ShouldEqual, 1)
		So(*quotes[0].RequestID, ShouldEqual, requestID)
	})
}

// Test UpdateQuote, ExecuteQuote, CancelQuote
func TestChangeStatusQuote(t *testing.T) {
	pdg := PostgreDatabaseClient{
		Host:     host,
		Port:     port,
		User:     dbuser,
		Password: dbpassword,
		Dbname:   dbname,
	}
	err := pdg.CreateConnection()
	if err != nil {
		LOGGER.Error(err)
		t.Fail()
	}
	defer pdg.CloseConnection()
	// create testing data
	quoteInsert1 := defaultQuote()
	testQuoteResponse := types.JSONText([]byte(`{"test":"test"}`))
	testQuoteResponseBase64 := "test"
	testQuoteSig := "testSig"
	testRfiID1 := "rfi1"
	limitMaxRfi := decimal.NewFromFloat(10)
	limitMinRfi := decimal.NewFromFloat(1)
	quoteInsert1.QuoteResponse = &testQuoteResponse
	quoteInsert1.QuoteResponseBase64 = &testQuoteResponseBase64
	quoteInsert1.QuoteResponseSignature = &testQuoteSig
	quoteInsert1.RfiId = &testRfiID1
	quoteInsert1.LimitMaxRfi = &limitMaxRfi
	quoteInsert1.LimitMinRfi = &limitMinRfi
	err = pdg.InsertQuote(quoteInsert1)
	if err != nil {
		LOGGER.Error(err)
		t.FailNow()
	}
	LOGGER.Debug(*quoteInsert1.QuoteID)
	// t.FailNow()
	quoteInsert2 := defaultQuote()
	testRfiID2 := "rfi2"
	quoteInsert2.QuoteResponse = &testQuoteResponse
	quoteInsert2.QuoteResponseBase64 = &testQuoteResponseBase64
	quoteInsert2.QuoteResponseSignature = &testQuoteSig
	quoteInsert2.RfiId = &testRfiID2
	err = pdg.InsertQuote(quoteInsert2)
	if err != nil {
		LOGGER.Error(err)
		t.FailNow()
	}
	timeOfQuote := time.Now().Unix()
	Convey("Update Quote should update one and only one row quote", t, func(c C) {
		LOGGER.Debug(*quoteInsert1.QuoteID)
		err := pdg.UpdateQuote(quoteInsert1, timeOfQuote)
		if err != nil {
			LOGGER.Error(err)
		}
		quotes, err := pdg.GetQuotes(*quoteInsert1.RequestID, *quoteInsert1.OfiId)
		LOGGER.Debug(*quotes[0].StatusQuote)
		for _, quote := range quotes {
			if *quote.RfiId == *quoteInsert1.RfiId {
				c.So(*quote.StatusQuote, ShouldEqual, 2)
			} else {
				c.So(*quote.StatusQuote, ShouldEqual, 1)
			}
		}
	})
	Convey("Execute Quote should not update quote as the amount is out of RFI's limit range", t, func(c C) {
		quoteResponse := testQuoteResponse
		timeExecuting := time.Now().Unix()
		err := pdg.ExecutingQuote(*quoteInsert1.QuoteID, *quoteInsert1.OfiId, []byte(quoteResponse), timeExecuting, decimal.NewFromFloat(20))
		c.So(err, ShouldNotBeNil)
	})

	Convey("Execute Quote should update one and only one row quote", t, func(c C) {
		quoteResponse := testQuoteResponse
		timeExecuting := time.Now().Unix()
		err := pdg.ExecutingQuote(*quoteInsert1.QuoteID, *quoteInsert1.OfiId, []byte(quoteResponse), timeExecuting, decimal.NewFromFloat(5))
		if err != nil {
			LOGGER.Error(err)
			c.So(err, ShouldBeNil)
			t.FailNow()
		}
		quotes, err := pdg.GetQuotes(*quoteInsert1.RequestID, *quoteInsert1.OfiId)
		if err != nil {
			LOGGER.Error(err)
			t.FailNow()
		}
		for _, quote := range quotes {
			if *quote.RfiId == *quoteInsert1.RfiId {
				c.So(*quote.StatusQuote, ShouldEqual, 3)
			} else {
				c.So(*quote.StatusQuote, ShouldEqual, 1)
			}
		}
	})

	timeExecuted := time.Now().Unix()
	Convey("Execute Quote should update one and only one row quote", t, func(c C) {
		err := pdg.ExecutedQuote(*quoteInsert1.QuoteID, *quoteInsert1.RfiId, timeExecuted)
		if err != nil {
			LOGGER.Error(err)
		}
		quotes, err := pdg.GetQuotes(*quoteInsert1.RequestID, *quoteInsert1.OfiId)
		for _, quote := range quotes {
			if *quote.RfiId == *quoteInsert1.RfiId {
				c.So(*quote.StatusQuote, ShouldEqual, 4)
			} else {
				c.So(*quote.StatusQuote, ShouldEqual, 1)
			}
		}
	})

	timeCancel := time.Now().Unix()
	Convey("Cancel Quote should not succeed", t, func(c C) {
		err := pdg.CancelQuote(*quoteInsert1.QuoteID, *quoteInsert1.RfiId, timeCancel)
		c.So(err, ShouldNotEqual, nil)
	})
}

func TestCancelQuote(t *testing.T) {
	pdg := PostgreDatabaseClient{
		Host:     host,
		Port:     port,
		User:     dbuser,
		Password: dbpassword,
		Dbname:   dbname,
	}
	err := pdg.CreateConnection()
	if err != nil {
		LOGGER.Error(err)
		t.Fail()
	}
	// defer pdg.CloseConnection()
	// create testing data
	quoteInsert := defaultQuote()
	testStatusQuote := 2
	quoteInsert.StatusQuote = &testStatusQuote
	err = pdg.InsertQuote(quoteInsert)
	if err != nil {
		LOGGER.Error(err)
		t.FailNow()
	}

	//cancel the quote
	timeCancel := time.Now().Unix()
	Convey("Cancel Quote should update one and only one row quote", t, func(c C) {
		err := pdg.CancelQuote(*quoteInsert.QuoteID, *quoteInsert.RfiId, timeCancel)
		if err != nil {
			LOGGER.Error(err)
		}
		quotes, err := pdg.GetQuotes(*quoteInsert.RequestID, *quoteInsert.OfiId)
		for _, quote := range quotes {
			if *quote.RfiId == *quoteInsert.RfiId {
				c.So(*quote.StatusQuote, ShouldEqual, 99)
			} else {
				c.So(*quote.StatusQuote, ShouldEqual, 1)
			}
		}
	})
}

func TestGetQuotesByAttributes(t *testing.T) {
	pdg := PostgreDatabaseClient{
		Host:     host,
		Port:     port,
		User:     dbuser,
		Password: dbpassword,
		Dbname:   dbname,
	}
	err := pdg.CreateConnection()
	if err != nil {
		LOGGER.Error(err)
		t.FailNow()
	}
	defer pdg.CloseConnection()
	quoteInsert := defaultQuote()
	exchangeRate := decimal.NewFromFloat(1.51)
	quoteInsert.ExchangeRate = &exchangeRate
	err = pdg.InsertQuote(quoteInsert)
	if err != nil {
		LOGGER.Error(err)
		t.FailNow()
	}
	query := Query{}
	jsonStr := `        {
		"ofi_id": "hk.one.payments.worldwire.io",
		"status_quote": {"operator": "gt", "threshold":0},
		"time_expire_rfi": {"operator": "gt", "threshold":1},
		"exchange_rate": {"operator": "eq", "threshold": 1.51},
		"source_asset": {
			"asset_code": "GBPDO",
			"asset_issuer": "GBPD3QEYU5E6ORGOY54WS2BQIFSJBM3GUQQOEPPYB4DZ2B5GQ7QHUVAV",
			"asset_type": "DO"
		},
		"target_asset": {
			"asset_code": "HKDDO",
			"asset_issuer": "GA3Z5DS6GAPBI6EGRFRCEJKAFBZHESW62U5ME3OJYVD5VEREY5ENTGIK",
			"asset_type": "DO"
		}
	}`
	json.Unmarshal([]byte(jsonStr), &query)
	// json.Unmarshal([]byte(`{"time_expire_rfi":"2018-11-26 02:22:49.62+00"}`), &query)
	quotes, err := pdg.GetQuotesByAttributes(&query)
	if err != nil {
		LOGGER.Error(err)
	}
	if len(quotes) == 0 {
		LOGGER.Error("No quote results")
		t.FailNow()
	}
	num := 0
	for _, quote := range quotes {
		if *quote.QuoteID == *quoteInsert.QuoteID {
			num = num + 1
		}
	}
	Convey("Get Quote by attribute should contain inserted quote", t, func(c C) {
		c.So(num, ShouldEqual, 1)
	})

}

func TestCancelQuotesByAttributes(t *testing.T) {
	pdg := PostgreDatabaseClient{
		Host:     host,
		Port:     port,
		User:     dbuser,
		Password: dbpassword,
		Dbname:   dbname,
	}
	err := pdg.CreateConnection()
	if err != nil {
		LOGGER.Error(err)
		t.FailNow()
	}
	defer pdg.CloseConnection()
	// jsonStr := `{"quote_id":"1asdf", "rfi_id : "domain"}`
	quoteInsert := defaultQuote()
	exchangeRate := decimal.NewFromFloat(1.51)
	quoteInsert.ExchangeRate = &exchangeRate
	err = pdg.InsertQuote(quoteInsert)
	if err != nil {
		LOGGER.Error(err)
		t.FailNow()
	}
	query := Query{}
	jsonStr := `        {
		"ofi_id": "hk.one.payments.worldwire.io",
		"status_quote": {"operator": "eq", "threshold":1},
		"source_asset": {
			"asset_code": "GBPDO",
			"asset_issuer": "GBPD3QEYU5E6ORGOY54WS2BQIFSJBM3GUQQOEPPYB4DZ2B5GQ7QHUVAV",
			"asset_type": "DO"
		},
		"target_asset": {
			"asset_code": "HKDDO",
			"asset_issuer": "GA3Z5DS6GAPBI6EGRFRCEJKAFBZHESW62U5ME3OJYVD5VEREY5ENTGIK",
			"asset_type": "DO"
		}
	}`
	json.Unmarshal([]byte(jsonStr), &query)
	timeCancel := time.Now().Unix()
	quotes, err := pdg.CancelQuotesByAttributes(&query, timeCancel)
	if err != nil {
		LOGGER.Error(err)
		t.FailNow()
	}
	if len(quotes) == 0 {
		LOGGER.Error("No quote results")
		t.FailNow()
	}
	num := 0
	for _, quote := range quotes {
		if *quote.QuoteID == *quoteInsert.QuoteID {
			num = num + 1
		}
	}
	Convey("Cancel Quote by attribute should contain inserted quote", t, func(c C) {
		c.So(num, ShouldEqual, 1)
	})

}

func TestGetQuoteByQuoteID(t *testing.T) {
	pdg := PostgreDatabaseClient{
		Host:     host,
		Port:     port,
		User:     dbuser,
		Password: dbpassword,
		Dbname:   dbname,
	}
	err := pdg.CreateConnection()
	if err != nil {
		LOGGER.Error(err)
		t.FailNow()
	}
	defer pdg.CloseConnection()

	requestID := "a6376649-34e7-42a0-87ae-b1a2009b9a1a"
	quoteID := "a6376649-34e7-42a0-87ae-b1a2009b9ac1aie.one.payments.worldwire.io"
	rfiID := "testRfi"
	ofiID := "testOfi"
	limitMaxOfi := decimal.NewFromFloat(1000.0)
	limitMinOfi := decimal.NewFromFloat(100.0)
	sourceAsset := types.JSONText(`{"asset_id" :"test"}`)
	targetAsset := types.JSONText(`{"asset_id" :"test"}`)
	timeRequest := time.Now().Unix()
	statusQuote := 1
	timeExpireOfi := int64(1796876024000)

	quote := QuoteDB{
		RequestID:     &requestID,
		QuoteID:       &quoteID,
		RfiId:         &rfiID,
		OfiId:         &ofiID,
		LimitMaxOfi:   &limitMaxOfi,
		LimitMinOfi:   &limitMinOfi,
		SourceAsset:   &sourceAsset,
		TargetAsset:   &targetAsset,
		TimeRequest:   &timeRequest,
		StatusQuote:   &statusQuote,
		TimeExpireOfi: &timeExpireOfi,
	}
	err = pdg.InsertQuote(quote)
	if err != nil {
		LOGGER.Error(err)
		t.FailNow()
	}

	// jsonStr := `{"quote_id":"1asdf", "rfi_id : "domain"}`
	quotes, err := pdg.GetQuoteByQuoteID(quoteID, rfiID)
	if err != nil {
		LOGGER.Error(err)
		t.FailNow()
	}
	Convey("Cancel Quote should update one and only one row quote", t, func(c C) {
		c.So(*quotes[0].QuoteID, ShouldEqual, quoteID)
	})

}
