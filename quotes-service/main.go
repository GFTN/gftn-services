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

	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/response"

	"github.com/gorilla/mux"
	logging "github.com/op/go-logging"
	"github.com/stellar/go/clients/horizon"
	"github.com/urfave/negroni"
	gasserviceclient "github.com/GFTN/gftn-services/gas-service-client"
	"github.com/GFTN/gftn-services/global-whitelist-service/whitelistclient"
	"github.com/GFTN/gftn-services/quotes-service/environment"
	"github.com/GFTN/gftn-services/quotes-service/handler/exchangehandler"
	"github.com/GFTN/gftn-services/quotes-service/handler/quoteshandler"
	"github.com/GFTN/gftn-services/quotes-service/utility/cryptoservice"
	"github.com/GFTN/gftn-services/quotes-service/utility/nqsdbclient"
	"github.com/GFTN/gftn-services/quotes-service/utility/participantregistry"
	"github.com/GFTN/gftn-services/utility"
	"github.com/GFTN/gftn-services/utility/global-environment/services"
	"github.com/GFTN/gftn-services/utility/kafka"
	"github.com/GFTN/gftn-services/utility/logconfig"
	"github.com/GFTN/gftn-services/utility/message"
	middleware_checks "github.com/GFTN/gftn-services/utility/middleware"
	"github.com/GFTN/gftn-services/utility/status"
	"github.com/GFTN/gftn-services/utility/wwfirebase"
)

var LOGGER = logging.MustGetLogger("quotes-service")

var a App
var serviceVersion string

type App struct {
	Router          *mux.Router
	serviceCheck    status.ServiceCheck
	quoteHandler    quoteshandler.QuoteHandler
	exchangeHandler exchangehandler.ExchangeHandler
	mwHandler       *middleware_checks.MiddlewareHandler
}

func quoteHandlerBuilder() quoteshandler.QuoteHandler {
	LOGGER.Info("Using REST Participant Registry Service Client")
	prClient := participantregistry.Client{
		HTTP: &http.Client{Timeout: time.Second * 10},
		URL:  os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL),
	}
	LOGGER.Info("Using ** as horizon server")
	horizonClient := horizon.Client{
		URL:  os.Getenv(global_environment.ENV_KEY_HORIZON_CLIENT_URL),
		HTTP: &http.Client{Timeout: time.Second * 10}}
	LOGGER.Info("Using ** as White List Service")
	wlsClient := whitelistclient.Client{
		HTTPClient: &http.Client{Timeout: time.Second * 10},
		WLURL:      os.Getenv(global_environment.ENV_KEY_WL_SVC_URL),
	}
	dbPort, _ := strconv.Atoi(os.Getenv(environment.ENV_KEY_POSTGRESQLPORT))
	postgreDBC := nqsdbclient.PostgreDatabaseClient{
		Host:     os.Getenv(environment.ENV_KEY_POSTGRESQLHOST),
		Port:     dbPort,
		User:     os.Getenv(environment.ENV_KEY_POSTGRESQLUSER),
		Password: os.Getenv(environment.ENV_KEY_POSTGRESQLPASSWORD),
		Dbname:   os.Getenv(environment.ENV_KEY_POSTGRESQLDBNAME),
	}
	err := postgreDBC.CreateConnection() //initialize connection
	if err != nil {
		LOGGER.Error(err)
		LOGGER.Error("Error creating DB connection")
		os.Exit(0)
	}
	quoteHandler := quoteshandler.QuoteHandler{
		PRClient:      &prClient,
		HTTP:          &http.Client{Timeout: time.Second * 10},
		WLSClient:     &wlsClient,
		DBClient:      &postgreDBC,
		HorizonClient: &horizonClient,
	}
	LOGGER.Infof("Initiate Kafka producer for ww-gateway")
	quoteHandler.GatewayOperation, err = kafka.Initialize()
	if err != nil {
		LOGGER.Error(err)
		LOGGER.Error("Initialize Kafka producer for ww-gateway failed")
		os.Exit(0)
	}
	return quoteHandler
}

func exchangeHandlerBuilder() exchangehandler.ExchangeHandler {
	LOGGER.Infof("Using REST Participant Registry Service Client")
	prClient := participantregistry.Client{
		HTTP: &http.Client{Timeout: time.Second * 10},
		URL:  os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL),
	}
	LOGGER.Infof("Using ** as horizon server")
	horizonClient := horizon.Client{
		URL:  os.Getenv(global_environment.ENV_KEY_HORIZON_CLIENT_URL),
		HTTP: &http.Client{Timeout: time.Second * 10}}
	LOGGER.Infof("Using ** as White List Service")
	wlsClient := whitelistclient.Client{
		HTTPClient: &http.Client{Timeout: time.Second * 10},
		WLURL:      os.Getenv(global_environment.ENV_KEY_WL_SVC_URL),
	}
	dbPort, _ := strconv.Atoi(os.Getenv(environment.ENV_KEY_POSTGRESQLPORT))
	postgreDBC := nqsdbclient.PostgreDatabaseClient{
		Host:     os.Getenv(environment.ENV_KEY_POSTGRESQLHOST),
		Port:     dbPort,
		User:     os.Getenv(environment.ENV_KEY_POSTGRESQLUSER),
		Password: os.Getenv(environment.ENV_KEY_POSTGRESQLPASSWORD),
		Dbname:   os.Getenv(environment.ENV_KEY_POSTGRESQLDBNAME),
	}
	err := postgreDBC.CreateConnection() //initialize connection
	if err != nil {
		LOGGER.Error(err)
		LOGGER.Error("Error creating DB connection")
		os.Exit(0)
	}
	gasServiceClient := gasserviceclient.Client{
		HTTP: &http.Client{Timeout: time.Second * 20},
		URL:  os.Getenv(global_environment.ENV_KEY_GAS_SVC_URL),
	}
	csClient := cryptoservice.Client{
		HTTP:        &http.Client{Timeout: time.Second * 10},
		URLTemplate: os.Getenv(global_environment.ENV_KEY_CRYPTO_SVC_INTERNAL_URL),
	}
	exchangeHandler := exchangehandler.ExchangeHandler{
		HTTP:             &http.Client{Timeout: time.Second * 10},
		GasServiceClient: &gasServiceClient,
		HorizonClient:    &horizonClient,
		CSClient:         &csClient,
		PRClient:         &prClient,
		WLSClient:        &wlsClient,
		DBClient:         &postgreDBC,
	}
	return exchangeHandler
}

