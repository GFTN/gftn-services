// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AssetsRoutingModule } from './assets-routing.module';
import { FormsModule } from '@angular/forms';
import { FlexLayoutModule } from '@angular/flex-layout';
import { PlaceholderModule, LoadingModule, ButtonModule, DialogModule, ModalModule, DropdownModule, InputModule, NotificationModule } from 'carbon-components-angular';
import { CustomCarbonModule } from '../../shared/custom-carbon.module';
import { PortalSharedModule } from '../shared/portal-shared.module';

import { AssetsOverviewComponent } from './assets-overview/assets-overview.component';
import { AssetCardComponent } from './asset-card/asset-card.component';
import { AssetModalComponent } from '../shared/components/asset-modal/asset-modal.component';
import { FormValidatorsModule } from '../../shared/form-validators.module';
import { AssetDetailsModalComponent } from './asset-details-modal/asset-details-modal.component';

@NgModule({
  declarations: [
    AssetsOverviewComponent,
    AssetCardComponent,
    AssetModalComponent,
    AssetDetailsModalComponent,
  ],
  imports: [
    CommonModule,
    FormsModule,
    FormValidatorsModule,
    FlexLayoutModule,
    // Carbon Modules
    PlaceholderModule,
    LoadingModule,
    ButtonModule,
    ModalModule,
    DropdownModule,
    InputModule,
    DialogModule,
    NotificationModule,
    // Custom Modules
    PortalSharedModule,
    CustomCarbonModule,
    AssetsRoutingModule
  ],
  entryComponents: [
    AssetModalComponent,
    AssetDetailsModalComponent
  ]
})
export class AssetsModule { }
