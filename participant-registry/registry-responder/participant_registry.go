// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package participantregistry

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	prc "github.com/GFTN/gftn-services/participant-registry-client/pr-client"

	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/participant-registry/environment"
	"github.com/GFTN/gftn-services/utility/database"
	"github.com/GFTN/gftn-services/utility/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Operations struct {
	session      *mongo.Client
	dbName       string
	dbCollection string
}

func CreateParticipantRegistryOperations() (Operations, error) {

	LOGGER.Infof("\t* CreateParticipantRegistryOperations connecting Mongo DB ")
	pr := Operations{}
	dbUser := os.Getenv(environment.ENV_KEY_DB_USER)
	dbPwd := os.Getenv(environment.ENV_KEY_DB_PWD)
	pr.dbName = os.Getenv(environment.ENV_KEY_PR_DB_NAME)
	pr.dbCollection = os.Getenv(environment.ENV_KEY_PARTICIPANTS_DB_TABLE)
	mongoId := os.Getenv(environment.ENV_KEY_MONGO_ID)

	if pr.dbCollection == "" || pr.dbName == "" {
		errMsg := "Error reading DB table, environment variables PR_DB_NAME and or PARTICIPANTS_DB_TABLE not set"
		LOGGER.Errorf(errMsg)
		panic(errMsg)
	}

	client, err := database.InitializeAtlasConnection(dbUser, dbPwd, mongoId)
	if err != nil {
		LOGGER.Errorf("Mongo Atlas DB connection failed! %s", err)
		panic("Mongo Atlas DB connection failed! " + err.Error())
	}
	pr.session = client
	LOGGER.Infof("\t* CreateParticipantRegistryOperations DB is set")

	return pr, nil
}

func (pr Operations) GetCollection() (*mongo.Collection, context.Context) {
	dbTimeout, _ := strconv.Atoi(os.Getenv(environment.ENV_KEY_DB_TIMEOUT))
	LOGGER.Infof("\t* Getting collection: %s from DB %s", pr.dbCollection, pr.dbName)
	collection := pr.session.Database(pr.dbName).Collection(pr.dbCollection)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(dbTimeout))
	return collection, ctx
}

func (pr Operations) GetParticipants(w http.ResponseWriter, request *http.Request) {
	var results []model.Participant

	collection, ctx := pr.GetCollection()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		LOGGER.Debugf("Error during GetParticipants query")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	}

	bytes, err := database.ParseResult(cursor, ctx)
	if err != nil {
		LOGGER.Debugf("Error parsing mongo data")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1007", err)
		return
	}

	_ = json.Unmarshal(bytes, &results)
	LOGGER.Infof("GetParticipants called, found %d matches", len(results))
	b, _ := json.Marshal(results)

	response.Respond(w, http.StatusOK, b)
}

func (pr Operations) GetParticipantsByCountry(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	countryCode := vars["country_code"]
	var results []model.Participant

	if countryCode == "" {
		err := errors.New("The participant country code in the request should not be empty")
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1001", err)
		return
	}
	countryCode = strings.ToUpper(countryCode)

	collection, ctx := pr.GetCollection()
	cursor, err := collection.Find(ctx, bson.M{"country_code": countryCode})
	if err != nil {
		LOGGER.Debugf("Error during GetParticipants query")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	}

	bytes, err := database.ParseResult(cursor, ctx)
	if err != nil {
		LOGGER.Debugf("Error parsing mongo data")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1007", err)
		return
	}

	_ = json.Unmarshal(bytes, &results)

	LOGGER.Infof("GetParticipantsByCountry called, found %d matches", len(results))

	if len(results) == 0 {
		LOGGER.Warning("No matching participant found for country:", countryCode)
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	}

	b, _ := json.Marshal(results)
	response.Respond(w, http.StatusOK, b)
}

