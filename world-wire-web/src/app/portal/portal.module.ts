// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// Module Imports
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { PortalRoutingModule } from './portal-routing.module';
import { FlexLayoutModule } from '@angular/flex-layout';
import { FormsModule } from '@angular/forms';
import { HeaderFooterModule } from '../shared/header-footer.module';
import { NotificationsModule } from '../shared/notification.module';
import {
    TableModule, PlaceholderModule,
    LoadingModule, ModalModule,
    DialogModule, DropdownModule,
    ToggleModule, CodeSnippetModule,
    TabsModule, NotificationModule
} from 'carbon-components-angular';
import { CustomCarbonModule } from '../shared/custom-carbon.module';
import { PortalSharedModule } from './shared/portal-shared.module';

// Component Imports
import { PortalComponent } from './portal.component';
import { AccountComponent } from './account/account.component';
import { AccountsComponent } from './accounts/accounts.component';
import { DiligenceComponent } from './diligence/diligence.component';
import { UsersComponent } from './users/users.component';
import { BillingComponent } from './billing/billing.component';
import { TransactionsComponent } from './transactions/transactions.component';
import { PortalSideNavComponent } from './portal-side-nav/portal-side-nav.component';
import { OverviewComponent } from './overview/overview.component';
import { SettingsComponent } from './settings/settings.component';
import { ExportModalComponent } from './export-modal/export-modal.component';
import { TransactionsSettingsComponent } from './transactions/transactions-settings/transactions-settings.component';
import { TransactionsFiltersComponent } from './transactions/transactions-filters/transactions-filters.component';
import { TokenComponent } from './token/token.component';
import { TokenDialogComponent } from './token/token-dialog/token-dialog.component';
import { SelectComponent } from './select/select.component';
import { StatusHistoryModalComponent } from './status-history-modal/status-history-modal.component';
import { StatusDetailsModalComponent } from './status/status-details-modal/status-details-modal.component';
import { StatusDetailsComponent } from './status/status-details/status-details.component';
import { StatusHistoryComponent } from './status/status-history/status-history.component';
import { StatusOverviewComponent } from './status/status-overview/status-overview.component';

import { LogService } from './shared/services/log.service';
import { ExportService } from './shared/services/export.service';
import { TransactionService } from './shared/services/transaction.service';
import { ParticipantPermissionsService } from '../shared/services/participant-permissions.service';
import { AccountService } from './shared/services/account.service';
import { AssetsComponent } from './assets/assets.component';

@NgModule({
    declarations: [
        PortalComponent,
        OverviewComponent,
        AccountComponent,
        AccountsComponent,
        AssetsComponent,
        DiligenceComponent,
        UsersComponent,
        BillingComponent,
        TransactionsComponent,
        SettingsComponent,
        PortalSideNavComponent,
        ExportModalComponent,
        TransactionsSettingsComponent,
        TransactionsFiltersComponent,
        TokenComponent,
        TokenDialogComponent,
        SelectComponent,
        StatusDetailsModalComponent,
        StatusDetailsComponent,
        StatusHistoryComponent,
        StatusHistoryModalComponent,
        StatusOverviewComponent,
    ],
    imports: [
        CommonModule,
        FormsModule,
        FlexLayoutModule,
        HeaderFooterModule,
        // Carbon Components Angular
        NotificationsModule,
        // Carbon Angular Modules
        NotificationModule,
        CodeSnippetModule,
        TableModule,
        PlaceholderModule,
        LoadingModule,
        ModalModule,
        DialogModule,
        DropdownModule,
        ToggleModule,
        TabsModule,
        // Custom Carbon Components
        CustomCarbonModule,
        PortalSharedModule,
        // Routing Module must come last
        PortalRoutingModule,
    ],
    providers: [
        LogService,
        ExportService,
        TransactionService,
        AccountService,
        ParticipantPermissionsService
    ],
    entryComponents: [
        TokenDialogComponent,
        ExportModalComponent,
        StatusHistoryModalComponent,
        StatusDetailsModalComponent
    ]
})
export class PortalModule { }
