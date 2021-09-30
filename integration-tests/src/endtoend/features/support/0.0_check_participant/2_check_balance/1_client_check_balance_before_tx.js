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
const environment = require('../../../../environment/env')
    // ../../../../environment/env
const logRequest = require('../../../../utility/logRequest')
const appendToken = require('../../../../utility/appendToken')
var bigDecimal = require('js-big-decimal');

Given('id: {string} check asset_code: {string} issuer: {string} issued do balance before transaction, participant_api_url:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, issuer, participant_api_url, done) {
    setTimeout(function() {
        var options = {
            method: 'GET',
            url: environment[participant_api_url] + '/client/obligations',
            qs: { asset_code: environment[asset_code] },
            headers: {}
        };

        options = appendToken(options, id)
        request(options, function(err, res, body) {
            if (err) {

                logRequest(err, options)
                done(err)
            } else {
                logRequest(res, options)
                if (res.statusCode == 200) {
                    JSON.parse(body).forEach(element => {
                        if (element.asset_code == environment[asset_code]) {
                            process.env[environment[id] + "_" + environment[asset_code] + "_" + environment[issuer] + "_DO_BALANCE"] = element.balance
                            should(element.issuer_id).be.exactly(environment[issuer], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                            should(element.asset_code).be.exactly(environment[asset_code], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                        }
                    });

                } else {
                    process.env[environment[id] + "_" + environment[asset_code] + "_" + environment[issuer] + "_DO_BALANCE"] = 0
                }
                done()
            }
        });
    }, 1000)
});

Given('id: {string} check asset_code: {string} issuer: {string} account_name: {string} balance before transaction, participant_api_url:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, issuer, account_name, participant_api_url, done) {
    setTimeout(function() {
        var options = {
            method: 'GET',
            url: environment[participant_api_url] + '/client/balances/accounts/' + environment[account_name],
            qs: { asset_code: environment[asset_code], issuer_id: environment[issuer] },
            headers: {}
        };

        options = appendToken(options, id)
        request(options, function(err, res, body) {
            if (err) {

                logRequest(err, options)
                done(err)
            } else {
                logRequest(res, options)
                if (res.statusCode == 200) {
                    process.env[environment[id] + "_" + environment[account_name] + "_" + environment[asset_code] + "_" + environment[issuer]] = JSON.parse(body).balance;
                    should(JSON.parse(body).issuer_id).be.exactly(environment[issuer], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    should(JSON.parse(body).asset_code).be.exactly(environment[asset_code], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

                } else {
                    process.env[environment[id] + "_" + environment[account_name] + "_" + environment[asset_code] + "_" + environment[issuer]] = 0
                }
                // console.log(process.env[environment[id] + "_" + environment[account_name] + "_" + environment[asset_code] + "_" + environment[issuer]]);
                done()
            }
        });
    }, 100)
});