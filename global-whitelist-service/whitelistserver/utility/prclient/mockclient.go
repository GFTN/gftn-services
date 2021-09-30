// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package prclient

import "github.com/GFTN/gftn-services/gftn-models/model"

type MockClient struct {
	GetAllParticipantsFunc func() ([]model.Participant, error)
}

func (mPR MockClient) GetAllParticipants() ([]model.Participant, error) {
	return mPR.GetAllParticipantsFunc()
}

func DefaultMock() MockClient {
	return MockClient{
		GetAllParticipantsFunc: func() ([]model.Participant, error) {
			return nil, nil
		},
	}
}
