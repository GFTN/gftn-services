// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
'use strict';

let should = require('should')
const timeOutSec = 3000000
const environment = require('../../environment/env')
const loadtest = require('loadtest/lib/loadtest')
const logFile = require('../../method/logTestResult')

const writeReport = new logFile('../../Report/txt/ConcurrentTestResult.txt')
describe( 'High value of calls to gas service - Current', function () {
    writeReport.Report()
    this.timeout(timeOutSec);
    it('/lockAccount , should get 200 or 500', function (done) {
        function statusCallback(error, result, latency) {
            if (result.statusCode==200){
                let obj = JSON.parse(result.body)
                    should.exist(obj.pkey, JSON.stringify(result) + "\n request data : " + JSON.stringify(options))
                    should.exist(obj.sequenceNumber, JSON.stringify(result) + "\n request data : " + JSON.stringify(options))
                    
            }
            else{
                should(result.statusCode).be.exactly(500, "\n response data : " + JSON.stringify(result))
            }
            console.log('Current latency %j, result %j, error %j', latency, result, error);
            console.log('----');
            console.log('Request elapsed milliseconds: ', result.requestElapsed);
            console.log('Request index: ', result.requestIndex);
            console.log('Request loadtest() instance index: ', result.instanceIndex);
        }
        
        const options = {
            url: environment.ENV_KEY_GAS_SERVICE_URL + ':' + environment.ENV_KEY_GAS_SERVICE_PORT ,
            concurrent: 5,
            method: 'GET',
            body: '',
            requestsPerSecond: 5,
            maxSeconds: 30,
            statusCallback: statusCallback,
            requestGenerator: (params, options, client, callback) => {
                options.headers['Content-Type'] = 'application/json';
                options.body = '';
                options.path = '/lockAccount';
                const request = client(options, callback);
                return request;
            }
        };

        loadtest.loadTest(options, (error, results) => {
            if (error) {
                return console.error('Got an error: %s', error);
            }
            console.log(results);
            console.log('Tests run successfully');
            done()
        });
    })
})

