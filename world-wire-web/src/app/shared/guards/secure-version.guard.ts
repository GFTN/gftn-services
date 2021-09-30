// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import * as _ from 'lodash';
import { Injectable } from '@angular/core';
import { CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot, Router } from '@angular/router';
import { Observable } from 'rxjs';
import { UtilsService } from '../utils/utils';

@Injectable()
export class SecureVersionGuard implements CanActivate {

    constructor(
        private router: Router,
        private utils: UtilsService
    ) { }

    canActivate(
        route: ActivatedRouteSnapshot,
        state: RouterStateSnapshot
    ): Observable<boolean> | Promise<boolean> | boolean {

        const accepted = this.utils.getCookie('wwdocspermission');

        if (accepted === 'true') {
            // if cookie with password does exist then proceed to docs
            return true;
        } else {
            // redirect user to input password generate cookie
            this.router.navigate(['/docs/secure']);
            // if cookie does not exist redirect user to secure component
            return false;
        }

    }
}
