// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package participants

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/go-openapi/validate"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/GFTN/gftn-services/api-service/environment"
	apiutil "github.com/GFTN/gftn-services/api-service/utility"
	crypto_client "github.com/GFTN/gftn-services/crypto-service-client/crypto-client"
	"github.com/GFTN/gftn-services/gftn-models/model"
	pr_client "github.com/GFTN/gftn-services/participant-registry-client/pr-client"
	"github.com/GFTN/gftn-services/utility"
	comn "github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/participant"
	"github.com/GFTN/gftn-services/utility/response"
	vauth "github.com/GFTN/gftn-services/utility/vault/auth"
	"github.com/GFTN/gftn-services/utility/vault/utils"
)

type ParticipantOperations struct {
	RestPRServiceClient pr_client.RestPRServiceClient
	MockPRServiceClient pr_client.MockPRServiceClient
	VaultSession        utils.Session
	CryptoServiceClient crypto_client.CryptoServiceClient
}

var (
	TYPE_BOTH    = "BOTH"
	TYPE_ISSUED  = "ISSUED"
	TYPE_TRUSTED = "TRUSTED"
)

func CreateParticipantOperations() (ParticipantOperations, error) {

	op := ParticipantOperations{}

	restPRServiceClient, err := pr_client.CreateRestPRServiceClient(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL))
	mockPRServiceClient := pr_client.MockPRServiceClient{}
	utility.ExitOnErr(LOGGER, err, "Unable to create Participant Registry Client")
	op.RestPRServiceClient = restPRServiceClient
	op.MockPRServiceClient = mockPRServiceClient
	var cClient crypto_client.CryptoServiceClient
	if os.Getenv(environment.ENV_KEY_CRYPTO_SERVICE_CLIENT) == "mock" {
		cClient, _ = crypto_client.CreateMockCryptoServiceClient()
	} else {
		err := errors.New("")
		cServiceInternalUrl, err := participant.GetServiceUrl(os.Getenv(global_environment.ENV_KEY_CRYPTO_SVC_INTERNAL_URL), os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME))
		if err != nil {
			return op, err
		}
		cClient, err = crypto_client.CreateRestCryptoServiceClient(cServiceInternalUrl)
		if err != nil {
			return op, err
		}
	}
	op.CryptoServiceClient = cClient

	op.VaultSession = utils.Session{}

	if os.Getenv(global_environment.ENV_KEY_SECRET_STORAGE_LOCATION) == comn.VAULT_SECRET {
		//Vault location
		err := errors.New("")
		op.VaultSession, err = vauth.GetSession()
		if err != nil {
			LOGGER.Errorf("Error reading account source environment settings")
			return op, err
		}
	}
	return op, nil
}

/*func (op ParticipantOperations) GetParticipantRegistryByCountry(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	countryCode := vars["country_code"]

	if countryCode == "" {
		response.NotifyWWError(w, request, http.StatusBadRequest, "API-1041",
			errors.New("country_code is required in query"))
		return
	}

	participants, err := op.RestPRServiceClient.GetParticipantsByCountry(countryCode)
	if err != nil {
		response.NotifyWWError(w, request, http.StatusNotFound, "API-1066", err)
		return
	}

	bytes, err := json.Marshal(participants)
	response.Respond(w, http.StatusOK, bytes)
}*/

func (op ParticipantOperations) GetParticipantByDomain(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	participantId := vars["participant_id"]

	if participantId == "" {
		response.NotifyWWError(w, request, http.StatusBadRequest, "API-1041",
			errors.New("participant_id is required in path"))
		return
	}

	if err := validate.Pattern("participant_id", "path", string(participantId), `^[a-zA-Z0-9-]{5,32}$`); err != nil {
		response.NotifyWWError(w, request, http.StatusBadRequest, "API-1041",
			err)
		return
	}

	participant, err := op.RestPRServiceClient.GetParticipantForDomain(participantId)
	if err != nil {
		response.NotifyWWError(w, request, http.StatusNotFound, "API-1066", err)
		return
	}

	bytes, err := json.Marshal(participant)
	response.Respond(w, http.StatusOK, bytes)
}

