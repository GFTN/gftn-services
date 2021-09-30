// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package hsm

import (
	"fmt"

	"github.com/miekg/pkcs11"
	"github.com/op/go-logging"
)

var LOGGER = logging.MustGetLogger("xdr-sign")

const CKM_EC_EDWARDS_KEY_PAIR_GEN uint = pkcs11.CKM_VENDOR_DEFINED + 0xC01
const CKM_EDDSA uint = pkcs11.CKM_VENDOR_DEFINED + 0xC03

var ed25519ECParams = []byte{0x06, 0x09, 0x2B, 0x06, 0x01, 0x04, 0x01, 0xDA, 0x47, 0x0F, 0x01}

type Crypto struct {
	p *pkcs11.Ctx
}

func NewCrypto(lib string) *Crypto {
	p := pkcs11.New(lib)
	err := p.Initialize()
	if err != nil {
		LOGGER.Error("%s", err)
	}
	return &Crypto{p}
}

func (c *Crypto) FinalizeCrypto() error {
	err := c.p.Finalize()
	return err
}

func (c *Crypto) OpenSession(slot uint) (pkcs11.SessionHandle, error) {
	session, err := c.p.OpenSession(slot, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		LOGGER.Errorf("Error while login to an opened login-session: %v\n", err)
	}
	return session, err
}

func (c *Crypto) CloseSession(session pkcs11.SessionHandle) error {
	err := c.p.CloseSession(session)
	return err
}

func (c *Crypto) Login(session pkcs11.SessionHandle, password string) error {
	return c.p.Login(session, pkcs11.CKU_USER, password)
}

func (c *Crypto) Logout(session pkcs11.SessionHandle) error {
	return c.p.Logout(session)
}

func (c *Crypto) GenerateED25519KeyPair(session pkcs11.SessionHandle, publicKeyLabel, privateKeyLabel string) (string, string, pkcs11.ObjectHandle, pkcs11.ObjectHandle, error) {

	publicKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
		pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_DERIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
		pkcs11.NewAttribute(pkcs11.CKA_MODIFIABLE, false),
		pkcs11.NewAttribute(pkcs11.CKA_EC_PARAMS, ed25519ECParams),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, publicKeyLabel),
	}
	privateKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
		pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true),
		pkcs11.NewAttribute(pkcs11.CKA_DERIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
		pkcs11.NewAttribute(pkcs11.CKA_MODIFIABLE, false),
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, false),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, privateKeyLabel),
	}
	mech := []*pkcs11.Mechanism{pkcs11.NewMechanism(CKM_EC_EDWARDS_KEY_PAIR_GEN, nil)}
	publicKeyObject, privateKeyObject, err := c.p.GenerateKeyPair(session, mech, publicKeyTemplate, privateKeyTemplate)
	return publicKeyLabel, privateKeyLabel, publicKeyObject, privateKeyObject, err
}

func (c *Crypto) Sign(session pkcs11.SessionHandle, privateKey pkcs11.ObjectHandle, data []byte) ([]byte, error) {
	err := c.p.SignInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(CKM_EDDSA, nil)}, privateKey)
	if err != nil {
		return nil, err
	}
	a, b := c.p.Sign(session, data)
	return a, b
}

func (c *Crypto) Verify(session pkcs11.SessionHandle, publicKey pkcs11.ObjectHandle, data []byte, signature []byte) error {
	err := c.p.VerifyInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(CKM_EDDSA, nil)}, publicKey)
	if err != nil {
		return err
	}
	return c.p.Verify(session, data, signature)
}

func (c *Crypto) FindObject(session pkcs11.SessionHandle, label string) ([]uint8, pkcs11.ObjectHandle, error) {
	findTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, label),
	}

	if err := c.p.FindObjectsInit(session, findTemplate); err != nil {
		fmt.Println(err)
		return nil, 0, err
	}

	// get the object with findTemplate which label's name match "label"
	obj, _, err := c.p.FindObjects(session, 1)
	if err != nil || len(obj) == 0 {
		fmt.Println(err)
		return nil, 0, err
	}

	findAttrTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_EC_POINT, nil),
	}

	// get the attributes defined in the findAttrTemplate
	attr, err := c.p.GetAttributeValue(session, obj[0], findAttrTemplate)
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}

	var result []uint8

	for _, a := range attr {
		if a.Type == pkcs11.CKA_EC_POINT {
			result = a.Value
		}
	}

	if err := c.p.FindObjectsFinal(session); err != nil {
		fmt.Println(err)
		return nil, 0, err
	}

	return result, obj[0], nil
}
