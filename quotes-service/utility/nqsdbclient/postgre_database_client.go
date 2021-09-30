// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package nqsdbclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/shopspring/decimal"

	_ "github.com/lib/pq"
)

var (
	QUOTESTATUS_PENDING   = 1
	QUOTESTATUS_UPDATED   = 2
	QUOTESTATUS_EXECUTING = 3
	QUOTESTATUS_EXECUTED  = 4
	QUOTESTATUS_FAILED    = 98
	QUOTESTATUS_CANCELLED = 99
)

type QuoteDB struct {
	RequestID         *string          `json:"request_id,omitempty" db:"requestid"`
	QuoteID           *string          `json:"quote_id,omitempty" db:"quoteid"`
	RfiId             *string          `json:"rfi_id,omitempty" db:"rfiid"`
	OfiId             *string          `json:"ofi_id,omitempty" db:"ofiid"`
	LimitMaxOfi       *decimal.Decimal `json:"limit_max_ofi,omitempty" db:"limitmaxofi"`
	LimitMinOfi       *decimal.Decimal `json:"limit_min_ofi,omitempty" db:"limitminofi"`
	LimitMaxRfi       *decimal.Decimal `json:"limit_max_rfi,omitempty" db:"limitmaxrfi"`
	LimitMinRfi       *decimal.Decimal `json:"limit_min_rfi,omitempty" db:"limitminrfi"`
	Amount            *decimal.Decimal `json:"amount,omitempty" db:"amount"`
	ExchangeRate      *decimal.Decimal `json:"exchange_rate,omitempty" db:"exchangerate"`
	SourceAsset       *types.JSONText  `json:"source_asset,omitempty" db:"sourceasset"`
	TargetAsset       *types.JSONText  `json:"target_asset,omitempty" db:"targetasset"`
	TimeRequest       *int64           `json:"time_request,omitempty" db:"timerequest"`
	TimeQuote         *int64           `json:"time_quote,omitempty" db:"timequote"`
	TimeExpireOfi     *int64           `json:"time_expire_ofi,omitempty" db:"timeexpireofi"`
	TimeStartRfi      *int64           `json:"time_start_rfi,omitempty" db:"timestartrfi"`
	TimeExpireRfi     *int64           `json:"time_expire_rfi,omitempty" db:"timeexpirerfi"`
	StatusQuote       *int             `json:"status_quote,omitempty" db:"statusquote"`
	TimeExecuted      *int64           `json:"time_executed,omitempty" db:"timeexecuted"`
	TimeCancel        *int64           `json:"time_cancel,omitempty" db:"timecancel"`
	AddressReceiveRfi *string          `json:"address_receive_rfi,omitempty" db:"addressreceiverfi"`
	AddressSendRfi    *string          `json:"address_send_rfi,omitempty" db:"addresssendrfi"`
	// QuoteRequest           *types.JSONText `json:"quote_request,omitempty" db:"quoterequest"`
	QuoteRequestSignature  *string         `json:"quote_request_signature,omitempty" db:"quoterequestsignature"`
	QuoteResponse          *types.JSONText `json:"quote_response,omitempty" db:"quoteresponse"`
	QuoteResponseBase64    *string         `json:"quote_response_base64,omitempty" db:"quoteresponsebase64"`
	QuoteResponseSignature *string         `json:"quote_response_signature,omitempty" db:"quoteresponsesignature"`
}

type RequestDB struct {
	RequestID     *string          `json:"request_id,omitempty" db:"requestid"`
	TimeExpireOfi *int64           `json:"time_expire_ofi,omitempty" db:"timeexpireofi"`
	LimitMaxOfi   *decimal.Decimal `json:"limit_max_rfi,omitempty" db:"limitmaxrfi"`
	LimitMinOfi   *decimal.Decimal `json:"limit_min_rfi,omitempty" db:"limitminrfi"`
	SourceAsset   *types.JSONText  `json:"source_asset,omitempty" db:"sourceasset"`
	TargetAsset   *types.JSONText  `json:"target_asset,omitempty" db:"targetasset"`
	TimeRequest   *int64           `json:"time_request,omitempty" db:"timerequest"`
	OfiId         *string          `json:"ofi_id,omitempty" db:"ofiid"`
}

