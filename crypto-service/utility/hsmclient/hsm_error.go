// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package hsmclient

import "strings"

const SESSION_HANDLE_INVALID = "0xB3"
const SIGNING_ERROR = "0x8000001A"
const UNIDENTIFIED_ERROR = "Unknown Error Message from HSM"

const NULL_HANDLE_ID = "NULL_HANDLE_ID"

func ParseErrorMsg(err error) string {
	errString := err.Error()
	LOGGER.Infof("HSM Error Message: %s", errString)
	errSubstrings := strings.Split(errString, ":")
	if len(errSubstrings) < 2 {
		LOGGER.Warningf(UNIDENTIFIED_ERROR)
		return UNIDENTIFIED_ERROR
	}
	errString = strings.TrimSpace(errSubstrings[1])
	return errString
}
