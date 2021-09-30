// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package secret_manager

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/GFTN/gftn-services/utility/aws/golang/utility"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
)

func Getenv(credential utility.CredentialInfo) (string, error) {
	LOGGER.Infof("Getting environment variable: %s", credential.Variable)
	credential.Variable = strings.ToLower(credential.Variable)
	env, exists := os.LookupEnv(credential.Variable)

	if exists && credential.Variable != "initialize" {
		LOGGER.Infof("Environment variable %s already defined", credential.Variable)
		return env, nil
	} else {
		LOGGER.Infof("Retrieving environment variable %s from secret manager", credential.Variable)

		res, err := GetSecret(credential)
		if err != nil {
			LOGGER.Errorf("Cannot get the specified environment variable: %s", err)
			return "", err
		}

		secretResult := map[string]string{}
		err = json.Unmarshal([]byte(res), &secretResult)
		if err != nil {
			errMsg := errors.New("Error parsing secret object format from AWS")
			LOGGER.Errorf("%s", errMsg)
			return "", errMsg
		}

		var result string
		for key, _ := range secretResult {
			os.Setenv(key, secretResult[key])
			result = secretResult[key]
		}
		os.Setenv(credential.Variable, "true")
		return result, nil
	}
}

func Setenv(credential utility.CredentialInfo, secret utility.SecretContent) error {
	LOGGER.Infof("Setting environment variable: %s", credential.Variable)
	credential.Variable = strings.ToLower(credential.Variable)
	err := CreateSecret(credential, secret)
	if err != nil {
		return err
	}
	if len(secret.Entry) > 1 {
		for _, entry := range secret.Entry {
			os.Setenv(entry.Key, entry.Value)
		}
		os.Setenv(credential.Variable, "true")
	} else {
		os.Setenv(credential.Variable, secret.Entry[0].Value)
	}

	return nil

}

func Updateenv(credential utility.CredentialInfo, secret utility.SecretContent) error {
	LOGGER.Infof("Updating environment variable: %s", credential.Variable)
	credential.Variable = strings.ToLower(credential.Variable)
	err := UpdateSecret(credential, secret)
	if err != nil {
		return err
	}
	if len(secret.Entry) > 1 {
		for _, entry := range secret.Entry {
			os.Setenv(entry.Key, entry.Value)
		}
		os.Setenv(credential.Variable, "true")
	} else {
		os.Setenv(credential.Variable, secret.Entry[0].Value)
	}
	return nil

}

func UpdateAccount(credential utility.CredentialInfo, secret utility.SecretContent) error {
	credential.Variable = strings.ToLower(credential.Variable)
	LOGGER.Infof("Updating account: %s in secret manager & environment variable...", credential.Variable)
	err := UpdateSecret(credential, secret)
	if err != nil {
		return err
	}
	for _, entry := range secret.Entry {
		LOGGER.Debugf("account: %+v \n", entry)
		os.Setenv(credential.Variable+"_"+entry.Key, entry.Value)
	}
	os.Setenv(credential.Variable, "true")
	return nil

}

func CreateAccount(credential utility.CredentialInfo, secret utility.SecretContent) error {

	credential.Variable = strings.ToLower(credential.Variable)
	LOGGER.Infof("Creating account: %s in secret manager & environment variable...", credential.Variable)
	err := CreateSecret(credential, secret)
	if err != nil {
		return err
	}
	for _, entry := range secret.Entry {
		LOGGER.Debugf("account: %+v \n", entry)
		os.Setenv(credential.Variable+"_"+entry.Key, entry.Value)
	}
	os.Setenv(credential.Variable, "true")
	return nil

}

func Deleteenv(credential utility.CredentialInfo) error {
	LOGGER.Infof("Deleting environment variable: %s", credential.Variable)
	credential.Variable = strings.ToLower(credential.Variable)
	err := DeleteSecret(credential, 7)
	if err != nil {
		return errors.New("Encounter error when removing secret from AWS")
	}
	os.Unsetenv(credential.Variable)
	return nil
}

func InitEnv() {
	//getGlobalEnv("initialize")
	getParticipantEnv("initialize")
	getServiceEnv("initialize")
}

func getServiceEnv(var_name string) {
	domainId := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	svcName := os.Getenv(global_environment.ENV_KEY_SERVICE_NAME)
	envVersion := os.Getenv(global_environment.ENV_KEY_ENVIRONMENT_VERSION)
	LOGGER.Infof("Initializing service env variables with Domain name: %s, Service name: %s, Environment version: %s", domainId, svcName, envVersion)
	_, err := Getenv(utility.CredentialInfo{
		Environment: envVersion,
		Domain:      domainId,
		Service:     svcName,
		Variable:    var_name,
	})
	if err != nil {
		panic("Error initializing service with AWS secret manager")
	}
}

func getParticipantEnv(var_name string) {
	domainId := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	envVersion := os.Getenv(global_environment.ENV_KEY_ENVIRONMENT_VERSION)
	LOGGER.Infof("Initializing participant: %s specific env variables", domainId)
	_, err := Getenv(utility.CredentialInfo{
		Environment: envVersion,
		Domain:      domainId,
		Service:     "participant",
		Variable:    var_name,
	})
	if err != nil {
		panic("Error initializing service with AWS secret manager")
	}
}

func getGlobalEnv(var_name string) {
	envVersion := os.Getenv(global_environment.ENV_KEY_ENVIRONMENT_VERSION)
	LOGGER.Infof("Initializing global env variables")
	_, err := Getenv(utility.CredentialInfo{
		Environment: envVersion,
		Domain:      "ww",
		Service:     "global",
		Variable:    var_name,
	})
	if err != nil {
		panic("Error initializing service with AWS secret manager")
	}
}
