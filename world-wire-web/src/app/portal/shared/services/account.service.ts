// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable, NgZone } from '@angular/core';
import { startCase, filter, sortBy } from 'lodash';
import { SessionService } from '../../../shared/services/session.service';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Participant } from '../../../shared/models/participant.interface';
import { ParticipantAccount, AccountRequest } from '../../../shared/models/account.interface';
import { Asset, AssetBalance, Obligation } from '../../../shared/models/asset.interface';
import { ENVIRONMENT } from '../../../shared/constants/general.constants';
import { AuthService } from '../../../shared/services/auth.service';
import { Observable, Observer, Subject } from 'rxjs';
import { TrustRequest, TrustRequestStatus } from '../models/trust-request.interface';
import { AngularFireDatabase } from '@angular/fire/database';

@Injectable()
export class AccountService {

  // slug of the current operation account being viewed
  accountSlug: string;

  issuedAssets: Asset[];

  // details of the current participant (node) being viewed
  participantDetails: Participant;

  // all participant accounts on the network
  allParticipants: Participant[];

  whitelistedParticipants: string[];

  // all accounts held by this participant
  allAccounts: ParticipantAccount[];

  // list of all assets on the network
  allAssets: { [issuer_id: string]: Asset[] };

  // Store root to the Client API for re-use in requests from same root
  apiRoot: string;

  globalRoot: string;

  public currentParticipantChanged: Subject<Participant>;

  constructor(
    private authService: AuthService,
    private sessionService: SessionService,
    private db: AngularFireDatabase,
    private ngZone: NgZone,
    private http: HttpClient
  ) {
    this.currentParticipantChanged = new Subject();
  }

  async getParticipant(): Promise<Participant> {

    try {

      // get account details of current participant node.
      // Client-facing endpoint (needs participant permissions) takes precedence
      let accountsRequest = `https://${this.sessionService.currentNode.participantId}.${ENVIRONMENT.envApiRoot}/v1/client/participants/${this.sessionService.currentNode.participantId}`;


      if (
        !this.authService.userIsParticipantManagerOrHigher(this.sessionService.currentNode.institutionId) &&
        this.authService.userIsSuperUser()
      ) {
        // use different URL to get participant details for users with only super admin permissions
        accountsRequest = `https://admin.${ENVIRONMENT.envGlobalRoot}/pr/v1/admin/pr/domain/${this.sessionService.currentNode.participantId}`;
      }

      // anchors/issuers also use a different endpoint entirely. grabbed from anchor service
      if (this.sessionService.currentNode.role === 'IS') {
        accountsRequest = `https://${this.sessionService.currentNode.participantId}.${ENVIRONMENT.envGlobalRoot}/anchor/v1/anchor/participants/${this.sessionService.currentNode.participantId}`;
      }

      const h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

      const options = {
        headers: h
      };

      let participant: Participant = null;

      // Get all accounts and details for this participant
      try {
        participant = await this.http.get(
          accountsRequest,
          options
        ).toPromise() as Participant;
      } catch (err) {
        participant = null;
      }

      // SAMPLE DATA. UNCOMMENT FOR TESTING
      // participant = {
      //   bic: 'USA1202323',
      //   country_code: 'USA',
      //   id: this.sessionService.currentNode.participantId,
      //   issuing_account: 'XXXXXXXXXXXX',
      //   role: 'MM',
      //   status: 'complete',
      // };

      // set participant in service
      this.participantDetails = participant;

      // reset API roots
      this.apiRoot = `https://${this.sessionService.currentNode.participantId}.${ENVIRONMENT.envApiRoot}`;

      this.globalRoot = `https://${this.sessionService.currentNode.participantId}.${ENVIRONMENT.envGlobalRoot}`;

      // reset accounts
      if (!this.allAccounts) {
        this.allAccounts = [];
      }

      if (participant.issuing_account) {
        this.allAccounts['issuing'] = {
          name: 'issuing',
          address: participant.issuing_account
        };
      }

      if (participant.operating_accounts) {
        for (const account of participant.operating_accounts) {
          this.allAccounts[account.name] = account;
        }
      }

      return this.participantDetails;
    } catch (err) {
      return null;
    }
  }

