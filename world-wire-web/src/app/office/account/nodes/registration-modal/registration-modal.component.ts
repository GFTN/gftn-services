// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Inject } from '@angular/core';
import { INodeAutomation } from '../../../../shared/models/node.interface';
import { BaseModal, NotificationService } from 'carbon-components-angular';
import { HttpClient, HttpHeaders, HttpErrorResponse } from '@angular/common/http';
import { AuthService } from '../../../../shared/services/auth.service';
import { ENVIRONMENT } from '../../../../shared/constants/general.constants';
import { Participant } from '../../../../shared/models/participant.interface';
import { NgForm } from '@angular/forms';
import { WorldWireError } from '../../../../shared/models/error.interface';
import { ApprovalInfo, ApprovalPermission } from '../../../../shared/models/approval.interface';
import { SuperApprovalsModel } from '../../../../shared/models/super-approval.model';
import { AngularFireDatabase } from '@angular/fire/database';

@Component({
  selector: 'app-registration-modal',
  templateUrl: './registration-modal.component.html',
  styleUrls: ['./registration-modal.component.scss'],
  providers: [
    SuperApprovalsModel
  ]
})

export class RegistrationModalComponent extends BaseModal implements OnInit {

  node: INodeAutomation;

  loaded = false;

  accountAddress: string;

  // will check if participant is registered in PR
  // This is a prerequisite for any account registration as the participant
  // needs to exist before any accounts can be added.
  officialParticipant: Participant;

  latestApproval: ApprovalInfo;

  dbRef: firebase.database.Reference;

  saving = false;

  showLoader = false;

  success = false;

  constructor(
    @Inject('MODAL_DATA') public data: any,
    private superApprovals: SuperApprovalsModel,
    private authService: AuthService,
    private notificationService: NotificationService,
    private http: HttpClient,
    private db: AngularFireDatabase
  ) {

    super();
    this.node = data.node;

    this.dbRef = this.db.database.ref(`/participants/${this.node.institutionId}/nodes`);
  }

  async ngOnInit() {

    const accountsRequest = `https://admin.${ENVIRONMENT.envGlobalRoot}/pr/v1/admin/pr/domain/${this.node.participantId}`;

    const h: HttpHeaders = await this.authService.getFirebaseIdToken(this.node.institutionId, this.node.participantId);

    const options = {
      headers: h
    };

    try {

      // get official account registration from PR
      this.officialParticipant = await this.http.get(
        accountsRequest,
        options
      ).toPromise() as Participant;

      this.accountAddress = this.officialParticipant && this.officialParticipant.issuing_account ? this.officialParticipant.issuing_account : this.node.issuingAccount;

      this.latestApproval = this.node.accountApprovalId ? await this.superApprovals.getApprovalInfo(this.node.accountApprovalId) : null;

      this.loaded = true;

    } catch (err) {
      this.loaded = true;

      // throw error
      this.throwError(err, 'view this registration');
    }
  }

  /**
   * Returns if this participant has an officially registered account with World Wire
   *
   * @returns {boolean}
   * @memberof RegistrationModalComponent
   */
  isRegistered(): boolean {
    return (this.officialParticipant && this.officialParticipant.issuing_account) ? true : false;
  }

  /**
   * Validates form input
   *
   * @param {NgForm} f
   * @returns
   * @memberof RegistrationModalComponent
   */
  validForm(f: NgForm) {


    if (!this.officialParticipant) {
      return false;
    }

    // trim value of empty spaces
    const addressVal: string = this.accountAddress ? this.accountAddress.trim() : null;

    // init registration
    if (f.valid && addressVal && !this.latestApproval) {
      return true;
    }

    // valid form submission for maker/checker process
    if (this.latestApproval && this.latestApproval.requestInitiatedBy && addressVal) {
      return true;
    }

    return false;
  }

  /**
   * Rejects account registration request
   *
   * @returns
   * @memberof RegistrationModalComponent
   */
  async rejectRegistration() {

    // Fallback error for if someone else already approved/successfully registered an account
    if (this.officialParticipant.issuing_account) {
      const error: HttpErrorResponse = new HttpErrorResponse({
        status: 400,
      });

      this.throwError(error, 'Account is already registered. Could not reject this request.');
      return;
    }

    // disable save button to prevent double submission
    this.saving = true;

    // show saving loader
    this.showLoader = true;

    if (this.latestApproval) {

      // reset approval
      await this.dbRef.child(this.node.participantId)
        .update({
          issuingAccount: null,
          accountApprovalId: null
        });
    }

    // emit success
    this.success = true;

    // reload form state
    this.saving = false;

    this.closeModal();
  }

