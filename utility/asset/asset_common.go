// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package asset

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/xdr"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"strings"
)

func DecodeStellarError(err error) (error) {
	if hz_err, ok := err.(*horizon.Error); ok {
		result_bytes := hz_err.Problem.Extras["result_xdr"]
		var b64 string
		err_1 := json.Unmarshal(result_bytes, &b64)
		rawr := strings.NewReader(b64)
		b64r := base64.NewDecoder(base64.StdEncoding, rawr)
		var result xdr.TransactionResult
		count, err_1 := xdr.Unmarshal(b64r, &result)
		LOGGER.Debug("*** Submit transaction error: ", result.Result.Code, result.FeeCharged, count)
		if err_1 == nil {
			errMsg := fmt.Sprintf("Code: %v", result.Result.Code)
			err = errors.New(errMsg)
		}
	}
	return err
}

func DecodeStellarPaymentError(err error) (error) {
	if hz_err, ok := err.(*horizon.Error); ok {
		result_bytes := hz_err.Problem.Extras["result_xdr"]
		LOGGER.Debug("ERROR: %v", string(result_bytes))
		var b64 string
		err_1 := json.Unmarshal(result_bytes, &b64)
		rawr := strings.NewReader(b64)
		b64r := base64.NewDecoder(base64.StdEncoding, rawr)
		var result xdr.TransactionResult
		_, err_1 = xdr.Unmarshal(b64r, &result)
		r := result.Result
		r2 := *r.Results
		if r2 != nil {
			LOGGER.Debugf("Error Code: %v", r2[0].Tr.PaymentResult.Code)
			errMsg := fmt.Sprintf("Error Code: %v", r2[0].Tr.PaymentResult.Code)
			err = errors.New(errMsg)
			return err
		}

		if err_1 == nil {
			errMsg := fmt.Sprintf("Code: %v", result.Result.Code)
			err = errors.New(errMsg)
		}
	}
	return err
}

func GetAssetType(assetCode string) (string) {
	if assetCode == "XLM" || assetCode == "xlm" {
		return model.AssetAssetTypeNative
	} else if strings.HasSuffix(assetCode, model.AssetAssetTypeDO) {
		return model.AssetAssetTypeDO
	}
	return model.AssetAssetTypeDA
}
