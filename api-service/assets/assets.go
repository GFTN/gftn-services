// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package assets

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/GFTN/gftn-services/utility/blockchain-adaptor/ww_stellar"

	"github.com/GFTN/gftn-services/utility/wwfirebase"

	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"
	hClient "github.com/stellar/go/clients/horizonclient"
	"github.com/GFTN/gftn-services/api-service/environment"
	apiutil "github.com/GFTN/gftn-services/api-service/utility"
	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"
	crypto_client "github.com/GFTN/gftn-services/crypto-service-client/crypto-client"
	gasserviceclient "github.com/GFTN/gftn-services/gas-service-client"
	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/global-whitelist-service/whitelistclient"
	pr_client "github.com/GFTN/gftn-services/participant-registry-client/pr-client"
	"github.com/GFTN/gftn-services/utility/asset"
	ast "github.com/GFTN/gftn-services/utility/asset"
	"github.com/GFTN/gftn-services/utility/common"
	comn "github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/participant"
	"github.com/GFTN/gftn-services/utility/response"
	vauth "github.com/GFTN/gftn-services/utility/vault/auth"
	"github.com/GFTN/gftn-services/utility/vault/utils"
)

type AssetOperations struct {
	cryptoServiceClient crypto_client.CryptoServiceClient
	VaultSession        utils.Session
	PRServiceClient     pr_client.PRServiceClient
	GasServiceClient    gasserviceclient.GasServiceClient
}

func CreateAssetOperations() (AssetOperations, error) {
	ap := AssetOperations{}
	ap.VaultSession = utils.Session{}

	if os.Getenv(global_environment.ENV_KEY_SECRET_STORAGE_LOCATION) == common.VAULT_SECRET {
		//Vault location
		err := errors.New("")
		ap.VaultSession, err = vauth.GetSession()
		if err != nil {
			LOGGER.Errorf("Error reading account source environment settings")
			return ap, err
		}
	}
	var prClient pr_client.PRServiceClient
	if os.Getenv(environment.ENV_KEY_PARTICIPANT_REGISTRY_SERVICE_CLIENT) == "mock" {
		LOGGER.Warningf("USING MOCK PARTICIPANT REGISTRY SERVICE CLIENT")
		prClient = pr_client.MockPRServiceClient{}
	} else {
		LOGGER.Infof("Using REST Participant Registry Service Client")
		cl, _ := pr_client.CreateRestPRServiceClient(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL))
		prClient = cl
	}
	ap.PRServiceClient = prClient

	var cClient crypto_client.CryptoServiceClient
	if os.Getenv(environment.ENV_KEY_PARTICIPANT_REGISTRY_SERVICE_CLIENT) == "mock" {
		cClient, _ = crypto_client.CreateMockCryptoServiceClient()
	} else {
		err := errors.New("")
		cServiceInternalUrl, err := participant.GetServiceUrl(os.Getenv(global_environment.ENV_KEY_CRYPTO_SVC_INTERNAL_URL), os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME))
		if err != nil {
			return ap, err
		}
		cClient, err = crypto_client.CreateRestCryptoServiceClient(cServiceInternalUrl)
		if err != nil {
			return ap, err
		}
	}
	ap.cryptoServiceClient = cClient

	gasServiceClient := gasserviceclient.Client{
		HTTP: &http.Client{Timeout: time.Second * 20},
		URL:  os.Getenv(global_environment.ENV_KEY_GAS_SVC_URL),
	}
	ap.GasServiceClient = &gasServiceClient

	return ap, nil
}

func (ap AssetOperations) SuccessResponse(code, issuer, AssetType string) *model.Asset {
	asset := model.Asset{
		AssetCode: &code,
		IssuerID:  issuer,
		AssetType: &AssetType,
	}
	return &asset
}

/*
// FundAccountDeprecated - Fund (mint) asset - not for Digital Obligations
// this operation must be deprecated
// add next
func (ap AssetOperations) FundAccountDeprecated(w http.ResponseWriter, req *http.Request) {

	LOGGER.Infof("api-service:assets:assets.go:FundAccount....")
	LOGGER.Errorf("This operation is deprecated")
	response.NotifyWWError(w, req, http.StatusBadRequest, "API-1266", errors.New("FundAccount"))
	return

	var fundRequest model.FundRequest
	json.NewDecoder(req.Body).Decode(&fundRequest)
	LOGGER.Infof("Going to validate")
	fundRequest.Validate(strfmt.Default)

	distAccountName := *fundRequest.AccountName
	assetCode := *fundRequest.AssetCode
	amount := *fundRequest.Amount

	if asset.GetAssetType(assetCode) == model.AssetAssetTypeDO {
		LOGGER.Errorf("Cannot Fund Operating Account with DO: %v", assetCode)
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1244", nil)
		return
	}
	LOGGER.Infof("Params: distAccountName %v, assetCode %v, amount : %v", distAccountName, assetCode, amount)
	if len(distAccountName) == 0 || len(assetCode) == 0 || (amount) <= 0.0 {
		LOGGER.Errorf("Error: distAccountName, assetCode, amount must not be empty")
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1105", nil)
		return
	}

	//get DA from nodeconfig.toml
	accountExists, distAccount, err := nc.OperatingAccountExists(distAccountName)
	if err != nil {
		LOGGER.Errorf("Error while checking if operating account exists: %v", err)
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1095", nil)
		return
	}
	if !accountExists {
		//handle account exists
		LOGGER.Errorf("Operating Account Not Found - name: %v", distAccountName)
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1098", nil)
		return
	}

	//get IA from nodeconfig.toml
	issuingAccount, err := nc.GetIssuingAccount()
	if err != nil {
		LOGGER.Errorf("No Issuing account. Cannot create fund operating account", err)
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1097", nil)
		return
	}
	// check if trust exists for DA to asset
	trustedAssets, err := ast.GetTrustedWWAssets(distAccount.NodeAddress)
	trusted := false
	for _, ast := range trustedAssets {
		if *ast.AssetCode == assetCode {
			trusted = true
			break
		}
	}
	if !trusted {
		LOGGER.Errorf("No trust for this asset: %v", assetCode)
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1245", nil)
		return
	}

	LOGGER.Infof("Going to create payment")
	//submit payment from IA to DA
	exeStatus, respStr := ast.CreatePaymentInStellar(issuingAccount, distAccount.NodeAddress, assetCode, amount)
	LOGGER.Infof("Created Payment %v: ", exeStatus)
	if exeStatus {
		LOGGER.Infof(respStr)
		response.NotifySuccess(w, req, "Account has been funded")
	} else {
		LOGGER.Error(respStr)
		LOGGER.Error("Account has not been funded.")
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1106", nil)
	}
}*/

