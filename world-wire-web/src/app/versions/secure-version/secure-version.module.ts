// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DocsRoutingModule } from './secure-version-routing.module';
import { FormsModule } from '@angular/forms';
import { VersionSharedModule } from '../shared/version-shared.module';

import { FlexLayoutModule } from '@angular/flex-layout';
import { SecureVersionComponent } from './secure-version.component';

import { SecureVersionGuard } from '../../shared/guards/secure-version.guard';

import { HeaderFooterModule } from '../../shared/header-footer.module';

@NgModule({
    declarations: [
        SecureVersionComponent,
    ],
    entryComponents: [
    ],
    imports: [
        CommonModule,
        DocsRoutingModule,
        FlexLayoutModule,
        VersionSharedModule,
        HeaderFooterModule,
        FormsModule
    ],
    providers: [
        SecureVersionGuard
    ]
})
export class VersionModule { }
