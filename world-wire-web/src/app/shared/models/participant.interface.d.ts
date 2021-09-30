// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { INodeAutomation } from "./node.interface";
import { IRolesOptions, IUserProfile } from "./user.interface";
import { Asset } from "./asset.interface";
import { ParticipantAccount } from "./account.interface";


export interface IOrganization extends IOrganizationMeta {

    // Company Name
    name: string;

    // Company Short Name
    // short?: string;

    // GeoPoint
    geo_lat?: number;
    geo_lon?: number;

    // head quarters
    country: string; // use ISO 3166-1 alpha-3 codes
    address1: string;
    address2: string;
    city?: string;
    state?: string;
    zip?: string;

    logo_url?: string;
    site_url?: string;

    // the type of entity this is.
    // current types are only club or system.
    kind: 'Money Transfer Operator' | 'Financial Institution' | 'Bank' | 'Central Bank' | 'Credit Union' | 'Other';

}

interface IOrganizationMeta {

    // NOTE: 'value', 'slug', 'type', and 'active' all need to be in the first level
    // of this object so that firebase.orderByChild() search can be performed on all entities

    // firebase id
    institutionId: string;

    // readable unique identifier for url only lower characters, no spaces
    // used to query club from office component and part of office URL
    slug: string;

    // phase: development phase
    // phase: 'dev' | 'qa' | 'st' | 'tn' | 'prod';
    // TODO : auto deploy participant sandbox so that their
    // sandbox usage doesn't collide with others
    // sandbox: root url for personal sandbox use
    // sandboxConfig: string;

    // status: state of account if account is suspended or active
    // pending = prior to prod-net activation
    // active = prod-net account activated (whitelisted in PR)
    // suspended = suspended account (ie: blacklisted from PR)
    status: 'pending' | 'active' | 'suspended';

    // admins associated with account
    // admins: IUser[];

    // Includes all accounts (ie: Members, guests, and staff) associated with this club
    // accounts: {
    //     [uid: string]: IEntityAccount
    // };

    // TODO: Determine the best way to store corridors - perhaps use WW API to query
    // corridors_receive: {
    //     [id: string]: string;
    // };

    // corridors_send: {
    //     [id: string]: string;
    // };

    // TODO: Need to determine if there needs to be some type of "Billing Plan"
    // associated with the entity
    // keeping 'plan' separate from info because this
    // would likely be updated in firebase on it's own
    // plan: IPlan;
}

/**
 * List of users and permissions for that user
 * NOTE: plural 'users'
 *
 * @export
 * @interface IParticipantUsers
 */
export interface IParticipantUsers {
    users?: {
        [userId: string]: IParticipantUser;
    };
}

/**
 * NOTE: singular 'user'
 *
 * @export
 * @interface IParticipantUser
 * @extends {IUserProfile}
 */
export interface IParticipantUser extends IUserProfile {
    // roles are used in view to for route guards
    roles: IRolesOptions
}

export interface IInstitution extends IParticipantUsers {
    info: IOrganization;
    // Institutions may or may not have nodes yet
    nodes?: INodeAutomation[];
}

/**
 * Particpant returned by the client API
 *
 * @export
 * @interface Participant
 */
export interface Participant {
    bic: string;
    country_code: string;
    id: string;
    issuing_account?: string;
    operating_accounts?: ParticipantAccount[];
    // Role of the registered participant on the network
    // 'MM' = Market Maker/regular participant, issues only DOs
    // 'IS' = Issuer of real world DAs/stable coins in addition to DOs
    role: ParticipantRole;
    status: string;
}

export type ParticipantRole = 'MM' | 'IS';
