// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package middlewares

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"

	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/wwfirebase"

	"github.com/dgrijalva/jwt-go/request"
	"github.com/GFTN/gftn-services/auth-service/authorization/authutility"
	"github.com/GFTN/gftn-services/auth-service/authorization/middleware/token"
	"github.com/GFTN/gftn-services/utility/response"
	authtesting "github.com/GFTN/gftn-services/utility/testing"
)

// middleware docs: https://github.com/gorilla/mux#middleware

// LogURI : example basic middleware that just logs the uri requested
// TODO: Remove this function eventually
func LogURI(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Do stuff here ...

		LOGGER.Info("Middleware - Running LogURI")
		LOGGER.Info(r.RequestURI + " is your requested route")

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)

	})
}

/*
* SuperAuthorization : Authorization for client portal for super users
* If JWT is not enabled, the next handler is served.
* If JWT is enabled, firebase ID, institution ID, permission (request/approve), requestID (if permission is approve), participantID are expected in the headers.
* Participant ID and Institution ID are no longer mandatory because at the time it wont be necessary that those are available.
* All GET requests are direct access, if there is access :- No maker checker
* All POSTS are maker-checker except payout point which needs the current security lead/team member to validate before it gets merged in
 */
func SuperAuthorization(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Record the name of this middleware function for testing purpose
	authtesting.SaveFuncName()

	enable, _ := os.LookupEnv(global_environment.ENV_KEY_ENABLE_JWT)

	if enable == "false" {
		next.ServeHTTP(w, r)
		return
	}

	// All Headers that are expected
	fID := r.Header.Get("X-Fid")

	if fID == "" {
		LOGGER.Error("fID is empty")
		response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1012", errors.New("Header missing - Firebase token"))
		return
	}

	iID := r.Header.Get("X-Iid")

	if iID == "" {
		LOGGER.Info("X-Iid is empty, it is not mandatory for super user access.")
		// response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1003", errors.New("Header missing - Institution ID"))
		// return
	}

	uri, err := url.Parse(r.RequestURI)
	if err != nil {
		LOGGER.Error("URL parse error: ", err.Error())
		response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1021", errors.New("URL could not be parsed"))
		return
	}

	// Dealing with User ID here (Extraction from FID and checking it)
	userID, err := authutility.DecodeFID(fID)
	if err != nil {
		LOGGER.Error("Could not Decode FID %s ", err.Error())
		response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1011", errors.New("Firebase token could not be parsed"))
		return
	}

	path := uri.Path

	permission := r.Header.Get("X-Permission")
	makerChecker := false

	if permission != "" {
		makerChecker = true
	}

	var rolesForSuperUser map[string]interface{}
	errFirebaseRef := wwfirebase.FbRef.Child("/super_permissions/").Child(userID).Get(wwfirebase.AppContext, &rolesForSuperUser)
	if errFirebaseRef != nil {
		LOGGER.Error("Error getting super user permissions info from Firebase %s", errFirebaseRef.Error())
		response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1001", errFirebaseRef)
		return
	}

	superRoles := rolesForSuperUser["roles"].(map[string]interface{})

	var roleToPass string

	for role, exists := range superRoles {
		if exists.(bool) {
			roleToPass = role
		}
	}

	LOGGER.Info("Super: roleToPass: ", roleToPass)
	LOGGER.Info("Super: makerChecker: ", makerChecker)
	LOGGER.Info("Super: Method: ", r.Method)
	LOGGER.Info("Super: path", path)

	routePath, err := authutility.ExtractRoutePath(r)
	if err != nil {
		LOGGER.Error("Error while extracting route path from request")
		response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1001", errors.New("SuperAuthorization failed"))
		return
	}
	LOGGER.Infof("Extracting route path %s from request", routePath)
	authorized, errCheckAccess := CheckAccess("Super_permissions", roleToPass, makerChecker, r.Method, routePath)

	LOGGER.Info("authorized in SuperUser", authorized)

	if !authorized {
		LOGGER.Error("Could not authorize endpoint", errCheckAccess.Error())
		response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1001", errors.New("SuperAuthorization failed"))
		return
	}

	if authorized && !makerChecker {
		LOGGER.Info("Going into next handler for super auth:")
		next.ServeHTTP(w, r)
		return
	}

	if authorized && makerChecker && permission == "request" {
		LOGGER.Info("Before Maker Request Super ")
		requestID, err := token.MakerRequest("", iID, path, r.Method, "super", userID)
		LOGGER.Info("After Maker Request Super ", requestID)
		if err != nil {
			LOGGER.Info("Maker Request failed for Super Auth")
			response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1002", errors.New("SuperAuthorization Maker failed"))
			return
		}
		response.NotifySuccess(w, r, requestID)
		return
	}

	if authorized && makerChecker && permission == "approve" {
		requestID := r.Header.Get("X-Request")
		LOGGER.Info("Before Checker Approve Super")
		approved, err := token.CheckerApprove(requestID, iID, userID, "super")
		LOGGER.Info("After Checker Approve Super")

		if err != nil {
			response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1003", errors.New("SuperAuthorization failed Checker approve"))
			return

		}

		if approved {
			LOGGER.Info("Going into next handler for Super auth Checker Approve")
			next.ServeHTTP(w, r)
			return
		}
	}

	// Ideally the function will never reach here, in case it does, return
	// We have managed to find ourselves a loop hole or in case of an attack good resilient logic
	LOGGER.Error("Could not authorize %s", errCheckAccess.Error())
	response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1001", errors.New("SuperAuthorization failed"))
	return
}

