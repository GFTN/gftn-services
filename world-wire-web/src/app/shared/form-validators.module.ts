// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ParticipantIdValidator } from './directives/participant-id-valid.directive';
import { ConfirmInputValidator } from './directives/confirm-input.directive';
import { ParticipantService } from './services/participant.service';
import { EmailValidator } from './directives/valid-ibm-email.directive';
import { InstitutionIdValidator } from './directives/institution-id-valid.directive';
import { AssetValidator } from './directives/asset-valid.directive';

@NgModule({
  declarations: [
    ParticipantIdValidator,
    ConfirmInputValidator,
    EmailValidator,
    InstitutionIdValidator,
    AssetValidator
  ],
  imports: [
    CommonModule
  ],
  exports: [
    ParticipantIdValidator,
    ConfirmInputValidator,
    EmailValidator,
    InstitutionIdValidator,
    AssetValidator
  ],
  providers: [
    ParticipantService
  ],
})
export class FormValidatorsModule { }
