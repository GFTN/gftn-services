// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
let should = require('should')
const environment = require('../../environment/env')
const request = require('request')
const timeOutSec = 300000

describe.stress(1, 'High value of calls to gas service - /lockAccount ', function () {
    this.timeout(timeOutSec);
    it.stress(1, '/lockAccount , should get 200 or 500', function (done) {
        let options = {
            contentType: 'application/json',
            method: 'GET',
            body: '',
            url: environment.ENV_KEY_GAS_SERVICE_URL+':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/lockaccount'
        }

        request(options, function (err, res, body) {
            if (err) done(err)
            else {
                let obj = JSON.parse(body)
                if(res.statusCode==200){
                    should.exist(obj.pkey, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    should.exist(obj.sequenceNumber, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    should(res.statusCode).be.exactly(200, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    done()}
                else{
                    should.exist(obj.failure_reason + "\n request data : " + JSON.stringify(options))
                    should(res.statusCode).be.exactly(500, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    done()
                }
            }
        });
      })
  })