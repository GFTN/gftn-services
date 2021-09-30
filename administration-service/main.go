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
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	logging "github.com/op/go-logging"
	b "github.com/stellar/go/build"
	"github.com/urfave/negroni"
	"github.com/GFTN/gftn-services/administration-service/blocklist"
	"github.com/GFTN/gftn-services/administration-service/killswitch"
	rr "github.com/GFTN/gftn-services/administration-service/persistence"
	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"
	"github.com/GFTN/gftn-services/utility"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/global-environment/services"
	"github.com/GFTN/gftn-services/utility/logconfig"
	"github.com/GFTN/gftn-services/utility/message"
	"github.com/GFTN/gftn-services/utility/response"
	"github.com/GFTN/gftn-services/utility/status"
	"github.com/GFTN/gftn-services/utility/wwfirebase"
)

type App struct {
	Router       *mux.Router
	AdminRouter  *mux.Router
	po           rr.MongoDBOperations
	serviceCheck status.ServiceCheck
	killSwitch   killswitch.KillSwitch
	blocklistOps blocklist.BlocklistOperations
	HTTPHandler  func(http.Handler) http.Handler
}

var LOGGER = logging.MustGetLogger("administration-service")
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

	//Set Defaults for stellar network
	b.DefaultNetwork.Passphrase = os.Getenv(global_environment.ENV_KEY_STELLAR_NETWORK)

	errorCodes := os.Getenv(global_environment.ENV_KEY_SERVICE_ERROR_CODES_FILE)
	err := message.LoadErrorConfig(errorCodes)
	utility.ExitOnErr(LOGGER, err, "Unable to set up error message config")

	serviceVersion = os.Getenv(global_environment.ENV_KEY_SERVICE_VERSION)
	servicePort := os.Getenv(global_environment.ENV_KEY_SERVICE_INTERNAL_PORT)
	LOGGER.Infof("Setting up administration service to listen on %v", fmt.Sprintf(":%v", servicePort))

	LOGGER.Infof("Setting Admin service API")
	a.po, err = rr.CreateAdminServicePersistenceOperations()
	utility.ExitOnErr(LOGGER, err, "Unable to set up Administration Service API")

	a.serviceCheck, err = status.CreateServiceCheck()

	LOGGER.Info("Setting up Kill Switch API")
	a.killSwitch, err = killswitch.CreateKillSwitch()
	utility.ExitOnErr(LOGGER, err, "Unable to setup Kill Switch API")

	LOGGER.Infof("Setting Blocklist API")
	a.blocklistOps, err = blocklist.CreateBlocklistOperations()
	utility.ExitOnErr(LOGGER, err, "Unable to set up blocklist API")

}

