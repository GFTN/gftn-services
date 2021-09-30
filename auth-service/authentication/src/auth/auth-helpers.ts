// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 

// PURPOSE: these are helper functions that are utilized by IBMId and Firebase Authentication via ./src/controller/ibmid.controller.ts.

import * as admin from 'firebase-admin';
// import { Email } from '../email/email';
// import * as bcrypt from 'bcrypt';
import { includes } from 'lodash';
import * as moment from 'moment';
import { IJWTTokenClaimsAndPayloadSecure } from '../shared/models/token.interface';
import { IUserProfile } from '../shared/models/user.interface';

export class AuthHelpers {

    // private email: Email;
    
    constructor() {
        // this.email = new Email();
    }

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

    // /**
    //  * Generates a hash of plain text password (or secret)
    //  *
    //  * @param {string} text
    //  * @returns {Promise<{ hash: string, salt: string }>}
    //  * @memberof DeveloperController
    //  */
    // async genHash(text: string): Promise<{ hash: string, salt: string }> {

    //     // gen salt
    //     const salt = bcrypt.genSaltSync();

    //     // use bycrypt with salt and pepper
    //     const hash: string = await bcrypt.hash(text + this.env.dev_key_pepper, salt);
    //     // .then((hash) => {
    //     //     // Store hash in your password DB.
    //     // });

    //     return { hash: hash, salt: salt };
    // }

    // /**
    //  * compares hashes of plain text plain text password (or secret) to output hash
    //  *
    //  * @param {string} text the plain text password (or secret)
    //  * @param {string} hash the hash to match the text too
    //  * @returns {Promise<boolean>}
    //  * @memberof DeveloperController
    //  */
    // async compareHash(text: string, hash: string): Promise<boolean> {
    //     // https://stackoverflow.com/questions/41564229/why-isnt-salt-required-to-compare-whether-password-is-correct-in-bcrypt
    //     return await bcrypt.compare(text + this.env.dev_key_pepper, hash);
    // }

    /**
     * runs custom validation checks to determine if a jwt token is valid
     *
     * @param {IJWTTokenClaimsAndPayloadSecure} decodedToken
     * @param {string} jtiFromDb the jti stored at /jwt_secure/{kid}/i
     * @param {number} nFromDb the jti stored at /jwt_secure/{kid}/n
     * @param {string} compareIncomingIp the requests IP Address
     * @param {string} [compareEndpoint] optional parameter to check if token is valid against a specified endpoint
     * @param {string} [compareAccount]  optional parameter to check if token is valid against a specified account
     * @returns {boolean}
     * @memberof AuthHelpers
     */
    verifyWWTokenCustom(decodedToken: IJWTTokenClaimsAndPayloadSecure, jtiFromDb: string, nFromDb: number, compareIncomingIp: string, compareEndpoint?: string, compareAccount?: string): { pass: boolean, msg: string } {

        let msg = 'failed: ';

        // check expiration date
        let isNotExpired = false;
        if (decodedToken.exp >= Number(moment.utc().format('X'))) {
            isNotExpired = true;
        } else {
            msg = msg + 'isNotExpired:now=' + moment.utc().format('X') + ',exp=' + decodedToken.exp + '; ';
        }

        // not before date 
        let isNotBefore = false;
        if (decodedToken.nbf <= Number(moment.utc().format('X'))) {
            isNotBefore = true;
        } else {
            msg = msg + 'decodedToken:now=' + moment.utc().format('X') + ',nbf=' + decodedToken.nbf + '; ';
        }

        // check if token has account in array
        let hasAccount = false;
        // only check account if provided
        if (compareAccount) {
            if (includes(decodedToken.acc, compareAccount)) {
                hasAccount = true;
            } else {
                msg = msg + 'compareAccount:comparedTo=' + compareAccount + '; ';
            }
        } else {
            hasAccount = true;
        }

        // check if token has IP in array
        let hasIP = false;
        // only check ip if provided 
        if (compareIncomingIp) {
            if (includes(decodedToken.ips, compareIncomingIp)) {
                hasIP = true;
            } else {
                msg = msg + 'compareIncomingIp:incomming=' + compareIncomingIp + '; ';
            }
        } else {
            hasIP = true;
        }

        // check if token has Endpoint in array
        let hasEndpoint = false;
        // only check endpoint if provided 
        if (compareEndpoint) {
            if (includes(decodedToken.enp, compareEndpoint)) {
                hasEndpoint = true;
            } else {
                msg = msg + 'compareEndpoint=' + compareEndpoint + '; ';
            }
        } else {
            hasEndpoint = true;
        }

        // check if decoded jti matches the jti stored in the db
        let matchingJti = false;
        if (decodedToken.jti === jtiFromDb) {
            matchingJti = true;
        } else {
            msg = msg + 'matchingJti' + '; ';
        }

        // check if token is on count
        let isOnCount = false;
        // check if the count in db is same as the one in the db
        if (decodedToken.n === nFromDb) {
            isOnCount = true;
        } else {
            msg = msg + 'isOnCount' + '; ';
        }

        // send response
        if (isNotExpired &&
            isNotBefore &&
            hasAccount &&
            hasIP &&
            hasEndpoint &&
            isOnCount &&
            matchingJti
            // matchSubStr
        ) {
            // token valid
            return {pass: true, msg: 'clear'};
        } else {
            // token invalid
            return {pass: false, msg: msg};
        }

    }

