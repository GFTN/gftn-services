// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Inject } from '@angular/core';
import { BaseModal, ListItem, NotificationService } from 'carbon-components-angular';
import { AssetType, AssetRequest, Asset } from '../../../../shared/models/asset.interface';
import { AccountService } from '../../services/account.service';
import { HttpClient, HttpHeaders, HttpErrorResponse } from '@angular/common/http';
import { AuthService } from '../../../../shared/services/auth.service';
import { SessionService } from '../../../../shared/services/session.service';
import { WorldWireError } from '../../../../shared/models/error.interface';
import { AngularFireDatabase } from '@angular/fire/database';
import { ApprovalInfo } from '../../../../shared/models/approval.interface';
import { ParticipantApprovalModel } from '../../../../shared/models/participant-approval.model';
import * as cc from 'currency-codes';
import { CUSTOM_REGEXES, CustomRegex } from '../../../../shared/constants/regex.constants';
import { NgForm } from '@angular/forms';

export interface AssetForm {
  assetCode: string;
  assetType: AssetType;
}

@Component({
  selector: 'app-asset-modal',
  templateUrl: './asset-modal.component.html',
  styleUrls: ['./asset-modal.component.scss'],
  providers: [
    ParticipantApprovalModel
  ]
})
export class AssetModalComponent extends BaseModal implements OnInit {

  // current asset type available for view
  currentAssetType: AssetType = 'DO';

  // current asset request available for approve/reject
  currentAssetRequest: AssetRequest;

  latestApproval: ApprovalInfo;

  // main store for asset types information
  assetTypes: { [key: string]: ListItem } = {
    'DO': {
      label: 'Digital Obligation',
      content: 'Digital Obligation (DO)',
      value: 'DO',
      selected: false
    },
    'DA': {
      label: 'Asset',
      content: 'Digital Asset (DA)',
      value: 'DA',
      selected: false
    },
  };

  // stores list of asset types for dropdown
  assetTypeOptions: ListItem[];

  formData: AssetForm;

  // id reference to the notification element in the UI to display errors
  notificationRef = '#assets-error-notication';

  assetDataLoaded = false;

  loadingRequest = false;

  registered = false;

  regexes: { [key: string]: CustomRegex };

  constructor(
    @Inject('MODAL_DATA') public data: any,
    private http: HttpClient,
    private accountService: AccountService,
    private sessionService: SessionService,
    private authService: AuthService,
    private notificationService: NotificationService,
    private participantApprovals: ParticipantApprovalModel,
    private db: AngularFireDatabase,
  ) {
    // MUST do super() to extend BaseModal
    super();

    this.regexes = CUSTOM_REGEXES;

    // OPTIONAL: get asset type if a specific option is provided. Defaults to 'DO' otherwise
    this.currentAssetType = data.assetType || this.currentAssetType;

    // get current asset request if provided. Can be nil for a new request form
    this.currentAssetRequest = data.assetRequest || this.currentAssetRequest;

  }

  async ngOnInit() {

    this.assetTypeOptions = Object.values(this.assetTypes);

    // init form variables
    const initAssetCode = this.currentAssetRequest ? this.currentAssetRequest.asset_code : '';
    const initAssetType = this.currentAssetRequest ? this.currentAssetRequest.asset_type : this.assetTypes[this.currentAssetType].value;

    // init form data
    this.formData = {
      assetCode: initAssetCode,
      assetType: initAssetType
    };

    // check to make sure asset details were passed into the modal
    if (this.currentAssetRequest) {

      // get approval if it exists (maker/checker process done through the portal)
      if (this.currentAssetRequest.approvalIds) {
        const latestApprovalId = this.currentAssetRequest.approvalIds[this.currentAssetRequest.approvalIds.length - 1];
        this.latestApproval = await this.participantApprovals.getApprovalInfo(latestApprovalId);
      }

      // check if asset is already registered
      const foundAsset: Asset = this.accountService.issuedAssets && this.accountService.issuedAssets.find((issuedAsset: Asset) => {
        return issuedAsset.asset_code === this.currentAssetRequest.asset_code;
      });

      // asset is registered if already approved through portal/added via API
      if (foundAsset || (this.latestApproval && this.latestApproval.requestApprovedBy)) {
        this.registered = true;
      }
    }

    this.assetDataLoaded = true;

    // TODO: add validation against already issued assets for asset code
  }

