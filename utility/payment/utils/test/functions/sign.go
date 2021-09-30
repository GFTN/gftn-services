// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package functions

import (
	"encoding/base64"
	"errors"
	"log"

	"github.com/beevik/etree"
	"github.com/stellar/go/keypair"
	"github.com/GFTN/gftn-services/utility/nodeconfig"
	"github.com/GFTN/gftn-services/utility/xmldsig"
)

func Sign() []byte {

	cbMsg := `<?xml version="1.0" encoding="UTF-8"?>
	<Message xmlns="urn:worldwire">
		<AppHdr>
			<Fr xmlns="urn:iso:std:iso:20022:tech:xsd:head.001.001.01">
				<FIId>
					<FinInstnId>
						<BICFI>IBMAUSA1001</BICFI>
						<Othr>
							<Id>ibmanchor</Id>
						</Othr>
					</FinInstnId>
				</FIId>
			</Fr>
			<To xmlns="urn:iso:std:iso:20022:tech:xsd:head.001.001.01">
				<FIId>
					<FinInstnId>
						<BICFI>WORLDWIRE00</BICFI>
						<Othr>
							<Id>WW</Id>
						</Othr>
					</FinInstnId>
				</FIId>
			</To>
			<BizMsgIdr xmlns="urn:iso:std:iso:20022:tech:xsd:head.001.001.01">B20190404WORLDWIRE00BHU0928374</BizMsgIdr>
			<MsgDefIdr xmlns="urn:iso:std:iso:20022:tech:xsd:head.001.001.01">pacs.002.001.09</MsgDefIdr>
			<CreDt xmlns="urn:iso:std:iso:20022:tech:xsd:head.001.001.01">2019-05-02T16:00:00Z</CreDt>
		</AppHdr>
		<FIToFIPmtStsRpt>
			<GrpHdr xmlns="urn:iso:std:iso:20022:tech:xsd:pacs.002.001.09">
				<MsgId>SGDDO20190819SGPTTEST003B3889747664</MsgId>
				<CreDtTm>2019-05-24T07:23:32</CreDtTm>
				<InstgAgt>
					<FinInstnId>
						<BICFI>IBMAUSA1001</BICFI>
						<Othr>
							<Id>ibmanchor</Id>
						</Othr>
					</FinInstnId>
				</InstgAgt>
				<InstdAgt>
					<FinInstnId>
						<BICFI>SGPTTEST003</BICFI>
						<Othr>
							<Id>testparticipant3dev</Id>
						</Othr>
					</FinInstnId>
				</InstdAgt>
			</GrpHdr>
			<TxInfAndSts xmlns="urn:iso:std:iso:20022:tech:xsd:pacs.002.001.09">
				<OrgnlInstrId>SGDDO20190819SGPTTEST003B3889747663</OrgnlInstrId>
				<OrgnlEndToEndId>SGDDO19082019SGPTTEST00377793380333</OrgnlEndToEndId>
				<OrgnlTxId>SGDDO19082019SGPTTEST00377793380333</OrgnlTxId>
				<TxSts>ACTC</TxSts>
				<ChrgsInf>
					<Amt Ccy="TWD">1.00</Amt>
					<Agt>
						<FinInstnId>
							<BICFI>IBMAUSA1001</BICFI>
							<Othr>
								<Id>ibmanchor</Id>
							</Othr>
						</FinInstnId>
					</Agt>
				</ChrgsInf>
				<OrgnlTxRef>
					<IntrBkSttlmAmt Ccy="TWD">2</IntrBkSttlmAmt>
				</OrgnlTxRef>
				<SplmtryData>
					<PlcAndNm>/Message/FIToFIPmtStsRpt/TxInfAndSts/Issr</PlcAndNm>
					<Envlp>
						<Id>ibmanchor</Id>
					</Envlp>
				</SplmtryData>
				<SplmtryData>
					<PlcAndNm>/Message/FIToFIPmtStsRpt/TxInfAndSts/SttlmAcctAddr</PlcAndNm>
					<Envlp>
						<Id>GBKV4YBX44WX3N4UDQKDR7UURSOD3P3QSWD2EKCAGCSLYEKBJSALVP5V</Id>
					</Envlp>
				</SplmtryData>
				<SplmtryData>
					<PlcAndNm>/Message/FIToFIPmtStsRpt/TxInfAndSts/PayToRef</PlcAndNm>
					<Envlp>
						<Id>7777</Id>
					</Envlp>
				</SplmtryData>
			</TxInfAndSts>
		</FIToFIPmtStsRpt>
	</Message>`
	var signErr error
	log.Printf("Signing with utility function")
	var signedRes string
	signedRes, signErr = SignPayload(string(cbMsg))
	signedMessage := []byte(signedRes)
	if signErr != nil {
		log.Fatalf("Failed to sign payload: %v", signErr.Error())
		return []byte{}
	}

	return signedMessage
}

