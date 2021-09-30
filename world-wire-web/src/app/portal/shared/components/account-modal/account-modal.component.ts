// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Inject } from '@angular/core';
import { BaseModal, NotificationService } from 'carbon-components-angular';
import { AngularFireDatabase } from '@angular/fire/database';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { AccountService } from '../../services/account.service';
import { SuperApprovalsModel } from '../../../../shared/models/super-approval.model';
import { ApprovalInfo } from '../../../../shared/models/approval.interface';
import { AccountRequest, ParticipantAccount } from '../../../../shared/models/account.interface';
import { kebabCase } from 'lodash';
import { WorldWireError } from '../../../../shared/models/error.interface';
import { SessionService } from '../../../../shared/services/session.service';
import { AuthService } from '../../../../shared/services/auth.service';

export interface AccountForm {
  accountName: string;
  accountId: string;
}

@Component({
  selector: 'app-account-modal',
  templateUrl: './account-modal.component.html',
  styleUrls: ['./account-modal.component.scss'],
  providers: [
    SuperApprovalsModel
  ]
})
export class AccountModalComponent extends BaseModal implements OnInit {

  // current account request available for approval/rejection
  currentAccountRequest: AccountRequest;

  latestApproval: ApprovalInfo;

  formData: AccountForm;

  loadingRequest = false;

  // id reference to the notification element in the UI to display errors
  notificationRef = '#error-notication';

  registered = false;

  constructor(
    @Inject('MODAL_DATA') public data: AccountRequest,
    private http: HttpClient,
    private accountService: AccountService,
    private sessionService: SessionService,
    private authService: AuthService,
    private notificationService: NotificationService,
    private superApprovals: SuperApprovalsModel,
    private db: AngularFireDatabase,
  ) {
    super();

    this.currentAccountRequest = data;
  }

  async ngOnInit() {

    this.formData = {
      accountName: '',
      accountId: ''
    };

    if (this.currentAccountRequest) {

      // get approval if it exists
      if (this.currentAccountRequest.approvalIds) {
        const latestApprovalId = this.currentAccountRequest.approvalIds[this.currentAccountRequest.approvalIds.length - 1];

        this.latestApproval = await this.superApprovals.getApprovalInfo(latestApprovalId);
      }

      // check if asset is already registered
      const foundAccount: ParticipantAccount = this.accountService.allAccounts[this.currentAccountRequest.name];

      // asset is registered if already approved through portal/added via API
      if (foundAccount || (this.latestApproval && this.latestApproval.requestApprovedBy)) {
        this.registered = true;
      }
    }
  }

  /**
   * Converts account name to url-safe form
   *
   * @memberof AccountModalComponent
   */
  getAccountId() {
    this.formData.accountId = kebabCase(this.formData.accountName);
  }

  /**
   * Makes request to API to request account creation request
   *
   * @returns
   * @memberof AccountModalComponent
   */
  public async requestAccountRequest() {
    this.loadingRequest = true;

    if (!this.formData.accountId) {
      const error: HttpErrorResponse = new HttpErrorResponse({
        status: 400,
      });
      this.throwError(error);

      return;
    }

    // Secondary check to prevent unauthorized requests
    if (!this.authService.userIsSuperUser()) {
      const error: HttpErrorResponse = new HttpErrorResponse({
        status: 401,
      });
      this.throwError(error);

      return;
    }

    try {

      const requestUrl = `${this.accountService.apiRoot}/v1/admin/accounts/${this.formData.accountId}`;

      let h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

      h = this.authService.addMakerCheckerHeaders(h, 'request');

      const options = {
        headers: h
      };


      const response: WorldWireError = await this.http.post(requestUrl, null, options).toPromise() as WorldWireError;

      const newRequest: AccountRequest = {
        name: this.formData.accountId,
        approvalIds: [response.msg]
      };

      this.db.database.ref('account_requests')
        .child(this.sessionService.currentNode.participantId)
        .child(this.formData.accountId)
        .update(newRequest);

      this.loadingRequest = false;

      this.closeModal();
    } catch (err) {
      this.throwError(err);
    }
  }

