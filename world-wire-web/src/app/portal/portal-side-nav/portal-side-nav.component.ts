// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, AfterViewInit } from '@angular/core';
import { Accordion } from 'carbon-components';
import { SessionService } from '../../shared/services/session.service';
import { IInstitution } from '../../shared/models/participant.interface';
import { AuthService } from '../../shared/services/auth.service';
import { ListItem } from 'carbon-components-angular';

@Component({
  selector: 'app-portal-side-nav',
  templateUrl: './portal-side-nav.component.html',
  styleUrls: ['./portal-side-nav.component.scss']
})
export class PortalSideNavComponent implements OnInit, AfterViewInit {

  public open = false;
  private accordionInitialized = false;

  public institution: IInstitution;

  showList: boolean;

  environmentOptions: ListItem[] = [];

  constructor(
    public sessionService: SessionService,
    private authService: AuthService,
  ) { }

  ngOnInit() {

    // get institution/participant from session/slug
    if (this.sessionService.institution) {
      this.institution = this.sessionService.institution;
      // show list if permissions exist
      this.showList = this.institution ? this.show(this.institution.info.institutionId, ['admin', 'manager'], ['admin', 'manager']) : false;
    }
  }

  ngAfterViewInit() {
    // initialize accordion only once
    if (this.accordionInitialized === false) {
      const elements = document.getElementsByClassName('bx--accordion');
      for (let i = 0; i < elements.length; i++) {
        Accordion.create(elements[i]);
      }

      this.accordionInitialized = true;
    }
  }

  /**
   * Toggles showing link if the user does or does not have elevated permissions
   *
   * @param {string[]} superPermission
   * @param {string[]} participantPermission
   * @memberof PortalSideNavComponent
   */
  public show(institutionId: string, superPermission: string[], participantPermission: string[]) {

    // check if participant permissions exist
    if (this.authService.hasParticipantPermissions(this.authService.userProfile, institutionId, participantPermission)) {
      return true;
    }

    // check if super permissions exist
    if (this.authService.hasSpecificSuperPermissions(this.authService.userProfile, superPermission)) {
      return true;
    }

    return false;

  }
}
