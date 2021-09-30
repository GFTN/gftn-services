// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const fs = require('fs');
const util = require('util');
const readFile = util.promisify(fs.readFile);
const encoder = require('nodejs-base64-encode');

module.exports = async(
    SEND_ASSET_CODE,
    RFI_BIC,
    RFI_ID,
    OFI_BIC,
    OFI_ID,
    ORI_INSTR_ID,
    ORI_END_TO_END_ID,
    ORI_TX_ID,
    CHARGS_ASSET_CODE,
    CHARGS_AMOUNT,
    ORI_SETTLE_ASSET_CODE,
    ORI_SETTLE_AMOUNT,
    RECEIVE_ACCOUNT_ADDRESS
) => {
    let pac002 = await readFile('./file/pacs002_template.xml', 'utf8')
    let today = new Date();
    let DD = ('0' + today.getDate()).slice(-2);
    let MM = ('0' + (today.getMonth() + 1)).slice(-2);
    let YYYY = today.getFullYear();

    let RFI_MSG_ID = SEND_ASSET_CODE + DD + MM + YYYY + RFI_BIC
    let randomNumLen = 35 - RFI_MSG_ID.length
    let pow = Math.pow(10, (randomNumLen - 1))
    let randomNum = Math.floor(Math.random() * ((9 * pow) - 1) + pow)

    RFI_MSG_ID = RFI_MSG_ID + randomNum.toString()


    let randomNumLen2 = 7
    let pow2 = Math.pow(10, (randomNumLen2 - 1))
    randomNum = Math.floor(Math.random() * ((9 * pow2) - 1) + pow2)
    let TX_CREATE_TIME = new Date().toISOString().replace(/\..+/, '')
    process.env['OFI_' + OFI_ID + '_ORI_TX_CREATE_DATE_TIME'] = TX_CREATE_TIME;

    let BUSINESS_MSG_ID = 'B' + YYYY + MM + DD + OFI_BIC + 'BAA' + randomNum.toString()
    let newpac002 = pac002.replace("$HEADER_BIC", RFI_BIC);
    newpac002 = newpac002.replace("$HEADER_SENDER_ID", RFI_ID);
    newpac002 = newpac002.replace("$BUSINESS_MSG_ID", BUSINESS_MSG_ID);
    newpac002 = newpac002.replace("$MSG_DEF_ID", RFI_MSG_ID);
    newpac002 = newpac002.replace("$HEADER_TX_CREATE_TIME", TX_CREATE_TIME);

    newpac002 = newpac002.replace("$MSG_ID", RFI_MSG_ID);
    newpac002 = newpac002.replace("$TX_CREATE_TIME", TX_CREATE_TIME);
    newpac002 = newpac002.replace("$INSTG_BIC", RFI_BIC);
    newpac002 = newpac002.replace("$INSTG_ID", RFI_ID);
    newpac002 = newpac002.replace("$INSTD_BIC", OFI_BIC);
    newpac002 = newpac002.replace("$INSTD_ID", OFI_ID);
    newpac002 = newpac002.replace("$ORI_INSTR_ID", ORI_INSTR_ID);
    newpac002 = newpac002.replace("$ORI_END_TO_END_ID", ORI_END_TO_END_ID);
    newpac002 = newpac002.replace("$ORI_TX_ID", ORI_TX_ID);
    newpac002 = newpac002.replace("$CHARGS_ASSET_CODE", CHARGS_ASSET_CODE);
    newpac002 = newpac002.replace("$CHARGS_AMOUNT", CHARGS_AMOUNT);
    newpac002 = newpac002.replace("$ISSUER_BIC", RFI_BIC);
    newpac002 = newpac002.replace("$ISSUER_ID", RFI_ID);
    newpac002 = newpac002.replace("$ORI_SETTLE_ASSET_CODE", ORI_SETTLE_ASSET_CODE);
    newpac002 = newpac002.replace("$ORI_SETTLE_AMOUNT", ORI_SETTLE_AMOUNT);
    newpac002 = newpac002.replace("$RECEIVER_ID", RFI_ID);
    newpac002 = newpac002.replace("$RECEIVE_ACCOUNT_ADDRESS", RECEIVE_ACCOUNT_ADDRESS);
    console.log(newpac002);
    let message = encoder.encode(newpac002, 'base64')
        // let message = ''

    return message

}