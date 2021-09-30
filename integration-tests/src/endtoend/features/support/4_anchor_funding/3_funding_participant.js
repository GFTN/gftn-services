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
const fs = require('fs');
const environment = require('../../../environment/env')

const logRequest = require('../../../utility/logRequest')
const appendToken = require('../../../utility/appendToken')
var StellarSdk = require('stellar-sdk');
StellarSdk.Network.use(new StellarSdk.Network(environment.ENV_KEY_STELLAR_NETWORK));

let details_funding, instruction_unsigned

Then('id: {string} get instruction to fund account_name: {string} amount_funding: {string} anchor_id: {string} asset_code_issued: {string} end_to_end_id: {string} participant_id: {string} memo_transaction: {string}, anchor_service_url: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, account_name, amount_funding, anchor_id, asset_code_issued, end_to_end_id, participant_id, memo_transaction, anchor_service_url, done) {
    var options = {
        method: 'POST',
        url: environment[anchor_service_url] + '/anchor/fundings/instruction',
        headers: {},
        body: {
            account_name: environment[account_name],
            amount_funding: parseInt(environment[amount_funding]),
            anchor_id: environment[anchor_id],
            asset_code_issued: environment[asset_code_issued],
            end_to_end_id: environment[end_to_end_id],
            participant_id: environment[participant_id],
            memo_transaction: environment[memo_transaction]
        },
        json: true
    };
    options = appendToken(options, id, true)
    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            logRequest(res, options)

            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

            should(body.details_funding.account_name).be.exactly(environment[account_name])
            should(body.details_funding.amount_funding).be.exactly(parseInt(environment[amount_funding]))
            should(body.details_funding.anchor_id).be.exactly(environment[anchor_id])
            should(body.details_funding.asset_code_issued).be.exactly(environment[asset_code_issued])
            should(body.details_funding.end_to_end_id).be.exactly(environment[end_to_end_id])
            should(body.details_funding.participant_id).be.exactly(environment[participant_id])
            should(body.details_funding.memo_transaction).be.exactly(environment[memo_transaction])

            details_funding = JSON.stringify(body.details_funding)
            instruction_unsigned = body.instruction_unsigned
            done()
        }
    });
});


Then('id: {string} signed instruction and details_funding to fund account_name: {string} amount_funding: {string} anchor_id: {string} asset_code_issued: {string} end_to_end_id: {string} participant_id: {string} memo_transaction: {string} with anchor_seed: {string}, anchor_service_url: {string}', {
        timeout: parseInt(environment.MAX_TIMEOUT)
    },
    function(id, account_name, amount_funding, anchor_id, asset_code_issued, end_to_end_id, participant_id, memo_transaction, anchor_seed, anchor_service_url, done) {

        const source = StellarSdk.Keypair.fromSecret(environment[anchor_seed])

        let afrString = details_funding
        let buf = Buffer.from(afrString, 'ascii');
        let base64afrString = buf.toString('base64')
        let signereq = source.sign(base64afrString)
        let funding_signed = signereq.toString('base64')

        let transaction = new StellarSdk.Transaction(instruction_unsigned)
        transaction.sign(source)

        let instruction_signed = transaction.toEnvelope().toXDR('base64')
        var options = {
            method: 'POST',
            url: environment[anchor_service_url] + '/anchor/fundings/send',
            qs: {
                funding_signed: encodeURI(funding_signed),
                instruction_signed: encodeURI(instruction_signed)
            },
            headers: {},
            body: {
                account_name: environment[account_name],
                amount_funding: parseInt(environment[amount_funding]),
                anchor_id: environment[anchor_id],
                asset_code_issued: environment[asset_code_issued],
                end_to_end_id: environment[end_to_end_id],
                participant_id: environment[participant_id],
                memo_transaction: environment[memo_transaction]
            },
            json: true
        };

        options = appendToken(options, id, true)
        request(options, function(err, res, body) {
            if (err) {
                logRequest(err, options)
                done(err)
            } else {
                logRequest(res, options)
                should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                done()
            }
        });

    });