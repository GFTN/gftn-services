// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package nodeconfig

import (
	bs "github.com/BurntSushi/toml"
	pel "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/utility/common"
	comn "github.com/GFTN/gftn-services/utility/common"
	"github.com/GFTN/gftn-services/utility/global-environment"
	"io/ioutil"
	"os"
)


type OnboardingOperations struct {
}

type NodeConfig struct {
	OperatingAccounts map[string]Account `toml:"DISTRIBUTION_ACCOUNTS"`
	IssuingAccount    Account            `toml:"ISSUING_ACCOUNT"`
	IBMTokenAccount   Account            `toml:"IBM_TOKEN_ACCOUNT"`
}

type Account struct {
	NodeAddress     string `toml:"NODE_ADDRESS"`
	NodeSeed        string `toml:"NODE_SEED"`
	PrivateKeyLabel string
	PublicKeyLabel  string
}
type WhiteList struct {
	ParticipantDomains string `toml:"PARTICIPANT_DOMAINS"`
}

func GetIBMTokenAccount() (Account, error) {

	_, account, err := IBMTokenAccountExists()
	if err != nil {
		LOGGER.Errorf("Error while checking if IBM Token account exists: %v", err)
		return Account{}, err
	}
	return account, nil
}

// if account name = comn.ISSUING, return issuing account, else return operating account
func GetAccount(accountName string) (Account, error) {
	if accountName == comn.ISSUING {
		return GetIssuingAccount()
	} else {
		return GetOperatingAccount(accountName)
	}
}

func GetOperatingAccount(accountName string) (Account, error) {

	_, account, err := OperatingAccountExists(accountName)
	if err != nil {
		LOGGER.Errorf("Error while checking if operating account exists: %v", err)
		return Account{}, err
	}
	return account, nil

}
func GetIssuingAccount() (Account, error) {

	/*
		1. check if there is a operating account with the given name, if yes, return operating account address
		2. create new operating account (in both stellar and NodeConfig, return the new operating account address
	*/

	LOGGER.Infof("utility:node_config.go:GetIssuingAccount....")

	accountExists, account, err := IssuingAccountExists()
	if err != nil {
		LOGGER.Errorf("Error while checking if issuing account exists: %v", err)
		return Account{}, err
	}
	if accountExists {
		return account, nil
	}
	return Account{}, errors.New("No Issuing Account")
}

// verify if IBM Token account already exists
func IBMTokenAccountExists() (bool, Account, error) {
	nodeconfig, err := loadNodeConfig()
	if err != nil {
		LOGGER.Errorf("Error loading nodeconfig: %v", err)
		return false, Account{}, err
	}

	if nodeconfig.IBMTokenAccount.NodeAddress == "" {
		LOGGER.Infof("IBM Token Account does not exist")
		return false, Account{}, nil
	}
	LOGGER.Infof("IBM Token Account already exists")
	return true, nodeconfig.IBMTokenAccount, nil
}

// verify if issuing account already exists
func IssuingAccountExists() (bool, Account, error) {
	nodeconfig, err := loadNodeConfig()
	if err != nil {
		LOGGER.Errorf("Error loading nodeconfig: %v", err)
		return false, Account{}, err
	}

	if nodeconfig.IssuingAccount.NodeAddress == "" {
		LOGGER.Infof("Issuing Account does not exist")
		return false, Account{}, nil
	}
	LOGGER.Infof("Issuing Account already exists")
	return true, nodeconfig.IssuingAccount, nil
}

// verify if operating account with given name already exists
func OperatingAccountExists(accountName string) (bool, Account, error) {

	config, err := loadNodeConfig()
	if err != nil {
		LOGGER.Errorf("Error loading nodeconfig: %v", err)
		return false, Account{}, err
	}

	for name, account := range config.OperatingAccounts {
		if name == accountName {
			LOGGER.Infof("Operating Account: %v already exists", name)
			return true, account, nil
		}
	}
	return false, Account{}, nil
}

// verify if operating account with given stellar already exists
func OperatingAddressExists(accountAddress string) (bool, Account, error) {
	LOGGER.Infof("accountAddress: ", accountAddress)
	config, err := loadNodeConfig()
	if err != nil {
		LOGGER.Errorf("Error loading nodeconfig: %v", err)
		return false, Account{}, err
	}

	for name, account := range config.OperatingAccounts {
		LOGGER.Infof("account check:%v", name)
		if account.NodeAddress == accountAddress {
			LOGGER.Infof("Operating Address: %v already exists", name)
			return true, account, nil
		}
	}
	LOGGER.Infof("accountAddress did not find operating address", accountAddress)
	return false, Account{}, nil
}

func loadNodeConfig() (NodeConfig, error) {

	configFile, _ := SafeGetConfigFile()
	var config NodeConfig
	if _, err := bs.DecodeFile(configFile, &config); err != nil {
		LOGGER.Errorf("errordecoding nodeconfig: %v", err)
		return NodeConfig{}, err
	}
	return config, nil
}

//Load NodeConfig: expose
func LoadNodeConfig() (NodeConfig, error) {
	return loadNodeConfig()
}

func AddOperatingAccountToNodeConfig(accountID string, account Account) error {
	configFile, _ := SafeGetConfigFile()
	config, err := loadNodeConfig()
	if err != nil {
		LOGGER.Errorf("Error loading nodeconfig: %v", err)
		return err
	}
	distAccounts := config.OperatingAccounts
	distAccounts[accountID] = account
	configBytes, err := pel.Marshal(config)
	if err != nil {
		LOGGER.Errorf("Error marshalling nodeconfig: %v", err)
		return err
	}
	ioutil.WriteFile(configFile, configBytes, 333)
	return nil
}

func AddIssuingAccountToNodeConfig(account Account) error {
	configFile, _ := SafeGetConfigFile()
	config, err := loadNodeConfig()
	if err != nil {
		LOGGER.Errorf("Error loading nodeconfig: %v", err)
		return err
	}
	config.IssuingAccount = account
	configBytes, err := pel.Marshal(config)
	if err != nil {
		LOGGER.Errorf("Error marshalling nodeconfig: %v", err)
		return err
	}
	ioutil.WriteFile(configFile, configBytes, 333)

	return nil
}

func GetAllOperatingAccounts() (accounts []model.Account, err error) {
	accounts = []model.Account{}
	config, err := loadNodeConfig()
	if err != nil {
		LOGGER.Errorf("Error loading nodeconfig: %v", err)
		return accounts, err
	}
	i := 1
	for name, account := range config.OperatingAccounts {
		da := model.Account{}
		da.Name = name
		nodeAddress := account.NodeAddress
		da.Address = &nodeAddress
		accounts = append(accounts, da)
		LOGGER.Infof("DA%v: %v", i, account.NodeAddress)
		i = i + 1
	}
	return accounts, nil
}

func SafeGetConfigFile() (string, error) {
	configFile := os.Getenv(global_environment.ENV_KEY_NODE_CONFIG)
	/*Check Marx fix : Absolute paths without double dots are recommended */
	if !common.IsSafePath(configFile) {
		err := errors.New("file path may be vulnerable")
		LOGGER.Errorf("Error loading node config: %v", err)
		return "", err
	}
	return configFile, nil
}
