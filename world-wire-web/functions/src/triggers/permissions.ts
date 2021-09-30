// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 

import { IUserProfile, IParticipantRoles } from '../../../src/app/shared/models/user.interface';
import { IInstitution, IParticipantUser } from '../../../src/app/shared/models/participant.interface';
import * as functions from 'firebase-functions';
import * as admin from 'firebase-admin';
import { set } from 'lodash';

export const writeParticipantPermissions = (uid: string, institutionId: string, permissions: IParticipantRoles) => {

    const promiseArr = [];

    // get user by uid
    const userPromise: Promise<IUserProfile> = new Promise((resolve, reject) => {

        console.log('userPromise uid:', uid);

        admin.database().ref('/users/' + uid)
            .once('value', (userData: admin.database.DataSnapshot) => {

                // get user profile
                const user: IUserProfile = userData.val();

                if (user) {
                    resolve(user);
                } else {
                    // handle null return
                    console.log('User not returned: ', uid);
                    reject();
                }

            }).then(() => {
                // do nothing
            }, (err: any) => {
                console.error(err);
                reject();
            });

    });

    // get participant by institutionId
    const participantPromise: Promise<IInstitution> = new Promise((resolve, reject) => {

        admin.database().ref('/participants/' + institutionId)
            .once('value', (participantData: admin.database.DataSnapshot) => {

                const participant: IInstitution = participantData.val();

                if (participant) {
                    resolve(participant);
                } else {
                    // handle null return
                    console.log('Participant not returned ', institutionId);
                    reject();
                }

            }).then(() => {
                // do nothing
            }, (err: any) => {
                console.error(err);
                reject();
            });
    });

    // create array of promises to resolve at once
    promiseArr.push(userPromise);
    promiseArr.push(participantPromise);

    // for performance purposes execute both promises at the same time
    return new Promise((resolve, reject) => {
        Promise.all(promiseArr).then((results) => {

            const user: IUserProfile = results[0];
            const participant: IInstitution = results[1];

            if (participant && user) {

                // set the participant name in the user node for readability in ui
                let participantWithName = set(permissions, 'name', participant.info.name);

                // set the participant slug in the user node for readability in ui
                participantWithName = set(participantWithName, 'slug', participant.info.slug);

                // add participant to user node to display the user's permissions for participants
                try {
                    admin.database().ref('/users/' + uid + '/participant_permissions/' + institutionId)
                        .update(participantWithName)
                        .then(() => {
                            // do nothing
                        }, (err: any) => {
                            console.error(err);
                        });

                } catch (err) {
                    console.error(err);
                    reject(err);
                }

                // add user to the participant node to track users and their permissions for a participant
                const userWithPermissions: IParticipantUser = {
                    roles: permissions.roles,
                    profile: user.profile
                };

                try {
                    admin.database().ref('/participants/' + institutionId + '/users/' + uid)
                        .update(userWithPermissions)
                        .then(() => {
                            // do nothing
                        }, (err: any) => {
                            console.error(err);
                        });
                } catch (err) {
                    console.error(err);
                    reject(err);
                }

                resolve(true);
            } else {

                // reject if participant or user is not set
                reject();
            }
        }, (err: any) => {
            console.error(err);
            reject(false);
        });
    });
};

// 1) copy - add user 'permissions' by userId to /users/{uid}
// 2) copy - add user 'permissions' by userId to /participant/{institutionId}
// TEST: updateParticipantPermissions({ after: { manager: true } }, { params: { institutionId: '-LKZvz9-SA9-YaYRL9Mb', uid:'ZdInKi2VNkddBWVkSWfCNwjAWF52'} })
export const updateParticipantPermissions = functions.database.ref('/participant_permissions/{uid}/{institutionId}')
    .onUpdate((snapshot: functions.Change<functions.database.DataSnapshot>, context: functions.EventContext) => {
        const data: IParticipantRoles = snapshot.after.val();
        return writeParticipantPermissions(context.params.uid, context.params.institutionId, data);
    });

// 1) copy - add user 'permissions' by userId to /users/{uid}
// 2) copy - add user 'permissions' by userId to /participant/{institutionId}
export const createParticipantPermissions = functions.database.ref('/participant_permissions/{uid}/{institutionId}')
    .onCreate((snapshot: functions.database.DataSnapshot, context: functions.EventContext) => {
        const data: IParticipantRoles = snapshot.val();
        return writeParticipantPermissions(context.params.uid, context.params.institutionId, data);
    });

// 1) copy - remove user 'permissions' by userId to /users/{uid}
// 2) copy - remove user 'permissions' by userId to /participant/{institutionId}
// TEST: removeParticipantPermissions({ manager: true }, { params: { institutionId: '-LOtWnngmoIu8otdvsUk', uid: 'ydU4CpByomPNq5kjHsx32wwUuU12' } })
export const removeParticipantPermissions = functions.database.ref('/participant_permissions/{uid}/{institutionId}')
    .onDelete((snapshot: functions.database.DataSnapshot, context: functions.EventContext) => {

        return new Promise((resolve, reject) => {

            const promiseArr = [];
            // add participant to user node to display the user's permissions for participants
            promiseArr.push(admin.database().ref('/users/' + context.params.uid + '/participant_permissions/' + context.params.institutionId).remove());
            // add user to the participant node to track users and their permissions for a participant
            promiseArr.push(admin.database().ref('/participants/' + context.params.institutionId + '/users/' + context.params.uid).remove());


            Promise.all(promiseArr).then((results) => {
                resolve(true);
            }).catch((err) => {
                reject(err);
            });
        });
    });

// copy - add super user 'permissions' by userId to /users/{uid}
export const updateSuperPermissions = functions.database.ref('/super_permissions/{uid}')
    .onUpdate((snapshot: functions.Change<functions.database.DataSnapshot>, context) => {
        // add participant to user node to display the user's permissions for participants
        admin.database().ref('/users/' + context.params.uid + '/super_permissions/').update(snapshot.after.val());
        return false;
    });

// copy - add super user 'permissions' by userId to /users/{uid}
export const createSuperPermissions = functions.database.ref('/super_permissions/{uid}')
    .onCreate((snapshot: functions.database.DataSnapshot, context) => {
        // add participant to user node to display the user's permissions for participants
        admin.database().ref('/users/' + context.params.uid + '/super_permissions/').update(snapshot.val());
        return false;

    });

// copy - remove super user 'permissions' by userId to /users/{uid}
export const removeSuperPermissions = functions.database.ref('/super_permissions/{uid}')
    .onDelete((snapshot: functions.database.DataSnapshot, context) => {
        // add participant to user node to display the user's permissions for participants
        admin.database().ref('/users/' + context.params.uid + '/super_permissions/').remove();
        return false;
    });
