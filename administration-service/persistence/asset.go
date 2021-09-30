// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package persistence


type Asset struct {

	// alphanumeric code for the asset - USD, XLM, etc
	// Required: true
	AssetCode *string `json:"asset_code" bson:"asset_code"`

	// the stellar address for the asset issuer
	// Required: true
	IssuerID *string `json:"issuer_id" bson:"issuer_id"`

	// native or credit
	// Required: true
	AssetType *string `json:"asset_type" bson:"asset_type"`
}
