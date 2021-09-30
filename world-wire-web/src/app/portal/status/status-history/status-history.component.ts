// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, HostBinding } from '@angular/core';
import { LogService } from '../../shared/services/log.service';
import { ActivatedRoute } from '@angular/router';
import { Service, StatusByDate } from '../../../shared/models/log.interface';
import { startCase, toLower, find } from 'lodash';
import { SessionService } from '../../../shared/services/session.service';

@Component({
  selector: 'app-status-history',
  templateUrl: './status-history.component.html',
  styleUrls: ['./status-history.component.scss']
})
export class StatusHistoryComponent implements OnInit {

  serviceName: string;

  displayName: string;

  currentService: Service;

  existingErrors: StatusByDate[] = [];

  constructor(
    private logService: LogService,
    private activatedRoute: ActivatedRoute,
    public sessionService: SessionService
  ) { }

  @HostBinding('attr.class') cls = 'flex-fill';

  ngOnInit() {

    // grab service name from url
    this.serviceName = this.activatedRoute.snapshot.params.serviceName;

    this.displayName = this.serviceName.replace('/-/g', ' ');
    this.displayName = startCase(toLower(this.displayName));

    // convert Api to API
    this.displayName = this.displayName.replace('Api', 'API');

    this.currentService = find(this.logService.services, (s: Service) => {
      return s.name === this.serviceName;
    });

    if (this.currentService) {
      // iterate through logs to see if errors exists
      for (const day of this.currentService.errorHistory) {
        if (day.errors.length > 0) {

          // list only existing errors
          this.existingErrors.push(day);
        }
      }

      this.existingErrors.reverse();
    }
  }
}
