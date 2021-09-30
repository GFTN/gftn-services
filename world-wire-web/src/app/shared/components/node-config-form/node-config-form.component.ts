// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Input, OnChanges, Output, EventEmitter, ViewChild, ElementRef, isDevMode } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient, HttpHeaders, HttpErrorResponse } from '@angular/common/http';
import { NgForm } from '@angular/forms';
import { AngularFireDatabase } from '@angular/fire/database';
import * as _ from 'lodash';

// import constants
import { VERSION_DETAILS } from '../../constants/versions.constant';
import { ENVIRONMENT } from '../../constants/general.constants';
import { CUSTOM_REGEXES, CustomRegex } from '../../constants/regex.constants';

// import necessary custom interfaces
import { INodeAutomation, NodeConfigData } from '../../models/node.interface';
import { IInstitution } from '../../models/participant.interface';

import { ListItem, NotificationService } from 'carbon-components-angular';
import { AuthService } from '../../services/auth.service';
import { ApprovalPermission, ApprovalInfo } from '../../models/approval.interface';
import { WorldWireError } from '../../models/error.interface';
import { SuperApprovalsModel } from '../../models/super-approval.model';
import { NodeService } from '../../../office/account/shared/node.service';

@Component({
  selector: 'app-node-config-form',
  templateUrl: './node-config-form.component.html',
  styleUrls: ['./node-config-form.component.scss'],
  providers: [
    SuperApprovalsModel
  ]
})
export class NodeConfigFormComponent implements OnInit, OnChanges {

  // OPTIONAL: Information (labels, ids) for Configuration Steps
  // This is only necessary for creation of a purely new node
  // in order to enforce the proper validation for certain fields
  @Input() configSteps: any[];

  // OPTIONAL: Can pass in the config data of a node if it exists
  @Input() configData?: NodeConfigData;

  // hack to pass in the institution from SessionService since
  // this component is shared between portal and office
  @Input() institution: IInstitution;

  @Input() theme: 'light' | 'dark' = 'light';

  @ViewChild('participantIdBaseElement') participantIdBase: ElementRef;

  public currentStep = 0;

  // stores old configuration to track any form changes
  // only applies to existing nodes upon update/view
  public defaultData: NodeConfigData;

  latestApproval: ApprovalInfo = {};

  // dropdown options
  regionOptions: ListItem[];
  apiVersionOptions: ListItem[] = [];
  participantRoleOptions: ListItem[] = [];

  @Output()
  submitted = new EventEmitter<string>();
  saving = false;

  showLoader = false;

  success = false;

  dbRef: firebase.database.Reference;

  // TODO: remove when callback is refactored
  defaultCallback = 'http://34.216.224.221:31002/v1/callback';

  defaultRdoClient = 'http://34.216.224.221:11003/v1';

  latestApiVersion = 'latest';

  regexes: { [key: string]: CustomRegex; };

  constructor(
    private router: Router,
    private db: AngularFireDatabase,
    public http: HttpClient,
    private superApprovals: SuperApprovalsModel,
    public authService: AuthService,
    public nodeService: NodeService,
    private notificationService: NotificationService,
  ) {

    this.regexes = CUSTOM_REGEXES;

    // https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html

    // initialize static options for region dropdown
    this.regionOptions = [
      {
        content: 'Asia Pacific',
        value: 'ap-southeast-1',
        selected: false
      },
      {
        content: 'Europe',
        value: 'eu-central-1',
        selected: false
      },
      {
        content: 'North America',
        value: 'us-east-2',
        selected: false
      }
    ];

    this.participantRoleOptions = [
      {
        content: 'Market Maker',
        value: 'MM',
        selected: false
      },
      {
        content: 'Issuer',
        value: 'IS',
        selected: false
      }
    ];
  }

  ngOnInit() {

    this.dbRef = this.db.database.ref(`/participants/${this.institution.info.institutionId}/nodes`);

    // initialize version options
    for (const i of _.toArray(VERSION_DETAILS)) {

      const option: ListItem = {
        value: i.releaseTag,
        content: i.version,
        selected: false,
      };

      this.apiVersionOptions.push(option);

    }

    // sort API versions from newest to oldest
    this.apiVersionOptions.reverse();

    // initializing current form data
    if (this.configData) {

      // set passed in configuration as default to compare against
      this.getDefaultData();

    } else {

      const participantRole = this.institution.info.kind === 'Bank' ? this.participantRoleOptions[1] : this.participantRoleOptions[0];

      // new empty form if there is no existing configuration to view/edit
      this.configData = {
        institutionId: this.institution.info.institutionId,
        countryCode: this.institution.info.country,
        initialized: false,
        role: participantRole.value,
        participantIdBase: '',
        participantId: '',
      };
    }
  }

  ngOnChanges() {
    // console.log('form configs: ', this.configData);
  }

  /**
   * GOTO previous step in the form
   * */
  public previousStep() {
    if (this.currentStep >= 0) {
      this.currentStep--;
    }
  }

