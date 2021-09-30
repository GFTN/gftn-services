# World Wire internal signing API Service for individual Participant
Internal API service which takes in unsigned transaction envelope as input and returns back a signed transaction envelop


## Version: 1.0.0

### /account/{account_name}

#### POST
##### Summary:

Retreive an account.

##### Description:

Retrieves an Issuing or Operating Account after it is newly created.


##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| account_name | path | name of the new account to be created | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Successfully retrieved account. | [account](#account) |
| 400 | Invalid create account request | [error](#error) |
| 409 | conflict or error | [error](#error) |
| 424 | conflict or error | [error](#error) |

### /admin/account

#### POST
##### Summary:

Retreive IBM account public address.

##### Description:

Retrieves IBM account public address.


##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Successfully retrieved account address. | [account](#account) |
| 400 | Invalid create account request | [error](#error) |
| 409 | conflict or error | [error](#error) |
| 424 | conflict or error | [error](#error) |

### /admin/sign

#### POST
##### Summary:

returns signed envelope with IBM token account signature

##### Description:

This API service which takes in unsigned transaction envelope as input and returns back a signed transaction envelope


##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Draft | body | This is a internal request model for signing request | Yes | [internalDraft](#internaldraft) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Successfully signed the transaction. Here you go. | [signature](#signature) |
| 404 | Invalid signing request | [error](#error) |

### /request/sign

#### POST
##### Summary:

Create a signature

##### Description:

Accepts a draft payload as input and returns a signature (signed version) back.


##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| payload | body | The payload that needs to be signed. | Yes | [requestPayload](#requestpayload) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Signed successfully. Here's you go. | [signature](#signature) |
| 404 | Invalid signing request | [error](#error) |

### /sign

#### POST
##### Summary:

returns signed envelope

##### Description:

This API service which takes in unsigned transaction envelope as input and returns back a signed transaction envelope


##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Draft | body | This is a request model for signing request | Yes | [draft](#draft) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Successfully signed the transaction. Here you go. | [signature](#signature) |
| 404 | Invalid signing request | [error](#error) |

### Models


#### account

Account

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | string | The address that was created on the ledger. | Yes |
| name | string | A name to identify this account. | No |

#### draft

draft

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| account_name | string | The name of the account with which the transactions needs to be signed | Yes |
| id_signed | byte | This will be signed reference envelope to verify against partcipant's signature for authenticity. | Yes |
| id_unsigned | byte | This will be unsigned reference envelope to verify against partcipant's signature for authenticity. | Yes |
| transaction_id | string | Identifier for transaction. | No |
| transaction_unsigned | byte | The unsigned transaction envelope to be signed by the participant. | Yes |

#### error

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

#### internalDraft

draft

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| account_name | string | The name of the account with which the transactions needs to be signed | Yes |
| transaction_id | string | Identifier for transaction. | No |
| transaction_unsigned | byte | The unsigned transaction envelope to be signed by IBM account. | Yes |

#### requestPayload

requestPayload

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| account_name | string | The name of the account with which the payload needs to be signed | Yes |
| payload | byte | unsigned request payload to be signed | Yes |

#### signature

signature

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| transaction_signed | byte | Transaction signed by Participant. | Yes |