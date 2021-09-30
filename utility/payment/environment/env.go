// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package environment

var (
	ENV_KEY_SERVICE_FILE = "SERVICE_FILE"

	//********AWS participant store
	// participant specific env variables
	ENV_KEY_KAFKA_BROKER_URL           = "KAFKA_BROKER_URL"
	ENV_KEY_DYNAMO_DB_REGION           = "DYNAMO_DB_REGION"
	ENV_KEY_DYNAMO_DB_TABLE_NAME       = "DYNAMO_DB_TABLE_NAME"
	ENV_KEY_KAFKA_ENABLE_SSL           = "KAFKA_ENABLE_SSL"
	ENV_KEY_KAFKA_CA_LOCATION          = "KAFKA_CA_LOCATION"
	ENV_KEY_KAFKA_CERTIFICATE_LOCATION = "KAFKA_CERTIFICATE_LOCATION"
	ENV_KEY_KAFKA_KEY_LOCATION         = "KAFKA_KEY_LOCATION"
	ENV_KEY_KAFKA_KEY_PASSWORD         = "KAFKA_KEY_PASSWORD"
	// BIC code of world wire
	ENV_KEY_WW_BIC = "WW_BIC"
	// participant id of world wire
	ENV_KEY_WW_ID              = "WW_ID"
	ENV_KEY_PARTICIPANT_BIC    = "PARTICIPANT_BIC"
	ENV_KEY_ANCHOR_SH_ASSET_ID = "ANCHOR_SH_ASSET_ID"
)
