// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, RouterStateSnapshot, Resolve, CanActivate } from '@angular/router';
import { WindowService } from '../services/window.service';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';

/**
 * Redirects to an external link that is set in the router "data" definition
 * (but does NOT allow the user to come back to the page they were previously viewing)
 *
 * @export
 * @class ExternalRedirectResolver
 * @implements {Resolve<any>}
 */
@Injectable()
export class ExternalRedirectResolver implements Resolve<any> {
    constructor(
        private windowService: WindowService,
    ) {
        // console.log('external redirect 1');
    }

    /**
     * Catch navigate event and redirect to external link (set per data field in angular router)
     *
     * @param {ActivatedRouteSnapshot} route
     * @param {RouterStateSnapshot} state
     * @returns {(Observable<any> | Promise<any> | any)}
     * @memberof ExternalRedirectResolver
     */
    resolve(
        route: ActivatedRouteSnapshot,
        state: RouterStateSnapshot
    ) {

        // console.log('external redirect 2');

        return this.windowService.windowRef.location.href = (route.data as any).externalUrl;
    }

}

// @Injectable()
// export class CanDeactivateGuard implements CanDeactivate<any> {

//     canDeactivate(
//         component: any,
//         route: ActivatedRouteSnapshot,
//         state: RouterStateSnapshot
//     ): Observable<boolean> | boolean {

//         console.log('leaving terms and conditions component');

//         return false;

//     }
// }

/**
 * Redirects to an external link that is set in the router "data" definition
 * (but DOES allow the user to come back to the page they were previously viewing)
 *
 * @export
 * @class CanActivateExternalUrlGuard
 * @implements {CanActivate}
 */
@Injectable()
export class ExternalRedirectGuard implements CanActivate {

    constructor(
        private windowService: WindowService
    ) { }

    canActivate(
        route: ActivatedRouteSnapshot
        // state: RouterStateSnapshot
    ): Observable<boolean> | Promise<boolean> | boolean {

        if (environment.production) {

            // only do redirect in production

            console.log('redirecting to external page');

            setTimeout(() => {
                this.windowService.windowRef.location.href = (route.data as any).externalUrl;
            });

            return false;

        }

    }

}
