// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { OnInit, Component, Inject, isDevMode, ViewChild } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { last, keys, forEach, includes, remove, get, trimEnd, find } from 'lodash';
import { IInstitution } from '../../../shared/models/participant.interface';
import { VERSION_DETAILS } from '../../../shared/constants/versions.constant';
import { HttpClient, HttpHeaders, HttpRequest, HttpErrorResponse } from '@angular/common/http';
import { Spec, Path } from 'swagger-schema-official';
import { IJWTTokenInfoPublic, IJWTPublic } from '../../../shared/models/token.interface';
import { UtilsService } from '../../../shared/utils/utils';
import { Confirm2faService } from '../../../shared/services/confirm2fa.service';
import { environment } from '../../../../environments/environment';
import { AuthService } from '../../../shared/services/auth.service';
import { NotificationService as AppNotificationService } from '../../../shared/services/notification.service';
import { NotificationService } from 'carbon-components-angular';
import { SessionService } from '../../../shared/services/session.service';
import { ENVIRONMENT } from '../../../shared/constants/general.constants';
import { CUSTOM_REGEXES, RegexMap } from '../../../shared/constants/regex.constants';
import { NgForm } from '@angular/forms';

export type ITokenActions = 'request' | 'approve' | 'reject' | 'revoke' | 'generate';
export interface ITokenDialogData {
  action: ITokenActions;

  // name not per the PR but the institution id per firebase
  institution: IInstitution;

  // for actions 'approve, revoke, reject'
  jwt_info?: IJWTPublic;
}

interface IPaths {
  [apiName: string]: { value: Path, checked: boolean, name: string }[];
}

@Component({
  templateUrl: './token-dialog.component.html',
  styleUrls: ['./token-dialog.component.scss']
})
export class TokenDialogComponent implements OnInit {

  @ViewChild('tokenForm') public tokenForm: NgForm;

  // when action 'generate' and generate button clicked
  // the one-time string displaying the token
  jwtCode: string;

  // error message to display to the ui
  errorMsg: string;

  // state indicating whether or not
  // data is being loaded to the view
  loading = false;

  // supporting form fields
  versions = VERSION_DETAILS;
  module: string;
  paths: IPaths;
  endpoints: string[];

  // token fields to be saved
  description: string;
  ver: string;
  aud: string;
  // text areas that will be converted to string[]
  accounts: string;
  ips: string;

  // all environments for related deployments
  envOptions = environment.supported_env;

  checkAllEndpoints: boolean;

  initNode = false;

  // Custom regexes for this particular form
  tokenRegexes: RegexMap = {
    description: {
      pattern: CUSTOM_REGEXES.text.pattern,
      validationText: CUSTOM_REGEXES.text.validationText
    },
    accountNameList: {
      pattern: '^[-a-z]+(, *[-a-z]+,?)*$',
      validationText: 'Invalid format. Account names must be comma separated.'
    },
    ipV4List: {
      pattern: `^${CUSTOM_REGEXES.ipV4.pattern}(, *(${CUSTOM_REGEXES.ipV4.pattern}),?)*$`,
      validationText: 'Invalid list. IP addresses must be valid and, if more than one listed, they must also be comma-separated.'
    }
  };

  constructor(
    public dialogRef: MatDialogRef<TokenDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: ITokenDialogData,
    public http: HttpClient,
    public utils: UtilsService,
    private confirm2Fa: Confirm2faService,
    private authService: AuthService,
    private notificationService: AppNotificationService,
    private notificationService2: NotificationService,
    private sessionService: SessionService
  ) {

    this.getVersions();

    // init default values
    this.endpoints = [];
    this.aud = this.sessionService.currentNode.participantId;
    this.accounts = 'issuing, default';
    this.ips = '';
    this.errorMsg = '';
    this.jwtCode = '';
    this.checkAllEndpoints = false;
    this.initNode = this.sessionService.currentNode ? true : false;
  }

  ngOnInit(): void { }

  /**
   * get list of release versions for dropdown so that user
   * can select which version api they are using and
   * retrieve a related set of endpoints
   *
   * @memberof TokenDialogComponent
   */
  getVersions() {

    // init to most recent version
    this.module = last(keys(this.versions));

    // update endpoints
    this.listEndpoints(this.module)
      .then(() => {

        // finished loading endpoints
        this.loading = false;
      }).catch(() => {

        // error loading endpoints. reset state
        this.loading = false;

        this.showErrorNotification(null, 'Error loading endpoints');
      });

  }

