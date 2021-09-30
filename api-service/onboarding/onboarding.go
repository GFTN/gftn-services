// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package onboarding

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	b "github.com/stellar/go/build"
	"github.com/GFTN/gftn-services/api-service/client"
	"github.com/GFTN/gftn-services/api-service/environment"
	"github.com/GFTN/gftn-services/api-service/fitoficct"
	crypto_client "github.com/GFTN/gftn-services/crypto-service-client/crypto-client"
	gasserviceclient "github.com/GFTN/gftn-services/gas-service-client"
	"github.com/GFTN/gftn-services/gftn-models/model"
	pr_client "github.com/GFTN/gftn-services/participant-registry-client/pr-client"
	"github.com/GFTN/gftn-services/utility"
	ast "github.com/GFTN/gftn-services/utility/asset"
	"github.com/GFTN/gftn-services/utility/common"
	comn "github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/participant"
	"github.com/GFTN/gftn-services/utility/response"
	vauth "github.com/GFTN/gftn-services/utility/vault/auth"
	"github.com/GFTN/gftn-services/utility/vault/utils"
)

// Operations - init struct
type Operations struct {
	ParticipantRegistryClient pr_client.PRServiceClient
	CryptoServiceClient       crypto_client.CryptoServiceClient
	VaultSession              utils.Session
	GasServiceClient          gasserviceclient.GasServiceClient
}

//CreateOnboardingOperations - Init opeartions
func CreateOnboardingOperations() (Operations, error) {

	op := Operations{}
	var prClient pr_client.PRServiceClient
	if os.Getenv(environment.ENV_KEY_PARTICIPANT_REGISTRY_SERVICE_CLIENT) == "mock" {
		LOGGER.Warningf("USING MOCK PARTICIPANT REGISTRY SERVICE CLIENT")
		prClient = pr_client.MockPRServiceClient{}
	} else {
		LOGGER.Infof("Using REST Participant Registry Service Client")
		cl, err := pr_client.CreateRestPRServiceClient(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL))
		utility.ExitOnErr(LOGGER, err, "Unable to create REST Participant Registry Service Client")
		prClient = cl
	}
	op.ParticipantRegistryClient = prClient
	var cClient crypto_client.CryptoServiceClient
	if os.Getenv(environment.ENV_KEY_PARTICIPANT_REGISTRY_SERVICE_CLIENT) == "mock" {
		cClient, _ = crypto_client.CreateMockCryptoServiceClient()
	} else {
		err := errors.New("")
		cServiceInternalURL, err := participant.GetServiceUrl(os.Getenv(global_environment.ENV_KEY_CRYPTO_SVC_INTERNAL_URL), os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME))
		if err != nil {
			return op, err
		}
		cClient, err = crypto_client.CreateRestCryptoServiceClient(cServiceInternalURL)
		if err != nil {
			return op, err
		}
	}
	op.CryptoServiceClient = cClient

	op.VaultSession = utils.Session{}

	if os.Getenv(global_environment.ENV_KEY_SECRET_STORAGE_LOCATION) == common.VAULT_SECRET {
		//Vault location
		err := errors.New("")
		op.VaultSession, err = vauth.GetSession()
		if err != nil {
			LOGGER.Errorf("Error reading account source environment settings")
			return op, err
		}
	}

	gasServiceClient := gasserviceclient.Client{
		HTTP: &http.Client{Timeout: time.Second * 20},
		URL:  os.Getenv(global_environment.ENV_KEY_GAS_SVC_URL),
	}
	op.GasServiceClient = &gasServiceClient

	return op, nil

	//PostParticipantDistAccount(domain string, account model.OperatingAccount) (string, error)

}

