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
const encoder = require('nodejs-base64-encode');
const parser = new xml2js.Parser({ attrkey: "ATTR" });
var bigDecimal = require('js-big-decimal');

When('id: {string} agree cancel payment sender_bank_name: {string} sender_id: {string} sender_bic: {string} sending_account_name: {string} settlement_method:{string} sending_asset_code: {string} receiver_id: {string} receiver_bic: {string} receive_account_name: {string} cancel_reason: {string}, return return_asset_code: {string} return_asset_issuer: {string} return_asset_amount: {string} cancel_agree_info: {string}, sendURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(
    id,
    sender_bank_name,
    sender_id,
    sender_bic,
    sending_account_name,
    settlement_method,
    sending_asset_code,
    receiver_id,
    receiver_bic,
    receive_account_name,
    cancel_reason,
    return_asset_code,
    return_asset_issuer,
    return_asset_amount,
    cancel_agree_info,
    sendURL,
    done
) {
    // let receiver_account_address = process.env[environment[receiver_id] + '_' + environment[receive_account_name] + '_ADDRESS_']

    let OFI_MSG_ID = process.env['OFI_MSG_ID_' + environment[sender_id]]
    let OFI_E2E_ID = process.env['OFI_E2E_ID_' + environment[sender_id]]
    let OFI_TX_ID = process.env['OFI_TX_ID_' + environment[sender_id]]
    let ori_tx_date = process.env['OFI_' + environment[sender_id] + '_CREATE_DATE']
    let ori_instr_id = process.env['OFI_' + environment[sender_id] + '_ORI_INSTR_ID']
    let ori_tx_date_time = process.env['OFI_' + environment[sender_id] + '_ORI_TX_CREATE_DATE_TIME']
    let fee_asset_code = process.env[environment[sender_id] + "_REQUEST_FEE_ASSET_CODE_" + environment[return_asset_code] + "_" + environment[return_asset_issuer]]


    let settle_amount = new bigDecimal(process.env[environment[sender_id] + "_REQUEST_SETTLE_AMOUNT_" + environment[return_asset_code] + "_" + environment[return_asset_issuer]]);
    let fee_amount = new bigDecimal(process.env[environment[sender_id] + "_REQUEST_FEE_AMOUNT_" + environment[return_asset_code] + "_" + environment[return_asset_issuer]]);


    let sum_send_fee = settle_amount.getValue()


    createCancelSettlePayloadAgree_pacs004(
        environment[settlement_method],
        environment[receiver_id],
        environment[receive_account_name],
        environment[receiver_bic],
        environment[sender_bic],
        environment[sender_id],
        ori_instr_id,
        OFI_MSG_ID,
        OFI_E2E_ID,
        OFI_TX_ID,
        'pacs.008.001.07',
        ori_tx_date_time,
        environment[sending_asset_code],
        sum_send_fee,
        environment[return_asset_code],
        environment[return_asset_amount],
        ori_tx_date,
        environment[return_asset_code],
        environment[return_asset_amount],
        fee_asset_code,
        fee_amount.getValue(),
        environment[sender_bank_name],
        environment[cancel_reason],
        environment[cancel_agree_info],
        environment[sending_account_name],
        environment[return_asset_issuer],
        1000,
        "Payment cancellation accepted"
    ).then(function(msg) {



        var options = {
            method: 'POST',
            url: environment[sendURL] + '/client/transactions/reply',
            headers: {},
            body: {
                message_type: 'iso20022:pacs.004.001.09',
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


When('id: {string} sign agree cancel payment sender_bank_name: {string} sender_id: {string} sender_bic: {string} sending_account_name: {string} settlement_method:{string} sending_asset_code: {string} receiver_id: {string} receiver_bic: {string} receive_account_name: {string} cancel_reason: {string}, return return_asset_code: {string} return_asset_issuer: {string} return_asset_amount: {string} cancel_agree_info: {string}, cryptoURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(
    id,
    sender_bank_name,
    sender_id,
    sender_bic,
    sending_account_name,
    settlement_method,
    sending_asset_code,
    receiver_id,
    receiver_bic,
    receive_account_name,
    cancel_reason,
    return_asset_code,
    return_asset_issuer,
    return_asset_amount,
    cancel_agree_info,
    cryptoURL,
    done
) {
    // let receiver_account_address = process.env[environment[receiver_id] + '_' + environment[receive_account_name] + '_ADDRESS_']

    let OFI_MSG_ID = process.env['OFI_MSG_ID_' + environment[sender_id]]
    let OFI_E2E_ID = process.env['OFI_E2E_ID_' + environment[sender_id]]
    let OFI_TX_ID = process.env['OFI_TX_ID_' + environment[sender_id]]
    let ori_tx_date = process.env['OFI_' + environment[sender_id] + '_CREATE_DATE']
    let ori_instr_id = process.env['OFI_' + environment[sender_id] + '_ORI_INSTR_ID']
    let ori_tx_date_time = process.env['OFI_' + environment[sender_id] + '_ORI_TX_CREATE_DATE_TIME']
    let fee_asset_code = process.env[environment[sender_id] + "_REQUEST_FEE_ASSET_CODE_" + environment[return_asset_code] + "_" + environment[return_asset_issuer]]


    let settle_amount = new bigDecimal(process.env[environment[sender_id] + "_REQUEST_SETTLE_AMOUNT_" + environment[return_asset_code] + "_" + environment[return_asset_issuer]]);
    let fee_amount = new bigDecimal(process.env[environment[sender_id] + "_REQUEST_FEE_AMOUNT_" + environment[return_asset_code] + "_" + environment[return_asset_issuer]]);


    let sum_send_fee = settle_amount.getValue()


    createCancelSettlePayloadAgree_pacs004(
        environment[settlement_method],
        environment[receiver_id],
        environment[receive_account_name],
        environment[receiver_bic],
        environment[sender_bic],
        environment[sender_id],
        ori_instr_id,
        OFI_MSG_ID,
        OFI_E2E_ID,
        OFI_TX_ID,
        'pacs.008.001.07',
        ori_tx_date_time,
        environment[sending_asset_code],
        sum_send_fee,
        environment[return_asset_code],
        environment[return_asset_amount],
        ori_tx_date,
        environment[return_asset_code],
        environment[return_asset_amount],
        fee_asset_code,
        fee_amount.getValue(),
        environment[sender_bank_name],
        environment[cancel_reason],
        environment[cancel_agree_info],
        environment[sending_account_name],
        environment[return_asset_issuer],
        1000,
        "Payment cancellation accepted"
    ).then(function(msg) {



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
                    // process.env[environment[id] + "_iso20022:pacs.004.001.09_PAYLOAD"] = encoder.encode(body, 'base64')
                process.env[environment[id] + "_iso20022:pacs.004.001.09_PAYLOAD"] = body.payload_with_signature
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