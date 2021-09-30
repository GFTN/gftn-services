// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, ComponentRef, OnDestroy } from '@angular/core';
import { ModalService } from 'carbon-components-angular';
import { WhitelistModalComponent } from '../whitelist-modal/whitelist-modal.component';
import { WhitelistRequest, WhitelistRequestStatus } from '../../shared/models/whitelist-request.interface';
import { AngularFireDatabase } from '@angular/fire/database';
import { SessionService } from '../../../shared/services/session.service';
import { UtilsService } from '../../../shared/utils/utils';
import { AuthService } from '../../../shared/services/auth.service';
import { NotificationService } from '../../../shared/services/notification.service';
import { AccountService } from '../../shared/services/account.service';
import { Approval } from '../../../shared/models/approval.interface';
import { IUserProfile } from '../../../shared/models/user.interface';
import { HttpHeaders, HttpClient, HttpResponse } from '@angular/common/http';
import { WorldWireError } from '../../../shared/models/error.interface';
import { Subscription } from 'rxjs';
import { isEmpty } from 'lodash';

@Component({
  selector: 'app-whitelist',
  templateUrl: './whitelist.component.html',
  styleUrls: ['./whitelist.component.scss']
})
export class WhitelistComponent implements OnInit, OnDestroy {

  loaded = false;

  whitelistRequests: WhitelistRequest[];

  dbRef: firebase.database.Reference;

  requestUrl: string;

  loadingRequest = false;

  selectedParticipant: string;

  participantSubscription: Subscription;

  currentOpenModal: ComponentRef<WhitelistModalComponent>;

  constructor(
    private modalService: ModalService,
    public sessionService: SessionService,
    private accountService: AccountService,
    public authService: AuthService,
    public utils: UtilsService,
    private notificationService: NotificationService,
    private http: HttpClient,
    private db: AngularFireDatabase
  ) {
    this.dbRef = this.db.database.ref('whitelist_requests').child(this.sessionService.currentNode.participantId);
  }

  ngOnInit() {

    this.participantSubscription = this.accountService.currentParticipantChanged.subscribe(() => {
      this.requestUrl = `${this.accountService.globalRoot}/whitelist/v1/client/participants/whitelist`;

      this.getWhitelistRequests();
    });
  }

  ngOnDestroy() {
    if (this.currentOpenModal) {
      this.currentOpenModal.instance.closeModal();
    }
    this.participantSubscription.unsubscribe();
  }

  /**
   * Get class for text color based on status code
   *
   * @param {WhitelistRequestStatus} status
   * @returns
   * @memberof WhitelistComponent
   */
  getTextColor(status: WhitelistRequestStatus): string {
    switch (status) {
      case 'approved': return 'text-ok';
      case 'rejected': return 'text-danger';
      default: return 'text-warning';
    }
  }

  /**
   * Get/refresh all whitelist requests
   *
   * @memberof WhitelistComponent
   */
  async getWhitelistRequests(): Promise<void> {

    // reset loaded
    this.loaded = false;

    // Get whitelisted participants from API. This needs to
    // be reloaded every time whitelist is refreshed
    await this.accountService.getWhitelistedParticipants();

    // refresh list of requests
    this.whitelistRequests = null;

    // lazy get data from firebase in a promise
    return new Promise((resolve) => {
      this.dbRef
        .once('value', async (data: firebase.database.DataSnapshot) => {

          let allRequests: { [key: string]: WhitelistRequest } = data.val() ? data.val() : {};

          // process all requests before comparing against official whitelist from API
          if (!isEmpty(allRequests)) {
            allRequests = await this.processAllData(allRequests);
          }

          // sync requests with whitelist through the API, if any
          if (this.accountService.whitelistedParticipants) {
            for (const participant of this.accountService.whitelistedParticipants) {
              // account for already whitelisted participants throug the API
              if (!allRequests[participant]) {
                allRequests[participant] = {
                  whitelisterId: this.sessionService.currentNode.participantId,
                  whitelistedId: participant,
                  requestedBy: 'N/A',
                  approvedBy: 'API',
                  timeUpdated: 0,
                  status: 'approved',
                  approvalIds: null
                };
              } else if (allRequests[participant] && allRequests[participant].approvedBy) {
                continue;
              } else {
                // whitelisted participant was added by API request
                allRequests[participant].status = 'approved';
                allRequests[participant].approvedBy = 'API';
              }
            }
          }

          this.whitelistRequests = allRequests ? Object.values(allRequests) : null;

          this.loaded = true;

          resolve();
        });
    });
  }

