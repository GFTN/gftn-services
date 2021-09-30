// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/global-whitelist-service/whitelistserver/database"
	"github.com/GFTN/gftn-services/global-whitelist-service/whitelistserver/utility/prclient"
)

func strToPtr(str string) *string {
	var strPtr *string
	strPtr = &str
	return strPtr
}

func TestCreateWLParticipant(t *testing.T) {
	mockDBClient := database.DefaultMock()
	mockPRClient := prclient.DefaultMock()
	wlh := WhitelistHandler{
		DBClient: &mockDBClient,
		PRClient: &mockPRClient,
	}
	mockPRClient.GetAllParticipantsFunc = func() ([]model.Participant, error) {
		part := model.Participant{ID: strToPtr("testIDwhobeingWL")}
		list := []model.Participant{}
		list = append(list, part)
		return list, nil
	}
	req, err := http.NewRequest("POST", "/{participant_id}", bytes.NewBufferString(`{"participant_id":"testIDwhobeingWL"}`))
	if err != nil {
		t.Fatal(err)
	}
	mux.SetURLVars(req, map[string]string{"participant_id": "testIDwhoWL"})
	rw := httptest.NewRecorder()
	handler := http.HandlerFunc(wlh.CreateWLParticipant)
	handler.ServeHTTP(rw, req)
	status := rw.Code
	Convey("Successful create WL", t, func() {
		So(status, ShouldEqual, http.StatusOK)
		So(mockDBClient.AddWLFlag, ShouldEqual, true)
	})
}

func TestCreateWLParticipantNegative(t *testing.T) {
	mockDBClient := database.DefaultMock()
	mockPRClient := prclient.DefaultMock()
	wlh := WhitelistHandler{
		DBClient: &mockDBClient,
		PRClient: &mockPRClient,
	}
	mockPRClient.GetAllParticipantsFunc = func() ([]model.Participant, error) {
		part := model.Participant{ID: strToPtr("failedID")}
		list := []model.Participant{}
		list = append(list, part)
		return list, nil
	}
	req, err := http.NewRequest("POST", "/{participant_id}", bytes.NewBufferString(`{"participant_id":"testIDwhobeingWL"}`))
	if err != nil {
		t.Fatal(err)
	}
	mux.SetURLVars(req, map[string]string{"participant_id": "testIDwhoWL"})
	rw := httptest.NewRecorder()
	handler := http.HandlerFunc(wlh.CreateWLParticipant)
	handler.ServeHTTP(rw, req)
	status := rw.Code

	Convey("Successful failed create WL", t, func() {
		So(status, ShouldEqual, http.StatusBadRequest)
		So(mockDBClient.AddWLFlag, ShouldEqual, false)

	})
}

func TestDeleteWLParticipant(t *testing.T) {
	mockDBClient := database.DefaultMock()
	mockPRClient := prclient.DefaultMock()
	wlh := WhitelistHandler{
		DBClient: &mockDBClient,
		PRClient: &mockPRClient,
	}
	req, err := http.NewRequest("POST", "/{participant_id}", bytes.NewBufferString(`{"participant_id":"testIDwhobeingWL"}`))
	if err != nil {
		t.Fatal(err)
	}
	mux.SetURLVars(req, map[string]string{"participant_id": "testIDwhoWL"})
	rw := httptest.NewRecorder()
	handler := http.HandlerFunc(wlh.DeleteWLParticipant)
	handler.ServeHTTP(rw, req)
	status := rw.Code

	Convey("Successful delete WL", t, func() {
		So(status, ShouldEqual, http.StatusOK)
		So(mockDBClient.DeleteWLFlag, ShouldEqual, true)

	})
}

func TestGetWLParticipant(t *testing.T) {
	mockDBClient := database.DefaultMock()
	mockPRClient := prclient.DefaultMock()
	wlh := WhitelistHandler{
		DBClient: &mockDBClient,
		PRClient: &mockPRClient,
	}
	req, err := http.NewRequest("GET", "/{participant_id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	mux.SetURLVars(req, map[string]string{"participant_id": "testIDwhoWL"})
	rw := httptest.NewRecorder()
	handler := http.HandlerFunc(wlh.GetWLParticipants)
	handler.ServeHTTP(rw, req)
	status := rw.Code

	Convey("Successful delete WL", t, func() {
		So(status, ShouldEqual, http.StatusOK)
		So(mockDBClient.GetWLFlag, ShouldEqual, true)

	})
}