    // /**
    //  * Check if IP Address is associated with user account
    //  *
    //  * @param {Request} _req
    //  * @memberof AuthHelpers
    //  */
    // async checkIp(_req: Request) {

    //     // /**
    //     //  * // TODO: Send across machine info in payload
    //     //  * // to be implemented in the browser client 
    //     //  * 
    //     //  * https://stackoverflow.com/questions/11219582/how-to-detect-my-browser-version-and-operating-system-using-javascript
    //     const OSName = "Unknown OS";
    //     try {
    //         if (navigator.appVersion.indexOf("Win") != -1) OSName = "Windows";
    //         if (navigator.appVersion.indexOf("Mac") != -1) OSName = "MacOS";
    //         if (navigator.appVersion.indexOf("X11") != -1) OSName = "UNIX";
    //         if (navigator.appVersion.indexOf("Linux") != -1) OSName = "Linux";

    //     } catch (error) {

    //     }
    //     interface IIpReq extends Request {
    //         connection: any;
    //     }

    //     const req = _req as IIpReq;

    //     // get request IP Address
    //     const ip = req.headers['x-forwarded-for'] ||
    //         req.connection.remoteAddress ||
    //         req.socket.remoteAddress ||
    //         (req.connection.socket ? req.connection.socket.remoteAddress : null);

    // }

    // /**
    //  * Welcome email inviting user notifying them 
    //  * that their account has been established
    //  *
    //  * @param {string} toEmailAddr
    //  * @memberof UserModel
    //  */
    // sendWelcomeEmail(first: string, email: string) {
    //     // send confirmation email to user and cc admin
    //     return this.email.sendEmail(
    //         {
    //             sender: {
    //                 name: 'IBM Blockchain Wold Wire',
    //                 email: 'noreply@ibm.com'
    //             },
    //             to: [{
    //                 name: first,
    //                 email: email
    //             }],
    //             // bcc so that admin is alerted
    //             bcc: [{
    //                 name: 'your.user',
    //                 email: 'your.user@your.domain'
    //             }],
    //             htmlContent: this.email.htmlTemplates.addUnsubscribeBtn(
    //                 this.email.htmlTemplates.welcome(first),
    //                 email,
    //                 ['na'],
    //                 true
    //             ),
    //             textContent: this.email.plainTemplates.welcome(first),
    //             subject: 'Welcome to IBM Blockchain World Wire',
    //         }
    //     );
    // }

    
    /**
     * Gets (or creates) firebase user for email provided
     *
     * @param {string} email
     * @returns {Promise<string>}
     * @memberof AuthHelpers
     */
    public setFirebaseUser(email: string): Promise<string> {

        // IMPORTANT: DO NOT PASS IN "PassportUser" as PARAMETERS
        // the user must be setup via email address first 
        // later the user profile can be updated on first login

        return new Promise((resolve, reject) => {

            admin.auth().getUserByEmail(email)
                .then((user: admin.auth.UserRecord) => {

                    // if (user.displayName) {

                    //     // if no displayName exists for user
                    //     // it is likely the first time they have logged in
                    //     // as a result the displayName from IBMId needs to
                    //     // added to the user's firebase auth profile
                    //     admin.auth().updateUser(user.uid,
                    //         { displayName: displayName });

                    //     // update firebase database profile
                    //     admin.database().ref('/users/' + user.uid).update({
                    //         displayName: displayName
                    //     });

                    // }

                    // user already has been created so send back firebase uid
                    resolve(user.uid);

                }, (err) => {

                    // no user with this email exists, so create a new user
                    admin.auth()
                        .createUser({
                            // display name will not be known when permissions are created
                            // displayName: displayName,
                            email: email
                        }).then((_user: admin.auth.UserRecord) => {

                            const data: IUserProfile = {
                                profile: {
                                    email: email
                                }
                            };

                            // create user profile in firebase 
                            admin.database().ref('/users/' + _user.uid).update(data);

                            // return new user id
                            resolve(_user.uid);
                        }, (err2: any) => {
                            console.error('failed to create user: ', err);
                            reject();
                        });

                }).catch((err: any) => {
                    console.error('failed to create user: ', err);
                    reject();
                });
        });

    }

    /**
     * Creates a firebase custom auth token to be used to sign-in at the portal
     *
     * @private
     * @param {string} uid
     * @param {Object} additionalClaims
     * @returns {Promise<string>}
     * @memberof PassportConfig
     */
    public async createCustomFirebaseToken(uid: string, additionalClaims: Object): Promise<string> {

        // const uid = 'some-uid';
        // const additionalClaims = {
        //     premiumAccount: true
        // };

        // create user with custom auth token
        return await admin.auth().createCustomToken(uid, additionalClaims)
            .then((customToken) => {

                // console.log(customToken);
                // returns the token for a new (or existing) firebase user
                return customToken;
            })
            .catch((err) => {
                // console.log(err);
                return Promise.reject(err);
            });

    }

}
