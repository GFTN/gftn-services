// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import * as fs from 'fs';
import * as BoxSDK from 'box-node-sdk';

export class WWBox {

    private client: any;

    async init() {

        const sdkConfig = await fs.readFileSync('.credentials/box/admin.json', "utf8")
        // const sdk = BoxSDK.getPreconfiguredInstance(sdkConfig);
        const sdkConfigObj = JSON.parse(sdkConfig);

        const sdk = new BoxSDK(sdkConfigObj.boxAppSettings);

        // Get the service account client, used to create and manage app user accounts
        // The enterprise ID is pre-populated by the JSON configuration,
        // so you don't need to specify it here
        // var serviceAccountClient = sdk.getAppAuthClient('enterprise');
        this.client = await sdk.getAppAuthClient('enterprise', sdkConfigObj.enterpriseID);

        // // Get an app user client
        // var appUserClient = sdk.getAppAuthClient('user', 'YOUR-APP-USER-ID');

    }

    async upload() {
        var stream = fs.createReadStream('README.md');
        await this.client.files.uploadFile('78644084284', 'chasetestfile', stream);
    }

}