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


Then('id: {string} check participant: {string} was in the whitelist - whitelistURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, target_participant, url_wl_service, done) {
    var options = {
        method: 'GET',
        url: environment[url_wl_service] + '/client/participants/whitelist',
        headers: {}
    };
    options = appendToken(options, id, true)
    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            logRequest(res, options)
            body.should.containEql(environment[target_participant], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            done()
        }
    });

});


Then('id: {string} check participant: {string} was not in whitelist - whitelistURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, target_participant, url_wl_service, done) {
    var options = {
        method: 'GET',
        url: environment[url_wl_service] + '/client/participants/whitelist',
        headers: {}
    };
    options = appendToken(options, id, true)
    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            logRequest(res, options)
            body.should.not.containEql(environment[target_participant], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            done()
        }
    });

});


Given('id: {string} check there has no participants in whitelist - whitelistURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, url_wl_service, done) {
    var options = {
        method: 'GET',
        url: environment[url_wl_service] + '/client/participants/whitelist',
        headers: {}
    };
    options = appendToken(options, id, true)
    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            logRequest(res, options)
                // body.should.not.containEql(environment[target_participant], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            should(JSON.parse(body)).be.exactly(null, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            done()
        }
    });
});