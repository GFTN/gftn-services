// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Request, Response, Application, NextFunction } from 'express';
import * as passport from 'passport';
import { IFirebaseUserRequest, RequestEnv, IPassportUser } from '../models/auth.model';
// import { TwoFactorOperations } from '../controllers/two-factor.controller';
import * as admin from 'firebase-admin';
import { IRoles } from '../shared/models/user.interface';
import { isEmpty, get, set } from 'lodash';
import { AuthHelpers } from '../auth/auth-helpers';
import * as jwt from 'jsonwebtoken';
import { IWWJWTSecure } from '../shared/models/token.interface';
import { TOTPTwoFactor } from '../controllers/totp-two-factor.controller';
import { IRoutePermissions, RoutePermissions } from './route-user-permissions.constants';
import { IGlobalEnvs } from '../environment';


interface IPassportRequest extends RequestEnv {
    session?: any;
    user?: IPassportUser;
    isAuthenticated(): boolean;
}

export function AuthMiddleware(app: Application) {

    const env: IGlobalEnvs = global['envs'];

    const authHelpers = new AuthHelpers();

    // /**
    //  *  uses authenticated firebase user to find the email to check 2fa
    //  *
    //  * @param {IFirebaseUserRequest} _req
    //  * @param {Response} res
    //  * @param {NextFunction} next
    //  * @returns
    //  */
    // const confirm2FAPushWithFirebaseEmail = async (req: Request, res: Response, next: NextFunction) => {

    //     const _req = req as any as IFirebaseUserRequest;

    //     if (env.enable_2fa === 'true') {

    //         // send 2fa push notification to IBM Verify App
    //         await new TwoFactorOperations().push(_req.email, env, req.connection.remoteAddress)
    //             .then((result: { success: boolean, msg: string }) => {

    //                 if (result.success === true) {

    //                     next();

    //                 } else {
    //                     res.status(401);
    //                     res.send('two-factor authentication failed');
    //                 }

    //             }, () => {
    //                 res.status(401);
    //                 res.send('two-factor authentication failed');
    //             });

    //     } else {
    //         // skip 2fa
    //         next();
    //     }

    // };

    // /**
    //  * uses authenticated IBMId user to find the email to check 2fa
    //  *
    //  * @param {Request} _req
    //  * @param {Response} res
    //  * @param {NextFunction} next
    //  */
    // const confirm2FAPushWithIBMIdEmail = async (_req: Request, res: Response, next: NextFunction) => {

    //     const req = _req as any as IPassportRequest;

    //     if (env.enable_2fa === 'true') {

    //         // send 2fa push notification to IBM Verify App
    //         await new TwoFactorOperations().push(req.user.emailAddress, env, req.connection.remoteAddress)
    //             .then((result: { success: boolean, msg: string }) => {

    //                 if (result.success === true) {

    //                     next();

    //                 } else {
    //                     res.status(401);
    //                     res.send('two-factor authentication failed');
    //                 }

    //             }, () => {
    //                 res.status(401);
    //                 res.send('two-factor authentication failed');
    //             });

    //     } else {
    //         // skip 2fa
    //         next();
    //     }

    // };

    const authenticateIbmIdRedirect = (_req: Request, res: Response, next: NextFunction) => {

        const req = _req as any as IPassportRequest;

        // console.info('Middleware cookie (authenticate(), redirect:' + req.session.originalUrl + '): ', req.headers.cookie);

        // "/sso/token"
        const redirect_url = req.session.originalUrl;

        // A redirect is commonly issued after authenticating a request.
        // In this case, the redirect options override the default behavior.
        // Upon successful authentication, the user will be redirected
        // to the originating url. If authentication fails, the user will be
        // redirected back to failure endpoint to deal with the failed attempt.
        passport.authenticate('openidconnect', {
            successRedirect: redirect_url,
            failureRedirect: '/failure',
        })(req, res, next);

    };

    const authenticateIbmId = (req: Request, res: Response, next: NextFunction) => {

        // IBMId Auth Token:
        // console.info('Middleware cookie (/sso/token): ', req.headers.cookie);

        // use passport middleware used to determine if user is already signed-in
        if (!req.isAuthenticated()) {

            // IMPORTANT: must redirect to login to get IBMId because of OpenId CORS protections
            req.session.originalUrl = req['originalUrl'];
            res.redirect('/sso/login');

        } else {
            return next();
        }

    };

    /**
     * Used to associate a request from the portal
     * with a user to determine permissions or if
     * and action can be taken
     */
    const authenticateFirebaseUser = async (req: Request, res: Response, next: NextFunction) => {

        try {

            // cast type to set properties
            const _req = req as any as IFirebaseUserRequest;

            // check if header exists
            if (isEmpty(_req.headers['x-fid'])) {

                // missing header
                res.status(401);
                res.send('unauthorized');

            } else {

                // verify the firebase token
                const decodedToken: admin.auth.DecodedIdToken = await admin.auth().verifyIdToken(_req.headers['x-fid'] as string)
                    .catch((error) => {
                        // not able to verify firebase user id token
                        // res.status(401);
                        // res.send('fid token is has expired');
                        console.error('fid token is not valid, check environment :' + error);
                        throw new Error('fid token is not valid');
                    });

                if (!isEmpty(get(decodedToken, 'uid'))) {
                    // set uid
                    _req['uid'] = decodedToken.uid;

                    // set email too for push notification
                    const user: admin.auth.UserRecord = await admin.auth().getUser(decodedToken.uid)
                        .catch(() => {
                            // res.status(401);
                            // res.send('unknown user id');
                            throw new Error('unknown user id');
                        });

                    if (user) {
                        _req['email'] = user.email;
                        next();
                    }

                }

            }

        } catch (error) {
            res.status(401);
            res.send(error.message);
        }

    };

    const checkPermissions = async (req: Request, res: Response, next: NextFunction) => {

        // set required permissions for routes
        // describes firebase route permissions required for continuing to route
        const routePermissions: IRoutePermissions = RoutePermissions;

        // get required permissions for accessing this route
        const requiredPermissionsForAccessingRoute = routePermissions[req.route.path];

        // default access init to false
        let allowAccess = false;

        let uid = '';

        if (isEmpty(req['uid'])) {
            // get uid by the signed in IBMId user
            const data = await admin.auth().getUserByEmail(req.user['emailAddress']);
            uid = data.uid;
        } else {
            // get uid by the already set by firebase user
            uid = req['uid'];
        }

        /**
         * gets super users permission roles from firebase database (as currently defined)
         *
         * @param {IPassportRequest} req
         * @returns {Promise<IRoles>}
         */
        const setUserSuperPermissions = async (): Promise<IRoles> => {

            return await admin.database().ref(
                'super_permissions/' +
                uid
            ).once('value').then((data: admin.database.DataSnapshot) => {

                const superPermissions: IRoles = data.val() as IRoles;

                // set permissions on the request object so that they can be checked
                // in more granularity in the endpoint logic
                req['super_permissions'] = superPermissions;

                return superPermissions;

            });

        };

        /**
        * gets participant users permission roles from firebase database (as currently defined)
        *
        * @param {IFirebaseUserRequest} req
        * @returns {Promise<IRoles>}
        */
        const setUserParticipantPermissions = async (): Promise<IRoles> => {

            return await admin.database().ref(
                'participant_permissions/' +
                uid + '/' +
                req['iid']
            ).once('value').then((data: admin.database.DataSnapshot) => {

                const participantPermissions: IRoles = data.val() as IRoles;

                // set permissions on the request object so that they can be checked
                // in more granularity in the endpoint logic
                req['participant_permissions'] = participantPermissions;

                return participantPermissions;

            });

        };

        // check if the route allows users with super permissions to access route
        if (get(requiredPermissionsForAccessingRoute, 'super_permissions')) {

            // get super Permissions
            const superPermissions = await setUserSuperPermissions();
            if (superPermissions) {

                // check if the user has been provisioned a super_permissions role in firebase
                // that matches one of the required permissions set levels required for this route
                for (let i = 0; i < requiredPermissionsForAccessingRoute.super_permissions.length; i++) {
                    if (get(req['super_permissions'], 'roles.' + requiredPermissionsForAccessingRoute.super_permissions[i])) {
                        allowAccess = true;
                    }
                }

            }

        }

        // if super permissions have already permitted access to this route then skip checking participant permissions
        if (allowAccess === false) {

            // user was not granted super permissions to access this route, so check if they have permissable participant permissions

            // check if the route allows users with participant permissions to access route
            if (get(requiredPermissionsForAccessingRoute, 'participant_permissions')) {

                // iid header is required to determine if the user has access rights for a specific institution "participant" permissions
                // iid = institution id which is required
                if (isEmpty(req.headers['x-iid'])) {

                    res.status(400);
                    res.send('unknown institution id');

                } else {

                    // set iid on request (consistent with how developer token will set iid on req)
                    req['iid'] = req.headers['x-iid'];

                    // get participant permissions
                    const participantPermissions = await setUserParticipantPermissions();
                    if (participantPermissions) {

                        // check if the user has been provisioned a participant_permissions role in firebase
                        // that matches one of the required permissions set levels required for this route
                        for (let i = 0; i < requiredPermissionsForAccessingRoute.participant_permissions.length; i++) {
                            if (get(req['participant_permissions'], 'roles.' + requiredPermissionsForAccessingRoute.participant_permissions[i])) {
                                allowAccess = true;
                            }
                        }

                    }

                }

            }

        }

        // evaluate if the user has access rights to proceed
        if (allowAccess) {
            // user has access rights to access route
            next();
        } else {
            // user does not have access rights to access route
            res.status(403);
            res.send('no user permissions to manage this institution');
        }

    };

    /**
     * IMPORTANT: This middleware will not be used for production.
     * This is only used by /verify and /refresh which are "proof-of-concept" endpoint.
     * Instead this will be implemented in the golang code.
     * Checks the validity of a token and returns the data using the
     * secret key for Developer Access
     *
     * @param {IJWTRequest} req
     * @memberof JwtController
     */
    const decodeWWToken = async (req: Request, res: Response, next: NextFunction) => {

        const header = req.headers['authorization'];

        // console.log('Disecting your request:', req);

        if (typeof header !== 'undefined') {

            const bearer = header.split(' ');
            const encodedToken = bearer[1];

            // jwt.verify() will throw error is basic validation fails, such as expired or nfb fails
            try {

                // decode token header and get keys to find secret stored in db and env
                const tokenHeader: { alg: string; kid: string } = JSON.parse(Buffer.from(encodedToken.split('.')[0], "base64").toString('utf8'));
                const dbKey = tokenHeader.kid.split('.')[0];
                // setting dbKey on ref so that the old token info can be deleted from the db
                req['dbKey'] = dbKey;
                const pepperKey = tokenHeader.kid.split('.')[1];

                // get db /jwt_secure data and pepper secret
                const data = await admin.database().ref('jwt_secure/' + dbKey).once('value');
                const jwt_secure: IWWJWTSecure = data.val();
                const pepperSecret = env.ww_jwt_pepper_obj.v[pepperKey];

                // if both secrets are available then continue
                if (!isEmpty(get(jwt_secure, 's')) && !isEmpty(pepperSecret)) {

                    // decode token (will throw error if fails to verify)
                    const decodedToken = await jwt.verify(
                        encodedToken,
                        jwt_secure.s +
                        pepperSecret
                    ) as any;

                    let originatingIp = '';

                    if (env.build === 'prod') {

                        // Production = get originating ip from x-forwared-for header

                        // (ie: 'x-forwarded-for': '170.225.9.145, 172.217.6.84',)
                        const xForwardString = req.headers['x-forwarded-for'] as string;

                        // get the first ip address in the x-forwarded-for array
                        originatingIp = xForwardString.split(',').map(item => item.trim())[0];

                        if (isEmpty(originatingIp) || originatingIp === '0.0.0.0' || originatingIp === req.connection.remoteAddress) {
                            // ip address must exist, must not be all ips, and originatingIp must not equal loadbalancer (if one is being used)
                            res.status(403);
                            res.send('invalid token x1: failed to get request ip: originatingIp=' + originatingIp + ', remoteAddress=' + req.connection.remoteAddress);
                        }

                    } else {
                        // Development = set local ip to originating ip
                        originatingIp = req.connection.remoteAddress;
                    }

                    // check if the token passes validation checks
                    const clear = authHelpers.verifyWWTokenCustom(decodedToken, jwt_secure.i, jwt_secure.n, originatingIp);

                    // run custom validation on developer token
                    if (clear.pass) {

                        // set token on the request
                        req['decodedToken'] = decodedToken;

                        // set payload vars on request object:

                        // set uid = for permissions
                        if (!isEmpty(get(decodedToken, 'uid'))) {
                            req['uid'] = decodedToken.uid;
                        }

                        // set iid = institutionId for permissions
                        if (!isEmpty(get(decodedToken, 'iid'))) {
                            req['iid'] = decodedToken.iid;
                        }

                        // set email for identifying a user
                        if (!isEmpty(get(decodedToken, 'email'))) {
                            req['email'] = decodedToken.email;
                            set(req, 'user.emailAddress', decodedToken.email);
                        }

                        next();

                    } else {
                        res.status(403);
                        res.send('invalid token x2: ' + clear.msg);
                    }

                } else {
                    res.status(403);
                    res.send('invalid token x3: missing decryption info');
                }

            } catch (error) {
                res.status(403);
                res.send('invalid token x4');
            }

        } else {
            res.status(401);
            res.send('invalid token x5: no authentication token provided');
        }

    };

    //TOTP
    const TOTPTwoFactorObj = new TOTPTwoFactor();

    // middleware to specific routes
    app.post('/permissions/participant', authenticateFirebaseUser, checkPermissions, TOTPTwoFactorObj.checkTOTPMiddleWareFirebaseUser);
    app.post('/permissions/super', authenticateFirebaseUser, checkPermissions, TOTPTwoFactorObj.checkTOTPMiddleWareFirebaseUser);
    app.post('/jwt/request', authenticateFirebaseUser, checkPermissions, TOTPTwoFactorObj.checkTOTPMiddleWareFirebaseUser);
    app.post('/jwt/approve', authenticateFirebaseUser, checkPermissions, TOTPTwoFactorObj.checkTOTPMiddleWareFirebaseUser);
    app.post('/jwt/revoke', authenticateFirebaseUser, checkPermissions, TOTPTwoFactorObj.checkTOTPMiddleWareFirebaseUser);
    app.post('/jwt/reject', authenticateFirebaseUser, checkPermissions, TOTPTwoFactorObj.checkTOTPMiddleWareFirebaseUser);
    app.post('/jwt/generate', authenticateFirebaseUser, checkPermissions, TOTPTwoFactorObj.checkTOTPMiddleWareFirebaseUser);
    app.post('/jwt/refresh', decodeWWToken);
    app.post('/jwt/verify', authenticateFirebaseUser, checkPermissions, TOTPTwoFactorObj.checkTOTPMiddleWareFirebaseUser);
    app.post('/jwt/rotate-pepper', authenticateFirebaseUser, checkPermissions, TOTPTwoFactorObj.checkTOTPMiddleWareFirebaseUser);
    // NOTE: No 'decodeWWToken' middleware needed on '/jwt/verify' because 'verifyWWTokenCustom()' is run in the endpoint logic directly
    // app.post('/jwt/verify');

    // IBMId login
    app.get('/sso/token', authenticateIbmId);
    // app.get('/sso/portal-login', authenticateIbmId, confirm2FAPushWithIBMIdEmail);
    app.get('/sso/callback', authenticateIbmIdRedirect);
    app.get('/sso/login', passport.authenticate('openidconnect'));

    // 2FA
    app.get('/2fa/push', authenticateIbmId);
    // register is handled in a redirect from IBMId /token endpoint

    app.get('/totp/:accountName', authenticateIbmId, TOTPTwoFactorObj.checkAccountNameMiddleWare);
    app.post('/totp/:accountName/confirm', authenticateIbmId, TOTPTwoFactorObj.checkAccountNameMiddleWare);
    // endpoint specifically for deployment services
    app.post('/totp/check', authenticateFirebaseUser);
    // endpoint pointing to /sso/portal-login with the difference of applied middleware
    app.post('/sso/portal-login-totp', authenticateIbmId, TOTPTwoFactorObj.checkTOTPMiddleWareIBMIdUser);
}
