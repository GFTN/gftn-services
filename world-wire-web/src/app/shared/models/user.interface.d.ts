// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
export interface IUserProfile {

    // NOTE: decided not to flatten because this allows us to put
    // firebase rules on /profile and lock off other inline firebase
    // nodes, namely /participant_permissions /super_permissions
    profile: {
        /**
         * User email (and username)
         *
         * @type {string}
         * @memberof IUser
         */
        email: string;
    };

    // used for /users/participant_permissions/{institutionId}
    participant_permissions?: {
        [institutionId: string]: IParticipantRoles;
    };

    super_permissions?: IRoles;
}

export interface IUserParticipantPermissions {

    // Object Array of roles associated with a institutionId
    // NOTE: This structure allows an super admin to update
    // permissions for a user for multiple participants at once.
    // used for /participant_permissions/{uid}/...
    participant_permissions?: {
        [userId: string]: IParticipantPermissions;
    };

}

export interface IParticipantPermissions {
    [institutionId: string]: IParticipantRoles;
}

export interface IUserSuperPermissions {

    // NOTE: This structure allows an super admin to update
    // permissions for a user for multiple participants at once.
    super_permissions?: ISuperPermissions

}

export interface ISuperPermissions {
    [userId: string]: IRoles
}

export interface IParticipantRoles extends IRoles {
    name?: string;
    slug?: string;
    // email is needed to display users by InstitutionId
}

export interface IRoles {

    // structured with plural roles rather than a single role
    // so that we can scale out more permissions if needed
    roles: IRolesOptions
    // in user management
    email: string;
}

export interface IRolesOptions {
    admin?: boolean;
    manager?: boolean;
    viewer?: boolean;
}

// export interface IUser {

//     /**
//      * Firebase ID/key for identifying user
//      *
//      * @type {string}
//      * @memberof IUser
//      */
//     userId: string; // duplicate of firebase key

//     /**
//      * User email (and username)
//      *
//      * @type {string}
//      * @memberof IUser
//      */
//     email: string;

//     /**
//      * User email (and username)
//      *
//      * @type {string}
//      * @memberof IUser
//      */
//     displayName: string;

//     /**
//      * 2 factor registration
//      * User must setup 2fa once they register for their account
//      *
//      * @type {boolean}
//      * @memberof IUser
//      */
//     registered?: boolean;

// }

export interface ITerms {
    terms: {
        [userId: string]: IUserTerms
    }
}


export interface IUserTerms {
    /**
     * if accepted by user
     *
     * @type {boolean}
     * @memberof IUserTerms
     */
    accepted: boolean;
    /**
     * title of terms accepted
     *
     * @type {string}
     * @memberof IUserTerms
     */
    title: string;
    /**
     * version or date of contract accepted
     *
     * @type {string}
     * @memberof IUserTerms
     */
    version: string;
    /**
     *timestamp of terms accepted
     *
     * @type {string}
     * @memberof IUserTerms
     */
    timestamp: string;
    /**
     * ip address of the user accepting the terms
     *
     * @type {string}
     * @memberof IUserTerms
     */
    ip: string;
}
