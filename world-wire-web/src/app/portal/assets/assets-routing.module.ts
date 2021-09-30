// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { ParticipantPermissionsGuard } from '../../shared/guards/permissions-redirect.guard';
import { NodeResolve } from '../../shared/guards/node.resolve';
import { AssetsOverviewComponent } from './assets-overview/assets-overview.component';
import { AccountsResolve } from '../accounts/shared/guards/accounts.resolve';


const assetsRoutes: Routes = [
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
        component: AssetsOverviewComponent,
        canActivate: [ParticipantPermissionsGuard],
        data: {
          title: 'Assets Management',
          shortTitle: 'Assets Overview',
          participant_permissions: ['manager', 'admin'],
        },
      },
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
    RouterModule.forChild(assetsRoutes)
  ],
  exports: [
    RouterModule
  ],
  providers: [
    AccountsResolve
  ]
})
export class AssetsRoutingModule { }
