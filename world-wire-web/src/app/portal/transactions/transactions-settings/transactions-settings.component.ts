// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Input, AfterViewInit } from '@angular/core';
import { ContentSwitcherOption } from '../../../shared/models/content-switcher-option.model';
import { ContentSwitcher } from 'carbon-components';
import { ListItem } from 'carbon-components-angular';
import { TransactionService } from '../../shared/services/transaction.service';

@Component({
  selector: 'app-transactions-settings',
  templateUrl: './transactions-settings.component.html',
  styleUrls: ['./transactions-settings.component.scss'],
})
export class TransactionsSettingsComponent implements OnInit, AfterViewInit {

  // initial load
  loaded = false;

  public dropDownOptions: ListItem[] = [];

  constructor(
    public transactionService: TransactionService
  ) { }

  ngOnInit() {

    // initializing dropdown options for Records per Load settings
    for (const option of this.transactionService.perLoadOptions) {

      const selected = (this.transactionService.perLoad === option) ? true : false;
      this.dropDownOptions.push({
        content: option.toString(),
        selected: selected,
        value: option
      });
    }
  }

  ngAfterViewInit() {

    const contentSwitcher = document.getElementById('display-settings');

    // initializing content switcher for display settings
    if (contentSwitcher) {
      ContentSwitcher.create(contentSwitcher);
    }
  }

  selectDisplayOption(optionId: string) {

    let selectedOption: ContentSwitcherOption;

    // set all options to false
    for (const option of this.transactionService.tableDisplayOptions) {
      option.selected = false;
      if (option.id === optionId) {
        selectedOption = option;
      }
    }

    // save selected option to service
    if (selectedOption) {
      selectedOption.selected = true;
      this.transactionService.tableSize = selectedOption.id;
    }
  }

  /**
   * Called by dropdown to select option for perLoad amount
   * @param event
   */
  selectLoadOption(event) {

    // this event is called on initialization
    // so need to check for initial load vs. actual event
    if (!this.loaded) {
      this.loaded = true;
    } else {
      this.transactionService.perLoad = event.item.value;
    }
  }

  /**
   * Toggles table view for striped rows
   */
  updateStripedRows() {
    // toggle boolean
    this.transactionService.stripedRows = !this.transactionService.stripedRows;
  }
}
