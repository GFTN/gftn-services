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



When('id: {string} check the asset_code: {string} asset_type: {string} issuer: {string} already issue by worldwire, apiURL:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, asset_type, issuer_id, url_api_service, done) {

    var options = {
        method: 'GET',
        url: environment[url_api_service] + '/client/assets',
        headers: {}
    };

    options = appendToken(options, id)
    let asset = {
        "asset_code": environment[asset_code],
        "asset_type": environment[asset_type],
        "issuer_id": environment[issuer_id]
    }
    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            logRequest(res, options)
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            JSON.parse(body).length.should.be.aboveOrEqual(0, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            JSON.parse(body).should.containEql(asset, "JSON.parse(body): " + JSON.parse(body) + "asset: " + asset)
            done()
        }
    });
});


Given('id: {string} asset_code: {string} asset_type: {string} issuer: {string} in trust_account {string} trusted list, apiURL:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, asset_type, issuer_id, trust_account, url_api_service, done) {
    var options = {
        method: 'GET',
        url: environment[url_api_service] + '/client/assets/accounts/' + environment[trust_account],
        headers: {}
    };
    let asset = {
        "asset_code": environment[asset_code],
        "asset_type": environment[asset_type],
        "issuer_id": environment[issuer_id]
    }
    options = appendToken(options, id)

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


Then('id: {string} query {string} asset_code: {string} asset_type: {string} issuer: {string} in trust_account {string} trusted list, apiURL:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(tokenid, id, asset_code, asset_type, issuer_id, trust_account, url_api_service, done) {
    var options = {
        method: 'GET',
        url: environment[url_api_service] + '/client/assets/participants/' + environment[id] + '?type=trusted',
        headers: {}
    };
    let asset = {
        "asset_code": environment[asset_code],
        "asset_type": environment[asset_type],
        "issuer_id": environment[issuer_id]
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


When('id: {string} asset_code: {string} asset_type: {string} issuer: {string} in issuing account trusted list, apiURL:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, asset_type, issuer_id, url_api_service, done) {
    var options = {
        method: 'GET',
        url: environment[url_api_service] + '/client/assets/accounts/issuing',
        headers: {}
    };
    let asset = {
        "asset_code": environment[asset_code],
        "asset_type": environment[asset_type],
        "issuer_id": environment[issuer_id]
    }
    options = appendToken(options, id)

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