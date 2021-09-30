// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package fitoficct

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/GFTN/gftn-services/api-service/environment"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/GFTN/gftn-services/api-service/client"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/utility"
	comn "github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/response"
)

type FItoFICustomerCreditTransferOperations struct {
	AdministrationServiceClient client.AdministrationServiceClient
}

// end to end id and stellar transaction id
const (
	INSTRUCTION_ID = "INSTRUCTION_ID"
	TRANSACTION_ID = "TRANSACTION_ID"
	DATE_RANGE     = "DATE_RANGE"
)

func CreateFItoFICustomerCreditTransferOperation() (FItoFICustomerCreditTransferOperations, error) {

	op := FItoFICustomerCreditTransferOperations{}
	administrationServiceClient, err := client.CreateRestAdministrationServiceClient()
	utility.ExitOnErr(LOGGER, err, "Unable to create internal compliance service client")

	op.AdministrationServiceClient = administrationServiceClient

	return op, nil

}

func isOperatingAccount(da string) bool {
	if da == "" || da == comn.ISSUING {
		return false
	}
	return true
}

// ProcessTxnDetails -
/*
 * This function gets query_type and query_data as parameters and communicate to
 * Administration_Service to get the transaction details for the ID received as
 * query_data.
 */
func (op FItoFICustomerCreditTransferOperations) ProcessTxnDetails(w http.ResponseWriter, request *http.Request) {
	LOGGER.Infof("API-Service:cash_clearing:ProcessTxnDetails")
	var txnDetailRequest model.FItoFITransactionRequest
	homeDomain := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)

	queryParams := request.URL.Query()
	queryType := strings.ToUpper(strings.TrimSpace(strings.Join(queryParams["query_type"], " ")))
	queryData := strings.TrimSpace(strings.Join(queryParams["query_data"], " "))
	startDate := strings.TrimSpace(strings.Join(queryParams["start_date"], " "))
	endDate := strings.TrimSpace(strings.Join(queryParams["end_date"], " "))

	LOGGER.Infof("queryType: %s", queryType)
	txnDetailRequest.QueryType = &queryType
	txnDetailRequest.OfiID = &homeDomain

	if queryType == INSTRUCTION_ID || queryType == TRANSACTION_ID {
		LOGGER.Infof("queryData: %s", queryData)
		if strings.TrimSpace(queryType) == "" || strings.TrimSpace(queryData) == "" {
			LOGGER.Error("queryType & queryData cannot be empty")
			response.NotifyWWError(w, request, http.StatusBadRequest, "API-1100", errors.New("queryType & queryData cannot be empty"))
			return
		}
		txnDetailRequest.QueryData = queryData
	} else if queryType == DATE_RANGE {
		LOGGER.Infof("startDate: %s", startDate)
		LOGGER.Infof("endDate: %s", endDate)
		if strings.TrimSpace(startDate) == "" || strings.TrimSpace(endDate) == "" {
			LOGGER.Error("startDate & endDate cannot be empty")
			response.NotifyWWError(w, request, http.StatusBadRequest, "API-1100", errors.New("startDate & endDate cannot be empty"))
			return
		}
		layout := "2006-01-02"
		date, _ := time.Parse(layout, startDate)
		txnDetailRequest.StartDate = strfmt.Date(date)
		date, _ = time.Parse(layout, endDate)
		txnDetailRequest.EndDate = strfmt.Date(date)

		transactionBatch := strings.TrimSpace(strings.Join(queryParams["batch"], " "))
		pageNumber := strings.TrimSpace(strings.Join(queryParams["page"], " "))
		if transactionBatch == "" || pageNumber == "" {
			LOGGER.Error("page & batch cannot be empty")
			response.NotifyWWError(w, request, http.StatusBadRequest, "API-1100", errors.New("page & batch cannot be empty"))
			return
		}
		batchLimit, _ := strconv.ParseInt(os.Getenv(environment.ENV_KEY_TRANSACTION_BATCH_LIMIT), 10, 64)
		txnDetailRequest.TransactionBatch, _ = strconv.ParseInt(transactionBatch, 10, 64)
		if txnDetailRequest.TransactionBatch > batchLimit {
			LOGGER.Error("transaction_batch has exceeded the batch limit of Worldwire")
			response.NotifyWWError(w, request, http.StatusBadRequest, "API-1310", errors.New("transaction_batch has exceeded the batch limit of Worldwire"))
			return
		}
		txnDetailRequest.PageNumber, _ = strconv.ParseInt(pageNumber, 10, 64)

	} else {
		response.NotifyWWError(w, request, http.StatusBadRequest, "API-1059", errors.New("Unknown Query Type"))
		return
	}

	txnResponseJSON, status, err := op.AdministrationServiceClient.GetTxnDetails(txnDetailRequest)
	if err != nil {
		switch status {
		case http.StatusNotFound:
			response.NotifyWWError(w, request, status, "API-1233", err)
			return
		case http.StatusBadRequest:
			response.NotifyWWError(w, request, status, "API-1234", err)
			return
		case http.StatusInternalServerError:
			response.NotifyWWError(w, request, status, "API-1235", err)
			return
		default:
			LOGGER.Infof("Transaction query success")
		}
	}

	response.Respond(w, http.StatusOK, txnResponseJSON)

}
