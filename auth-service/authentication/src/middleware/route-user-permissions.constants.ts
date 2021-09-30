// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
export interface IRoutePermissions {
    [route: string]: {
        super_permissions?: string[],
        participant_permissions?: string[]
    };
}

// reusable permisisons
const developerPermissions = {
    super_permissions: ['admin', 'manager'],
    participant_permissions: ['admin', 'manager']
};

// set required permissions for routes
// describes firebase route permissions required for continuing to route
export const RoutePermissions: IRoutePermissions = {
    '/permissions/participant': {
        super_permissions: ['admin', 'manager'],
        participant_permissions: ['admin']
    },
    '/permissions/super': {
        super_permissions: ['admin']
    },
    '/jwt/request': developerPermissions,
    '/jwt/revoke': developerPermissions,
    '/jwt/generate': developerPermissions,
    '/jwt/reject': developerPermissions,
    '/jwt/verify': developerPermissions,
    '/jwt/rotate-pepper': {
        super_permissions: ['admin']
    },
    '/jwt/approve': {
        super_permissions: ['admin', 'manager'],
        participant_permissions: ['admin']
    }
};