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
const createSettleload_ibwf002 = require('../../../utility/createSettleload_ibwf002')
const encoder = require('nodejs-base64-encode');
const parser = new xml2js.Parser({ attrkey: "ATTR" });
var bigDecimal = require('js-big-decimal');
let settlePayload


Then('id: {string} send settlement message sending_amount {string} settle_amount {string} sender_bic {string} sending_account_name {string} sending_bank_name {string} sending_street_name {string} sending_building_number {string} sending_post_code {string} sending_town_name {string} sending_country {string} with settlement_method {string} asset_code {string} to receiver {string} recever_bic {string} receiver_bank_name {string} receiver_street_name {string} receiver_building_number {string} receiver_post_code {string} receiver_town_name {string} receiver_country {string}, sendURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id,
    sending_amount,
    settle_amount,
    sender_bic,
    sending_account_name,
    sending_bank_name,
    sending_street_name,
    sending_building_number,
    sending_post_code,
    sending_town_name,
    sending_country,
    settlement_method,
    asset_code,
    receiver,
    recever_bic,
    receiver_bank_name,
    receiver_street_name,
    receiver_building_number,
    receiver_post_code,
    receiver_town_name,
    receiver_country,
    sendURL,
    done) {
    createSettleload_ibwf002(
        environment[settlement_method],
        environment[id],
        environment[sending_account_name],
        environment[sender_bic],
        environment[id],
        environment[recever_bic],
        environment[receiver],
        environment[asset_code],
        environment[sending_amount],
        environment[asset_code],
        // environment[settle_amount],
        process.env[environment[id] + "_REQUEST_SETTLE_AMOUNT_" + environment[asset_code] + "_" + environment[id]],
        environment[sending_bank_name],
        environment[sending_street_name],
        environment[sending_building_number],
        environment[sending_post_code],
        environment[sending_town_name],
        environment[sending_country],
        environment[sender_bic],
        environment[receiver_bank_name],
        environment[receiver_street_name],
        environment[receiver_building_number],
        environment[receiver_post_code],
        environment[receiver_town_name],
        environment[receiver_country],
        'pacs.008.001.07',
        process.env['OFI_MSG_ID_' + environment[id]],
        process.env['OFI_' + environment[id] + '_ORI_INSTR_ID'],
        process.env['OFI_MSG_ID_' + environment[id]],
        process.env['OFI_' + environment[id] + '_ORI_TX_CREATE_DATE_TIME']
        // ORI_MSG_ID,
        // ORI_INSTR_ID,
        // ORI_END_TO_END_ID,
        // ORI_TX_CREATE_TIME
    ).then(function(msg) {

        var options = {
            method: 'POST',
            url: environment[sendURL] + '/client/transactions/send',
            headers: {},
            body: {
                message_type: 'iso20022:ibwf.002.001.01',
                message: msg
            },
            json: true
        };
        options = appendToken(options, id)
        request(options, function(err, res, body) {
            if (err) {
                logRequest(err, options)
                done(err)
            } else {

                parser.parseString(encoder.decode(body.message, 'base64'), function(error, result) {
                    // if (error === null) {

                    // logRequest(res, options)
                    // console.log(JSON.stringify(result));
                    // console.log(result.Document.FIToFIPmtStsRpt[0].TxInfAndSts[0].StsRsnInf[0].Rsn[0].Prtry);
                    // let msg = result.Document.FIToFIPmtStsRpt[0].TxInfAndSts[0].StsRsnInf[0].Rsn[0].Prtry[0]
                    // should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    // should(msg).be.exactly(environment[error_message], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

                    // done()
                    // } else {
                    logRequest(res, options)
                    should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    done()
                        // }
                });
            }
        });
    })
});




