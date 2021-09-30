// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package xmldsig

import (
	"bufio"
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"fknsrs.biz/p/xml/c14n"
	"github.com/beevik/etree"
	"github.com/op/go-logging"
	"github.com/stellar/go/keypair"
	"io"
	"io/ioutil"
	"strings"
)

var LOGGER = logging.MustGetLogger("xmldsig")

/*
The C14N library expect a WriteCloser type
*/
type XmlWriteCloser struct {
	io.Writer
}

/*
Implement the Close() method of the WriteCloser Interface
*/
func (XmlWriteCloser) Close() error {
	return nil
}

func NewSignature() *Sgntr {
	sgntr := &Sgntr{}
	sgntr.Signature.SignedInfo.CanonicalizationMethod.Algorithm = "http://www.w3.org/2001/10/xml-exc-c14n#"
	transforms := &sgntr.Signature.SignedInfo.Reference.Transforms.Transform
	*transforms = append(*transforms, Algorithm{"http://www.w3.org/2000/09/xmldsig#enveloped-signature"})
	*transforms = append(*transforms, Algorithm{"http://www.w3.org/2001/10/xml-exc-c14n#"})
	sgntr.Signature.SignedInfo.SignatureMethod.Algorithm = "http://www.w3.org/2009/xmldsig11#rsa-sha256"
	sgntr.Signature.SignedInfo.Reference.DigestMethod.Algorithm = "http://www.w3.org/2001/04/xmlenc#sha256"
	return sgntr
}

func VerifySignature(xml string) bool {
	//Get the document object model
	signedDocument := etree.NewDocument()
	signedDocument.ReadFromString(xml)

	//Get appheader node
	appHeaderElement := GetElementByName(signedDocument.Root(), "AppHdr")
	if appHeaderElement == nil {
		LOGGER.Infof("AppHdr is missing")
		return false
	}

	//Get Sgntr node
	sgntrElement := GetElementByName(appHeaderElement, "Sgntr")
	if sgntrElement == nil {
		LOGGER.Infof("Sgntr is missing")
		return false
	}

	//Get Signature node
	signatureElement := GetElementByName(sgntrElement, "Signature")
	if signatureElement == nil {
		LOGGER.Infof("Signature is missing")
		return false
	}

	//Get signinfo node
	signInfoElement := GetElementByName(signatureElement, "SignedInfo")
	if signInfoElement == nil {
		LOGGER.Infof("SignedInfo is missing")
		return false
	}

	//Get keyinfo node
	keyInfoElement := GetElementByName(signatureElement, "KeyInfo")
	if keyInfoElement == nil {
		LOGGER.Infof("KeyInfo is missing")
		return false
	}

	//
	x509DataElement := GetElementByName(keyInfoElement, "X509Data")
	if x509DataElement == nil {
		LOGGER.Infof("X509Data is missing")
		return false
	}

	//
	x509Certificate := GetElementByName(x509DataElement, "X509Certificate")
	if x509Certificate == nil {
		LOGGER.Infof("X509Certificate is missing")
		return false
	}

	//
	signatureValueElement := GetElementByName(signatureElement, "SignatureValue")
	if signatureValueElement == nil {
		LOGGER.Infof("SignatureValue is missing")
		return false
	}

	//
	referenceElement := GetElementByName(signInfoElement, "Reference")
	if referenceElement == nil {
		LOGGER.Infof("Reference is missing")
		return false
	}

	//
	digentValueElement := GetElementByName(referenceElement, "DigestValue")
	if digentValueElement == nil {
		LOGGER.Infof("DigestValue is missing")
		return false
	}

	//Read the values
	publicKay := x509Certificate.Text()
	digest := digentValueElement.Text()
	signature := signatureValueElement.Text()

	//Conver the SignedInfo node to String to calculate the hash for verification
	signedInfoDocument := etree.NewDocument()
	signedInfoDocument.AddChild(signInfoElement)
	signedInfoString, err := signedInfoDocument.WriteToString()
	if err != nil {
		LOGGER.Errorf("Error occured while converting the SignedInfo node to String")
		return false
	}

	//Canonicalise the signedinfo
	signedInfoCanonicalisedString, err := Canonicalise(signedInfoString)

	//Remove the Signature element from AppHdr
	appHeaderElement.RemoveChild(sgntrElement)

	payloadStringWithoutSignature, err := signedDocument.WriteToString()
	if err != nil {
		LOGGER.Errorf("Error occured while converting the payload with removed sigature node to String")
		return false
	}

	payloadWithoutSignatureCanonicalised, err := Canonicalise(payloadStringWithoutSignature)
	if err != nil {
		LOGGER.Errorf("Error occured while canonicalising the payloadStringWithoutSignature")
		return false
	}

	payloadDigestWithoutSignature, err := GenerateSHA256Hash(payloadWithoutSignatureCanonicalised)
	if err != nil {
		LOGGER.Errorf("Error occured while hashing the payloadStringWithoutSignature")
		return false
	}

	payloadDigestWithoutSignatureBase64 := base64.StdEncoding.EncodeToString(payloadDigestWithoutSignature)
	LOGGER.Infof("Calculated digest of payload XML: %s", payloadDigestWithoutSignatureBase64)

	if digest != payloadDigestWithoutSignatureBase64 {
		LOGGER.Warningf("Digest in document and digest calculated doesn't match")
		return false
	}

	return doVerifySignatureUsingHSM(signature, signedInfoCanonicalisedString, publicKay)
	//return doVerifySignatureUsingPEM(signature, signedInfoCanonicalisedString, publicKay)
}