func (pr Operations) CreateParticipant(w http.ResponseWriter, request *http.Request) {

	var participantRq model.Participant

	err := json.NewDecoder(request.Body).Decode(&participantRq)
	if err != nil {
		LOGGER.Warningf("Error while validating Participant :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1001", err)
		return
	}

	err = participantRq.Validate(strfmt.Default)

	if err != nil {
		LOGGER.Warningf("Error while validating Participant :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1001", err)
		return
	}

	if err != nil {
		LOGGER.Warningf("Create participant was not successful  %v", err)
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1001", err)
		return
	}

	collection, ctx := pr.GetCollection()
	cursor, err := collection.Find(ctx, bson.M{"issuing_account": participantRq.IssuingAccount, "bic": participantRq.Bic})
	if err != nil {
		LOGGER.Debugf("Error during CreateParticipants query")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	}
	var results []model.Participant

	bytes, err := database.ParseResult(cursor, ctx)
	if err != nil {
		LOGGER.Debugf("Error parsing mongo data")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1007", err)
		return
	}

	_ = json.Unmarshal(bytes, &results)

	if len(results) > 0 {
		LOGGER.Errorf("The participant exists with given participant address or BIC code")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1003", err)
		return
	}

	*participantRq.ID = strings.ToLower(*participantRq.ID)
	*participantRq.CountryCode = strings.ToUpper(*participantRq.CountryCode)

	cursor, err = collection.Find(ctx, bson.M{"id": *participantRq.ID})
	if err != nil {
		LOGGER.Debugf("Error during CreateParticipants query")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	}
	bytes, err = database.ParseResult(cursor, ctx)
	if err != nil {
		LOGGER.Debugf("Error parsing mongo data")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1007", err)
		return
	}
	_ = json.Unmarshal(bytes, &results)

	if len(results) > 0 {
		LOGGER.Errorf("The participant exists with given participant domain")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1003", errors.New("The participant exists with given participant domain"))
		return
	}

	//Set the active status to inactive first time you create a participant
	inActiveStatus := prc.PR_INACTIVE
	participantRq.Status = inActiveStatus
	_, err = collection.InsertOne(ctx, participantRq)
	if err != nil {
		LOGGER.Warningf("Create participant was not successful  %v", err)
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1003", err)
		return
	}

	LOGGER.Infof("Create participant was successful")
	b, _ := json.Marshal(participantRq)
	response.Respond(w, http.StatusOK, b)
}

func (pr Operations) GetParticipantDomain(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	participantDomain := vars["participant_id"]
	var results []model.Participant

	if participantDomain == "" {
		err := errors.New("The participant domain should not be empty")
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1001", err)
		return
	}
	participantDomain = strings.ToLower(participantDomain)

	collection, ctx := pr.GetCollection()
	cursor, err := collection.Find(ctx, bson.M{"id": participantDomain})
	if err != nil {
		LOGGER.Warning("No matching participant found for domain:", participantDomain)
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1002", err)
		return
	}
	bytes, err := database.ParseResult(cursor, ctx)
	if err != nil {
		LOGGER.Debugf("Error parsing mongo data")
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1007", err)
		return
	}

	_ = json.Unmarshal(bytes, &results)

	LOGGER.Infof("GetParticipantsParticipantDomain called, found %d Matches for %s", len(results), participantDomain)

	if len(results) > 1 {
		LOGGER.Warning("The participant address should be unique, something is wrong for requested domain", participantDomain)
		err = errors.New("The participant address should be unique, something is wrong for requested domain")
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1002", err)
		return
	} else if len(results) == 0 {
		LOGGER.Warning("No matching participant found for domain:", participantDomain)
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	}

	b, _ := json.Marshal(results[0])
	response.Respond(w, http.StatusOK, b)
}

