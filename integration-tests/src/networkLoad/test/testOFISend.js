// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
'use strict';

let should = require('should')
const timeOutSec = 3000000
const environment = require('../environment/env')
const loadtest = require('loadtest/lib/loadtest')
const logFile = require('../method/logTestResult')

const writeReport = new logFile('../Report/txt/networkLoadTestResult.txt')
describe( 'High value of calls to gas service - Current', function () {
    writeReport.Report()
    this.timeout(timeOutSec);
    it('/denf , should get 200 or 500', function (done) {
        function statusCallback(error, result, latency) {
                should(result.statusCode).be.exactly(200, "\n response data : " + JSON.stringify(result))
            
            console.log('Current latency %j, result %j, error %j', latency, result, error);
            console.log('----');
            console.log('Request elapsed milliseconds: ', result.requestElapsed);
            console.log('Request index: ', result.requestIndex);
            console.log('Request loadtest() instance index: ', result.instanceIndex);
        }
        
        const options = {
            url: "http://"+environment.ENV_KEY_OFI_URL,
            concurrent: 5,
            method: 'POST',
            body: {
                "message_type": "pacs.008.001.07",
                "message": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPERvY3VtZW50IHhtbG5zPSJ1cm46aXNvOnN0ZDppc286MjAwMjI6dGVjaDp4c2Q6cGFjcy4wMDguMDAxLjA3Ij4KCTxGSVRvRklDc3RtckNkdFRyZj4KCQk8R3JwSGRyPgoJCQk8TXNnSWQ+QkJCQi8xNTA5MjgtQ0NUL0pQWS8xMjM8L01zZ0lkPgoJCQk8Q3JlRHRUbT4yMDE1LTA5LTI4VDE2OjAwOjAwPC9DcmVEdFRtPgoJCQk8TmJPZlR4cz4xPC9OYk9mVHhzPgoJCQk8U3R0bG1JbmY+CgkJCQk8U3R0bG1NdGQ+Q09WRTwvU3R0bG1NdGQ+CgkJCQk8SW5zdGdSbWJyc21udEFndD4KCQkJCQk8RmluSW5zdG5JZD4KCQkJCQkJPEJJQ0ZJPkNDQ0NKUEpUPC9CSUNGST4KCQkJCQk8L0Zpbkluc3RuSWQ+CgkJCQk8L0luc3RnUm1icnNtbnRBZ3Q+CgkJCQk8SW5zdGRSbWJyc21udEFndD4KCQkJCQk8RmluSW5zdG5JZD4KCQkJCQkJPEJJQ0ZJPkFBQUFKUEpUPC9CSUNGST4KCQkJCQk8L0Zpbkluc3RuSWQ+CgkJCQk8L0luc3RkUm1icnNtbnRBZ3Q+CgkJCTwvU3R0bG1JbmY+CgkJCTxJbnN0Z0FndD4KCQkJCTxGaW5JbnN0bklkPgoJCQkJCTxCSUNGST5CQkJCVVMzMzwvQklDRkk+CgkJCQk8L0Zpbkluc3RuSWQ+CgkJCTwvSW5zdGdBZ3Q+CgkJCTxJbnN0ZEFndD4KCQkJCTxGaW5JbnN0bklkPgoJCQkJCTxCSUNGST5BQUFBR0IyTDwvQklDRkk+CgkJCQk8L0Zpbkluc3RuSWQ+CgkJCTwvSW5zdGRBZ3Q+CgkJPC9HcnBIZHI+CgkJPENkdFRyZlR4SW5mPgoJCQk8UG10SWQ+CgkJCQk8SW5zdHJJZD5CQkJCLzE1MDkyOC1DQ1QvSlBZLzEyMy8xPC9JbnN0cklkPgoJCQkJPEVuZFRvRW5kSWQ+QUJDLzQ1NjIvMjAxNS0wOS0wODwvRW5kVG9FbmRJZD4KCQkJCTxUeElkPkJCQkIvMTUwOTI4LUNDVC9KUFkvMTIzLzE8L1R4SWQ+CgkJCTwvUG10SWQ+CgkJCTxQbXRUcEluZj4KCQkJCTxJbnN0clBydHk+Tk9STTwvSW5zdHJQcnR5PgoJCQk8L1BtdFRwSW5mPgoJCQk8SW50ckJrU3R0bG1BbXQgQ2N5PSJKUFkiPjEwPC9JbnRyQmtTdHRsbUFtdD4KCQkJPEludHJCa1N0dGxtRHQ+MjAxNS0wOS0yOTwvSW50ckJrU3R0bG1EdD4KCQkJPENocmdCcj5TSEFSPC9DaHJnQnI+CgkJCTxEYnRyPgoJCQkJPE5tPkFCQyBDb3Jwb3JhdGlvbjwvTm0+CgkJCQk8UHN0bEFkcj4KCQkJCQk8U3RydE5tPlRpbWVzIFNxdWFyZTwvU3RydE5tPgoJCQkJCTxCbGRnTmI+NzwvQmxkZ05iPgoJCQkJCTxQc3RDZD5OWSAxMDAzNjwvUHN0Q2Q+CgkJCQkJPFR3bk5tPk5ldyBZb3JrPC9Ud25ObT4KCQkJCQk8Q3RyeT5VUzwvQ3RyeT4KCQkJCTwvUHN0bEFkcj4KCQkJPC9EYnRyPgoJCQk8RGJ0ckFjY3Q+CgkJCQk8SWQ+CgkJCQkJPE90aHI+CgkJCQkJCTxJZD4wMDEyNTU3NDk5OTwvSWQ+CgkJCQkJPC9PdGhyPgoJCQkJPC9JZD4KCQkJPC9EYnRyQWNjdD4KCQkJPERidHJBZ3Q+CgkJCQk8RmluSW5zdG5JZD4KCQkJCQk8QklDRkk+QkJCQlVTMzM8L0JJQ0ZJPgoJCQkJPC9GaW5JbnN0bklkPgoJCQk8L0RidHJBZ3Q+CgkJCTxDZHRyQWd0PgoJCQkJPEZpbkluc3RuSWQ+CgkJCQkJPEJJQ0ZJPkFBQUFHQjJMPC9CSUNGST4KCQkJCTwvRmluSW5zdG5JZD4KCQkJPC9DZHRyQWd0PgoJCQk8Q2R0cj4KCQkJCTxObT5ERUYgRWxlY3Ryb25pY3M8L05tPgoJCQkJPFBzdGxBZHI+CgkJCQkJPFN0cnRObT5NYXJrIExhbmU8L1N0cnRObT4KCQkJCQk8QmxkZ05iPjU1PC9CbGRnTmI+CgkJCQkJPFBzdENkPkVDM1I3TkU8L1BzdENkPgoJCQkJCTxUd25ObT5Mb25kb248L1R3bk5tPgoJCQkJCTxDdHJ5PkdCPC9DdHJ5PgoJCQkJCTxBZHJMaW5lPkNvcm4gRXhjaGFuZ2UgNXRoIEZsb29yPC9BZHJMaW5lPgoJCQkJPC9Qc3RsQWRyPgoJCQk8L0NkdHI+CgkJCTxDZHRyQWNjdD4KCQkJCTxJZD4KCQkJCQk8T3Rocj4KCQkJCQkJPElkPjIzNjgzNzA3OTk0MjE1PC9JZD4KCQkJCQk8L090aHI+CgkJCQk8L0lkPgoJCQk8L0NkdHJBY2N0PgoJCQk8UHVycD4KCQkJCTxDZD5HRERTPC9DZD4KCQkJPC9QdXJwPgoJCQk8Um10SW5mPgoJCQkJPFN0cmQ+CgkJCQkJPFJmcmREb2NJbmY+CgkJCQkJCTxUcD4KCQkJCQkJCTxDZE9yUHJ0cnk+CgkJCQkJCQkJPENkPkNJTlY8L0NkPgoJCQkJCQkJPC9DZE9yUHJ0cnk+CgkJCQkJCTwvVHA+CgkJCQkJCTxOYj40NTYyPC9OYj4KCQkJCQkJPFJsdGREdD4yMDE1LTA5LTA4PC9SbHRkRHQ+CgkJCQkJPC9SZnJkRG9jSW5mPgoJCQkJPC9TdHJkPgoJCQk8L1JtdEluZj4KCQk8L0NkdFRyZlR4SW5mPgoJPC9GSVRvRklDc3RtckNkdFRyZj4KPC9Eb2N1bWVudD4=",
                "transaction_details": {
                "amount_beneficiary": 1,
                "amount_settlement": 3,
                "asset_code_beneficiary": "EURDO",
                "asset_settlement": {
                  "asset_code": "EURDO",
                  "asset_type": "DO",
                  "issuer_id": "ie.one.payments.worldwire.io"
                },
                "fee_admin": {
                  "cost": 1,
                  "cost_asset": {
                    "asset_code": "HKDDO",
                    "asset_type": "DO",
                    "issuer_id": "hk.one.payments.worldwire.io"
                  }
                },
                "fee_creditor": {
                  "cost": 1,
                  "cost_asset": {
                    "asset_code": "EURDO",
                    "asset_type": "DO",
                    "issuer_id": "ie.one.payments.worldwire.io"
                  }
                },
                "fee_debtor": {
                  "cost": 1,
                  "cost_asset": {
                    "asset_code": "HKDDO",
                    "asset_type": "DO",
                    "issuer_id": "hk.one.payments.worldwire.io"
                  }
                },
                "ofi_id": "hk.one.payments.worldwire.io",
                "rfi_id": "ie.one.payments.worldwire.io",
                "settlement_method": "DO"
              }
            },
            requestsPerSecond: 5,
            maxSeconds: 30,
            statusCallback: statusCallback,
            requestGenerator: (params, options, client, callback) => {
                options.headers['Content-Type'] = 'application/json';
                options.body = '';
                options.path = '/v1/client/transactions/send';
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