// IssueAsset -
// if asset type is DA, reject this request. Only Anchors can issue a Digital Asset
/*	Receives asset_code as the input. "asset_code" is a string parameter (example: USD, CAD, EUR NZD)
	Queries Issuing Account from Node Config file
	Queries IBM Token Account from Node Config file
	Invoke change trust operation using Stellar SDK. This will make the IBM Token Account trust the asset
	This asset will have issuing account as the issuer.
	After executing this function, results can be verified by invoking the following URL:
	https://horizon-testnet.stellar.org/accounts/{IBM_Token_account_id} (For test net)
	https://horizon.stellar.org/accounts/{IBM_Token_account_id} (For Pub net)
*/
func (ap AssetOperations) IssueAsset(w http.ResponseWriter, req *http.Request) {

	LOGGER.Infof("api-service:Asset Operations :IssueAsset")
	queryParams := req.URL.Query()

	assetCode := queryParams["asset_code"]
	assetType := queryParams["asset_type"]

	if assetCode == nil || assetType == nil || assetType[0] == "" || assetCode[0] == "" {
		LOGGER.Errorf("asset_code or asset_type missing in request")
		response.NotifyWWError(w, req, http.StatusBadRequest, "API-1124", nil)
		return
	}
	homeDomain := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	asst := model.Asset{}
	asst.AssetType = &assetType[0]
	asst.AssetCode = &assetCode[0]
	asst.IssuerID = homeDomain
	err := asst.Validate(strfmt.Default)

	if err != nil {
		LOGGER.Errorf("Error validating issue asset request: ", err.Error())
		response.NotifyWWError(w, req, http.StatusBadRequest, "API-1124", err)
		return
	}
	if assetType[0] != model.AssetAssetTypeDO {
		response.NotifyWWError(w, req, http.StatusBadRequest, "API-1265", errors.New(assetType[0]))
		return
	}
	if asset.GetAssetType(assetCode[0]) != assetType[0] {
		if assetType[0] == model.AssetAssetTypeDO {
			LOGGER.Errorf("Error: asset_code should end with DO")
			response.NotifyWWError(w, req, http.StatusBadRequest, "API-1124", errors.New(
				"for digial obligation, asset_code should end with DO"))
			return
		} else if assetType[0] == model.AssetAssetTypeNative {
			LOGGER.Errorf("Error validating issue asset request: Native asset cannot be issued ")
			response.NotifyWWError(w, req, http.StatusBadRequest, "API-1124", errors.New(
				" Native asset cannot be issued"))
			return
		} else {
			response.NotifyWWError(w, req, http.StatusBadRequest, "API-1124", errors.New(
				"for digital asset, asset_code should not end with DO"))
			return
		}
		return
	}

	err = model.IsValidAssetCode(*asst.AssetCode)
	if err != nil {
		LOGGER.Debug("asset code is invalid:", err)
		response.NotifyWWError(w, req, http.StatusBadRequest, "API-1124", err)
		return
	}

	if os.Getenv(global_environment.ENV_KEY_ENABLE_JWT) != "false" {
		//For issuing asset authenticated token should have permission to use issuing account
		if !middlewares.HasAccount(comn.ISSUING, req) {
			response.NotifyWWError(w, req, http.StatusUnauthorized, "API-1267",
				errors.New("issuing Account name is not same as authenticated account token for issuing asset"))
			return
		}
	}

	issuingAccount, err := participant.GenericGetAccount(ap.VaultSession, comn.ISSUING)
	if err != nil {
		LOGGER.Errorf("Error while checking if issuing account exists: %v", err)
		LOGGER.Errorf("The Asset could NOT be issued due to error retrieving Issuing Account")
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1054", nil)
		return
	}

	limit := "1"
	LOGGER.Infof("Issuer: %s", issuingAccount.NodeAddress)
	LOGGER.Infof("Asset Code received : %v", assetCode)

	wwAdminAccount, errorMsg, status, _ := ap.cryptoServiceClient.GetIBMAccount()
	if status != http.StatusOK || errorMsg != nil {
		LOGGER.Errorf("Error getting IBM account: %v", errorMsg.Error())
		return
	}
	LOGGER.Infof("IBM Token Account: %s", wwAdminAccount)

	xdr, err := asset.ChangeTrust(ap.GasServiceClient,
		*wwAdminAccount.Address,
		issuingAccount.NodeAddress,
		assetCode[0],
		limit)
	xdrB, _ := base64.StdEncoding.DecodeString(xdr)

	//Get signed by IBM account on crypto service
	sigXdr, errorMsg, status, code := ap.cryptoServiceClient.AddIBMSign(xdrB)

	if status != http.StatusOK {
		LOGGER.Errorf("Error issuing new asset: %v", errorMsg.Error())
		code = "API-1090"
		response.NotifyWWError(w, req, http.StatusConflict, code, errors.New("Error issuing new asset"))
		return
	}
	if errorMsg != nil {
		LOGGER.Errorf("Error issuing new asset: %v", errorMsg.Error())
		code = "API-1090"
		response.NotifyWWError(w, req, http.StatusConflict, code, errorMsg)
		return
	}
	LOGGER.Debugf("signed transaction: %v", sigXdr)

	//submit transaction to stellar
	//step 3: create new operating account  in stellar
	b64Xdr := base64.StdEncoding.EncodeToString(sigXdr)
	LOGGER.Debugf("signed transaction: %v", b64Xdr)

	//Submit on Gas service
	hash, ledger, err := ap.GasServiceClient.SubmitTxe(b64Xdr)
	if err != nil {
		err = ast.DecodeStellarError(err)
		LOGGER.Error("The Asset could not be issued. Error Communicating with Stellar.")
		response.NotifyWWError(w, req, http.StatusInternalServerError, "API-1008", err)
		return
	}
	LOGGER.Debugf("submitTransaction  %v, %v", hash, ledger)

	//Asset type is credit since we are not issuing any native asset
	asType := *asst.AssetType
	assetBytes, _ := json.Marshal(ap.SuccessResponse(assetCode[0], homeDomain, asType))
	response.Respond(w, http.StatusOK, assetBytes)

}

