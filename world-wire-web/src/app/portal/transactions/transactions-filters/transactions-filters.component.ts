// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, ElementRef, AfterViewInit, ViewChild, OnDestroy } from '@angular/core';
import { NgForm } from '@angular/forms';
import { TransactionService } from '../../shared/services/transaction.service';
import { CheckboxOption } from '../../../shared/models/checkbox-option.model';
import { ListItem } from 'carbon-components-angular';
import { DatePicker } from 'carbon-components';
import * as _ from 'lodash';
import { Filter } from '../../shared/models/filter.model';
import { TransactionTypeDetail, TransactionKeyBook } from '../../../shared/models/transaction.interface';

/**
 * Data Model to store options for asset filters
 * Types of assets can be: 'settlement', 'sent', 'received'
 */
export class AssetFilterData {
  amountChecked: boolean;

  amountLogic: string;

  amount?: number;

  assetChecked: boolean;

  // stores selected asset options to filter by
  selectedAssets?: ListItem[];

  // stores list of available asset filter options
  assetOptions: ListItem[] = [];
}

/**
 * Data Model for form used to
 * create and/or select transaction filters
 */
export class FilterFormData {

  transactionStatuses: CheckboxOption[];

  // Payment DIRECTION of transaction (one-way only)
  transactionDirections: CheckboxOption[];

  assetFilters: { [key: string]: AssetFilterData };

  dateChecked: boolean;

  // stores date ranges for filtering
  dateFrom?: string;

  dateTo?: string;
}

/**
 * Transactions Filters Menu.
 * Creates and stores a form used by a participant user
 * to filter through transactions.
 *
 * @export
 * @class TransactionsFiltersComponent
 * @implements {OnInit}
 * @implements {AfterViewInit}
 * @implements {OnDestroy}
 */
@Component({
  selector: 'app-transactions-filters',
  templateUrl: './transactions-filters.component.html',
  styleUrls: ['./transactions-filters.component.scss']
})

export class TransactionsFiltersComponent implements OnInit, AfterViewInit, OnDestroy {

  logicDropOptions: ListItem[] = [];

  @ViewChild('dateFromRef') dateFrom: ElementRef;

  @ViewChild('dateToRef') dateTo: ElementRef;

  // stores data for new filter selections
  filterFormData: FilterFormData;

  // use for visual comparison to enable 'save' button
  defaultData: FilterFormData;

  transactionTypeDetails: TransactionTypeDetail;

  // stores filter keys for current transaction type
  transactionKeys: TransactionKeyBook;

  calendarEl: any;

  constructor(
    protected transactionService: TransactionService,
  ) { }

  ngOnInit() {

    // get current transaction type details
    this.transactionTypeDetails = _.find(this.transactionService.transactionTypes,
      (type: TransactionTypeDetail) => {
        return type.key === this.transactionService.transactionType;
      });

    // get filter keys for current transaction type
    this.transactionKeys = this.transactionService
      .transactionKeys[this.transactionService.transactionType];

    const transactionStatuses: CheckboxOption[] = [];

    // get transaction statuses for the selected transaction tiype
    for (const status of this.transactionTypeDetails.transactionStatuses) {
      transactionStatuses.push({
        name: status.name,
        label: status.label,
        checked: false
      });
    }

    const transactionDirections: CheckboxOption[] = [{
      name: '+',
      label: 'Incoming Transactions (+)',
      checked: false
    }, {
      name: '-',
      label: 'Outgoing Transactions (-)',
      checked: false
    }];

    for (const option of this.transactionService.logicOptions) {

      this.logicDropOptions.push({
        content: option,
        selected: false
      });
    }

    this.logicDropOptions[0].selected = true;

    // initializing empty form data
    this.filterFormData = {
      transactionStatuses: transactionStatuses,
      transactionDirections: transactionDirections,
      assetFilters: {},
      dateChecked: false,
    };

    // initializing empty or saved asset options
    this.initializeAssetOptions();

    if (this.transactionTypeDetails.filters) {
      this.loadDefaultFilters();
    }

    // set default data to track form changes for visual feedback
    this.defaultData = _.cloneDeep(this.filterFormData);
  }

  ngAfterViewInit() {

    // get date picker element from DOM to initialize Carbon date-picker component
    // This is neccessary for Vanilla Carbon, but we will change this later
    // so we don't have to rely on the DOM for component initialization
    const datepickerEl = document.getElementById('filter-date-picker');

    // check for element, so it doesn't error out when initializing the datepicker
    if (datepickerEl) {
      // date: mm/dd/yyyy is the format selected for the Carbon spec
      const dateOptions = { 'day': '2-digit', 'month': '2-digit', 'year': 'numeric' };

      // initializing date picker via Vanilla Carbon Components
      // TODO: Change to Angular components when it's available
      this.calendarEl = DatePicker.create(document.getElementById('filter-date-picker'), {
        maxDate: Date.now(),
        onChange: (selectedDates) => {

          // add date form options to reformat date model upon date selection
          // from the UI since the model doesn't bind properly
          this.filterFormData.dateFrom = selectedDates[0]
            .toLocaleDateString(window.navigator.language, dateOptions);
          this.filterFormData.dateTo = selectedDates[1] ? selectedDates[1]
            .toLocaleDateString('en-US', dateOptions) : '';

        },
        onOpen: () => {

          // sync date filter option selection upon click
          if (!this.filterFormData.dateChecked) {
            this.filterFormData.dateChecked = true;
          }
        },
        onClose: () => {

          // set default end date to today
          // if start date is selected and end date is missing
          if (this.filterFormData.dateFrom && !this.filterFormData.dateTo) {
            this.filterFormData.dateTo = new Date(Date.now())
              .toLocaleDateString(window.navigator.language, dateOptions);
          }

        }
      });
    }
  }

