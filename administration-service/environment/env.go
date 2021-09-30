// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package environment

var (
	//********AWS service store
	// admin service will have specific env variables
	ENV_KEY_DB_USER            = "DB_USER"
	ENV_KEY_DB_PWD             = "DB_PWD"
	ENV_KEY_ADMIN_DB_NAME      = "ADMIN_DB_NAME"
	ENV_KEY_PAYMENTS_DB_TABLE  = "PAYMENTS_DB_TABLE"
	ENV_KEY_BLOCKLIST_DB_TABLE = "BLOCKLIST_DB_TABLE"
	ENV_KEY_DB_TIMEOUT         = "DB_TIMEOUT"
	ENV_KEY_MONGO_ID           = "MONGO_ID"

	//used for local testing only
	ENV_KEY_PARTICIPANT_REGISTRY_SERVICE_CLIENT = "PARTICIPANT_REGISTRY_SERVICE_CLIENT"
)
