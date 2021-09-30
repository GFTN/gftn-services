// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit } from '@angular/core';
import { Modal } from 'carbon-components';
import { DocumentService } from '../../shared/services/document.service';
import * as _ from 'lodash';
import { UtilsService } from '../../shared/utils/utils';

@Component({
  templateUrl: './version.component.html',
  styleUrls: ['./version.component.scss']
})
export class VersionComponent implements OnInit {

  sideNavModal: any;
  termsModal: any;
  apiOptions = false;

  constructor(
    private utils: UtilsService
  ) { }

  ngOnInit() {

    const modalSideNavElement = document.getElementById('modal-side-nav');
    this.sideNavModal = Modal.create(modalSideNavElement);

    const modalTermsElement = document.getElementById('modal-terms');
    modalTermsElement.addEventListener('modal-beinghidden', (event: any) => {

      // console.log('modal-beingshown fired:', event);

      // use _.get() to prevent error when launching element not present
      if (_.get(event, 'detail.launchingElement.id') === 'acceptTermsBtn') {
        this.acceptOptInCookie(true);
      } else {
        // prevent closing modal unless it is button click
        event.preventDefault();
      }

      return false;

    });

    this.termsModal = Modal.create(modalTermsElement);

    this.checkOptInCookie();

    // console.log('version');

  }

  checkOptInCookie() {

    const accepted = this.utils.getCookie('wwOptIn');

    if (accepted !== 'true') {
      // console.log('cookie', this.documentService.documentRef.cookie);
      this.openTermsModal();
    }

  }

  acceptOptInCookie(accept: boolean) {

    const dateString: string = new Date(new Date().getFullYear(), new Date().getMonth() + 3, new Date().getDate()).toUTCString();

    // create cookie
    this.utils.setCookie('wwOptIn', dateString);

    this.termsModal.hide();

  }

  openSideNavModal() {
    this.sideNavModal.show();
  }

  closeSideNavModal() {
    this.sideNavModal.hide();
  }

  openTermsModal() {
    this.termsModal.show();
  }
}
