// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"
	b "github.com/stellar/go/build"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/xdr"
	gasserviceclient "github.com/GFTN/gftn-services/gas-service-client"
	"github.com/GFTN/gftn-services/gftn-models/model"
	pr_client "github.com/GFTN/gftn-services/participant-registry-client/pr-client"
	"github.com/GFTN/gftn-services/utility/asset"
	ast "github.com/GFTN/gftn-services/utility/asset"
	comn "github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/participant"
	"github.com/GFTN/gftn-services/utility/response"
	vauth "github.com/GFTN/gftn-services/utility/vault/auth"
	vutils "github.com/GFTN/gftn-services/utility/vault/utils"
)

type OnBoardingHandler struct {
	prClient         pr_client.RestPRServiceClient
	VaultSession     vutils.Session
	GasServiceClient gasserviceclient.GasServiceClient
}

func CreateOnBoardingHandler() (OnBoardingHandler, error) {
	oh := OnBoardingHandler{}
	prClient, err := pr_client.CreateRestPRServiceClient(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL))
	if err != nil {
		LOGGER.Errorf(" Error getParticipantForDomain CreateRestPRServiceClient failed  %v", err)
		return oh, err
	}
	oh.prClient = prClient
	if os.Getenv(global_environment.ENV_KEY_SECRET_STORAGE_LOCATION) == comn.VAULT_SECRET {
		//Vault location
		oh.VaultSession, err = vauth.GetSession()
		if err != nil {
			LOGGER.Errorf("Error reading account source environment settings")
			return oh, err
		}
	}

	gasServiceClient := gasserviceclient.Client{
		HTTP: &http.Client{Timeout: time.Second * 20},
		URL:  os.Getenv(global_environment.ENV_KEY_GAS_SVC_URL),
	}
	oh.GasServiceClient = &gasServiceClient
	return oh, nil
}

