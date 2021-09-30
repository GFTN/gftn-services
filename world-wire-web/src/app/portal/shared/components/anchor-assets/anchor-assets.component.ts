// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Input } from '@angular/core';
import { AccountService } from '../../services/account.service';

@Component({
  selector: 'app-anchor-assets',
  templateUrl: './anchor-assets.component.html',
  styleUrls: ['./anchor-assets.component.scss']
})
export class AnchorAssetsComponent implements OnInit {

  @Input() issuedAssetsLoaded = false;

  @Input() error = false;

  constructor(
    public accountService: AccountService,
  ) { }

  ngOnInit() {
  }

}