func doVerifySignatureUsingPEM(signatureInPayload, signedInfoCanonicalisedString, publicKey string) bool {
	a, _ := ioutil.ReadFile("certificate.pem")
	block, _ := pem.Decode(a)
	certificate, err := x509.ParseCertificate(block.Bytes)
	rsaPublicKey := certificate.PublicKey.(*rsa.PublicKey)

	signatureByte, err := base64.StdEncoding.DecodeString(signatureInPayload)
	if err != nil {
		return false
	}

	hash, err := GenerateSHA256Hash(signedInfoCanonicalisedString)
	if err != nil {
		return false
	}

	err = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hash[:], signatureByte)
	if err != nil {
		LOGGER.Warningf("Signature verification failed. %s", err)
		return false
	} else {
		LOGGER.Infof("Signature verification success")
		return true
	}
}

func doVerifySignatureUsingHSM(signatureInPayload, signedInfoCanonicalisedString, publicKey string) bool {
	kp, err := keypair.Parse(publicKey)
	if err != nil {
		LOGGER.Error("Wrong Public Key for verification")
		return false
	}
	signedInfoCanonicalisedStringHashed, err := GenerateSHA256Hash(signedInfoCanonicalisedString)
	if err != nil {
		LOGGER.Error("Error while generate sha256 hash")
		return false
	}
	signatureByte, err := base64.StdEncoding.DecodeString(signatureInPayload)
	if err != nil {
		LOGGER.Error("error while decoding the signature in payload")
		return false
	}
	verErr := kp.Verify([]byte(signedInfoCanonicalisedStringHashed), signatureByte)
	if verErr != nil {
		LOGGER.Error("Wrong signature for verification: " + verErr.Error())
		return false
	}

	return true
}

func GetElementByName(element *etree.Element, tag string) *etree.Element {
	for _, child := range element.ChildElements() {
		if child.Tag == tag {
			return child
		}
	}

	return nil
}

/*
Canonicalise the XML based on W3C standard, using another library
*/
func Canonicalise(xmlString string) (string, error) {
	decoder := xml.NewDecoder(strings.NewReader(xmlString))
	buf := bytes.NewBufferString("")

	//The library expect a WriteCloser, so create a type to implement the Close method
	xmlWriter := XmlWriteCloser{buf}
	err := c14n.Canonicalise(decoder, xmlWriter, true)
	if err != nil {
		LOGGER.Infof("Error occured while canonicalising the XML payload. %s", err)
		return "", err
	}

	return buf.String(), nil
}

func GenerateSHA256Hash(val string) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(val))
	if err != nil {
		LOGGER.Infof("Error occured while calculating the digest. %s", err)
		return nil, err
	}
	checksum := hash.Sum(nil)
	return checksum, nil
}

func MarshallToXML(data interface{}) (string, error) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	encoder := xml.NewEncoder(writer)
	err := encoder.Encode(data)
	if err != nil {
		LOGGER.Infof("Error occured while marshalling the SignedInfo tag to XML. %s", err)
		return "", err
	}
	err = encoder.Flush()
	if err != nil {
		LOGGER.Infof("Error occured while Flush. %s", err)
		return "", err
	}
	return buffer.String(), nil
}

// Signature element is the root element of an XML Signature.
type Sgntr struct {
	XMLName   xml.Name `xml:"urn:iso:std:iso:20022:tech:xsd:head.001.001.01 Sgntr"`
	Signature Signature
}

type Signature struct {
	XMLName        xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# Signature"`
	SignedInfo     SignedInfo
	SignatureValue string `xml:"http://www.w3.org/2000/09/xmldsig# SignatureValue"`
	KeyInfo        KeyInfo
}

// Algorithm describes the digest or signature used when digest or signature.
type Algorithm struct {
	Algorithm string `xml:",attr"`
}

// SignedInfo includes a canonicalization algorithm, a signature algorithm, and a reference.
type SignedInfo struct {
	XMLName                xml.Name  `xml:"http://www.w3.org/2000/09/xmldsig# SignedInfo"`
	CanonicalizationMethod Algorithm `xml:"http://www.w3.org/2000/09/xmldsig# CanonicalizationMethod"`
	SignatureMethod        Algorithm `xml:"http://www.w3.org/2000/09/xmldsig# SignatureMethod"`
	Reference              Reference
}

// Reference specifies a digest algorithm and digest value, and optionally an identifier of the object being signed, the type of the object, and/or a list of transforms to be applied prior to digesting.
type Reference struct {
	XMLName      xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# Reference"`
	URI          string   `xml:",attr,omitempty"`
	Transforms   Transforms
	DigestMethod Algorithm `xml:"http://www.w3.org/2000/09/xmldsig# DigestMethod"`
	DigestValue  string    `xml:"http://www.w3.org/2000/09/xmldsig# DigestValue"`
}

// Transforms is an optional ordered list of processing steps that were applied to the resource's content before it was digested.
type Transforms struct {
	XMLName   xml.Name    `xml:"http://www.w3.org/2000/09/xmldsig# Transforms"`
	Transform []Algorithm `xml:"http://www.w3.org/2000/09/xmldsig# Transform"`
}

// KeyInfo is an optional element that enables the recipient(s) to obtain the key needed to validate the signature.
type KeyInfo struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# KeyInfo"`
	X509Data *X509Data
	Children []interface{}
}

// X509Data element within KeyInfo contains one an X509 certificate
type X509Data struct {
	XMLName         xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# X509Data"`
	X509Certificate string   `xml:"http://www.w3.org/2000/09/xmldsig# X509Certificate"`
}
