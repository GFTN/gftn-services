// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

import { DropdownComponent } from './components/dropdown/dropdown.component';
import { IconTooltipComponent } from './components/icon-tooltip/icon-tooltip.component';
import { BackToTopComponent } from './components/back-to-top/back-to-top.component';
import { PopoverComponent } from './components/popover/popover.component';
import { PopoverDirective } from './components/popover/popover.directive';
import { FlexLayoutModule } from '@angular/flex-layout';
import { NotificationTileComponent } from './components/notification-tile/notification-tile.component';

export { DropdownComponent } from './components/dropdown/dropdown.component';
export { IconTooltipComponent } from './components/icon-tooltip/icon-tooltip.component';
export { BackToTopComponent } from './components/back-to-top/back-to-top.component';
export { PopoverComponent } from './components/popover/popover.component';

@NgModule({
  declarations: [
    DropdownComponent,
    IconTooltipComponent,
    BackToTopComponent,
    PopoverComponent,
    PopoverDirective,
    NotificationTileComponent
  ],
  exports: [
    DropdownComponent,
    IconTooltipComponent,
    BackToTopComponent,
    PopoverComponent,
    PopoverDirective,
    NotificationTileComponent
  ],
  entryComponents: [
    PopoverComponent
  ],
  imports: [
    CommonModule,
    FlexLayoutModule,
    RouterModule
  ]
})
export class CustomCarbonModule { }
