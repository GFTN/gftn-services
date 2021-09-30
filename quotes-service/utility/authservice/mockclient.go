// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package authservice

import (
	"net/http"
)

type MockClient struct {
	HTTP *http.Client
}

func (asc *MockClient) VerifyTokenAndEndpoint(jwt string, endpoint string) (bool, error) {
	return true, nil
}