//CreateIssuingAccount - Create issuing account
func (op Operations) CreateIssuingAccount(w http.ResponseWriter, req *http.Request) {
	/*
		1. check if there is an issuing account, if yes, return issuing account address
		2. create new issuing account (in both stellar and NodeConfig, return the new issuing account address
	*/
	LOGGER.Infof("Creating Issuing Account")
	homeDomain := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)

	issuingAccount, err := op.ParticipantRegistryClient.GetParticipantIssuingAccount(homeDomain)
	if err == nil && issuingAccount != "" {
		LOGGER.Errorf("Issuing account already exists in PR, Account address: %v", issuingAccount)
		response.NotifyWWError(w, req, http.StatusBadRequest, "API-1277", errors.New(comn.ISSUING))
		return
	}
	if err != nil {
		LOGGER.Errorf("Error connecting to PR, Account address: %v", err)
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1119", err)
		return
	}

	randString := participant.GetSecretPhrase()
	//Call new signing service
	account, errorMsg, status, code := op.CryptoServiceClient.CreateAccount(comn.ISSUING)

	if status != http.StatusCreated || errorMsg != nil {
		LOGGER.Errorf("Error creating new account: %s", errorMsg)
		code = "API-1090"
		response.NotifyWWError(w, req, status, code, errorMsg)
		return
	}

	LOGGER.Debugf("Created issuing account: %v", account.NodeAddress)
	// create issuing account in stellar

	err = op.createIssuingAccountInStellar(account.NodeAddress)
	if err != nil {
		LOGGER.Errorf("Error creating new account on stellar: %v", err)
		response.NotifyWWError(w, req, http.StatusConflict, "API-1092", err)
		return
	}

	modelAccount := model.Account{
		Name:    comn.ISSUING,
		Address: &account.NodeAddress,
	}

	LOGGER.Infof("Now storing account: %s to secret manager", modelAccount.Address)
	err = participant.GenericStoreAccount(op.VaultSession, comn.ISSUING, account, randString)
	if err != nil {
		LOGGER.Debugf("Error  %v", err.Error())
		response.NotifyWWError(w, req, http.StatusConflict, "API-1092", err)
		return
	}

	// setOptions for the issuing account
	//This is done in the background thread to keep the endpoint responsive
	//Will need to investigate logs for failures in these steps
	go func() {

		LOGGER.Infof("Setting options for issuing account %s", modelAccount.Address)
		err = op.setOptionsForIssuingAccount(&modelAccount, randString)
		if err != nil {
			LOGGER.Errorf("Error Creating set options for account on stellar, Account address: %v, %v", modelAccount.Address, err)
			response.NotifyWWError(w, req, http.StatusConflict, "API-1094", err)
			return
		}

		LOGGER.Infof("Storing account to participant registry")
		err = op.ParticipantRegistryClient.PostParticipantIssuingAccount(homeDomain, modelAccount)
		if err != nil {
			LOGGER.Errorf("Error adding issuing account to PR, Account address: %v, %v", modelAccount.Address, err)
			response.NotifyWWError(w, req, http.StatusNotFound, "API-1094", err)
			return
		}

		LOGGER.Infof("Starting payment listener for issuing account")
		//start payment listener for issuing account
		paymentLClient, err := client.CreateRestPaymentListenerClient()
		if err != nil {
			LOGGER.Errorf("CreateRestPaymentListenerClient Error: %v", err.Error())
			response.NotifyWWError(w, req, http.StatusNotFound, "API-1123", err)
			return
		}
		err = paymentLClient.SubscribePayments(comn.ISSUING)
		if err != nil {
			//Account is created but was error starting payment listener
			LOGGER.Errorf("CreateRestPaymentListenerClient Error: %v", err.Error())
			response.NotifyWWError(w, req, http.StatusCreated, "API-1123", errors.New(comn.ISSUING))
			return
		}

	}()

	LOGGER.Infof("Issuing account creation complete!")
	acrJSON, _ := json.Marshal(modelAccount)
	response.Respond(w, http.StatusOK, acrJSON)

}

//CreateAccount - Create generic account handler
func (op Operations) CreateAccount(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	accountName := vars["account_name"]
	LOGGER.Infof("Creating Account: %s", accountName)

	if accountName == comn.ISSUING {
		// invoke issuing account creation
		op.CreateIssuingAccount(w, req)
		return
	}
	//invoke operating account creation
	op.CreateOperatingAccount(w, req)

}

