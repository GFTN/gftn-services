// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// PURPOSE: These endpoints are used to manage "internal" IBMId logins for ibm.com
// Unfortunately, the documentation for use of IBMId internally is lacking, see links
// used below. The IBMId implementation here makes use of passport.js (http://www.passportjs.org/),
// and implements the standard for OpenId and OAuth. This login flow requires an additional authentication
// provider, which is firebase authentication (https://firebase.google.com/docs/auth/) which secures front-end
// access to resources via the firebase database. Security with firebase is implemented in "BOLT"
// (https://github.com/FirebaseExtended/bolt) - security rules for bolt can be
// viewed here: https://github.com/GFTN/gftn-web/blob/development/database.rules.bolt.

// IMPORTANT: the environment variables for IBMId in the .vscode/launch.json are
// different for production vs. development. The production credentials utilize the real, live
// running instance of IBMId. Where as the development env var for IBMId run a sandbox environment.

// IBMId resources:
// http://ibm.biz/IBMid_main
// http://ibm.biz/IBMid_start
// https://w3.innovate.ibm.com/tools/sso/application/list.html
// https://w3.innovate.ibm.com/tools/sso
// https://www.ibm.com/support/knowledgecenter/SSPREK_9.0.0/com.ibm.isam.doc/config/concept/con_oidc_support.html

/*
// IBMId Redirect urls
================ CURRENT ===============
https://auth.worldwire-dev.io/sso/callback
https://auth.worldwire-qa.io/sso/callback
https://auth.worldwire-st.io/sso/callback
https://auth.worldwire.io/sso/callback
*/

import { Route, Controller, Get, Request, Post, OperationId
    // , Security
} from 'tsoa';
import { IPassportRequest } from '../models/auth.model';
import { Response } from 'express';
import { AuthHelpers } from '../auth/auth-helpers';
import * as admin from 'firebase-admin';
import { isEmpty, isUndefined, get } from 'lodash';
// import { TwoFactorOperations, I2faProfile } from './two-factor.controller';
import { TOTPProfile, TOTPRegistrationData } from '../models/totp.interface';
import { IGlobalEnvs } from '../environment';

@Route('sso')
export class IBMIdController extends Controller {

    private env: IGlobalEnvs = global['envs'];
    private authHelpers = new AuthHelpers();

    /**
     * Called to force login with IBMId
     * NOTE: This endpoint uses middleware via passport.js
     *
     * @returns
     * @memberof IBMIdController
     */
    @Get('login')
    @OperationId('authSSOLogin')
    // @Security('api_key')
    public async login() {
        // https://localhost:6001/sso/login
        return null;
    }

    // /**
    //  * Generates a new firebase users and returns a firebase userId
    //  *
    //  * @returns
    //  * @memberof IBMIdController
    //  */
    // @Post('new-user')
    // public async setFirebaseUser(@Body() body: { email: string }): Promise<string> {
    //     return await new AuthHelpers().setFirebaseUser(body.email);
    // }

