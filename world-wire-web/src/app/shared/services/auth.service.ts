// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable, NgZone, isDevMode } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { AngularFireAuth } from '@angular/fire/auth';
import { UserInfo } from 'firebase';
import { DocumentService } from './document.service';
import { environment } from '../../../environments/environment';
import { AngularFireDatabase } from '@angular/fire/database';
import { IUserProfile } from '../models/user.interface';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { has, get, isEmpty } from 'lodash';
import { auth as fbAuth } from 'firebase';
import { SessionService } from '../../shared/services/session.service';

// import { HttpClient } from '@angular/common/http';
// TODO: 'firebase/auth' can't find auth module,
// need to look at firebase docs to see if there
// is a way to just import firebase auth module


import { OnDestroy } from '@angular/core';
import { Subject, timer, Subscription, Observable } from 'rxjs';
import { takeUntil, take } from 'rxjs/operators';

@Injectable()
export class AuthService implements OnDestroy {

    // IMPORTANT: instead of 'user'
    // user firebase.auth.currentUser
    // user: UserInfo;

    // IMPORTANT: firebase.auth.currentUser should be in scope
    // where needed since the "auth.guard.ts" waits to load view
    // until user is loaded so "watchUser()" is no longer needed

    // IMPORTANT: As a best practice use angularFire as much as possible
    // when checking permissions as abstracting new functions in this
    // service could result in inconsistency between the
    // auth.service.ts and the angularFire.auth service

    // user$: Observable<User | null>;
    // isInitialized: boolean;
    loginError: string;
    registrationError: string;
    email: string;
    submitted: boolean;

    // firebase profile as an observable that can be unsubscribed when
    // a user logs out
    private profileRef: firebase.database.Reference;

    /**
     * Set current user of the session.
     * Use to gather permissions, etc. other
     * necessary information about the user
     * for viewing certain routes
     * NOTE: see watchUser for getting user profile
     *
     * @type {IUserProfile}
     */
    userProfile: IUserProfile;

    // stores user language (short/long). defaults to browser default
    // TODO: can also store user's language preference.
    userLang = window.navigator.language;
    userLangShort = window.navigator.language.slice(0, 2);

    constructor(
        private session: SessionService,
        // auth is public so that components can reference
        // fields like "this.auth.auth.currentUser.uid"
        // in the html view
        public auth: AngularFireAuth,
        private router: Router,
        private activatedRoute: ActivatedRoute,
        private http: HttpClient,
        private ngZone: NgZone,
        private document: DocumentService,
        private db: AngularFireDatabase
    ) {

        // initialize user session watcher
        // to execute a custom action on sign-in and sign-out
        this.watchUser();

    }

    /**
     * Redirect from portal to node.js api to login to IBMId
     * upon successful sign in the user will be redirected to
     * /login_token with a firebase custom token which can be
     * used to login from the portal by calling
     * firebaseCustomTokenSignIn()
     *
     * @memberof AuthService
     */
    signInIbmId() {

        // go to IBMId sign-in page
        this.document.documentRef.location.href = environment.apiRootUrl + '/sso/token';
    }

    /**
     * Login custom firebase user and create session
     * and set local/long-lived auth token.
     * This authentication mechanism allows a firebase
     * user to be associated with an IBMId user and
     * therefore is afforded all the native firebase
     * features in a secure environment using firebase
     * security rules. This token is created in custom
     * node.js server and is redirected with the short lived
     * token to /login_token which calls this function to login
     * the user and create a session.
     *
     * @memberof AuthService
     */
    async firebaseCustomTokenSignIn(): Promise<IUserProfile> {

        try {

            if (!this.session.firebaseTempAuthToken) {
                console.log('missing temp auth token when trying to login');
                this.router.navigate(['/unauthorized']);
            }

            // auth.setPersistence = Existing and future Auth states are persisted in the current
            // session only. Closing the window would clear any existing state even
            // if a user forgets to sign out.
            // ...
            // New sign-in will be persisted with session persistence.
            await this.auth.auth.setPersistence(fbAuth.Auth.Persistence.SESSION);

            const userCred: firebase.auth.UserCredential = await this.auth.auth.signInWithCustomToken(this.session.firebaseTempAuthToken);

            // remove temp auth token from the browsers memory
            this.session.firebaseTempAuthToken = '';

            // Sign in success
            console.log('Successful login:', userCred);

            const userProfile = await this.getUserProfile(this);

            return userProfile;

        } catch (error) {

            // Sign in failure

            console.log('Failed login', error);

            // Handle Errors here.
            // var errorCode = error.code;
            // var errorMessage = error.message;
            // this.router.navigate(['/']);

            return null;

        }

    }

