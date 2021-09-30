// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package asset

import (
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"strings"

	b "github.com/stellar/go/build"
	hProtocol "github.com/stellar/go/protocols/horizon"
	"github.com/GFTN/gftn-services/crypto-service-client/crypto-client"
	gasservice "github.com/GFTN/gftn-services/gas-service-client"
	"github.com/GFTN/gftn-services/gftn-models/model"
	comn "github.com/GFTN/gftn-services/utility/common"
	"github.com/GFTN/gftn-services/utility/global-environment"
)

/*func GetFriendBotBaseUrl(passphrase string) string {
	if passphrase == n.PublicNetworkPassphrase {
		return horizon.DefaultPublicNetClient.URL
	}
	return horizon.DefaultTestNetClient.URL
}

func GetHorizonClientByPassphrase(passphrase string) horizon.Client {
	if passphrase == n.PublicNetworkPassphrase {
		return *horizon.DefaultPublicNetClient
	}
	return *horizon.DefaultTestNetClient
}*/

//ChangeTrust: Creates change trust transaction and returns xdr
func ChangeTrust(gClient gasservice.GasServiceClient, trustorAddr, trustee, assetCode, limit string) (string, error) {

	stellarNetwork := comn.GetStellarNetwork(os.Getenv(global_environment.ENV_KEY_STELLAR_NETWORK))

	ibmAccount, sequenceNum, err := gClient.GetAccountAndSequence()
	if ibmAccount == "" || sequenceNum == 0 || err != nil {
		LOGGER.Debugf("IBM Gas account failed to load")
		return "", errors.New("Error getting IBM Gas account")
	}

	tx, err := b.Transaction(
		b.SourceAccount{AddressOrSeed: ibmAccount},
		b.Sequence{Sequence: sequenceNum},
		stellarNetwork,
		b.Trust(assetCode, trustee, b.Limit(limit), b.SourceAccount{AddressOrSeed: trustorAddr}),
	)
	if err != nil {
		LOGGER.Errorf(err.Error())
		return "", err
	}
	var txeb b.TransactionEnvelopeBuilder
	err = txeb.Mutate(tx)
	txeB64, err := txeb.Base64()
	if err != nil {
		LOGGER.Errorf(err.Error())
		return "", err
	}
	LOGGER.Debug("txeB64: ", txeB64)
	return txeB64, nil

}

func GetStellarAccount(accountId string) hProtocol.Account {
	LOGGER.Infof("TrustedStellarAsset")
	horizonClient := comn.GetHorizonClient(os.Getenv(global_environment.ENV_KEY_HORIZON_CLIENT_URL))
	account, err := horizonClient.LoadAccount(accountId)
	if err != nil {
		LOGGER.Error("Error while getting the account from Stellar : ", accountId)
	}
	LOGGER.Infof("Got the account")
	return account
}

func IsNative(assetCode string) bool {
	lowercaseCode := strings.ToLower(assetCode)
	return strings.Compare("native", lowercaseCode) == 0 ||
		strings.Compare("xlm", lowercaseCode) == 0 ||
		strings.Compare("lumen", lowercaseCode) == 0
}

func GetBalance(accountId, assetCode, assetIssuer string) string {
	account := GetStellarAccount(accountId)
	if IsNative(assetCode) {
		balStr, err := account.GetNativeBalance()
		if err != nil {
			LOGGER.Errorf(err.Error())
		}
		return balStr
	} else {
		return account.GetCreditBalance(assetCode, assetIssuer)
	}
}

func IsBalanceExist(accountId, assetCode, assetIssuer string) bool {
	account := GetStellarAccount(accountId)
	LOGGER.Debugf("IsBalanceExist accountId %v, assetIssuer %v", accountId, assetIssuer)
	if IsNative(assetCode) {
		_, err := account.GetNativeBalance()
		if err != nil {
			LOGGER.Errorf(err.Error())
		}
		return true
	} else {
		for _, balance := range account.Balances {
			if balance.Asset.Code == assetCode && balance.Asset.Issuer == assetIssuer {
				return true
			}
		}
	}
	return false
}

func GetTrustedAssets(accountAddress string) (bool, []hProtocol.Asset) {
	exeStatus := false
	horizonClient := comn.GetHorizonClient(os.Getenv(global_environment.ENV_KEY_HORIZON_CLIENT_URL))

	account, err := horizonClient.LoadAccount(accountAddress)

	var assets []hProtocol.Asset

	if err != nil {
		return exeStatus, assets
	}

	for _, balance := range account.Balances {
		if balance.Asset.Type != "native" {
			assets = append(assets, hProtocol.Asset{Type: balance.Type, Code: balance.Code, Issuer: balance.Issuer})
		}
	}

	exeStatus = true
	return exeStatus, assets
}