func (pr Operations) GetParticipantDistAccount(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	participantDomain := vars["participant_id"]

	accountName := vars["account_name"]
	var results []model.Participant

	if participantDomain == "" || accountName == "" {
		err := errors.New("The participant domain and account name in the request should not be empty")
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1001", err)
		return
	}
	participantDomain = strings.ToLower(participantDomain)

	collection, ctx := pr.GetCollection()
	cursor, err := collection.Find(ctx, bson.M{
		"id": participantDomain,
		"operating_accounts": bson.M{
			"$elemMatch": bson.M{
				"name": accountName,
			},
		},
	})
	if err != nil {
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	}
	bytes, err := database.ParseResult(cursor, ctx)
	if err != nil {
		LOGGER.Debugf("Error parsing mongo data")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1007", err)
		return
	}
	_ = json.Unmarshal(bytes, &results)

	LOGGER.Infof("GetParticipantsParticipantDomain called, found %d Matches for %s", len(results), participantDomain)

	if len(results) > 1 {
		LOGGER.Warning("The participant address should be unique, something is wrong for requested domain", participantDomain)
		err = errors.New("The participant address should be unique, something is wrong for requested domain")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	} else if len(results) == 0 {
		LOGGER.Warning("No matching participant found for domain:", participantDomain)
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	}
	if len(results[0].OperatingAccounts) == 0 {
		err = errors.New(accountName)
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1005", err)
		return
	}

	var disPubKey string
	for _, account := range results[0].OperatingAccounts {
		if account.Name == accountName {
			disPubKey = *account.Address
		}
	}
	if disPubKey == "" {
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1005", err)
		return
	}
	response.Respond(w, http.StatusOK, []byte(disPubKey))
}

func (pr Operations) GetParticipantForIssuingAccount(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	accountAddress := vars["account_address"]

	var results []model.Participant

	if accountAddress == "" {
		err := errors.New("The issuing account address in the request should not be empty")
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1001", err)
		return
	}
	collection, ctx := pr.GetCollection()
	cursor, err := collection.Find(ctx, bson.M{"issuing_account": accountAddress})

	if err != nil {
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	}

	bytes, err := database.ParseResult(cursor, ctx)
	if err != nil {
		LOGGER.Debugf("Error parsing mongo data")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1007", err)
		return
	}
	_ = json.Unmarshal(bytes, &results)

	LOGGER.Infof("GetParticipantsParticipantForIssuingAccount called, found %d Matches for %s", len(results), accountAddress)

	if len(results) > 1 {
		LOGGER.Warning("The participant issuing address should be unique, something is wrong for requested address", accountAddress)
		err = errors.New("The participant issuing address should be unique, something is wrong for requested address")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	} else if len(results) == 0 {
		LOGGER.Warning("No matching participant found for issuing address:", accountAddress)
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	}

	b, _ := json.Marshal(results[0])
	response.Respond(w, http.StatusOK, b)
}

func (pr Operations) SaveParticipantDistAccount(w http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	participantDomain := vars["participant_id"]
	var participant model.Participant
	var participantRq model.Account
	var results []model.Participant

	err := json.NewDecoder(request.Body).Decode(&participantRq)
	if err != nil {
		LOGGER.Warningf("Error while validating Participant Operating Account request :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1001", err)
		return
	}
	err = participantRq.Validate(strfmt.Default)

	participantDomain = strings.ToLower(participantDomain)
	if err != nil {
		LOGGER.Warningf("Error while validating Participant Operating Account request :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1001", err)
		return
	}

	collection, ctx := pr.GetCollection()
	cursor, err := collection.Find(ctx,
		bson.M{
			"id": participantDomain,
			"operating_accounts": bson.M{
				"$elemMatch": bson.M{
					"name": participantRq.Name,
				},
			},
		},
	)

	if err != nil {
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1001", err)
		return
	}

	bytes, err := database.ParseResult(cursor, ctx)
	if err != nil {
		LOGGER.Debugf("Error parsing mongo data")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1007", err)
		return
	}
	_ = json.Unmarshal(bytes, &results)

	LOGGER.Infof("GetParticipantsParticipantDomain called, found %d Matches for %s", len(results), participantDomain)
	if len(results) > 0 {
		err = errors.New(participantRq.Name)
		response.NotifyWWError(w, request, http.StatusConflict, "PR-1004", err)
		return
	}

	cursor, err = collection.Find(ctx,
		bson.M{
			"id": participantDomain,
			"operating_accounts": bson.M{
				"$elemMatch": bson.M{
					"address": *participantRq.Address,
				},
			},
		},
	)

	if err != nil {
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1001", err)
		return
	}

	bytes, err = database.ParseResult(cursor, ctx)
	if err != nil {
		LOGGER.Debugf("Error parsing mongo data")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1007", err)
		return
	}
	_ = json.Unmarshal(bytes, &results)

	LOGGER.Infof("GetParticipantsParticipantDomain called, found %d Matches for %s", len(results), participantDomain)
	if len(results) > 0 {
		err = errors.New(participantRq.Name)
		response.NotifyWWError(w, request, http.StatusConflict, "PR-1004", err)
		return
	}

	err = collection.FindOne(ctx,
		bson.M{"id": participantDomain}).Decode(&participant)
	participant.OperatingAccounts = append(participant.OperatingAccounts, &model.Account{
		participantRq.Address, participantRq.Name})

	_, err = collection.UpdateOne(ctx, bson.M{"id": participantDomain}, bson.M{"$set": participant})
	if err != nil {
		LOGGER.Warningf("Saving participant operating account was not successful  %v", err)
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1003", err)
		return
	}

	LOGGER.Infof("Saving participant operating account was successful")
	response.NotifySuccess(w, request, "Saving participant operating account was successful")
}

