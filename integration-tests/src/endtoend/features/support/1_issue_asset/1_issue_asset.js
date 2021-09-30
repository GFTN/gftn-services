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


When('id: {string} issue asset asset_code: {string}, asset_type: {string} - apiURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, asset_type, url_api_service, done) {

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
            should(res.statusCode).be.exactly(200, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            done()
        }
    });
});


Then('id: {string} check id: {string} has asset_code: {string}, asset_type: {string} in issued asset list - apiURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(tokenid, id, asset_code, asset_type, url_api_service, done) {

    let options = {
        contentType: 'application/json',
        method: 'GET',
        headers: {},
        url: environment[url_api_service] + '/client/assets/issued'
    }

    let asset = {
        "asset_code": environment[asset_code],
        "asset_type": environment[asset_type],
        "issuer_id": environment[id]
    }
    options = appendToken(options, tokenid)
    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            logRequest(res, options)
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            JSON.parse(body).should.containDeep([asset], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            done()
        }
    });
});


Then('id: {string} query id: {string} has asset_code: {string}, asset_type: {string} in issued asset list - apiURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(tokenid, id, asset_code, asset_type, url_api_service, done) {

    let options = {
        contentType: 'application/json',
        method: 'GET',
        headers: {},
        url: environment[url_api_service] + '/client/assets/participants/' + environment[id] + '?type=issued'
    }

    let asset = {
        "asset_code": environment[asset_code],
        "asset_type": environment[asset_type],
        "issuer_id": environment[id]
    }
    options = appendToken(options, tokenid)
    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            logRequest(res, options)
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            JSON.parse(body).should.containDeep([asset], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            done()
        }
    });
});