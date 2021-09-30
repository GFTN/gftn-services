// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// PURPOSE: passport.js (http://www.passportjs.org/docs/) helper functions that are utilized by IBMId via ./src/controller/ibmid.controller.ts.

import { Application } from 'express';
import * as cookieParser from 'cookie-parser';
import * as bodyParser from 'body-parser';
import * as passport from 'passport';
import * as session from 'express-session';
import { IGlobalEnvs } from '../environment';

export class PassportConfig {

    private env: IGlobalEnvs = global['envs'];

    constructor(
        public app: Application
    ) {
        this.initPassport();
    }

    private initPassport() {

        // console.info('Initializing IBMId using PassportJs');

        // http://www.passportjs.org/docs/
        this.configMiddleware();
        this.configSessions();
        this.configStrategy();

    }

    /**
     *  http://www.passportjs.org/docs/configure/
     *
     * @private
     * @memberof PassportConfig
     */
    private configMiddleware() {

        // Parse Cookie header and populate req.cookies with an object keyed by the cookie 
        // names. Optionally you may enable signed cookie support by passing a secret 
        // string, which assigns req.secret so it may be used by other middleware.
        this.app.use(cookieParser(this.env.passport_secret));

        // Parse incoming request bodies in a middleware before your 
        // handlers, available under the req.body property.
        // this.app.use(bodyParser());
        this.app.use(bodyParser.json());

        // Create a session middleware with the given options.
        // Note Session data is not saved in the cookie itself, 
        // just the session ID. Session data is stored server-side.
        // Note Since version 1.5.0, the cookie-parser middleware no longer needs 
        // to be used for this module to work. This module now directly reads and 
        // writes cookies on req/res. Using cookie-parser may result in issues 
        // if the secret is not the same between this module and cookie-parser.
        this.app.use(session({
            resave: true,
            saveUninitialized: true,
            secret: this.env.passport_secret
        }));

        // In a Connect or Express-based application, passport.initialize() middleware
        // is required to initialize Passport. If your application uses persistent
        // login sessions, passport.session() middleware must also be used.
        this.app.use(passport.initialize());
        this.app.use(passport.session());

    }

    /**
     * http://www.passportjs.org/docs/configure/
     *
     * @private
     * @memberof PassportConfig
     */
    private configSessions() {

        // In a typical web application, the credentials used to authenticate a user will
        // only be transmitted during the login request. If authentication succeeds, a 
        // session will be established and maintained via a cookie set in the user's browser.
        // Each subsequent request will not contain credentials, but rather the unique cookie
        // that identifies the session. In order to support login sessions, Passport will 
        // serialize and deserialize user instances to and from the session.
        // The serialization and deserialization logic is supplied by the application, 
        // allowing the application to choose an appropriate database and/or object 
        // mapper, without imposition by the authentication layer.
        passport.serializeUser((user, done) => {
            done(null, user);
        });

        passport.deserializeUser((obj, done) => {
            done(null, obj);
        });
    }

    /**
     * http://www.passportjs.org/docs/configure/
     *
     * @private
     * @memberof PassportConfig
     */
    private configStrategy() {

        try {

            // do not precede certPath with / this important in gcloud file system
            // let certPath = '/tmp';
            // let certPath = '/.ibmid-certs/ibmid';

            // let ibmIdCert = '/idaas_iam_ibm_com.crt';

            // if (this.env.build === 'dev') {
            //     ibmIdCert = '/prepiam_toronto_ca_ibm_com.crt';
            // }

            // Passport uses what are termed strategies to authenticate requests. Strategies range 
            // from verifying a username and password, delegated authentication using OAuth or 
            // federated authentication using OpenID. Before asking Passport to authenticate a 
            // request, the strategy (or strategies) used by an application must be configured.
            const OpenIDConnectStrategy = require('passport-idaas-openidconnect-ww').IDaaSOIDCStrategy;
            const Strategy = new OpenIDConnectStrategy({
                authorizationURL: this.env.ibmid_authorization_url,
                tokenURL: this.env.ibmid_token_url,
                clientID: this.env.ibmid_client_id,
                scope: 'openid',
                response_type: 'code',
                clientSecret: this.env.ibmid_client_secret,
                callbackURL: this.env.ibmId_callback_url,
                skipUserProfile: true,
                issuer: this.env.ibmid_issuer_id,
                // according to https://developer.ibm.com/answers/questions/310837/passport-idaas-openidconnect/
                // after openId connect passport lib version update to 2.0.x ibm certs should be enabled and added
                // certs are found here: https://w3-connections.ibm.com/wikis/home?lang=en-us#!/wiki/W89b23bf7ad80_4411_822f_2a6dc171c6b3/page/Choosing%20a%20w3id%20or%20IBMid%20Provider?section=IBMidcerts
                addCACert: true,
                // CACertPathList: [
                //     certPath + '/digicert-root.pem',
                //     certPath + '/IBMid-server.crt',
                //     certPath + ibmIdCert
                // ],
                CACertPathListBase64: this.env.ibmId_encoded_certs
            }, (iss, sub, profile, accessToken, refreshToken, params, done) => {
                process.nextTick(() => {
                    profile.accessToken = accessToken;
                    profile.refreshToken = refreshToken;
                    // console.info('strategy tick profile: ', profile);
                    done(null, profile);
                });
            });

            passport.use(Strategy);

        } catch (error) {

            console.error('failed passport strategy: ', error);
            
        }
    }

}

