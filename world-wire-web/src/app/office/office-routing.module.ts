// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { OfficeComponent } from './office.component';
import { AccountsComponent } from './accounts/accounts.component';
import { UsersComponent } from './users/users.component';
import { BlocklistManagementComponent } from './blocklist-management/blocklist-management.component';

// Route Guards and Resolvers
import { AuthCanActivateGuard } from '../shared/guards/auth.guard';
import { SuperPermissionsGuard } from '../shared/guards/super-permissions-redirect.guard';
import { AuthResolver } from '../shared/guards/auth.resolver';
import { UserProfileResolver } from '../shared/guards/user-profile.resolver';

const officeRoutes: Routes = [
    {
        path: '', component: OfficeComponent,
        canActivate: [AuthCanActivateGuard],
        resolve: [AuthResolver, UserProfileResolver],
        children: [
            {
                path: '',
                // component: OverviewComponent,
                // temporary redirect until Overview is finished
                redirectTo: 'accounts',
                pathMatch: 'full',
                data: {
                    title: 'Overview',
                    super_permissions: ['viewer', 'manager', 'admin']
                },
            },
            {
                path: 'account',
                canActivate: [SuperPermissionsGuard],
                data: {
                    title: 'Manage Participant Account',
                    super_permissions: ['manager', 'admin']
                },
                loadChildren: './account/account.module#AccountModule'
            },

            // Controls all actions related to participant accounts
            {
                path: 'accounts',
                component: AccountsComponent,
                canActivate: [SuperPermissionsGuard],
                data: {
                    title: 'Participant Accounts',
                    super_permissions: ['viewer', 'manager', 'admin']
                },
            },

            // modify super user permissions
            {
                path: 'users',
                component: UsersComponent,
                canActivate: [SuperPermissionsGuard],
                data: {
                    title: 'Super Users',
                    super_permissions: ['admin']
                },
            },

            // Blocklist management
            {
                path: 'blocklist',
                component: BlocklistManagementComponent,
                canActivate: [SuperPermissionsGuard],
                data: {
                    title: 'Blocklist Management',
                    super_permissions: ['viewer', 'manager', 'admin']
                },
            },
        ]
    }

];

@NgModule({
    imports: [
        RouterModule.forChild(officeRoutes)
    ],
    exports: [
        RouterModule
    ],
    providers: [
        SuperPermissionsGuard
    ]
})
export class OfficeRoutingModule { }
