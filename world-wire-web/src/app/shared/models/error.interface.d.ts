// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
export interface WorldWireError {
    build_version: string;
    msg?: string;
    message?: string;
    details?: string;
    participant_id: string;
    service: string;
    time_stamp: number;
    url: string;
}
