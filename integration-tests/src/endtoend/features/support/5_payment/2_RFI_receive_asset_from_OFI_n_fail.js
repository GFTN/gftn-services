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
const createReceivePayload_ibwf001 = require('../../../utility/createReceivePayload_ibwf001')
const encoder = require('nodejs-base64-encode');
const parser = new xml2js.Parser({ attrkey: "ATTR" });


Given('id: {string} get account_name: {string} address from WW, participant_api_url:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, account_name, url_api_service, done) {

    let options = {
        method: 'GET',
        headers: {},
        url: environment[url_api_service] + "/client/participants/" + environment[id]
    };

    options = appendToken(options, id)
    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            logRequest(res, options)
            if (environment[account_name] == 'issuing') {
                process.env[environment[id] + '_' + environment[account_name] + '_ADDRESS_'] = JSON.parse(body).issuing_account;
            } else {
                console.log(environment[account_name]);

                JSON.parse(body).operating_accounts.find(function(element) {
                    if (element.name == environment[account_name]) {
                        console.log(element.address);
                        process.env[environment[id] + '_' + environment[account_name] + '_ADDRESS_'] = element.address

                    }
                });

            }
            should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            should(process.env[environment[id] + '_' + environment[account_name] + '_ADDRESS_'].length).be.exactly(56)
                // process.env[environment[id] + '_' + environment[account_name] + '_ADDRESS_'] = JSON.parse(body).account.address;
            done()
        }
    });

});

Given('id: {string} check account_name: {string} address not exist in WW, participant_api_url:{string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, account_name, url_api_service, done) {



    let options = {
        method: 'GET',
        headers: {},
        url: environment[url_api_service] + "/onboarding/accounts/" + environment[account_name]
    };

    options = appendToken(options, id)
    request(options, function(err, res, body) {
        if (err) {
            logRequest(err, options)
            done(err)
        } else {
            logRequest(res, options)
            should(res.statusCode).be.exactly(404, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
            process.env[environment[id] + '_' + environment[account_name] + '_ADDRESS_'] = 'notExistAddress'
            done()
        }
    });


});


Then('id: {string} participant_bic: {string} finished federation and compliance check response to send_participant: {string} send_participant_bic: {string} with federation_status: {string} compliance_status: {string} compliance_status: {string} receive asset_code: {string} account_name: {string} settlement_method: {string}, sendURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, participant_bic, send_participant, send_participant_bic, federation_status, compliance_status_1, compliance_status_2, asset_code, account_name, settlement_method, sendURL, done) {
    let OFI_MSG_ID = process.env['OFI_MSG_ID_' + environment[send_participant]];
    let OFI_E2E_ID = process.env['OFI_E2E_ID_' + environment[send_participant]];
    let OFI_TX_ID = process.env['OFI_TX_ID_' + environment[send_participant]];
    let OFI_INSTR_ID = process.env['OFI_' + environment[send_participant] + '_ORI_INSTR_ID']
    let TX_CREATE_TIME = new Date().toISOString().replace(/\..+/, '')

    // let str = 'COMPLIANCEE2E' + environment[participant_bic]
    // let randomNumLen = 35 - str.length
    // let pow = Math.pow(10, (randomNumLen - 1))
    // let randomNum = Math.floor(Math.random() * ((9 * pow) - 1) + pow)
    // let COMPLIANCEE2E = str + randomNum.toString()
    let COMPLIANCEE2E = OFI_MSG_ID

    str = 'COMPLIANCETX' + environment[participant_bic]
    randomNumLen = 35 - str.length
    pow = Math.pow(10, (randomNumLen - 1))
    randomNum = Math.floor(Math.random() * ((9 * pow) - 1) + pow)
        // let COMPLIANCETX = str + randomNum.toString()
    let COMPLIANCEINSTR = OFI_INSTR_ID


    createReceivePayload_ibwf001(
        environment[asset_code],
        environment[participant_bic],
        TX_CREATE_TIME,
        environment[settlement_method],
        environment[id],
        environment[account_name],
        environment[participant_bic],
        environment[id],
        environment[send_participant_bic],
        environment[send_participant],
        process.env[environment[id] + '_' + environment[account_name] + '_ADDRESS_'],
        environment[federation_status],
        OFI_E2E_ID,
        OFI_INSTR_ID,
        environment[compliance_status_1],
        environment[compliance_status_2],
        COMPLIANCEE2E,
        COMPLIANCEINSTR).then(function(msg) {

        var options = {
            method: 'POST',
            url: environment[sendURL] + '/client/transactions/reply',
            headers: {},
            body: {
                message_type: 'iso20022:ibwf.001.001.01',
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
                        console.log(error);
                    }
                });
            }
        });


    })
});


