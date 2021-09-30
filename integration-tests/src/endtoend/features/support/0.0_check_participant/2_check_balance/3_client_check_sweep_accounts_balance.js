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

const logRequest = require('../../../../utility/logRequest')
const appendToken = require('../../../../utility/appendToken')
var bigDecimal = require('js-big-decimal');


Then('id: {string} check asset_code: {string} issuer: {string} account_name: {string} balance is greater than sweep amount: {string}, participant_api_url:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, issuer, account_name, sweep_amount, participant_api_url, done) {
    var options = {
        method: 'GET',
        url: environment[participant_api_url] + '/client/balances/accounts/' + environment[account_name],
        qs: { asset_code: environment[asset_code], issuer_id: environment[issuer] },
        headers: {}
    };

    options = appendToken(options, id)
    setTimeout(function() {
        request(options, function(err, res, body) {
            if (err) {

                logRequest(err, options)
                done(err)
            } else {
                logRequest(res, options)
                should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).issuer_id).be.exactly(environment[issuer], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).asset_code).be.exactly(environment[asset_code], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                JSON.parse(body).balance.should.be.above(parseInt(environment[sweep_amount]), "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

                done()
            }
        });
    }, 1000)
});


Then('id: {string} check asset_code: {string} issuer: {string} account_name: {string} balance from sourceList, participant_api_url:{string}', {
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
                let ori_account_balance = new bigDecimal(process.env[environment[id] + "_" + environment[account_name] + "_" + environment[asset_code] + "_" + environment[issuer]]);
                let sum_increase_amount = new bigDecimal(process.env["increaseList_" + environment[asset_code] + "_" + environment[issuer]]);
                let aft_account_balance = ori_account_balance.add(sum_increase_amount)

                let balance = new bigDecimal(JSON.parse(body).balance)

                should(balance.getValue()).be.exactly(aft_account_balance.getValue(), "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).issuer_id).be.exactly(environment[issuer], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).asset_code).be.exactly(environment[asset_code], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

                done()
            }
        });
    }, 100)
});