    /**
     * Sign-out user to landing page
     *
     * @memberof AuthService
     */
    signOut(redirect?: '/home' | '/inactive') {

        // NOTE: see watch user to perform other
        // actions on signOut

        const self = this;

        let _redirect = redirect;

        // if no redirect is set default to home
        if (!redirect) {
            _redirect = '/home';
        }

        // sign-out firebase session
        this.auth.auth.signOut().then(() => {

            // redirect to homepage after user logs out
            self.ngZone.run(() => {
                self.router.navigate([_redirect]);
            });

        });

        // sign-out IBMId session
        this.http.get(
            environment.apiRootUrl + '/sso/logout',
            { withCredentials: true });

    }

    /**
     * Generates an id to send to the back end to authenticate a request
     * against firebase when the identity of the user is needed to take
     * an action on the server side.
     *
     * @param {string} [institutionId]
     * @returns {Promise<string>}
     * @memberof AuthService
     */
    async getFirebaseIdToken(institutionId?: string, participantId?: string): Promise<HttpHeaders> {

        const fid = await this.auth.auth.currentUser.getIdToken(/* forceRefresh */ true);

        // firebase auth id used by back-end to
        // determine the firebase user making the http request
        let h = new HttpHeaders().set('X-fid', fid);
        h.append('x-fid', fid);

        // institutionId used by participant permissions lookup in back-end
        if (institutionId) {
            h = h.set('x-iid', institutionId);
        }

        if (participantId) {
            h = h.set('x-pid', participantId);
        }

        // usage example: this.http.post('api/items/add', body, { headers: h })

        return h;
    }

    /**
     * Adds necessary headers for maker/checker requests
     *
     * @param {HttpHeaders} h
     * @param {('request' | 'approve')} permission
     * @param {string} [approvalId]
     * @returns {HttpHeaders}
     * @memberof AuthService
     */
    addMakerCheckerHeaders(
        h: HttpHeaders,
        permission: 'request' | 'approve',
        approvalId?: string
    ): HttpHeaders {
        h = h.set('x-permission', permission);

        if (permission === 'approve') {
            h = h.set('x-request', approvalId);
        }
        return h;
    }

    /**
     * Used as middleware to do an action on sign-in and sign-out
     *
     * @private
     * @memberof AuthService
     */
    private watchUser() {

        // const user$ = Observable.create((observer: Observer<UserInfo>) => {

        // watch logged in user
        this.auth.auth.onAuthStateChanged((user: UserInfo) => {

            if (user) {

                // this.auth.auth.currentUser.getIdToken(false).then((fid) => {
                //     console.log('fid', fid);
                // });

                this.startTimer();

                // login success, so remove error message if exists
                this.registrationError = '';
                this.loginError = '';

                // The signed-in user info
                console.log('Signed in: ', user.email);

                // set user data on auth service
                this.getUserProfile(this);

                // tell angular about changes to detect
                // this.ngZone.run(() => {
                //     observer.next(user);
                // });

            } else {

                // clear inactivity timer by stopping
                this.stopTimer();

                // The signed-out user info.
                console.log('Signed out.');

                // tell angular about changes to detect
                // this.ngZone.run(() => {
                //     observer.next(null);
                // });

                // remove user session data
                this.userProfile = null;

                // stop observable to watch user profile when the user logs out
                if (this.profileRef) {
                    // if the profile ref is not null then stop observing
                    // because it will start re-watching below
                    this.profileRef.off();
                }

                // IMPORTANT: Better to handle navigation on SignOut()
                // so that navigation only occurs when user
                // clicks logout, instead of prompting the
                // user to login every time the user navigates to any page

            }

        });

        // });

        // return user$;

    }