type Query struct {
	DeleteAllQuotes   *bool            `json:"delete_all_quotes,omitempty"`
	RequestID         *string          `json:"request_id,omitempty" db:"requestid"`
	QuoteID           *string          `json:"quote_id,omitempty"`
	RfiID             *string          `json:"rfi_id,omitempty"`
	OfiID             *string          `json:"ofi_id,omitempty"`
	LimitMaxOfi       *decimal.Decimal `json:"limit_max_ofi,omitempty"`
	LimitMinOfi       *decimal.Decimal `json:"limit_min_ofi,omitempty"`
	ExchangeRate      *Comparison      `json:"exchange_rate,omitempty"`
	SourceAsset       *types.JSONText  `json:"source_asset,omitempty"`
	TargetAsset       *types.JSONText  `json:"target_asset,omitempty"`
	TimeRequest       *int64           `json:"time_request,omitempty"`
	TimeQuote         *int64           `json:"time_quote,omitempty"`
	TimeExpireRfi     *Comparison      `json:"time_expire_rfi,omitempty"`
	StatusQuote       *Comparison      `json:"status_quote,omitempty"`
	TimeExecuted      *int64           `json:"time_executed,omitempty"`
	TimeCancel        *int64           `json:"time_cancel,omitempty"`
	AddressReceiveRfi *string          `json:"address_receive_rfi,omitempty"`
	AddressSendRfi    *string          `json:"address_send_rfi,omitempty"`
}
type Comparison struct {
	Threshold *interface{} `json:"threshold"`
	Operator  *string      `json:"operator"` //comparison operator
}
type PostgreDatabaseClient struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	db       *sqlx.DB
}

//CreateConnection opens DB connection
func (dbc *PostgreDatabaseClient) CreateConnection() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=require",
		dbc.Host, dbc.Port, dbc.User, dbc.Password, dbc.Dbname)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	LOGGER.Info("Successfully connected!")
	dbc.db = db
	return nil
}

//CloseConnection closes DB connection
func (dbc *PostgreDatabaseClient) CloseConnection() {
	dbc.db.Close()
}

