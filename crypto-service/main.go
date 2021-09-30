// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package main

import (
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	global_environment "github.com/GFTN/gftn-services/utility/global-environment"

	"github.com/gorilla/mux"
	"github.com/op/go-logging"
	b "github.com/stellar/go/build"
	"github.com/urfave/negroni"
	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"
	crypto_handler "github.com/GFTN/gftn-services/crypto-service/crypto-handler"
	"github.com/GFTN/gftn-services/utility"
	"github.com/GFTN/gftn-services/utility/global-environment/services"
	"github.com/GFTN/gftn-services/utility/logconfig"
	"github.com/GFTN/gftn-services/utility/message"
	middleware_checks "github.com/GFTN/gftn-services/utility/middleware"
	"github.com/GFTN/gftn-services/utility/status"
	"github.com/GFTN/gftn-services/utility/wwfirebase"
	"golang.org/x/net/context"
)

type App struct {
	Router         *mux.Router
	InternalRouter *mux.Router
	serviceCheck   status.ServiceCheck
	cryptoHandler  crypto_handler.CryptoOperations
	mwHandler      *middleware_checks.MiddlewareHandler
}

var LOGGER = logging.MustGetLogger("crypto-service")
var serviceVersion string

func (a *App) Initialize() {

	services.VariableCheck()
	services.InitEnv()
	a.Router = mux.NewRouter()
	a.InternalRouter = mux.NewRouter()

	serviceVersion = os.Getenv(global_environment.ENV_KEY_SERVICE_VERSION)
	networkPassphrase := os.Getenv(global_environment.ENV_KEY_STELLAR_NETWORK)
	//Set Defaults for stellar network
	b.DefaultNetwork.Passphrase = networkPassphrase

	errorCodes := os.Getenv(global_environment.ENV_KEY_SERVICE_ERROR_CODES_FILE)
	err := message.LoadErrorConfig(errorCodes)
	utility.ExitOnErr(LOGGER, err, "Unable to set up error message config")

	LOGGER.Infof("Setting up service status check")
	a.serviceCheck, err = status.CreateServiceCheck()
	utility.ExitOnErr(LOGGER, err, "Unable to set up Service Check API")

	a.cryptoHandler, err = crypto_handler.CreateCryptoOperations()
	utility.ExitOnErr(LOGGER, err, "Unable to set up Account Finder Internal API")

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

	LOGGER.Infof("\t* Internal API:  signing endpoint")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/sign", a.cryptoHandler.SignXdr).Methods("POST")

	LOGGER.Infof("\t* Internal API:  signing endpoint for participant with out verification")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/participant/sign", a.cryptoHandler.ParticipantSignXdr).Methods("POST")

	LOGGER.Infof("\t* client API:  sign payload endpoint")
	a.Router.Handle("/"+serviceVersion+"/client/sign", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.cryptoHandler.SignPayload),
	)).Methods("POST")

	LOGGER.Infof("\t* client API:  Sign ISO20022 XML")
	a.Router.Handle("/"+serviceVersion+"/client/payload/sign", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.cryptoHandler.SignXML),
	)).Methods("POST")

	LOGGER.Infof("\t* Internal API:  sign payload endpoint, does same operation without JWT tokens check")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/payload/wwsign", a.cryptoHandler.SignXMLUsingStellar).Methods("POST")

	LOGGER.Infof("\t* Internal API:  sign payload endpoint, does same operation without JWT tokens check")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/request/sign", a.cryptoHandler.SignPayload).Methods("POST")

	LOGGER.Infof("\t* internal API:  create account endpoint")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/account/{account_name}", a.cryptoHandler.CreateAccount).Methods("POST")

	LOGGER.Infof("\t* internal API:  Get IBM admin account endpoint")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/admin/account", a.cryptoHandler.GetIBMAccount).Methods("GET")

	LOGGER.Infof("\t* internal API:sign with IBM admin account endpoint")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal/admin/sign", a.cryptoHandler.AddIBMSign).Methods("POST")

	//add router for internal endpoints and these endpoints don't need authorization
	a.InternalRouter.PathPrefix("/" + serviceVersion + "/internal").Handler(negroni.New(
		// set middleware on a group of routes:
		negroni.Wrap(internalApiRoutes),
	))

}

func main() {
	app := App{}
	app.Initialize()

	serviceLogs := os.Getenv(global_environment.ENV_KEY_SERVICE_LOG_FILE)
	f, err := logconfig.SetupLogging(serviceLogs, LOGGER)
	if err != nil {
		utility.ExitOnErr(LOGGER, err, "Unable to set up logging")
	}
	defer f.Close()

	// JWT environment variables dependency
	if os.Getenv(global_environment.ENV_KEY_ENABLE_JWT) != "false" {
		pepperObject := os.Getenv(global_environment.ENV_KEY_WW_JWT_PEPPER_OBJ)
		if pepperObject == "" {
			utility.ExitOnErr(LOGGER, errors.New("Pepper object must be set if jwt is enabled"), "Error in environment variable pepper object")
			return
		}
	}

	wwfirebase.FbClient, wwfirebase.FbAuthClient, err = wwfirebase.AuthenticateWithAdminPrivileges()
	wwfirebase.FbRef = wwfirebase.GetRootRef()

	if err != nil {
		LOGGER.Error("Error initializing firebase: %s", err.Error())
	}

	app.initializeRoutes()

	servicePort := os.Getenv(global_environment.ENV_KEY_SERVICE_PORT)
	serviceInternalPort := os.Getenv(global_environment.ENV_KEY_SERVICE_INTERNAL_PORT)
	var handler http.Handler = app.Router
	var internalHandler http.Handler = app.InternalRouter

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
	_ = srv.Shutdown(ctx)
	_ = intSrv.Shutdown(ctx)

	//Close Crypto Session
	_ = crypto_handler.CYPTO_OPERATIONS.HSMInstance.C.Logout(crypto_handler.CYPTO_OPERATIONS.HSMInstance.Session)
	_ = crypto_handler.CYPTO_OPERATIONS.HSMInstance.C.CloseSession(crypto_handler.CYPTO_OPERATIONS.HSMInstance.Session)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	LOGGER.Errorf("shutting down")
	os.Exit(0)

}
