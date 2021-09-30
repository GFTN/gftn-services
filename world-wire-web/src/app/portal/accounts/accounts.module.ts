// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { InputModule, ButtonModule, DropdownModule, ModalModule, PlaceholderModule, NotificationModule, TabsModule, DialogModule, TableModule, LoadingModule, StructuredListModule, InlineLoadingModule } from 'carbon-components-angular';
import { AccountsRoutingModule } from './accounts-routing.module';
import { TrustedAssetsComponent } from './trusted-assets/trusted-assets.component';
import { TrustlinesComponent } from './trustlines/trustlines.component';
import { AccountsOverviewComponent } from './accounts-overview/accounts-overview.component';
import { FlexLayoutModule } from '@angular/flex-layout';
import { TrustlineModalComponent } from './trustline-modal/trustline-modal.component';
import { PortalSharedModule } from '../shared/portal-shared.module';
import { WhitelistComponent } from './whitelist/whitelist.component';
import { CustomCarbonModule } from '../../shared/custom-carbon.module';
import { WhitelistModalComponent } from './whitelist-modal/whitelist-modal.component';
import { ParticipantAccountComponent } from './participant-account/participant-account.component';
import { AccountModalComponent } from '../shared/components/account-modal/account-modal.component';

/**
 * Holds all Modules related to a participant's
 * accounts on the World Wire network
 *
 * @export
 * @class AccountsModule
 */
@NgModule({
  declarations: [
    TrustedAssetsComponent,
    TrustlinesComponent,
    AccountsOverviewComponent,
    TrustlineModalComponent,
    WhitelistComponent,
    WhitelistModalComponent,
    ParticipantAccountComponent
  ],
  imports: [
    CommonModule,
    FormsModule,
    FlexLayoutModule,
    PlaceholderModule,
    InputModule,
    DropdownModule,
    ButtonModule,
    ModalModule,
    NotificationModule,
    TabsModule,
    DialogModule,
    TableModule,
    StructuredListModule,
    InlineLoadingModule,
    LoadingModule,
    PlaceholderModule,
    CustomCarbonModule,
    PortalSharedModule,
    AccountsRoutingModule
  ],
  entryComponents: [
    TrustlineModalComponent,
    WhitelistModalComponent,
    AccountModalComponent
  ],
  providers: []
})
export class AccountsModule { }
