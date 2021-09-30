// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package wwfirebase

import (
	"encoding/json"
	"github.com/go-errors/errors"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/utility/global-environment"
	"os"
	"time"
)

type FbTrustLogSuccess struct {
	// Log for successful trust operation
	TimeUpdated    *int64  `json:"time_updated"`
	RequestorID    *string `json:"requestor_id"`
	IssuerID       *string `json:"issuer_id"`
	AccountName    *string `json:"account_name"`
	AssetCode      *string `json:"asset_code"`
	Limit          int64   `json:"limit"`
	Status         *string `json:"status"`
	ReasonRejected *string `json:"reason_rejected"`
}

//This function transmits the trust line request result directly to firebase database.
// Since requestor and issuer are two actors on the same trust request, the firebase DB connected the portal to show the updated
// state on the UI using firebase DB as common data store
func SendFBTrustSuccess(trustRequest model.Trust) error {
	LOGGER.Debugf("Writing success result to FireBase")
	release := os.Getenv(global_environment.ENV_KEY_SERVICE_VERSION)
	homeDomain := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	endtoendId := trustRequest.EndToEndID

	if endtoendId == "" {
		LOGGER.Error("endtoendId is not set with this trust request")
		return errors.New("endtoendId is not set with this trust request")
	}

	permission := *trustRequest.Permission

	//timestamp
	timeNow := time.Now().Unix()

	if permission == "request" {
		status := "requested"
		// Post Log to Firebase for trust request
		trustSuccess := FbTrustLogSuccess{
			RequestorID: &homeDomain,
			IssuerID:    trustRequest.ParticipantID,
			Limit:       trustRequest.Limit,
			AssetCode:   trustRequest.AssetCode,
			AccountName: trustRequest.AccountName,
			Status:      &status,
			TimeUpdated: &timeNow,
		}
		mapLog := map[string]interface{}{}
		str, _ := json.Marshal(trustSuccess)
		_ = json.Unmarshal([]byte(str), mapLog)
		// Post Log to Firebase for allow operation as issuer of asset
		err := FbRef.Child(release+"/trust_requests/"+endtoendId).Update(AppContext, mapLog)
		if err != nil {
			LOGGER.Debug("error updating trust line in firebase db: %s", err.Error())
		}

	} else if permission == "allow" || permission == "revoke" {
		status := ""
		if permission == "allow" {
			status = "approved"
		} else {
			status = "revoked"
		}
		trustSuccess := FbTrustLogSuccess{
			ReasonRejected: trustRequest.ParticipantID,
			IssuerID:       &homeDomain,
			Limit:          trustRequest.Limit,
			AccountName:    trustRequest.AssetCode,
			AssetCode:      trustRequest.AccountName,
			Status:         &status,
			TimeUpdated:    &timeNow,
		}
		mapLog := map[string]interface{}{}
		str, _ := json.Marshal(trustSuccess)
		_ = json.Unmarshal([]byte(str), mapLog)

		// Post Log to Firebase for allow operation as issuer of asset
		err := FbRef.Child(release+"/trust_requests/"+endtoendId).Update(AppContext, mapLog)
		if err != nil {
			LOGGER.Debug("error updating trust line in firebase db: %s", err.Error())
		}
	}
	return errors.New("invalid permission: " + permission)

}
