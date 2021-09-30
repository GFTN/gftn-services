World Wire Anchor Service API
=============================
Client Facing API for Anchors to interact with World Wire.


**Version:** 1.0.0

### /address
---
##### ***GET***
**Summary:** Retrieve a participant's ledger address

**Description:** Returns ledger address corresponding to the supplied identifier of a Participant. Learn more about the Ledger in the [Concepts](??base_url??/docs/??version??/concepts) section of World Wire.


**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| name | query | Account name concatenated with the World Wire Participant's ID.  (i.e. 1234554321*uk.barclays.payments.ibm.com) | No | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | The account identifier is recognized as being able to receive value, presenting a account address | [addressLedger](#addressledger) |
| 400 | Missing or invalid parameters in the request | [error](#error) |
| 401 | JWT token in header is invalid | [error](#error) |
| 404 | There is no matching record found for the participant domain | [error](#error) |

### /assets/issued/{anchor_id}
---
##### ***GET***
**Summary:** List your issued assets

**Description:** Returns a list of all your issued assets on World Wire. Learn more about Assets in the [Concepts](??base_url??/docs/??version??/concepts) section of World Wire.


**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| anchor_id | path | Identifier of a World Wire Anchor. To get a list of all participants, make a GET request to /participants.  | Yes | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | All assets issued on World Wire by this anchor participant | [ [asset](#asset) ] |
| 404 | No assets issued on World Wire by this anchor participant | [error](#error) |

### /fundings/instruction
---
##### ***POST***
**Summary:** Create a funding instruction

**Description:** Generates the bytecode instruction necessary to record your transaction on the ledger. Once you receive this instruction, you can use it on the /fundings/send endpoint to complete your funding to other Participants on the WorldWire network. Learn more about Fundings in the [Concepts](??base_url??/docs/??version??/concepts) section of World Wire.


**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Funding | body | Includes all necessary detail about the anchor funding. | Yes | [funding](#funding) |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Succesfully created a funding instruction for the ledger. Here you go. You'll need to sign this and supply it to the /funding/send endpoint before delivery to the Participant.  | [fundingInstruction](#fundinginstruction) |
| 400 | Missing or invalid parameters in the request | [error](#error) |
| 401 | JWT token in header is invalid | [error](#error) |
| 404 | invalid participant or fund request failed | [error](#error) |

### /fundings/send
---
##### ***POST***
**Summary:** Fund a participant

**Description:** Funds a Participant on the World Wire network with stablecoins by signing the ledger instruction you received from the /fundings/instruction endpoint with your secret key. IBM doesn't charge for this, hooray! Learn more about Fundings in the [Concepts](??base_url??/docs/??version??/concepts) section of World Wire.


**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| funding | body | Details about the funding from the Anchor to the Participant.  | Yes | [funding](#funding) |
| funding_signed | query | Signed version of the funding details. | Yes | string |
| instruction_signed | query | You'll receive an unsigned version of this instruction when you first create it from the /fundings/instruction endpoint. Sign it with your secret key and supply it here.  | Yes | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Succesfully funded the Participant by submitting this to the ledger! Here's your receipt.  | [fundingReceipt](#fundingreceipt) |
| 400 | Something went wrong with your funding! You probably forgot or supplied invalid parameters.  | [error](#error) |
| 401 | Something went wrong with your funding! It looks like your JWT token in the header is invalid.  | [error](#error) |
| 404 | Something went wrong with your funding! You probably supplied an invalid participant.  | [error](#error) |

### /participants
---
##### ***GET***
**Summary:** List all participants

**Description:** Retrieves a list of all active Participants and associated data on World Wire. Learn more about Participants in the [Concepts](??base_url??/docs/??version??/concepts) section of World Wire.


**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| country_code | query | The 2-digit code of the country where the Participants are located. | No | string |
| asset_code | query | The identifier of the asset balance being queried. For a list of assets, retrieve all World Wire assets from the /assets endpoint. | No | string |
| issuer_id | query | Identifier of the Issuer of this asset.  | No | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Retrieved Participants on World Wire | [ [participant](#participant) ] |
| 404 | No Participants found on World Wire | [error](#error) |

### /participants/{participant_id}
---
##### ***GET***
**Summary:** Retrieve a specific participant

**Description:** Retrieves a specific Participant and their associated data on World Wire. Learn more about Participants in the [Concepts](??base_url??/docs/??version??/concepts) section of World Wire.


**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| participant_id | path | Identifier of a WorldWire Participant. To get a list of all participants, make a GET request to /participants.  | Yes | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Successfully retrieved a WorldWire Participant. | [participant](#participant) |
| 404 | Whoops, no Participant found with that ID on World Wire. | [error](#error) |

### /trust/{anchor_id}
---
##### ***POST***
**Summary:** Submit asset trust permissions

**Description:** Changes the trust relationship you have with an OFI Participant. As the Anchor, you can request, allow, or revoke permission to transact with an OFI by supplying a corresponding permission on this request. Learn more about Trust in the [Concepts](??base_url??/docs/??version??/concepts) section of World Wire.


**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| anchor_id | path | Identifier of a World Wire Anchor. To get a list of all participants, make a GET request to /participants.  | Yes | string |
| Trust | body | Indicate who you are trusting | Yes | [trust](#trust) |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Allow Trust was successful |  |
| 400 | Missing or invalid parameters in the request | [error](#error) |
| 401 | JWT token in header is invalid | [error](#error) |
| 404 | invalid participant or allow trust failed | [error](#error) |

### Models
---

### account  

Account

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | string | The address that was created on the ledger. | Yes |
| name | string | A name to identify this account. | No |

### addressLedger  

Address Ledger

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| account_name | string | Can be either 'issuing' or the Participants operating account's name. | Yes |
| address | string | The ledger address which is expected to be the recipient for this transaction, once compliance checks are complete. | Yes |

### asset  

Details of the asset being transacted

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| asset_code | string | Alphanumeric code for the asset - USD, XLM, etc | Yes |
| asset_type | string | The type of asset. Options include digital obligation, "DO", digital asset "DA", or a cryptocurrency "native". | Yes |
| issuer_id | string | The asset issuer's participant id. | No |

### error  

Error

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| build_version | string | build version string of micro service | No |
| code | string | Error code reference. | No |
| details | string | Detailed message description about error. | Yes |
| message | string | Short message description about error. | Yes |
| participant_id | string | participant id, same as home domain as set by environment variables | No |
| service | string | name of micro service | No |
| time_stamp | number (int64) | The timestamp of the occurance. | Yes |
| type | string | Type is for query purposes, it an identifier to assist with troubleshooting where an error came from (eg, containing func name) tells us if it is originating from NotifyWWError vs. NotFound vs. some other spot | No |
| url | string | Url of endpoint that failed. | No |

### funding  

Details about a Funding

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| account_name | string | The name of an operating account or "issuing" for an issuing account.
 | No |
| amount_funding | number | The amount that the Anchor is funding the Participant. | Yes |
| anchor_id | string | Identifier of the World Wire Anchor that will fund the Participant with stablecoins. (i.e., "thebestbankintheUK")
 | Yes |
| asset_code_issued | string | Identifier of the stable coin asset issued by the Anchor. For a list of assets, retrieve all World Wire assets from the /assets endpoint.
 | Yes |
| end_to_end_id | string | Generated by the anchor, a unique ID for this funding request | Yes |
| memo_transaction | string | An optional way for anchor to name a transaction. | No |
| participant_id | string | Identifier of the World Wire Participant that will receive the funding.
 | Yes |

### fundingInstruction  

Funding Instruction

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| details_funding | [funding](#funding) |  | No |
| instruction_unsigned | string | Unsigned transaction xdr related to the funding. This will need to be signed in the next step.
 | No |

### fundingReceipt  

Funding Receipt

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| details_funding | [funding](#funding) |  | No |
| receipt_funding | [transactionReceipt](#transactionreceipt) |  | No |

### participant  

Participant

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| bic | string | The business identifier code of each participant | Yes |
| country_code | string | Participant's country of residence, country code in ISO 3166-1 format | Yes |
| id | string | The participant id for the participant | Yes |
| issuing_account | string | The ledger address belonging to the issuing account. | No |
| operating_accounts | [ [account](#account) ] | Accounts | No |
| role | string | The Role of this registered participant, it can be MM for Market Maker and IS for Issuer or anchor | Yes |
| status | string | Participant active status on WW network, inactive, active, suspended | No |
| url_callback | string | Callback url of the finiancial institute's backend system. | Yes |

### transactionReceipt  

Transaction Receipt

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| time_stamp | number (int64) | The timestamp of the transaction. | Yes |
| transaction_id | string | A unique transaction identifier generated by the ledger. | Yes |
| transaction_status | string | For DA (digital asset) or DO (digital obligation) ops, this will be "cleared".  For cryptocurrencies, this will be "settled". | Yes |

### trust  

Trust

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| account_name | string | This is account name that is trusting of the asset. Options include "default", "issuing", or another string that identifies another operating account.
 | Yes |
| asset_code | string | The 3-letter code identifying the asset being trusted. For a list of assets retrieve all World Wire assets from the [/assets](??base_url??/docs/??version??/api/participant-client-api?jump=path_get__assets) endpoint.
 | Yes |
| end_to_end_id | string | Generated by requester, a unique ID for this entire trust flow | No |
| limit | integer | The trust limit for this asset between source and issuer. This parameter is only necessary when the trust permission you are submitting is "request".
 | No |
| participant_id | string | When the permission submitted by an OFI is "request", this is the identifier of the RFI who issued the asset. However, when the permission submitted by an RFI is "allow", this is the OFI's identifier (i.e., uk.yourbankintheUK.payments.ibm.com). Make sure you request trust first to the RFI's issuing account, and then also their operating account.
 | Yes |
| permission | string | This string identifier represents the level of trust you would like to set with another participant in your trust object. Options are "request", "allow", or "revoke".
 | Yes |