/*func (op ParticipantOperations) GetAllParticipants(w http.ResponseWriter, request *http.Request) {

	participants, err := op.RestPRServiceClient.GetAllParticipants()
	if err != nil {
		response.NotifyWWError(w, request, http.StatusNotFound, "API-1114", err)
		return
	}

	bytes, err := json.Marshal(participants)
	response.Respond(w, http.StatusOK, bytes)
}*/

func (op ParticipantOperations) GetParticipantByQuery(w http.ResponseWriter, request *http.Request) {
	queryParams := request.URL.Query()
	countryCode := queryParams["country_code"]
	assetCode := queryParams["asset_code"]
	assetIssuer := queryParams["issuer_id"]

	if (assetCode != nil || assetIssuer != nil) && (assetCode == nil || assetIssuer == nil) {
		response.NotifyWWError(w, request, http.StatusNotAcceptable, "API-1042",
			errors.New("Both asset_code and issuer_id are required in query"))
		return
	}
	if (assetCode != nil && assetIssuer != nil) && ((assetCode[0] == "") || (assetIssuer[0] == "")) {
		response.NotifyWWError(w, request, http.StatusNotAcceptable, "API-1042",
			errors.New("Both asset_code and issuer_id are required in query"))
		return
	}

	if assetCode != nil {
		assetCodeStr := string(assetCode[0])
		err := model.IsValidAssetCode(assetCodeStr)
		if err != nil {
			LOGGER.Debug("asset code is invalid:", err)
			response.NotifyWWError(w, request, http.StatusBadRequest, "API-1124", err)
			return
		}
	}
	country := ""
	participants := []model.Participant{}
	participantResponse := []model.Participant{}
	err := errors.New("")

	if countryCode != nil {
		country = string(countryCode[0])
		participants, err = op.RestPRServiceClient.GetParticipantsByCountry(country)

		if err != nil {
			LOGGER.Debugf("Error during country_code query")
			response.NotifyWWError(w, request, http.StatusNotFound, "API-1119",
				err)
			return
		}
	} else {
		participants, err = op.RestPRServiceClient.GetAllParticipants()
		if err != nil {
			LOGGER.Debugf("Error during getting all participants query")
			response.NotifyWWError(w, request, http.StatusNotFound, "API-1119",
				err)
			return
		}
	}
	if assetCode != nil || assetIssuer != nil {
		if len(participants) > 0 {
			for i := 0; i < len(participants); i++ {

				//check for assets
				tcMatch := false
				//check if this participant is issuer
				if *participants[i].ID == assetIssuer[0] {
					LOGGER.Infof("Target asset issue account Match %v, %v", assetIssuer[0], participants[i].IssuingAccount)
					tcMatch = true
				}

				if tcMatch == false {
					// now check for trusted asset
					trustedAssets, err := apiutil.GetTrustedWWAssets(string(participants[i].IssuingAccount), op.RestPRServiceClient)
					if err != nil {
						LOGGER.Debugf(" Error calling GetTrustedWWAssets %v", err.Error())
					}
					for j := 0; j < len(trustedAssets); j++ {
						if trustedAssets[j].IssuerID == assetIssuer[0] &&
							*trustedAssets[j].AssetCode == assetCode[0] {
							LOGGER.Infof("Target asset trusted Asset Match %v", assetCode[0])
							tcMatch = true
						}
					}
				}

				if tcMatch == true {
					raw, err := json.Marshal(participants[i])
					if err != nil {
						LOGGER.Warningf("Error Marshaling query result", err)
					}
					if raw != nil {
						var participant model.Participant
						err := participant.UnmarshalBinary(raw)
						if err != nil {
							LOGGER.Warningf("Error unMarshaling query result", err)
						}

						participantResponse = append(participantResponse, participant)
					}
				}
			}
		}
	} else {
		participantResponse = participants
	}

	if len(participantResponse) <= 0 {
		response.NotifyWWError(w, request, http.StatusNotFound, "API-1120",
			errors.New("Given asset_code issuer_id and country in query"))
		return
	}

	bytes, err := json.Marshal(participantResponse)
	response.Respond(w, http.StatusOK, bytes)
}

