// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, OnDestroy, isDevMode, Input, Output, EventEmitter, ComponentRef } from '@angular/core';
import { ParticipantAccountDetail, AccountRequest } from '../../../shared/models/account.interface';
import { SessionService } from '../../../shared/services/session.service';
import { AuthService } from '../../../shared/services/auth.service';
import { NotificationService } from '../../../shared/services/notification.service';
import { ParticipantAccountService } from '../shared/services/participant-account-service.service';
import { AccountService } from '../../shared/services/account.service';
import { AccountModalComponent } from '../../shared/components/account-modal/account-modal.component';
import { ModalService } from 'carbon-components-angular';

@Component({
  selector: 'app-participant-account',
  templateUrl: './participant-account.component.html',
  styleUrls: ['./participant-account.component.scss'],
  providers: [
    ParticipantAccountService
  ]
})
export class ParticipantAccountComponent implements OnInit, OnDestroy {

  @Input() account: ParticipantAccountDetail;

  // default view: list
  @Input() view: 'list' | 'grid' = 'list';

  isSubmitting = false;

  @Output() accountChanged = new EventEmitter<string>();

  suspendReactivateDbRef: firebase.database.Reference;

  participantAuthorized = false;

  superAuthorized = false;

  currentOpenModal: ComponentRef<AccountModalComponent>;

  constructor(
    private authService: AuthService,
    public sessionService: SessionService,
    public participantAccountService: ParticipantAccountService,
    private notificationService: NotificationService,
    private modalService: ModalService
  ) { }

  ngOnInit() {

    this.participantAuthorized = this.sessionService.institution && this.authService.userIsParticipantManagerOrHigher(
      this.sessionService.institution.info.institutionId
    );

    this.superAuthorized = this.authService.userIsSuperUser();

    // TODO: Enable for regular participant admins. For now, this is enabled only for IBM/super users.
    if (this.superAuthorized) {
      this.participantAccountService.getSupendReactivateRequest(this.account);
    }

    if (!this.account.request) {
      this.account.request = {
        name: this.account.name,
        approvalIds: null
      };
    }
  }

  ngOnDestroy() {

    // clean up memory from observable
    if (this.participantAccountService.suspendReactivateDbRef) {
      this.participantAccountService.suspendReactivateDbRef.off();
    }
  }

  /**
   * Used in view to request suspension/reactivation of account
   *
   * @param {string} accountAddress
   * @param {boolean} suspend
   * @returns {Promise<void>}
   * @memberof ParticipantAccountComponent
   */
  public async requestSuspendReactivateAccount(suspend: boolean): Promise<void> {
    this.isSubmitting = true;
    try {
      await this.participantAccountService.requestSuspendReactivateAccount(this.account, suspend);

      this.notificationService.show('success', 'Request Submitted');

    } catch (err) {
      if (isDevMode) {
        console.log(err);
      }

      const errMsg = err.error ? err.error.details : 'Unexpected error found when submit suspend account request.';

      this.notificationService.show('error', errMsg);
    }
    this.isSubmitting = false;
  }

  /**
   * Used in view to approve suspension/reactivation request of account
   *
   * @param {string} accountAddress
   * @returns {Promise<void>}
   * @memberof ParticipantAccountComponent
   */
  public async approveSuspendReactivateAccount(): Promise<void> {
    this.isSubmitting = true;

    try {

      await this.participantAccountService.approveSuspendReactivateAccount(this.account);

      this.notificationService.show('success', 'Request Submitted');

    } catch (err) {
      if (isDevMode) {
        console.log(err);
      }

      this.notificationService.show('error', err.error.details || 'Unexpected error found when submit suspend account request.');

    }
    this.isSubmitting = false;
  }

  /**
   * Used in view to reject suspension/reactivation request of account
   *
   * @param {string} accountAddress
   * @returns {Promise<void>}
   * @memberof ParticipantAccountComponent
   */
  public async rejectSuspendReactivateAccount(): Promise<void> {
    this.isSubmitting = true;

    try {

      await this.participantAccountService.rejectSuspendReactivateAccount(this.account);

      this.notificationService.show('success', 'Request Submitted');

    } catch (err) {
      if (isDevMode) {
        console.log(err);
      }

      this.notificationService.show('error', err.error.details || 'Unexpected error found when submit suspend account request.');
    }
    this.isSubmitting = false;
  }

  /**
   * Used in view to toggle whether or not this person is checker (maker/checker flow)
   *
   * @param {string} accountAddress
   * @returns {boolean}
   * @memberof ParticipantAccountComponent
   */
  public getApproveRejectButton(): boolean {

    // user not logged in. Default to false
    if (!this.authService.userProfile) {
      return false;
    }

    return this.participantAccountService.getRequesterEmail(this.account) === this.authService.userProfile.profile.email;
  }

  public viewAccountApproval(request: AccountRequest) {

    this.currentOpenModal = this.modalService.create({
      component: AccountModalComponent,
      inputs: {
        MODAL_DATA: request
      }
    });
  }
}
