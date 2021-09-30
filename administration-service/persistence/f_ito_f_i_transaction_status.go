// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package persistence

// FItoFITransactionStatus fitoFITransactionStatus
//
// Transaction's change of status information
type FItoFITransactionStatus struct {

	// This would capture the new status of a transaction while transaction travel through payment flow.
	// Required: true
	Status *string `json:"status" bson:"status"`

	// Timestamp of the status change according to World Wire
	// Required: true
	TimeStamp *int64 `json:"time_stamp" bson:"time_stamp"`
}
