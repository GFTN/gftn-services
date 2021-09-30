// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package persistence

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/GFTN/gftn-services/administration-service/environment"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/utility/database"
	"github.com/GFTN/gftn-services/utility/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBOperations struct {
	session      *mongo.Client
	dbName       string
	dbCollection string
}

const (
	//INSTRUCTION_ID ID for message
	INSTRUCTION_ID = "INSTRUCTION_ID"
	//TRANSACTION_ID id from ledger
	TRANSACTION_ID = "TRANSACTION_ID"
	//DATE_RANGE type for range query
	DATE_RANGE = "DATE_RANGE"
)

func CreateAdminServicePersistenceOperations() (MongoDBOperations, error) {

	mo := MongoDBOperations{}
	dbUser := os.Getenv(environment.ENV_KEY_DB_USER)
	dbPwd := os.Getenv(environment.ENV_KEY_DB_PWD)
	mo.dbName = os.Getenv(environment.ENV_KEY_ADMIN_DB_NAME)
	mo.dbCollection = os.Getenv(environment.ENV_KEY_PAYMENTS_DB_TABLE)
	mongoId := os.Getenv(environment.ENV_KEY_MONGO_ID)

	if mo.dbCollection == "" || mo.dbName == "" {
		LOGGER.Warningf("Error reading DB table, environment variables ADMIN_DB_NAME and or PAYMENTS_DB_TABLE not set")
		os.Exit(1)
	}

	LOGGER.Infof("\t* CreateAdminServicePersistenceOperations dialing DB:%s collection:%s ", mo.dbName, mo.dbCollection)

	client, err := database.InitializeAtlasConnection(dbUser, dbPwd, mongoId)
	if err != nil {
		LOGGER.Errorf("Mongo Atlas DB connection failed! %s", err)
		panic("Mongo Atlas DB connection failed! " + err.Error())
	}
	mo.session = client

	LOGGER.Infof("\t* CreateAdminServicePersistenceOperations DB is set")

	return mo, nil
}

func (mo MongoDBOperations) GetCollection() (*mongo.Collection, context.Context) {
	dbTimeout, _ := strconv.Atoi(os.Getenv(environment.ENV_KEY_DB_TIMEOUT))
	LOGGER.Infof("\t* Getting collection: %s from DB %s", mo.dbCollection, mo.dbName)
	collection := mo.session.Database(mo.dbName).Collection(mo.dbCollection)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(dbTimeout))
	return collection, ctx
}

