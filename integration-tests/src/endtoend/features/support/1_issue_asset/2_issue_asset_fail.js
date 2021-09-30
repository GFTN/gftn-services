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


When('id: {string} issue asset with incalid asset_code: {string}, asset_type: {string} - apiURL: {string} - should return error code {string} error message {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, asset_type, url_api_service, error_code, error_message, done) {

    let options = {
        method: 'POST',
        headers: {},
        url: environment[url_api_service] + '/client/assets',
        qs: { asset_code: environment[asset_code], asset_type: environment[asset_type] }
    }

    options = appendToken(options, id)
    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            logRequest(res, options)
            should(res.statusCode).be.exactly(parseInt(environment[error_code]), "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            should(JSON.parse(body).code).be.exactly(environment[error_message], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            done()
        }
    });
});