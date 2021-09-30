// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Input, ViewChild, ElementRef } from '@angular/core';

@Component({
  selector: 'app-icon-tooltip',
  templateUrl: './icon-tooltip.component.html',
  styleUrls: ['./icon-tooltip.component.scss']
})
export class IconTooltipComponent implements OnInit {

  iconPath = '/assets/icons/ibm/carbon-icons.svg#';
  active = false;

  @Input() icon = 'icon--info--outline';

  // required
  @Input() tooltipContent = 'Tooltip not defined.';

  constructor() { }

  ngOnInit() {
    if (this.icon != null) {
      this.iconPath = this.iconPath + this.icon;
    }
  }

  toggleTooltip() {
    this.active = !this.active;
  }
}