//DoChangeTrust: support DA / DO
func DoChangeTrust(gClient gasservice.GasServiceClient, trustorAccount model.Account, assetCode string, issuerAddress string, limit int64,
	cClient crypto_client.CryptoServiceClient) (error, string) {
	LOGGER.Debugf("DoChangeTrust")
	stellarNetwork := comn.GetStellarNetwork(os.Getenv(global_environment.ENV_KEY_STELLAR_NETWORK))

	LOGGER.Debugf("DoChangeTrust build tx: trustorAccountName %v", trustorAccount.Name)

	var tx *b.TransactionBuilder
	var txe b.TransactionEnvelopeBuilder

	//Get IBM gas account
	ibmAccount, sequenceNum, err := gClient.GetAccountAndSequence()

	tx, err = b.Transaction(
		b.SourceAccount{AddressOrSeed: ibmAccount},
		stellarNetwork,
		b.Sequence{Sequence: sequenceNum},
		b.ChangeTrust(
			b.SourceAccount{*trustorAccount.Address},
			b.Asset{assetCode, issuerAddress, false},
			b.Limit(comn.IntToString(limit)),
		),
	)
	if err != nil {
		LOGGER.Errorf("Error creating change trust transaction: %v", err.Error())
		return errors.New("1206"), err.Error()
	}

	err = txe.Mutate(tx)

	if err != nil {
		msg := "Error during building Mutate"
		LOGGER.Error(msg)
		return errors.New("1208"), msg
	}

	txeB64, err := txe.Base64()
	//TBD: will have to integrate with gas service
	xdrB, _ := base64.StdEncoding.DecodeString(txeB64)

	//Get signed by issuing account on crypto service
	sigXdr, errorMsg, status, _ := cClient.ParticipantSignXdr(trustorAccount.Name, xdrB)

	if status != http.StatusCreated {
		LOGGER.Errorf("Error creating change trust %v", errorMsg.Error())
		return errors.New("1208"), errorMsg.Error()
	}
	LOGGER.Debugf("signed transaction: %v", base64.StdEncoding.EncodeToString(sigXdr))

	if errorMsg != nil {
		msg := "Signing trust went through. Error during encoding"
		LOGGER.Error(msg)
		return errors.New("1209"), msg
	}

	txeB64 = base64.StdEncoding.EncodeToString(sigXdr)

	//Post to gas service
	hash, ledger, err := gClient.SubmitTxe(txeB64)
	if err != nil {
		LOGGER.Warningf("ChangeTrust failed gas service error... %v ", err.Error())
		return err, "ChangeTrust failed:" + err.Error()
	}
	LOGGER.Debugf("Hash:%v, Ledger:%v", hash, ledger)

	msg := "Transaction posted in ledger: " + hash
	return nil, msg
}

//return true if AuthRequired && AuthRevocable
func isAccountFlagsValid(accountId string) bool {
	account := GetStellarAccount(accountId)
	return account.Flags.AuthRequired && account.Flags.AuthRevocable
}

//DoAllowTrust support DO / DA
func DoAllowTrust(gClient gasservice.GasServiceClient, trustorAccountAddress string, assetCode string, issuerAddress string,
	authorize bool, cClient crypto_client.CryptoServiceClient) (error, string) {

	LOGGER.Debugf("DoAllowTrust")
	if GetAssetType(assetCode) != model.AssetAssetTypeDO && GetAssetType(assetCode) != model.AssetAssetTypeDA {
		msg := "Asset code is not DO nor DA: " + assetCode
		LOGGER.Error(msg)
		return errors.New("1201"), msg
	}

	LOGGER.Debugf("DoAllowTrust Transaction issuingAccount =%v, assetCode =%v, trustorAddr=%v, authorize=%v",
		issuerAddress, assetCode, trustorAccountAddress, authorize)

	//check AccountFlags
	if !isAccountFlagsValid(issuerAddress) {
		msg := "Issuing account's flags is not valid: " + issuerAddress
		LOGGER.Errorf(msg)
		return errors.New("1211"), msg
	}

	/*
		Submit a AllowTrust operation with authorize=boolean
	*/

	stellarNetwork := comn.GetStellarNetwork(os.Getenv(global_environment.ENV_KEY_STELLAR_NETWORK))

	//Get IBM gas account
	ibmAccount, sequenceNum, err := gClient.GetAccountAndSequence()

	tx, err := b.Transaction(
		b.SourceAccount{AddressOrSeed: ibmAccount},
		stellarNetwork,
		b.Sequence{Sequence: sequenceNum},
		//	Submit a AllowTrust operation with authorize=boolean
		b.AllowTrust(
			b.SourceAccount{AddressOrSeed: issuerAddress},
			b.Trustor{Address: trustorAccountAddress},
			b.AllowTrustAsset{Code: assetCode},
			b.Authorize{Value: authorize}),
	)
	if err != nil {
		msg := "Error while allowing trust: " + err.Error()
		LOGGER.Error(msg)
		return errors.New("1222"), msg
	}
	var txe b.TransactionEnvelopeBuilder
	err = txe.Mutate(tx)

	if err != nil {
		msg := "Error during building Mutate"
		LOGGER.Error(msg)
		return errors.New("1208"), msg
	}

	txeB64, err := txe.Base64()
	//TBD: will have to integrate with gas service
	xdrB, _ := base64.StdEncoding.DecodeString(txeB64)

	//Get signed by issuing account on crypto service
	sigXdr, errorMsg, status, _ := cClient.ParticipantSignXdr(comn.ISSUING, xdrB)

	if status != http.StatusCreated {
		LOGGER.Errorf("Error creating allow trust %v", errorMsg.Error())
		return errors.New("1208"), errorMsg.Error()
	}
	LOGGER.Debugf("signed transaction: %v", base64.StdEncoding.EncodeToString(sigXdr))

	if errorMsg != nil {
		msg := "Signing trust went through. Error during encoding"
		LOGGER.Error(msg)
		return errors.New("1209"), msg
	}

	txeB64 = base64.StdEncoding.EncodeToString(sigXdr)

	//Post to gas service
	hash, ledger, err := gClient.SubmitTxe(txeB64)
	if err != nil {
		LOGGER.Warningf("AllowTrust failed gas service error... %v ", err.Error())
		return err, "AllowTrust failed:" + err.Error()
	}
	LOGGER.Debugf("Hash:%v, Ledger:%v", hash, ledger)

	msg := "Transaction posted in ledger: " + hash
	return nil, msg
}
