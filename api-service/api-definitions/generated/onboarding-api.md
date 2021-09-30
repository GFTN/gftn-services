# GFTN API for onboarding a new client
API endpoints for creating issuing, operating accounts, new assets

## Version: 1.0.0

### /accounts/{account_name}

#### GET
##### Summary:

retrieve a operating or issuing account

##### Description:

retrieve the operating or issuing account for a participant


##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| account_name | path | a name to identify this account, use "issuing" as account_name for issuing account | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Account found | [operatingAccount](#operatingaccount) |
| 400 | missing or invalid parameters in the request |  |
| 404 | The operating account could not be found |  |

#### POST
##### Summary:

Called when an originator wants to create a new issuing account

##### Description:

During client onboarding, each client will require one issuing account. this API creates the issuing account


##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| account_name | path | a name to identify this account | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 208 | Issuing Account created or already exists | [operatingAccount](#operatingaccount) |
| 400 | missing or invalid parameters in the request |  |
| 404 | The issuing account could not be created |  |

### Models


#### account

Account

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | string | The address that was created on the ledger. | Yes |
| name | string | A name to identify this account. | No |

#### operatingAccount

Account with the token

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| account | [account](#account) |  | No |
| token | string | auth token | No |