//CreateOperatingAccount - create operating account handler
func (op Operations) CreateOperatingAccount(w http.ResponseWriter, req *http.Request) {

	/*
		1. check if there is a operating account with the given name, if yes, return operating account address
		2. get the issuing account from participant registry. issuing account will create and fund the operating account
		3. create new operating account (in both stellar, participant registry, & NodeConfig, return the new operating account address
		4. use SetOptions to set the flags, weights and thresholds for the DA
	*/

	vars := mux.Vars(req)
	accountName := vars["account_name"]
	homeDomain := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	if accountName == "" {
		response.NotifyWWError(w, req, http.StatusConflict, "API-1100", errors.New("account name missing"))
		return
	}
	LOGGER.Infof("Creating Operating Account: %s", accountName)

	opAccount, err := op.ParticipantRegistryClient.GetParticipantDistAccount(homeDomain, accountName)
	if err == nil && opAccount != "" {
		LOGGER.Errorf("Operating account already exists in PR, Account address: %v", accountName)
		response.NotifyWWError(w, req, http.StatusBadRequest, "API-1277", errors.New(accountName))
		return
	}
	// if err != nil {
	// 	LOGGER.Errorf("Error connecting to PR, Account address: %v", err)
	// 	response.NotifyWWError(w, req, http.StatusNotFound, "API-1119", err)
	// 	return
	// }

	randString := participant.GetSecretPhrase()

	//step1: Call new signing service to create new operating account
	account, errorMsg, status, code := op.CryptoServiceClient.CreateAccount(accountName)

	if status != http.StatusCreated || errorMsg != nil {
		LOGGER.Errorf("Error creating new account: %v", errorMsg.Error())
		code = "API-1090"
		response.NotifyWWError(w, req, status, code, errorMsg)
		return
	}

	//step 2: get the issuing account from participant registry. issuing account will create and fund the operating account
	issuingAccount, err := op.ParticipantRegistryClient.GetParticipantIssuingAccount(homeDomain)
	if err != nil {
		LOGGER.Errorf("No Issuing account. Cannot create Operating account", err)
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1097", err)
		return
	}

	err = op.createNewOperatingAccountInStellar(issuingAccount, account.NodeAddress)
	if err != nil {
		LOGGER.Errorf("SetOptions error for new operating account: %v", err.Error())
		response.NotifyWWError(w, req, http.StatusConflict, "API-1092", err)
		return
	}

	modelAccount := model.Account{
		Name:    accountName,
		Address: &account.NodeAddress,
	}

	LOGGER.Infof("Now storing operating account: %s to secret manager", modelAccount.Address)
	err = participant.GenericStoreAccount(op.VaultSession, accountName, account, randString)
	if err != nil {
		LOGGER.Debugf("Error  %v", err.Error())
		response.NotifyWWError(w, req, http.StatusConflict, "API-1092", err)
		return
	}

	//This is done in the background thread to keep the endpoint responsive
	//Will need to investigate logs for failures in these steps
	go func() {

		LOGGER.Infof("Setting options for operating account %s", modelAccount.Address)
		// set options for the Operating account
		err = op.setOptionsForOperatingAccount(modelAccount, randString)
		if err != nil {
			LOGGER.Errorf("SetOptions error for new issuing account: %v", err.Error())
			response.NotifyWWError(w, req, http.StatusNotFound, "API-1092", err)
			return
		}

		LOGGER.Infof("Storing account to participant registry")
		//register operating account in participant registry
		err = op.ParticipantRegistryClient.PostParticipantDistAccount(homeDomain, modelAccount)
		if err != nil {
			LOGGER.Errorf("Error adding Operating account to PR, Account Name: %v, %v", accountName, err)
			response.NotifyWWError(w, req, http.StatusNotFound, "API-1113", err)
			return
		}

		LOGGER.Infof("Starting payment listener for operating account: %s", modelAccount.Name)
		//start payment listener
		paymentLClient, err := client.CreateRestPaymentListenerClient()
		if err != nil {
			LOGGER.Errorf("CreateRestPaymentListenerClient Error: %v", err.Error())
			response.NotifyWWError(w, req, http.StatusNotFound, "API-1123", err)
			return
		}
		err = paymentLClient.SubscribePayments(accountName)
		if err != nil {
			//Account is created but was error starting payment listener
			LOGGER.Errorf("CreateRestPaymentListenerClient Error: %v", err.Error())
			response.NotifyWWError(w, req, http.StatusCreated, "API-1123", errors.New(accountName))
			return
		}

	}()

	LOGGER.Infof("Operating account: %s creation complete!", modelAccount.Name)
	//return account created response
	acrJSON, _ := json.Marshal(modelAccount)
	response.Respond(w, http.StatusOK, acrJSON)
}

