// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const fs = require('fs');
const util = require('util');
const readFile = util.promisify(fs.readFile);
const encoder = require('nodejs-base64-encode');

module.exports = async(
    RECEIVER_ID,
    RECEIVER_BIC,
    PAYMENT_SENDER_BIC,
    PAYMENT_SENDER_ID,
    ORI_TX_CREATE_DATE_TIME,
    ORI_MSG_ID,
    SEND_REQUEST_FILE_NAME,
    SETTLE_METHOD,
    SENDER_ID,
    SENDER_ACCOUNT_NAME,
    ORI_OFI_SEND_ASSET_CODE,
    ORI_INSTR
) => {
    try {

        let camt029 = await readFile('./file/camt029_template.xml', 'utf8')
        let today = new Date();
        let DD = ('0' + today.getDate()).slice(-2);
        let MM = ('0' + (today.getMonth() + 1)).slice(-2);
        let YYYY = today.getFullYear();
        let AAAAA
        if (ORI_OFI_SEND_ASSET_CODE.length < 5) {
            AAAAA = ORI_OFI_SEND_ASSET_CODE + 'XX'
        } else {
            AAAAA = ORI_OFI_SEND_ASSET_CODE
        }
        let randomNumLen = 10
        let pow = Math.pow(10, (randomNumLen - 1))
        let randomNum = Math.floor(Math.random() * ((9 * pow) - 1) + pow)
        let TX_CREATE_TIME = new Date().toISOString().replace(/\..+/, '')
        let RFI_CANCEL_REJECT_MSG_ID = AAAAA + YYYY + MM + DD + RECEIVER_BIC + 'B' + randomNum


        let newcamt029 = camt029.replace("$MSG_ID", RFI_CANCEL_REJECT_MSG_ID);
        newcamt029 = newcamt029.replace("$RECEIVER_ID", RECEIVER_ID);;
        newcamt029 = newcamt029.replace("$RECEIVER_BIC", RECEIVER_BIC);
        newcamt029 = newcamt029.replace("$PAYMENT_SENDER_BIC", PAYMENT_SENDER_BIC);
        newcamt029 = newcamt029.replace("$PAYMENT_SENDER_ID", PAYMENT_SENDER_ID);
        newcamt029 = newcamt029.replace("$ORI_TX_CREATE_DATE_TIME", ORI_TX_CREATE_DATE_TIME);
        newcamt029 = newcamt029.replace("$ORI_MSG_ID", ORI_MSG_ID);
        newcamt029 = newcamt029.replace("$SEND_REQUEST_FILE_NAME", SEND_REQUEST_FILE_NAME);
        newcamt029 = newcamt029.replace("$SETTLE_METHOD", SETTLE_METHOD);
        newcamt029 = newcamt029.replace("$SENDER_ID", SENDER_ID);
        newcamt029 = newcamt029.replace("$SENDER_ACCOUNT_NAME", SENDER_ACCOUNT_NAME);
        newcamt029 = newcamt029.replace("$HEADER_BIC", RECEIVER_BIC);
        newcamt029 = newcamt029.replace("$HEADER_SENDER_ID", RECEIVER_ID);
        newcamt029 = newcamt029.replace("$ORI_INSTR", ORI_INSTR);


        let randomNumLen2 = 7
        let pow2 = Math.pow(10, (randomNumLen2 - 1))
        randomNum = Math.floor(Math.random() * ((9 * pow2) - 1) + pow2)

        let BUSINESS_MSG_ID = 'B' + YYYY + MM + DD + PAYMENT_SENDER_BIC + 'BAA' + randomNum.toString()

        newcamt029 = newcamt029.replace("$BUSINESS_MSG_ID", BUSINESS_MSG_ID);
        newcamt029 = newcamt029.replace("$MSG_DEF_ID", RFI_CANCEL_REJECT_MSG_ID);
        newcamt029 = newcamt029.replace("$HEADER_TX_CREATE_TIME", TX_CREATE_TIME);
        console.log(newcamt029);

        let message = encoder.encode(newcamt029, 'base64')


        return message
    } catch (error) {
        return error
    }

}