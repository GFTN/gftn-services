// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable } from '@angular/core';
import { CanActivate, ActivatedRouteSnapshot, Router } from '@angular/router';
import { isEmpty } from 'lodash';
import { AuthService } from '../services/auth.service';


/***
 * SuperPermissionsGuard
 * is a route guard that protects routes that require
 * super permissions to view
 */
@Injectable()
export class SuperPermissionsGuard implements CanActivate {

    constructor(
        private authService: AuthService,
        private router: Router
    ) { }

    async canActivate(
        route: ActivatedRouteSnapshot
    ) {

        try {

            // set userProfile on authService
            const userProfile = await this.authService.getUserProfile(this.authService);

            // check if user profile exists otherwise redirect to login
            if (isEmpty(userProfile)) {

                // unable to determine if the user has permissions for this institution
                // redirect user to login
                this.router.navigate(['/login']);
                return false;

            } else {

                if (!route.data.super_permissions) {
                    this.router.navigate(['/not-found']);
                    return false;
                }

                if (
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

        } catch (error) {
            // most likely one of the promises related
            // to getting the institution or user permissions failed
            console.log(error);
            this.router.navigate(['/not-found']);
            return false;
        }

    }

}
