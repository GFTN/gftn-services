// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import * as admin from 'firebase-admin';
import * as permissions from './triggers/permissions';
import * as institution from './triggers/institution';

admin.initializeApp();

export const updateParticipantPermissions = permissions.updateParticipantPermissions;
export const createParticipantPermissions = permissions.createParticipantPermissions;
export const removeParticipantPermissions = permissions.removeParticipantPermissions;
export const updateSuperPermissions = permissions.updateSuperPermissions;
export const createSuperPermissions = permissions.createSuperPermissions;
export const removeSuperPermissions = permissions.removeSuperPermissions;
export const updateSlugMap = institution.updateSlugMap;
export const createSlugMap = institution.createSlugMap;
export const createNodeMap = institution.createNodeMap;
export const removeNodeMap = institution.removeNodeMap;

// // Test function to run locally:
// import * as functions from 'firebase-functions';
// // Run locally https://firebase.google.com/docs/functions/local-emulator#invoke_realtime_database_functions
// // ** One time Setup Step: ** $ set GOOGLE_APPLICATION_CREDENTIALS=path\to\key.json
// // ** Step 1: ** $ tsc -p functions/tsconfig.dev.json ; firebase functions:shell
// // ** Step 2: ** $ test({before: 'old_data', after: 'new_data' })
// export const test = functions.database.ref('/test').onUpdate(async (snapshot: functions.Change<functions.database.DataSnapshot>, context: functions.EventContext) => {
//     await setTimeout(() => {
//     const c = context;
//     const s = snapshot.after.val();
//     console.log('this here is your context:', c);
//     console.log('here is your snapshot:', s);
//         return false;
//     }, 1500);
// });
