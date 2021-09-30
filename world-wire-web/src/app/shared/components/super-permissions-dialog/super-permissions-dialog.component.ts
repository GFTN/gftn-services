// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { OnInit, Component, Inject, ViewChild } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { IRolesOptions } from '../../models/user.interface';
import { NgForm } from '@angular/forms';
import { NotificationService } from '../../services/notification.service';
import { UtilsService } from '../../utils/utils';
import { toArray, keys } from 'lodash';
import { SuperPermissionsService } from '../../services/super-permissions.service';

export interface ISuperPermissionsDialogData {
  action: 'add' | 'edit' | 'remove';

  // email and user id provided for when action 'edit' or 'remove'
  email?: string;
  userId?: string;
  role?: IRolesOptions;
}

@Component({
  templateUrl: './super-permissions-dialog.component.html',
  styleUrls: ['./super-permissions-dialog.component.scss'],
  providers: [
    SuperPermissionsService
  ]
})

export class SuperPermissionsDialogComponent implements OnInit {

  @ViewChild('userForm') public userForm: NgForm;
  userId: string;
  email: string;
  role: 'admin' | 'manager' | 'viewer' | '' = '';
  invalidForm: boolean;

  constructor(
    public dialogRef: MatDialogRef<SuperPermissionsDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: ISuperPermissionsDialogData,
    public superPermissionsService: SuperPermissionsService,
    private notification: NotificationService,
    public utils: UtilsService
  ) {

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


  ngOnInit(): void {


  }

  async update(userForm: NgForm): Promise<any> {

    const self = this;


    // check if form is valid
    if (this.email && this.role) {

      // normalize email
      this.email = this.email.toLowerCase().trim();

      // get new user userId for email
      // this.permissionsService.setUserId(this.email).then((uid) => {

      return this.superPermissionsService.update(this.role, this.email).then(() => {

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

    return this.superPermissionsService.remove(this.userId)
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
