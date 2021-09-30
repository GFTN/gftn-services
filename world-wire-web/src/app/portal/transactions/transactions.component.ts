// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, HostBinding, ViewChild, OnDestroy, TemplateRef, ElementRef } from '@angular/core';
// import { ElasticsearchService } from '../shared/services/elasticsearch.service';
import { isEmpty, find, startCase } from 'lodash';
import { ITransfer, IExchange, TransactionTypeDetail, TransactionType, TransactionStatus } from '../../shared/models/transaction.interface';
import { HttpClient } from '@angular/common/http';

import { TableModel, TableItem, ModalService } from 'carbon-components-angular';
import { SessionService } from '../../shared/services/session.service';

// Necessary imports for data export
import { ExportModalComponent } from '../export-modal/export-modal.component';
import { TransactionService } from '../shared/services/transaction.service';
import { CustomTableHeaderItem, CustomTableItem } from '../shared/models/table-item.model';
import { Subject, Subscription } from 'rxjs';
import { map, debounceTime, distinctUntilChanged, tap } from 'rxjs/operators';
// import { INodeAutomation } from '../../shared/models/node.interface';
import { StatusHistoryModalComponent } from '../status-history-modal/status-history-modal.component';
import { UtilsService } from '../../shared/utils/utils';

@Component({
  selector: 'app-transactions',
  templateUrl: './transactions.component.html',
  styleUrls: ['./transactions.component.scss'],
})

export class TransactionsComponent implements OnInit, OnDestroy {

  // number of records loaded for this current view
  private currLoadLength = 0;

  public loaded = false;
  Math: any;
  Array: any;

  // Settings for IBM Data Table
  model = new TableModel();
  size = 'sm';
  showSelectionColumn = false;
  striped = false;

  currSortIndex: number;

  public keyUp = new Subject<string>();

  // store current and previous query for comparison
  public currSearchText;
  private prevSearchText;

  public transactionTypeDetails: TransactionTypeDetail;

  participantSubscription: Subscription;

  @ViewChild('customHeaderTemplate') customHeaderTemplate: TemplateRef<any>;
  @ViewChild('customTableItemTemplate') customTableItemTemplate: TemplateRef<any>;
  @ViewChild('tableItemTemplate') tableItemTemplate: TemplateRef<any>;
  @ViewChild('statusTemplate') statusTemplate: TemplateRef<any>;
  @ViewChild('overflowMenuItemTemplate') overflowMenuItemTemplate: TemplateRef<any>;
  @ViewChild('container') containerRef: ElementRef;

  constructor(
    public http: HttpClient,
    private sessionService: SessionService,
    protected modalService: ModalService,
    public transactionService: TransactionService,
    private utils: UtilsService
  ) {

    this.Math = Math;
    this.Array = Array;
  }

  @HostBinding('attr.class') cls = 'flex-fill';

  ngOnInit() {

    // initialize perLoad to first option onInit
    if (!this.transactionService.perLoad) {
      this.transactionService.perLoad = this.transactionService.perLoadOptions[0];
    }

    if (!this.transactionService.statusColorMap) {
      this.transactionService.initializeTransactionStatuses();
    }

    // initialize table display options
    if (!this.transactionService.tableDisplayOptions) {
      this.transactionService.initializeTableDisplayOptions();
    }

    // initializing model header for new empty view
    this.model.header = [];
    /*
    // TBD once ES Data is done
    this.elasticsearchService.search()
      .subscribe(
        (data: any) => {
          if (this.elasticsearchService.getIsSet() == true) {
            let transactionData = this.elasticsearchService.getData();
            console.log(transactionData);
            let messages = _.map(transactionData, '_source');
            messages = _.map(messages, 'message');

            console.log(this.results);
          }
        },
        (err) => { }
      );
    */

    // init search subscription
    this.initDebounceSearch();

    this.participantSubscription = this.sessionService.currentNodeChanged.subscribe(() => {
      // initial load of transactions
      this.refreshTransactions();
    });
  }

  ngOnDestroy() {

    this.loaded = false;

    this.transactionService.reloaded = true;

    this.participantSubscription.unsubscribe();
  }