//GetOperatingAccount - get account function
func (op Operations) GetOperatingAccount(w http.ResponseWriter, req *http.Request) {

	/*
		1. check if there is a Operating account with the given name, if yes, return Operating account address
		2. create new Operating account (in both stellar and NodeConfig, return the new Operating account address
	*/

	vars := mux.Vars(req)
	accountName := vars["account_name"]
	if accountName == "" {
		response.NotifyWWError(w, req, http.StatusConflict, "API-1100", errors.New("account name missing"))
		return
	}
	LOGGER.Infof("Getting Operating Account: %s", accountName)

	account, err := participant.GenericGetAccount(op.VaultSession, accountName)
	if err != nil {
		LOGGER.Errorf("Error while checking if Operating account exists: %v", err)
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1095", err)
		return
	}
	if account.NodeAddress != "" {
		//handle account exists
		LOGGER.Infof("Operating Account Exists - name: %v, address: %v", accountName, account.NodeAddress)
		acr := model.AuthAccount{Account: &model.Account{Address: &account.NodeAddress, Name: accountName}}
		acrJSON, _ := json.Marshal(acr)
		response.Respond(w, http.StatusOK, acrJSON)
		return
	}
	response.NotifyWWError(w, req, http.StatusNotFound, "API-1098", err)
}

// newKeyPair generates a new private/public keypair
/*func newKeyPair() (*keypair.Full, error) {
	pair, err := keypair.Random()
	if err != nil {
		LOGGER.Errorf("Error Creating a new keypair: %v", err)
	}
	return pair, err
}*/

// createIssuingAccountInStellar - creates a new issuing account for the given addr
// use the ww_admin_key to create this account instead of friendbot

