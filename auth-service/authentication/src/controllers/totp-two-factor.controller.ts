// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import * as express from 'express';
import * as admin from 'firebase-admin';
import * as crypto from 'crypto';
import { get, isEmpty } from 'lodash';
import { Request, Body, Controller, Get, Post, Route } from 'tsoa';
import * as notp from 'notp';
import * as base32 from 'thirty-two';
import { TokenBody, TOTPResponse, TOTPProfile } from '../models/totp.interface';
import { IFirebaseUserRequest } from '../models/auth.model';
import { IGlobalEnvs } from '../environment';

export class Database {
    async createTOTP(accountName: string, key: string): Promise<boolean> {
        // console.debug('database.createTOTP called');
        const firebaseUid = await this.getUser(accountName);
        const profile: TOTPProfile = {
            key: key,
            registered: false
        };
        await admin.database().ref('totp/' + firebaseUid).update(profile);
        return true;
    }
    async getUser(email: string): Promise<string> {
        // console.debug('database.getUser called');
        // console.debug('user email: ' + email);
        // get firebase user id by email
        const firebaseUser = await admin.auth().getUserByEmail(email);
        // console.debug('firebase user:' + firebaseUser);
        return firebaseUser.uid;
    }
    async getTOTPKey(accountName: string): Promise<string> {
        // console.debug('database.getTOTPKey called');
        const firebaseUid = await this.getUser(accountName);
        const profile: admin.database.DataSnapshot = await admin.database().ref('totp/' + firebaseUid).once('value');
        const res: TOTPProfile = profile.val();
        return res.key;
    }
    async getTOTPStatus(accountName: string): Promise<boolean> {
        // console.debug('database.getTOTPStatus called');
        if (!accountName) {
            return false;
        }
        const firebaseUid = await this.getUser(accountName);
        // console.debug(firebaseUid);
        const profile: admin.database.DataSnapshot = await admin.database().ref('totp/' + firebaseUid).once('value');
        const res: TOTPProfile = profile.val();
        if (!res) {
            return false;
        }
        return res.registered;
    }
    async confirmTOTP(accountName: string): Promise<boolean> {
        // console.debug('database.confirmTOTP called');
        const firebaseUid = await this.getUser(accountName);
        await admin.database().ref('totp/' + firebaseUid).update({ 'registered': true });
        return true;
    }


}

/**
 * private methods for TOTPTwoFactor
 *
 * @class TOTPTwoFactorPrivate
 */

export class TOTPTwoFactorPrivate {
    /**
     * get user email from request passed through IBMid middleware
     *
     * @param {express.Request} req request object
     * @returns {string} email address
     * @memberof TOTPTwoFactorPrivate
     */
    getUserEmail(req: express.Request): string {
        if (isEmpty(get(req, 'user'))) {
            throw new Error('User must be logged in with IBMId first');
        }
        return req.user._json.email;
    }
    /**
     * random string generator for TOTP seed
     *
     * @param {number} len length of random string
     * @param {boolean} includeSpecialCharacters include special char or not
     * @returns {boolean}
     * @memberof TOTPTwoFactorPrivate
     */
    randStr(len: number, includeSpecialCharacters?: boolean) {
        let charSet = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXTZabcdefghiklmnopqrstuvwxyz';
        if (includeSpecialCharacters) {
            charSet += '!@#$%^&*()<>?/[]{},.:;';
        }
        const bytes = crypto.randomBytes(len || 32);
        const str = bytes.reduce((result, byte) => {
            return result + charSet[Math.floor(byte / 255.0 * (charSet.length - 1))];
          }, '');
        return str;
    }
}

/**
 * Endpoints for TOTP
 *
 * @export
 * @class TOTPTwoFactor
 */
export class TOTPTwoFactor {

    database: Database;

    private env: IGlobalEnvs = global['envs'];
    private privateMethod : TOTPTwoFactorPrivate;

    constructor() {
        this.database = new Database();
        this.privateMethod = new TOTPTwoFactorPrivate();
    }