func (op ParticipantOperations) GetAssetsForParticipant(w http.ResponseWriter, request *http.Request) {
	queryParams := request.URL.Query()
	vars := mux.Vars(request)
	domain := vars["participant_id"]
	if domain == "" {
		LOGGER.Warningf("Participant Domain is required in query")
		response.NotifyWWError(w, request, http.StatusBadRequest, "API-1132", nil)
		return
	}
	domainString := domain
	if strings.TrimSpace(domainString) == "" {
		LOGGER.Warningf("Participant Domain cannot be empty")
		response.NotifyWWError(w, request, http.StatusBadRequest, "API-1132", nil)
		return
	}
	if err := validate.Pattern("participant_id", "path", string(domainString), `^[a-zA-Z0-9-]{5,32}$`); err != nil {
		response.NotifyWWError(w, request, http.StatusBadRequest, "API-1041",
			err)
		return
	}
	queryType, found := queryParams["type"]
	if !found {
		LOGGER.Warningf("Type parameter is required in query")
		response.NotifyWWError(w, request, http.StatusBadRequest, "API-1254", nil)
		return
	}
	queryTypeString := queryType[0]
	if strings.TrimSpace(queryTypeString) == "" {
		LOGGER.Warningf("Type parameter cannot be empty")
		response.NotifyWWError(w, request, http.StatusBadRequest, "API-1254", nil)
		return
	}
	upperType := strings.ToUpper(strings.TrimSpace(queryTypeString))

	if upperType != TYPE_BOTH && upperType != TYPE_ISSUED && upperType != TYPE_TRUSTED {
		LOGGER.Warningf("Query Type must be one of: both, issued, trusted ")
		response.NotifyWWError(w, request, http.StatusBadRequest, "API-1255", nil)
		return
	}

	issuingDomain := domainString

	if upperType == TYPE_BOTH {
		assets, err := op.GetAllAssetsForParticipant(issuingDomain)
		if err != nil {
			newError := errors.New(err.Error() + "Error getting all assets for account: " + issuingDomain)
			LOGGER.Warningf("Error getting all assets for account: %v ", issuingDomain)
			response.NotifyWWError(w, request, http.StatusNotFound, "API-1258", newError)
			return
		}
		if assets == nil {
			as := model.Asset{}
			assets = []*model.Asset{&as}
		}
		assetBytes, _ := json.Marshal(assets)
		response.Respond(w, http.StatusOK, assetBytes)
		return
	}
	if upperType == TYPE_TRUSTED {
		assets, err := op.GetTrustedAssetsForParticipant(issuingDomain)
		if err != nil {
			newError := errors.New(err.Error() + "Error getting trusted assets for account: " + issuingDomain)
			LOGGER.Warningf("Error getting trusted assets for account: %v ", issuingDomain)
			response.NotifyWWError(w, request, http.StatusNotFound, "API-1258", newError)
			return
		}
		if assets == nil {
			as := model.Asset{}
			assets = []*model.Asset{&as}
		}
		assetBytes, _ := json.Marshal(assets)
		response.Respond(w, http.StatusOK, assetBytes)
		return
	}
	if upperType == TYPE_ISSUED {
		assets, err := op.GetIssuedAssetsForParticipant(issuingDomain)
		if err != nil {
			newError := errors.New(err.Error() + "Error getting issued assets for account: " + issuingDomain)
			LOGGER.Warningf("Error getting issued assets for account: %v ", issuingDomain)
			response.NotifyWWError(w, request, http.StatusNotFound, "API-1258", newError)
			return
		}
		if assets == nil {
			as := model.Asset{}
			assets = []*model.Asset{&as}
		}
		assetBytes, _ := json.Marshal(assets)
		response.Respond(w, http.StatusOK, assetBytes)
		return
	}
}