Then('id: {string} sign send settlement message sending_amount {string} settle_amount {string} sender_bic {string} sending_account_name {string} sending_bank_name {string} sending_street_name {string} sending_building_number {string} sending_post_code {string} sending_town_name {string} sending_country {string} with settlement_method {string} asset_code {string} to receiver {string} recever_bic {string} receiver_bank_name {string} receiver_street_name {string} receiver_building_number {string} receiver_post_code {string} receiver_town_name {string} receiver_country {string}, cryptoURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id,
    sending_amount,
    settle_amount,
    sender_bic,
    sending_account_name,
    sending_bank_name,
    sending_street_name,
    sending_building_number,
    sending_post_code,
    sending_town_name,
    sending_country,
    settlement_method,
    asset_code,
    receiver,
    recever_bic,
    receiver_bank_name,
    receiver_street_name,
    receiver_building_number,
    receiver_post_code,
    receiver_town_name,
    receiver_country,
    cryptoURL,
    done) {
    createSettleload_ibwf002(
        environment[settlement_method],
        environment[id],
        environment[sending_account_name],
        environment[sender_bic],
        environment[id],
        environment[recever_bic],
        environment[receiver],
        environment[asset_code],
        environment[sending_amount],
        environment[asset_code],
        // environment[settle_amount],
        process.env[environment[id] + "_REQUEST_SETTLE_AMOUNT_" + environment[asset_code] + "_" + environment[id]],
        environment[sending_bank_name],
        environment[sending_street_name],
        environment[sending_building_number],
        environment[sending_post_code],
        environment[sending_town_name],
        environment[sending_country],
        environment[sender_bic],
        environment[receiver_bank_name],
        environment[receiver_street_name],
        environment[receiver_building_number],
        environment[receiver_post_code],
        environment[receiver_town_name],
        environment[receiver_country],
        'pacs.008.001.07',
        process.env['OFI_MSG_ID_' + environment[id]],
        process.env['OFI_' + environment[id] + '_ORI_INSTR_ID'],
        process.env['OFI_MSG_ID_' + environment[id]],
        process.env['OFI_' + environment[id] + '_ORI_TX_CREATE_DATE_TIME']
        // ORI_MSG_ID,
        // ORI_INSTR_ID,
        // ORI_END_TO_END_ID,
        // ORI_TX_CREATE_TIME
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
        request(options, function(err, res, body) {
            if (err) {
                logRequest(err, options)
                done(err)
            } else {

                logRequest(res, options)
                should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                process.env[environment[id] + "_iso20022:ibwf.002.001.01_PAYLOAD"] =
                    body.payload_with_signature

                done()


                // parser.parseString(encoder.decode(body.message, 'base64'), function(error, result) {
                //     // if (error === null) {

                //     // logRequest(res, options)
                //     // console.log(JSON.stringify(result));
                //     // console.log(result.Document.FIToFIPmtStsRpt[0].TxInfAndSts[0].StsRsnInf[0].Rsn[0].Prtry);
                //     // let msg = result.Document.FIToFIPmtStsRpt[0].TxInfAndSts[0].StsRsnInf[0].Rsn[0].Prtry[0]
                //     // should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                //     // should(msg).be.exactly(environment[error_message], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

                //     // done()
                //     // } else {
                //     logRequest(res, options)
                //     should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                //     done()
                //         // }
                // });
            }
        });
    })
});


When('id: {string} response settle sender_bank_name: {string} sender_id: {string} sender_bic: {string} sending_account_name: {string} settlement_method:{string} sending_asset_code: {string} receiver_id: {string} receiver_bic: {string} receive_account_name: {string} settlement_reason: {string}, return return_asset_code: {string} return_asset_issuer: {string} settlement_info: {string}, sendURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id,
    sender_bank_name,
    sender_id,
    sender_bic,
    sending_account_name,
    settlement_method,
    sending_asset_code,
    receiver_id,
    receiver_bic,
    receive_account_name,
    settlement_reason,
    return_asset_code,
    return_asset_issuer,
    settlement_info,
    sendURL,
    done
) {

    setTimeout(function() {



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
            sum_send_fee,
            ori_tx_date,
            environment[return_asset_code],
            sum_send_fee,
            fee_asset_code,
            fee_amount.getValue(),
            environment[sender_bank_name],
            environment[settlement_reason],
            environment[settlement_info],
            environment[sending_account_name],
            environment[return_asset_issuer],
            1001,
            "Payment settlement accepted"
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
                            console.log(JSON.stringify(result));
                            // console.log(result.Document.FIToFIPmtStsRpt[0].TxInfAndSts[0].StsRsnInf[0].Rsn[0].Prtry);
                            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                            done()
                        } else {
                            console.log(error);
                        }
                    });
                }
            });

        })
    }, 2000)
});



