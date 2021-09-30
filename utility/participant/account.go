// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package participant

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/GFTN/gftn-services/gftn-models/model"
	pr_client "github.com/GFTN/gftn-services/participant-registry-client/pr-client"

	comn "github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"

	secret_manager "github.com/GFTN/gftn-services/utility/aws/golang/secret-manager"
	"github.com/GFTN/gftn-services/utility/aws/golang/utility"
	"github.com/GFTN/gftn-services/utility/nodeconfig"
	VaultApi "github.com/GFTN/gftn-services/utility/vault/api"
	"github.com/GFTN/gftn-services/utility/vault/utils"
)

//retrieve either issuing account or operating account
func GenericGetAccount(session utils.Session, accountName string) (nodeconfig.Account, error) {
	storage := strings.ToUpper(os.Getenv(global_environment.ENV_KEY_SECRET_STORAGE_LOCATION))
	if storage == comn.LOCAL_SECRET {
		LOGGER.Infof("Getting account from nodeconfig")
		return nodeconfig.GetAccount(accountName)
	} else if storage == comn.VAULT_SECRET {
		LOGGER.Infof("Getting account from vault")
		return GetAccountFromVault(session, "LOCAL", accountName)
	} else if storage == comn.AWS_SECRET {
		LOGGER.Infof("Getting account from AWS")
		domainId := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
		envVersion := os.Getenv(global_environment.ENV_KEY_ENVIRONMENT_VERSION)
		return getAccountFromAws(utility.CredentialInfo{
			Environment: envVersion,
			Domain:      domainId,
			Service:     "account",
			Variable:    accountName,
		})
	}
	return nodeconfig.Account{}, errors.New("Cannot fetch correct env variables for GenericGetAccount function")
}

//retrieve IBM token account from either Vault or local node_config
func GenericGetIBMTokenAccount(session utils.Session) (nodeconfig.Account, error) {
	storage := strings.ToUpper(os.Getenv(global_environment.ENV_KEY_SECRET_STORAGE_LOCATION))
	if storage == comn.LOCAL_SECRET {
		LOGGER.Infof("Getting IBM account from nodeconfig")
		return nodeconfig.GetIBMTokenAccount()
	} else if storage == comn.VAULT_SECRET {
		LOGGER.Infof("Getting IBM account from vault")
		return GetAccountFromVault(session, "IBM", "IBM")
	} else if storage == comn.AWS_SECRET {
		LOGGER.Infof("Getting IBM account from AWS")
		envVersion := os.Getenv(global_environment.ENV_KEY_ENVIRONMENT_VERSION)
		domainId := os.Getenv(global_environment.ENV_KEY_IBM_TOKEN_DOMAIN_ID)
		return getAccountFromAws(utility.CredentialInfo{
			Environment: envVersion,
			Domain:      domainId,
			Service:     "account",
			Variable:    "token",
		})
	}
	return nodeconfig.Account{}, errors.New("Cannot fetch correct env variables for GenericGetIBMTokenAccount function")
}

func GetAccountFromVault(session utils.Session, safeName string, accountName string) (nodeconfig.Account, error) {
	var account = nodeconfig.Account{}
	accountName = strings.ToUpper(accountName)

	if address, ok := os.LookupEnv(accountName + "_NODE_ADDRESS"); ok {
		LOGGER.Infof(accountName + " account already set as env variables")
		account.NodeAddress = address
		account.PublicKeyLabel = os.Getenv(accountName + "_PUBLIC_LABEL")
		account.PrivateKeyLabel = os.Getenv(accountName + "_PRIVATE_LABEL")
		LOGGER.Debugf("Getting account from local env: %v", account.NodeAddress)
		return account, nil

	} else {
		LOGGER.Infof(accountName + " account haven't set as env variables, now retrieving...")
		//initVaultSession()
		var err error
		var wg sync.WaitGroup
		if safeName == "LOCAL" {
			safeName = os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
		}
		wg.Add(4)
		go func() {
			defer wg.Done()
			account.NodeAddress = getAccountAttribute(session, safeName, accountName+"_NODE_ADDRESS")
		}()
		go func() {
			defer wg.Done()
			account.PublicKeyLabel = getAccountAttribute(session, safeName, accountName+"_PUBLIC_LABEL")
		}()
		go func() {
			defer wg.Done()
			account.PrivateKeyLabel = getAccountAttribute(session, safeName, accountName+"_PRIVATE_LABEL")
		}()
		wg.Wait()

		LOGGER.Infof("Setting " + accountName + " account as environment variable")
		setEnvVariablesFromVault(accountName+"_NODE_ADDRESS", account.NodeAddress)
		setEnvVariablesFromVault(accountName+"_PUBLIC_LABEL", account.PublicKeyLabel)
		setEnvVariablesFromVault(accountName+"_PRIVATE_LABEL", account.PrivateKeyLabel)
		return account, err
	}

}

