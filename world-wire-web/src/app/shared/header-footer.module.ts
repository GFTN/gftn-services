// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FlexLayoutModule } from '@angular/flex-layout';
import { RouterModule } from '@angular/router';
import {SiteHeaderComponent } from '../site-header/site-header.component';
import { PortalNavComponent } from '../portal-nav/portal-nav.component';
import { PublicNavComponent } from '../public-nav/public-nav.component';
import { PublicFooterComponent } from '../public-footer/public-footer.component';
import { CustomMaterialModule } from './custom-material.module';
import { SupportButtonComponent } from './components/support-button/support-button.component';

/**
 * This shared module imports all default header and footer components
 *
 * @export
 * @class HeaderFooterModule
 */
@NgModule({
    imports: [
        CommonModule,
        RouterModule,
        FlexLayoutModule,
        CustomMaterialModule
    ],
    declarations: [
        SiteHeaderComponent,
        PortalNavComponent,
        PublicNavComponent,
        PublicFooterComponent,
        SupportButtonComponent
    ],
    exports: [
        SiteHeaderComponent,
        PortalNavComponent,
        PublicNavComponent,
        PublicFooterComponent,
        SupportButtonComponent
    ]
})
export class HeaderFooterModule { }
