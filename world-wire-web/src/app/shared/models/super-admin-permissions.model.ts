// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import * as _ from 'lodash';
import { Injectable } from '@angular/core';
import { AngularFireDatabase } from '@angular/fire/database';

@Injectable()
export class SuperAdminPermissionsModel {

    /**
     * Used to determine if a user is a super admin
     * Security: Is a SECURE node for permissions
     * @memberof SuperAdminPermissionsModel
     */
    route = 'super_permissions/{userId}/role';

    constructor(
        private db: AngularFireDatabase
    ) { }

    /**
     * Add user permissions.
     * NOTE: permissions are added and removed one at at time.
     *
     * @returns {Promise<void>}
     * @memberof SuperAdminPermissionsModel
     */
    add(
        userId: string,
        roleType: 'admin' | 'manager'
    ): Promise<void> {
        // return promise since result is a single success or failure
        return new Promise((resolve, reject) => {

            // update record in firebase
            this.db.database.ref(
                'super_permissions/' +
                userId +
                '/role'
            ).update({ [roleType]: true })
                .then(() => {
                    resolve();
                }, (error) => {
                    console.log('Error: Unable to add super admin permissions.', error);
                    alert('Error: Unable to add super admin permissions.');
                    reject();
                });

        });

    }

    /**
     * Remove user permissions.
     * NOTE: permissions are added and removed one at at time.
     *
     * @returns {Promise<void>}
     * @memberof SuperAdminPermissionsModel
     */
    remove(
        userId: string,
        roleType: 'admin' | 'manager'
    ): Promise<void> {

        // return promise since result is a single success or failure
        return new Promise((resolve, reject) => {

            // update record in firebase
            this.db.database.ref(
                'super_permissions/' +
                userId +
                '/role/' +
                roleType
            ).remove()
                .then(() => {
                    resolve();
                }, (error) => {
                    console.log('Error: Unable to remove super admin permissions.', error);
                    alert('Error: Unable to remove super admin permissions.');
                    reject();
                });

        });

    }

}


