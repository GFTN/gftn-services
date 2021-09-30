// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, Inject, isDevMode, OnInit } from '@angular/core';
import { BaseModal, ListItem, NotificationService } from 'carbon-components-angular';
import { TrustRequest } from '../../shared/models/trust-request.interface';
import { WorldWireError } from '../../../shared/models/error.interface';
import { AccountService } from '../../shared/services/account.service';
import { SessionService } from '../../../shared/services/session.service';
import { CheckboxOption } from '../../../shared/models/checkbox-option.model';
import { AngularFireDatabase } from '@angular/fire/database';
import { AuthService } from '../../../shared/services/auth.service';
import { HttpHeaders, HttpClient, HttpErrorResponse } from '@angular/common/http';
import { NgForm } from '@angular/forms';

export interface TrustRequestForm {
  issuer_id: string;
  assets: AssetOption[];
}

export interface AssetOption extends CheckboxOption {
  limit: number;
}

@Component({
  selector: 'app-trustline-modal',
  templateUrl: './trustline-modal.component.html',
  styleUrls: ['./trustline-modal.component.scss']
})

export class TrustlineModalComponent extends BaseModal implements OnInit {

  // controls all the data in the trust request form
  formData: TrustRequestForm;

  // list of all issuers (participants who have issued assets) on the network
  issuerOptions: ListItem[];

  // toggles whether or not all checkboxes are selected
  checkAll = false;

  // toggles whether or not there is a
  // mixture of selected/non-selected checkboxes
  checkMixed = false;

  // loader - indicates whether or not the request for getting a participant's assets has finished
  assetsLoading = true;

  dbRef: firebase.database.Reference;

  // id reference to the notification element in the UI to display errors
  notificationRef = '#assets-error-notication';

  loadingRequest = false;

  constructor(
    @Inject('MODAL_DATA') public data: any,
    private db: AngularFireDatabase,
    public accountService: AccountService,
    private sessionService: SessionService,
    private notificationService: NotificationService,
    private authService: AuthService,
    private http: HttpClient
  ) {

    // MUST do super() to extend BaseModal
    super();
  }

  async ngOnInit() {

    this.dbRef = this.db.database.ref('trust_requests');

    // initializing form data
    this.formData = {
      issuer_id: '',
      assets: []
    };

    await this.getIssuerOptions();
  }

  async getIssuerOptions(): Promise<void> {
    // get list of whitelisted participants that
    // this participant can trust assets from
    if (!this.accountService.whitelistedParticipants) {

      // necessary for issuer dropdown
      try {
        await this.accountService.getWhitelistedParticipants();
      } catch (err) {
        console.log(err);
      }
    }

    // initialize list of issuer options
    this.issuerOptions = [];

    // push all whitelisted participants to list of issuers for dropdown
    if (this.accountService.whitelistedParticipants) {
      for (const issuerId of this.accountService.whitelistedParticipants) {

        // exclude current participant
        if (issuerId !== this.sessionService.currentNode.participantId) {
          this.issuerOptions.push({
            content: issuerId,
            value: issuerId,
            selected: false
          });
        }
      }
    }


    // DEFAULT: init selection to first issuer option
    this.formData.issuer_id = this.issuerOptions[0] ? this.issuerOptions[0].value : null;

    if (this.issuerOptions.length === 0) {
      this.getAssetsForIssuer();
    }
  }

  /**
   * Gets list of selectable assets
   * based on the selected Issuer
   *
   * @memberof TrustlineModalComponent
   */
  async getAssetsForIssuer(): Promise<void> {

    // reset list of assets
    this.formData.assets = [];
    this.checkAll = false;
    this.checkMixed = false;
    this.assetsLoading = true;

    // only get assets if issuer is selected
    if (this.formData.issuer_id) {

      // filter DOs for operating accounts
      const filterDO: boolean = this.accountService.accountSlug === 'issuing' ? false : true;

      // Get list of issued assets by participant
      const issuerAssets = await this.accountService.getAssetsForParticipant(this.formData.issuer_id, 'issued', filterDO);

      for (const asset of issuerAssets) {
        this.formData.assets.push({
          name: asset.asset_code,
          label: asset.asset_code,
          limit: null,
          checked: false
        });
      }
    }

    // data request is finished
    this.assetsLoading = false;
  }


  /**
   * Checks if all asset options are checked
   * or if there are a mixture of checked/unchecked options
   *
   * @returns
   * @memberof TrustlineModalComponent
   */
  checkIfMixed() {
    this.checkMixed = false;

    // finds the first unchecked option in the list
    const unchecked = this.formData.assets.find((asset: AssetOption) => !asset.checked);


    // Set 'indeterminate'/mixed state for "Select All" checkbox
    // because of the existence of at least one unchecked option
    if (unchecked && this.checkAll) {
      this.checkMixed = true;
    }

    // Set 'checked' state for "Select All" checkbox
    // since no unchecked options were found
    if (!unchecked) {
      this.checkAll = true;
    }
  }

