// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
//
// import * as assert from 'assert';
import { assert, expect } from 'chai';
import { after } from 'mocha';
// import { User } from './models/user';
import { ParticipantsModel } from '../src/app/shared/models/participants.model';
import { IInstitution } from '../src/app/shared/models/participant.interface';
import * as _ from 'lodash';
import * as admin from 'firebase-admin';
// import * as test from 'firebase-functions-test';
// import * as fs from 'fs';
import * as permissions from './src/triggers/permissions';

// IMPORTANT: **Never** use production credentials here
const config = {
    apiKey: "your.api.key.goes.here",
    authDomain: "your.project.id.firebaseapp.com",
    databaseURL: "https://your.project.id.firebaseio.com",
    projectId: "your.project.id",
    storageBucket: "your.project.id.appspot.com",
    messagingSenderId: "your.messaging.sender.id"
};

const serviceAccount = './../src/_next-gftn-firebase-adminsdk-wvjz8-67ea263932.json';

// There are two ways to use firebase-functions-test:

// Offline mode: Write siloed and offline unit tests with no side effects.
// This means that any method calls that interact with a Firebase product
// (e.g. writing to the database or creating a user) needs to be stubbed.
// At the top of test/index.test.js

const test = require('firebase-functions-test')(config, serviceAccount);

// Online mode: Write tests that interact with a Firebase project
// dedicated to testing so that database writes, user creates, etc.
// would actually happen, and your test code can inspect the results.
// This also means that other Google SDKs used in your functions will work as well.

// Ideally, we would want to use 'offline' unit testing
// for improved performance and isolation. However,
// since we are relying on firebase real-time database
// and our goal is to test simulated production behavior
// we have chosen to do 'online' unit testing. This will
// allow us to create new users and modify them in firebase's
// real time database and observe the resulting behavior.

// using staging credentials for testing
// const firebaseTest = test(config, serviceAccount);

// set environment variables:
// Using 'development' so that affects ONLY impact
// staging env (ie: next-gftn) firebase project
// firebaseTest.mockConfig({ env: { build: 'development' } });

// initializing app so that admin sdk can be used
// const app = admin.initializeApp({
//     credential: admin.credential.cert(require(serviceAccount)),
//     databaseURL: 'https://next-gftn.firebaseio.com'
// });

// // Test Example
// describe('sample test', () => {
//     let val = false;
//     val = true;
//     it('expect result to be true', () => {
//         assert.equal(val, true);
//     });
// });



permissions.writeParticipantPermissions('6TPij0ZJHMOq2PuqkH9Zjc4IXRA2', '-LPmrF0Qgg6fshzLc48i',
    {
        email: 'email@email.com',
        roles: { admin: true }
    });

// const triggers = require('./index');
const triggers = require('./triggers/permissions');

