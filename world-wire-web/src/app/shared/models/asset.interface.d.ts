// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { ApprovalRequest } from "./approval.interface";

/**
 * Asset Model returned by the client API
 *
 * @export
 * @interface IAsset
 */
export interface Asset {

    // 3-letter ISO code of the asset
    // (DOs have 5, with 'DO' appended)
    asset_code: string;

    // type of asset:
    // DO - digital obligation (owed amount)
    // DA - stable asset on WW backed by fiat
    asset_type: AssetType;

    /**
     * OPTIONAL fields used in the view
     * to consolidate asset details
     */
    // Issuer of the asset
    issuer_id?: string;

    // current overall balance of the asset
    balance?: number;

    // 3-letter ISO code that this asset pertains to in the real world
    currency?: string;
}

export type AssetType = 'DO' | 'DA';

export interface AssetBalance {
    account_name: string;
    asset_code: string;
    balance: string;
    issuer_id?: string;
}

/**
 * Stores details about a DO balance
 * by each partipant owed
 *
 * @export
 * @interface Obligation
 */
export interface Obligation {

    // ID of participant who the obligation is owed
    participant_id: string;

    // amount/balance that is owed
    balance: AssetBalance;
}

export interface AssetRequest extends Asset, ApprovalRequest {

    status?: AssetRequestStatus;
}

export type AssetRequestStatus = 'requested' | 'rejected' | 'approved';