func (pr Operations) SaveParticipantIssuingAccount(w http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	participantDomain := vars["participant_id"]
	var participant model.Participant
	var participantRq model.Account
	var results []model.Participant

	err := json.NewDecoder(request.Body).Decode(&participantRq)

	if err != nil {
		LOGGER.Warningf("Error while validating Participant Issuing Account request :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1001", err)
		return
	}
	err = participantRq.Validate(strfmt.Default)

	if err != nil {
		LOGGER.Warningf("Error while validating Participant Issuing Account request :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1001", err)
		return
	}

	participantDomain = strings.ToLower(participantDomain)
	collection, ctx := pr.GetCollection()
	cursor, err := collection.Find(ctx,
		bson.M{"id": participantDomain})

	if err != nil {
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	}

	bytes, err := database.ParseResult(cursor, ctx)
	if err != nil {
		LOGGER.Debugf("Error parsing mongo data")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1007", err)
		return
	}

	_ = json.Unmarshal(bytes, &results)
	LOGGER.Infof("SaveParticipantsParticipantIssuingAccount called, found %d Matches for %s", len(results), participantDomain)

	if len(results) > 1 {
		err = errors.New("The participant address should be unique, something is wrong for requested domain")
		response.NotifyWWError(w, request, http.StatusConflict, "PR-1002", err)
		return
	} else if len(results) == 0 {
		err = errors.New("The participant does not exists with given participant domain")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	}

	cursor, err = collection.Find(ctx,
		bson.M{"issuing_account": *participantRq.Address})

	if err != nil {
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	}

	bytes, err = database.ParseResult(cursor, ctx)
	if err != nil {
		LOGGER.Debugf("Error parsing mongo data")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1007", err)
		return
	}
	_ = json.Unmarshal(bytes, &results)
	LOGGER.Infof("SaveParticipantsParticipantIssuingAccount called, found %d Matches for %s", len(results), *participantRq.Address)

	if len(results) > 0 {
		err = errors.New("The participant issuing should be unique, something is wrong for requested issuing account")
		response.NotifyWWError(w, request, http.StatusConflict, "PR-1002", err)
		return
	}

	err = collection.FindOne(ctx,
		bson.M{"id": participantDomain}).Decode(&participant)

	participant.IssuingAccount = *participantRq.Address
	_, err = collection.UpdateOne(ctx, bson.M{"id": participantDomain}, bson.M{"$set": participant})
	if err != nil {
		LOGGER.Warningf("Saving participant issuing account was not successful  %v", err)
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1006", err)
		return
	}

	LOGGER.Infof("Saving participant issuing account was successful")
	response.NotifySuccess(w, request, "Saving participant issuing account was successful")
}

