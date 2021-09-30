WW Administration Service Internal; API
=======================================
Private API for inter-component communication within the Participant's World Wire instance.

**Version:** 1.0.0

### /blocklist
---
##### ***GET***
**Summary:** Called when a participant wants to lookup if a certain currency/institution/country is in the blocklist or not

**Description:** This endpoint will search for the existing record in the blocklist that meets the searching type


**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| type | query | The type of the blocklist record. | No | string |

**Responses**

| Code | Description |
| ---- | ----------- |
| 200 | Blocklist record found |
| 400 | Blocklist record not found due to malformed payload |
| 404 | Blocklist record not found |
| 500 | Internal server error |

##### ***POST***
**Summary:** Called when a currency/country/institution needs to be added into the blocklist

**Description:** This endpoint will create a new record in the blocklist if there doesn't have an existing blocklist record


**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| blocklist | body | The blocklist record that needs to be added. | Yes | [blocklist](#blocklist) |

**Responses**

| Code | Description |
| ---- | ----------- |
| 200 | New blocklist record created |
| 400 | Blocklist record could not be created due to the record already exists or malformed payload |
| 404 | Blocklist record could not be created |
| 500 | Internal server error |

##### ***DELETE***
**Summary:** Called when a currency/country/institution needs to be removed from the blocklist

**Description:** This endpoint will remove an existing record in the blocklist


**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| blocklist | body | The blocklist record that needs to be removed. | Yes | [blocklist](#blocklist) |

**Responses**

| Code | Description |
| ---- | ----------- |
| 200 | Blocklist record removed |
| 400 | Blocklist record could not be removed due to malformed payload |
| 404 | No blocklist record found |
| 500 | Internal server error |

### /blocklist/validate
---
##### ***POST***
**Summary:** Called when a currency/country/institution needs to be validated before transaction

**Description:** This endpoint will check if the query value is in the blocklist or not.


**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| blocklist | body | The blocklist record that needs to be added. | Yes | [ [blocklist](#blocklist) ] |

**Responses**

| Code | Description |
| ---- | ----------- |
| 200 | Validation complete |
| 400 | Could not validate due to malformed payload |
| 404 | Blocklist record not found |
| 500 | Internal server error |

### /fitoficct
---
##### ***POST***
**Summary:** send the fitoficct transaction to WW Admin Service for storage

**Description:** sends the details of the recently submitted fitoficct transaction to the WW Admin Service. the hash value of FItoFICCTMemoData is stored in the ledger memo field The PII data is hashed for security


**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| FItoFICCTMemoData | body | The fitoficct data | Yes | [fitoFICCTMemoData](#fitoficctmemodata) |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | the transactionMemoData is successfully persisted by Admin Service |  |
| 400 | input parameters not acceptable or some error happened | [error](#error) |

### /reactivate/{participant_id}/{account_name}
---
##### ***POST***
**Summary:** Undoes suspension of a suspended Participant in WW.

**Description:** sends transaction to Stellar Network using IBM account and raises the Participant's master key weight to 2, signing thresholds to [1,2,3] and adds a new SHA256 signer to the signing list

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| account_name | path | The address of the account to be re-activated. | Yes | string |
| participant_id | path | The id of the participant | Yes | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Account has been activated successfully. |  |
| 400 | Input parameter not acceptable or some error happened | [error](#error) |
| 500 | Internal Server Error happened. | [error](#error) |

### /suspend/{participant_id}/{account_name}
---
##### ***POST***
**Summary:** Suspend Participant from doing any activities in WW

**Description:** sends transaction to Stellar Network using IBM account and SHA256 signer and make Participant's master key weight 0, threshold to [1,1,1] and removes SHA256 signer from the signing list

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| account_name | path | The address of the account to be suspended | Yes | string |
| participant_id | path | The id of the participant | Yes | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Account has been suspended successfully. |  |
| 400 | Input parameter not acceptable or some error happened | [error](#error) |
| 500 | Internal Server Error happened. | [error](#error) |

### /transaction
---
##### ***POST***
**Summary:** Query transactions

**Description:** Query transaction's details by End-to-End Id or Stellar Transaction Id

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| FItoFITransactionRequest | body | Request containing ID, type and domain details to query transactions. | Yes | [fitoFITransactionRequest](#fitofitransactionrequest) |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Valid Transaction details according to the ID provided | [transaction](#transaction) |
| 400 | Missing or invalid parameters in the request | [error](#error) |
| 404 | No data found for the criteria quried. | [error](#error) |

### Models
---

### asset  

Details of the asset being transacted

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| asset_code | string | Alphanumeric code for the asset - USD, XLM, etc | Yes |
| asset_type | string | The type of asset. Options include digital obligation, "DO", digital asset "DA", or a cryptocurrency "native". | Yes |
| issuer_id | string | The asset issuer's participant id. | No |

### blocklist  

A blocklist that records all the currencies/countries/particpants that is forbidden to transact with

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string | The id of the block type | No |
| name | string | The name of the block type | No |
| type | string | The type of the blocklist element | Yes |
| value | [ string ] | The value of the block type | Yes |

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

### fee  

Fee

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| cost | number | The fee amount, should be a float64 number | Yes |
| cost_asset | [asset](#asset) |  | Yes |

### fitoFICCTMemoData  

FI to FI CCT Memo Data - the hash value of this will be stored in the transaction memo field

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| fitoficct_non_pii_data | [fitoFICCTNonPiiData](#fitoficctnonpiidata) |  | Yes |
| fitoficct_pii_hash | string | The hash value of the FI to FI CCT Pii Data | Yes |
| id | string | Unique autogenerate ID for mongoDB primary key | No |
| message_type | string | This is the message type of the transaction request | Yes |
| ofi_id | string | The participant id of the OFI (payment sender) | Yes |
| time_stamp | number (int64) | The timestamp for this transaction | Yes |
| transaction_identifier | [ string ] | This is the unique id for this transaction generated by the distributed ledger (but not in txn memo hash) | No |
| transaction_status | [ [transactionReceipt](#transactionreceipt) ] | This would capture the new status of a transaction while transaction travel through payment flow. | Yes |

### fitoFICCTNonPiiData  

FI to FI CCT Non-Pii Data

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| account_name_send | string | The name of the operating or issuing account from which the payment is to be sent | Yes |
| creditor_payment_address | string | The RFI address where the payment is to be sent - received during federation protocol | No |
| end_to_end_id | string | Generated by originator, a unique ID for this entire use case | Yes |
| exchange_rate | number | The exchange rate between settlement asset and beneficiary asset. not required if asset is same | Yes |
| original_message_id | string | This is the reference to the original credit transfer message | Yes |
| transaction_details | [transactionDetails](#transactiondetails) |  | Yes |

### fitoFITransactionRequest  

Transaction GET request parameters

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| end_date | date | End Date of the range in which transactions are being quried. | No |
| ofi_id | string | A name to identify from which OFI the request is coming from | Yes |
| page_number | long | Page number for pagination. | No |
| query_data | string | A name to identify the transaction | No |
| query_type | string | A type to identify what kind of data is passed | Yes |
| start_date | date | Start Date of the range in which transactions are being quried. | No |

### transaction  

Transaction

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| transaction_details | [transactionDetails](#transactiondetails) |  | Yes |
| transaction_receipt | [transactionReceipt](#transactionreceipt) |  | Yes |

### transactionDetails  

Transaction Details

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| amount_beneficiary | number | The amount the beneficiary should receive in beneficiary currency | Yes |
| amount_settlement | number | The amount of the settlement. | Yes |
| asset_code_beneficiary | string | The asset code for the beneficiary | Yes |
| asset_settlement | [asset](#asset) |  | Yes |
| fee_creditor | [fee](#fee) |  | Yes |
| ofi_id | string | The ID that identifies the OFI Participant on the WorldWire network (i.e. uk.yourbankintheUK.payments.ibm.com). | Yes |
| rfi_id | string | The ID that identifies the RFI Participant on the WorldWire network (i.e. uk.yourbankintheUK.payments.ibm.com). | Yes |
| settlement_method | string | The preferred settlement method for this payment request (DA, DO, or XLM) | Yes |

### transactionReceipt  

Transaction Receipt

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| time_stamp | number (int64) | The timestamp of the transaction. | Yes |
| transaction_id | string | A unique transaction identifier generated by the ledger. | Yes |
| transaction_status | string | For DA (digital asset) or DO (digital obligation) ops, this will be "cleared".  For cryptocurrencies, this will be "settled". | Yes |