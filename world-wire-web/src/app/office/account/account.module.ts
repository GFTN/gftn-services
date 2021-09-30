// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FlexLayoutModule } from '@angular/flex-layout';
import { InputModule, ButtonModule, ModalModule, PlaceholderModule, InlineLoadingModule } from 'carbon-components-angular';
import { AccountRoutingModule } from './account-routing.module';
import { AccountComponent } from './account.component';
import { NodesComponent } from './nodes/nodes.component';
import { NodeConfigComponent } from './nodes/nodeconfig/nodeconfig.component';
import { NodeConfigFormModule } from '../../shared/node-config-form.module';

import { EditNodeComponent } from './nodes/edit-node/edit-node.component';
import { DeleteNodeModalComponent } from './nodes/delete-node-modal/delete-node-modal.component';
import { FormsModule } from '@angular/forms';
import { FormValidatorsModule } from '../../shared/form-validators.module';
import { AccountOverviewComponent } from './account-overview/account-overview.component';
import { PipeModule } from '../../shared/pipe-module.module';
import { RegistrationModalComponent } from './nodes/registration-modal/registration-modal.component';

import { NodeService } from './shared/node.service';

@NgModule({
  declarations: [
    AccountComponent,
    NodesComponent,
    NodeConfigComponent,
    EditNodeComponent,
    DeleteNodeModalComponent,
    AccountOverviewComponent,
    RegistrationModalComponent
  ],
  imports: [
    CommonModule,
    PipeModule,
    FormsModule,
    FlexLayoutModule,
    // Carbon Angular Modules
    InputModule,
    ButtonModule,
    ModalModule,
    PlaceholderModule,
    InlineLoadingModule,
    // Custom Modules for AccountModule
    FormValidatorsModule,
    NodeConfigFormModule,
    AccountRoutingModule
  ],
  providers: [
    NodeService
  ],
  // necessary for dialogs
  entryComponents: [
    EditNodeComponent,
    DeleteNodeModalComponent,
    RegistrationModalComponent
  ]
})
export class AccountModule { }
