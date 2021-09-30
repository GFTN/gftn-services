// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const LOGGER = require('../method/logger')
const environment = require('../environment/env')
const SM = require('../utility/aws/javascript/build/awsSecret')

module.exports = {
    /**
     * 
     * @param {String} accountInfos 
     */
    getDataFromAws: async function(accountInfos) {

        let accountPromise = []
        for (let index = 0; index < accountInfos.length; index++) {

            let keyTag = accountInfos[index].key.Object
            let seedTag = accountInfos[index].seed.Object

            let title = {
                environment: process.env[environment.ENV_KEY_ENVIRONMENT_VERSION],
                domain: process.env[environment.ENV_KEY_HOME_DOMAIN_NAME],
                service: process.env[environment.ENV_KEY_SERVICE_NAME],
                variable: "gas-account"
            }
            console.log(title)
            let res = await SM.getSecret(title)
            console.log(res)
            let obj = JSON.parse(res)
            let keys = Object.keys(obj)
            let account = {
                pkey: obj[keyTag],
                secret: obj[seedTag],
                accountStatus: accountInfos[index].accountStatus,
                topicName: accountInfos[index].topicName
            }
            accountPromise.push(account)
        }
        Promise.all(accountPromise)
        if (accountPromise.length == 0) {
            throw ({
                statusCode: 400,
                Message: {
                    ErrorMsg: 'All accounts exist'
                }
            })
        }
        return accountPromise
    }
}