When('id: {string} sign response settle sender_bank_name: {string} sender_id: {string} sender_bic: {string} sending_account_name: {string} settlement_method:{string} sending_asset_code: {string} receiver_id: {string} receiver_bic: {string} receive_account_name: {string} settlement_reason: {string}, return return_asset_code: {string} return_asset_issuer: {string} settlement_info: {string}, cryptoURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id,
    sender_bank_name,
    sender_id,
    sender_bic,
    sending_account_name,
    settlement_method,
    sending_asset_code,
    receiver_id,
    receiver_bic,
    receive_account_name,
    settlement_reason,
    return_asset_code,
    return_asset_issuer,
    settlement_info,
    cryptoURL,
    done
) {

    setTimeout(function() {



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
            sum_send_fee,
            ori_tx_date,
            environment[return_asset_code],
            sum_send_fee,
            fee_asset_code,
            fee_amount.getValue(),
            environment[sender_bank_name],
            environment[settlement_reason],
            environment[settlement_info],
            environment[sending_account_name],
            environment[return_asset_issuer],
            1001,
            "Payment settlement accepted"
        ).then(function(msg) {



            var options = {
                method: 'POST',
                url: environment[cryptoURL] + '/client/payload/sign',
                headers: {},
                body: {
                    account_name: environment[receive_account_name],
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
                        //         console.log(JSON.stringify(result));
                        //         // console.log(result.Document.FIToFIPmtStsRpt[0].TxInfAndSts[0].StsRsnInf[0].Rsn[0].Prtry);
                        //         should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                        //         done()
                        //     } else {
                        //         console.log(error);
                        //     }
                        // });
                }
            });

        })
    }, 2000)
});


Given('id: {string} using endtoend ID get transaction detail ofi_id: {string} rfi_id: {string} asset_code: {string} issuer_id: {string} sending_account_name: {string} receive_account_name: {string} status should be: {string}, participant_rdo_client_url: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, ofi_id, rfi_id, asset_code, issuer_id, sending_account_name, receive_account_name, status, participant_rdo_client_url, done) {
    // setTimeout(function() {
    var options = {
        method: 'GET',
        url: environment[participant_rdo_client_url] + '/callback/transactions/do/' + process.env['OFI_MSG_ID_' + environment[ofi_id]],
        headers: {}
    };

    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {

            logRequest(res, options)
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            JSON.parse(body).should.not.be.empty("response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            should(JSON.parse(body).asset_code).be.exactly(environment[asset_code])
            should(JSON.parse(body).issuer_id).be.exactly(environment[issuer_id])
            should(JSON.parse(body).ofi_id).be.exactly(environment[ofi_id])
            should(JSON.parse(body).ofi_sending_account_name).be.exactly(environment[sending_account_name])
            should(JSON.parse(body).ofi_sending_account_address).be.exactly(process.env[environment[ofi_id] + '_' + environment[sending_account_name] + '_ADDRESS_'])
            should(JSON.parse(body).rfi_id).be.exactly(environment[rfi_id])
            should(JSON.parse(body).rfi_receiving_account_address).be.exactly(process.env[environment[rfi_id] + '_' + environment[receive_account_name] + '_ADDRESS_'])
            should(JSON.parse(body).rfi_receiving_account_name).be.exactly(environment[receive_account_name])
            should(JSON.parse(body).status).be.exactly(environment[status])
            process.env[environment[id] + "_" + environment[sending_account_name] + "_" + environment[asset_code] + "_" + environment[issuer_id] + "_DO_SETTLE_AMOUNT"] = JSON.parse(body).amount
            settlePayload = JSON.parse(body)
            done()
        }
    });
    // }, 5000)
});


When('id: {string} settle payment from outside WW, participant_rdo_client_url: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, participant_rdo_client_url, done) {
    // setTimeout(function() {
    var options = {
        method: 'POST',
        url: environment[participant_rdo_client_url] + '/simulation/rdo/notify',
        headers: {},
        body: settlePayload,
        json: true
    };

    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            logRequest(res, options)
            done()
        }
    });
    // }, 5000)
});



Then('id: {string} settle payment from WW, participant_rdo_url: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, participant_rdo_url, done) {

    var options = {
        method: 'POST',
        url: environment[participant_rdo_url] + '/client/transactions/settle/do',
        headers: {},
        body: settlePayload,
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
            done()
        }
    });
    // }, 5000)
});