  /**
   * Selects all available asset options
   *
   * @memberof TrustlineModalComponent
   */
  selectAll() {

    // this.checkAll is 2-way data bound (ngModel)
    // so we don't have to update this boolean manually.
    // This updates the 'checked' value for all asset options
    // accordingly to match this.checkAll
    for (const asset of this.formData.assets) {
      asset.checked = this.checkAll;
    }
  }

  /**
   * Checks if form options and data are valid for submission
   *
   * @param {NgForm} trustlineForm
   * @returns {boolean}
   * @memberof TrustlineModalComponent
   */
  validForm(form: NgForm): boolean {
    return (form.valid && this.issuerOptions && this.issuerOptions.length > 0 && this.formData.assets && this.formData.assets.length > 0);
  }

  /**
   * Aggregates all selected assets from the form
   * and creates a separate Trustline request for
   * each requested asset
   *
   * @memberof TrustlineModalComponent
   */
  async submitForm(form: NgForm) {

    // Form is not valid. Show user error.
    if (!this.validForm(form)) {
      this.notificationService.showNotification({
        type: 'error',
        title: 'Invalid Form',
        message: 'Form is not valid for submission. You may have missing or incorrect inputs.',
        target: this.notificationRef
      });
      return;
    }

    const selectedAssets: AssetOption[] = this.formData.assets.filter((asset: AssetOption) => asset.checked);

    // no assets were selected
    if (selectedAssets.length === 0) {

      // show error notification
      this.notificationService.showNotification({
        type: 'error',
        title: 'No Assets Selected',
        message: 'At least one asset must be selected to create a trustline request.',
        target: this.notificationRef
      });
      return;
    }

    if (!this.authService.userProfile) {
      this.notificationService.showNotification({
        type: 'error',
        title: 'No User Found',
        message: 'Could not get details of user creating this request.',
        target: this.notificationRef
      });
      return;
    }

    this.loadingRequest = true;

    const trustRequests: { [key: string]: TrustRequest } = {};

    const trustRequest = `${this.accountService.apiRoot}/v1/client/trust`;

    const requests: TrustRequest[] = [];

    let h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

    h = this.authService.addMakerCheckerHeaders(h, 'request');

    const options = {
      headers: h
    };

    // create a trust request for each selected asset
    for (const asset of selectedAssets) {

      // generate unix timestamp for the request
      const timeStamp = Math.floor(Date.now() / 1000);

      const body = {
        'permission': 'request',
        'asset_code': asset.name,
        'account_name': this.accountService.accountSlug,
        'participant_id': this.formData.issuer_id,
        'limit': asset.limit
      };

      try {
        const response: WorldWireError = await this.http.post(
          trustRequest,
          body,
          options
        ).toPromise() as WorldWireError;

        // check to make sure response is correctly formatted
        if (response && response.msg) {

          // check to make sure approvalId actually exists
          const approvalIdExists = await this.db.database.ref('/participant_approvals').child(response.msg).once('value');

          if (approvalIdExists) {
            requests.push({
              requestor_id: this.sessionService.currentNode.participantId,
              issuer_id: this.formData.issuer_id,
              account_name: this.accountService.accountSlug,
              asset_code: asset.name,
              limit: asset.limit,
              time_updated: timeStamp,
              status: 'initiated',
              approval_ids: [response.msg]
            });
          }
        }
      } catch (err) {

        let error: HttpErrorResponse = err;

        if (!error.error) {
          error = new HttpErrorResponse({
            status: 500,
          });
        }

        this.throwError(err);

        return;
      }
    }

    // only handle successful requests
    for (const request of requests) {
      // push up trust request data for a newly created key
      trustRequests[this.dbRef.push().key] = request;
    }


    if (isDevMode) {
      console.log('formData', this.formData);
      console.log('trustRequests', trustRequests);
      console.log('requests', requests);
    }

    // Post trust requests to firebase
    this.dbRef.update(trustRequests).then(() => {

      this.loadingRequest = false;

      // close modal
      this.closeModal();
    }).catch((err: HttpErrorResponse) => {
      console.log(err);

      // show error message if unexpected error
      this.throwError(err);

    });

  }

  /**
   * Throw error message to the screen for user
   *
   * @param {HttpErrorResponse} error
   * @memberof TrustlineModalComponent
   */
  public throwError(error: HttpErrorResponse) {

    this.loadingRequest = false;

    const wwError: WorldWireError = error.error ? error.error : null;

    switch (error.status) {
      case 401:
      case 0: {
        // Unauthorized ERROR notification
        this.notificationService.showNotification({
          type: 'error',
          title: 'Unauthorized',
          message: 'You are not authorized to make a trustline request. Please contact administrator.',
          target: '#notification'
        });
        break;
      }
      case 400: {
        const message = wwError && wwError.details ? wwError.details : 'Invalid form data was submitted and could not be processed. Please try again.';
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
        const message = wwError && wwError.details ? wwError.details : 'Network could not be reached to make this request. Please contact administrator.';

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
        // for 404s and other errors: Unexpected ERROR notification
        this.notificationService.showNotification({
          type: 'error',
          title: 'Unexpected Error',
          message: 'Unexpected error found when creating this trust request. Please contact administrator.',
          target: this.notificationRef
        });
        break;
      }
    }
  }
}
