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
const environment = require('../../../../../environment/env')

const logRequest = require('../../../../../utility/logRequest')
const appendToken = require('../../../../../utility/appendToken')

var bigDecimal = require('js-big-decimal');


Then('id: {string} check asset_code: {string} issuer: {string} account_name: {string} balance does not increase, participant_api_url:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, issuer, account_name, participant_api_url, done) {
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

                let ori_account_balance = new bigDecimal(process.env[environment[id] + "_" + environment[account_name] + "_" + environment[asset_code] + "_" + environment[issuer]]);

                should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).issuer_id).be.exactly(environment[issuer], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).asset_code).be.exactly(environment[asset_code], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).balance).be.exactly(ori_account_balance.getValue(), "ori_account_balance: " + ori_account_balance + "\n response balance: " + JSON.stringify(res.body))

                done()
            }
        });
    }, 5000)
})


Given('id: {string} check asset_code: {string} issuer: {string} issued do balance does not increase, participant_api_url:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, issuer, participant_api_url, done) {
    setTimeout(function() {
        var options = {
            method: 'GET',
            url: environment[participant_api_url] + '/client/obligations',
            qs: { asset_code: environment[asset_code] },
            headers: {}
        };

        let asset = {
            "account_name": "issuing",
            "asset_code": environment[asset_code],
            "issuer_id": environment[issuer],
            "balance": process.env[environment[id] + "_" + environment[asset_code] + "_" + environment[issuer] + "_DO_BALANCE"]
        }


        options = appendToken(options, id)
        request(options, function(err, res, body) {
            if (err) {

                logRequest(err, options)
                done(err)
            } else {
                logRequest(res, options)

                if (res.statusCode == 200) {
                    should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    JSON.parse(body).should.containDeep([asset], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    done()

                } else {
                    should(0).be.exactly(parseInt(process.env[environment[id] + "_" + environment[asset_code] + "_" + environment[issuer] + "_DO_BALANCE"]), "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    done()
                }
            }
        });
    }, 10)
});