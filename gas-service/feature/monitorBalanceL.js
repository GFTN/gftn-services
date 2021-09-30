// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const environment = require('../environment/env')
const stellar = require('../method/stellar')
const LOGGER = require('../method/logger')
const log = new LOGGER('Monitor LT')
const notifyEmail = require('../method/notifyEmal')
const notifySMS = require('../method/notifySMS')
const AWS = require('../method/AWS')

const accountstablename = process.env[environment.ENV_KEY_DYNAMODB_ACCOUNTS_TABLE_NAME]
const sendEmail = process.env[environment.ENV_KEY_GAS_SERVICE_EMAIL_NOTIFICATION]
const sendSMS = process.env[environment.ENV_KEY_GAS_SERVICE_SMS_NOTIFICATION]
    /**
     * 0. if there has no accounts in the low threshold list stop
     * 1. read account balance from stellar
     * 2. check whetheer the account balancs is lower the NOTIFY_BALANCE
     * 3. if $balance>MONITOR_LOW_THRESHOLD_BALANCE move to another list (high threshold list)
     * 4. if lower than the NOTIFY_BALANCE , send SMS,EMAIL,SLACK
     */
module.exports =
    async function monitor(highThresholdAccounts, lowThresholdAccounts, lowThresholdBalance, lowThresholdTimeout) {

        /**
         * check low threshold length
         * if there has no accounts in the low threshold list stop
         * else check each account
         */

        log.logger("Low threshold monitoring ", 'Accounts Number: ' + lowThresholdAccounts.length);
        if (lowThresholdAccounts.length > 0) {

            setTimeout(function() {

                lowThresholdAccounts.forEach(async(account, index) => {

                    let balancesSet = await stellar.getBalance(account)
                    let nativeBalanceIndex = balancesSet.map(function(item) { return item.asset_type; }).indexOf('native');
                    log.info(account[0] + " (" + parseFloat(balancesSet[nativeBalanceIndex].balance) + ' < ' + lowThresholdBalance + ') ', (parseFloat(balancesSet[nativeBalanceIndex].balance) < lowThresholdBalance))
                    let pkey = account[0]
                    if (parseFloat(balancesSet[nativeBalanceIndex].balance) < lowThresholdBalance) {

                        /**
                         * send email (to account email group)
                         */
                        if (sendEmail == 'true' || sendSMS == 'true') {
                            let params = {
                                TableName: accountstablename,
                                KeyConditionExpression: "#pk = :pkey",
                                ExpressionAttributeNames: {
                                    "#pk": "pkey"
                                },
                                ExpressionAttributeValues: {
                                    ":pkey": pkey
                                }
                            };
                            let item = await AWS.queryData(params)
                            let topicName = item.topicName

                            if (sendEmail == 'true') {
                                notifyEmail(topicName, ' Balance Notification', 'Please top up to ' + account, '<strong>  Low balance notification  Please top up to ' + account + '</strong>')
                            }
                            if (sendSMS == 'true') {
                                notifySMS(account[0], topicName)
                            }

                        }

                    } else {
                        let popAccountIndex = lowThresholdAccounts.indexOf(account[0]);
                        let highBalanceAccount = lowThresholdAccounts.splice(popAccountIndex, 1)
                        highThresholdAccounts.push(highBalanceAccount)
                    }
                })
                monitor(highThresholdAccounts, lowThresholdAccounts, lowThresholdBalance, lowThresholdTimeout)
            }, lowThresholdTimeout)
        } else {
            log.logger("Ligh threshold monitoring ", 'MONITOR LOW THRESHOLD STOP ');
        }



    }