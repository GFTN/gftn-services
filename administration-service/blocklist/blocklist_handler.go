// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package blocklist

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/GFTN/gftn-services/administration-service/environment"
	"github.com/GFTN/gftn-services/administration-service/utility"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/utility/common"
	"github.com/GFTN/gftn-services/utility/database"
	"github.com/GFTN/gftn-services/utility/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ValidateType = []string{"CURRENCY", "COUNTRY", "INSTITUTION"}

type BlocklistOperations struct {
	session      *mongo.Client
	dbName       string
	dbCollection string
	conn         database.MongoDBConnect
}

func CreateBlocklistOperations() (BlocklistOperations, error) {
	bo := BlocklistOperations{}
	dbUser := os.Getenv(environment.ENV_KEY_DB_USER)
	dbPwd := os.Getenv(environment.ENV_KEY_DB_PWD)
	bo.dbName = os.Getenv(environment.ENV_KEY_ADMIN_DB_NAME)
	bo.dbCollection = os.Getenv(environment.ENV_KEY_BLOCKLIST_DB_TABLE)
	mongoId := os.Getenv(environment.ENV_KEY_MONGO_ID)

	if bo.dbCollection == "" || bo.dbName == "" {
		LOGGER.Warningf("Error reading DB table, environment variables ADMIN_DB_NAME and or BLOCKLIST_DB_TABLE not set")
		os.Exit(1)
	}

	LOGGER.Infof("\t* CreateBlocklistOperations dialing DB:%s collection:%s ", bo.dbName, bo.dbCollection)

	client, err := database.InitializeAtlasConnection(dbUser, dbPwd, mongoId)
	if err != nil {
		LOGGER.Errorf("Mongo Atlas DB connection failed! %s", err)
		panic("Mongo Atlas DB connection failed! " + err.Error())
	}
	bo.session = client

	LOGGER.Infof("\t* CreateBlocklistOperations DB is set")
	return bo, nil
}

func (bo BlocklistOperations) GetCollection() (*mongo.Collection, context.Context) {
	dbTimeout, _ := strconv.Atoi(os.Getenv(environment.ENV_KEY_DB_TIMEOUT))
	LOGGER.Infof("\t* Getting collection: %s from DB %s", bo.dbCollection, bo.dbName)
	collection := bo.session.Database(bo.dbName).Collection(bo.dbCollection)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(dbTimeout))
	return collection, ctx
}

