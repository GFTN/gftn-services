// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit } from '@angular/core';
import { Modal } from 'carbon-components';
import { AuthService } from '../shared/services/auth.service';

@Component({
  selector: 'app-site-header',
  templateUrl: './site-header.component.html',
  styleUrls: ['./site-header.component.scss']
})
export class SiteHeaderComponent implements OnInit {

  sideNavModal: any;

  constructor(public afAuth: AuthService) { }

  ngOnInit() {
    const modalSideNavElement = document.getElementById('modal-side-nav');
    this.sideNavModal = Modal.create(modalSideNavElement);
  }

  openSideNavModal() {
    this.sideNavModal.show();
  }

  closeSideNavModal() {
    this.sideNavModal.hide();
  }

}