/*
* ParticipantAuthorization : Authorization for client portal
* If JWT is not enabled, the next handler is served.
* If JWT is enabled, firebase ID, institution ID, permission (request/approve), requestID (if permission is approve), participantID are expected in the headers.
* The error message can be relayed back with NotifyWWError but it seems sensible to log it.
 */
func ParticipantAuthorization(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Record the name of this middleware function for testing purpose
	authtesting.SaveFuncName()

	enable, _ := os.LookupEnv(global_environment.ENV_KEY_ENABLE_JWT)

	uri, err := url.Parse(r.RequestURI)
	if err != nil {
		LOGGER.Error("URL parse error: ", err.Error())
		response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1021", errors.New("URL could not be parsed"))
		return
	}

	path := uri.Path
	LOGGER.Info(path)

	fID := r.Header.Get("X-Fid")

	if enable == "false" {
		next.ServeHTTP(w, r)
		return
	} else if r.Header.Get("Authorization") == "" && fID == "" && enable == "true" {
		LOGGER.Error("Token is missing")
		response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1018", errors.New("Header missing - either JWT or FID"))
		return
	} else if r.Header.Get("Authorization") != "" && enable == "true" {

		routePath, err := authutility.ExtractRoutePath(r)
		if err != nil {
			LOGGER.Error("Error while extracting route path from request")
			response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1001", errors.New("ParticipantAuthorization failed"))
			return
		}
		LOGGER.Infof("Extracting route path %s from request", routePath)
		isEndpointValid, err := CheckAccess("Jwt", "allow", false, r.Method, routePath)

		if !isEndpointValid {
			LOGGER.Error("JWT not authorized for endpoint the current endpoint, error message: %s", err.Error())
			response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1001", errors.New("Insufficient JWT permissions"))
			return
		}

		jwtAuthorization(w, r, next)
		return
	}

	// At this point of execution, neither FID nor JWT got missing

	iID := r.Header.Get("X-Iid")

	if iID == "" {
		LOGGER.Error("X-Iid is empty")
		response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1014", errors.New("Header missing - Institution ID"))
		return
	}

	// Dealing with User ID here (Extraction from FID and checking it)
	userID, err := authutility.DecodeFID(fID)
	if err != nil {
		LOGGER.Error("Could not Decode FID %s ", err.Error())
		response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1011", errors.New("FID could not be decoded"))
		return
	}

	pID := r.Header.Get("X-Pid")

	if pID == "" {
		LOGGER.Error("X-Pid is empty")
		response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1016", errors.New("Header missing pID"))
		return
	}

	var participantsList map[string]interface{}
	var rolesForUserInstitution map[string]interface{}
	var nodes map[string]interface{}

	participantIDFromEnv, exists := os.LookupEnv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)

	if exists {
		if participantIDFromEnv != "ww" {
			if participantIDFromEnv != pID {
				LOGGER.Error("Environment Home domain name does not match the participant ID from the header")
				response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1001", errors.New("Environment Home domain name does not match the participant ID from the header"))
				return
			}

			if err := wwfirebase.FbRef.Child("/nodes/").
				Get(wwfirebase.AppContext, &nodes); err != nil {
				LOGGER.Error("Error getting nodes info from Firebase %s", err.Error())
				response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1001", errors.New("Error getting nodes from firebase"))
				return

			}

			if len(nodes) > 0 {
				institutionIDFromNode := nodes[pID]
				if institutionIDFromNode != "" {
					if institutionIDFromNode != iID {
						LOGGER.Error("Institution From node and institution ID in header does not match")
						response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1001", errors.New("Institution From node and institution ID in header does not match"))
						return
					}

				}
			}

		}
	}

	if err := wwfirebase.FbRef.Child("/participants/"+iID).
		Get(wwfirebase.AppContext, &participantsList); err != nil {
		LOGGER.Error("Error getting participant info from Firebase %s", err.Error())
		response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1001", errors.New("Error getting participant info from Firebase"))
		return
	}

	if errRoles := wwfirebase.FbRef.Child("/participant_permissions/"+userID).Child(iID).
		Get(wwfirebase.AppContext, &rolesForUserInstitution); errRoles != nil {
		LOGGER.Error("Institution doesn't exist for this user")
		response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1001", errors.New("Institution doesn't exist for this user"))
		return
	}

	var roleToPass string

	if len(rolesForUserInstitution) > 0 {
		roles := rolesForUserInstitution["roles"].(map[string]interface{})
		for role, exists := range roles {
			if exists.(bool) {
				roleToPass = role
			}
		}
	}

	permission := r.Header.Get("X-Permission")
	makerChecker := false

	if permission != "" {
		makerChecker = true
	}

	LOGGER.Info("Participant: roleToPass: ", roleToPass)
	LOGGER.Info("Participant: makerChecker: ", makerChecker)
	LOGGER.Info("Participant: Method: ", r.Method)
	LOGGER.Info("Participant: path", path)

	routePath, err := authutility.ExtractRoutePath(r)
	if err != nil {
		LOGGER.Error("Error while extracting route path from request")
		response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1001", errors.New("ParticipantAuthorization failed"))
		return
	}
	LOGGER.Infof("Extracting route path %s from request", routePath)
	authorized, errCheckAccess := CheckAccess("Participant_permissions", roleToPass, makerChecker, r.Method, routePath)

	LOGGER.Info("authorized in participant", authorized)

	if !authorized {
		LOGGER.Error("Could not authorize endpoint", errCheckAccess.Error())
		response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1001", errors.New("ParticipantAuthorization failed"))
		return
	}

	if authorized && !makerChecker {
		LOGGER.Info("next handler for participant")
		next.ServeHTTP(w, r)
		return
	}

	if authorized && makerChecker && permission == "request" {
		LOGGER.Info("Participant Auth before maker request")
		requestID, err := token.MakerRequest(pID, iID, path, r.Method, "participant", userID)
		LOGGER.Info("Participant Auth after maker request", requestID)
		if err != nil {
			LOGGER.Info("Maker Request failed for Participant Auth")
			response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1002", errors.New("ParticipantAuthorization Maker failed"))
			return
		}
		response.NotifySuccess(w, r, requestID)
		return
	}

	if authorized && makerChecker && permission == "approve" {
		requestID := r.Header.Get("X-Request")
		LOGGER.Info("Participant Auth before Checker approve")
		approved, err := token.CheckerApprove(requestID, iID, userID, "participant")
		LOGGER.Info("Participant Auth after Checker approve")

		if err != nil {
			response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1003", errors.New("ParticipantAuthorization failed Checker approve"))
			return

		}

		if approved {
			LOGGER.Info("Going into next handler for Participant auth Checker Approve")
			next.ServeHTTP(w, r)
			return
		}
	}

	response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1001", errors.New("ParticipantAuthorization failed"))
	return
}

