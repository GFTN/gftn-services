// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const fs = require('fs');
const util = require('util');
const readFile = util.promisify(fs.readFile);
const encoder = require('nodejs-base64-encode');

module.exports = async(
    SETTLE_METHOD,
    OFI_ID,
    OFI_ACCOUNT_NAME,
    INSTG_BIC,
    INSTG_ID,
    RECEIVER_BIC,
    RECEIVER_ID,
    SETTLE_ASSET_CODE,
    SETTLE_AMOUNT,
    ORI_SETTLE_ASSET_CODE,
    ORI_SETTLE_AMOUNT,
    SENDER_BANK_NAME,
    SENDER_STREET_NAME,
    SENDER_BUILDING_NUMBER,
    SENDER_POST_CODE,
    SENDER_TOWN_NAME,
    SENDER_COUNTRY,
    SENDER_BIC,
    RECEIVER_BANK_NAME,
    RECEIVER_STREET_NAME,
    RECEIVER_BUILDING_NUMBER,
    RECEIVER_POST_CODE,
    RECEIVER_TOWN_NAME,
    RECEIVER_COUNTRY,
    SEND_REQUEST_FILE_NAME,
    ORI_MSG_ID,
    ORI_INSTR_ID,
    ORI_END_TO_END_ID,
    ORI_TX_CREATE_TIME
) => {
    let ibwf002 = await readFile('./file/ibwf002_template.xml', 'utf8')

    let today = new Date();
    let DD = ('0' + today.getDate()).slice(-2);
    let MM = ('0' + (today.getMonth() + 1)).slice(-2);
    let YYYY = today.getFullYear();
    let AAAAA
    if (SETTLE_ASSET_CODE.length < 5) {
        AAAAA = SETTLE_ASSET_CODE + 'XX'
    } else {
        AAAAA = SETTLE_ASSET_CODE
    }
    let randomNumLen = 10
    let pow = Math.pow(10, (randomNumLen - 1))
    let randomNum = Math.floor(Math.random() * ((9 * pow) - 1) + pow)
    let OFI_MSG_ID = AAAAA + YYYY + MM + DD + INSTG_BIC + 'B' + randomNum
    let TX_CREATE_TIME = new Date().toISOString().replace(/\..+/, '')

    process.env['OFI_' + OFI_ID + '_ORI_TX_CREATE_DATE_TIME'] = TX_CREATE_TIME;

    let newibwf002 = ibwf002.replace("$MSG_ID", OFI_MSG_ID);
    newibwf002 = newibwf002.replace("$TX_CREATE_TIME", TX_CREATE_TIME);
    newibwf002 = newibwf002.replace("$SETTLE_METHOD", SETTLE_METHOD);
    newibwf002 = newibwf002.replace("$OFI_ID", OFI_ID);
    newibwf002 = newibwf002.replace("$OFI_ACCOUNT_NAME", OFI_ACCOUNT_NAME);
    newibwf002 = newibwf002.replace("$INSTG_BIC", INSTG_BIC);
    newibwf002 = newibwf002.replace("$INSTG_ID", INSTG_ID);
    newibwf002 = newibwf002.replace("$INSTD_BIC", RECEIVER_BIC);
    newibwf002 = newibwf002.replace("$INSTD_ID", RECEIVER_ID);
    newibwf002 = newibwf002.replace("$RECEIVER_BIC", RECEIVER_BIC);
    newibwf002 = newibwf002.replace("$RECEIVER_ID", RECEIVER_ID);

    newibwf002 = newibwf002.replace("$ORI_MSG_ID", ORI_MSG_ID);
    newibwf002 = newibwf002.replace("$ORI_TX_CREATE_TIME", ORI_TX_CREATE_TIME);
    newibwf002 = newibwf002.replace("$OFI_FINANCIAL_ID", OFI_MSG_ID);
    newibwf002 = newibwf002.replace("$ORI_INSTR_ID", ORI_INSTR_ID);
    newibwf002 = newibwf002.replace("$ORI_END_TO_END_ID", ORI_END_TO_END_ID);
    newibwf002 = newibwf002.replace("$ORI_TX_ID", ORI_END_TO_END_ID);

    newibwf002 = newibwf002.replace("$SETTLE_ASSET_CODE", SETTLE_ASSET_CODE);
    newibwf002 = newibwf002.replace("$SETTLE_AMOUNT", SETTLE_AMOUNT);
    newibwf002 = newibwf002.replace("$ORI_SETTLE_ASSET_CODE", ORI_SETTLE_ASSET_CODE);
    newibwf002 = newibwf002.replace("$ORI_SETTLE_AMOUNT", ORI_SETTLE_AMOUNT);

    newibwf002 = newibwf002.replace("$SENDER_BANK_NAME", SENDER_BANK_NAME);
    newibwf002 = newibwf002.replace("$SENDER_STREET_NAME", SENDER_STREET_NAME);
    newibwf002 = newibwf002.replace("$SENDER_BUILDING_NUMBER", SENDER_BUILDING_NUMBER);
    newibwf002 = newibwf002.replace("$SENDER_POST_CODE", SENDER_POST_CODE);
    newibwf002 = newibwf002.replace("$SENDER_TOWN_NAME", SENDER_TOWN_NAME);
    newibwf002 = newibwf002.replace("$SENDER_COUNTRY", SENDER_COUNTRY);

    newibwf002 = newibwf002.replace("$RECEIVER_BANK_NAME", RECEIVER_BANK_NAME);
    newibwf002 = newibwf002.replace("$RECEIVER_STREET_NAME", RECEIVER_STREET_NAME);
    newibwf002 = newibwf002.replace("$RECEIVER_BUILDING_NUMBER", RECEIVER_BUILDING_NUMBER);
    newibwf002 = newibwf002.replace("$RECEIVER_POST_CODE", RECEIVER_POST_CODE);
    newibwf002 = newibwf002.replace("$RECEIVER_TOWN_NAME", RECEIVER_TOWN_NAME);
    newibwf002 = newibwf002.replace("$RECEIVER_COUNTRY", RECEIVER_COUNTRY);
    newibwf002 = newibwf002.replace("$ORI_SEND_REQUEST_FILE_NAME", SEND_REQUEST_FILE_NAME);
    newibwf002 = newibwf002.replace("$PMT_PARTICIPANT_ID", OFI_ID);

    randomNumLen = 7
    pow = Math.pow(10, (randomNumLen - 1))
    randomNum = Math.floor(Math.random() * ((9 * pow) - 1) + pow)

    let BUSINESS_MSG_ID = 'B' + YYYY + MM + DD + SENDER_BIC + 'BAA' + randomNum.toString()
        // been add since v2.9.3.12_RC
    newibwf002 = newibwf002.replace("$HEADER_BIC", INSTG_BIC);
    newibwf002 = newibwf002.replace("$HEADER_SENDER_ID", OFI_ID);
    newibwf002 = newibwf002.replace("$BUSINESS_MSG_ID", BUSINESS_MSG_ID);
    newibwf002 = newibwf002.replace("$MSG_DEF_ID", OFI_MSG_ID);
    newibwf002 = newibwf002.replace("$HEADER_TX_CREATE_TIME", TX_CREATE_TIME);
    // ---------------------------------------

    console.log(newibwf002);

    let message = encoder.encode(newibwf002, 'base64')
    return message


}