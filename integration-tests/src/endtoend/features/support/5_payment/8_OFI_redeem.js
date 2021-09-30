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
const xml2js = require('xml2js');
const logRequest = require('../../../utility/logRequest')
const appendToken = require('../../../utility/appendToken')
const createCancelSettlePayloadAgree_pacs004 = require('../../../utility/createCancelSettlePayloadAgree_pacs004')
const createRedeemRequest_pacs009 = require('../../../utility/createRedeemRequest_pacs009')
const encoder = require('nodejs-base64-encode');
const parser = new xml2js.Parser({ attrkey: "ATTR" });
var bigDecimal = require('js-big-decimal');
let settlePayload


Then('id: {string} sign payload for redeem asset from sender_bic {string} sending_account_name {string} amount {string} with settlement_method {string} asset_code {string} asset_issuer {string} to receiver {string} recever_bic {string}, cryptoURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id,
    sender_bic,
    sending_account_name,
    amount,
    settlement_method,
    asset_code,
    asset_issuer,
    receiver,
    recever_bic,
    cryptoURL,
    done) {
    createRedeemRequest_pacs009(
        environment[settlement_method],
        environment[id],
        environment[sending_account_name],
        environment[asset_issuer],
        environment[sender_bic],
        environment[id],
        environment[recever_bic],
        environment[receiver],
        environment[asset_code],
        environment[amount],
        environment[sender_bic]).then(function(msg) {

        var options = {
            method: 'POST',
            url: environment[cryptoURL] + '/client/payload/sign',
            headers: {},
            body: {
                account_name: environment[sending_account_name],
                payload: msg
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
                should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                process.env[environment[id] + "_iso20022:pacs.009.001.08_PAYLOAD"] = body.payload_with_signature

                done()
                    // parser.parseString(encoder.decode(body.message, 'base64'), function(error, result) {
                    //     if (error === null) {

                //         console.log(JSON.stringify(result));
                //         // console.log(result.Document.FIToFIPmtStsRpt[0].TxInfAndSts[0].StsRsnInf[0].Rsn[0].Prtry);
                //         should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

                //         done()
                //     } else {
                //         logRequest(res, options)
                //         should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                //         done()
                //     }
                // });
            }
        });
    })

});