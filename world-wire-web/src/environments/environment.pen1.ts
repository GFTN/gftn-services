// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
export const environment = {
    production: true,
    firebase: // Copy and paste this into your JavaScript code to initialize the Firebase SDK.
    // You will also need to load the Firebase SDK.
    // See https://firebase.google.com/docs/web/setup for more details.
    {
        projectId: 'your.project.id',
        appId: 'your.appid.arn',
        databaseURL: 'https://your.project.id.firebaseio.com',
        storageBucket: 'your.project.id.appspot.com',
        apiKey: 'your.api.key.goes.here',
        authDomain: 'your.project.id.firebaseapp.com',
        messagingSenderId: 'your.messaging.sender.id'
    },
    apiRootUrl: 'https://auth.worldwire-pen.io',
    supported_env: {
        text: 'Pentesting1',
        val: 'pen1',
        envApiRoot: 'worldwire-pen.io/local/api',
        envGlobalRoot: 'worldwire-pen.io/global',
    },
    inactivityTimeout: 15
};