func (a *App) initRoutes() {

	LOGGER.Infof("Setting up router")
	a.Router = mux.NewRouter()
	a.AdminRouter = mux.NewRouter()

	internalAPIRoutes := mux.NewRouter()

	internalAPIRoutes.NotFoundHandler = http.HandlerFunc(response.NotFound)

	LOGGER.Infof("\t* API:  Service Check")
	internalAPIRoutes.HandleFunc("/"+serviceVersion+"/gftn/service_check", a.serviceCheck.ServiceCheck).Methods("GET")

	LOGGER.Infof("\t* internal API:  persist FiToFiCCTMemo")
	internalAPIRoutes.HandleFunc("/"+serviceVersion+"/internal/fitoficct", a.po.StoreFiToFiCCTMemo).Methods("POST")

	LOGGER.Infof("\t* internal API:  query Transaction details")
	internalAPIRoutes.HandleFunc("/"+serviceVersion+"/internal/transaction", a.po.GetTxnDetails).Methods("POST")

	LOGGER.Infof("\t* internal API:  Suspend Account")
	internalAPIRoutes.HandleFunc("/"+serviceVersion+"/internal/suspend/{participant_id}/{account_name}", a.killSwitch.SuspendAccount).Methods("POST")

	LOGGER.Infof("\t* internal API:  Re activate Account")
	internalAPIRoutes.HandleFunc("/"+serviceVersion+"/internal/reactivate/{participant_id}/{account_name}", a.killSwitch.ReactivateAccount).Methods("POST")

	LOGGER.Infof("\t* internal API:  Check if target currency/institution/country is in the blocklist")
	internalAPIRoutes.HandleFunc("/"+serviceVersion+"/internal/blocklist/validate", a.blocklistOps.Validate).Methods("POST")

	LOGGER.Infof("\t* internal API:  Add new currency/institution/country into the blocklist")
	internalAPIRoutes.HandleFunc("/"+serviceVersion+"/internal/blocklist", a.blocklistOps.Add).Methods("POST")

	LOGGER.Infof("\t* internal API:  Remove certain currency/institution/country from the blocklist")
	internalAPIRoutes.HandleFunc("/"+serviceVersion+"/internal/blocklist", a.blocklistOps.Remove).Methods("DELETE")

	LOGGER.Infof("\t* internal API:  Get blocklist")
	internalAPIRoutes.HandleFunc("/"+serviceVersion+"/internal/blocklist", a.blocklistOps.Get).Methods("GET")

	LOGGER.Infof("\t* admin API:  Check if target currency/institution/country is in the blocklist")
	internalAPIRoutes.HandleFunc("/"+serviceVersion+"/internal/blocklist/validate", a.blocklistOps.Validate).Methods("POST")

	LOGGER.Infof("\t* admin API:  Add new currency/institution/country into the blocklist")
	a.AdminRouter.Handle("/"+serviceVersion+"/admin/blocklist", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.blocklistOps.Add),
	)).Methods("POST")

	LOGGER.Infof("\t* admin API:  Remove certain currency/institution/country from the blocklist")
	a.AdminRouter.Handle("/"+serviceVersion+"/admin/blocklist", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.blocklistOps.Remove),
	)).Methods("DELETE")

	LOGGER.Infof("\t* admin API:  Get blocklist")
	a.AdminRouter.Handle("/"+serviceVersion+"/admin/blocklist", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.blocklistOps.Get),
	)).Methods("GET")

	LOGGER.Infof("\t* internal API:  query Transaction details")
	a.AdminRouter.Handle("/"+serviceVersion+"/admin/transaction", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.po.GetTxnDetails),
	)).Methods("POST")

	LOGGER.Infof("\t* internal API:  Suspend Account")
	a.AdminRouter.Handle("/"+serviceVersion+"/admin/suspend/{participant_id}/{account_name}", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.killSwitch.SuspendAccount),
	)).Methods("POST")

	LOGGER.Infof("\t* internal API:  Re activate Account")
	a.AdminRouter.Handle("/"+serviceVersion+"/admin/reactivate/{participant_id}/{account_name}", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.killSwitch.ReactivateAccount),
	)).Methods("POST")

	//add router for internal endpoints and these endpoints don't need authorization
	a.Router.PathPrefix("/" + serviceVersion + "/internal").Handler(negroni.New(
		// set middleware on a group of routes:
		negroni.Wrap(internalAPIRoutes),
	))

}

func main() {

	services.VariableCheck()
	services.InitEnv()
	a = App{}
	serviceLogs := os.Getenv(global_environment.ENV_KEY_SERVICE_LOG_FILE)
	f, err := logconfig.SetupLogging(serviceLogs, LOGGER)
	defer f.Close()

	if err != nil {
		LOGGER.Error("Error setting up logging: ", err.Error())
	}

	errorCodes := os.Getenv(global_environment.ENV_KEY_SERVICE_ERROR_CODES_FILE)
	err = message.LoadErrorConfig(errorCodes)
	utility.ExitOnErr(LOGGER, err, "Unable to set up error message config")

	servicePort := os.Getenv(global_environment.ENV_KEY_SERVICE_PORT)
	serviceInternalPort := os.Getenv(global_environment.ENV_KEY_SERVICE_INTERNAL_PORT)

	// initiate firebase
	wwfirebase.FbClient, wwfirebase.FbAuthClient, err = wwfirebase.AuthenticateWithAdminPrivileges()
	wwfirebase.FbRef = wwfirebase.GetRootRef()

	if err != nil {
		LOGGER.Error("Error initializing firebase: %s", err.Error())
	}

	a.InitApp()
	a.initRoutes()
	LOGGER.Infof("Listening on :%s, internal :%s", servicePort, serviceInternalPort)

	var adminHandler http.Handler = a.AdminRouter // Admin endpoints to connect UI
	var internalHandler http.Handler = a.Router   //internal

	//if CORS is set
	if a.HTTPHandler != nil {
		adminHandler = a.HTTPHandler(a.AdminRouter)
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
		Handler: adminHandler, // Pass our instance of gorilla/mux in.
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

	//This low on gas jobs will happen in gas service instead
	//batch.KickOffBatchJobs()

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
