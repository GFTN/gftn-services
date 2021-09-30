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

	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"

	"github.com/GFTN/gftn-services/utility/wwfirebase"

	global_environment "github.com/GFTN/gftn-services/utility/global-environment"

	"github.com/gorilla/mux"
	logging "github.com/op/go-logging"
	"github.com/urfave/negroni"

	"github.com/GFTN/gftn-services/global-whitelist-service/whitelistserver/database"
	"github.com/GFTN/gftn-services/global-whitelist-service/whitelistserver/environment"
	"github.com/GFTN/gftn-services/global-whitelist-service/whitelistserver/handler"
	"github.com/GFTN/gftn-services/global-whitelist-service/whitelistserver/utility/prclient"
	"github.com/GFTN/gftn-services/utility/status"

	"github.com/GFTN/gftn-services/utility"
	"github.com/GFTN/gftn-services/utility/global-environment/services"
	"github.com/GFTN/gftn-services/utility/logconfig"
	"github.com/GFTN/gftn-services/utility/message"
	middleware_checks "github.com/GFTN/gftn-services/utility/middleware"
	"github.com/GFTN/gftn-services/utility/response"
)

type App struct {
	Router         *mux.Router
	InternalRouter *mux.Router
	wlh            handler.WhitelistHandler
	mwHandler      *middleware_checks.MiddlewareHandler
	serviceCheck   status.ServiceCheck
}

func whitelistHandlerBuilder() handler.WhitelistHandler {
	dc := database.DynamoClient{
		AWS_SECRET_KEY: os.Getenv(global_environment.ENV_KEY_AWS_SECRET_ACCESS_KEY),
		AWS_KEY_ID:     os.Getenv(global_environment.ENV_KEY_AWS_ACCESS_KEY_ID),
		TableName:      os.Getenv(environment.ENV_KEY_WL_SVC_TABLE_NAME),
		Region:         os.Getenv(environment.ENV_KEY_GWL_DYNAMO_DB_REGION),
	}
	err := dc.CreateConnection()
	if err != nil {
		LOGGER.Error("Error establishing dynamo DB connection")
		LOGGER.Error(err)
	}
	wl, err := dc.GetWhiteListParicipants("hk")
	LOGGER.Info(wl)
	prc := prclient.Client{
		HTTPClient: &http.Client{Timeout: time.Second * 10},
		URL:        os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL),
	}
	wlh := handler.WhitelistHandler{
		DBClient: &dc,
		PRClient: &prc,
	}
	return wlh
}

func MiddleWareBuilder() handler.MiddleWare {
	return handler.MiddleWare{}
}

var LOGGER = logging.MustGetLogger("whilelistservice")
var serviceVersion = ""

func (a *App) InitApp() {

	LOGGER.Infof("Setting up error codes")
	errorCodes := os.Getenv(global_environment.ENV_KEY_SERVICE_ERROR_CODES_FILE)
	err := message.LoadErrorConfig(errorCodes)
	utility.ExitOnErr(LOGGER, err, "Unable to set up error message config")

	LOGGER.Infof("Setting up handler")
	a.wlh = whitelistHandlerBuilder()
	serviceVersion = os.Getenv(global_environment.ENV_KEY_SERVICE_VERSION)

	a.serviceCheck, _ = status.CreateServiceCheck()

	// Create middleware handler
	a.mwHandler = middleware_checks.CreateMiddlewareHandler()

}