func (op Operations) createIssuingAccountInStellar(addr string) error {

	LOGGER.Infof("Creating issuing account in stellar: %s", addr)
	stellarNetwork := comn.GetStellarNetwork(os.Getenv(global_environment.ENV_KEY_STELLAR_NETWORK))
	ibmAccount, sequenceNum, err := op.GasServiceClient.GetAccountAndSequence()
	if ibmAccount == "" || sequenceNum == 0 || err != nil {
		LOGGER.Debugf("IBM Gas account failed to load")
		return errors.New("Error getting IBM Gas account")
	}

	//Get IBM account
	wwAdminAccount, errorMsg, status, _ := op.CryptoServiceClient.GetIBMAccount()
	if errorMsg != nil {
		LOGGER.Errorf("Error getting IBM account: %v", errorMsg.Error(), wwAdminAccount)
		return errorMsg
	}
	if status != http.StatusOK {
		LOGGER.Errorf("createIssuingAccountInStellar: Error getting IBM account: status %v", status)
		return errors.New("Error creating new account in stellar")
	}

	//Initial fund is sourced from IBM admin account
	initialFund := os.Getenv(environment.ENV_KEY_ACCOUNT_INITIAL_FUND)
	tx, err := b.Transaction(
		b.SourceAccount{AddressOrSeed: ibmAccount},
		b.Sequence{Sequence: sequenceNum},
		stellarNetwork,
		b.CreateAccount(
			//initial fund is sourced from ibmadmin account
			b.SourceAccount{AddressOrSeed: *wwAdminAccount.Address},
			b.Destination{AddressOrSeed: addr},
			b.NativeAmount{Amount: initialFund},
		),
	)

	if err != nil {
		return err
	}

	var txeb b.TransactionEnvelopeBuilder
	err = txeb.Mutate(tx)
	txeB64, err := txeb.Base64()

	//Get IBM signature from crypto service
	xdrB, _ := base64.StdEncoding.DecodeString(txeB64)
	sigXdr, errorMsg, status, _ := op.CryptoServiceClient.AddIBMSign(xdrB)

	if errorMsg != nil {
		LOGGER.Errorf("createIssuingAccountInStellar: Error creating new account Error: %v", errorMsg.Error())
		return errorMsg
	}
	if status != http.StatusOK {
		LOGGER.Errorf("createIssuingAccountInStellar: Error creating new account")
		return errors.New("Error creating new account in stellar")
	}
	LOGGER.Debugf("signed transaction: %v", sigXdr)

	b64Xdr := base64.StdEncoding.EncodeToString(sigXdr)
	LOGGER.Debugf("signed transaction: %v", b64Xdr)

	//Submit on Gas service
	hash, ledger, err := op.GasServiceClient.SubmitTxe(b64Xdr)
	if err != nil {
		err = ast.DecodeStellarError(err)
		LOGGER.Error("Create issuing account failed, Error Communicating with Stellar.", err.Error())
		return err
	}
	LOGGER.Debugf("submitTransaction  %v, %v", hash, ledger)
	LOGGER.Infof("CreateAccount Transaction submitted in stellar")
	return nil
}

// createNewOperatingAccountInStellar - creates a new Operating account for the given addr
func (op Operations) createNewOperatingAccountInStellar(issuingAccount string, addr string) error {

	LOGGER.Infof("Creating operating account in stellar: %s", addr)
	stellarNetwork := comn.GetStellarNetwork(os.Getenv(global_environment.ENV_KEY_STELLAR_NETWORK))
	ibmAccount, sequenceNum, err := op.GasServiceClient.GetAccountAndSequence()
	if ibmAccount == "" || sequenceNum == 0 || err != nil {
		LOGGER.Debugf("IBM Gas account failed to load")
		return errors.New("Error getting IBM Gas account")
	}

	wwAdminAccount, errorMsg, status, _ := op.CryptoServiceClient.GetIBMAccount()
	if errorMsg != nil {
		LOGGER.Errorf("Error getting IBM account: %v", errorMsg.Error(), wwAdminAccount)
		return errorMsg
	}
	if status != http.StatusOK {
		LOGGER.Errorf("createNewOperatingAccountInStellar: Error getting IBM account: status %v", status)
		return errors.New("Error creating new account in stellar")
	}

	initialFund := os.Getenv(environment.ENV_KEY_ACCOUNT_INITIAL_FUND)

	//Initial fund is sourced from IBM admin account
	tx, err := b.Transaction(
		b.SourceAccount{AddressOrSeed: ibmAccount},
		b.Sequence{Sequence: sequenceNum},
		stellarNetwork,
		b.CreateAccount(
			//initial fund is sourced from ibmadmin account
			b.SourceAccount{AddressOrSeed: *wwAdminAccount.Address},
			b.Destination{AddressOrSeed: addr},
			b.NativeAmount{Amount: initialFund},
		),
	)

	if err != nil {
		LOGGER.Errorf("createNewOperatingAccountInStellar: Error creating CreateAccount stellar transaction")
		return err
	}
	var txeb b.TransactionEnvelopeBuilder
	err = txeb.Mutate(tx)
	if err != nil {
		LOGGER.Errorf("createNewOperatingAccountInStellar: Error creating CreateAccount stellar transaction")
		return err
	}
	txeB64, err := txeb.Base64()
	LOGGER.Debug("txeB64: ", txeB64)
	//Get IBM signature from crypto service
	xdrB, _ := base64.StdEncoding.DecodeString(txeB64)
	sigXdr, errorMsg, status, _ := op.CryptoServiceClient.AddIBMSign(xdrB)

	if errorMsg != nil {
		LOGGER.Errorf("createNewOperatingAccountInStellar: Error creating new account Error: %v", errorMsg.Error())
		return errorMsg
	}
	if status != http.StatusOK {
		LOGGER.Errorf("createNewOperatingAccountInStellar: Error creating new account")
		return errors.New("Error creating new account in stellar")
	}

	b64Xdr := base64.StdEncoding.EncodeToString(sigXdr)
	LOGGER.Debugf("signed transaction: %v", b64Xdr)

	//Submit on Gas service
	hash, ledger, err := op.GasServiceClient.SubmitTxe(b64Xdr)
	if err != nil {
		err = ast.DecodeStellarError(err)
		LOGGER.Error("CreateAccount Failed. Error Communicating with Stellar.", err.Error())
		return err
	}
	LOGGER.Debugf("submitTransaction  %v, %v", hash, ledger)
	//Success
	LOGGER.Infof("CreateAccount Transaction submitted in stellar")
	return nil
}

