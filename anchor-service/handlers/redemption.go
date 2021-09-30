// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-resty/resty"
	"github.com/GFTN/gftn-services/anchor-service/environment"
	crypto_client "github.com/GFTN/gftn-services/crypto-service-client/crypto-client"
	gasserviceclient "github.com/GFTN/gftn-services/gas-service-client"
	"github.com/GFTN/gftn-services/gftn-models/model"
	pr_client "github.com/GFTN/gftn-services/participant-registry-client/pr-client"
	util "github.com/GFTN/gftn-services/utility"
	"github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/participant"
	vauth "github.com/GFTN/gftn-services/utility/vault/auth"
	"github.com/GFTN/gftn-services/utility/vault/utils"
)

type Operations struct {
	ParticipantRegistryClient pr_client.PRServiceClient
	CryptoClient              crypto_client.CryptoServiceClient
	VaultSession              utils.Session
	GasServiceClient          gasserviceclient.GasServiceClient
	StrongHoldAnchorID        string
}

func CreateAnchorOperations() (Operations, error) {

	op := Operations{}
	shAnchorDomain := os.Getenv(global_environment.ENV_KEY_STRONGHOLD_ANCHOR_ID)

	op.VaultSession = utils.Session{}

	if os.Getenv(global_environment.ENV_KEY_SECRET_STORAGE_LOCATION) == common.VAULT_SECRET {
		//Vault location
		err := errors.New("")
		op.VaultSession, err = vauth.GetSession()
		if err != nil {
			LOGGER.Errorf("Error reading account source environment settings")
			return op, err
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

	if shAnchorDomain == "" {
		msg := "ENV variable STRONGHOLD_ANCHOR_ID should be set to strong hold domain name and api url should be set correctly on PR"
		LOGGER.Errorf(msg)
		util.ExitOnErr(LOGGER, errors.New(msg), msg)
	}

	op.StrongHoldAnchorID = shAnchorDomain

	op.ParticipantRegistryClient = prClient

	var cClient crypto_client.CryptoServiceClient
	err := errors.New("")
	cServiceInternalUrl, err := participant.GetServiceUrl(os.Getenv(global_environment.ENV_KEY_CRYPTO_SVC_INTERNAL_URL), os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME))
	if err != nil {
		return op, err
	}
	cClient, err = crypto_client.CreateRestCryptoServiceClient(cServiceInternalUrl)
	if err != nil {
		return op, err
	}

	gasServiceClient := gasserviceclient.Client{
		HTTP: &http.Client{Timeout: time.Second * 20},
		URL:  os.Getenv(global_environment.ENV_KEY_GAS_SVC_URL),
	}
	op.GasServiceClient = &gasServiceClient
	op.CryptoClient = cClient
	return op, nil

}

func (op Operations) WithDraw(withdrawRequest model.StrongholdWithdrawRequest) (*model.StrongholdWithdrawResponse, error) {

	LOGGER.Debugf("WithDraw")
	err := withdrawRequest.Validate(strfmt.Default)
	if err != nil {
		msg := "Unable to validate Anchor request: " + err.Error()
		LOGGER.Warningf(msg)
		return nil, err
	}

	var ia = ""
	var anchorResponse = model.StrongholdWithdrawResponse{}

	anchorId := os.Getenv(global_environment.ENV_KEY_STRONGHOLD_ANCHOR_ID)
	//get anchor issuing account
	ia, err = op.ParticipantRegistryClient.GetParticipantIssuingAccount(anchorId)
	if err != nil || ia == "" {
		return nil, err
	}

	anchorResponse, err = op.getSHTransactionReference(withdrawRequest)
	if err != nil {
		LOGGER.Debugf("error getting getSHTransactionReference %v", err)
		return nil, err
	}

	return &anchorResponse, nil
}

func computeSHHmac256(body string) (hmacStr string, timeStr string) {
	LOGGER.Debugf("computeSHHmac256: %v", body)
	secret := os.Getenv(environment.ENV_KEY_ANCHOR_SH_SEC)
	key, _ := base64.StdEncoding.DecodeString(secret)
	h := hmac.New(sha256.New, key)
	strTime := strconv.FormatInt(time.Now().Unix(), 10)
	LOGGER.Debugf("\nUnixTime: %v", strTime)
	venue := os.Getenv(environment.ENV_KEY_ANCHOR_SH_VENEU)
	message := strTime + "POST" + "/v1/venues/" + venue + "/direct/withdrawals" + body
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), strTime
}

