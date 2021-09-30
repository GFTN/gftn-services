// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const Env = require('../utility/aws/javascript/build/env')

function env_init() {
    Env.CheckVariable()
    Env.InitEnv()
}

env_init()
module.exports = {
    ENV_KEY_DYNAMODB_ACCOUNTS_TABLE_NAME: "DYNAMODB_ACCOUNTS_TABLE_NAME",
    ENV_KEY_GAS_SERVICE_EMAIL_NOTIFICATION: "GAS_SERVICE_EMAIL_NOTIFICATION",
    ENV_KEY_GAS_SERVICE_SMS_NOTIFICATION: "GAS_SERVICE_SMS_NOTIFICATION",
    ENV_KEY_DYNAMODB_CONTACTS_TABLE_NAME: "DYNAMODB_CONTACTS_TABLE_NAME",
    ENV_KEY_DYNAMODB_GROUPS_TABLE_NAME: "DYNAMODB_GROUPS_TABLE_NAME",
    ENV_KEY_DYNAMODB_REGION: "DYNAMODB_REGION",
    ENV_KEY_DYNAMODB_ENDPOINT_URL: "DYNAMODB_ENDPOINT_URL",
    ENV_KEY_DYNAMODB_ACCESSKEYID: "AWS_ACCESS_KEY_ID",
    ENV_KEY_DYNAMODB_SECRECT_ACCESS_KEY: "AWS_SECRET_ACCESS_KEY",
    ENV_KEY_VAULT_APPID: "VAULT_APPID",
    ENV_KEY_VAULT_SAFE: "VAULT_SAFE",
    ENV_KEY_VAULT_FOLDER: "VAULT_FOLDER",
    ENV_KEY_VAULT_URL: "VAULT_URL",
    ENV_KEY_VAULT_CERT_PATH: "VAULT_CERT_PATH",
    ENV_KEY_VAULT_KEY_PATH: "VAULT_KEY_PATH",
    ENV_KEY_GAS_SERVICE_PORT: "GAS_SERVICE_PORT",
    ENV_KEY_GAS_SERVICE_MONITOR_LOCKACCOUNT_FEQ: "GAS_SERVICE_MONITOR_LOCKACCOUNT_FEQ",
    ENV_KEY_GAS_SERVICE_EXPIRE_TIME: "GAS_SERVICE_EXPIRE_TIME",
    ENV_KEY_HIGH_THRESHOLD_BALANCE: "HIGH_THRESHOLD_BALANCE",
    ENV_KEY_HIGH_THRESHOLD_TIMEOUT: "HIGH_THRESHOLD_TIMEOUT",
    ENV_KEY_LOW_THRESHOLD_BALANCE: "LOW_THRESHOLD_BALANCE",
    ENV_KEY_LOW_THRESHOLD_TIMEOUT: "LOW_THRESHOLD_TIMEOUT",
    ENV_KEY_EMAIL_SENDER: "EMAIL_SENDER",
    ENV_KEY_GAS_SERVICE_URL: "GAS_SERVICE_URL",
    ENV_KEY_SENDGRID_API_KEY: "SENDGRID_API_KEY",
    ENV_KEY_SENDGRIID_ENDPINT: "SENDGRIID_ENDPINT",
    ENV_KEY_SERVICE_LOG_FILE: "SERVICE_LOG_FILE",
    ENV_KEY_STELLAR_NETWORK: "STELLAR_NETWORK",
    ENV_KEY_HORIZON_CLIENT_URL: "HORIZON_CLIENT_URL",
    ENV_KEY_SERVICE_NAME: "SERVICE_NAME",
    ENV_KEY_ENVIRONMENT_VERSION: "ENV_VERSION",
    ENV_KEY_HOME_DOMAIN_NAME: "HOME_DOMAIN_NAME"
}