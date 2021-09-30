// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Inject } from '@angular/core';
import { BaseModal } from 'carbon-components-angular';
import { TransactionStatus, TransactionMessageType } from '../../shared/models/transaction.interface';
import { UtilsService } from '../../shared/utils/utils';
import { TransactionService } from '../shared/services/transaction.service';

/**
 * Data necessary for status history modal
 *
 * @export
 * @interface StatusModalData
 */
export interface StatusModalData {
  statuses: TransactionStatus[];
  message_type: TransactionMessageType;
}

@Component({
  selector: 'app-status-history-modal',
  templateUrl: './status-history-modal.component.html',
  styleUrls: ['./status-history-modal.component.scss']
})
export class StatusHistoryModalComponent extends BaseModal implements OnInit {

  statuses: TransactionStatus[];

  message_type: TransactionMessageType;

  constructor(
    @Inject('MODAL_DATA') public data: StatusModalData,
    public utils: UtilsService,
    public transactionService: TransactionService
  ) {
    super();

    this.statuses = data.statuses;
    this.message_type = data.message_type;
  }

  ngOnInit() {
  }

}
