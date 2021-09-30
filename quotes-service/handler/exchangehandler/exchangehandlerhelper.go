// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package exchangehandler

import (
	"github.com/shopspring/decimal"

	b "github.com/stellar/go/build"
)

func GetCreditAsset(assetCode string, assetAddressIssuer string, amount decimal.Decimal) (creditAsset interface{}) {
	creditAsset = b.CreditAmount{}
	if assetCode == "xlm" || assetCode == "XLM" {
		creditAsset = b.NativeAmount{Amount: amount.Round(7).String()}
	} else {
		creditAsset = b.CreditAmount{Code: assetCode, Issuer: assetAddressIssuer,
			Amount: amount.Round(7).String()}
	}
	return creditAsset
}

func CheckQuoteStatus() {
	//call database

}
