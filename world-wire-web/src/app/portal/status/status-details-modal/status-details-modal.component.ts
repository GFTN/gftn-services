// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Inject } from '@angular/core';
import { BaseModal } from 'carbon-components-angular/modal/base-modal.class';
import { StatusByDate } from '../../../shared/models/log.interface';

@Component({
  selector: 'app-status-details-modal',
  templateUrl: './status-details-modal.component.html',
  styleUrls: ['./status-details-modal.component.scss']
})
export class StatusDetailsModalComponent extends BaseModal implements OnInit {

  detailsObject: StatusByDate;

  constructor(@Inject('MODAL_DATA') public data: StatusByDate) {
    super();

    this.detailsObject = data;

  }

  ngOnInit() {
  }

}
