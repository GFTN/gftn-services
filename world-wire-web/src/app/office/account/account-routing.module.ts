// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { NodesComponent } from './nodes/nodes.component';
import { NodeConfigComponent } from './nodes/nodeconfig/nodeconfig.component';
import { AccountComponent } from './account.component';
import { ParticipantResolve } from '../../shared/guards/participant.resolve';
import { SuperPermissionsGuard } from '../../shared/guards/super-permissions-redirect.guard';
import { AccountOverviewComponent } from './account-overview/account-overview.component';

const accountRoutes: Routes = [

  // no slug. redirect back to main account page
  { path: '', redirectTo: '/office/accounts', pathMatch: 'full' },
  {
    path: ':slug',
    component: AccountComponent,
    canActivate: [SuperPermissionsGuard],
    data: {
      super_permissions: ['admin']
    },
    resolve: {
      participant: ParticipantResolve
    },
    children: [
      {
        path: '',
        component: AccountOverviewComponent,
        canActivate: [SuperPermissionsGuard],
        data: {
          title: 'Manage Participant Account',
          shortTitle: 'Manage Account',
          super_permissions: ['viewer', 'manager', 'admin']
        },
      },
      {
        path: 'nodes',
        component: NodesComponent,
        canActivate: [SuperPermissionsGuard],
        data: {
          title: 'Manage Participant Nodes',
          shortTitle: 'Nodes',
          super_permissions: ['admin']
        },

      },
      {
        path: 'nodes/add',
        component: NodeConfigComponent,
        canActivate: [SuperPermissionsGuard],
        data: {
          title: 'Add Participant Node',
          shortTitle: 'Add Node',
          super_permissions: ['admin']
        },
      },
    ]
  }

];

@NgModule({
  imports: [
    RouterModule.forChild(accountRoutes)
  ],
  exports: [
    RouterModule
  ],
  providers: [
    ParticipantResolve,
    SuperPermissionsGuard
  ]
})
export class AccountRoutingModule { }
