// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { ComingComponent } from './coming.component';
import { CommonModule } from '@angular/common';
import { ComingRoutingModule } from './coming-routing.module';
import { FlexLayoutModule } from '@angular/flex-layout';

@NgModule({
    declarations: [
        ComingComponent
    ],
    imports: [
        CommonModule,
        FlexLayoutModule,
        ComingRoutingModule
    ]
})
export class ComingModule { }