  ngOnDestroy() {

    // Cleanup: destroy vanilla js datepicker element from page
    const element = this.calendarEl ? this.calendarEl.calendar.calendarContainer : null;
    if (element) {
      element.parentNode.removeChild(element);
    }
  }

  initAssetFilterData(): AssetFilterData {
    return {
      amountChecked: false,
      amountLogic: '=',
      assetChecked: false,
      assetOptions: [],
      selectedAssets: []
    };
  }

  /**
   * Initialization of asset options for the view
   */
  initializeAssetOptions(): void {

    // Each object MUST BE INITIALIZED SEPARATELY because
    // assigning to object is always by reference.
    this.filterFormData.assetFilters['settlement'] = this.initAssetFilterData();
    this.filterFormData.assetFilters['sent'] = this.initAssetFilterData();
    this.filterFormData.assetFilters['received'] = this.initAssetFilterData();

    // initializing 'settlement' asset options
    for (const option of this.transactionService.settlementAssetOptions) {
      this.filterFormData
        .assetFilters['settlement']
        .assetOptions.push(
          {
            content: option,
            selected: false
          }
        );
    }

    // initializing 'sent' asset options
    for (const option of this.transactionService.sentAssetOptions) {
      this.filterFormData
        .assetFilters['sent']
        .assetOptions.push(
          {
            content: option,
            selected: false
          }
        );
    }

    // initializing 'received' asset options
    for (const option of this.transactionService.receivedAssetOptions) {
      this.filterFormData
        .assetFilters['received']
        .assetOptions.push(
          {
            content: option,
            selected: false
          }
        );
    }
  }

  /**
   * Loads cached filters from transaction service.
   * Speeds load time when switching between views in SPA
   */
  loadDefaultFilters() {

    _.forEach(this.transactionTypeDetails.filters,
      (filterArray: Filter[], key: string) => {

        // Load saved Transaction Status filters from TransactionService
        if (key === this.transactionKeys.transactionStatus) {
          filterArray.forEach((filter: Filter) => {
            const option: CheckboxOption = _.filter(this.filterFormData.transactionStatuses, (o: CheckboxOption) => {
              return o.name === filter.value;
            })[0];

            // set default
            if (option) {
              option.checked = true;
            }
          });
        }

        // Load saved Payment Type (transactionDirection) filters from TransactionService
        if (key === this.transactionKeys.paymentType) {
          filterArray.forEach((filter: Filter) => {
            const option: CheckboxOption = _.filter(this.filterFormData.transactionDirections, (o: CheckboxOption) => {
              return o.name === filter.value;
            })[0];

            // set default
            if (option) {
              option.checked = true;
            }
          });
        }

        // load saved 'Sent' Asset filters from TransactionService
        if (key === this.transactionKeys.sentAmount) {
          this.setDefaultAssetAmount(this.filterFormData.assetFilters['sent'], filterArray);
        }

        if (key === this.transactionKeys.sentAssetCode) {
          this.setDefaultAssetCodes(this.filterFormData.assetFilters['sent'], filterArray);
        }

        // load saved 'received' Asset filters from TransactionService
        if (key === this.transactionKeys.receivedAmount) {
          this.setDefaultAssetAmount(this.filterFormData.assetFilters['received'], filterArray);
        }

        if (key === this.transactionKeys.receivedAssetCode) {
          this.setDefaultAssetCodes(this.filterFormData.assetFilters['received'], filterArray);
        }

        // load saved 'Settlement' Asset filters from TransactionService
        if (key === this.transactionKeys.settlementAmount) {
          this.setDefaultAssetAmount(this.filterFormData.assetFilters['settlement'], filterArray);
        }

        if (key === this.transactionKeys.settlementAssetCode) {
          this.setDefaultAssetCodes(this.filterFormData.assetFilters['settlement'], filterArray);
        }

        // Load saved Date filters from TransactionService
        if (key === this.transactionKeys.timeStamp) {
          this.filterFormData.dateChecked = true;

          const filter = filterArray[0];
          filter.value = filter.value as string;
          const range = filter.value.split('|');

          // initialize date values
          this.filterFormData.dateFrom = range[0];
          this.filterFormData.dateTo = range[1];
        }
      });
  }

