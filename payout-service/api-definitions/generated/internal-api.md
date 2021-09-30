# Payout location API
API endpoints for querying the details of payout location locations

## Version: 1.0.0

### /

#### POST
##### Summary:

Called when participant wants to create a new payout location

##### Description:

If there doesn't have an existing payout location. participants can use this endpoint to create the payout location


##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| payout location | body | The payout location that needs to be added. | Yes | [payoutLocation](#payoutlocation) |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | Payout location created |
| 400 | Payout location could not be created due to payout location already exists or malformed payload |
| 404 | Payout location could not be created |
| 500 | Internal server error |

#### DELETE
##### Summary:

delete a payout location

##### Description:

delete the payout location with given id


##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | query | the id of the payout location you want to delete | Yes | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | Payout location deleted |
| 400 | missing or invalid parameters in the request |
| 404 | The payout location could not be found |
| 500 | Internal server error |

#### PATCH
##### Summary:

Called when participant wants to update an existing payout location

##### Description:

If there is an existing payout location. participants can use this endpoint to update the payout location


##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| payout location | body | The payout location that needs to be updated. | Yes | [payoutLocationUpdateRequest](#payoutlocationupdaterequest) |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | Payout location updated |
| 400 | Failed updating the payout location |
| 404 | Cannot find the payout location |
| 500 | Internal server error |

### Models


#### address

Address

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| building_number | string | The building number or identifier. | Yes |
| city | string | Name of the city or town. | Yes |
| country | string | Country code of the location. | Yes |
| postal_code | string | Postal code for the location. | Yes |
| state | string | Name of the state. | Yes |
| street | string | The street name. | Yes |

#### coordinate

Geographic coordinates for a location. Based on https://schema.org/geo

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| lat | number | The latitude of the geo coordinates | Yes |
| long | number | The longitude of the geo coordinates | Yes |

#### geo

Geographic coordinates for a location. Based on https://schema.org/geo

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| coordinates | [ [coordinate](#coordinate) ] | The geo coordinates | Yes |
| type | string | The type of location. Options include "point" if the location is a single pickup location, or "area" if it's a region.
 | Yes |

#### payoutLocation

Details of each payout location. Based on https://schema.org/LocalBusiness

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | [address](#address) |  | No |
| category | [payoutLocationCategory](#payoutlocationcategory) |  | Yes |
| currencies_accepted | [ string ] | The currency accepted. | Yes |
| geo | [geo](#geo) |  | Yes |
| id | string | The unique identifier of the location. | No |
| image | string | An image of the item. This can be a URL or a fully described ImageObject. | Yes |
| member_of | [ string ] | The financial institute that this location belongs to. | Yes |
| name | string | The name of the location. | Yes |
| opening_hours | [ [payoutLocationOpeningHours](#payoutlocationopeninghours) ] | The opening hours of the location. | Yes |
| payout_child | [ string ] | The collection of identifiers for locations which belong to the location - these can include areas, and points.
 | Yes |
| payout_parent | [ string ] | The collection of identifiers for the parents of the locations - it can be only areas.
 | Yes |
| routing_number | string | Optional routing information, also known as BIC (bank id code). | No |
| telephone | string | The phone number of the location. | Yes |
| type | string | The type of location. Options include: "Bank", "Non-Bank Financial Institution", "Mobile Network Operator", or "Other".
 | Yes |
| url | string | The URL of the location. | Yes |

#### payoutLocationCategory

Details of each payout location offer category. Based on https://schema.org/hasOfferCatalog

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| name | string | name of the category | Yes |
| options | [ [payoutLocationOption](#payoutlocationoption) ] | offer list of the category | Yes |

#### payoutLocationOpeningHours

The opening hours of each payout location. Based on https://schema.org/OpeningHoursSpecification

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| closes | string | The closing hour of the payout location on the given day(s) of the week | Yes |
| day_of_week | [ string ] | The day of the week for which these opening hours are valid | Yes |
| opens | string | The opening hour of the payout location on the given day(s) of the week | Yes |

#### payoutLocationOption

Details of each payout location offer. Based on https://schema.org/hasOfferCatalog

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| description | string | name of the service | Yes |
| terms | string | service detail | Yes |

#### payoutLocationUpdateRequest

List of updated payout location attributes. Based on https://schema.org/LocalBusiness

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string | The identifier of the payout location | Yes |
| updated_payload | [payoutLocation](#payoutlocation) |  | Yes |