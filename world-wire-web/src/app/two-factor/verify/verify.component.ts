// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit } from '@angular/core';
import { Confirm2faService } from '../../shared/services/confirm2fa.service';
import { environment } from '../../../environments/environment';
import { Router } from '@angular/router';
import { HttpRequest } from '@angular/common/http';
import { SessionService } from '../../shared/services/session.service';

@Component({
  selector: 'app-verify',
  templateUrl: './verify.component.html',
  styleUrls: ['./verify.component.scss']
})
export class VerifyComponent implements OnInit {

  constructor(
    private session: SessionService,
    private confirm2fa: Confirm2faService,
    private router: Router
  ) { }

  ngOnInit() {

    // get authentication code IBMId
    const r = new HttpRequest(
      'POST',
      environment.apiRootUrl + '/sso/portal-login-totp',
      null,
      { withCredentials: true }
    );

    this.confirm2fa.go(r)
      .then((data: any) => {
        if (data.token) {

          // save token to service
          this.session.firebaseTempAuthToken = data.token;

          // navigate to /login with custom firebase login token
          this.router.navigate(['/login'], { queryParams: { token: 'check' } });
        } else {
          // failed to get auth token so redirect to unauthorized
          this.unauthorized();
        }
      }, (err: any) => {
        console.log(err);
        // failed 2fa so redirect to unauthorized
        this.unauthorized();
      });

  }

  /**
   * Redirects page to 'unauthorized' route
   *
   * @memberof VerifyComponent
   */
  unauthorized() {
    this.router.navigate(['/unauthorized']);
  }

}
