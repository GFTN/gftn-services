// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/GFTN/gftn-services/payout-service/handler"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/op/go-logging"
	"github.com/urfave/negroni"
	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"
	"github.com/GFTN/gftn-services/utility"
	"github.com/GFTN/gftn-services/utility/global-environment/services"
	"github.com/GFTN/gftn-services/utility/logconfig"
	"github.com/GFTN/gftn-services/utility/message"
	middleware_checks "github.com/GFTN/gftn-services/utility/middleware"
	"github.com/GFTN/gftn-services/utility/response"
	"github.com/GFTN/gftn-services/utility/status"
	"github.com/GFTN/gftn-services/utility/wwfirebase"
)

type App struct {
	Router       *mux.Router
	serviceCheck status.ServiceCheck
	payoutOps    handler.PayoutPointOperations
	mwHandler    *middleware_checks.MiddlewareHandler
	HTTPHandler  func(http.Handler) http.Handler
}

var LOGGER = logging.MustGetLogger("payout-service")

func (a *App) Initialize() {
	services.VariableCheck()
	services.InitEnv()

	a.HTTPHandler = nil
	if os.Getenv(global_environment.ENV_KEY_ORIGIN_ALLOWED) == "true" {
		headersOk := handlers.AllowedHeaders([]string{"Access-Control-Allow-Headers", "Origin", "Content-Type", "X-Auth-Token", "Authorization"})
		originsOk := handlers.AllowedOrigins([]string{"*"})
		methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
		LOGGER.Infof("Setting up CORS")
		a.HTTPHandler = handlers.CORS(
			headersOk, originsOk, methodsOk)
	}

	serviceVersion = os.Getenv(global_environment.ENV_KEY_SERVICE_VERSION)

	// JWT environment variables dependency
	if os.Getenv(global_environment.ENV_KEY_ENABLE_JWT) != "false" {
		pepperObject := os.Getenv(global_environment.ENV_KEY_WW_JWT_PEPPER_OBJ)
		if pepperObject == "" {
			utility.ExitOnErr(LOGGER, errors.New("Pepper object must be set if jwt is enabled"), "Error in environment variable pepper object")
			return
		}
	}

	//Init firebase
	err := errors.New("")
	wwfirebase.FbClient, wwfirebase.FbAuthClient, err = wwfirebase.AuthenticateWithAdminPrivileges()
	wwfirebase.FbRef = wwfirebase.GetRootRef()

	if err != nil {
		LOGGER.Error("Error initializing firebase: %s", err.Error())
	}

	errorCodes := os.Getenv(global_environment.ENV_KEY_SERVICE_ERROR_CODES_FILE)
	err = message.LoadErrorConfig(errorCodes)
	utility.ExitOnErr(LOGGER, err, "Unable to set up error message config")

	LOGGER.Infof("Setting up service status check")
	a.serviceCheck, err = status.CreateServiceCheck()
	utility.ExitOnErr(LOGGER, err, "Unable to set up Service Check API")

	LOGGER.Infof("Setting up Payout Ops Internal API")
	a.payoutOps, err = handler.CreatePayoutPointOperations()
	utility.ExitOnErr(LOGGER, err, "Unable to set up Payout Ops API")

	LOGGER.Infof("Setting up router")

	a.mwHandler = middleware_checks.CreateMiddlewareHandler()
}

