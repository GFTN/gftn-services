// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"
	"github.com/GFTN/gftn-services/automation-service/automate/participant"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/op/go-logging"
	"github.com/urfave/negroni"
	"github.com/GFTN/gftn-services/automation-service/automate"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/message"
	"github.com/GFTN/gftn-services/utility/wwfirebase"
)

type App struct {
	Router          *mux.Router
	serviceCheckOps automate.ServiceCheck
	automationOps   participant.DeploymentOperations
	HTTPHandler     func(http.Handler) http.Handler
}

var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02T15:04:05Z07:00} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

/* sets up logging directing it to the given log file */
func SetupLogging(serviceLogs string, LOGGER *logging.Logger) (*os.File, error) {
	LOGGER.Infof("Log File: %s", serviceLogs)
	f, err := os.OpenFile(serviceLogs, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v: %v", err, serviceLogs)
		return nil, err
	}
	logWriter := io.MultiWriter(f, os.Stdout)

	backend1 := logging.NewLogBackend(logWriter, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend1, format)
	logging.SetBackend(backend2Formatter)
	return f, nil
}

var LOGGER = logging.MustGetLogger("automation-service")

func (a *App) InitApp() {

	//TODO get service version from env variable
	//serviceVersion = os.Getenv(global_environment.ENV_KEY_SERVICE_VERSION)
	errorCodes := os.Getenv(global_environment.ENV_KEY_SERVICE_ERROR_CODES_FILE)
	message.LoadErrorConfig(errorCodes)

	if os.Getenv(global_environment.ENV_KEY_ORIGIN_ALLOWED) == "true" {
		headersOk := handlers.AllowedHeaders([]string{"Access-Control-Allow-Headers", "Origin", "Content-Type", "X-Auth-Token", "Authorization", "X-Fid", "X-Verify-Token"})
		originsOk := handlers.AllowedOrigins([]string{"*"})
		methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
		LOGGER.Infof("Setting up CORS")
		a.HTTPHandler = handlers.CORS(
			headersOk, originsOk, methodsOk)
	}

	LOGGER.Infof("Setting Service Check API")
	serviceCheckOps, serviceCheckErr := automate.InitiateServiceCheck()
	if serviceCheckErr != nil {
		LOGGER.Errorf("Unable to set up service check API:  %v", serviceCheckErr.Error())
		os.Exit(1)
	}

	a.serviceCheckOps = serviceCheckOps

	LOGGER.Infof("Setting Automate Deployment API")
	automationOps, automateErr := participant.InitiateDeploymentOperations()
	if automateErr != nil {
		LOGGER.Errorf("Unable to set up automation API:  %v", automateErr.Error())
		os.Exit(1)
	}

	a.automationOps = automationOps

}

func (a *App) initRoutes() {

	a.Router = mux.NewRouter()

	serviceVersion := "v1"
	//TODO get service version from env variable
	//serviceVersion = os.Getenv(global_environment.ENV_KEY_SERVICE_VERSION)
	LOGGER.Infof("\t* internal API: 	Service check")
	a.Router.HandleFunc("/"+serviceVersion+"/check", a.serviceCheckOps.Check).Methods(http.MethodGet)

	// Deploy new participant
	LOGGER.Infof("\t* internal API: 	Deploy participant configurations and services")
	a.Router.Handle("/"+serviceVersion+"/deploy/participant", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.automationOps.DeployParticipantServicesAndConfigs),
	)).Methods(http.MethodPost)
}

func main() {
	// Set up logging file
	serviceLogs := os.Getenv(global_environment.ENV_KEY_SERVICE_LOG_FILE)
	f, err := SetupLogging(serviceLogs, LOGGER)
	if err != nil {
		LOGGER.Errorf("Unable to set up logging: %s", err.Error())
		return
	}
	defer f.Close()

	// Set up FireBase connection
	wwfirebase.FbClient, wwfirebase.FbAuthClient, err = wwfirebase.AuthenticateWithAdminPrivileges()
	wwfirebase.FbRef = wwfirebase.GetRootRef()

	APP := App{}
	APP.InitApp()
	APP.initRoutes()

	servicePort := "5566"
	LOGGER.Infof("Listening on :%s", servicePort)

	var handler http.Handler
	handler = APP.Router

	if APP.HTTPHandler != nil {
		handler = APP.HTTPHandler(APP.Router)
	}

	srv := &http.Server{
		Addr: ":" + servicePort,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 180,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 15,
		Handler:      handler, // Pass our instance of gorilla/mux in.
	}

	go func() {
		LOGGER.Error(srv.ListenAndServe().Error())
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
