// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, HostBinding } from '@angular/core';
import { ActivatedRoute, Router, NavigationEnd, Data, UrlSegment } from '@angular/router';
import { map, filter } from 'rxjs/operators';
import { IInstitution } from '../../shared/models/participant.interface';
import { SessionService } from '../../shared/services/session.service';
import { startCase } from 'lodash';

@Component({
  templateUrl: './account.component.html',
  styleUrls: ['./account.component.scss']
})
export class AccountComponent implements OnInit {

  title: string;
  urlSegments: UrlSegment[];

  constructor(
    private activatedRoute: ActivatedRoute,
    private router: Router,
    public sessionService: SessionService
  ) {

    // subscribing to router events
    // to set page title
    this.router.events.pipe(
      filter(event => event instanceof NavigationEnd),
      map(() => this.activatedRoute),
      map((route) => {
        while (route.firstChild) {
          route = route.firstChild;
        }
        return route;
      })
    )
      .pipe(
        filter((route) => route.outlet === 'primary'),
      )
      .subscribe((route: ActivatedRoute) => {
        this.getTitle(route.snapshot.data);
        this.urlSegments = route.snapshot.url;
      });
  }

  @HostBinding('attr.class') cls = 'flex-fill';

  ngOnInit() {
  }

  getCurrentSlug() {
    return this.sessionService.institution ? this.sessionService.institution.info.slug : '';
  }

  /**
   * Reusable function to get title and shortTile from route data
   * @param data
   */
  getTitle(data: Data) {

    if (data.shortTitle) {
      this.title = data.shortTitle;
    } else if (data.title && !data.shortTitle) {
      this.title = data.title;
    }
  }

  /**
   * Construct url string for breadcrumb
   * @param maxIndex
   */
  constructUrl(maxIndex: number): string {
    let url = '/office/account/' + this.sessionService.institution.info.slug;

    for (let i = 0; i < (maxIndex + 1); i++) {
      url = url + '/' + this.urlSegments[i].path;
    }
    return url;
  }

  /**
   * Get title for breadcrumb
   * @param path
   * @param index
   */
  getBreadcrumbTitle(path: string, index: number = 0): string {
    if (this.urlSegments.length === 0 || index === (this.urlSegments.length - 1)) {
      return this.title;
    }
    return startCase(path);
  }
}