// set thresholds, weights, flags and signers for the Operating account
func (op Operations) setOptionsForOperatingAccount(da model.Account, randString string) error {

	stellarNetwork := comn.GetStellarNetwork(os.Getenv(global_environment.ENV_KEY_STELLAR_NETWORK))
	//Get admin account from crypto service
	wwAdminAccount, errorMsg, status, _ := op.CryptoServiceClient.GetIBMAccount()
	if errorMsg != nil {
		LOGGER.Errorf("setOptionsForOperatingAccount: Error getting IBM account: %v", errorMsg.Error())
		return errorMsg
	}
	if status != http.StatusOK {
		LOGGER.Errorf("setOptionsForOperatingAccount: Error getting IBM account: Status %v", status)
		return errors.New("Error getting IBM account")
	}
	ibmAccount, sequenceNum, err := op.GasServiceClient.GetAccountAndSequence()
	if ibmAccount == "" || sequenceNum == 0 {
		LOGGER.Debugf("IBM Gas account failed to load")
		return errors.New("Error getting IBM Gas account")
	}

	//step1: add signer for kill switch
	homeDomain := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	tx, err := b.Transaction(
		b.SourceAccount{AddressOrSeed: ibmAccount},
		b.Sequence{Sequence: sequenceNum},
		stellarNetwork,
		b.SetOptions(
			b.SourceAccount{AddressOrSeed: *da.Address},
			b.AddSigner(participant.GenerateSHA256Hash(randString), fitoficct.SHA_WEIGHT),
		),
	)

	if err != nil {
		LOGGER.Errorf("setOptionsForOperatingAccount: Step1 Error creating SetOptions stellar transaction: %v", err.Error())
		return err
	}
	var txeb b.TransactionEnvelopeBuilder
	err = txeb.Mutate(tx)
	if err != nil {
		LOGGER.Errorf("setOptionsForOperatingAccount: Step1 Error creating Mutating SetOptions stellar transaction: %v", err.Error())
		return err
	}
	txeB64, err := txeb.Base64()
	LOGGER.Debug("txeB64: ", txeB64)

	code, err := op.submitTransaction(da.Name, txeB64)

	if err != nil {
		LOGGER.Errorf("setOptionsForOperatingAccount: Step1 Error Creating set options for account on stellar, Account address: %v, %v, %v", da.Address, err, code)
		return err
	}
	//step2: add thresholds and signer
	//Get Gas service account
	ibmAccount, sequenceNum, err = op.GasServiceClient.GetAccountAndSequence()
	if ibmAccount == "" || sequenceNum == 0 {
		LOGGER.Debugf("setOptionsForOperatingAccount: IBM Gas account failed to load")
		return errors.New("Error getting IBM Gas account")
	}
	tx2, err := b.Transaction(
		b.SourceAccount{AddressOrSeed: ibmAccount},
		b.Sequence{Sequence: sequenceNum},
		stellarNetwork,
		b.SetOptions(
			b.SourceAccount{AddressOrSeed: *da.Address},
			b.SetAuthRevocable(),
			b.SetAuthRequired(),
			b.SetAuthImmutable(),
			b.SetLowThreshold(fitoficct.LOW_THRESHOLD),
			b.SetMediumThreshold(fitoficct.MEDIUM_THRESHOLD),
			b.SetHighThreshold(fitoficct.HIGH_THRESHOLD),
			b.AddSigner(*wwAdminAccount.Address, fitoficct.WW_ADMIN_WEIGHT),
			b.MasterWeight(fitoficct.MASTER_WEIGHT),
			b.HomeDomain(homeDomain),
		),
	)

	if err != nil {
		LOGGER.Errorf("setOptionsForOperatingAccount: Step2 Error creating SetOptions stellar transaction: %v", err.Error())
		return err
	}
	var txeb2 b.TransactionEnvelopeBuilder
	err = txeb2.Mutate(tx2)
	if err != nil {
		LOGGER.Errorf("setOptionsForOperatingAccount: Step2 Error creating Mutating SetOptions stellar transaction: %v", err.Error())
		return err
	}
	txeB64, err = txeb2.Base64()
	LOGGER.Debug("txeB64: ", txeB64)

	code, err = op.submitTransaction(da.Name, txeB64)

	if err != nil {
		LOGGER.Errorf("setOptionsForOperatingAccount: Step2 Error Creating set options for account on stellar, Account address: %v, %v, %v", da.Address, err, code)
		return err
	}

	return nil
}

