// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
export interface IChange {

    /**
     * timestamp
     *
     * @type {string}
     * @memberof IChange
     */
    deleted_at: string;

    /**
     * user id
     *
     * @type {string}
     * @memberof IChange
     */
    deleted_by: string;

    /**
     *timestamp
     *
     * @type {string}
     * @memberof IChange
     */
    modified_at: string;

    /**
     *user id
     *
     * @type {string}
     * @memberof IChange
     */
    modified_by: string;

}