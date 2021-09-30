// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package database

type MockDatabase struct {
	DeleteWLFlag                   bool
	DeleteWhitelistParticipantFunc func() error
	AddWLFlag                      bool
	AddWhitelistParticipantFunc    func() error
	GetWLFlag                      bool
	GetWhiteListParicipantsFunc    func() ([]string, error)
}

func (mDB *MockDatabase) DeleteWhitelistParticipant(participantID, wlParticipant string) error {
	mDB.DeleteWLFlag = true
	return mDB.DeleteWhitelistParticipantFunc()
}
func (mDB *MockDatabase) AddWhitelistParticipant(participant, wlparticipant string) error {
	mDB.AddWLFlag = true
	return mDB.AddWhitelistParticipantFunc()
}
func (mDB *MockDatabase) GetWhiteListParicipants(participantID string) ([]string, error) {
	mDB.GetWLFlag = true
	return mDB.GetWhiteListParicipantsFunc()
}

func DefaultMock() MockDatabase {
	return MockDatabase{
		DeleteWLFlag: false,
		DeleteWhitelistParticipantFunc: func() error {
			return nil
		},
		AddWLFlag: false,
		AddWhitelistParticipantFunc: func() error {
			return nil
		},
		GetWLFlag: false,
		GetWhiteListParicipantsFunc: func() ([]string, error) {
			return nil, nil
		},
	}
}
