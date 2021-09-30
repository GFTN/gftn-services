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
const createCancelPayload_camt056 = require('../../../utility/createCancelPayload_camt056')
const encoder = require('nodejs-base64-encode');
const parser = new xml2js.Parser({ attrkey: "ATTR" });
var bigDecimal = require('js-big-decimal');

Given('id: {string} request cancel payment send_id: {string} sender_bic: {string} receiver_id: {string} receiver_bic: {string} asset_code: {string} settlement_method: {string} sending_account_name: {string}, send_service_url: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id,
    send_id,
    sender_bic,
    receiver_id,
    receiver_bic,
    asset_code,
    settlement_method,
    sending_account_name,
    sendURL,
    done) {

    let ori_instr_id = process.env['OFI_' + environment[send_id] + '_ORI_INSTR_ID']
    let OFI_MSG_ID = process.env['OFI_MSG_ID_' + environment[send_id]]
    let OFI_E2E_ID = process.env['OFI_E2E_ID_' + environment[send_id]]
    let OFI_TX_ID = process.env['OFI_TX_ID_' + environment[send_id]]
    let SEND_CREATE_DATE = process.env['OFI_' + environment[send_id] + '_CREATE_DATE']
    let CANCEL_REASON = 'Beneficiary information is wrong'

    let send_amount = new bigDecimal(process.env['OFI_' + environment[send_id] + '_SEND_AMOUNT']);
    let fee_amount = new bigDecimal(process.env['OFI_' + environment[send_id] + '_FEE_AMOUNT']);
    let sum_send_fee = send_amount.add(fee_amount).getValue()

    createCancelPayload_camt056(

        environment[sender_bic],
        environment[send_id],
        environment[receiver_bic],
        environment[receiver_id],
        environment[sender_bic],
        environment[send_id],
        ori_instr_id,
        OFI_MSG_ID,
        OFI_E2E_ID,
        OFI_TX_ID,
        'pacs.008.001.07',
        environment[asset_code],
        sum_send_fee,
        SEND_CREATE_DATE,
        CANCEL_REASON,
        environment[settlement_method],
        environment[sending_account_name]).then(function(msg) {

        var options = {
            method: 'POST',
            url: environment[sendURL] + '/client/transactions/send',
            headers: {},
            body: {
                message_type: 'iso20022:camt.056.001.08',
                message: msg
            },
            json: true
        };
        options = appendToken(options, id)

        // done()
        request(options, function(err, res, body) {
            if (err) {
                logRequest(err, options)
                done(err)
            } else {

                parser.parseString(encoder.decode(body.message, 'base64'), function(error, result) {
                    if (error === null) {

                        logRequest(res, options)
                        console.log(JSON.stringify(result));
                        // console.log(result.Document.FIToFIPmtStsRpt[0].TxInfAndSts[0].StsRsnInf[0].Rsn[0].Prtry);
                        should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

                        done()
                    } else {
                        logRequest(res, options)
                        should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                        done()
                    }
                });
            }
        });

    })
});



Given('id: {string} sign request cancel payment send_id: {string} sender_bic: {string} receiver_id: {string} receiver_bic: {string} asset_code: {string} settlement_method: {string} sending_account_name: {string}, cryptoURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id,
    send_id,
    sender_bic,
    receiver_id,
    receiver_bic,
    asset_code,
    settlement_method,
    sending_account_name,
    cryptoURL,
    done) {

    let ori_instr_id = process.env['OFI_' + environment[send_id] + '_ORI_INSTR_ID']
    let OFI_MSG_ID = process.env['OFI_MSG_ID_' + environment[send_id]]
    let OFI_E2E_ID = process.env['OFI_E2E_ID_' + environment[send_id]]
    let OFI_TX_ID = process.env['OFI_TX_ID_' + environment[send_id]]
    let SEND_CREATE_DATE = process.env['OFI_' + environment[send_id] + '_CREATE_DATE']
    let CANCEL_REASON = 'Beneficiary information is wrong'

    let send_amount = new bigDecimal(process.env['OFI_' + environment[send_id] + '_SEND_AMOUNT']);
    let fee_amount = new bigDecimal(process.env['OFI_' + environment[send_id] + '_FEE_AMOUNT']);
    let sum_send_fee = send_amount.add(fee_amount).getValue()



    createCancelPayload_camt056(

        environment[sender_bic],
        environment[send_id],
        environment[receiver_bic],
        environment[receiver_id],
        environment[sender_bic],
        environment[send_id],
        ori_instr_id,
        OFI_MSG_ID,
        OFI_E2E_ID,
        OFI_TX_ID,
        'pacs.008.001.07',
        environment[asset_code],
        sum_send_fee,
        SEND_CREATE_DATE,
        CANCEL_REASON,
        environment[settlement_method],
        environment[sending_account_name]).then(function(msg) {

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

        // done()
        request(options, function(err, res, body) {
            if (err) {
                logRequest(err, options)
                done(err)
            } else {
                logRequest(res, options)
                should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    // process.env[environment[id] + "_iso20022:camt.056.001.08_PAYLOAD"] = encoder.encode(body, 'base64')
                process.env[environment[id] + "_iso20022:camt.056.001.08_PAYLOAD"] = body.payload_with_signature

                done()

                // parser.parseString(encoder.decode(body.message, 'base64'), function(error, result) {
                //     if (error === null) {

                //         logRequest(res, options)
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