// verify if issuing account already exists
func GenericIssuingAccountExists(session utils.Session) (bool, nodeconfig.Account, error) {
	storage := strings.ToUpper(os.Getenv(global_environment.ENV_KEY_SECRET_STORAGE_LOCATION))
	if storage == comn.LOCAL_SECRET {
		return nodeconfig.IssuingAccountExists()
	} else if storage == comn.VAULT_SECRET {
		account, err := GetAccountFromVault(session, "LOCAL", comn.ISSUING)
		if account.NodeAddress == "" {
			return false, nodeconfig.Account{}, err
		}
		return true, account, nil
	} else if storage == comn.AWS_SECRET {
		account, err := GenericGetAccount(session, comn.ISSUING)
		if err != nil {
			return false, nodeconfig.Account{}, err
		}
		return true, account, nil
	}
	return false, nodeconfig.Account{}, errors.New("Cannot fetch correct env variables for GenericIssuingAccountExists function")
}

// verify if operating account with given name already exists
func GenericOperatingAccountExists(session utils.Session, accountName string) (bool, nodeconfig.Account, error) {
	storage := strings.ToUpper(os.Getenv(global_environment.ENV_KEY_SECRET_STORAGE_LOCATION))
	if storage == comn.LOCAL_SECRET {
		return nodeconfig.OperatingAccountExists(accountName)
	} else if storage == comn.VAULT_SECRET {
		account, err := GetAccountFromVault(session, "LOCAL", accountName)
		if account.NodeAddress == "" {
			return false, nodeconfig.Account{}, err
		}
		return true, account, nil
	} else if storage == comn.AWS_SECRET {
		account, err := GenericGetAccount(session, accountName)
		if err != nil {
			return false, nodeconfig.Account{}, err
		}
		return true, account, nil
	}
	return false, nodeconfig.Account{}, errors.New("Cannot fetch correct env variables for GenericOperatingAccountExists function")
}

func setEnvVariablesFromVault(key string, value string) {
	if value == "" {
		LOGGER.Warningf("Error while setting env variables, value is empty")
	} else {
		os.Setenv(key, value)
	}
}

func getAccountAttribute(session utils.Session, domainName string, keyword string) string {

	//cyberark safe name does not allow '.' but '_'
	safeName := strings.Replace(domainName, ".", "_", -1)
	result, err := VaultApi.GetPasswordWithAim(session, safeName, keyword)
	// no account found on the vault
	if result == "" {
		LOGGER.Warningf("No account:%s found in safe %s", keyword, safeName)
		return ""
	}

	// something goes wrong during fetching data from the vault
	if err != nil {
		return ""
	}
	return result
}

