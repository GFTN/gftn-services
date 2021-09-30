// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Inject, isDevMode } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { FormGroup } from '@angular/forms';
import * as cc from 'currency-codes';
import * as countries from 'i18n-iso-countries';
import { BlocklistRequest, BlocklistType } from '../../../shared/models/blocklist.interface';
import { BlocklistService } from '../../../shared/services/blocklist.service';
import { WorldWireError } from '../../../shared/models/error.interface';
import { AngularFireDatabase } from '@angular/fire/database';
import { find, startCase } from 'lodash';
import { HttpErrorResponse } from '@angular/common/http';
import { NotificationService } from 'carbon-components-angular';
import { AuthService } from '../../../shared/services/auth.service';

export interface IBlockListDialogData {
  action: 'add' | 'remove';
  type: BlocklistType;
  name?: string;
  isoCode?: string;
  alternateName?: string;
}

export interface BlocklistOption {
  name: string;
  isoCode: string;
}

@Component({
  selector: 'app-blocklist-dialog',
  templateUrl: './blocklist-dialog.component.html',
  providers: [
    BlocklistService
  ]
})
export class BlocklistDialogComponent implements OnInit {

  public addCurrencyFormGroup: FormGroup;
  public addCountryFormGroup: FormGroup;

  public blocklistOptions: BlocklistOption[] = [];

  request: BlocklistRequest;

  public loadingRequest = false;

  public constructor(
    private authService: AuthService,
    public dialogRef: MatDialogRef<BlocklistDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: IBlockListDialogData,
    private blocklistService: BlocklistService,
    private notificationService: NotificationService,
    private db: AngularFireDatabase
  ) {
    countries.registerLocale(require(`i18n-iso-countries/langs/${this.authService.userLangShort}.json`));
  }

  public ngOnInit(): void {

    // get dropdown options
    switch (this.data.type) {
      case 'country': {
        this.getCountries();

        break;
      }
      case 'currency': {
        this.getCurrencies();

        break;
      }
      default:
        break;
    }





    const requestValue = this.data.isoCode ? this.data.isoCode : '';

    this.request = {
      type: this.data.type,
      value: requestValue
    };

  }

  private getCountries() {
    for (const [key, value] of Object.entries(countries.getNames(this.authService.userLangShort))) {
      this.blocklistOptions.push({
        name: value,
        isoCode: key,
      });
    }

    this.blocklistOptions.sort((a: BlocklistOption, b: BlocklistOption) => {
      return a.name.localeCompare(b.name);
    });
  }

  private getCurrencies() {
    for (const code of cc.codes()) {
      const found = find(this.blocklistOptions, (option) => {
        return option.isoCode === code;
      });

      // prevent duplicate codes
      if (!found) {
        this.blocklistOptions.push({
          name: cc.code(code).currency,
          isoCode: code,
        });
      }
    }

    // sort currencies alphabetically by:
    // 1. isoCode
    // 2. full name
    this.blocklistOptions.sort((a, b) => {
      // compare isoCodes alphabetically first
      const codeA = a.isoCode;
      const codeB = b.isoCode;

      if (codeA < codeB) {
        return -1;
      }
      if (codeA > codeB) {
        return 1;
      }

      // move on to comparing names if isoCode is the same
      const nameA = a.name;
      const nameB = b.name;

      if (nameA < nameB) {
        return -1;
      }
      if (nameA > nameB) {
        return 1;
      }

      // names must be equal
      return 0;
    });
  }

  public setModalTitle(): string {
    switch (this.data.action) {
      case 'add':
        return 'Add New ' + startCase(this.data.type) + ' to Blocklist';
      case 'remove':
        return 'Remove ' + startCase(this.data.type) + ' from Blocklist';
      default:
        break;
    }
  }

  public closeModal(): void {
    this.dialogRef.close();
  }

  /**
   * Save initial blocklist request
   *
   * @memberof BlocklistDialogComponent
   */
  public submitRequest() {

    this.blocklistService.addToBlocklist(this.request, 'request').then((response: WorldWireError) => {
      const approvalId: string = response.msg ? response.msg : '';

      // log reference to approval in blocklist request
      if (approvalId) {
        this.request.approvalIds = [approvalId];

        this.db.database.ref('blocklist_requests').child(this.request.type).update({
          [this.request.value]: this.request
        });
      }

      this.loadingRequest = true;

      this.closeModal();
    }).catch((err: HttpErrorResponse) => {
      if (isDevMode) {
        console.log('err', err);
      }

      this.loadingRequest = true;

      // temp handling for 200 response
      if (err.status === 200) {
        this.closeModal();
      } else {
        // throw error for other responses
        this.throwError(err);
      }
    });

  }

  /**
   * Remove from blocklist
   *
   * @memberof BlocklistDialogComponent
   */
  public removalRequest() {

    this.blocklistService.removeFromBlocklist(this.request, 'request').then((response: WorldWireError) => {
      const approvalId: string = response.msg ? response.msg : '';

      // log reference to approval in blocklist request
      if (approvalId) {
        this.request.approvalIds = [approvalId];

        this.db.database.ref('blocklist_requests').child(this.request.type).update({
          [this.request.value]: this.request
        });
      }

      this.loadingRequest = true;

      this.closeModal();
    }).catch((err: HttpErrorResponse) => {
      if (isDevMode) {
        console.log('err', err);
      }

      this.loadingRequest = true;

      // temp handling for 200 response
      if (err.status === 200) {
        this.closeModal();
      } else {
        // throw error for other responses
        this.throwError(err);
      }
    });

  }

  /**
   * Error handling for Blocklist requests
   *
   * @private
   * @param {HttpErrorResponse} error
   * @memberof BlocklistDialogComponent
   */
  private throwError(error: HttpErrorResponse) {
    this.loadingRequest = false;

    let errorText = 'add this ' + this.data.type.toLowerCase() + ' to the blocklist';

    if (this.data.action === 'remove') {
      errorText = 'remove this ' + this.data.type.toLowerCase() + ' from the blocklist';
    }

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
        // for 500s and other errors: Unexpected ERROR notification
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

}
