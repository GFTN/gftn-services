// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit } from '@angular/core';
import { AuthService } from '../../shared/services/auth.service';

@Component({
  selector: 'app-office-side-nav',
  templateUrl: './office-side-nav.component.html',
  styleUrls: ['./office-side-nav.component.scss']
})
export class OfficeSideNavComponent implements OnInit {

  constructor(
    private authService: AuthService
  ) { }

  ngOnInit() {
  }

  /**
   * Toggles showing link if the user does or does not have elevated permissions
   *
   * @param {string[]} superPermission
   * @memberof OfficeSideNavComponent
   */
  public show(superPermission: string[]) {

    // check if super permissions exist
    if (this.authService.hasSpecificSuperPermissions(this.authService.userProfile, superPermission)) {
      return true;
    }

    return false;

  }

}