    /**
     * Init TOTP 2FA registration
     *
     * @param {string} accountName  user email
     * @returns {Promise<TOTPResponse>}
     * @memberof TOTPTwoFactor
     */
    async createTOTP(accountName: string): Promise<TOTPResponse> {
        // console.debug('TOTPTwoFactor.createTOTP called', accountName);
        if (await this.checkRegistered(accountName)) {
            console.info(accountName + ' registered TOTP already!');
            return {
                success: true,
                registered: true,
                msg: accountName + ' registered TOTP!',
                data: {
                    qrcodeURI: null,
                    accountName: accountName,
                }
            };
        }
        const key = this.randStr(32, false);
        await this.database.createTOTP(accountName, key);
        // encoded will be the secret key, base32 encoded
        const encoded = base32.encode(key);

        // Google authenticator doesn't like equal signs
        const encodedForGoogle = encoded.toString().replace(/=/g, '');

        const label = encodeURIComponent(this.env.totp_label);
        // to create a URI for a qr code (change totp to hotp if using hotp)
        const qrcodeURI = 'otpauth://totp/' + label + ':' + accountName + '?secret=' + encodedForGoogle + '&issuer=' + label;
        const res = {
            success: true,
            registered: false,
            msg: 'creating TOTP, please confirm',
            data: {
                qrcodeURI: qrcodeURI,
                accountName: accountName,
            }
        };
        // console.debug("TOTP registration: ", res.data.accountName, res.success);
        return res;
    }

    /**
     * confirm TOTP to finish 2FA registration
     *
     * @param {string} accountName  user email
     * @param {string} token  TOTP
     * @returns {Promise<TOTPResponse>}
     * @memberof TOTPTwoFactor
     */
    async confirmTOTP(accountName: string, token: string): Promise<TOTPResponse> {
        // console.debug('TOTPTwoFactor.confirmTOTP called');

        if (await this.checkTOTP(accountName, token)) {
            const success = await this.database.confirmTOTP(accountName);
            const res = {
                success: success,
                msg: 'totp confirmed',
            };
            // console.debug(res);
            return res;
        }
        else {
            const res = {
                success: false,
                msg: 'totp failed to pass'
            };
            // console.debug(res);
            return res;
        }
    }
    /**
     * check if the supplied TOTP passes the verification
     *
     * @param {string} accountName  user email
     * @param {string} token  TOTP
     * @returns {Promise<boolean>}
     * @memberof TOTPTwoFactor
     */
    async checkTOTP(accountName: string, token: string): Promise<boolean> {
        // console.debug('TOTPTwoFactor.checkTOTP called');

        const key = await this.database.getTOTPKey(accountName);
        // console.debug(token);
        // console.debug(key);
        // options for totp verification
        const opt = {
            window: 1,
        };
        // Check TOTP is correct (HOTP if hotp pass type)
        const login = notp.totp.verify(token, key, opt);

        // invalid token if login is null
        if (!login) {
            console.info('checkTOTP: Token invalid, account: %s', accountName);
            return false;
        }
        // valid token
        console.info('checkTOTP: Token valid, sync value is %s, account: %s', login.delta, accountName);
        return true;

    }

    /**
     * middleware for checking if the supplied TOTP passes the verification using IBMId user
     *
     * @param {express.Request} req  request object
     * @param {express.Response} req  response object
     * @param {express.NextFunction} next  next function
     * @memberof TOTPTwoFactor
     */
    checkTOTPMiddleWareIBMIdUser = async (req: express.Request, res: express.Response, next: express.NextFunction) => {

        try {
            if (this.env.enable_2fa === 'false') {
                next();
                return;
            } else {

                // check if request contains user information from IBMId via passport.js
                const accountName = this.getUserEmail(req);

                // const token = req.body.token;
                if (isEmpty(req.headers['x-verify-code'])) {
                    // missing header
                    res.status(401);
                    res.send('unauthorized');
                }

                const token = req.headers['x-verify-code'] as string;
                if (await this.checkTOTP(accountName, token)) {
                    // console.debug('checkTOTP passed');
                    next();
                } else {
                    res.status(401);
                    res.send('totp two-factor authentication failed - IBMId');
                    return;
                }
            }

        } catch (e) {
            console.error(e);
            res.status(500);
            res.send('totp two-factor authentication failed - IBMId');
            return;
        }
    }


    /**
     * middleware for checking if the supplied TOTP passes the verification using firebase user
     *
     * @param {express.Request} req  request object
     * @param {express.Response} req  response object
     * @param {express.NextFunction} next  next function
     * @memberof TOTPTwoFactor
     */
    checkTOTPMiddleWareFirebaseUser = async (req: express.Request, res: express.Response, next: express.NextFunction) => {

        try {
            if (this.env.enable_2fa === 'false') {
                next();
                return;
            } else {

                // check if request contains user information from firebase user session authentication
                const accountName = req['email'];

                // const token = req.body.token;
                if (isEmpty(req.headers['x-verify-code'])) {
                    // missing header
                    res.status(401);
                    res.send('unauthorized');
                }

                const token = req.headers['x-verify-code'] as string;
                if (await this.checkTOTP(accountName, token)) {
                    // console.debug('checkTOTP passed');
                    next();
                } else {
                    res.status(401);
                    res.send('totp two-factor authentication failed - fb');
                    return;
                }
            }

        } catch (e) {
            console.error(e);
            res.status(500);
            res.send('totp two-factor authentication failed - fb');
            return;
        }
    }

