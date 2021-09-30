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
    TX_CREATE_TIME,
    SETTLE_METHOD,
    RECEIVER_ID,
    RECEIVER_ACCOUNT_NAME,
    INSTG_BIC,
    INSTG_ID,
    PAYMENT_SENDER_BIC,
    PAYMENT_SENDER_ID,
    RECEIVER_ACCOUNT_ADDRESS,
    FEDERATIONS_STATUS,
    FEDERATIONS_END_TO_END_ID,
    FEDERATIONS_INSTR_ID,
    COMPLIANCE_STATUS_1,
    COMPLIANCE_STATUS_2,
    COMPLICANCE_END_TO_END_ID,
    COMPLICANCE_INSTR_ID

) => {
    let ibwf001 = await readFile('./file/ibwf001_template.xml', 'utf8')
    let today = new Date();
    let DD = ('0' + today.getDate()).slice(-2);
    let MM = ('0' + (today.getMonth() + 1)).slice(-2);
    let YYYY = today.getFullYear();
    let randomNumLen = 10
    let pow = Math.pow(10, (randomNumLen - 1))
    let randomNum = Math.floor(Math.random() * ((9 * pow) - 1) + pow)

    let AAAAA
    if (SEND_ASSET_CODE.length < 5) {
        AAAAA = SEND_ASSET_CODE + 'XX'
    } else {
        AAAAA = SEND_ASSET_CODE
    }
    // let RFI_MSG_ID = SEND_ASSET_CODE + DD + MM + YYYY + RFI_BIC
    let RFI_MSG_ID = AAAAA + YYYY + MM + DD + RFI_BIC + 'B' + randomNum
    let newibwf001 = ibwf001.replace("$MSG_ID", RFI_MSG_ID);
    newibwf001 = newibwf001.replace("$RES_ID", RFI_MSG_ID);
    newibwf001 = newibwf001.replace("$TX_CREATE_TIME", TX_CREATE_TIME);
    newibwf001 = newibwf001.replace("$SETTLE_METHOD", SETTLE_METHOD);
    newibwf001 = newibwf001.replace("$RECEIVER_ID", RECEIVER_ID);
    newibwf001 = newibwf001.replace("$RECEIVER_ACCOUNT_NAME", RECEIVER_ACCOUNT_NAME);
    newibwf001 = newibwf001.replace("$INSTG_BIC", INSTG_BIC);
    newibwf001 = newibwf001.replace("$INSTG_ID", INSTG_ID);
    newibwf001 = newibwf001.replace("$PAYMENT_SENDER_BIC", PAYMENT_SENDER_BIC);
    newibwf001 = newibwf001.replace("$PAYMENT_SENDER_ID", PAYMENT_SENDER_ID);
    newibwf001 = newibwf001.replace("$RECEIVER_ACCOUNT_ADDRESS", RECEIVER_ACCOUNT_ADDRESS);
    newibwf001 = newibwf001.replace("$FEDERATIONS_STATUS", FEDERATIONS_STATUS);
    newibwf001 = newibwf001.replace("$FEDERATIONS_END_TO_END_ID", FEDERATIONS_END_TO_END_ID);
    newibwf001 = newibwf001.replace("$FEDERATIONS_INSTR_ID", FEDERATIONS_INSTR_ID);
    newibwf001 = newibwf001.replace("$COMPLIANCE_STATUS_1", COMPLIANCE_STATUS_1);
    newibwf001 = newibwf001.replace("$COMPLIANCE_STATUS_2", COMPLIANCE_STATUS_2);
    newibwf001 = newibwf001.replace("$COMPLICANCE_END_TO_END_ID", COMPLICANCE_END_TO_END_ID);
    newibwf001 = newibwf001.replace("$COMPLICANCE_INSTR_ID", COMPLICANCE_INSTR_ID);
    // been add since v2.9.3.12_RC

    let randomNumLen3 = 7
    let pow3 = Math.pow(10, (randomNumLen3 - 1))
    let randomNum3 = Math.floor(Math.random() * ((9 * pow3) - 1) + pow3)

    let BUSINESS_MSG_ID = 'B' + YYYY + MM + DD + PAYMENT_SENDER_BIC + 'BAA' + randomNum3.toString()
    newibwf001 = newibwf001.replace("$HEADER_BIC", INSTG_BIC);
    newibwf001 = newibwf001.replace("$HEADER_SENDER_ID", RECEIVER_ID);
    newibwf001 = newibwf001.replace("$BUSINESS_MSG_ID", BUSINESS_MSG_ID);
    newibwf001 = newibwf001.replace("$MSG_DEF_ID", RFI_MSG_ID);
    newibwf001 = newibwf001.replace("$HEADER_TX_CREATE_TIME", TX_CREATE_TIME);
    // ---------------------------------------
    console.log(newibwf001);
    let message = encoder.encode(newibwf001, 'base64')
        // let message = ''

    return message

}