  /**
   * Processes all requests
   *
   * @param {{ [key: string]: WhitelistRequest }} allRequests
   * @returns {Promise<{ [key: string]: WhitelistRequest }>}
   * @memberof WhitelistComponent
   */
  processAllData(allRequests: { [key: string]: WhitelistRequest })
    : Promise<{ [key: string]: WhitelistRequest }> {

    return new Promise((resolve) => {

      const length = Object.keys(allRequests).length;

      let i = 0;

      // transform request data
      for (const [key, request] of Object.entries(allRequests)) {
        request.key = key;

        if (request.approvalIds && request.approvalIds.length > 0) {
          this.processRequestData(request)
            .then(() => {

              i++;

              if (i === length) {
                resolve(allRequests);
              }
            });
        } else {
          i++;

          request.status = 'deleted';

          if (i === length) {
            resolve(allRequests);
          }
        }
      }

    });

  }

  processRequestData(request: WhitelistRequest): Promise<void> {

    return new Promise(async (resolve) => {
      request = await this.getApprovalInfo(request, request.approvalIds[0], 'add');

      if (request.approvalIds.length > 1) {
        request = await this.getApprovalInfo(request, request.approvalIds[0], 'delete');
      }

      resolve();
    });
  }

  /**
   * Get approval information
   *
   * @param {WhitelistRequest} request
   * @returns {Promise<WhitelistRequest>}
   * @memberof WhitelistComponent
   */
  getApprovalInfo(request: WhitelistRequest, approvalId: string, action: 'add' | 'delete'): Promise<WhitelistRequest> {
    return new Promise((resolve) => {
      this.db.database.ref('/participant_approvals')
        .child(approvalId)
        .once('value', async (data: firebase.database.DataSnapshot) => {
          const approval: Approval = data.val();

          if (approval) {

            const userRequests: Promise<any>[] = [];

            // get users ref for this institution
            const dbUsersRef = `participants/${this.sessionService.institution.info.institutionId}/users`;

            if (approval.uid_request) {
              userRequests.push(this.db.database.ref(dbUsersRef).child(approval.uid_request).once('value'));
            }

            if (approval.uid_approve) {
              userRequests.push(this.db.database.ref(dbUsersRef).child(approval.uid_approve).once('value'));
            }

            const getUsers = await Promise.all(userRequests);

            const users: IUserProfile[] = [];

            for (const userData of getUsers) {
              if (userData.val()) {
                users.push(userData.val());
              }
            }

            switch (action) {
              case 'add':
                request.requestedBy = users ? users[0].profile.email : null;
                request.approvedBy = users.length > 1 ? users[1].profile.email : null;
                if (request.status === 'rejected') {
                  break;
                }
                request.status = users.length > 1 ? 'approved' : 'pending';
                break;
              case 'delete':
                request.deleteRequestedBy = users ? users[0].profile.email : null;
                request.deleteApprovedBy = users.length > 1 ? users[1].profile.email : null;
                request.status = users.length > 1 ? 'deleted' : 'pending';
            }

            request.timeUpdated = approval.timestamp_request;
          }

          resolve(request);
        });
    });
  }


  async approveWhitelistRequest(request: WhitelistRequest) {
    // approver of whitelist creation must be an admin
    if (!this.authService.userIsParticipantAdmin(
      this.sessionService.institution.info.institutionId
    )) {
      this.notificationService.show(
        'error',
        'You must be an admin of this participant to approve this action',
        null,
        'Unauthorized Action',
        'top'
      );
      return;
    }

    // do secondary check to prevent UI hacks
    if (request.requestedBy !== this.authService.userProfile.profile.email) {

      // Call whitelist endpoint
      try {
        this.loadingRequest = true;

        let h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

        // checker action for whitelist request
        h = this.authService.addMakerCheckerHeaders(h, 'approve', request.approvalIds[0]);

        const options = {
          headers: h,
          responseType: 'text' as 'text'
        };

        const body = {
          participant_id: request.whitelistedId
        };

        const response = await this.http.post(
          this.requestUrl,
          body,
          options
        ).toPromise();


        // Send notification of approval
        this.notificationService.show(
          'success', response
        );

        this.loadingRequest = false;

        // refresh whitelist
        this.getWhitelistRequests();
      } catch (err) {
        this.loadingRequest = false;

        if (err.error) {
          const error: WorldWireError = err.error;

          console.log('err', err);

          this.notificationService.show(
            'error',
            error.details,
            null,
            error.message,
            'top'
          );
        } else {
          this.notificationService.show(
            'error',
            'Unexpected error found when creating this whitelist request.',
            null,
            'Unexpected Error',
            'top'
          );
        }

        // reset approvals
        this.resetApprovals(request);

        // refresh whitelist
        this.getWhitelistRequests();
      }
    } else {
      this.loadingRequest = false;

      this.notificationService.show(
        'error',
        'Approver cannot be the same person as the creator of the request.',
        null,
        'Unauthorized Action',
        'top'
      );
    }
  }

