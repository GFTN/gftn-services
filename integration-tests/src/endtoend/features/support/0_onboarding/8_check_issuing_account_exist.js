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


Then('id: {string} check id: {string} issuing account exist, participantAPIURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(tokenid, id, url_api_service, done) {
    let options = {
        method: 'GET',
        headers: {},
        url: environment[url_api_service] + "/client/participants/" + environment[id]
    };

    options = appendToken(options, tokenid)
    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            logRequest(res, options)
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            should(JSON.parse(body).issuing_account.length).be.exactly(56)
            done()
        }
    });
});