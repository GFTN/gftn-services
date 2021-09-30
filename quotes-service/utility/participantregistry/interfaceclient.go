// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package participantregistry

import "github.com/GFTN/gftn-services/gftn-models/model"

type InterfaceClient interface {
	GetAllParticipants() ([]model.Participant, error)
	GetParticipantForDomain(participantID string) (model.Participant, error)
	GetParticipantAccount(domain string, account string) (string, error)
}
