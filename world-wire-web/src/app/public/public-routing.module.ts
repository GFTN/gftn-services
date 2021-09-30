// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
//
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { CareersComponent } from './careers/careers.component';
import { LandingComponent } from './landing/landing.component';
import { PublicComponent } from './public.component';
import { TermsComponent } from './terms/terms.component';
import { PrivacyComponent } from './privacy/privacy.component';
import { UnauthorizedComponent } from './unauthorized/unauthorized.component';
import { UnsubscribeComponent } from './unsubscribe/unsubscribe.component';
import { ContactComponent } from './contact/contact.component';
import { OnlinePrivacyComponent } from './online-privacy/online-privacy.component';
import { ExternalRedirectGuard } from '../shared/guards/external-redirect.guard';
import { LoginTokenComponent } from './login-token/login-token.component';
import { AuthResolver } from '../shared/guards/auth.resolver';
import { InactivityComponent } from './inactivity/inactivity.component';
import { UnrecognizedComponent } from './unrecognized/unrecognized.component';
import { FidComponent } from './fid/fid.component';

const publicRoutes: Routes = [
    {
        path: '', component: PublicComponent,
        children: [
            {
                path: '',
                // redirectTo: '/soon',
                pathMatch: 'full',
                canActivate: [ExternalRedirectGuard],
                data: {
                    externalUrl: 'https://www.ibm.com/blockchain/solutions/gftn'
                }
            }, // redirects to "Coming Soon" module
            { path: 'home', component: LandingComponent },
            { path: 'careers', component: CareersComponent },
            { path: 'login/:token', component: LoginTokenComponent },
            {
                path: 'login',
                component: LoginTokenComponent,
                resolve: [AuthResolver]
            },
            // fid component used to get x-fid header for use in back-end
            {
                path: 'fid',
                component: FidComponent,
                resolve: [AuthResolver]
            },
            // 'unrecognized:' for users who have no permissions provisioned (tells user to contact support)
            { path: 'unrecognized', component: UnrecognizedComponent, resolve: [AuthResolver] },
            // 'unauthorized:' for users trying to access a restricted part of the site without access rights
            { path: 'unauthorized', component: UnauthorizedComponent, resolve: [AuthResolver] },
            { path: 'inactive', component: InactivityComponent, resolve: [AuthResolver] },
            { path: 'unsubscribe/:emailHash/:mailingList/:all', component: UnsubscribeComponent },
            { path: 'unsubscribe/:emailHash/:mailingList', component: UnsubscribeComponent },
            { path: 'contact', component: ContactComponent },
            // TODO: use redirect for terms, privacy and online privacy
            // for now until we have customer T&C for WW
            {
                path: 'terms',
                component: TermsComponent,
                canActivate: [ExternalRedirectGuard],
                // resolve: {
                //     url: ExternalRedirectResolver
                // },
                data: {
                    externalUrl: 'https://www.ibm.com/legal/us/en/'
                }
            },
            {
                path: 'privacy',
                component: PrivacyComponent,
                canActivate: [ExternalRedirectGuard],
                // resolve: {
                //     url: ExternalRedirectResolver
                // },
                data: {
                    externalUrl: 'https://www.ibm.com/privacy/us/en/'
                }
            },
            {
                path: 'online-privacy',
                component: OnlinePrivacyComponent,
                canActivate: [ExternalRedirectGuard],
                // resolve: {
                //     url: ExternalRedirectResolver
                // },
                data: {
                    externalUrl: 'https://www.ibm.com/privacy/details/us/en/'
                }
            }
        ]
    }
];

@NgModule({
    imports: [
        RouterModule.forChild(publicRoutes)
    ],
    exports: [
        RouterModule
    ]
})
export class PublicRoutingModule { }
