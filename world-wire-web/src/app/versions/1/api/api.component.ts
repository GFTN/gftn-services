// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// NOTE: Cannot parse WW yaml api definitions
// properly because the yaml contains "$ref" to other yaml
// files in repo as such all Yaml files must be complied
// to usable json as part of repo

// TODO:  research how to dynamically "get" json using
// https://github.com/k33g/gh3, using github API (or using angular http.get request)
// this.getApiDef('./assets/api-def/' + this.apiVersion + '/' + this.apiName + '.json');

import { Component, OnInit, OnDestroy, HostBinding, AfterViewInit, OnChanges } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { SwaggerService } from '../../shared/services/swagger.service';
import { Dropdown } from 'carbon-components';
import { ViewportScroller } from '@angular/common';

@Component({
  templateUrl: './api.component.html',
  styleUrls: ['./api.component.scss'],
})
export class ApiComponent implements AfterViewInit {

  // fileName: string;
  // loaded: boolean;

  constructor(
    public swaggerService: SwaggerService,
    private activatedRoute: ActivatedRoute,
    private viewportScroller: ViewportScroller
  ) { }

  @HostBinding('attr.class') cls = 'flex-fill';

  ngAfterViewInit() {

    // init carbon component dropdown
    // const dropDownElm = document.getElementById('participant-dropdown');
    // Dropdown.create(dropDownElm);

    this.activatedRoute
      .queryParams
      .subscribe(queryParams => {

        // console.log('Query Params:', queryParams);
        if (queryParams['jump']) {

          setTimeout(() => {

            // jump to anchor with the related jump queryParams
            this.viewportScroller.scrollToAnchor(queryParams['jump']);

          });

        }

      });

  }

}
