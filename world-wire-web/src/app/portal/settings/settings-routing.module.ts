// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AccountAssetsComponent } from './account-assets/account-assets.component';
import { NodeConfigComponent } from './nodeconfig/nodeconfig.component';

const settingsRoutes: Routes = [
  {
    path: 'account',
    component: AccountAssetsComponent,
    data: {
      title: 'Accounts and Assets Management'
    },
  },
  {
    path: 'nodes',
    component: NodeConfigComponent,
    data: {
      title: 'Node Configuration'
    },
  }
];

@NgModule({
  imports: [
      RouterModule.forChild(settingsRoutes)
  ],
  exports: [
      RouterModule
  ]
})
export class SettingsRoutingModule { }
