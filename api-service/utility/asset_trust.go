// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package utility

import (
	"os"

	hClient "github.com/stellar/go/clients/horizonclient"
	"github.com/GFTN/gftn-services/gftn-models/model"
	pr_client "github.com/GFTN/gftn-services/participant-registry-client/pr-client"
	"github.com/GFTN/gftn-services/utility/asset"
	"github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
)

//GetAssets : This function combines the list of requested as well as trusted(allowed trust) assets by an account
func GetAssets(accountAddress string, prclient pr_client.PRServiceClient) ([]*model.Asset, error) {
	var assets []*model.Asset
	horizonClient := common.GetHorizonClient(os.Getenv(global_environment.ENV_KEY_HORIZON_CLIENT_URL))
	account, err := horizonClient.LoadAccount(accountAddress)
	if err != nil {
		LOGGER.Errorf("Encounter error while loading account from horizon client: %s", err)
		return nil, err
	}
	for i := range account.Balances {
		horizonAsset := account.Balances[i].Asset
		assetType := asset.GetAssetType(horizonAsset.Code)
		//Get domain name from PR for given issuing address
		if horizonAsset.Type != "native" {
			participant, err := prclient.GetParticipantForIssuingAddress(horizonAsset.Issuer)
			if err == nil && *participant.ID != "" {
				assets = append(assets, &model.Asset{
					AssetCode: &horizonAsset.Code,
					IssuerID:  *participant.ID,
					AssetType: &assetType,
				})

			} else {
				LOGGER.Errorf(err.Error()+"Error retrieving participant for issuer id %v", horizonAsset.Issuer)
			}
		}
	}
	return assets, nil

}

//GetTrustedWWAssets : Gets account information from stellar and parses trusted assets to return
func GetTrustedWWAssets(accountAddress string, prclient pr_client.PRServiceClient) ([]*model.Asset, error) {

	client := common.GetNewHorizonClient(os.Getenv(global_environment.ENV_KEY_HORIZON_CLIENT_URL))

	accountRequest := hClient.AccountRequest{AccountID: accountAddress}

	// Load the account detail from the network
	account, err := client.AccountDetail(accountRequest)
	if err != nil {
		LOGGER.Errorf(err.Error())
		return nil, err
	}

	var assets []*model.Asset
	for _, b := range account.Balances {
		LOGGER.Debug(b.Code, b.Issuer)
		assetCode := b.Code
		issuer := b.Issuer
		flag := false
		if b.IsAuthorized != nil {
			flag = *b.IsAuthorized
		}
		if flag {
			prUrl := os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL)
			prc, _ := pr_client.CreateRestPRServiceClient(prUrl)
			participant, err := prc.GetParticipantForIssuingAddress(issuer)
			if err == nil {
				assetType := asset.GetAssetType(assetCode)
				assets = append(assets, &model.Asset{
					AssetCode: &assetCode,
					IssuerID:  *participant.ID,
					AssetType: &assetType,
				})
			}
		}
	}
	return assets, nil
}
