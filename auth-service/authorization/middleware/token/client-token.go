// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package token

import (
	"errors"
	"time"

	"github.com/GFTN/gftn-services/utility/wwfirebase"
)

// MakerRequestData is a json-serializable type.
// This is the data structure that gets committed to firebase.
// The ApproveUserId is added from the beginning since altering  types later is not possible.
type MakerRequestData struct {
	RequestUserID string `json:"uid_request,omitempty"`
	ApproveUserID string `json:"uid_approve"`
	InstitutionID string `json:"iid,omitempty"`
	ParticipantID string `json:"pid,omitempty"`
	Status        string `json:"status,omitempty"`
	Endpoint      string `json:"endpoint,omitempty"`
	Method        string `json:"method,omitempty"`
	Timestamp     int64  `json:"timestamp_request,omitempty"`
}

/*
 * MakerRequest : After authorization, write the request to firebase in approvals (participant_approvals/super_approvals) node
 *@param{ w: http response writer, r: http request, level: string , userID: string }
 */
func MakerRequest(participantID string, institutionID string, path string, method string, level string, userID string) (string, error) {

	approvalRef := wwfirebase.FbRef.Child("/participant_approvals/")

	if level == "super" {
		approvalRef = wwfirebase.FbRef.Child("/super_approvals/")

	} else if level == "participant" {
		approvalRef = wwfirebase.FbRef.Child("/participant_approvals/")
	}

	newPostRef, err := approvalRef.Push(wwfirebase.AppContext, nil)

	if err != nil {
		LOGGER.Error("Error sending approvals to Firebase %s", err.Error())
		return "", err
	}

	timestamp := time.Now().Unix()

	if err := newPostRef.Set(wwfirebase.AppContext, &MakerRequestData{
		RequestUserID: userID,
		ApproveUserID: " ",
		InstitutionID: institutionID,
		ParticipantID: participantID,
		Status:        "request",
		Endpoint:      path,
		Method:        method,
		Timestamp:     timestamp,
	}); err != nil {
		return "", err
	}

	return newPostRef.Key, nil
}

// CheckerApprove : After authorization, this is used for approving the request (internal) and making changes in firebase approvals node.
func CheckerApprove(requestID string, iIDFromHeader string, userID string, level string) (bool, error) {

	if requestID == "" {
		LOGGER.Info("Request ID is nil", requestID)
		return false, errors.New("RequestID must be provided in a header for checker to approve")
	}

	// For default, make it participant_approvals
	approvalRef := wwfirebase.FbRef.Child("/participant_approvals/")

	if level == "super" {
		approvalRef = wwfirebase.FbRef.Child("/super_approvals/")

	} else if level == "participant" {
		approvalRef = wwfirebase.FbRef.Child("/participant_approvals/")

	}

	var requestMap MakerRequestData

	if err := approvalRef.Child(requestID).
		Get(wwfirebase.AppContext, &requestMap); err != nil {
		LOGGER.Error("Error getting approvals from Firebase %s", err.Error())
		return false, errors.New("Request wasn't made")
	}

	// get original maker of the request
	requestUser := requestMap.RequestUserID
	requestRef := approvalRef.Child(requestID)

	if requestUser == "" {
		return false, errors.New("Error getting uid_request from Firebase record")
	}

	if requestUser == userID {
		LOGGER.Error("Approver cannot be the same person as the creator of the request")
		return false, errors.New("Approver cannot be the same person as the creator of the request")
	}

	if level == "participant" {

		if iIDFromHeader == requestMap.InstitutionID {
			LOGGER.Info("Header Institution matches with the institution in the request")
		}

		/*
		 *1) the maker should not be allowed to approve the request
		 *2) the IID in the header should be the same as the one stored in firebase request object for a participant
		 *3) the status of the request in firebase should be "request"
		 */
		if requestUser != userID && requestMap.Status == "request" && iIDFromHeader == requestMap.InstitutionID {

			if err := requestRef.Update(wwfirebase.AppContext, map[string]interface{}{
				"status":      "approved",
				"uid_approve": userID,
			}); err != nil {
				LOGGER.Info("Firebase update failed")
				return false, err
			}
			return true, nil
		}
	} else {

		/*
		 *1) the maker should not be allowed to approve the request
		 *2) the IID in the header should be the same as the one stored in firebase (which was originally in the request)
		 *3) the status of the request in firebase should be "request"
		 */
		if requestUser != userID && requestMap.Status == "request" {

			if err := requestRef.Update(wwfirebase.AppContext, map[string]interface{}{
				"status":      "approved",
				"uid_approve": userID,
			}); err != nil {
				LOGGER.Info("Firebase update failed")
				return false, err
			}
			return true, nil
		}

	}

	return false, errors.New("Reached end of function, could not authorize")
}
