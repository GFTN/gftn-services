// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package signing

import (
	"encoding/base64"
	"errors"
	"os"

	"github.com/beevik/etree"
	"github.com/op/go-logging"
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
	crypto_client "github.com/GFTN/gftn-services/crypto-service-client/crypto-client"
	"github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/participant"
	"github.com/GFTN/gftn-services/utility/payment/client"
	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"
	"github.com/GFTN/gftn-services/utility/vault/utils"
	"github.com/GFTN/gftn-services/utility/xmldsig"
)

var LOGGER = logging.MustGetLogger("sign-handler")

type CreateSignOperations struct {
	prServiceURL             string
	cryptoServiceURL         string
	cryptoServiceInternalURL string
	CryptoServiceClient      crypto_client.RestCryptoServiceClient
}

func InitiateSignOperations(pr string) (op CreateSignOperations) {
	err := errors.New("")
	//Construct participant specific crypto service URLs for initialization
	op.cryptoServiceInternalURL, err = participant.GetServiceUrl(os.Getenv(global_environment.ENV_KEY_CRYPTO_SVC_INTERNAL_URL), os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME))
	if err != nil {
		LOGGER.Error("Error initializing CRYPTO_SERVICE_INTERNAL_URL for participant")
		return op
	}
	csClient, csClientErr := crypto_client.CreateRestCryptoServiceClient(op.cryptoServiceInternalURL)
	if csClientErr != nil {
		LOGGER.Error("Can not connect to crypto service client, please check if crypto service is running")
	}
	op.CryptoServiceClient = csClient
	op.prServiceURL = pr

	return op
}

func (op *CreateSignOperations) SignPayload(data []byte, accountName string) ([]byte, error) {
	signedPayload, signErr, statusCode, errCode := op.CryptoServiceClient.SignPayload(accountName, data)
	if signErr != nil {
		LOGGER.Errorf("Sign payload Error: %s, %d", errCode, statusCode)
		return nil, signErr
	}

	return signedPayload, nil
}

func (op *CreateSignOperations) SignPayloadByMasterAccount(data []byte) ([]byte, error) {
	LOGGER.Infof("Signing message payload with IBM master account")
	signedPayload, signErr := op.CryptoServiceClient.SignPayloadByMasterAccount(data)
	if signErr != nil {
		LOGGER.Errorf("Sign payload by Master account Error: %v", signErr)
		return nil, signErr
	}
	LOGGER.Infof("Message payload signing success!")

	return signedPayload, nil
}

func (op *CreateSignOperations) SignTx(ibmAccount, sendingAccount, receivingAccount, settlementAccountName string, seqNum uint64, dbData *sendmodel.SignData, memoHash xdr.Memo) (string, error) {
	var tx *build.TransactionBuilder
	var err error

	// two different scenario for XLM and DO DA
	if dbData.CurrencyCode == "XLM" {
		// default fee amount is zero
		LOGGER.Info("0 fee")
		tx, err = build.Transaction(
			build.SourceAccount{AddressOrSeed: ibmAccount},
			build.Sequence{Sequence: seqNum},
			build.Payment(
				build.SourceAccount{AddressOrSeed: sendingAccount},
				build.Destination{AddressOrSeed: receivingAccount},
				build.NativeAmount{
					Amount: dbData.SettlementAmount,
				},
			),
		)
		if err != nil {
			LOGGER.Errorf("Error while creating the transaction: %v\n", err.Error())
			return "", err
		}
	} else {
		// query the asset issuer's account address using asset issuer's id from participant-registry
		issuerAccount := client.GetParticipantAccount(op.prServiceURL, dbData.AssetIssuerId, common.ISSUING)
		if issuerAccount == nil {
			LOGGER.Errorf("Can't find asset issuer's account in participant-registry: %s", dbData.AssetIssuerId)
			return "", errors.New("can not find asset issuer's account in participant-registry")
		}

		// default fee amount is zero
		LOGGER.Info("0 fee")
		tx, err = build.Transaction(
			build.SourceAccount{AddressOrSeed: ibmAccount},
			//build.AutoSequence{SequenceProvider: &horizonClient},
			build.Sequence{Sequence: seqNum},
			build.Payment(
				build.SourceAccount{AddressOrSeed: sendingAccount},
				build.Destination{AddressOrSeed: receivingAccount},
				build.CreditAmount{
					Code:   dbData.CurrencyCode,
					Issuer: *issuerAccount,
					Amount: dbData.SettlementAmount,
				},
			),
		)

		if err != nil {
			LOGGER.Errorf("Error while creating the transaction: %v\n", err.Error())
			return "", err
		}
	}

	tx.TX.Memo = memoHash
	sEnc, signErr := op.Signing(tx, settlementAccountName)
	if signErr != nil {
		LOGGER.Errorf("Sign transaction Error: %s", signErr.Error())
		return "", signErr
	}

	return sEnc, nil
}

