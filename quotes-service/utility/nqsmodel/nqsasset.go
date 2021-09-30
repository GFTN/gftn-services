// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package nqsmodel

// Asset asset
//
// Details of the asset being transacted
// swagger:model Asset
type Asset struct {

	// Alphanumeric code for the asset - USD, XLM, etc
	// Required: true
	AssetCode *string `json:"asset_code"`

	// Asset type can be native or digital asset(DA) or digital obligation(DO)
	// Required: true
	// Enum: [DO DA native]
	AssetType *string `json:"asset_type"`

	// The stellar address for the asset issuer
	// Required: true
	IssuerID *string `json:"issuer_id"`
}

const (

	// AssetAssetTypeDO captures enum value "DO"
	AssetAssetTypeDO string = "DO"

	// AssetAssetTypeDA captures enum value "DA"
	AssetAssetTypeDA string = "DA"

	// AssetAssetTypeNative captures enum value "native"
	AssetAssetTypeNative string = "native"
)
