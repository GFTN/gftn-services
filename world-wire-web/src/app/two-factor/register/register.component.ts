// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { DomSanitizer, SafeResourceUrl } from '@angular/platform-browser';
import * as QRCode from 'qrcode';
import { TOTPResponse, TOTPRegistrationData, TokenBody } from '../../shared/models/totp.interface';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.scss']
})
export class RegisterComponent implements OnInit {

  public steps: { text: string, desc: string, state: string[] }[];
  public currentStep: number;
  public registered: boolean;
  public registrationCode: string;
  public qrCodeError: boolean;
  public timer: number | string;
  private interval: any;

  // for old 2FA method with IBM Verify Push
  // private queryParams: IRegistrationResponse;

  private registerData: TOTPRegistrationData;

  @ViewChild('qrCodeCanvas') qrCodeCanvas: ElementRef;

  // used in model to confirm 2FA code
  public confirmCode: string;
  confirmError = false;

  constructor(
    private router: Router,
    private http: HttpClient,
    private route: ActivatedRoute,
    private _sanitizer: DomSanitizer
  ) {
    this.registered = false;
    this.qrCodeError = true;
    this.currentStep = 0;
  }

  ngOnInit() {

    this.defineSteps();
    this.getQrCode();
  }

  defineSteps() {
    this.steps = [
      {
        // step 1
        text: 'Download',
        desc: 'Download Authenticator App',
        state: ['current']
      },
      {
        // step 2
        text: 'Configure',
        desc: 'Scan QR Code',
        state: ['incomplete']
      },
      {
        // step 3
        text: 'Confirm',
        desc: 'Confirm Two Factor Registration',
        state: ['incomplete']
      },
      {
        // step 4
        text: 'Login',
        desc: 'Success! Now you can login',
        state: ['incomplete']
      }
    ];
  }

  /**
   * update view for progress indicator
   * by progessing the steps forward or backward
   * @param option
   */
  selectStep(option: 'next' | 'prev') {

    // progress forward
    if (option === 'next') {
      this.currentStep = this.currentStep + 1;

      if (this.currentStep >= this.steps.length) {
        this.currentStep = this.steps.length;
      } else {
        // update the progress indicator at the top of the page
        this.updateStage();
      }
    }

    // progress backward
    if (option === 'prev') {
      this.currentStep = this.currentStep - 1;
      if (this.currentStep <= 0) {
        // min limit check
        this.currentStep = 0;
      } else {
        // update the progress indicator at the top of the page
        this.updateStage();
      }

    }

  }

  updateStage() {

    // update the check marks in the progress indicator
    for (let i = 0; i < this.steps.length; i++) {

      // current step
      if (i === this.currentStep) {
        this.steps[i].state = ['current'];
      }

      // completed steps
      if (i < this.currentStep) {
        this.steps[i].state = ['complete'];
      }

      // not-started steps
      if (i > this.currentStep) {
        this.steps[i].state = ['incomplete'];
      }

    }

  }

  /**
   * Generates QR Code for user to scan into
   * their authenticator app
   */
  getQrCode() {

    // get qrCode from query param
    const queryData = this.route.snapshot.queryParams.data ? this.route.snapshot.queryParams.data : null;

    this.registerData = queryData ? JSON.parse(atob(queryData)) : null;

    try {
      this.http.get(environment.apiRootUrl + '/totp/' + this.registerData.email,
        { withCredentials: true })
        .toPromise()
        .then(async (response: TOTPResponse) => {

          const qrCodeUri = response.data.qrcodeURI;

          // only display QR Code if QRCode Uri  was successfulLY generated
          if (qrCodeUri) {

            // QR Code exists for generation
            this.qrCodeError = false;

            try {

              // rendering options for the QR Code
              const qrOptions = {
                errorCorrectionLevel: 'H',
                type: 'svg',
                rendererOpts: {
                  quality: 0.3
                }
              };

              // convert URI to QR Code in the view
              await QRCode.toDataURL(this.qrCodeCanvas.nativeElement, qrCodeUri, qrOptions);
            } catch (err) {
              console.error(err);
            }
          }

        }, (err: any) => {
          // error unable to get qr code...
          console.log('error unable to get qr code:', err);
        });
    } catch (error) {
      console.log(error);
    }

  }

  /**
   * Called by form to register and confirm the
   * authentication token after scanning the QR code
   *
   * @memberof RegisterComponent
   */
  submitConfirmForm() {

    // reset error to false when submitting form
    this.confirmError = false;

    // constructing body to pass into 2FA registration confirmation
    const body: TokenBody = {
      token: this.confirmCode.toString()
    };

    this.http.post(`${environment.apiRootUrl}/totp/${this.registerData.email}/confirm`,
      body,
      { withCredentials: true })
      .toPromise()
      .then((response: TOTPResponse) => {

        if (response.success) {

          // Successful 2fa registration
          this.selectStep('next');
          this.registered = true;

        } else {

          // Failed 2fa registration
          this.confirmError = true;
        }
      }, (err: any) => {

        // Unexpected error when hitting the endpoint
        console.log(err);
        this.confirmError = true;
      });
  }

  /**
   * Checks if the user has registered the IBM Verify QR Code successfully
   * NOTE: Will wait n minutes for a response while the user opens up the IBM Verify
   * app and scans the code from their mobile device
   *
   * @memberof RegisterComponent
   */
  // checkRegistrationStatus() {
  //   const check = this.http.post(
  //     environment.apiRootUrl + '/2fa/status/' +
  //     // queryParams.code is not the ICE registration code
  //     // this.queryParams.code + '/' +
  //     this.queryParams.ci + '/' +
  //     this.queryParams.fid, {})
  //     .toPromise();

  //   check.then((data: { status: number, desc: string }) => {

  //     if (data.status === 1) {
  //       setTimeout(() => {
  //         // wait n seconds and retry
  //         this.checkRegistrationStatus();
  //       }, 2000);
  //     } else if (data.status === 2) {
  //       this.success();
  //     } else {
  //       this.failed('request returned false');
  //     }

  //   }, (err) => {
  //     this.failed(err);
  //   });
  // }

  /**
   * Updates view when the user successfully registered 2fa
   *
   * @memberof RegisterComponent
   */
  success() {
    // end timer
    clearInterval(this.interval);

    // set interval text to empty string
    this.timer = '';
    this.registered = true;

    this.currentStep = 2;
    this.updateStage();
  }

  /**
   * Updates view when something went wrong registering 2fa
   *
   * @param {*} [err]
   * @memberof RegisterComponent
   */
  failed(err?: any) {
    console.log('failed to register 2fa: ' + err);
    this.currentStep = 1;
    this.updateStage();
    this.qrCodeError = true;
  }

  login() {
    this.router.navigate(['/login']);
  }

}
