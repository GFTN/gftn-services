// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 

const AWS = require("../method/AWS")
const LOGGER = require("../method/logger")
const log = new LOGGER('Monitor BL')

/**
 * check lock per monitoringTime
 * if lock timestamp + expireTime > now
 * unlock
 */

module.exports = 
/**
 * 
 * @param {Integer} timeout 
 * @param {Integer} expireTime 
 * @param {Array} lockAccounts 
 * @param {Array} unlockAccounts 
 * @param {String} accountstablename 
 */
function monitor(timeout, expireTime, lockAccounts, unlockAccounts, accountstablename) {


    setTimeout(function () {

        
        // log.data(lockAccounts,'')
        /**
         * if lockArray has something then do chck
         */
        if (lockAccounts.length > 0) {
            
            lockAccounts.forEach((account, accountInunlockArrayIndex) => {
                

                let now = new Date().getTime()
                /**
                 * if expire then unlock
                 */
                
                if (account.lockTimestamp + parseInt(expireTime)*1000 < now) {
                    /**
                     * updateDB
                     */
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
                     AWS.updateItem(params)

                    /**
                     * update memory
                     */
                    lockAccounts.splice(accountInunlockArrayIndex, 1);
                    unlockAccounts.push(account.pkey)
                    
                    log.info('Updated','Successful');
                }
            });
        }
        
        /**
         * monitor
         */
        log.data("Lock Accounts",JSON.stringify(lockAccounts))
        log.data("Unlock Accounts",JSON.stringify(unlockAccounts))

        monitor(timeout, expireTime, lockAccounts, unlockAccounts, accountstablename)
    }, timeout);
}

