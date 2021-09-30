// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 

import { Observable } from '@firebase/util';
import { ActivatedRouteSnapshot, Resolve } from '@angular/router';
import { Injectable } from '@angular/core';

import { SessionService } from '../services/session.service';

import { INodeAutomation } from '../models/node.interface';
import { IInstitution } from '../models/participant.interface';
import { find, filter } from 'lodash';
import { UtilsService } from '../utils/utils';

@Injectable()
export class NodeResolve implements Resolve<INodeAutomation> {

    constructor(
        private utils: UtilsService,
        private sessionService: SessionService
    ) { }

    resolve(
        route: ActivatedRouteSnapshot,
        // state: RouterStateSnapshot
    ): Observable<any> | Promise<any> | any {

        // resolve with participant details

        return new Promise((resolve) => {

            // check to make sure there is a participant/institution attached
            if (!this.sessionService.institution) {
                resolve(0);
            } else {

                // Nodes already set in Session Service for this currently set institution
                if (this.sessionService.institutionNodes && this.sessionService.institutionNodes.length > 0
                    && this.sessionService.institution.info.slug === route.params.slug
                    && this.sessionService.institutionNodes[0].institutionId === this.sessionService.institution.info.institutionId
                ) {
                    // set current node if not set
                    if (!this.sessionService.currentNode && this.sessionService.institutionNodes) {
                        this.getCurrentNode(this.sessionService.institution);
                    }
                    return resolve(this.sessionService.currentNode);
                } else {

                    // reset when navigating to a new participant/institution
                    this.sessionService.institutionNodes = null;

                    // get institution from service and/or slug
                    const institution: IInstitution = this.sessionService.institution;

                    // set array of nodes for this current institution
                    this.sessionService.institutionNodes = institution.nodes ? Object.values(institution.nodes) : null;

                    // only get list of fully configured/active nodes
                    this.sessionService.institutionNodes = filter(this.sessionService.institutionNodes, (node: INodeAutomation) => {
                        return node.status[0] === 'complete';
                    });

                    // set current node
                    this.getCurrentNode(institution);

                    return resolve(this.sessionService.currentNode);
                }
            }

        });

    }

    /**
     * Gets current node for this institution in this session
     *
     * @param {IInstitution} institution
     * @memberof NodeResolve
     */
    getCurrentNode(institution: IInstitution) {
        // get first available node in array for comparision
        const firstNode = this.sessionService.institutionNodes ? this.sessionService.institutionNodes[0] : null;

        // construct cookie name based on Institution slug (don't expose uuid)
        const cookieName = 'pid-' + institution.info.slug;

        const cachedNodeId = this.utils.getCookie(cookieName);

        let cachedNode: INodeAutomation = null;

        if (this.sessionService.institutionNodes) {
            cachedNode = this.sessionService.institutionNodes.find((node: INodeAutomation) => {
                return node.participantId === cachedNodeId;
            });
        }

        if (cachedNode) {

            // set current node as cached node from cookie
            this.sessionService.setCurrentNode(cachedNode);

        } else {

            // cached node doesn't exist. get first ever node
            this.sessionService.setCurrentNode(firstNode);
        }
    }
}
