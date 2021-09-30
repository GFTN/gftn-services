// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const fs = require('fs');
const util = require('util');
const readFile = util.promisify(fs.readFile);
const encoder = require('nodejs-base64-encode');

module.exports = async(
    CANCEL_AGREE_SETTLE_METHOD,
    RECEIVER_ID,
    RECEIVER_ACCOUNT_NAME,
    PAYMENT_RECEIVER_BIC,
    PAYMENT_SENDER_BIC,
    PAYMENT_SENDER_ID,
    ORI_INSTR_ID,
    ORI_MSG_ID,
    OFI_E2E_ID,
    ORI_TX_ID,
    SEND_REQUEST_FILE_NAME,
    ORI_TX_CREATE_DATE_TIME,
    ORI_OFI_SEND_ASSET_CODE,
    ORI_OFI_SEND_ASSET_AMOUNT,
    RETURN_ASSET_CODE,
    RETURN_ASSET_AMOUNT,
    ORI_TX_CREATE_DATE,
    REFUNDED_OFI_CLIENT_ASSET_CODE,
    REFUNDED_OFI_CLIENT_ASSET_AMOUNT,
    FEE_ASSET_CODE,
    FEE_AMOUNT,
    SENDER_BANK_NAME,
    REPKY_REASON,
    REPLY_AGREE_INFO,
    SENDER_ACCOUNT_NAME,
    RETURN_ASSET_ISSUER_ID,
    RETURN_REASON_INFO_NUMBER,
    RETURN_REASON_INFO
) => {
    try {

        let pac004 = await readFile('./file/pacs004_template.xml', 'utf8')
        let today = new Date();
        let DD = ('0' + today.getDate()).slice(-2);
        let MM = ('0' + (today.getMonth() + 1)).slice(-2);
        let YYYY = today.getFullYear();


        let randomNumLen = 10
        let pow = Math.pow(10, (randomNumLen - 1))
        let randomNum = Math.floor(Math.random() * ((9 * pow) - 1) + pow)
        let TX_CREATE_TIME = new Date().toISOString().replace(/\..+/, '')
        let AAAAA
        if (ORI_OFI_SEND_ASSET_CODE.length < 5) {
            AAAAA = ORI_OFI_SEND_ASSET_CODE + 'XX'
        } else {
            AAAAA = ORI_OFI_SEND_ASSET_CODE
        }
        let RFI_CANCEL_AGREE_MSG_ID = AAAAA + YYYY + MM + DD + PAYMENT_RECEIVER_BIC + 'B' + randomNum


        let newpac004 = pac004.replace("$MSG_ID", RFI_CANCEL_AGREE_MSG_ID);
        newpac004 = newpac004.replace("$TX_CREATE_TIME", TX_CREATE_TIME);
        newpac004 = newpac004.replace("$CANCEL_AGREE_SETTLE_METHOD", CANCEL_AGREE_SETTLE_METHOD);
        newpac004 = newpac004.replace("$RECEIVER_ID", RECEIVER_ID);;
        newpac004 = newpac004.replace("$RECEIVER_ACCOUNT_NAME", RECEIVER_ACCOUNT_NAME);
        newpac004 = newpac004.replace("$PAYMENT_RECEIVER_BIC", PAYMENT_RECEIVER_BIC);
        newpac004 = newpac004.replace("$PAYMENT_RECEIVER_ID", RECEIVER_ID);
        newpac004 = newpac004.replace("$PAYMENT_SENDER_BIC", PAYMENT_SENDER_BIC);
        newpac004 = newpac004.replace("$PAYMENT_SENDER_ID", PAYMENT_SENDER_ID);
        newpac004 = newpac004.replace("$ORI_MSG_ID", ORI_MSG_ID);
        newpac004 = newpac004.replace("$SEND_REQUEST_FILE_NAME", SEND_REQUEST_FILE_NAME);
        newpac004 = newpac004.replace("$ORI_TX_CREATE_DATE_TIME", ORI_TX_CREATE_DATE_TIME);
        newpac004 = newpac004.replace("$TX_MSG_ID", RFI_CANCEL_AGREE_MSG_ID);
        newpac004 = newpac004.replace("$ORI_INSTR_ID", ORI_INSTR_ID);
        newpac004 = newpac004.replace("$OFI_E2E_ID", OFI_E2E_ID);
        newpac004 = newpac004.replace("$ORI_TX_ID", ORI_TX_ID);
        newpac004 = newpac004.replace("$ORI_OFI_SEND_ASSET_CODE", ORI_OFI_SEND_ASSET_CODE);
        newpac004 = newpac004.replace("$ORI_OFI_SEND_ASSET_AMOUNT", ORI_OFI_SEND_ASSET_AMOUNT);
        newpac004 = newpac004.replace("$RETURN_ASSET_CODE", RETURN_ASSET_CODE);
        newpac004 = newpac004.replace("$RETURN_ASSET_AMOUNT", RETURN_ASSET_AMOUNT);
        newpac004 = newpac004.replace("$ORI_TX_CREATE_DATE", ORI_TX_CREATE_DATE);
        newpac004 = newpac004.replace("$REFUNDED_OFI_CLIENT_ASSET_CODE", REFUNDED_OFI_CLIENT_ASSET_CODE);
        newpac004 = newpac004.replace("$REFUNDED_OFI_CLIENT_ASSET_AMOUNT", REFUNDED_OFI_CLIENT_ASSET_AMOUNT);
        newpac004 = newpac004.replace("$FEE_ASSET_CODE", FEE_ASSET_CODE);
        newpac004 = newpac004.replace("$FEE_AMOUNT", FEE_AMOUNT);
        newpac004 = newpac004.replace("$CHARGER_BIC", PAYMENT_RECEIVER_BIC);
        newpac004 = newpac004.replace("$CHARGER_ID", RECEIVER_ID);
        newpac004 = newpac004.replace("$SENDER_BANK_NAME", SENDER_BANK_NAME);
        newpac004 = newpac004.replace("$REPKY_REASON", REPKY_REASON);
        newpac004 = newpac004.replace("$REPLY_AGREE_INFO", REPLY_AGREE_INFO);
        newpac004 = newpac004.replace("$SETTLE_METHOD", CANCEL_AGREE_SETTLE_METHOD);
        newpac004 = newpac004.replace("$SENDER_ID", PAYMENT_SENDER_ID);
        newpac004 = newpac004.replace("$SENDER_ACCOUNT_NAME", SENDER_ACCOUNT_NAME);
        newpac004 = newpac004.replace("$RETURN_ASSET_ISSUER_ID", RETURN_ASSET_ISSUER_ID);
        newpac004 = newpac004.replace("$HEADER_BIC", PAYMENT_RECEIVER_BIC);
        newpac004 = newpac004.replace("$HEADER_SENDER_ID", RECEIVER_ID);
        newpac004 = newpac004.replace("$RETURN_REASON_INFO_NUMBER", RETURN_REASON_INFO_NUMBER);
        newpac004 = newpac004.replace("$RETURN_REASON_INFO", RETURN_REASON_INFO);

        let randomNumLen2 = 7
        let pow2 = Math.pow(10, (randomNumLen2 - 1))
        randomNum = Math.floor(Math.random() * ((9 * pow2) - 1) + pow2)

        let BUSINESS_MSG_ID = 'B' + YYYY + MM + DD + PAYMENT_SENDER_BIC + 'BAA' + randomNum.toString()

        newpac004 = newpac004.replace("$BUSINESS_MSG_ID", BUSINESS_MSG_ID);
        newpac004 = newpac004.replace("$HEADER_TX_CREATE_TIME", TX_CREATE_TIME);
        console.log(newpac004);

        let message = encoder.encode(newpac004, 'base64')


        return message
    } catch (error) {
        console.log(error);

        throw error
    }

}