  /**
 * Used by other components to notify service
 * that current node has changed
 */
  public propogateParticipantChange() {

    this.currentParticipantChanged.next(this.participantDetails);
  }


  /**
   * Get account requests
   *
   * @returns {Observable<AccountRequest[]>}
   * @memberof AccountService
   */
  getAccountRequests(): Observable<AccountRequest[]> {
    return new Observable((observer: Observer<AccountRequest[]>) => {
      this.db.database.ref('account_requests')
        .child(this.sessionService.currentNode.participantId)
        .on('value', (data: firebase.database.DataSnapshot) => {

          if (data.val()) {

            const allRequests: AccountRequest[] = Object.values(data.val());

            // refresh the view
            this.ngZone.run(() => {
              // update observer value
              observer.next(
                allRequests,
              );
            }); // end ngZone
          }
        });
    });
  }

  /**
   * Get list of trust requests
   * Optional: filter by status
   *
   * @param {string} requestField
   * @param {TrustRequestStatus[]} [statusFilters]
   * @returns {Observable<TrustRequest[]>}
   * @memberof AccountService
   */
  getTrustRequests(requestField: string, statusFilters?: TrustRequestStatus[]): Observable<TrustRequest[]> {
    return new Observable((observer: Observer<TrustRequest[]>) => {
      this.db.database.ref('trust_requests')
        .orderByChild(requestField)
        .equalTo(this.sessionService.currentNode.participantId)
        .on('value', (data: firebase.database.DataSnapshot) => {

          const allRequests: { [key: string]: TrustRequest } = data.val() ? data.val() : {};

          let filteredRequests: TrustRequest[] = [];

          for (const [key, request] of Object.entries(allRequests)) {

            // temporarily store key for future querying
            request.key = key;

            // set loaded to false for the view
            request.loaded = false;

            filteredRequests.push(request);
          }

          // OPTIONAL: filter out all rejected requests by status
          if (statusFilters) {
            for (const statusFilter of statusFilters) {
              filteredRequests = filter(filteredRequests, (request: TrustRequest) => {
                return request.status !== statusFilter;
              });
            }
          }

          // refresh the view
          this.ngZone.run(() => {
            // update observer value
            observer.next(
              filteredRequests,
            );
          }); // end ngZone
        });
    });
  }

  /**
   * Gets all participants on the network
   *
   * @returns {Promise<Participant[]>}
   * @memberof AccountService
   */
  async getAllAssets(): Promise<any> {

    this.allAssets = {};

    // init list of available issuers/participants
    // available on this environment
    const accountsRequest = `${this.apiRoot}/v1/client/assets`;

    const h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

    const options = {
      headers: h
    };


    const assets: Asset[] = await this.http.get(
      accountsRequest,
      options
    ).toPromise() as Asset[];

    // initializing key map for assets and issuers
    // TODO: filter by DAs and DOs
    for (const asset of assets) {
      if (!this.allAssets[asset.issuer_id]) {
        this.allAssets[asset.issuer_id] = [];
      }
      this.allAssets[asset.issuer_id].push(asset);
    }

    return;
  }

  /**
   * Gets list of issued assets for this participant
   *
   * @returns {Promise<Asset[]>}
   * @memberof AccountService
   */
  async getIssuedAssets(): Promise<Asset[]> {

    // no official participant registered. return empty
    if (!this.participantDetails) {
      return null;
    }

    // request for regular participants - role === 'MM' (maker makers)
    let request = `${this.apiRoot}/v1/client/assets/issued`;

    if (this.participantDetails && this.participantDetails.role === 'IS') {
      // request for anchor participants - role === 'IS' (issuers)
      request = `${this.globalRoot}/anchor/v1/anchor/assets/issued/${this.participantDetails.id}`;
    }

    try {
      const h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

      const options = {
        headers: h
      };

      // SAMPLE DATA. Uncomment for testing in view
      // let assets: Asset[] = [{
      //   asset_code: 'SGDDO',
      //   asset_type: 'DO',
      //   issuer_id: this.sessionService.currentNode.participantId,
      //   balance: 1.0324234,
      // }, {
      //   asset_code: 'USDDO',
      //   asset_type: 'DO',
      //   issuer_id: this.sessionService.currentNode.participantId,
      //   balance: 1.34543,
      // }, {
      //   asset_code: 'HKDDO',
      //   asset_type: 'DO',
      //   issuer_id: this.sessionService.currentNode.participantId,
      //   balance: 2.45364,
      // }];

      // return type is list[] of Assets
      let assets: Asset[] = await this.http.get(
        request,
        options
      ).toPromise() as Asset[];

      // handle and filter out undefined/null assets
      assets = filter(assets, (asset: Asset) => {
        return asset.asset_code !== null;
      });

      return assets;
    } catch (err) {
      return null;
    }
  }

