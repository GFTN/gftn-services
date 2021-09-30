// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, TemplateRef, ElementRef, ViewChild, AfterViewInit, HostListener } from '@angular/core';
import { Dialog } from 'carbon-components-angular';
import { position } from 'carbon-components-angular/utils/position';
import { slideDownAnimation } from '../../animations';
import { TransactionService } from '../../../portal/shared/services/transaction.service';
import * as _ from 'lodash';

@Component({
  selector: 'app-popover',
  templateUrl: './popover.component.html',
  styleUrls: ['./popover.component.scss'],
  animations: [
    slideDownAnimation
  ]
})

/**
 * Custom Popover/Dropdown Menu
 * to present external filters, settings, etc.
 */
export class PopoverComponent extends Dialog implements AfterViewInit {

  // content can be a template
  hasContentTemplate = false;

  // toggles whether the initial positioning of the Dialog
  // overflowed onto the page
  overFlow = false;

  showOverflow = false;

  scrolledDown = false;

  adjustedYPos = 40;

  show = true;

  // dropdown arrow pointing to the clicked element
  @ViewChild('arrow') arrow: ElementRef;

  constructor(
    elementRef: ElementRef,
    protected transactionService: TransactionService
  ) {
    super(elementRef);
    this.transactionService.popoverRef = this;
  }

  onDialogInit() {

    this.hasContentTemplate = this.dialogConfig.content instanceof TemplateRef;

    // calculate width from settings,
    // else it will default to full width of container
    if (this.dialogConfig.data['dialogWidth']) {
      this.dialog.nativeElement.style.width = this.dialogConfig.data['dialogWidth'];
    }

    // reposition to center arrow
    this.addGap['left-bottom'] = pos => {
      return position.addOffset(pos, this.adjustedYPos, 50);
    };

    // programmatically show overflow
    setTimeout(() => {
      this.showOverflow = true;
    }, 300);
  }

  ngAfterViewInit() {

    // inherit from dialog for initial placement
    super.ngAfterViewInit();

    const containerEl: ElementRef = this.dialogConfig.data['container'];

    const containerPos = (containerEl !== undefined) ? containerEl.nativeElement.offsetLeft : 0;


    // get absolute positioning of dialog
    // has to be gathered after DOM loads (ngAfterViewInit)
    // in order to get correctly computed values from DOM
    const absPos = position.findPosition(this.dialogConfig.parentRef.nativeElement, this.dialog.nativeElement, 'left-bottom');
    const dialogLeft = position.getPlacementBox(this.dialog.nativeElement, absPos).left;

    const parentLeft = this.dialogConfig.parentRef.nativeElement.offsetLeft;
    const containerWidth = containerEl.nativeElement.offsetWidth;

    const dialogWidth = this.dialog.nativeElement.offsetWidth;

    // get left coordinate to account for container margin
    if (dialogLeft < containerPos) {
      this.overFlow = true;
    }

    // reposition for overflow
    if (this.overFlow) {

      // dialog is bigger than or equal in size to the container
      if (dialogWidth >= containerWidth) {
        this.addGap['left-bottom'] = pos => {

          // account for initial positioning (relative to parent)
          return position.addOffset(pos, this.adjustedYPos, containerWidth - parentLeft);
        };
        this.arrow.nativeElement.style.right = (containerWidth - parentLeft - 24) + 'px';
      } else {
        this.addGap['left-bottom'] = pos => {

          // account for amount overflowed off viewport
          return position.addOffset(pos, this.adjustedYPos, -(dialogLeft));
        };
      }

      // reverse arrow if dialog is on other side
      if (Math.abs(dialogLeft) >= dialogWidth / 2 && this.dialogConfig.data['dialogWidth']) {
        this.arrow.nativeElement.classList.add('left');
      }

      // re-position dialog manually after computation
      this.placeDialog();
    }
  }

  // override default doc.click listener
  // to account for clicks emitted from nested components
  @HostListener('document:click', ['$event'])
  clickClose(event) {

    const container = this.dialog.nativeElement;

    // account for date pickers appended to body
    const datePickers = document.getElementsByClassName('bx--date-picker__calendar');

    let datePickerClicked = false;

    for (const picker of _.toArray(datePickers)) {
      if (picker.contains(event.target)) {
        datePickerClicked = true;
      }
    }

    // Pulled from Carbon Angular and extended to exclude additional elements
    if (!(container === event.target)
      && !container.contains(event.target)
      // exclude click events from dropdowns since they emit the same close event
      && !event.target.classList.contains('bx--list-box__menu-item')
      // exclude click events from datepicker since they propogate outside of container
      && !datePickerClicked
      && !this.dialogConfig.parentRef.nativeElement.contains(event.target)
    ) {
      // destroy component only after animation finished
      this.doClose();
    }
  }

  doClose() {

    // toggle show for animation
    this.show = false;

    // stagger destruction of component for animation
    setTimeout(() => {
      super.doClose();
    }, 300);
  }
}
