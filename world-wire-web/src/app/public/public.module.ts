// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { PublicRoutingModule } from './public-routing.module';
import { FormsModule } from '@angular/forms';
import { FlexLayoutModule } from '../../../node_modules/@angular/flex-layout';
import { HeaderFooterModule } from '../shared/header-footer.module';
import { PhonePipe } from '../shared/pipes/phone.pipe';
import { PublicComponent } from './public.component';
import { LandingComponent } from './landing/landing.component';
import { CareersComponent } from './careers/careers.component';
import { TermsComponent } from './terms/terms.component';
import { PrivacyComponent } from './privacy/privacy.component';
import { UnsubscribeComponent } from './unsubscribe/unsubscribe.component';
import { ContactComponent } from './contact/contact.component';
import { OnlinePrivacyComponent } from './online-privacy/online-privacy.component';
import { LoginTokenComponent } from './login-token/login-token.component';
import { UnauthorizedComponent } from './unauthorized/unauthorized.component';
import { InactivityComponent } from './inactivity/inactivity.component';
import { UnrecognizedComponent } from './unrecognized/unrecognized.component';
import { FidComponent } from './fid/fid.component';
import { CodeSnippetModule } from 'carbon-components-angular';

@NgModule({
    declarations: [
        PublicComponent,
        LandingComponent,
        CareersComponent,
        PhonePipe,
        TermsComponent,
        PrivacyComponent,
        OnlinePrivacyComponent,
        UnsubscribeComponent,
        ContactComponent,
        LoginTokenComponent,
        UnauthorizedComponent,
        InactivityComponent,
        UnrecognizedComponent,
        FidComponent
    ],
    imports: [
        CommonModule,
        FormsModule,
        HeaderFooterModule,
        PublicRoutingModule,
        FlexLayoutModule,
        CodeSnippetModule
    ]
})
export class PublicModule { }
