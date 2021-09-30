// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const fs = require('fs');
const util = require('util');
const readFile = util.promisify(fs.readFile);
const encoder = require('nodejs-base64-encode');

module.exports = async(
    TX_SENDER_BIC,
    TX_SENDER_ID,
    TX_RECEIVER_BIC,
    TX_RECEIVER_ID,
    ORI_INSTR_ID,
    ORI_CURRENCY,
    ORI_AMOUNT) => {
    try {

        let camt026 = await readFile('./file/camt026_template.xml', 'utf8')
        let today = new Date();
        let DD = ('0' + today.getDate()).slice(-2);
        let MM = ('0' + (today.getMonth() + 1)).slice(-2);
        let YYYY = today.getFullYear();

        // let OFI_CANCEL_MSG_ID = SEND_ASSET_CODE + DD + MM + YYYY + PAYMENT_SENDER_BIC
        let randomNumLen = 10
        let pow = Math.pow(10, (randomNumLen - 1))
        let randomNum = Math.floor(Math.random() * ((9 * pow) - 1) + pow)
        let TX_CREATE_TIME = new Date().toISOString().replace(/\..+/, '')
        let AAAAA
        if (ORI_CURRENCY.length < 5) {
            AAAAA = ORI_CURRENCY + 'XX'
        } else {
            AAAAA = ORI_CURRENCY
        }
        let MSG_ID = AAAAA + YYYY + MM + DD + TX_SENDER_BIC + 'B' + randomNum

        let randomNumLen2 = 7
        let pow2 = Math.pow(10, (randomNumLen2 - 1))
        randomNum = Math.floor(Math.random() * ((9 * pow2) - 1) + pow2)

        let BUSINESS_MSG_ID = 'B' + YYYY + MM + DD + TX_SENDER_BIC + 'BAA' + randomNum.toString()

        let newcamt026 = camt026.replace("$HEADER_BIC", TX_SENDER_BIC);
        newcamt026 = newcamt026.replace("$HEADER_SENDER_ID", TX_SENDER_ID);
        newcamt026 = newcamt026.replace("$BUSINESS_MSG_ID", BUSINESS_MSG_ID);
        newcamt026 = newcamt026.replace("$HEADER_TX_CREATE_TIME", TX_CREATE_TIME);

        newcamt026 = newcamt026.replace("$MSG_ID", MSG_ID);
        newcamt026 = newcamt026.replace("$ASSIGNER_BIC", TX_SENDER_BIC);
        newcamt026 = newcamt026.replace("$ASSIGNER_ID", TX_SENDER_ID);
        newcamt026 = newcamt026.replace("$ASSIGNEE_BIC", TX_RECEIVER_BIC);
        newcamt026 = newcamt026.replace("$ADDIGNEE_ID", TX_RECEIVER_ID);

        newcamt026 = newcamt026.replace("$TX_CREATE_DATE_TIME", TX_CREATE_TIME);
        newcamt026 = newcamt026.replace("$ORI_INSTR_ID", ORI_INSTR_ID);
        newcamt026 = newcamt026.replace("$ORI_CURRENCY", ORI_CURRENCY);
        newcamt026 = newcamt026.replace("$ORI_AMOUNT", ORI_AMOUNT);
        console.log(newcamt026);

        let message = encoder.encode(newcamt026, 'base64')


        return message
    } catch (error) {
        return error
    }

}