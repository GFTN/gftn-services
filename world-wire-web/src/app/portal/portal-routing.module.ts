// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { PortalComponent } from './portal.component';
import { AccountComponent } from './account/account.component';
import { OverviewComponent } from './overview/overview.component';
import { StatusOverviewComponent } from './status/status-overview/status-overview.component';
import { StatusHistoryComponent } from './status/status-history/status-history.component';
import { DiligenceComponent } from './diligence/diligence.component';
import { UsersComponent } from './users/users.component';
import { BillingComponent } from './billing/billing.component';
import { TransactionsComponent } from './transactions/transactions.component';
import { SettingsComponent } from './settings/settings.component';
import { AuthCanActivateGuard } from '../shared/guards/auth.guard';
import { ParticipantPermissionsGuard } from '../shared/guards/permissions-redirect.guard';
import { ParticipantResolve } from '../shared/guards/participant.resolve';
import { NodeResolve } from '../shared/guards/node.resolve';
import { TokenComponent } from './token/token.component';
import { SelectComponent } from './select/select.component';
import { UserProfileResolver } from '../shared/guards/user-profile.resolver';
import { AuthResolver } from '../shared/guards/auth.resolver';
import { AccountsComponent } from './accounts/accounts.component';
import { AssetsComponent } from './assets/assets.component';


const portalRoutes: Routes = [
    // {
    //     // this route is needed to check redirect
    //     // if no participant is provides since a
    //     // redirect auth cookie may be present
    //     path: '',
    //     component: PortalComponent,
    //     canActivate: [AuthCanActivateGuard],
    //     data: {
    //         redirectBase: '/portal',
    //         redirectUnauthorized: '/unauthorized'
    //     },
    // },
    {
        path: '',
        redirectTo: 'select',
        pathMatch: 'full'
    },
    {
        path: 'select',
        component: PortalComponent,
        // component: PortalComponent,
        resolve: [AuthResolver, UserProfileResolver],
        children: [
            {
                path: '',
                canActivate: [AuthCanActivateGuard],
                component: SelectComponent
            }
        ]
    },
    {
        path: ':slug',
        component: PortalComponent,
        resolve: {
            participant: ParticipantResolve,
            node: NodeResolve
        },
        canActivate: [AuthCanActivateGuard],
        data: {
            redirectBase: '/portal',
            redirectUnauthorized: '/unauthorized'
        },
        children: [
            { path: '', redirectTo: 'transactions', pathMatch: 'full' },
            // {
            //     path: 'overview',
            //     component: OverviewComponent,
            //     canActivate: [ParticipantPermissionsGuard],
            //     data: {
            //         title: 'Overview',
            //         participant_permissions: ['viewer', 'manager', 'admin'],
            //         super_permissions: ['viewer', 'manager', 'admin']
            //     },
            // },
            {
                path: 'status',
                canActivate: [ParticipantPermissionsGuard],
                data: {
                    title: 'System Status',
                    participant_permissions: ['viewer', 'manager', 'admin'],
                    super_permissions: ['viewer', 'manager', 'admin']
                },
                children: [
                    {
                        path: '',
                        component: StatusOverviewComponent,
                    },
                    {
                        path: 'history/:serviceName',
                        component: StatusHistoryComponent,
                        data: {
                            title: 'Status Error History',
                        },
                    },
                ]
            },
            // {
            //     path: 'diligence',
            //     component: DiligenceComponent,
            //     canActivate: [ParticipantPermissionsGuard],
            //     data: {
            //         title: 'Diligence',
            //         participant_permissions: ['viewer', 'manager', 'admin'],
            //         super_permissions: ['manager', 'admin']
            //     },
            // },
            {
                path: 'users',
                component: UsersComponent,
                canActivate: [ParticipantPermissionsGuard],
                data: {
                    title: 'Users',
                    participant_permissions: ['admin'],
                    super_permissions: ['manager', 'admin']
                },
            },
            {
                path: 'account',
                component: AccountComponent,
                canActivate: [ParticipantPermissionsGuard],
                data: {
                    title: 'Account',
                    participant_permissions: ['viewer', 'manager', 'admin'],
                    super_permissions: ['manager', 'admin']
                },
            },
            {
                path: 'token',
                component: TokenComponent,
                canActivate: [ParticipantPermissionsGuard],
                data: {
                    title: 'Access Tokens',
                    participant_permissions: ['manager', 'admin'],
                    super_permissions: ['manager', 'admin']
                },
            },
            // {
            //     path: 'vpn',
            //     component: VpnComponent,
            //     canActivate: [ParticipantPermissionsGuard],
            //     data: {
            //         title: 'VPN Setup',
            //         participant_permissions: ['manager', 'admin'],
            //         super_permissions: ['manager', 'admin']
            //     },
            // },
            // {
            //     path: 'billing',
            //     component: BillingComponent,
            //     canActivate: [ParticipantPermissionsGuard],
            //     data: {
            //         title: 'Billing',
            //         participant_permissions: ['viewer', 'manager', 'admin'],
            //         super_permissions: ['viewer', 'manager', 'admin']
            //     },
            // },
            {
                path: 'transactions',
                component: TransactionsComponent,
                canActivate: [ParticipantPermissionsGuard],
                data: {
                    title: 'Transactions',
                    participant_permissions: ['viewer', 'manager', 'admin'],
                    super_permissions: ['viewer', 'manager', 'admin']
                },
            },
            {
                // handles routes for a participant's accounts and assets on World Wire
                path: 'accounts',
                component: AccountsComponent,
                data: {
                    title: 'Accounts Management',
                },
                loadChildren: './accounts/accounts.module#AccountsModule'
            },
            {
                // handles routes for a participant's accounts and assets on World Wire
                path: 'assets',
                component: AssetsComponent,
                data: {
                    title: 'Accounts Management',
                },
                loadChildren: './assets/assets.module#AssetsModule'
            },
            // {
            //     path: 'settings',
            //     component: SettingsComponent,
            //     canActivate: [ParticipantPermissionsGuard],
            //     data: {
            //         title: 'Settings',
            //         participant_permissions: ['manager', 'admin'],
            //         super_permissions: ['manager', 'admin']
            //     },
            //     loadChildren: './settings/settings.module#SettingsModule'
            // },
        ]
    }
];

@NgModule({
    imports: [
        RouterModule.forChild(portalRoutes)
    ],
    exports: [
        RouterModule
    ],
    providers: [
        ParticipantResolve,
        NodeResolve,
        ParticipantPermissionsGuard
    ]
})
export class PortalRoutingModule { }
