// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, NgZone, HostBinding, OnDestroy } from '@angular/core';
import { SessionService } from '../../shared/services/session.service';
import { ITokenDialogData, TokenDialogComponent, ITokenActions } from './token-dialog/token-dialog.component';
import { MatDialog } from '@angular/material';
import { AngularFireDatabase } from '@angular/fire/database';
import { IJWTPublic } from '../../shared/models/token.interface';
import { UtilsService } from '../../shared/utils/utils';
import { AuthService } from '../../shared/services/auth.service';
import { map, filter, clone, isNumber } from 'lodash';
import { NotificationService } from '../../shared/services/notification.service';
import { Subscription, Observable, Observer } from 'rxjs';

interface IJWTInfoListItem extends IJWTPublic {
  key: string;
}

@Component({
  selector: 'app-token',
  templateUrl: './token.component.html',
  styleUrls: ['./token.component.scss']
})
export class TokenComponent implements OnInit, OnDestroy {

  showNotification: boolean;

  // unsorted and filtered
  jwt_infosAll: IJWTInfoListItem[];

  // sorted and filtered
  jwt_infos: IJWTInfoListItem[];

  viewAll: boolean;

  // toggles whether token results have finished loading to show for the view
  loaded = false;

  currentParticipantId: string;

  participantSubscription: Subscription;

  tokenSubscription: Subscription;

  constructor(
    public dialog: MatDialog,
    public sessionService: SessionService,
    private notificationService: NotificationService,
    private db: AngularFireDatabase,
    private ngZone: NgZone,
    public utils: UtilsService,
    public authService: AuthService
  ) {
    this.viewAll = false;
  }

  @HostBinding('attr.class') cls = 'flex-fill';

  ngOnInit() {
    this.tokenSubscription = this.getTokens().subscribe();

    this.currentParticipantId = this.sessionService.currentNode ? this.sessionService.currentNode.participantId : null;

    this.participantSubscription = this.sessionService.currentNodeChanged.subscribe(() => {

      // only propgate change if successful request was made
      if (this.currentParticipantId &&
        (this.currentParticipantId !== this.sessionService.currentNode.participantId)) {

        this.currentParticipantId = this.sessionService.currentNode.participantId;

        this.tokenSubscription = this.getTokens().subscribe();
      }
    });
  }

  ngOnDestroy() {
    this.participantSubscription.unsubscribe();

    this.tokenSubscription.unsubscribe();
  }

  /**
   * display list of tokens in the ui
   *
   * @memberof TokenComponent
   */
  getTokens(): Observable<void> {

    // Get by institutionId and filter query by currently viewed env
    return new Observable((observer: Observer<void>) => {
      this.db.database.ref('/jwt_info/' + this.sessionService.institution.info.institutionId)
        .on('value', (data: firebase.database.DataSnapshot) => {
          const jwt_infos: { [iid: string]: IJWTPublic } = data.val();

          if (jwt_infos) {

            // add keys inline with each obejct
            // because we must turn it into an array to sort
            let transform: IJWTInfoListItem[] = map(jwt_infos, (val: IJWTPublic, key: string) => {
              const v = val as IJWTInfoListItem;
              v.key = key;
              return v;
            });

            // reverse order so that the most recent tokens appear first
            transform = transform.reverse();

            this.jwt_infosAll = transform;
          }

          this.ngZone.run(() => {
            this.updateView();

            // update observer value
            observer.next();
          });
        });
    });
  }

  /**
   * create a request to generate a new JWT token
   *
   * @memberof TokenComponent
   */
  request() {
    // only open new request dialog if there is a participant
    // node AVAILABLE to create tokens for.
    // Show error dialog if there isn't one for this environment
    if (this.sessionService.currentNode) {
      this.openDialog('request');
    } else {
      this.notificationService.show(
        'error',
        'Cannot generate tokens. No participant node found for this environment.',
        null,
        'Error Generating Tokens',
        'top'
      );
    }
  }

  /**
   * 2nd admin approves the creation of a token
   * (think of it like a submarine captain
   * needing another officer to turn the key
   * to arm the nukes)
   *
   * @memberof TokenComponent
   */
  approve(jwt_info: IJWTPublic) {
    this.openDialog('approve', jwt_info);
  }

  /**
   * inactivating a token by revoking it
   *
   * @memberof TokenComponent
   */
  reject(jwt_info: IJWTPublic) {
    this.openDialog('reject', jwt_info);
  }

  /**
   * inactivating a token by revoking it
   *
   * @memberof TokenComponent
   */
  revoke(jwt_info: IJWTPublic) {
    this.openDialog('revoke', jwt_info);
  }

  /**
  * generates a token for one time viewing
  *
  * @memberof TokenComponent
  */
  generate(jwt_info: IJWTPublic) {
    this.openDialog('generate', jwt_info);
  }

  /**
   * Update list of tokens in the view
   */
  updateView(): void {

    // reset view load
    this.loaded = false;

    if (this.sessionService.institutionNodes && this.sessionService.currentNode) {
      // includes revoked
      this.jwt_infos = clone(this.jwt_infosAll);
    }

    // filter for currently-viewed participant
    if (this.jwt_infos && this.sessionService.currentNode) {
      this.jwt_infos = filter(this.jwt_infos, (val: IJWTInfoListItem) => {
        return val.aud === this.sessionService.currentNode.participantId;
      });
    }

    // only filter through results if there are any available for this environment
    if (this.jwt_infos && !this.viewAll) {

      // filters out revoked
      this.jwt_infos = filter(this.jwt_infos, (val: IJWTInfoListItem) => {
        return !isNumber(val.revokedAt);
      });
    }

    // view has completely loaded
    this.loaded = true;
  }

  openDialog(action: ITokenActions, jwt_info?: IJWTPublic) {

    const data: ITokenDialogData = {
      action: action,
      institution: this.sessionService.institution,
      jwt_info: jwt_info
    };

    const dialogRef = this.dialog.open(TokenDialogComponent, {
      disableClose: true,
      data: data
    });

    dialogRef.afterClosed().subscribe(result => {
    });

  }

  /**
   * Controls whether or not logged in user can view this token's details
   *
   * @param {string} isCreator
   * @returns
   * @memberof TokenComponent
   */
  userCanView(isCreator: string): boolean {

    // Allowed users that are able to view token details:
    // 1. creator of token
    // 2. participant admin (approver/rejector)
    return this.authService.userIsParticipantAdmin(this.sessionService.institution.info.institutionId) || this.authService.userProfile.profile.email === isCreator;
  }

  /**
   * Returns whether or not a user is a valid approver
   * Helper function for the view
   *
   * @param {string} isCreator
   * @returns {boolean}
   * @memberof TokenComponent
   */
  isValidApprover(isCreator: string): boolean {

    // Valid approver of a token must be:
    // 1. a participant admin
    // 2. NOT the creator of the token
    return this.authService.userIsParticipantAdmin(this.sessionService.institution.info.institutionId) && this.authService.userProfile.profile.email !== isCreator;
  }
}
