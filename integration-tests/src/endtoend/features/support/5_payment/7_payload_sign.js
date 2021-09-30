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


Then('id: {string} sending type: {string} signed payload, sendURL: {string}', {
    timeout: parseInt(environment.MAX_TIMEOUT)
}, function(id, filetype, sendURL, done) {
    // setTimeout(function() {
    switch (filetype) {
        case "iso20022:pacs.008.001.07":
            postURL = environment[sendURL] + '/client/transactions/send'
            msg = process.env[environment[id] + "_iso20022:pacs.008.001.07_PAYLOAD"]
            break;
        case "iso20022:ibwf.001.001.01":
            postURL = environment[sendURL] + '/client/transactions/reply'
            msg = process.env[environment[id] + "_iso20022:ibwf.001.001.01_PAYLOAD"]
            break;
        case "iso20022:camt.056.001.08":
            postURL = environment[sendURL] + '/client/transactions/send'
            msg = process.env[environment[id] + "_iso20022:camt.056.001.08_PAYLOAD"]
            break;
        case "iso20022:camt.029.001.09":
            postURL = environment[sendURL] + '/client/transactions/reply'
            msg = process.env[environment[id] + "_iso20022:camt.029.001.09_PAYLOAD"]
            break;
        case "iso20022:pacs.004.001.09":
            postURL = environment[sendURL] + '/client/transactions/reply'
            msg = process.env[environment[id] + "_iso20022:pacs.004.001.09_PAYLOAD"]
            console.log(environment[id] + "_iso20022:pacs.004.001.09_PAYLOAD")
            break;
        case "iso20022:ibwf.002.001.01":
            postURL = environment[sendURL] + '/client/transactions/send'
            msg = process.env[environment[id] + "_iso20022:ibwf.002.001.01_PAYLOAD"]
            break;
        case "iso20022:pacs.009.001.08":
            postURL = environment[sendURL] + '/client/transactions/redeem'
            msg = process.env[environment[id] + "_iso20022:pacs.009.001.08_PAYLOAD"]
            break;
        case "iso20022:camt.026.001.07":
            postURL = environment[sendURL] + '/client/transactions/send'
            msg = process.env[environment[id] + "_iso20022:camt.026.001.07_PAYLOAD"]
            break;
        default:
            break;
    }

    var options = {
        method: 'POST',
        url: postURL,
        headers: {},
        body: {
            message_type: filetype,
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
                    // logRequest(res, options)
                    should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                    done()
                }
            });
        }
    });
    // }, 20000)
});