// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
//
import { Injectable, NgZone } from '@angular/core';
import { ContentSwitcherOption } from '../../../shared/models/content-switcher-option.model';
import {
  ITransfer, IExchange,
  TransactionTypeDetail, TransactionType,
  TransactionKeyBook,
  TransactionStatus, TransactionStatusMapping
} from '../../../shared/models/transaction.interface';
import { INodeAutomation } from '../../../shared/models/node.interface';
import { Filter } from '../models/filter.model';

import { PopoverComponent } from '../../../shared/custom-carbon.module';

import { Observable, Observer, BehaviorSubject } from 'rxjs';
import { AngularFireDatabase } from '@angular/fire/database';

import * as _ from 'lodash';

@Injectable()
export class TransactionService {

  /**
   * Transaction dashboard load settings and data stores
   */
  public perLoadOptions: number[] = [20, 40, 100, 200];

  public perLoad: number;

  public tableSize = 'md';

  public tableDisplayOptions: ContentSwitcherOption[];

  public stripedRows = false;

  // default logic options available for comparison
  public logicOptions: string[] = ['=', '>=', '>', '<', '<='];

  // keep track of current popover opened
  // to reference when getting data,
  // doing actions such as closing the popover/etc.
  public popoverRef: PopoverComponent;

  // source/target asset based on OFI perspective
  public sentAssetOptions: string[] = [];

  // source/target asset based on OFI perspective
  public receivedAssetOptions: string[] = [];

  // TRANSFERS ONLY: asset the transaction was settled in
  public settlementAssetOptions: string[] = [];

  public $transactions: Observable<(ITransfer | IExchange)[]>;

  public transactionType: TransactionType = 'transfer';

  public reloaded = true;

  // keeps track of data load/reload
  public reloadedSubject: BehaviorSubject<boolean>;

  // set status colors and names from transaction status codes
  public statusColorMap: Map<string, string>;
  public statusReadableMappings: {
    [key: string]: TransactionStatusMapping | { [key: string]: TransactionStatusMapping };
  };

  // keep track of transaction keys for consistency
  // so if these ever change, we only have to change it in one place
  // readonly: can never be modified by any component that uses it
  readonly transactionDetails = 'fitoficctnonpiidata.transactiondetails';

  readonly exchangeReceipt = 'ExchangeReceipt';
  readonly exchange = this.exchangeReceipt + '.exchange';
  readonly exchangeQuote = this.exchange + 'quote';

  public readonly transactionKeys: { [key: string]: TransactionKeyBook } = {
    'transfer': {
      id: 'fitoficctnonpiidata.instruction_id',
      originalId: 'fitoficctnonpiidata.original_instruction_id',
      transactionStatus: 'current_status',
      paymentType: 'payment_type',
      accountName: 'account_name',
      sentAmount: null,
      sentAssetCode: null,
      receivedAmount: null,
      receivedAssetCode: null,
      settlementAmount: this.transactionDetails + '.amount_settlement',
      settlementAssetCode: this.transactionDetails + '.assetsettlement.asset_code',
      timeStamp: 'time_stamp',
      stellarId: 'transaction_identifier',
      counterparty: 'counterparty',
    },
    'exchange': {
      id: this.exchange + '.quote.quote_id',
      transactionStatus: this.exchangeReceipt + '.status_exchange',
      accountName: 'account_name',
      paymentType: null,
      sentAmount: 'amount_sent',
      sentAssetCode: 'asset_sent.asset_code',
      receivedAmount: 'amount_received',
      receivedAssetCode: 'asset_received.asset_code',
      settlementAmount: null,
      settlementAssetCode: null,
      timeStamp: 'time_stamp',
      stellarId: this.exchangeReceipt + '.transaction_hash',
      counterparty: 'counterparty',
    }
  };