  /**
   * Get custom regex with pattern and
   * validation text based on asset type
   *
   * @returns {CustomRegex}
   * @memberof AssetModalComponent
   */
  getAssetRegex(): CustomRegex {
    return (this.formData.assetType === 'DO') ? this.regexes.assetDO : this.regexes.assetDA;
  }

  /**
   * Get full name of asset based on asset code
   *
   * @returns {string}
   * @memberof AssetModalComponent
   */
  getAssetCodeName(): string {
    const code = (this.formData.assetCode && this.formData.assetCode.length > 2) ? this.formData.assetCode.substring(0, 3) : '';

    return cc.code(code) ? cc.code(code).currency : '';
  }

  /**
   * Makes request to API to request asset issuance request
   *
   * @memberof AssetModalComponent
   */
  public async requestAssetRequest(assetForm: NgForm) {

    if (assetForm.valid) {

      this.loadingRequest = true;

      // force asset code to upper case for submission
      const assetCode = this.formData.assetCode.toUpperCase();

      try {
        const requestUrl = `${this.accountService.apiRoot}/v1/client/assets?asset_code=${assetCode}&asset_type=${this.formData.assetType}`;

        let h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

        h = this.authService.addMakerCheckerHeaders(h, 'request');

        const options = {
          headers: h
        };

        const response: WorldWireError = await this.http.post(requestUrl, null, options).toPromise() as WorldWireError;

        const newRequest: AssetRequest = {
          asset_code: assetCode,
          asset_type: this.formData.assetType,
          approvalIds: [response.msg],
        };

        this.db.database.ref('asset_requests')
          .child(this.sessionService.currentNode.participantId)
          .child(assetCode)
          .update(newRequest);

        this.loadingRequest = false;

        this.closeModal();
      } catch (err) {

        this.throwError(err);
      }
    }
  }

  /**
   * Makes request to API to approve asset issuance request
   *
   * @returns
   * @memberof AssetModalComponent
   */
  public async approveAssetRequest() {

    // Secondary Check: maker CANNOT be same as checker. throw 401 unauthorized error.
    if (this.latestApproval && this.latestApproval.requestInitiatedBy === this.authService.userProfile.profile.email) {
      const error: HttpErrorResponse = new HttpErrorResponse({
        status: 401,
      });
      this.throwError(error, 'Approver of request cannot be the same as requestor.');
      return;
    }

    this.loadingRequest = true;

    try {
      const requestUrl = `${this.accountService.apiRoot}/v1/client/assets?asset_code=${this.currentAssetRequest.asset_code}&asset_type=${this.currentAssetRequest.asset_type}`;

      let h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

      const approveId = this.latestApproval.key;

      h = this.authService.addMakerCheckerHeaders(h, 'approve', approveId);

      const options = { headers: h };

      await this.http.post(requestUrl, null, options).toPromise();

      this.loadingRequest = false;

      this.closeModal();
    } catch (err) {

      // error occurred. reset approval process
      if (this.latestApproval) {
        await this.participantApprovals.resetApprovals(this.latestApproval.key);
      }

      this.throwError(err);
    }
  }

  /**
   * Rejects asset issuance request
   *
   * @returns
   * @memberof AssetModalComponent
   */
  public async rejectAssetRequest() {

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
      const newApprovalIds: string[] = this.currentAssetRequest.approvalIds.filter((id: string) => {
        return id !== this.latestApproval.key;
      });

      // remove approval to reset request
      this.db.database.ref('asset_requests')
        .child(this.sessionService.currentNode.participantId)
        .child(this.currentAssetRequest.asset_code)
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
   * Generic error handler for asset issance requests
   *
   * @private
   * @param {HttpErrorResponse} err
   * @param {string} [optionalErrMsg]
   * @memberof AssetModalComponent
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
        // for 503s and other errors: Unexpected ERROR notification
        this.notificationService.showNotification({
          type: 'error',
          title: 'Unexpected Error',
          message: 'Unexpected error found when making this request. Please contact administrator.',
          target: this.notificationRef
        });
        break;
      }
    }
  }
}