//GetAB : get asset balance
func (ap AssetOperations) GetAB(req *http.Request) (string, string, error) {
	queryValues := req.URL.Query()
	assetCode := strings.TrimSpace(queryValues.Get("asset_code"))
	assetIssuer := strings.TrimSpace(queryValues.Get("issuer_id"))

	LOGGER.Infof("Going to validate")
	if assetCode == "" || (!ast.IsNative(assetCode) && assetIssuer == "") {
		return assetCode, assetIssuer, errors.New("mandatory request parameters are missing")
	}
	LOGGER.Infof("AssetCode: %s", assetCode)
	LOGGER.Infof("IssuerID: %s", assetIssuer)

	participant, err := getParticipantForDomain(assetIssuer)
	if err != nil {
		LOGGER.Errorf("Error while getting Participant for Domain: Participant Registry not available", err.Error())
		return assetCode, assetIssuer, errors.New("error querying PR for given issuer ID")
	}
	return assetCode, participant.IssuingAccount, nil
}

func (ap AssetOperations) AssetBalance(w http.ResponseWriter, req *http.Request) {

	LOGGER.Infof("Inside AssetBalance Func")
	vars := mux.Vars(req)
	accountName := vars["account_name"]
	assetCode, assetIssuer, err := ap.GetAB(req)
	if err != nil {
		response.NotifyWWError(w, req, http.StatusBadRequest, "API-1100", nil)
		return
	}
	balStr := ""

	account, err := participant.GenericGetAccount(ap.VaultSession, accountName)
	if err == nil {
		balStr = ast.GetBalance(account.NodeAddress, assetCode, assetIssuer)
	} else {
		LOGGER.Errorf("There is no matching account")
		response.NotifyWWError(w, req, http.StatusBadRequest, "API-1055", nil)
		return
	}

	/*
		Stellar returns "0" if there is no match for the given issuer and asset code. If there is a match but contains zero balance
		stellar returns "0.0000000".
	*/

	LOGGER.Infof("Balance string: %s", balStr)
	if "" == balStr || "0" == balStr {
		LOGGER.Errorf("There is no trust relation for account and asset")
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1056", nil)
	} else {

		prURL := os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL)
		prc, _ := pr_client.CreateRestPRServiceClient(prURL)
		participant, err := prc.GetParticipantForIssuingAddress(assetIssuer)
		if err != nil {
			LOGGER.Errorf("There is no matching participant")
			response.NotifyWWError(w, req, http.StatusBadRequest, "API-1058", nil)
			return
		}
		assetBalance := model.AssetBalance{
			AccountName: &accountName,
			Balance:     &balStr,
			AssetCode:   &assetCode,
			IssuerID:    *participant.ID,
		}
		balanceBytes, _ := json.Marshal(assetBalance)
		response.Respond(w, http.StatusOK, balanceBytes)
	}

}

// IssuedAssets -
/*	Get a list of all the assets that have been issued by this participant on World Wire
 */
