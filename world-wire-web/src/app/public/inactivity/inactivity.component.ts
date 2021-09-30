// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, HostBinding } from '@angular/core';
import { AuthService } from '../../shared/services/auth.service';
import { get, isEmpty } from 'lodash';
import { Router } from '@angular/router';

/**
 * Displays inactivity timer page
 *
 * @export
 * @class InactivityComponent
 * @implements {OnInit}
 */
@Component({
  selector: 'app-inactivity',
  templateUrl: './inactivity.component.html',
  styleUrls: ['./inactivity.component.scss']
})
export class InactivityComponent implements OnInit {

  constructor(
    private authService: AuthService,
    private router: Router
  ) { }

  @HostBinding('attr.class') cls = 'flex-fill';

  ngOnInit() {

    // presume that this page can only be reached after inactivity-timer.component.ts
    // calls authService.signOut('/inactive') and user is logged out
    if (!isEmpty(get(this.authService.auth, 'auth.currentUser'))) {
      // if the user calls the page directly and user exists they should
      // be redirected to either /office or /portal depending on permissions
      // this logic is handled by /login
      this.router.navigate(['/login'], { queryParams: { token: 'empty' } });
    }

  }

}
