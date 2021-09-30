// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, HostBinding, NgZone } from '@angular/core';
import { AuthService } from '../../shared/services/auth.service';
import { ActivatedRoute, Router } from '@angular/router';
import { IUserProfile } from '../../shared/models/user.interface';
import { get, isEmpty } from 'lodash';

@Component({
  selector: 'app-login',
  templateUrl: './login-token.component.html',
  styleUrls: ['./login-token.component.scss']
})
export class LoginTokenComponent implements OnInit {

  constructor(
    public authService: AuthService,
    private route: ActivatedRoute,
    private ngZone: NgZone,
    private router: Router
  ) { }

  @HostBinding('attr.class') cls = 'flex-fill';

  ngOnInit() {

    // NOTE: no need to put this in a route guard or parent component
    // because this is only called after IBMId login callback is
    // fired from back-end node.js server.

    // check if firebase user exists
    if (isEmpty(get(this.authService, 'auth.auth.currentUser'))) {

      // get data from route params
      const data = this.route.snapshot.queryParams as { token: string };

      // if token is provided in query param
      // try to use token to login to firebase
      if (data.token) {

        // login user using custom firebase token
        this.authService.firebaseCustomTokenSignIn()
          .then((userPermissions: IUserProfile) => {
            // user now logged in so redirect
            this.redirect(userPermissions);
          });

      } else {

        // redirect user to login via IBM Id
        this.authService.signInIbmId();

      }

    } else {

      // already logged in so redirect
      this.authService.getUserProfile(this.authService)
        .then((profile: IUserProfile) => {
          this.redirect(profile);
        });

    }

  }

  /**
   * redirect user to where they have permissions
   *
   * @param {IUser} userProfile
   * @memberof LoginTokenComponent
   */
  redirect(userProfile: IUserProfile) {

    if (!isEmpty(get(userProfile, 'super_permissions'))) {

      // if user is super admin redirect to office
      this.ngZone.run(() => {
        return this.router.navigate(['/office']);
      });

    } else if (!isEmpty(get(userProfile, 'participant_permissions'))) {

      // if user is client admin redirect to portal
      this.ngZone.run(() => {
        return this.router.navigate(['/portal']);
      });

    } else {

      // no known permissions
      this.ngZone.run(() => {
        return this.router.navigate(['/unrecognized']);
      });

    }

  }

}
