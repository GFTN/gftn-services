// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package authservice

type InterfaceClient interface {
	VerifyTokenAndEndpoint(jwt string, endpoint string) (bool, error)
}
