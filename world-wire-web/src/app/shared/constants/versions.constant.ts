// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
//
import * as _ from 'lodash';

export interface IWWApi {
    // ww release version
    version: string;
    // name of module mapped in angular (increment number value)
    module: string;
    // Github release tag version
    releaseTag: string;

    // names of openAPI (swagger) files
    config: {
        name: string;
        // show sandbox
        sandbox: boolean;
        // API endpoints are hosted on World Wire
        wwHosted: boolean;
        // true if public facing endpoint or false if internally facing endpoint (like onboarding endpoints which needs a token)
        public: boolean;
    }[];

    // paths: IPaths[]
    paths?: {
        [apiFileUrl: string]: {
            [mdRouteUrl: string]: string
        }
    };
}

export interface IWWApis {
    [moduleName: string]: IWWApi;
}

// Available Release versions corresponding to release at:
// https://github.com/GFTN/api-service/releases
// NOTE: creating as object rather than array
// so that module details can been keyed off of
export const VERSION_DETAILS: IWWApis = {

    1: {
        module: '1',
        version: 'v2.11.3',
        releaseTag: '2.11.3',
        config: [
            {
                name: 'admin-api',
                sandbox: false,
                wwHosted: true,
                public: false
            },
            {
                name: 'anchor-api',
                sandbox: false,
                wwHosted: true,
                public: true
            },
            //     AT SOME POINT WHEN WE GET ANCHORS GOING, WE'LL NEED THIS
            // {
            //     name: 'anchor-onboarding-api',
            //     sandbox: false,
            //     wwHosted: true,
            //     public: false
            // }
            {
                name: 'participant-api',
                sandbox: false,
                wwHosted: true,
                public: true
            },
            {
                name: 'participant-onboarding-api',
                sandbox: false,
                wwHosted: true,
                public: false
            }
        ]
    }
};
