// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
export const environment = {
  production: true,
  firebase: {
    projectId: 'your.project.id',
    appId: 'your.appid.arn',
    databaseURL: 'https://your.project.id.firebaseio.com',
    storageBucket: 'your.project.id.appspot.com',
    apiKey: 'your.api.key.goes.here',
    authDomain: 'your.project.id.firebaseapp.com',
    messagingSenderId: 'your.messaging.sender.id'
  },
  apiRootUrl: 'https://auth.worldwire.io',
  supported_env: {
    text: 'Production',
    name: 'prod',
    val: 'prod',
    envApiRoot: 'worldwire.io/local/api',
    envGlobalRoot: 'worldwire.io/global',
  },
  inactivityTimeout: 15
};