  /**
   * get endpoints from the api for a specific version
   *
   * @memberof TokenDialogComponent
   */
  listEndpoints(mod: string): Promise<void> {

    // start loading state
    this.loading = true;

    const self = this;

    // reset endpoints arr and reset selectAll bool
    this.endpoints = [];
    this.checkAllEndpoints = false;

    // set version
    this.ver = this.versions[Number(mod)].releaseTag;

    // reset path to empty obj
    this.paths = {};

    return new Promise((resolve) => {

      // get all apis associated with this version
      const apiArr = this.versions[this.module].config;

      // create object of paths
      for (let i = 0; i < apiArr.length; i++) {
        const apiName = apiArr[i].name;

        let displayEndpoints = false;

        // if api endpoints require token
        if (apiArr[i].wwHosted) {

          if (apiArr[i].public === true) {
            // will show public facing endpoints to a participant admin and super admin
            displayEndpoints = true;
          }

          // is super admin and not public facing
          if (get(this.authService, 'userProfile.super_permissions.roles.admin') === true && apiArr[i].public === false) {
            // will show onboarding related endpoints to a super admin
            displayEndpoints = true;
          }

        }

        if (displayEndpoints) {

          this.http.get(
            '/assets/open-api/' +
            this.module +
            '/' + apiName + '.json'
          ).toPromise().then((swaggerDef: Spec) => {

            // set default paths into array
            forEach(swaggerDef.paths, (val: Path, path: string) => {

              if (!self.paths[apiName]) {
                // if no array add array
                self.paths[apiName] = [];
              }

              // create item with path name to add to array of paths to display:

              // // the below includes the path params (ie: /test/{somePathParams} to by output to ui as /test/{somePathParams} )
              // const item = { name: /* toArray(val)[0]['x-base-url'] + */ swaggerDef.basePath + path, value: val, checked: false };

              // // the below includes the path params (ie: /test/{somePathParams} to by output to ui as /test )
              const displayPath = trimEnd((swaggerDef.basePath + path).split('{', 1)[0], '/');
              const item = {
                name: /* toArray(val)[0]['x-base-url'] + */ displayPath,
                value: val,
                checked: false
              };

              // only add if unique path (prevents the same path from
              // being included twice since certain paths share the
              // same root and different path params)
              const isUniq = find(self.paths[apiName], { 'name': displayPath });

              if (!isUniq) {
                // add path
                self.paths[apiName].push(item);
              }

            });

            if (i >= apiArr.length) {
              // no more api paths to add
              resolve();
            }

          });

        }

      }

    });

  }

  selectAll() {

    this.checkAllEndpoints = !this.checkAllEndpoints;

    // set default paths into array
    forEach(this.paths, (val: { value: Path; checked: boolean; name: string; }[], apiName: string) => {

      for (let i = 0; i < val.length; i++) {
        this.paths[apiName][i].checked = this.checkAllEndpoints;

        if (this.checkAllEndpoints) {
          // add all endpoints to list
          this.updateEndpoints(this.paths[apiName][i].name, this.paths[apiName][i].value);
        } else {
          // reset endpoints arr
          this.endpoints = [];
        }
      }

    });

  }

  updateEndpoints(name: string, val: Path) {

    if (!includes(this.endpoints, name)) {
      // add
      this.endpoints.push(name);
    } else {
      // remove
      remove(this.endpoints, (n: string) => {
        return n === name;
      });
    }

    // sort array by values
    this.endpoints = this.endpoints.sort();

    if (isDevMode) {
      console.log(this.endpoints);
    }
  }

  /**
   * closes the modal
   *
   * @memberof TokenDialogComponent
   */
  close(): void {
    this.dialogRef.close();
  }

  convertTextToArray(text: string): string[] {

    let _text = text;

    // replace newlines with commas, if any
    _text.replace(/\n/g, ',').replace(/,/g, ',');

    // remove spaces
    _text = _text.replace(/\s/g, '');

    // remove last ',' if exists
    if (last(_text) === ',') {
      _text = _text.slice(0, -1);
    }

    // convert to array
    return _text.split(',');
  }

