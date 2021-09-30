// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, OnDestroy, ComponentRef } from '@angular/core';
import { AccountService } from '../../shared/services/account.service';
import { AssetRequest, Asset, AssetType } from '../../../shared/models/asset.interface';
import { AngularFireDatabase } from '@angular/fire/database';
import { SessionService } from '../../../shared/services/session.service';
import { Subscription } from 'rxjs';
import { AssetModalComponent } from '../../shared/components/asset-modal/asset-modal.component';
import { ModalService } from 'carbon-components-angular';
import { AuthService } from '../../../shared/services/auth.service';

@Component({
  selector: 'app-assets-overview',
  templateUrl: './assets-overview.component.html',
  styleUrls: ['./assets-overview.component.scss']
})
export class AssetsOverviewComponent implements OnInit, OnDestroy {

  participantSubscription: Subscription;

  // stores all asset requests
  allAssetRequests: AssetRequest[];

  // stores only pending asset requests
  pendingAssetRequests: AssetRequest[];

  assetRequestsLoaded = false;

  issuedAssetsLoaded = false;

  // references currently opened modal for closing/dereferencing later
  currentOpenModal: ComponentRef<AssetModalComponent>;

  issuedAssetType: AssetType = 'DO';

  error = false;

  isAnchor = false;

  participantAuthorized = false;

  superAuthorized = false;

  constructor(
    private db: AngularFireDatabase,
    private authService: AuthService,
    public accountService: AccountService,
    private sessionService: SessionService,
    private modalService: ModalService
  ) { }

  ngOnInit() {

    this.participantAuthorized = this.sessionService.institution && this.authService.userIsParticipantManagerOrHigher(
      this.sessionService.institution.info.institutionId
    );

    const nodeIsAnchor = this.sessionService.currentNode ? this.sessionService.currentNode.role === 'IS' : this.isAnchor;

    // check if isAnchor from current node configuration
    this.isAnchor = nodeIsAnchor;

    this.participantSubscription = this.accountService.currentParticipantChanged.subscribe(() => {

      if (!this.accountService.participantDetails) {
        this.error = true;
        return;
      }

      this.issuedAssetType = this.accountService.participantDetails && this.accountService.participantDetails.role === 'IS' ? 'DA' : 'DO';

      // check if isAnchor from official entry in PR if it exists
      this.isAnchor = this.accountService.participantDetails ? this.accountService.participantDetails.role === 'IS' : nodeIsAnchor;

      this.loadAssets();
    });
  }

  ngOnDestroy() {

    // programmatically close modal if open
    if (this.currentOpenModal) {
      this.currentOpenModal.instance.closeModal();
    }

    this.participantSubscription.unsubscribe();
  }

  private async loadAssets(refresh?: boolean) {
    await this.getIssuedAssets(refresh);

    // wait for issued assets request to come back before getting pending requests
    this.getAssetRequests();
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
      this.pendingAssetRequests = null;
    }

    // alway reset loader
    this.issuedAssetsLoaded = false;

    // store in service to prevent requests on every page navigation
    try {
      this.accountService.issuedAssets = this.accountService.issuedAssets ? this.accountService.issuedAssets : await this.accountService.getIssuedAssets();
    } catch (err) {
      this.error = true;
    }

    // promise returned null = request failed
    if (!this.accountService.issuedAssets) {
      this.error = true;
    }

    // data loaded
    this.issuedAssetsLoaded = true;
  }


  /**
   * Gets list of asset issuance requests pending approval/rejection
   *
   * @private
   * @returns {Promise<void>}
   * @memberof AccountsOverviewComponent
   */
  private async getAssetRequests(): Promise<void> {

    this.assetRequestsLoaded = false;

    try {

      this.db.database.ref('asset_requests')
        .child(this.sessionService.currentNode.participantId)
        .once('value', (data: firebase.database.DataSnapshot) => {
          if (data.val()) {
            this.allAssetRequests = Object.values(data.val());

            // only get unapproved, pending requests
            const filteredAssetRequests = this.allAssetRequests.filter((asset: AssetRequest) => {
              if (this.accountService.issuedAssets) {
                const foundAsset: Asset = this.accountService.issuedAssets.find((issuedAsset: Asset) => {
                  return issuedAsset.asset_code === asset.asset_code;
                });

                // asset already issued. remove from list of pending requests
                if (foundAsset) {
                  return false;
                }
              }

              return asset.approvalIds;
            });

            this.pendingAssetRequests = filteredAssetRequests;
          }

          this.assetRequestsLoaded = true;
        });

    } catch (err) {

      this.error = true;

      this.assetRequestsLoaded = true;
    }
  }

  /**
   * Gets request from existing list of requests
   *
   * @param {Asset} asset
   * @returns {(AssetRequest | Asset)}
   * @memberof AssetsOverviewComponent
   */
  getAssetRequest(asset: Asset): AssetRequest {
    if (!this.allAssetRequests) {
      return null;
    }
    return this.allAssetRequests.find((req: AssetRequest) => {
      return req.asset_code === asset.asset_code;
    });
  }

  /**
   * Opens new modal for issuing new asset
   *
   * @param {AssetType} type
   * @param {AssetRequest} [assetRequest]
   * @memberof AssetsOverviewComponent
   */
  public openAssetModal(type: AssetType, assetRequest?: AssetRequest) {

    let request = assetRequest ? this.getAssetRequest(assetRequest) : null;

    if (!request) {
      request = assetRequest;
    }

    // creates new modal
    this.currentOpenModal = this.modalService.create({
      component: AssetModalComponent,
      inputs: {
        MODAL_DATA: {
          [!assetRequest && 'assetType']: type,
          [assetRequest && 'assetRequest']: request
        }
      }
    });

    // listen to close event of modal
    this.currentOpenModal.instance.close.subscribe(() => {

      this.loadAssets(true);
    });
  }
}
