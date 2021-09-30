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
    ASSET_ISSUER,
    INSTG_BIC,
    INSTG_ID,
    RECEIVER_BIC,
    RECEIVER_ID,
    SETTLE_ASSET_CODE,
    SETTLE_AMOUNT,
    EXCHANGE_RATE,
    FEE_ASSET_CODE,
    FEE_AMOUNT,
    CHARGER_BIC,
    SENDER_BANK_NAME,
    SENDER_STREET_NAME,
    SENDER_BUILDING_NUMBER,
    SENDER_POST_CODE,
    SENDER_TOWN_NAME,
    SENDER_COUNTRY,
    SENDER_BIC,
    RECEIVE_BIC,
    RECEIVER_BANK_NAME,
    RECEIVER_STREET_NAME,
    RECEIVER_BUILDING_NUMBER,
    RECEIVER_POST_CODE,
    RECEIVER_TOWN_NAME,
    RECEIVER_COUNTRY,
    RECEIVER_ADDRESS_LINE,
    // been add since v2.9.3.12_RC  ---------------------------------------
    FEE_ID,
    PAYOUT_ID,
    PAYOUT_ASSET_CODE,
    PAYOUT_AMOUNT
) => {
    // let pacs008 = await readFile('./file/send_sample.xml', 'utf8')
    let pacs008 = await readFile('./file/pacs008_template.xml', 'utf8')
        // ---------------------------------------

    let today = new Date();
    let DD = ('0' + today.getDate()).slice(-2);
    let MM = ('0' + (today.getMonth() + 1)).slice(-2);
    let YYYY = today.getFullYear();

    let OFI_MSG_ID_STR = SETTLE_ASSET_CODE + DD + MM + YYYY + SENDER_BIC

    let randomNumLen = 35 - OFI_MSG_ID_STR.length
    let pow = Math.pow(10, (randomNumLen - 1))
    let randomNum = Math.floor(Math.random() * ((9 * pow) - 1) + pow)
    process.env['OFI_MSG_ID_' + OFI_ID] = OFI_MSG_ID_STR + randomNum.toString()

    let TX_CREATE_DATE = YYYY + '-' + MM + '-' + DD
    process.env['OFI_' + OFI_ID + '_CREATE_DATE'] = TX_CREATE_DATE;


    randomNumLen = 10
    pow = Math.pow(10, (randomNumLen - 1))
    randomNum = Math.floor(Math.random() * ((9 * pow) - 1) + pow)
    let AAAAA
    if (SETTLE_ASSET_CODE.length < 5) {
        AAAAA = SETTLE_ASSET_CODE + 'XX'
    } else {
        AAAAA = SETTLE_ASSET_CODE
    }
    let INSTR_ID = AAAAA + YYYY + MM + DD + INSTG_BIC + 'B' + randomNum
    process.env['OFI_' + OFI_ID + '_ORI_INSTR_ID'] = INSTR_ID;
    process.env['OFI_E2E_ID_' + OFI_ID] = process.env['OFI_MSG_ID_' + OFI_ID]
    process.env['OFI_TX_ID_' + OFI_ID] = process.env['OFI_MSG_ID_' + OFI_ID]
        // process.env['OFI_E2E_ID_' + OFI_ID] = 'E2EID' + YYYY + MM + DD + INSTG_BIC + 'O' + randomNum
        // process.env['OFI_TX_ID_' + OFI_ID] = 'TXID' + YYYY + MM + DD + INSTG_BIC + 'OO' + randomNum

    let OFI_MSG_ID = process.env['OFI_MSG_ID_' + OFI_ID]
    let E2EID = process.env['OFI_E2E_ID_' + OFI_ID]
    let TX_ID = process.env['OFI_TX_ID_' + OFI_ID]
        // let TX_ID = INSTR_ID
    let TX_CREATE_TIME = new Date().toISOString().replace(/\..+/, '')
    process.env['OFI_' + OFI_ID + '_ORI_TX_CREATE_DATE_TIME'] = TX_CREATE_TIME;

    let newpacs008 = pacs008.replace("$MSG_ID", OFI_MSG_ID);
    newpacs008 = newpacs008.replace("$TX_CREATE_TIME", TX_CREATE_TIME);
    newpacs008 = newpacs008.replace("$SETTLE_METHOD", SETTLE_METHOD);
    newpacs008 = newpacs008.replace("$OFI_ID", OFI_ID);
    newpacs008 = newpacs008.replace("$OFI_ACCOUNT_NAME", OFI_ACCOUNT_NAME);
    newpacs008 = newpacs008.replace("$ASSET_ISSUER", ASSET_ISSUER);
    newpacs008 = newpacs008.replace("$INSTG_BIC", INSTG_BIC);
    newpacs008 = newpacs008.replace("$INSTG_ID", INSTG_ID);
    newpacs008 = newpacs008.replace("$INSTD_BIC", RECEIVER_BIC);
    newpacs008 = newpacs008.replace("$INSTD_ID", RECEIVER_ID);
    newpacs008 = newpacs008.replace("$RECEIVER_BIC", RECEIVER_BIC);
    newpacs008 = newpacs008.replace("$RECEIVER_ID", RECEIVER_ID);
    newpacs008 = newpacs008.replace("$INSTR_ID", INSTR_ID);
    newpacs008 = newpacs008.replace("$END_TO_END_ID", E2EID);
    newpacs008 = newpacs008.replace("$TX_ID", TX_ID);
    newpacs008 = newpacs008.replace("$SETTLE_ASSET_CODE", SETTLE_ASSET_CODE);
    newpacs008 = newpacs008.replace("$SETTLE_AMOUNT", SETTLE_AMOUNT);
    newpacs008 = newpacs008.replace("$BANK_SETTLEMENT_DATE", TX_CREATE_DATE);
    newpacs008 = newpacs008.replace("$EXCHANGE_RATE", EXCHANGE_RATE);
    newpacs008 = newpacs008.replace("$FEE_ASSET_CODE", FEE_ASSET_CODE);
    newpacs008 = newpacs008.replace("$FEE_AMOUNT", FEE_AMOUNT);
    newpacs008 = newpacs008.replace("$CHARGER_BIC", CHARGER_BIC);
    newpacs008 = newpacs008.replace("$CHARGER_ID", OFI_ID);
    newpacs008 = newpacs008.replace("$SENDER_BANK_NAME", SENDER_BANK_NAME);
    newpacs008 = newpacs008.replace("$SENDER_STREET_NAME", SENDER_STREET_NAME);
    newpacs008 = newpacs008.replace("$SENDER_BUILDING_NUMBER", SENDER_BUILDING_NUMBER);
    newpacs008 = newpacs008.replace("$SENDER_POST_CODE", SENDER_POST_CODE);
    newpacs008 = newpacs008.replace("$SENDER_TOWN_NAME", SENDER_TOWN_NAME);
    newpacs008 = newpacs008.replace("$SENDER_COUNTRY", SENDER_COUNTRY);
    newpacs008 = newpacs008.replace("$SENDER_BIC", SENDER_BIC);
    newpacs008 = newpacs008.replace("$RECEIVE_BIC", RECEIVE_BIC);
    newpacs008 = newpacs008.replace("$SENDER_ID", OFI_ID);
    newpacs008 = newpacs008.replace("$RECEIVER_BANK_NAME", RECEIVER_BANK_NAME);
    newpacs008 = newpacs008.replace("$RECEIVER_STREET_NAME", RECEIVER_STREET_NAME);
    newpacs008 = newpacs008.replace("$RECEIVER_BUILDING_NUMBER", RECEIVER_BUILDING_NUMBER);
    newpacs008 = newpacs008.replace("$RECEIVER_POST_CODE", RECEIVER_POST_CODE);
    newpacs008 = newpacs008.replace("$RECEIVER_TOWN_NAME", RECEIVER_TOWN_NAME);
    newpacs008 = newpacs008.replace("$RECEIVER_COUNTRY", RECEIVER_COUNTRY);
    newpacs008 = newpacs008.replace("$RECEIVER_ADDRESS_LINE", RECEIVER_ADDRESS_LINE);

    randomNumLen = 7
    pow = Math.pow(10, (randomNumLen - 1))
    randomNum = Math.floor(Math.random() * ((9 * pow) - 1) + pow)

    let BUSINESS_MSG_ID = 'B' + YYYY + MM + DD + SENDER_BIC + 'BAA' + randomNum.toString()
        // been add since v2.9.3.12_RC
    newpacs008 = newpacs008.replace("$HEADER_BIC", INSTG_BIC);
    newpacs008 = newpacs008.replace("$HEADER_SENDER_ID", OFI_ID);
    newpacs008 = newpacs008.replace("$BUSINESS_MSG_ID", BUSINESS_MSG_ID);
    newpacs008 = newpacs008.replace("$MSG_DEF_ID", OFI_MSG_ID);
    newpacs008 = newpacs008.replace("$HEADER_TX_CREATE_TIME", TX_CREATE_TIME);
    newpacs008 = newpacs008.replace("$FEE_ID", FEE_ID);
    newpacs008 = newpacs008.replace("$PAYOUT_ID", PAYOUT_ID);
    newpacs008 = newpacs008.replace("$PAYOUT_ASSET_CODE", PAYOUT_ASSET_CODE);
    newpacs008 = newpacs008.replace("$PAYOUT_AMOUNT", PAYOUT_AMOUNT);
    // ---------------------------------------

    console.log(newpacs008);

    let message = encoder.encode(newpacs008, 'base64')
    fs.appendFile('./file/send_endtoendID.txt', today + ' -  ' + OFI_MSG_ID + '\n', function(err) {
        if (err) throw err;
        // console.log('Saved!');
    });
    return message


}