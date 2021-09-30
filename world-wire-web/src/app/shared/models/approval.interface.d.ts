// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
export interface Approval {
    uid_request: string;
    uid_approve: string;
    iid: string;
    status: 'request' | 'approved';
    timestamp_request: number;
    endpoint: string;
    method: string;
}

export interface ApprovalInfo {
    key?: string;
    requestInitiatedBy?: string;
    requestApprovedBy?: string;
}


/**
 * Base request for approvals
 *
 * @export
 * @interface ApprovalRequest
 */
export interface ApprovalRequest {
    // UNIX timestamp of the last action for this request
    timeUpdated?: number;

    // list of ids of the maker/checker request
    approvalIds: string[];
}

export type ApprovalPermission = 'request' | 'approve';
