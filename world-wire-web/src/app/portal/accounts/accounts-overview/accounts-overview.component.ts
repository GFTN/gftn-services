// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, isDevMode, OnDestroy, ComponentRef } from '@angular/core';
import { SessionService } from '../../../shared/services/session.service';
import { AccountService } from '../../shared/services/account.service';
import { ParticipantAccountDetail, AccountRequest } from '../../../shared/models/account.interface';
import { Asset } from '../../../shared/models/asset.interface';
import { Subscription } from 'rxjs';
import { AuthService } from '../../../shared/services/auth.service';
import { ModalService } from 'carbon-components-angular';
import { AccountModalComponent } from '../../shared/components/account-modal/account-modal.component';

/**
 * Main handler for accounts overview
 *
 * @export
 * @class AccountsOverviewComponent
 * @implements {OnInit}
 */
@Component({
  selector: 'app-accounts-overview',
  templateUrl: './accounts-overview.component.html',
  styleUrls: ['./accounts-overview.component.scss']
})
export class AccountsOverviewComponent implements OnInit, OnDestroy {

  participantSubscription: Subscription;

  // stores overview account details about operating accounts
  operatingAccounts: ParticipantAccountDetail[];

  // stores overview account details about issuing account
  issuingAccount: ParticipantAccountDetail;

  accountsLoaded = false;

  issuedAssetsLoaded = false;

  issuedAssetsError = false;

  accountRequestsLoaded = false;

  // toggle if current node is an issuer/anchor node
  isAnchor = false;

  // references currently opened modal for closing/dereferencing later
  currentOpenModal: ComponentRef<AccountModalComponent>;

  accountRequestSubscription: Subscription;

  pendingAccountRequests: AccountRequest[];

  participantAuthorized = false;

  superAuthorized = false;

  constructor(
    public sessionService: SessionService,
    public accountService: AccountService,
    private authService: AuthService,
    private modalService: ModalService
  ) { }

  ngOnInit() {

    this.participantAuthorized = this.sessionService.institution && this.authService.userIsParticipantManagerOrHigher(
      this.sessionService.institution.info.institutionId
    );

    this.superAuthorized = this.authService.userIsSuperUser();


    const nodeIsAnchor = this.sessionService.currentNode ? this.sessionService.currentNode.role === 'IS' : this.isAnchor;

    // check if isAnchor from current node configuration
    this.isAnchor = nodeIsAnchor;

    this.participantSubscription = this.accountService.currentParticipantChanged.subscribe(() => {

      // reset accounts data stored in view
      this.operatingAccounts = null;
      this.issuingAccount = null;

      // check if isAnchor from official entry in PR if it exists
      this.isAnchor = this.accountService.participantDetails ? this.accountService.participantDetails.role === 'IS' : nodeIsAnchor;

      // PR entry was successful
      if (this.accountService.participantDetails) {

        // PROMISE: get issuing account
        this.getIssuingAccountDetails(true);

        if (!this.isAnchor) {
          // get asset details for operation accounts
          this.getOperatingAccountDetails();
        }
      }

      // PROMISE: get issued assets
      this.getIssuedAssets(true);

      this.getAccountRequests();

      this.accountsLoaded = true;
    });
  }

  ngOnDestroy() {

    // cleanup observables for view
    this.participantSubscription.unsubscribe();

    if (this.accountRequestSubscription) {
      this.accountRequestSubscription.unsubscribe();
    }

    // programmatically close modal if open
    if (this.currentOpenModal) {
      this.currentOpenModal.instance.closeModal();
    }
  }

  private getAccountRequests(): void {
    // don't get requests if PR entry was unsuccessful
    if (!this.accountService.participantDetails) {
      return;
    }

    if (this.superAuthorized) {
      // don't start subscribing to account requests until AFTER participant details are retrieved
      this.accountRequestSubscription = this.accountService.getAccountRequests()
        .subscribe((allRequests: AccountRequest[]) => {
          this.pendingAccountRequests = allRequests.filter((request: AccountRequest) => {

            // filter out accounts that are already registered
            let foundAccount: ParticipantAccountDetail = (this.issuingAccount && this.issuingAccount.name === request.name) ? this.issuingAccount : null;

            foundAccount = (this.operatingAccounts && !foundAccount) ? this.operatingAccounts
              .find((account: ParticipantAccountDetail) => {
                return account.name === request.name;
              }) : null;

            if (foundAccount) {
              // set request in details
              foundAccount.request = request;
              return false;
            }

            // get requests with actual approvals pending
            return request.approvalIds;
          });
        });
    }
  }

  /**
   * Get all trusted assets for issuing account
   *
   * @returns {Promise<void>}
   * @memberof AccountsOverviewComponent
   */
  private async getIssuingAccountDetails(refresh?: boolean): Promise<void> {

    if (refresh) {
      this.issuingAccount = null;
    }

    try {
      if (this.accountService.allAccounts && 'issuing' in this.accountService.allAccounts) {
        this.issuingAccount = this.accountService.allAccounts['issuing'];

        if (!this.isAnchor) {
          this.issuingAccount.assets = this.issuingAccount.assets ? this.issuingAccount.assets : await this.accountService.getTrustedAssetBalances('issuing');
        }

        this.issuingAccount.loaded = true;
      }
    } catch (err) {
      this.issuingAccount.loaded = true;
    }
  }

  /**
   * Get all issued DOs
   *
   * @returns {Promise<void>}
   * @memberof AccountsOverviewComponent
   */
  private async getIssuedAssets(refresh?: boolean): Promise<void> {

    if (refresh) {
      this.accountService.issuedAssets = null;
    }

    this.issuedAssetsLoaded = false;

    try {

      // store in service to prevent requests on every page navigation
      this.accountService.issuedAssets = this.accountService.issuedAssets ? this.accountService.issuedAssets : await this.accountService.getIssuedAssets();

    } catch (err) {

      this.issuedAssetsError = true;
    }

    // promise returned null = request failed
    if (!this.accountService.issuedAssets) {
      this.issuedAssetsError = true;
    }

    this.issuedAssetsLoaded = true;
  }

  /**
   * Get all asset details for each operating account
   *
   * @private
   * @memberof AccountsOverviewComponent
   */
  private getOperatingAccountDetails() {
    this.operatingAccounts = (this.accountService.participantDetails && this.accountService.participantDetails.operating_accounts) ? this.accountService.participantDetails.operating_accounts : [];

    for (const account of this.operatingAccounts) {

      // init to false
      account.loaded = false;

      const savedAccount = this.accountService.allAccounts[account.name];

      // grab asset data from API, if not already retrieved
      if (savedAccount.assets) {
        account.assets = savedAccount.assets;

        account.loaded = true;
      } else {
        // get asset balances in parallel
        this.accountService.getTrustedAssetBalances(account.name).then((assets: Asset[]) => {
          savedAccount.assets = account.loaded = true;
          account.assets = assets;
        });
      }
    }
  }

  /**
   * Opens new modal for issuing new/viewing account request
   *
   * @param {AccountRequest} [request]
   * @memberof AccountsOverviewComponent
   */
  openAccountModal(request?: AccountRequest) {

    // creates new modal
    this.currentOpenModal = this.modalService.create({
      component: AccountModalComponent,
      inputs: {
        MODAL_DATA: request
      }
    });

    // listen to close event of modal
    this.currentOpenModal.instance.close.subscribe(async () => {

      // reset the accounts to grab them anew
      this.accountsLoaded = false;

      this.issuingAccount = null;
      this.operatingAccounts = null;

      this.accountService.getParticipant().then(() => {
        this.accountService.propogateParticipantChange();
      });
    });
  }
}