/*
*clientReadOnly : Authorization for client portal
*  This is currently for all read only requests (GET), this may be extended for special super user calls for endpoints that don't require maker/checker role
 */
/*func clientReadOnly(pID string, path string, iID string, userID string, method string) bool {

	authorizedParticipant, err := token.ParticipantAuthorize(method, pID, path, iID, userID)

	if err != nil {
		// Not returning from here because either of two permissions can be used for reading
		LOGGER.Info("Participant authorization failed")
	}

	authorizedSuper, err := token.SuperUserAuthorize(method, path, userID)

	if err != nil {
		// Not returning from here because either of two permissions can be used for reading
		LOGGER.Info("Super User authorization failed")
	}

	authorized := authorizedParticipant || authorizedSuper

	if !authorized {
		return false
	}

	return true

}
*/

/*
* clientAuthorization : Authorization for client portal
* This is where magic happens. All maker checker requests and authorizations come here. The participants and super users are authorized here.
* Once authorization happens, the call is forwarded to maker or checker respectively
* If a request is not authorized or should have been authorized, but not authorized, this is where you might wanna look unless you are missing a header,
* in which case, the request never came here.
 */
/*
func clientAuthorization(pID string, path string, iID string, userID string, method string, level string, permission string, requestIDFromHeader string) (bool, string) {

	LOGGER.Info("Middleware - Running Client Authorization")

	var authorized bool
	var err error

	if level == "super" {
		authorized, err = token.SuperUserAuthorize(method, path, userID)
	} else if level == "participant" {
		authorized, err = token.ParticipantAuthorize(method, pID, path, iID, userID)
	}

	if err != nil {
		LOGGER.Error("Authorize error %s", err)
	}

	if !authorized {
		LOGGER.Error("Not authorized because of : %s", err.Error())
		return false, ""
	}

	if permission == "request" {
		requestID, err := token.MakerRequest(pID, iID, path, method, level, userID)
		LOGGER.Info(requestID)
		if err != nil {
			return false, ""
		}

		// TODO: send back request ID
		// response.NotifySuccess(w, r, requestID)
		return true, requestID
	}

	if permission == "approve" {
		approved, err := token.CheckerApprove(requestIDFromHeader, iID, userID, level)
		if err != nil || !approved {
			LOGGER.Error("Error in Checker Approve: %s ", err.Error())
			return false, ""
		}

		return true, ""

	}

	// The function should never reach here, in case there is some combination of inputs that make it reach here, it is safer to return a false at the end.
	return false, ""

}
*/

