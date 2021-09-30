// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FlexLayoutModule } from '@angular/flex-layout';
import { TwoFactorComponent } from './two-factor.component';
import { TwoFactorRoutingModule } from './two-factor-routing.module';
import { RegisterComponent } from './register/register.component';
import { MatDialogModule } from '@angular/material';
import { FormsModule } from '@angular/forms';
import { ProgressIndicatorModule, ButtonModule, InputModule } from 'carbon-components-angular';
import { VerifyComponent } from './verify/verify.component';

@NgModule({
    declarations: [
        TwoFactorComponent,
        RegisterComponent,
        VerifyComponent
    ],
    imports: [
        CommonModule,
        FlexLayoutModule,
        FormsModule,
        MatDialogModule,
        ProgressIndicatorModule,
        ButtonModule,
        InputModule,
        TwoFactorRoutingModule
    ]
})
export class TwoFactorModule { }
