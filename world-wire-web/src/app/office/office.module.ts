// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FlexLayoutModule } from '../../../node_modules/@angular/flex-layout';
import { OfficeRoutingModule } from './office-routing.module';

import { OfficeComponent } from './office.component';
import { OverviewComponent } from './overview/overview.component';
import { AccountsComponent } from './accounts/accounts.component';

import { CustomMaterialModule } from '../shared/custom-material.module';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HeaderFooterModule } from '../shared/header-footer.module';
import { OfficeSideNavComponent } from './office-side-nav/office-side-nav.component';
// import { AccountComponent } from './account/account.component';
import { NotificationsModule } from '../shared/notification.module';
import { AccountModule } from './account/account.module';
import { UsersComponent } from './users/users.component';
import { BlocklistManagementComponent } from './blocklist-management/blocklist-management.component';
import { BlocklistDialogComponent } from './blocklist-management/blocklist-dialog/blocklist-dialog.component';
import { ButtonModule, LoadingModule } from 'carbon-components-angular';

@NgModule({
    declarations: [
        OfficeComponent,
        OverviewComponent,
        AccountsComponent,
        OfficeSideNavComponent,
        // AccountComponent,
        UsersComponent,
        BlocklistManagementComponent,
        BlocklistDialogComponent
    ],
    imports: [
        CommonModule,
        OfficeRoutingModule,
        CustomMaterialModule,
        // Carbon Angular Modules
        ButtonModule,
        LoadingModule,
        FlexLayoutModule,
        NotificationsModule,
        FormsModule,
        AccountModule,
        HeaderFooterModule,
        ReactiveFormsModule
    ],
    entryComponents: [BlocklistDialogComponent],
})
export class OfficeModule { }
