// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { ApprovalRequest } from "../../../shared/models/approval.interface";

/**
 * Request interface for adding a participant to one's own whitelist
 *
 * @export
 * @interface WhitelistRequest
 */
export interface WhitelistRequest extends ApprovalRequest {

    // temp field to store db key of request from firebase
    // DON'T SAVE TO THE RECORD
    key?: string;

    // participantId of initial whitelister
    whitelisterId: string;

    // id of participant being whitelisted
    whitelistedId: string;

    // email of user who initiated the 'add' action for a whitelist request
    requestedBy?: string;

    // email of user who initiated the 'add' action for a whitelist request
    approvedBy?: string;

    // email of user who rejected the whitelist request
    rejectedBy?: string;

    // email of user who initiated the 'delete' action for a whitelist request
    deleteRequestedBy?: string;

    // email of user who approved the 'delete' action for a whitelist request
    deleteApprovedBy?: string;

    status?: WhitelistRequestStatus;
}

export type WhitelistRequestStatus = 'pending' | 'approved' | 'rejected' | 'deleted';
