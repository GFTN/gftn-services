// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// PURPOSE: these endpoints provide the means to request and generate
// a jwt token for a participant stack. A description of how security keys
// are secured is detailed here: https://github.com/GFTN/gftn-services/issues/77#issuecomment-10976927.
// This process implements a maker (aka: requester), checker (aka: approver) flow to
// issue a jwt token. jwt tokens can only be issued to users who belong to a institution from which the
// participantId was generated via the ./src/controllers/automation.controller.ts.

// Other Docs:
// https://medium.com/@maison.moa/using-jwt-json-web-tokens-to-authorize-users-and-protect-api-routes-3e04a1453c3e
// JWS = https://tools.ietf.org/html/rfc7515
// JWT = https://tools.ietf.org/html/rfc7519

import {
    Route, Controller, Request, Post, Body, Get, Header, OperationId
    // , Security
} from 'tsoa';
import { IJWTRequest, IFirebaseUserRequest, IVerifyCompare, IRandomPepperObj } from '../models/auth.model';
import { size, isEmpty, merge, isUndefined, map, cloneDeep, forEach, set, keys, get } from 'lodash';
import * as admin from 'firebase-admin';
import { IJWTTokenInfoPublic, IJWTPublic, IJWTTokenInfoGeneratedPublic, IJWTTokenClaimsSecure, IJWTTokenClaimsAndPayloadSecure, IJWTTokenPayloadSecure, IWWJWTSecure } from '../shared/models/token.interface';
import * as moment from 'moment';
import { AuthHelpers } from '../auth/auth-helpers';
import { sign, verify } from 'jsonwebtoken';
import { INodeAutomation } from '../shared/models/node.interface';
import { IGlobalEnvs } from '../environment';

@Route('jwt')
export class JwtController extends Controller {

    public env: IGlobalEnvs = global['envs'];
    public refs = {
        info: 'jwt_info',
        secure: 'jwt_secure'
    } 
    authHelpers = new AuthHelpers();
    

    /**
     * Refreshes a JWT Token for an
     * initial session (default refresh time is defined as env variable)
     *
     * @param {IJWTRequest} _req
     * @returns {Promise<any>}
     * @memberof JwtController
     */
    @Post('refresh')
    @OperationId('jwtRefresh')
    // @Security('api_key')
    public async refresh(
        @Request() req: IJWTRequest,
        @Header('Authorization') authorization: string
    ): Promise<any> {

        try {

            const _req = req;

            // set firebase ref to look up and modify token info visible to the ui
            const ref_jwt_info = admin.database().ref(
                '/'+ this.refs.info +'/'
                + _req.decodedToken.sub
                + '/' +
                _req.decodedToken.jti
            );

            const data = await ref_jwt_info.once('value');

            // get the token's general information
            const token_info: IJWTPublic = data.val();

            // if token exists
            if (!isEmpty(token_info.jti)) {

                // only approve if stage is currently ready or approved
                if (token_info.stage === 'ready' || token_info.stage === 'initialized') {

                    // ================= Update general information about the token visible to the ui =================

                    // update token info
                    token_info.stage = 'initialized';
                    token_info.refreshedAt = Number(moment.utc().format('X'));
                    token_info.active = true;

                    // ================= Generate the new refresh token  =================

                    // increment high watermark
                    _req.decodedToken.n = _req.decodedToken.n + 1;

                    // now + n minutes
                    _req.decodedToken.exp = Number(moment.utc().format('X')) + (60000 * Number(this.env.refresh_mins));

                    // cannot use until n milliseconds earlier (allows a buffer for un-synced clocks)
                    _req.decodedToken.nbf = Number(moment.utc().format('X')) - 500;

                    //  ================= Write updated data to firebase  =================

                    // delete out "old" secure token data
                    await admin.database().ref('/'+ this.refs.secure +'/' + req.dbKey)
                        .remove();

                    // update secure "new" token data and generate newly refreshed encoded token
                    const refreshedToken = await this._generateToken(_req.decodedToken, this.env.refresh_mins + 'm');

                    // update general token info visible to the ui
                    await ref_jwt_info
                        .update(token_info);

                    // console.log(refreshedToken);

                    return refreshedToken;

                } else {
                    this.setStatus(404);
                    return 'token is not currently under review';
                }

            } else {
                this.setStatus(404);
                return 'unable to find token info';
            }

        } catch (error) {
            this.setStatus(500);
            return 'failed to refresh: ' + error;
        }

    }

