// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, EventEmitter, Output, AfterViewInit, ViewChild, ElementRef } from '@angular/core';
import { VersionService } from '../../../shared/services/version.service';
import { SwaggerService } from '../../../shared/services/swagger.service';
import { ViewportScroller } from '@angular/common';
import { RegexPipe } from '../../../../shared/pipes/regex.pipe';
import { Accordion } from 'carbon-components';
import { Router } from '@angular/router';
// import { Accordion } from 'carbon-components-angular';

/**
 * Used to display the side-nav in the docs
 * containing the version-select and docs links
 *
 * @export
 * @class VersionNavComponent
 * @implements {OnInit}
 * @implements {AfterViewInit}
 */
@Component({
  selector: 'app-version-nav',
  templateUrl: './version-nav.component.html',
  styleUrls: ['./version-nav.component.scss'],
  providers: [RegexPipe]
})
export class VersionNavComponent implements OnInit, AfterViewInit {

  private accordionInitialized = false;

  // close the version side-nav in version.component.ts
  @Output() closeSideNavModal = new EventEmitter<any>();

  // list of apis in in side-nav
  @ViewChild('accordionApis') accordionApis: ElementRef;

  // list of models in associated with the selected api
  @ViewChild('accordionModels') accordionModels: ElementRef;

  constructor(
    public versionService: VersionService,
    public swaggerService: SwaggerService
  ) { }

  ngOnInit() {
    this.accordionInitialized = false;
    // console.log(this.swaggerService._jsonSpec.definitions);
  }

  ngAfterViewInit() {

    // initialize accordion only once
    if (this.accordionInitialized === false) {

      Accordion.create(this.accordionApis.nativeElement);
      Accordion.create(this.accordionModels.nativeElement);

      this.accordionInitialized = true;

    }

  }

  /**
   * close side nav in mobile view
   *
   * @memberof VersionNavComponent
   */
  closeNav() {
    this.closeSideNavModal.emit();
  }

}
