// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, HostBinding, Inject, isDevMode } from '@angular/core';
import { InlineLoading } from 'carbon-components';
import { clone, isEmpty } from 'lodash';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';
import { HttpClient, HttpRequest, HttpResponse, HttpErrorResponse } from '@angular/common/http';
import { ActivatedRoute } from '@angular/router';

export interface IConfirm2faDialogData {
  httpRequest2fa: HttpRequest<any>;
}

export interface IConfirm2faResult {
  // if 2fa was a success or not
  success: boolean;
  // response payload after successful 2fa confirmation
  result: any;
}

@Component({
  selector: 'app-confirm2fa-dialog',
  templateUrl: './confirm2fa-dialog.component.html',
  styleUrls: ['./confirm2fa-dialog.component.scss']
})
export class Confirm2faDialogComponent implements OnInit {

  timerText: number;
  disableResend: boolean;
  status:
    // waiting 2fa push confirmation
    'active' |
    // failed 2fa push confirmation
    'inactive' |
    // successful 2fa push confirmation
    'finished';

  // carbon design spinner
  private inlineLoadingInstance: any;

  private dialogDisplayTime = 3000;

  confirmCode: string;

  disabled = false;

  constructor(
    public dialogRef: MatDialogRef<Confirm2faDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: IConfirm2faDialogData,
    private http: HttpClient,
    private route: ActivatedRoute,
  ) {
    this.disabled = this.route.snapshot.queryParams.disabled;
  }

  @HostBinding('attr.class') cls = 'flex-fill';

  ngOnInit() {

    // `#my-inline-loading` is an element with `[data-inline-loading]` attribute
    this.inlineLoadingInstance = InlineLoading.create(document.getElementById('my-inline-loading'));

    if (this.disabled) {
      this.confirmCode = '';
      this.confirm();
    }
  }

  /**
   * Prevents multiple calls to send push notification
   *
   * @memberof VerifyDialogComponent
   */
  startTimer() {

    this.status = 'active';

    this.disableResend = true;
    const seconds = 10;
    this.timerText = clone(seconds);

    const timer = setInterval(() => {
      if (this.timerText >= 1) {
        this.timerText = this.timerText - 1;
      }
    }, 1000);

    setTimeout(() => {
      this.disableResend = false;
      this.timerText = null;
      clearInterval(timer);
    }, 1000 * seconds);

  }

  /**
   * Confirms 2FA verification
   * with the inputted 6-digit code
   *
   * @memberof VerifyDialogComponent
   */
  confirm() {

    // if timer has expired then retry
    if (isEmpty(this.timerText)) {

      // restart timer
      this.startTimer();

      // get authentication code IBMId
      try {

        // MUST USE .clone() to overwrite the request
        // when updating the object since HttpRequest is readonly
        const request = this.data.httpRequest2fa.clone({
          // set 2FA verification code in header to validate in middleware
          headers: this.data.httpRequest2fa.headers.set('x-verify-code', this.confirmCode)
        });

        // Send request to 2FA protected route
        const promise = this.http.request(request).toPromise();

        promise.then((response: HttpResponse<any>) => {

          // only handle responses returned in the body
          const result = response.body ? response.body : response;

          // handle result
          this.success(result);

        }, (err: HttpErrorResponse) => {

          if (isDevMode) {
            console.log('err', err);
          }

          // no permissions, forbidden to access endpoint
          if (err.status !== 403) {

            const message = err.message ? err.message : err;

            // 2FA success, but unsuccessful response
            this.success(err, err.ok);
          } else {
            // Failed 2fa
            this.failed(err);
          }
        });

      } catch (tryErr) {

        // Failed 2fa with bad error,
        // so redirect to unauthorized
        this.failed(tryErr, true);
      }

    }

  }

  /**
   * Display to user that 2fa failed
   *
   * @memberof VerifyDialogComponent
   */
  failed(err: any, close = false) {

    this.status = 'inactive';
    this.inlineLoadingInstance.setState(InlineLoading.states.INACTIVE);

    // allow for auto-closing of the dialog
    // auto-set to false to allow retry of 2FA verification
    if (close) {
      setTimeout(() => {
        this.close(false, err);
      }, this.dialogDisplayTime);
    }
  }

  /**
   * Handles successful 2FA Verification
   *
   * @param {*} result
   * @param {boolean} [responseOK=true]
   * @memberof Confirm2faDialogComponent
   */
  success(result: any, responseOK: boolean = true) {

    this.status = 'finished';
    this.inlineLoadingInstance.setState(InlineLoading.states.FINISHED);

    setTimeout(() => {

      // handles whether or not the actual request
      // NOT the 2FA verification) was successful
      this.close(responseOK, result);

    }, this.dialogDisplayTime);

  }

  /**
   * Closes the 2FA dialog
   * @param success
   * @param result
   *
   * @memberof VerifyDialogComponent
   */
  close(success: boolean, result: any) {

    const res: IConfirm2faResult = {
      result: result,
      success: success
    };

    this.dialogRef.close(res);
  }

}
