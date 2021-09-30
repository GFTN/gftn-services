// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { TwoFactorComponent } from './two-factor.component';
import { RegisterComponent } from './register/register.component';
import { VerifyComponent } from './verify/verify.component';

const routes: Routes = [
    {
        path: '', component: TwoFactorComponent,
        children: [
            { path: 'register', component: RegisterComponent },
            { path: 'verify', component: VerifyComponent }
        ]
    }
];

@NgModule({
    imports: [
        RouterModule.forChild(routes)
    ],
    exports: [
        RouterModule
    ]
})
export class TwoFactorRoutingModule { }
