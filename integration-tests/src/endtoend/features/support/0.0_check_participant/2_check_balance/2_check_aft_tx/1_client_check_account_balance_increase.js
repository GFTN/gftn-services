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

Then('id: {string} check asset_code: {string} issuer: {string} account_name: {string} balance increase, participant_api_url:{string}', {
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
                console.log(environment[id] + "_" + environment[account_name] + "_" + environment[asset_code] + "_" + environment[issuer] + " = " + process.env[environment[id] + "_" + environment[account_name] + "_" + environment[asset_code] + "_" + environment[issuer]])
                let ori_account_balance = new bigDecimal(process.env[environment[id] + "_" + environment[account_name] + "_" + environment[asset_code] + "_" + environment[issuer]]);

                should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options));
                should(JSON.parse(body).issuer_id).be.exactly(environment[issuer], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).asset_code).be.exactly(environment[asset_code], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

                parseInt(JSON.parse(body).balance).should.be.above(ori_account_balance.getValue(), "ori_balance=" + ori_account_balance.getValue() + " , after_balance=" + JSON.parse(body).balance + "\n response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

                done()
            }
        });
    }, 20000)
})

Then('id: {string} check asset_code: {string} issuer: {string} account_name: {string} balance increase: {string}, participant_api_url:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, issuer, account_name, increase_amount, participant_api_url, done) {
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
                let aft_account_balance = ori_account_balance.add(new bigDecimal(environment[increase_amount]))

                // console.log(aft_account_balance.getValue());
                let balance = new bigDecimal(JSON.parse(body).balance)
                    // console.log(balance.getValue());
                should(balance.getValue()).be.exactly(aft_account_balance.getValue(), "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).issuer_id).be.exactly(environment[issuer], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).asset_code).be.exactly(environment[asset_code], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

                done()
            }
        });
    }, 20000)
});


Given('id: {string} check asset_code: {string} issuer: {string} issued do balance increase settle_amount ofi: {string}, participant_api_url: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, issuer, OFIid, participant_api_url, done) {
    setTimeout(function() {
        var options = {
            method: 'GET',
            url: environment[participant_api_url] + '/client/obligations',
            qs: { asset_code: environment[asset_code] },
            headers: {}
        };



        let ori_account_balance = new bigDecimal(process.env[environment[id] + "_" + environment[asset_code] + "_" + environment[issuer] + "_DO_BALANCE"]);
        let settle_amount = new bigDecimal(process.env[environment[OFIid] + "_REQUEST_SETTLE_AMOUNT_" + environment[asset_code] + "_" + environment[issuer]]);
        let aft_account_balance = ori_account_balance.add(settle_amount)


        let asset = {
            "account_name": "issuing",
            "asset_code": environment[asset_code],
            "issuer_id": environment[issuer],
            // "balance": aft_account_balance.getValue()
        }

        options = appendToken(options, id)
        request(options, function(err, res, body) {
            if (err) {

                logRequest(err, options)
                done(err)
            } else {

                logRequest(res, options)

                JSON.parse(body).find(function(asset) {
                    if (asset.asset_code == environment[asset_code] && asset.issuer_id == environment[issuer]) {
                        should(aft_account_balance.getValue()).be.exactly(asset.balance, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    }
                });

                should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                JSON.parse(body).should.containDeep([asset], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                done()
            }
        });
    }, 20000)
});


Then('id: {string} check asset_code: {string} issuer: {string} account_name: {string} balance increase settle amount ofi: {string} participant_api_url:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, issuer, account_name, OFIid, participant_api_url, done) {
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
                let settle_amount = new bigDecimal(process.env[environment[OFIid] + "_REQUEST_SETTLE_AMOUNT_" + environment[asset_code] + "_" + environment[issuer]]);
                let aft_account_balance = ori_account_balance.add(settle_amount)

                let balance = new bigDecimal(JSON.parse(body).balance)
                should(balance.getValue()).be.exactly(aft_account_balance.getValue(), "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).issuer_id).be.exactly(environment[issuer], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).asset_code).be.exactly(environment[asset_code], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

                done()
            }
        });
    }, 20000)
});