// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { NodeConfigFormComponent } from './components/node-config-form/node-config-form.component';
import { FlexLayoutModule } from '@angular/flex-layout';
import { NotificationModule, ProgressIndicatorModule, DropdownModule, InlineLoadingModule, InputModule, ButtonModule } from 'carbon-components-angular';
import { FormValidatorsModule } from './form-validators.module';
export { NodeConfigFormComponent } from './components/node-config-form/node-config-form.component';

@NgModule({
  declarations: [
    NodeConfigFormComponent
  ],
  exports: [
    NodeConfigFormComponent
  ],
  imports: [
    CommonModule,
    FormsModule,
    // Carbon Angular Modules
    NotificationModule,
    ProgressIndicatorModule,
    DropdownModule,
    InlineLoadingModule,
    InputModule,
    ButtonModule,
    FlexLayoutModule,
    FormValidatorsModule
  ]
})
export class NodeConfigFormModule { }
