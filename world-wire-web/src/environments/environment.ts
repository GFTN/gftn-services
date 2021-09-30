// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// The file contents for the current environment will overwrite these during build.
// The build system defaults to the dev environment which uses `environment.ts`, but if you do
// `ng build --env=prod` then `environment.prod.ts` will be used instead.
// The list of which env maps to which file can be found in `.angular-cli.json`.

// Next site
export const environment = {
  production: false,
  firebase: {
    apiKey: 'your.api.key.goes.here',
    authDomain: 'your.project.id.firebaseapp.com',
    databaseURL: 'https://your.project.id.firebaseio.com',
    projectId: 'your.project.id',
    storageBucket: '',
    messagingSenderId: 'your.messaging.sender.id',
    appId: 'your.appid.arn'
  },
  apiRootUrl: 'https://localhost:6001',
  supported_env: {
    text: 'Development',
    name: 'dev',
    val: 'eksdev',
    envApiRoot: 'worldwire-dev.io/local/api',
    envGlobalRoot: 'worldwire-dev.io/global',
  },
  inactivityTimeout: 100 // increase so that view doesn't timeout during development
};