Then('id: {string} sign participant_bic: {string} finished federation and compliance check response to send_participant: {string} send_participant_bic: {string} with federation_status: {string} compliance_status: {string} compliance_status: {string} receive asset_code: {string} account_name: {string} settlement_method: {string}, cryptoURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, participant_bic, send_participant, send_participant_bic, federation_status, compliance_status_1, compliance_status_2, asset_code, account_name, settlement_method, cryptoURL, done) {

    let OFI_MSG_ID = process.env['OFI_MSG_ID_' + environment[send_participant]];
    let OFI_E2E_ID = process.env['OFI_E2E_ID_' + environment[send_participant]];
    let OFI_TX_ID = process.env['OFI_TX_ID_' + environment[send_participant]];
    let OFI_INSTR_ID = process.env['OFI_' + environment[send_participant] + '_ORI_INSTR_ID']
    let TX_CREATE_TIME = new Date().toISOString().replace(/\..+/, '')

    // let str = 'COMPLIANCEE2E' + environment[participant_bic]
    // let randomNumLen = 35 - str.length
    // let pow = Math.pow(10, (randomNumLen - 1))
    // let randomNum = Math.floor(Math.random() * ((9 * pow) - 1) + pow)
    // let COMPLIANCEE2E = str + randomNum.toString()
    let COMPLIANCEE2E = OFI_MSG_ID

    str = 'COMPLIANCETX' + environment[participant_bic]
    randomNumLen = 35 - str.length
    pow = Math.pow(10, (randomNumLen - 1))
    randomNum = Math.floor(Math.random() * ((9 * pow) - 1) + pow)
        // let COMPLIANCETX = str + randomNum.toString()
    let COMPLIANCEINSTR = OFI_INSTR_ID


    createReceivePayload_ibwf001(
        environment[asset_code],
        environment[participant_bic],
        TX_CREATE_TIME,
        environment[settlement_method],
        environment[id],
        environment[account_name],
        environment[participant_bic],
        environment[id],
        environment[send_participant_bic],
        environment[send_participant],
        process.env[environment[id] + '_' + environment[account_name] + '_ADDRESS_'],
        environment[federation_status],
        OFI_E2E_ID,
        OFI_INSTR_ID,
        environment[compliance_status_1],
        environment[compliance_status_2],
        COMPLIANCEE2E,
        COMPLIANCEINSTR).then(function(msg) {

        var options = {
            method: 'POST',
            url: environment[cryptoURL] + '/client/payload/sign',
            headers: {},
            body: {
                account_name: environment[account_name],
                payload: msg
            },
            json: true
        };


        options = appendToken(options, id)


        // iso20022:ibwf.001.001.01
        request(options, function(err, res, body) {
            if (err) {
                logRequest(err, options)
                done(err)
            } else {

                logRequest(res, options)
                    // console.log(body);
                should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    // process.env[environment[id] + "_iso20022:ibwf.001.001.01_PAYLOAD"] = encoder.encode(body, 'base64')
                process.env[environment[id] + "_iso20022:ibwf.001.001.01_PAYLOAD"] = body.payload_with_signature
                    // console.log(process.env[environment[id] + "_iso20022:pacs.008.001.07_PAYLOAD"]);

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

});