describe('firebase functions - http, cron job, and database triggers', () => {

    describe('create new firebase user and participant', () => {

        // mock user information
        const userProfile = {
            email: 'your.user@your.domain',
            uid: '' // set uid after creating new firebase user
        };

        // mock new participant
        const participantInfo: IInstitution = {
            info: {
                institutionId: '',
                slug: 'big-test-bank',
                status: 'pending',
                name: 'Big Test Bank',
                country: 'USA',
                address1: '600 Anton',
                address2: '',
                city: 'Costa Mesa',
                state: 'CA',
                zip: '92626',
                kind: 'Bank'
            }
        };

        describe('setUser', () => {

            // create new firebase user and participant to run tests against
            // const user = new User();

            // describe('create user', () => {

            //     it('should create new firebase user and resolve the new uid', async () => {

            //         await user.setUser(userProfile.email)
            //             .then((userId: any) => {
            //                 userProfile.uid = userId;
            //                 expect(userId).to.be.a('string');
            //             }, (err: any) => {
            //                 console.error(err);
            //             });

            //     });

            // });

            // describe('get existing user', () => {

            //     it('should look-up user by email address resolve the previously created uid', async () => {

            //         await user.setUser(userProfile.email)
            //             .then((userId: any) => {
            //                 expect(userId).to.equal(userProfile.uid);
            //             }, (err) => {
            //                 console.error(err);
            //             });

            //     });

            // });

        });

        const participant = new ParticipantsModel();

        describe('setParticipant', () => {

            // create new firebase user and participant to run tests against
            const mockAngularFireDb: { database: firebase.database.Database } = {
                database: admin.database() as any
            };

            participant.setDb(mockAngularFireDb.database);

            participant.model = participantInfo;

            describe('create participant', () => {

                it('should create new firebase participant/institution and resolve the new institutionId', async () => {

                    await participant.create()
                        .then((participantCreated: boolean) => {
                            participantInfo.info.institutionId = participant.model.info.institutionId;

                            // participantInfo.info.institutionId should have been set
                            expect(participantInfo.info.institutionId).to.equal(participant.model.info.institutionId);

                            assert.isTrue(participantCreated);
                        }, (err: any) => {
                            console.error(err, participant.model);
                        });

                });

            });

            describe('get existing participant', () => {

                it('should look-up participant by institutionId', async () => {

                    await participant.get(participantInfo.info.institutionId)
                        .then((resultData: any) => {
                            expect(resultData.institutionId).to.equal(participantInfo.info.institutionId);
                        }, (err: any) => {
                            console.error(err, participant.model);
                        });

                });

            });

        });

        // let triggers;

        // before(() => {
        //     // Require permissions.ts and save the exports inside a
        //     // namespace called permissions.
        //     // This includes our cloud functions, which can now be
        //     // accessed at permissions.myFunction
        //     triggers = require('./index');
        // });

        describe('updateParticipantPermissions', () => {

            // const params = {
            //     institutionId: participantInfo.info.institutionId,
            //     uid: userProfile.uid
            // };

            it('should create user permissions for user',
                async () => {

                    const roles = { manager: true };

                    const snapshot = {
                        val: () => (roles),
                        after: {
                            val: () => (roles)
                        }
                    };


                    // make sure the parameters being passed
                    // in to the trigger exists
                    assert.isNotNull(participantInfo.info.institutionId);
                    assert.isNotNull(userProfile.uid);

                    // testing background functions using firebase-functions-test wrapper
                    const wrapped = test.wrap(triggers.createParticipantPermissions);

                    await wrapped(snapshot, {
                        params: {
                            institutionId: participantInfo.info.institutionId,
                            uid: userProfile.uid
                        }
                    }).then((result: boolean) => {

                        // create permissions should be successful
                        assert.isTrue(result);
                    });
                });

            it('should update user permissions for user',
                async () => {

                    const roles = { viewer: true };

                    const snapshot = {
                        val: () => (roles),
                        after: {
                            val: () => (roles)
                        }
                    };

                    // testing background functions using firebase-functions-test wrapper
                    const wrapped = test.wrap(triggers.updateParticipantPermissions);

                    await wrapped(snapshot, {
                        params: {
                            institutionId: participantInfo.info.institutionId,
                            uid: userProfile.uid
                        }
                    }).then((result: boolean) => {

                        // update permissions should be successful
                        assert.isTrue(result);
                    });
                });

            it('should delete user permissions for user',
                async () => {

                    const roles = { viewer: true };

                    const snapshot = {
                        val: () => (roles),
                        after: {
                            val: () => (roles)
                        }
                    };

                    // testing background functions using firebase-functions-test wrapper
                    const wrapped = test.wrap(triggers.removeParticipantPermissions);

                    await wrapped(snapshot, {
                        params: {
                            institutionId: participantInfo.info.institutionId,
                            uid: userProfile.uid
                        }
                    }).then((result: boolean) => {

                        // removal of permissions should be successful
                        assert.isTrue(result);
                    });
                });
        });

        after(async () => {

            console.log('Cleaning up mock data');
            test.cleanup();

            const promiseArr = [];

            // delete mock user from firebase
            if (!_.isEmpty(userProfile.uid)) {
                promiseArr.push(admin.auth().deleteUser(userProfile.uid));
                promiseArr.push(admin.database().ref('users/' + userProfile.uid).remove());
            } else {
                console.error('uid is missing', { user: userProfile, participant: participantInfo });
            }

            // delete the mock participant from firebase
            if (!_.isEmpty(participantInfo.info.institutionId)) {
                promiseArr.push(admin.database().ref('participants/' + participantInfo.info.institutionId).remove());
            } else {
                console.error('missing uid or institutionId', { user: userProfile, participant: participantInfo });
            }

            // execute promises
            await Promise.all(promiseArr)
                .catch((err) => {
                    console.error('unable to clean-up mock data', err);
                });
        });

    });

    describe('should NOT be able to delete root nodes in firebase', () => {

        // ONLY if security roles are enabled
        // check that certain roots cannot be written to or deleted

        // TODO: write tests for /participants /users /super_permissions /participant_permissions

    });

});

