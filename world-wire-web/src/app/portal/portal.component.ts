// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit } from '@angular/core';
import { Modal } from 'carbon-components';
import { ActivatedRoute, Router, NavigationEnd } from '@angular/router';
import { map, filter, mergeMap } from 'rxjs/operators';
import { SessionService } from '../shared/services/session.service';

@Component({
  templateUrl: './portal.component.html',
  styleUrls: ['./portal.component.scss']
})
export class PortalComponent implements OnInit {

  modal: any;
  title: string;

  constructor(
    public sessionService: SessionService,
    private activatedRoute: ActivatedRoute,
    private router: Router
  ) {

    // set page title
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
        mergeMap((route) => route.data)
      )
      .subscribe((event) => {
        if (event['title']) {
          this.title = event['title'];
        }
      });
  }

  ngOnInit() {
    const modalElement = document.getElementById('modal-side-nav');
    this.modal = Modal.create(modalElement);
  }

  openModal() {
    this.modal.show();
  }
}