func (op ParticipantOperations) GetAllAssetsForParticipant(issuingDomain string) ([]*model.Asset, error) {
	var assets []*model.Asset
	issuedAssets, err := op.GetIssuedAssetsForParticipant(issuingDomain)
	if err != nil {
		return nil, err
	}
	trustedAssets, err := op.GetTrustedAssetsForParticipant(issuingDomain)
	if err != nil {
		return nil, err
	}
	if issuedAssets != nil {
		assets = append(assets, issuedAssets...)
	}
	if trustedAssets != nil {
		assets = append(assets, trustedAssets...)
	}
	if len(assets) == 0 {
		return nil, nil
	}
	return assets, nil
}

func (op ParticipantOperations) GetIssuedAssetsForParticipant(issuingDomain string) ([]*model.Asset, error) {
	var assets []*model.Asset

	//Get IBM admin account from crypto service
	wwAdminAccount, errorMsg, status, _ := op.CryptoServiceClient.GetIBMAccount()
	/*
		dev test
		Address := "GDMVNODEQVXQL4GN5FXYPQ3BH6ATTKLURXEENOJAXNFRZTCT2CGX4DHV"
		wwAdminAccount.Address = &Address
		status = http.StatusOK
		errorMsg = nil
	*/

	if status != http.StatusOK || errorMsg != nil {
		LOGGER.Errorf("Error getting IBM account: %v", errorMsg.Error())
		return nil, errorMsg
	}
	LOGGER.Infof("IBM Token Account: %s", wwAdminAccount)

	LOGGER.Infof("IBM Token Account: ", *wwAdminAccount.Address)

	if strings.TrimSpace(*wwAdminAccount.Address) == "" {
		newError := errors.New("No IBM Token Account in NodeConfig")
		return nil, newError
	}

	//Get trusted assets by IBM account
	wwAssets, err := apiutil.GetAssets(*wwAdminAccount.Address, op.RestPRServiceClient)
	if err != nil {
		LOGGER.Errorf("Error getting WW Assets", err)
		return nil, nil
	}
	if len(wwAssets) == 0 {
		return nil, nil
	}
	for _, ast := range wwAssets {
		if ast.IssuerID != "" && ast.IssuerID == issuingDomain {
			assets = append(assets, ast)
			continue
		} else {
			continue
		}
	}
	if len(assets) == 0 {
		return nil, nil
	}
	return assets, nil
}

func (op ParticipantOperations) GetTrustedAssetsForParticipant(domain string) ([]*model.Asset, error) {
	var assets []*model.Asset
	accountAddress, err := op.RestPRServiceClient.GetParticipantIssuingAccount(domain)
	if err != nil {
		return nil, err
	}
	assets, err = op.TrustedAssets(accountAddress)
	if err != nil {
		return nil, err
	}
	if len(assets) == 0 {
		return nil, nil
	}
	return assets, nil
}