func (ap AssetOperations) IssuedAssets(w http.ResponseWriter, req *http.Request) {

	// first, get the IBM Token Account
	var assets []*model.Asset
	wwAdminAccount, errorMsg, status, _ := ap.cryptoServiceClient.GetIBMAccount()

	if status != http.StatusOK || errorMsg != nil {
		LOGGER.Errorf("Error getting IBM account: %v", errorMsg.Error())
		return
	}

	LOGGER.Infof("IBM Token Account: %s", *wwAdminAccount.Address)

	if strings.TrimSpace(*wwAdminAccount.Address) == "" {
		LOGGER.Errorf("No IBM Token Account", errorMsg)
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1230", errorMsg)
		return
	}

	homeDomain := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	if strings.TrimSpace(homeDomain) == "" {
		LOGGER.Errorf("Participant does not have home domain")
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1253", errors.New("participant does not have home domain"))
		return
	}

	//get trusted assets trusted by IBM Token Account
	wwAssets, err := apiutil.GetAssets(*wwAdminAccount.Address, ap.PRServiceClient)
	if err != nil {
		LOGGER.Errorf("Error getting WW Assets", err)
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1068", err)
		return
	}

	if len(wwAssets) == 0 {

		assets = append(assets, &model.Asset{})
		assetBytes, _ := json.Marshal(assets)
		response.Respond(w, http.StatusOK, assetBytes)
		return
	}

	for _, ast := range wwAssets {

		if ast.IssuerID == homeDomain {

			assets = append(assets, ast)
			continue
		}
	}

	if len(assets) == 0 {
		assets = append(assets, &model.Asset{})
		assetBytes, _ := json.Marshal(assets)
		response.Respond(w, http.StatusOK, assetBytes)
		return
	}

	assetBytes, _ := json.Marshal(assets)
	response.Respond(w, http.StatusOK, assetBytes)
	return

}

// AssetBalances - Get input parameters (assetCode, assetIssuer) from request
/*	Load issuer from node config and compare it with the issuer in request parameter.
	Don't proceed if the issuers are different
	Query the stellar db for all accounts with non-zero balance for the asset
*/
func (ap AssetOperations) AssetBalances(w http.ResponseWriter, req *http.Request) {

	LOGGER.Infof("AssetBalances Func")
	queryValues := req.URL.Query()
	assetCode := strings.TrimSpace(queryValues.Get("asset_code"))

	account, err := participant.GenericGetAccount(ap.VaultSession, comn.ISSUING)
	if err != nil {
		LOGGER.Errorf("Error while checking if issuing account exists: %v", err)
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1102", err)
		return
	}
	assetIssuer := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	prc, _ := pr_client.CreateRestPRServiceClient(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL))
	participant, err := prc.GetParticipantForDomain(assetIssuer)
	if err != nil {
		LOGGER.Errorf("Error while checking if issuing account exists: %v", err)
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1102", err)
		return
	}

	//Check if secrets manager and PR are aligned
	if participant.IssuingAccount != account.NodeAddress {
		LOGGER.Errorf("Error while checking if issuing account correct: %v", err)
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1056", err)
		return
	}

	LOGGER.Debug("AssetCode: %s", assetCode, account.NodeAddress)

	client := common.GetNewHorizonClient(os.Getenv(global_environment.ENV_KEY_HORIZON_CLIENT_URL))

	assetRequest := hClient.AssetRequest{ForAssetCode: assetCode,
		ForAssetIssuer: account.NodeAddress}

	// Load the asset List from the network
	assetsList, err := client.Assets(assetRequest)
	if err != nil {
		LOGGER.Errorf("Error loading asset list", err)
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1309", err)
		return
	}
	LOGGER.Debugf("assets:", assetsList)

	/*
		create list of Asset Balances by iterating over the result set and return the list
	*/

	var assets []*model.AssetBalance
	for _, b := range assetsList.Embedded.Records {
		LOGGER.Debug(b.Code, b.Issuer)
		assetCode := b.Code
		//Regular participants return only DO balances, doesn't include DAs
		if ast.GetAssetType(assetCode) == model.AssetAssetTypeDO {
			issuer := b.Issuer
			balance, _ := strconv.ParseFloat(b.Amount, 32)
			balanceString := b.Amount
			LOGGER.Info("balance    : %s", balance)
			LOGGER.Info("issuer     : %s", assetIssuer)

			accountName := ""
			if issuer == participant.IssuingAccount {
				accountName = comn.ISSUING
			}

			if accountName == "" {
				response.NotifyWWError(w, req, http.StatusNotFound, "API-1070", err)
				return
			}

			assets = append(assets, &model.AssetBalance{
				AccountName: &accountName,
				Balance:     &balanceString,
				AssetCode:   &assetCode,
				IssuerID:    assetIssuer,
			})
		}
	}

	if len(assets) > 0 {
		balanceBytes, _ := json.Marshal(assets)
		response.Respond(w, http.StatusOK, balanceBytes)
		return
	}
	response.NotifyWWError(w, req, http.StatusNotFound, "API-1126", err)
}

// TrustedAssets - return trusted assets for the given accountAddres
func (ap AssetOperations) TrustedAssets(accountAddress string) ([]*model.Asset, error) {
	//get trusted assets

	assets, err := apiutil.GetTrustedWWAssets(accountAddress, ap.PRServiceClient)
	if err != nil {
		return nil, err
	}
	if len(assets) == 0 {
		assets = append(assets, &model.Asset{})
		return assets, nil
	}
	return assets, nil
}

