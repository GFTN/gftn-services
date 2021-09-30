// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';

// IMPORTANT: Do not import BrowserAnimationModule here as it will try to re-load
// thi module every time an new route is lazy loaded and throw an error while routing to new module.
// import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

// import { FlexLayoutModule } from '@angular/flex-layout';
// https://github.com/angular/flex-layout/wiki
// https://tburleson-layouts-demos.firebaseapp.com/#/docs

// IMPORTANT: import modules 'INDIVIDUALLY' for improved performance
// import { MatButtonModule } from '@angular/material/button';
// import { MatCheckboxModule } from '@angular/material/checkbox';
// import { MatInputModule } from '@angular/material/input';
// import { MatFormFieldModule } from '@angular/material/form-field';
// import { MatIconModule } from '@angular/material/icon';
// import { MatMenuModule } from '@angular/material/menu';
// import { MatToolbarModule } from '@angular/material/toolbar';
// import { MatTabsModule } from '@angular/material/tabs';
// import { MatRippleModule } from '@angular/material/core';
// import { MatCardModule } from '@angular/material/card';
import { MatDialogModule } from '@angular/material/dialog';
// import { ReactiveFormsModule } from '@angular/forms';
// import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
// import { MatSidenavModule } from '@angular/material/sidenav';
// import { MatRadioModule } from '@angular/material/radio';
// import {MatExpansionModule} from '@angular/material/expansion';

const MAT_MODULES = [
    // FlexLayoutModule,
    // MatButtonModule,
    // MatCheckboxModule,
    // MatInputModule,
    // MatFormFieldModule,
    // ReactiveFormsModule,
    // MatIconModule,
    // MatMenuModule,
    // MatToolbarModule,
    // MatTabsModule,
    // MatRippleModule,
    MatDialogModule,
    // MatProgressSpinnerModule,
    // MatSidenavModule,
    // MatRadioModule,
    // MatExpansionModule,
    // MatCardModule
];

@NgModule({
    exports: MAT_MODULES
})
export class CustomMaterialModule { }
