// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, HostBinding, NgZone } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { toArray, find, filter, cloneDeep, isEmpty, startCase } from 'lodash';
import { Subject } from 'rxjs';
import { map, debounceTime, distinctUntilChanged, tap } from 'rxjs/operators';
import { ParticipantsModel } from '../../shared/models/participants.model';
import { IInstitution } from '../../shared/models/participant.interface';
import { ManageParticipantDialogComponent } from './manage-participant-dialog/manage-participant-dialog.component';
import { SessionService } from '../../shared/services/session.service';
import { Router } from '@angular/router';
import { AuthService } from '../../shared/services/auth.service';
/**
 * used by a participant admin user to manage participant user accounts
 *
 * @export
 * @class UsersComponent
 * @implements {OnInit}
 * @implements {AfterViewInit}
 */


export interface IInstitutionManageData {
  action: 'add' | 'edit';
  // institutionId: string;
  institution: IInstitution;
}
@Component({
  templateUrl: './accounts.component.html',
  styleUrls: ['./accounts.component.scss'],
  providers: [
    ParticipantsModel
  ]
})
export class AccountsComponent implements OnInit {

  public multiParticipants: boolean;
  public participants: IInstitution[];

  public searchedParticipants: IInstitution[];
  public keyUp = new Subject<string>();
  public searchText: string;

  constructor(
    public dialog: MatDialog,
    public _participants: ParticipantsModel,
    public sessionService: SessionService,
    private authService: AuthService,
    public router: Router,
    public ngZone: NgZone
  ) { }

  @HostBinding('attr.class') cls = 'flex-fill';

  ngOnInit() {
    this.getParticipants();
  }

  addParticipant() {
    this.openParticipantDialog('add');
  }

  private openParticipantDialog(action: 'add' | 'edit', institution?: IInstitution) {

    const data: IInstitutionManageData = {
      action: action,
      institution: institution,
    };

    const dialogRef = this.dialog.open(ManageParticipantDialogComponent, {
      disableClose: true,
      data: data
    });

    dialogRef.afterClosed().subscribe(result => {
      console.log('Participant dialog was closed');
      this.getParticipants();
    });
  }

  getParticipants() {

    // init empty participants array
    this.participants = [];
    this.searchedParticipants = [];

    this._participants.allPromise()
      .then((participants: IInstitution[]) => {

        this.ngZone.run(() => {
          this.participants = participants;

          this.searchedParticipants = participants;

          // initialize search
          this.initDebounceSearch();
        });

      });

  }

  getParticipantDetails(institutionId: string) {

    return find(this.participants, (p: IInstitution) => {
      return p.info.institutionId === institutionId;
    });

  }

  initDebounceSearch() {

    this.keyUp.pipe(
      // value to tap into
      map((event: any) => (event.target).value),
      // delay 1000
      debounceTime(1000),
      // don't fire unless new value
      distinctUntilChanged(),
      // update searched items
      tap(
        searchText => {
          this._search(searchText);
        }
      )
    ).subscribe(console.log);

  }

  private _search(searchText: string) {

    this.searchText = searchText;

    let searchedParticipants = [];

    // list only searched participants
    if (!isEmpty(searchText)) {

      console.log('search text: ', searchText);

      searchedParticipants = filter(this.participants, (institution: IInstitution) => {
        return institution.info.name.toLowerCase().includes(searchText.toLowerCase());
      });

    } else {
      // list all searched participants
      searchedParticipants = cloneDeep(this.participants);
    }

    this.searchedParticipants = searchedParticipants;

  }


}
