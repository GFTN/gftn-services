// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package utility

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"googlemaps.github.io/maps"
)

func CreateNode(session neo4j.Session, payoutRq model.PayoutLocation) error {
	LOGGER.Infof("Creating node %s", payoutRq.ID)

	rawBytes, _ := json.Marshal(payoutRq)
	input, err := Neo4jQueryPreprocess(rawBytes)
	if err != nil {
		return nil
	}

	_, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		err = checkUniqueStatement(transaction, payoutRq)
		if err != nil {
			return nil, err
		}

		err := createNodeStatement(transaction, input)
		if err != nil {
			return nil, err
		}

		if len(payoutRq.PayoutParent) > 0 || len(payoutRq.PayoutChild) > 0 {
			err = createRelationshipStatement(transaction, payoutRq)
			if err != nil {
				return nil, err
			}
		}
		err = detectCycleStatement(transaction, payoutRq)
		if err != nil {
			return nil, err
		}

		err = createAttributeNodeStatement(transaction, rawBytes, payoutRq.ID)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	LOGGER.Infof("Node %s created!", payoutRq.ID)
	return err
}

func UpdateNode(session neo4j.Session, updateRq model.PayoutLocationUpdateRequest) error {
	LOGGER.Infof("Updating node %s", *updateRq.ID)

	rawBytes, _ := json.Marshal(*updateRq.UpdatedPayload)
	input, err := Neo4jQueryPreprocess(rawBytes)
	if err != nil {
		return nil
	}
	input["id"] = *updateRq.ID

	_, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		var payoutRq model.PayoutLocation
		payoutRq = *updateRq.UpdatedPayload
		payoutRq.ID = *updateRq.ID

		if len(payoutRq.PayoutParent) > 0 {
			err := validateStatement(transaction, payoutRq)
			if err != nil {
				return nil, err
			}
		}

		// delete attribute first, node itself will be the last,
		//or else we will not able to find it's relatives in this statement
		err = deleteAttributeNodeStatement(transaction, *updateRq.ID)
		if err != nil {
			return nil, err
		}

		err = deleteNodeStatement(transaction, *updateRq.ID)
		if err != nil {
			return nil, err
		}

		err = checkUniqueStatement(transaction, *updateRq.UpdatedPayload)
		if err != nil {
			return nil, err
		}

		err = createNodeStatement(transaction, input)
		if err != nil {
			return nil, err
		}

		if len(payoutRq.PayoutParent) > 0 || len(payoutRq.PayoutChild) > 0 {
			err = createRelationshipStatement(transaction, payoutRq)
			if err != nil {
				return nil, err
			}
		}
		err = detectCycleStatement(transaction, payoutRq)
		if err != nil {
			return nil, err
		}

		err = createAttributeNodeStatement(transaction, rawBytes, *updateRq.ID)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	LOGGER.Infof("Node %s updated!", *updateRq.ID)
	return err
}

func DeleteNode(session neo4j.Session, id string) error {
	LOGGER.Infof("Deleting node %s", id)

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		// delete attribute first, node itself will be the last,
		//or else we will not able to find it's relatives in this statement
		err := deleteAttributeNodeStatement(transaction, id)
		if err != nil {
			return nil, err
		}

		err = deleteNodeStatement(transaction, id)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	LOGGER.Infof("Node %s deleted!", id)
	return err
}

