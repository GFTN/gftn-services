// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, HostBinding, OnDestroy } from '@angular/core';
import { CheckboxOption } from '../../../shared/models/checkbox-option.model';
import { TrustRequest, TrustRequestStatus, TrustRequestPermission } from '../../shared/models/trust-request.interface';
import { AngularFireDatabase } from '@angular/fire/database';
import { SessionService } from '../../../shared/services/session.service';
import { includes, pickBy, filter } from 'lodash';
import { UtilsService } from '../../../shared/utils/utils';
import { AuthService } from '../../../shared/services/auth.service';
import { NotificationService } from '../../../shared/services/notification.service';
import { Approval, ApprovalInfo } from '../../../shared/models/approval.interface';
import { IUserProfile } from '../../../shared/models/user.interface';
import { AccountService } from '../../shared/services/account.service';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { WorldWireError } from '../../../shared/models/error.interface';
import { Subscription } from 'rxjs';
import { ParticipantApprovalModel } from '../../../shared/models/participant-approval.model';


export interface TrustRequestType {
  key: 'outgoing' | 'incoming';
  requestField: string;
  statuses: StatusMapping[];
}

export interface StatusMapping {
  // user-friendly name/label of status
  name: string;
  type: 'info' | 'success' | 'warning' | 'error';
  // actual values of the status recorded in firebase
  value?: TrustRequestStatus[];
}

/**
 * Main handler for managing trustlines
 * of a participant account
 *
 * @export
 * @class TrustlinesComponent
 * @implements {OnInit}
 */
@Component({
  selector: 'app-trustlines',
  templateUrl: './trustlines.component.html',
  styleUrls: ['./trustlines.component.scss'],
  providers: [
    ParticipantApprovalModel
  ]
})

export class TrustlinesComponent implements OnInit, OnDestroy {

  /**
   * Incoming - requests made by other participants
   * to trust this participant's issued asset
   * Outgoing - requests sent out by this participant
   * to trust other participants' issued assets
   *
   * @type {TrustRequestType[]}
   * @memberof TrustlinesComponent
   */
  requestTypes: TrustRequestType[] = [
    {
      key: 'outgoing',
      requestField: 'requestor_id',
      statuses: [{
        name: 'initiated',
        type: 'info',
      }, {
        name: 'pending',
        type: 'warning',
        value: ['requested', 'allowed', 'rejectPending']
      }, {
        name: 'approved',
        type: 'success',
        // revokePending stays in 'approved' state
        // until Issuer admin approves revocation
        value: ['approved', 'revokePending']
      }, {
        name: 'rejected',
        type: 'error'
      },
      {
        name: 'revoked',
        type: 'error'
      }],

    },
    {
      key: 'incoming',
      requestField: 'issuer_id',
      statuses: [{
        name: 'requested',
        type: 'info'
      }, {
        name: 'pending',
        type: 'warning',
        value: ['allowed', 'rejectPending', 'revokePending']
      }, {
        name: 'approved',
        type: 'success'
      }, {
        name: 'rejected',
        type: 'error'
      },
      {
        name: 'revoked',
        type: 'error'
      }],
    }
  ];

  supportedRequestTypes: TrustRequestType[];

  // stores the current request type being viewed
  currRequestType: 'outgoing' | 'incoming';

  // stores the status filter optiosn for the
  // current request type being viewed
  statusFilters: CheckboxOption[];

  // stores date filter options (ascending/descending)
  dateFilters: CheckboxOption[] = [{
    name: 'ascending',
    checked: false
  }, {
    name: 'descending',
    checked: false
  }];

  // stores the current date filter being applied
  currDateFilter: string;

  // controls "View all" toggle for status filters
  statusAll = true;

  // stores all trust requests
  allRequests: TrustRequest[];

  // current list of requests, filtered for the view
  filteredRequests: TrustRequest[];

  dbRef: firebase.database.Reference;

  requestUrl: string;

  loadingRequest = false;

  participantSubscription: Subscription;

  trustRequestsSubscription: Subscription;

  isAnchor = false;

  constructor(
    private sessionService: SessionService,
    private db: AngularFireDatabase,
    private http: HttpClient,
    public utils: UtilsService,
    public authService: AuthService,
    private notificationService: NotificationService,
    private accountService: AccountService,
    private participantApprovals: ParticipantApprovalModel
  ) { }

  @HostBinding('attr.class') cls = 'flex-fill';

