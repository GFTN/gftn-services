// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable } from '@angular/core';
import { find, toArray, keys, max } from 'lodash';
import { VERSION_DETAILS, IWWApi, IWWApis } from '../../../shared/constants/versions.constant';

/**
 * Used to set the current version selected in the ui (ie: /docs)
 * Shared across the different versions of the API to maintain
 * the current state/version of the API being viewed
 *
 * @export
 * @class VersionService
 */
@Injectable()
export class VersionService {

    /**
     * currently selected api version details (per constant)
     *
     * @type {IWWApi}
     * @memberof VersionService
     */
    current: IWWApi;
    versions: IWWApis;

    constructor() {

        this.versions = VERSION_DETAILS;
    }

    /**
     * Get version details by version name (ie: look up version details from url)
     *
     * @param {string} version
     * @memberof versionService
     */
    getVersion(version: string) {

        return find(toArray(VERSION_DETAILS), (details: IWWApi) => {
            return details.version === version;
        });

    }

    /**
     * Sets the api version to the newest version in the sequence
     *
     * @memberof versionService
     */
    getNewestVersion(): IWWApi {
        // assume module ids are number in sequential order
        // starting at 1 (largest number would be the most recent)
        const newestModuleId = max(
            // get array of key values
            keys(VERSION_DETAILS)
        );

        return VERSION_DETAILS[newestModuleId];

    }

}
