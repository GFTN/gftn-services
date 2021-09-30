// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// TODO: Currently a work-in-progress - not being used. 

// PURPOSE: To manage institutions. Institutions represent a top-level account
// that can be used to group all resources related to that institutions such as 
// participant stacks (ie: world wire participant deployment on kubernetes), and 
// participant_users.

import {
    Route, Controller, Post, OperationId
    , Security
} from 'tsoa';
// import { google } from 'googleapis';
// import { OAuth2Client } from 'googleapis-common';
// import { JWT } from 'google-auth-library';
// import { Compute } from 'google-auth-library';
// import { UserRefreshClient } from 'google-auth-library';

@Route('institution')
export class InstitutionController extends Controller {
    
    @Post('create')
    @OperationId('createInstitution')
    @Security('api_key', [])
    public async create() {

        // this.newApiKey();
        return 'test';

    }

    public async newApiKey() {

        // BEFORE RUNNING:
        // ---------------
        // 1. If not already done, enable the Identity and Access Management (IAM) API
        //    and check the quota for your project at
        //    https://console.developers.google.com/apis/api/iam
        // 2. This sample uses Application Default Credentials for authentication.
        //    If not already done, install the gcloud CLI from
        //    https://cloud.google.com/sdk and run
        //    `gcloud beta auth application-default login`.
        //    For more information, see
        //    https://developers.google.com/identity/protocols/application-default-credentials
        // 3. Install the Node.js client library by running
        //    `npm install googleapis --save`

        // const iam = google.iam('v1');
        // const projectId = 'chase-endpoints';
        // const serviceAccount = 'chase-endpoints@appspot.gserviceaccount.com';

        // const authorize = (callback) => {

        //     google.auth.getApplicationDefault((err, authClient: OAuth2Client | any) => {

        //         if (err) {
        //             console.error('authentication failed: ', err);
        //             return;
        //         }

        //         if (authClient.createScopedRequired && authClient.createScopedRequired()) {
        //             const scopes = ['https://www.googleapis.com/auth/cloud-platform'];
        //             authClient = authClient.createScoped(scopes);
        //         }

        //         callback(authClient);

        //     });

    }

    //     authorize((authClient: string | OAuth2Client | JWT | Compute | UserRefreshClient) => {

    //         const request = {

    //             /**
    //              * The resource name of the service account in the following format:
    //              * `projects/{PROJECT_ID}/serviceAccounts/{ACCOUNT}`. Using `-` as a
    //              * wildcard for the `PROJECT_ID` will infer the project from the account.
    //              * The `ACCOUNT` value can be the `email` address or the `unique_id` of the
    //              * service account.
    //              */
    //             name: `projects/${projectId}/serviceAccounts/${serviceAccount}`,

    //             resource: {
    //                 // TODO: Add desired properties to the request body.
    //             },

    //             auth: authClient

    //         } as any;

    //         // create a new service account key
    //         iam.projects.serviceAccounts.keys.create(request, (err, response) => {

    //             if (err) {
    //                 console.error(err);
    //                 return;
    //             }

    //             // TODO: Change code below to process the `response` object:
    //             console.log(JSON.stringify(response, null, 2));

    //         });

    //     });

    // }

}
