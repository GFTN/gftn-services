// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, HostBinding } from '@angular/core';

/**
 * Used by server to redirect users who login with either
 * the wrong email address or an account that has not
 * yet been provisioned
 *
 * @export
 * @class UnrecognizedComponent
 * @implements {OnInit}
 */
@Component({
  selector: 'app-unrecognized',
  templateUrl: './unrecognized.component.html',
  styleUrls: ['./unrecognized.component.scss']
})
export class UnrecognizedComponent implements OnInit {

  constructor() { }

  @HostBinding('attr.class') cls = 'flex-fill';

  ngOnInit() {
  }

}