func (mo MongoDBOperations) StoreFiToFiCCTMemo(w http.ResponseWriter, request *http.Request) {

	var fitoficctMemoRQ model.FitoFICCTMemoData

	err := json.NewDecoder(request.Body).Decode(&fitoficctMemoRQ)

	if err != nil {
		LOGGER.Warningf("Unable to parse body of REST call to store FioTFiCCTMemo:  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "ADMIN-0002", err)
		return
	}

	LOGGER.Infof("Storing FiToFiCCTMemo in DB, instruction ID is %v", *fitoficctMemoRQ.Fitoficctnonpiidata.InstructionID)

	collection, ctx := mo.GetCollection()
	var results []model.FitoFICCTMemoData

	cursor, err := collection.Find(ctx,
		bson.M{
			"fitoficctnonpiidata.instruction_id": fitoficctMemoRQ.Fitoficctnonpiidata.InstructionID,
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

		_, err = collection.UpdateOne(
			ctx,
			bson.M{"fitoficctnonpiidata.instruction_id": fitoficctMemoRQ.Fitoficctnonpiidata.InstructionID},
			bson.M{
				"$set": bson.M{
					"fitoficctnonpiidata":    fitoficctMemoRQ.Fitoficctnonpiidata,
					"fitoficct_pii_hash":     fitoficctMemoRQ.FitoficctPiiHash,
					"transaction_identifier": fitoficctMemoRQ.TransactionIdentifier,
					"transaction_status":     fitoficctMemoRQ.TransactionStatus,
				},
			},
		)
		if err != nil {
			LOGGER.Warningf("Update of FITOFICCTMEMO was not successful  %v", err)
			response.NotifyWWError(w, request, http.StatusBadRequest, "ADMIN-0003", err)
		}
	} else {
		fitoficctMemoRQ.ID = primitive.NewObjectID().Hex()
		_, err = collection.InsertOne(ctx, fitoficctMemoRQ)
		if err != nil {
			LOGGER.Warningf("Insert of FITOFICCTMEMO was not successful  %v", err)
			response.NotifyWWError(w, request, http.StatusBadRequest, "ADMIN-0003", err)
		}

	}

}

func (mo MongoDBOperations) GetTxnDetails(w http.ResponseWriter, request *http.Request) {
	LOGGER.Infof("Administration-Service:persistance_handler:GetTxnDetails")
	var txnRequest model.FItoFITransactionRequest
	var fitoficctMemos []*model.FitoFICCTMemoData
	var txnResponse []model.FItoFITransaction
	var txn model.FItoFITransaction

	err := json.NewDecoder(request.Body).Decode(&txnRequest)

	if err != nil {
		LOGGER.Warningf("Unable to parse body of REST call to query transaction details:  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "ADMIN-0004", err)
		return
	}

	queryType := *txnRequest.QueryType
	queryData := txnRequest.QueryData
	startDate := txnRequest.StartDate
	endDate := txnRequest.EndDate
	homeDomain := *txnRequest.OfiID

	collection, ctx := mo.GetCollection()

	if strings.EqualFold(queryType, TRANSACTION_ID) {
		cursor, err := collection.Find(
			ctx,
			bson.M{
				"transaction_identifier": queryData,
				"$or": []bson.M{
					bson.M{"fitoficctnonpiidata.transactiondetails.rfi_id": homeDomain},
					bson.M{"fitoficctnonpiidata.transactiondetails.ofi_id": homeDomain},
				},
			})
		if err != nil {
			response.NotifyWWError(w, request, http.StatusNotFound, "ADMIN-0005", err)
			return
		}
		bytes, err := database.ParseResult(cursor, ctx)
		if err != nil {
			LOGGER.Debugf("Error parsing mongo data")
			response.NotifyWWError(w, request, http.StatusNotFound, "ADMIN-0022", err)
			return
		}
		_ = json.Unmarshal(bytes, &fitoficctMemos)

	} else if strings.EqualFold(queryType, INSTRUCTION_ID) {
		cursor, err := collection.Find(
			ctx,
			bson.M{
				"fitoficctnonpiidata.instruction_id": queryData,
				"$or": []bson.M{
					bson.M{"fitoficctnonpiidata.transactiondetails.rfi_id": homeDomain},
					bson.M{"fitoficctnonpiidata.transactiondetails.ofi_id": homeDomain},
				},
			})
		if err != nil {
			response.NotifyWWError(w, request, http.StatusNotFound, "ADMIN-0005", err)
			return
		}

		bytes, err := database.ParseResult(cursor, ctx)
		if err != nil {
			LOGGER.Debugf("Error parsing mongo data")
			response.NotifyWWError(w, request, http.StatusNotFound, "ADMIN-0022", err)
			return
		}
		_ = json.Unmarshal(bytes, &fitoficctMemos)

	} else if strings.EqualFold(queryType, DATE_RANGE) {
		layout := "2006-01-02"

		start, _ := time.Parse(layout, startDate.String())
		sd := start.Unix()

		end, _ := time.Parse(layout, endDate.String())
		ed := end.Unix()

		paginationSkip := int64((txnRequest.PageNumber - 1) * txnRequest.TransactionBatch)

		mongoOptions := &options.FindOptions{
			Skip:  &paginationSkip,
			Limit: &txnRequest.TransactionBatch,
		}

		cursor, err := collection.Find(
			ctx,
			bson.M{
				"time_stamp": bson.M{
					"$gt": sd,
					"$lt": ed,
				},
				"$or": []bson.M{
					bson.M{"fitoficctnonpiidata.transactiondetails.rfi_id": homeDomain},
					bson.M{"fitoficctnonpiidata.transactiondetails.ofi_id": homeDomain},
				},
			},
			mongoOptions,
		)
		if err != nil {
			response.NotifyWWError(w, request, http.StatusNotFound, "ADMIN-0005", err)
			return
		}
		bytes, err := database.ParseResult(cursor, ctx)

		if err != nil {
			LOGGER.Debugf("Error parsing mongo data")
			response.NotifyWWError(w, request, http.StatusNotFound, "ADMIN-0022", err)
			return
		}
		_ = json.Unmarshal(bytes, &fitoficctMemos)
	}

	LOGGER.Debugf("Number of Result found: %v", len(fitoficctMemos))

	for _, fitoficctMemo := range fitoficctMemos {
		if fitoficctMemo.Fitoficctnonpiidata == nil {
			continue
		}

		var ofiId, rfiId string
		if fitoficctMemo.Fitoficctnonpiidata.Transactiondetails != nil {
			ofiId = *fitoficctMemo.Fitoficctnonpiidata.Transactiondetails.OfiID
			rfiId = *fitoficctMemo.Fitoficctnonpiidata.Transactiondetails.RfiID
		}

		if ofiId != homeDomain && rfiId != homeDomain {
			//LOGGER.Infof("Transaction details request received from a third party, who did not participated in the transaction.")
			continue
		}

		txn.TransactionReceipt = fitoficctMemo.TransactionStatus
		txn.TransactionDetails = fitoficctMemo.Fitoficctnonpiidata.Transactiondetails

		txnResponse = append(txnResponse, txn)
	}

	txnResponseJSON, err := json.Marshal(txnResponse)
	if err != nil {
		response.NotifyWWError(w, request, http.StatusNotFound, "ADMIN-0006", err)
		return
	}

	response.Respond(w, http.StatusOK, txnResponseJSON)
}
