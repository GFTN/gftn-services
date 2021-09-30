// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package handler

import (
	"net/http"

	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"
)

type MiddleWare struct {
}

func (mw *MiddleWare) VerifyTokenAndEndpoints(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	LOGGER.Info("checking jwt validaty")
	middlewares.ParticipantAuthorization(w, r, next)
}
