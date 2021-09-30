// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
let should = require('should')
const environment = require('../../environment/env')
const request = require('request')
const timeOutSec = 300000

describe.stress(1, 'High value of calls to gas service - /signXDRAndExecuteXDR ',  function () {
    this.timeout(timeOutSec);
    let oneSignedXDR;
    let pkey;
    let seq;
    it.stress(10, 'if /lockaccount return 200 should go through /getMockTx then /signXDRAndExecuteXDR ', function (done) {
        let options = {
            contentType: 'application/json',
            method: 'GET',
            body: '',
            url: environment.ENV_KEY_GAS_SERVICE_URL + ':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/lockaccount'
        }

        request(options, function (err, res, body) {
            if (err) done(err)
            else {
                let obj = JSON.parse(body)
                if (res.statusCode == 200) {
                    should.exist(obj.pkey, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    should.exist(obj.sequenceNumber, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    should(res.statusCode).be.exactly(200, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

                    pkey = obj.pkey
                    seq = parseInt(obj.sequenceNumber) - 1
                    // isGetSequenceNum = true
                     options = {
                        contentType: 'application/json',
                        method: 'GET',
                        body: {
                            sequenceNumber: seq,
                            from: {
                                pkey: "GA6ANHY4BOJZFHWO3M56PNNP5OD6SLCF4EQ4KUOFVU6ISAEQ2G5JKRIZ",
                                secret: "SCSQZEMCE4JW23JW22YI7OZLZB4SRMOGHPOVK3XBC2AJHA4RP7ST3TKD",
                                asset: {
                                    code: "",
                                    issuer: "GARQZQKXTOTWP22UFUEHSYU7BEJIPP7TK2EM27P55HA3GH5E6SSPIDFE",
                                    amount: "1"
                                }
                            },
                            to: {
                                pkey: pkey,
                                asset: {
                                    code: "",
                                    issuer: "SCIOZJEUGAO7PYYYHPBU7O7VKG7FZCA2RS3AR66L4UNTE4DW64PDKL4W",
                                    amount: "2"
                                }
                            }
                        },
                        json: true,
                        url: environment.ENV_KEY_GAS_SERVICE_URL + ':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/getMockTx'
                    }
                    request(options, function (err, res, body) {
                        if (err) done(err)
                        else {
                            should.exist(body.oneSignedXDR, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                            oneSignedXDR = body.oneSignedXDR
                            should(res.statusCode).be.exactly(200, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                            
                             options = {
                                contentType: 'application/json',
                                method: 'POST',
                                body: {
                                    oneSignedXDR: oneSignedXDR
                                },
                                json: true,
                                url: environment.ENV_KEY_GAS_SERVICE_URL + ':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/signXDRAndExecuteXDR'
                            }
                
                            request(options, function (err, res, body) {
                                if (err) (done(err))
                                else {
                                    should(res.statusCode).be.exactly(200, JSON.stringify(body) + "\n request data : " + JSON.stringify(options))
                                    done()
                                }
                            });
                            
                        }
                    });

                }
                if (res.statusCode == 500) {
                    should.exist(obj.failure_reason + "\n request data : " + JSON.stringify(options))
                    should(res.statusCode).be.exactly(500, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    isGetSequenceNum = false
                    done()
                }
            }
        });
    })
})