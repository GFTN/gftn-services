// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { OnInit, Component, Inject, ViewChild, ElementRef } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import * as countries from 'i18n-iso-countries';
import { cloneDeep, map, isEqual } from 'lodash';
import { ParticipantsModel } from '../../../shared/models/participants.model';
import { IInstitutionManageData } from '../accounts.component';
import { ModalService, ModalButtonType, AlertModalType } from 'carbon-components-angular';
import { AuthService } from '../../../shared/services/auth.service';
import { KeyValue } from '@angular/common';

/**
 * Used by IBM Super Admin to manage a participant account
 *
 * @export
 * @class ManageParticipantDialogComponent
 * @implements {OnInit}
 */
@Component({
  templateUrl: './manage-participant-dialog.component.html',
  styleUrls: ['./manage-participant-dialog.component.scss'],
  providers: [
    ParticipantsModel
  ]
})
export class ManageParticipantDialogComponent implements OnInit {

  public roles;
  public countries: countries.LocalizedCountryNames;

  @ViewChild('nameElement') nameElement: ElementRef;


  // sorts object by value
  valueAscOrder = (a: KeyValue<string, string>, b: KeyValue<string, string>): number => {
    return a.value.localeCompare(b.value);
  }

  constructor(
    private authService: AuthService,
    private modalService: ModalService,
    public dialogRef: MatDialogRef<ManageParticipantDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: IInstitutionManageData,
    public institution: ParticipantsModel
  ) {

    countries.registerLocale(require(`i18n-iso-countries/langs/${this.authService.userLangShort}.json`));

    if (!data.institution) {

      // create new participant (with empty fields) to add
      this.data.institution = cloneDeep(this.institution.model);

    } else {

      // set to existing participant to edit
      this.institution.model = cloneDeep(this.data.institution);
      // this.institution.model.institutionId = data.institutionId;

    }

  }

  ngOnInit(): void {
    this.countries = countries.getNames(this.authService.userLangShort);
  }

  /**
   * Converts 2 letter country code to 3 letter (compatibility version)
   *
   * @param {string} alpha2Code
   * @returns
   * @memberof ManageParticipantDialogComponent
   */
  getAlpha3Code(alpha2Code: string) {
    return countries.alpha2ToAlpha3(alpha2Code);
  }

  /**
   * Opens a modal to confirm slug change
   *
   * @memberof ManageParticipantDialogComponent
   */
  confirmSlugChange() {
    this.modalService.show({
      type: AlertModalType.danger,
      label: 'Warning',
      title: 'Confirm Identifier Change',
      content: 'This will change the page route for this participant. Please CONFIRM that you want to change the identifier for this participant.',
      buttons: [{
        text: 'Cancel',
        type: ModalButtonType.tertiary,
      }, {
        text: 'Confirm',
        type: ModalButtonType.danger_primary,
        click: () => {
          // updates slug
          if (this.institution.model) {
            this.institution.genSlug(this.institution.model.info.name);

            // trigger event to run slug validation
            const element: HTMLElement = this.nameElement.nativeElement;
            element.dispatchEvent(new Event('input'));
          }
        }
      }]
    });
  }

  updateSlug(event: KeyboardEvent, edit: boolean = false) {

    // only update/generate slug if this is a new institution,
    // or action to manually update it was enabled
    if (!this.data.institution.info.slug || edit) {
      this.institution.genSlug((<HTMLInputElement>event.target).value);
    }
  }

  /**
   * Validates if form data has changed
   *
   * @returns
   * @memberof ManageParticipantDialogComponent
   */
  isFormUpdated() {
    return !isEqual(this.institution.model, this.data.institution);
  }

  submit() {
    // edit() or add()
    this[this.data.action]();
  }

  /**
   * called by submit if action is edit
   *
   * @memberof ManageParticipantDialogComponent
   */
  edit() {
    // save participant info to firebase
    if (this.isFormUpdated()) {
      this.institution.update();
    }

    this.close();
  }

  /**
   * called by submit() if action is add
   *
   * @memberof ManageParticipantDialogComponent
   */
  add() {

    // save participant info to firebase
    this.institution.create();
    this.close();
  }

  /**
   * cancel button event
   *
   * @memberof ManageParticipantDialogComponent
   */
  cancel() {
    this.close();
  }

  /**
   * closes the modal
   *
   * @private
   * @memberof ManageParticipantDialogComponent
   */
  private close(): void {
    this.dialogRef.close();
  }

}