func (oh OnBoardingHandler) RegisterAnchor(w http.ResponseWriter, request *http.Request) {
	LOGGER.Infof("anchor-service:Asset Operations :RegisterAnchor")

	urlVars := mux.Vars(request)
	anchorDomain := urlVars["anchor_domain"]
	if anchorDomain == "" {
		response.NotifyWWError(w, request, http.StatusBadRequest, "ANCHOR-0028", errors.New("anchor domain should not be empty"))
		return
	}
	regRequest := model.RegisterAnchorRequest{}
	err := json.NewDecoder(request.Body).Decode(&regRequest)
	if err != nil {
		LOGGER.Warningf("Error while validating Participant Status :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "ANCHOR-0022", err)
		return
	}
	err = regRequest.Validate(strfmt.Default)

	if err != nil {
		LOGGER.Warningf("Error while validating Participant Status :  %v", err)
		response.NotifyWWError(w, request, http.StatusBadRequest, "ANCHOR-0022", err)
		return
	}

	//This is admin endpoint JWT check is not needed on these endpoints
	//Check JWT token
	/*if os.Getenv(global_environment.ENV_KEY_ENABLE_JWT) != "false" {
		participantID, err := middlewares.GetIdentity(req)
		//Check if requesting anchor id is same as participant id in the token
		if participantID != anchorDomain {
			response.NotifyWWError(w, request, http.StatusUnauthorized, "ANCHOR-0067",
				err)
			return
		}
	}*/

	account := ast.GetStellarAccount(*regRequest.Address)
	if account.AccountID == "" {
		LOGGER.Warningf("error validating account on stellar nw :  %v", *regRequest.Address)
		response.NotifyWWError(w, request, http.StatusNotFound, "ANCHOR-0023", errors.New(*regRequest.Address))
		return
	}

	hasIBMSigner := oh.hasIBMSigner(account.Signers)
	if hasIBMSigner == false {
		LOGGER.Warningf("error validating account on stellar nw :  %v", *regRequest.Address)
		response.NotifyWWError(w, request, http.StatusNotFound, "ANCHOR-0024", errors.New(*regRequest.Address))
		return
	}

	anchorAccount := model.Account{Address: regRequest.Address, Name: "issuing"}
	err = oh.prClient.PostParticipantIssuingAccount(anchorDomain, anchorAccount)
	if err != nil {
		LOGGER.Errorf("Error adding issuing account to PR, Account address: %v, %v", anchorAccount.Address, err)
		response.NotifyWWError(w, request, http.StatusNotFound, "ANCHOR-0025", err)
		return
	}

	//return back generated pass code
	acrJSON, _ := json.Marshal(anchorAccount)
	response.Respond(w, http.StatusOK, acrJSON)
}

func (oh OnBoardingHandler) OnBoardAsset(w http.ResponseWriter, req *http.Request) {

	LOGGER.Debugf("anchor-service:Asset Operations :OnBoardAsset")
	queryParams := req.URL.Query()

	urlVars := mux.Vars(req)
	anchorDomain := urlVars["anchor_domain"]
	if anchorDomain == "" {
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0028", errors.New("anchor domain should not be empty"))
		return
	}

	assetCode := queryParams["asset_code"][0]
	assetType := queryParams["asset_type"][0]

	if assetCode == "" || assetType == "" {
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0028", errors.New("asset code and type should not be empty"))
		return
	}
	if assetType != model.AssetAssetTypeDA {
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0065", errors.New(assetType))
		return
	}
	//This is admin endpoint JWT check is not needed on these endpoints
	//Check JWT token
	/*if os.Getenv(global_environment.ENV_KEY_ENABLE_JWT) != "false" {
		participantID, err := middlewares.GetIdentity(req)
		//Check if requesting anchor id is same as participant id in the token
		if participantID != anchorDomain {
			response.NotifyWWError(w, req, http.StatusUnauthorized, "ANCHOR-0067",
				err)
			return
		}
	}*/

	issuer, err := oh.prClient.GetParticipantIssuingAccount(anchorDomain)
	if err != nil {
		LOGGER.Errorf("Error Getting issuing account for anchor domain: ", err.Error())
		response.NotifyWWError(w, req, http.StatusNotFound, "ANCHOR-0028", err)
		return
	}

	asst := model.Asset{}
	asst.AssetType = &assetType
	asst.AssetCode = &assetCode
	asst.IssuerID = anchorDomain
	err = asst.Validate(strfmt.Default)

	if err != nil {
		LOGGER.Errorf("Error validating issue asset request: ", err.Error())
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0028", err)
		return
	}

	if asset.GetAssetType(assetCode) != assetType {
		if assetType == model.AssetAssetTypeDO {
			LOGGER.Errorf("Error: asset_code should not end with DO")
			response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0028", errors.New(
				"for digial obligation, asset_code should end with DO"))
		} else if assetType == model.AssetAssetTypeNative {
			LOGGER.Errorf("Error validating issue asset request: Native asset cannot be issued ")
			response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0028", errors.New(
				" Native asset cannot be issued"))
		} else {
			response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0028", errors.New(
				"for digital asset, asset_code should not end with DO"))
		}
		return
	}

	limit := "1"
	LOGGER.Debugf("Issuer: %s", issuer)
	LOGGER.Debugf("Asset Code received : %s", assetCode)

	err = model.IsValidDACode(*asst.AssetCode)
	if err != nil {
		LOGGER.Debug("OnBoardAsset:", err)
		response.NotifyWWError(w, req, http.StatusBadRequest, "ANCHOR-0028", err)
		return
	}

	//Get IBM token account from nc, vault or AWS secret mngr
	ibmAccount, err := participant.GenericGetIBMTokenAccount(oh.VaultSession)
	if err != nil {
		msg := "Error getting IBM account"
		code := "ANCHOR-0077"
		response.NotifyWWError(w, req, http.StatusConflict, code, errors.New(msg))
		return
	}

	xdrStr, err := asset.ChangeTrust(oh.GasServiceClient,
		ibmAccount.NodeAddress,
		issuer,
		assetCode,
		limit)

	raw := strings.NewReader(xdrStr)
	b64r := base64.NewDecoder(base64.StdEncoding, raw)

	var tx xdr.TransactionEnvelope
	bytesRead, err := xdr.Unmarshal(b64r, &tx)
	LOGGER.Debugf("read %d bytes from Xdr \n", bytesRead)
	LOGGER.Debugf("This tx has %d operations, source Account %s\n", len(tx.Tx.Operations), tx.Tx.SourceAccount)
	LOGGER.Debugf("This tx has paid %d fee\n", tx.Tx.Fee)
	LOGGER.Debugf("This tx has sequence number %d \n", tx.Tx.SeqNum)
	stellarPassphrase := os.Getenv(global_environment.ENV_KEY_STELLAR_NETWORK)

	LOGGER.Infof("Stellar Network: %v", stellarPassphrase)

	txeb := &b.TransactionEnvelopeBuilder{E: &tx}
	txeb.Init()
	stellarNetwork := comn.GetStellarNetwork(stellarPassphrase)
	err = txeb.MutateTX(stellarNetwork)

	//Sign with IBM account
	sig := b.Sign{Seed: ibmAccount.NodeSeed}
	err = sig.MutateTransactionEnvelope(txeb)
	txeB64, err := txeb.Base64()

	LOGGER.Debugf("signed transaction: %v", txeB64)
	if err != nil {
		LOGGER.Error("The Asset could not be issued. Error signing transaction.")
		response.NotifyWWError(w, req, http.StatusInternalServerError, "ANCHOR-0029", err)
		return
	}
	//submit transaction to stellar
	//step 3: create new operating account  in stellar

	//Submit on Gas service
	hash, ledger, err := oh.GasServiceClient.SubmitTxe(txeB64)
	if err != nil {
		err = ast.DecodeStellarError(err)
		LOGGER.Error("The Asset could not be issued. Error Communicating with Stellar.")
		response.NotifyWWError(w, req, http.StatusInternalServerError, "ANCHOR-0029", err)
		return
	}
	LOGGER.Debugf("submitTransaction  %v, %v", hash, ledger)

	astType := *asst.AssetType
	assetBytes, _ := json.Marshal(oh.successResponse(assetCode, anchorDomain, astType))
	response.Respond(w, http.StatusOK, assetBytes)
}

func (oh OnBoardingHandler) successResponse(code, issuer, AssetType string) *model.Asset {
	asset := model.Asset{
		AssetCode: &code,
		IssuerID:  issuer,
		AssetType: &AssetType,
	}
	return &asset
}

//Validate if account has valid IBM signer
func (oh OnBoardingHandler) hasIBMSigner(signers []horizon.Signer) bool {
	//Get IBM token account from nc, vault or AWS secret mngr
	ibmAccount, err := participant.GenericGetIBMTokenAccount(oh.VaultSession)
	if err != nil {
		LOGGER.Debugf("hasIBMSigner: Error getting IBM account")
		return false
	}
	for _, n := range signers {
		//New horizon SDK used key instead of PublicKey
		if ibmAccount.NodeAddress == n.Key {
			return true
		}
	}
	return false
}