  /**
   * GOTO next step in the form
   * */
  public nextStep(f: NgForm) {
    if (this.currentStep < (this.configSteps.length - 1) && f.valid) {

      // set statuses of config steps
      this.configSteps[this.currentStep].state = 'complete';
      this.configSteps[this.currentStep + 1].state = 'current';

      // increment current Step
      this.currentStep++;

    }
  }

  /**
   * Helper function to sync validation for participant id
   * when dropdowns, inputs, etc. other related elements have changed
   */
  public registerParticipantIdChange() {
    if (this.configData.participantIdBase !== '') {
      const element: HTMLElement = this.participantIdBase.nativeElement;
      element.dispatchEvent(new Event('input'));
    }
  }

  /**
   * Construct Participant ID from participantIdBase
   * */
  public getParticipantId() {

    // const env = '';
    // const env = this.getSuffix();

    // set participantId
    this.configData.participantId = _.kebabCase(this.configData.participantIdBase);

    const matches = this.configData.participantId.match(/[a-z]-[0-9]/g);

    const splits = this.configData.participantId.split(/[a-z]-[0-9]/g);

    if (matches && splits) {
      for (let i = 0; i < matches.length; i++) {

        let match = matches[i];

        match = match.replace('-', '');

        splits[i] = splits[i] + match;
      }
    }

    this.configData.participantId = splits ? splits.join('') : this.configData.participantId;

    // this.configData.participantId = _.camelCase(this.configData.participantId).toLowerCase();

    return this.configData.participantId;
  }

  /**
   * Get name suffix
   *
   * @returns {string}
   * @memberof NodeConfigFormComponent
   */
  getSuffix(): string {

    let suffix = '';

    suffix = (ENVIRONMENT.name && ENVIRONMENT.name !== 'prod') ? ENVIRONMENT.name : '';

    return suffix;
  }

  /**
   * Append url to participant ID for Participant/Client API Domain
   * */
  public getClientDomain() {

    const suffix = this.getSuffix();

    return `${this.configData.participantId}.worldwire-${suffix}.io`;
  }

  public async getDefaultData() {

    // get approvals
    if (this.configData.approvalIds) {
      const latestApprovalId = this.configData.approvalIds[this.configData.approvalIds.length - 1];

      this.latestApproval = await this.superApprovals.getApprovalInfo(latestApprovalId);
    }

    // clone and save default data for comparison
    this.defaultData = _.cloneDeep(this.configData);

    // if failed, set up view to retry
    if (this.nodeService.getStatus(this.defaultData.status).includes('failed')) {
      this.configData.update = this.defaultData;
    }

    if (this.configData.update) {
      this.configData = this.configData.update;
    }
  }

  /**
   * Checks if form has changed from default values
   * true = form data has changes
   * false = form data is same as defualt
   *
   * @returns {boolean}
   * @memberof NodeConfigFormComponent
   */
  public isDataUpdated(): boolean {
    return !_.isEqual(this.configData, this.defaultData);
  }

  public validForm(f: NgForm) {
    if (this.configData.update) {
      return f.valid && !this.saving && !this.success;
    }
    return f.valid && this.isDataUpdated() && !this.saving && !this.success;
  }