  /**
   * Try to send request to register a Stellar address account
   *
   * @param {NgForm} f
   * @returns
   * @memberof RegistrationModalComponent
   */
  async submitForm(f: NgForm, permission: ApprovalPermission) {

    const errorText = 'register this account';

    // DNS was not set up. Throw Error.
    if (!this.officialParticipant) {
      const error: HttpErrorResponse = new HttpErrorResponse({
        status: 404,
      });
      this.throwError(error, errorText);
      return;
    }

    // Invalid form. Throw Error
    if (!this.validForm(f)) {

      const error: HttpErrorResponse = new HttpErrorResponse({
        status: 400,
      });
      this.throwError(error, errorText);
      return;
    }

    // disable save button to prevent double submission
    this.saving = true;

    // show saving loader
    this.showLoader = true;

    const request = `https://${this.officialParticipant.id}.${ENVIRONMENT.envGlobalRoot}/anchor/v1/admin/anchor/${this.officialParticipant.id}/register`;


    let h: HttpHeaders = await this.authService.getFirebaseIdToken(this.node.institutionId, this.node.participantId);

    if (permission === 'approve') {

      h = this.authService.addMakerCheckerHeaders(h, 'approve', this.node.accountApprovalId);
    } else {
      h = this.authService.addMakerCheckerHeaders(h, 'request');
    }

    const options = {
      headers: h
    };

    const body = {
      address: this.accountAddress
    };

    // get official account registration from PR
    const apiRequest: Promise<any> = this.http.post(
      request,
      body,
      options
    ).toPromise();

    apiRequest.then(async (response: any) => {
      // update record for checker to see approval
      if (permission === 'request' && response.msg) {

        // create new record, if does not exist
        if (!this.node.accountApprovalId) {
          await this.dbRef.child(this.node.participantId).update({
            'accountApprovalId': response.msg,
            'issuingAccount': this.accountAddress
          });
        }
      }

      // emit success
      this.success = true;

      // reload form state
      this.saving = false;

      this.closeModal();

    }).catch((error: HttpErrorResponse) => {
      this.throwError(error, errorText);

      // Error in approval request. Revert approval back to 'request' state to allow retries
      if (permission === 'approve') {
        // remove approval if errored out, to allow retry
        this.db.database.ref('/super_approvals').child(this.node.accountApprovalId).update({
          uid_approve: null,
          status: 'request'
        });
      }
    });

  }

  /**
   * Helper function for user error handling
   *
   * @private
   * @param {HttpErrorResponse} error
   * @param {string} [optionalText]
   * @memberof RegistrationModalComponent
   */
  private throwError(error: HttpErrorResponse, optionalText?: string) {
    // reload form state
    this.saving = false;

    // disable loader
    this.showLoader = false;

    const errorText = optionalText ? optionalText : ' make this request.';

    const wwError: WorldWireError = error.error ? error.error : null;

    switch (error.status) {
      case 401:
      case 0: {
        // Unauthorized ERROR notification
        this.notificationService.showNotification({
          type: 'error',
          title: 'Unauthorized',
          message: 'You are not authorized to ' + errorText + '. Please contact administrator.',
          target: '#notification'
        });
        break;
      }
      case 400: {

        const defaultText: string = optionalText ? optionalText : 'Invalid form data was submitted and could not be processed. Please try again.';

        const message = wwError && wwError.details ? wwError.details + wwError.message : defaultText;
        // Unauthorized ERROR notification
        this.notificationService.showNotification({
          type: 'error',
          title: 'Bad Request',
          message: message,
          target: '#notification'
        });
        break;
      }
      case 404: {
        const message = wwError && wwError.details ? wwError.details + wwError.message : 'Network could not be reached to make this request. Please contact administrator.';
        // Network error notification
        this.notificationService.showNotification({
          type: 'error',
          title: 'Network Error',
          message: message,
          target: '#notification'
        });
        break;
      }
      default: {
        // for 500s and other errors: Unexpected ERROR notification
        this.notificationService.showNotification({
          type: 'error',
          title: 'Unexpected Error',
          message: 'Unable to ' + errorText + '. Please contact administrator.',
          target: '#notification'
        });
        break;
      }
    }
  }

  /**
   * Handles successful response from form
   */
  onSuccess() {
    this.showLoader = false;

    this.success = true;

    // close modal
    this.closeModal();
  }
}
