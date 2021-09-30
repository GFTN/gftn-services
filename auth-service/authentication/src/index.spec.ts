// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { assert, expect } from 'chai';
import { JwtController } from './controllers/jwt.controller';
import *  as admin from 'firebase-admin';
import { Environment, IGlobalEnvs } from './environment';
import { IJWTTokenInfoPublic, IJWTTokenClaimsAndPayloadSecure } from './models/token.interface';
import { IInstitution } from './shared/models/participant.interface';
import { IFirebaseUserRequest } from './models/auth.model';
import * as _ from 'lodash';
import * as moment from 'moment';

describe('run auth-service unit tests', () => {

    let env: IGlobalEnvs;
    let db: admin.database.Database;
    const jwt = new JwtController();
    let mockParticipant: IInstitution;
    let req: IFirebaseUserRequest

    before(async () => {

        // // env overrides
        // process.env.initial_mins='15';
        // process.env.refresh_mins='15';

        // set environment variables to init firebase
        await new Environment().init();
        env = global['envs'];
        jwt.env = env; // set envs on the jwt.controller.ts

        // use a test database so dev data does get messed up
        env.firebaseConfig.databaseURL = 'https://dev-test-1b7da.firebaseio.com/'
        // initialize firebase app:
        admin.initializeApp(env.firebaseConfig);
        db = admin.database();

        const nodes = {
            participantid1: {
                bic: "IBMTESTA001",
                callbackUrl: "http://34.216.224.221:31002/v1/callback",
                countryCode: "USA",
                createAccount: true,
                createNode: true,
                createPREntry: true,
                institutionId: "institutionid1",
                participantId: "participantid1",
                rdoClientUrl: "http://34.216.224.221:11003/v1",
                role: "MM",
                status: ["complete"],
                version: "latest"
            }
        } as any;

        mockParticipant = {
            info: {
                address1: "600 Anton Blvd",
                address2: "",
                city: "Costa Mesa",
                country: "USA",
                geo_lat: 0,
                geo_lon: 0,
                institutionId: "institutionid1",
                kind: "Money Transfer Operator",
                logo_url: "",
                name: "participantid1",
                site_url: "",
                slug: "slugid1",
                state: "California",
                status: "active",
                zip: "123123"
            },
            nodes: nodes,
            users: {
                userid1: {
                    profile: {
                        email: "user1@ibm.com"
                    },
                    roles: {
                        admin: true
                    }
                },
                userid2: {
                    profile: {
                        email: "user2@ibm.com"
                    },
                    roles: {
                        admin: true
                    }
                },
                userid3: {
                    profile: {
                        email: "user3@ibm.com"
                    },
                    roles: {
                        manager: true
                    }
                }
            }

        }

        req = {
            headers: {
                'x-iid': mockParticipant.info.institutionId
            },
            email: ''
        } as any;

        // seed the database with mock participant
        const ref = 'participants/' + mockParticipant.info.institutionId
        let res;
        await db.ref(ref).set(mockParticipant).then(async () => {

            await db.ref(ref).once("value", (snapshot) => {
                res = snapshot.val();
            });

            if (mockParticipant.info.institutionId !== res.info.institutionId) {
                console.error('firebase envs must be misconfigured, check application (or launch.json) envs')
                process.exit(1)
            }

        });

    });

    let requestedJWT: IJWTTokenInfoPublic;

    describe('successfully complete JWT lifecycle', async () => {

        let jti: string;

        it('user requests creation of JWT (JWT Request)', async () => {

            req.email = _.toArray(mockParticipant.users)[0].profile.email

            requestedJWT = {
                jti: 'willBeReplacedLater',
                aud: _.toArray(mockParticipant.nodes)[0].participantId,
                description: 'unit test jwt',
                // allowable 'human readable' stellar accounts
                acc: ['default', 'issuing'],
                // versions
                ver: '',
                // allowable ips
                ips: ['0.0.0.0', '1.1.1.1'],
                // environment this token can be used for.
                // generally pulled from environment variable
                env: 'dev',
                enp: ['/testendpoint1', '/testendpoint2']
            }

            const res = await jwt.request(req as any, requestedJWT, '', '')

            const msg = 'successfully requested token. Your new jti: '
            const containsMsg = res.includes(msg);
            assert.equal(containsMsg, true);

            // set the jti for approval flow
            if (containsMsg) {
                jti = res.split(msg)[1]
            }

        })

        it('requesting user should not be able to approve creation of JWT (JWT approval)', async () => {

            let res;
            req.email = _.toArray(mockParticipant.users)[0].profile.email

            try {
                res = await jwt.approve(req as any, { jti: jti }, '', '')
            } catch (error) {
                res = error
            }

            assert.equal(res, 'same user who created the token request cannot also approve');

        })

        it('another user approves creation of JWT (JWT approval)', async () => {

            req.email = _.toArray(mockParticipant.users)[1].profile.email

            const res = await jwt.approve(req as any, { jti: jti }, '', '')

            assert.equal(res, 'success, token approved');

        })

        it('approver user should not be able to generate one-time JWT (JWT creation)', async () => {

            let res;
            req.email = _.toArray(mockParticipant.users)[1].profile.email;

            try {
                res = await jwt.generate(req as any, { jti: jti }, '', '')
            } catch (error) {
                res = error
            }

            assert.equal(res, 'user who requested the token must be the same user to generate the token');

        })

        let encodedToken: string;
        let decodedTokenBody: IJWTTokenClaimsAndPayloadSecure;
        let jwtCreatedAt = moment();
        it('requesting user should generate one-time JWT (JWT creation)', async () => {

            let res;
            req.email = _.toArray(mockParticipant.users)[0].profile.email;
            let header: { alg: string, typ: string, kid: string };
            let signature: string;
            let secretLookup: string;
            let pepperLookup: string;

            try {
                encodedToken = await jwt.generate(req as any, { jti: jti }, '', '')
                const tokenParts = encodedToken.split('.')

                header = JSON.parse(Buffer.from(tokenParts[0], 'base64').toString('utf8'));
                decodedTokenBody = JSON.parse(Buffer.from(tokenParts[1], 'base64').toString('utf8'));
                signature = tokenParts[2];

                const kid = header.kid.split('.')
                secretLookup = kid[0];
                pepperLookup = kid[1];

            } catch (error) {
                res = error
                console.error(error)
            }

            assert.equal(decodedTokenBody.jti, jti);

        })

        it('expiration time should be properly calculated', async () => {

            
            // moment.unix(decodedTokenBody.exp).format("dddd, MMMM Do YYYY, h:mm:ss a")
            var date1 = moment.unix(decodedTokenBody.exp)
            const actualExpInMins = date1.diff(moment.now(), 'minutes')
            const expectedExpInMins = _.toNumber(env.initial_mins)
            
            // console.log(jwtCreatedAt);
            // console.log(expectedExpInMins);
            // console.log(actualExpInMins);

            expect(actualExpInMins).to.be.oneOf([expectedExpInMins, expectedExpInMins - 1]);
        })

        it('verify newly generated JWT', async () => {

            // // wait 2 seconds as the jwt takes 1 - 2 seconds to become valid
            // await new Promise((resolve, reject)=>{
            //     setTimeout(() => {
            //         resolve();
            //     }, 5000);
            // })

            let res;
            req.email = _.toArray(mockParticipant.users)[0].profile.email;
            req.headers.authorization = 'bearer ' + encodedToken;
            req.connection = { remoteAddress: requestedJWT.ips[0] } as any;
            const compare = {
                endpoint: requestedJWT.enp[0],
                ip: requestedJWT.ips[0],
                account: requestedJWT.acc[0]
            }

            try {
                res = await jwt.verify(req as any, compare, encodedToken, '', '')
            } catch (error) {
                res = error
                console.error(error)
            }

            assert.equal(res, 'Success! Token is valid for the supplied body parameters.');

        })

        it('mismatch ip for newly generated JWT should fail', async () => {

            let res;
            req.email = _.toArray(mockParticipant.users)[0].profile.email;
            req.headers.authorization = 'bearer ' + encodedToken;
            req.connection = { remoteAddress: requestedJWT.ips[0] } as any;
            const compare = {
                endpoint: requestedJWT.enp[0],
                ip: '9.9.9.9',
                account: requestedJWT.acc[0]
            }

            try {
                res = await jwt.verify(req as any, compare, encodedToken, '', '')
            } catch (error) {
                res = error
            }

            const failedVerification = res.includes('failed to pass one (or more) of the many token validation checks:')
            assert.equal(failedVerification, true);

        })

        it('mismatch account for newly generated JWT should fail', async () => {

            let res;
            req.email = _.toArray(mockParticipant.users)[0].profile.email;
            req.headers.authorization = 'bearer ' + encodedToken;
            req.connection = { remoteAddress: requestedJWT.ips[0] } as any;
            const compare = {
                endpoint: requestedJWT.enp[0],
                ip: requestedJWT.ips[0],
                account: 'bankaccount'
            }

            try {
                res = await jwt.verify(req as any, compare, encodedToken, '', '')
            } catch (error) {
                res = error
            }

            const failedVerification = res.includes('failed to pass one (or more) of the many token validation checks:')
            assert.equal(failedVerification, true);

        })

        it('mismatch endpoint for newly generated JWT should fail', async () => {

            let res;
            req.email = _.toArray(mockParticipant.users)[0].profile.email;
            req.headers.authorization = 'bearer ' + encodedToken;
            req.connection = { remoteAddress: requestedJWT.ips[0] } as any;
            const compare = {
                endpoint: '/somemaliciousendpoint',
                ip: requestedJWT.ips[0],
                account: requestedJWT.acc[0]
            }

            try {
                res = await jwt.verify(req as any, compare, encodedToken, '', '')
            } catch (error) {
                res = error
            }

            const failedVerification = res.includes('failed to pass one (or more) of the many token validation checks:')
            assert.equal(failedVerification, true);

        })



    });

    after(() => {
        process.exit(0);
    })

})

