// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const AWS = require('../method/AWS')
const environment = require('../environment/env')
let accountstablename = process.env[environment.ENV_KEY_DYNAMODB_ACCOUNTS_TABLE_NAME]

/**
 * 
 * @param {String} accountstablename 
 */
async function createTable(accountstablename) {
    
    var params = {
        TableName: accountstablename,
        KeySchema: [
            { AttributeName: "pkey", KeyType: "HASH" }
        ],
        AttributeDefinitions: [
            { AttributeName: "pkey", AttributeType: "S" }
        ],
        ProvisionedThroughput: {
            ReadCapacityUnits: 5,
            WriteCapacityUnits: 5
        }
    };
    AWS.createTable(params)


}

createTable(accountstablename)