    /**
     * Used to test a token to determine if it is valid given a defined
     * endpoint, ip, and/or account.
     *
     * @param {IJWTRequest} req
     * @param {{ endpoints: string[] }} body - "Endpoints" are used to check if the token provided is good for those specified endpoints
     * @returns {Promise<any>}
     * @memberof JwtController
     */
    @Post('verify')
    @OperationId('jwtVerify')
    // @Security('api_key')
    public async verify(
        @Request() req: IJWTRequest,
        @Body() body: IVerifyCompare,
        @Header('Authorization') authorization: string,
        @Header('x-fid') fid: string,
        @Header('x-iid') iid: string
    ): Promise<any> {

        const header = req.headers['authorization'];

        if (typeof header !== 'undefined') {

            const bearer = header.split(' ');
            const encodedToken = bearer[1];

            // decode token will error out if jwt is expired or
            try {

                // decode token header and get keys to find secret stored in db and env
                const tokenHeader: { alg: string; kid: string } = JSON.parse(Buffer.from(encodedToken.split('.')[0], "base64").toString('utf8'));
                const dbKey = tokenHeader.kid.split('.')[0];
                const pepperKey = tokenHeader.kid.split('.')[1];

                // get db /jwt_secure data and pepper secret
                const data = await admin.database().ref(this.refs.secure + '/' + dbKey).once('value');
                const jwt_secure: IWWJWTSecure = data.val();
                const pepperSecret = this.env.ww_jwt_pepper_obj.v[pepperKey];

                // if both secrets are available then continue
                if (!isEmpty(jwt_secure.s) && !isEmpty(pepperSecret)) {

                    // decode token
                    const decodedToken = await verify(
                        encodedToken,
                        jwt_secure.s +
                        pepperSecret
                    ) as any;

                    // default ip to actual request ip
                    let compareIncomingIp = req.connection.remoteAddress;
                    if (body.ip) {
                        // if IP is provided in body, then compare the encoded token
                        // ip to the one in the request body
                        compareIncomingIp = body.ip;
                    }

                    // run custom validation on developer token
                    const res = this.authHelpers.verifyWWTokenCustom(decodedToken, jwt_secure.i, jwt_secure.n, compareIncomingIp, body.endpoint, body.account);
                    if (res.pass) {
                        return 'Success! Token is valid for the supplied body parameters.';
                    } else {
                        this.setStatus(403);
                        return 'failed to pass one (or more) of the many token validation checks: ' + res.msg;
                    }

                } else {
                    this.setStatus(403);
                    return 'unable to get token keyId';
                }

            } catch (error) {
                this.setStatus(403);
                return 'failed to decode the token';
            }

        } else {
            this.setStatus(401);
            return 'no authentication token provided';
        }

    }

    /**
     * Request creation of a JWT Token, generates request for
     * approval to create token. NOTE: This endpoint does not
     * create a jwt token only 'request'. User must request the creation
     * of a token, and must be 'approved' by another admin user before the
     * token can be 'generated'.
     *
     * @param {IJWTTokenInfoPublic} body
     * @returns {Promise<void>}
     * @memberof JwtController
     */
    @Post('request')
    @OperationId('jwtRequest')
    // @Security('api_key')
    public async request(
        @Request() req: IFirebaseUserRequest,
        @Body() body: IJWTTokenInfoPublic,
        @Header('x-fid') fid: string,
        @Header('x-iid') iid: string
    ): Promise<string> {

        // set institution id
        const institutionId: string = req.headers['x-iid'] as string;

        if (isEmpty(institutionId)) {
            this.setStatus(400);
            return 'missing institution id';
        }

        // check if the participantId is associated with the institutionId
        const institutionOwnsParticipantId = await this.institutionOwnsParticipantId(institutionId, body.aud);
        if (!institutionOwnsParticipantId) {
            this.setStatus(400);
            return 'participantId must be associated with institution';
        }

        const jwtServerInfo: IJWTTokenInfoGeneratedPublic = {
            stage: 'review',
            active: false,
            createdAt: Number(moment.utc().format('X')),
            createdBy: req.email || 'system',
            // approve field will be set by approve endpoint
            approvedAt: null,
            approvedBy: null,
            // refreshedAt field will be set at refresh endpoint
            refreshedAt: null
        };

        // combine models
        const data: IJWTPublic = merge(body, jwtServerInfo, { sub: institutionId });

        // generate ref
        const newRef = await admin.database().ref(
            '/'+ this.refs.info +'/' +
            institutionId
        ).push();

        // set public claim token properties
        // get new key and set id to the
        data.jti = newRef.key;

        return newRef.set(data).then(() => {
            return 'successfully requested token. Your new jti: ' + data.jti;
        }, () => {
            this.setStatus(500);
            return 'unable to save token request';
        });

    }

