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


Then('allower anchor: {string} send alow trust request to allow trust DA asset_code: {string} trustSender: {string} trust_account {string}, allowURL:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(anchor_id, asset_code, trustSender, trust_account, allowURL, done) {

    var options = {
        method: 'POST',
        url: environment[allowURL] + '/anchor/trust/' + environment[anchor_id],
        headers: {},
        body: {
            permission: 'allow',
            asset_code: environment[asset_code],
            account_name: environment[trust_account],
            participant_id: environment[trustSender]
        },
        json: true
    };

    options = appendToken(options, anchor_id, true)

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