  public transactionTypes: TransactionTypeDetail[] = [
    {
      'key': 'transfer',
      'name': 'Transfers',
      'headers': [
        { name: 'Instruction Id' },
        { name: 'Account' }, // Formerly Distribution Account
        { name: 'Counterparty' },
        { name: 'Settlement Amount', tooltip: 'Incoming(+)/Outgoing(-) Amount' },
        { name: 'Settlement Currency' },
        { name: 'Type' },
        { name: 'Time Initiated', metadata: { date: true } },
        { name: 'Status' },
        { name: null, sortable: false },
      ],
      'transactionStatuses': [
        { name: 'initiated', label: 'Initiated' },
        { name: 'processing', label: 'Processing' },
        { name: 'validated', label: 'Validated' },
        { name: 'cleared', label: 'Cleared' },
        { name: 'settled', label: ' Settled' },
        { name: 'rejected', label: 'Rejected' },
        { name: 'failed', label: 'Failed' }
      ],
      'search_terms': [
        this.transactionKeys.transfer.id,
        this.transactionKeys.transfer.originalId,
        this.transactionKeys.transfer.stellarId,
        this.transactionKeys.transfer.accountName,
        this.transactionKeys.transfer.counterparty
      ],
      'searchPlaceholderText': 'Search by Account, Instruction Id, Original Instruction Id, Stellar Transaction ID, Counterparty',
      'filters': {}
    },
    {
      'key': 'exchange',
      'name': 'Exchanges',
      'headers': [
        { name: 'Quote ID' },
        { name: 'Account' },
        { name: 'Counterparty' },
        { name: 'Sent Amount', displayName: 'Sent', tooltip: 'Sent Amount' },
        { name: 'Sent Currency', displayName: 'CCY', tooltip: 'Sent Currency' },
        { name: 'Received Amount', displayName: 'Received', tooltip: 'Received Amount' },
        { name: 'Received Currency', displayName: 'CCY', tooltip: 'Received Currency' },
        { name: 'Time Initiated', metadata: { date: true } },
        { name: 'Status' },
      ],
      'transactionStatuses': [
        { name: 'success', label: 'Success' },
        { name: 'rejected', label: 'Rejected' }
      ],
      'search_terms': [
        this.transactionKeys.exchange.id,
        this.transactionKeys.exchange.stellarId,
        this.transactionKeys.exchange.accountName,
        this.transactionKeys.exchange.counterparty
      ],
      'searchPlaceholderText': 'Search by Account, Quote ID, Stellar Transaction ID, Counterparty',
      'filters': {}
    },
  ];

  constructor(
    private db: AngularFireDatabase,
    private ngZone: NgZone,
  ) {
    this.reloadedSubject = new BehaviorSubject<boolean>(true);
  }

  /**
   * initializing table display options if not set
   * possible options for table display density
   * pulled from carbon-components-angular
   * md = normal, sm = compact, lg = tall
   */
  public initializeTableDisplayOptions() {
    this.tableDisplayOptions = [{
      id: 'md',
      name: 'Normal',
      selected: false
    }, {
      id: 'sm',
      name: 'Compact',
      selected: false
    },
    {
      id: 'lg',
      name: 'Tall',
      selected: false
    }];

    // set selected option
    for (const option of this.tableDisplayOptions) {
      if (this.tableSize === option.id) {
        option.selected = true;
      }
    }
  }