  /**
   * Send form data to the API
   * to construct or update node configuration
   * from deployment service
   *
   * @param {NgForm} f
   * @returns
   * @memberof NodeConfigFormComponent
   */
  public async sendFormData(f: NgForm): Promise<void> {

    // Secondary Check: Invalid form. Throw 400 error.
    if (!this.validForm(f)) {

      const error: HttpErrorResponse = new HttpErrorResponse({
        status: 400,
      });
      this.throwError(error);
      return;
    }

    // Secondary Check: requestor CANNOT be same as approver. throw 401 unauthorized error.
    if (this.latestApproval && !this.latestApproval.requestApprovedBy && this.latestApproval.requestInitiatedBy === this.authService.userProfile.profile.email) {
      const error: HttpErrorResponse = new HttpErrorResponse({
        status: 401,
      });
      this.throwError(error);
      return;
    }

    // make sure node is not empty in case of accidental form submission (enter key)
    if (this.validForm(f)) {

      // disable save button to prevent double submission
      this.saving = true;

      // show saving loader
      this.showLoader = true;

      const statuses: string[] = this.configData.status ? this.configData.status : [];

      let permission: ApprovalPermission = 'request';

      if (statuses.length > 0 && this.nodeService.getStatus(statuses).includes('failed')) {
        permission = 'request';
      } else {
        permission = this.latestApproval && this.latestApproval.requestInitiatedBy && !this.latestApproval.requestApprovedBy ? 'approve' : 'request';
      }

      // maker = request, checker = approve
      const status = permission === 'request' ? 'pending' : 'configuring';

      if (statuses.length > 0) {
        // status array contains 'failed' in code
        if (this.nodeService.getStatus(statuses).includes('failed')) {
          statuses.unshift(status);
        } else {
          // replace only first item in array (contains current status of node configuration)
          statuses[0] = status;
        }
      } else {
        statuses.push(status);
      }

      const automationUrl = `https://admin.${ENVIRONMENT.envGlobalRoot}/deployment/v1/deploy/participant`;

      const latestApprovalId = this.latestApproval ? this.latestApproval.key : null;

      let h: HttpHeaders = await this.authService.getFirebaseIdToken(this.configData.institutionId, this.configData.participantId);

      if (permission === 'approve') {

        h = this.authService.addMakerCheckerHeaders(h, 'approve', latestApprovalId);
      } else {
        h = this.authService.addMakerCheckerHeaders(h, 'request');
      }

      // constructing request options
      const options = {
        headers: h
      };


      const node: INodeAutomation = {
        institutionId: this.configData.institutionId,
        status: statuses,
        bic: this.configData.bic,
        role: this.configData.role,
        participantId: this.configData.participantId,
        initialized: this.configData.initialized,
        countryCode: this.configData.countryCode,
      };


      if (isDevMode) {
        console.log('automationUrl', automationUrl);
        console.log('options', options);
        console.log('node', node);
      }

      // Call Automation endpoint to spin up and configure node (maker/checker process)
      // This endpoint will also create the node in firebase
      const apiRequest: Promise<any> = this.http
        .post(automationUrl,
          node,
          options
        ).toPromise();

      // Empty Promise for testing API
      // const apiRequest = new Promise((resolve) => { resolve(); });

      // delay 2.5 secs to create visual
      // feedback of the node posting to firebase
      setTimeout(() => {
        apiRequest.then(async (response: any) => {
          if (isDevMode) {
            console.log('response', response);
          }

          // update record for checker to see approval
          if (permission === 'request' && response.msg) {

            // create new record, if does not exist
            if (!this.configData.approvalIds) {
              await this.dbRef.child(node.participantId).update(node);
            }

            // init empty approvals array if not set
            const newApprovals = this.configData.approvalIds ? this.configData.approvalIds : [];

            newApprovals.push(response.msg);

            const requestStatuses: string[] = node.status;

            if (requestStatuses.length > 0 && !requestStatuses[0].includes('pending')) {
              requestStatuses.unshift('pending');
            }

            // keep track of approval ID in node log
            await this.dbRef.child(node.participantId).update({
              status: requestStatuses,
              approvalIds: newApprovals,
              update: node,
            });
          }

          if (permission === 'approve') {
            await this.dbRef.child(node.participantId).update({
              update: null,
            });
          }

          // emit success
          this.success = true;

          // reload form state
          this.saving = false;

        }, async (error: HttpErrorResponse) => {

          if (isDevMode) {
            console.log(error);
          }

          this.throwError(error);

          if (permission === 'approve') {

            const data = await this.dbRef.child(node.participantId).once('value');

            const currentNodeConfig: INodeAutomation = data.val() ? data.val() : null;

            if (currentNodeConfig) {
              const requestStatuses: string[] = currentNodeConfig.status;

              // set status to failed if no failures returned back
              if (requestStatuses.length > 0 && !requestStatuses[0].includes('failed')) {
                requestStatuses[0] = 'configuration_failed';

                this.dbRef.child(node.participantId).update({
                  status: requestStatuses
                });
              }
            }
          }

          // reload form state
          this.saving = false;

          this.showLoader = !this.showLoader;

        }); // end Promise.all()
      }, 2500);
    }

    return;
  }

  public rejectFormRequest(approvalId: string) {
    if (isDevMode) {
      console.log('rejected ', approvalId);
    }

    this.configData.approvalIds = _.remove(this.configData.approvalIds, (id: string) => {
      return id === approvalId;
    });
    this.dbRef.child(this.configData.participantId).update({
      status: ['deleted'],
      update: null,
      approvalIds: this.configData.approvalIds
    });

    this.submitted.emit('complete');
  }

  /**
   * Throw error message to the screen for user
   *
   * @param {HttpErrorResponse} error
   * @memberof NodeConfigFormComponent
   */
  public throwError(error: HttpErrorResponse) {

    const errorText = this.configSteps ? 'create a new participant' : 'update node configuration';

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
        const message = wwError && wwError.details ? wwError.details : 'Invalid form data was submitted and could not be processed. Please try again.';
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
        const message = wwError && wwError.details ? wwError.details : 'Network could not be reached to make this request. Please contact administrator.';
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
        // for 404s and other errors: Unexpected ERROR notification
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
    this.showLoader = !this.showLoader;

    this.success = !this.success;

    if (this.defaultData) {

      // refresh default data
      this.getDefaultData();

      // notify parent component of form submission success
      // emit event to exit out of form (if in modal)
      this.submitted.emit('complete');

    } else {
      // redirect back to main nodes page
      this.router.navigate([`/office/account/${this.institution.info.slug}/nodes`]);
    }
  }

}
