// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package environment

var (
	//********AWS service store
	// anchor service will have specific env variables

	//used for mock tests
	ENV_KEY_PARTICIPANT_REGISTRY_SERVICE_CLIENT = "PARTICIPANT_REGISTRY_CLIENT"

	ENV_KEY_KAFKA_BROKER_URL           = "KAFKA_BROKER_URL"
	ENV_KEY_DYNAMO_DB_REGION           = "DYNAMO_DB_REGION"
	ENV_KEY_DYNAMO_DB_TABLE_NAME       = "DYNAMO_DB_TABLE_NAME"
	ENV_KEY_KAFKA_ENABLE_SSL           = "KAFKA_ENABLE_SSL"
	ENV_KEY_KAFKA_CA_LOCATION          = "KAFKA_CA_LOCATION"
	ENV_KEY_KAFKA_CERTIFICATE_LOCATION = "KAFKA_CERTIFICATE_LOCATION"
	ENV_KEY_KAFKA_KEY_LOCATION         = "KAFKA_KEY_LOCATION"
	ENV_KEY_KAFKA_KEY_PASSWORD         = "KAFKA_KEY_PASSWORD"
	ENV_KEY_WW_BIC                     = "WW_BIC"
	ENV_KEY_WW_ID                      = "WW_ID"

	//***** AWS secrets Participant specific environments
	//2
	//Strong-hold Callback Credentials, should be participant specific
	ENV_KEY_ANCHOR_SH_PASS     = "ANCHOR_SH_PASS"
	ENV_KEY_ANCHOR_SH_SEC      = "ANCHOR_SH_SEC"
	ENV_KEY_ANCHOR_SH_CRED     = "ANCHOR_SH_CRED"
	ENV_KEY_ANCHOR_SH_VENEU    = "ANCHOR_SH_VENUE"
	ENV_KEY_ANCHOR_SH_ROOT_URL = "ANCHOR_SH_ROOT_URL"
)