    /**
     * Approve creation of a JWT Token
     *
     * @param {IJWTRequest} req
     * @param {{ endpoints: string[] }} body
     * @returns {Promise<any>}
     * @memberof JwtController
     */
    @Post('approve')
    @OperationId('jwtApprove')
    // @Security('api_key')
    public async approve(
        @Request() req: IFirebaseUserRequest,
        @Body() body: { jti: string },
        @Header('x-fid') fid: string,
        @Header('x-iid') iid: string
    ): Promise<string> {

        const self = this;

        // get the token by jti
        return new Promise((resolve, reject) => {

            // generate ref to look up token
            const jwt_infoRef = admin.database().ref('/'+ this.refs.info +'/' + req.headers['x-iid'] + '/' + body.jti);

            // get the token
            jwt_infoRef.once('value', (data: admin.database.DataSnapshot) => {

                const token_info: IJWTPublic = data.val();

                self.setHeader('Content-type', 'text/plain');

                // if token exists
                if (token_info) {

                    // only approve if stage is currently request
                    if (token_info.stage === 'review') {

                        // check to make sure the approver is not the same user (by email)
                        if (!isUndefined(req.email) && req.email !== token_info.createdBy) {

                            // ensure that the returned token request includes a jti
                            if (token_info.jti === body.jti) {

                                // update token info
                                token_info.stage = 'approved';
                                token_info.approvedAt = Number(moment.utc().format('X'));
                                token_info.approvedBy = req.email;

                                // save updated token info
                                jwt_infoRef.update(token_info).then(() => {
                                    resolve('Success, token approved.');
                                }, () => {
                                    this.setStatus(500);
                                    reject(new Error('Unable to approve token.'));
                                });

                            } else {
                                self.setStatus(403);
                                resolve('unknown token id');
                            }

                        } else {
                            self.setStatus(403);
                            resolve('Same user who created the token request cannot also approve.');
                        }

                    } else {
                        self.setStatus(403);
                        resolve('Token is not currently under review.');
                    }

                } else {
                    self.setStatus(404);
                    resolve('Token id info not found.');
                }

            }).catch(() => {
                reject(new Error('Unexpected database error'));
            });

        });

    }

    /**
     * Reject is an alias for revoke
     * NOTE: needed to create reject so that the action from
     * the ui calls reject and revoke separately - this allows
     * us to handle reject different in the future if we so choose
     *
     * @param {IFirebaseUserRequest} req
     * @param {{ jti: string }} body
     * @returns {Promise<any>}
     * @memberof JwtController
     */
    @Post('reject')
    @OperationId('jwtReject')
    // @Security('api_key')
    public async reject(
        @Request() req: IFirebaseUserRequest,
        @Body() body: { jti: string },
        @Header('x-fid') fid: string,
        @Header('x-iid') iid: string
    ): Promise<any> {
        return this.revoke(req, body, fid, iid);
    }

    /**
     * Invalidates a token session by removing it from the database (identified by ID)
     *
     * @param {IJWTRequest} req
     * @param {{ endpoints: string[] }} body
     * @returns {Promise<any>}
     * @memberof JwtController
     */
    @Post('revoke')
    @OperationId('jwtRevoke')
    // @Security('api_key')
    public async revoke(
        @Request() req: IFirebaseUserRequest,
        @Body() body: { jti: string },
        @Header('x-fid') fid: string,
        @Header('x-iid') iid: string
    ): Promise<any> {

        const self = this;

        // generate ref to look up token
        const jwt_infoRef = admin.database().ref('/'+ this.refs.info +'/' + req.headers['x-iid'] + '/' + body.jti);

        // remove secure jwt
        await admin.database()
            .ref('/'+ this.refs.secure +'/' + req.headers['x-iid'] + '/' + body.jti)
            .remove();

        // get the token
        const data = await jwt_infoRef.once('value');

        const token_info: IJWTPublic = data.val();

        // if token exists
        if (token_info) {

            // update token info
            token_info.stage = 'revoked';
            token_info.revokedAt = Number(moment.utc().format('X'));
            token_info.revokedBy = req.email;

            // save updated token info
            return await jwt_infoRef.update(token_info).then(() => {
                return 'Success. Token revoked.';
            }, () => {
                this.setStatus(500);
                return 'Unable to approve token.';
            });

        } else {
            self.setStatus(404);
            return 'Token id info not found.';
        }

    }