    /**
     * Gets user profile including permissions if the user is logged in
     *
     * @param {AuthService} authService Needed for context in situations where context is lost due to async calls
     * @returns {Promise<IUserProfile>}
     * @memberof AuthService
     */
    getUserProfile(authService: AuthService): Promise<IUserProfile> {

        return new Promise((resolve, reject) => {

            // get user profile if there is a currently logged in user
            if (!isEmpty(get(authService.auth, 'auth.currentUser'))) {

                // check if the user already exists on authService
                if (authService.userProfile) {
                    // return current user
                    resolve(authService.userProfile);
                } else {

                    if (this.profileRef) {
                        // if the profile ref is not null then stop observing
                        // because it will start re-watching below
                        this.profileRef.off();
                    }

                    // get current user info from database
                    this.profileRef = authService.db.database.ref(`/users/${authService.auth.auth.currentUser.uid}`);

                    try {

                        // watch for changes to profile (necessary for situations when a super admin revokes permissions
                        // while the user is already logged-in, this will update the users profile to add or remove permissions
                        // even while the user is logged in)
                        this.profileRef.on('value', (userData: firebase.database.DataSnapshot) => {

                            const userProfile: IUserProfile = userData.val();

                            // set current user permissions
                            authService.userProfile = userProfile;

                            if (isDevMode) {
                                console.log('User Permissions: ', authService.userProfile);
                            }

                            resolve(userProfile);

                        });

                    } catch (error) {
                        console.log(error);
                        resolve(null);
                    }

                }

            } else {
                // no logged in user so reject
                resolve(null);
            }

        });

    }

    // /**
    //  * Gets user profile including permissions if the user is logged in
    //  *
    //  * @param {AuthService} authService Needed for context in situations where context is lost due to async calls
    //  * @returns {Promise<IUserProfile>}
    //  * @memberof AuthService
    //  */
    // getUserProfile(authService: AuthService): Promise<IUserProfile> {

    //     return new Promise((resolve, reject) => {

    //         // get user profile if there is a currently logged in user
    //         if (!isEmpty(get(authService.auth, 'auth.currentUser'))) {

    //             // check if the user already exists on authService
    //             if (authService.userProfile) {
    //                 // return current user
    //                 resolve(authService.userProfile);
    //             } else {

    //                 // get current user info from database
    //                 authService.db.database.ref(`/users/${authService.auth.auth.currentUser.uid}`)
    //                     .once('value', (userData: firebase.database.DataSnapshot) => {

    //                         const userProfile: IUserProfile = userData.val();

    //                         // set current user permissions
    //                         authService.userProfile = userProfile;
    //                         console.log('User Profile: ', authService.userProfile);

    //                         resolve(userProfile);

    //                     });

    //             }

    //         } else {
    //             // no logged in user so reject
    //             reject(null);
    //         }

    //     });
    // }

    /**
     * checks if the user has "SPECIFIC" super permissions defined
     *
     * @returns
     * @memberof AuthService
     */
    hasSpecificSuperPermissions(
        userProfile: IUserProfile,
        allowableSuperPermissions: ('admin' | 'manager' | 'viewer' | string)[]
    ) {

        // for each allowable Super User permissions type
        for (let i = 0; i < allowableSuperPermissions.length; i++) {
            // check if allowable permission exists on user profile object
            if (has(userProfile, 'super_permissions.roles.' + allowableSuperPermissions[i])) {
                // permissions allowed!
                return true;
            }
        }

        // no permissions found so return false
        return false;
    }

    /**
     * checks if the user has "ANY" super permissions defined
     *
     * @memberof AuthService
     */
    hasAnySuperPermissions(userProfile: IUserProfile) {
        // checks if the user has any participant permission roles
        const isSuper = has(userProfile, 'super_permissions.roles');
        return isSuper;
    }

    /**
     * Check if permissions to view participant related information exists for a user
     * NOTE: data is secured by firebase security rules, this logic
     * is only used for evaluating route guard permission.
     *
     * returns true or false
     * if user has permissions for the child route
     *
     * @param {IUserProfile} userProfile
     * @param {string} checkInstitutionId
     * @param {(('admin' | 'manager' | 'viewer')[])} allowablePermissions
     * @returns {boolean}
     * @memberof AuthService
     */
    hasParticipantPermissions(
        userProfile: IUserProfile,
        checkInstitutionId: string,
        allowableParticipantPermissions: ('admin' | 'manager' | 'viewer' | string)[]
    ): boolean {

        // check to make sure passed in route permissions exist
        if (!allowableParticipantPermissions) {
            return false;
        }

        // for each allowable Participant User permissions type
        for (let i = 0; i < allowableParticipantPermissions.length; i++) {
            // check if allowable permission exists on user profile object
            if (has(userProfile, 'participant_permissions.' + checkInstitutionId + '.roles.' + allowableParticipantPermissions[i])) {
                // permissions allowed!
                return true;
            }
        }

        // restrict permissions, block user from proceeding to route
        // since user does not have permissions needed
        return false;

    }

