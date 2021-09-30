// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// PURPOSE: the purpose of this controller is for managing users (aka: User Management)
// Endpoints exist to provision and revoke permissions for both
//     "participant users"  - users associated with a participant and logged into the portal as without an @ibm.com email
//     "super users"        - users with an @ibm.com email that have been provisioned super access rights

import {
    Route, Controller, Post, Body, Request, Header, OperationId
    // , Security
} from 'tsoa';
import * as admin from 'firebase-admin';
import { isEmpty, includes } from 'lodash';
import { IFirebaseUserRequest } from '../models/auth.model';
import { IParticipantRoles, IUserProfile } from '../shared/models/user.interface';
import { AuthHelpers } from '../auth/auth-helpers';

@Route('permissions')
export class PermissionsController extends Controller {

    authHelpers = new AuthHelpers();

    /**
    * Updates permissions for a user associated with a participant
    *
    * @private
    * @param {string} institutionId
    * @param {string} firebaseUserId
    * @param {('admin' | 'manager' | 'viewer')} role
    * @returns
    * @memberof IBMIdController
    */
    @Post('participant')
    @OperationId('permissionsParticipant')
    // @Security('api_key')
    async updateParticipantUserPermissions(
        @Request() req: IFirebaseUserRequest,
        @Body() body: { institutionId: string, email: string, role: 'admin' | 'manager' | 'viewer' },
        @Header('x-fid') fid: string,
        @Header('x-iid') iid: string
    ): Promise<string> {

        return await new Promise<string>((resolve, reject) => {

            // requester should not be able to delete themselves
            if (req.email === body.email) {
                this.setStatus(409);
                reject('user may not delete themselves');
            } else {

                // validate required fields are present
                if (!isEmpty(body.email) && !isEmpty(body.institutionId) && !isEmpty(body.role)) {

                    // get new user userId for email
                    this.authHelpers.setFirebaseUser(body.email)
                        .then(async (firebaseUserId: string) => {

                            const data: IParticipantRoles = {
                                email: body.email,
                                roles: { [body.role]: true }
                            };

                            await this.checkInitUser(firebaseUserId, body.email);

                            // set permissions
                            // update permissions node in firebase
                            return await admin.database().ref(
                                'participant_permissions/' +
                                firebaseUserId +
                                '/' +
                                body.institutionId
                            ).update(data).then(() => {
                                resolve(firebaseUserId);
                            }, (error) => {
                                console.log('Error: Unable to add participant permissions.', error);
                                reject();
                            });

                        }, (err: any) => {
                            console.log('Failed to set user: ', err);
                            reject();
                        });

                } else {
                    this.setStatus(400);
                    reject('unknown email, role, and/or institution');
                }

            }

        });
    }

    /**
    * Updates permissions for a super user (IBM Super Administrator)
    *
    * @private
    * @param {string} firebaseUserId
    * @param {('admin' | 'manager' | 'viewer')} role
    * @returns
    * @memberof IBMIdController
    */
    @Post('super')
    @OperationId('permissionsSuper')
    // @Security('api_key')
    async updateSuperPermissions(
        @Request() req: IFirebaseUserRequest,
        @Body() body: { email: string, role: 'admin' | 'manager' | 'viewer' },
        @Header('x-fid') fid: string
    ) {

        try {

            // requester should not be able to delete themselves
            if (req.email === body.email) {
                this.setStatus(409);
                return 'user may not delete themselves';
            } else {
                // check that email ending ends in an allowable super admin address
                const ending: string = body.email.split('@')[1];
                const allowableEndings: string[] = ['ibm.com', 'us.ibm.com', 'sg.ibm.com', 'in.ibm.com'];
                if (!includes(allowableEndings, ending)) {
                    this.setStatus(403);
                    return 'email ending must be ' + allowableEndings.toString();
                } else {

                    // validate required fields are present
                    if (!isEmpty(body.email) && !isEmpty(body.role)) {

                        // get new user userId for email
                        const firebaseUserId = await this.authHelpers.setFirebaseUser(body.email)
                            .catch((err: any) => {
                                console.error(err);
                                // this.setStatus(500);
                                // return 'unknown user';
                                throw new Error('unknown user');
                            });

                        await this.checkInitUser(firebaseUserId, body.email);

                        const data: IParticipantRoles = {
                            email: body.email,
                            roles: { [body.role]: true }
                        };

                        // set permissions
                        // update permissions node in firebase
                        return await admin.database().ref(
                            'super_permissions/' +
                            firebaseUserId
                        ).update(data).then(() => {
                            return firebaseUserId;
                        }, (err: any) => {
                            console.error(err);
                            // this.setStatus(500);
                            // return 'unknown user';
                            throw new Error('unable to add user by email');
                        });

                    } else {
                        this.setStatus(400);
                        return 'unknown email and/or role';
                    }

                }

            }

        } catch (error) {
            this.setStatus(500);
            return error.message;
        }

    }

    async checkInitUser(firebaseUserId: string, email: string): Promise<void> {
        // NOTE: user would not exist on initial create
        // and as such the permissions trigger will error out
        const userRef = admin.database().ref('users');
        const _user = await userRef.child(firebaseUserId).once('value');
        const user = _user.val();


        // check if user profile exists ie: ref.('/user'), if not create user profile
        // must have email and uid to write
        if (isEmpty(user) && !isEmpty(firebaseUserId)) {

            // create new user profile (even thought they may have never logged in before)
            const newUser: IUserProfile = {
                profile: {
                    email: email
                }
                // // NOTE: super and participant permissions will
                // // be set by firebase trigger for security purposes
                // participant_permissions: ...,
                // super_permissions: ...
            };

            // write new user to firebase
            await userRef.child(firebaseUserId).update(newUser);

            return;

        } else {

            return;

        }
    }

}
