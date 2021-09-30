// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { IRolesOptions } from '../../shared/models/user.interface';
import {
  ParticipantPermissionsDialogComponent, IParticipantPermissionsDialogData
} from '../../shared/components/participant-permissions-dialog/participant-permissions-dialog.component';
import { SessionService } from '../../shared/services/session.service';
import { ParticipantPermissionsService } from '../../shared/services/participant-permissions.service';
import { AuthService } from '../../shared/services/auth.service';
import { Observable, Subscription } from 'rxjs';
import { KeyValue } from '@angular/common/src/pipes/keyvalue_pipe';
import { IParticipantUsers, IParticipantUser } from '../../shared/models/participant.interface';

@Component({
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.scss']
})
export class UsersComponent implements OnInit, OnDestroy {

  users: IParticipantUsers;
  disable: any;
  humanizeRoles: any;

  public $users: Observable<IParticipantUsers>;

  usersSubscription: Subscription;

  emailSortDesc = (a: KeyValue<string, IParticipantUser>, b: KeyValue<string, IParticipantUser>) => {

    const emailA = a.value.profile.email.toLowerCase();
    const emailB = b.value.profile.email.toLowerCase();

    if (emailA > emailB) {
      return 1;
    } else if (emailA < emailB) {
      return -1;
    }
    return 0;
  }

  constructor(
    public dialog: MatDialog,
    public participantPermissionsService: ParticipantPermissionsService,
    public authService: AuthService,
    private session: SessionService
  ) {
    this.disable = this.participantPermissionsService.disable;
    this.humanizeRoles = this.participantPermissionsService.humanizeRoles;
  }

  ngOnInit() {

    // change to observable to pull data from trigger run
    this.$users = this.participantPermissionsService
      .getAllUsersObservable(this.session.institution.info.institutionId);

    this.usersSubscription = this.$users.subscribe((users: IParticipantUsers) => {
      this.users = users;
    });

  }

  ngOnDestroy() {

    // cleanup subscriptions
    if (this.usersSubscription) {
      this.usersSubscription.unsubscribe();
    }
  }

  /**
   * add user button click to open modal
   *
   * @memberof UsersComponent
   */
  addUser() {
    this.openUserDialog('add');
  }

  /**
   * add user button click to open modal
   *
   * @memberof UsersComponent
   */
  editUser(email: string, userId: string, role: IRolesOptions) {
    this.openUserDialog('edit', email, userId, role);
  }

  deleteUser(email: string, userId: string, role: IRolesOptions) {
    this.openUserDialog('remove', email, userId, role);
  }

  /**
   * open the user dialog for edit and save
   *
   * @private
   * @param {('add' | 'edit')} action
   * @param {string} [email] used to 'remove' or 'edit'
   * @param {string} [userId] used to 'remove' or 'edit'
   * @memberof UsersComponent
   */
  private openUserDialog(action: 'add' | 'edit' | 'remove', email?: string, userId?: string, role?: IRolesOptions) {

    const data: IParticipantPermissionsDialogData = {
      action: action,
      institution: this.session.institution,
      email: email,
      userId: userId,
      role: role
    };

    const dialogRef = this.dialog.open(ParticipantPermissionsDialogComponent, {
      disableClose: true,
      data: data
    });

    dialogRef.afterClosed().subscribe(result => {});

  }

}