  /**
   * Gets selected transactions based on transaction type.
   * This assumes this is for the current node in the session
   * unless otherwise specified
   * @param transactionType
   * @param isNewNode
   */
  getSelectedTransactions(transactionType: TransactionType, isNewNode = false) {

    // Guard against double click
    if (!this.transactionTypeDetails
      || this.transactionTypeDetails.key !== transactionType
      || isNewNode) {

      this.transactionService.reloaded = true;

      // re-initializing empty model header
      this.model.header = [];

      // clear asset option lists
      this.transactionService.sentAssetOptions = [];
      this.transactionService.receivedAssetOptions = [];
      this.transactionService.settlementAssetOptions = [];

      // get current transaction type details
      this.transactionTypeDetails = find(this.transactionService.transactionTypes,
        (type: TransactionTypeDetail) => {
          return type.key === transactionType;
        });

      // set as currently-viewed transaction
      this.transactionService.transactionType = transactionType;

      if (this.transactionTypeDetails.filters.text) {
        // set to last searched term if previously saved/set
        this.currSearchText = this.transactionTypeDetails.filters.text[0].value;
      } else {
        // reset reference to current search term
        this.currSearchText = '';
      }

      // initializing header columns (sortable)
      this.transactionTypeDetails.headers.forEach((header: any) => {
        const headerItem = new CustomTableHeaderItem({
          data: {
            value: header.name,
            displayName: header.displayName,
            tooltip: header.tooltip,
            class: header.class,
            hide: header.hide
          },
          sortable: header.sortable ? header.sortable : true,
          template: this.customHeaderTemplate
        });

        if (header.metadata) {
          headerItem.metadata = header.metadata;
        }
        this.model.header.push(headerItem);
      });

      // check to make sure there is a participant/institution/node attached
      if (isEmpty(this.sessionService.institutionNodes) || !this.sessionService.currentNode) {
        this.loaded = true;
      } else {

        if (this.model.data.length === 1) {
          this.transactionService.reloadedSubject.asObservable().subscribe((reloaded: boolean) => {

            // reset view data
            this.model.data = [[]];

            if (reloaded) {

              this.getTransactionsForCurrentNode();

              // this.transactionService.toggleReloaded();
              this.transactionService.reloaded = false;
            }
          }, (err: any) => {
            this.loaded = true;
          });
        } else {

          // set to false to reset to true when reloading
          this.transactionService.reloaded = false;
          this.transactionService.toggleReloaded();
        }
      }
    }

  }

  /**
   * Loads transactions for the current participant node.
   * This is used for both exchanges and transfers
   */
  getTransactionsForCurrentNode() {

    // reset loaded state when loading transactions
    // for another institution/participant, transaction type, filter
    if (this.loaded === true) {
      this.loaded = false;
    }
    // Subscribe to transactions
    this.transactionService.$transactions = this.transactionService.loadTransactions(
      this.sessionService.currentNode,
      this.transactionService.transactionType
    );

    this.transactionService.$transactions.subscribe((transactions: (ITransfer | IExchange)[]) => {

      // initializing current length of records
      // to the records per load setting
      if (!this.currLoadLength) {
        this.currLoadLength = this.transactionService.perLoad;
      }

      switch (this.transactionService.transactionType) {
        case 'exchange': this.loadExchangeView(transactions as IExchange[]);
          break;
        default: this.loadTransferView(transactions as ITransfer[]);
          break;
      }

      // finished loading transactions
      this.loaded = true;
    }, (err: any) => {
      console.log(err);
      this.loaded = true;
    });
  }

  /**
   * Load View for Exchange transactions
   *
   * @param {IExchange[]} transactions
   * @memberof TransactionsComponent
   */
  loadExchangeView(transactions: IExchange[]) {

    for (let iterator = 0; iterator < transactions.length; iterator++) {

      if (iterator === this.currLoadLength) {
        break;
      }

      const previousDataLength = this.model.totalDataLength;

      if (iterator < previousDataLength) {
        continue;
      }

      const transaction: IExchange = transactions[iterator];

      const expandedDataFields: CustomTableItem[] = [];

      if (transaction.ExchangeReceipt.transaction_hash) {
        expandedDataFields.push({
          name: 'Stellar Transaction Id',
          value: transaction.ExchangeReceipt.transaction_hash
        });
      } else {
        expandedDataFields.push({
          name: 'Error Message',
          value: ''
        });
      }

      const row = [
        new TableItem({
          data: {
            class: 'break-word max-data-width',
            value: transaction.ExchangeReceipt.exchange.quote.quote_id
          },
          template: this.tableItemTemplate,
        }),
        new TableItem({
          data: {
            value: transaction.account_name
          },
          template: this.tableItemTemplate,
        }),
        new TableItem({
          data: {
            value: transaction.counterparty
          },
          template: this.tableItemTemplate,
        }),
        new TableItem({
          data: {
            // Sent Amount
            value: transaction.reversed ?
              transaction.ExchangeReceipt.transacted_amount_target :
              transaction.ExchangeReceipt.transacted_amount_source
          },
          template: this.tableItemTemplate,
        }),
        new TableItem({
          data: {
            // Sent Asset
            value: transaction.asset_sent.asset_code
          },
          template: this.tableItemTemplate,
        }),
        new TableItem({
          data: {
            // Received Amount
            value: transaction.reversed ?
              transaction.ExchangeReceipt.transacted_amount_source :
              transaction.ExchangeReceipt.transacted_amount_target
          },
          template: this.tableItemTemplate,
        }),
        new TableItem({
          data: {
            // Received Asset
            value: transaction.asset_received.asset_code
          },
          template: this.tableItemTemplate,
        }),
        new TableItem({
          data: {
            // timestamp
            value: this.utils.toLocaleDateTime(transaction.time_stamp)
          },
          template: this.tableItemTemplate,
        }),
        new TableItem({
          data: {
            // Final Status
            value: transaction.ExchangeReceipt.status_exchange
          },
          template: this.statusTemplate,
          expandedData: {
            class: 'sub-table',
            fields: expandedDataFields
          },
          expandedTemplate: this.customTableItemTemplate
        })
      ];

      this.model.addRow(row);

      // add asset codes to existing lists
      if (!this.transactionService.sentAssetOptions.includes(transaction.asset_sent.asset_code)) {
        this.transactionService.sentAssetOptions.push(transaction.asset_sent.asset_code);
      }

      if (!this.transactionService.receivedAssetOptions.includes(transaction.asset_received.asset_code)) {
        this.transactionService.receivedAssetOptions.push(transaction.asset_received.asset_code);
      }
    }
  }

