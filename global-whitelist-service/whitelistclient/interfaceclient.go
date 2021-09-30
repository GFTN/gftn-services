// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package whitelistclient

import "github.com/GFTN/gftn-services/gftn-models/model"

type InterfaceClient interface {
	GetWhiteListParticipantDomains(participantID string) ([]string, error)
	GetWhiteListParticipants(participantID string) ([]model.Participant, error)
	CreateWhiteListParticipants(participantID, wlparitcipantID string) error
	IsParticipantWhiteListed(participantID string, targetDomain string) (bool, error)
	DeleteWhiteListParticipants(participantID, wlparitcipantID string) error
	GetMutualWhiteListParticipantDomains(participantID string) ([]string, error)
	GetMutualWhiteListParticipants(participantID string) ([]model.Participant, error)
}
