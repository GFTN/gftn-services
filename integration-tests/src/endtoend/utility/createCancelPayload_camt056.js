// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const fs = require('fs');
const util = require('util');
const readFile = util.promisify(fs.readFile);
const encoder = require('nodejs-base64-encode');

module.exports = async(
    PAYMENT_SENDER_BIC,
    PAYMENT_SENDER_ID,
    PAYMENT_RECEIVER_BIC,
    PAYMENT_RECEIVER_ID,
    PAYMENT_CANCELER_BIC,
    PAYMENT_CANCELER_ID,
    ORI_INSTR_ID,
    ORI_MSG_ID,
    ORI_END_TO_END_ID,
    ORI_TX_ID,
    SEND_REQUEST_FILE_NAME,
    SEND_ASSET_CODE,
    SEND_AMOUNT_WITH_FEE,
    SEND_BANK_SETTLEMENT_DATE,
    CANCEL_REASON,
    SETTLE_METHOD,
    PAYMENT_SENDER_ACCOUNT_NAME
) => {
    try {

        let camt056 = await readFile('./file/camt056_template.xml', 'utf8')
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
        if (SEND_ASSET_CODE.length < 5) {
            AAAAA = SEND_ASSET_CODE + 'XX'
        } else {
            AAAAA = SEND_ASSET_CODE
        }
        let OFI_CANCEL_MSG_ID = AAAAA + YYYY + MM + DD + PAYMENT_SENDER_BIC + 'B' + randomNum

        let newcamt056 = camt056.replace("$OFI_CANCEL_MSG_ID", OFI_CANCEL_MSG_ID);

        let randomNumLen2 = 7
        let pow2 = Math.pow(10, (randomNumLen2 - 1))
        randomNum = Math.floor(Math.random() * ((9 * pow2) - 1) + pow2)

        let BUSINESS_MSG_ID = 'B' + YYYY + MM + DD + PAYMENT_SENDER_BIC + 'BAA' + randomNum.toString()

        // been add since v2.9.3.12_RC
        newcamt056 = newcamt056.replace("$HEADER_BIC", PAYMENT_SENDER_BIC);
        newcamt056 = newcamt056.replace("$HEADER_SENDER_ID", PAYMENT_SENDER_ID);
        newcamt056 = newcamt056.replace("$BUSINESS_MSG_ID", BUSINESS_MSG_ID);
        // newcamt056 = newcamt056.replace("$MSG_DEF_ID", END_TO_END_ID);
        newcamt056 = newcamt056.replace("$HEADER_TX_CREATE_TIME", TX_CREATE_TIME);
        newcamt056 = newcamt056.replace("$CASE_OFI_CANCEL_MSG_ID", OFI_CANCEL_MSG_ID);
        newcamt056 = newcamt056.replace("$PAYMENT_SENDER_BIC", PAYMENT_SENDER_BIC);
        newcamt056 = newcamt056.replace("$PAYMENT_SENDER_ID", PAYMENT_SENDER_ID);
        newcamt056 = newcamt056.replace("$PAYMENT_RECEIVER_BIC", PAYMENT_RECEIVER_BIC);
        newcamt056 = newcamt056.replace("$PAYMENT_RECEIVER_ID", PAYMENT_RECEIVER_ID);
        newcamt056 = newcamt056.replace("$TX_CREATE_TIME", TX_CREATE_TIME);
        newcamt056 = newcamt056.replace("$PAYMENT_CANCELER_BIC", PAYMENT_CANCELER_BIC);
        newcamt056 = newcamt056.replace("$PAYMENT_CANCELER_ID", PAYMENT_CANCELER_ID);
        newcamt056 = newcamt056.replace("$ORI_MSG_END_TO_END_ID", ORI_MSG_ID);
        newcamt056 = newcamt056.replace("$ORI_INSTR_ID", ORI_INSTR_ID);
        newcamt056 = newcamt056.replace("$ORI_END_TO_END_ID", ORI_END_TO_END_ID);
        newcamt056 = newcamt056.replace("$ORI_TX_ID", ORI_TX_ID);
        newcamt056 = newcamt056.replace("$SEND_REQUEST_FILE_NAME", SEND_REQUEST_FILE_NAME);
        newcamt056 = newcamt056.replace("$SEND_ASSET_CODE", SEND_ASSET_CODE);
        newcamt056 = newcamt056.replace("$SEND_BANK_SETTLEMENT_DATE", SEND_BANK_SETTLEMENT_DATE);
        newcamt056 = newcamt056.replace("$SEND_AMOUNT_WITH_FEE", SEND_AMOUNT_WITH_FEE);
        newcamt056 = newcamt056.replace("$CANCEL_REASON", CANCEL_REASON);
        newcamt056 = newcamt056.replace("$SETTLE_METHOD", SETTLE_METHOD);
        newcamt056 = newcamt056.replace("$ORI_PAYMENT_SENDER_ID", PAYMENT_SENDER_ID);
        newcamt056 = newcamt056.replace("$ORI_PAYMENT_SENDER_ACCOUNT_NAME", PAYMENT_SENDER_ACCOUNT_NAME);
        // ---------------------------------------
        console.log(newcamt056);

        let message = encoder.encode(newcamt056, 'base64')


        return message
    } catch (error) {
        return error
    }

}