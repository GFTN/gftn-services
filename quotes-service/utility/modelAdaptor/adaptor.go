// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package modeladaptor

import (
	"encoding/base64"
	"encoding/json"

	"github.com/go-openapi/strfmt"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/quotes-service/utility/nqsdbclient"
	"github.com/GFTN/gftn-services/quotes-service/utility/nqsmodel"
)

func ExchangeRequestEnvelopeToNqs(exchangeRequestEnv model.ExchangeEnvelope) (nqsmodel.NqsExchangeRequest, error) {
	// expTime, _ := time.Parse(time.RFC3339, exchangeRequest.Quote.TimeExpire.String())
	// expTimeUnix := expTime.Unix()
	nqsExchangeRequest := nqsmodel.NqsExchangeRequest{}
	nqsExchangeRequest.ExchangeRequestBase64 = exchangeRequestEnv.Exchange
	nqsExchangeRequest.ExchangeRequestSignBase64 = exchangeRequestEnv.Signature
	exchangeRequest := model.Exchange{}
	exchangeRequestByte, _ := base64.StdEncoding.DecodeString(*exchangeRequestEnv.Exchange)
	nqsExchangeRequest.ExchangeRequestDecode = &exchangeRequest
	err := json.Unmarshal(exchangeRequestByte, &exchangeRequest)
	if err != nil {
		return nqsExchangeRequest, err
	}
	// validate exchange model
	err = exchangeRequest.Validate(strfmt.Default)
	if err != nil {
		msg := "Error validating structure of exchange request" + err.Error()
		LOGGER.Warningf(msg)
		return nqsExchangeRequest, err
	}
	// TODO: change requestID in quoteID
	nqsExchangeRequest.QuoteID = exchangeRequest.Quote.QuoteID
	nqsExchangeRequest.RfiId = exchangeRequest.Quote.RfiID
	nqsExchangeRequest.OfiId = exchangeRequest.Quote.QuoteRequest.OfiID
	// nqsExchangeRequest.LimitMaxOfi * float64
	// nqsExchangeRequest.LimitMinOfi * float64
	// nqsExchangeRequest.LimitMaxRfi * float64
	// nqsExchangeRequest.LimitMinRfi * float64
	exchangeRate := *exchangeRequest.Quote.ExchangeRate
	nqsExchangeRequest.ExchangeRate = &exchangeRate
	amount := *exchangeRequest.Amount
	nqsExchangeRequest.Amount = &amount
	nqsExchangeRequest.SourceAsset = copyAsset(exchangeRequest.Quote.QuoteRequest.SourceAsset)
	nqsExchangeRequest.TargetAsset = copyAsset(exchangeRequest.Quote.QuoteRequest.TargetAsset)
	// nqsExchangeRequest.TimeRequest * int64
	// nqsExchangeRequest.TimeQuote * int64
	// nqsExchangeRequest.TimeExpireOfi * int64

	// tempTime, _ := time.Parse(time.RFC3339, exchangeRequest.Quote.TimeStart.String())
	// tempTimeUnix := tempTime.Unix()
	// nqsExchangeRequest.TimeStartRfi = &tempTimeUnix
	nqsExchangeRequest.TimeStartRfi = exchangeRequest.Quote.TimeStart
	// tempTime, _ = time.Parse(time.RFC3339, exchangeRequest.Quote.TimeExpire.String())
	// tempTimeUnix = tempTime.Unix()
	// nqsExchangeRequest.TimeExpireRfi = &tempTimeUnix
	nqsExchangeRequest.TimeExpireRfi = exchangeRequest.Quote.TimeExpire

	// nqsExchangeRequest.StatusQuote * int
	// nqsExchangeRequest.TimeExecuted * int64
	// nqsExchangeRequest.TimeCancel * int64
	// nqsExchangeRequest.AddressReceiveRfi = exchangeRequest.Quote.AddressReceive
	// nqsExchangeRequest.AddressSendRfi = exchangeRequest.Quote.AddressSend
	// nqsExchangeRequest.AddressReceiveOfi * string
	// nqsExchangeRequest.AddressSendOfi * string
	nqsExchangeRequest.AccountSendOfi = exchangeRequest.AccountNameSend
	nqsExchangeRequest.AccountReceiveOfi = exchangeRequest.AccountNameReceive
	nqsExchangeRequest.AccountSendRfi = exchangeRequest.Quote.AccountNameSend
	nqsExchangeRequest.AccountReceiveRfi = exchangeRequest.Quote.AccountNameReceive
	return nqsExchangeRequest, nil

}

func QuoteRequestToNqs(quoteRequest *model.QuoteRequest) nqsmodel.NqsQuoteRequest {
	nqsQuoteRequest := nqsmodel.NqsQuoteRequest{}
	nqsQuoteRequest.OfiId = quoteRequest.OfiID
	nqsQuoteRequest.SourceAsset = copyAsset(quoteRequest.SourceAsset)
	nqsQuoteRequest.TargetAsset = copyAsset(quoteRequest.TargetAsset)
	// tempTime, _ := time.Parse(time.RFC3339, quoteRequest.TimeExpire.String())
	// tempTimeUnix := tempTime.Unix()
	nqsQuoteRequest.TimeExpireOfi = quoteRequest.TimeExpire
	nqsQuoteRequest.LimitMaxOfi = quoteRequest.LimitMax
	nqsQuoteRequest.LimitMinOfi = quoteRequest.LimitMin
	return nqsQuoteRequest
}