  /**
   * Setting default Asset Amount from saved filters
   * based on the filter key. Uses 0th index since
   * there is typically only one filter per asset option
   * @param assetKey
   * @param filterArray
   */
  setDefaultAssetAmount(assetKey, filterArray) {
    assetKey.amountChecked = true;
    assetKey.amountLogic = filterArray[0].logic;
    assetKey.amount = filterArray[0].value as number;
  }

  /**
   * Setting default Asset Code from saved filters
   * based on the filter key. Uses 0th index since
   * there is typically only one filter per asset option
   * @param assetKey
   * @param filterArray
   */
  setDefaultAssetCodes(assetKey, filterArray) {
    assetKey.assetChecked = true;

    // initialize array
    assetKey.selectedAssets = [];

    for (const filter of filterArray) {
      const option: ListItem[] = _.filter(assetKey.assetOptions, (o: ListItem) => {
        return o.content.toLowerCase() === filter.value.toLowerCase();
      });

      option[0].selected = true;
    }
  }


  /**
   * Closes the Transaction Filter Menu manually
   */
  closeMenu() {
    this.transactionService.popoverRef.doClose();
  }

  onSubmit(form: NgForm) {
    // console.log(form.value);
    // console.log(this.filterFormData);

    // close filter menu upon successful submission of form data
    this.closeMenu();

    // clear old filters upon save of new filters
    this.transactionTypeDetails.filters = {};

    // saving filters for transaction statuses
    this.filterFormData.transactionStatuses.forEach((option: CheckboxOption) => {
      if (option.checked === true) {

        // initializing filter for transaction status by key
        if (!this.transactionTypeDetails.filters[this.transactionKeys.transactionStatus]) {
          this.transactionTypeDetails.filters[this.transactionKeys.transactionStatus] = [];
        }

        // creating and appending filter data to transaction status filter
        const filter = new Filter();
        filter.value = option.name;

        this.transactionTypeDetails.filters[this.transactionKeys.transactionStatus].push(filter);
      }
    });

    // saving filters for transaction types (incoming/outgoing)
    this.filterFormData.transactionDirections.forEach((option: CheckboxOption) => {
      if (option.checked === true) {
        if (!this.transactionTypeDetails.filters[this.transactionKeys.paymentType]) {
          this.transactionTypeDetails.filters[this.transactionKeys.paymentType] = [];
        }

        const filter = new Filter();
        filter.value = option.name;
        this.transactionTypeDetails.filters[this.transactionKeys.paymentType].push(filter);
      }
    });

    // save asset options
    _.forEach(this.filterFormData.assetFilters, (filterData: AssetFilterData, key: string) => {
      let filterAmountKey: string;
      let filterAssetKey: string;

      // get filter key and save to filter
      switch (key) {
        case 'settlement': {
          filterAmountKey = this.transactionKeys.settlementAmount;
          filterAssetKey = this.transactionKeys.settlementAssetCode;
        }
          break;
        case 'sent': {
          filterAmountKey = this.transactionKeys.sentAmount;
          filterAssetKey = this.transactionKeys.sentAssetCode;
        }
          break;
        case 'received': {
          filterAmountKey = this.transactionKeys.receivedAmount;
          filterAssetKey = this.transactionKeys.receivedAssetCode;
        }
          break;
        default:
          break;
      }

      // save asset amount to filter
      if (filterData.amountChecked) {
        this.transactionTypeDetails.filters[filterAmountKey] = [{
          logic: filterData.amountLogic,
          value: filterData.amount
        }];
      }

      // Create new asset filter
      if (filterData.assetChecked && filterData.selectedAssets.length > 0) {
        filterData.selectedAssets.forEach((option: ListItem) => {
          if (option.selected === true) {
            if (!this.transactionTypeDetails.filters[filterAssetKey]) {
              this.transactionTypeDetails.filters[filterAssetKey] = [];
            }

            const assetCodeFilter = new Filter();
            assetCodeFilter.value = option.content;
            this.transactionTypeDetails.filters[filterAssetKey].push(assetCodeFilter);
          }
        });
      }
    });

    if (this.filterFormData.dateChecked) {

      // Secondary check to make sure dateFrom is BEFORE dateTo
      const tempDate = this.filterFormData.dateFrom;

      if (new Date(this.filterFormData.dateFrom) > new Date(this.filterFormData.dateTo)) {
        this.filterFormData.dateFrom = this.filterFormData.dateTo;
        this.filterFormData.dateTo = tempDate;
      }

      // create Date Filter
      this.transactionTypeDetails.filters[this.transactionKeys.timeStamp] = [{
        logic: 'date',
        value: `${this.filterFormData.dateFrom}|${this.filterFormData.dateTo}`
      }];
    }

    // reload transaction data to account for filters
    this.transactionService.toggleReloaded();
  }

  /**
   * Compare whether filter data has changed in form
   * */
  public isDataUpdated(): boolean {
    return _.isEqual(this.filterFormData, this.defaultData);
  }

  /**
   * Sync selected Assets to form data
   * @param selectedAssets
   */
  selectAssets(selectedAssets: ListItem[], assetType: string) {
    this.filterFormData.assetFilters[assetType].selectedAssets = selectedAssets;
  }
}
