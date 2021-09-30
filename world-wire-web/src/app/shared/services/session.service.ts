// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable } from '@angular/core';
import { IInstitution } from '../models/participant.interface';
import { INodeAutomation } from '../models/node.interface';
import { ENVIRONMENT } from '../constants/general.constants';
import { BehaviorSubject } from 'rxjs';
import { UtilsService } from '../utils/utils';

@Injectable()
export class SessionService {

    /**
     * Sets the current participant for the user's session,
     * so that the relevant angular router guard does not have to
     * make a call to the database on every page change.
     * This is also important since our system enables a user
     * to have permissions with multiple participants.
     *
     * @type {IInstitution}
     * @memberof SessionService
     */
    institution: IInstitution;

    /**
     * stores all nodes belonging to the institution.
     * typically should be ONE node for each environment,
     * but supports multiple nodes in the future
     *
     * @type {INodeAutomation[]}
     * @memberof SessionService
     */
    institutionNodes: INodeAutomation[];

    // stores current participant node being viewed
    currentNode: INodeAutomation;

    // received from authentication service to login user to firebase
    firebaseTempAuthToken: string;

    // keeps track of environment change
    public currentNodeChanged: BehaviorSubject<INodeAutomation>;

    constructor(
        private utils: UtilsService
    ) {
        this.currentNodeChanged = new BehaviorSubject<INodeAutomation>(null);
    }

    /**
     * Helper function to synchronize the node environment
     * when the current node is being set in the session
     *
     * @param {INodeAutomation} node
     * @memberof SessionService
     */
    setCurrentNode(node: INodeAutomation) {
        this.currentNode = node;

        // construct cookie name based on Institution slug (don't expose uuid)
        const cookieName = 'pid-' + this.institution.info.slug;

        // cache participant node in cookie for 1 month, if not set
        if (node && (!this.utils.getCookie(cookieName) || this.utils.getCookie(cookieName) !== node.participantId)) {

            const dateString: string = new Date(new Date().getFullYear(), new Date().getMonth() + 1, new Date().getDate()).toUTCString();

            this.utils.setCookie(cookieName, dateString, node.participantId);
        }
    }

    /**
     * Used by other components to notify service
     * that current node has changed
     */
    public propogateNodeChange() {

        this.currentNodeChanged.next(this.currentNode);
    }
}
