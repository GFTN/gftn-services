// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const stellar = require("../method/stellar.js")
const LOGGER = require('../method/logger')
const log = new LOGGER('Create Tx')
/**
 * 1. careate transaction : build a transaction , 
 * 2. use from_pair secrect to sign the transaction
 * 3. return a signedXDR
 */
module.exports = 
/**
 * 
 * @param {Integer} sequenceNumber 
 * @param {String} from_pair_pkey 
 * @param {String} from_pair_secrect 
 * @param {String} to_pair_pkey 
 * @param {Object} from_pay_asset 
 * @param {Object} to_pay_asset 
 * @param {Float} from_pay_amount 
 * @param {Float} to_pay_amount 
 */
async function (sequenceNumber, from_pair_pkey, from_pair_secrect, to_pair_pkey, from_pay_asset, to_pay_asset, from_pay_amount, to_pay_amount) {

    let account = await stellar.getBuilderAccount(to_pair_pkey, sequenceNumber)

    /**
     * build transaction
     * Need to check what if the source was something else?
     */

    let transaction = await stellar.transactionBuilder(account)
    
    from_pay_asset = await stellar.getAsset(from_pay_asset)
    to_pay_asset = await stellar.getAsset(to_pay_asset)


    transaction = await stellar.addPaymentOperation(transaction, from_pair_pkey, to_pair_pkey, from_pay_asset, from_pay_amount)
    transaction = await stellar.addPaymentOperation(transaction,to_pair_pkey, from_pair_pkey, to_pay_asset, to_pay_amount)    

    transaction = await stellar.buildTransaction(transaction)
    
    log.logger('Unsigned Tx',transaction.toEnvelope().toXDR('base64'))
    

    /**
     * sign transaction
     * this is where we'd change the code to sign from the third party, 
     * i.e. IBM gas source account?
     */
    

    log.logger('Using Secrect',from_pair_secrect)
    
    transaction = await stellar.signTx(transaction, from_pair_secrect)
    log.logger('Signed Tx',transaction.toEnvelope().toXDR('base64'))
    
    return transaction.toEnvelope().toXDR('base64')
}
