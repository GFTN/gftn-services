// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package authconstants

import (
	"errors"
	"strings"

	util "github.com/GFTN/gftn-services/utility/common"
)

func endpointForAuthConstant(endpointMap map[string][]string, endpoint string) (string, error) {

	for k := range endpointMap {
		if strings.Contains(endpoint, k) {
			return k, nil
		}
	}
	return "", errors.New("Endpoint not in map")
}

/*
 * GetEndpointPermission :  return the list of permission that the endpoint has
 *
 * @param {{ endpoint : string, method : string  }}
 * @returns { list of strings i.e the permissions associated with this specific endpoint. For example /trust should only be allowed by participant_manager role.}
 * @memberOf AuthService
 */
func GetEndpointPermissions(endpoint string, method string, level string) []string {

	if method == "GET" {
		if level == "super" {
			endpointToCheck, _ := endpointForAuthConstant(EndpointSuperPermissionsForGet(), endpoint)
			return EndpointSuperPermissionsForGet()[endpointToCheck]
		} else if level == "participant" {
			endpointToCheck, _ := endpointForAuthConstant(EndpointParticipantPermissionsForGet(), endpoint)
			return EndpointParticipantPermissionsForGet()[endpointToCheck]
		}
	} else if method == "POST" {
		if level == "super" {
			endpointToCheck, _ := endpointForAuthConstant(EndpointSuperPermissionsForPost(), endpoint)
			return EndpointSuperPermissionsForPost()[endpointToCheck]
		} else if level == "participant" {
			endpointToCheck, _ := endpointForAuthConstant(EndpointParticipantPermissionsForPost(), endpoint)
			return EndpointParticipantPermissionsForPost()[endpointToCheck]
		}
	}

	return []string{}
}

// EndpointJwtPermissions : for jwt-authentication
// This are for controlling the admin endpoints. This needs to be reviewed by the security team
// This will also require require strict control because these are endpoints that the api will hit.
// We don't want to allow unnecessary api-requests pass through
var EndpointJwtPermissions = func() map[string]bool {
	return map[string]bool{
		"/accounts":                   true,
		"/assets":                     true,
		"/assets/participants":        true,
		"/assets/issued":              true,
		"/assets/accounts":            true,
		"/balances/accounts":          true,
		"/obligations":                true,
		"/exchange":                   true,
		"/fees":                       true,
		"/participants":               true,
		"/participants/whitelist":     true,
		"/payout":                     true,
		"/quotes":                     true,
		"/quotes/request":             true,
		"/sign":                       true,
		"/token/refresh":              true,
		"/transactions":               true,
		"/transactions/reply":         true,
		"/transactions/send":          true,
		"/transactions/settle/da":     true,
		"/transactions/settle/do":     true,
		"/trust":                      true,
		"/auth/participant-auth-test": true,
		"/address":                    true,
		"/fundings/instruction":       true,
		"/fundings/send":              true,
		"/message":                    true,
	}

}

// EndpointSuperPermissionsForGet : Permission based on Chase's permission branch on world wire services
var EndpointSuperPermissionsForGet = func() map[string][]string {
	return map[string][]string{
		"/v1/registry/participants":     []string{"admin", "manager"},
		"/v1/onboarding/accounts":       []string{"admin", "manager"},
		"v1/admin/anchor/assets/issued": []string{"admin", "manager"},
		"/v1/admin/blocklist":           []string{"admin", "manager"},
		"/maker-checker":                []string{"admin", "manager"},
	}
}

// EndpointSuperPermissionsForPost lists endpoints for super user permissions for post
// Will go through maker checker flow
var EndpointSuperPermissionsForPost = func() map[string][]string {
	// TODO: endpoints with variables in between need to be fixed.

	return map[string][]string{
		"/v1/admin/pr":            []string{"admin", "manager"},
		"/v1/admin/anchor":        []string{"admin", "manager"},
		"/v1/onboarding/accounts": []string{"admin", "manager"},
		"/v1/deploy/participant":  []string{"admin", "manager"},
		"/v1/admin/blocklist":     []string{"admin", "manager"},
		"/v1/admin/suspend":       []string{"admin", "manager"},
		"/v1/admin/reactivate":    []string{"admin", "manager"},
		"/maker-checker":          []string{"admin", "manager"},
	}
}

// EndpointParticipantPermissionsForGet lists endpoints for participants to get
var EndpointParticipantPermissionsForGet = func() map[string][]string {

	return map[string][]string{
		"/v1/admin/pr":                      []string{"admin", "manager"},
		"/v1/admin/pr/domain":               []string{"admin", "manager"},
		"/v1/anchor/assets/issued":          []string{"admin", "manager"},
		"/v1/anchor/address":                []string{"admin", "manager"},
		"/v1/client/accounts":               []string{"admin", "manager"},
		"/v1/client/assets":                 []string{"admin", "manager"},
		"/v1/client/assets/accounts":        []string{"admin", "manager"},
		"/v1/client/assets/issued":          []string{"admin", "manager"},
		"/v1/client/assets/participants":    []string{"admin", "manager"},
		"/v1/client/balances/accounts":      []string{"admin", "manager"},
		"/v1/client/obligations":            []string{"admin", "manager"},
		"/v1/client/participants/whitelist": []string{"admin", "manager"},
		"/v1/client/participants":           []string{"admin", "manager"},
		"/v1/client/transactions":           []string{"admin", "manager"},
		"/maker-checker":                    []string{"admin", "manager"},
	}
}

// EndpointParticipantPermissionsForPost lists endpoints for participants to post
// Will go through maker checker flow
var EndpointParticipantPermissionsForPost = func() map[string][]string {

	return map[string][]string{
		"/v1/client/participants":           []string{"admin", "manager"},
		"/v1/client/participants/whitelist": []string{"admin"},
		"/v1/client/transactions/settle/da": []string{"admin", "manager"},
		"/v1/client/transactions/settle/do": []string{"admin", "manager"},
		"/v1/client/trust":                  []string{"admin", "manager"},
		"/v1/anchor/fundings/instruction":   []string{"admin", "manager"},
		"/v1/anchor/fundings/send":          []string{"admin", "manager"},
		"/v1/anchor/trust":                  []string{"admin"},
		"/v1/client/assets":                 []string{"admin", "manager"},
		"/maker-checker":                    []string{"admin", "manager"},
	}
}

// EndpointSuperNoMakerChecker : set of endpoints that don't need maker checker for POST requests
// Updating this list requires critical review from Security team and all functional leaders at the time.
var EndpointSuperNoMakerChecker = func() map[string]bool {

	return map[string]bool{
		"/v1/admin/payout":                   true,
		"/v1/admin/payout/csv":               true,
		"/v1/admin/accounts/" + util.ISSUING: true,
		"/v1/admin/accounts/" + util.DEFAULT: true,
		"/direct-maker-checker":              true,
	}

}

// EndpointParticipantNoMakerChecker : set of endpoints that don't need maker checker for POST
// Updating this list requires critical review from Security team and all functional leaders at the time.
var EndpointParticipantNoMakerChecker = func() map[string]bool {

	return map[string]bool{
		"/direct-maker-checker": true,
	}

}
