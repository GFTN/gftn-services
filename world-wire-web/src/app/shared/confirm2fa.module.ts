// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CustomMaterialModule } from './custom-material.module';
import { FlexLayoutModule } from '@angular/flex-layout';
import { RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { Confirm2faService } from './services/confirm2fa.service';
import { Confirm2faDialogComponent } from './components/confirm2fa-dialog/confirm2fa-dialog.component';
import { FormsModule } from '@angular/forms';
import { InputModule } from 'carbon-components-angular';

@NgModule({
  declarations: [
    Confirm2faDialogComponent
  ],
  imports: [
    CustomMaterialModule,
    CommonModule,
    RouterModule,
    FlexLayoutModule,
    FormsModule,
    InputModule
  ],
  providers: [
    Confirm2faService
  ],
  entryComponents: [
    Confirm2faDialogComponent
  ]
})
export class Confirm2faModule { }