func QuoteResponseEnvelopeToNqs(quoteResponseE *model.QuoteEnvelope) (nqsmodel.NqsAssetPriceQuote, error) {

	nqsAssetPriceQuote := nqsmodel.NqsAssetPriceQuote{}
	nqsAssetPriceQuote.QuoteResponseBase64 = quoteResponseE.Quote
	nqsAssetPriceQuote.QuoteResponseSignature = quoteResponseE.Signature
	//decode quoteResponse
	quoteResponse := model.Quote{}
	quoteResponseByte, _ := base64.StdEncoding.DecodeString(*quoteResponseE.Quote)
	err := json.Unmarshal(quoteResponseByte, &quoteResponse)
	if err != nil {
		return nqsAssetPriceQuote, err
	}
	// validate quote model:
	err = quoteResponse.Validate(strfmt.Default)
	if err != nil {
		msg := "Error validating structure of quote response" + err.Error()
		LOGGER.Warningf(msg)
		return nqsAssetPriceQuote, err
	}
	nqsAssetPriceQuote.TimeExpireRfi = quoteResponse.TimeExpire
	nqsAssetPriceQuote.LimitMaxOfi = quoteResponse.QuoteRequest.LimitMax
	nqsAssetPriceQuote.LimitMinOfi = quoteResponse.QuoteRequest.LimitMin
	nqsAssetPriceQuote.LimitMaxRfi = quoteResponse.LimitMax
	nqsAssetPriceQuote.LimitMinRfi = quoteResponse.LimitMin
	nqsAssetPriceQuote.AddressReceiveRfi = quoteResponse.AccountNameReceive
	nqsAssetPriceQuote.AddressSendRfi = quoteResponse.AccountNameSend
	nqsAssetPriceQuote.OfiId = quoteResponse.QuoteRequest.OfiID
	nqsAssetPriceQuote.ExchangeRate = quoteResponse.ExchangeRate
	nqsAssetPriceQuote.QuoteID = quoteResponse.QuoteID
	nqsAssetPriceQuote.RfiId = quoteResponse.RfiID
	nqsAssetPriceQuote.SourceAsset = copyAsset(quoteResponse.QuoteRequest.SourceAsset)
	nqsAssetPriceQuote.TargetAsset = copyAsset(quoteResponse.QuoteRequest.TargetAsset)
	nqsAssetPriceQuote.TimeStartRfi = quoteResponse.TimeStart
	nqsAssetPriceQuote.QuoteResponse = &quoteResponse
	return nqsAssetPriceQuote, nil
}

func copyAsset(asset *model.Asset) *nqsmodel.Asset {
	nqsAsset := nqsmodel.Asset{}
	nqsAsset.AssetCode = asset.AssetCode
	nqsAsset.AssetType = asset.AssetType
	nqsAsset.IssuerID = &asset.IssuerID
	return &nqsAsset
}

func QueryToQueryDB(apiQuery *model.QuoteFilter) nqsdbclient.Query {
	nqsQuery := nqsdbclient.Query{}
	apiQueryByte, _ := json.Marshal(apiQuery)
	temp := make(map[string]interface{})
	json.Unmarshal(apiQueryByte, &temp)
	// Replace the map key
	temp["status_quote"] = temp["status"]
	delete(temp, "status")
	tempByte, _ := json.Marshal(temp)
	json.Unmarshal(tempByte, &nqsQuery)
	return nqsQuery
}

func QuoteDBToQuoteStatus(quoteDB nqsdbclient.QuoteDB) model.QuoteStatus {
	quoteStatus := model.QuoteStatus{}

	quoteDBByte, _ := json.Marshal(quoteDB)
	temp := make(map[string]interface{})
	json.Unmarshal(quoteDBByte, &temp)
	// Replace the map key
	temp["status"] = temp["status_quote"]
	delete(temp, "status_quote")
	tempByte, _ := json.Marshal(temp)
	json.Unmarshal(tempByte, &quoteStatus)
	return quoteStatus
}

func NqsAssetPriceQuoteToQuoteDB(nqsQuote nqsmodel.NqsAssetPriceQuote) nqsdbclient.QuoteDB {
	quoteDB := nqsdbclient.QuoteDB{}
	nqsQuoteByte, _ := json.Marshal(nqsQuote)
	json.Unmarshal(nqsQuoteByte, &quoteDB)
	return quoteDB
}

// outbound
func QuoteDBToNqsAssetPriceQuote(quoteDB nqsdbclient.QuoteDB) nqsmodel.NqsAssetPriceQuote {
	nqsQuote := nqsmodel.NqsAssetPriceQuote{}
	quoteDBByte, _ := json.Marshal(quoteDB)
	json.Unmarshal(quoteDBByte, &nqsQuote)
	return nqsQuote
}