  /**
   * Makes request to API to approve account creation request
   *
   * @returns
   * @memberof AccountModalComponent
   */
  public async approveAccountRequest() {
    // Secondary Check: maker CANNOT be same as checker. throw 401 unauthorized error.
    if (this.latestApproval && this.latestApproval.requestInitiatedBy === this.authService.userProfile.profile.email) {
      const error: HttpErrorResponse = new HttpErrorResponse({
        status: 401,
      });
      this.throwError(error, 'Approver of request cannot be the same as requestor.');
      return;
    }

    // Secondary Check: checker must be super user
    if (!this.authService.userIsSuperUser()) {
      const error: HttpErrorResponse = new HttpErrorResponse({
        status: 401,
      });
      this.throwError(error, 'Approver must be a super user');
      return;
    }

    this.loadingRequest = true;

    try {

      const requestUrl = `${this.accountService.apiRoot}/v1/admin/accounts/${this.currentAccountRequest.name}`;

      let h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

      const approveId = this.latestApproval.key;

      h = this.authService.addMakerCheckerHeaders(h, 'approve', approveId);

      const options = {
        headers: h
      };


      await this.http.post(requestUrl, null, options).toPromise();


      this.loadingRequest = false;

      this.closeModal();

    } catch (err) {

      // error occurred. reset approval process
      if (this.latestApproval) {
        await this.superApprovals.resetApprovals(this.latestApproval.key);
      }

      this.throwError(err);
    }
  }

  /**
   * Rejects account creation request
   *
   * @returns
   * @memberof AccountModalComponent
   */
  public async rejectAccountRequest() {

    // Secondary Check: maker CANNOT be same as checker. throw 401 unauthorized error.
    if (this.latestApproval && this.latestApproval.requestInitiatedBy === this.authService.userProfile.profile.email) {
      const error: HttpErrorResponse = new HttpErrorResponse({
        status: 401,
      });
      this.throwError(error, 'Rejector of request cannot be the same as requestor.');
      return;
    }

    this.loadingRequest = true;

    try {

      // remove approvalId from current list since this is no longer valid.
      // forces maker request to be remade for this request
      const newApprovalIds: string[] = this.currentAccountRequest.approvalIds.filter((id: string) => {
        return id !== this.latestApproval.key;
      });

      // remove approval to reset request
      this.db.database.ref('account_requests')
        .child(this.sessionService.currentNode.participantId)
        .child(this.currentAccountRequest.name)
        .update({
          approvalIds: newApprovalIds
        });

      this.loadingRequest = false;

      this.closeModal();
    } catch (err) {

      this.throwError(err);
    }
  }

  /**
   * Generic error handler
   *
   * @private
   * @param {HttpErrorResponse} err
   * @param {string} [optionalErrMsg]
   * @memberof AccountModalComponent
   */
  private throwError(err: HttpErrorResponse, optionalErrMsg?: string) {
    this.loadingRequest = false;

    const wwError: WorldWireError = err.error ? err.error : null;

    switch (err.status) {
      case 401:
      case 0: {
        // Unauthorized ERROR notification
        this.notificationService.showNotification({
          type: 'error',
          title: 'Unauthorized',
          message: 'You are not authorized to make this request. ' + optionalErrMsg + ' Please contact administrator.',
          target: this.notificationRef
        });
        break;
      }
      case 400: {
        const message = wwError && wwError.details ? wwError.details + wwError.message : 'Invalid form data was submitted and could not be processed. Please try again.';
        // Unauthorized ERROR notification
        this.notificationService.showNotification({
          type: 'error',
          title: 'Bad Request',
          message: message,
          target: this.notificationRef
        });
        break;
      }
      case 404: {
        const message = wwError && wwError.details ? wwError.details + wwError.message : 'Network could not be reached to make this request. Please contact administrator.';

        // Network error notification
        this.notificationService.showNotification({
          type: 'error',
          title: 'Bad Request',
          message: message,
          target: this.notificationRef
        });
        break;
      }
      default: {

        const message = wwError && wwError.details ? wwError.details + wwError.message : 'Unexpected error found when making this request. Please contact administrator.';

        // for 503s and other errors: Unexpected ERROR notification
        this.notificationService.showNotification({
          type: 'error',
          title: 'Unexpected Error',
          message: message,
          target: this.notificationRef
        });
        break;
      }
    }
  }
}