func (a *App) initRoutes() {

	// Initatilize a.Router
	a.Router = mux.NewRouter()
	a.InternalRouter = mux.NewRouter()

	LOGGER.Infof("Setting up client API router...")
	a.Router.NotFoundHandler = http.HandlerFunc(response.NotFound)

	LOGGER.Infof("\t* API:  Service Check")
	a.Router.Handle("/"+serviceVersion+"/client/service_check", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.serviceCheck.ServiceCheck),
	)).Methods("GET")

	LOGGER.Infof("\t* Get wlparticipantIDs for given participant ID")
	a.Router.Handle("/"+serviceVersion+"/client/participants/whitelist", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.wlh.GetWLParticipantIDsClient),
	)).Methods("GET")

	LOGGER.Infof("\t* Get wlparticipants for a given participant ID")
	a.Router.Handle("/"+serviceVersion+"/client/participants/whitelist/object", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.wlh.GetWLParticipantsClient),
	)).Methods("GET")

	LOGGER.Infof("\t* Create wlparticipants(body) for a given participant ID")
	a.Router.Handle("/"+serviceVersion+"/client/participants/whitelist", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.wlh.CreateWLParticipantClient),
	)).Methods("POST")

	LOGGER.Infof("\t* Delete wlparticipants(body) for a given participant ID")
	a.Router.Handle("/"+serviceVersion+"/client/participants/whitelist", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.wlh.DeleteWLParticipantClient),
	)).Methods("Delete")

	// a.Router.PathPrefix("/" + serviceVersion + "/").Handler(negroni.New(
	// 	negroni.Wrap(clientAPIRouter),
	// ))
	LOGGER.Info("Setting up internal API router...")
	internalAPIRouter := mux.NewRouter()
	internalAPIRouter.NotFoundHandler = http.HandlerFunc(response.NotFound)

	LOGGER.Infof("\t* Get wlparticipantIDs for given participant ID")
	internalAPIRouter.HandleFunc("/"+serviceVersion+"/internal/participants/whitelist/{participant_id}", a.wlh.GetWLParticipantIDs).Methods("GET")

	LOGGER.Infof("\t* Get wlparticipants for a given participant ID")
	internalAPIRouter.HandleFunc("/"+serviceVersion+"/internal/participants/whitelist/{participant_id}/object", a.wlh.GetWLParticipants).Methods("GET")

	LOGGER.Infof("\t* Create wlparticipants(body) for a given participant ID")
	internalAPIRouter.HandleFunc("/"+serviceVersion+"/internal/participants/whitelist/{participant_id}", a.wlh.CreateWLParticipant).Methods("POST")

	LOGGER.Infof("\t* Delete wlparticipants(body) for a given participant ID")
	internalAPIRouter.HandleFunc("/"+serviceVersion+"/internal/participants/whitelist/{participant_id}", a.wlh.DeleteWLParticipant).Methods("Delete")

	LOGGER.Infof("\t* Get Mutual wlparticipantIDs(body) for a given participant ID")
	internalAPIRouter.HandleFunc("/"+serviceVersion+"/internal/participants/whitelist/{participant_id}/mutual", a.wlh.GetMutualWLParticipantIDs).Methods("GET")

	LOGGER.Infof("\t* Get Mutual wlparticipants(body) for a given participant ID")
	internalAPIRouter.HandleFunc("/"+serviceVersion+"/internal/participants/whitelist/{participant_id}/mutual/object", a.wlh.GetMutualWLParticipants).Methods("GET")

	//add router for internal endpoints and these endpoints don't need authorization
	a.InternalRouter.PathPrefix("/" + serviceVersion + "/internal").Handler(negroni.New(
		negroni.Wrap(internalAPIRouter),
	))
}

func main() {
	app := App{}
	services.VariableCheck()
	services.InitEnv()
	serviceLogs := os.Getenv(global_environment.ENV_KEY_SERVICE_LOG_FILE)
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

	app.InitApp()
	app.initRoutes()

	servicePort := os.Getenv(global_environment.ENV_KEY_SERVICE_PORT)
	serviceInternalPort := os.Getenv(global_environment.ENV_KEY_SERVICE_INTERNAL_PORT)
	var clientHandler http.Handler = app.Router
	var internalHandler http.Handler = app.InternalRouter

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
		Handler:      clientHandler, // Pass our instance of gorilla/mux in.
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
