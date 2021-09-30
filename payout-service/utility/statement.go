// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package utility

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/GFTN/gftn-services/gftn-models/model"
)

func deleteNodeStatement(transaction neo4j.Transaction, id string) error {
	result, err := transaction.Run(
		`MATCH (node:Payout {id: $id}) DETACH DELETE node
			RETURN {status: "success"}`,
		map[string]interface{}{
			"id": id,
		})
	if err != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", err)
		return err
	}

	found := false
	if result.Next() {
		found = true
		returnedMap := result.Record().GetByIndex(0).(map[string]interface{})
		if returnedMap["status"].(string) != "success" {
			errMsg := "Encounter while deleting the node"
			LOGGER.Errorf("Encounter error during DB query: %s", errMsg)
			return errors.New(errMsg)
		}
	}
	if !found && result.Err() == nil {
		errMsg := "Cannot find the specified payout location"
		LOGGER.Errorf("Encounter error during DB query: %s", errMsg)
		return errors.New(errMsg)
	}

	if result.Err() != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", result.Err().Error())
		return result.Err()
	}
	LOGGER.Infof("Delete node statement success!")
	return nil
}

func deleteAttributeNodeStatement(transaction neo4j.Transaction, id string) error {
	result, err := transaction.Run(
		`MATCH (Payout:Payout {id:$id}) 
		MATCH (x) -[:has*1..1]- (Payout)
		OPTIONAL MATCH (y) -[:comprise*1..1]- (x)
		OPTIONAL MATCH (Payout) -[:opens*1..1]- (z)
		DETACH DELETE x, y, z`,
		map[string]interface{}{
			"id": id,
		})
	if err != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", err)
		return err
	}

	if result.Next() {
		returnedMap := result.Record().GetByIndex(0).(map[string]interface{})
		if returnedMap["status"].(string) != "success" {
			errMsg := "Encounter while deleting the attribute nodes"
			LOGGER.Errorf("Encounter error during DB query: %s", errMsg)
			return errors.New(errMsg)
		}
	}

	if result.Err() != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", result.Err().Error())
		return result.Err()
	}
	LOGGER.Infof("Delete attribute nodes statement success!")
	return nil
}

func createNodeStatement(transaction neo4j.Transaction, input map[string]interface{}) error {

	result, err := transaction.Run(
		`CREATE (currentNode:Payout $props) RETURN {status: "success"}`,
		map[string]interface{}{
			"props": input,
		})
	if err != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", err)
		return err
	}

	if result.Next() {
		returnedMap := result.Record().GetByIndex(0).(map[string]interface{})
		if returnedMap["status"].(string) != "success" {
			errMsg := "Encounter while creating the node"
			LOGGER.Errorf("Encounter error during DB query: %s", errMsg)
			return errors.New(errMsg)
		}
	}

	if result.Err() != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", result.Err().Error())
		return result.Err()
	}
	LOGGER.Infof("Create node statement success!")
	return nil
}

func createNodesStatement(transaction neo4j.Transaction, inputs []map[string]interface{}) error {

	var query string
	inputInterface := make(map[string]interface{}, len(inputs))
	query += "CREATE "
	for index, input := range inputs {
		i := strconv.Itoa(index)
		query += `(:Payout $props` + i + `)`
		if index != len(inputs)-1 {
			query += `,`
		}
		inputInterface["props"+i] = input
	}

	result, err := transaction.Run(query, inputInterface)
	if err != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", err)
		return err
	}

	if result.Next() {
		returnedMap := result.Record().GetByIndex(0).(map[string]interface{})
		if returnedMap["status"].(string) != "success" {
			errMsg := "Encounter while creating the node"
			LOGGER.Errorf("Encounter error during DB query: %s", errMsg)
			return errors.New(errMsg)
		}
	}

	if result.Err() != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", result.Err().Error())
		return result.Err()
	}
	LOGGER.Infof("Create node statement success!")
	return nil
}

