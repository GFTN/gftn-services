// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const {
    Given,
    When,
    Then
} = require('cucumber')
const request = require('request')
const should = require('should');
const fs = require('fs');
const environment = require('../../../environment/env')

const logRequest = require('../../../utility/logRequest')
const appendToken = require('../../../utility/appendToken')
const sendReq = require('../../../utility/asyncSendReq')
var bigDecimal = require('js-big-decimal');

let quote_requestID, quote_id, objJsonB64QuoteR, signed_payload, rfi_quote_obj

Given('id: {string} request a exchange_requset with time_expire: {string} limit_max: {string} limit_min: {string} from source_asset: asset_code: {string} asset_type: {string} issuer_id: {string} to target_asset: asset_code {string} asset_type: {string} issuer_id: {string}, quote_service_url: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, time_expire, limit_max, limit_min, source_asset_code, source_asset_type, source_issuer_id, targert_asset_code, target_asset_type, target_issuer_id, quote_service_url, done) {

    var options = {
        method: 'POST',
        url: environment[quote_service_url] + '/client/quotes/request',
        headers: {},
        body: {
            time_expire: parseInt(environment[time_expire]),
            limit_max: parseInt(environment[limit_max]),
            limit_min: parseInt(environment[limit_min]),
            source_asset: {
                asset_code: environment[source_asset_code],
                asset_type: environment[source_asset_type],
                issuer_id: environment[source_issuer_id]
            },
            target_asset: {
                asset_code: environment[targert_asset_code],
                asset_type: environment[target_asset_type],
                issuer_id: environment[target_issuer_id]
            },
            ofi_id: environment[id]
        },
        json: true
    };

    options = appendToken(options, id, true)
    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            let today = new Date();
            fs.appendFile('./file/exchange_requestID.txt', today + ' -  ' + JSON.stringify(res.body) + '\n', function(err) {
                if (err) throw err;
                // console.log('Saved!');
            });
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            quote_requestID = body.request_id;
            logRequest(res, options)
            done()
        }
    });

});


When('id: {string} get quote from requestID check rfi: {string} was in the quote list and status: {string} , quote_service_url: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, rfi_id, status, quote_service_url, done) {
    setTimeout(function() {
        var options = {
            method: 'GET',
            url: environment[quote_service_url] + '/client/quotes/request/' + quote_requestID,
            headers: {}
        };

        let quote = {
            "request_id": quote_requestID,
            "rfi_id": environment[rfi_id],
            "status": parseInt(environment[status]),
        }

        options = appendToken(options, id, true)
        request(options, function(err, res, body) {
            if (err) {
                logRequest(err, options)
                done(err)
            } else {
                should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                JSON.parse(body).should.containDeep([quote], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                JSON.parse(body).forEach(function(element) {

                    if (element.rfi_id == environment[rfi_id] && element.status == environment[status]) {
                        quote_id = element.quote_id
                    }
                });
                logRequest(res, options)
                done()
            }
        });
    }, 10000)

});

When('id: {string} check quotes status: {string} from quotes_id, quote_service_url: {string}', {
        timeout: parseInt(environment.MAX_TIMEOUT)
    },
    function(rfi_id, status, quote_service_url, done) {
        var options = {
            method: 'GET',
            url: environment[quote_service_url] + '/client/quotes/' + quote_id,
            headers: {}
        };

        options = appendToken(options, rfi_id, true)
        request(options, function(err, res, body) {
            if (err) {
                logRequest(err, options)
                done(err)
            } else {
                should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).rfi_id).be.exactly(environment[rfi_id], JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).status).be.exactly(parseInt(environment[status]), +"\n request data : " + JSON.stringify(options))
                logRequest(res, options)
                done()
            }
        });
    });


