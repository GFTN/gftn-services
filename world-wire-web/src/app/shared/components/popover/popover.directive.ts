// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Directive, ElementRef, ViewContainerRef, Input, TemplateRef } from '@angular/core';
import { DialogDirective, DialogService } from 'carbon-components-angular';
import { PopoverComponent } from './popover.component';

@Directive({
  selector: '[appPopover]',
  exportAs: 'appPopover',
  providers: [
    DialogService
  ]
})
export class PopoverDirective extends DialogDirective {

  @Input() appPopover: string | TemplateRef<any>;

  @Input() container: ElementRef;

  // Optional: user can specify a width for their popover
  @Input() dialogWidth: string;

  constructor(
    protected elementRef: ElementRef,
    protected viewContainerRef: ViewContainerRef,
    protected dialogService: DialogService
  ) {
    super(elementRef, viewContainerRef, dialogService);
    dialogService.create(PopoverComponent);
  }

  onDialogInit() {

    this.dialogConfig.content = this.appPopover;

    // custom data attributes to set container ref
    // and dialog width for re-positioning
    this.dialogConfig.data['container'] = this.container;

    if (this.dialogWidth) {
      this.dialogConfig.data['dialogWidth'] = this.dialogWidth;
    }
  }

  toggle() {
    if (this.expanded) {
      // trigger close animation
      this.dialogService.dialogRef.instance.show = false;

      // stagger destruction of component for animation
      setTimeout(() => {
        super.toggle();
      }, 300);
    } else {
      super.toggle();
    }
  }
}
