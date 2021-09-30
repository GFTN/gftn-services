// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import * as functions from 'firebase-functions';
import * as admin from 'firebase-admin';
import { IInstitution } from '../../../src/app/shared/models/participant.interface';
import { IParticipantPermissions } from '../../../src/app/shared/models/user.interface';

const writeSlug = async (institution: IInstitution): Promise<void> => {

    // update new slug
    admin.database().ref('slugs')
        .update({ [institution.info.slug]: institution.info.institutionId });

};

const writeNodeSlug = async (participantId: string, institutionId: string): Promise<void> => {
    // updates participant node to the institution it belongs tos
    admin.database().ref('nodes').update({
        [participantId]: institutionId
    });
};

// updates all permissions for users affected by the slug change
const updatePermissionSlug = async (institution: IInstitution): Promise<void> => {

    // get all users
    const data = await admin.database().ref('participant_permissions').once('value');

    const participant_permissions: IParticipantPermissions = data.val();

    // update the slug in each participant's affected profile
    for (const [key, value] of Object.entries(participant_permissions)) {

        // filter users by those who have permissions for this institution
        const permissionedInstitutions: string[] = Object.keys(value);

        const institutionFound = permissionedInstitutions.includes(institution.info.institutionId);

        // only update users that have permissions for the changed institution
        if (institutionFound) {
            // update specific user for participant's info (ie: slug)
            await admin.database().ref('users')
                .child(key)
                .child('participant_permissions')
                .child(institution.info.institutionId)
                .update({
                    slug: institution.info.slug,
                    name: institution.info.name
                });
        }
    }
};

// run locally: // updateSlugMap({ after: {info: {slug: 'hong-kong-bank-1', institutionId: '-LPmrF0Qgg6fshzLc48i'}}, before: {info: {slug: 'hong-kong-bank', institutionId: '-LPmrF0Qgg6fshzLc48i'}} })
export const updateSlugMap = functions.database.ref('/participants/{institutionId}')
    .onUpdate(async (snapshot: functions.Change<functions.database.DataSnapshot>) => {

        const afterInstitution: IInstitution = snapshot.after.val();
        const beforeInstitution: IInstitution = snapshot.before.val();

        // console.log(institution);
        // create a human readable slug for the participant in the ui
        // needed for querying see - app/shared/guards/permissions-redirect.guard.ts

        // update slugs if changed
        if (beforeInstitution.info.slug !== afterInstitution.info.slug) {

            // delete out old slug
            if (beforeInstitution.info.slug) {
                admin.database().ref('slugs').child(beforeInstitution.info.slug)
                    .remove();
            }

            await writeSlug(afterInstitution);

            await updatePermissionSlug(afterInstitution);
        }

        return false;

    });

export const createSlugMap = functions.database.ref('/participants/{institutionId}')
    .onCreate(async (snapshot: functions.database.DataSnapshot) => {

        const institution: IInstitution = snapshot.val();

        await writeSlug(institution);

        return false;

    });

// maps new participant node to the institution it belongs to upon creation
export const createNodeMap = functions.database.ref('/participants/{institutionId}/nodes/{participantId}')
    .onCreate(async (snapshot: functions.database.DataSnapshot, context) => {

        await writeNodeSlug(context.params.participantId, context.params.institutionId);

        return false;

    });

// remove mapping if node is deleted from the list of nodes for an institution.
// Node mappings can only ever be created or deleted,
// since update to the participant ID cannot actually be done in the PR at the moment.
export const removeNodeMap = functions.database.ref('/participants/{institutionId}/nodes/{participantId}')
    .onDelete(async (snapshot: functions.database.DataSnapshot, context) => {

        await writeNodeSlug(context.params.participantId, null);

        return false;

    });