func (op Operations) getSHTransactionReference(withdrawRequest model.StrongholdWithdrawRequest) (model.StrongholdWithdrawResponse, error) {
	// shAnchorDomain := *withdrawRequest.AnchorID
	//Anchor callback should resolve from callback URL
	// participant, err := op.ParticipantRegistryClient.GetParticipantForDomain(shAnchorDomain)
	// if err != nil {
	// 	LOGGER.Errorf("Error finding given anchor domain on PR: %v ", shAnchorDomain)
	// 	return anchorResponse, "", nil, errors.New("error finding given anchor domain on PR")
	// }

	// shAnchorServiceWithdrawURL := *participant.URLCallback + "/" + os.Getenv(environment.ENV_KEY_ANCHOR_SH_VENEU) + "/direct/withdrawals"

	// TODO: define anchor callback
	url := os.Getenv(environment.ENV_KEY_ANCHOR_SH_ROOT_URL)
	if url == "" {
		err := errors.New("ANCHOR_SH_ROOT_URL not set")
		return model.StrongholdWithdrawResponse{}, err
	}

	bodyBytes, err := json.Marshal(&withdrawRequest)
	credID := os.Getenv(environment.ENV_KEY_ANCHOR_SH_CRED)
	passPhrase := os.Getenv(environment.ENV_KEY_ANCHOR_SH_PASS)
	venueId := os.Getenv(environment.ENV_KEY_ANCHOR_SH_VENEU)
	if credID == "" || passPhrase == "" || venueId == "" {
		err = errors.New("Stronnghold credetials are not set for this participant, please set ANCHOR_SH_CRED, ANCHOR_SH_PASS, ANCHOR_SH_SEC and ANCHOR_SH_VENUE")
		LOGGER.Debug("WithDraw:", err)
		return model.StrongholdWithdrawResponse{}, err
	}

	url = url + "/v1/venues/" + venueId + "/direct/withdrawals"
	LOGGER.Debugf("Anchor service WithdrawDigitalAsset url:  %v", url)

	var responseBody []byte
	hmacString, timeStr := computeSHHmac256(string(bodyBytes))

	aResponse, err := resty.R().SetHeader(
		"Content-type", "application/json").SetHeader(
		"SH-CRED-ID", credID).SetHeader(
		"SH-CRED-SIG", hmacString).SetHeader(
		"SH-CRED-TIME", timeStr).SetHeader(
		"SH-CRED-PASS", passPhrase).SetBody(bodyBytes).Post(url)

	if err != nil {
		LOGGER.Errorf("Error while making request to Anchor:  %v", err)
		return model.StrongholdWithdrawResponse{}, err
	}

	responseBody = aResponse.Body()
	var SHAnchorResponse = model.StrongholdWithdrawResponse{}
	err = json.Unmarshal(responseBody, &SHAnchorResponse)

	if err != nil {
		LOGGER.Debug("StrongHold error response: %v", string(responseBody))
		return SHAnchorResponse, errors.New(err.Error() + "")
	}
	if SHAnchorResponse.Result == nil {
		LOGGER.Debug("StrongHold response: %v", string(responseBody))
		return SHAnchorResponse, errors.New("received an empty response from anchor" + "")
	}

	//set transaction id to stellar hash received

	LOGGER.Debugf("anchor withdraw pay to return address: %s, reference: %s ", SHAnchorResponse.Result.PaymentMethodInstructions.PayToVenueSpecific, SHAnchorResponse.Result.PaymentMethodInstructions.PayToReference)
	return SHAnchorResponse, nil

}