// set thresholds, weights, flags and signers for the Operating account
func (op Operations) setOptionsForIssuingAccount(ia *model.Account, randString string) error {

	stellarNetwork := comn.GetStellarNetwork(os.Getenv(global_environment.ENV_KEY_STELLAR_NETWORK))

	//Get admin account from crypto service
	wwAdminAccount, errorMsg, status, _ := op.CryptoServiceClient.GetIBMAccount()
	if errorMsg != nil {
		LOGGER.Errorf("setOptionsForIssuingAccount: Error getting IBM account: %v", errorMsg)
		return errorMsg
	}
	if status != http.StatusOK {
		LOGGER.Errorf("setOptionsForIssuingAccount: Error getting IBM account: Status: %v, account:%v", status, ia.Name)
		return errors.New("Error getting IBM account")
	}
	homeDomain := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	//LOGGER.Debugf("GenerateSHA256Hash %v", participant.GenerateSHA256Hash())

	//Get Gas service account
	ibmAccount, sequenceNum, err := op.GasServiceClient.GetAccountAndSequence()
	if ibmAccount == "" || sequenceNum == 0 {
		LOGGER.Debugf("setOptionsForOperatingAccount: IBM Gas account failed to load")
		return errors.New("Error getting IBM Gas account")
	}

	//Breaking setoptions into two steps

	//step1: add signer for kill switch
	tx, err := b.Transaction(
		b.SourceAccount{AddressOrSeed: ibmAccount},
		b.Sequence{Sequence: sequenceNum},
		stellarNetwork,
		b.SetOptions(
			b.SourceAccount{AddressOrSeed: *ia.Address},
			b.AddSigner(participant.GenerateSHA256Hash(randString), fitoficct.SHA_WEIGHT),
		),
	)

	if err != nil {
		LOGGER.Errorf("setOptionsForIssuingAccount: Error creating SetOptions step1  stellar transaction")
		return err
	}

	var txeb b.TransactionEnvelopeBuilder
	err = txeb.Mutate(tx)
	if err != nil {
		LOGGER.Errorf("setOptionsForIssuingAccount: Error creating SetOptions step1 stellar transaction ")
		return err
	}
	txeB64, err := txeb.Base64()
	LOGGER.Debug("txeB64: ", txeB64)

	code, err := op.submitTransaction("issuing", txeB64)

	if err != nil {
		LOGGER.Errorf("setOptionsForIssuingAccount: Error creating SetOptions step1 stellar transaction, code:%v, error:v", code, err.Error())
		return err
	}

	//Get Gas service account
	ibmAccount, sequenceNum, err = op.GasServiceClient.GetAccountAndSequence()
	if ibmAccount == "" || sequenceNum == 0 {
		LOGGER.Debugf("setOptionsForOperatingAccount: IBM Gas account failed to load")
		return errors.New("Error getting IBM Gas account")
	}
	//step2: add thresholds and signer
	tx2, err := b.Transaction(
		b.SourceAccount{AddressOrSeed: ibmAccount},
		b.Sequence{Sequence: sequenceNum},
		stellarNetwork,
		b.SetOptions(
			b.SourceAccount{AddressOrSeed: *ia.Address},
			b.SetAuthRevocable(),
			b.SetAuthRequired(),
			b.SetAuthImmutable(),
			b.SetLowThreshold(fitoficct.LOW_THRESHOLD),
			b.SetMediumThreshold(fitoficct.MEDIUM_THRESHOLD),
			b.SetHighThreshold(fitoficct.HIGH_THRESHOLD),
			b.AddSigner(*wwAdminAccount.Address, fitoficct.WW_ADMIN_WEIGHT),
			b.MasterWeight(fitoficct.MASTER_WEIGHT),
			b.HomeDomain(homeDomain),
		),
	)

	if err != nil {
		LOGGER.Errorf("setOptionsForIssuingAccount: Error creating SetOptions step2  stellar transaction")
		return err
	}

	var txeb2 b.TransactionEnvelopeBuilder
	err = txeb2.Mutate(tx2)
	if err != nil {
		LOGGER.Errorf("setOptionsForIssuingAccount: Error creating SetOptions step2 stellar transaction ")
		return err
	}
	txeB64, err = txeb2.Base64()
	LOGGER.Debug("txeB64: ", txeB64)

	code, err = op.submitTransaction("issuing", txeB64)

	if err != nil {
		LOGGER.Errorf("setOptionsForIssuingAccount: Error creating SetOptions step2 stellar transaction, code:%v, error:v", code, err.Error())
		return err
	}
	return nil
}

