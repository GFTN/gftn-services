// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit } from '@angular/core';
import { AuthService } from '../../shared/services/auth.service';
import { IParticipantRoles } from '../../shared/models/user.interface';
import { size, get } from 'lodash';
import { Router } from '@angular/router';

/**
 * Used when user initially logs in with one or more participant_permissions.
 * The logic here will allow a non-super admin user to select which participant
 * they want to manage. If it's a super admin they should be redirected to /office
 * where they can directly select a participant to manage.
 *
 * @export
 * @class SelectComponent
 * @implements {OnInit}
 */
@Component({
  selector: 'app-select',
  templateUrl: './select.component.html',
  styleUrls: ['./select.component.scss']
})
export class SelectComponent implements OnInit {

  // used to display participant options to select from if
  // the user has permissions to more than one participant
  participant_permissions: {
    [institutionId: string]: IParticipantRoles;
  };

  constructor(
    private authService: AuthService,
    private router: Router
  ) { }

  ngOnInit() {

    // set participant permissions
    this.participant_permissions = get(this.authService, 'userProfile.participant_permissions');

    // if super admin redirect to office
    if (this.authService.hasAnySuperPermissions(this.authService.userProfile)) {

      // show selection of participants,
      // if any specific participant assocations are assigned
      if (this.participant_permissions) {
        return;
      }

      // is a super admin, so no specific participant associations
      // redirect to /office
      return this.router.navigate(['/office']);
    }

    if (size(this.participant_permissions) > 1) {

      // user has permissions to more than one participant
      // to choose from, so display which participants the
      // user has to choose from in the view
      return; // do nothing, just display '/portal/select' view to select participant

    } else if (size(this.participant_permissions) === 1) {

      // redirect to the only participant the user has permission for
      // IE: force their selection
      const iid = Object.keys(this.participant_permissions)[0];
      const navigate = ['/portal/' + this.participant_permissions[iid].slug];
      return this.router.navigate(navigate);

    } else {

      // no participant permissions associated with this non-super user
      // so redirect to /unrecognized
      return this.router.navigate(['/unrecognized']);

    }

  }

}
