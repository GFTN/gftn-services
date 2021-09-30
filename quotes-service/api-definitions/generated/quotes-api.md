# World Wire API Service Callbacks
API endpoints clients are expected to implement in order to receive notifications of transactions

## Version: 1.0.0

### /quote

#### POST
##### Summary:

Create a quote

##### Description:

Provides a quote in response to requests for a given target asset in exchange for a source asset, using the source asset as its price.


##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| QuoteRequestNotification | body | Quote request to RFI detailing quoteID target asset, source asset, and amount desired to exchange.  | Yes | [quoteRequestNotification](#quoterequestnotification) |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | Successfully receive a valid quote request.  |
| 404 | Unsuccessfully receive a valid quote request. |

### Models


#### quoteRequestNotification

Quote Request

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| quote_id | string | Unique id for this quote as set by the quote service | Yes |
| quote_request | object | Quote Request | Yes |