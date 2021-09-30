// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package automate

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/GFTN/gftn-services/automation-service/environment"
	"github.com/GFTN/gftn-services/utility/response"
	"github.com/GFTN/gftn-services/utility/wwfirebase"
)

type Client struct {
	HttpClient    *http.Client
	twoFACheckURL string
}

type TokenBody struct {
	Token string `json:"token"`
}

// AuthorizeDeployment : Authorize the user id in the firebase id and the TOTP token
func AuthorizeDeployment(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	fID := r.Header.Get("X-Fid")

	// Dealing with User ID here (Extraction from FID and checking it)
	userID, err := decodeFID(fID)
	if err != nil {
		LOGGER.Error("Unable to extract to user ID from the X-Fid")
		response.NotifyWWError(w, r, http.StatusUnauthorized, "API-1267", errors.New("not Authorized"))
		return
	}

	authorized := authorizePermission(userID)

	if !authorized {
		response.NotifyWWError(w, r, http.StatusUnauthorized, "API-1267", errors.New("not Authorized"))
		return
	}

	next.ServeHTTP(w, r)
}

// decodeFID : Decode FID and return the user ID
func decodeFID(fID string) (string, error) {
	// VerifyIDToken helps verify token (source, timestamp, claims) through firebase.
	// It is used here to extract UID.
	token, err := wwfirebase.FbAuthClient.VerifyIDToken(wwfirebase.AppContext, fID)
	if err != nil {
		LOGGER.Errorf("%s", err.Error())
		return "", err
	}

	userID := token.UID

	return userID, nil
}

// authorizePermission : Check the super permissions for the portal user
func authorizePermission(userID string) bool {
	var rolesForSuperUser map[string]interface{}

	err := wwfirebase.FbRef.Child("/super_permissions/").Child(userID).Get(wwfirebase.AppContext, &rolesForSuperUser)
	if err != nil {
		LOGGER.Error("Error getting permission info from Firebase %s", err.Error())
		return false
	}

	if len(rolesForSuperUser) > 0 {
		superRoles := rolesForSuperUser["roles"].(map[string]interface{})
		for role, exists := range superRoles {
			// check if `admin:true` was contained in the superRoles
			if role == "admin" && exists.(bool) {
				return true
			}
		}

	}

	return false
}

// Confirm2FA : Confirm if the TOTP token from the portal is a valid token
func Confirm2FA(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	xfid := r.Header.Get("X-Fid")
	totpToken := r.Header.Get("X-Verify-Token")
	twoFACheckURL := os.Getenv(environment.ENV_KEY_AUTH_SERVICE_URL)

	httpClient := Client{
		HttpClient:    &http.Client{Timeout: time.Second * 15},
		twoFACheckURL: twoFACheckURL + "/auth/totp/check",
	}

	tokenBody := TokenBody{
		Token: totpToken,
	}

	tokenByte, _ := json.Marshal(tokenBody)

	// need to pass the TOTP token as payload to the auth-service
	req, _ := http.NewRequest(http.MethodPost, httpClient.twoFACheckURL, bytes.NewBuffer(tokenByte))
	req.Header.Add("x-fid", xfid)
	req.Header.Add("Content-Type", "application/json")
	res, restErr := http.DefaultClient.Do(req)
	if restErr != nil {
		LOGGER.Errorf("Error when calling the auth-service: %s", restErr.Error())
		response.NotifyWWError(w, r, http.StatusForbidden, "API-1267", restErr)
		return
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		if res.StatusCode == http.StatusOK {
			LOGGER.Debugf("Successfully validate the TOTP token: %s", string(body))
			//next.ServeHTTP(w, r)
		} else {
			LOGGER.Errorf("Can not verify the TOTP token: %s", string(body))
			response.NotifyWWError(w, r, http.StatusUnauthorized, "API-1267", errors.New("can not verify the TOTP token"))
			return
		}
	}
}
