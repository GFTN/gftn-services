// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package prclient

import "github.com/GFTN/gftn-services/gftn-models/model"

type InterfaceClient interface {
	GetAllParticipants() ([]model.Participant, error)
}
