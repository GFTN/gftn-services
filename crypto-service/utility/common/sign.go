// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package common

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/GFTN/gftn-services/crypto-service/utility/constant"

	"github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"

	"github.com/miekg/pkcs11"
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
	"github.com/GFTN/gftn-services/crypto-service/environment"
	"github.com/GFTN/gftn-services/crypto-service/utility/hsmclient"
	"github.com/GFTN/gftn-services/utility/nodeconfig"
)

//GenericSign for choosing between HSM/Native Stellar signing
func (obj *HsmObject) GenericSign(txeBuilder *build.TransactionEnvelopeBuilder, accounts ...nodeconfig.Account) (build.TransactionEnvelopeBuilder, error) {
	source := strings.ToUpper(os.Getenv(environment.ENV_KEY_ACCOUNT_SOURCE))

	var txBuilder build.TransactionBuilder
	txBuilder.TX = &txeBuilder.E.Tx

	txBuilder.NetworkPassphrase = os.Getenv(global_environment.ENV_KEY_STELLAR_NETWORK)
	err := errors.New("")
	if source == constant.ACCOUNT_FROM_STELLAR {
		//in case of nodeconfig based account mutate transaction env to retain the original signatures
		LOGGER.Infof("Signing transaction with Stellar SDK")
		seeds := make([]string, len(accounts))
		for key, val := range accounts {
			seeds[key] = val.NodeSeed
			sig := build.Sign{Seed: seeds[key]}
			err = sig.MutateTransactionEnvelope(txeBuilder)
		}
		return *txeBuilder, err

	} else if source == constant.ACCOUNT_FROM_HSM {

		txHash32, err := txBuilder.Hash()
		if err != nil {
			LOGGER.Errorf("Error hashing transaction: %v\n", err)
			return *txeBuilder, err
		}
		var txHash []byte = txHash32[:]

		if err != nil {
			LOGGER.Errorf("Error while converting TransactionBuilder to TransactionEnvelopeBuilder: %v\n", err)
			return build.TransactionEnvelopeBuilder{}, err
		}
		LOGGER.Infof("Signing transaction with HSM")
		privateObjectHandles := make([]pkcs11.ObjectHandle, len(accounts))
		publicKeys := make([]string, len(accounts))
		LOGGER.Debugf("Accounts list: %+v", accounts)

		for key, account := range accounts {
			//get private key id
			LOGGER.Infof("Retrieving private key handle ID")
			if value, ok := os.LookupEnv(common.HandleIdName + account.PrivateKeyLabel); ok && value != "0" {
				LOGGER.Infof("Private key handle ID already defined in environment variables")
				LOGGER.Debugf("Private key handle ID: %+v", value)
				objectHandle, _ := strconv.ParseUint(value, 10, 32)
				privateObjectHandles[key] = pkcs11.ObjectHandle(objectHandle)
			} else {
				LOGGER.Infof("Unable to find private key handle ID in environment variables, retrieving from HSM now")
				privateHandleId, _ := obj.retrievePrivateHandleIdFromHsm(account.PrivateKeyLabel)
				os.Setenv(common.HandleIdName+account.PrivateKeyLabel, fmt.Sprint(privateHandleId))
				LOGGER.Debugf("Private key handle ID: %+v", fmt.Sprint(privateHandleId))
				privateObjectHandles[key] = privateHandleId
			}
			//get public key
			publicKeys[key] = account.NodeAddress
		}
		LOGGER.Debugf("privateObjectHandles[0]: %s, obj.Slot: %s, txHash: %s, publicKeys[0]: %s", privateObjectHandles[0], obj.Slot, txHash, publicKeys[0])
		txe, err := hsmclient.GetSignatureAndAddToTransaction(*obj.C, privateObjectHandles[0], obj.Slot, txHash, *txeBuilder, publicKeys[0])
		if err != nil {
			LOGGER.Debugf("Error getting signatures from HSM: %s", err)
			errCode := hsmclient.ParseErrorMsg(err)
			if errCode == hsmclient.SIGNING_ERROR && fmt.Sprint(privateObjectHandles[0]) != "0" {
				obj.C, obj.Session, err = ReinitializeHSM(obj.C, obj.Session)
				if err != nil {
					LOGGER.Errorf("Encounter error while re-intializing HSM login session")
					return build.TransactionEnvelopeBuilder{}, err
				}
				txe, err = hsmclient.GetSignatureAndAddToTransaction(*obj.C, privateObjectHandles[0], obj.Slot, txHash, *txeBuilder, publicKeys[0])
			} else {
				LOGGER.Errorf("Invalid private object handle Id: %v", privateObjectHandles)
			}

			return build.TransactionEnvelopeBuilder{}, err
		}
		return txe, nil
	}
	return build.TransactionEnvelopeBuilder{}, errors.New("cannot fetch correct env variables for GenericSign function")
}

