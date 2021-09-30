// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/urfave/negroni"
	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/op/go-logging"
	"github.com/GFTN/gftn-services/participant-registry/environment"
	rr "github.com/GFTN/gftn-services/participant-registry/registry-responder"
	"github.com/GFTN/gftn-services/utility"
	"github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/global-environment/services"
	"github.com/GFTN/gftn-services/utility/logconfig"
	"github.com/GFTN/gftn-services/utility/message"
	"github.com/GFTN/gftn-services/utility/response"
	"github.com/GFTN/gftn-services/utility/status"
	"github.com/GFTN/gftn-services/utility/wwfirebase"
)

type App struct {
	Router           *mux.Router
	OnboardingRouter *mux.Router
	prAPI            rr.Operations
	serviceCheck     status.ServiceCheck
	HTTPHandler      func(http.Handler) http.Handler
}

var LOGGER = logging.MustGetLogger("participant-registry")
var a App
var serviceVersion = ""

func (a *App) InitApp() {

	a.HTTPHandler = nil
	if os.Getenv(global_environment.ENV_KEY_ORIGIN_ALLOWED) == "true" {
		headersOk := handlers.AllowedHeaders([]string{"Access-Control-Allow-Headers", "Origin", "Content-Type", "X-Auth-Token", "Authorization"})
		originsOk := handlers.AllowedOrigins([]string{"*"})
		methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
		LOGGER.Infof("Setting up CORS")
		a.HTTPHandler = handlers.CORS(
			headersOk, originsOk, methodsOk)
	}

	errorCodes := os.Getenv(global_environment.ENV_KEY_SERVICE_ERROR_CODES_FILE)
	/*Check Marx fix : Absolute paths without double dots are recommended */
	if !common.IsSafePath(errorCodes) {
		utility.ExitOnErr(LOGGER, errors.New("file path may be vulnerable"), "Unable to set up error message config")
	}
	err := message.LoadErrorConfig(errorCodes)
	utility.ExitOnErr(LOGGER, err, "Unable to set up error message config")

	serviceVersion = os.Getenv(global_environment.ENV_KEY_SERVICE_VERSION)
	servicePort := os.Getenv(global_environment.ENV_KEY_SERVICE_PORT)
	LOGGER.Infof("Setting up participant-registry to listen on %v", fmt.Sprintf(":%v", servicePort))
	// TODO:  For Market Maker API
	LOGGER.Infof("Setting Participant registry API")
	isUnitTest := os.Getenv(environment.ENV_KEY_IS_UNIT_TEST)
	a.prAPI, err = rr.CreateParticipantRegistryOperations()
	utility.ExitOnErr(LOGGER, err, "Unable to set up Participant Registry API")
	if isUnitTest == "true" {
		//This set up is only for unit tests
		//This starts with empty data in collection for unit test uses "test" DB
		//prAPI, err = rr.CreateParticipantRegistryOperationsForTest()
		utility.ExitOnErr(LOGGER, err, "Unable to set up Participant Registry API for Unit test")
	}

	a.serviceCheck, err = status.CreateServiceCheck()
}

func (a *App) initRoutes() {
	LOGGER.Infof("Setting up router")
	a.Router = mux.NewRouter()
	a.OnboardingRouter = mux.NewRouter()
	internalApiRoutes := mux.NewRouter()

	LOGGER.Infof("\t* API:  Service Check")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/service_check", a.serviceCheck.ServiceCheck).Methods("GET")

	LOGGER.Infof("\t* Get participants for given participant domain")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/pr/domain/{participant_id}", a.prAPI.GetParticipantDomain).Methods("GET")

	LOGGER.Infof("\t* Get participants for a given country")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/pr/country/{country_code}", a.prAPI.GetParticipantsByCountry).Methods("GET")

	LOGGER.Infof("\t* Get participants pub key for given account name")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/pr/account/{participant_id}/{account_name}",
		a.prAPI.GetParticipantDistAccount).Methods("GET")

	LOGGER.Infof("\t* Save participants pub key for given account name")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/pr/account/{participant_id}",
		a.prAPI.SaveParticipantDistAccount).Methods("POST")

	LOGGER.Infof("\t* Save participants pub key for issuing account ")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/pr/issuingaccount/{participant_id}",
		a.prAPI.SaveParticipantIssuingAccount).Methods("POST")

	LOGGER.Infof("\t* Get participant for given issuing account address ")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/pr/issuingaccount/{account_address}",
		a.prAPI.GetParticipantForIssuingAccount).Methods("GET")

	LOGGER.Infof("\t* Create and save new Participant")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/pr", a.prAPI.CreateParticipant).Methods("POST")

	LOGGER.Infof("\t* Get all Participants")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/pr", a.prAPI.GetParticipants).Methods("GET")

	LOGGER.Infof("\t* Update a Participant")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/pr/{participant_id}", a.prAPI.UpdateParticipant).Methods("PUT")

	LOGGER.Infof("\t* Update a Participant's status")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/pr/{participant_id}/status", a.prAPI.UpdateParticipantStatus).Methods("PUT")

	LOGGER.Infof("\t* Get participant by either issuing address or operating address")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/pr/account/{account_address}", a.prAPI.GetParticipantByAddress).Methods("GET")

	/* Onboarding apis*/

	LOGGER.Infof("\t* API:  Service Check")
	a.OnboardingRouter.Handle("/"+serviceVersion+"/admin/service_check", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.serviceCheck.ServiceCheck),
	)).Methods("GET")

	LOGGER.Infof("\t* Create and save new Participant")
	a.OnboardingRouter.Handle("/"+serviceVersion+"/admin/pr", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.prAPI.CreateParticipant),
	)).Methods("POST")

	LOGGER.Infof("\t* Get all Participants")
	a.OnboardingRouter.Handle("/"+serviceVersion+"/admin/pr", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.prAPI.GetParticipants),
	)).Methods("GET")

	LOGGER.Infof("\t* Update a Participant")
	a.OnboardingRouter.Handle("/"+serviceVersion+"/admin/pr/{participant_id}", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.prAPI.UpdateParticipant),
	)).Methods("PUT")

	LOGGER.Infof("\t* Update a Participant's status")
	a.OnboardingRouter.Handle("/"+serviceVersion+"/admin/pr/{participant_id}/status", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.prAPI.UpdateParticipantStatus),
	)).Methods("PUT")

	LOGGER.Infof("\t* Get participants for given participant domain")
	a.OnboardingRouter.Handle("/"+serviceVersion+"/admin/pr/domain/{participant_id}", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.prAPI.GetParticipantDomain),
	)).Methods("GET")

	//add router for internal endpoints and these endpoints don't need authorization
	a.Router.PathPrefix("/" + serviceVersion + "/internal").Handler(negroni.New(
		// set middleware on a group of routes:
		negroni.Wrap(internalApiRoutes),
	))

	a.Router.NotFoundHandler = http.HandlerFunc(response.NotFound)
	a.OnboardingRouter.NotFoundHandler = http.HandlerFunc(response.NotFound)
}

