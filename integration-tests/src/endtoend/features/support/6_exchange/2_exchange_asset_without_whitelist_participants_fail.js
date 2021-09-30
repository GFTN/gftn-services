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

When('id: {string} request a exchange_requset with time_expire: {string} limit_max: {string} limit_min: {string} from source_asset: asset_code: {string} asset_type: {string} issuer_id: {string} to target_asset: asset_code {string} asset_type: {string} issuer_id: {string}, quote_service_url: {string} - should return error code {string} error message {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, time_expire, limit_max, limit_min, source_asset_code, source_asset_type, source_issuer_id, targert_asset_code, target_asset_type, target_issuer_id, quote_service_url, error_code, error_message, done) {

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
            logRequest(res, options)
            should(res.statusCode).be.exactly(parseInt(environment[error_code]), "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            should(body.code).be.exactly(environment[error_message], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            done()
        }
    });

});