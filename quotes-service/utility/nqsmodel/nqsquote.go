// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package nqsmodel

import (
	"github.com/shopspring/decimal"
	"github.com/GFTN/gftn-services/gftn-models/model"
)

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// NqsAssetPriceQuote nqsAssetPriceQuote
//
// Asset Price Quote
// swagger:model NqsAssetPriceQuote
type NqsAssetPriceQuote struct {
	TimeExpireRfi            *int64           `json:"time_expire_rfi"`
	LimitMaxOfi              *decimal.Decimal `json:"limit_max_ofi"`
	LimitMinOfi              *decimal.Decimal `json:"limit_min_ofi"`
	LimitMaxRfi              *decimal.Decimal `json:"limit_max_rfi"`
	LimitMinRfi              *decimal.Decimal `json:"limit_min_rfi"`
	OfiId                    *string          `json:"ofi_id"`
	ExchangeRate             *decimal.Decimal `json:"exchange_rate"`
	QuoteID                  *string          `json:"quote_id"`
	RfiId                    *string          `json:"rfi_id"`
	AddressReceiveRfi        *string          `json:"address_receive_rfi"`
	AddressSendRfi           *string          `json:"address_send_rfi"`
	SourceAsset              *Asset           `json:"source_asset"`
	TargetAsset              *Asset           `json:"target_asset"`
	TimeStartRfi             *int64           `json:"time_start_rfi"`
	IssuerAddressTargetAsset *string
	IssuerAddressSourceAsset *string
	QuoteResponse            *model.Quote `json:"quote_response,omitempty"`
	QuoteResponseBase64      *string      `json:"quote_response_base64,omitempty"`
	QuoteResponseSignature   *string      `json:"quote_response_signature,omitempty"`
}
