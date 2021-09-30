// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { FaqComponent } from './faq/faq.component';
import { OnboardingComponent } from './onboarding/onboarding.component';
import { IntroductionComponent } from './introduction/introduction.component';
import { ConceptsComponent } from './concepts/concepts.component';
import { GuidesComponent } from './guides/guides.component';
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { VersionComponent } from './version.component';
import { ApiComponent } from './api/api.component';
import { SamplesComponent } from './samples/samples.component';
import { GettingStartedComponent } from './getting-started/getting-started.component';
import { ChangeLogComponent } from './change-log/change-log.component';
import { VersionGuard } from '../shared/services/version.guard';

const docsRoutes: Routes = [
    {
        path: '',
        component: VersionComponent,
        children: [
            {
                path: '',
                redirectTo: 'introduction',
                pathMatch: 'full',
                canActivate: [VersionGuard]
            },
            {
                path: 'introduction',
                component: IntroductionComponent,
                canActivate: [VersionGuard]
            },
            {
                path: 'onboarding',
                component: OnboardingComponent,
                canActivate: [VersionGuard]
            },
            {
                path: 'concepts',
                component: ConceptsComponent,
                canActivate: [VersionGuard]
            },
            {
                path: 'guides',
                component: GuidesComponent,
                canActivate: [VersionGuard]
            },
            {
                path: 'getting-started',
                component: GettingStartedComponent,
                canActivate: [VersionGuard]
            },
            {
                path: 'faq',
                component: FaqComponent,
                canActivate: [VersionGuard]
            },
            {
                path: 'changes',
                component: ChangeLogComponent,
                canActivate: [VersionGuard]
            },
            {
                path: 'api',
                redirectTo: 'api/client-api',
                pathMatch: 'full',
                canActivate: [VersionGuard]
            },
            {
                path: 'api/:name',
                component: ApiComponent,
                canActivate: [VersionGuard]
            },
            {
                path: 'api/:name/:jump',
                component: ApiComponent,
                canActivate: [VersionGuard]
            },
            { path: 'samples', component: SamplesComponent }
        ]
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
