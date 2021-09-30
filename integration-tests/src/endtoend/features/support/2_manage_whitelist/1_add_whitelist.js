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


When('id: {string} add participant: {string} - whitelistURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, target_participant, url_wl_service, done) {
    var options = {
        method: 'POST',
        url: environment[url_wl_service] + '/client/participants/whitelist',
        headers: {},
        body: '{\n"participant_id":"' + environment[target_participant] + '"\n}'

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