func getAccountFromAws(credential utility.CredentialInfo) (nodeconfig.Account, error) {

	credential.Variable = strings.ToLower(credential.Variable)
	_, exists := os.LookupEnv(credential.Variable)

	if exists {
		LOGGER.Infof("Environment variable: \"%s\" already defined in environment", credential.Variable)
	} else {
		LOGGER.Infof("Retrieving account: \"%s\" from secret manager", credential.Variable)

		res, err := secret_manager.GetSecret(credential)
		if err != nil {
			LOGGER.Errorf("Cannot get the specified environment variable: %s", err)
			return nodeconfig.Account{}, err
		}

		secretResult := map[string]string{}
		err = json.Unmarshal([]byte(res), &secretResult)
		if err != nil {
			errMsg := errors.New("Error parsing secret object format from AWS")
			LOGGER.Errorf("%s", errMsg)
			return nodeconfig.Account{}, errMsg
		}

		for key, val := range secretResult {
			if strings.TrimSpace(val) == "" {
				LOGGER.Infof("Found %s account on AWS secret manager, but the value is empty", credential.Variable)
				return nodeconfig.Account{}, nil
			}
			os.Setenv(credential.Variable+"_"+key, secretResult[key])
		}
		os.Setenv(credential.Variable, "true")
	}

	var account nodeconfig.Account
	account.NodeAddress = os.Getenv(credential.Variable + "_" + global_environment.ENV_KEY_NODE_ADDRESS)
	account.NodeSeed = os.Getenv(credential.Variable + "_" + global_environment.ENV_KEY_NODE_SEED)
	account.PrivateKeyLabel = os.Getenv(credential.Variable + "_" + global_environment.ENV_KEY_PRIVATE_KEY_LABEL)
	account.PublicKeyLabel = os.Getenv(credential.Variable + "_" + global_environment.ENV_KEY_PUBLIC_KEY_LABEL)
	return account, nil
}

func GenericStoreAccount(session utils.Session, accountName string, account nodeconfig.Account, randString string) error {
	location := strings.ToUpper(os.Getenv(global_environment.ENV_KEY_SECRET_STORAGE_LOCATION))

	if location == comn.LOCAL_SECRET {
		LOGGER.Infof("Storing account info into local nodeconfig file...")
		if strings.ToUpper(accountName) == strings.ToUpper(comn.ISSUING) {
			nodeconfig.AddIssuingAccountToNodeConfig(account)
		} else {
			nodeconfig.AddOperatingAccountToNodeConfig(accountName, account)
		}
	} else if location == comn.VAULT_SECRET {
		LOGGER.Infof("Storing account info into the vault...")
		err := addAccountToVault(session, account, accountName)
		if err != nil {
			LOGGER.Debugf("Error  %s", err)
			return err
		}
	} else if location == comn.AWS_SECRET {
		LOGGER.Infof("Storing account info into aws secret manager...")
		domainId := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
		envVersion := os.Getenv(global_environment.ENV_KEY_ENVIRONMENT_VERSION)
		var accountCredential = utility.CredentialInfo{
			Environment: envVersion,
			Domain:      domainId,
			Service:     "account",
			Variable:    accountName,
		}

		var entries = []utility.SecretEntry{
			utility.SecretEntry{
				Key:   global_environment.ENV_KEY_NODE_ADDRESS,
				Value: account.NodeAddress,
			},
			utility.SecretEntry{
				Key:   global_environment.ENV_KEY_NODE_SEED,
				Value: account.NodeSeed,
			},
			utility.SecretEntry{
				Key:   global_environment.ENV_KEY_PUBLIC_KEY_LABEL,
				Value: account.PublicKeyLabel,
			},
			utility.SecretEntry{
				Key:   global_environment.ENV_KEY_PRIVATE_KEY_LABEL,
				Value: account.PrivateKeyLabel,
			},
		}

		var accountSecretContent = utility.SecretContent{
			Entry:       entries,
			Description: accountName + " of " + domainId,
		}

		err := secret_manager.CreateAccount(accountCredential, accountSecretContent)
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == secretsmanager.ErrCodeResourceExistsException {
				LOGGER.Warningf("%s Now switching to update account process.", aerr.Error())
				err = secret_manager.UpdateAccount(accountCredential, accountSecretContent)
			}
		}

		if err != nil {
			LOGGER.Debugf("Encounter error while storing account to AWS: %s", err)
			return err
		}

		// store secret string
		var accountCredential2 = utility.CredentialInfo{
			Environment: envVersion,
			Domain:      domainId,
			Service:     "killswitch",
			Variable:    "accounts",
		}

		var entries2 = []utility.SecretEntry{
			utility.SecretEntry{
				Key:   account.NodeAddress,
				Value: randString,
			},
		}

		var accountSecretContent2 = utility.SecretContent{
			Entry:       entries2,
			Description: accountName + " of " + domainId,
		}
		err = secret_manager.CreateSecret(accountCredential2, accountSecretContent2)
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == secretsmanager.ErrCodeResourceExistsException {
				LOGGER.Warningf("%s Now switching to append killswitch account process", aerr.Error())
				err = secret_manager.AppendSecret(accountCredential2, accountSecretContent2)
			}
		}

		if err != nil {
			LOGGER.Debugf("Encounter error while storing random string of account to AWS: %s", err)
			return err
		}

	} else {
		LOGGER.Errorf("Error reading account storage location environment settings")
	}
	return nil
}