    /**
     * Helper function for Particpant Permissions.
     * Checks if user is Manager or Admin of this participant
     *
     * @param {string} institutionId
     * @returns {boolean}
     * @memberof AuthService
     */
    public userIsParticipantManagerOrHigher(institutionId: string): boolean {
        return this.hasParticipantPermissions(
            this.userProfile,
            institutionId,
            ['manager', 'admin']);
    }

    /**
     * Helper function for Particpant Permissions.
     * Checks if user is Admin of this participant
     *
     * @param {string} institutionId
     * @returns {boolean}
     * @memberof AuthService
     */
    public userIsParticipantAdmin(institutionId: string): boolean {
        return this.hasParticipantPermissions(
            this.userProfile,
            institutionId,
            ['admin']);
    }

    /**
     * Helper function for Super Permissions.
     * Checks if current user has super permissions.
     *
     * @returns {boolean}
     * @memberof AuthService
     */
    public userIsSuperUser(): boolean {
        return this.hasAnySuperPermissions(
            this.userProfile
        );
    }

    private handleLoginError(err: firebase.FirebaseError, email: string) {

        this.registrationError = ''; // reset registration error

        // If there is login error display the error message
        if (err) {

            switch (err.code) {
                case 'auth/invalid-email':
                    this.loginError = err.message;
                    break;

                case 'auth/user-disabled':
                    this.loginError = 'The user corresponding to the given email has been disabled.';
                    break;

                case 'auth/user-not-found':
                    this.loginError = 'Invalid email or password';
                    break;

                case 'auth/wrong-password':
                    this.loginError = 'Invalid email or password';
                    break;

                case 'auth/internal-error':
                    // this typically fires when the user redirects back to the site
                    // from the OAuth provider pop-up
                    this.loginError = 'Unable to login using this method at this time.';
                    break;

                default:
                    // for all OAuth error messages show full message because user may
                    // need to know detailed info as to why they weren't logged in
                    this.loginError = err.message;
                    break;
            }

        }

        console.log('Cannot sign in via oAuth: ', this.loginError);

        // transition to public login page with login error message
        // this.$state.go('site.public.login', { loginError: err, email: email });
        this.router.navigate(['/login']);

    }

    private loginRedirect() {

        // check if there is a route in query param to redirect
        this.activatedRoute
            .queryParams
            .subscribe(queryParams => {

                // console.log('Query Params:', queryParams);

                if (queryParams['redirect']) {
                    // if redirect query param is present
                    this.router.navigate(queryParams['redirect']);

                } else {
                    // if redirect query param is NOT present
                    this.router.navigate(['/portal']);

                }

            });

    }

    // ================== TIMEOUT SEVRVICE METHODS (Start) ====================

    minutesDisplay = 0;
    secondsDisplay = 0;

    unsubscribe$: Subject<void> = new Subject();
    timerSubscription: Subscription;

    // if the timer is started
    init: boolean;

    _userActionOccurred: Subject<void> = new Subject();
    get userActionOccurred(): Observable<void> { return this._userActionOccurred.asObservable(); }

    ngOnDestroy() {
        this.stopTimer();
    }

    notifyUserAction() {
        this._userActionOccurred.next();
    }

    startTimer() {

        // if timer is not initialized startTimer
        if (!this.init) {

            this.init = true;

            this.resetInactivityTimer();
            this.userActionOccurred.pipe(
                takeUntil(this.unsubscribe$)
            ).subscribe(() => {
                if (this.timerSubscription) {
                    this.timerSubscription.unsubscribe();
                }
                this.resetInactivityTimer();
            });

        }

    }

    stopTimer() {
        this.init = false;
        this.unsubscribe$.next();
        this.unsubscribe$.complete();
    }

    private resetInactivityTimer(endTime: number = environment.inactivityTimeout) {
        const interval = 1000;
        const duration = endTime * 60;
        this.timerSubscription = timer(0, interval).pipe(
            take(duration)
        ).subscribe(value =>
            this.render((duration - +value) * interval),
            err => { },
            () => {

                //  navigate to inactive page
                this.signOut('/inactive');

            }
        );
    }

    private render(count) {
        this.secondsDisplay = this.getSeconds(count);
        this.minutesDisplay = this.getMinutes(count);
    }

    private getSeconds(ticks: number) {
        const seconds = ((ticks % 60000) / 1000).toFixed(0);
        return this.pad(seconds);
    }

    private getMinutes(ticks: number) {
        const minutes = Math.floor(ticks / 60000);
        return this.pad(minutes);
    }

    private pad(digit: any) {
        return digit <= 9 ? '0' + digit : digit;
    }

    // ================== TIMEOUT SEVRVICE METHODS (End) ====================

}