func (op Operations) submitTransaction(accountName string, xdrStr string) (code string, err error) {
	xdrB, _ := base64.StdEncoding.DecodeString(xdrStr)
	sig, errorMsg, status, code := op.CryptoServiceClient.SignPayload(accountName, xdrB)

	LOGGER.Debugf("signed payload: %v", base64.StdEncoding.EncodeToString(sig))

	if status != http.StatusOK {
		LOGGER.Errorf("Error creating new account: %v", errorMsg.Error())
		return code, errors.New("Error creating new account")
	}
	if errorMsg != nil {
		return code, errorMsg
	}

	//TBD: will have to integrate with gas service

	sigXdr, errorMsg, status, code := op.CryptoServiceClient.SignXdr(accountName, xdrB, sig, xdrB)

	if errorMsg != nil {
		LOGGER.Errorf("Error creating new account: %v", errorMsg.Error())
		return code, errorMsg
	}
	if status != http.StatusCreated {
		LOGGER.Errorf("Error creating new account")
		return code, errors.New("Error creating new account")
	}

	b64Xdr := base64.StdEncoding.EncodeToString(sigXdr)
	LOGGER.Debugf("signed transaction: %v", b64Xdr)

	//Submit on Gas service
	hash, ledger, err := op.GasServiceClient.SubmitTxe(b64Xdr)
	if err != nil {
		ast.DecodeStellarError(err)
		LOGGER.Warningf("submitTransaction  error: %v", err.Error())
		return "", err
	}
	LOGGER.Debugf("submitTransaction  %v, %v", hash, ledger)

	return "", nil
}
