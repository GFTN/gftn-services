// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import * as _ from 'lodash';
import { ROLES } from '../constants/general.constants';
import { Injectable, NgZone } from '@angular/core';
import { AngularFireDatabase } from '@angular/fire/database';
import { HttpClient, HttpHeaders, HttpRequest } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { AuthService } from './auth.service';
import { IUserParticipantPermissions, IUserSuperPermissions } from '../models/user.interface';
import { Confirm2faService } from './confirm2fa.service';
import { ParticipantPermissionsService } from './participant-permissions.service';
import { Observable, Observer } from 'rxjs';

@Injectable()
export class SuperPermissionsService {

    roles = ROLES;
    environment = environment;

    constructor(
        private db: AngularFireDatabase,
        public http: HttpClient,
        private authService: AuthService,
        private ngZone: NgZone,
        private confirm2Fa: Confirm2faService,
        private participantPermissionsService: ParticipantPermissionsService
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
    disable = this.participantPermissionsService.disable;

    /**
     * formats permissions array to human readable string in view
     *
     * @param {IRolesOptions} roles
     * @returns
     * @memberof UsersComponent
     */
    humanizeRoles = this.participantPermissionsService.humanizeRoles;

    /**
     * Get all super users
     *
     * @returns {Promise<IUserParticipantPermissions>}
     * @memberof SuperPermissionsService
     */
    getAllUsers(): Promise<IUserParticipantPermissions> {

        const ngZone = this.ngZone;

        return new Promise((resolve, reject) => {

            this.db.database.ref('/super_permissions/')
                .once('value', (data: firebase.database.DataSnapshot) => {
                    const users = data.val();
                    if (users) {
                        // return users object
                        ngZone.run(() => {
                            resolve(users);
                        });
                    }
                    // return null if no participants
                    ngZone.run(() => {
                        resolve(null);
                    });
                });

        });

    }

    /**
     * Get all super users
     *
     * @returns {Observable<IUserSuperPermissions>}
     * @memberof SuperPermissionsService
     */
    getAllUsersObservable(): Observable<IUserSuperPermissions> {
        return new Observable((observer: Observer<IUserSuperPermissions>) => {
            this.db.database.ref('/super_permissions/')
                .on('value', (data: firebase.database.DataSnapshot) => {
                    const users = data.val() ? data.val() : null;

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
     * Alias: SuperPermissionsService.setUserId()
     *
     * @param {string} institutionId
     * @param {('admin' | 'manager' | 'viewer')} role
     * @param {string} email
     * @returns
     * @memberof ParticipantPermissionsService
     */
    update(role: 'admin' | 'manager' | 'viewer', email: string) {
        return new Promise((resolve, reject) => {

            const self = this;

            this.authService.getFirebaseIdToken().then((h: HttpHeaders) => {

                // update user's permissions

                // validate required fields are present
                if (!_.isEmpty(email) && !_.isEmpty(role)) {

                    // update permissions node in firebase
                    // NOTE: doing as an http post instead of calling firebase db
                    // directly because this requires the success of two separate calls to
                    // firebase (in the case of adding a user -> one call to create the user,
                    // and one call to add their permissions). Creating a http post request
                    // ensures graceful failure
                    const r = new HttpRequest(
                        'POST',
                        self.environment.apiRootUrl + '/permissions/super',
                        {
                            email: email,
                            role: role
                        },
                        { headers: h }
                    );

                    self.confirm2Fa.go(r)
                        .then((uid: string) => {
                            resolve(uid);
                        }, (error) => {
                            console.log('Error: Unable to add participant permissions.', error);
                            reject();
                        });

                }

            });

        });

    }


    /**
     * Remove user permissions.
     * NOTE: permissions are added and removed one at at time.
     *
     * @param {string} userId
     * @returns {Promise<void>}
     * @memberof SuperPermissionsService
     */
    remove(userId: string): Promise<void> {

        // validate required fields are present
        if (!_.isEmpty(userId)) {

            // TODO: create a DELETE endpoint in auth-service to handle this instead, similar to update()
            // return promise since result is a single success or failure
            return new Promise((resolve, reject) => {

                // update record in firebase
                this.db.database.ref(
                    'super_permissions/' +
                    userId
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