/*
//GetWhiteListParticipants gets Participant domains in whitelist
func (op ParticipantOperations) GetWhiteListParticipants(w http.ResponseWriter, request *http.Request) {
	LOGGER.Infof("GetWhiteListParticipants")
	wwError, bytes := nc.GetWhiteListParticipants()
	if wwError.Code != "" {
		response.NotifyWWError(w, request, http.StatusNotFound, "API-"+wwError.Code, nil)
		return
	}

	response.Respond(w, http.StatusOK, bytes)
}

//AddWhiteListParticipant adds participant domain to whitelist
//add next
func (op ParticipantOperations) AddWhiteListParticipant(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	participantDomain := vars["participant_id"]

	LOGGER.Infof("AddWhiteListParticipant: participantDomain %v", participantDomain)

	if participantDomain == "" {
		response.NotifyWWError(w, request, http.StatusNotAcceptable, "API-1132",
			errors.New("Participant domain is required in query"))
		return
	}
	wwError := nc.AddWhiteListParticipant(participantDomain)
	if wwError.Code != "" {
		response.NotifyWWError(w, request, http.StatusNotAcceptable, "API-"+wwError.Code, nil)
		return
	}

	response.NotifySuccess(w, request, "Add WhiteList Participant OK")
}

//RemoveWhiteListParticipant -  Removes participant domain from whitelist
//add next
func (op ParticipantOperations) RemoveWhiteListParticipant(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	participantDomain := vars["participant_id"]

	LOGGER.Infof("RemoveWhiteListParticipant: participantDomain %v", participantDomain)

	if participantDomain == "" {
		response.NotifyWWError(w, request, http.StatusNotAcceptable, "API-1132",
			errors.New("Participant domain is required in query"))
		return
	}

	wwError := nc.RemoveWhiteListParticipant(participantDomain)
	if wwError.Code != "" {
		response.NotifyWWError(w, request, http.StatusNotAcceptable, "API-"+wwError.Code, nil)
		return
	}

	response.NotifySuccess(w, request, "Remove WhiteList Participant OK")
}

// IsParticipantWhiteListed - checks if a participant is in white list
func (op ParticipantOperations) IsParticipantWhiteListed(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	participantDomain := vars["participant_id"]

	LOGGER.Infof("IsParticipantWhiteListed: participantDomain %v", participantDomain)

	if participantDomain == "" {
		response.NotifyWWError(w, request, http.StatusNotAcceptable, "API-1132",
			errors.New("Participant domain is required in query"))
		return
	}

	isWhiteListed, wwError := nc.IsParticipantWhiteListed(participantDomain)
	if wwError.Code != "" {
		response.NotifyWWError(w, request, http.StatusNotFound, "API-"+wwError.Code, nil)
		return
	}

	response.NotifySuccess(w, request, strconv.FormatBool(isWhiteListed))
}
*/
// TrustedAssets - return trusted assets for the given accountAddres
func (op ParticipantOperations) TrustedAssets(accountAddress string) ([]*model.Asset, error) {

	//get trusted assets
	assets, err := apiutil.GetTrustedWWAssets(accountAddress, op.RestPRServiceClient)
	if err != nil {
		return nil, err
	}
	if len(assets) == 0 {
		assets = append(assets, &model.Asset{})
		return assets, nil
	}
	return assets, nil
}

func (op ParticipantOperations) GetAccount(w http.ResponseWriter, req *http.Request) {

	/*
		1. check if there is a Operating account with the given name, if yes, return Operating account address
		2. create new Operating account (in both stellar and NodeConfig, return the new Operating account address
	*/

	LOGGER.Infof("api-service:ParticipantOperations:GetAccount....")
	vars := mux.Vars(req)
	accountName := vars["account_name"]
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
	err = errors.New(accountName)
	response.NotifyWWError(w, req, http.StatusNotFound, "API-1055", err)
}

//Get account information from participant registry, parse it into list of accounts
func (op ParticipantOperations) GetAccounts(w http.ResponseWriter, req *http.Request) {

	LOGGER.Infof("api-service:ParticipantOperations:GetAccounts....")
	var accounts []model.Account
	participantObj, err := op.RestPRServiceClient.GetParticipantForDomain(os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME))
	if err != nil {
		LOGGER.Errorf(" Error GetParticipantForDomain failed: %v", err)
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1256", err)
		return
	}

	//add issuing account to list of accounts
	if participantObj.IssuingAccount != "" {
		issueAccount := model.Account{}
		issueAccount.Name = comn.ISSUING
		issueAccount.Address = &participantObj.IssuingAccount
		accounts = append(accounts, issueAccount)
	}

	for _, account := range participantObj.OperatingAccounts {
		accounts = append(accounts, *account)
	}

	if len(accounts) == 0 {
		err = errors.New(" no issuing or operating created for participant")
		LOGGER.Error(err)
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1070", err)
		return
	}

	acrJSON, _ := json.Marshal(accounts)
	response.Respond(w, http.StatusOK, acrJSON)

}
