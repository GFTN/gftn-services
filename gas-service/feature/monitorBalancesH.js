// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const stellar = require('../method/stellar')
const monitorBalanceL = require('../feature/monitorBalanceL')
const LOGGER = require('../method/logger')
const log = new LOGGER('Monitor HT')


module.exports =

    /**
     * 0. read all IBM accounts from dynamoDB
     * 1. get the balance from stellar
     * 2. check each accounts whether balance is ($balance<MONITOR_HIGH_THRESHOLD_BALANCE)
     * 3. if ($balance<MONITOR_HIGH_THRESHOLD_BALANCE) and ($lowThresholdAccounts.length<1) { create thread (call monitorBalanceL)}
     * 4. else ($balance<MONITOR_HIGH_THRESHOLD_BALANCE) { push accounts to lowThresholdAccounts }
     */
    async function monitor(highThresholdAccounts, highThresholdBalance, highThresholdTimeout,
        lowThresholdAccounts, lowThresholdBalance, lowThresholdTimeout) {


        setTimeout(function () {
            log.logger("Monitoring account's up (High threshold) , Account Number - ",  highThresholdAccounts.length);            
            /**
             * monitor all account's balance(for loop)
             */
            
            highThresholdAccounts.forEach(async (account, index) => {
                log.info(index, account)
                /**
                 * get account balance from stellar
                 * get native account index ( for know index to check belance)
                 * get account balance
                 */
                let balancesSet = await stellar.getBalance(account)
                let nativeBalanceIndex = balancesSet.map(function (item) { return item.asset_type; }).indexOf('native');


                /**
                 * check native account balance
                 * if balance is lower then High threshold balnce
                 *  {if low threshold monitoring was not exist then create monitoring
                 *   else just delete account from queue , and add account to low threshold queue}
                 */

                if (parseFloat(balancesSet[nativeBalanceIndex].balance) < highThresholdBalance) {
                    
                    if (lowThresholdAccounts.length == 0) {
                    let popAccountIndex = highThresholdAccounts.indexOf(account);
                    let lowBalanceAccount = highThresholdAccounts.splice(popAccountIndex, 1)
                    if (lowBalanceAccount != '') {
                        lowThresholdAccounts.push(lowBalanceAccount)
                        log.info('highThresholdAccounts array', highThresholdAccounts)

                    }
                        log.info("Create lowThreashold monitor", 'NOW CALL MONITORL');
                        monitorBalanceL(highThresholdAccounts, lowThresholdAccounts, lowThresholdBalance, lowThresholdTimeout)
                    }
                    else {
                    let popAccountIndex = highThresholdAccounts.indexOf(account);
                    let lowBalanceAccount = highThresholdAccounts.splice(popAccountIndex, 1)
                        if (lowBalanceAccount != '') {
                            lowThresholdAccounts.push(lowBalanceAccount)
                            log.info('highThresholdAccounts array', highThresholdAccounts)
    
                        }
                        log.info("Add accounts to lowThreshold accounts", account);
                    }

                }

            });


            monitor(highThresholdAccounts, highThresholdBalance, highThresholdTimeout,
                lowThresholdAccounts, lowThresholdBalance, lowThresholdTimeout)
        }, highThresholdTimeout);
    }