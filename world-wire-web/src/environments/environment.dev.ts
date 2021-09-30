// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
export const environment = {
  production: true,
  firebase: {
    apiKey: 'your.api.key.goes.here',
    authDomain: 'your.project.id.firebaseapp.com',
    databaseURL: 'https://your.project.id.firebaseio.com',
    projectId: 'your.project.id',
    storageBucket: '',
    messagingSenderId: 'your.messaging.sender.id',
    appId: 'your.appid.arn'
  },
  apiRootUrl: 'https://auth.worldwire-dev.io',
  supported_env: {
    text: 'Development',
    name: 'dev',
    val: 'eksdev',
    envApiRoot: 'worldwire-dev.io/local/api',
    envGlobalRoot: 'worldwire-dev.io/global',
  },
  inactivityTimeout: 15
};
