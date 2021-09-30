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



Given('id: {string} check id: {string} operating account:{string} exist, participantAPIURL: {string}', {
        timeout: parseInt(environment.MAX_TIMEOUT)
    },
    function(tokenid, id, account_name, url_api_service, done) {
        let options = {
            method: 'GET',
            headers: {},
            url: environment[url_api_service] + "/client/participants/" + environment[id]
        };

        options = appendToken(options, tokenid)
        request(options, function(err, res, body) {
            if (err) {
                logRequest(err, options)
                done(err)
            } else {
                logRequest(res, options)
                let acc;
                should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                JSON.parse(body).operating_accounts.find(function(element) {
                    if (element.name == environment[account_name]) {
                        acc = element.address
                    }
                });
                should(acc.length).be.exactly(56)
                done()
            }
        });
    });


Then('id: {string} check participant: {string} operating account: {string} not exist - apiURL: {string}', {
        timeout: parseInt(environment.MAX_TIMEOUT)
    },
    function(tokenid, id, account_name, url_api_service, done) {
        let options = {
            method: 'GET',
            headers: {},
            url: environment[url_api_service] + "/client/participants/" + environment[id]
        };

        options = appendToken(options, tokenid)
        request(options, function(err, res, body) {
            if (err) {
                logRequest(err, options)
                done(err)
            } else {
                logRequest(res, options)
                let acc;
                should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                JSON.parse(body).operating_accounts.find(function(element) {
                    if (element.name == environment[account_name]) {
                        acc = element.address
                    }
                });
                should(acc.length).be.exactly(0)
                done()
            }
        });
    });