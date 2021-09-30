// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/payout-service/environment"
	"github.com/GFTN/gftn-services/payout-service/utility"
	"github.com/GFTN/gftn-services/utility/response"
	"googlemaps.github.io/maps"
)

type PayoutPointOperations struct {
	driver    neo4j.Driver
	mapClient *maps.Client
}

func CreatePayoutPointOperations() (PayoutPointOperations, error) {
	operation := PayoutPointOperations{}
	dbUser := os.Getenv(environment.ENV_KEY_DB_USER)
	dbPwd := os.Getenv(environment.ENV_KEY_DB_PWD)
	dbUrl := os.Getenv(environment.ENV_KEY_DB_URL)
	apikey := os.Getenv(environment.ENV_KEY_GEOCODING_API_KEY)

	if dbUser == "" || dbPwd == "" || dbUrl == "" || apikey == "" {
		LOGGER.Warningf("Environment variables DB_USER, DB_PWD, DB_URL, or GEOCODING_API_KEY not set")
		os.Exit(1)
	}

	LOGGER.Infof("* CreatePayoutPointOperations dialing Graph DB:%s user:%s ", dbUrl, dbUser)
	driver, err := neo4j.NewDriver(dbUrl, neo4j.BasicAuth(dbUser, dbPwd, ""))
	if err != nil {
		LOGGER.Errorf("Neo4j graph DB connection failed! %s", err)
		panic("Neo4j graph DB connection failed! " + err.Error())
	}
	//defer driver.Close()

	operation.driver = driver
	LOGGER.Infof("* Initializing graph DB constraints & index")
	err = utility.GraphDbInitialize(driver)
	if err != nil {
		LOGGER.Errorf("Neo4j graph DB initializing failed! %s", err)
		panic("Neo4j graph DB initializing failed! " + err.Error())
	}
	LOGGER.Infof("* Graph DB Initialized!")

	LOGGER.Infof("* Initializing Google Map client")
	operation.mapClient, err = maps.NewClient(maps.WithAPIKey(apikey))
	if err != nil {
		LOGGER.Errorf("Google Map client initializing failed! %s", err)
		panic("Google Map client initializing failed: " + err.Error())
	}
	LOGGER.Infof("* Google Map API initialized!")

	LOGGER.Infof("* CreatePayoutPointOperations DB is set")
	return operation, nil
}

