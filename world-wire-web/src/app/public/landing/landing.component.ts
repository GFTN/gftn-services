// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import {
  Component,
  OnInit,
  HostBinding
} from '@angular/core';


@Component({
  templateUrl: './landing.component.html',
  styleUrls: ['./landing.component.scss']
})
export class LandingComponent implements OnInit {

  @HostBinding('attr.class') cls = 'flex-fill';

  ngOnInit() { }

}
