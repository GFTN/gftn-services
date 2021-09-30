// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, NgZone } from '@angular/core';
import { AuthService } from '../../shared/services/auth.service';

/**
 * Used by developer to get an fid header to authenticate to
 * middleware in World Wire backend auth-service
 *
 * @export
 * @class FidComponent
 * @implements {OnInit}
 */
@Component({
  selector: 'app-fid',
  templateUrl: './fid.component.html',
  styleUrls: ['./fid.component.scss']
})
export class FidComponent implements OnInit {

  // the firebase user id for use in x-fid header to make calls against the backend
  fid: string;

  constructor(
    public authService: AuthService,
    private ngZone: NgZone
  ) { }

  ngOnInit() {
    this.init();
  }

  init() {
    if (this.authService.auth.auth.currentUser) {
      // get fid to pass along in header
      this.ngZone.run(async () => {
        this.fid = await this.authService.auth.auth.currentUser.getIdToken(false);
      });
    } else {
      // redirect user to login via IBM Id
      this.authService.signInIbmId();
    }
  }

}
