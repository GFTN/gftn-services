// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package authutility

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/op/go-logging"
	"github.com/GFTN/gftn-services/auth-service/authorization/middleware/authconstants"
	"github.com/GFTN/gftn-services/utility/response"
	"github.com/GFTN/gftn-services/utility/wwfirebase"
)

// LOGGER : logs middleware package
var LOGGER = logging.MustGetLogger("middlewares")

type msgStruct struct {
	Msg string
}

// DecodeFID : Decode FID and return the user
func DecodeFID(fID string) (string, error) {
	// VerifyIDToken helps verify token (source, timestamp, claims) through firebase.
	// It is used here to extract UID.
	var userID string
	token, err := wwfirebase.FbAuthClient.VerifyIDToken(wwfirebase.AppContext, fID)
	if err != nil {
		return "", err
	}

	userID = token.UID

	return userID, nil
}

// CheckEndpointPermission : checks for permissions in auth constants file
//  Params {{ endpoint: string, role : string, method : string., level : string }}
func CheckEndpointPermission(endpoint string, role string, method string, level string) bool {
	permissions := authconstants.GetEndpointPermissions(endpoint, method, level)

	for i := range permissions {
		if role == permissions[i] {
			return true
		}
	}
	return false
}

// DoSomething : final test endpoint logic (executes after passing jwt middleware)
// $  curl -X POST -d '{"msg":"some really cool message"}' http://localhost:8080/test -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJuYW1lIjoidXNlcjEyMyIsInBhc3N3b3JkIjoiMTIzNCIsImZpcnN0TmFtZSI6IkpvbiIsImxhc3ROYW1lIjoiRG9lIiwiZG9iIjoiMTIvMTEvMTk5MSIsImVtYWlsIjoidXNlckBnbWFpbC5jb20iLCJhZGRyZXNzIjp7InN0cmVldCI6IjU1NSBCYXlzaG9yZSBCbHZkIiwiY2l0eSI6IlRhbXBhIiwic3RhdGUiOiJGbG9yaWRhIiwiemlwIjoiMzM4MTMifX0sImNvdW50IjoxLCJpYXQiOjE1NDYwNzMyMzcsImF1ZCI6IltcIi90ZXN0XCIsXCIvdGVzdDFcIixcIi90ZXN0MlwiXSJ9.opHpsnB4Glrnyqm5_pFXN-OuSyRde8a_-l1uB5qA56g"
func DoSomething(w http.ResponseWriter, r *http.Request) {
	LOGGER.Info("Token passed! - hitting logic in endpoint")

	j, _ := json.Marshal("Your JWT token passed!")
	w.Write(j)
}

// ClientToken : final test endpoint logic (executes after passing client token middleware)
// $  curl -X POST -d '{"msg":"some really cool message"}' http://localhost:8080/test -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJuYW1lIjoidXNlcjEyMyIsInBhc3N3b3JkIjoiMTIzNCIsImZpcnN0TmFtZSI6IkpvbiIsImxhc3ROYW1lIjoiRG9lIiwiZG9iIjoiMTIvMTEvMTk5MSIsImVtYWlsIjoidXNlckBnbWFpbC5jb20iLCJhZGRyZXNzIjp7InN0cmVldCI6IjU1NSBCYXlzaG9yZSBCbHZkIiwiY2l0eSI6IlRhbXBhIiwic3RhdGUiOiJGbG9yaWRhIiwiemlwIjoiMzM4MTMifX0sImNvdW50IjoxLCJpYXQiOjE1NDYwNzMyMzcsImF1ZCI6IltcIi90ZXN0XCIsXCIvdGVzdDFcIixcIi90ZXN0MlwiXSJ9.opHpsnB4Glrnyqm5_pFXN-OuSyRde8a_-l1uB5qA56g"
func ClientToken(w http.ResponseWriter, r *http.Request) {
	LOGGER.Info("Client Token being hit")
}

// ServiceCheck : Service check for auth service
func ServiceCheck(w http.ResponseWriter, r *http.Request) {
	// w.WriteHeader(http.StatusOK)
	response.NotifySuccess(w, r, "Successful")
}

// ComparePaths : compares two input routes and trims surrounding whitespace
func ComparePaths(expr string, target string) bool {
	expr = strings.TrimSpace(expr)
	target = strings.TrimSpace(target)
	if expr != target {
		return false
	}
	return true
}

// ExtractRoutePath : uses mux function to get the requested raw path used by the mux router
// (ie: with path params formated with in format /some/route/{path}/{params})
// returns path, err
func ExtractRoutePath(r *http.Request) (string, error) {
	route := mux.CurrentRoute(r)
	routePath, err := route.GetPathTemplate()
	if err != nil {
		return "", err
	}
	return routePath, nil
}
