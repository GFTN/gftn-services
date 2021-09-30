// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { ApprovalRequest } from "../../../shared/models/approval.interface";

/**
 * Request interface for kill switch request
 *
 * @export
 * @interface KillSwitchRequest
 */
export interface KillSwitchRequest extends ApprovalRequest {

    // temp field to store db key of request from firebase
    // DON'T SAVE TO THE RECORD
    key?: string;

    // participant id
    participantId: string;

    // account address to be suspended/reactivated
    accountAddress: string;

    // email of user who initiated suspend request
    suspendRequestedBy?: string;

    // email of user who approved suspend request
    suspendApprovedBy?: string;

    // email of user who rejected killswitch request
    suspendRejectedBy?: string;

    // email of user who initiated reactivate request
    reactivateRequestedBy?: string;

    // email of user who approved reactivate request
    reactivateApprovedBy?: string;

    // email of user who rejected reactivate request
    reactivateRejectedBy?: string;

    status?: KillSwitchRequestStatus;
}

/**
 * Used only for view to toggle whether or not this request has loaded
 *
 * @export
 * @interface KillSwitchRequestDetail
 * @extends {KillSwitchRequest}
 */
export interface KillSwitchRequestDetail extends KillSwitchRequest {
    loaded?: boolean;
}

export type KillSwitchRequestStatus = 'normal' | 'suspended' | 'suspend_requested' | 'reactivate_requested';
