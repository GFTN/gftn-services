// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable } from '@angular/core';
import { AngularFireDatabase } from '@angular/fire/database';
import { ApprovalInfo, Approval } from './approval.interface';
import { IRoles } from './user.interface';

@Injectable()
export class SuperApprovalsModel {

    constructor(
        private db: AngularFireDatabase
    ) { }

    /**
     * Get super_approval info based on approval Id
     *
     * @param {string} approvalId
     * @returns {Promise<ApprovalInfo>}
     * @memberof SuperApprovalsModel
     */
    async getApprovalInfo(approvalId: string): Promise<ApprovalInfo> {

        const data = await this.db.database.ref('/super_approvals').child(approvalId).once('value');

        // get users ref for this institution
        const dbUsersRef = 'super_permissions';

        const approval: Approval = data.val() ? data.val() : null;

        const approvalInfo: ApprovalInfo = {
            key: approvalId,
        };

        if (approval) {
            // get user emails of uids
            const userRequests: Promise<any>[] = [];

            if (approval.uid_request) {
                userRequests.push(this.db.database.ref(dbUsersRef).child(approval.uid_request).once('value'));
            }

            if (approval.uid_approve) {
                userRequests.push(this.db.database.ref(dbUsersRef).child(approval.uid_approve).once('value'));
            }
            const getUsers = await Promise.all(userRequests);

            const users: IRoles[] = [];

            for (const userData of getUsers) {
                if (userData.val()) {
                    users.push(userData.val());
                }
            }

            approvalInfo.requestInitiatedBy = users ? users[0].email : null;
            approvalInfo.requestApprovedBy = users.length > 1 ? users[1].email : null;
        }

        return approvalInfo;

    }

    /**
     * Resets approval in case of error
     *
     * @param {string} approvalId
     * @returns {Promise<any>}
     * @memberof SuperApprovalsModel
     */
    async resetApprovals(approvalId: string): Promise<any> {

        const updateFields = {
            status: 'request',
            uid_approve: '',
        };

        // reset approval Id
        return await this.db.database.ref('super_approvals')
            .child(approvalId)
            .update(updateFields);
    }
}
