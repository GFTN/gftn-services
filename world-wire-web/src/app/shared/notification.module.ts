// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NotificationComponent } from './components/notifications/notifications.component';
import { NotificationService } from './services/notification.service';
import { FlexLayoutModule } from '@angular/flex-layout';

/**
 * This shared module is used to import components and services related to managing users
 *
 * @export
 * @class UserModule
 */
@NgModule({
    imports: [
        CommonModule,
        FlexLayoutModule,
    ],
    declarations: [
        NotificationComponent
    ],
    exports: [
        NotificationComponent
    ],
    // entryComponents: [
    //     NotificationComponent
    // ],
    providers: [
        NotificationService
    ]
})
export class NotificationsModule { }
