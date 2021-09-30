// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package participant

import (
	"errors"
	"strings"
)

func GetServiceUrl(urlTemplate string, participantId string) (string, error) {

	replaceTarget := "{participant_id}"

	if urlTemplate == "" || participantId == "" {
		return "", errors.New("GetServiceUrl: parameter missing")
	}
	if !strings.Contains(urlTemplate, replaceTarget) {
		return "", errors.New("GetServiceUrl: Environment variable: " + urlTemplate + " should contains " + replaceTarget + "")
	}
	res := strings.Replace(urlTemplate, replaceTarget, participantId, 1)
	return res, nil
}
