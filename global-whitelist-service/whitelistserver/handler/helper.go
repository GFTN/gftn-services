// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package handler

import (
	"net/http"
	"time"

	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"
)

func (wlh *WhitelistHandler) getMutualWL(participantID string) ([]string, error) {
	// get wl belongs to participantID
	participants, err := wlh.DBClient.GetWhiteListParicipants(participantID)

	ch := make(chan Result)
	// get wl belongs to participnatID's wl
	for idx, _ := range participants {
		go func(ch chan Result, idx int) {
			wlparticipantIDs, err := wlh.DBClient.GetWhiteListParicipants(participants[idx])
			ch <- Result{participants[idx], wlparticipantIDs, err}
		}(ch, idx)
	}
	// mutual wl screening
	var mutualwl []string
	for _ = range participants {
		select {
		case result := <-ch:
			if result.Error != nil {
				err = result.Error
				LOGGER.Error(err)
				return nil, err
			}
			for _, wl := range result.Wlparticipants {
				if participantID == wl {
					mutualwl = append(mutualwl, result.ParticipantID)
					break
				}
			}
		case <-time.After(time.Second * 5):
			// call timed out
		}
	}
	close(ch)
	LOGGER.Info("Mutual whitelist for participant", participantID, ":", mutualwl)
	return mutualwl, nil
}

func GetIdentity(req *http.Request) (string, error) {
	sessionContext, err := middlewares.GetSessionContext(req)
	if err != nil {
		LOGGER.Error(err)
		return "", err
	}
	identity := sessionContext.ParticipantID
	LOGGER.Info("Caller Identity: ", identity)
	return identity, nil
}
