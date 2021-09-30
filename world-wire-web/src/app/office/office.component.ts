// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit } from '@angular/core';
import { Modal } from 'carbon-components';
import { AuthService } from '../shared/services/auth.service';
import { ActivatedRoute, Router, NavigationEnd } from '@angular/router';
import { map, filter, mergeMap } from 'rxjs/operators';

@Component({
  // selector: 'app-root',
  templateUrl: './office.component.html',
  styleUrls: ['./office.component.scss']
})
export class OfficeComponent implements OnInit {

  modal: any;
  title: string;

  constructor(
    public authService: AuthService,
    private activatedRoute: ActivatedRoute,
    private router: Router
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