  async rejectWhitelistRequest(request: WhitelistRequest) {
    // approver of whitelist creation must be an admin
    if (!this.authService.userIsParticipantAdmin(
      this.sessionService.institution.info.institutionId
    )) {
      this.notificationService.show(
        'error',
        'You must be an admin of this participant to approve this action',
        null,
        'Unauthorized Action',
        'top'
      );
      return;
    }
    this.loadingRequest = true;
    // do secondary check to prevent UI hacks
    if (request.requestedBy !== this.authService.userProfile.profile.email) {
      // TODO: Send notification of rejection

      // update record in firebase
      await this.dbRef.child(request.key).update({
        status: 'rejected',
        rejectedBy: this.authService.userProfile.profile.email,
        timeUpdated: this.utils.getTimestampInSecs(),
      });

      this.loadingRequest = false;
      // refresh whitelist
      this.getWhitelistRequests();
    } else {
      this.notificationService.show(
        'error',
        'Approver cannot be the same person as the creator of the request.',
        null,
        'Unauthorized Action',
        'top'
      );
    }
  }

  selectParticipant(participantId: string) {
    this.selectedParticipant = participantId;
  }

  /**
   * Initiates deletion of a whitelisted participant
   *
   * @param {WhitelistRequest} request
   * @memberof WhitelistComponent
   */
  async initiateDeleteWhitelistRequest(request: WhitelistRequest) {

    try {
      this.loadingRequest = true;

      let h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

      // maker action for delete whitelist request
      h = this.authService.addMakerCheckerHeaders(h, 'approve', request.approvalIds[0]);

      const options = {
        headers: h
      };

      //  TODO: CANNOT USE BODY IN DELETE. NEEDS ENDPOINT REFACTORED
      // const body = {
      //   participant_id: request.whitelistedId
      // };

      await this.http.delete(
        this.requestUrl,
        options
      );

      this.loadingRequest = false;
    } catch (err) {
      this.loadingRequest = false;

      this.notificationService.show(
        'error',
        'Unexpected error when making delete request.',
        null,
        'Unexpected Error',
        'top'
      );
    }
  }

  /**
   * Approves deletion of a whitelisted participant
   *
   * @param {WhitelistRequest} request
   * @returns
   * @memberof WhitelistComponent
   */
  async approveDeleteWhitelistRequest(request: WhitelistRequest) {
    // approver of whitelist deletion must be an admin
    if (!this.authService.userIsParticipantAdmin(
      this.sessionService.institution.info.institutionId
    )) {
      this.notificationService.show(
        'error',
        'You must be an admin of this participant to approve this action',
        null,
        'Unauthorized Action',
        'top'
      );
      return;
    }

    // do secondary check to prevent UI hacks
    if (request.deleteRequestedBy !== this.authService.userProfile.profile.email) {
      try {
        this.loadingRequest = true;

        let h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

        // checker action for delete whitelist request
        h = this.authService.addMakerCheckerHeaders(h, 'approve', request.approvalIds[0]);

        const options = {
          headers: h
        };

        //  TODO: CANNOT USE BODY IN DELETE. NEEDS ENDPOINT REFACTORED
        // const body = {
        //   participant_id: request.whitelistedId
        // };

        const response: any = await this.http.delete(
          this.requestUrl,
          options
        ).toPromise();

        // Remove delete request from firebase
        if (response.ok) {
          await this.dbRef
            .child(request.whitelisterId)
            .child(request.whitelisterId)
            .remove();
        }

        this.loadingRequest = false;
      } catch (err) {
        this.loadingRequest = false;

        this.notificationService.show(
          'error',
          'Unexpected error when making delete request.',
          null,
          'Unexpected Error',
          'top'
        );

        // reset approvals
        this.resetApprovals(request);
      }
    } else {
      this.notificationService.show(
        'error',
        'Approver cannot be the same person as the creator of the request.',
        null,
        'Unauthorized Action',
        'top'
      );
    }
  }

  /**
   * Reset maker/checker array to capture failed requests.
   * Maker will be forced to re-send their request,in case of a system error.
   * This is a worst-case solution, so that we can keep the approval data
   * in sync with the request notification.
   *
   * @param {WhitelistRequest} request
   * @memberof WhitelistComponent
   */
  resetApprovals(request: WhitelistRequest) {

    // remove last approvalId
    request.approvalIds.pop();

    // record timestamp of the reset
    const timeStamp = Math.floor(Date.now() / 1000);

    // set to null if array is empty for posting properly to firebase
    if (request.approvalIds && request.approvalIds.length === 0) {
      request.approvalIds = null;
    }

    this.dbRef.child(request.key).update({
      time_updated: timeStamp,
      approvalIds: request.approvalIds
    });
  }

  /**
   * Open modal for creating a new whitelist request
   *
   * @memberof WhitelistComponent
   */
  addToWhitelist() {
    // creates and opens the modal for adding a new whitelist request
    this.currentOpenModal = this.modalService.create({
      component: WhitelistModalComponent,
      inputs: {
        MODAL_DATA: {}
      }
    });

    // detect close event
    this.currentOpenModal.instance.close.subscribe(() => {
      // refresh whitelist
      this.getWhitelistRequests();
    });
  }
}