func SignPayload(xml string) (string, error) {

	account := nodeconfig.Account{
		NodeAddress: "<anchor_addr>",
		NodeSeed:    "<anchor_seed>",
	}

	publicKey := account.NodeAddress
	seed := account.NodeSeed
	// Read the XML string and convert it to DOM object in memory
	payloadDocument := etree.NewDocument()
	err := payloadDocument.ReadFromString(xml)
	if err != nil {
		log.Fatalf("%v", err)
		return "", err
	}

	appHeaderElement := xmldsig.GetElementByName(payloadDocument.Root(), "AppHdr")
	if appHeaderElement == nil {
		log.Println("AppHdr is missing")
		return "", errors.New("AppHdr is missing")

	}

	sgntrElement := xmldsig.GetElementByName(appHeaderElement, "Sgntr")
	if sgntrElement != nil {
		log.Fatalf("Sgntr tag is already available")
		return "", errors.New("Sgntr tag is already available")
	}

	//Initialise the Signature object with properties, values are hard coded only support SHA-256
	sgntr := xmldsig.NewSignature()

	//Canonicalise the original xml payload
	canonicalisedPayload, err := xmldsig.Canonicalise(xml)
	if err != nil {
		log.Fatalf("%v", err)
		return canonicalisedPayload, err
	}

	//Get the digest of the Canonicalised payload
	canonicalisedPayloadHashed, err := xmldsig.GenerateSHA256Hash(canonicalisedPayload)
	if err != nil {
		log.Fatalf("%v", err)
		return "", err
	}

	//Based64 encode the canonicalised hashed payload
	canonicalisedPayloadHashedEncoded := base64.StdEncoding.EncodeToString(canonicalisedPayloadHashed)
	//set the hash to the Signature struct
	sgntr.Signature.SignedInfo.Reference.DigestValue = canonicalisedPayloadHashedEncoded

	//Now generate the SignedInfo tag by marshalling the struct to XML string
	signedInfoString, err := xmldsig.MarshallToXML(sgntr.Signature.SignedInfo)
	if err != nil {
		log.Fatalf("%v", err)
		return signedInfoString, err
	}

	//Canonicalise the SignedInfo
	signedInfoStringCanonicalised, err := xmldsig.Canonicalise(signedInfoString)

	checkSumToSign, err := xmldsig.GenerateSHA256Hash(signedInfoStringCanonicalised)
	if err != nil {
		log.Fatalf("%v", err)
		return "", err
	}

	kp, err := keypair.Parse(seed)
	if err != nil {
		log.Fatalf("There was an error while getting node seed")
		return "", err
	}
	signedInfoSignature, err := kp.Sign(checkSumToSign)
	if err != nil {
		log.Fatalf("%v", err)
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
		log.Fatalf("%v", err)
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
		log.Fatalf("%v", err)
		return s, err
	}

	xmldsig.VerifySignature(s)
	return s, nil
}
