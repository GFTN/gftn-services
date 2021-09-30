// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import * as crypto from 'crypto'
import * as _ from 'lodash'
import { Buffer } from 'buffer';

export interface IDecrypt {
    iv: string;
    salt: string;
    pass: string;
    iterate: string;
}

export class Encrypt {

    async encrypt(plainText: string): Promise<{
        enc: string,
        dec: IDecrypt
    }> {

        return new Promise(async (resolve, reject) => {

            try {

                const iv = await this.genIV()
                const salt = await this.genSalt()
                const iterate = await this.genIterate()
                const pass = await this.genPass()

                const pbkdf2Key: Buffer = await crypto.pbkdf2Sync(pass, salt, Number(iterate), 32, 'sha512');

                // create cipher
                const cipher = crypto.createCipheriv('aes-256-cbc', pbkdf2Key, Buffer.alloc(16, iv));
                const encodedText = Buffer.from(plainText).toString('base64');

                // encrypted output
                let encrypted: string = cipher.update(encodedText, 'utf8', 'hex');
                encrypted += cipher.final('hex');

                resolve({
                    enc: encrypted,
                    dec: {
                        iv: iv,
                        salt: salt,
                        pass: pass,
                        iterate: iterate
                    }
                })

            } catch (error) {

                console.error('encryption failed: ', error)
                reject(error);

            }

        });

    }

    async decrypt(pass: string, salt: string, iterate: number, iv: string, encText: string) {

        try {

            const pbkdf2Key: Buffer = await crypto.pbkdf2Sync(pass, salt, iterate, 32, 'sha512');

            const decipher = crypto.createDecipheriv('aes-256-cbc', pbkdf2Key, Buffer.alloc(16, iv));

            // decipher text
            let dec = decipher.update(encText, 'hex', 'utf8');
            dec += decipher.final('utf8');

            const decodedText = Buffer.from(dec, 'base64').toString();

            return decodedText;

        } catch (error) {
            console.error('decryption failed: ', error)
            return error
        }

    }

    async genPass() {
        const randPass = this.genRandomStr(32);
        return randPass;
    }

    async genSalt() {
        const randSalt = await this.genPass()
        return randSalt;
    }

    async genIterate() {
        const randIterate = Math.floor(Math.random() * 100);
        // to string because iterate will go into json to be consumed as system env (which requires type string)
        return String(randIterate);
    }

    async genIV() {
        const rand = this.genRandomStr(96);
        const iv = crypto.createHash('md5').update(rand).digest('hex');
        return iv;
    }

    public genRandomStr(length: number) {
        let result = '';
        const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
        const charactersLength = characters.length;
        for (let i = 0; i < length; i++) {
            result += characters.charAt(Math.floor(Math.random() * charactersLength));
        }
        return result;
    }

}        
        
