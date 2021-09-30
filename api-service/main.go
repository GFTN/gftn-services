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

	"github.com/GFTN/gftn-services/api-service/environment"

	"github.com/GFTN/gftn-services/api-service/fitoficct"
	"github.com/GFTN/gftn-services/api-service/sweeping"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	logging "github.com/op/go-logging"
	b "github.com/stellar/go/build"
	"github.com/urfave/negroni"
	"github.com/GFTN/gftn-services/api-service/assets"
	"github.com/GFTN/gftn-services/api-service/onboarding"
	"github.com/GFTN/gftn-services/api-service/participants"
	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"
	"github.com/GFTN/gftn-services/utility"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/global-environment/services"
	"github.com/GFTN/gftn-services/utility/logconfig"
	"github.com/GFTN/gftn-services/utility/message"
	middleware_checks "github.com/GFTN/gftn-services/utility/middleware"
	"github.com/GFTN/gftn-services/utility/response"
	"github.com/GFTN/gftn-services/utility/status"
	"github.com/GFTN/gftn-services/utility/wwfirebase"
)

type App struct {
	Router         *mux.Router
	InternalRouter *mux.Router
	serviceCheck   status.ServiceCheck
	assetOps       assets.AssetOperations
	participantOps participants.ParticipantOperations
	onboardingOps  onboarding.Operations
	clearingOp     fitoficct.FItoFICustomerCreditTransferOperations
	sweepingOps    sweeping.Operations
	mwHandler      *middleware_checks.MiddlewareHandler
	HTTPHandler    func(http.Handler) http.Handler
}

var LOGGER = logging.MustGetLogger("api-service")

func (a *App) Initialize() {

	services.VariableCheck()
	services.InitEnv()

	a.HTTPHandler = nil
	if os.Getenv(global_environment.ENV_KEY_ORIGIN_ALLOWED) == "true" {
		headersOk := handlers.AllowedHeaders([]string{"Access-Control-Allow-Headers", "Origin", "Content-Type", "X-Auth-Token", "Authorization"})
		originsOk := handlers.AllowedOrigins([]string{"*"})
		methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
		LOGGER.Infof("Setting up CORS")
		a.HTTPHandler = handlers.CORS(
			headersOk, originsOk, methodsOk)
	}

	serviceVersion = os.Getenv(global_environment.ENV_KEY_SERVICE_VERSION)

	//Set Defaults for stellar network
	b.DefaultNetwork.Passphrase = os.Getenv(global_environment.ENV_KEY_STELLAR_NETWORK)

	if _, exists := os.LookupEnv(environment.ENV_KEY_TRANSACTION_BATCH_LIMIT); !exists {
		panic("Environment variable TRANSACTION_BATCH_LIMIT is empty")
		return
	}

	errorCodes := os.Getenv(global_environment.ENV_KEY_SERVICE_ERROR_CODES_FILE)
	err := message.LoadErrorConfig(errorCodes)
	utility.ExitOnErr(LOGGER, err, "Unable to set up error message config")

	LOGGER.Infof("Setting up service status check")
	a.serviceCheck, err = status.CreateServiceCheck()
	utility.ExitOnErr(LOGGER, err, "Unable to set up Service Check API")

	LOGGER.Infof("Setting up Asset Ops API")
	a.assetOps, err = assets.CreateAssetOperations()
	utility.ExitOnErr(LOGGER, err, "Unable to set up Asset Ops API")

	// Participant operations
	LOGGER.Infof("Setting up Participant Ops API")
	a.participantOps, err = participants.CreateParticipantOperations()
	utility.ExitOnErr(LOGGER, err, "Unable to set up Participant Registry API")

	// Onboarding operations
	LOGGER.Infof("Setting up Onboarding Ops API")
	a.onboardingOps, err = onboarding.CreateOnboardingOperations()
	utility.ExitOnErr(LOGGER, err, "Unable to set up Onboarding Ops API")

	LOGGER.Infof("Setting up FI to FI Customer Credit Transfer Clearing External API")
	a.clearingOp, err = fitoficct.CreateFItoFICustomerCreditTransferOperation()
	utility.ExitOnErr(LOGGER, err, "Unable to set up FI to FI Customer Credit Transfer External API")

	LOGGER.Infof("Setting up Sweeping Ops API")
	a.sweepingOps, err = sweeping.CreateSweepingOperations()
	utility.ExitOnErr(LOGGER, err, "Unable to set up Sweeping Ops API")

	// Create middleware handler
	a.mwHandler = middleware_checks.CreateMiddlewareHandler()
}