  /**
   * Load Transactions based on the institution name given
   * This is here in a service so other components (like Overview)
   * can use the same data
   * @param institution
   */
  public loadTransactions(
    // institution: string,
    currentNode: INodeAutomation,
    transactionType: TransactionType = 'transfer'
  ): Observable<(ITransfer | IExchange)[]> {

    const source = new Observable((observer: Observer<(ITransfer | IExchange)[]>) => {

      const dbRef: firebase.database.Reference = this.db.database.ref(`v1/txn/${transactionType}/${currentNode.participantId}`);

      // Get transfer transactions from firebase
      // by participant (as OFI or RFI)
      dbRef.on('value', (transactionData: any) => {

        let transactions: (ITransfer | IExchange)[];

        switch (transactionType) {
          case 'exchange': transactions = this.processExchangeTransactions(currentNode.participantId, transactionData);
            break;
          default: transactions = this.processTransferTransactions(currentNode.participantId, transactionData);
            break;
        }

        const transactionTypeDetails: TransactionTypeDetail = _.find(this.transactionTypes,
          (type: TransactionTypeDetail) => {
            return type.key === transactionType;
          });

        _.forEach(transactionTypeDetails.filters, (filterArray: Filter[], key: string) => {

          transactions = _.filter(transactions, (t: ITransfer | IExchange) => {

            return this.filterTransactions(t, key, filterArray);
          });

        });

        // Sort transactions by timestamp
        transactions = _.sortBy(
          // return in array format
          transactions,
          (t: ITransfer | IExchange) => {
            // sort case insensitive
            return new Date(t.time_stamp);
          }
        );

        // reverse to order from most to least recent
        transactions = transactions.reverse();

        // refresh the view
        this.ngZone.run(() => {
          // update observer value
          observer.next(
            transactions,
          );
        }); // end ngZone
      }); // end dbRef
    });

    return source;
  }

  filterTransactions(transaction: any, key: string, filterArray: Filter[]): boolean {
    let filteredObject: any = transaction;

    if (key !== 'text') {
      filteredObject = this.getFilteredObject(key, filteredObject);
    }

    let filterIterator = 0;

    // go through each option to check if the object matches filter
    for (const filter of filterArray) {

      let filteredObjectIncluded = false;

      switch (filter.logic) {
        case 'date': {

          // use date range as string for parsing
          filter.value = filter.value as string;
          const range = filter.value.split('|');
          const startAt = new Date(range[0]);
          const endAt = new Date(range[1]);

          // get end of day
          endAt.setDate(endAt.getDate() + 1);

          const tDate = new Date(filteredObject * 1000);

          if (tDate >= startAt && tDate < endAt) {
            filteredObjectIncluded = true;
          }
        }
          break;
        case 'text': {

          // reset filteredObject reference for each filter key
          filteredObject = transaction;

          filteredObject = this.getFilteredObject(filter.key, filteredObject);

          if (filteredObject) {
            // array of string objects
            if (Array.isArray(filteredObject)) {
              for (const object of filteredObject) {
                filteredObjectIncluded = this.filterText(object, filter);

                // return at first instance in array when found
                if (filteredObjectIncluded) {
                  return filteredObjectIncluded;
                }
              }
            } else if (typeof filteredObject === 'string') {
              // string object
              filteredObjectIncluded = this.filterText(filteredObject, filter);
            } else {
              // default false
              filteredObjectIncluded = false;
            }
          } else {
            filteredObjectIncluded = false;
          }
        }
          break;
        case '<':
          filteredObjectIncluded = filteredObject < filter.value;
          break;
        case '>':
          filteredObjectIncluded = filteredObject > filter.value;
          break;
        case '<=':
          filteredObjectIncluded = filteredObject <= filter.value;
          break;
        case '>=':
          filteredObjectIncluded = filteredObject >= filter.value;
          break;
        default:

          // allow for case-insensitive value comparisons if text
          if (typeof filter.value === 'string') {
            filteredObject = filteredObject as string;
            filteredObject = filteredObject.toLowerCase();

            filter.value = filter.value as string;
            filter.value = filter.value.toLowerCase();
          }
          filteredObjectIncluded = filteredObject === filter.value;
          break;
      }

      if (!filteredObjectIncluded && filterIterator < filterArray.length) {
        filterIterator++;
        continue;
      } else {
        return filteredObjectIncluded;
      }
    }
  }

