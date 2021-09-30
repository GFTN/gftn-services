// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package middlewares

import (
	"encoding/json"
	"errors"

	"github.com/GFTN/gftn-services/auth-service/authorization/authutility"
	authPermissions "github.com/GFTN/gftn-services/auth-service/authorization/middleware/permissions"
)

// Roles : defines user & JWT permissions needed to access an endpoint
type Permissions struct {
	Permissions Groups
}

type Groups struct {
	Jwt                     Default
	Participant_permissions Default
	Super_permissions       Default
}

type Default struct {
	Default       Method
	Maker_checker Method
}

type Method struct {
	Method Operation
}

type Operation map[string]Endpoint

// type Operation struct {
// 	GET    Endpoint
// 	POST   Endpoint
// 	PUT    Endpoint
// 	DELETE Endpoint
// }

type Endpoint struct {
	// Endpoint Endpoint
	Endpoint map[string]Role
}

type Role struct {
	Role Permit
}

type Permit struct {
	Allow   bool
	Admin   bool
	Manager bool
	Viewer  bool
}

// CheckAccess : gets roles needed for an endpoint and permission type
// permissionGroup = Jwt | Super_permissions | Participant_permissions
// hasRole = admin | manager | viewer | allow (for jwt only)
// makerChecker = true (ie: maker/checker required) | false (ie: maker/checker NOT required) NOTE: does not matter if the value is true or false for JWT related endpoint group since JWT does not implement a maker/checker flow
// requestedMethod = GET | PUT | POST | DELETE
// requestedEndpoint = path that the inbound request is attempting to reach
// userRole = admin | manager | viewer | allow (for jwt only)
func CheckAccess(permissionGroup string, hasRole string, makerChecker bool, requestedMethod string, requestedEndpoint string) (bool, error) {

	p := authPermissions.Permissions()
	// permissionsByte, err := ioutil.ReadFile("../permissions.json")
	var permissions Permissions
	// if err != nil {
	// 	panic(err)
	// }
	err := json.Unmarshal([]byte(p), &permissions)
	// fmt.Print(string(permissionsByte))

	if err != nil {
		LOGGER.Error("Error while parsing JSON")
		return false, errors.New("not authorized, no matching permissions")
	}

	if permissionGroup == "Jwt" {

		// endpoints requiring JWT
		// fmt.Print("\n\n" + requestedMethod + " | jwt endpoints:\n")
		jwtEndp := permissions.Permissions.Jwt.Default.Method[requestedMethod].Endpoint
		for key, value := range jwtEndp {
			// fmt.Println(key+" - Allow: ", value.Role.Allow)
			if hasRole == "allow" && value.Role.Allow == true && authutility.ComparePaths(key, requestedEndpoint) {
				// fmt.Println("Permissions Succeeded! "+key+" - Allow: ", value.Role.Allow)
				return true, nil
			}
		}

	}

	if permissionGroup == "Super_permissions" && makerChecker == false {

		// super user endpoints
		// fmt.Print("\n\n" + requestedMethod + " | super user endpoints:\n")
		superEndpDef := permissions.Permissions.Super_permissions.Default.Method[requestedMethod].Endpoint
		for key, value := range superEndpDef {
			// fmt.Println(key+" - Admin: ", value.Role.Admin)
			// fmt.Println(key+" - Manager: ", value.Role.Manager)
			if hasRole == "admin" && value.Role.Admin == true && authutility.ComparePaths(key, requestedEndpoint) {
				// fmt.Println("\nPermissions Succeeded! "+key+" - Admin: ", value.Role.Admin)
				return true, nil
			}
			if hasRole == "manager" && value.Role.Manager == true && authutility.ComparePaths(key, requestedEndpoint) {
				// fmt.Println("\nPermissions Succeeded! "+key+" - Manager: ", value.Role.Manager)
				return true, nil
			}
		}

	}

	if permissionGroup == "Super_permissions" && makerChecker == true {

		// super user endpoints requiring maker/checker
		// fmt.Print("\n\n" + requestedMethod + " | super user + maker/checker endpoints:\n")
		superEndpMC := permissions.Permissions.Super_permissions.Maker_checker.Method[requestedMethod].Endpoint
		for key, value := range superEndpMC {
			// fmt.Println(key+" - Admin: ", value.Role.Admin)
			// fmt.Println(key+" - Manager: ", value.Role.Manager)
			if hasRole == "admin" && value.Role.Admin == true && authutility.ComparePaths(key, requestedEndpoint) {
				// fmt.Println("\nPermissions Succeeded! "+key+" - Admin: ", value.Role.Admin)
				return true, nil
			}
			if hasRole == "manager" && value.Role.Manager == true && authutility.ComparePaths(key, requestedEndpoint) {
				// fmt.Println("\nPermissions Succeeded! "+key+" - Manager: ", value.Role.Manager)
				return true, nil
			}
		}

	}

	if permissionGroup == "Participant_permissions" && makerChecker == false {

		// participant user endpoints
		// fmt.Print("\n\n" + requestedMethod + " | participant endpoints:\n")
		participantEndpDef := permissions.Permissions.Participant_permissions.Default.Method[requestedMethod].Endpoint
		for key, value := range participantEndpDef {
			// fmt.Println(key+" - Admin: ", value.Role.Admin)
			// fmt.Println(key+" - Manager: ", value.Role.Manager)
			if hasRole == "admin" && value.Role.Admin == true && authutility.ComparePaths(key, requestedEndpoint) {
				// fmt.Println("\nPermissions Succeeded! "+key+" - Admin: ", value.Role.Admin)
				return true, nil
			}
			if hasRole == "manager" && value.Role.Manager == true && authutility.ComparePaths(key, requestedEndpoint) {
				// fmt.Println("\nPermissions Succeeded! "+key+" - Manager: ", value.Role.Manager)
				return true, nil
			}
		}

	}

	if permissionGroup == "Participant_permissions" && makerChecker == true {

		// participant user endpoints requiring maker/checker
		// fmt.Print("\n\n" + requestedMethod + " | participant + maker/checker endpoints:\n")
		participantEndpMC := permissions.Permissions.Participant_permissions.Maker_checker.Method[requestedMethod].Endpoint
		for key, value := range participantEndpMC {
			// fmt.Println(key+" - Admin: ", value.Role.Admin)
			// fmt.Println(key+" - Manager: ", value.Role.Manager)
			if hasRole == "admin" && value.Role.Admin == true && authutility.ComparePaths(key, requestedEndpoint) {
				// fmt.Println("\nPermissions Succeeded! "+key+" - Admin: ", value.Role.Admin)
				return true, nil
			}
			if hasRole == "manager" && value.Role.Manager == true && authutility.ComparePaths(key, requestedEndpoint) {
				// fmt.Println("\nPermissions Succeeded! "+key+" - Manager: ", value.Role.Manager)
				return true, nil
			}
		}

	}

	return false, errors.New("not authorized, no matching permissions")

}
