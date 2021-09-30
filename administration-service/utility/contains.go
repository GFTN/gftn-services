// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package utility

import "strings"

func Contains(a []string, x string) int {
	for key, n := range a {
		if x == strings.ToUpper(n) {
			return key
		}
	}
	return -1
}
