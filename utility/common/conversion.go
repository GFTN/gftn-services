// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package common

import (
	"math"
	"strconv"
)

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 7, 64)
}

func FloatToFixedPrecisionString(input_num float64, precision int) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', precision, 64)
}

func IntToString(input_num int64) string {
	// to convert a int number to a string
	return strconv.FormatInt(input_num, 10)
}

func ToFixedPrecision(num float64, precision int) float64 {
	multiple := float64(math.Pow10(precision))
	x := float64(math.Ceil(num*multiple)) / multiple
	return x
}

func StrPtr(str string) *string {
	return &str
}
