// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, ElementRef, ViewChild, OnInit } from '@angular/core';
import { VersionService } from '../services/version.service';
import { ListItem } from 'carbon-components-angular';
import { toArray } from 'lodash';
import { Router } from '@angular/router';

@Component({
  selector: 'app-version-select',
  templateUrl: './version-select.component.html',
  styleUrls: ['./version-select.component.scss']
})
export class VersionSelectComponent implements OnInit {

  private accordionInitialized = false;
  @ViewChild('accordionVersion') accordionVersion: ElementRef;

  routeVersion: string;
  open = false;

  versionOptions: ListItem[] = [];

  constructor(
    private router: Router,
    public versionService: VersionService
  ) { }

  ngOnInit(): void {
    // init version to current version per service
    this.routeVersion = this.versionService.current ? this.versionService.current.version : null;

    for (const version of toArray(this.versionService.versions)) {
      const option: ListItem = {
        content: version.version,
        selected: false,
      };

      // initializing to current version
      option.selected = version === this.versionService.current ? true : false;

      this.versionOptions.push(option);
    }
  }

  /**
   *
   *
   * @param {string} version
   * @memberof VersionSelectComponent
   */
  onSelect(event) {

    const item: ListItem = event.item;

    // prevent re-route to on refresh of page
    if (this.versionService.current.version !== this.versionService.current.version) {
      this.setVersion(item.content);
    }
  }

  /**
   * Update current version stored per version
   *
   * @param {string} version
   * @memberof VersionSelectComponent
   */
  private setVersion(version: string) {

    this.versionService.current = this.versionService.getVersion(version);

    // route to current version
    this.router.navigate([`/docs/${version}`]);
  }

}