//TrustedAssetsForAccount : get trusted assets for a operating account.
func (ap AssetOperations) TrustedAssetsForAccount(w http.ResponseWriter, req *http.Request) {
	LOGGER.Infof("api-service:assets:assets.go:TrustedAssets for operating account....")
	var assets []*model.Asset
	vars := mux.Vars(req)
	distAccountName := vars["account_name"]
	LOGGER.Infof("account_name in request: %s", distAccountName)
	//get Operating Account from nodeconfig.toml
	LOGGER.Infof("Getting trusted assets for Operating Account")

	account, err := participant.GenericGetAccount(ap.VaultSession, distAccountName)
	if err != nil {
		LOGGER.Errorf("Error while checking if operating account exists: %v", err)
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1067", err)
		return
	}

	// dev test
	//account.NodeAddress = "GCLKOYF4AJ2EXVRYJUNYLU5RYBJDT4N2DQEGVP7UL247UJAQDWZVDBZM"

	LOGGER.Infof("Getting trusted assets for Operating account %v", account.NodeAddress)
	assets, err = ap.TrustedAssets(account.NodeAddress)
	if err != nil {
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1067", err)
		return
	}
	if len(assets) == 0 {
		assets = append(assets, &model.Asset{})
	}
	b, marshalErr := json.Marshal(assets)
	if marshalErr != nil {
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1257", err)
	}
	response.Respond(w, http.StatusOK, b)
}

/*
TrustedAssetsForIA - get trusted assets for a issuing account.
*/
func (ap AssetOperations) TrustedAssetsForIA(w http.ResponseWriter, req *http.Request) {
	LOGGER.Infof("api-service:assets:assets.go:TrustedAssets for issuing account....")
	var assets []*model.Asset
	//get IA from nodeconfig.toml
	LOGGER.Infof("Getting trusted assets for IA")
	issuingAccount, err := participant.GenericGetAccount(ap.VaultSession, comn.ISSUING)
	if err != nil {
		LOGGER.Errorf("Issuing account not found", err)
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1067", err)
	}
	LOGGER.Infof("Getting trusted assets for Issuing account %v", issuingAccount.NodeAddress)
	assets, err = ap.TrustedAssets(issuingAccount.NodeAddress)
	if err != nil {
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1067", err)
		return
	}
	if len(assets) == 0 {
		assets = append(assets, &model.Asset{})
	}
	b, marshalErr := json.Marshal(assets)
	if marshalErr != nil {
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1257", err)
	}
	response.Respond(w, http.StatusOK, b)
}

// WorldWireAssets : handler to process WorldWireAssets
func (ap AssetOperations) WorldWireAssets(w http.ResponseWriter, req *http.Request) {
	var assets []*model.Asset
	LOGGER.Infof("api-service:Asset Operations :WorldWireAssets")
	wwAdminAccount, errorMsg, status, _ := ap.cryptoServiceClient.GetIBMAccount()
	/*
		dev test
		Address := "GDMVNODEQVXQL4GN5FXYPQ3BH6ATTKLURXEENOJAXNFRZTCT2CGX4DHV"
		wwAdminAccount.Address = &Address
		status = http.StatusOK
		errorMsg = nil
	*/
	if status != http.StatusOK || errorMsg != nil {
		LOGGER.Errorf("Error getting IBM account: %v", errorMsg.Error())
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1068", errorMsg)
		return
	}
	LOGGER.Infof("IBM Token Account: %s", *wwAdminAccount.Address)

	assets, err := apiutil.GetAssets(*wwAdminAccount.Address, ap.PRServiceClient)
	if err != nil {
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1068", err)
		return
	}
	if len(assets) == 0 {
		assets = append(assets, &model.Asset{})
	}
	b, marshalErr := json.Marshal(assets)
	if marshalErr != nil {
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1257", err)
	}
	response.Respond(w, http.StatusOK, b)
}

//	if the participant is not on the white list, reject the allowTrust operation

