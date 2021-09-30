// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ComingComponent } from './coming.component';

const comingRoutes: Routes = [
    {
        path: '', component: ComingComponent,
    }
];

@NgModule({
    imports: [
        RouterModule.forChild(comingRoutes)
    ],
    exports: [
        RouterModule
    ]
})
export class ComingRoutingModule { }
