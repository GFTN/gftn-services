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
const createDeliveryFailPayload_camt026 = require('../../../utility/createDeliveryFailPayload_camt026')
const encoder = require('nodejs-base64-encode');
const parser = new xml2js.Parser({ attrkey: "ATTR" });
var bigDecimal = require('js-big-decimal');
let settlePayload


Then('id: {string} sign participant_bic: {string} not able to complete the payment from send_participant: {string} send_participant_bic: {string} with sending_asset_code: {string} sending_amount: {string} by signing_account: {string}, cryptoURL: {string}', {
        timeout: parseInt(environment.MAX_TIMEOUT)
    },
    function(id, participant_bic, receiver_id, receiver_bic, asset_code, sending_amount, signing_account_name, cryptoURL, done) {
        let OFI_INSTR_ID = process.env['OFI_' + environment[receiver_id] + '_ORI_INSTR_ID']
        createDeliveryFailPayload_camt026(
            environment[participant_bic],
            environment[id],
            environment[receiver_bic],
            environment[receiver_id],
            OFI_INSTR_ID,
            environment[asset_code],
            environment[sending_amount]
        ).then(function(msg) {

            var options = {
                method: 'POST',
                url: environment[cryptoURL] + '/client/payload/sign',
                headers: {},
                body: {
                    account_name: environment[signing_account_name],
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
                    process.env[environment[id] + "_iso20022:camt.026.001.07_PAYLOAD"] = body.payload_with_signature

                    done()
                }
            });
        })

    });