func (ap AssetOperations) allowTrust(trustRequest model.Trust, assetCode string, authorize bool, httpMethod string) error {
	LOGGER.Debugf("api-service:assets:assets.go:allowTrust....")
	LOGGER.Debugf("Note: Issuer Account should have auth_required: true, auth_revocable: true")
	LOGGER.Debugf("Note: Source Account should trust the asset first")

	LOGGER.Debugf("allowTrust: Going to validate")
	err := trustRequest.Validate(strfmt.Default)
	if err != nil {
		msg := "Unable to validate allow trust request: " + err.Error()
		LOGGER.Debugf(msg)
		return errors.New("API-1110")
	}

	// domain is the domain of the trustor
	domain := *trustRequest.ParticipantID
	accountName := *trustRequest.AccountName
	LOGGER.Debugf("allowTrust : domain %v, accountName %v ", domain, accountName)

	if len(domain) == 0 || len(accountName) == 0 {
		LOGGER.Errorf("Error: domain, accountName must not be empty")
		return errors.New("API-1101")
	}

	// if the domain is not local and the participant is not on the whitelist, reject this request
	client := whitelistclient.Client{
		HTTPClient: &http.Client{Timeout: time.Second * 10},
		WLURL:      os.Getenv(global_environment.ENV_KEY_WL_SVC_URL),
	}
	homeDomain := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	isWhitelisted, err2 := client.IsParticipantWhiteListed(homeDomain, domain)
	if err2 != nil {
		LOGGER.Errorf("Error: %s", err)
		return errors.New("API-FITOFICCT-1025")
	}

	if !isLocalDomain(domain) && !isWhitelisted {
		LOGGER.Errorf("allowTrust: Participant domain: %v is not in the whitelist", domain)
		return errors.New("API-FITOFICCT-1025")
	}

	// if asset type is DO, then domain cannot be local, else it is OK
	isLocal := isLocalDomain(domain)
	if isLocal && ast.GetAssetType(assetCode) == model.AssetAssetTypeDO {
		LOGGER.Errorf("Error: domain name for allow trust should not be local if asset type is DO")
		return errors.New("API-1242")
	}

	// if the asset type is a DO, then account name must be an issuing account
	if ast.GetAssetType(assetCode) == model.AssetAssetTypeDO && accountName != comn.ISSUING {
		LOGGER.Errorf("Operating Account: %v cannot enter into a trust relationship for a DO: %v", accountName, assetCode)
		return errors.New("API-1248")
	}

	participantObj, err := getParticipantForDomain(domain)
	if err != nil {
		LOGGER.Errorf("Error while getting Participant for Domain: Participant Registry not available", err.Error())
		return errors.New("API-1241")
	}

	// Check if participant(trustor) is active before trusting
	LOGGER.Info("Check participant active")
	err = participant.CheckStatusActive(participantObj)
	if err != nil {
		LOGGER.Error(err.Error())
		return errors.New("API-1104")
	}

	//find  the account for which trust is to be allowed
	// if domain is local, then account must be a local operating account
	// if the domain is remote, then the account could be either a remote issuing account or remote operating account
	trustorAddress := ""
	if isLocal {
		trustorAccount, _ := participant.GenericGetAccount(ap.VaultSession, accountName)
		trustorAddress = trustorAccount.NodeAddress
	} else {
		if accountName == "" || accountName == comn.ISSUING {
			trustorAddress = participantObj.IssuingAccount
		} else {
			for _, distAccount := range participantObj.OperatingAccounts {
				if distAccount.Name == accountName {
					trustorAddress = *distAccount.Address
					break
				}
			}
		}

	}

	if trustorAddress == "" {
		LOGGER.Errorf("Invalid Trustor Account for Allow Trust")
		return errors.New("API-1251")
	}

	issuingAccount, err := participant.GenericGetAccount(ap.VaultSession, comn.ISSUING)
	if err != nil {
		msg := "No Issuing account found"
		LOGGER.Error(msg)
		return errors.New("API-1202")
	}
	err, errMsg := ast.DoAllowTrust(ap.GasServiceClient, trustorAddress, assetCode, issuingAccount.NodeAddress, authorize, ap.cryptoServiceClient)

	LOGGER.Debugf(errMsg)
	if err == nil {
		if httpMethod == "PUT" {
			LOGGER.Debugf("modified a allowtrust for the given issued asset")
		} else {
			LOGGER.Debugf("created a allowtrust for the given issued asset")
		}
		return nil
	}
	return errors.New("API-" + err.Error())

}
func isLocalDomain(domain string) bool {
	homeDomain := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	return domain == homeDomain
}

// getParticipantForDomain : Get participant for domain
func getParticipantForDomain(domain string) (model.Participant, error) {
	LOGGER.Debugf("getParticipantForDomain domain = %v", domain)
	var participantObj model.Participant

	prClient, err := pr_client.CreateRestPRServiceClient(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL))
	//prClient, err := pr_client.CreateMockPRServiceClient()
	if err != nil {
		LOGGER.Errorf(" Error getParticipantForDomain CreateRestPRServiceClient failed  %v", err)
		return participantObj, err
	}
	participantObj, err = prClient.GetParticipantForDomain(domain)
	if err != nil {
		LOGGER.Errorf(" Error getParticipantForDomain failed: %v", err)
		return participantObj, err
	}

	return participantObj, err
}

