// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package testing

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"runtime"
	"strings"

	"github.com/gorilla/mux"
	"github.com/GFTN/gftn-services/utility/ds"
)

var (
	funcName   string
	ignoreSet  ds.Set
	adminSet   ds.Set
	clientSet  ds.Set
	isInit     = false                                                                      // This flag signals whether or not the testing has been initialized before running `AuthWalker` or other functions in this testing module.
	ignoreKeys = [5]string{"internal", "service_check", "helloworldwire", "check", "debug"} // Keywords that are included in any route to be ignored for the testing
	adminKeys  = [3]string{"admin", "onboarding", "deploy"}                                 // Keywords that indicate routes to be accessed by a super user
	clientKeys = [2]string{"client", "anchor"}                                              // Keywords that indicate routes to be accessed by a user
)

const (
	// Define authorization middleware names
	SUPER       = "SuperAuthorization"
	PARTICIPANT = "ParticipantAuthorization"
	// `ROUTELEVEL` indicates the depth of routePath where a specific keyword is located.
	ROUTELEVEL = 2
)

/*
	This function is used by mux.Router to walk through routes that are registered with mux.Router.
	AuthWalker retrieves routes and validates that auth middleware is properly implemented with the routes for different accessibility levels.
*/
func AuthWalker(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {

	// Ensure that auth test was initialized
	if !isInit {
		return errors.New("Auth test should be initialized first")
	}

	// Accessing the path from the current route
	routePath, err := route.GetPathTemplate()
	if err != nil {
		fmt.Errorf("Error while getting route path: %s", err.Error())
		return errors.New("Cannot access current route")
	}
	fmt.Println("Testing route " + routePath)

	/*
		Assume that any `routePath` follows the current convention of `/service_version/xxx/xxx',
		and the first `xxx` defines a specific accessibility level for users.
	*/
	keywords := strings.Split(routePath, "/")

	targetKeyword := keywords[ROUTELEVEL]

	// If the current route is to be ignored for, e.g., internal use, skip to the next route or completion.
	fmt.Println("Identifying keyword: " + targetKeyword)

	if ignoreSet.Contains(targetKeyword) {
		fmt.Println("Auth not needed for current route")
		fmt.Println("Continue testing other routes...")
		return nil
	}

	// Flag to decide the user accessibility level
	isAdmin := false

	// If a keyword in the routePath, is defined in either adminSet or clientSet, isDefined will eventually evaluate to `true`
	isDefined := false

	/*
		The following is to set the flag of user accessibility levels for the currently accessed route.
	*/

	fmt.Println("Looking for client keywords...")
	if clientSet.Contains(targetKeyword) {
		isAdmin = false
		isDefined = true
		fmt.Println("Client keyword defined")
	}

	fmt.Println("Looking for admin keywords...")
	if adminSet.Contains(targetKeyword) {
		isAdmin = true
		isDefined = true
		fmt.Println("Admin keyword defined")
	}

	// For security reason, we cannot accept a route with undefined keywords.
	if !isDefined {
		return errors.New("Cannot identify RoutePath with defined keywords for admin and client")
	}

	/*
		At this point, we can consider the route should have been implemented with the auth middlewares.
		The rest of code ensures if the route is implemented appropriately with the middlewares.
	*/

	handler := route.GetHandler()

	server := httptest.NewServer(handler)
	defer server.Close()

	client := &http.Client{}
	methods, err := route.GetMethods()
	if err != nil {
		fmt.Errorf("Error while getting route methods; %s", err.Error())
		return errors.New("Cannot fetch route methods")
	}

	// Requests are made given different methods with the same route.
	for _, m := range methods {

		req, err := http.NewRequest(m, server.URL, nil)
		if err != nil {
			fmt.Errorf("Error while building request; %s", err.Error())
			return errors.New("Cannot build request")
		}

		res, err := client.Do(req)

		if err != nil {
			fmt.Errorf("Error while sending %s request; %s", m, err.Error())
			return errors.New("Cannot send request")
		}

		if res.StatusCode != http.StatusUnauthorized {
			fmt.Errorf("Missing authorization; auth needed for current route")
			return errors.New("Missing authorization")
		}

		/*
			From now on, it is certain that this route is implemented with the authorization middleware.
			The following snippet is to check if the middleware implemented is correct for different user accessibility levels
			with the saved name of the authorization middleware.
		*/

		if !isAdmin && funcName == PARTICIPANT {
			fmt.Println("Participant authorization is properly set.")
		} else if isAdmin && funcName == SUPER {
			fmt.Println("Super authorization is properly set.")
		} else {
			fmt.Errorf("Authorization is not properly set.")
			return errors.New("Authorization is not properly set for " + routePath)
		}
	}

	return nil
}

/*
	This function returns the calling function name.
*/
func SaveFuncName() {
	// Omit accessing function names from the call stack if this testing was not initialized.
	if !isInit {
		return
	}

	// We use the `runtime` package to save the name of the previously called function.
	pc, _, _, ok := runtime.Caller(1)

	if !ok {
		panic("Cannot retrieve the function name")
	}
	// Acquire the function name given a program counter
	f := runtime.FuncForPC(pc)
	re := regexp.MustCompile(`\.`).Split(f.Name(), -1)
	funcName = re[len(re)-1]
}

/*
	This function does all the things and sets `isInit` to true before the testing is run.
*/
func InitAuthTesting() {
	fmt.Println("Initializing auth test...")
	ignoreSet = ds.NewSet()
	adminSet = ds.NewSet()
	clientSet = ds.NewSet()
	// Inserting all the keywords that are defined for this testing above
	for _, e := range ignoreKeys {
		ignoreSet.Add(e)
	}

	for _, e := range adminKeys {
		adminSet.Add(e)
	}

	for _, e := range clientKeys {
		clientSet.Add(e)
	}

	isInit = true
}
