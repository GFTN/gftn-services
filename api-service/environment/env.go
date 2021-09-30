// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package environment

var (
	//********AWS service store
	// api service will have specific env variables
	//Constant for all API-services only
	//1
	ENV_KEY_ACCOUNT_INITIAL_FUND    = "ACCOUNT_INITIAL_FUND"
	ENV_KEY_TRANSACTION_BATCH_LIMIT = "TRANSACTION_BATCH_LIMIT"

	//env variables for mock test
	ENV_KEY_CRYPTO_SERVICE_CLIENT               = "CRYPTO_SERVICE_CLIENT"
	ENV_KEY_PARTICIPANT_REGISTRY_SERVICE_CLIENT = "PARTICIPANT_REGISTRY_CLIENT"
)
