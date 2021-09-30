// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { TrustedAssetsComponent } from './trusted-assets/trusted-assets.component';
import { TrustlinesComponent } from './trustlines/trustlines.component';
import { ParticipantPermissionsGuard } from '../../shared/guards/permissions-redirect.guard';
import { AccountsOverviewComponent } from './accounts-overview/accounts-overview.component';
import { WhitelistComponent } from './whitelist/whitelist.component';
import { NodeResolve } from '../../shared/guards/node.resolve';
import { AccountsResolve } from './shared/guards/accounts.resolve';


const accountsRoutes: Routes = [
  {
    path: '',
    resolve: [NodeResolve, AccountsResolve],
    children: [
      {
        path: '',
        redirectTo: 'overview',
        pathMatch: 'full'
      },
      {
        path: 'overview',
        component: AccountsOverviewComponent,
        canActivate: [ParticipantPermissionsGuard],
        data: {
          title: 'Accounts Management',
          shortTitle: 'Accounts Overview',
          participant_permissions: ['manager', 'admin'],
          super_permissions: ['manager', 'admin']
        },
      },
      {
        path: 'trustlines',
        component: TrustlinesComponent,
        data: {
          title: 'Trustline Requests',
          shortTitle: 'Trustline Requests'
        },
      },
      {
        path: 'whitelist',
        component: WhitelistComponent,
        data: {
          title: 'Whitelist Management',
          shortTitle: 'Manage Whitelist'
        },
      },
      {
        // '/accounts/:account_name'
        path: ':slug',
        canActivate: [ParticipantPermissionsGuard],
        data: {
          participant_permissions: ['viewer', 'manager', 'admin'],
          super_permissions: ['viewer', 'manager', 'admin']
        },
        children: [
          {
            path: '',
            // temporary redirect until we have overview page
            // for a single operating account
            redirectTo: 'assets',
            pathMatch: 'full'
          },
          {
            path: 'assets',
            component: TrustedAssetsComponent,
            // TODO: Implement new guard to lock down permissions for an operating account
            // canActivate: [AccountPermissionsGuard],
            data: {
              title: 'Trusted Assets Management',
              shortTitle: 'Trusted Assets',
            },
          },
        ]
      }
    ]
  }
];

/**
 *
 *
 * @export
 * @class AccountsRoutingModule
 */
@NgModule({
  imports: [
    RouterModule.forChild(accountsRoutes)
  ],
  exports: [
    RouterModule
  ],
  providers: [
    AccountsResolve
  ]
})
export class AccountsRoutingModule { }
