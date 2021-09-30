// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package utils

import "strings"

func Contains(a []string, x string) bool {
	for _, n := range a {
		if strings.ToUpper(x) == strings.ToUpper(n) {
			return true
		}
	}
	return false
}

// Check if two strings is the same with case insensitive comparison
func StringsEqual(a, b string) bool {
	if strings.TrimSpace(strings.ToUpper(a)) == strings.TrimSpace(strings.ToUpper(b)) {
		return true
	}
	return false
}
