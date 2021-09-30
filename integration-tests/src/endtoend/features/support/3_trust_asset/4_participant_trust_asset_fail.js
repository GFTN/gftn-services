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
const environment = require('../../../environment/env')

const logRequest = require('../../../utility/logRequest')
const appendToken = require('../../../utility/appendToken')


Then('id: {string} send trust request to trust asset_code: {string} asset_issuer: {string} with trust_limit {string} in trust_account {string}, apiURL:{string} - should return error code {string} error message {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, asset_issuer, trust_limit, trust_account, url_api_service, error_code, error_message, done) {
    var options = {
        method: 'POST',
        url: environment[url_api_service] + '/client/trust',
        headers: {},
        body: {
            permission: 'request',
            asset_code: environment[asset_code],
            account_name: environment[trust_account],
            participant_id: environment[asset_issuer],
            limit: parseInt(environment[trust_limit])
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
            should(res.statusCode).be.exactly(parseInt(environment[error_code]), "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            should(body.code).be.exactly(environment[error_message], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            done()
        }
    });

});