func (a *App) initializeRoutes() {

	a.Router = mux.NewRouter()
	// Code Block added by Operations team for debugging/testing http headers
	a.Router.HandleFunc("/"+serviceVersion+"/helloworldwire", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("BODY:", req.Body)
		// tester := {"test"}
		// response.Respond(w, http.StatusOK, JSON.Marshall(tester))
		type TestGroup struct {
			ID         int
			TestString string
			TestArray  []string
		}
		test := TestGroup{
			ID:         1,
			TestString: "Test",
			TestArray:  []string{"Value1", "Value2"},
		}
		payload, _ := json.Marshal(test)
		response.Respond(w, http.StatusOK, payload)

	}).Methods("POST")

	a.Router.NotFoundHandler = http.HandlerFunc(response.NotFound)

	// External & Admin API Service Endpoints

	LOGGER.Infof("\t* Admin API: Service Check")
	a.Router.Handle("/"+serviceVersion+"/client/service_check", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.serviceCheck.ServiceCheck),
	)).Methods("GET")

	/*
		Create a new payout location with json payload
	*/

	LOGGER.Infof("\t* Admin API: Create a new payout location with json payload")
	a.Router.Handle("/"+serviceVersion+"/admin/payout", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.payoutOps.AddPayout),
	)).Methods("POST")

	/*
		Create a new payout location with CSV file.
	*/

	LOGGER.Infof("\t* Admin API: Create a new payout location with CSV file")
	a.Router.Handle("/"+serviceVersion+"/admin/payout/csv", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.payoutOps.AddPayoutCSV),
	)).Methods("POST")

	/*
		To update the payout point location according to the entry id
	*/

	LOGGER.Infof("\t* Admin API: Update an entry that is already exists for a location")
	a.Router.Handle("/"+serviceVersion+"/admin/payout", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.payoutOps.UpdatePayout),
	)).Methods("PATCH")

	/*
		Returns the payout locations that meet the filters and criteria
	*/
	LOGGER.Infof("\t* External API: Get the payout locations that meet the filters and criteria")
	a.Router.Handle("/"+serviceVersion+"/client/payout", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.payoutOps.GetPayout),
	)).Methods("GET")

	/*
		To remove a payout location given the @id
	*/

	LOGGER.Infof("\t* Admin API: Delete a payout location provided by given id")
	a.Router.Handle("/"+serviceVersion+"/admin/payout", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.payoutOps.DeletePayout),
	)).Methods("DELETE")

}

var serviceVersion = ""

func main() {
	app := App{}
	app.Initialize()
	serviceLogs := os.Getenv(global_environment.ENV_KEY_SERVICE_LOG_FILE)
	servicePort := os.Getenv(global_environment.ENV_KEY_SERVICE_PORT)

	if serviceLogs == "" || servicePort == "" {
		utility.ExitOnErr(LOGGER, errors.New("Environment variables missing: SERVICE_LOG_FILE, SERVICE_PORT"), "Unable to initialize service")
	}

	f, err := logconfig.SetupLogging(serviceLogs, LOGGER)
	if err != nil {
		utility.ExitOnErr(LOGGER, err, "Unable to set up logging")
	}
	defer f.Close()

	app.initializeRoutes()

	var handler http.Handler = app.Router

	//if CORS is set
	if app.HTTPHandler != nil {
		handler = app.HTTPHandler(app.Router)
	}

	writeTimeout, _ := strconv.ParseInt(os.Getenv(global_environment.ENV_KEY_WRITE_TIMEOUT), 10, 64)
	readTimeout, _ := strconv.ParseInt(os.Getenv(global_environment.ENV_KEY_READ_TIMEOUT), 10, 64)
	idleTimeout, _ := strconv.ParseInt(os.Getenv(global_environment.ENV_KEY_IDLE_TIMEOUT), 10, 64)

	if writeTimeout == 0 || readTimeout == 0 || idleTimeout == 0 {
		panic("Service timeout should not be zero, please check if the environment variables WRITE_TIMEOUT, READ_TIMEOUT, IDLE_TIMEOUT are being set correctly")
	}

	LOGGER.Infof("Listening on external port:%s", servicePort)

	srv := &http.Server{
		Addr: ":" + servicePort,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * time.Duration(writeTimeout),
		ReadTimeout:  time.Second * time.Duration(readTimeout),
		IdleTimeout:  time.Second * time.Duration(idleTimeout),
		Handler:      handler, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		LOGGER.Error(srv.ListenAndServe().Error())
	}()

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*60, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s")
	flag.Parse()
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
	_ = srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	LOGGER.Errorf("shutting down")
	os.Exit(0)

}