func (obj *HsmObject) GenericVerifySignatureIdentity(unsignedMsg []byte, signedMsg []byte, account nodeconfig.Account) (bool, error) {

	err := hsmclient.VerifyDataWithPublicNodeAddress(account.NodeAddress, unsignedMsg, signedMsg)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (obj *HsmObject) GenericSignPayload(payload []byte, acc nodeconfig.Account) ([]byte, error) {
	source := strings.ToUpper(os.Getenv(environment.ENV_KEY_ACCOUNT_SOURCE))
	LOGGER.Infof("Signing payload data...")

	if source == constant.ACCOUNT_FROM_STELLAR {

		kp, err := keypair.Parse(acc.NodeSeed)
		if err != nil {
			LOGGER.Errorf("There was an error while getting node seed")
		}
		return kp.Sign(payload)

	} else if source == constant.ACCOUNT_FROM_HSM {

		LOGGER.Infof("Start finding HSM object handle with label %s", acc.PrivateKeyLabel)
		privateObjectHandleID, err := obj.retrievePrivateHandleIdFromHsm(acc.PrivateKeyLabel)
		if err != nil {
			return nil, err
		}
		//LOGGER.Debugf("Private Hanldle: %v", privateObjectHandleID)
		signedPayload, err := hsmclient.SignDataWithPrivateKey(*obj.C, privateObjectHandleID, obj.Slot, payload)

		if err != nil {
			errCode := hsmclient.ParseErrorMsg(err)
			if errCode == hsmclient.SIGNING_ERROR && privateObjectHandleID != 0 {
				obj.C, obj.Session, err = ReinitializeHSM(obj.C, obj.Session)
				signedPayload, err := hsmclient.SignDataWithPrivateKey(*obj.C, privateObjectHandleID, obj.Slot, payload)
				return signedPayload, err
			}
		}

		return signedPayload, nil
	}
	return nil, nil
}

func (obj *HsmObject) SignUsingHSMPrivateKeyHandle(payload []byte, privateKeyLabel string) ([]byte, error) {

	var privateObjectHandleID pkcs11.ObjectHandle
	//get private key id
	LOGGER.Infof("Retrieving private key handle ID")
	if value, ok := os.LookupEnv(common.HandleIdName + privateKeyLabel); ok && value != "0" {
		LOGGER.Infof("Private key handle ID already defined in environment variables")
		LOGGER.Debugf("Private key handle ID: %+v", value)
		objectHandle, _ := strconv.ParseUint(value, 10, 32)
		privateObjectHandleID = pkcs11.ObjectHandle(objectHandle)
	} else {
		LOGGER.Infof("Unable to find private key handle ID in environment variables, retrieving from HSM now")
		privateHandleId, err := obj.retrievePrivateHandleIdFromHsm(privateKeyLabel)
		if err != nil {
			LOGGER.Errorf("Error while retrieving the private key handle: %+v", err)
			return nil, err
		}
		os.Setenv(common.HandleIdName+privateKeyLabel, fmt.Sprint(privateHandleId))
		LOGGER.Debugf("Private key handle ID: %+v", fmt.Sprint(privateHandleId))
		privateObjectHandleID = privateHandleId
	}

	return hsmclient.SignDataWithPrivateKey(*obj.C, privateObjectHandleID, obj.Slot, payload)
}

func (obj *HsmObject) retrievePrivateHandleIdFromHsm(privateKeyLabel string) (pkcs11.ObjectHandle, error) {
	_, privateObjectHandleID, err := hsmclient.FindHSMObject(obj.C, obj.Slot, privateKeyLabel)
	if err != nil {
		LOGGER.Errorf("There was an error while getting hsm private object handle id: %s", err.Error())
		if err.Error() == hsmclient.NULL_HANDLE_ID && privateKeyLabel != "" {
			obj.C, obj.Session, err = ReinitializeHSM(obj.C, obj.Session)
			if err != nil {
				LOGGER.Errorf("Encounter error while re-intializing HSM login session: %v", err)
				return 0, err
			}
			_, privateObjectHandleID, err = hsmclient.FindHSMObject(obj.C, obj.Slot, privateKeyLabel)
		}
	}

	if fmt.Sprint(privateObjectHandleID) == "0" {
		LOGGER.Errorf("Unable to retrieve object Handle ID with label %v", privateKeyLabel)
		return 0, err
	}
	return privateObjectHandleID, err
}