// Create quote request to DB
func (dbc *PostgreDatabaseClient) CreateRequest(requestID string, ofiID string, LimitMaxOfi decimal.Decimal, LimitMinOfi decimal.Decimal, sourceAsset []byte, targetAsset []byte, timeRequest int64, timeExpireOfi int64) error {
	sqlStatement := `
		INSERT INTO REQUESTS ( requestID, ofiID, LimitMaxOfi, LimitMinOfi, sourceAsset, targetAsset, timeRequest, timeExpireOfi )
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING requestID
		`
	id := ""
	err := dbc.db.QueryRow(sqlStatement,
		requestID,
		ofiID,
		LimitMaxOfi,
		LimitMinOfi,
		sourceAsset,
		targetAsset,
		timeRequest,
		timeExpireOfi,
	).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

// Get Quote Request
func (dbc PostgreDatabaseClient) GetRequest(requestID string, ofiID string) ([]RequestDB, error) {
	sqlStatement := `SELECT * FROM requests WHERE requestID=$1 AND ofiID=$2;`
	rows, err := dbc.db.Queryx(sqlStatement, requestID, ofiID)
	var requests []RequestDB
	if err != nil {
		return requests, err
	}
	for rows.Next() {
		request := RequestDB{}
		err = rows.StructScan(&request)
		if err != nil {
			return requests, err
		}
		requests = append(requests, request)
	}
	return requests, nil
}

// CreateQuote creates quote to DB
func (dbc *PostgreDatabaseClient) CreateQuote(requestID string, quoteID string, rfiID string, ofiID string, LimitMaxOfi decimal.Decimal, LimitMinOfi decimal.Decimal, sourceAsset []byte, targetAsset []byte, timeRequest int64, statusQuote int, timeExpireOfi int64) error {
	sqlStatement := `
		INSERT INTO QUOTES ( requestID, quoteID, rfiID, ofiID, LimitMaxOfi, LimitMinOfi, sourceAsset, targetAsset, timeRequest, statusQuote, timeExpireOfi )
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING requestID
		`
	id := ""
	err := dbc.db.QueryRow(sqlStatement,
		requestID,
		quoteID,
		rfiID,
		ofiID,
		LimitMaxOfi,
		LimitMinOfi,
		sourceAsset,
		targetAsset,
		timeRequest,
		statusQuote,
		timeExpireOfi,
	).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

// Get Quote Response
func (dbc PostgreDatabaseClient) GetQuotes(requestID string, ofiID string) ([]QuoteDB, error) {
	var quotesResponse []QuoteDB
	sqlStatement := `SELECT * FROM quotes WHERE requestID=$1 AND ofiID=$2;`
	rows, err := dbc.db.Queryx(sqlStatement, requestID, ofiID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		quote := QuoteDB{}
		err = rows.StructScan(&quote)
		if err != nil {
			return nil, err
		}
		quotesResponse = append(quotesResponse, quote)
	}
	// check if OFI has actually made request with requestID
	//if no quote returned, check request table.
	if len(quotesResponse) == 0 {
		requests, err := dbc.GetRequest(requestID, ofiID)
		if err != nil {
			LOGGER.Error(err)
			msg := "Error fetching request from DB, requestID: " + requestID
			LOGGER.Error(msg)
			return nil, err
		}
		if len(requests) == 0 {
			return nil, errors.New(fmt.Sprintf("requestID %v does not exist", requestID))
		} else {
			return quotesResponse, nil
		}
	}
	return quotesResponse, nil
}

//called by ofi
func (dbc PostgreDatabaseClient) GetQuotesByAttributes(query *Query) ([]QuoteDB, error) {
	var m map[string]interface{}
	queryJSON, _ := json.Marshal(*query)
	json.Unmarshal(queryJSON, &m)
	var values []interface{}
	var where []string
	i := 1 // for $num in postgresql varbinding
	comOperator := "="
	for _, k := range []string{"request_id", "quote_id", "rfi_id", "ofi_id",
		"limit_max_rfi", "limit_min_rfi", "exchange_rate",
		"source_asset", "target_asset", "time_expire_rfi", "status_quote"} {
		if v, ok := m[k]; ok {
			if k == "time_expire_rfi" ||
				k == "exchange_rate" ||
				k == "status_quote" { // operator
				vJSON, _ := json.Marshal(v)
				var tempMap map[string]interface{}
				json.Unmarshal(vJSON, &tempMap)
				tempOperator := tempMap["operator"].(string)
				comOperator = getOperator(tempOperator)
				values = append(values, tempMap["threshold"])
				where = append(where, fmt.Sprintf("%s"+comOperator+"$"+strconv.Itoa(i), strings.Replace(k, "_", "", -1)))
				i++
				continue
			} else {
				comOperator = "="
			}
			if k == "source_asset" || k == "target_asset" { // mashal []interface{} into json blob
				vJSON, err := json.Marshal(v)
				if err != nil {
					return nil, err
				}
				values = append(values, vJSON)
				where = append(where, fmt.Sprintf("%s"+comOperator+"$"+strconv.Itoa(i), strings.Replace(k, "_", "", -1)))
				i++
				continue
			} else {
				values = append(values, v)
				where = append(where, fmt.Sprintf("%s"+comOperator+"$"+strconv.Itoa(i), strings.Replace(k, "_", "", -1)))
				i++
			}
		}
	}
	LOGGER.Debug("Get quote by attributes statement: " + "SELECT * FROM quotes WHERE " + strings.Join(where, " AND "))
	rows, err := dbc.db.Queryx("SELECT * FROM quotes WHERE "+strings.Join(where, " AND "), values...)
	// rows, err := dbc.db.Queryx("SELECT * FROM quotes WHERE timeexpirerfi>$1", 1232)
	if err != nil {
		return nil, err
	}
	var quotesResponse []QuoteDB
	for rows.Next() {
		quote := QuoteDB{}
		err = rows.StructScan(&quote)
		quotesResponse = append(quotesResponse, quote)
	}
	return quotesResponse, nil
}

// statusQuote : 1 = pending, 2 = updated, 3 = executing , 4 = executed, 98=failed, 99 = canceled
func (dbc PostgreDatabaseClient) UpdateQuote(quote QuoteDB, timequote int64) error {

	sqlStatement := "SELECT updateQuote($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26)"

	result, err := dbc.db.Exec(
		sqlStatement,
		quote.QuoteID,
		quote.RfiId,
		QUOTESTATUS_UPDATED,
		quote.ExchangeRate,
		timequote,
		quote.TimeStartRfi,
		quote.TimeExpireRfi,
		quote.AddressReceiveRfi,
		quote.AddressSendRfi,
		quote.QuoteResponse,
		quote.QuoteResponseBase64,
		quote.QuoteResponseSignature,
		quote.LimitMaxRfi,
		quote.LimitMinRfi,
		quote.SourceAsset,
		quote.TargetAsset,
		quote.LimitMaxOfi,
		quote.LimitMinOfi,
		quote.TimeExpireOfi,
		quote.OfiId,
		QUOTESTATUS_EXECUTED,
		QUOTESTATUS_CANCELLED,
		QUOTESTATUS_EXECUTING,
		QUOTESTATUS_PENDING,
		QUOTESTATUS_UPDATED,
		QUOTESTATUS_FAILED,
	)
	if err != nil {
		return err
	}
	rowAffeced, err := result.RowsAffected()
	if rowAffeced > 1 {
		return errors.New("Updated more than one row. Expecting one or none")
	}
	return nil
}

//TODO time
func (dbc PostgreDatabaseClient) CancelQuote(quoteID string, rfiID string, timeCancel int64) error {

	sqlStatement := "SELECT cancelQuote($1,$2,$3,$4,$5,$6,$7,$8,$9)"

	result, err := dbc.db.Exec(
		sqlStatement,
		quoteID,
		rfiID,
		QUOTESTATUS_CANCELLED,
		timeCancel,
		QUOTESTATUS_UPDATED,
		QUOTESTATUS_FAILED,
		QUOTESTATUS_PENDING,
		QUOTESTATUS_EXECUTED,
		QUOTESTATUS_EXECUTING,
	)
	if err != nil {
		return err
	}
	rowAffeced, err := result.RowsAffected()
	if rowAffeced > 1 {
		return errors.New("Updated more than one quote. Expecting one and only one")
	}
	if rowAffeced < 1 {
		return errors.New("Updated no quote.")
	}
	return nil
}

//called by rfi
func (dbc PostgreDatabaseClient) CancelQuotesByAttributes(query *Query, timeCancel int64) ([]QuoteDB, error) {
	var m map[string]interface{}
	queryJSON, _ := json.Marshal(*query)
	json.Unmarshal(queryJSON, &m)
	var values []interface{}
	var where []string
	// set the first values in binding
	values = append(values, timeCancel)

	i := 2 // for $num in postgresql varbinding
	// check if delete_all_quotes flag exist
	if query.DeleteAllQuotes == (*bool)(nil) {
		if len(m) < 2 {
			return nil, errors.New("delete_all_flag=null with no other query fields")
		}
	} else if *query.DeleteAllQuotes == false {
		//ensure not all other fields are empty
		if len(m) < 3 {
			return nil, errors.New("delete_all_flag=false with no other query fields")
		}
	} else if *query.DeleteAllQuotes == true {
		//overwrite other fields
		temp := make(map[string]interface{})
		temp["rfi_id"] = m["rfi_id"]
		m = temp
	}
	comOperator := "="
	for _, k := range []string{"quote_id", "rfi_id", "ofi_id",
		"limit_max_ofi", "limit_min_ofi", "exchange_rate",
		"source_asset", "target_asset", "time_expire_rfi", "status_quote"} {

		if v, ok := m[k]; ok {
			if k == "time_expire_rfi" ||
				k == "exchange_rate" ||
				k == "status_quote" { // operator
				vJSON, _ := json.Marshal(v)
				var tempMap map[string]interface{}
				json.Unmarshal(vJSON, &tempMap)
				tempOperator := tempMap["operator"].(string)
				comOperator = getOperator(tempOperator)
				values = append(values, tempMap["threshold"])
				where = append(where, fmt.Sprintf("%s"+comOperator+"$"+strconv.Itoa(i), strings.Replace(k, "_", "", -1)))
				i++
				continue
			} else {
				comOperator = "="
			}
			if k == "source_asset" || k == "target_asset" { // mashal []interface{} into json blob
				vJSON, err := json.Marshal(v)
				if err != nil {
					return nil, err
				}
				values = append(values, vJSON)
				where = append(where, fmt.Sprintf("%s"+comOperator+"$"+strconv.Itoa(i), strings.Replace(k, "_", "", -1)))
				i++
				continue
			} else {
				values = append(values, v)
				where = append(where, fmt.Sprintf("%s"+comOperator+"$"+strconv.Itoa(i), strings.Replace(k, "_", "", -1)))
				i++
			}
		}
	}
	rows, err := dbc.db.Queryx(" UPDATE quotes SET statusQuote=99, timeCancel=$1 WHERE "+strings.Join(where, " AND ")+"  AND (statusQuote <=2 OR statusQuote = 98) returning *", values...)

	if err != nil {
		return nil, err
	}
	var quotesResponse []QuoteDB
	for rows.Next() {
		quote := QuoteDB{}
		err = rows.StructScan(&quote)
		quotesResponse = append(quotesResponse, quote)
	}
	return quotesResponse, nil
}

//TODO time
func (dbc PostgreDatabaseClient) ExecutingQuote(quoteID string, ofiID string, quoteResponse []byte, timeExecuting int64, amount decimal.Decimal) error {

	sqlStatement := "SELECT ExecutingQuote($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)"

	result, err := dbc.db.Exec(sqlStatement,
		quoteID,
		ofiID,
		quoteResponse,
		QUOTESTATUS_EXECUTED,
		QUOTESTATUS_CANCELLED,
		QUOTESTATUS_EXECUTING,
		QUOTESTATUS_PENDING,
		QUOTESTATUS_UPDATED,
		QUOTESTATUS_FAILED,
		timeExecuting,
		amount,
	)
	if err != nil {
		return err
	}
	rowAffeced, err := result.RowsAffected()
	if rowAffeced > 1 {
		return errors.New("Execute more than one quote. Expecting one and only one")
	}
	if rowAffeced < 1 {
		return errors.New("Quote Not Found.")
	}
	return nil
}

func (dbc PostgreDatabaseClient) ExecutedQuote(quoteID string, rfiID string, timeExecuted int64) error {
	sqlStatement := `
	UPDATE quotes
	SET statusQuote = $3, timeExecuted = $4
	WHERE quoteID = $1 AND rfiID = $2 AND statusQuote=3;`
	result, err := dbc.db.Exec(sqlStatement, quoteID, rfiID, QUOTESTATUS_EXECUTED, timeExecuted)
	if err != nil {
		return err
	}
	rowAffeced, err := result.RowsAffected()
	if rowAffeced > 1 {
		return errors.New("Updated more than one quote. Expecting one and only one")
	}
	if rowAffeced < 1 {
		return errors.New("Valid Quote Not Found.")
	}

	return nil
}

func (dbc PostgreDatabaseClient) FailedQuote(quoteID string, rfiID string) error {
	sqlStatement := `
	UPDATE quotes
	SET statusQuote = $3
	WHERE quoteID = $1 AND rfiID = $2;`
	result, err := dbc.db.Exec(sqlStatement, quoteID, rfiID, QUOTESTATUS_FAILED)
	if err != nil {
		return err
	}
	rowAffeced, err := result.RowsAffected()
	if rowAffeced > 1 {
		return errors.New("Updated more than one quote. Expecting one and only one")
	}
	if rowAffeced < 1 {
		return errors.New("Valid Quote Not Found.")
	}
	return nil
}

func (dbc PostgreDatabaseClient) GetQuoteByQuoteID(QuoteID string, rfiID string) ([]QuoteDB, error) {
	var quotesDB []QuoteDB
	sqlStatement := `SELECT * FROM quotes WHERE quoteID=$1 AND rfiID=$2;`
	rows, err := dbc.db.Queryx(sqlStatement, QuoteID, rfiID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		quote := QuoteDB{}
		err = rows.StructScan(&quote)
		if err != nil {
			return nil, err
		}
		quotesDB = append(quotesDB, quote)
	}
	if len(quotesDB) == 0 {
		return quotesDB, errors.New("Quote Not Found")
	}
	return quotesDB, nil

}

// for testing purpose only
func (dbc PostgreDatabaseClient) InsertQuote(quote QuoteDB) error {
	sqlStatement := `
		INSERT INTO QUOTES ( requestID, quoteID, rfiID, ofiID, LimitMaxOfi, LimitMinOfi, sourceAsset, targetAsset, timeRequest, statusQuote, timeExpireOfi, timeExpireRfi, exchangeRate, QuoteResponse, QuoteResponseBase64, QuoteResponseSignature, Amount)
		VALUES (:requestid, :quoteid, :rfiid, :ofiid, :limitmaxofi, :limitminofi, :sourceasset, :targetasset, :timerequest, :statusquote, :timeexpireofi, :timeexpirerfi, :exchangerate, :quoteresponse, :quoteresponsebase64, :quoteresponsesignature, :amount)
		`
	_, err := dbc.db.NamedExec(sqlStatement, quote)
	if err != nil {
		return err
	}
	return nil
}

func getOperator(opStr string) string {
	switch opStr {
	case "eq":
		return "="
	case "gt":
		return ">"
	case "lt":
		return "<"
	case "ge":
		return ">="
	case "le":
		return "<="
	default:
		return "="
	}
}
