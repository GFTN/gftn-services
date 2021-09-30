// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { OnInit, Component, Inject, ViewChild } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { IRolesOptions } from '../../models/user.interface';
import { ParticipantPermissionsService } from '../../services/participant-permissions.service';
import { IInstitution } from '../../models/participant.interface';
import { NgForm } from '@angular/forms';
import { NotificationService } from '../../services/notification.service';
import { UtilsService } from '../../utils/utils';
import { toArray, keys } from 'lodash';
import { CUSTOM_REGEXES, RegexMap } from '../../constants/regex.constants';

export interface IParticipantPermissionsDialogData {
  action: 'add' | 'edit' | 'remove';
  institution: IInstitution;

  // email and user id provided for when action 'edit' or 'remove'
  email?: string;
  userId?: string;
  role?: IRolesOptions;
}

@Component({
  templateUrl: './participant-permissions-dialog.component.html',
  styleUrls: ['./participant-permissions-dialog.component.scss'],
  providers: [
    ParticipantPermissionsService
  ]
})
export class ParticipantPermissionsDialogComponent implements OnInit {

  @ViewChild('userForm') public userForm: NgForm;
  userId: string;
  email: string;
  role: 'admin' | 'manager' | 'viewer' | '' = '';
  invalidForm: boolean;
  institution: IInstitution;

  regexes: RegexMap;

  constructor(
    public dialogRef: MatDialogRef<ParticipantPermissionsDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: IParticipantPermissionsDialogData,
    public permissionsService: ParticipantPermissionsService,
    private notification: NotificationService,
    public utils: UtilsService
  ) {
    this.regexes = CUSTOM_REGEXES;

    this.institution = data.institution;
    if (data.userId && data.email) {
      // set user and set role
      this.email = data.email;
      this.userId = data.userId;

      // // set user role
      // * NOTE: Even though our design allows for multiple roles
      // * this only selects the first one and presumable the only
      // * role that a participant should have since our system is
      // * currently designed so that user only have a single role
      this.role = keys(data.role)[0] as any;
    }
  }

  ngOnInit(): void { }

  async update(userForm: NgForm): Promise<any> {

    const self = this;

    // check if form is valid
    if (this.email && this.role && userForm.valid) {

      // normalize email
      this.email = this.email.toLowerCase().trim();

      // get new user userId for email
      // this.permissionsService.setUserId(this.email).then((uid) => {

      return await this.permissionsService.update(this.institution.info.institutionId, this.role, this.email).then(() => {

        // display success notification
        self.notification.show('success',
          self.utils.capitalizeFirstLetter(self.data.action) + ' user permissions for ' + self.email
        );

        // close modify view
        self.close();

      }, (err: any) => {
        console.log('Failed to update permissions for user: ', err);
      });

    } else {

      // indicates that an email is required
      this.invalidForm = true;

      return;
    }

  }

  async remove(): Promise<any> {

    return await this.permissionsService.remove(this.userId, this.data.institution.info.institutionId)
      .then(() => {
        this.notification.show('success', 'Removed user permissions for ' + this.email);
        this.close();
      });

  }

  /**
   * closes the modal
   *
   * @private
   * @memberof ManageUserDialogComponent
   */
  close(): void {
    this.dialogRef.close();
  }

}
