// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
let should = require('should')
const environment = require('../../environment/env')
const request = require('request')
const timeOutSec = 300000

describe('/signXDRAndExecuteXDR', function () {
    this.timeout(timeOutSec);
    let pkey;
    let seq;
    let oneSignedXDR;
    describe('successful case ', function () {
        it('using /unlockAccount to unlock at least one accout', (done) => {

            let options = {
                contentType: 'application/json',
                method: 'POST',
                body: {
                    "pkey": "GARAU4JHLLMGMNKSNOFGGJCF3LV5GJHSHK5BRNNVQQ7TGNTEW4ABB4ZC"
                },
                json: true,
                headers: { "Connection": "keep-alive" },
                url: environment.ENV_KEY_GAS_SERVICE_URL + ':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/unlockAccount'
            }
            request(options, function (err, res, body) {
                if (err) {
                    should(typeof err).be.exactly('undefined', 'throw err : '+JSON.stringify(err))
                    done(err)
                }
                else {
                    should(res.statusCode).be.exactly(200, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    done()
                }
            });
        })
        it('/lockaccount get public key and sequence number ', function (done) {

            let options = {
                contentType: 'application/json',
                method: 'GET',
                body: '',
                json: false,
                url: environment.ENV_KEY_GAS_SERVICE_URL + ':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/lockaccount'
            }

            request(options, function (err, res, body) {
                if (err) {
                    should(typeof err).be.exactly('undefined', 'throw err : '+JSON.stringify(err))
                    done(err)
                }
                else {
                    let obj = JSON.parse(body)
                    should.exist(obj.pkey, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    pkey = obj.pkey
                    should.exist(obj.sequenceNumber, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    seq = parseInt(obj.sequenceNumber) - 1
                    should(res.statusCode).be.exactly(200, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    done()
                }
            });
        });
        it('/getMockTx, using Mock test to make a tx ', function (done) {

            let options = {
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
                if (err) {
                    should(typeof err).be.exactly('undefined', 'throw err : '+JSON.stringify(err))
                    done(err)
                }
                else {
                    should.exist(body.oneSignedXDR, res.body + "\n request data : " + JSON.stringify(options))
                    oneSignedXDR = body.oneSignedXDR
                    should(res.statusCode).be.exactly(200, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    done()
                }
            });

        });
        it('/signXDRAndExecuteXDR, should successful , when using the correct sequence number and account ', function (done) {

            let options = {
                contentType: 'application/json',
                method: 'POST',
                body: {
                    oneSignedXDR: oneSignedXDR
                },
                json: true,
                url: environment.ENV_KEY_GAS_SERVICE_URL + ':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/signXDRAndExecuteXDR'
            }

            request(options, function (err, res, body) {
                if (err) {
                    should(typeof err).be.exactly('undefined', 'throw err : '+JSON.stringify(err))
                    done(err)
                }
                else {
                    should(res.statusCode).be.exactly(200, JSON.stringify(body) + "\n request data : " + JSON.stringify(options))
                    done()
                }
            });

        });
    })
    describe('failing case : 400 - signing fail / account lock', function () {
        it('/lockaccount get public key and sequence number ', function (done) {

            let options = {
                contentType: 'application/json',
                method: 'GET',
                body: '',
                json: false,
                url: environment.ENV_KEY_GAS_SERVICE_URL + ':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/lockaccount'
            }

            request(options, function (err, res, body) {
                if (err) {
                    should(typeof err).be.exactly('undefined', 'throw err : '+JSON.stringify(err))
                    done(err)
                }
                else {
                    let obj = JSON.parse(body)
                    should.exist(obj.pkey, res.body + "\n request data : " + JSON.stringify(options))
                    pkey = obj.pkey
                    should.exist(obj.sequenceNumber, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    seq = parseInt(obj.sequenceNumber) - 1
                    should(res.statusCode).be.exactly(200, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    done()
                }
            });
        });
        it('/getMockTx, using Mock test to make a tx ', function (done) {

            let options = {
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
                if (err) {
                    should(typeof err).be.exactly('undefined', 'throw err : '+JSON.stringify(err))
                    done(err)
                }
                else {
                    should.exist(body.oneSignedXDR, res.body + "\n request data : " + JSON.stringify(options))
                    oneSignedXDR = body.oneSignedXDR
                    should(res.statusCode).be.exactly(200, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    done()
                }
            });

        });
        it('/signXDRAndExecuteXDR, should successful , when using the correct sequence number and account ', function (done) {

            let options = {
                contentType: 'application/json',
                method: 'POST',
                body: {
                    oneSignedXDR: oneSignedXDR
                },
                json: true,
                url: environment.ENV_KEY_GAS_SERVICE_URL + ':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/signXDRAndExecuteXDR'
            }

            request(options, function (err, res, body) {
                if (err) {
                    should(typeof err).be.exactly('undefined', 'throw err : '+JSON.stringify(err))
                    done(err)
                }
                else {

                    should(res.statusCode).be.exactly(200, JSON.stringify(body) + "\n request data : " + JSON.stringify(options))
                    done()
                }
            });

        });

        it('/signXDRAndExecuteXDR, send the same XDR again , should get status : 400 fail cause of the account should be lock again', function (done) {

            let options = {
                contentType: 'application/json',
                method: 'POST',
                body: {
                    oneSignedXDR: oneSignedXDR
                },
                json: true,
                url: environment.ENV_KEY_GAS_SERVICE_URL + ':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/signXDRAndExecuteXDR'
            }


            request(options, function (err, res, body) {
                if (err) {
                    should(typeof err).be.exactly('undefined', 'throw err : '+JSON.stringify(err))
                    done(err)
                }
                else {

                    should(body.title).be.exactly("Source Account Expire", JSON.stringify(body))
                    should(body.failure_reason).be.exactly("source account not availible", JSON.stringify(body))
                    should(res.statusCode).be.exactly(400, JSON.stringify(body) + "\n request data : " + JSON.stringify(options))
                    done()
                }
            });

        });

    })

    describe('failing case : 400 - Signer not Exist', function () {
        it('/getMockTx, using Mock test and using wrong public key to make a tx ', function (done) {

            let options = {
                contentType: 'application/json',
                method: 'GET',
                body: {
                    sequenceNumber: "7036144273326080",
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
                        pkey: "GBVQDOCFTX3N54G2IWOBISDKRMDOZPXWPBRDNMN2MN3ORD3AKND37R46",
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
                
                
                if (err) {
                    should(typeof err).be.exactly('undefined', 'throw err : '+JSON.stringify(err))
                    done(err)
                }
                else {
                    should.exist(body.oneSignedXDR, res.body + "\n request data : " + JSON.stringify(options))
                    oneSignedXDR = body.oneSignedXDR
                    should(res.statusCode).be.exactly(200, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    done()
                }
            });

        });
        it('/signXDRAndExecuteXDR,  should get status : 400 fail cause of the wrong public key', function (done) {

            let options = {
                contentType: 'application/json',
                method: 'POST',
                body: {
                    oneSignedXDR: oneSignedXDR
                },
                json: true,
                url: environment.ENV_KEY_GAS_SERVICE_URL + ':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/signXDRAndExecuteXDR'
            }

            request(options, function (err, res, body) {
                if (err) {
                    should(typeof err).be.exactly('undefined', 'throw err : '+JSON.stringify(err))
                    done(err)
                }
                else {

                    should(body.title).be.exactly("Source Account Not IBM Account", JSON.stringify(body))
                    should(res.statusCode).be.exactly(400, JSON.stringify(body) + "\n request data : " + JSON.stringify(options))
                    done()
                }
            });

        });

    })
    describe('failing case : 403 - Transaction fail', function () {
        it('/lockaccount get public key and sequence number ', function (done) {

            let options = {
                contentType: 'application/json',
                method: 'GET',
                body: '',
                json: false,
                url: environment.ENV_KEY_GAS_SERVICE_URL + ':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/lockaccount'
            }

            request(options, function (err, res, body) {
                if (err) {
                    should(typeof err).be.exactly('undefined', 'throw err : '+JSON.stringify(err))
                    done(err)
                }
                else {
                    let obj = JSON.parse(body)
                    should.exist(obj.pkey, res.body + "\n request data : " + JSON.stringify(options))
                    pkey = obj.pkey
                    should.exist(obj.sequenceNumber, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    seq = parseInt(obj.sequenceNumber) - 1
                    should(res.statusCode).be.exactly(200, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    done()
                }
            });
        });
        it('/getMockTx, using Mock test and using wrong sequence number to make a tx ', function (done) {

            let options = {
                contentType: 'application/json',
                method: 'GET',
                body: {
                    sequenceNumber: 9999999999999999,
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
                if (err) {
                    should(typeof err).be.exactly('undefined', 'throw err : '+JSON.stringify(err))
                    done(err)
                }
                else {
                    should.exist(body.oneSignedXDR, res.body + "\n request data : " + JSON.stringify(options))
                    oneSignedXDR = body.oneSignedXDR
                    should(res.statusCode).be.exactly(200, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    done()
                }
            });

        });
        it('/signXDRAndExecuteXDR, should get status : 403 fail cause of the wrong sequence number', function (done) {

            let options = {
                contentType: 'application/json',
                method: 'POST',
                body: {
                    oneSignedXDR: oneSignedXDR
                },
                json: true,
                url: environment.ENV_KEY_GAS_SERVICE_URL + ':' + environment.ENV_KEY_GAS_SERVICE_PORT + '/signXDRAndExecuteXDR'
            }

            request(options, function (err, res, body) {
                if (err) {
                    should(typeof err).be.exactly('undefined', 'throw err : '+JSON.stringify(err))
                    done(err)
                }
                else {
                    should(body.title).be.exactly("Transaction Failed", JSON.stringify(body))
                    should(res.statusCode).be.exactly(403, JSON.stringify(body) + "\n request data : " + JSON.stringify(options))
                    done()
                }
            });

        });

    })

});