// perform ChangeTrust operation
// if the domain is not local and if the domain is not on the whitelist, reject this request
// if the asset is a digital obligation, then trust can only be between 2 issuing accounts
// if the asset is a digital asset, then there is no restriction on the accounts
func (ap AssetOperations) changeTrust(trustRequest model.Trust, assetCode string, httpMethod string) error {

	err := trustRequest.Validate(strfmt.Default)
	if err != nil {
		msg := "Unable to validate change trust request: " + err.Error()
		LOGGER.Debugf(msg)
		return errors.New("API-1250")
	}

	if ast.GetAssetType(assetCode) != model.AssetAssetTypeDO && ast.GetAssetType(assetCode) != model.AssetAssetTypeDA {
		msg := "Error Asset code is not DO nor DA: " + assetCode
		LOGGER.Error(msg)
		return errors.New("API-1201")
	}
	/*
		domain: this is the domain of the asset issuer
		accountName: this is the source account (trustor) for the change trust
	*/
	limit := trustRequest.Limit
	limitString := common.IntToString(limit)
	LOGGER.Infof("Limit amount is: %v", limitString)
	domain := *trustRequest.ParticipantID
	sourceAccountName := *trustRequest.AccountName
	if sourceAccountName == "" {
		sourceAccountName = comn.ISSUING
	}
	//if the dmain is not a local domain and if the domain is not on the whitelist, reject this request
	client := whitelistclient.Client{
		HTTPClient: &http.Client{Timeout: time.Second * 10},
		WLURL:      os.Getenv(global_environment.ENV_KEY_WL_SVC_URL),
	}
	homeDomain := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	isWhitelisted, _ := client.IsParticipantWhiteListed(homeDomain, domain)
	if !isLocalDomain(domain) && !isWhitelisted {
		LOGGER.Errorf("changeTrust: Participant domain: %v is not in the whitelist ", domain)
		return errors.New("API-FITOFICCT-1025")
	}
	// if the asset type is DO and the domain is a local domain, reject this request
	if ast.GetAssetType(assetCode) == model.AssetAssetTypeDO && isLocalDomain(domain) {
		LOGGER.Errorf("Domain: %v cannot be home domain if asset: %v is a DO ", domain, assetCode)
		return errors.New("API-1247")
	}

	// if the asset type is DO and the account  is a Operating account, reject this request
	if ast.GetAssetType(assetCode) == model.AssetAssetTypeDO && sourceAccountName != comn.ISSUING {
		LOGGER.Errorf("Operating Account: %v cannot enter into a trust relationship for a DO: %v ", sourceAccountName, assetCode)
		return errors.New("API-1248")
	}

	if sourceAccountName != comn.ISSUING {

		acc, err := participant.GenericGetAccount(utils.Session{}, common.ISSUING)
		if err != nil {
			LOGGER.Errorf("Encounter error while finding issuing account")
			return errors.New("API-1054")
		}

		prc, _ := pr_client.CreateRestPRServiceClient(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL))

		assets, err := apiutil.GetTrustedWWAssets(acc.NodeAddress, prc)
		if err != nil {
			LOGGER.Errorf("Encounter error while retrieving trusted assets")
			return errors.New("API-1067")
		}

		exists := false
		for _, asset := range assets {
			if *asset.AssetCode == assetCode && asset.IssuerID == domain {
				exists = true
				break
			}
		}

		if !exists {
			LOGGER.Errorf("Issuing account should be trusted first")
			return errors.New("API-1276")
		}
	}

	// get the participant for the given domain
	participantObj, err := getParticipantForDomain(domain)
	if err != nil {
		LOGGER.Errorf("Invalid Participant Domain: %v", domain)
		return errors.New("API-1136")
	}

	// Check if participant(trustor) is active before trusting
	LOGGER.Info("Check participant active")
	err = participant.CheckStatusActive(participantObj)
	if err != nil {
		LOGGER.Error(err.Error())
		return errors.New("API-1104")
	}

	// get the issuing address for the asset
	assetIssuer := participantObj.IssuingAccount

	// ensure that the asset has already been issued in world wire
	wwAdminAccount, errorMsg, status, _ := ap.cryptoServiceClient.GetIBMAccount()
	if status != http.StatusOK || errorMsg != nil {
		LOGGER.Errorf("Error getting IBM account: %v", errorMsg.Error())
		return errors.New("API-1230")
	}
	LOGGER.Infof("IBM Token Account: %s", *wwAdminAccount.Address)
	if !ast.IsBalanceExist(*wwAdminAccount.Address, assetCode, assetIssuer) {
		LOGGER.Errorf("Asset was not issued by IBM: %v", assetCode)
		return errors.New("API-1231")
	}

	// get the source address (trustor address) for the change trust
	sourceAccount, err := participant.GenericGetAccount(ap.VaultSession, sourceAccountName)
	if err != nil {
		LOGGER.Errorf("Invalid Account Name: %v", sourceAccountName)
		return errors.New("API-1249")
	}

	LOGGER.Debugf("changeTrust:trust request : domain %v, limit %v, accountName %v ", domain, limit, sourceAccount)
	LOGGER.Infof("Source address for trust is: %v", sourceAccount.NodeAddress)

	//Pass the trustor account information to change trust operation
	trustorAccount := model.Account{Name: sourceAccountName, Address: &sourceAccount.NodeAddress}
	err, errMsg := ast.DoChangeTrust(ap.GasServiceClient, trustorAccount, assetCode, assetIssuer, limit, ap.cryptoServiceClient)

	LOGGER.Debugf(errMsg)
	if err == nil {
		if httpMethod == "PUT" {
			LOGGER.Debugf("modified a changetrust for the given DO issued asset")
		} else {
			LOGGER.Debugf("created a changetrust for the given DO issued asset")
		}
	} else {
		return errors.New("API-" + err.Error())
	}
	return nil
}

