// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FlexLayoutModule } from '@angular/flex-layout';
import { RouterModule } from '@angular/router';
import { HeaderFooterModule } from '../../shared/header-footer.module';
import { RegexPipe } from '../../shared/pipes/regex.pipe';
import { DocumentService } from '../../shared/services/document.service';
import { DocsFeedbackComponent } from './docs-feedback/docs-feedback.component';
import { TooltipComponent } from './tooltip/tooltip.component';
import { VersionSelectComponent } from './version-select/version-select.component';
import { VersionService } from './services/version.service';
import { DropdownModule } from 'carbon-components-angular';
import { SwaggerService } from './services/swagger.service';
import { VersionGuard } from './services/version.guard';

@NgModule({
    imports: [
        CommonModule,
        RouterModule,
        FlexLayoutModule,
        HeaderFooterModule,
        DropdownModule
    ],
    providers: [
        RegexPipe,
        VersionService,
        SwaggerService,
        VersionGuard,
        DocumentService
    ],
    declarations: [
        TooltipComponent,
        RegexPipe,
        DocsFeedbackComponent,
        VersionSelectComponent
    ],
    exports: [
        TooltipComponent,
        RegexPipe,
        DocsFeedbackComponent,
        VersionSelectComponent
    ]
})
export class VersionSharedModule { }
