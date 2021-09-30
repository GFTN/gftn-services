World Wire Anchor onboarding Service API
========================================
Onboarding admin API for Anchor service setup World Wire.


**Version:** 1.0.0

### /admin/anchor/assets/issued/{anchor_id}
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

### /admin/anchor/{anchor_id}/onboard/assets
---
##### ***POST***
**Summary:** Registers asset issued by an anchor on world wire

**Description:** Creates trust line with IBM admin account

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| anchor_id | path | anchor id | Yes | string |
| asset_code | query | Asset code of the Digital Asset, should be a 3-letter ISO currency code | Yes | string |
| asset_type | query | Asset type can only be a digital obligation (DA) since issued by a anchor participant | Yes | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | The asset has been issued | [asset](#asset) |
| 400 | Missing or invalid parameters in the request | [error](#error) |
| 404 | The asset could NOT be issued due to error retrieving Issuing Account | [error](#error) |
| 500 | The asset could NOT be issued due to error communicating with ledger | [error](#error) |

### /admin/anchor/{anchor_id}/register
---
##### ***POST***
**Summary:** Registers anchor domain to ww

**Description:** Registers anchor domain on WW participant registry with given issuing account address and generates and returns authentication token, used to generate JWT token to access apis

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| anchor_id | path | anchor domain name | Yes | string |
| registerAnchorRequest | body | Anchor regsitration request | Yes | [registerAnchorRequest](#registeranchorrequest) |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | achor domain is registered | [operatingAccount](#operatingaccount) |
| 400 | Missing or invalid parameters in the request | [error](#error) |
| 404 | The registration failed due to error communicating with ledger | [error](#error) |
| 409 | The registration failed due to conflict with ledger | [error](#error) |

### Models
---

### account  

Account

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | string | The address that was created on the ledger. | Yes |
| name | string | A name to identify this account. | No |

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

### operatingAccount  

Account with the token

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| account | [account](#account) |  | No |
| token | string | auth token | No |

### registerAnchorRequest  

register anchor id on WW

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | string | The stellar address of the issuing account of the anchor, should have ibm admin account as signatory | Yes |