func GetNode(session neo4j.Session, mapClient *maps.Client, id string, name string, payoutType string, image string, url string, telephone string, currency string, child string, parent string, member string, receiveMode string, country string, state string, city string, street string, postalCode string, address string, geo string) (interface{}, error) {
	LOGGER.Infof("Getting node %s", id)
	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		var validLocations []model.PayoutLocation
		payoutLocations, err := getNodeStatement(transaction, id, name, payoutType, image, url, telephone, currency, child, parent, member, receiveMode, country, state, city, street, postalCode, geo)
		if err != nil {
			return nil, err
		}

		payoutLocations, err = getAttributeNodeStatement(transaction, payoutLocations, id)
		if err != nil {
			return nil, err
		}

		if address != "" && len(payoutLocations) > 0 {

			r := &maps.GeocodingRequest{
				Address: address,
			}

			resp, err := mapClient.Geocode(context.Background(), r)

			if len(resp) != 1 {
				LOGGER.Errorf("[Geocoding] Expected length of response is 1, was %+v", len(resp))
				return nil, errors.New("Error getting the geo coordinates of the specified address")
			}
			if err != nil {
				LOGGER.Errorf("[Geocoding] %v", err)
				return nil, err
			}

			var targetPoint model.Coordinate
			targetPoint.Long = &resp[0].Geometry.Location.Lng
			targetPoint.Lat = &resp[0].Geometry.Location.Lat
			LOGGER.Infof("[Geocoding] Successfully retrieve the geo code of %s -> Lng: %+v, Lat: %+v", address, *targetPoint.Long, *targetPoint.Lat)
			var point_candidates []string
			var area_candidates []string
			var raw_candidates []string
			for _, location := range payoutLocations {
				raw_candidates = append(raw_candidates, location.ID)
			}

			/* --------------
				payout area
			-----------------*/
			LOGGER.Infof("Filtering payout area that contains the target coordinate")

			var candidate_coordinates map[string][]model.Coordinate
			// the ids are the candidate fileted by all the previous criteria
			// retrieve the location which is only area and its coordinate details
			candidate_coordinates, err = retrieveAreaCoordinatesStatement(transaction, raw_candidates)
			if len(candidate_coordinates) == 0 {
				LOGGER.Warningf("Cannot find any specified payout area candidates")
			}

			if len(candidate_coordinates) > 0 {
				// check for all area
				//if contains the target point, include it
				area_candidates, err = ContainsPoint(raw_candidates, candidate_coordinates, targetPoint)
				if len(area_candidates) == 0 {
					LOGGER.Warningf("Cannot find any specified payout area candidates")
				}
			}

			/* --------------
			   payout point
			-----------------*/

			LOGGER.Infof("Sorting payout point that nears the target coordinate")
			point_candidates, err = retrieveNearestPointStatement(transaction, raw_candidates, 5, *targetPoint.Long, *targetPoint.Lat)
			if len(point_candidates) == 0 {
				LOGGER.Warningf("Cannot find any specified payout point candidates")
			}

			/* -----------------
			  merge & finalize
			-------------------*/

			var totalValidCandidates []string
			totalValidCandidates = append(totalValidCandidates, point_candidates...)
			totalValidCandidates = append(totalValidCandidates, area_candidates...)

			for index, location := range payoutLocations {
				if !Contains(totalValidCandidates, location.ID) {
					continue
				}
				//append the payout locations that are qualified to the final result
				validLocations = append(validLocations, payoutLocations[index])
			}
			if len(validLocations) == 0 {
				return nil, errors.New("Cannot find any specified payout location")
			}
			return validLocations, nil
		}

		return payoutLocations, nil
	})
	LOGGER.Infof("Node %s retrieved!", id)
	return result, err
}

func CreateNodesByCSV(session neo4j.Session, payoutRqs []model.PayoutLocation) error {
	LOGGER.Infof("Creating %+v nodes", len(payoutRqs))

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		err := checkDuplicateNodesStatement(transaction, payoutRqs)
		if err != nil {
			return nil, err
		}

		var inputs []map[string]interface{}
		var rawBytes []byte
		for _, payoutRq := range payoutRqs {
			rawBytes, _ = json.Marshal(payoutRq)
			input, err := Neo4jQueryPreprocess(rawBytes)
			if err != nil {
				return nil, err
			}
			inputs = append(inputs, input)
		}

		err = createNodesStatement(transaction, inputs)
		if err != nil {
			return nil, err
		}
		for _, payoutRq := range payoutRqs {

			rawBytes, _ := json.Marshal(payoutRq)
			if len(payoutRq.PayoutParent) > 0 || len(payoutRq.PayoutChild) > 0 {
				err = createRelationshipStatement(transaction, payoutRq)
				if err != nil {
					return nil, err
				}
			}
			err = detectCycleStatement(transaction, payoutRq)
			if err != nil {
				return nil, err
			}

			err = createAttributeNodeStatement(transaction, rawBytes, payoutRq.ID)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil

	})
	LOGGER.Infof("%+v Nodes from CSV created", len(payoutRqs))
	return err
}
