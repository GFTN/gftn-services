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


Then('id: {string} send trust request to trust asset_code: {string} asset_issuer: {string} with trust_limit {string} in trust_account {string}, apiURL:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, asset_issuer, trust_limit, trust_account, url_api_service, done) {
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
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            done()
        }
    });

});

Given('allower participant: {string} send alow trust request to allow trust DO asset_code: {string} trustSender: {string} trust_account {string}, allowURL:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, trustSender, trust_account, url_api_service, done) {

    var options = {
        method: 'POST',
        url: environment[url_api_service] + '/client/trust',
        headers: {},
        body: {
            permission: 'allow',
            asset_code: environment[asset_code],
            account_name: environment[trust_account],
            participant_id: environment[trustSender]
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
            done()
        }
    });
});


Then('allower participant: {string} send revoke trust request to revoke trust DO asset_code: {string} trustSender: {string} trust_account {string}, allowURL:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, trustSender, trust_account, url_api_service, done) {

    var options = {
        method: 'POST',
        url: environment[url_api_service] + '/client/trust',
        headers: {},
        body: {
            permission: 'revoke',
            asset_code: environment[asset_code],
            account_name: environment[trust_account],
            participant_id: environment[trustSender]
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
            done()
        }
    });
});