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
const createCancelRejectPayload_camt029 = require('../../../utility/createCancelRejectPayload_camt029')
const encoder = require('nodejs-base64-encode');
const parser = new xml2js.Parser({ attrkey: "ATTR" });
var bigDecimal = require('js-big-decimal');

Given('id: {string} reject cancel payment send_asset: {string} sender_id: {string} sender_bic: {string} sending_account_name: {string} settlement_method:{string} receiver_id: {string} receiver_bic: {string} , sendURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(
    id,
    send_asset,
    sender_id,
    sender_bic,
    sending_account_name,
    settlement_method,
    receiver_id,
    receiver_bic,
    sendURL,
    done
) {
    // let receiver_account_address = process.env[environment[receiver_id] + '_' + environment[receive_account_name] + '_ADDRESS_']
    // let ORI_E2E_ID = process.env['OFI_E2E_ID_' + environment[sender_id]]
    let ORI_MSG_ID = process.env['OFI_MSG_ID_' + environment[sender_id]]
    let ori_tx_date_time = process.env['OFI_' + environment[sender_id] + '_ORI_TX_CREATE_DATE_TIME']
    let ORI_INSTR = process.env['OFI_' + environment[sender_id] + '_ORI_INSTR_ID']
    createCancelRejectPayload_camt029(
        environment[receiver_id],
        environment[receiver_bic],
        environment[sender_bic],
        environment[sender_id],
        ori_tx_date_time,
        ORI_MSG_ID,
        'pacs.008.001.07',
        environment[settlement_method],
        environment[sender_id],
        environment[sending_account_name],
        environment[send_asset],
        ORI_INSTR
    ).then(function(msg) {



        var options = {
            method: 'POST',
            url: environment[sendURL] + '/client/transactions/reply',
            headers: {},
            body: {
                message_type: 'iso20022:camt.029.001.09',
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
                logRequest(res, options)
                parser.parseString(encoder.decode(body.message, 'base64'), function(error, result) {
                    if (error === null) {
                        console.log(JSON.stringify(result));
                        // console.log(result.Document.FIToFIPmtStsRpt[0].TxInfAndSts[0].StsRsnInf[0].Rsn[0].Prtry);
                        should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

                        done()
                    } else {
                        // logRequest(res, options)
                        should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                        done()
                    }
                });
            }
        });

    })
});


Given('id: {string} sign reject cancel payment send_asset: {string} sender_id: {string} sender_bic: {string} sending_account_name: {string} settlement_method:{string} receiver_id: {string} receiver_bic: {string} , cryptoURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(
    id,
    send_asset,
    sender_id,
    sender_bic,
    sending_account_name,
    settlement_method,
    receiver_id,
    receiver_bic,
    cryptoURL,
    done
) {
    // let receiver_account_address = process.env[environment[receiver_id] + '_' + environment[receive_account_name] + '_ADDRESS_']
    // let ORI_E2E_ID = process.env['OFI_E2E_ID_' + environment[sender_id]]
    let ORI_MSG_ID = process.env['OFI_MSG_ID_' + environment[sender_id]]
    let ori_tx_date_time = process.env['OFI_' + environment[sender_id] + '_ORI_TX_CREATE_DATE_TIME']
    let ORI_INSTR = process.env['OFI_' + environment[sender_id] + '_ORI_INSTR_ID']
    createCancelRejectPayload_camt029(
        environment[receiver_id],
        environment[receiver_bic],
        environment[sender_bic],
        environment[sender_id],
        ori_tx_date_time,
        ORI_MSG_ID,
        'pacs.008.001.07',
        environment[settlement_method],
        environment[sender_id],
        environment[sending_account_name],
        environment[send_asset],
        ORI_INSTR
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
                    // process.env[environment[id] + "_iso20022:camt.029.001.09_PAYLOAD"] = encoder.encode(body, 'base64')
                process.env[environment[id] + "_iso20022:camt.029.001.09_PAYLOAD"] = body.payload_with_signature
                done()

                // parser.parseString(encoder.decode(body.message, 'base64'), function(error, result) {
                //     if (error === null) {
                //         console.log(JSON.stringify(result));
                //         // console.log(result.Document.FIToFIPmtStsRpt[0].TxInfAndSts[0].StsRsnInf[0].Rsn[0].Prtry);
                //         should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

                //         done()
                //     } else {
                //         // logRequest(res, options)
                //         should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                //         done()
                //     }
                // });
            }
        });

    })
});