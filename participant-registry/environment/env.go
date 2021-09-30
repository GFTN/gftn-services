// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package environment

var (
	//********AWS service store
	// participant registry service will have specific env variables
	ENV_KEY_PR_DB_NAME            = "PR_DB_NAME"
	ENV_KEY_PARTICIPANTS_DB_TABLE = "PARTICIPANTS_DB_TABLE"
	ENV_KEY_DB_USER               = "DB_USER"
	ENV_KEY_DB_PWD                = "DB_PWD"
	ENV_KEY_DB_TIMEOUT            = "DB_TIMEOUT"
	ENV_KEY_MONGO_ID              = "MONGO_ID"

	//used for local testing only
	ENV_KEY_IS_UNIT_TEST = "IS_UNIT_TEST"
)
