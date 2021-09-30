// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package participant

import (
	"bytes"

	"github.com/GFTN/gftn-services/gftn-models/model"
)

type ParticipantStatusError struct {
	notActiveParticipants *[]model.Participant
}

func (err *ParticipantStatusError) Error() string {
	var msgBuffer bytes.Buffer

	msgBuffer.WriteString("Participant: ")

	for _, p := range *err.notActiveParticipants {
		msgBuffer.WriteString(*p.ID)
		msgBuffer.WriteString(", ")
	}

	msgBuffer.WriteString("status not active")

	return msgBuffer.String()
}

func (err *ParticipantStatusError) GetNotActiveParticipants() []model.Participant {
	return *err.notActiveParticipants
}

func CheckStatusActive(ps ...model.Participant) error {

	var notActiveParticipants []model.Participant

	for _, p := range ps {
		if p.Status != "active" {
			notActiveParticipants = append(notActiveParticipants, p)
		}
	}

	if len(notActiveParticipants) > 0 {
		return &ParticipantStatusError{&notActiveParticipants}
	}

	return nil
}

func ExtractActiveParticipants(ps []model.Participant) []model.Participant {

	var activeParticipants []model.Participant

	for _, p := range ps {
		if p.Status == "active" {
			activeParticipants = append(activeParticipants, p)
		}
	}

	return activeParticipants
}
