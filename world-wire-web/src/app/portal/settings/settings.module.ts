// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AccountAssetsComponent } from './account-assets/account-assets.component';
import { NodeConfigComponent } from './nodeconfig/nodeconfig.component';
import { SettingsRoutingModule } from './settings-routing.module';
import { NodeConfigFormModule } from '../../shared/node-config-form.module';
// Carbon Components
import { NotificationModule, ContentSwitcherModule } from 'carbon-components-angular';

@NgModule({
  declarations: [
    AccountAssetsComponent,
    NodeConfigComponent,
  ],
  imports: [
    CommonModule,
    NodeConfigFormModule,
    // Carbon Angular Modules
    NotificationModule,
    ContentSwitcherModule,
    SettingsRoutingModule
  ]
})
export class SettingsModule { }
