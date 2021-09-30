// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package utility

import (
	"encoding/base64"
	"errors"
	"net/http"
	"os"

	"github.com/op/go-logging"
	b "github.com/stellar/go/build"
	crypto_client "github.com/GFTN/gftn-services/crypto-service-client/crypto-client"
	gasserviceclient "github.com/GFTN/gftn-services/gas-service-client"
	ast "github.com/GFTN/gftn-services/utility/asset"
	util "github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
)

var LOGGER = logging.MustGetLogger("utility")

func CreatePaymentTransaction(gClient gasserviceclient.Client, sender string, receiver string, assetCode string, assetIssuer string, amount float64) (*b.TransactionBuilder, error) {

	horizonClient := util.GetHorizonClient(os.Getenv(global_environment.ENV_KEY_HORIZON_CLIENT_URL))
	stellarNetwork := util.GetStellarNetwork(os.Getenv(global_environment.ENV_KEY_STELLAR_NETWORK))

	LOGGER.Debugf("BuildPaymentTransaction hc=%v, horizonClient=%v", os.Getenv(global_environment.ENV_KEY_HORIZON_CLIENT_URL), horizonClient)

	//Get IBM gas account
	ibmAccount, sequenceNum, err := gClient.GetAccountAndSequence()

	tx, err := b.Transaction(
		b.SourceAccount{AddressOrSeed: ibmAccount},
		stellarNetwork,
		b.Sequence{Sequence: sequenceNum},
		b.Payment(
			b.SourceAccount{AddressOrSeed: sender},
			b.Destination{AddressOrSeed: receiver},
			b.CreditAmount{Code: assetCode, Issuer: assetIssuer, Amount: util.FloatToString(amount)},
		),
	)

	if err != nil {
		LOGGER.Warningf(" Error building transaction %s", err)
		return tx, err
	}
	return tx, nil
}

func SubmitPaymentTransaction(gClient gasserviceclient.Client, sourceAccount string, tx *b.TransactionBuilder,
	cClient crypto_client.CryptoServiceClient) (hash string, ledger uint64, err error) {
	LOGGER.Debugf("SubmitPaymentTransaction start... sourceSeed = %v %v", sourceAccount, tx)

	var txeb b.TransactionEnvelopeBuilder
	err = txeb.Mutate(tx)
	txeB64, err := txeb.Base64()
	//TBD: will have to integrate with gas service
	xdrB, _ := base64.StdEncoding.DecodeString(txeB64)

	sigXdr_, errorMsg, status, _ := cClient.ParticipantSignXdr(sourceAccount, xdrB)

	if status != http.StatusCreated {
		LOGGER.Errorf("Error creating transaction: %v", errorMsg.Error())
		return "", 0, errors.New("Error creating transaction")
	}
	if errorMsg != nil {
		LOGGER.Errorf("Error creating new account: %v", errorMsg.Error())
		return "", 0, errorMsg
	}
	LOGGER.Debugf("signed transaction: %v", base64.StdEncoding.EncodeToString(sigXdr_))

	//Post to gas service
	hash, ledger, err = gClient.SubmitTxe(base64.StdEncoding.EncodeToString(sigXdr_))
	if err != nil {
		LOGGER.Warningf("SubmitPaymentTransaction failed... %v ", err.Error())
		newErr := ast.DecodeStellarPaymentError(err)
		LOGGER.Warningf(newErr.Error())
		return hash, ledger, err
	}

	LOGGER.Debugf("SubmitPaymentTransaction OK")
	return hash, ledger, err
}
