// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable, isDevMode } from '@angular/core';
import { CanActivate, Router, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { map, take, tap } from 'rxjs/operators';
import { AngularFireAuth } from '@angular/fire/auth';
import { Observable } from '@firebase/util';
import { DocumentService } from '../services/document.service';
import { AuthService } from '../services/auth.service';

export class AuthRedirectCookie {

    cookieName = 'wwAuthRedirect';

    check(documentService: DocumentService) {

        const redirectRoute = this.get(this.cookieName, documentService);

        if (redirectRoute) {
            //  delete cookie
            this.delete();

            // if redirectRoute exists redirect to that route
            return redirectRoute;

        } else {
            // no redirect route set
            return false;
        }

    }

    /**
     * Delete the cookie
     *
     * @memberof AuthRedirectCookie
     */
    delete() {
        // delete cookie by setting the expire date in the past
        // document.cookie = this.cookieName + '=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;';
        // You can delete a cookie by simply updating its expiration time to zero.
        document.cookie = this.cookieName + '=; Path=/; Expires=0;';
    }

    /**
     * Get cookie by its name
     *
     * @param {string} name
     * @param {DocumentService} documentService
     * @returns
     * @memberof AuthRedirectCookie
     */
    get(name: string, documentService: DocumentService) {
        const value = '; ' + documentService.documentRef.cookie;
        const parts = value.split('; ' + name + '=');
        if (parts.length === 2) {
            return parts.pop().split(';').shift();
        }
    }

    /**
     * Creates a cookie to redirect the user upon successful authentication
     *
     * @param {DocumentService} documentService
     * @param {string} route
     * @memberof AuthRedirectCookie
     */
    set(documentService: DocumentService, route: string) {

        documentService.documentRef.cookie = this.cookieName + '=' + route + '; expires=' +
            new Date(new Date().getFullYear(), new Date().getMonth() + 3, new Date().getDate()).toUTCString() +
            ';path=/';

    }

}

@Injectable()
export class AuthCanActivateGuard implements CanActivate {

    constructor(
        private auth: AngularFireAuth,
        private router: Router,
        private documentService: DocumentService,
        private authRedirectCookie: AuthRedirectCookie,
        private authService: AuthService
    ) { }

    canActivate(
        next: ActivatedRouteSnapshot,
        state: RouterStateSnapshot
    ): Observable<any> | Promise<any> | any {

        if (isDevMode) {
            console.log(next);
            console.log(state);
        }

        // if NOT page refresh - get currently cached/loaded user
        if (this.auth.auth.currentUser) {

            const redirectRoute = this.authRedirectCookie.check(this.documentService);

            if (redirectRoute) {
                // check if auth redirect cookie exists and go to that location
                this.router.navigate([redirectRoute]);
            } else {
                // continue to selected route if no redirect is present
                return true;
            }

        }

        // if is page refresh - wait for cached/loaded user
        return this.auth.authState.pipe(
            take(1),
            map(user => {

                // if (user) {
                //     // The signed-in user info
                //     console.log('Signed in (from route-guard): ', user.email);
                // }

                // return bool related to user existence
                // true = user exists, false = user does not exist
                return !!user;
            }),
            tap(
                // rxjs tap operator === rxjs do operator
                // loggedIn result of the rxjs map operator above
                async (loggedIn) => {

                    // await this.authService.getUserProfile(this.authService);

                    // create a redirect cookie to redirect user to the url
                    // they were originally attempting to visit after login
                    // is successful

                    // checks if the user is signed-in
                    if (!loggedIn) {

                        // create a redirect token to use upon successful login
                        if (state.url) {
                            this.authRedirectCookie.set(this.documentService, state.url);
                        }

                        console.log('access denied');

                        // if not signed-in - redirect to login screen
                        this.router.navigate(['/login']);

                    }

                }
            )
        );

    }

}
