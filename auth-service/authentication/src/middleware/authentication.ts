// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// see usage docs here https://github.com/lukeautry/tsoa#security
// See routes.ts for implementation of the expressAuthentication() function

import * as express from 'express';

export function expressAuthentication(request: express.Request, securityName: string, scopes?: string[]): Promise<any> {

    // TODO: move all middleware checks here from auth.middleware.ts

    if (securityName === 'api_key') {

        return new Promise((resolve) => {
            // Google Endpoints uses ESP to check permissions of the API Key
            // Therefore, if the request reaches this point in the code the 
            // default response is to allow the request to continue through 
            // remaining middleware.
            return resolve();
        });

    }

    // default response to reject if not previously resolved
    return Promise.reject({});

};