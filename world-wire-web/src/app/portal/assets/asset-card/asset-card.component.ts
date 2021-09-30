// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { Asset, AssetBalance } from '../../../shared/models/asset.interface';
import { AccountService } from '../../shared/services/account.service';
import { ModalService } from 'carbon-components-angular';
import { AssetDetailsModalComponent } from '../asset-details-modal/asset-details-modal.component';

@Component({
  selector: 'app-asset-card',
  templateUrl: './asset-card.component.html',
  styleUrls: ['./asset-card.component.scss']
})

export class AssetCardComponent implements OnInit {

  @Input() asset: Asset;

  @Output() modalOpened = new EventEmitter<true>();

  loaded = false;

  constructor(
    private accountService: AccountService,
    private modalService: ModalService,
  ) { }

  async ngOnInit() {
    const balanceObject: AssetBalance[] = await this.accountService.getDOAssetBalances(this.asset.asset_code);

    // get init balance for asset overview
    this.asset.balance = balanceObject ? parseFloat(balanceObject[0].balance) : 0.00;

    this.loaded = true;
  }

  /**
   * Opens modal with asset balance details
   *
   * @memberof AssetCardComponent
   */
  openDetailsModal() {
    this.modalService.create({
      component: AssetDetailsModalComponent,
      inputs: {
        MODAL_DATA: this.asset
      }
    });
  }

  /**
   * Opens modal with asset approval detail
   *
   * @memberof AssetCardComponent
   */
  openApprovalModal() {
    this.modalOpened.emit();
  }
}