    /**
     * A JWT token can only be generated by the requestor after an approver
     * approved the creation of the token
     *
     * @param {IFirebaseUserRequest} req
     * @param {{ jti: string }} body
     * @returns {Promise<any>}
     * @memberof JwtController
     */
    @Post('generate')
    @OperationId('jwtGenerate')
    // @Security('api_key')
    public async generate(
        @Request() req: IFirebaseUserRequest,
        @Body() body: { jti: string },
        @Header('x-fid') fid: string,
        @Header('x-iid') iid: string
    ): Promise<any> {

        const self = this;

        // generate ref to look up token
        const jwt_infoRef = admin.database().ref('/'+ this.refs.info +'/' + req.headers['x-iid'] + '/' + body.jti);

        // get the token
        const data = await jwt_infoRef.once('value');

        const token_info: IJWTPublic = data.val();

        // if token exists
        if (token_info) {

            // only generate if stage is currently approved
            if (token_info.stage === 'approved') {

                // check to make sure the authenticated user is the same user who requested the token
                if (!isUndefined(req.email) && req.email === token_info.createdBy) {

                    // ensure that the approved request includes a jti
                    if (token_info.jti === body.jti) {

                        const count = 0;

                        // define payload
                        const payload: IJWTTokenPayloadSecure = {
                            acc: token_info.acc,
                            ver: token_info.ver,
                            ips: token_info.ips,
                            env: token_info.env,
                            enp: token_info.enp,
                            n: count
                        };

                        // define claims
                        const claims: IJWTTokenClaimsSecure = {
                            jti: token_info.jti,
                            aud: token_info.aud,
                            sub: req.headers['x-iid'] as string
                        };

                        // update token info
                        token_info.stage = 'ready';

                        try {

                            let init_exp = this.env.initial_mins + 'm'

                            // set default expiration time
                            if (!init_exp) {
                                // default: initial expiration time (format https://github.com/zeit/ms)
                                init_exp = '1h';
                            }

                            // generate the token with payload and claims
                            // initialize to expire in n1 hrs and not before n2 seconds from now
                            const encodedToken = await self._generateToken(merge(payload, claims), init_exp, '0s');

                            // save updated token info
                            return await jwt_infoRef.update(token_info).then(async () => {
                                return encodedToken;
                            }, () => {
                                this.setStatus(500);
                                return 'unable to initialize newly generate token';
                            });

                        } catch (error) {
                            return 'Unable to generate token.';
                        }

                    } else {
                        self.setStatus(403);
                        return 'Unknown token id';
                    }

                } else {
                    self.setStatus(403);
                    return 'User who requested the token must be the same user to generate the token.';
                }

            } else {
                self.setStatus(403);
                return 'Token is not currently approved.';
            }

        } else {
            self.setStatus(404);
            return 'Token info not found.';
        }

    }

    /**
     * Used to randomly generate a new set of pepper values for the aws-secret store
     * DISCUSSION: These values can be de-serialized and referenced using keyId
     * and '.' to split on the keyid for looking up the value from the database
     * and looking up the correct env var. This allows us to rotate the pepper
     * secrets because we can append to the old list of pepper values to
     * the new result of randomly generate pepper values and then periodically (ie: 15 minutes)
     * when all the a values are presumed to be expired we can delete our the old values and rotate again.
     * Thereby preventing any disruption to service but effectively rotating secret keys
     *
     * @returns {Promise<any>}
     * @memberof JwtController
     */
    @Get('rotate-pepper')
    @OperationId('jwtRotatePepper')
    // @Security('api_key')
    public async rotatePepper(
        @Header('x-fid') fid: string
    ): Promise<any> {

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
        const pepperObj = this.env.ww_jwt_pepper_obj || init;

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
            const val = this.authHelpers.randStr(len, false);

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
        // by un-commenting below line
        // console.log(newEnvVar);

        // check if the length is less than 4096
        // equation => old values + new values = 2 x n
        if (newEnvVar.length <= 4096) {

            return {
                msg: 'success, rotated keys',
                data: pepperObj
            };
        } else {
            this.setStatus(500);
            return 'failed to create a env value less than 4096 characters';
        }

    }


