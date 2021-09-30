// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, HostBinding } from '@angular/core';
import { Router } from '@angular/router';
import { SessionService } from '../../shared/services/session.service';

@Component({
  selector: 'app-overview',
  templateUrl: './overview.component.html',
  styleUrls: ['./overview.component.scss']
})
export class OverviewComponent implements OnInit {

  constructor(
    private router: Router,
    protected sessionService: SessionService
  ) {
    /**
     * Remove once Overview page is done
     */
    this.router.navigate([`/portal/${this.sessionService.institution.info.slug}/status`]);
  }

  @HostBinding('attr.class') cls = 'flex-fill';

  ngOnInit() {
  }

}
