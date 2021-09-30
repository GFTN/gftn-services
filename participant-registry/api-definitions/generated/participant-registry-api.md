# The Participant Registry API
This is a registry of Stellar Participants that contains configuration details for each participant

## Version: 1.0.0

### /internal/pr

#### GET
##### Summary:

Get list of all participants on WW

##### Description:

Get list of all participants on WW


##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | all Participants on World Wire | [ [participant](#participant) ] |
| 404 | there is no participant for this country |  |

#### POST
##### Summary:

Create a new participant

##### Description:

Sends a request to the Participant Registry to create an participant


##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| body | body | The participant data | Yes | [participant](#participant) |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | Participant created successfully |
| 404 | Participant could not be created |

### /internal/pr/account/{participant_id}

#### POST
##### Summary:

Save Participant Operating account

##### Description:

Saves Participant Operating account


##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| participant_id | path | participant domain for the participant | Yes | string |
| body | body | The participant Operating data | Yes | [operatingAccount](#operatingaccount) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Participant operating account saved successfully |  |
| 404 | there is no participant for this participant_id | string |
| 409 | Participant operating account already exists | string |

### /internal/pr/account/{participant_id}/{account_name}

#### GET
##### Summary:

Get the pub key for the given operating account name

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| participant_id | path | the participant domain for this participant | Yes | string |
| account_name | path | the participant's operating account name | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Found participant operating account key for given account name | string |
| 404 | there is no participant disctribution account for given name | string |

### /internal/pr/country/{country_code}

#### GET
##### Summary:

Get List of participants operating in the given country

##### Description:

Get List of participants operating in the given country

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| country_code | path | country code | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Found participants for this country | [ [participant](#participant) ] |
| 404 | there is no participant for this country |  |

### /internal/pr/domain/{participant_id}

#### GET
##### Summary:

Get the configuration details for the participant idenfied by his participant domain

##### Description:

Get the configuration details for the participant idenfied by his participant domain

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| participant_id | path | the participant domain for this participant | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Found participant for this domain | [participant](#participant) |
| 404 | there is no participant for this participant_id |  |

### /internal/pr/issuingaccount/{participant_id}

#### POST
##### Summary:

Save Participant Issuing account

##### Description:

Saves Participant Issuing account


##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| participant_id | path | participant domain for the participant | Yes | string |
| body | body | The participant Issuing data | Yes | [operatingAccount](#operatingaccount) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Participant Issuing account saved successfully |  |
| 404 | there is no participant for this participant_id | string |
| 409 | Participant Issuing account already exists | string |

### /internal/pr/{participant_id}

#### PUT
##### Summary:

Update an existing participant

##### Description:

Sends a request to the Participant Registry to to update an existing participant


##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| participant_id | path | participant domain for the participant | Yes | string |
| body | body | The participant data | Yes | [participant](#participant) |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | participant successfully updated |
| 404 | Participant not found |

### /internal/pr/{participant_id}/status

#### PUT
##### Summary:

Save Participant WW network status, its a admin api

##### Description:

Saves Participant WW network status, its a admin api


##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| participant_id | path | participant domain for the participant | Yes | string |
| body | body | The participant status | Yes | [participantStatus](#participantstatus) |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | Participant network status updated successfully |

### Models


#### operatingAccount

Account with the token

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| account | object | Account | No |
| token | string | auth token | No |

#### participant

Participant

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| bic | string | The business identifier code of each participant | Yes |
| country_code | string | Participant's country of residence, country code in ISO 3166-1 format | Yes |
| id | string | The participant id for the participant | Yes |
| issuing_account | string | The ledger address belonging to the issuing account. | No |
| operating_accounts | [ object ] | Accounts | No |
| role | string | The Role of this registered participant, it can be MM for Market Maker and IS for Issuer or anchor | Yes |
| status | string | Participant active status on WW network, inactive, active, suspended | No |

#### participantStatus

ParticipantStatus

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| status | string | Participant active status on WW network, inactive, active, suspended | Yes |