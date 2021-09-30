// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package client

type PaymentListenerClient interface {
	SubscribePayments(distAccountName string) (err error)
}