  /**
   * Check if the text matches the filter
   *
   * @param {string} filteredObject
   * @param {*} filter
   * @returns {boolean}
   * @memberof TransactionService
   */
  filterText(filteredObject: string, filter): boolean {
    // allow for case-insensitive searches
    filteredObject = filteredObject.toLowerCase();
    filter.value = filter.value as string;

    // allow partial searches
    return filteredObject.includes(filter.value);
  }

  /**
   * Parses nested object for proper value comparison with filter key
   * @param keyString
   * @param filteredObject
   */
  getFilteredObject(keyString: string, filteredObject) {
    // store in cases when we get the initial transaction
    // as the object (when searching for search terms)
    const transaction = filteredObject;

    const keyArray = keyString.split('.');

    keyArray.forEach((keyPiece: string) => {
      filteredObject = filteredObject[keyPiece];
    });

    // get readable status name for status
    if (keyString === this.transactionKeys[this.transactionType].transactionStatus) {
      const messageType = transaction.message_type ? transaction.messageType : 'credit_transfer';
      filteredObject = this.getStatusName(filteredObject, messageType);
    }

    return filteredObject;
  }

  /**
   * Processes transactions as Exchanges
   * @param institution
   * @param transactionData
   */
  processExchangeTransactions(institution: string, transactionData: any): IExchange[] {
    const transactions: IExchange[] = [];

    transactionData.forEach((node) => {
      const transaction: IExchange = node.val();

      // if receipt is successful
      if (transaction.ExchangeReceipt) {
        // add to array of displayed transactions
        const ofi = transaction.OFIID;
        const rfi = transaction.RFIID;

        // get all transactions where participant is originator or receiver
        if (institution === ofi || institution === rfi) {

          // set vars based on OFI's POV (if viewing from single participant)
          if (institution === rfi) {
            transaction.reversed = true;
          }

          transaction.account_name = transaction.reversed ?
            transaction.ExchangeReceipt.exchange.account_name_receive :
            transaction.ExchangeReceipt.exchange.account_name_send;

          transaction.counterparty = transaction.reversed ? ofi : rfi;

          transaction.asset_sent = transaction.reversed ?
            transaction.ExchangeReceipt.exchange.quote.quote_request.target_asset :
            transaction.ExchangeReceipt.exchange.quote.quote_request.source_asset;

          transaction.asset_received = transaction.reversed ?
            transaction.ExchangeReceipt.exchange.quote.quote_request.source_asset :
            transaction.ExchangeReceipt.exchange.quote.quote_request.target_asset;

          transaction.amount_sent = transaction.reversed ?
            transaction.ExchangeReceipt.transacted_amount_target :
            transaction.ExchangeReceipt.transacted_amount_source;

          transaction.amount_received = transaction.reversed ?
            transaction.ExchangeReceipt.transacted_amount_source :
            transaction.ExchangeReceipt.transacted_amount_target;

          // save to same field as transfers for common filtering purposes
          transaction.time_stamp = transaction.ExchangeReceipt.time_executed;

          transactions.push(transaction);
        }
      }
    });
    return transactions;
  }

