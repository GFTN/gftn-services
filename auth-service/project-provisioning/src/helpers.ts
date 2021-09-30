// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { cloneDeep, forEach, set } from 'lodash';

export interface IEncodeFirebaseCred {
    type: string;
    project_id: string;
    private_key_id: string;
    private_key: string;
    client_email: string;
    client_id: string;
    auth_uri: string;
    token_uri: string;
    auth_provider_x509_cert_url: string;
    client_x509_cert_url: string;
};

export interface IRandomPepperObj {
    // the old prefix key value for signing tokens
    o: number;
    // c = the current prefix key to use for signing tokens
    // the current prefix that should be used to generated new pepper values
    // format - prefix should be a single character a-z; convention = "{1-1:XXXXXRandomStringHereXXXXX}"
    c: number;
    // v = array of random values must be less than {4096 characters for aws secrets manager to store serialized data}
    // old values should be provided in body and appended to new values array for lookup of new and old values
    v: { [prefix: string]: string };
}

export class Helpers {

    /**
    * Generates a random string for a specified length.
    * Can be used to generated random passwords
    *
    * @param {number} len the length of the output string
    * @param {boolean} [includeSpecialCharacters] if the random string should include special characters
    * @returns
    * @memberof TwoFactorAPI
    */
    randStr(len: number, includeSpecialCharacters?: boolean) {
        let text = '';
        let possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

        // include special characters in random string (random password)
        if (includeSpecialCharacters) {
            possible = possible + "!@#$%^&*()_+<>?,./:[]{};'";
        }

        for (let i = 0; i < len; i++)
            text += possible.charAt(Math.floor(Math.random() * possible.length));

        return text;
    }

    public async encodeFbCred(
        body: IEncodeFirebaseCred
    ) {

        // set to escaped json string
        const sanitizedText = JSON.stringify({
            type: body.type,
            project_id: body.project_id,
            private_key_id: body.private_key_id,
            private_key: body.private_key,
            client_email: body.client_email,
            client_id: body.client_id,
            auth_uri: body.auth_uri,
            token_uri: body.token_uri,
            auth_provider_x509_cert_url: body.auth_provider_x509_cert_url,
            client_x509_cert_url: body.client_x509_cert_url,
        });

        // encode to base 64
        const encodedStr = Buffer.from(sanitizedText).toString('base64');
        return encodedStr;
    }

    getPepperObj() {

        // init object can be used to initialize the first
        // IRandomPepperObj if none exists
        const init: IRandomPepperObj = {
            o: 0,
            c: 0,
            v: {}
        };

        // length of secret value
        // best length would be 64 char string but since
        // this will be combined with another 64 char string (ie: the val stored in db)
        // and additional 32 chars (or more) seems sufficiently hard to crack
        // for total of 96 char len decryption key
        const len = 32;

        // qty of random pepper values in array
        const qty = 45;

        // get the env ww_jwt_pepper_obj
        const pepperObj = init;

        // delete out old prefixed values
        const newValues = {};
        forEach(pepperObj.v, (val, key) => {
            // set only values that don't equal the old values
            if (key.split('-')[0] !== pepperObj.o.toString()) {
                set(newValues, key, val);
            }
        });

        // update old prefix
        pepperObj.o = cloneDeep(pepperObj.c);

        // increment current
        pepperObj.c = pepperObj.c + 1;

        // start count over at 10 (preventing long prefixes that
        // would increase json length unnecessarily)
        if (pepperObj.c >= 10) {
            pepperObj.c = 0;
        }

        // add in new prefixed values
        for (let i = 0; i < qty; i++) {
            // set new prefix
            const prefix = pepperObj.c + '-' + i;
            // set new value
            const val = this.randStr(len, false);

            // push to existing array
            // pepperObj.v.push({ [prefix]: val });

            // set on newValue object
            set(newValues, prefix, val);
        }

        // set new values
        pepperObj.v = newValues;

        // serialize new object to json string
        const newEnvVar = JSON.stringify(pepperObj);

        // too manually copy into env var print to log
        // by uncommenting below line
        // console.log(newEnvVar);

        // check if the length is less than 4096
        // equation => old values + new values = 2 x n
        if (newEnvVar.length <= 4096) {

            return pepperObj;
        } else {
            console.error('failed to create a env value less than 4096 characters');
            return false;
        }

    }

}