func (a *App) initializeRoutes() {

	LOGGER.Infof("Setting up router")
	a.Router = mux.NewRouter()
	a.InternalRouter = mux.NewRouter()
	internalApiRoutes := mux.NewRouter()

	// Code Block added by Operations team for debugging/testing http headers
	/*a.Router.HandleFunc("/"+serviceVersion+"/helloworldwire", func(w http.ResponseWriter, req *http.Request) {
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
	*/

	a.Router.NotFoundHandler = http.HandlerFunc(response.NotFound)

	// External & Internal API Service Endpoints

	LOGGER.Infof("\t* Internal API:  Service Check")
	a.Router.Handle("/"+serviceVersion+"/client"+"/service_check", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.serviceCheck.ServiceCheck),
	)).Methods(http.MethodGet)

	/*
		Clearing and Settlement Endpoints
	*/

	LOGGER.Infof("\t* External API:  FI Transaction Details Processor")
	a.Router.Handle("/"+serviceVersion+"/client"+"/transactions", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.clearingOp.ProcessTxnDetails),
	)).Methods(http.MethodGet)

	/*
		Asset Endpoints
	*/
	/*
			API Security Vulnerability
		 	39846: Missing Authorization Header (Medium)
		 	add RequireAccountTokenAuthentication
			change to use Handle not HandleFun
	*/
	LOGGER.Infof("\t* External API:  issue asset")
	a.Router.Handle("/"+serviceVersion+"/client"+"/assets", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.assetOps.IssueAsset),
	)).Methods(http.MethodPost)

	/*


		/*
			API Security Vulnerability
		 	39846: Missing Authorization Header (Medium)
		 	add RequireAccountTokenAuthentication
			change to use Handle not HandleFun
	*/

	//Taking out fund account endpoint as it is added in anchor service
	/*LOGGER.Infof("\t* External API:  fund account")
	a.Router.Handle("/"+serviceVersion+"/client"+"/account/asset/fund",
		assetOps.FundAccountDeprecated).Methods(http.MethodPost)
	*/

	LOGGER.Infof("\t* External API:  query asset balance")
	a.Router.Handle("/"+serviceVersion+"/client"+"/balances/accounts/{account_name}", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.assetOps.AssetBalance),
	)).Methods(http.MethodGet)

	LOGGER.Infof("\t* External API:  query all accounts with non-zero balances for a given asset")
	a.Router.Handle("/"+serviceVersion+"/client"+"/obligations", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.assetOps.AssetBalances),
	)).Methods(http.MethodGet)

	LOGGER.Infof("\t* External API:  get all assets issued by this participant on World Wire")
	a.Router.Handle("/"+serviceVersion+"/client"+"/assets/issued", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.assetOps.IssuedAssets),
	)).Methods(http.MethodGet)

	LOGGER.Infof("\t* External API:  Assets issued by IBM account")
	a.Router.Handle("/"+serviceVersion+"/client"+"/assets", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.assetOps.WorldWireAssets),
	)).Methods(http.MethodGet)

	LOGGER.Infof("\t* External API:  query all accounts with outstanding balances for a given asset")
	a.Router.Handle("/"+serviceVersion+"/client"+"/obligations/{asset_code}", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.assetOps.GetOutstandingBalance),
	)).Methods(http.MethodGet)

	//Taking out this endpoint as issue #307
	//LOGGER.Infof("\t* External API:  get trusted assets for Issuing account")
	//a.Router.HandleFunc("/"+serviceVersion+"/client/issuingaccount/assets", assetOps.TrustedAssetsForIA).Methods(http.MethodGet)

	LOGGER.Infof("\t* External API:  get trusted assets for Operating account")
	a.Router.Handle("/"+serviceVersion+"/client"+"/assets/accounts/{account_name}", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.assetOps.TrustedAssetsForAccount),
	)).Methods(http.MethodGet)

	/*
		Change/Allow Trust Endpoints
	*/
	LOGGER.Infof("\t* External API:  change DO trust")
	a.Router.Handle("/"+serviceVersion+"/client"+"/trust", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.assetOps.CreateOrAllowTrust),
	)).Methods(http.MethodPost)
	/*
		Onboarding Endpoints
		change to use Handle not HandleFun
	*/
	LOGGER.Infof("\t* Onboarding API:  create issuing account")
	a.Router.Handle("/"+serviceVersion+"/onboarding"+"/accounts/{account_name}", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.onboardingOps.CreateAccount),
	)).Methods(http.MethodPost)

	LOGGER.Infof("\t* Admin API:  create issuing account")
	a.Router.Handle("/"+serviceVersion+"/admin"+"/accounts/{account_name}", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.onboardingOps.CreateAccount),
	)).Methods(http.MethodPost)

	LOGGER.Infof("\t* Onboarding API:  get participant account")
	a.Router.Handle("/"+serviceVersion+"/onboarding"+"/accounts/{account_name}", negroni.New(
		negroni.HandlerFunc(middlewares.SuperAuthorization),
		negroni.WrapFunc(a.onboardingOps.GetOperatingAccount),
	)).Methods(http.MethodGet)

	/*
		internal Endpoints
		to be available for WW services only
	*/
	LOGGER.Infof("\t* Internal API:  create issuing account")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal"+"/accounts/{account_name}",
		a.onboardingOps.CreateAccount).Methods(http.MethodPost)

	LOGGER.Infof("\t* Internal API:  get participant account")
	internalApiRoutes.HandleFunc("/"+serviceVersion+"/internal"+"/accounts/{account_name}", a.onboardingOps.GetOperatingAccount).Methods(http.MethodGet)

	/*
		Participant Endpoints
	*/

	LOGGER.Infof("\t* External API:  get participants on WW using query")
	a.Router.Handle("/"+serviceVersion+"/client"+"/participants", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.participantOps.GetParticipantByQuery),
	)).Methods(http.MethodGet)

	LOGGER.Infof("\t* External API:  get participants on WW using query")
	a.Router.Handle("/"+serviceVersion+"/client"+"/participants/{participant_id}", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.participantOps.GetParticipantByDomain),
	)).Methods(http.MethodGet)

	LOGGER.Infof("\t* External API:  get assets of the participant from participant's domain")
	a.Router.Handle("/"+serviceVersion+"/client"+"/assets/participants/{participant_id}", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.participantOps.GetAssetsForParticipant),
	)).Methods(http.MethodGet)

	/*
		account endpoints

	*/

	LOGGER.Infof("\t* External API:  get participant account")
	a.Router.Handle("/"+serviceVersion+"/client"+"/accounts/{account_name}", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.participantOps.GetAccount),
	)).Methods(http.MethodGet)

	LOGGER.Infof("\t* External API:  get participant list of accounts")
	a.Router.Handle("/"+serviceVersion+"/client"+"/accounts", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.participantOps.GetAccounts),
	)).Methods(http.MethodGet)

	//add router for internal endpoints and these endpoints don't need authorization
	a.InternalRouter.PathPrefix("/" + serviceVersion + "/internal").Handler(negroni.New(
		// set middleware on a group of routes:
		negroni.Wrap(internalApiRoutes),
	))

	/* sweeping endpoint
	 */
	LOGGER.Infof("\t* External API:  accounts sweeping")
	a.Router.Handle("/"+serviceVersion+"/client"+"/accounts/{account_name}/sweep", negroni.New(
		negroni.HandlerFunc(middlewares.ParticipantAuthorization),
		negroni.HandlerFunc(a.mwHandler.ParticipantStatusCheck),
		negroni.WrapFunc(a.sweepingOps.Sweep),
	)).Methods(http.MethodPost)
}

var APP App

var serviceVersion = ""

func main() {
	APP = App{}
	APP.Initialize()
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

	APP.initializeRoutes()

	servicePort := os.Getenv(global_environment.ENV_KEY_SERVICE_PORT)
	serviceInternalPort := os.Getenv(global_environment.ENV_KEY_SERVICE_INTERNAL_PORT)

	var internalHandler http.Handler = APP.InternalRouter //internal Routes

	var handler http.Handler = APP.Router

	//if CORS is set
	if APP.HTTPHandler != nil {
		handler = APP.HTTPHandler(APP.Router)
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

	LOGGER.Infof("Listening on :%s, %s", servicePort, serviceInternalPort)

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