  /**
   * Load View for Transfer transactions
   *
   * @param {ITransfer[]} transactions
   * @memberof TransactionsComponent
   */
  loadTransferView(transactions: ITransfer[]) {
    for (let iterator = 0; iterator < transactions.length; iterator++) {
      const transaction: ITransfer = transactions[iterator];

      if (iterator === this.currLoadLength) {
        break;
      }

      const previousDataLength = this.model.totalDataLength;

      if (iterator < previousDataLength) {
        continue;
      }
      const transactionDetails = transaction.fitoficctnonpiidata
        .transactiondetails;

      const expandedDataFields: CustomTableItem[] = [
        {
          name: 'Pay-in Currency',
          value: transaction.payIn ? transaction.payIn : 'N/A'
        },
        {
          name: 'Pay-out Currency',
          value: transaction.payOut ? transaction.payOut : 'N/A'
        },
        {
          name: 'Quoted Price',
          value: transaction.fitoficctnonpiidata.exchange_rate
        },
        {
          name: 'Counterparty Fee',
          delimiter: '+',
          value: (transactionDetails.feecreditor && transactionDetails.feecreditor.cost) ? transactionDetails.feecreditor.cost : 'N/A'
        }
      ];

      // for payment 'cancellation' transaction type
      if (transaction.fitoficctnonpiidata.instruction_id &&
        transaction.fitoficctnonpiidata.original_instruction_id !== transaction.fitoficctnonpiidata.instruction_id) {

        // Move to front of array of expanded data fields
        expandedDataFields.unshift({
          name: 'Original Instruction Id',
          value: transaction.fitoficctnonpiidata.original_instruction_id,
          class: 'break-word',
          containerClass: 'status-details-container'
        });
      }

      // SHow Stellar/Ledger ID for those with one generated
      if (transaction.transactionid) {

        expandedDataFields.push({
          name: 'Stellar Id/s',
          tooltip: 'Block ID of transaction on Stellar ledger',
          value: transaction.transactionid,
          class: 'break-word',
          containerClass: 'status-details-container'
        });
      }

      // Show Status Detail for ongoing and unsuccessful transactions
      if (this.transactionService.getStatusName(transaction.current_status) !== 'settled') {
        expandedDataFields.push({
          name: 'Status Details',
          value: this.transactionService.getStatusDetail(transaction.current_status, transaction.message_type),
          containerClass: 'status-details-container'
        });
      }

      // constructing data row
      const row = [
        new TableItem({
          data: {
            value: transaction.fitoficctnonpiidata.instruction_id ? transaction.fitoficctnonpiidata.instruction_id : 'MISSING - ' + this.utils.toLocaleDateTime(transaction.time_stamp)
          },
          template: this.tableItemTemplate,
        }),
        new TableItem({
          data: {
            value: transaction.account_name ? transaction.account_name : 'Not Specified'
          },
          template: this.tableItemTemplate
        }),
        new TableItem({
          data: {
            value: transaction.counterparty
          },
          template: this.tableItemTemplate
        }),
        new TableItem({
          data: {
            pre: transaction.payment_type,
            space: true,
            value: transaction.fitoficctnonpiidata.transactiondetails.amount_settlement ? transaction.fitoficctnonpiidata.transactiondetails.amount_settlement : 0
          },
          template: this.tableItemTemplate
        }),
        new TableItem({
          data: {
            value: transactionDetails.assetsettlement ? transactionDetails.assetsettlement.asset_code : 'N/A'
          },
          template: this.tableItemTemplate
        }),
        new TableItem({
          data: {
            value: transaction.message_type ? startCase(transaction.message_type) : 'N/A'
          },
          template: this.tableItemTemplate
        }),
        new TableItem({
          data: {
            value: this.utils.toLocaleDateTime(transaction.time_stamp)
          },
          template: this.tableItemTemplate
        }),
        new TableItem({
          data: {
            type: transaction.message_type,
            value: startCase(transaction.current_status)
          },
          template: this.statusTemplate
        }),
        new TableItem({
          data: {
            status: transaction.transaction_status,
            type: transaction.message_type
          },
          template: this.overflowMenuItemTemplate,
          expandedData: {
            class: 'sub-table',
            fields: expandedDataFields
          },
          expandedTemplate: this.customTableItemTemplate
        })
      ];

      this.model.addRow(row);

      // add asset code to existing list
      if (transactionDetails.assetsettlement && !this.transactionService.settlementAssetOptions
        .includes(transactionDetails.assetsettlement.asset_code)) {
        this.transactionService.settlementAssetOptions
          .push(transactionDetails.assetsettlement.asset_code);
      }
    }
  }


