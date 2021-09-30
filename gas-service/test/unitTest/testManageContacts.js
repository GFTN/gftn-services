// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
let should = require('should')
const environment = require('../../environment/env')
const request = require('request')
const AWS = require('../../method/AWS')
const timeOutSec = 300000

    describe('Manage Topics', function () {
        this.timeout(timeOutSec);
        describe('successful case', function () {
            let topicArn
            it('/createTopics should return status: 200 ',  function (done) {
            
                let options = {
                    contentType: 'application/json',
                    method: 'POST',
                    body: [
                        {
                            "TopicName":"TestGroup",
                            "DisplayName":"TestGroup"
                        }
                    
                    ],
                    json: true,
                    url: environment.ENV_KEY_GAS_SERVICE_URL+':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/createTopics'
                }
    
                request(options, function (err, res, body) {
                    if (err) {done(err)}
                    else {
                        should(res.statusCode).be.exactly(200, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                        done()
                    }
                });
            });
            it('check dynamoDB should return topicArn ', async function () {
                try{
                    const topicstablename = environment.ENV_KEY_DYNAMODB_GROUPS_TABLE_NAME
                    topicArn = await AWS.getTopicArn(topicstablename, "TestGroup")
                    should.notEqual(topicArn,null, topicArn)
                }
                catch(err){
                    should.equal(err,null, err)
                }
            });
            it('/createContacts should return status: 200  ',  function (done) {
            
                let options = {
                    contentType: 'application/json',
                    method: 'POST',
                    body: [
                        {
                            "topicName": "TestGroup",
                            "email":"test@ibm.com",
                            "phoneNumber":"+6599999999"
                        }
                    ],
                    json: true,
                    url: environment.ENV_KEY_GAS_SERVICE_URL+':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/createContacts'
                }
    
                request(options, function (err, res, body) {
                    if (err) {done(err)}
                    else {
                        should(res.statusCode).be.exactly(200, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                        done()
                    }
                });
            });
            it('/deleteContact should return status: 200  ',  function (done) {
            
                let options = {
                    contentType: 'application/json',
                    method: 'DELETE',
                    body: {
                        "topicName": "TestGroup",
                        "email": "test@ibm.com"
                },
                    json: true,
                    url: environment.ENV_KEY_GAS_SERVICE_URL+':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/deleteContact'
                }
    
                request(options, function (err, res, body) {
                    if (err) {done(err)}
                    else {
                        should(res.statusCode).be.exactly(200, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                        done()
                    }
                });
            });

            it('/deleteTopic should return status: 200  ', async function () {
            
                let options = {
                    contentType: 'application/json',
                    method: 'DELETE',
                    body: {
                        "TopicArn": topicArn,
                        "TopicName": "TestGroup"
                },
                    json: true,
                    url: environment.ENV_KEY_GAS_SERVICE_URL+':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/deleteTopic'
                }
    
                request(options, function (err, res, body) {
                    if (err) {done(err)}
                    else {
                        should(res.statusCode).be.exactly(200, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    }
                });
            });
        })
    });
