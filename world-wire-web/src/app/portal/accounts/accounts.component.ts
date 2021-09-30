// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, HostBinding, OnDestroy, isDevMode } from '@angular/core';
import { UrlSegment, ActivatedRoute, Router, Data, NavigationEnd } from '@angular/router';
import { map, filter } from 'rxjs/operators';
import { startCase } from 'lodash';
import { AccountService } from '../shared/services/account.service';
import { SessionService } from '../../shared/services/session.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-accounts',
  templateUrl: './accounts.component.html',
  styleUrls: ['./accounts.component.scss']
})


export class AccountsComponent implements OnInit, OnDestroy {

  title: string;
  urlSegments: UrlSegment[];

  participantSubscription: Subscription;

  init = false;

  constructor(
    private activatedRoute: ActivatedRoute,
    private router: Router,
    public accountService: AccountService,
    private sessionService: SessionService
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
        if (isDevMode) {
          // console.log('changed route', route);
        }
        this.getTitle(route.snapshot.data);

        this.accountService.accountSlug = route.parent.snapshot.params.slug ? route.parent.snapshot.params.slug : null;

        this.urlSegments = route.snapshot.url;

        // propogate init event to listeners, if participant is already retrieved
        if (this.init) {
          setTimeout(() => {
            this.accountService.propogateParticipantChange();
          });
        }
      });
  }

  @HostBinding('attr.class') cls = 'flex-fill';

  ngOnInit() {

    // listen for switching of participant
    this.participantSubscription = this.sessionService.currentNodeChanged.subscribe(() => {

      this.init = true;

      // only propgate change if successful request was made
      if (!this.accountService.participantDetails || (this.accountService.participantDetails &&
        this.sessionService.currentNode.participantId !== this.accountService.participantDetails.id)) {
        // retrieve by promise to wait for data to come back from server
        this.accountService.getParticipant().then(() => {
          this.accountService.propogateParticipantChange();
        });
      }
    });
  }

  ngOnDestroy() {
    this.participantSubscription.unsubscribe();
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
    let url = '';

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
  getBreadcrumbTitle(path: string, index: number): string {
    if (index === (this.urlSegments.length - 1)) {
      return this.title;
    }
    return startCase(path);
  }

}