  /**
   * Sort model data by index.
   * Pulled from carbon-components-angular table story
   *
   * @param {number} index
   * @param {boolean} [reload=false]
   * @memberof TransactionsComponent
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
   * Refresh entire page of transaction
   *
   * @memberof TransactionsComponent
   */
  refreshTransactions(): void {

    // reset load length
    this.currLoadLength = 0;

    // get transactions for new node
    this.getSelectedTransactions(this.transactionService.transactionType, true);
  }

  /**
   * Load more transactions onto the page
   *
   * @memberof TransactionsComponent
   */
  loadMore(): void {
    this.currLoadLength = this.currLoadLength + this.transactionService.perLoad;

    this.getTransactionsForCurrentNode();

    // resort based on currently sorted index, if any
    if (this.currSortIndex) {
      this.customSort(this.currSortIndex, true);
    }
  }

  /**
   * Brings up modal for export options
   *
   * @memberof TransactionsComponent
   */
  openExportModal(): void {

    // creates and opens the modal
    this.modalService.create({
      component: ExportModalComponent,
      inputs: {
        MODAL_DATA: this.model
      }
    });
  }

  /**
   * Copied from Office/AccountsComponent
   * Initializes subscription to search event
   *
   * @memberof TransactionsComponent
   */
  initDebounceSearch() {

    this.keyUp.pipe(
      // value to tap into
      map((event: any) => (<HTMLInputElement>event.target).value),
      // delay 1000
      debounceTime(1000),
      // don't fire unless new value
      distinctUntilChanged(),
      // update searched items
      tap(
        searchText => {
          this.searchTransactions(searchText);
        }
      )
    ).subscribe();

  }

  /**
   * Clears search input
   *
   * @memberof TransactionsComponent
   */
  clearSearch() {
    this.currSearchText = '';

    // reload search results
    this.searchTransactions(this.currSearchText);
  }

  /**
   * Make request to service to reload transactions
   * based on search terms and/or filters
   *
   * @private
   * @param {string} searchText
   * @memberof TransactionsComponent
   */
  private searchTransactions(searchText: string) {

    // trim whitespace to enable clean comparison
    searchText = searchText.toLowerCase().trim();

    // prevent duplicate search requests after trimmed whitespace
    if (searchText !== this.prevSearchText) {

      this.transactionTypeDetails.filters['text'] = [];

      this.transactionTypeDetails.search_terms.forEach((term: string) => {
        this.transactionTypeDetails.filters['text'].push({
          key: term,
          logic: 'text',
          value: searchText
        });
      });

      this.prevSearchText = searchText;

      this.transactionService.toggleReloaded();
    }
  }

  /**
   * View History of Statuses for a given transaction in a pop up modal
   *
   * @param {TransactionStatus[]} statuses
   * @param {string} message_type
   * @memberof TransactionsComponent
   */
  viewStatusHistory(statuses: TransactionStatus[], message_type: string) {

    // creates and opens the modal
    this.modalService.create({
      component: StatusHistoryModalComponent,
      inputs: {
        MODAL_DATA: {
          statuses: statuses,
          message_type: message_type
        }
      }
    });
  }
}
