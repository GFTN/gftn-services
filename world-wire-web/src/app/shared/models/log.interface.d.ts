// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
/**
 * Main Interface for
 * Log Object
 *  **/
export interface Log {
    code: string;
    details: string;
    message: string;
    participant_id: string;
    time_stamp: string;
    type: string;
    url: string
}

/**
 * Stores error status and history
 * of errors by date
 */
export interface StatusByDate {
  date: string;
  status: number;
  errors?: Log[];
}

/**
 * Current running service for this particular
 * version of the API
 */
export interface Service {
  name: string;
  displayName: string;
  url: string;
  errorHistory: StatusByDate[];
}
