// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Observable } from '@firebase/util';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { ActivatedRouteSnapshot, Resolve } from '@angular/router';
import { Injectable } from '@angular/core';

import { SessionService } from '../../../../shared/services/session.service';
import { AccountService } from '../../../shared/services/account.service';
import { Participant } from '../../../../shared/models/participant.interface';
import { ENVIRONMENT } from '../../../../shared/constants/general.constants';
import { AuthService } from '../../../../shared/services/auth.service';
import { UtilsService } from '../../../../shared/utils/utils';
import { INodeAutomation } from '../../../../shared/models/node.interface';




@Injectable()
export class AccountsResolve implements Resolve<Participant> {

    constructor(
        private utils: UtilsService,
        private sessionService: SessionService,
        private accountService: AccountService
    ) { }

    resolve(
        route: ActivatedRouteSnapshot,
        // state: RouterStateSnapshot
    ): Observable<any> | Promise<any> | any {

        return new Promise(async (resolve) => {

            if (!this.sessionService.currentNode && this.sessionService.institutionNodes) {

                const firstNode = this.sessionService.institutionNodes[0];

                // construct cookie name based on Institution slug (don't expose uuid)
                const cookieName = 'pid-' + this.sessionService.institution.info.institutionId;

                const cachedNodeId = this.utils.getCookie(cookieName);

                const cachedNode: INodeAutomation = this.sessionService.institutionNodes.find((node: INodeAutomation) => {
                    return node.participantId === cachedNodeId;
                });

                if (cachedNode) {
                    // set current node as cached node from cookie
                    this.sessionService.setCurrentNode(cachedNode);

                } else {
                    // cached node doesn't exist. get first ever node
                    this.sessionService.setCurrentNode(firstNode);
                }
            }

            if (this.sessionService.currentNode) {

                this.accountService.apiRoot = `https://${this.sessionService.currentNode.participantId}.${ENVIRONMENT.envApiRoot}`;

                this.accountService.globalRoot = `https://${this.sessionService.currentNode.participantId}.${ENVIRONMENT.envGlobalRoot}`;

                if (this.accountService.participantDetails) {
                    resolve(this.accountService.participantDetails);
                } else {
                    // always resolve to prevent page from breaking unexpectedly
                    resolve(null);
                }

            } else {
                // always resolve to prevent page from breaking unexpectedly
                resolve(null);
            }
        });

    }
}
