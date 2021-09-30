// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package authservice

import (
	"net/http"
)

type Client struct {
	HTTP *http.Client
}