//CreateOrAllowTrust - create trust for an asset
func (ap AssetOperations) CreateOrAllowTrust(w http.ResponseWriter, req *http.Request) {
	LOGGER.Debugf("api-service:assets:assets.go:CreateOrAllowTrust....")

	var trustRequest model.Trust
	err := json.NewDecoder(req.Body).Decode(&trustRequest)
	if err != nil {
		LOGGER.Debugf("Error  %v", err.Error())
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1042", nil)
		return
	}
	err = trustRequest.Validate(strfmt.Default)
	if err != nil {
		LOGGER.Debugf("Error  %v", err.Error())
		response.NotifyWWError(w, req, http.StatusNotFound, "API-1042", nil)
		return
	}

	permission := *trustRequest.Permission
	assetCode := *trustRequest.AssetCode

	if os.Getenv(global_environment.ENV_KEY_ENABLE_JWT) != "false" {
		participantID, err := middlewares.GetIdentity(req)
		//skip checking auth token for dev envs
		//Check if requested ofi account is the same as account token
		if participantID != os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME) {
			response.NotifyWWError(w, req, http.StatusUnauthorized, "API-1267",
				err)
			return
		}
	}

	if permission == "request" {
		if os.Getenv(global_environment.ENV_KEY_ENABLE_JWT) != "false" {
			//For change trust authenticated token should have permission to use given account name in the request
			if !middlewares.HasAccount(*trustRequest.AccountName, req) {
				err = errors.New("given Account name is not same as authenticated account in change trust request")
				response.NotifyWWError(w, req, http.StatusUnauthorized, "API-1267",
					err)
				return
			}
		}
		if trustRequest.Limit == 0 {
			//set the default limit
			trustRequest.Limit = common.DefaultTrustLimit
		}
		err := ap.changeTrust(trustRequest, assetCode, req.Method)
		if err != nil {
			LOGGER.Debugf("Error  %v", err.Error())
			response.NotifyWWError(w, req, http.StatusNotFound, err.Error(), errors.New(*trustRequest.ParticipantID))
			return
		}
		response.NotifySuccess(w, req, "Request Trust OK")
		err = wwfirebase.SendFBTrustSuccess(trustRequest)
		if err != nil {
			LOGGER.Debug("Update to Firebase db failed: ", err.Error())
		}
		return
	} else if permission == "allow" || permission == "revoke" {
		if os.Getenv(global_environment.ENV_KEY_ENABLE_JWT) != "false" {
			//For allow trust authenticated token should have permission to use issuing account
			err = errors.New("issuing Account name is not same as authenticated account for allow trust request")
			if !middlewares.HasAccount(comn.ISSUING, req) {
				response.NotifyWWError(w, req, http.StatusUnauthorized, "API-1267",
					err)
				return
			}
		}

		if permission == "allow" {
			err := ap.allowTrust(trustRequest, assetCode, true, req.Method)
			if err != nil {
				response.NotifyWWError(w, req, http.StatusNotFound, err.Error(), errors.New(*trustRequest.ParticipantID))
				return
			}
			response.NotifySuccess(w, req, "Allow Trust OK")
			err = wwfirebase.SendFBTrustSuccess(trustRequest)
			if err != nil {
				LOGGER.Debug("Update to Firebase db failed: ", err.Error())
			}
			return
		}
		err := ap.allowTrust(trustRequest, assetCode, false, req.Method)
		if err != nil {
			response.NotifyWWError(w, req, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.NotifySuccess(w, req, "Revoke Trust OK")
		err = wwfirebase.SendFBTrustSuccess(trustRequest)
		if err != nil {
			LOGGER.Debug("Update to Firebase db failed: ", err.Error())
		}
		return
	}
	err = errors.New("permission can be request, allow or revoke")
	response.NotifyWWError(w, req, http.StatusNotFound, "API-1042", err)
}

// GetOutstandingBalance : return outstanding balance of an asset issued by the caller
func (ap AssetOperations) GetOutstandingBalance(w http.ResponseWriter, req *http.Request) {

	LOGGER.Infof("GetOutstandingBalance...")
	vars := mux.Vars(req)
	assetCode := vars["asset_code"]
	//check if asset_code is valid
	err := model.IsValidAssetCode(assetCode)
	if err != nil {
		LOGGER.Debug("asset code is invalid:", err)
		response.NotifyWWError(w, req, http.StatusBadRequest, "API-1201", err)
		return
	}
	// establish the asset
	assetIssuer := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	asset := model.Asset{}
	asset.AssetCode = &assetCode
	asset.IssuerID = assetIssuer
	// get the balances
	nhc := common.GetNewHorizonClient(os.Getenv(global_environment.ENV_KEY_HORIZON_CLIENT_URL)) // new horizon client
	outstandingBalances, err := ww_stellar.GetOutstandingBalances(asset, ap.PRServiceClient, &nhc)
	if err != nil {
		LOGGER.Debug("Error ", err)
		if err.(ww_stellar.Error).Code() == ww_stellar.ERROR_ASSET_DOES_NOT_EXIST {
			response.NotifyWWError(w, req, http.StatusNotFound, "API-1401", nil)
			return
		}
		response.NotifyFailure(w, req, http.StatusInternalServerError, "internal error")
		return
	}

	resultOutstandingBalances := []model.Obligation{}
	// mapping to output object type
	for _, outstandingBalance := range outstandingBalances {
		resultOutstandingBalance := model.Obligation{}
		assetBalance := model.AssetBalance{}
		resultOutstandingBalance.Balance = &assetBalance
		resultOutstandingBalance.Balance.AccountName = common.StrPtr(outstandingBalance.Account.Account)
		resultOutstandingBalance.Balance.IssuerID = assetIssuer
		resultOutstandingBalance.ParticipantID = outstandingBalance.Account.ParticipantID
		resultOutstandingBalance.Balance.AssetCode = common.StrPtr(assetCode)
		resultOutstandingBalance.Balance.Balance = common.StrPtr(outstandingBalance.Balance.Amount.String())
		resultOutstandingBalances = append(resultOutstandingBalances, resultOutstandingBalance)
	}
	balancesBytes, _ := json.Marshal(resultOutstandingBalances)
	response.Respond(w, http.StatusOK, balancesBytes)
	return
}
