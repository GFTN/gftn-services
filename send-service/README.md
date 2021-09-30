# Send Service
Participant service for sending a transaction request (ISO20022 standard) to other participant.

## ISO20022 payment messages
**Payments Clearing and Settlement (pacs)**\
Please find more information about the ISO20022 payment messages [here](https://www.iso20022.org/payments_messages.page).

## ISO4217 currency code (alpha)
The currency code used in settlement amount and transaction fee will follow the [ISO4217](https://en.wikipedia.org/wiki/ISO_4217) standard.

## API
> **/transactions/send**

For OFI to send a transaction request which include transaction details and a encoded base 64 ISO20022 pacs008 message
#### payload
```
{
	"message_type": "iso20022:pacs.008.001.07",
	"message": "encoded pacs008 XML message"
}
```

#### pacs008 example
```

```

> **/transactions/receive**

For RFI to send a encoded base 64 ibwf001 federation and compliance result message back to the OFI.
#### payload
```
{
	"message_type": "iso20022:ibwf.001.001.01",
	"message": "encoded ibwf001 XML message"
}
```

#### ibwf001 example
```
<?xml version="1.0" encoding="UTF-8"?>
<Document xmlns="urn:iso:std:iso:20022:tech:xsd:ibwf.001.001.01">
    <FedCompRes>
        <GrpHdr>
            <MsgId>USDDO26012019USDVWXYZ77700000000001</MsgId>
            <CreDtTm>2015-09-28T16:00:00</CreDtTm>
            <NbOfTxs>1</NbOfTxs>
            <SttlmInf>
                <SttlmMtd>{WWDO WWDA XLM}</SttlmMtd>
                <SttlmAcct>
                    <Id>
                        <Othr>
                            <Id>{RFI id}</Id>
                        </Othr>
                    </Id>
                    <Nm>{issuing or name of operating account}</Nm>
                </SttlmAcct>
            </SttlmInf>
            <InstgAgt>
                <FinInstnId>
                    <BICFI>USDVWXYZ777</BICFI>
                    <Othr>
                        <Id>{RFI id}</Id>
                    </Othr>
                </FinInstnId>
            </InstgAgt>
            <InstdAgt>
                <FinInstnId>
                    <BICFI>TWDABCDE101</BICFI>
                    <Othr>
                        <Id>{OFI id}</Id>
                    </Othr>
                </FinInstnId>
            </InstdAgt>
        </GrpHdr>
        <FedRes>
            <AccId>{RFI public key}</AccId>
            <FedSts>ACTC</FedSts>
            <PmtId>
                <EndToEndId>TWDDO25012019TWDABCDE10100000000001</EndToEndId>
                <TxId>TWDDO25012019TWDABCDE10100000000001</TxId>
            </PmtId>
        </FedRes>
        <CmpRes>
            <InfSts>ACTC</InfSts>
            <TxnSts>ACTC</TxnSts>
            <PmtId>
                <EndToEndId>TWDDO25012019TWDABCDE10100000000001</EndToEndId>
                <TxId>TWDDO25012019TWDABCDE10100000000001</TxId>
            </PmtId>
        </CmpRes>
    </FedCompRes>
</Document>
```
## Send Status Response
After sending a request to send service, OFI will receive a encoded base 64 ISO20022 pacs002 message.\
Decode the pacs002 message, you will see the status under the tag "Rsn(reason)".

#### status type

| Status | Code| Description |
| :--- | :--- | :--- |
| `RJCT` | `1001` | Can not parse the request body |
| `RJCT` | `1002` | Unable to validate the send request |
| `RJCT` | `1003` | Unable to set XML file for validation |
| `RJCT` | `1004` | Error validate the XML file |
| `RJCT` | `1005` | Error finding account address from participant registry |
| `RJCT` | `1006` | RFI was not whitelisted by OFI |
| `RJCT` | `1007` | OFI was not whitelisted by RFI |
| `RJCT` | `1008` | Unable to sign the payload |
| `RJCT` | `1009` | Error parsing xml message into Go structure or ProtoBuffer |
| `RJCT` | `1010` | Unable to send ProtoBuffer to Kafka broker |
| `RJCT` | `1011` | Error verifying OFI signature |
| `RJCT` | `1012` | Error verifying RFI signature |
| `RJCT` | `1013` | Unable to get IBM account from gas-service |
| `RJCT` | `1014` | Error signing the transaction envelope |
| `RJCT` | `1015` | Error sending transaction to Stellar network |
| `RJCT` | `1016` | Unable to communicate with RFI backend system |
| `RJCT` | `1017` | Can not use the same message ID as the undergoing or previous request |
| `RJCT` | `1018` | Please use the same settlement method as the payment request |
| `RJCT` | `1019` | Failed to get transaction related information from database |
| `ACTC` | `1101` | Payment request had been successfully sent to Kafka broker, waiting for the response form RFI |
| `ACTC` | `1102` | Payment response had been successfully sent to Kafka broker |
| `RCVD` | `1103` | Transaction successfully sent to Stellar network |
| `RJCT` | `1201` | RFI reject the payment request |

#### response json
```
{
  "message_type": "pacs.002.001.09",
  "message": "encoded pacs002 XML message"
}
```

#### pacs002 example

```
<?xml version="1.0" encoding="UTF-8"?>
<Document xmlns="urn:iso:std:iso:20022:tech:xsd:pacs.002.001.09">
	<FIToFIPmtStsRpt>
		<GrpHdr>
			<MsgId>TWDDO25012019TWDABCDE10100000000002</MsgId>
			<CreDtTm>2019-01-26T00:58:51</CreDtTm>
			<InstgAgt>
				<FinInstnId>
					<BICFI>TWDABCDE101</BICFI>
				</FinInstnId>
			</InstgAgt>
			<InstdAgt>
				<FinInstnId>
					<BICFI>USDVWXYZ777</BICFI>
				</FinInstnId>
			</InstdAgt>
		</GrpHdr>
		<TxInfAndSts>
			<StsId>AB/8568</StsId>
			<OrgnlEndToEndId>TWDDO25012019TWDABCDE10100000000001</OrgnlEndToEndId>
			<OrgnlTxId>TWDDO25012019TWDABCDE10100000000001</OrgnlTxId>
			<TxSts>ACTC</TxSts>
			<StsRsnInf>
				<Rsn>
					<Cd>1006</Cd>
					<Prtry>Payment successfully sent to Kafka broker, waiting for RFI response</Prtry>
				</Rsn>
			</StsRsnInf>
		</TxInfAndSts>
	</FIToFIPmtStsRpt>
</Document>
```