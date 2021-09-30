// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package persistence

type ChargesInformation struct {

	// amount of fee charged
	// Required: true
	Amount *float64 `json:"amount" bson:"amount"`

	// asset
	// Required: true
	Asset *Asset `json:"asset" bson:"asset"`
}
