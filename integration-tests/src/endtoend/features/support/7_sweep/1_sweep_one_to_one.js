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

When('id: {string} sweep asset_code: {string} issuer: {string} asset_type: {string} amount: {string} account_name: {string} to account_name: {string}, participant_api_url:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, issuer, asset_type, amount, source_account_name, targer_account_name, url_api_service, done) {
    var options = {
        method: 'POST',
        url: environment[url_api_service] + '/client/account/sweep',
        headers: {},
        body: {
            target_account: environment[targer_account_name],
            source_accounts: [{
                account: environment[source_account_name],
                amount: parseFloat(environment[amount]),
                asset: { asset_code: environment[asset_code], issuer_id: environment[issuer], asset_type: environment[asset_type] }
            }]
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
            should(res.statusCode).be.exactly(200, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            done()
        }
    });

});