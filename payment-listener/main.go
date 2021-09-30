// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/op/go-logging"
	"github.com/urfave/negroni"
	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/payment-listener/listeners"
	"github.com/GFTN/gftn-services/utility"
	comn "github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/global-environment/services"
	"github.com/GFTN/gftn-services/utility/logconfig"
	"github.com/GFTN/gftn-services/utility/message"
	middleware_checks "github.com/GFTN/gftn-services/utility/middleware"
	"github.com/GFTN/gftn-services/utility/participant"
	"github.com/GFTN/gftn-services/utility/status"
	"github.com/GFTN/gftn-services/utility/vault/utils"
	"github.com/GFTN/gftn-services/utility/wwfirebase"
)

type App struct {
	Router              *mux.Router
	InternalRouter      *mux.Router
	HTTPHandler         func(http.Handler) http.Handler
	InternalHTTPHandler func(http.Handler) http.Handler
	ListenerOps         listeners.PaymentListenerOperation
	serviceCheck        status.ServiceCheck
	mwHandler           *middleware_checks.MiddlewareHandler
}

var LOGGER = logging.MustGetLogger("api-service")

func (a *App) Initialize() {

	a.HTTPHandler = nil
	a.InternalHTTPHandler = nil
	if os.Getenv(global_environment.ENV_KEY_ORIGIN_ALLOWED) == "true" {
		headersOk := handlers.AllowedHeaders([]string{"Access-Control-Allow-Headers", "Origin", "Content-Type", "X-Auth-Token", "Authorization"})
		originsOk := handlers.AllowedOrigins([]string{"*"})
		methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
		LOGGER.Infof("Setting up CORS")
		a.HTTPHandler = handlers.CORS(
			headersOk, originsOk, methodsOk)
		a.InternalHTTPHandler = handlers.CORS(
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

	a.ListenerOps = listeners.CreatePaymentListenerOperation()

	// get all operating accounts
	accounts, err := a.ListenerOps.GetParticipantOperatingAccounts()
	if err != nil {
		LOGGER.Warningf("Error getting Operating accounts")
		utility.ExitOnErr(LOGGER, err, "Error GetParticipantOperatingAccounts failed")
		return
	}

	issuingAccount, err := participant.GenericGetAccount(utils.Session{}, comn.ISSUING)
	//add issuing account to list of payment listeners
	if issuingAccount.NodeAddress != "" {
		issueAccount := model.Account{}
		issueAccount.Name = comn.ISSUING
		issueAccount.Address = &issuingAccount.NodeAddress
		accounts = append(accounts, issueAccount)
	}

	// create payment listener for each of my existing operating accounts
	/*
		start the http listener
	*/
	a.ListenerOps.CreatePaymentListeners(accounts)

	// Create middleware handler
	a.mwHandler = middleware_checks.CreateMiddlewareHandler()
}

func (a *App) initializeRoutes() {

	a.Router = mux.NewRouter()
	a.InternalRouter = mux.NewRouter()
	internalApiRoutes := mux.NewRouter()

	LOGGER.Infof("\t* Internal API:  Service Check")
	a.Router.Handle("/"+serviceVersion+"/client/service_check", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.serviceCheck.ServiceCheck),
	)).Methods("GET")

	/*
		Payment Listener
	*/

	LOGGER.Infof("\t* client Payment API:  Subscribe to Notification payment from an account")
	a.Router.Handle("/"+serviceVersion+"/client/accounts/{account_name}/{cursor}", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.ListenerOps.ReStartListener),
	)).Methods("POST")

	LOGGER.Infof("\t* internal Payment API:  Subscribe to Notification payment from an account")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/accounts/{account_name}/{cursor}", a.ListenerOps.ReStartListener).Methods("POST")

	//add router for internal endpoints and these endpoints don't need authorization
	a.InternalRouter.PathPrefix("/" + serviceVersion + "/internal").Handler(negroni.New(
		// set middleware on a group of routes:
		negroni.Wrap(internalApiRoutes),
	))

}

var APP App

var serviceVersion = ""

func main() {
	APP = App{}
	services.VariableCheck()
	services.InitEnv()

	serviceLogs := os.Getenv(global_environment.ENV_KEY_SERVICE_LOG_FILE)
	f, err := logconfig.SetupLogging(serviceLogs, LOGGER)
	if err != nil {
		utility.ExitOnErr(LOGGER, err, "Unable to set up logging")
	}
	defer f.Close()
	APP.Initialize()
	APP.initializeRoutes()

	servicePort := os.Getenv(global_environment.ENV_KEY_SERVICE_PORT)
	serviceInternalPort := os.Getenv(global_environment.ENV_KEY_SERVICE_INTERNAL_PORT)

	var handler http.Handler = APP.Router
	var internalHandler http.Handler = APP.InternalRouter

	//if CORS is set
	if APP.HTTPHandler != nil {
		handler = APP.HTTPHandler(APP.Router)
		internalHandler = APP.InternalHTTPHandler(APP.InternalRouter)
	}

	writeTimeout, _ := strconv.ParseInt(os.Getenv(global_environment.ENV_KEY_WRITE_TIMEOUT), 10, 64)
	readTimeout, _ := strconv.ParseInt(os.Getenv(global_environment.ENV_KEY_READ_TIMEOUT), 10, 64)
	idleTimeout, _ := strconv.ParseInt(os.Getenv(global_environment.ENV_KEY_IDLE_TIMEOUT), 10, 64)

	if writeTimeout == 0 || readTimeout == 0 || idleTimeout == 0 {
		panic("Service timeout should not be zero, please check if the environment variables WRITE_TIMEOUT, READ_TIMEOUT, IDLE_TIMEOUT are being set correctly")
	}

	LOGGER.Infof("Listening on :%s, internalport:%v", servicePort, serviceInternalPort)

	srv := &http.Server{
		Addr: ":" + servicePort,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * time.Duration(writeTimeout),
		ReadTimeout:  time.Second * time.Duration(readTimeout),
		IdleTimeout:  time.Second * time.Duration(idleTimeout),
		Handler:      handler, // Pass our instance of gorilla/mux in.
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
	_ = srv.Shutdown(ctx)
	_ = intSrv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	LOGGER.Errorf("shutting down")
	os.Exit(0)

}