  /**
   * Get list of overall DO balances:
   * overall amount of an asset owed to other participants.
   * Optional: filter/query by asset code
   *
   * @param {string} [assetCode]
   * @returns {Promise<AssetBalance[]>}
   * @memberof AccountService
   */
  async getDOAssetBalances(assetCode?: string): Promise<AssetBalance[]> {

    try {

      let balanceRequest = `${this.apiRoot}/v1/client/obligations`;

      if (assetCode) {

        // include asset code in query if provided
        balanceRequest = `${balanceRequest}?asset_code=${assetCode}`;
      }

      const h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

      const options = {
        headers: h
      };

      // returns list of balances by issued asset (DO)
      const balances: AssetBalance[] = await this.http.get(
        balanceRequest,
        options
      ).toPromise() as AssetBalance[];


      return balances;
    } catch (err) {
      return null;
    }
  }

  /**
   * Gets a list of participants that hold a particular DO
   * issued by the current participant
   *
   * @param {string} assetCode
   * @param {string} [issuerId] (optional)
   * @returns {Promise<Participant[]>}
   * @memberof AccountService
   */
  async getAssetHolders(assetCode: string, issuerId?: string): Promise<Participant[]> {

    const requestedIssuerId = issuerId || this.sessionService.currentNode.participantId;

    try {
      const request = `${this.apiRoot}/v1/client/participants?asset_code=${assetCode}&issuer_id=${requestedIssuerId}`;

      const h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

      const options = {
        headers: h
      };

      return await this.http.get(
        request,
        options
      ).toPromise() as Participant[];

    } catch (err) {
      return null;
    }
  }

  /**
   * Gets breakdown of participants that
   * a participant has obligations to pay back,
   * as well as the associated "balance" (amount) owed
   *
   * @param {string} assetCode
   * @returns {Promise<AssetBalance[]>}
   * @memberof AccountService
   */
  async getDOBalanceDetails(assetCode: string): Promise<Obligation[]> {

    try {
      const request = `${this.apiRoot}/v1/client/obligations/${assetCode}`;

      const h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

      const options = {
        headers: h
      };

      return await this.http.get(
        request,
        options
      ).toPromise() as Obligation[];
    } catch (err) {
      return null;
    }
  }

  async getAssetsForParticipant(
    participantId: string,
    type: 'issued' | 'trusted' | 'both' = 'issued', filterDO: boolean = true)
    : Promise<Asset[]> {

    try {
      const request = `${this.apiRoot}/v1/client/assets/participants/${participantId}?participant_id=${participantId}&type=${type}`;

      const h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

      const options = {
        headers: h
      };

      let assets: Asset[] = await this.http.get(
        request,
        options
      ).toPromise() as Asset[];

      // handle and filter out undefined/null assets
      assets = filter(assets, (asset: Asset) => {
        return asset.asset_code !== null;
      });

      // filter by DO, if enabled (default setting)
      if (filterDO) {
        assets = filter(assets, (asset: Asset) => {
          return asset.asset_type !== 'DO';
        });
      }

      return assets;

    } catch (err) {
      return [];
    }
  }

