// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const AWS = require('../method/AWS')
const environment = require('../environment/env')
let emailgroupTablename = process.env[environment.ENV_KEY_DYNAMODB_GROUPS_TABLE_NAME]

/**
 * 
 * @param {String} tablename 
 */
async function createTable(emailgroupTablename) {
    
    var params = {
        TableName: emailgroupTablename,
        KeySchema: [
            { AttributeName: "TopicName", KeyType: "HASH" },
            { AttributeName: "TopicArn", KeyType: "RANGE" }
        ],
        AttributeDefinitions: [
            { AttributeName: "TopicName", AttributeType: "S" },
            { AttributeName: "TopicArn", AttributeType: "S" }
        ],
        ProvisionedThroughput: {
            ReadCapacityUnits: 5,
            WriteCapacityUnits: 5
        }
    };
    AWS.createTable(params)

    // await AWS.scanData(emailgroupTablename)


}

createTable(emailgroupTablename)