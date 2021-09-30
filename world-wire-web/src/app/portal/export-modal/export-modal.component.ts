// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Inject } from '@angular/core';
import { BaseModal } from 'carbon-components-angular/modal/base-modal.class';
import { ExportService } from '../shared/services/export.service';
import { ExportData } from '../shared/models/export-data.model';
import { BookType } from 'xlsx';
import { TableModel } from 'carbon-components-angular';

@Component({
  selector: 'app-export-modal',
  templateUrl: './export-modal.component.html',
  styleUrls: ['./export-modal.component.scss']
})
export class ExportModalComponent extends BaseModal implements OnInit {

  dataModel: TableModel;

  constructor(
    @Inject('MODAL_DATA') public data: TableModel,
    public exportService: ExportService
  ) {
    super();

    this.dataModel = data;
  }

  ngOnInit() {
  }

  async exportData(extension: BookType) {
    this.exportService.exportInProgress = true;

    // close modal after quick visual feedback

    const exportDataObj: ExportData = await this.exportService.processExportData(this.dataModel);

    this.exportService.processWorkSheet(exportDataObj, extension).then(() => {

      this.closeModal();

      // export is finished
      this.exportService.exportInProgress = false;
    });

  }


}