func createRelationshipStatement(transaction neo4j.Transaction, payoutRq model.PayoutLocation) error {
	if len(payoutRq.PayoutParent) > 0 {
		result, err := transaction.Run(
			`MATCH (currentNode:Payout {id:$id})
			MATCH (parent:Payout) WHERE parent.id IN $parent
			MERGE (parent) -[:contains]-> (currentNode) 
			RETURN {parent: length(collect(distinct parent))} `,
			map[string]interface{}{
				"id":     payoutRq.ID,
				"parent": payoutRq.PayoutParent,
			})
		if err != nil {
			LOGGER.Errorf("Encounter error during DB query: %s", err)
			return err
		}

		if result.Next() {
			returnedMap := result.Record().GetByIndex(0).(map[string]interface{})
			parent := returnedMap["parent"].(int64)
			LOGGER.Infof("Found %+v parent", parent)
			if parent != int64(len(payoutRq.PayoutParent)) {
				errMsg := "Cannot find the given parent of the payout location"
				LOGGER.Errorf("Encounter error during DB query: %s", errMsg)
				return errors.New(errMsg)
			}
		}

		if result.Err() != nil {
			LOGGER.Errorf("Encounter error during DB query: %s", result.Err().Error())
			return result.Err()
		}
	}
	if len(payoutRq.PayoutChild) > 0 {
		result, err := transaction.Run(
			`MATCH (currentNode:Payout {id:$id})
			WITH currentNode
			MATCH (children:Payout) WHERE children.id IN $children
			MERGE (currentNode)-[:contains]->(children) 
			WITH *
			RETURN {children: length(collect(distinct children))} `,
			map[string]interface{}{
				"id":       payoutRq.ID,
				"children": payoutRq.PayoutChild,
			})
		if err != nil {
			LOGGER.Errorf("Encounter error during DB query: %s", err)
			return err
		}

		if result.Next() {
			returnedMap := result.Record().GetByIndex(0).(map[string]interface{})
			children := returnedMap["children"].(int64)
			LOGGER.Infof("Found %+v children", children)
			if children != int64(len(payoutRq.PayoutChild)) {
				errMsg := "Cannot find the given children of the payout location"
				LOGGER.Errorf("Encounter error during DB query: %s", errMsg)
				return errors.New(errMsg)
			}
		}

		if result.Err() != nil {
			LOGGER.Errorf("Encounter error during DB query: %s", result.Err().Error())
			return result.Err()
		}
	}
	LOGGER.Infof("Create relationship statement success!")
	return nil
}

func createAttributeNodeStatement(transaction neo4j.Transaction, bytes []byte, nodeId string) error {

	// we only allow at most two layer of nested json,
	//or else the whole structure will become too complex and too slow
	// for graph database

	var input map[string]interface{}
	err := json.Unmarshal(bytes, &input)
	if err != nil {
		return err
	}

	var propsMap = map[string]interface{}{
		"id": nodeId,
	}
	query := `MATCH (node:Payout {id:$id}) WITH * `

	// for struct
	// first layer attribute iteration
	for index, value := range input {
		valueType := reflect.ValueOf(value).Kind()
		if valueType == reflect.Map {

			bytes, _ := json.Marshal(value)
			input[index] = string(bytes)
			query += `CREATE (` + index + `:` + index + ` $` + index + `_props) 
				MERGE (node) -[:has]-> (` + index + `) `
			// second layer attribute iteration
			for inner_index, inner_value := range value.(map[string]interface{}) {
				inner_valueType := reflect.ValueOf(inner_value).Kind()
				if inner_valueType == reflect.Slice {
					switch x := inner_value.(type) {
					case []interface{}:
						// check if the array inside second layer is a normal string array or array of struct
						if len(x) > 0 && reflect.ValueOf(x[0]).Kind() == reflect.Map {
							delete(value.(map[string]interface{}), inner_index)
							for i, k := range x {
								innerId := inner_index + strconv.Itoa(i)
								query += `CREATE (` + innerId + `:` + inner_index + ` $` + innerId + `_props) 
									MERGE (` + innerId + `) -[:comprise]-> (` + index + `) `
								propsMap[innerId+"_props"] = k.(map[string]interface{})
							}
						}
					}
				}

			}
			propsMap[index+"_props"] = value
		} else if valueType == reflect.Slice {
			// for array of struct
			arr := value.([]interface{})
			if len(arr) > 0 && reflect.ValueOf(arr[0]).Kind() == reflect.Map {
				for i, k := range arr {
					innerId := index + strconv.Itoa(i)
					query += `CREATE (` + innerId + `:` + index + ` $` + innerId + `_props) 
					MERGE (node) -[:opens]-> (` + innerId + `) `
					propsMap[innerId+"_props"] = k.(map[string]interface{})
					propsMap[innerId] = i
				}
			}
		}
	}

	query += ` RETURN {status:"success"}`

	result, err := transaction.Run(query, propsMap)
	if err != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", err)
		return err
	}

	if result.Next() {
		returnedMap := result.Record().GetByIndex(0).(map[string]interface{})
		if returnedMap["status"].(string) != "success" {
			errMsg := "Failed creating attribute nodes"
			LOGGER.Errorf("Encounter error during DB query: %s", errMsg)
			return errors.New(errMsg)
		}
	}

	if result.Err() != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", result.Err().Error())
		return result.Err()
	}

	LOGGER.Infof("Adding attribute node statement success!")
	return nil
}