  ngOnInit() {

    this.dbRef = this.db.database.ref('trust_requests');

    const nodeIsAnchor = this.sessionService.currentNode ? this.sessionService.currentNode.role === 'IS' : this.isAnchor;

    this.isAnchor = this.accountService.participantDetails ? this.accountService.participantDetails.role === 'IS' : nodeIsAnchor;

    this.participantSubscription = this.accountService.currentParticipantChanged.subscribe(() => {

      this.isAnchor = this.accountService.participantDetails ? this.accountService.participantDetails.role === 'IS' : nodeIsAnchor;

      if (this.isAnchor) {

        // check for successful request
        this.requestUrl = this.accountService.participantDetails ? `${this.accountService.globalRoot}/anchor/v1/anchor/trust/${this.accountService.participantDetails.id}` : null;

        // only show incoming rquests for Anchor.
        // Anchors can only allow trust requests from other participants, not create any themselves
        this.supportedRequestTypes = [this.requestTypes[1]];
      } else {
        this.requestUrl = `${this.accountService.apiRoot}/v1/client/trust`;

        this.supportedRequestTypes = this.requestTypes;
      }

      // default initial view to display outgoing requests
      this.getRequestsByType(this.supportedRequestTypes[0].key);
    });
  }

  ngOnDestroy() {
    this.participantSubscription.unsubscribe();

    if (this.trustRequestsSubscription) {
      this.trustRequestsSubscription.unsubscribe();
    }
  }

  /**
   * Get List of Trust Requests by request type
   *
   * @param {string} requestType
   * @memberof TrustlinesComponent
   */
  async getRequestsByType(requestType: 'outgoing' | 'incoming'): Promise<void> {

    const refreshed = this.currRequestType === requestType;

    // set current request type
    this.currRequestType = requestType;

    const requestTypeObject: TrustRequestType = this.requestTypes.find((type) =>
      type.key === this.currRequestType
    );

    // reset all requests
    this.allRequests = null;

    if (!refreshed) {
      // reset status filters based on new request type
      this.setStatusFilters(requestTypeObject.statuses);

      // reset date filter. default to return by descending date
      this.currDateFilter = this.dateFilters[1].name;

      // reset status filters
      this.statusAll = true;
    }

    // wait for data to come back before filtering
    this.trustRequestsSubscription = this.accountService.getTrustRequests(requestTypeObject.requestField)
      .subscribe(async (allRequests: TrustRequest[]) => {

        this.allRequests = allRequests;

        // transform request data
        await this.processAllData();

        this.filterByStatus();
      });
  }

  processAllData(): Promise<void> {
    return new Promise((resolve) => {

      if (this.allRequests.length > 0) {

        let i = 0;

        for (const request of this.allRequests) {

          // process each request simultaneously in promises
          this.processRequestData(request).then(() => {

            request.loaded = true;

            i++;

            if (i === this.allRequests.length) {

              resolve();
            }
          });
        }
      } else {

        resolve();
      }
    });
  }

  /**
   * Process metadata for an individual trust request
   *
   * @param {string} key
   * @param {TrustRequest} request
   * @returns {Promise<void>}
   * @memberof TrustlinesComponent
   */
  processRequestData(request: TrustRequest): Promise<void> {
    return new Promise(async (resolve) => {

      if (request.approval_ids && request.approval_ids.length > 0) {

        const approvalId1: string = request.approval_ids[0];

        request = await this.getApprovalInfo(approvalId1, 'request', request);


        if (request.approval_ids.length === 2) {
          const approvalId2: string = request.approval_ids.length === 2 ? request.approval_ids[1] : null;

          request = await this.getApprovalInfo(approvalId2, 'allow', request);

        }

        if (request.approval_ids.length === 3) {
          const approvalId3: string = request.approval_ids.length === 3 ? request.approval_ids[2] : null;

          request = await this.getApprovalInfo(approvalId3, 'revoke', request);

        }
        resolve();
      } else {
        // No approvals or current status found. Mark request as 'rejected' as final state in the view
        if (request.status !== 'rejected') {
          request.status = 'rejected';
          request.reason_rejected = 'Request failed due to system error.';
        }

        resolve();
      }
    });
  }