  /**
   * Processes transactions as Transfers
   * @param institution
   * @param transactionData
   */
  processTransferTransactions(institution: string, transactionData: any): ITransfer[] {
    let transactions: ITransfer[] = [];

    transactionData.forEach((node) => {
      const transaction: ITransfer = node.val().transaction_memo ? node.val().transaction_memo : node.val().TransactionMemo;

      // secondary check to prevent breakage of the structure
      if (transaction) {

        const ofi = (transaction.fitoficctnonpiidata && transaction.fitoficctnonpiidata.transactiondetails) ?
          transaction.fitoficctnonpiidata.transactiondetails.ofi_id : null;

        const rfi = (transaction.fitoficctnonpiidata && transaction.fitoficctnonpiidata.transactiondetails) ?
          transaction.fitoficctnonpiidata.transactiondetails.rfi_id : null;

        // get all transactions where participant is originator or receiver
        if (institution === ofi || institution === rfi) {

          // set payment vars based on OFI's POV (if viewing from single participant)
          if (institution === rfi) {
            transaction.reversed = true;
          }

          transaction.payment_type = transaction.reversed ? '-' : '+';

          // only .account_name_send exists right now
          // TODO: change this to .account_name_receive when https://github.com/GFTN/gftn-services/issues/1071 is done
          transaction.account_name = transaction.reversed ?
            transaction.fitoficctnonpiidata.account_name_send :
            transaction.fitoficctnonpiidata.account_name_receive;

          transaction.counterparty = transaction.reversed ? ofi : rfi;

          // get asset code of creditor
          // TODO: modify model in send-service so this doesn't have to depend on the fee
          const asset_code_creditor = transaction.fitoficctnonpiidata
            .transactiondetails.feecreditor ?
            transaction.fitoficctnonpiidata
              .transactiondetails.feecreditor.costasset.asset_code : '';

          // Asset of incoming asset being transferred
          transaction.payIn = transaction.reversed ?
            asset_code_creditor :
            transaction.fitoficctnonpiidata
              .transactiondetails.asset_code_beneficiary;

          // Asset of outgoing asset being transferred
          transaction.payOut = transaction.reversed ?
            transaction.fitoficctnonpiidata
              .transactiondetails.asset_code_beneficiary :
            asset_code_creditor;

          // get latest status
          const statusArray: TransactionStatus[] = transaction.transaction_status;

          transaction.current_status = statusArray[statusArray.length - 1].transactionstatus.trim();

          // add to array of displayed transactions
          transactions.push(transaction);
        }
      }
    }); // end foreach

    // filter by 'credit_transfer' and 'cancellation' types
    transactions = _.filter(transactions, (transaction: ITransfer) => {
      return transaction.message_type === 'cancellation' || transaction.message_type === 'credit_transfer';
    });

    return transactions;
  }

  /**
  * Initializes readable status names and tag colors
  * for transaction status codes
  */
  initializeTransactionStatuses() {
    this.statusReadableMappings = {
      // transfers
      'initialized': {
        'credit_transfer': {
          name: 'initiated',
          detail: 'Payment Initiated'
        },
        'cancellation': {
          name: 'initiated',
          detail: 'Cancellation Request Initiated'
        }
      },
      'validation_success': {
        'credit_transfer': {
          name: 'validated',
          detail: 'Successful Message Sent'
        },
        'cancellation': {
          name: 'validated',
          detail: 'Successful Cancellation Request Sent'
        }
      },
      'ofi_processing': {
        name: 'processing',
        detail: 'Message Successfully Received by OFI'
      },
      'ofi_validation_success': {
        name: 'validated',
        detail: 'Successful Message Sent'
      },
      'rfi_processing': {
        name: 'processing',
        detail: 'Message Successfully Received by RFI'
      },
      'rfi_validation_success': {
        name: 'validated',
        detail: 'Message Successfully Validated by RFI'
      },
      'cleared': {
        name: 'cleared',
        detail: 'Successful Federation & Compliance Check'
      },
      'unable_to_execute': {
        name: 'rejected',
        detail: 'RFI unable to disburse payment to Ultimate Beneficiary'
      },
      'request_to_modify_payment': {
        name: 'rejected',
        detail: 'Payment modification requested'
      },
      // RDO success will turn into 'settled' in the final step of the RDO flow.
      // This just will indicate if the request for RDO was successfully sent
      'rdo_initialized': {
        name: 'cleared',
        detail: 'DO Return Initiated'
      },
      'returned': {
        name: 'cleared',
        detail: 'Reversal Completed'
      },
      'settled': {
        'credit_transfer': {
          name: 'settled',
          // Fallback details, just in case ledger id doesn't get tracked
          detail: 'Settled'
        },
        'cancellation': {
          name: 'settled',
          detail: 'Reversal Completed'
        }
      },
      'validation_fail': {
        'credit_transfer': {
          name: 'rejected',
          detail: 'Invalid Message'
        },
        'cancellation': {
          name: 'rejected',
          detail: 'Invalid Cancellation Request'
        }
      },
      'failed': {
        name: 'failed',
        detail: 'Message Failed to Deliver',
      },
      // Cancellation Statuses
      'cancellation_initialized': {
        name: 'initiated',
        detail: 'Cancellation Requested by OFI'
      },
      'cancellation_rejected': {
        name: 'rejected',
        detail: 'Payment Cancellation Rejected'
      },
      'payment_returned': {
        name: 'cleared',
        detail: 'Payment Cancellation Completed'
      },
      'payment_rejected': {
        name: 'rejected',
        detail: 'Payment Rejected'
      },
      // exchanges
      'ok': {
        name: 'success',
        detail: 'Successful Exchange Completed'
      },
      'denied': {
        name: 'rejected',
        detail: 'Insuccessful Exchange Made'
      },
    };

    this.statusColorMap = new Map([
      // transfers
      ['initial', 'neutral-1'],
      ['processing', 'warning'],
      ['validated', 'neutral-1'],
      ['settled', 'success'],
      ['cleared', 'neutral-2'],
      ['rejected', 'failure'],
      ['failed', 'failure'],
      // exchanges
      ['success', 'success'],
      ['rejected', 'failure'],
    ]);
  }