func (op *CreateSignOperations) Signing(tx *build.TransactionBuilder, settlementAccountName string) (string, error) {
	var txe xdr.TransactionEnvelope
	txe.Tx = *tx.TX
	txeB, _ := txe.MarshalBinary()

	txSigned, signTxErr, statusCode, errCode := op.CryptoServiceClient.ParticipantSignXdr(settlementAccountName, txeB)
	if signTxErr != nil {
		LOGGER.Errorf("Sign transaction Error: %s, %d", errCode, statusCode)
		return "", signTxErr
	}

	sEnc := base64.StdEncoding.EncodeToString(txSigned)

	return sEnc, nil
}

func (op *CreateSignOperations) VerifyPayload(data, signature []byte, publicKey string) bool {
	LOGGER.Infof("Verifying request payload ...")

	kp, perr := keypair.Parse(publicKey)
	if perr != nil {
		LOGGER.Error("Wrong Public Key for verification")
		return false
	}

	verErr := kp.Verify(data, signature)
	if verErr != nil {
		LOGGER.Error("Wrong signature for verification: " + verErr.Error())
		return false
	}

	LOGGER.Infof("Signature is correct")

	return true
}

func SignPayloadByMasterAccount(xml string) (string, error) {
	account, err := participant.GenericGetIBMTokenAccount(utils.Session{})
	if err != nil {
		LOGGER.Errorf("Error occured retrieving the account: %+v", err)
		return "", err
	}
	publicKey := account.NodeAddress
	seed := account.NodeSeed
	// Read the XML string and convert it to DOM object in memory
	payloadDocument := etree.NewDocument()
	err = payloadDocument.ReadFromString(xml)
	if err != nil {
		LOGGER.Errorf("%v", err)
		return "", err
	}

	appHeaderElement := xmldsig.GetElementByName(payloadDocument.Root(), "AppHdr")
	if appHeaderElement == nil {
		LOGGER.Infof("AppHdr is missing")
		return "", errors.New("AppHdr is missing")

	}

	sgntrElement := xmldsig.GetElementByName(appHeaderElement, "Sgntr")
	if sgntrElement != nil {
		LOGGER.Errorf("Sgntr tag is already available")
		return "", errors.New("Sgntr tag is already available")
	}

	//Initialise the Signature object with properties, values are hard coded only support SHA-256
	sgntr := xmldsig.NewSignature()

	//Canonicalise the original xml payload
	canonicalisedPayload, err := xmldsig.Canonicalise(xml)
	if err != nil {
		LOGGER.Errorf("%v", err)
		return canonicalisedPayload, err
	}

	//Get the digest of the Canonicalised payload
	canonicalisedPayloadHashed, err := xmldsig.GenerateSHA256Hash(canonicalisedPayload)
	if err != nil {
		LOGGER.Errorf("%v", err)
		return "", err
	}

	//Based64 encode the canonicalised hashed payload
	canonicalisedPayloadHashedEncoded := base64.StdEncoding.EncodeToString(canonicalisedPayloadHashed)
	//set the hash to the Signature struct
	sgntr.Signature.SignedInfo.Reference.DigestValue = canonicalisedPayloadHashedEncoded

	//Now generate the SignedInfo tag by marshalling the struct to XML string
	signedInfoString, err := xmldsig.MarshallToXML(sgntr.Signature.SignedInfo)
	if err != nil {
		LOGGER.Errorf("%v", err)
		return signedInfoString, err
	}

	//Canonicalise the SignedInfo
	signedInfoStringCanonicalised, err := xmldsig.Canonicalise(signedInfoString)

	checkSumToSign, err := xmldsig.GenerateSHA256Hash(signedInfoStringCanonicalised)
	if err != nil {
		LOGGER.Errorf("%v", err)
		return "", err
	}

	kp, err := keypair.Parse(seed)
	if err != nil {
		LOGGER.Errorf("There was an error while getting node seed")
		return "", err
	}
	signedInfoSignature, err := kp.Sign(checkSumToSign)
	if err != nil {
		LOGGER.Errorf("%v", err)
		return "", err
	}

	//Populate the Sigature object with hashed signaure value
	sgntr.Signature.SignatureValue = base64.StdEncoding.EncodeToString(signedInfoSignature)

	//Populate the public key
	x509Data := &xmldsig.X509Data{
		X509Certificate: publicKey,
	}
	sgntr.Signature.KeyInfo.X509Data = x509Data

	//Generate the XML from Signature object
	signature, err := xmldsig.MarshallToXML(sgntr)
	if err != nil {
		LOGGER.Errorf("%v", err)
		return signature, err
	}

	//Create new XML node Signature
	sgnrDocument := etree.NewDocument()
	sgnrDocument.ReadFromString(signature)

	//Add the signatureDocument under AppHdr tag
	for _, element := range payloadDocument.Root().ChildElements() {
		if element.Tag == "AppHdr" {
			element.AddChild(sgnrDocument.Root())
		}
	}

	//Return the Signed XML document as String
	s, err := payloadDocument.WriteToString()
	if err != nil {
		LOGGER.Errorf("%v", err)
		return s, err
	}

	xmldsig.VerifySignature(s)
	return s, nil
}