// verify if operating account with given stellar already exists
func GenericOperatingAddressExists(accountAddress string) (bool, nodeconfig.Account, error) {
	storage := strings.ToUpper(os.Getenv(global_environment.ENV_KEY_SECRET_STORAGE_LOCATION))
	//TODO
	if storage == comn.LOCAL_SECRET {
		return nodeconfig.OperatingAddressExists(accountAddress)
	} else if storage == comn.VAULT_SECRET {

	}
	return false, nodeconfig.Account{}, errors.New("Cannot fetch correct env variables for GenericOperatingAddressExists function")
}

func GenericIsIssuer(session utils.Session, asset model.Asset) bool {
	var iaNc = nodeconfig.Account{}

	prUrl, err := GetServiceUrl(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL), os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME))
	if err != nil {
		LOGGER.Errorf("Encounter error when retrieving participant registry url")
		return false
	}
	var PRClient, _ = pr_client.CreateRestPRServiceClient(prUrl)
	storage := strings.ToUpper(os.Getenv(global_environment.ENV_KEY_SECRET_STORAGE_LOCATION))
	if storage == comn.LOCAL_SECRET {
		iaNc, _ = nodeconfig.GetIssuingAccount()
	} else if storage == comn.VAULT_SECRET {
		iaNc, _ = GetAccountFromVault(session, "LOCAL", comn.ISSUING)
	} else if storage == comn.AWS_SECRET {
		iaNc, _ = GenericGetAccount(session, comn.ISSUING)
	}
	iaPr, _ := PRClient.GetParticipantIssuingAccount(asset.IssuerID)
	if iaNc.NodeAddress == iaPr {
		return true
	}
	return false
}

func addAccountToVault(session utils.Session, account nodeconfig.Account, accountID string) error {
	//sessions.InitVaultSession()
	storage := strings.ToUpper(os.Getenv(global_environment.ENV_KEY_SECRET_STORAGE_LOCATION))
	if storage != comn.LOCAL_SECRET {
		// in case of storage location vault, save these account attributes to vault, skip storage for nodeconfig
		var wg sync.WaitGroup
		wg.Add(4)
		domainName := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
		//cyberark safe name does not allow '.' but '_'
		go func() {
			defer wg.Done()
			addAccountAttribute(session, domainName, accountID+"_NODE_ADDRESS", account.NodeAddress)
		}()
		go func() {
			defer wg.Done()
			addAccountAttribute(session, domainName, accountID+"_PUBLIC_LABEL", account.PublicKeyLabel)
		}()
		go func() {
			defer wg.Done()
			addAccountAttribute(session, domainName, accountID+"_PRIVATE_LABEL", account.PrivateKeyLabel)
		}()

		wg.Wait()
	}
	return nil
}

func addAccountAttribute(vaultSession utils.Session, domainName string, accountAttr string, data string) {
	underscoreDomainName := strings.Replace(domainName, ".", "_", -1)
	err := VaultApi.AddAccount(vaultSession, underscoreDomainName, accountAttr, data, "")

	// something goes wrong during fetching data from the vault
	if err != nil {
		LOGGER.Errorf("Error while storing " + accountAttr + " to the vault")
	}
}