  /**
   * Clean Status for reading from map
   * just in case the status gets malformed
   *
   * @param {string} statusCode
   * @returns {string}
   * @memberof TransactionService
   */
  cleanStatusCode(statusCode: string): string {
    // added null check
    if (statusCode) {
      let cleanStatus = statusCode.toLowerCase().trim();

      // enforce snake case in case it gets lost
      cleanStatus = _.snakeCase(cleanStatus);

      return cleanStatus;
    } else {
      return null;
    }
  }

  /**
   * Get user-friendly status name
   * based on status codes returned from API
   *
   * @param {string} status
   * @param {string} [messageType='credit_transfer']
   * @returns {string}
   * @memberof TransactionService
   */
  getStatusName(status: string, messageType: string = 'credit_transfer'): string {

    status = this.cleanStatusCode(status);

    let statusName = '';

    if (this.statusReadableMappings[status]) {

      if (this.statusReadableMappings[status].name) {
        const primaryName = this.statusReadableMappings[status].name as string;
        statusName = primaryName ? primaryName : statusName;
      } else {
        statusName = this.statusReadableMappings[status][messageType] ? this.statusReadableMappings[status][messageType].name : statusName;
      }
    }

    return statusName;
  }

  /**
   * Get color of status tag
   * Based on status name
   *
   * @param {string} statusCode
   * @returns {string}
   * @memberof TransactionService
   */
  getStatusColor(statusCode: string): string {
    // get common status based on code
    const status = this.getStatusName(statusCode);

    if (this.statusColorMap.get(status)) {
      return this.statusColorMap.get(status);
    }

    // default to initial color if status doesn't exist
    return this.statusColorMap.get('initial');
  }

  /**
   * Get user-friendly status name
   * based on status codes returned from API
   * @param status
   */
  getStatusDetail(status: string, messageType: string = 'credit_transfer'): string {
    status = this.cleanStatusCode(status);

    let statusDetail = 'N/A';

    if (this.statusReadableMappings[status]) {

      if (this.statusReadableMappings[status].detail) {
        const primaryDetail = this.statusReadableMappings[status].detail as string;
        statusDetail = primaryDetail ? primaryDetail : statusDetail;
      } else {
        statusDetail = this.statusReadableMappings[status][messageType] ? this.statusReadableMappings[status][messageType].detail : statusDetail;
      }
    }

    return statusDetail;
  }

  /**
   * Lets transaction service know
   * if the data needs to be reloaed
   */
  public toggleReloaded() {
    this.reloaded = !this.reloaded;
    this.reloadedSubject.next(this.reloaded);
  }
}
