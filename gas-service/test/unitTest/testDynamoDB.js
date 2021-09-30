// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
let should = require('should')
const AWS = require('../../method/AWS')
const timeOutSec = 300000

const delay = (interval) => {
    return new Promise((resolve) => {
        setTimeout(resolve, interval);
    });
};
const logFile = require('../../method/logTestResult')
const writeReport = new logFile('../../Report/txt/UnitTestResult.txt')
writeReport.Report()

describe('DynamoDB', function () {
    this.timeout(timeOutSec);
    describe('successful case ', function () {
        it('create test Table ', async function () {
            var params = {
                TableName: "testTable",
                KeySchema: [
                    { AttributeName: "test", KeyType: "HASH" }
                ],
                AttributeDefinitions: [
                    { AttributeName: "test", AttributeType: "S" }
                ],
                ProvisionedThroughput: {
                    ReadCapacityUnits: 5,
                    WriteCapacityUnits: 5
                }
            };
            let result = await AWS.createTable(params)
            should(result.TableDescription.TableStatus).be.exactly("CREATING",JSON.stringify(result))
        });
        it('create item to test Table ', async function () {
            await delay(8000);
            let item = {
                test: 't1'
            }
            let result = await AWS.createItem("testTable", item)
            should(result).be.exactly(item, JSON.stringify(result, 2, null))
        });
        it('get all item from test Table ', async function () {
            await delay(5000);
            let item = { test: 't1' } 
            let result = await AWS.getAllDatas("testTable")

            should.deepEqual(result[0],item, JSON.stringify(result, 2, null))
        });
        it('delete test Table ', async function () {
            await delay(5000);
            let result =await  AWS.deleteTable("testTable")
            should(result.TableDescription.TableStatus).be.exactly("DELETING", JSON.stringify(result))
        });
    })

    describe('failing case - create table', function () {
        it('create test Table2 ', async function () {
            var params = {
                TableName: "testTable2",
                KeySchema: [
                    { AttributeName: "test", KeyType: "HASH" }
                ],
                AttributeDefinitions: [
                    { AttributeName: "test", AttributeType: "S" }
                ],
                ProvisionedThroughput: {
                    ReadCapacityUnits: 5,
                    WriteCapacityUnits: 5
                }
            };
            let result = await AWS.createTable(params)
            should(result.TableDescription.TableStatus).be.exactly("CREATING",JSON.stringify(result))
        });
        it('create test Table2 , should return get: Table already exists ', async function () {
            await delay(10000);
            var params = {
                TableName: "testTable2",
                KeySchema: [
                    { AttributeName: "test", KeyType: "HASH" }
                ],
                AttributeDefinitions: [
                    { AttributeName: "test", AttributeType: "S" }
                ],
                ProvisionedThroughput: {
                    ReadCapacityUnits: 5,
                    WriteCapacityUnits: 5
                }
            };
            let result = await AWS.createTable(params)
            should(result.message).be.exactly("Table already exists: testTable2",JSON.stringify(result))
        });
        it('delete test Table ', async function () {
            await delay(5000);
            let result =await  AWS.deleteTable("testTable2")
            should(result.TableDescription.TableStatus).be.exactly("DELETING", JSON.stringify(result))
        });
    })


    describe('failing case - read item', function () {
        it('read not exist table ', async function () {
            let result = await AWS.getAllDatas("notexist")
            should(result.message).be.exactly("Requested resource not found", JSON.stringify(result))
        });
    })
    describe('failing case - using wrong crt and key', function () {
        it('read not exist table ', async function () {
            let result = await AWS.getAllDatas("notexist")
            should(result.message).be.exactly("Requested resource not found", JSON.stringify(result))
        });
    })

});

