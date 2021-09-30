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


Then('id: {string} check asset_code: {string} issuer: {string} account_name: {string} balance decrease: {string}, participant_api_url:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, issuer, account_name, decrease_amount, participant_api_url, done) {
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
                let aft_account_balance = ori_account_balance.subtract(new bigDecimal(environment[decrease_amount]))

                let balance = new bigDecimal(JSON.parse(body).balance)
                should(balance.getValue()).be.exactly(aft_account_balance.getValue(), "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).issuer_id).be.exactly(environment[issuer], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).asset_code).be.exactly(environment[asset_code], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

                done()
            }
        });
    }, 20000)
});




Then('id: {string} check asset_code: {string} issuer: {string} issued do balance decrease: {string} by RFI send back to OFI: {string}, participant_api_url:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, issuer, decrease_amount, OFIid, participant_api_url, done) {
    setTimeout(function() {
        var options = {
            method: 'GET',
            url: environment[participant_api_url] + '/client/obligations',
            qs: { asset_code: environment[asset_code] },
            headers: {}
        };

        let ori_account_balance = new bigDecimal(process.env[environment[id] + "_" + environment[asset_code] + "_" + environment[issuer] + "_DO_BALANCE"]);
        let decrease = new bigDecimal(environment[decrease_amount]);
        let aft_account_balance = ori_account_balance.subtract(decrease)

        let asset = {
            "account_name": "issuing",
            "asset_code": environment[asset_code],
            "issuer_id": environment[issuer],
            "balance": aft_account_balance.getValue()
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
    }, 20000)
});


Then('id: {string} check asset_code: {string} issuer: {string} issued do balance decrease settle_amount by RFI send back to OFI: {string}, participant_api_url:{string}', {
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
        let aft_account_balance = ori_account_balance.subtract(settle_amount)

        let asset = {
            "account_name": "issuing",
            "asset_code": environment[asset_code],
            "issuer_id": environment[issuer],
            "balance": aft_account_balance.getValue()
        }

        options = appendToken(options, id)
        request(options, function(err, res, body) {
            if (err) {

                logRequest(err, options)
                done(err)
            } else {
                logRequest(res, options)
                if (res.statusCode == 200) {
                    JSON.parse(body).should.containDeep([asset], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                }
                // should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

                done()
            }
        });
    }, 20000)
});

Then('id: {string} check asset_code: {string} issuer: {string} account_name: {string} balance decrease settle_amount which RFI send back to OFI: {string}, participant_api_url:{string}', {
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
                let aft_account_balance = ori_account_balance.subtract(settle_amount)


                let balance = new bigDecimal(JSON.parse(body).balance)
                should(balance.getValue()).be.exactly(aft_account_balance.getValue(), "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).issuer_id).be.exactly(environment[issuer], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).asset_code).be.exactly(environment[asset_code], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

                done()
            }
        });
    }, 20000)
});

Then('id: {string} check asset_code: {string} issuer: {string} account_name: {string} balance decrease: {string} which RFI send back to OFI: {string}, participant_api_url:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, asset_code, issuer, account_name, decrease_amount, OFIid, participant_api_url, done) {
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
                let decrease = new bigDecimal(environment[decrease_amount]);
                let aft_account_balance = ori_account_balance.subtract(decrease)


                let balance = new bigDecimal(JSON.parse(body).balance)
                should(balance.getValue()).be.exactly(aft_account_balance.getValue(), "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).issuer_id).be.exactly(environment[issuer], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                should(JSON.parse(body).asset_code).be.exactly(environment[asset_code], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

                done()
            }
        });
    }, 20000)
});