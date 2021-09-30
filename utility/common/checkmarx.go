// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package common

import (
	"net/http"
	"regexp"
	"strings"
	"path/filepath"
)

/*This header is set to prevent HSTS attacks*/
func SetHSTS(w http.ResponseWriter) {
	w.Header().Set("Strict-Transport-Security", "max-age=31536000")
}

/* Setting content Security policy*/
func SetCSP(w http.ResponseWriter) {
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
}

/* 	39840: Missing X-XSS-Protection Header (Informational)
	Setting X-XSS-Protection
*/
func SetXXSS(w http.ResponseWriter) {
	w.Header().Set("X-XSS-Protection", "1")
}

/* 	39842: Missing X-Content-Type-Options Header (Informational)
	Setting X Content Type Options
*/
func SetXCTO(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
}
/* 	39844: Cacheable HTTPS response (Informational)
	Setting Cache-control
*/
func SetCacheControl(w http.ResponseWriter) {
	w.Header().Set("Cache-control", "no-store")
}
/* 	39844: Cacheable HTTPS response (Informational)
	Setting Pragma
*/
func SetPragma(w http.ResponseWriter) {
	w.Header().Set("Pragma", "no-cache")
}

/* 	39845: Input returned in response (Reflected) (Low)
	Setting Pragma
*/
func SetContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

/* validate a hash for checkmarx */
func ValidateHash(h string) bool {
	return (len(h) == 64) && IsAlphanumeric(h)
}

func IsAlphanumeric(h string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9]*$")
	return re.MatchString(h)
}

/*It is important to use absolute /canonical paths in order to avoid Path Traversal attacks*/
func IsSafePath(p string) bool {
	return !strings.Contains(p, "..") && filepath.IsAbs(p)
}

func Abs(p string) string {
	absPath, _ := filepath.Abs(p)
	return absPath
}