func (pr Operations) UpdateParticipant(w http.ResponseWriter, request *http.Request) {

	var participantRq model.Participant
	vars := mux.Vars(request)
	participantDomain := vars["participant_id"]

	err := json.NewDecoder(request.Body).Decode(&participantRq)
	if err != nil {
		LOGGER.Warningf("Error while validating Participant :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1001", err)
		return
	}
	err = participantRq.Validate(strfmt.Default)
	if err != nil {
		LOGGER.Warningf("Error while validating Participant :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1001", err)
		return
	}
	participantDomain = strings.ToLower(participantDomain)

	if err != nil {
		LOGGER.Warningf("Update participant was not successful  %v", err)
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1006", err)
		return
	}

	collection, ctx := pr.GetCollection()

	_, err = collection.UpdateOne(ctx, bson.M{"id": participantDomain}, bson.M{"$set": participantRq})
	if err != nil {
		LOGGER.Warningf("Update participant was not successful  %v", err)
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1006", err)
		return
	}

	LOGGER.Infof("Update participant was successful")
	b, _ := json.Marshal(participantRq)
	response.Respond(w, http.StatusOK, b)

}

func (pr Operations) UpdateParticipantStatus(w http.ResponseWriter, request *http.Request) {

	var participant model.Participant
	var participantRq model.ParticipantStatus
	vars := mux.Vars(request)
	participantDomain := vars["participant_id"]
	var results []model.Participant

	err := json.NewDecoder(request.Body).Decode(&participantRq)
	if err != nil {
		LOGGER.Warningf("Error while validating Participant Status :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1001", err)
		return
	}
	err = participantRq.Validate(strfmt.Default)

	if err != nil {
		LOGGER.Warningf("Error while validating Participant Status :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1001", err)
		return
	}
	participantDomain = strings.ToLower(participantDomain)

	collection, ctx := pr.GetCollection()
	cursor, err := collection.Find(ctx,
		bson.M{"id": participantDomain})

	if err != nil {
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	}

	bytes, err := database.ParseResult(cursor, ctx)
	if err != nil {
		LOGGER.Debugf("Error parsing mongo data")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1007", err)
		return
	}

	_ = json.Unmarshal(bytes, &results)
	LOGGER.Infof("SaveParticipantsParticipantIssuingAccount called, found %d Matches for %s", len(results), participantDomain)

	if len(results) > 1 {
		err = errors.New("The participant address should be unique, something is wrong for requested domain")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	} else if len(results) == 0 {
		err = errors.New("The participant does not exists with given participant domain")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	}

	//copy participant data and update status
	participant = results[0]
	participant.Status = *participantRq.Status

	LOGGER.Debugf("Updating participant %v status to %v", *participant.ID, participant.Status)

	_, err = collection.UpdateOne(ctx, bson.M{"id": participantDomain}, bson.M{"$set": participant})
	if err != nil {
		LOGGER.Warningf("Update participantStatus was not successful  %v", err)
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1006", err)
		return
	}

	LOGGER.Infof("Update participantStatus was successful")
	b, _ := json.Marshal(participant)
	response.Respond(w, http.StatusOK, b)

}

func (pr Operations) GetParticipantByAddress(w http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	address := vars["account_address"]
	if address == "" {
		response.NotifyWWError(w, request, http.StatusBadRequest, "PR-1001", errors.New("account_address missing when calling GetParticipantByAddress endpoint"))
		return
	}

	var result model.Participant
	collection, ctx := pr.GetCollection()
	err := collection.FindOne(ctx,
		bson.D{
			{"$or",
				bson.A{
					bson.M{
						"operating_accounts": bson.M{
							"$elemMatch": bson.M{
								"address": address,
							},
						},
					},
					bson.M{
						"issuing_account": address,
					},
				},
			},
		},
	).Decode(&result)

	if err != nil {
		LOGGER.Debugf("Error during GetParticipantByAddress DB query")
		response.NotifyWWError(w, request, http.StatusNotFound, "PR-1002", err)
		return
	}

	LOGGER.Infof("GetParticipantByAddress success")
	LOGGER.Infof("Update participantStatus was successful")
	b, _ := json.Marshal(result)
	response.Respond(w, http.StatusOK, b)
}