describe('dependencies', () => {

    // init allRequiredPackagesExistInPackageJson
    let allRequiredPackagesExistInPackageJson = true;

    const compareDependencies = (depGroup1: { name: string, dep: {} }, depGroup2: { name: string, dep: {} }, errMsg: string) => {

        _.forEach(depGroup1.dep, (val: string, key: string) => {

            // check to see that each package is include in requiredPackages
            if (depGroup1.dep[key] !== depGroup2.dep[key]) {
                console.error(errMsg, {
                    [depGroup1.name]: {
                        _key: key,
                        _val: depGroup1.dep[key],
                    },
                    [depGroup2.name]: {
                        _key: key,
                        _val: depGroup2.dep[key],
                    },
                });
                allRequiredPackagesExistInPackageJson = false;
            }

        });

    };

    describe('functions dependencies', () => {

        // get production dependencies from package.json
        // console.info(functionPackages);
        // const packageJson = fs.readFileSync('./../package.json');
        const packageJson = require('./../package.json');
        const requiredPackages = {
            "aws-sdk": "^2.323.0",
            "body-parser": "^1.18.3",
            "cookie-parser": "^1.4.3",
            "cors": "^2.8.4",
            "express": "^4.16.3",
            "express-session": "^1.15.6",
            "firebase-admin": "^6.0.0",
            "firebase-functions": "^2.0.4",
            "lodash": "^4.17.10",
            "moment": "^2.22.2",
            "passport": "0.2.x",
            "passport-idaas-openidconnect": "1.0.0",
            "tsoa": "^2.1.8",
            "sib-api-v3-sdk": "^7.0.1"
        };

        // console.info(JSON.stringify(packageJson.dependencies));

        allRequiredPackagesExistInPackageJson = true;

        it('should not include unnecessary or missing dependencies', () => {

            // If error, check if there are any accidental new dependencies added
            // to bundle that are not necessary or any if versions changed.
            // If there is a new required dependency, add it to 'requiredPackages'

            // check for newly added dependencies
            compareDependencies(
                { name: 'packageJson', dep: packageJson.dependencies },
                { name: 'requiredPackages', dep: requiredPackages },
                'possible bloat error: '
            );

            // check for missing dependencies
            compareDependencies(
                { name: 'requiredPackages', dep: requiredPackages },
                { name: 'packageJson', dep: packageJson.dependencies },
                'missing required dependencies: '
            );

            expect(allRequiredPackagesExistInPackageJson).to.equal(true);
        });

    });

    describe('angular project dependencies', () => {

        // get production dependencies from package.json
        // console.info(functionPackages);
        // const packageJson = fs.readFileSync('./../package.json');
        const packageJson = require('./../../package.json');
        const requiredPackages = {
            "@angular/animations": "^6.1.0",
            "@angular/cdk": "^6.4.3",
            "@angular/common": "^6.1.0",
            "@angular/compiler": "^6.1.0",
            "@angular/core": "^6.1.0",
            "@angular/fire": "^5.2.1",
            "@angular/flex-layout": "^6.0.0-beta.17",
            "@angular/forms": "^6.1.0",
            "@angular/http": "^6.1.0",
            "@angular/material": "^6.4.5",
            "@angular/platform-browser": "^6.1.2",
            "@angular/platform-browser-dynamic": "^6.1.0",
            "@angular/router": "^6.1.0",
            "carbon-components": "^9.28.0",
            "core-js": "^2.5.4",
            "elasticsearch": "^15.1.1",
            "elasticsearch-browser": "^15.1.1",
            "firebase": "^5.3.1",
            "highlight.js": "^9.12.0",
            "lodash": "^4.17.10",
            "memoizee": "0.4.14",
            "moment": "^2.22.2",
            "ngx-markdown": "^6.0.1",
            "rxjs": "^6.2.2",
            "zone.js": "^0.8.26"
        };

        // console.info(JSON.stringify(packageJson.dependencies));

        allRequiredPackagesExistInPackageJson = true;

        it('should not include unnecessary or missing dependencies', () => {

            // If error, check if there are any accidental new dependencies added
            // to bundle that are not necessary or any if versions changed.
            // If there is a new required dependency, add it to 'requiredPackages'

            // check for newly added dependencies
            compareDependencies(
                { name: 'packageJson', dep: packageJson.dependencies },
                { name: 'requiredPackages', dep: requiredPackages },
                'possible bloat error: '
            );

            // check for missing dependencies
            compareDependencies(
                { name: 'requiredPackages', dep: requiredPackages },
                { name: 'packageJson', dep: packageJson.dependencies },
                'missing required dependencies: '
            );

            expect(allRequiredPackagesExistInPackageJson).to.equal(true);
        });

    });





});
