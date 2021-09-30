// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package middleware

import (
	"net/http"

	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"
	"github.com/GFTN/gftn-services/quotes-service/utility/authservice"
)

type MiddleWare struct {
	AuthClient authservice.InterfaceClient
}

func (mw *MiddleWare) VerifyTokenAndEndpoints(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	LOGGER.Info("checking jwt validaty")
	middlewares.ParticipantAuthorization(w, r, next)

	/* keep as reference
	jwt, err := helper.ExtractJwt(r)
	if err != nil {
		LOGGER.Error("Error extracting jwt token")
		response.NotifyWWError(w, r, http.StatusBadRequest, "EXCHANGE-1005", err)
		return
	}
	flag, err := mw.AuthClient.VerifyTokenAndEndpoint(jwt, r.URL.Path)
	if err != nil {
		LOGGER.Error("Error Verifying Token And Endpoints")
		response.NotifyWWError(w, r, http.StatusBadRequest, "EXCHANGE-1003", err)
		return
	}
	if flag == false {
		LOGGER.Error("Endpoint access not permitted")
		response.NotifyWWError(w, r, http.StatusForbidden, "EXCHANGE-1004", err)
		return
	}
	LOGGER.Info("Access granted to endpoint: " + r.URL.Path)
	next.ServeHTTP(w, r)
	*/
}
