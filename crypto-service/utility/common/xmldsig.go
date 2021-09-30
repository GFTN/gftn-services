// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package common

import (
	"crypto"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"errors"

	"github.com/beevik/etree"
	"github.com/stellar/go/keypair"
	"github.com/GFTN/gftn-services/utility/xmldsig"
)

func (obj *HsmObject) generateSignatureTag(xml, publicKey string, privateKeyHandle string, signUsingStellar bool) (string, error) {
	//Initialise the Signature object with properties, values are hard coded only support SHA-256
	sgntr := xmldsig.NewSignature()

	//Canonicalise the original xml payload
	canonicalisedPayload, err := xmldsig.Canonicalise(xml)
	if err != nil {
		return "", err
	}

	//Get the digest of the Canonicalised payload
	canonicalisedPayloadHashed, err := xmldsig.GenerateSHA256Hash(canonicalisedPayload)
	if err != nil {
		return "", err
	}

	//Based64 encode the canonicalised hashed payload
	canonicalisedPayloadHashedEncoded := base64.StdEncoding.EncodeToString(canonicalisedPayloadHashed)
	//set the hash to the Signature struct
	sgntr.Signature.SignedInfo.Reference.DigestValue = canonicalisedPayloadHashedEncoded

	//Now generate the SignedInfo tag by marshalling the struct to XML string
	signedInfoString, err := xmldsig.MarshallToXML(sgntr.Signature.SignedInfo)
	if err != nil {
		return "", err
	}

	//Canonicalise the SignedInfo
	signedInfoStringCanonicalised, err := xmldsig.Canonicalise(signedInfoString)

	signedInfoSignature, err := obj.signSignedInfo(signedInfoStringCanonicalised, privateKeyHandle, signUsingStellar)
	if err != nil {
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
	s, err := xmldsig.MarshallToXML(sgntr)
	if err != nil {
		return "", err
	}

	return s, nil
}

func (obj *HsmObject) SignXML(xml, privateKeyHandle string, publicKey string, signUsingStellar bool) (string, error) {
	// Read the XML string and convert it to DOM object in memory
	payloadDocument := etree.NewDocument()
	err := payloadDocument.ReadFromString(xml)
	if err != nil {
		return "", errors.New("Error while parsing XML into DOM")
	}

	appHeaderElement := xmldsig.GetElementByName(payloadDocument.Root(), "AppHdr")
	if appHeaderElement == nil {
		LOGGER.Infof("AppHdr is missing")
		return "", errors.New("Payload doesn't have AppHdr")
	}

	sgntrElement := xmldsig.GetElementByName(appHeaderElement, "Sgntr")
	if sgntrElement != nil {
		LOGGER.Infof("Sgntr tag is already available")
		return "", errors.New("Sgntr tag is already available")
	}

	signature, err := obj.generateSignatureTag(xml, publicKey, privateKeyHandle, signUsingStellar)
	if err != nil {
		return "", err
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
		return "", err
	}
	return s, nil
}

/*
Messaged will be hashed and signed
*/
func (obj *HsmObject) signSignedInfo(data, privateKayHandle string, signUsingStellar bool) ([]byte, error) {
	//Signing has to sign the hash of the SignedInfo
	checkSumToSign, err := xmldsig.GenerateSHA256Hash(data)
	if err != nil {
		return nil, err
	}
	if signUsingStellar {
		return doSignUsingStellar(checkSumToSign, privateKayHandle)
	} else {
		return obj.doSignUsingHSM(checkSumToSign, privateKayHandle)
	}
	//return doSignUsingPEM(checkSumToSign, publicKey)
}

func (obj *HsmObject) doSignUsingHSM(data []byte, privateKayHandle string) ([]byte, error) {
	LOGGER.Infof("Signing using privateKayHandle: %s", privateKayHandle)

	signedPayload, err := obj.SignUsingHSMPrivateKeyHandle(data, privateKayHandle)
	if err != nil {
		return nil, err
	}

	return signedPayload, nil
}

func doSignUsingStellar(data []byte, privateKayHandle string) ([]byte, error) {
	LOGGER.Infof("Signing using privateKayHandle: %s", privateKayHandle)
	kp, err := keypair.Parse(privateKayHandle)
	if err != nil {
		LOGGER.Errorf("There was an error while getting node seed")
		return nil, err
	}
	return kp.Sign(data)
}

func doSignUsingPEM(data []byte, accountName string) ([]byte, error) {
	//This is using the Self signed pem files, has to use HSM here
	cert, _ := tls.LoadX509KeyPair("certificate.pem", "key.pem")
	signer := cert.PrivateKey.(crypto.Signer)
	signature, err := signer.Sign(rand.Reader, data, crypto.SHA256)
	if err != nil {
		LOGGER.Errorf("Error occured while signing using PEM files. %s", err)
		return nil, err
	}

	return signature, nil
}
