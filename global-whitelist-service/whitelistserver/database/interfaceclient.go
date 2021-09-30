// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package database

type InterfaceClient interface {
	DeleteWhitelistParticipant(participantID, wlParticipant string) error
	AddWhitelistParticipant(participant, wlparticipant string) error
	GetWhiteListParicipants(participantID string) ([]string, error)
}
