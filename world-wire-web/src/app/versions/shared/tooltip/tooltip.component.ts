// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Tooltip } from 'carbon-components';
import { Component, OnInit, Input, AfterViewInit } from '@angular/core';

@Component({
  selector: 'app-tooltip',
  templateUrl: './tooltip.component.html',
  styleUrls: ['./tooltip.component.scss']
})
export class TooltipComponent implements OnInit, AfterViewInit {
  @Input() termName: string;
  @Input() termDefinition: string;
  @Input() icon = false;

  termId: string;
  termsTooltip: any;
  toolTipDefinition: HTMLElement;
  tooltipTermParent: HTMLElement;
  overflow = false;
  active = false;

  constructor() { }

  ngOnInit() {
    this.termId = this.termName
      .replace(/\s+/g, '-') // remove spaces
      .replace(/\(.+?\)/g, '') // remove abbreviations
      .toLowerCase();
  }

  ngAfterViewInit() {
    const toolTipTerm = document.getElementById(`${this.termId}-term`);
    this.toolTipDefinition = document.getElementById(`${this.termId}-def`);

    // initialize tooltip
    this.termsTooltip = Tooltip.create(toolTipTerm);

    const termDefLength = this.termDefinition.length;
    // get parent element of tooltip
    this.tooltipTermParent = this.toolTipDefinition.parentElement;

    // set default width dynamically -
    // width = ratio of term defintion length to term word length
    let width = termDefLength / this.tooltipTermParent.offsetWidth;

    // minimum width = 100%
    if (width <= 1) {
      width = 1;
    }

    this.toolTipDefinition.style.width = (width * 100) + '%';
  }

  toggleTooltip() {
    const rect = this.tooltipTermParent.getBoundingClientRect();
    const rect2 = this.toolTipDefinition.getBoundingClientRect();
    // get right positioning of tooltip term and defintion
    const tooltipOffset = this.tooltipTermParent.offsetWidth + rect.left;
    const tooltipDefOffset = this.toolTipDefinition.offsetWidth + rect2.left;

    // reposition tooltip based on window size
    if (tooltipDefOffset > window.innerWidth && this.overflow === false) {
      this.toolTipDefinition.style.right = -(window.innerWidth - tooltipOffset) + 'px';
      this.overflow = true;
    } else {
      this.toolTipDefinition.style.right = 'auto';
      this.overflow = false;
    }

    // show tooltip
    this.active = !this.active;
  }
}
