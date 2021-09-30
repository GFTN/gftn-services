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
const createSendPayload_pacs008 = require('../../../utility/createSendPayload_pacs008')
const encoder = require('nodejs-base64-encode');
const parser = new xml2js.Parser({ attrkey: "ATTR" });
const sendReq = require('../../../utility/asyncSendReq')
var bigDecimal = require('js-big-decimal');

let payout_detail
let feeReqMsg, feeResMsg

Given('id: {string} query payout area by {string} - payoutURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, payoutquery, payoutURL, done) {
    var options = {
        method: 'GET',
        url: environment[payoutURL] + '/client/payout',
        headers: {}
    };
    options = appendToken(options, id, true)
    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            logRequest(res, options)
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            payout_detail = JSON.parse(body)[0]
            done()
        }
    });
});


Then('id:  {string} request fee with fee_request_id {string} receiver {string} sending_asset_code {string} sending_asset_type {string} sending_asset_issuer_id {string} amount_gross {string} asset_payout {string}, feeURL:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, fee_request_id, receiver_id, sending_asset_code, sending_asset_type, sending_asset_issuer_id, amount_gross, asset_payout, feeURL, done) {
    var options = {
        method: 'POST',
        url: environment[feeURL] + '/client/fees/request/' + environment[receiver_id],
        headers: {},
        body: {
            request_id: environment[fee_request_id],
            participant_id: environment[receiver_id],
            asset_settlement: {
                asset_type: environment[sending_asset_type],
                asset_code: environment[sending_asset_code],
                issuer_id: environment[sending_asset_issuer_id]
            },
            amount_payout: parseFloat(environment[amount_gross]),
            asset_payout: environment[asset_payout],
            details_payout_location: payout_detail
        },
        json: true
    };
    options = appendToken(options, id, true)
    process.env[environment[id] + "_REQUEST_FEE_ID_" + environment[sending_asset_code] + "_" + environment[sending_asset_issuer_id]] = environment[fee_request_id]
    process.env[environment[id] + "_REQUEST_PAYOUT_ID_" + environment[sending_asset_code] + "_" + environment[sending_asset_issuer_id]] = payout_detail.id

    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            logRequest(res, options)
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            should(body.status).be.exactly("Success", "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

            done()
        }
    });

});

Then('id: {string} pick up message from RFI_FEE topic sent by ofi {string} should get fee_request_id {string} receiver {string} sending_asset_code {string} sending_asset_type {string} sending_asset_issuer_id {string} amount_gross {string} asset_payout {string}, wwGatewayURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, ofi_participant_id, fee_request_id, receiver_id, sending_asset_code, sending_asset_type, sending_asset_issuer_id, amount_gross, asset_payout, wwGateWayURL, done) {
    setTimeout(async function() {
        var options = {
            method: 'GET',
            url: environment[wwGateWayURL] + '/client/message',
            headers: {},
            qs: { type: 'fee' }
        };

        options = appendToken(options, id)

        try {
            // console.log(options);

            let retry = environment.ENV_KEY_GATEWAY_RETRY_TIMES
            let res
            while (retry > 0) {
                res = await sendReq(options)
                    // console.log(JSON.parse(res.body))
                if (JSON.parse(res.body).data == null) {
                    retry--
                } else {
                    break
                }
            }
            let feeReqMsg = {
                request_id: environment[fee_request_id],
                participant_id: environment[receiver_id],
                asset_settlement: {
                    asset_type: environment[sending_asset_type],
                    asset_code: environment[sending_asset_code],
                    issuer_id: environment[sending_asset_issuer_id]
                },
                amount_payout: parseFloat(environment[amount_gross]),
                asset_payout: environment[asset_payout],
                details_payout_location: payout_detail
            }
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            JSON.parse(res.body).data.should.containDeep([feeReqMsg], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))

            done()
        } catch (err) {
            console.log(err);
            // done()

        }
    }, 2000)
});


