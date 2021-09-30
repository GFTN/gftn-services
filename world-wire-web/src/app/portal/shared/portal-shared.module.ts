// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { FlexLayoutModule } from '@angular/flex-layout';
import { DropdownModule, ModalModule, LoadingModule, ButtonModule, InputModule, NotificationModule } from 'carbon-components-angular';

import { NodeSelectComponent } from './components/node-select/node-select.component';
import { QuickFilterComponent } from './components/quick-filter/quick-filter.component';
import { AccountModalComponent } from './components/account-modal/account-modal.component';
import { AnchorAssetsComponent } from './components/anchor-assets/anchor-assets.component';

@NgModule({
  declarations: [
    NodeSelectComponent,
    QuickFilterComponent,
    AccountModalComponent,
    AnchorAssetsComponent,
  ],
  exports: [
    NodeSelectComponent,
    QuickFilterComponent,
    AnchorAssetsComponent
  ],
  imports: [
    FormsModule,
    FlexLayoutModule,
    DropdownModule,
    ModalModule,
    ButtonModule,
    InputModule,
    LoadingModule,
    NotificationModule,
    CommonModule
  ],
})
export class PortalSharedModule { }
