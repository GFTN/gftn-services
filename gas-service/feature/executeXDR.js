// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const stellar = require("../method/stellar.js")
const AWS = require("../method/AWS")
const LOGGER = require('../method/logger')
const log = new LOGGER('Execute Tx')
    /**
     * 1. decode signedXDRin , get source account 
     * 2. check whether account is in $lockAccounts , if account is not in $lockAccounts , return can not execute 
     * 3. unlock account
     * 4. execute signedXDR
     * 5. return result ,$lockAccounts,$unlockAccounts
     */
module.exports =
    /**
     * 
     * @param {String} signedXDRin 
     * @param {String} lockAccounts 
     * @param {String} unlockAccounts 
     * @param {String} accountstablename 
     */

    async function(signedXDRin, lockAccounts, unlockAccounts, accountstablename) {
        try {
            let account = {}

            /**
             * 1. decode signedXDRin , get source account 
             */
            let transaction = await stellar.newTransaction(signedXDRin);
            account.pkey = transaction.source

            /**
             * 2. check whether account is in $lockAccounts
             */
            // var accountInunlockArrayIndex = lockAccounts.indexOf(account.pkey)
            let accountInunlockArrayIndex = lockAccounts.map(function(item) { return item.pkey; }).indexOf(account.pkey);
            /**
             * if account is not in $lockAccounts , return tx fail ,  
             */
            log.logger('Account Public Key', account.pkey)
            log.logger('Lock Accounts', JSON.stringify(lockAccounts))

            if (accountInunlockArrayIndex < 0) {
                let rejMsg = {
                    title: "Source Account Expire",
                    failure_reason: "source account not availible"
                }
                return rejMsg
            } else {

                let params = {
                    TableName: accountstablename,
                    Key: {
                        "pkey": account.pkey
                    },
                    UpdateExpression: "set accountStatus = :st, lockTimestamp=:ts",
                    ExpressionAttributeValues: {
                        ":st": true,
                        ":ts": null
                    },
                    ReturnValues: "UPDATED_NEW"
                };

                let data = await AWS.updateItem(params)
                log.logger('Updated', JSON.stringify(data))

                /**
                 * 3. execute signedXDR
                 */
                let result = await stellar.submitTransaction(transaction)

                /**
                 *  unlock account
                 */
                lockAccounts.splice(accountInunlockArrayIndex, 1);
                unlockAccounts.push(account.pkey)

                /**
                 * 5. return result
                 */
                return result
            }
        } catch (err) {
            return err
        }


    }