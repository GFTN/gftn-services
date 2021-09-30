// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { ApprovalInfo } from "./approval.interface";

/**
 * Blocklist response returned back from the Admin API
 *
 * @export
 * @interface Blocklist
 */
export interface Blocklist {
    // The id of the blocked element
    id?: string;

    // The name of the blocked element
    name?: string;

    // The type of the blocklist element
    // Required: true
    // Enum: [CURRENCY COUNTRY INSTITUTION]
    type: BlocklistType;

    // The value of the block type
    // Required: true
    value: string[];
}

export interface BlocklistRequest {

    // type of blocklisted element
    type: BlocklistType;

    // value of the blocklisted element
    value: string;

    approvalIds?: string[];
}

export type BlocklistType = 'currency' | 'country' | 'institution';
