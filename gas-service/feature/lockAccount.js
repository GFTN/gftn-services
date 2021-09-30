// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 

const stellar = require("../method/stellar")
const AWS = require("../method/AWS")
const LOGGER = require('../method/logger')
const log = new LOGGER('Lock Accoounts')
/**
    * 1.pop from unlockAccounts ($account)
    * 2.get timestamp
    * 3.update to dynamoDB ,
    * [ update account from table where pkey = account.pkey and status = unlock (status,timestamp) ]
    * 3.add timestamp to $account 
    * 4.push $account to lockArray (lockAccounts.push($account))
    * 5.get sequence number $account.seqNum
    * 5.return $account.seqNum,$account.pkey , $lockAccounts , $unlockAccounts
    */


module.exports = 

/**
 * 
 * @param {Array} lockAccounts 
 * @param {Array} unlockAccounts 
 * @param {String} accountstablename 
 */
async function (lockAccounts, unlockAccounts, accountstablename) {
    try {
        let account={}
        /**
        * 1.pop from unlockAccounts ($account)
        */
       
        if (unlockAccounts.length == 0) {
            return null
        }
        account.pkey = await unlockAccounts.shift()

        /**
         * 2.get timestamp
         */
        let timestamp = new Date().getTime()
        /**
         * 3.update to dynamoDB ,
         * [ update account from table where pkey = account.pkey and status = unlock (status,timestamp) ]
         */
        account.accountStatus=false
        let params = {
            TableName: accountstablename,
            Key: {
                "pkey": account.pkey
            },
            UpdateExpression: "set accountStatus = :st, lockTimestamp=:ts",
            ExpressionAttributeValues: {
                ":st": account.accountStatus,
                ":ts": timestamp
            },
            ReturnValues: "UPDATED_NEW"
        };
        let data = await AWS.updateItem(params)
        log.logger('UPDATED','')
        console.log(data)
        
        

        /** 
         * 4.push $account to lockArray (lockAccounts.push($account))
         */
        account.accountStatus = false
        let pushdata = {
            pkey:account.pkey,
            lockTimestamp:timestamp
        }
        lockAccounts.push(pushdata)


        /**
         * 5.get sequence number $account.seqNum
         * 5.return $account.seqNum,$account.pkey , $lockAccounts , $unlockAccounts
         */

        account.sequenceNumber = await stellar.getAccountSequenceNumber(account.pkey)
        return account


    } catch (err) {
        return err
    }

};