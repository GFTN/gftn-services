// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package nqsmodel

import (
	"github.com/shopspring/decimal"
	"github.com/GFTN/gftn-services/gftn-models/model"
)

type NqsExchangeRequest struct {
	QuoteID                   *string          `json:"quote_id,omitempty"`
	RequestID                 *string          `json:"request_id,omitempty" `
	RfiId                     *string          `json:"rfi_id,omitempty"`
	OfiId                     *string          `json:"ofi_id,omitempty"`
	LimitMaxOfi               *decimal.Decimal `json:"limit_max_ofi,omitempty"`
	LimitMinOfi               *decimal.Decimal `json:"limit_min_ofi,omitempty"`
	LimitMaxRfi               *decimal.Decimal `json:"limit_max_rfi,omitempty"`
	LimitMinRfi               *decimal.Decimal `json:"limit_min_rfi,omitempty"`
	ExchangeRate              *decimal.Decimal `json:"exchange_rate,omitempty"`
	Amount                    *decimal.Decimal `json:"amount,omitempty"`
	SourceAsset               *Asset           `json:"source_asset,omitempty"`
	TargetAsset               *Asset           `json:"target_asset,omitempty"`
	TimeRequest               *int64           `json:"time_request,omitempty"`
	TimeQuote                 *int64           `json:"time_quote,omitempty"`
	TimeExpireOfi             *int64           `json:"time_expire_ofi,omitempty"`
	TimeStartRfi              *int64           `json:"time_start_rfi,omitempty"`
	TimeExpireRfi             *int64           `json:"time_expire_rfi,omitempty"`
	StatusQuote               *int             `json:"status_quote,omitempty"`
	TimeExecuted              *int64           `json:"time_executed,omitempty"`
	TimeCancel                *int64           `json:"time_cancel,omitempty"`
	AddressReceiveRfi         *string          `json:"address_receive_rfi,omitempty"`
	AddressSendRfi            *string          `json:"address_send_rfi,omitempty"`
	AddressReceiveOfi         *string          `json:"address_receive_ofi,omitempty"`
	AddressSendOfi            *string          `json:"address_send_ofi,omitempty"`
	AccountReceiveRfi         *string          `json:"account_receive_rfi,omitempty"`
	AccountSendRfi            *string          `json:"account_send_rfi,omitempty"`
	AccountReceiveOfi         *string          `json:"account_receive_ofi,omitempty"`
	AccountSendOfi            *string          `json:"account_send_ofi,omitempty"`
	IssuerAddressTargetAsset  *string
	IssuerAddressSourceAsset  *string
	ExchangeRequestBase64     *string
	ExchangeRequestSignBase64 *string
	ExchangeRequestDecode     *model.Exchange
}
