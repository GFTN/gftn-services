// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Input } from '@angular/core';
import { StatusByDate } from '../../../shared/models/log.interface';

@Component({
  selector: 'app-status-details',
  templateUrl: './status-details.component.html',
  styleUrls: ['./status-details.component.scss']
})
export class StatusDetailsComponent implements OnInit {

  @Input() detailsObj: StatusByDate;

  constructor() { }

  ngOnInit() {
  }

}