  /**
   * Creates a new JWT request in the firebase
   *
   * @memberof TokenDialogComponent
   */
  request() {

    // reset error msg
    this.errorMsg = '';

    // convert ip text to ip arr
    const ipsArr = this.convertTextToArray(this.ips);

    // // check if ip addresses conform
    // for (let i = 0; i < ipsArr.length; i++) {
    //   if (!this.utils.isIpv4(ipsArr[i])) {
    //     passValidation = false;
    //     this.errorMsg = 'IPv4 addresses are not formated properly';
    //   }
    // }

    if (!this.tokenForm.valid) {

      // General form validation error. Place in FRONT of other errors
      this.errorMsg = 'Invalid form data was submitted and could not be processed. Please check errors and try again.';
    }

    // check if there are endpoints selected
    if (this.endpoints.length <= 0) {
      this.errorMsg = this.errorMsg + '\nPlease select endpoint(s).';
    }

    // check for a description
    if (!this.description) {
      this.errorMsg = this.errorMsg + '\nPlease provide a description.';
    }

    // check for participant id
    if (!this.aud) {
      this.errorMsg = this.errorMsg + '\nPlease provide a participant id as listed in this environment\'s participant registry.';
    }

    // check if the user stipulated an IP Address
    if (this.ips.length <= 0) {
      this.errorMsg = this.errorMsg + '\nIP Address is required.';
    }

    // convert account text to account array;
    const accountArr = this.convertTextToArray(this.accounts);

    if (this.errorMsg) {

      this.showErrorNotification();

      return;

    } else {

      // creat obj and send request if passes validation

      const self = this;

      // create 'jwt_info' token object from from
      const token: IJWTTokenInfoPublic = {
        description: this.description,
        jti: '', // to be assigned by firebase in back-end auth service
        aud: this.aud,
        acc: accountArr,
        ver: this.ver,
        ips: ipsArr,
        env: ENVIRONMENT.val,
        enp: this.endpoints
      };

      // get server side auth token for this request
      this.authService.getFirebaseIdToken(this.data.institution.info.institutionId).then((h: HttpHeaders) => {

        // create request
        const r = new HttpRequest(
          'POST',
          environment.apiRootUrl + '/jwt/request',
          token,
          { headers: h }
        );

        self.confirm2Fa.go(r)
          .then(() => {
            self.notificationService.show('success', 'Token request created');
            self.close();
          }, (err) => {
            console.log('unable to create jwt request', err);
            self.notificationService.show('error', 'Unable to create jwt request. Please contact support');
          });

      });

    }

  }

  /**
   * calls approve reject and revoke
   *
   * @memberof TokenDialogComponent
   */
  action() {

    if (this.tokenForm && !this.tokenForm.valid) {

      this.errorMsg = 'Invalid form data was submitted and could not be processed. Please try again.';
      this.showErrorNotification();

      return;
    }

    // get server side auth token for this request
    this.authService.getFirebaseIdToken(this.data.institution.info.institutionId).then((h: HttpHeaders) => {

      // create request
      const r = new HttpRequest(
        'POST',
        environment.apiRootUrl + '/jwt/' + this.data.action,
        { jti: this.data.jwt_info.jti },
        { headers: h }
      );

      this.confirm2Fa.go(r)
        .then(() => {
          // closes dialog
          this.close();

          this.notificationService.show(
            'success',
            `Successfully ${this.utils.convertVerbToPastTense(this.data.action)} token.`,
            null,
            'Success'
          );
        }, (err: HttpErrorResponse) => {

          // closes dialog
          this.close();

          this.showErrorToast(err);
        });

    });

  }

  generate() {

    // get server side auth token for this request
    this.authService.getFirebaseIdToken(this.data.institution.info.institutionId).then((h: HttpHeaders) => {

      // create request
      const r = new HttpRequest(
        'POST',
        environment.apiRootUrl + '/jwt/generate',
        { jti: this.data.jwt_info.jti },
        { headers: h }
      );

      this.confirm2Fa.go(r)
        .then((jwtCode: string) => {

          this.jwtCode = jwtCode;

        }, (err) => {

          // closes dialog
          this.close();

          this.showErrorToast(err);
        });

    });
  }

  /**
   * Shows inline error notification for modal validation
   *
   * @param {string} [title]
   * @param {string} [errorText]
   * @returns
   * @memberof TokenDialogComponent
   */
  showErrorNotification(title?: string, errorText?: string) {

    this.notificationService2.showNotification({
      type: 'error',
      title: title ? title : 'Bad Request',
      message: errorText ? errorText : this.errorMsg,
      target: '#notification'
    });

    // reset error message after error is shown
    this.errorMsg = null;

    return;
  }

  /**
   * Shows toast notification in response to
   * an error/bad request on a token action
   *
   * @param {HttpErrorResponse} err
   * @memberof TokenDialogComponent
   */
  showErrorToast(err: HttpErrorResponse) {
    let errMsg = 'Unable to ' + this.data.action + ' token. Please contact support.';

    errMsg = err.error ? errMsg + ' ' + err.error : errMsg;

    // notify user if there is an error with the request
    this.notificationService.show(
      'error',
      errMsg
    );
  }
}
