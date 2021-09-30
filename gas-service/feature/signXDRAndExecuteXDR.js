// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const signXDR = require("./signXDR.js")
const executeXDR = require("./executeXDR.js")

/**
 * 1. signedXDR
 * 2. execute signedXDR
 */
module.exports =
    /**
     * 
     * @param {String} signedXDRin 
     * @param {Array} lockAccounts 
     * @param {Array} unlockAccounts 
     * @param {String} accountstablename 
     */
    async function(signedXDRin, lockAccounts, unlockAccounts, accountstablename) {

        let signedXDR = await signXDR(signedXDRin, accountstablename)
        let result = await executeXDR(signedXDR, lockAccounts, unlockAccounts, accountstablename)
        return result

    }