    /**
     * For IBMId User this endpoint creates a users in
     * firebase (if they do not previously exist) and creates a
     * custom firebase auth Token to be consumed by the
     * client to log in a user to the portal. NOTE: This
     * endpoint uses middleware via passport.js
     *
     * @param {IPassportRequest} req
     * @returns {Promise<any>}
     * @memberof IBMIdController
     */
    @Get('token')
    @OperationId('authSSOToken')
    // @Security('api_key')
    public async token(@Request() req: IPassportRequest): Promise<any> {
        return new Promise((resolve, reject) => {

            const res = (<any>req).res as Response;

            // check if user exists on request via passport.js
            if (req.user) {

                // get firebase uid for user email
                this.authHelpers.setFirebaseUser(req.user._json.email)
                    .then((firebaseUid: string) => {

                        if (this.env.enable_2fa === 'false') {

                            // if 2fa not enabled, should redirect to login
                            resolve(this.redirect(
                                req,
                                '2fa/verify',
                                { disabled: true }
                            ));

                        } else {

                            // 1: check if user has registered 2FA username
                            admin.database().ref('totp/' + firebaseUid)
                                .once('value', (data: admin.database.DataSnapshot) => {

                                    // const twoFactor: I2faProfile = data.val();

                                    const twoFactor: TOTPProfile = data.val();

                                    if (get(twoFactor, 'registered') === true) {

                                        // User previously registered IBM Verify 2fa
                                        // and successfully authenticated with IBMId
                                        // so redirect user to the portal with IBMId
                                        // Token and call push notification to get
                                        // firebase token to login via firebase

                                        // redirect to a push notification page in portal
                                        // portal will call firebaseToken endpoint to send custom
                                        // firebase token down to the ui

                                        // console.info('Already 2fa registered IBMId cookie: ', req.headers.cookie);

                                        resolve(this.redirect(
                                            req,
                                            '2fa/verify',
                                            null
                                        ));

                                    }

                                    // If NOT previously registered 2FA - create new user with Username
                                    // if twoFactor.registered is undefined or false then the user has not register yet
                                    else if (isUndefined(get(twoFactor, 'registered')) || !get(twoFactor, 'registered')) {
                                        // 4: Require the user to confirm 2FA code to register

                                        // store temporary QR Code in firebase until
                                        // user successfully logs in for the first time
                                        // console.log('Registering 2FA for username');

                                        // register using IBM Verify 2fa
                                        // const twoFactorOperations = new TwoFactorOperations();
                                        // twoFactorOperations.initRegistration(firebaseUid, req)
                                        //     .then(result => {

                                        //         if (result.success === true) {

                                        //             // redirect to 2fa register page
                                        //             // this.getFirebaseToken(firebaseUid)
                                        //             // .then((token: string) => {
                                        //             resolve(this.redirect(
                                        //                 req,
                                        //                 '2fa/register',
                                        //                 result.data
                                        //             ));
                                        //             // });

                                        //         } else {
                                        //             res.status(500);
                                        //             reject(result.msg);
                                        //         }

                                        //     }, (err: any) => {
                                        //         console.error(err);
                                        //         res.status(500);
                                        //         reject('2fa registration failed');
                                        //     });

                                        const registerData: TOTPRegistrationData = {
                                            email: req.user._json.email
                                        };

                                        // encode sensitive personal data used for generating the QR Code
                                        // so that it is not exposed in the URL
                                        const registerDataEncoded = {
                                            data: Buffer.from(JSON.stringify(registerData)).toString('base64')
                                        };

                                        console.log(registerDataEncoded);

                                        resolve(this.redirect(
                                            req,
                                            '2fa/register',
                                            registerDataEncoded
                                        ));
                                    }

                                });

                        }

                    });

            } else {
                // error:
                // This code should never be reached because
                // this.ensureAuthenticated as middleware should always check to see if
                // the user exists before getting to this point. If this code is reached
                // something is likely wrong with the IBMId Passport Strategy
                res.status(500);
                reject('unable to identify user');
            }

        });

    }

    /**
     * login from portal
     *
     * @param {IPassportRequest} req
     * @returns {Promise<any>}
     * @memberof AuthController
     */
    @Get('portal-login')
    @OperationId('authSSOPortalLogin')
    // @Security('api_key')
    public async portalToken(@Request() req: IPassportRequest): Promise<any> {

        return new Promise((resolve, reject) => {

            const res = (<any>req).res as Response;

            if (req.user) {

                try {
                    // get firebase UserId from IBMId Token
                    this.authHelpers.setFirebaseUser(req.user._json.email)
                        .then((firebaseUid: string) => {

                            // send custom firebase token to be consumed by the
                            // portal to login the user to a portal session
                            this.getFirebaseToken(firebaseUid, this.authHelpers)
                                .then((token: string) => {
                                    // resolve(this.redirect(req, 'login', { token: token }));
                                    resolve({ token: token });
                                });

                        });
                } catch (error) {

                    console.error('Unable to authenticate user', error);
                    this.setStatus(401);
                    res.send('Unable to authenticate user');

                }

            }

        });

    }