func detectCycleStatement(transaction neo4j.Transaction, payoutRq model.PayoutLocation) error {
	result, err := transaction.Run(
		`MATCH (currentNode:Payout {id:$id})
			WITH currentNode
			OPTIONAL MATCH path=(currentNode)-[*]-> (currentNode)
			WHERE length(path) >0
			RETURN {cycle: length(path)} LIMIT 1`,
		map[string]interface{}{
			"id": payoutRq.ID,
		})
	if err != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", err)
		return err
	}

	if result.Next() {
		returnedMap := result.Record().GetByIndex(0).(map[string]interface{})
		if returnedMap["cycle"] != nil && returnedMap["cycle"].(int64) > 0 {
			errMsg := "Cycle detected! Please Change the payout_child/payout_parent of the payout location " + *payoutRq.Name
			LOGGER.Errorf("Encounter error during DB query: %s", errMsg)
			return errors.New(errMsg)
		}
	}

	if result.Err() != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", result.Err().Error())
		return result.Err()
	}
	LOGGER.Infof("Detect cycle statement success!")
	return nil
}

func getNodeStatement(
	transaction neo4j.Transaction,
	id string,
	name string,
	payoutType string,
	image string,
	url string,
	telephone string,
	currency string,
	child string,
	parent string,
	member string,
	receiveMode string,
	country string,
	state string,
	city string,
	street string,
	postalCode string,
	geo string,
) ([]model.PayoutLocation, error) {
	query := `MATCH (location:Payout) WHERE location.id =~ $id AND 
		location.name =~ $name AND
		location.type =~ $type AND
		location.image =~ $image AND
		location.url =~ $url AND
		location.telephone =~ $telephone AND
		location.image =~ $image AND
		ANY(member IN location.member_of WHERE member =~ $member) AND
		ANY(currency IN location.currencies_accepted WHERE currency =~ $currency)
		MATCH (category:category) -[:has*1]-(location)
		WHERE category.name =~ $receiveMode
		MATCH (geo:geo) -[:has*1]-(location)
		WHERE geo.type =~ $geo
		MATCH (location) -[:has*1]-> (addr:address)
		WHERE addr.country =~ $country AND addr.state =~ $state AND addr.city =~ $city AND addr.postal_code =~ $postal_code AND addr.street =~ $street`
	if parent != "" {
		query += ` MATCH (parent_query:Payout) -[:contains*1]-> (location) WHERE parent_query.id =~ $parent`
	}
	if child != "" {
		query += ` MATCH (location) -[:contains*1]-> (child_query:Payout) WHERE child_query.id =~ $child`
	}

	query += ` OPTIONAL MATCH (location)-[:contains]->(children)
		OPTIONAL MATCH (parent)-[:contains]->(location)
		WITH *
		MATCH (x) -[:has*1..1]- (location)
		OPTIONAL MATCH (location) -[:opens*1]- (y)
		RETURN {location:location, parent:collect(distinct parent.id), children:collect(distinct children.id),attr:collect(x), opening_hours:collect(distinct y)}`

	result, err := transaction.Run(
		query,
		map[string]interface{}{
			"id":          id,
			"name":        name,
			"type":        payoutType,
			"image":       image,
			"url":         url,
			"telephone":   telephone,
			"currency":    currency,
			"child":       child,
			"parent":      parent,
			"member":      member,
			"receiveMode": receiveMode,
			"country":     country,
			"state":       state,
			"city":        city,
			"street":      street,
			"postal_code": postalCode,
			"geo":         geo,
		})
	if err != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", err)
		return nil, err
	}

	var payoutLocations []model.PayoutLocation
	found := false
	for result.Next() {
		found = true
		returnedMap := result.Record().GetByIndex(0).(map[string]interface{})
		location := returnedMap["location"].(neo4j.Node).Props()

		//restoring payout_child, payout_parent
		raw_children := returnedMap["children"].([]interface{})
		raw_parent := returnedMap["parent"].([]interface{})
		children := make([]string, len(raw_children))
		for i, v := range raw_children {
			children[i] = fmt.Sprint(v)
		}
		parent := make([]string, len(raw_parent))
		for i, v := range raw_parent {
			parent[i] = fmt.Sprint(v)
		}

		// restoring geo, address, category, opening hours attribute
		var payoutGeo model.Geo
		var payoutAddress model.Address
		var payoutCategory model.PayoutLocationCategory

		attr := returnedMap["attr"].([]interface{})
		for _, v := range attr {
			b, _ := json.Marshal(v.(neo4j.Node).Props())
			switch v.(neo4j.Node).Labels()[0] {
			case "geo":
				json.Unmarshal(b, &payoutGeo)
			case "category":
				json.Unmarshal(b, &payoutCategory)
			case "address":
				json.Unmarshal(b, &payoutAddress)
			}
		}

		// restoring opening_hours
		var payoutOpeningHrs []*model.PayoutLocationOpeningHour
		opening := returnedMap["opening_hours"].([]interface{})
		for _, time := range opening {
			var payoutOpeningHr model.PayoutLocationOpeningHour
			b, _ := json.Marshal(time.(neo4j.Node).Props())
			json.Unmarshal(b, &payoutOpeningHr)
			payoutOpeningHrs = append(payoutOpeningHrs, &payoutOpeningHr)
		}

		var payoutLocation model.PayoutLocation
		b, _ := json.Marshal(location)
		json.Unmarshal(b, &payoutLocation)
		payoutLocation.PayoutParent = parent
		payoutLocation.PayoutChild = children
		payoutLocation.Category = &payoutCategory
		payoutLocation.Address = &payoutAddress
		payoutLocation.Geo = &payoutGeo
		payoutLocation.OpeningHours = payoutOpeningHrs
		payoutLocations = append(payoutLocations, payoutLocation)
	}
	if !found && result.Err() == nil {
		errMsg := "Cannot find the specified payout location"
		LOGGER.Errorf("Encounter error during DB query: %s", errMsg)
		return nil, errors.New(errMsg)
	}

	if result.Err() != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", result.Err().Error())
		return nil, result.Err()
	}
	LOGGER.Infof("Get node statement success!")
	return payoutLocations, nil
}

