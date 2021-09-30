// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const createDeliveryFailPayload_camt026 = require('../../../utility/createDeliveryFailPayload_camt026')

// var SignedXml = require('xml-crypto').SignedXml,
var fs = require('fs')
const encoder = require('nodejs-base64-encode');
createDeliveryFailPayload_camt026(
    'IBMQAIBM002',
    'ibm02',
    'IBMQAIBM001',
    'ibm01',
    'USDDO20200106IBMQAIBM001B2857664323',
    'USDDO',
    '50'
).then(function(msg) {
    console.log(msg);
    var xml = encoder.decode(msg, 'base64')

    // var sig = new SignedXml()
    // sig.addReference("//*[local-name(.)='book']")
    // sig.signingKey = "SAJHQPYYZIS6GYFVFSTPOMAVGJDYKXDABXXOXG67TDUNTYIIGJRNLZKL"
    // sig.canonicalizationAlgorithm = "http://www.w3.org/2001/10/xml-exc-c14n#WithComments"
    // sig.computeSignature(xml)
    // console.log(sig.getSignedXml());

})