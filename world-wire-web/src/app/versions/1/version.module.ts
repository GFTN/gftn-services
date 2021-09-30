// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DocsRoutingModule } from './version-routing.module';
import { FormsModule } from '@angular/forms';
import { VersionSharedModule } from '../shared/version-shared.module';

import { ApiComponent } from './api/api.component';
import { SamplesComponent } from './samples/samples.component';

import { FlexLayoutModule } from '@angular/flex-layout';
import { AccordionModule } from 'carbon-components-angular';
import { ChangeLogComponent } from './change-log/change-log.component';

import { VersionComponent } from './version.component';
import { VersionNavComponent } from './shared/version-nav/version-nav.component';

import { PathComponent } from './api/path/path.component';
import { FaqComponent } from './faq/faq.component';
import { ResponseComponent } from './api/response/response.component';
import { OperationComponent } from './api/operation/operation.component';
import { ModelComponent } from './api/model/model.component';
import { ParametersComponent } from './api/parameters/parameters.component';


import { SwaggerService } from '../shared/services/swagger.service';
import { WwHttpService } from './shared/services/ww-http.service';

import { DemoSandboxComponent } from './getting-started/demo-sandbox/demo-sandbox.component';
import { OnboardingComponent } from './onboarding/onboarding.component';
import { IntroductionComponent } from './introduction/introduction.component';
import { ConceptsComponent } from './concepts/concepts.component';
import { GuidesComponent } from './guides/guides.component';
import { GettingStartedComponent } from './getting-started/getting-started.component';
import { HeaderFooterModule } from '../../shared/header-footer.module';
import { VersionGuard } from '../shared/services/version.guard';
import { GroupEndpointsComponent } from './shared/group-endpoints/group-endpoints.component';
import { MarkdownModule } from 'ngx-markdown';

@NgModule({
    declarations: [
        VersionComponent,
        ApiComponent,
        SamplesComponent,
        VersionNavComponent,
        ChangeLogComponent,
        PathComponent,
        FaqComponent,
        ResponseComponent,
        OperationComponent,
        ModelComponent,
        ParametersComponent,
        DemoSandboxComponent,
        IntroductionComponent,
        ConceptsComponent,
        GuidesComponent,
        OnboardingComponent,
        GettingStartedComponent,
        GroupEndpointsComponent,
    ],
    entryComponents: [
    ],
    imports: [
        CommonModule,
        DocsRoutingModule,
        FlexLayoutModule,
        AccordionModule,
        VersionSharedModule,
        HeaderFooterModule,
        FormsModule,
        MarkdownModule.forRoot()
    ],
    providers: [
        SwaggerService,
        WwHttpService,
        VersionGuard
    ]
})
export class VersionModule { }
