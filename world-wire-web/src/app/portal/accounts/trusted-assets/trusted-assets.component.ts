// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
//
import { Component, OnInit, ViewChild, ElementRef, OnDestroy, ComponentRef } from '@angular/core';
import { ModalService, TableModel, TableHeaderItem, TableItem } from 'carbon-components-angular';
import { TrustlineModalComponent } from '../trustline-modal/trustline-modal.component';
import { AccountService } from '../../shared/services/account.service';
import { CheckboxOption } from '../../../shared/models/checkbox-option.model';
import { ParticipantAccount } from '../../../shared/models/account.interface';
import { Asset } from '../../../shared/models/asset.interface';
import { filter } from 'lodash';
import { CheckboxGroupFilter } from '../../shared/models/filter.model';
import { TrustRequest } from '../../shared/models/trust-request.interface';
import { AngularFireDatabase } from '@angular/fire/database';
import { SessionService } from '../../../shared/services/session.service';
import { Approval } from '../../../shared/models/approval.interface';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-trusted-assets',
  templateUrl: './trusted-assets.component.html',
  styleUrls: ['./trusted-assets.component.scss']
})

export class TrustedAssetsComponent implements OnInit, OnDestroy {

  // holds current account details and full asset data
  currentAccount: ParticipantAccount;

  // holds currently viewable asset data in the table
  currentAssetData: Asset[];

  // toggles if asset data is ready
  public loaded = false;

  currSortIndex: number;

  // data model used for ibm-table
  model = new TableModel();

  // all possible filters
  filters: CheckboxGroupFilter;

  // necessary container ref for filter menu dropdown
  @ViewChild('container') containerRef: ElementRef;

  pendingRequests: TrustRequest[];

  participantSubscription: Subscription;

  notificationSubscription: Subscription;

  currentOpenModal: ComponentRef<TrustlineModalComponent>;

  constructor(
    public sessionService: SessionService,
    public accountService: AccountService,
    private modalService: ModalService,
    private db: AngularFireDatabase
  ) { }

  async ngOnInit() {

    // init asset filters. fixed list of types
    this.filters = {
      'issuer_id': {
        name: 'Asset Issuer',
        options: []
      },
      'asset_code': {
        name: 'Asset Code',
        options: []
      },
      'currency': {
        name: 'Currency',
        options: []
      }
    };

    // initializing model header for new empty view
    this.model.header = [
      new TableHeaderItem({
        data: 'Issuer'
      }),
      new TableHeaderItem({
        data: 'Asset Code'
      }),
      new TableHeaderItem({
        data: 'Currency'
      }),
      // TODO: add asset limit after https://github.com/GFTN/gftn-services/issues/893 is completed
      // new TableHeaderItem({
      //   data: 'Limit'
      // }),
      new TableHeaderItem({
        data: 'Current Balance'
      }),
    ];

    this.notificationSubscription = this.getNotifications();

    this.participantSubscription = this.accountService.currentParticipantChanged.subscribe(async () => {
      // get asset data for this page
      await this.initAssetData();

      this.loadDataView();
    });
  }

  ngOnDestroy() {
    // programmatically close modal if open
    if (this.currentOpenModal) {
      this.currentOpenModal.instance.closeModal();
    }
    this.participantSubscription.unsubscribe();

    this.notificationSubscription.unsubscribe();
  }

  /**
   * Get details for every trusted asset in this account
   *
   * @memberof TrustedAssetsComponent
   */
  async initAssetData(): Promise<void> {

    // SAMPLE DATA. Uncomment for testing in view
    // this.accountService.allAccounts = [];
    // this.accountService.allAccounts[this.accountService.accountSlug] = {
    //   address: 'XXXXXXXXXXXXX',
    //   name: this.accountService.accountSlug,
    // };

    if (this.accountService.allAccounts) {
      // get cached details of current account (if available)
      this.currentAccount = this.accountService.allAccounts[this.accountService.accountSlug];

      // no account found!
      if (!this.currentAccount) {
        return;
      }

      // try to get asset balances data again if not available
      if (!this.currentAccount.assets) {
        this.currentAccount.assets = await this.accountService.getTrustedAssetBalances(this.accountService.accountSlug);
      }

      // init all possible filters
      for (const asset of this.currentAccount.assets) {

        // issuer filter
        const foundIssuer = this.filters['issuer_id'].options.find((option: CheckboxOption) => {
          return option.name === asset.issuer_id;
        });

        // add to list if filter does not already exist
        if (!foundIssuer) {
          this.filters['issuer_id'].options.push({
            name: asset.issuer_id,
            checked: false
          });
        }

        // Asset Code filter
        const foundAssetCode = this.filters['asset_code'].options.find((option: CheckboxOption) => {
          return option.name === asset.asset_code;
        });

        // add to list if filter does not already exist
        if (!foundAssetCode) {
          this.filters['asset_code'].options.push({
            name: asset.asset_code,
            checked: false
          });
        }

        // set currency field (used only in UI, for now)
        asset.currency = asset.asset_code.substring(0, 3);

        // Currency Filter
        const foundCurrency = this.filters['currency'].options.find((option: CheckboxOption) => {
          return option.name === asset.currency;
        });

        // add to list if filter does not already exist
        if (!foundCurrency) {
          this.filters['currency'].options.push({
            name: asset.currency,
            checked: false
          });
        }
      }

      this.currentAssetData = this.currentAccount.assets;
    }
  }