    /**
     * Checks if a participantId is associated with an institution
     *
     * @private
     * @param {string} institutionId
     * @param {string} env
     * @param {string} participantId
     * @returns {Promise<boolean>}
     * @memberof JwtController
     */
    private async institutionOwnsParticipantId(institutionId: string, participantId: string): Promise<boolean> {

        // check if the selected institutionId contains the selected participantId
        // ie: prevents a user from generating a token for a
        // participantId for which they don't have access rights
        // NOTE 1: for this endpoint to issue a token it assumes that the stack details
        // exist in firebase under /participants/{participantId}/nodes/{env}/...
        // NOTE 2: middleware prevents a user that is not associate with this institution
        // from calling this endpoint.
        const _nodeDetails = await admin.database().ref('participants/' + institutionId + '/nodes/' + participantId)
            .once('value');

        const nodeDetails = _nodeDetails.val() as INodeAutomation;

        // if nodeDetails is not null or empty, and participantId of the returned recorded for the institution matches return true
        if (get(nodeDetails, 'participantId') === participantId && !isEmpty(nodeDetails)) {
            // participantId is associated with this institutionId
            return true;
        } else {
            // participantId does not belong to the institution
            return false;
        }

    }

    /**
     * Generates a JWT Token
     *
     * @private
     * @returns {string}
     * @memberof JwtController
     */
    private async _generateToken(decodedTokenValues: IJWTTokenClaimsAndPayloadSecure, exp?: string, nbf?: string): Promise<string> {

        // save a randomly generated secret key in db (found by keyid in token header)
        const dbSecret = this.authHelpers.randStr(64, false);

        // return a random number between min and max (min and max included)
        const genRandNum = (min: number, max: number) => {
            return Math.floor(Math.random() * (max - min + 1) + min);
        };

        const randNum: number = genRandNum(0, size(this.env.ww_jwt_pepper_obj.v) - 1);

        // convert pepper values to an array w/ keys intact
        // NOTE: don't use _.toArray() since this will not preserve keys in the output array
        const arrPepper = map(this.env.ww_jwt_pepper_obj.v, (val, key) => {
            // returns object array
            return { [key]: val };
        });

        const selectedPepper = arrPepper[randNum];
        const pepperKey = keys(selectedPepper)[0];
        const pepperSecret = selectedPepper[pepperKey];

        // generate ref
        const newRef = await admin.database().ref(
            '/'+ this.refs.secure +'/'
        ).push();

        // set public claim token properties
        // get new key and set id to the
        const dbKey = newRef.key;

        // write to firebase the jti and token secret
        if (!isEmpty(decodedTokenValues.jti)) {
            const data: IWWJWTSecure = {
                s: dbSecret,
                i: decodedTokenValues.jti,
                n: decodedTokenValues.n
            };
            // set data to check against in firebase
            await newRef.set(data);
        }

        let _exp = exp;
        let _nbf = nbf;

        // set default expiration time
        if (!exp) {
            // default: expiration time (format https://github.com/zeit/ms)
            _exp = '15m';
        }

        // set default not before time
        if (!nbf) {
            // not before now (format https://github.com/zeit/ms)
            _nbf = '0s';
        }

        const payload: IJWTTokenPayloadSecure = {
            acc: decodedTokenValues.acc,
            ver: decodedTokenValues.ver,
            ips: decodedTokenValues.ips,
            env: decodedTokenValues.env,
            enp: decodedTokenValues.enp,
            // NOTE: increment watermark for token before decoding
            n: decodedTokenValues.n
        };

        // create token
        // === Synchronous Sign with RSA SHA256 ===
        const encodedToken = sign(
            payload,
            dbSecret + // secret stored in db
            pepperSecret, // pepper secret stored in env var
            {
                keyid: dbKey + '.' + pepperKey, // changes for every new token
                jwtid: decodedTokenValues.jti, // stays the same for session
                audience: decodedTokenValues.aud, //'urn:' + decodedToken.aud,
                subject: decodedTokenValues.sub, // 'urn:' + decodedToken.sub,
                expiresIn: _exp,
                notBefore: _nbf,
            }
        );

        return encodedToken;

    }

}
