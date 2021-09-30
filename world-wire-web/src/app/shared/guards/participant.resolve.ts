// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 

import { Observable } from '@firebase/util';
import { AngularFireDatabase } from '@angular/fire/database';
import { IInstitution } from '../../shared/models/participant.interface';
import { ActivatedRouteSnapshot, Resolve } from '@angular/router';
import { Injectable } from '@angular/core';
import { SessionService } from '../services/session.service';
import * as _ from 'lodash';
import { AuthService } from '../services/auth.service';
import { IUserProfile } from '../models/user.interface';

@Injectable()
export class ParticipantResolve implements Resolve<string> {

    constructor(
        private authService: AuthService,
        private db: AngularFireDatabase,
        private sessionService: SessionService
    ) { }

    resolve(
        route: ActivatedRouteSnapshot,
        // state: RouterStateSnapshot
    ): Observable<any> | Promise<any> | any {

        // resolve with participant details

        return new Promise(async (resolve) => {

            const getPreviousSlug: string =
                _.has(this.sessionService, 'institution.info.slug') ? this.sessionService.institution.info.slug : '';

            // check if current route matches the previously saved participant in session.service.ts
            if (getPreviousSlug === route.params.slug) {

                // no need to query institution if it is already stored in the service
                resolve(this.sessionService.institution);

            } else {

                // const slug = route.params.slug ? route.params.slug : null;
                const slug = route.paramMap.get('slug');

                const _iid = await this.db.database.ref('slugs/' + slug).once('value');

                let iid = _iid.val();


                await this.getUserPermissions().then(async (resolvedUser: IUserProfile) => {

                    if (!slug) {
                        // get user permissions
                        let participantPermissions = this.authService.userProfile ?
                            this.authService.userProfile.participant_permissions : null;

                        if (!participantPermissions) {
                            // check if permissions exist on user
                            if (resolvedUser.participant_permissions) {
                                // set user permissions
                                participantPermissions = resolvedUser.participant_permissions;
                            } else {
                                return resolve(false);
                            }
                        }

                        // get list of institutions for which
                        // the user has valid permissions for
                        iid = Object.keys(participantPermissions)[0];
                    }

                    // set current participant
                    // this.db.database.ref('participants/' + route.params.institutionId)
                    await this.db.database.ref('participants/' + iid).once('value')
                        .then((snapshot: firebase.database.DataSnapshot) => {

                            // const participant: IInstitution = snapshot.val();
                            const participant: IInstitution = snapshot.val();

                            // store participant so that we don't have to wait on each refresh
                            // this.sessionService.institution = participant;
                            this.sessionService.institution = participant;

                            // console.log('resolved: ', this.sessionService.institution);
                            return resolve(participant);

                        });
                });
            }

        });

    }

    /**
     * getUserPermissions
     *
     * dependency for determining which participant
     * to resolve to if no slug exists
     */
    private getUserPermissions(): Promise<IUserProfile> {
        return new Promise((resolve) => {

            const getPreviousSlug: IUserProfile = _.has(this.authService, 'user') ? this.authService.userProfile : null;

            // check if current route matches the previously saved participant in session.service.ts
            if (getPreviousSlug !== null) {

                // no need to query institution if it is already stored in the service
                resolve(this.authService.userProfile);

            } else {

                const uid = this.authService.auth.auth.currentUser.uid;

                // set current participant
                // get user permissions for the view
                this.db.database.ref(`/users/${uid}`)
                    .once('value', (userData: firebase.database.DataSnapshot) => {

                        const user: IUserProfile = userData.val();
                        // set user info in session
                        this.authService.userProfile = user;

                        // console.log('got user permissions: ', this.authService.userPermissions);
                        resolve(user);
                    });

            }

        });
    }
}
