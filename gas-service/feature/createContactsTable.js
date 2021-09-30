// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const AWS = require('../method/AWS')
const environment = require('../environment/env')
let emailgroupTablename = process.env[environment.ENV_KEY_DYNAMODB_CONTACTS_TABLE_NAME]

/**
 * 
 * @param {String} tablename 
 */
async function createTable(emailgroupTablename) {
    
    var params = {
        TableName: emailgroupTablename,
        KeySchema: [
            { AttributeName: "topicName", KeyType: "HASH" },
            { AttributeName: "email", KeyType: "RANGE" }
        ],
        AttributeDefinitions: [
            { AttributeName: "topicName", AttributeType: "S" },
            { AttributeName: "email", AttributeType: "S" }
        ],
        ProvisionedThroughput: {
            ReadCapacityUnits: 5,
            WriteCapacityUnits: 5
        }
    };
    AWS.createTable(params)

}

createTable(emailgroupTablename)