// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package hsmclient

import (
	"errors"

	"github.com/miekg/pkcs11"
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
	"github.com/GFTN/gftn-services/crypto-service/utility/hsm"
)

// Signing the transaction hash with the handle id of private key. Need to initiate the HSM and get the private
// jey object handle id first.
// Add both the signature and the hint parse from public key to the transaction.
func GetSignatureAndAddToTransaction(c hsm.Crypto, privateKeyObject pkcs11.ObjectHandle, slot uint, txHash []byte, tx build.TransactionEnvelopeBuilder, publicKey string) (txe build.TransactionEnvelopeBuilder, err error) {
	LOGGER.Debugf("c: %s, privateKeyObject: %s, slot: %s, txHash: %s", c, privateKeyObject, slot, txHash)
	signature, err := SignDataWithPrivateKey(c, privateKeyObject, slot, txHash)
	if err != nil {
		LOGGER.Errorf("Error while signing a data with private key: %v\n", err)
		return txe, err
	}
	LOGGER.Infof("Signing transaction hash with private key. %v\n", signature)

	txe, err = AddSignatureToTransaction(tx, signature, publicKey)

	if err != nil {
		LOGGER.Errorf("Error while creating a TransactionEnvelopeBuilder: %v\n", err)
		return txe, err
	}
	LOGGER.Infof("Add signature to transaction.\n")

	return txe, nil
}

func SignDataWithPrivateKey(c hsm.Crypto, privateKeyObject pkcs11.ObjectHandle, slot uint, data []byte) ([]byte, error) {
	session, err := c.OpenSession(slot)
	if err != nil {
		LOGGER.Errorf("Error while opening a session: %v\n", err)
	}
	//LOGGER.Debugf("SignDataWithPrivateKey: %v, %v, %v, %v", c, slot, data, privateKeyObject)
	signature, err := c.Sign(session, privateKeyObject, data[:])

	defer c.CloseSession(session)
	return signature, err
}

func VerifyDataWithPublicKey(c hsm.Crypto, publicKeyObject pkcs11.ObjectHandle, slot uint, data []byte, sig []byte) error {
	session, err := c.OpenSession(slot)
	if err != nil {
		LOGGER.Errorf("Error while opening a session: %v\n", err)
	}

	//LOGGER.Debugf("data: %v, sig: %v", data, sig)

	err = c.Verify(session, publicKeyObject, data[:], sig)
	defer c.CloseSession(session)
	return err
}

func VerifyDataWithPublicNodeAddress(publicKey string, data []byte, sig []byte) error {
	LOGGER.Debugf("verifying with node address: %v", publicKey)
	kp, err := keypair.Parse(publicKey)
	if err != nil {
		return err
	}

	err = kp.Verify(data[:], sig)

	return err
}

func AddSignatureToTransaction(tx build.TransactionEnvelopeBuilder, signature []byte, publicKey string) (build.TransactionEnvelopeBuilder, error) {
	kp, err := keypair.Parse(publicKey)

	if err != nil {
		return tx, err
	}

	ds0 := xdr.DecoratedSignature{
		Hint:      xdr.SignatureHint(kp.Hint()),
		Signature: xdr.Signature(signature),
	}
	//LOGGER.Debugf("Sigs: %v, %v", len(tx.E.Signatures), tx.E.Signatures)
	tx.E.Signatures = append(tx.E.Signatures, ds0)

	//LOGGER.Debugf("New Sigs: %v, %v", len(tx.E.Signatures), tx.E.Signatures)

	return tx, nil
}

// Using public key's CKA_EC_POINT from the HSM to generate the Stellar account.
func GenerateStellarAccount(ecPoints []uint8) string {
	publicKey, err := strkey.Encode(strkey.VersionByteAccountID, ecPoints[2:])
	if err != nil {
		LOGGER.Errorf("Error while generating Stellar account using public key's ec point: %v\n", err)
	}
	LOGGER.Infof("Public Key: %s\n", publicKey)

	return publicKey
}

// Generate key pair and get the label of public key and private key from HSM. Use the public key label to get
// the CKA_EC_POINT. Need to initiate the HSM first.
func GenerateKeyPair(c *hsm.Crypto, slot uint, label1, label2 string) (string, string, pkcs11.ObjectHandle, pkcs11.ObjectHandle, []uint8, error) {
	session, err := c.OpenSession(slot)
	defer c.CloseSession(session)
	if err != nil {
		LOGGER.Errorf("Error while opening a session: %v\n", err)
		return "", "", 0, 0, nil, err
	}
	publicKeyLabel, privateKeyLabel, publicKeyObjectHandle, privateKeyObjectHandle, err := c.GenerateED25519KeyPair(session, label1, label2)
	if err != nil {
		LOGGER.Errorf("Error while generating ED25519 key pair: %v\n", err)
		return "", "", 0, 0, nil, err
	}
	LOGGER.Infof("Public Key label='%s', Private Key label='%s'\n", publicKeyLabel, privateKeyLabel)
	ecPoints, _, ecErr := FindHSMObject(c, slot, publicKeyLabel)
	if ecErr != nil {
		LOGGER.Errorf("Error while getting ecPoints from HSM: %v\n", err)
		return "", "", 0, 0, nil, ecErr
	}
	return publicKeyLabel, privateKeyLabel, publicKeyObjectHandle, privateKeyObjectHandle, ecPoints, nil
}

func FindHSMObject(c *hsm.Crypto, slot uint, label string) ([]byte, pkcs11.ObjectHandle, error) {
	session, err := c.OpenSession(slot)
	defer c.CloseSession(session)

	if err != nil {
		LOGGER.Warningf("Error while opening a session: %v\n", err)
	}
	attr, objectHandleId, err := c.FindObject(session, label)
	if uint(objectHandleId) == 0 {
		LOGGER.Errorf("Unable to retrieve HSM object")
		return attr, objectHandleId, errors.New(NULL_HANDLE_ID)
	}
	LOGGER.Infof("Retrieving HSM object ID success!\n")
	return attr, objectHandleId, err
}
