// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// let should = require('should')
// const environment = require('../../environment/env')
// const request = require('request')
// const timeOutSec = 300000

//     describe('/createAccounts', function () {
//         this.timeout(timeOutSec);
//         describe('failing case - All the accounts already exist', function () {
//             it('should return status: 400, All accounts exist ', function (done) {
            
//                 let options = {
//                     contentType: 'application/json',
//                     method: 'POST',
//                     body: [
//                         {
//                             "key": {
//                                 "Object": "IBM_TOKEN_ACCOUNT_ADDRESS_1"
//                             },
//                             "seed": {
//                                 "Object": "IBM_TOKEN_ACCOUNT_SEED_1"
//                             },
//                             "accountStatus": true,
//                             "topicName": "Group4"
                            
//                         }
//                     ],
//                     json: true,
//                     url: environment.ENV_KEY_GAS_SERVICE_URL+':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/createAccounts'
//                 }
    
//                 request(options, function (err, res, body) {
//                     if (err) {done(err)}
//                     else {
//                         for (let index = 0; index < body.length; index++) {
//                             should.exist(body[index].pkey, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
//                             should.exist(body[index].secret, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
//                             should.exist(body[index].accountStatus, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
//                             should.exist(body[index].groupName, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
//                         }
                        
//                         should(res.statusCode).be.exactly(400, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
//                         done()
//                     }
//                 });
//             });
//         })
        
//         describe('failing case - Using wrong tag to get from Vault', function () {
//             it('should return status: 500 and fail Info.', function (done) {
            
//                 let options = {
//                     contentType: 'application/json',
//                     method: 'POST',
//                     body: [
//                         {
//                             "key": {
//                                 "Object": "Wrong Tag"
//                             },
//                             "seed": {
//                                 "Object": "Wrong Tag"
//                             },
//                             "accountStatus": true,
//                             "groupName": "Group4"
                            
                            
//                         }
//                     ],
//                     json: true,
//                     url: environment.ENV_KEY_GAS_SERVICE_URL+':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/createAccounts'
//                 }
    
//                 request(options, function (err, res, body) {
//                     if (err) done(err)
//                     else {
//                         should.exist(body.body.ErrorMsg, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
//                         should(res.statusCode).be.exactly(500, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
//                         done()
//                     }
//                 });
//             });
//         })
//     });
