// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
/**
 * Request for trusting an asset
 *
 * @export
 * @interface TrustRequest
 */
export interface TrustRequest {

    // temp field to store db key of request from firebase
    // DON'T SAVE TO THE RECORD
    key?: string;

    // participant id of the participant requesting trust for an asset
    requestor_id: string;

    // participant id of the asset issuer
    issuer_id: string;

    // name of the operating account requesting to hold the asset
    account_name: string;

    // code of the asset being requested
    asset_code: string;

    // max amount requested to be held in an operating account
    // necessary for permission=request
    limit: number;

    status?: TrustRequestStatus;

    // Unix timestamp of last status update
    time_updated: number;

    approval_ids: string[];

    // OPTIONAL: reason why the request was rejected
    reason_rejected?: string;

    // uuid of the requestor user who initiated
    // the 'request' action on the trustline
    requestInitiatedBy?: string;

    // uuid of the requestor's user who approved
    // the 'request' action on the trustline
    requestApprovedBy?: string;

    // uuid of the issuer's user who initiated
    // the 'allow' action on the trustline
    allowInitiatedBy?: string;

    // uuid of the issuer's user who approved
    // the 'allow' action on the trustline
    allowApprovedBy?: string;

    // uuid of the issuer's user who initiated
    // the 'reject' (ignores the request) action on the trustline
    rejectInitiatedBy?: string;

    // uuid of the issuer's user who approved
    // the 'reject' (ignores the request) action on the trustline
    rejectApprovedBy?: string;

    // uuid of the issuer's user who initiated
    // the 'revoke' action on the trustline
    revokeInitiatedBy?: string;

    // uuid of the issuer's user who approved
    // the 'revoke' action on the trustline
    revokeApprovedBy?: string;

    loaded?: boolean;
}

export type TrustRequestPermission = 'request' | 'allow' | 'revoke';

export type TrustRequestStatus =

    // 'REQUEST' action: initiated = pending, requested = approved
    'initiated' | 'requested' |

    // 'ALLOW' action: allowed = pending, approved = approved (final stage)
    'allowed' | 'approved' |

    // 'REJECT' action: rejectPending = pending, rejected = approved
    'rejectPending' | 'rejected' |

    // 'REVOKE' action: revokePending = pending, revoked = approved (final stage)
    'revokePending' | 'revoked';
