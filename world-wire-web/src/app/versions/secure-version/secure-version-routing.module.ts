// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { SecureVersionComponent } from './secure-version.component';

const docsRoutes: Routes = [
    {
        path: '',
        component: SecureVersionComponent
    }
];

@NgModule({
    imports: [
        RouterModule.forChild(docsRoutes)
    ],
    exports: [
        RouterModule
    ]
})
export class DocsRoutingModule { }
