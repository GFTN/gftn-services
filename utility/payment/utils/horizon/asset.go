// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package horizon

import (
	"errors"
	"os"

	"github.com/stellar/go/clients/horizonclient"
	pr_client "github.com/GFTN/gftn-services/participant-registry-client/pr-client"
	"github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
)

func IsIssuer(pId, assetCode string) bool {
	LOGGER.Infof("Checking if %v is the issuer of %v", pId, assetCode)

	// get issuing address using participant ID
	prclient, _ := pr_client.CreateRestPRServiceClient(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL))
	issuingAccount, err := prclient.GetParticipantAccount(pId, common.ISSUING)
	if err != nil {
		LOGGER.Errorf(err.Error())
		LOGGER.Error("Validation Failed!")
		return false
	}

	client := common.GetNewHorizonClient(os.Getenv(global_environment.ENV_KEY_HORIZON_CLIENT_URL))

	assetRequest := horizonclient.AssetRequest{ForAssetCode: assetCode,
		ForAssetIssuer: issuingAccount}

	// Load the asset List from the network
	assetsList, err := client.Assets(assetRequest)
	if err != nil {
		LOGGER.Errorf("Error loading asset list", err)
		return false
	}

	for _, b := range assetsList.Embedded.Records {
		if b.Code == assetCode && b.Issuer == issuingAccount {
			LOGGER.Infof("Issuer validation success")
			return true
		}
	}

	LOGGER.Infof("Issuer validation failed")
	return false

}

func CheckBalance(pId, assetCode, accountName, assetIssuer string) (string, error) {
	LOGGER.Infof("Checking the %v balance of participant %v", assetCode, pId)
	// get issuing address using participant ID
	prclient, _ := pr_client.CreateRestPRServiceClient(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL))
	issuingAccount, err := prclient.GetParticipantAccount(pId, accountName)
	if err != nil {
		LOGGER.Errorf(err.Error())
		LOGGER.Error("Validation Failed!")
		return "", err
	}

	// query all issued asset of the issuing address
	horizonClient := common.GetHorizonClient(os.Getenv(global_environment.ENV_KEY_HORIZON_CLIENT_URL))
	stellarAccount, err := horizonClient.LoadAccount(issuingAccount)
	// check if the asset exists in the list
	for i := range stellarAccount.Balances {
		horizonAsset := stellarAccount.Balances[i]
		//Get domain name from PR for given issuing address
		if horizonAsset.Asset.Code == assetCode && horizonAsset.Issuer == assetIssuer {
			LOGGER.Infof("Validation Success!")
			return horizonAsset.Balance, nil
		}
	}
	LOGGER.Error("Validation Failed!")
	return "", errors.New("Cannot find the asset " + assetCode + " from " + pId + "'s stellar account")
}