// jwtAuthorization : authorization middleware checks against JWT token to authorize a user for the enpdoint
func jwtAuthorization(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	LOGGER.Info("Middleware - Running JwtAuthorization")

	// TODO : Remove comments for cleaner code

	// get JWT bearer token from authorization header
	encodedToken, err := request.OAuth2Extractor.ExtractToken(r)
	if err != nil {
		response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1017", errors.New("authentication failed, token is invalid"))
		return
	}

	// Fatal because the pod should show red on launch failure
	sBase64 := os.Getenv(global_environment.ENV_KEY_WW_JWT_PEPPER_OBJ)
	if sBase64 == "" {
		response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1001", errors.New("Pepper object is nil"))
		return
	}

	// parse token and claimsAndPayload and payload
	claimsAndPayload, valid := token.ExtractJWTClaims(encodedToken, r)

	// check for err parsing token
	// valid automatically checks exp and nfb
	if !valid {

		// w.WriteHeader(http.StatusForbidden)
		// // invalid token (likely invalid signature)
		// w.Write([]byte("authentication failed"))

		// log error to console
		// panic(err)

		// log error
		response.NotifyWWError(w, r, http.StatusUnauthorized, "AUTH-1017", errors.New("auth token is not valid"))

	} else {

		// Parse and set context to pass some session data to the handler function call.
		// Using gorilla mux and context here to share context between middleware and handler function
		// Reference: https://stackoverflow.com/questions/41876310/negroni-passing-context-from-middleware-to-handlers
		// and https://www.nicolasmerouze.com/share-values-between-middlewares-context-golang/
		parsedToken, err := ParseContext(r, claimsAndPayload)
		if err != nil {
			LOGGER.Error("ParseContext failed:" + err.Error())
			response.NotifyWWError(w, r, http.StatusForbidden, "AUTH-1017", errors.New("authentication failed, invalid token"))
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, contextKey, parsedToken)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r.WithContext(ctx))

		// }

	}

}
