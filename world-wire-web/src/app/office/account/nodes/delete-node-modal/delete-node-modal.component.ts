// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, Inject, isDevMode } from '@angular/core';
import { BaseModal } from 'carbon-components-angular';
import { environment } from '../../../../../environments/environment';
import { AuthService } from '../../../../shared/services/auth.service';
import { HttpHeaders, HttpClient } from '@angular/common/http';
import { INodeAutomation } from '../../../../shared/models/node.interface';

@Component({
  selector: 'app-delete-node-modal',
  templateUrl: './delete-node-modal.component.html',
  styleUrls: ['./delete-node-modal.component.scss']
})
export class DeleteNodeModalComponent extends BaseModal {

  node: INodeAutomation;

  confirmParticipantId: string;

  constructor(
    @Inject('MODAL_DATA') public data: any,
    private authService: AuthService,
    private http: HttpClient
  ) {

    super();
    this.node = data.node;
  }

  /**
 * Deletes Node in both AWS and Firebase
 * Functionality should only available
 * in /office when logged in as super-admin
 * @param nodeId
 */
  public async deleteNode() {

    if (this.confirmParticipantId === this.node.participantId) {
      if (isDevMode) {
        console.log('deleting: ', this.node.participantId);
      }

      // const deleteEndpoint = `${environment.apiRootUrl}/manage/market-maker/aws/stack/delete/`
      //   + `${this.node.institutionId}/${this.node.participantId}`;

      // // get fid from user to authenticate endpoint
      // const headers: HttpHeaders = await this.authService.getFirebaseIdToken(this.node.institutionId);

      // // constructing request options
      // const options = {
      //   headers: headers
      // };

      // // Call API to delete node
      // const apiRequest: Promise<any> = this.http
      //   .delete(deleteEndpoint,
      //     options
      //   ).toPromise();

      // await apiRequest.then((response: any) => {

      //   if (isDevMode) {
      //     console.log('response: ', response);
      //     console.log('removed node: ', this.node);
      //   }

      //   this.closeModal();
      // });
    }

    return;
  }
}
