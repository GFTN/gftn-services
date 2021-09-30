// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const LOGGER = require('../method/logger')
const log = new LOGGER('unlockAct')
const AWS = require('../method/AWS')
module.exports = 
async function(account,accountstablename,lockAccounts,unlockAccounts){
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
    if (item == null) {
        throw ({
            statusCode: 400,
            Message: {
                title: "Account Not IBM Account",
                failure_reason: "can not find account from DynamoDB"
            }
        })
    }
    else {
        let accountInunlockArrayIndex = unlockAccounts.map(function (item) { return item.pkey; }).indexOf(account.pkey);
        let accountInlockArrayIndex = lockAccounts.map(function (item) { return item.pkey; }).indexOf(account.pkey);
        
        if (accountInunlockArrayIndex!= -1){
            unlockAccounts.splice(accountInlockArrayIndex, 1);
        }if (accountInlockArrayIndex!= -1){
            lockAccounts.splice(accountInunlockArrayIndex, 1);
            unlockAccounts.push(account.pkey)
            let params = {
                TableName: accountstablename,
                Key: {
                    "pkey": account.pkey
                },
                UpdateExpression: "set accountStatus = :st, lockTimestamp=:ts",
                ExpressionAttributeValues: {
                    ":st": account.accountStatus,
                    ":ts": null
                },
                ReturnValues: "UPDATED_NEW"
            };
            
             result = await AWS.updateItem(params)
        }

        
        return unlockAccounts
    }

}