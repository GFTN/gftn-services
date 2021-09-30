// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package utility

import (
	"encoding/json"
	"errors"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/GFTN/gftn-services/gftn-models/model"
)

func Neo4jQueryPreprocess(bytes []byte) (map[string]interface{}, error) {
	var input map[string]interface{}
	err := json.Unmarshal(bytes, &input)
	if err != nil {
		return map[string]interface{}{}, err
	}

	//skip writing these two attributes since we've already build the relationship inside Neo4j
	delete(input, "payout_parent")
	delete(input, "payout_child")
	for index, value := range input {
		valueType := reflect.ValueOf(value).Kind()

		// neo4j does not allow array of struct, so we have to process those struct into string
		if valueType == reflect.Slice {
			switch x := value.(type) {
			case []interface{}:
				if len(x) > 0 && reflect.ValueOf(x[0]).Kind() == reflect.Map {
					/*
						for i, e := range x {
							bytes, _ := json.Marshal(e)
							x[i] = string(bytes)
						}
					*/
					delete(input, index)
				}
			}
		}

		if valueType == reflect.Map {
			delete(input, index)
		}
	}
	return input, nil
}

func GraphDbInitialize(driver neo4j.Driver) error {
	var session neo4j.Session
	var err error
	if session, err = driver.Session(neo4j.AccessModeWrite); err != nil {
		return err
	}
	defer session.Close()

	// create index & constraints
	//The MATCH will use an index if possible.
	//If there is no index, it will lookup up all nodes carrying the label and see if the property matches.
	_, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`CREATE CONSTRAINT ON (a:Payout) ASSERT a.id IS UNIQUE`, map[string]interface{}{})
		if err != nil {
			return nil, err
		}
		if result.Next() {
			return result.Record().GetByIndex(0), nil
		}
		return nil, result.Err()
	})
	if err != nil {
		return err
	}
	return nil
}

func IsValid(payload model.PayoutLocation) error {
	geoType := strings.ToLower(*payload.Geo.Type)
	category := strings.ToLower(*payload.Category.Name)
	// check if geo coordinates match the payout point type
	if (len(payload.Geo.Coordinates) > 1 && geoType == GeoType[0]) || (len(payload.Geo.Coordinates) == 1 && geoType == GeoType[1]) {
		return errors.New("Geo coordinates does not match with the payout geo type")
	}

	// check area type match the category
	if (geoType == GeoType[0] && category != ReceiveMode[0]) || (geoType == GeoType[1] && category == ReceiveMode[0]) {
		return errors.New("Geo type does not match with the payout location category. 'cash_pickup' can only be 'point' type, and the rest will be 'area' type")
	}

	//no duplicate parent/child
	for _, parent := range payload.PayoutParent {
		if Contains(payload.PayoutChild, parent) {
			return errors.New("Payout parent cannot be payout child at the same time")
		}
	}

	// if category is point, no child should be contained
	if *payload.Geo.Type == GeoType[0] && len(payload.PayoutChild) > 0 {
		return errors.New("Payout point cannot contain another payout location")
	}

	return nil
}

func GenerateNodeID() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	b := make([]byte, id_length)
	for i := range b {
		b[i] = letterBytes[r1.Intn(len(letterBytes))]
	}
	return string(b)
}
