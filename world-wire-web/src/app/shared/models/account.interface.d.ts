// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Asset } from "./asset.interface";
import { ApprovalRequest } from "./approval.interface";

export interface AccountRequest extends ApprovalRequest {
    name: string;
}

/**
 * Account returned by the client API
 *
 * @export
 * @interface ParticipantAccount
 */
export interface ParticipantAccount {
    address: string;
    name?: string;
    assets?: Asset[];
}

/**
 * Holds overview details of a participant's
 * operating account for the view only, including:
 * - list of assets
 * - loading state
 *
 * @export
 * @interface ParticipantAccountDetail
 * @extends {ParticipantAccount}
 */
export interface ParticipantAccountDetail extends ParticipantAccount {
    assets?: Asset[];
    loaded?: boolean;

    // OPTIONAL: approval request, if it exists
    request?: AccountRequest;
}
