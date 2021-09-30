// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Inject } from '@angular/core';
import { BaseModal } from 'carbon-components-angular';
import { Asset, Obligation } from '../../../shared/models/asset.interface';
import { AccountService } from '../../shared/services/account.service';
import { SessionService } from '../../../shared/services/session.service';
import { UtilsService } from '../../../shared/utils/utils';

@Component({
  selector: 'app-asset-details-modal',
  templateUrl: './asset-details-modal.component.html',
  styleUrls: ['./asset-details-modal.component.scss']
})

export class AssetDetailsModalComponent extends BaseModal implements OnInit {

  // stores current asset details being viewed
  currentAsset: Asset = {
    asset_code: '',
    asset_type: 'DO'
  };

  obligationsLoaded = false;

  obligationBalances: Obligation[];

  constructor(
    @Inject('MODAL_DATA') public data: Asset,
    private sessionService: SessionService,
    private accountService: AccountService,
    public utils: UtilsService
  ) {
    // MUST do super() to extend BaseModal
    super();

    this.currentAsset = data;
  }

  ngOnInit() {
    this.getOutstandingBalances();
  }

  /**
   * Get outstanding balances of the asset held by other participants
   *
   * @returns {Promise<void>}
   * @memberof AssetDetailsModalComponent
   */
  async getOutstandingBalances(): Promise<void> {

    this.obligationsLoaded = false;

    // Get list of holders of this issued asset
    try {
      this.obligationBalances = this.obligationBalances || await this.accountService.getDOBalanceDetails(this.currentAsset.asset_code);

      // Filter out current participant being viewed
      this.obligationBalances = this.obligationBalances.filter((obligation: Obligation) => {
        return obligation.participant_id !== this.sessionService.currentNode.participantId;
      });

      this.obligationsLoaded = true;

    } catch (err) {

      this.obligationsLoaded = true;
    }
  }
}
