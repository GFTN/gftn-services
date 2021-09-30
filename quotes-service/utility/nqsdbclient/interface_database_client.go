// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package nqsdbclient

import "github.com/shopspring/decimal"

type DatabaseClient interface {
	CreateConnection() error
	CloseConnection()
	CreateRequest(requestID string, ofiID string, LimitMaxOfi decimal.Decimal, LimitMinOfi decimal.Decimal, sourceAsset []byte, targetAsset []byte, timeRequest int64, timeExpireOfi int64) error
	GetRequest(requestID string, ofiID string) ([]RequestDB, error)
	CreateQuote(requestID string, quoteID string, rfiID string, ofiID string, maxLimit decimal.Decimal, minLimit decimal.Decimal, sourceAsset []byte, targetAsset []byte, timeOfRequest int64, quoteStatus int, timeExpireOfi int64) error
	GetQuotes(requestID string, ofiID string) ([]QuoteDB, error)
	GetQuotesByAttributes(query *Query) ([]QuoteDB, error)
	CancelQuotesByAttributes(query *Query, TimeCancel int64) ([]QuoteDB, error)
	UpdateQuote(quote QuoteDB, timeOfQuote int64) error
	CancelQuote(quoteID string, rfiID string, TimeCancel int64) error
	ExecutingQuote(quoteID string, ofiID string, quoteResponse []byte, executingTime int64, amount decimal.Decimal) error
	ExecutedQuote(quoteID string, ofiID string, timeExecute int64) error
	FailedQuote(quoteID string, ofiID string) error
	GetQuoteByQuoteID(QuoteID string, rfiID string) ([]QuoteDB, error)
}
