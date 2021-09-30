// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package global_environment

var (
	//Should come from deployment parameters, comes from kubernets/k8
	// participant specific needed only for per participant service
	ENV_KEY_HOME_DOMAIN_NAME      = "HOME_DOMAIN_NAME"
	ENV_KEY_AWS_REGION            = "AWS_REGION"
	ENV_KEY_AWS_ACCESS_KEY_ID     = "AWS_ACCESS_KEY_ID"
	ENV_KEY_AWS_SECRET_ACCESS_KEY = "AWS_SECRET_ACCESS_KEY"

	// every service will have specific env variables
	ENV_KEY_SERVICE_NAME         = "SERVICE_NAME"
	ENV_KEY_WRITE_TIMEOUT        = "WRITE_TIMEOUT"
	ENV_KEY_READ_TIMEOUT         = "READ_TIMEOUT"
	ENV_KEY_IDLE_TIMEOUT         = "IDLE_TIMEOUT"
	ENV_KEY_WAIT_UNLOCK_DURATION = "WAIT_UNLOCK_DURATION"
	ENV_KEY_GAS_ACCOUNT_ATTEMPTS = "GAS_ACCOUNT_ATTEMPTS"

	//********AWS global
	// 1 same for for all services
	ENV_KEY_ENABLE_JWT              = "ENABLE_JWT"
	ENV_KEY_ORIGIN_ALLOWED          = "ORIGIN_ALLOWED"
	ENV_KEY_HORIZON_CLIENT_URL      = "HORIZON_CLIENT_URL"
	ENV_KEY_STELLAR_NETWORK         = "STELLAR_NETWORK"
	ENV_KEY_FIREBASE_CREDENTIALS    = "FIREBASE_CREDENTIALS" //json file will come from k8
	ENV_KEY_FIREBASE_DB_URL         = "FIREBASE_DB_URL"
	ENV_KEY_VAULT_BASE_URL          = "VAULT_BASE_URL"
	ENV_KEY_VAULT_CERT              = "VAULT_CERT"
	ENV_KEY_VAULT_CERT_PRIVATE_KEY  = "VAULT_CERT_PRIVATE_KEY"
	ENV_KEY_SECRET_STORAGE_LOCATION = "SECRET_STORAGE_LOCATION"
	ENV_KEY_SERVICE_VERSION         = "SERVICE_VERSION"
	ENV_KEY_IBM_TOKEN_DOMAIN_ID     = "IBM_TOKEN_DOMAIN_ID"

	// ENV_KEY_ENVIRONMENT_VERSION = dev | qa | st | tn | prod
	// IMPORTANT: Used to write to firebase so that the write data structure can
	// be observed to distinguish the environment from which logs
	// and transactions were genererated from
	ENV_KEY_ENVIRONMENT_VERSION = "ENV_VERSION"

	ENV_KEY_STRONGHOLD_ANCHOR_ID = "STRONGHOLD_ANCHOR_ID"

	//2 service URLs
	//singleton apis
	ENV_KEY_PARTICIPANT_REGISTRY_URL = "PARTICIPANT_REGISTRY_URL"
	ENV_KEY_QUOTE_SVC_URL            = "QUOTE_SVC_URL"
	ENV_KEY_WL_SVC_URL               = "WL_SVC_URL"
	ENV_KEY_GAS_SVC_URL              = "GAS_SVC_URL"
	ENV_KEY_ADMIN_SVC_URL            = "ADMIN_SVC_URL"
	ENV_KEY_ANCHOR_SVC_URL           = "ANCHOR_SVC_URL"
	ENV_KEY_FEE_SVC_URL              = "FEE_SVC_URL"
	//Participant apis
	//these are taken as a template and utility function resolves url based on participant ids
	ENV_KEY_RDO_SVC_URL             = "RDO_SVC_URL"        //should be something like https://{participant_id}-rdo-service.worldwire-io
	ENV_KEY_SEND_SVC_URL            = "SEND_SVC_URL"       //should be something like https://{participant_id}-send-service.worldwire-io
	ENV_KEY_API_SVC_URL             = "API_SVC_URL"        //should be something like https://{participant_id}-api-service.worldwire-io
	ENV_KEY_PAYMENT_SVC_URL         = "PAYMENT_SVC_URL"    //should be something like https://{participant_id}-payment-listener.worldwire-io
	ENV_KEY_CRYPTO_SVC_INTERNAL_URL = "CRYPTO_SVC_INT_URL" //should be something like https://{participant_id}-crypto-service.worldwire-io

	//3 Constants, the value of these env vars are these strings already
	ENV_KEY_NODE_ADDRESS      = "NODE_ADDRESS"
	ENV_KEY_NODE_SEED         = "NODE_SEED"
	ENV_KEY_PUBLIC_KEY_LABEL  = "PUBLIC_LABEL"
	ENV_KEY_PRIVATE_KEY_LABEL = "PRIVATE_LABEL"

	//Service Sepcific Variables
	ENV_KEY_SERVICE_LOG_FILE      = "SERVICE_LOG_FILE"
	ENV_KEY_SERVICE_PORT          = "SERVICE_PORT"
	ENV_KEY_SERVICE_INTERNAL_PORT = "SERVICE_INTERNAL_PORT"

	//Stored and Taken from docker image
	ENV_KEY_SERVICE_ERROR_CODES_FILE = "SERVICE_ERROR_CODES_FILE"

	//********AWS participant store/service store
	// every participant specific service will this key have it based on participant
	// For global services it will come from AWS service specific store
	// ENV_KEY_JWT_SECRET_KEY  = "JWT_SECRET_KEY" // deprecating due to usage of
	ENV_KEY_WW_JWT_PEPPER_OBJ = "WW_JWT_PEPPER_OBJ"

	// Should be participant specific
	//Needed for tests with account keys stored in nodeconfig only
	ENV_KEY_NODE_CONFIG = "STELLAR_NODE_CONFIG"

	//Not used
	ENV_KEY_STELLAR_CONFIG_REMAP  = "STELLAR_CONFIG_REMAP"
	ENV_KEY_STELLAR_CONFIG_SCHEME = "STELLAR_CONFIG_SCHEME"
)
