// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// // =============== usage example start ===================

//  // get a https request promise
//  function testPromise(): Promise<any> {
//     return new Promise((resolve, reject) => {
//       setTimeout(() => {
//         resolve('this is a test result');
//       }, 2000);
//     });
//   }

//   // reference the promise
//   const t = testPromise();

//   // thenable example here
//   t.then(() => {
//     console.log('from here');
//   });

//   // thenable after modal/dialog notification
//   this.confirm2fa.go(t).then((res) => {
//     console.log(res);
//   });

// // =============== usage example end ===================

import { Injectable, NgZone } from '@angular/core';
import { Router } from '@angular/router';
import { MatDialog } from '@angular/material';
import { Confirm2faDialogComponent, IConfirm2faResult } from '../components/confirm2fa-dialog/confirm2fa-dialog.component';
import { HttpRequest } from '@angular/common/http';

@Injectable()
export class Confirm2faService {

    constructor(
        public dialog: MatDialog,
        public ngZone: NgZone,
        public router: Router
    ) { }

    /**
     * Shows two-factor verify dialog and returns the result after 2fa verification
     *
     * @param {Promise<any>} httpRequest2fa must be a route that requires 2fa on the server
     * @returns {*} result after confirming 2fa
     * @memberof Confirm2faService
     */
    go(httpRequest2fa?: HttpRequest<any>): Promise<any> {

        const self = this;

        return new Promise((resolve, reject) => {

            // use setTimeout to prevent rendering issue "ExpressionChangedAfterItHasBeenCheckedError"
            setTimeout(() => {

                const dialog = this.dialog.open(Confirm2faDialogComponent, {
                    data: {
                        httpRequest2fa: httpRequest2fa
                    },
                    disableClose: true
                });

                dialog.afterClosed().subscribe((res: IConfirm2faResult) => {

                    // 2FA-protected request was successful
                    if (res && res.success) {
                        // success return response payload
                        self.ngZone.run(() => {
                            resolve(res.result);
                        });
                    } else {
                        // 2FA-protected request was a failure.
                        // Return appropriate error info for the user
                        self.ngZone.run(() => {
                            reject(res.result);
                        });
                    }

                });

            });

        });

    }

}
