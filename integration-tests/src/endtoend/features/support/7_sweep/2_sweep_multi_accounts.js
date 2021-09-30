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
var bigDecimal = require('js-big-decimal');
let sourceList = []
Then('id: {string} add asset_code: {string} issuer: {string} asset_type: {string} account_name: {string} amount: {string} to source list', function(id, asset_code, issuer, asset_type, account_name, amount, done) {

    let source = {
        "account": environment[account_name],
        "amount": parseFloat(environment[amount]),
        "asset": {
            "asset_code": environment[asset_code],
            "issuer_id": environment[issuer],
            "asset_type": environment[asset_type]
        }
    }
    if (process.env["increaseList_" + environment[asset_code] + "_" + environment[issuer]] == undefined) {
        process.env["increaseList_" + environment[asset_code] + "_" + environment[issuer]] = environment[amount]
    } else {
        let ori_increase = new bigDecimal(process.env["increaseList_" + environment[asset_code] + "_" + environment[issuer]]);
        let sum_increase = ori_increase.add(new bigDecimal(environment[amount]))
        process.env["increaseList_" + environment[asset_code] + "_" + environment[issuer]] = sum_increase.getValue()
    }

    sourceList.push(source)
    done()
});

When('source list is ready', function(done) {
    sourceList.length.should.be.above(0, "source list :" + sourceList + "\n source list is empty")
    done()
});

Then('id: {string} do sweep from source list to account:{string}, participantAPIURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, targer_account_name, url_api_service, done) {
    var options = {
        method: 'POST',
        url: environment[url_api_service] + '/client/account/sweep',
        headers: {},
        body: {
            target_account: environment[targer_account_name],
            source_accounts: sourceList
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
            should(res.statusCode).be.exactly(200, JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            done()
        }
    });

});