func (bo BlocklistOperations) Add(w http.ResponseWriter, request *http.Request) {
	LOGGER.Infof("admin-service:Blocklist Operations :Add currency/country/participant into blocklist")

	var blocklistUpdateRq model.Blocklist

	err := json.NewDecoder(request.Body).Decode(&blocklistUpdateRq)
	if err != nil {
		LOGGER.Warningf("Error while decoding Blocklist update payload :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "ADMIN-0013", err)
		return
	}

	if blocklistUpdateRq.Type != nil {
		*blocklistUpdateRq.Type = strings.ToUpper(*blocklistUpdateRq.Type)
	}

	if len(blocklistUpdateRq.Value) == 0 {
		LOGGER.Warningf("Error while validating Blocklist payload : Value should contains at least 1 element")
		response.NotifyWWError(w, request, http.StatusBadRequest, "ADMIN-0014", err)
		return
	}

	err = blocklistUpdateRq.Validate(strfmt.Default)

	if err != nil {
		LOGGER.Warningf("Error while validating Blocklist payload :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "ADMIN-0014", err)
		return
	}

	collection, ctx := bo.GetCollection()

	var results []model.Blocklist
	cursor, err := collection.Find(ctx,
		bson.M{
			"type": blocklistUpdateRq.Type,
		},
	)
	if err != nil {
		LOGGER.Debugf("Error during Get blocklist query")
		response.NotifyWWError(w, request, http.StatusInternalServerError, "ADMIN-0015", err)
		return
	}
	bytes, err := database.ParseResult(cursor, ctx)
	if err != nil {
		LOGGER.Debugf("Error parsing mongo data")
		response.NotifyWWError(w, request, http.StatusNotFound, "ADMIN-0022", err)
		return
	}
	_ = json.Unmarshal(bytes, &results)

	if len(results) > 0 {
		LOGGER.Infof("The blocklist type already exists")
		existValues := results[0].Value

		for _, value := range blocklistUpdateRq.Value {
			value = strings.ToUpper(value)
			if utility.Contains(existValues, value) == -1 {
				existValues = append(existValues, value)
			} else {
				LOGGER.Warningf(value + " already exists in " + *blocklistUpdateRq.Type + " blocklist")
				response.NotifyWWError(w, request, http.StatusBadRequest, "ADMIN-0017", errors.New(value+" already exists in "+*blocklistUpdateRq.Type+" blocklist"))
				return
			}
		}

		_, err = collection.UpdateOne(ctx,
			bson.M{"type": results[0].Type},
			bson.M{
				"$set": bson.M{
					"value": existValues,
				},
			},
		)
		if err != nil {
			LOGGER.Warningf("Failed updating blocklist record %v", err)
			response.NotifyWWError(w, request, http.StatusInternalServerError, "ADMIN-0015", err)
			return
		}
		LOGGER.Infof("Update new blocklist element was successful")

	} else {
		LOGGER.Infof("Creating blocklist record type :%v", blocklistUpdateRq.Type)
		for key, value := range blocklistUpdateRq.Value {
			blocklistUpdateRq.Value[key] = strings.ToUpper(value)
		}
		blocklistUpdateRq.ID = primitive.NewObjectID().Hex()
		_, err = collection.InsertOne(ctx, blocklistUpdateRq)
		if err != nil {
			LOGGER.Warningf("Failed creating blocklist record %v", err)
			response.NotifyWWError(w, request, http.StatusInternalServerError, "ADMIN-0015", err)
			return
		}
		LOGGER.Infof("Create new blocklist record successfully")
	}

	response.Respond(w, http.StatusOK, []byte(`{"status":"Success"}`))

}

func (bo BlocklistOperations) Validate(w http.ResponseWriter, request *http.Request) {
	LOGGER.Infof("admin-service:Blocklist Operations :Validate currency/country/participant")

	var blocklistValidateRq []model.Blocklist

	err := json.NewDecoder(request.Body).Decode(&blocklistValidateRq)
	if err != nil {
		LOGGER.Warningf("Error while decoding Blocklist validate payload :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "ADMIN-0013", err)
		return
	}

	var valueArr = [][]string{[]string{}, []string{}, []string{}}
	for _, data := range blocklistValidateRq {

		if data.Type != nil {
			*data.Type = strings.ToUpper(*data.Type)
		} else {
			continue
		}

		if len(data.Value) == 0 {
			continue
		} else {
			for key, dataElem := range data.Value {
				data.Value[key] = strings.ToUpper(dataElem)
			}
		}

		err = data.Validate(strfmt.Default)

		if err != nil {
			LOGGER.Warningf("Error while validating Blocklist payload :  %v", err)
			response.NotifyWWError(w, request, http.StatusBadRequest, "ADMIN-0014", err)
			return
		}

		switch *data.Type {
		case ValidateType[0]:
			valueArr[0] = append(valueArr[0], data.Value...)
		case ValidateType[1]:
			valueArr[1] = append(valueArr[1], data.Value...)
		case ValidateType[2]:
			valueArr[2] = append(valueArr[2], data.Value...)
		default:
			LOGGER.Warningf("Error while parsing Blocklist validate payload : cannot recognize type %s", *data.Type)
			response.NotifyWWError(w, request, http.StatusBadRequest, "ADMIN-0013", err)
			return
		}
	}

	for key, validateType := range ValidateType {
		var results []model.Blocklist
		collection, ctx := bo.GetCollection()
		cursor, err := collection.Find(ctx, bson.M{
			"type": validateType,
			"value": bson.M{
				"$elemMatch": bson.M{
					"$in": valueArr[key],
				},
			},
		})
		if err != nil {
			LOGGER.Debugf("Error during Get blocklist query")
			response.NotifyWWError(w, request, http.StatusInternalServerError, "ADMIN-0015", err)
			return
		}
		bytes, err := database.ParseResult(cursor, ctx)
		if err != nil {
			LOGGER.Debugf("Error parsing mongo data")
			response.NotifyWWError(w, request, http.StatusNotFound, "ADMIN-0022", err)
			return
		}
		_ = json.Unmarshal(bytes, &results)

		if len(results) > 0 {
			response.Respond(w, http.StatusOK, []byte(common.BlocklistDeniedString))
			return
		}
	}

	LOGGER.Infof("transaction with %s: %s %s", valueArr[0], valueArr[1], valueArr[2])
	response.Respond(w, http.StatusOK, []byte(common.BlocklistApprovedString))
	return
}

func (bo BlocklistOperations) Remove(w http.ResponseWriter, request *http.Request) {
	LOGGER.Infof("admin-service:Blocklist Operations :Remove currency/country/participant from blocklist")

	var blocklistDeleteRq model.Blocklist

	err := json.NewDecoder(request.Body).Decode(&blocklistDeleteRq)
	if err != nil {
		LOGGER.Warningf("Error while decoding Blocklist update payload :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "ADMIN-0013", err)
		return
	}

	if blocklistDeleteRq.Type != nil {
		*blocklistDeleteRq.Type = strings.ToUpper(*blocklistDeleteRq.Type)
	}

	if len(blocklistDeleteRq.Value) == 0 {
		LOGGER.Warningf("Error while validating Blocklist payload : Value should contains at least 1 element")
		response.NotifyWWError(w, request, http.StatusBadRequest, "ADMIN-0014", err)
		return
	}

	err = blocklistDeleteRq.Validate(strfmt.Default)

	if err != nil {
		LOGGER.Warningf("Error while validating Blocklist payload :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "ADMIN-0014", err)
		return
	}

	collection, ctx := bo.GetCollection()
	var results []model.Blocklist

	cursor, err := collection.Find(ctx,
		bson.M{
			"type": blocklistDeleteRq.Type,
		},
	)
	if err != nil {
		LOGGER.Debugf("Error during Get blocklist query")
		response.NotifyWWError(w, request, http.StatusInternalServerError, "ADMIN-0015", err)
		return
	}
	bytes, err := database.ParseResult(cursor, ctx)
	if err != nil {
		LOGGER.Debugf("Error parsing mongo data")
		response.NotifyWWError(w, request, http.StatusNotFound, "ADMIN-0022", err)
		return
	}
	_ = json.Unmarshal(bytes, &results)

	if len(results) > 0 {
		LOGGER.Infof("The blocklist type already exists")
		existValues := results[0].Value

		for _, value := range blocklistDeleteRq.Value {
			value = strings.ToUpper(value)
			if index := utility.Contains(existValues, value); index > -1 {
				existValues[index] = existValues[len(existValues)-1] // Copy last element to index i
				existValues[len(existValues)-1] = ""                 // Erase last element (write zero value)
				existValues = existValues[:len(existValues)-1]       // Truncate slice
			} else {
				LOGGER.Warningf("Record not found")
				response.NotifyWWError(w, request, http.StatusNotFound, "ADMIN-0016", errors.New("Record not found"))
				return
			}
		}

		_, err = collection.UpdateOne(ctx,
			bson.M{"type": results[0].Type},
			bson.M{
				"$set": bson.M{
					"value": existValues,
				},
			},
		)
		if err != nil {
			LOGGER.Warningf("Delete blocklist element was not successful  %v", err)
			response.NotifyWWError(w, request, http.StatusInternalServerError, "ADMIN-0015", err)
			return
		}
		LOGGER.Infof("Delete blocklist element was successful")
		response.Respond(w, http.StatusOK, []byte(`{"status":"Success"}`))
	} else {
		LOGGER.Infof("Record not found")
		response.NotifyWWError(w, request, http.StatusNotFound, "ADMIN-0016", errors.New("Record not found"))
	}
	return
}

func (bo BlocklistOperations) Get(w http.ResponseWriter, request *http.Request) {
	LOGGER.Infof("admin-service:Blocklist Operations :Get currency/country/participant into blocklist")

	queryParams := request.URL.Query()

	var queryType interface{}

	if len(queryParams["type"]) <= 0 {
		queryType = bson.M{"$exists": true}
	} else {
		queryType = strings.ToUpper(queryParams["type"][0])
	}

	var results []model.Blocklist
	collection, ctx := bo.GetCollection()
	cursor, err := collection.Find(ctx,
		bson.M{
			"type": queryType,
		},
	)
	if err != nil {
		LOGGER.Debugf("Error during Get blocklist query")
		response.NotifyWWError(w, request, http.StatusInternalServerError, "ADMIN-0015", err)
		return
	}

	bytes, err := database.ParseResult(cursor, ctx)
	if err != nil {
		LOGGER.Debugf("Error parsing mongo data")
		response.NotifyWWError(w, request, http.StatusNotFound, "ADMIN-0022", err)
		return
	}
	_ = json.Unmarshal(bytes, &results)

	LOGGER.Infof("Get blocklist called, found %d matches", len(results))

	if len(results) == 0 {
		LOGGER.Warningf("No matching blocklist record found for type: %v", queryType)
	}
	b, _ := json.Marshal(results)
	response.Respond(w, http.StatusOK, b)

}