  /**
   * Get statuses and users based on approval information
   *
   * @param {string} approvalId
   * @param {TrustRequestPermission} permission
   * @param {TrustRequest} request
   * @returns {Promise<TrustRequest>}
   * @memberof TrustlinesComponent
   */
  getApprovalInfo(approvalId: string, permission: TrustRequestPermission, request: TrustRequest): Promise<TrustRequest> {
    return new Promise((resolve) => {
      this.db.database.ref('/participant_approvals')
        .child(approvalId)
        .once('value', async (data: firebase.database.DataSnapshot) => {
          const approval: Approval = data.val();

          // get users ref for this institution
          const dbUsersRef = 'users';

          if (approval) {
            const userRequests: Promise<any>[] = [];

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

            switch (permission) {
              case 'request':
                request.requestInitiatedBy = users ? users[0].profile.email : null;
                request.requestApprovedBy = users.length > 1 ? users[1].profile.email : null;

                if (request.status === 'rejected') {
                  break;
                }
                request.status = users.length > 1 && request ? 'requested' : 'initiated';
                break;
              case 'allow':
                request.allowInitiatedBy = users ? users[0].profile.email : null;
                request.allowApprovedBy = users.length > 1 ? users[1].profile.email : null;

                if (request.status === 'rejected') {
                  break;
                }
                request.status = users.length > 1 ? 'approved' : 'allowed';
                break;
              case 'revoke':
                request.revokeInitiatedBy = users ? users[0].profile.email : null;
                request.revokeApprovedBy = users.length > 1 ? users[1].profile.email : null;

                if (request.status === 'rejected') {
                  break;
                }

                request.status = users.length > 1 ? 'revoked' : 'revokePending';
                break;
            }

            // get latest timestamp, either from approval record or notification
            request.time_updated = (request.time_updated < approval.timestamp_request) ? approval.timestamp_request : request.time_updated;
          }

          resolve(request);
        });
    });
  }

  /**
   * Get user-friendly title/label for the status code
   *
   * @param {string} statusInput
   * @returns {string}
   * @memberof TrustlinesComponent
   */
  getStatusTitle(statusInput: string): string {
    const statusList = this.getStatusList();

    const statusObj: StatusMapping = statusList.find((status: StatusMapping) => {

      const checkValue: string[] = status.value ? status.value : [status.name];

      for (const value of checkValue) {
        if (value === statusInput) {
          return true;
        }
      }
      return false;
    });

    return statusObj ? statusObj.name : '';
  }

  /**
   * Get notification type for the requested status
   *
   * @param {string} statusInput
   * @returns {string}
   * @memberof TrustlinesComponent
   */
  getStatusType(statusInput: TrustRequestStatus): string {

    const statusList = this.getStatusList();

    const statusObj: StatusMapping = statusList.find((status: StatusMapping) => {

      const checkValue: string[] = status.value ? status.value : [status.name];

      for (const value of checkValue) {
        if (value === statusInput) {
          return true;
        }
      }
      return false;
    });

    return statusObj ? statusObj.type : 'info';
  }

  /**
   * Get list of statuses based on the current request type
   */
  getStatusList(): StatusMapping[] {
    const requestType: TrustRequestType = this.requestTypes.find((type) =>
      type.key === this.currRequestType
    );

    return requestType ? requestType.statuses : [];
  }

  /**
   * Initializes and sets list of status filters
   * based on the request type being viewed
   *
   * OPTIONAL: @param {TrustRequestStatus[]} [statusList]
   * @memberof TrustlinesComponent
   */
  setStatusFilters(statusList?: StatusMapping[]) {
    this.statusFilters = [];

    statusList = statusList ? statusList : this.getStatusList();

    for (const status of statusList) {

      const value = status.value ? status.value.join(' ') : status.name;
      // init to all statuses
      this.statusFilters.push({
        label: status.name,
        name: value,
        checked: true
      });
    }
  }

  /**
   * Filter results by descending/ascending date
   *
   * @memberof TrustlinesComponent
   */
  sortByDate() {
    this.filteredRequests.reverse();
  }

  /**
   * Filter results by status
   *
   * @memberof TrustlinesComponent
   */
  filterByStatus() {
    const selectedFilters: string[] = [];

    for (const status of this.statusFilters) {
      // get list of selected filters
      if (status.checked) {
        const vals: string[] = status.name.split(' ');
        for (const val of vals) {
          selectedFilters.push(val);
        }
      }
    }

    // get requests based on selected status filters
    this.filteredRequests = filter(this.allRequests, (request: TrustRequest) => {
      return includes(selectedFilters, request.status);
    });

    // reset sort. default from newest to oldest (descending)
    if (this.currDateFilter === 'descending') {
      this.sortByDate();
    }
  }

