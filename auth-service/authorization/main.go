// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	logging "github.com/op/go-logging"
	"github.com/urfave/negroni"
	authutility "github.com/GFTN/gftn-services/auth-service/authorization/authutility"
	"github.com/GFTN/gftn-services/auth-service/authorization/helper"
	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"
	"github.com/GFTN/gftn-services/utility/response"
	"github.com/GFTN/gftn-services/utility/wwfirebase"
)

// TODO: Permissions in firebase instead of constants file (decouple logic from values)
// TODO: SuperAuthorization : Make Participant ID and Institution ID optional. The function should not return errors for them not being there.

// App : type App
type App struct {
	Router *mux.Router
	//InternalRouter      *mux.Router
	HTTPHandler func(http.Handler) http.Handler
	//InternalHTTPHandler func(http.Handler) http.Handler
}

// LOGGER : logs middleware package
var LOGGER = logging.MustGetLogger("middlewares")

// Initialize app to add a new router and an HTTP handler
func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.HTTPHandler = nil
}

func (a *App) initializeRoutes() {

	// initialize gorilla mux router
	routes := mux.NewRouter()

	a.Router.NotFoundHandler = http.HandlerFunc(response.NotFound)

	// initialize public routes (routes without JWT Auth Middleware):
	routes.HandleFunc("/auth/service_check", authutility.ServiceCheck).Methods(http.MethodGet)

	//define middleware to be called for **all** routes:
	routes.Use(middlewares.LogURI)

	// initialize negroni for routing with middleware
	// negroni.Classic() provides some default middleware that is useful for most applications:
	// - negroni.Recovery - Panic Recovery Middleware.
	// - negroni.Logger - Request/Response Logger Middleware.
	// - negroni.Static - Static File serving under the "public" directory.
	n := negroni.Classic()

	/*
	* This route is for demoing participant authorization for Worldwire client portal
	* POST method.
	* The parameters are explained in ParticipantAuthorization in the middleware package.
	 */
	routes.Handle("/auth/participant-auth-test",
		negroni.New(

			negroni.HandlerFunc(middlewares.ParticipantAuthorization),

			// set endpoint (final logic to call)
			negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
				authutility.ServiceCheck(w, r)
			}))).
		Methods("POST")

	/*
	* This route is for demoing super user authorization for Worldwire client portal
	* POST method.
	* The parameters are explained in SuperUserAuthorization in the middleware package.
	 */
	routes.Handle("/auth/super-auth-test",
		negroni.New(

			negroni.HandlerFunc(middlewares.ParticipantAuthorization),

			// set endpoint (final logic to call)
			negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
				authutility.ServiceCheck(w, r)
			}))).
		Methods("POST")

	n.UseHandler(routes)

	a.Router.PathPrefix("/auth/").Handler(negroni.New(
		// set group of routes for test
		negroni.Wrap(routes),
	))

}

// APP : Declare a global APP of the type App
var APP App

func main() {

	APP = App{}
	APP.Initialize()

	// Following is example usage for Nakul to use
	// // should deny if enpoint is not named correctly
	// approved, err := middlewares.CheckAccess("Super_permissions", "manager", false, "GET", "/v1/admin/p")
	// // should succeed - super
	// approved, err := middlewares.CheckAccess("Super_permissions", "manager", false, "GET", "/v1/admin/pr")
	// should succeed - participant
	// approved, err := middlewares.CheckAccess("Participant_permissions", "admin", false, "GET", "/v1/anchor/address")
	// should succeed - participant + maker/checker
	// approved, err := middlewares.CheckAccess("Participant_permissions", "manager", true, "POST", "/v1/anchor/fundings/instruction")
	// // should not succeed, wrong method - participant + maker/checker
	// approved, err := middlewares.CheckAccess("Participant_permissions", "manager", true, "PUT", "/v1/anchor/fundings/instruction")
	// // should not succeed - insufficient permissions
	// approved, err := middlewares.CheckAccess("Participant_permissions", "manager", false, "POST", "/v1/anchor/fundings/instruction")
	// should succeed using delete - participant + maker/checker
	// approved, err := middlewares.CheckAccess("Participant_permissions", "admin", false, "DELETE", "/v1/client/participants/whitelist")
	// if err != nil && approved == false {
	// 	fmt.Println(err)
	// 	fmt.Println("Access denied!")
	// } else {
	// 	fmt.Println("Access approved!")
	// }

	// set in launch.json for debugging purposes
	credentialsDir, ok := os.LookupEnv("cred")
	if ok {
		fmt.Printf("overriding env vars from " + credentialsDir)
		helper.SetCustomEnvs(credentialsDir)
	} else {
		fmt.Printf("ok, please note cred env (.credentials/{env}) not set for debugging")
	}

	// start firebase app (admin sdk)
	FbClient, FbAuthClient, err := wwfirebase.AuthenticateWithAdminPrivileges()
	wwfirebase.FbClient = FbClient
	wwfirebase.FbAuthClient = FbAuthClient
	wwfirebase.FbRef = wwfirebase.GetRootRef()

	if err != nil {
		LOGGER.Error("Error initializing firebase: %s", err.Error())
	}

	// start server with graceful shutdown
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Initialize routes
	// Making this change in order for us to be able to test the application better
	APP.initializeRoutes()
	var handler http.Handler = APP.Router

	// if CORS is set
	// this can be added the same way as API service. Only requires copy paste of the code.
	// TODO : Check if CORS needs to be added here.
	if APP.HTTPHandler != nil {
		handler = APP.HTTPHandler(APP.Router)
	}

	// create server
	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler, // Pass our instance of mux/gorilla
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			LOGGER.Error(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	LOGGER.Info("Shutting down")
	os.Exit(0)
}