func getAttributeNodeStatement(transaction neo4j.Transaction, payoutLocations []model.PayoutLocation, id string) ([]model.PayoutLocation, error) {

	var ids []string
	for _, v := range payoutLocations {
		ids = append(ids, v.ID)
	}

	result, err := transaction.Run(
		`MATCH (location:Payout) WHERE location.id IN $ids
		MATCH (x) -[:has*1]- (location)
		OPTIONAL MATCH (y) -[:comprise*1]- (x)
		RETURN {id:location.id, attr:collect(x), sub_attr:collect(y)}`, map[string]interface{}{"ids": ids})
	if err != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", err)
		return nil, err
	}

	found := false
	for result.Next() {
		found = true

		returnedMap := result.Record().GetByIndex(0).(map[string]interface{})
		var id string
		var targetIndex int
		id = returnedMap["id"].(string)

		for i, v := range payoutLocations {
			if v.ID == id {
				targetIndex = i
				break
			}
		}

		// restoring sub attributes of geo, address, category

		sub_attr := returnedMap["sub_attr"].([]interface{})

		for _, sub_v := range sub_attr {
			b, _ := json.Marshal(sub_v.(neo4j.Node).Props())
			var payoutOption model.PayoutLocationOption
			var payoutCoordinate model.Coordinate
			switch sub_v.(neo4j.Node).Labels()[0] {
			case "options":
				json.Unmarshal(b, &payoutOption)
				payoutLocations[targetIndex].Category.Options = append(payoutLocations[targetIndex].Category.Options, &payoutOption)
			case "coordinates":
				json.Unmarshal(b, &payoutCoordinate)
				payoutLocations[targetIndex].Geo.Coordinates = append(payoutLocations[targetIndex].Geo.Coordinates, &payoutCoordinate)
			}
		}
	}
	if !found && result.Err() == nil {
		errMsg := "Cannot find the attribute node of the payout location"
		LOGGER.Errorf("Encounter error during DB query: %s", errMsg)
		return nil, errors.New(errMsg)
	}

	if result.Err() != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", result.Err().Error())
		return nil, result.Err()
	}
	LOGGER.Infof("Get attribute node statement success!")
	return payoutLocations, nil

}

