// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { ROLES } from '../constants/general.constants';
import { Injectable, NgZone } from '@angular/core';
import { AngularFireDatabase } from '@angular/fire/database';
import { HttpClient, HttpHeaders, HttpRequest } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { AuthService } from './auth.service';
import { IRolesOptions } from '../models/user.interface';
import { Confirm2faService } from './confirm2fa.service';
import { keys, startCase, isEmpty } from 'lodash';
import { Observable, Observer } from 'rxjs';
import { IParticipantUsers } from '../models/participant.interface';

@Injectable()
export class ParticipantPermissionsService {

    roles = ROLES;
    environment = environment;

    constructor(
        private db: AngularFireDatabase,
        public http: HttpClient,
        private authService: AuthService,
        private ngZone: NgZone,
        private confirm2Fa: Confirm2faService
    ) {
        // Do NOT create a model here rather pass the model into the angular
        // service method call. This makes it possible to use the methods as
        // service and models as objects
    }

    /**
     * disables button in view so that users cannot
     * edit or change their own permissions
     *
     * @param {string} email1
     * @param {string} email2
     * @returns
     * @memberof UsersComponent
     */
    disable(email1: string, email2: string) {
        if (email1.toLowerCase() === email2.toLowerCase()) {
            return true;
        } else {
            return false;
        }
    }

    /**
     * formats permissions array to human readable string in view
     *
     * @param {IRolesOptions} roles
     * @returns
     * @memberof UsersComponent
     */
    humanizeRoles(roles: IRolesOptions) {
        const str = keys(roles).toString();
        return startCase(str.replace(',', ', '));
    }

    /**
     * get all users for a participant
     *
     * @param {string} institutionId
     * @returns {Promise<IUserParticipantPermissions>}
     * @memberof ParticipantPermissionsService
     */
    getAllUsers(institutionId: string): Promise<IParticipantUsers> {

        const ngZone = this.ngZone;

        return new Promise((resolve, reject) => {

            this.db.database.ref('/participants/' + institutionId + '/users')
                .once('value', (data: firebase.database.DataSnapshot) => {
                    const users = data.val() ? data.val() : null;

                    // return users object
                    ngZone.run(() => {
                        resolve(users);
                    });
                });

        });

    }

    /**
     * get all users for a participant
     *
     * @param {string} institutionId
     * @returns {Observable<IUserParticipantPermissions>}
     * @memberof ParticipantPermissionsService
     */
    getAllUsersObservable(institutionId: string): Observable<IParticipantUsers> {

        return new Observable((observer: Observer<IParticipantUsers>) => {

            this.db.database.ref('/participants/' + institutionId + '/users')
                .on('value', (data: firebase.database.DataSnapshot) => {

                    const users: IParticipantUsers = data.val() ? data.val() : null;

                    this.ngZone.run(() => {
                        observer.next(
                            users
                        );
                    });
                });
        });

    }

    /**
     * Creates new (or updates) user permissions and returns uid
     * (or returns existing user uid for email)
     *
     * @param {string} institutionId
     * @param {('admin' | 'manager' | 'viewer')} role
     * @param {string} email
     * @returns
     * @memberof ParticipantPermissionsService
     */
    update(institutionId: string, role: 'admin' | 'manager' | 'viewer', email: string): Promise<any> {
        return new Promise(async (resolve, reject) => {

            const self = this;

            try {
                // update user's permissions
                const h: HttpHeaders = await this.authService.getFirebaseIdToken(institutionId);

                // validate required fields are present
                if (!isEmpty(email) && !isEmpty(institutionId) && !isEmpty(role)) {

                    // update permissions node in firebase
                    // NOTE: doing as an http post instead of calling firebase db
                    // directly because this requires the success of two separate calls to
                    // firebase (in the case of adding a user -> one call to create the user,
                    // and one call to add their permissions). Creating a http post request
                    // ensures graceful failure
                    const r = new HttpRequest(
                        'POST',
                        self.environment.apiRootUrl + '/permissions/participant',
                        {
                            email: email,
                            institutionId: institutionId,
                            role: role
                        },
                        { headers: h }
                    );

                    await self.confirm2Fa.go(r)
                        .then((uid: string) => {
                            resolve(uid);
                        }, (error) => {
                            console.log('Error: Unable to add participant permissions.', error);
                            reject();
                        });

                }
            } catch (error) {
                console.log('Error: Unable to add participant permissions.', error);
                reject();
            }

        });

    }

    /**
     * Remove user permissions.
     * NOTE: permissions are added and removed one at at time.
     *
     * @returns {Promise<void>}
     * @memberof ParticipantPermissionsModel
     */
    remove(userId: string, institutionId: string): Promise<void> {

        // validate required fields are present
        if (!isEmpty(userId) && !isEmpty(institutionId)) {

            // TODO: create a DELETE endpoint  in auth-service to handle this instead, similar to update()
            // return promise since result is a single success or failure
            return new Promise(async (resolve, reject) => {

                // update record in firebase
                await this.db.database.ref(
                    'participant_permissions/' +
                    userId +
                    '/' +
                    institutionId
                ).remove()
                    .then(() => {
                        resolve();
                    }, (error) => {
                        console.log('Error: Unable to remove participant permissions.', error);
                        reject();
                    });

            });

        }

    }

}
