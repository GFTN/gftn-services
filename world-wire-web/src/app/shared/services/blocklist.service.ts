// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import * as _ from 'lodash';
import { ENVIRONMENT } from '../constants/general.constants';
import { Injectable } from '@angular/core';
import { AngularFireDatabase } from '@angular/fire/database';
import { HttpClient, HttpHeaders, HttpRequest } from '@angular/common/http';

import { Blocklist, BlocklistRequest, BlocklistType } from '../models/blocklist.interface';
import { AuthService } from './auth.service';
import { ApprovalPermission } from '../models/approval.interface';


@Injectable()
export class BlocklistService {

  requestUrlBase = `https://admin.${ENVIRONMENT.envGlobalRoot}/admin/v1/admin/blocklist`;

  constructor(
    private db: AngularFireDatabase,
    public http: HttpClient,
    private authService: AuthService,
  ) { }

  /**
   * Gets a list of values from the blocklist.
   *
   * @param {('country' | 'currency' | 'institution')} type
   * @returns {Promise<string[]>}
   * @memberof BlocklistService
   */
  public async getBlocklist(blocklistType: BlocklistType): Promise<string[]> {

    const requestUrl = `${this.requestUrlBase}?type=${blocklistType}`;

    try {
      const h: HttpHeaders = await this.authService.getFirebaseIdToken();

      const options = {
        headers: h
      };

      const blocklist: Blocklist = await this.http.get(
        requestUrl,
        options
      ).toPromise() as Blocklist;

      return blocklist[0].value;

    } catch (err) {
      return null;
    }
  }

  /**
   * Add to blocklist/post request
   *
   * @memberof BlocklistService
   */
  public async addToBlocklist(request: BlocklistRequest, permission: ApprovalPermission): Promise<any> {

    const requestUrl = this.requestUrlBase;

    const latestApprovalId = request.approvalIds ? request.approvalIds[request.approvalIds.length - 1] : null;

    let h: HttpHeaders = await this.authService.getFirebaseIdToken();

    if (permission === 'request') {
      h = this.authService.addMakerCheckerHeaders(h, permission);
    } else {
      h = this.authService.addMakerCheckerHeaders(h, permission, latestApprovalId);
    }

    const options = {
      headers: h
    };

    const body: Blocklist = {
      type: request.type,
      value: [request.value]
    };

    return await this.http.post(
      requestUrl,
      body,
      options
    ).toPromise();

  }

  /**
   * Sends API request to remove currency, country,
   * institution from Blocklist.
   *
   * @param {BlocklistRequest} request
   * @param {ApprovalPermission} permission
   * @returns {Promise<any>}
   * @memberof BlocklistService
   */
  public async removeFromBlocklist(request: BlocklistRequest, permission: ApprovalPermission): Promise<any> {

    const requestUrl = this.requestUrlBase;

    const latestApprovalId = request.approvalIds ? request.approvalIds[request.approvalIds.length - 1] : null;

    let h: HttpHeaders = await this.authService.getFirebaseIdToken();

    if (permission === 'request') {
      h = this.authService.addMakerCheckerHeaders(h, permission);
    } else {
      h = this.authService.addMakerCheckerHeaders(h, permission, latestApprovalId);
    }

    const body: Blocklist = {
      type: request.type,
      value: [request.value]
    };

    // Per Angular 7+ spec, body for DELETE can be included in options
    const options = {
      headers: h,
      body: body
    };

    return await this.http.delete(
      requestUrl,
      options
    ).toPromise();
  }

  /**
   * Gets all blocklist requests for approval data
   *
   * @returns {Promise<BlocklistRequest>}
   * @memberof BlocklistService
   */
  public getBlocklistRequests(type: BlocklistType): Promise<any> {
    return new Promise((resolve, reject) => {
      this.db.database.ref('blocklist_requests')
        .child(type)
        .once('value', (data: firebase.database.DataSnapshot) => {

          resolve(data);
        });
    });
  }

}