func checkUniqueStatement(transaction neo4j.Transaction, payoutRq model.PayoutLocation) error {

	var query string
	if *payoutRq.Category.Name == "area" {
		query = `MATCH (location:Payout) WHERE location.name = $name`
	} else {
		query = `MATCH (location:Payout) WHERE location.name = $name
				MATCH (addr:address) WHERE addr.building_number = $buildingNumber AND addr.city = $city AND addr.country = $country AND addr.state=$state AND addr.street = $street AND addr.postal_code = $postalCode`
	}
	query += ` RETURN {status: "found"}`

	result, err := transaction.Run(
		query, map[string]interface{}{
			"name":           *payoutRq.Name,
			"buildingNumber": *payoutRq.Address.BuildingNumber,
			"city":           *payoutRq.Address.City,
			"country":        *payoutRq.Address.Country,
			"state":          *payoutRq.Address.State,
			"street":         *payoutRq.Address.Street,
			"postalCode":     *payoutRq.Address.PostalCode,
		})
	if err != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", err)
		return err
	}

	for result.Next() {
		errMsg := "Duplicate node found! Please change the name/address of the payout location " + *payoutRq.Name
		LOGGER.Errorf("Encounter error during DB query: %s", errMsg)
		return errors.New(errMsg)
	}

	if result.Err() != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", result.Err().Error())
		return result.Err()
	}
	LOGGER.Infof("No duplicate node detected!")
	return nil

}

func checkDuplicateNodesStatement(transaction neo4j.Transaction, payoutRqs []model.PayoutLocation) error {

	var query string
	returnStatement := `RETURN {`
	params := make(map[string]interface{}, len(payoutRqs)*7)
	for index, payoutRq := range payoutRqs {
		i := strconv.Itoa(index)
		if *payoutRq.Category.Name == "area" {
			query += `OPTIONAL MATCH (location` + i + `:Payout) WHERE location` + i + `.name = $name` + i + ` `
		} else {
			query += `OPTIONAL MATCH (location` + i + `:Payout) WHERE location` + i + `.name = $name` + i + `
				OPTIONAL MATCH (addr` + i + `:address) -[*1]- (location` + i + `) 
				WHERE addr` + i + `.building_number = $buildingNumber` + i + ` AND 
				addr` + i + `.city = $city` + i + ` AND 
				addr` + i + `.country = $country` + i + ` AND 
				addr` + i + `.state=$state` + i + ` AND 
				addr` + i + `.street = $street` + i + ` AND 
				addr` + i + `.postal_code = $postalCode` + i + ` `
		}
		returnStatement += `location` + i + `:location` + i + `.name`
		if index != len(payoutRqs)-1 {
			returnStatement += `,`
		}
		params["name"+i] = *payoutRq.Name
		params["buildingNumber"+i] = *payoutRq.Address.BuildingNumber
		params["city"+i] = *payoutRq.Address.City
		params["country"+i] = *payoutRq.Address.Country
		params["state"+i] = *payoutRq.Address.State
		params["street"+i] = *payoutRq.Address.Street
		params["postalCode"+i] = *payoutRq.Address.PostalCode
	}
	returnStatement += `}`
	query += returnStatement

	result, err := transaction.Run(query, params)
	if err != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", err)
		return err
	}

	for result.Next() {
		returnedMap := result.Record().GetByIndex(0).(map[string]interface{})
		for _, entry := range returnedMap {
			if entry != nil {
				res := entry.(string)
				errMsg := "Duplicate node found! Please change the name/address of the payout location(name:" + res + ")"
				LOGGER.Errorf("Encounter error during DB query: %s", errMsg)
				return errors.New(errMsg)
			}
		}
	}

	if result.Err() != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", result.Err().Error())
		return result.Err()
	}
	LOGGER.Infof("No duplicate node detected!")
	return nil

}

