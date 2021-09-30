// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { FlexLayoutModule } from '@angular/flex-layout';
import { AngularFireAuth } from '@angular/fire/auth';
import { AuthService } from './services/auth.service';
import { CustomMaterialModule } from './custom-material.module';
import { ManageParticipantDialogComponent } from '../office/accounts/manage-participant-dialog/manage-participant-dialog.component';
import { AuthCanActivateGuard, AuthRedirectCookie } from './guards/auth.guard';
import { AuthResolver } from './guards/auth.resolver';
import {
    ParticipantPermissionsDialogComponent,
} from './components/participant-permissions-dialog/participant-permissions-dialog.component';
import { SessionService } from './services/session.service';
import { SuperPermissionsDialogComponent } from './components/super-permissions-dialog/super-permissions-dialog.component';
import { ParticipantPermissionsService } from './services/participant-permissions.service';
import { SuperPermissionsService } from './services/super-permissions.service';
import { UserProfileResolver } from './guards/user-profile.resolver';
import { EmailValidator } from './directives/valid-ibm-email.directive';
import { FormValidatorsModule } from './form-validators.module';
import { ModalModule, PlaceholderModule } from 'carbon-components-angular';

/**
 * This shared module is used to import components and services related to managing users
 *
 * @export
 * @class UserModule
 */
@NgModule({
    imports: [
        ModalModule,
        PlaceholderModule,
        CustomMaterialModule,
        CommonModule,
        FormsModule,
        FormValidatorsModule, /* Put modules before Router module */
        RouterModule,
        FlexLayoutModule,

    ],
    declarations: [
        ManageParticipantDialogComponent,
        ParticipantPermissionsDialogComponent,
        SuperPermissionsDialogComponent
    ],
    entryComponents: [
        ManageParticipantDialogComponent,
        ParticipantPermissionsDialogComponent,
        SuperPermissionsDialogComponent
    ],
    providers: [
        AngularFireAuth,
        AuthService,
        AuthCanActivateGuard,
        AuthRedirectCookie,
        AuthResolver,
        UserProfileResolver,
        SessionService,
        ParticipantPermissionsService,
        SuperPermissionsService
    ]
})
export class AuthModule { }
