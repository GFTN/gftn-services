// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit } from '@angular/core';
import { Modal } from 'carbon-components';
import { DocumentService } from '../../shared/services/document.service';
import * as _ from 'lodash';
import { Router } from '@angular/router';
import { UtilsService } from '../../shared/utils/utils';

@Component({
  templateUrl: './secure-version.component.html',
  styleUrls: ['./secure-version.component.scss']
})
export class SecureVersionComponent implements OnInit {

  errorMsg: string;
  model: { password: string };

  constructor(
    private utils: UtilsService,
    private router: Router
  ) { }

  ngOnInit() {

    this.errorMsg = '';
    this.model = { password: '' };
  }

  onSubmit() {

    this.errorMsg = '';

    // if password matches
    if (this.model.password === 'NobleMission') {

      const dateString: string = new Date(new Date().getFullYear(), new Date().getMonth() + 3, new Date().getDate()).toUTCString();
      // create cookie
      this.utils.setCookie('wwdocspermission', dateString);

      this.router.navigate(['/docs']);
    } else {
      this.errorMsg = 'Incorrect password';
    }

  }

}
