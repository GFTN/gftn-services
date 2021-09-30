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


Given('id: {string} check participant: {string}, was in Workd Wire, apiURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, target_participant, url_api_service, done) {
    var options = {
        method: 'GET',
        url: environment[url_api_service] + '/client/participants/' + environment[target_participant],
        headers: {}
    };
    options = appendToken(options, id)

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


Given('id: {string} check participant: {string}, was in Workd Wire, anchorURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, target_participant, anchorURL, done) {
    var options = {
        method: 'GET',
        url: environment[anchorURL] + '/anchor/participants/' + environment[target_participant],
        headers: {}
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