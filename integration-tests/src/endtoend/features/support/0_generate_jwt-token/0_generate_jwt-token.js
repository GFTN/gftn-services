// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const {
    Given,
    When,
    Then
} = require('cucumber')
var request = require("request");
var speakeasy = require("speakeasy");
const environment = require('../../../environment/env')
var jwtDecode = require('jwt-decode');
const generateFID = require('../../../utility/generateFID')
const getIP = require('external-ip')();
const logRequest = require('../../../utility/logRequest')

Given('user: {string} generate a FID', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(userEmail, done) {
    generateFID(environment[userEmail]).then(function(fid) {
        console.log(fid);
        process.env[environment[userEmail] + '_FID'] = fid
            // var decoded = jwtDecode(fid);
            // console.log(decoded);
        done()
    })
});


Given('user: {string} send request to request a jwt-token for participant: {string} IID: {string} using totpkey: {string}, auth_url:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(userEmail, participantID, IID, totpkey, auth_url, done) {
    let fid = process.env[environment[userEmail] + '_FID']
    let totpCode = speakeasy.totp({
        secret: environment[totpkey],
        encoding: 'ascii'
    });


    getIP((err, ip) => {
        if (err) {
            // every service in the list has failed
            throw err;
        }
        var options = {
            method: 'POST',
            url: environment[auth_url] + '/jwt/request',
            headers: {
                'x-verify-code': totpCode,
                'x-fid': fid,
                'x-iid': environment[IID],
                'Content-Type': 'application/json'
            },
            body: {
                description: 'testing purpose',
                acc: ['issuing',
                    'default',
                    'default-2'
                ],
                ver: environment.ENV_VERSION,
                ips: [ip,
                    // "202.135.245.39",
                    // "202.135.245.4",
                    // "219.74.15.207",
                    // "202.135.245.2",
                    // "125.18.9.20",
                    // "169.63.140.106",
                    // "202.135.245.35",
                    // "36.226.251.220",
                    // "39.12.202.15",
                    // "202.135.245.3",
                    // "202.135.245.4",
                    // "27.104.245.116",
                    // "202.135.245.40"
                ],
                env: environment.ENV_ENVIRONMNET,
                enp: ["/v1/admin/pr",
                    "/v1/admin/pr/domain",
                    "/v1/anchor/address",
                    "/v1/anchor/assets/issued",
                    "/v1/anchor/assets/redeem",
                    "/v1/anchor/fundings/instruction",
                    "/v1/anchor/fundings/send",
                    "/v1/anchor/participants",
                    "/v1/anchor/trust",
                    "/v1/client/accounts",
                    "/v1/client/assets",
                    "/v1/client/assets/accounts",
                    "/v1/client/assets/issued",
                    "/v1/client/assets/participants",
                    "/v1/client/balances/accounts",
                    "/v1/client/exchange",
                    "/v1/client/fees/request",
                    "/v1/client/fees/response",
                    "/v1/client/message",
                    "/v1/client/obligations",
                    "/v1/client/participants",
                    "/v1/client/participants/whitelist",
                    "/v1/client/payload/sign",
                    "/v1/client/payout",
                    "/v1/client/quotes",
                    "/v1/client/quotes/request",
                    "/v1/client/sign",
                    "/v1/client/token/refresh",
                    "/v1/client/transactions",
                    "/v1/client/transactions/redeem",
                    "/v1/client/transactions/reply",
                    "/v1/client/transactions/send",
                    "/v1/client/trust",
                    "/v1/onboarding/accounts"
                ],
                jti: '',
                aud: environment[participantID]
            },
            json: true
        };

        console.log(options);

        request(options, function(error, res, body) {
            if (error) {
                logRequest(err, options)
                done(err)
            } else {
                logRequest(res, options)
                should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                var str = body.split(": ");
                process.env[environment[participantID] + '_jti'] = str[1]
                done()
            }
        });

    });
});

When('user: {string} send request to approve a jwt-token for participant: {string} IID: {string} using totpkey: {string}, auth_url:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(userEmail, participantID, IID, totpkey, auth_url, done) {
    let fid = process.env[environment[userEmail] + '_FID']
    console.log(fid);

    let totpCode = speakeasy.totp({
        secret: environment[totpkey],
        encoding: 'ascii'
    });


    var options = {
        method: 'POST',
        url: environment[auth_url] + '/jwt/approve',
        headers: {
            'x-verify-code': totpCode,
            'x-fid': fid,
            'x-iid': environment[IID],
            'Content-Type': 'application/json'
        },
        body: { jti: process.env[environment[participantID] + '_jti'] },
        json: true
    };

    request(options, function(error, res, body) {
        if (error) {
            logRequest(err, options)
            done(err)
        } else {
            console.log(body);

            logRequest(res, options)
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            done()
        }
    });


});


Then('user: {string} send request to get a jwt-token for participant: {string} IID: {string} using totpkey: {string} naming as: {string} , auth_url:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(userEmail, participantID, IID, totpkey, name, auth_url, done) {
    setTimeout(function() {
        let fid = process.env[environment[userEmail] + '_FID']
        console.log(fid);

        let totpCode = speakeasy.totp({
            secret: environment[totpkey],
            encoding: 'ascii'
        });


        var options = {
            method: 'POST',
            url: environment[auth_url] + '/jwt/generate',
            headers: {
                'x-verify-code': totpCode,
                'x-fid': fid,
                'x-iid': environment[IID],
                'Content-Type': 'application/json'
            },
            body: { jti: process.env[environment[participantID] + '_jti'] },
            json: true
        };

        request(options, function(error, res, body) {
            if (error) {
                logRequest(err, options)
                done(err)
            } else {
                logRequest(res, options)
                process.env[name] = body

                should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                done()
            }
        });

    }, 3000)

});