func (a *App) initialize() {
	services.VariableCheck()
	services.InitEnv()
	errorCodes := os.Getenv(global_environment.ENV_KEY_SERVICE_ERROR_CODES_FILE)
	err := message.LoadErrorConfig(errorCodes)
	utility.ExitOnErr(LOGGER, err, "Unable to set up error message config")
	serviceVersion = os.Getenv(global_environment.ENV_KEY_SERVICE_VERSION)

	servicePort := os.Getenv(global_environment.ENV_KEY_SERVICE_PORT)

	LOGGER.Infof("Setting up Quotes Service to listen on: %v", fmt.Sprintf(":%v", servicePort))
	LOGGER.Infof("Quotes Service Version:  %v", serviceVersion)
	LOGGER.Infof("Setting up service status check")
	a.serviceCheck, err = status.CreateServiceCheck()
	utility.ExitOnErr(LOGGER, err, "Unable to set up Service Check API")

	// get QuoteAPI
	LOGGER.Infof("Setting up Quotes Handler...")
	a.quoteHandler = quoteHandlerBuilder()
	LOGGER.Infof("Setting up Exchange Handler...")
	a.exchangeHandler = exchangeHandlerBuilder()

	LOGGER.Infof("Setting up middleware")
	// Create middleware handler
	a.mwHandler = middleware_checks.CreateMiddlewareHandler()
}

func (a *App) initializeRoutes() {

	// initialize router
	LOGGER.Infof("Initialize router...")
	a.Router = mux.NewRouter()

	a.Router.NotFoundHandler = http.HandlerFunc(response.NotFound)

	LOGGER.Infof("\t* Internal API:  Service Check")
	a.Router.Handle("/"+serviceVersion+"/client/service_check", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.serviceCheck.ServiceCheck),
	)).Methods("GET")

	LOGGER.Infof("\t* Protocol API:  Quote API for OFI to request quote and have request_id in return")
	url := "/" + serviceVersion + "/client/quotes/request"
	a.Router.Handle(url, negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.quoteHandler.RequestQuote),
	)).Methods("post")

	LOGGER.Infof("\t* Protocol API:  Quote API for OFI to get quotes results with request_id")
	url = "/" + serviceVersion + "/client/quotes/request/{request_id}"
	a.Router.Handle(url, negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.quoteHandler.GetQuotes),
	)).Methods("GET")

	LOGGER.Infof("\t* Protocol API:  Quote API for RFI to get quotes regarding quote_id")
	url = "/" + serviceVersion + "/client/quotes/{quote_id}"
	a.Router.Handle(url, negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.quoteHandler.GetQuotesByQuoteID),
	)).Methods("GET")

	LOGGER.Infof("\t* Protocol API:  Quote API for RFI to post quotes regarding quote_id")
	url = "/" + serviceVersion + "/client/quotes/{quote_id}"
	a.Router.Handle(url, negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.quoteHandler.UpdateQuote),
	)).Methods("POST")

	LOGGER.Infof("\t* Protocol API:  Quote API for RFI to cancel quote regrading quote_id, rfi_Domain")
	url = "/" + serviceVersion + "/client/quotes/{quote_id}"
	a.Router.Handle(url, negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.quoteHandler.CancelQuote),
	)).Methods("Delete")

	LOGGER.Infof("\t* Protocol API:  Quote API for OFI to get quotes results by attributes defined in request body")
	url = "/" + serviceVersion + "/client/quotes"
	a.Router.Handle(url, negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.quoteHandler.GetQuotesByAttributes),
	)).Methods("GET")

	LOGGER.Infof("\t* Protocol API:  Quote API for OFI to get quotes results by attributes defined in request body (POST Method")
	url = "/" + serviceVersion + "/client/quotes"
	a.Router.Handle(url, negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.quoteHandler.GetQuotesByAttributes),
	)).Methods("POST")

	LOGGER.Infof("\t* Protocol API:  Quote API for RFI to cancel quote by attributes defined in request body")
	url = "/" + serviceVersion + "/client/quotes"
	a.Router.Handle(url, negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.quoteHandler.CancelQuotesByAttributes),
	)).Methods("Delete")

	LOGGER.Infof("\t* Protocol API:  Quote API for RFI to cancel quote by attributes defined in request body")
	url = "/" + serviceVersion + "/client/exchange"
	a.Router.Handle(url, negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.exchangeHandler.CreateAtomicExchange),
	)).Methods("POST")

}

func main() {
	a = App{}
	a.initialize()
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

	servicePort := os.Getenv(global_environment.ENV_KEY_SERVICE_PORT)
	a.initializeRoutes()
	LOGGER.Infof("Setting up Middleware...")
	handler := a.Router

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
		Handler:      handler, // Pass our instance of gorilla/mux in.
	}
	LOGGER.Infof("Listening on :%s", servicePort)
	// Run our server in a goroutine so that it doesn't block.
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
