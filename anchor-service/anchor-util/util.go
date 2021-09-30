// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package anchor_util

// this util package is dedicated for anchor service so that there is no dependency on the other
// util package and can be deployed as a standalone service
import (
	"errors"
	"os"
	"strings"

	b "github.com/stellar/go/build"
	hProtocol "github.com/stellar/go/protocols/horizon"
	gasserviceclient "github.com/GFTN/gftn-services/gas-service-client"
	"github.com/GFTN/gftn-services/gftn-models/model"
	pr_client "github.com/GFTN/gftn-services/participant-registry-client/pr-client"
	util "github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/participant"
	vutils "github.com/GFTN/gftn-services/utility/vault/utils"
)

// getParticipantForDomain : Get participant for domain
func GetParticipantForDomain(domain string) (model.Participant, error) {
	LOGGER.Debugf("getParticipantForDomain domain = %v", domain)
	var participant model.Participant

	prClient, err := pr_client.CreateRestPRServiceClient(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL))
	//prClient, err := pr_client.CreateMockPRServiceClient()
	if err != nil {
		LOGGER.Errorf(" Error getParticipantForDomain CreateRestPRServiceClient failed  %v", err)
		return participant, errors.New("ANCHOR-0007")
	}
	participant, err = prClient.GetParticipantForDomain(domain)
	if err != nil {
		LOGGER.Errorf(" Error getParticipantForDomain failed: %v", err)
		return participant, errors.New("ANCHOR-0006")
	}

	return participant, nil
}

func ParseFederationName(name string) (string, string, error) {

	// name must be of the form:  accountIdentifier*participantDomain.  Look for "*"
	if !strings.Contains(name, "*") {
		return "", "", errors.New("ANCHOR-0002")
	}

	nameParts := strings.Split(name, "*")
	if len(nameParts) != 2 {
		LOGGER.Warningf("Unable to parse given federation name:  %v", name)
		return "", "", errors.New("ANCHOR-0003")
	}

	return nameParts[0], nameParts[1], nil
}

func GetAccountAddressForParticipant(participant model.Participant, accountName string) string {

	LOGGER.Infof("Account name is: %v", accountName)
	stellarAddress := ""
	if accountName == "" || accountName == util.ISSUING {
		stellarAddress = participant.IssuingAccount
	} else {
		for _, distAccount := range participant.OperatingAccounts {
			if distAccount.Name == accountName {
				stellarAddress = *distAccount.Address
				break
			}
		}
	}
	return stellarAddress
}

// for digital assets like stable coins, IBM account will be used to sign the transaction

func AllowTrustForDigitalAsset(gClient gasserviceclient.GasServiceClient, trustorAddress string, anchorAddress string,
	assetCode string, authorize bool, vsession vutils.Session) (error, string) {
	LOGGER.Debugf("AllowTrustForDigitalAsset")

	//check AccountFlags
	if !isAccountFlagsValid(anchorAddress) {
		msg := "Anchor's flags is not valid: " + anchorAddress
		LOGGER.Errorf(msg)
		return errors.New("1211"), msg
	}

	/*
		Submit a AllowTrust operation with authorize=boolean
	*/
	stellarNetwork := util.GetStellarNetwork(os.Getenv(global_environment.ENV_KEY_STELLAR_NETWORK))

	//Get IBM gas account
	ibmAccount, sequenceNum, err := gClient.GetAccountAndSequence()

	tx, err := b.Transaction(
		b.SourceAccount{AddressOrSeed: ibmAccount},
		stellarNetwork,
		b.Sequence{Sequence: sequenceNum},
		//	Submit a AllowTrust operation with authorize=boolean
		b.AllowTrust(
			b.SourceAccount{AddressOrSeed: anchorAddress},
			b.Trustor{Address: trustorAddress},
			b.AllowTrustAsset{Code: assetCode},
			b.Authorize{Value: authorize}),
	)
	if err != nil {
		msg := "Error while allowing trust: " + err.Error()
		LOGGER.Error(msg)
		return errors.New("1222"), msg
	}

	//Get IBM token account from nc, vault or AWS secret mngr
	ibmAdminAccount, err := participant.GenericGetIBMTokenAccount(vsession)
	if err != nil {
		msg := "Error getting IBM account"
		LOGGER.Error(msg)
		return errors.New("1209"), msg
	}

	txe, err := tx.Sign(ibmAdminAccount.NodeSeed)
	if err != nil {
		msg := "Error while signing trust"
		LOGGER.Error(msg)
		return errors.New("1207"), msg
	}

	txeB64, err := txe.Base64()
	if err != nil {
		msg := "Error during building Mutate"
		LOGGER.Error(msg)
		return errors.New("1208"), msg
	}

	//Post to gas service
	hash, ledger, err := gClient.SubmitTxe(txeB64)
	if err != nil {
		LOGGER.Warningf("AllowTrustForDigitalAsset failed gas service error... %v ", err.Error())
		return err, "AllowTrustForDigitalAsset failed:" + err.Error()
	}
	LOGGER.Debugf("Hash:%v, Ledger:%v", hash, ledger)

	msg := "AllowTrustForDigitalAsset Transaction posted in ledger: " + hash
	return nil, msg
}

//return true if AuthRequired && AuthRevocable
func isAccountFlagsValid(accountId string) bool {
	account := GetStellarAccount(accountId)
	return account.Flags.AuthRequired && account.Flags.AuthRevocable
}

func GetStellarAccount(accountId string) hProtocol.Account {
	LOGGER.Infof("TrustedStellarAsset")
	horizonClient := util.GetHorizonClient(os.Getenv(global_environment.ENV_KEY_HORIZON_CLIENT_URL))
	account, err := horizonClient.LoadAccount(accountId)
	if err != nil {
		LOGGER.Error("Error while getting the account from Stellar : ", accountId)
	}
	LOGGER.Infof("Got the account")
	return account
}
