// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate } from '@angular/router';
import { WindowService } from '../services/window.service';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';

/**
 * Redirects to an external link that is set in the router "data" definition
 * (and DOES allow the user to come back to the page they were previously viewing)
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
