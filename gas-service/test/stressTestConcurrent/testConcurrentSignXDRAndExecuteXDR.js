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
    it('/signXDRAndExecuteXDR , should get 200 ,400, 403', function (done) {
        function statusCallback(error, result, latency) {

            let obj = JSON.parse(result.body)
            if (result.statusCode==200){
                should(obj.title).be.exactly("Transaction successful", JSON.stringify(result))
            }
            if(result.statusCode==400){
                should(obj.title).be.exactly("Source Account Expire", JSON.stringify(result.body))
                should(obj.failure_reason).be.exactly("source account not availible", JSON.stringify(result))
            }
            else{
                should(result.statusCoderesult.statusCode).be.exactly(403, JSON.stringify(result) + "\n request data : " + JSON.stringify(options))
                should(obj.title).be.exactly("Transaction Failed", JSON.stringify(result))
            }
            console.log('Current latency %j, result %j, error %j', latency, result, error);
            console.log('----');
            console.log('Request elapsed milliseconds: ', result.requestElapsed);
            console.log('Request index: ', result.requestIndex);
            console.log('Request loadtest() instance index: ', result.instanceIndex);
        }
        
        const options = {
            url: environment.ENV_KEY_GAS_SERVICE_URL + ':' + environment.ENV_KEY_GAS_SERVICE_PORT  ,
            concurrent: 5,
            method: 'POST',
            body: {
                "oneSignedXDR": "AAAAACIKcSda2GY1UmuKYyRF2uvTJPI6uhi1tYQ/MzZktwAQAAAAyAANhhcAAACCAAAAAAAAAAAAAAACAAAAAQAAAAA8Bp8cC5OSns7bO+e1r+uH6SxF4SHFUcWtPIkAkNG6lQAAAAEAAAAAIgpxJ1rYZjVSa4pjJEXa69Mk8jq6GLW1hD8zNmS3ABAAAAAAAAAAAACYloAAAAABAAAAACIKcSda2GY1UmuKYyRF2uvTJPI6uhi1tYQ/MzZktwAQAAAAAQAAAAA8Bp8cC5OSns7bO+e1r+uH6SxF4SHFUcWtPIkAkNG6lQAAAAAAAAAAATEtAAAAAAAAAAABkNG6lQAAAECDIIDFai3QSssFQUPqoojO7BX/tCHB7TDKUALpHYPm5hy5/r/hswrkKGasUcxfku06JoYHqJqijmWXvDB76cAP"
            },
            requestsPerSecond: 5,
            maxSeconds: 30,
            statusCallback: statusCallback,
            requestGenerator: (params, options, client, callback) => {
                options.headers['Content-Type'] = 'application/json';
                options.body = '';
                options.path = '/signXDRAndExecuteXDR';
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