  /**
   * Get all assets and balances for account
   *
   * @param {string} accountName
   * @returns {Promise<Asset[]>}
   * @memberof AccountService
   */
  async getTrustedAssetBalances(accountName: string): Promise<Asset[]> {

    let assets: Asset[];

    try {
      const request = `${this.apiRoot}/v1/client/assets/accounts/${accountName}`;

      const h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

      const options = {
        headers: h
      };

      assets = await this.http.get(
        request,
        options
      ).toPromise() as Asset[];

      // SAMPLE DATA. Uncomment for testing in view
      //   assets = [{
      //   asset_code: 'USD',
      //   asset_type: 'DA',
      //   issuer_id: 'participant4',
      //   balance: 1.0324234,
      // }, {
      //   asset_code: 'USDDO',
      //   asset_type: 'DO',
      //   issuer_id: 'participant1',
      //   balance: 1.34543,
      // }, {
      //   asset_code: 'HKD',
      //   asset_type: 'DA',
      //   issuer_id: 'participant4',
      //   balance: 2.45364,
      // }];
    } catch (err) {
      // fall back to empty array in case of request error
      return [];
    }

    // handle and filter out undefined/null assets
    assets = filter(assets, (asset: Asset) => {
      return asset.asset_code !== null;
    });

    for (const asset of assets) {

      // handle undefined/null assets
      if (asset.asset_code && asset.issuer_id && !asset.balance) {
        const balanceRequest = `${this.apiRoot}/v1/client/balances/accounts/${accountName}?asset_code=${asset.asset_code}&issuer_id=${asset.issuer_id}`;

        const h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

        const options = {
          headers: h
        };

        const balanceObject: AssetBalance = await this.http.get(
          balanceRequest,
          options
        ).toPromise() as AssetBalance;

        if (balanceObject) {
          asset.balance = parseFloat(balanceObject.balance);
        }
      }
    }

    // sort assets by balance amount
    assets = sortBy(
      assets, (asset: Asset) => {
        // descending order
        return -asset.balance;
      }
    );

    return assets;
  }

  /**
   * Gets all active participants on the network
   *
   * @returns {Promise<Participant[]>}
   * @memberof AccountService
   */
  async getAllParticipants(): Promise<Participant[]> {

    // get list of all issuers/participants from PR
    let accountsRequest = `https://${this.sessionService.currentNode.participantId}.${ENVIRONMENT.envApiRoot}/v1/client/participants`;

    // different endpoint for IS/Anchor
    if (this.sessionService.currentNode.role === 'IS') {
      accountsRequest = `https://${this.sessionService.currentNode.participantId}.${ENVIRONMENT.envGlobalRoot}/anchor/v1/anchor/participants`;
    }

    const h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

    const options = {
      headers: h
    };


    const participants: Participant[] = await this.http.get(
      accountsRequest,
      options
    ).toPromise() as Participant[];

    // return only active participants
    this.allParticipants = participants.filter((participant: Participant) => participant.status === 'active');

    // sort assets by balance amount
    this.allParticipants = sortBy(
      this.allParticipants, (participant: Participant) => {
        // Sort Order: normal from A - Z
        return participant.id;
      }
    );

    return this.allParticipants;
  }

  /**
   * Get list of whitelisted participants
   *
   * @returns {Promise<string[]>}
   * @memberof AccountService
   */
  async getWhitelistedParticipants(): Promise<string[]> {
    try {
      const request = `${this.globalRoot}/whitelist/v1/client/participants/whitelist`;

      const h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId, this.sessionService.currentNode.participantId);

      const options = {
        headers: h
      };

      this.whitelistedParticipants = await this.http.get(
        request,
        options
      ).toPromise() as string[];

      // this.whitelistedParticipants = await new Promise((resolve) => {
      //   setTimeout(() => {
      //     resolve(['ibmanchor', 'participant3', 'participant4']);
      //   }, 3000);
      // }) as string[];

      return this.whitelistedParticipants;

    } catch (err) {
      return null;
    }
  }

  /**
   * Generates human-readable name of the
   * current operating account being viewed
   *
   * @returns
   * @memberof AccountService
   */
  public getAccountNameReadable() {
    return startCase(this.accountSlug);
  }
}
