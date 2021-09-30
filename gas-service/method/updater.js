// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 


module.exports = {
    updateUnlockQueue:  function (unlockAccounts,newAccounts) {
        return new Promise(function(res,rej){

            for (let index = 0; index < newAccounts.length; index++) {
                unlockAccounts.push(newAccounts[index].pkey)
            }
            Promise.all(unlockAccounts)
            res(unlockAccounts)
        })
    },
    addHighThresholdAccountsQueue:  function (addHighThresholdAccounts,newAccounts) {
        return new Promise(function(res,rej){

            for (let index = 0; index < newAccounts.length; index++) {
                addHighThresholdAccounts.push(newAccounts[index].pkey)
            }
            Promise.all(addHighThresholdAccounts)
            res(addHighThresholdAccounts)
        })
    }
}