Then('id: {string} caculate fee amount as: {string}, fee amount_payout as: {string} asset_payout {string} amount_settlement {string} sending_asset_code {string} sending_asset_type {string} sending_asset_issuer_id {string} and response fee message to OFI: {string} OFI_FEE topic, fee: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, fee_amount, amount_payout, asset_payout, amount_settlement, sending_asset_code, sending_asset_type, sending_asset_issuer_id, ofi_participantID, feeURL, done) {
    let amountFee = new bigDecimal(environment[fee_amount])
    let amountPayout = new bigDecimal(environment[amount_payout])
    let amountSettlement = new bigDecimal(environment[amount_settlement])
    var options = {
        method: 'POST',
        url: environment[feeURL] + '/client/fees/response/' + environment[ofi_participantID],
        headers: {},
        body: {
            amount_fee: parseFloat(amountFee.getValue()),
            amount_payout: parseFloat(amountPayout.getValue()),
            amount_settlement: parseFloat(amountSettlement.getValue()),
            asset_code_payout: environment[asset_payout],
            details_asset_settlement: {
                asset_code: environment[sending_asset_code],
                asset_type: environment[sending_asset_type],
                issuer_id: environment[sending_asset_issuer_id]
            },
            request_id: process.env[environment[ofi_participantID] + "_REQUEST_FEE_ID_" + environment[sending_asset_code] + "_" + environment[sending_asset_issuer_id]]
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
            should(body.status).be.exactly("Success", "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            done()
        }
    });

});


Then('id: {string} pick up message from OFI_FEE topic should get fee amount as: {string}, fee amount_payout as: {string} asset_payout {string} amount_settlement {string} sending_asset_code {string} sending_asset_type {string} sending_asset_issuer_id {string}, wwGatewayURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, fee_amount, amount_payout, asset_payout, amount_settlement, sending_asset_code, sending_asset_type, sending_asset_issuer_id, wwGateWayURL, done) {
    setTimeout(async function() {
        var options = {
            method: 'GET',
            url: environment[wwGateWayURL] + '/client/message',
            headers: {},
            qs: { type: 'fee' }
        };
        let amountFee = new bigDecimal(environment[fee_amount])
        let amountPayout = new bigDecimal(environment[amount_payout])
        let amountSettlement = new bigDecimal(environment[amount_settlement])

        feeResMsg = {
            amount_fee: parseFloat(amountFee.getValue()),
            amount_payout: parseFloat(amountPayout.getValue()),
            amount_settlement: parseFloat(amountSettlement.getValue()),
            asset_code_payout: environment[asset_payout],
            details_asset_settlement: {
                asset_code: environment[sending_asset_code],
                asset_type: environment[sending_asset_type],
                issuer_id: environment[sending_asset_issuer_id]
            },
            request_id: process.env[environment[id] + "_REQUEST_FEE_ID_" + environment[sending_asset_code] + "_" + environment[sending_asset_issuer_id]]
        }

        options = appendToken(options, id)

        try {

            let retry = environment.ENV_KEY_GATEWAY_RETRY_TIMES
            let res
            while (retry > 0) {
                res = await sendReq(options)
                if (JSON.parse(res.body).data == null) {
                    retry--
                } else {
                    break
                }
            }

            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            JSON.parse(res.body).data.should.containDeep([feeResMsg], "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            process.env[environment[id] + "_REQUEST_PAYOUT_ASSET_CODE_" + environment[sending_asset_code] + "_" + environment[sending_asset_issuer_id]] = feeResMsg.asset_code_payout
            process.env[environment[id] + "_REQUEST_PAYOUT_AMOUNT_" + environment[sending_asset_code] + "_" + environment[sending_asset_issuer_id]] = feeResMsg.amount_payout
            process.env[environment[id] + "_REQUEST_SETTLE_AMOUNT_" + environment[sending_asset_code] + "_" + environment[sending_asset_issuer_id]] = feeResMsg.amount_settlement
            process.env[environment[id] + "_REQUEST_FEE_ASSET_CODE_" + environment[sending_asset_code] + "_" + environment[sending_asset_issuer_id]] = feeResMsg.details_asset_settlement.asset_code
            process.env[environment[id] + "_REQUEST_FEE_AMOUNT_" + environment[sending_asset_code] + "_" + environment[sending_asset_issuer_id]] = feeResMsg.amount_fee
            done()

        } catch (err) {
            console.log(err);

        }


    }, 300)
});