    @Post('portal-login-totp')
    public async portalTokenTOTP(@Request() req: IPassportRequest): Promise<any> {

        return new Promise((resolve, reject) => {

            const res = (<any>req).res as Response;

            if (req.user) {

                try {
                    // get firebase UserId from IBMId Token
                    this.authHelpers.setFirebaseUser(req.user._json.email)
                        .then((firebaseUid: string) => {

                            // send custom firebase token to be consumed by the
                            // portal to login the user to a portal session
                            this.getFirebaseToken(firebaseUid, this.authHelpers)
                                .then((token: string) => {
                                    // resolve(this.redirect(req, 'login', { token: token }));
                                    resolve({ token: token });
                                    // res.send({ token: token })
                                });

                        });
                } catch (error) {

                    console.error('Unable to authenticate user', error);
                    this.setStatus(401);
                    res.send('Unable to authenticate user');

                }

            }

        });

    }

    /**
     * Used by IBMId as the specified callback upon successful login
     * Consider limiting CORs and IP to the IBMId server to thwart against attacks
     * API Key cannot be used on this route since the caller is the IBMId server
     * and setting the api key in the request to the callback is not possible
     *
     * @returns {Promise<any>}
     * @memberof IBMIdController
     */
    @Get('callback')
    @OperationId('authSSOCallback')
    public async callback(): Promise<any> {
       // unreachable code here - see auth.middleware.ts
       // added this method here ONLY for tsoa definition purposes
    }

    /**
     * Logout of the current IBMid user session
     *
     * @param {IPassportRequest} req
     * @returns {Promise<any>}
     * @memberof AuthController
     */
    @Post('logout')
    @OperationId('authSSOLogout')
    // @Security('api_key')
    public async logout(@Request() req: IPassportRequest): Promise<any> {
        req.logout();
        const res = (<any>req).res as Response;
        // res.redirect('/');  // TODO: redirect to a successful logout page
        res.send('Logout success!');
    }

    /**
     * Called internally by redirect if login failure is experienced
     *
     * @param {IPassportRequest} req
     * @returns {Promise<any>}
     * @memberof AuthController
     */
    @Post('failure')
    @OperationId('authSSOFailure')
    // @Security('api_key')
    public async failure(@Request() req: IPassportRequest): Promise<any> {
        const res = (<any>req).res as Response;
        res.send('login failed');
    }

    private async getFirebaseToken(firebaseUid: string, authHelpers: AuthHelpers): Promise<string> {

        // create custom firebase token with claims to view/set
        return await authHelpers.createCustomFirebaseToken(firebaseUid, {})
            .then((token: string) => {

                return token;

            }).catch((err: Error) => {

                // Failed to get authentication token for user (firebase auth might be down?)
                // throw {
                //     message: 'Login failed. Unable to set token:' + err.message,
                //     status: 500,
                // };

                // this.setStatus(500);
                // // Failed to get authentication token for user (firebase auth might be down?)
                // return 'Unable to set auth token:' + err.message;

                // res.status(500);
                // res.send('Unable to set auth token:' + err.message);

                throw new Error('Unable to set auth token:' + err.message);

            });
    }

    private async redirect(req: IPassportRequest, redirect: string, queryObj: any) {

        const res = (<any>req).res as Response;

        if (!isEmpty(queryObj)) {

            // create redirect with query string
            const queryString = Object.keys(queryObj)
                .map(key => key + '=' + queryObj[key])
                .join('&');

            res.status(200);
            res.redirect(this.env.site_root + '/' + redirect + '/?' + queryString);

        } else {

            // create redirect without query string
            res.status(200);
            res.redirect(this.env.site_root + '/' + redirect);

        }

    }

}
