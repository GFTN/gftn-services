// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
let StellarSdk = require('stellar-sdk');
let environment = require('../environment/env')
const LOGGER = require('../method/logger')
const log = new LOGGER('Stellar Call')
let server

StellarSdk.Network.use(new StellarSdk.Network(process.env[environment.ENV_KEY_STELLAR_NETWORK]));
server = new StellarSdk.Server(process.env[environment.ENV_KEY_HORIZON_CLIENT_URL], { allowHttp: true });

module.exports = {

    /* get account Info. from stellar */
    /**
     * 
     * @param {Object} pkey 
     */
    getBalance: function(account) {
        return new Promise(function(res, rej) {
            try {

                server.accounts()
                    .accountId(account.toString())
                    .call()
                    .then(function(accountResult) {
                        // accountResult.balances= accountResult.balances
                        res(accountResult.balances)
                        log.logger('stellar-API - server.accounts()', '')
                            // console.log(accountResult.balances);

                    })
                    .catch(function(err) {
                        log.error('stellar-API - server.accounts()', JSON.stringify(err))
                        rej(err);
                    })
            } catch (err) {
                log.error('stellar-API - server.accounts()', JSON.stringify(err))
                rej(err);

            }

        })
    },
    transactionBuilder:
    /**
     * 
     * @param {String} source 
     */
        function(source) {
        return new Promise(function(res, rej) {
            log.logger('SDK - TransactionBuilder(source)', 'Created')
            res(new StellarSdk.TransactionBuilder(source))
        }).catch(function(err) {

            log.logger('SDK - TransactionBuilder(source)', JSON.stringify(err))
            rej(err)
        })
    },
    getAccountSequenceNumber:
    /**
     * 
     * @param {String} account 
     */
        function(account) {
        return new Promise(function(res, rej) {
            server.loadAccount(account)
                .then(function(account) {
                    res(account.sequence)
                })

        })
    },

    getBuilderAccount: function(pkey, sequenceNumber) {
        return new Promise(function(res, rej) {
            try {
                res(new StellarSdk.Account(pkey, sequenceNumber.toString()))
            } catch (err) {
                log.error('SDK - getBuilderAccount', err)
                throw (err)
            }
        })
    },
    submitTransaction:
    /**
     * 
     * @param {Object} transaction 
     */
        function(transaction) {
        return new Promise(function(res, rej) {
            server.submitTransaction(transaction)
                .then(function(transactionResult) {
                    log.logger('stellar-API -  server.submitTransaction(transaction)', transactionResult)
                    res(JSON.stringify(transactionResult, null, 2))
                })
                .catch(function(err) {
                    log.error('stellar-API -  server.submitTransaction(transaction)', JSON.stringify(err))
                    rej(err)
                });
        })

    },
    getAsset:
    /**
     * 
     * @param {Object} asset 
     */
        function(asset) {
        return new Promise(function(res, rej) {
            if (asset.code == '') {
                log.logger('SDK -  StellarSdk.Asset.native()', StellarSdk.Asset.native())
                res(StellarSdk.Asset.native())
            } else {
                log.logger('SDK -  StellarSdk.Asset(asset.code, asset.issuer)', new StellarSdk.Asset(asset.code, asset.issuer))
                res(new StellarSdk.Asset(asset.code, asset.issuer))
            }
        })
    },
    addPaymentOperation:
    /**
     * 
     * @param {Object} transaction 
     * @param {String} source 
     * @param {String} destination 
     * @param {Object} asset 
     * @param Float*} amount 
     */
        function(transaction, source, destination, asset, amount) {
        return new Promise(function(res, rej) {
            try {
                log.logger('stellar-API - addOperation - StellarSdk.Operation.payment', 'Create')
                res(transaction
                    .addOperation(StellarSdk.Operation.payment({
                        source: source,
                        destination: destination,
                        // asset: StellarSdk.Asset.native(),
                        asset: asset,
                        amount: amount
                    }))
                    // .setTimeout(StellarSdk.TimeoutInfinite)
                )
            } catch (err) {
                log.error('SDK - addOperation', JSON.stringify(err))
                rej(err)
            }
        })
    },
    addSignerOperation:
    /**
     * 
     * @param {Object} transaction 
     * @param {String} destination 
     */
        function(transaction, destination) {
        return new Promise(function(res, rej) {
            try {
                log.logger('stellar-API - addOperation - StellarSdk.Operation.setOptions', 'Create')
                res(transaction
                    .addOperation(StellarSdk.Operation.setOptions({
                        source: destination,
                        signer: {
                            ed25519PublicKey: destination,
                            weight: 4
                        }
                    }))
                )
            } catch (err) {
                log.error('SDK - addOperation', JSON.stringify(err))
                rej(err)
            }

        })
    },
    addSetWeightOperation:
    /**
     * 
     * @param {Object} transaction 
     * @param {String} destination 
     */
        function(transaction, destination) {
        return new Promise(function(res, rej) {
            try {
                log.logger('SDK - addOperation - StellarSdk.Operation.setOptions', 'Create')
                res(transaction
                    .addOperation(StellarSdk.Operation.setOptions({
                        source: destination,
                        masterWeight: 1,
                        lowThreshold: 1,
                        medThreshold: 2,
                        highThreshold: 2
                    })))
            } catch (err) {
                log.error('SDK - addOperation', JSON.stringify(err))
                rej(err)
            }
        })
    },
    buildTransaction:
    /**
     * 
     * @param {Object} transaction 
     */
        function(transaction) {
        return new Promise(function(res, rej) {
            try {
                log.logger('Transaction Build', 'Promise Object')
                res(transaction.build())
            } catch (err) {
                log.error('Transaction Build', JSON.stringify(err))
                rej(err)
            }
        })
    },
    signTx:
    /**
     * 
     * @param {Object} transaction 
     * @param String*} secret 
     */
        function(transaction, secret) {
        return new Promise(function(res, rej) {
            try {
                transaction.sign(StellarSdk.Keypair.fromSecret(secret))
                log.logger('stellar-API - sign', JSON.stringify(transaction))
                res(transaction)
            } catch (err) {
                log.error('stellar-API - sign', JSON.stringify(err))
                rej(err)
            }
        })
    },
    decode:
    /**
     * 
     * @param {Object} transaction 
     */
        function(transaction) {
        return new Promise(function(res, rej) {
            try {
                res(StellarSdk.xdr.TransactionEnvelope.fromXDR(transaction, 'base64'))
            } catch (err) {
                log.error('SDK - .xdr.TransactionEnvelope.fromXDR', JSON.stringify(err))
                rej(err)
            }
        })
    },
    newTransaction:
    /**
     * 
     * @param {String} signedXDR 
     */
        function(signedXDR) {
        return new Promise(function(res, rej) {
            try {
                res(new StellarSdk.Transaction(signedXDR))
            } catch (err) {
                rej(err)
            }

        })
    },
    submitTransaction:
    /**
     * 
     * @param {Object} transaction 
     */
        function(transaction) {
        return new Promise(async function(res, rej) {
            try {
                let tx = await server.submitTransaction(transaction)
                let r = {
                    title: "Transaction successful",
                    hash: tx.hash,
                    ledger: tx.ledger
                }
                log.logger('stellar-API - submitTransaction', JSON.stringify(tx))
                res(r)
            } catch (err) {
                let errMsg = {
                    title: "Transaction Failed",
                    failure_reason: err.response.data.extras
                }
                log.error('stellar-API - submitTransaction Fail', err)
                rej(errMsg)
            }

        })
    }

}