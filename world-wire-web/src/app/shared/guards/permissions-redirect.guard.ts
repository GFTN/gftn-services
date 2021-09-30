// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable } from '@angular/core';
import { CanActivate, ActivatedRouteSnapshot, Router, RouterStateSnapshot } from '@angular/router';
import { SessionService } from '../services/session.service';
import { get, isEmpty } from 'lodash';
import { AuthService } from '../services/auth.service';
import { AngularFireDatabase } from '@angular/fire/database';
import { IInstitution } from '../models/participant.interface';
import { switchMap } from 'rxjs/operators';

/***
 * ParticipantPermissionsGuard
 * protects routes that require participant permissions
 * or super permissions to view
 */
@Injectable()
export class ParticipantPermissionsGuard implements CanActivate {

    // store slug from route params to use across functions
    slug: string;

    constructor(
        private authService: AuthService,
        private sessionService: SessionService,
        private db: AngularFireDatabase,
        private router: Router
    ) { }

    async canActivate(
        route: ActivatedRouteSnapshot,
        state: RouterStateSnapshot
    ) {

        try {

            // get slug of parent route to activate child route
            const url: string[] = state.url.split('/');

            this.slug = url.length > 2 ? url[2] : '';

            // set userProfile on authService
            const userProfile = await this.authService.getUserProfile(this.authService);

            // set institution on session
            const selectedInstitution = await this.getInstitution();

            // logical redirect based on if:
            // 1) user has permissions
            // 2) institution is available in the session
            // check if user profile exists otherwise redirect to login
            if (isEmpty(userProfile)) {

                // unable to determine if the user has permissions for this institution
                // redirect user to login
                this.router.navigate(['/login']);
                return false;

            } else {

                // check if institution exists otherwise redirect to not-found
                if (isEmpty(selectedInstitution)) {

                    // unable to set the institution
                    // redirect to unauthorized page
                    this.router.navigate(['/not-found']);
                    return false;

                } else {

                    // check if user has access rights to proceed to route based
                    // on specified permissions on portals 'route.data' in portal-routing.ts
                    if (
                        // check if user has "PARTICIPANT" specific permissions to access route
                        this.authService.hasParticipantPermissions(
                            userProfile,
                            selectedInstitution.info.institutionId,
                            route.data.participant_permissions) === true ||
                        // or
                        // check if user has "SUPER" permissions to access route
                        this.authService.hasSpecificSuperPermissions(userProfile, route.data.super_permissions) === true
                    ) {
                        // success resolve true and allow transition
                        // activate route since permissions exist
                        return true;
                    } else {
                        // unable to determine if the user has permissions for this institution
                        // redirect to unauthorized page
                        this.router.navigate(['/unauthorized']);
                        return false;
                    }

                }

            }

        } catch (error) {
            // most likely one of the promises related
            // to getting the institution or user permissions failed
            console.log(error);
            this.router.navigate(['/not-found']);
            return false;
        }

    }

    /**
     * checks if the session already includes the institution
     * associated with the selected routes current slug,
     * otherwise it retrieves the institution information from
     * the db and sets it to the session based
     * on the specified route's slug
     *
     * @private
     * @returns {Promise<IInstitution>}
     * @memberof ParticipantPermissionsGuard
     */
    private getInstitution(): Promise<IInstitution> {

        return new Promise(async (resolve) => {

            // prevent looking up a new institution from the
            // db if it already exists in the session
            if (get(this.sessionService, 'institution.info.slug') === this.slug) {
                resolve(this.sessionService.institution);
            }

            // get the institutionId for the current route's slug
            const _iid = await this.db.database.ref('slugs/' + this.slug).once('value');
            const iid = _iid.val();

            // Query to set institution in SessionService.
            // Moved from ParticipantResolve to Guard to check
            // route permissions before activating route.
            this.db.database.ref('participants/' + iid).once('value')
                .then((snapshot: firebase.database.DataSnapshot) => {

                    const institution: IInstitution = snapshot.val();

                    // store the institution in current session
                    this.sessionService.institution = institution;

                    // return the newly looked up institution
                    resolve(institution);

                }).catch((err) => {
                    // something unexpectedly went wrong, so return empty value...
                    resolve(null);
                });

        });
    }

}
