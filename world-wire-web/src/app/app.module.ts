// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// Modules
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgModule } from '@angular/core';
import { AppRoutingModule } from './app-routing.module';
import { CustomMaterialModule } from './shared/custom-material.module';

// Components
import { AppComponent } from './app.component';
import { PageNotFoundComponent } from './not-found/not-found.component';

import { environment } from '../environments/environment';

import { AngularFireModule } from '@angular/fire';

import { AngularFirestoreModule } from '@angular/fire/firestore';
import { AngularFireDatabaseModule } from '@angular/fire/database';
import { DocumentService } from './shared/services/document.service';
import { WindowService } from './shared/services/window.service';
import { ExternalRedirectGuard } from './shared/guards/external-redirect.guard';
import { SecureVersionGuard } from './shared/guards/secure-version.guard';
import { AuthModule } from './shared/auth.module';
import { HttpClientModule } from '@angular/common/http';
import { Confirm2faModule } from './shared/confirm2fa.module';
import { UtilsService } from './shared/utils/utils';

@NgModule({
  declarations: [
    AppComponent,
    PageNotFoundComponent
  ],
  imports: [
    BrowserModule,
    AngularFireModule.initializeApp(environment.firebase),
    AngularFirestoreModule,
    AngularFireDatabaseModule,
    BrowserAnimationsModule,
    CustomMaterialModule,
    HttpClientModule,
    AppRoutingModule, // <-- Important: AppRoutingModule must come last https://angular.io/guide/router#module-import-order-matters
    Confirm2faModule,
    AuthModule
  ],
  providers: [
    DocumentService,
    WindowService,
    ExternalRedirectGuard,
    SecureVersionGuard,
    UtilsService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
