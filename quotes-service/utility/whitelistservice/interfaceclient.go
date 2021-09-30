// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package whitelistservice

import "github.com/GFTN/gftn-services/gftn-models/model"

type InterfaceClient interface {
	GetWhiteListParticipantDomains(participantID string) ([]string, error)
	GetWhiteListParticipants(participantID string) ([]model.Participant, error)
	GetMutualWhiteListParticipants(participantID string) ([]model.Participant, error)
	GetMutualWhiteListParticipantDomains(participantID string) ([]string, error)
}