    /**
     * middleware for checking if the account name (user email) in the path is the same as email address set by IBMid middleware
     *
     * @param {express.Request} req  request object
     * @param {express.Response} req  response object
     * @param {express.NextFunction} next  next function
     * @memberof TOTPTwoFactor
     */
    checkAccountNameMiddleWare = async (req: express.Request, res: express.Response, next: express.NextFunction) => {
        try {
            const accountName = req.params.accountName;
            const userEmail = this.getUserEmail(req);
            if (accountName === userEmail) {
                // console.debug('Middleware: accountNameCheck Passed');
                next();
            } else {
                res.status(401);
                res.send('totp two-factor authentication failed');
                return;
            }
        } catch (e) {
            console.error(e);
            res.status(500);
            res.send('totp two-factor authentication failed');
            return;
        }
    }
    /**
     * check if the user has registered TOTP or not
     *
     * @param {string} accountName  user email
     * @returns {Promise<boolean>}
     * @memberof TOTPTwoFactor
     */
    async checkRegistered(accountName: string): Promise<boolean> {
        // console.debug('TOTPTwoFactor.checkRegistered called');

        return await this.database.getTOTPStatus(accountName);
    }
    /**
     * get user email from request passed through IBMid middleware
     *
     * @param {express.Request} req request object
     * @returns {string} email address
     * @memberof TOTPTwoFactor
     */

    private getUserEmail(req: express.Request): string {
        return this.privateMethod.getUserEmail(req);
    }
    /**
     * random string generator for TOTP seed
     *
     * @param {number} len length of random string
     * @param {boolean} includeSpecialCharacters include special char or not
     * @returns {boolean}
     * @memberof TOTPTwoFactor
     */
    private randStr(len: number, includeSpecialCharacters?: boolean) {
        return this.privateMethod.randStr(len, includeSpecialCharacters);
    }
}

/**
 * Endpoints for TOTP
 * @export
 * @class TOTPTwoFactorAuthController
 * @extends {Controller}
 */
@Route('totp')
export class TOTPTwoFactorAuthController extends Controller {

    handlers = new TOTPTwoFactor();

    /**
     * Get registration QR Code
     * Returns TOTPResponse which contain a uri string that can be converted to a QR Code
     *
     * @param {string} accountName user email
     * @returns {Promise<TOTPResponse>}
     * @memberof TOTPTwoFactorAuthController
     */
    @Get('{accountName}')
    public async createTOTP(accountName: string): Promise<TOTPResponse> {

        // console.debug(accountName);

        if (accountName) {
            return await this.handlers.createTOTP(accountName).then((res: TOTPResponse) => {
                if (res.success && !res.registered) {
                    this.setStatus(200);
                    return res;
                } else if (res.success && res.registered) {
                    this.setStatus(201);
                    return res;
                } else {
                    this.setStatus(403);
                    return res;
                }
            },
                (e) => {
                    console.error(e);
                    this.setStatus(500);
                    return null;
                });
        } else {
            return null;
        }

    }

    /**
     * complete TOTP registration
     * Returns TOTPResponse which contain the msg
     *
     * @param {string} accountName user email
     * @param {TokenBody} body TOTP token
     * @returns {Promise<TOTPResponse>}
     * @memberof TOTPTwoFactorAuthController
     */
    @Post('{accountName}/confirm')
    public async confirmTOTP(accountName: string, @Body() body: TokenBody): Promise<TOTPResponse> {
        // console.debug(accountName);
        // console.debug(body);
        return await this.handlers.confirmTOTP(accountName, body.token).then((res: TOTPResponse) => {
            if (res.success) {
                this.setStatus(200);
            } else {
                this.setStatus(401);
            }
            return res;
        },
            (e) => {
                console.error(e);
                this.setStatus(500);
                return null;
            });
    }

    /**
     * check if TOTP is correct this endpoint is for other services like deployment service
     *
     * @param {IFirebaseUserRequest} req request which contain user email
     * @param {TokenBody} body TOTP token
     * @returns {Promise<boolean>}
     * @memberof TOTPTwoFactorAuthController
     */
    @Post('check')
    public async checkTOTP(@Request() req: IFirebaseUserRequest, @Body() body: TokenBody): Promise<boolean> {
        const accountName: string = req.email;
        return await this.handlers.checkTOTP(accountName, body.token).then((res: boolean) => {
            if (res) {
                this.setStatus(200);
            } else {
                this.setStatus(401);
            }
            return res;
        },
            (e) => {
                console.error(e);
                this.setStatus(500);
                return null;
            });
    }
}