  getNotifications(): Subscription {
    return this.accountService.getTrustRequests('requestor_id', ['rejected'])
      .subscribe((requests: TrustRequest[]) => {
        this.pendingRequests = [];

        // filter by this current account
        requests = filter(requests, (request: TrustRequest) => {
          return request.account_name === this.accountService.accountSlug;
        });

        if (requests) {
          requests.forEach((request: TrustRequest) => {

            if (request.approval_ids && request.approval_ids.length === 1) {
              this.pendingRequests.push(request);
            }

            if (request.approval_ids && request.approval_ids.length > 1) {
              this.db.database.ref('participant_approvals')
                .child(request.approval_ids[1])
                .on('value', (data: firebase.database.DataSnapshot) => {

                  const approval: Approval = data.val();
                  if (approval && approval.status !== 'approved' && !this.pendingRequests.includes(request)) {
                    this.pendingRequests.push(request);
                  }
                });
            }
          });
        }
      });
  }

  /**
   * Creates table view from the gathered asset data
   *
   * @returns {void}
   * @memberof TrustedAssetsComponent
   */
  loadDataView(): void {

    // reset view data
    this.model.data = [[]];
    this.loaded = false;

    // Fallback in case of broken request. No assets or account details was found.
    if (!(this.currentAccount && this.currentAssetData)) {
      this.loaded = true;
      return;
    }

    for (const asset of this.currentAssetData) {
      const row = [
        new TableItem({
          data: asset.issuer_id
        }),
        new TableItem({
          data: asset.asset_code
        }),
        new TableItem({
          data: asset.currency
        }),
        // TODO: add asset limit after https://github.com/GFTN/gftn-services/issues/893 is completed
        // new TableItem({
        //   data: 'N/A'
        // }),
        new TableItem({
          data: asset.balance
        }),
      ];

      this.model.addRow(row);
    }

    // view has been fully created
    this.loaded = true;
  }

  /**
   * Filter viewable data by selected values
   *
   * @memberof TrustedAssetsComponent
   */
  filterData() {

    // reset asset data
    this.currentAssetData = this.currentAccount.assets;

    for (const [filterKey, filterVal] of Object.entries(this.filters)) {
      for (const option of filterVal.options) {
        if (option.checked) {
          this.currentAssetData = filter(this.currentAssetData,
            (asset: Asset) => {
              return asset[filterKey] === option.name;
            });
        }
      }
    }

    //  reload view
    this.loadDataView();
  }

  /**
   * Sort model data by index.
   * Pulled from carbon-components-angular table story
   * @param index
   */
  customSort(index: number, reload: boolean = false) {

    // only sort if there is data to sort
    // use '1' to account for empty array
    if (this.model.data.length > 1) {
      this.currSortIndex = index;
      this.sort(this.model, index, reload);
    }
  }

  private sort(model: TableModel, index: number, reload: boolean) {

    // take into account resort reload when loading more
    // so sort direction (asc/desc) doesn't get reversed
    if (model.header[index].sorted && !reload) {
      // if already sorted flip sorting direction
      model.header[index].ascending = model.header[index].descending;
    }
    model.sort(index);
  }

  /**
   * creates and opens the Modal for adding a new trustline request
   *
   * @memberof TrustedAssetsComponent
   */
  public addTrustLine() {

    // Utilizes Carbon Angular Modal
    this.currentOpenModal = this.modalService.create({
      component: TrustlineModalComponent,
      inputs: {
        MODAL_DATA: {}
      }
    });
  }
}