  /**
   * toggle view all on/off and update
   * the individual checkboxes accordingly
   *
   * @memberof TrustlinesComponent
   */
  selectAllStatuses() {
    for (const status of this.statusFilters) {
      status.checked = this.statusAll;
    }

    this.filterByStatus();
  }


  /**
   * Checker approves the created request.
   * Sets the status to 'requested' by sending
   * 'REQUEST' action to the API
   *
   * @param {TrustRequest} request
   * @memberof TrustlinesComponent
   */
  async approveCreatedTrustline(request: TrustRequest) {

    // approver of trustline creation must be an admin
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
    if (request.requestInitiatedBy !== this.authService.userProfile.profile.email) {

      // Call trust endpoint to request trust. Send notification of approval
      let h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

      h = this.authService.addMakerCheckerHeaders(h, 'approve', request.approval_ids[0]);

      const options = {
        headers: h
      };

      const body = {
        'permission': 'request',
        'asset_code': request.asset_code,
        'account_name': request.account_name,
        'participant_id': request.issuer_id,
        'limit': request.limit,
        'end_to_end_id': request.key
      };

      try {

        this.loadingRequest = true;

        await this.http.post(
          this.requestUrl,
          body,
          options
        ).toPromise();

        this.loadingRequest = false;
      } catch (err) {

        // reset approvals
        await this.resetApprovals(request.approval_ids[0]);

        this.loadingRequest = false;

        if (err.error) {
          const error: WorldWireError = err.error;

          let message = error.message ? error.message : error.msg;
          message = message ? message : 'System Error';

          this.notificationService.show(
            'error',
            error.details,
            null,
            message,
            'top'
          );
        } else {
          this.notificationService.show(
            'error',
            'Approval of trust request could not be completed due to system error.',
            null,
            'System Error',
            'top'
          );
        }
      }

      // refresh transactions
      this.getRequestsByType(this.currRequestType);
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
   * Perform initial allowance of trust (maker action)
   * Sets status of Trust Request to 'allow' as an intermediary step
   * before hitting the actual endpoint
   */
  async allowTrustline(request: TrustRequest) {

    // allower of trustline should have manager or admin permissions
    if (!this.authService.userIsParticipantManagerOrHigher(
      this.sessionService.institution.info.institutionId
    )) {
      this.notificationService.show(
        'error',
        'You have insufficient permissions to perform this action.',
        null,
        'Unauthorized Action',
        'top'
      );
      return;
    }

    // Call endpoint to allow trust - maker request
    try {
      this.loadingRequest = true;

      let h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

      h = this.authService.addMakerCheckerHeaders(h, 'request');

      const options = {
        headers: h
      };

      const body = {
        'permission': 'allow',
        'asset_code': request.asset_code,
        'account_name': request.account_name,
        'participant_id': request.requestor_id,
        'end_to_end_id': request.key
      };

      const response: WorldWireError = await this.http.post(
        this.requestUrl,
        body,
        options
      ).toPromise() as WorldWireError;

      // check to make sure response is correctly formatted
      if (response && response.msg) {

        // check to make sure approvalId actually exists
        const approvalIdExists = await this.db.database.ref('/participant_approvals').child(response.msg).once('value');

        if (approvalIdExists) {
          const newApprovalId: string = response.msg;

          request.approval_ids.push(newApprovalId);

          // update list of approvals
          await this.dbRef.child(request.key).update({
            'approval_ids': request.approval_ids
          });
        }
      }

      this.loadingRequest = false;
    } catch (err) {

      this.loadingRequest = false;

      if (err.error) {
        const error: WorldWireError = err.error;

        let message = error.message ? error.message : error.msg;
        message = message ? message : 'System Error';

        this.notificationService.show(
          'error',
          error.details,
          null,
          message,
          'top'
        );
      } else {
        this.notificationService.show(
          'error',
          'Approval of trust request could not be completed due to system error.',
          null,
          'System Error',
          'top'
        );
      }
    }

    // TODO: Send notification to admins to approve/reject this Allowal

    // refresh transactions
    this.getRequestsByType(this.currRequestType);
  }

  /**
   * Checker approves the allowed request.
   * Sets the status to 'approved' by sending
   * 'ALLOW' action to the API
   *
   *
   * @param {TrustRequest} request
   * @memberof TrustlinesComponent
   */
  async approveAllowedTrustline(request: TrustRequest) {

    // allower of trustline should be an admin
    if (!this.authService.userIsParticipantAdmin(
      this.sessionService.institution.info.institutionId
    )) {
      this.notificationService.show(
        'error',
        'You must be an admin to perform this action.',
        null,
        'Unauthorized Action',
        'top'
      );
      return;
    }

    // checker cannot be the same as the maker of the 'allow' action
    if (request.allowInitiatedBy !== this.authService.userProfile.profile.email) {
      // Call trust endpoint to approval allowal of the trust request

      try {

        this.loadingRequest = true;

        let h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

        h = this.authService.addMakerCheckerHeaders(h, 'approve', request.approval_ids[1]);

        const options = {
          headers: h
        };

        const body = {
          'permission': 'allow',
          'asset_code': request.asset_code,
          'account_name': request.account_name,
          'participant_id': request.requestor_id,
          'end_to_end_id': request.key
        };

        await this.http.post(
          this.requestUrl,
          body,
          options
        ).toPromise();

        this.loadingRequest = false;

      } catch (err) {

        // reset approvals
        await this.resetApprovals(request.approval_ids[1]);

        this.loadingRequest = false;

        if (err.error) {
          const error: WorldWireError = err.error;

          let message = error.message ? error.message : error.msg;
          message = message ? message : 'System Error';

          this.notificationService.show(
            'error',
            error.details,
            null,
            message,
            'top'
          );
        } else {
          this.notificationService.show(
            'error',
            'Approval of trust request could not be completed due to system error.',
            null,
            'System Error',
            'top'
          );
        }

      }

      // TODO: Send notification of approval

      // refresh transactions
      this.getRequestsByType(this.currRequestType);
    } else {
      this.notificationService.show(
        'error',
        'Approval cannot be performed by the same person as the initial allower of the request.',
        null,
        'Unauthorized Action',
        'top'
      );
    }
  }

  /**
   * Marks Trust Request as 'rejected' to ignore the request.
   * No request is made to the API.
   * TODO: Send notification to the participant about rejection.
   * TODO: This can be undone by an admin or by a new trust request.
   *
   * @param {TrustRequest} request
   * @param {string} [checker]
   * @param {string} [reasonRejected]
   * @memberof TrustlinesComponent
   */
  async rejectTrustline(request: TrustRequest, reasonRejected?: string, bypassChecker?: boolean) {

    this.loadingRequest = true;

    const latestApprovalId = request.approval_ids[request.approval_ids.length - 1];

    const latestApproval: ApprovalInfo = await this.participantApprovals.getApprovalInfo(latestApprovalId);

    const checker = latestApproval.requestApprovedBy;

    // Rejector of trustline must be checker and
    // checker cannot be the same as the maker of the action
    if (!checker || checker !== this.authService.userProfile.profile.email) {

      const newApprovalIds: string[] = request.approval_ids.filter((id: string) => {
        return id !== latestApprovalId;
      });

      if (!bypassChecker) {
        // remove approval to reset state to previous
        await this.dbRef.child(request.key)
          .update({
            approval_ids: newApprovalIds
          });
      }

      // final state of rejection
      if (!newApprovalIds || (newApprovalIds && newApprovalIds.length === 0) || !request.approval_ids) {
        const updates = {
          status: 'rejected',
          rejectApprovedBy: this.authService.userProfile.profile.email,
        };

        if (reasonRejected) {
          updates['reason_rejected'] = reasonRejected;
        }
        await this.dbRef.child(request.key).update(updates);
      }

      this.loadingRequest = false;

      // refresh transactions
      this.getRequestsByType(this.currRequestType);
    } else {

      this.loadingRequest = false;

      this.notificationService.show(
        'error',
        'Rejection cannot be performed by the same person as the creator of the request.',
        null,
        'Unauthorized Action',
        'top'
      );
    }
  }

  /**
   * Perform initial revocation of trust (maker action)
   * Sets status of Trust Request to 'revokePending' as an intermediary step
   * before hitting the actual endpoint
   * TODO: Send notification to the participant about rejection.
   *
   * @param {TrustRequest} request
   * @returns
   * @memberof TrustlinesComponent
   */
  async revokeTrustline(request: TrustRequest) {

    // revoker of trustline should have manager or admin permissions
    if (!this.authService.userIsParticipantManagerOrHigher(
      this.sessionService.institution.info.institutionId
    )) {
      this.notificationService.show(
        'error',
        'You have insufficient permissions to perform this action.',
        null,
        'Unauthorized Action',
        'top'
      );
      return;
    }

    try {
      this.loadingRequest = true;

      let h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

      h = this.authService.addMakerCheckerHeaders(h, 'request');

      const options = {
        headers: h
      };

      const body = {
        'permission': 'revoke',
        'asset_code': request.asset_code,
        'account_name': request.account_name,
        'participant_id': request.requestor_id,
        'end_to_end_id': request.key
      };

      const response: WorldWireError = await this.http.post(
        this.requestUrl,
        body,
        options
      ).toPromise() as WorldWireError;

      // check to make sure response is correctly formatted
      if (response && response.msg) {

        // check to make sure approvalId actually exists
        const approvalIdExists = await this.db.database.ref('/participant_approvals').child(response.msg).once('value');

        if (approvalIdExists) {
          const newApprovalId: string = response.msg;

          request.approval_ids.push(newApprovalId);

          // update list of approvals
          await this.dbRef.child(request.key).update({
            'approval_ids': request.approval_ids
          });
        }
      }

      this.loadingRequest = false;

      // TODO: Send notification to admins to approve/reject this revocation

    } catch (err) {

      this.loadingRequest = false;

      if (err.error) {
        const error: WorldWireError = err.error;

        let message = error.message ? error.message : error.msg;
        message = message ? message : 'System Error';

        this.notificationService.show(
          'error',
          error.details,
          null,
          message,
          'top'
        );
      } else {
        this.notificationService.show(
          'error',
          'Revokal of trust request could not be completed due to system error.',
          null,
          'System Error',
          'top'
        );
      }
    }

    // refresh transactions
    this.getRequestsByType(this.currRequestType);
  }

  /**
   * Checker approves the revoked request.
   * Sets the status to 'revoked' by sending
   * 'REVOKE' action to the API
   *
   *
   * @param {TrustRequest} request
   * @memberof TrustlinesComponent
   */
  async approveRevokedTrustline(request: TrustRequest) {

    // revoker of trustline should be an admin
    if (!this.authService.userIsParticipantAdmin(
      this.sessionService.institution.info.institutionId
    )) {
      this.notificationService.show(
        'error',
        'You must be an admin to perform this action.',
        null,
        'Unauthorized Action',
        'top'
      );
      return;
    }

    // checker cannot be the same as the maker of the 'revoke' action
    if (request.revokeInitiatedBy !== this.authService.userProfile.profile.email) {

      // Call trust endpoint
      try {
        this.loadingRequest = true;

        let h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

        h = this.authService.addMakerCheckerHeaders(h, 'approve', request.approval_ids[2]);

        const options = {
          headers: h
        };

        const body = {
          'permission': 'revoke',
          'asset_code': request.asset_code,
          'account_name': request.account_name,
          'participant_id': request.requestor_id,
          'end_to_end_id': request.key
        };

        await this.http.post(
          this.requestUrl,
          body,
          options
        ).toPromise();

        this.loadingRequest = false;

      } catch (err) {

        // reset approvals
        await this.resetApprovals(request.approval_ids[2]);

        this.loadingRequest = false;

        if (err.error) {
          const error: WorldWireError = err.error;

          let message = error.message ? error.message : error.msg;
          message = message ? message : 'System Error';

          this.notificationService.show(
            'error',
            error.details,
            null,
            message,
            'top'
          );
        } else {
          this.notificationService.show(
            'error',
            'Revokal of trust request could not be completed due to system error.',
            null,
            'System Error',
            'top'
          );
        }

      }

      // TODO: Send notification of revocation

      // refresh transactions
      this.getRequestsByType(this.currRequestType);
    } else {
      this.notificationService.show(
        'error',
        'Approval cannot be performed by the same person as the initial revoker of the request.',
        null,
        'Unauthorized Action',
        'top'
      );
    }
  }

  /**
   * Reset checker to capture failed requests upon approval
   *
   * @param {string} approvalId
   * @param {boolean} [newApproval]
   * @memberof TrustlinesComponent
   */
  async resetApprovals(approvalId: string, newApproval?: boolean): Promise<any> {

    const updateFields = {
      status: 'request',
      uid_approve: '',
    };

    // reset approval Id
    return await this.db.database.ref('participant_approvals')
      .child(approvalId)
      .update(updateFields);
  }
}
