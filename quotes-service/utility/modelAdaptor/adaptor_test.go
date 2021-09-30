// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package modeladaptor

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/GFTN/gftn-services/gftn-models/model"

	. "github.com/smartystreets/goconvey/convey"
)

// compare if 2 json string is deep equal
func IsEqualJson(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	err := json.Unmarshal([]byte(s1), &o1)

	if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(s1), &o2)

	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(o1, o2), nil
}
func TestQueryToQueryDB(t *testing.T) {
	quoteFilter := model.QuoteFilter{}
	jsonStr := `        {
		"ofi_id": "hk.one.payments.worldwire.io",
		"status_quote": {"operator": "gt", "threshold":0},
		"time_expire_rfi": {"operator": "gt", "threshold":0},
		"exchange_rate": {"operator": "eq", "threshold": 1.50},
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
	json.Unmarshal([]byte(jsonStr), &quoteFilter)
	queryDB := QueryToQueryDB(&quoteFilter)
	queryDBJson, _ := json.Marshal(queryDB)
	Convey("information shall be conserved", t, func(c C) {
		fmt.Println(string(queryDBJson))
		flag, _ := IsEqualJson(string(queryDBJson), jsonStr)
		c.So(flag, ShouldEqual, true)

	})
}
