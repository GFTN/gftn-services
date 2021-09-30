// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 

import * as crypto from 'crypto';

class _Encrypt {

    private alg: string;
    private secret: string;
    private iv: string;

    constructor() {
        // set the secret key and initialization vector
        this.genSecretIv('aes-256-cbc');
    }

    /**
     * generate the secret and initialization vector
     *
     * @private
     * @param {'aes-256-cbc'} alg
     * @memberof Encryption
     */
    private genSecretIv(alg: 'aes-256-cbc') {

        if (alg === 'aes-256-cbc') {
            this.alg = 'aes-256-cbc';
            // secret MUST be 32 characters using aes-256-cbc
            this.secret = new Buffer(crypto.randomBytes(16)).toString('hex');

            // iv MUST be 16 characters using using aes-256-cbc
            // do not use a global iv for production, 
            // generate a new one for each encryption
            this.iv = new Buffer(crypto.randomBytes(8)).toString('hex');
        }

        // if else {
        // // other algorithms here...
        // }

        else {
            throw new Error('Unrecognized algorithm');
        }

    }

    /**
     * encrypt text
     *
     * @param {string} text
     * @returns {string}
     * @memberof Encryption
     */
    encrypt(text: string): { encryptedText: string, secret: string, iv: string } {
        const cipher = crypto.createCipheriv(this.alg, this.secret, this.iv);
        let encrypted = cipher.update(text, 'utf8', 'hex');
        encrypted += cipher.final('hex');
        // console.log('encrypted ', encrypted);
        return {
            encryptedText: encrypted,
            secret: this.secret,
            iv: this.iv
        };
    }

}

export function encrypt(text: string): { encryptedText: string, secret: string, iv: string } {
    return new _Encrypt().encrypt(text);
}

/**
 * decrypt text
 *
 *
 * @export
 * @param {string} encryptedText
 * @param {string} secret Store separately from iv
 * @param {string} iv Store separately from secret 
 * @returns {string}
 */
export function decrypt(encryptedText: string, secret: string, iv: string, alg?: 'aes-256-cbc'): string {

    // set default algorithm
    const _alg = alg ? alg : 'aes-256-cbc';

    // decipher text
    const decipher = crypto.createDecipheriv(_alg, secret, iv);
    let dec = decipher.update(encryptedText, 'hex', 'utf8');
    dec += decipher.final('utf8');

    // console.log('decrypted ', dec);

    return dec;
}