When('id: {string} response quotes, account_name_receive: {string} account_name_send: {string} exchange_rate: {string} limit_max: {string} response_time_expire: {string} response_time_start: {string} to quote request time_expire: {string} limit_max: {string} limit_min: {string} from source_asset: asset_code: {string} asset_type: {string} issuer_id: {string} to target_asset: asset_code {string} asset_type: {string} issuer_id: {string} by ofi: {string}, crypto_service_url: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, account_name_receive, account_name_send, exchange_rate, rfi_limit_max, response_time_expire, response_time_start, ofi_time_expire, ofi_limit_max, ofi_limit_min, source_asset_code, source_asset_type, source_issuer_id, target_asset_code, target_asset_type, target_issuer_id, ofi_id, crypto_service_url, done) {
    rfi_quote_obj = {
        account_name_receive: environment[account_name_receive],
        account_name_send: environment[account_name_send],
        exchange_rate: new bigDecimal(environment[exchange_rate]).getValue(),
        // exchange_rate: parseInt(environment[exchange_rate]),
        limit_max: parseInt(environment[rfi_limit_max]),
        quote_id: quote_id,
        quote_request: {
            time_expire: parseInt(environment[ofi_time_expire]),
            limit_max: parseInt(environment[ofi_limit_max]),
            limit_min: parseInt(environment[ofi_limit_min]),
            source_asset: {
                asset_code: environment[source_asset_code],
                asset_type: environment[source_asset_type],
                issuer_id: environment[source_issuer_id]
            },
            target_asset: {
                asset_code: environment[target_asset_code],
                asset_type: environment[target_asset_type],
                issuer_id: environment[target_issuer_id]
            },
            ofi_id: environment[ofi_id]
        },
        rfi_id: environment[id],
        time_expire: parseInt(environment[response_time_expire]),
        time_start: parseInt(environment[response_time_start])
    }
    let objJsonStr = JSON.stringify(rfi_quote_obj);
    objJsonB64QuoteR = Buffer.from(objJsonStr).toString("base64");

    var options = {
        method: 'POST',
        url: environment[crypto_service_url] + '/client/sign',
        headers: {},
        body: {
            account_name: environment[account_name_send],
            payload: objJsonB64QuoteR
        },
        json: true
    };

    options = appendToken(options, id)
    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            signed_payload = body.transaction_signed
            logRequest(res, options)
            done()
        }
    });


});

When('id: {string} post response, quote_service_url: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, quote_service_url, done) {


    var options = {
        method: 'POST',
        url: environment[quote_service_url] + '/client/quotes/' + quote_id,
        headers: {},
        body: {
            quote: objJsonB64QuoteR,
            signature: signed_payload
        },
        json: true
    };

    options = appendToken(options, id, true)
    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            logRequest(res, options)
            done()
        }
    });

});

When('id: {string} sign quotes request with exchange_amount: {string} account_name_receive: {string} account_name_send: {string}, crypto_service_url: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, exchange_amount, account_name_receive, account_name_send, crypto_service_url, done) {
    let ofi_quote_obj = {
        account_name_receive: environment[account_name_receive],
        account_name_send: environment[account_name_send],
        amount: parseInt(environment[exchange_amount]),
        quote: rfi_quote_obj
    }
    let objJsonStr = JSON.stringify(ofi_quote_obj);
    objJsonB64QuoteR = Buffer.from(objJsonStr).toString("base64");

    var options = {
        method: 'POST',
        url: environment[crypto_service_url] + '/client/sign',
        headers: {},
        body: {
            account_name: environment[account_name_send],
            payload: objJsonB64QuoteR
        },
        json: true
    };

    options = appendToken(options, id)
    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            logRequest(res, options)
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            signed_payload = body.transaction_signed
            done()
        }
    });

});


When('id: {string} post exchange, quote_service_url: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, quote_service_url, done) {


    var options = {
        method: 'POST',
        url: environment[quote_service_url] + '/client/exchange',
        headers: {},
        body: {
            exchange: objJsonB64QuoteR,
            signature: signed_payload
        },
        json: true
    };

    options = appendToken(options, id, true)
    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            logRequest(res, options)
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            done()
        }
    });
});


Given('id: {string} pick up message from RFI_QUOTES topic sent by ofi {string} should get time_expire: {string} limit_max: {string} limit_min: {string} from asset_code: {string} asset_type: {string} issuer_id: {string} to asset_code {string} asset_type: {string} issuer_id: {string}, wwGatewayURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, ofi_id, time_expire, limit_max, limit_min, source_asset_code, source_asset_type, issuer_id, target_asset_code, target_asset_type, target_issuer_id, wwGateWayURL, done) {
    setTimeout(async function() {
        var options = {
            method: 'GET',
            url: environment[wwGateWayURL] + '/client/message',
            headers: {},
            qs: { type: 'quotes' }
        };

        options = appendToken(options, id)

        try {

            let data = {
                "quote_id": quote_id,
                "quote_request": {
                    "limit_max": environment[limit_max].toString(),
                    "limit_min": environment[limit_min].toString(),
                    "ofi_id": environment[ofi_id],
                    "source_asset": {
                        "asset_code": environment[source_asset_code],
                        "asset_type": environment[source_asset_type],
                        "issuer_id": environment[issuer_id]
                    },
                    "target_asset": {
                        "asset_code": environment[target_asset_code],
                        "asset_type": environment[target_asset_type],
                        "issuer_id": environment[target_issuer_id]
                    },
                    "time_expire": environment[time_expire]
                }
            }

            let retry = environment.ENV_KEY_GATEWAY_RETRY_TIMES
            let res
            while (retry > 0) {
                res = await sendReq(options)
                    // console.log(JSON.parse(res.body))
                if (JSON.parse(res.body).data == null) {
                    retry--
                } else {
                    break
                }
            }

            JSON.parse(res.body).data.should.containEql(data, "JSON.parse(body): " + JSON.parse(res.body))
            done()

        } catch (err) {
            console.log(err);
            // done()

        }
    }, 2000)
});