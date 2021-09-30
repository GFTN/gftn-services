// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package utility

import "strings"

func Contains(a []string, x string) bool {
	for _, n := range a {
		if strings.ToUpper(x) == strings.ToUpper(n) {
			return true
		}
	}
	return false
}
