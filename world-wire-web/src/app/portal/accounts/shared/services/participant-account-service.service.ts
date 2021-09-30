// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable, isDevMode } from '@angular/core';
import { ENVIRONMENT } from '../../../../shared/constants/general.constants';
import { AuthService } from '../../../../shared/services/auth.service';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { AngularFireDatabase } from '@angular/fire/database';
import { SessionService } from '../../../../shared/services/session.service';
import { ParticipantAccount } from '../../../../shared/models/account.interface';
import { WorldWireError } from '../../../../shared/models/error.interface';
import { KillSwitchRequest, KillSwitchRequestDetail, KillSwitchRequestStatus } from '../../../shared/models/killswitch-request.interface';

@Injectable()
export class ParticipantAccountService {

  suspendReactivateRequest: KillSwitchRequestDetail;

  suspendReactivateDbRef: firebase.database.Reference;

  constructor(
    private authService: AuthService,
    private sessionService: SessionService,
    private http: HttpClient,
    private db: AngularFireDatabase
  ) { }

  /**
   * Get current killswitch request if it exists
   *
   * @private
   * @memberof ParticipantAccountService
   */
  public getSupendReactivateRequest(account: ParticipantAccount) {

    if (this.sessionService.currentNode) {

      this.suspendReactivateDbRef = this.db.database.ref('killswitch_requests')
        .child(this.sessionService.currentNode.participantId)
        .child(account.address);

      this.suspendReactivateDbRef.on('value', (request: firebase.database.DataSnapshot) => {

        if (request.val()) {
          this.suspendReactivateRequest = request.val();

          this.suspendReactivateRequest.loaded = true;
        }
      });
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
  public async requestSuspendReactivateAccount(account: ParticipantAccount, suspend: boolean): Promise<KillSwitchRequestDetail> {

    const requestUrl = `https://${this.sessionService.currentNode.participantId}.${ENVIRONMENT.envGlobalRoot}/admin/v1/admin/${suspend ? 'suspend' : 'reactivate'}/${this.sessionService.currentNode.participantId}/${account.address}`;

    let h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

    h = this.authService.addMakerCheckerHeaders(h, 'request');

    const options = { headers: h };
    if (isDevMode) {
      console.log(requestUrl);
      console.log(options);
    }

    const response: WorldWireError = await this.http.post(requestUrl, null, options).toPromise() as WorldWireError;
    if (isDevMode) {
      console.log(response);
    }

    if (response.msg) {
      const newRequest: KillSwitchRequest = {
        participantId: this.sessionService.currentNode.participantId,
        accountAddress: account.address,
        [suspend && 'suspendRequestedBy']: this.authService.auth.auth.currentUser.email,
        [!suspend && 'reactivateRequestedBy']: this.authService.auth.auth.currentUser.email,
        approvalIds: [response.msg],
        status: suspend ? 'suspend_requested' : 'reactivate_requested',
      };

      await this.suspendReactivateDbRef.update(newRequest);

      let getAccountRequest: KillSwitchRequestDetail = this.suspendReactivateRequest;

      getAccountRequest = getAccountRequest ? Object.assign(getAccountRequest, newRequest) : newRequest;

      return getAccountRequest;
    }

    return null;
  }

  public async approveSuspendReactivateAccount(account: ParticipantAccount): Promise<KillSwitchRequestDetail> {
    const suspend = this.getAccountStatus(account.address) === 'suspend_requested' ? true : false;

    const requestUrl = `https://${this.sessionService.currentNode.participantId}.${ENVIRONMENT.envGlobalRoot}/admin/v1/admin/${suspend ? 'suspend' : 'reactivate'}/${this.sessionService.currentNode.participantId}/${account.address}`;

    let h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

    // error out if suspend/reactivate request was not found
    if (!this.suspendReactivateRequest || !this.suspendReactivateRequest.approvalIds) {
      return null;
    }

    const approveId = this.suspendReactivateRequest.approvalIds.slice(-1)[0];
    h = this.authService.addMakerCheckerHeaders(h, 'approve', approveId);


    const options = { headers: h };
    if (isDevMode) {
      console.log(requestUrl);
      console.log(options);
    }
    const response: WorldWireError = await this.http.post(requestUrl, null, options).toPromise() as WorldWireError;

    if (isDevMode) {
      console.log(response);
    }

    if (this.suspendReactivateDbRef) {
      const newRequest: KillSwitchRequest = {
        participantId: this.sessionService.currentNode.participantId,
        accountAddress: account.address,
        [suspend && 'suspendApprovedBy']: this.authService.auth.auth.currentUser.email,
        [!suspend && 'reactivateApprovedBy']: this.authService.auth.auth.currentUser.email,
        approvalIds: [],
        status: suspend === true ? 'suspended' : 'normal',
      };

      // capture requestId in meta data
      await this.suspendReactivateDbRef.update(newRequest);
      let getAccountRequest: KillSwitchRequestDetail = this.suspendReactivateRequest;

      getAccountRequest = getAccountRequest ? Object.assign(getAccountRequest, newRequest) : newRequest;

      return getAccountRequest;
    }

    return null;
  }

  public async rejectSuspendReactivateAccount(account: ParticipantAccount): Promise<KillSwitchRequestDetail> {

    const suspend = this.getAccountStatus(account.address) === 'suspend_requested' ? true : false;

    // error out if suspend/reactivate request was not found
    if (!this.suspendReactivateRequest || !this.suspendReactivateRequest.approvalIds) {
      return null;
    }

    if (this.suspendReactivateDbRef) {
      const newRequest: KillSwitchRequest = JSON.parse(JSON.stringify({
        participantId: this.sessionService.currentNode.participantId,
        accountAddress: account.address,
        [suspend && 'suspendRejectedBy']: this.authService.auth.auth.currentUser.email,
        [!suspend && 'reactivateRejectedBy']: this.authService.auth.auth.currentUser.email,
        approvalIds: [],
        status: suspend === true ? 'normal' : 'suspended',
      }));

      // capture requestId in meta data
      await this.suspendReactivateDbRef.update(newRequest);

      let getAccountRequest: KillSwitchRequestDetail = this.suspendReactivateRequest;

      getAccountRequest = getAccountRequest ? Object.assign(getAccountRequest, newRequest) : newRequest;

      return getAccountRequest;
    }

    return null;
  }

  /**
   * Gets current status of account based on existing request.
   *
   * @private
   * @param {string} accountAddress
   * @returns {KillSwitchRequestStatus}
   * @memberof ParticipantAccountComponent
   */
  public getAccountStatus(accountAddress: string): KillSwitchRequestStatus {
    if (this.suspendReactivateRequest && this.suspendReactivateRequest.accountAddress === accountAddress) {
      return this.suspendReactivateRequest.status;
    }

    // without a request in firebase, we will assume account is already active
    // since there's no way to get status of an account from the API at the moment
    return 'normal';
  }

  /**
 * Get suspend/reactivate request requestor
 *
 * @private
 * @param {string} accountAddress
 * @returns {string}
 * @memberof ParticipantAccountComponent
 */
  public getRequesterEmail(account: ParticipantAccount): string {
    if (this.suspendReactivateRequest && this.suspendReactivateRequest.accountAddress === account.address) {
      switch (this.getAccountStatus(account.address)) {
        case 'suspend_requested':
          return this.suspendReactivateRequest.suspendRequestedBy;
        case 'reactivate_requested':
          return this.suspendReactivateRequest.reactivateRequestedBy;
      }
    }
  }

  /**
   * Get user-friendly current action (suspend/reactivate) being done on request
   *
   * @param {string} accountAddress
   * @returns
   * @memberof ParticipantAccountService
   */
  public getSuspendReactivate(account: ParticipantAccount): string {
    if (this.suspendReactivateRequest && this.suspendReactivateRequest.accountAddress === account.address) {
      switch (this.getAccountStatus(account.address)) {
        case 'suspend_requested':
          return 'Suspension';
        case 'reactivate_requested':
          return 'Reactivation';
      }
    }

    return null;
  }
}