func (operation PayoutPointOperations) UpdatePayout(w http.ResponseWriter, request *http.Request) {
	LOGGER.Infof("payout-service:Payout Operations :update payout")

	var payoutUpdateRq model.PayoutLocationUpdateRequest

	err := json.NewDecoder(request.Body).Decode(&payoutUpdateRq)
	if err != nil {
		LOGGER.Warningf("Error while decoding Payout Point payload :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "PAYOUT-1001", err)
		return
	}

	err = payoutUpdateRq.Validate(strfmt.Default)
	if err != nil {
		LOGGER.Warningf("Error while validating Payout Point payload :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "PAYOUT-1002", err)
		return
	}

	err = utility.IsValid(*payoutUpdateRq.UpdatedPayload)
	if err != nil {
		LOGGER.Warningf("Error while validating Payout Point payload :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "PAYOUT-1013", err)
		return
	}

	// db operation

	var session neo4j.Session
	if session, err = operation.driver.Session(neo4j.AccessModeWrite); err != nil {
		LOGGER.Warningf("Error while establishing session to Neo4j graph database: %v", err)
		response.NotifyWWError(w, request, http.StatusInternalServerError, "PAYOUT-1018", err)
		return
	}
	defer session.Close()

	LOGGER.Infof("updatePayout with id: %s", *payoutUpdateRq.ID)

	// update node
	err = utility.UpdateNode(session, payoutUpdateRq)
	if err != nil {
		LOGGER.Errorf("Update payout location failed")
		response.NotifyWWError(w, request, http.StatusBadRequest, "PAYOUT-1007", err)
		return
	}
	LOGGER.Infof("Update new payout point successfully")
	response.Respond(w, http.StatusOK, []byte(`{"status":"Success"}`))
	return

}

func (operation PayoutPointOperations) AddPayout(w http.ResponseWriter, request *http.Request) {
	LOGGER.Infof("payout-service:Payout Operations :add payout")

	var payoutRq model.PayoutLocation

	err := json.NewDecoder(request.Body).Decode(&payoutRq)
	if err != nil {
		LOGGER.Warningf("Error while decoding Payout Point payload :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "PAYOUT-1001", err)
		return
	}

	err = payoutRq.Validate(strfmt.Default)
	if err != nil {
		LOGGER.Warningf("Error while validating Payout Point payload :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "PAYOUT-1002", err)
		return
	}

	err = utility.IsValid(payoutRq)
	if err != nil {
		LOGGER.Warningf("Error while validating Payout Point payload :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "PAYOUT-1013", err)
		return
	}

	// graph db operation
	var session neo4j.Session
	if session, err = operation.driver.Session(neo4j.AccessModeWrite); err != nil {
		LOGGER.Warningf("Error while establishing session to Neo4j graph database: %v", err)
		response.NotifyWWError(w, request, http.StatusInternalServerError, "PAYOUT-1018", err)
		return
	}
	defer session.Close()

	payoutRq.ID = utility.GenerateNodeID()
	LOGGER.Infof("Adding payout location with id: %s", payoutRq.ID)

	// create node
	err = utility.CreateNode(session, payoutRq)
	if err != nil {
		LOGGER.Errorf("Create payout location failed")
		response.NotifyWWError(w, request, http.StatusBadRequest, "PAYOUT-1007", err)
		return
	}

	LOGGER.Infof("Create new payout location was successful")
	response.Respond(w, http.StatusOK, []byte(`{"id":"`+payoutRq.ID+`"}`))
	return

}

func (operation PayoutPointOperations) DeletePayout(w http.ResponseWriter, request *http.Request) {
	LOGGER.Infof("payout-service:Payout Operations :delete payout")

	queryParams := request.URL.Query()
	if len(queryParams["id"]) <= 0 {
		LOGGER.Warningf("Failed deleting payout point")
		response.NotifyWWError(w, request, http.StatusBadRequest, "PAYOUT-1011", nil)
		return
	}
	id := queryParams["id"][0]

	// graph db operation
	var session neo4j.Session
	var err error
	if session, err = operation.driver.Session(neo4j.AccessModeWrite); err != nil {
		LOGGER.Warningf("Error while establishing session to Neo4j graph database: %v", err)
		response.NotifyWWError(w, request, http.StatusInternalServerError, "PAYOUT-1018", err)
		return
	}
	defer session.Close()

	err = utility.DeleteNode(session, id)
	if err != nil {
		LOGGER.Errorf("Delete payout location failed")
		response.NotifyWWError(w, request, http.StatusBadRequest, "PAYOUT-1007", err)
		return
	}
	LOGGER.Infof("Payout point %s successfully deleted", id)
	response.Respond(w, http.StatusOK, []byte(`{"status":"success"}`))
	return
}

func (operation PayoutPointOperations) GetPayout(w http.ResponseWriter, request *http.Request) {

	LOGGER.Infof("payout-service:Payout Operations :get payout")
	queryParams := request.URL.Query()
	var parent string
	var child string
	var address string
	id := ".*"
	name := ".*"
	payoutType := ".*"
	image := ".*"
	url := ".*"
	telephone := ".*"
	currency := ".*"
	member := ".*"
	receiveMode := ".*"
	country := ".*"
	city := ".*"
	state := ".*"
	postalCode := ".*"
	street := ".*"
	geo := ".*"

	if len(queryParams["id"]) > 0 {
		id = `(?i)` + queryParams["id"][0]
	}

	if len(queryParams["name"]) > 0 {
		name = `(?i)` + queryParams["name"][0]
	}

	if len(queryParams["type"]) > 0 {
		payoutType = `(?i)` + queryParams["type"][0]
	}

	if len(queryParams["image"]) > 0 {
		image = `(?i)` + queryParams["image"][0]
	}

	if len(queryParams["url"]) > 0 {
		url = `(?i)` + queryParams["url"][0]
	}

	if len(queryParams["telephone"]) > 0 {
		telephone = `(?i)` + queryParams["telephone"][0]
	}

	if len(queryParams["currency"]) > 0 {
		currency = `(?i)` + queryParams["currency"][0]
	}

	if len(queryParams["child"]) > 0 {
		child = `(?i)` + queryParams["child"][0]
	}
	if len(queryParams["parent"]) > 0 {
		parent = `(?i)` + queryParams["parent"][0]
	}
	if len(queryParams["member"]) > 0 {
		member = `(?i)` + queryParams["member"][0]
	}

	if len(queryParams["receive_mode"]) > 0 {
		receiveMode = `(?i)` + queryParams["receive_mode"][0]
	}

	if len(queryParams["city"]) > 0 {
		city = `(?i)` + queryParams["city"][0]
	}

	if len(queryParams["state"]) > 0 {
		state = `(?i)` + queryParams["state"][0]
	}

	if len(queryParams["street"]) > 0 {
		street = `(?i)` + queryParams["street"][0]
	}

	if len(queryParams["country"]) > 0 {
		country = `(?i)` + queryParams["country"][0]
	}

	if len(queryParams["postal_code"]) > 0 {
		postalCode = `(?i)` + queryParams["postal_code"][0]
	}

	if len(queryParams["address"]) > 0 {
		address = queryParams["address"][0]
	}

	if len(queryParams["geo"]) > 0 {
		geo = queryParams["geo"][0]
	}

	// graph db operation
	var session neo4j.Session
	var err error
	if session, err = operation.driver.Session(neo4j.AccessModeWrite); err != nil {
		LOGGER.Errorf("Error while establishing session to Neo4j graph database: %v", err)
		response.NotifyWWError(w, request, http.StatusInternalServerError, "PAYOUT-1018", err)
		return
	}
	defer session.Close()

	result, err := utility.GetNode(session, operation.mapClient, id, name, payoutType, image, url, telephone, currency, child, parent, member, receiveMode, country, state, city, street, postalCode, address, geo)
	if err != nil {
		LOGGER.Errorf("Get payout location failed: %s", err.Error())
		response.NotifyWWError(w, request, http.StatusBadRequest, "PAYOUT-1007", err)
		return
	}

	b, _ := json.Marshal(result)
	response.Respond(w, http.StatusOK, b)

}

func (operation PayoutPointOperations) AddPayoutCSV(w http.ResponseWriter, request *http.Request) {

	var Buf bytes.Buffer
	file, header, err := request.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	LOGGER.Infof("File name %s\n", name[0])
	// Copy the file data to buffer
	io.Copy(&Buf, file)
	content := Buf.String()
	// reset the buffer to reduce memory allocations in more intense projects
	Buf.Reset()

	contents := strings.Split(content, "\n")

	headersArr := strings.Split(contents[0], ",")

	if len(contents) < 1 {
		log.Fatal("Something wrong, the file maybe empty or length of the lines are not the same")
	}

	var results []model.PayoutLocation
	respond := "{"
	for row, csv := range contents {
		//skip head row
		if row == 0 {
			continue
		}
		var entry model.PayoutLocation
		if strings.Contains(csv, `"`) {
			response.NotifyWWError(w, request, http.StatusBadRequest, "PAYOUT-1020", errors.New("CSV value cannot contains comma"))
			return
		}

		attributes := strings.Split(csv, ",")

		for index, attribute := range attributes {

			attr := attribute

			headerType := strings.Split(headersArr[index], utility.CSV_Head_Delimiter)
			for i, headerElem := range headerType {
				headerType[i] = strings.TrimSpace(strings.ToLower(headerElem))
				headerType[i] = strings.Replace(headerType[i], "\uFEFF", "", -1)
			}

			switch headerType[0] {
			case "geo":
				if headerType[1] == "type" {
					if attr == "" {
						entry.Geo = nil
					} else {
						entry.Geo = &model.Geo{Type: &attr}
					}
				} else {
					if attribute == "" || entry.Geo == nil {
						continue
					}
					raw := strings.Split(attr, utility.CSV_Value_Delimiter)
					lng, _ := strconv.ParseFloat(raw[0], 64)
					lat, _ := strconv.ParseFloat(raw[1], 64)
					var coordinate = model.Coordinate{Long: &lng, Lat: &lat}
					entry.Geo.Coordinates = append(entry.Geo.Coordinates, &coordinate)
				}
			case "opening_hours":
				if attribute == "" {
					continue
				}
				raw := strings.Split(attr, utility.CSV_Value_Delimiter)
				open := raw[0]
				close := raw[1]
				timeExists := false
				if len(entry.OpeningHours) > 0 {
					for innerIndex, hr := range entry.OpeningHours {
						if *hr.Opens == open && *hr.Closes == close {
							entry.OpeningHours[innerIndex].DayOfWeek = append(entry.OpeningHours[innerIndex].DayOfWeek, headerType[1])
							timeExists = true
							break
						}
					}
				}
				if !timeExists {
					var openHour model.PayoutLocationOpeningHour
					openHour.DayOfWeek = append(openHour.DayOfWeek, headerType[1])
					openHour.Opens = &open
					openHour.Closes = &close
					entry.OpeningHours = append(entry.OpeningHours, &openHour)
				}

			case "category":
				if headerType[1] == "name" {
					if attr == "" {
						entry.Category = nil
					} else {
						var category model.PayoutLocationCategory
						category.Name = &attr
						entry.Category = &category
					}
				} else if headerType[1] == "description" {
					if attribute == "" || entry.Category == nil {
						continue
					}
					i, _ := strconv.Atoi(headerType[2])
					if len(entry.Category.Options) > i && entry.Category.Options[i].Terms != nil {
						entry.Category.Options[i].Description = &attr
					} else {
						var option model.PayoutLocationOption
						option.Description = &attr
						entry.Category.Options = append(entry.Category.Options, &option)
					}
				} else if headerType[1] == "terms" {
					if attribute == "" || entry.Category == nil {
						continue
					}
					i, _ := strconv.Atoi(headerType[2])
					if len(entry.Category.Options) > i && entry.Category.Options[i].Description != nil {
						entry.Category.Options[i].Terms = &attr
					} else {
						var option model.PayoutLocationOption
						option.Terms = &attr
						entry.Category.Options = append(entry.Category.Options, &option)
					}
				}
			case "type":
				if attr == "" {
					entry.Type = nil
				} else {
					entry.Type = &attr
				}
			case "address":
				if attribute == "" {
					continue
				}
				if entry.Address == nil {
					var address model.Address
					entry.Address = &address
				}
				if headerType[1] == "street" {
					entry.Address.Street = &attr
				} else if headerType[1] == "state" {
					entry.Address.State = &attr
				} else if headerType[1] == "city" {
					entry.Address.City = &attr
				} else if headerType[1] == "country" {
					entry.Address.Country = &attr
				} else if headerType[1] == "building_number" {
					entry.Address.BuildingNumber = &attr
				} else if headerType[1] == "postal_code" {
					entry.Address.PostalCode = &attr
				} else {
					response.NotifyWWError(w, request, http.StatusBadRequest, "PAYOUT-1021", errors.New("Unknown field "+headerType[0]+"/"+headerType[1]))
					return
				}

			case "payout_child":
				attrs := strings.Split(attr, utility.CSV_Value_Delimiter)
				if attrs[0] == "" {
					entry.PayoutChild = []string{}
					continue
				}
				for _, elem := range attrs {
					entry.PayoutChild = append(entry.PayoutChild, elem)
				}
			case "payout_parent":
				attrs := strings.Split(attr, utility.CSV_Value_Delimiter)
				if attrs[0] == "" {
					entry.PayoutParent = []string{}
					continue
				}
				for _, elem := range attrs {
					entry.PayoutParent = append(entry.PayoutParent, elem)
				}
			case "name":
				entry.Name = &attr
			case "currency_accepted":
				attrs := strings.Split(attr, utility.CSV_Value_Delimiter)
				if attrs[0] == "" {
					entry.CurrenciesAccepted = nil
					continue
				}
				for _, elem := range attrs {
					entry.CurrenciesAccepted = append(entry.CurrenciesAccepted, elem)
				}
			case "image":
				if attr == "" {
					entry.Image = nil
				} else {
					entry.Image = &attr
				}
			case "url":
				if attr == "" {
					entry.URL = nil
				} else {
					entry.URL = &attr
				}
			case "telephone":
				if attr == "" {
					entry.Telephone = nil
				} else {
					entry.Telephone = &attr
				}
			case "member_of":
				attrs := strings.Split(attr, utility.CSV_Value_Delimiter)
				if attrs[0] == "" {
					entry.MemberOf = nil
					continue
				}
				for _, elem := range attrs {
					entry.MemberOf = append(entry.MemberOf, elem)
				}
			default:
				response.NotifyWWError(w, request, http.StatusBadRequest, "PAYOUT-1021", errors.New("Unknown field "+strings.Join(headerType, "/")))
				return
			}

		}

		err = entry.Validate(strfmt.Default)
		if err != nil {
			errMsg := "Error while validating Payout location payload on row " + strconv.Itoa(row+1) + " : " + err.Error() + ""
			LOGGER.Errorf(errMsg)
			response.NotifyWWError(w, request, http.StatusBadRequest, "PAYOUT-1002", errors.New(errMsg))
			return
		}

		err = utility.IsValid(entry)
		if err != nil {
			errMsg := "Error while validating Payout location payload on row " + strconv.Itoa(row+1) + " : " + err.Error() + ""
			LOGGER.Errorf(errMsg)
			response.NotifyWWError(w, request, http.StatusBadRequest, "PAYOUT-1013", errors.New(errMsg))
			return
		}

		entry.ID = utility.GenerateNodeID()
		respond += `"Row #` + strconv.Itoa(row+1) + `":"` + entry.ID + `"`
		if row != len(contents)-1 {
			respond += ","
		}
		results = append(results, entry)
	}
	respond += "}"
	// graph db operation
	var session neo4j.Session
	if session, err = operation.driver.Session(neo4j.AccessModeWrite); err != nil {
		LOGGER.Errorf("Error while establishing session to Neo4j graph database: %v", err)
		response.NotifyWWError(w, request, http.StatusInternalServerError, "PAYOUT-1018", err)
		return
	}
	defer session.Close()

	err = utility.CreateNodesByCSV(session, results)
	if err != nil {
		LOGGER.Errorf("Create payout location failed")
		response.NotifyWWError(w, request, http.StatusBadRequest, "PAYOUT-1007", err)
		return
	}

	response.Respond(w, http.StatusOK, []byte(respond))
}