func validateStatement(transaction neo4j.Transaction, payoutRq model.PayoutLocation) error {

	var parent_ids []string
	for _, v := range payoutRq.PayoutParent {
		parent_ids = append(parent_ids, v)
	}

	result, err := transaction.Run(
		`MATCH (parent:Payout) WHERE parent.id IN $parent_id
		MATCH (parent) -[:has]-> (parent_geo:geo) WHERE parent_geo.type <> "point"
        RETURN {status:"success"}`, map[string]interface{}{"parent_id": parent_ids})
	if err != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", err)
		return err
	}

	if result.Next() {
		LOGGER.Infof("Payout location payload validated!")
		return nil
	}

	if result.Err() != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", result.Err().Error())
		return result.Err()
	} else {
		errMsg := "Parent node contains a geo type of 'point', but payout point cannot be a parent node"
		LOGGER.Errorf(errMsg)
		return errors.New(errMsg)
	}
	return nil

}

func retrieveAreaCoordinatesStatement(transaction neo4j.Transaction, ids []string) (map[string][]model.Coordinate, error) {

	finalCoordinate := make(map[string][]model.Coordinate, len(ids))

	result, err := transaction.Run(
		`MATCH (location:Payout) WHERE location.id IN $ids
		MATCH (geo:geo) -[*1]- (location) WHERE geo.type = "area"
		MATCH (x:coordinates) -[*2]- (location)
		RETURN {id:location.id, coordinates:collect(x)}`, map[string]interface{}{"ids": ids})
	if err != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", err)
		return nil, err
	}

	found := false
	for result.Next() {
		found = true

		returnedMap := result.Record().GetByIndex(0).(map[string]interface{})
		var id string
		var coordinates []model.Coordinate
		id = returnedMap["id"].(string)

		temp := returnedMap["coordinates"].([]interface{})
		for _, v := range temp {
			var coordinate model.Coordinate
			prop := v.(neo4j.Node).Props()
			b, _ := json.Marshal(prop)
			_ = json.Unmarshal(b, &coordinate)
			coordinates = append(coordinates, coordinate)
		}

		finalCoordinate[id] = coordinates

	}
	if !found && result.Err() == nil {
		errMsg := "Cannot retrieve the coordinates detail of the payout area"
		LOGGER.Errorf("Encounter error during DB query: %s", errMsg)
		return nil, errors.New(errMsg)
	}

	if result.Err() != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", result.Err().Error())
		return nil, result.Err()
	}
	LOGGER.Infof("Get attribute node statement success!")
	return finalCoordinate, nil

}

func retrieveNearestPointStatement(transaction neo4j.Transaction, ids []string, limit int, longitude float64, latitude float64) ([]string, error) {

	var chosen_locations []string

	result, err := transaction.Run(
		`MATCH (x:coordinates) -[*1]- (geo:geo) -[*1]- (node:Payout)
		WHERE geo.type="point" AND node.id IN $ids
		WITH *,AVG(x.lat) as candidateLat, AVG(x.long) as candidateLong
		WITH *, point({ longitude: candidateLong, latitude: candidateLat }) AS candidates, point({ longitude: $long, latitude: $lat }) AS source
		WITH *, round(distance(candidates, source)) AS dist
		return {dist: dist, id: node.id} ORDER BY dist ASC LIMIT $limit`,
		map[string]interface{}{
			"ids":   ids,
			"limit": limit,
			"long":  longitude,
			"lat":   latitude,
		})
	if err != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", err)
		return nil, err
	}

	found := false
	for result.Next() {
		found = true
		returnedMap := result.Record().GetByIndex(0).(map[string]interface{})
		var id string
		id = returnedMap["id"].(string)
		dist := returnedMap["dist"].(float64)
		LOGGER.Infof("For payout point %s, the distance to the target point is %+v", id, dist)

		chosen_locations = append(chosen_locations, id)
	}
	if !found && result.Err() == nil {
		errMsg := "Cannot find the nearest payout point of the target payout location"
		LOGGER.Errorf("Encounter error during DB query: %s", errMsg)
		return nil, errors.New(errMsg)
	}

	if result.Err() != nil {
		LOGGER.Errorf("Encounter error during DB query: %s", result.Err().Error())
		return nil, result.Err()
	}
	LOGGER.Infof("Get attribute node statement success!")
	return chosen_locations, nil

}