func main() {
	services.VariableCheck()
	services.InitEnv()
	a = App{}
	serviceLogs := os.Getenv(global_environment.ENV_KEY_SERVICE_LOG_FILE)
	if !common.IsSafePath(serviceLogs) {
		utility.ExitOnErr(LOGGER, errors.New("file path may be vulnerable to path traversal attacks"), "Error setting up logging")
		return
	}
	f, err := logconfig.SetupLogging(serviceLogs, LOGGER)
	defer f.Close()

	if err != nil {
		LOGGER.Error("Error setting up logging: ", err.Error())
	}

	// JWT environment variables dependency
	if os.Getenv(global_environment.ENV_KEY_ENABLE_JWT) != "false" {
		pepperObject := os.Getenv(global_environment.ENV_KEY_WW_JWT_PEPPER_OBJ)
		if pepperObject == "" {
			utility.ExitOnErr(LOGGER, errors.New("Pepper object must be set if jwt is enabled"), "Error in environment variable pepper object")
			return
		}
	}

	//Init firebase

	wwfirebase.FbClient, wwfirebase.FbAuthClient, err = wwfirebase.AuthenticateWithAdminPrivileges()
	wwfirebase.FbRef = wwfirebase.GetRootRef()

	if err != nil {
		LOGGER.Error("Error initializing firebase: %s", err.Error())
	}

	a.InitApp()
	a.initRoutes()

	servicePort := os.Getenv(global_environment.ENV_KEY_SERVICE_PORT)
	serviceInternalPort := os.Getenv(global_environment.ENV_KEY_SERVICE_INTERNAL_PORT)
	var handler http.Handler = a.OnboardingRouter // Onboarding
	var internalHandler http.Handler = a.Router   //internal

	//if CORS is set
	if a.HTTPHandler != nil {
		handler = a.HTTPHandler(a.OnboardingRouter)
	}

	writeTimeout, _ := strconv.ParseInt(os.Getenv(global_environment.ENV_KEY_WRITE_TIMEOUT), 10, 64)
	readTimeout, _ := strconv.ParseInt(os.Getenv(global_environment.ENV_KEY_READ_TIMEOUT), 10, 64)
	idleTimeout, _ := strconv.ParseInt(os.Getenv(global_environment.ENV_KEY_IDLE_TIMEOUT), 10, 64)

	if writeTimeout == 0 || readTimeout == 0 || idleTimeout == 0 {
		panic("Service timeout should not be zero, please check if the environment variables WRITE_TIMEOUT, READ_TIMEOUT, IDLE_TIMEOUT are being set correctly")
	}

	srv := &http.Server{
		Addr: ":" + servicePort,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * time.Duration(writeTimeout),
		ReadTimeout:  time.Second * time.Duration(readTimeout),
		IdleTimeout:  time.Second * time.Duration(idleTimeout),
		//TLSConfig:    &cfg,
		Handler: handler, // Pass our instance of gorilla/mux in.
	}

	intSrv := &http.Server{
		Addr: ":" + serviceInternalPort,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * time.Duration(writeTimeout),
		ReadTimeout:  time.Second * time.Duration(readTimeout),
		IdleTimeout:  time.Second * time.Duration(idleTimeout),
		//TLSConfig:    &cfg,
		Handler: internalHandler, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		LOGGER.Error(srv.ListenAndServe().Error())
	}()
	go func() {
		LOGGER.Error(intSrv.ListenAndServe().Error())
	}()

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s")
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
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	LOGGER.Errorf("shutting down")
	os.Exit(0)

}