Then('id: {string} send asset from sender_bic {string} sending_account_name {string} sending_bank_name {string} sending_street_name {string} sending_building_number {string} sending_post_code {string} sending_town_name {string} sending_country {string} with settlement_method {string} asset_code {string} asset_issuer {string} charger_bic {string} to receiver {string} recever_bic {string} receiver_bank_name {string} receiver_street_name {string} receiver_building_number {string} receiver_post_code {string} receiver_town_name {string} receiver_country {string} receiver_address_line {string}, sendURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id,
    sender_bic,
    sending_account_name,
    sending_bank_name,
    sending_street_name,
    sending_building_number,
    sending_post_code,
    sending_town_name,
    sending_country,
    settlement_method,
    send_asset_code,
    send_asset_issuer,
    charger_bic,
    receiver_id,
    recever_bic,
    receiver_bank_name,
    receiver_street_name,
    receiver_building_number,
    receiver_post_code,
    receiver_town_name,
    receiver_country,
    receiver_address_line,
    sendURL,
    done) {

    createSendPayload_pacs008(environment[settlement_method],
        environment[id],
        environment[sending_account_name],
        environment[send_asset_issuer],
        environment[sender_bic],
        environment[id],
        environment[recever_bic],
        environment[receiver_id],
        environment[send_asset_code],
        process.env[environment[id] + "_REQUEST_SETTLE_AMOUNT_" + environment[send_asset_code] + "_" + environment[send_asset_issuer]],
        1,
        // new bigDecimal(process.env[environment[id] + "_REQUEST_PAYOUT_AMOUNT_" + environment[send_asset_code] + "_" + environment[send_asset_issuer]]).divide(new bigDecimal(process.env[environment[id] + "_REQUEST_SETTLE_AMOUNT_" + environment[send_asset_code] + "_" + environment[send_asset_issuer]])).getValue(),
        process.env[environment[id] + "_REQUEST_FEE_ASSET_CODE_" + environment[send_asset_code] + "_" + environment[send_asset_issuer]],
        process.env[environment[id] + "_REQUEST_FEE_AMOUNT_" + environment[send_asset_code] + "_" + environment[send_asset_issuer]],
        environment[charger_bic],
        environment[sending_bank_name],
        environment[sending_street_name],
        environment[sending_building_number],
        environment[sending_post_code],
        environment[sending_town_name],
        environment[sending_country],
        environment[sender_bic],
        environment[recever_bic],
        environment[receiver_bank_name],
        environment[receiver_street_name],
        environment[receiver_building_number],
        environment[receiver_post_code],
        environment[receiver_town_name],
        environment[receiver_country],
        environment[receiver_address_line],
        // been add since v2.9.3.12_RC
        process.env[environment[id] + "_REQUEST_FEE_ID_" + environment[send_asset_code] + "_" + environment[send_asset_issuer]],
        process.env[environment[id] + "_REQUEST_PAYOUT_ID_" + environment[send_asset_code] + "_" + environment[send_asset_issuer]],
        process.env[environment[id] + "_REQUEST_PAYOUT_ASSET_CODE_" + environment[send_asset_code] + "_" + environment[send_asset_issuer]],
        process.env[environment[id] + "_REQUEST_PAYOUT_AMOUNT_" + environment[send_asset_code] + "_" + environment[send_asset_issuer]]
        // ---------------------------------------
    ).then(function(msg) {
        console.log(options);

        var options = {
            method: 'POST',
            url: environment[sendURL] + '/client/transactions/send',
            headers: {},
            body: {
                message_type: 'iso20022:pacs.008.001.07',
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
                logRequest(res, options)
                parser.parseString(encoder.decode(body.message, 'base64'), function(error, result) {
                    if (error === null) {

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



Then('id: {string} sign payload for send asset from sender_bic {string} sending_account_name {string} sending_bank_name {string} sending_street_name {string} sending_building_number {string} sending_post_code {string} sending_town_name {string} sending_country {string} with settlement_method {string} asset_code {string} asset_issuer {string} charger_bic {string} to receiver {string} recever_bic {string} receiver_bank_name {string} receiver_street_name {string} receiver_building_number {string} receiver_post_code {string} receiver_town_name {string} receiver_country {string} receiver_address_line {string}, cryptoURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id,
    sender_bic,
    sending_account_name,
    sending_bank_name,
    sending_street_name,
    sending_building_number,
    sending_post_code,
    sending_town_name,
    sending_country,
    settlement_method,
    send_asset_code,
    send_asset_issuer,
    charger_bic,
    receiver_id,
    recever_bic,
    receiver_bank_name,
    receiver_street_name,
    receiver_building_number,
    receiver_post_code,
    receiver_town_name,
    receiver_country,
    receiver_address_line,
    cryptoURL,
    done) {

    createSendPayload_pacs008(environment[settlement_method],
        environment[id],
        environment[sending_account_name],
        environment[send_asset_issuer],
        environment[sender_bic],
        environment[id],
        environment[recever_bic],
        environment[receiver_id],
        environment[send_asset_code],
        process.env[environment[id] + "_REQUEST_SETTLE_AMOUNT_" + environment[send_asset_code] + "_" + environment[send_asset_issuer]],
        1,
        // new bigDecimal(process.env[environment[id] + "_REQUEST_PAYOUT_AMOUNT_" + environment[send_asset_code] + "_" + environment[send_asset_issuer]]).divide(new bigDecimal(process.env[environment[id] + "_REQUEST_SETTLE_AMOUNT_" + environment[send_asset_code] + "_" + environment[send_asset_issuer]])).getValue(),
        process.env[environment[id] + "_REQUEST_FEE_ASSET_CODE_" + environment[send_asset_code] + "_" + environment[send_asset_issuer]],
        process.env[environment[id] + "_REQUEST_FEE_AMOUNT_" + environment[send_asset_code] + "_" + environment[send_asset_issuer]],
        environment[charger_bic],
        environment[sending_bank_name],
        environment[sending_street_name],
        environment[sending_building_number],
        environment[sending_post_code],
        environment[sending_town_name],
        environment[sending_country],
        environment[sender_bic],
        environment[recever_bic],
        environment[receiver_bank_name],
        environment[receiver_street_name],
        environment[receiver_building_number],
        environment[receiver_post_code],
        environment[receiver_town_name],
        environment[receiver_country],
        environment[receiver_address_line],
        // been add since v2.9.3.12_RC
        process.env[environment[id] + "_REQUEST_FEE_ID_" + environment[send_asset_code] + "_" + environment[send_asset_issuer]],
        process.env[environment[id] + "_REQUEST_PAYOUT_ID_" + environment[send_asset_code] + "_" + environment[send_asset_issuer]],
        process.env[environment[id] + "_REQUEST_PAYOUT_ASSET_CODE_" + environment[send_asset_code] + "_" + environment[send_asset_issuer]],
        process.env[environment[id] + "_REQUEST_PAYOUT_AMOUNT_" + environment[send_asset_code] + "_" + environment[send_asset_issuer]]
        // ---------------------------------------
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
                    // process.env[environment[id] + "_iso20022:pacs.008.001.07_PAYLOAD"] = encoder.encode(body, 'base64')
                process.env[environment[id] + "_iso20022:pacs.008.001.07_PAYLOAD"] = body.payload_with_signature

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





Then('id: {string} pick up transaction message , wwGatewayURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, wwGateWayURL, done) {
    setTimeout(async function() {
        var options = {
            method: 'GET',
            url: environment[wwGateWayURL] + '/client/message',
            headers: {},
            qs: { type: 'transactions' }
        };

        options = appendToken(options, id)

        try {
            // console.log(options);

            let retry = environment.ENV_KEY_GATEWAY_RETRY_TIMES
            let res
            while (retry > 0) {
                res = await sendReq(options)
                    // console.log(JSON.parse(res.body))
                if (JSON.parse(res.body).data == null) {
                    retry--
                } else {
                    break
                }
            }
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            JSON.parse(res.body).data.length.should.be.above(0, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            done()

        } catch (err) {
            console.log(err);
            // done()

        }
    }, 2000)
});



Then('id: {string} pick up payment message , wwGatewayURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, wwGateWayURL, done) {
    setTimeout(async function() {
        var options = {
            method: 'GET',
            url: environment[wwGateWayURL] + '/client/message',
            headers: {},
            qs: { type: 'payment' }
        };

        options = appendToken(options, id)

        try {
            // console.log(options);

            let retry = environment.ENV_KEY_GATEWAY_RETRY_TIMES
            let res
            while (retry > 0) {
                res = await sendReq(options)
                    // console.log(JSON.parse(res.body))
                if (JSON.parse(res.body).data == null) {
                    retry--
                } else {
                    break
                }
            }
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            JSON.parse(res.body).data.length.should.be.above(0, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            done()

        } catch (err) {
            console.log(err);
            // done()

        }
    }, 2000)
});


Then('id: {string} pick up rdo message , wwGatewayURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, wwGateWayURL, done) {
    setTimeout(async function() {
        var options = {
            method: 'GET',
            url: environment[wwGateWayURL] + '/client/message',
            headers: {},
            qs: { type: 'rdo' }
        };

        options = appendToken(options, id)

        try {
            // console.log(options);

            let retry = environment.ENV_KEY_GATEWAY_RETRY_TIMES
            let res
            while (retry > 0) {
                res = await sendReq(options)
                    // console.log(JSON.parse(res.body))
                if (JSON.parse(res.body).data == null) {
                    retry--
                } else {
                    break
                }
            }
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            JSON.parse(res.body).data.length.should.be.above(0, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            done()

        } catch (err) {
            console.log(err);
            // done()

        }
    }, 2000)
});