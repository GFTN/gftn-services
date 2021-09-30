// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const stellar = require("../method/stellar.js")
const AWS = require('../method/AWS')
const LOGGER = require('../method/logger')
const log = new LOGGER('Sign Transaction')

/**
 * 1. decode signedXDRin , get pkey 
 * x. (not do)check whether signedXDRin, sequence number is fit dynamoDB info
 * 3. use pkey to get secret
 * 4. signed by the secret 
 * 3. return txeB64 , using $account
 * 4. 
 */
module.exports =
    /**
     * 
     * @param {String} signedXDRin 
     * @param {String} accountstablename 
     */
    async function(signedXDRin, accountstablename) {
        let account = {}
        console.log(signedXDRin);
        /**
         * 1. decode signedXDRin , get pkey 
         */

        let decode = await stellar.decode(signedXDRin)
        log.logger('Existing Signature ', decode._attributes.signatures.length)


        let transaction = await stellar.newTransaction(signedXDRin);
        account.pkey = transaction.source
        log.logger('IBM Account', account.pkey)


        /**
         * 3. use pkey to get secret
         */

        let params = {
            TableName: accountstablename,
            KeyConditionExpression: "#pk = :pkey",
            ExpressionAttributeNames: {
                "#pk": "pkey"
            },
            ExpressionAttributeValues: {
                ":pkey": account.pkey
            }
        };

        let item = await AWS.queryData(params)
        log.logger('IBM Account', JSON.stringify(item))

        if (item == null) {
            throw ({
                statusCode: 400,
                Message: {
                    title: "Source Account Not IBM Account",
                    failure_reason: "can not find account from DynamoDB"
                }
            })
        }
        account.secret = item.secret
        log.logger('Using Secrect', account.secret)


        /**
         * 4. signed by the secret 
         */
        transaction = await stellar.signTx(transaction, account.secret)

        let signResult = transaction.toEnvelope().toXDR('base64')
        log.logger('Signed Result', signResult)


        /**
         * log how many signatures
         */
        decode = await stellar.decode(signResult)
        log.logger('Existing Signatures ', decode._attributes.signatures.length)


        /**
         * 3. return signResult , using $account
         */
        return signResult
    }