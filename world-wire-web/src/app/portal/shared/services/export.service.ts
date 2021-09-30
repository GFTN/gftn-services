// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable } from '@angular/core';
import * as XLSX from 'xlsx';
import { ExportData } from '../models/export-data.model';
import { SessionService } from '../../../shared/services/session.service';
import { TableModel, TableHeaderItem, TableItem } from 'carbon-components-angular';

@Injectable()
export class ExportService {

    public exportInProgress = false;

    private timer: any;

    constructor(private sessionService: SessionService) { }

    /**
     * Process transaction data from user selections (data model) for export
     *
     * @param {TableModel} model
     * @returns {ExportData}
     * @memberof ExportService
     */
    processExportData(model: TableModel): Promise<ExportData> {

        return new Promise<ExportData>((resolve) => {
            const fields: string[] = [];

            // get header info for export
            model.header.forEach((headerItem: TableHeaderItem) => {
                fields.push(headerItem.data.value);
            });

            // will hold final data for export
            // XLSX expects an Array<Array>
            const myData: Array<Array<any>> = [];

            const expandedKeyExists = [];

            // get first-level table row data for export
            model.data.forEach((tableRows: TableItem[]) => {
                const rowData = [];
                for (let i = 0; i < tableRows.length; i++) {
                    const tableRow = tableRows[i];

                    // header is not null
                    if (fields[i]) {
                        rowData[fields[i]] = tableRow.data.value;
                    }

                    // get expanded level data
                    if (tableRow.expandedData) {

                        const expandedRowData = tableRow.expandedData.fields;

                        // append expanded row data
                        for (let j = 0; j < expandedRowData.length; j++) {

                            this.processExpandedField(expandedKeyExists, expandedRowData[j], fields, rowData);
                        }
                    }
                }

                myData.push(rowData);
            });

            // Filename: Institution Name + today's date & time
            // Should be unique unless requested twice in the same second
            const fileNameBase = this.sessionService.institution.info.slug;

            // return export data for processing
            const exportData: ExportData = {
                jsonData: myData,
                fileName: fileNameBase,
                date: true
            };

            resolve(exportData);
        });
    }

    /**
     * Process expanded Field from table model
     *
     * @param {string[]} expandedKeyExists
     * @param {*} field
     * @param {string[]} fields
     * @param {*} rowData
     * @memberof ExportService
     */
    processExpandedField(expandedKeyExists: string[], field, fields: string[], rowData) {

        // Expanded field value is a list of fields
        if (Array.isArray(field.value) && typeof field.value[0] !== 'string') {
            for (const item of field.value) {
                this.processExpandedField(expandedKeyExists, item, fields, rowData);
            }
        } else {

            // add expanded data field key if it doesn't already exist
            if (!expandedKeyExists.includes(field.name)) {

                // track expanded row field key for local reference
                expandedKeyExists.push(field.name);

                // append expanded field key to current list of all keys
                fields.push(field.name);

            }

            // construct string value for export
            const rowDataValue: string = Array.isArray(field.value) ? field.value.join(';') : field.value;

            // use field key to map to value for that field
            rowData[field.name] = rowDataValue;
        }
    }

    /**
   * Processing export in a promise
   * to put process in background
   * @param extension
   */
    public processWorkSheet(
        exportDataObj: ExportData,
        extension: XLSX.BookType
    ): Promise<boolean> {
        return new Promise((resolve) => {

            clearTimeout(this.timer);

            // delay validation to allow for data processing
            this.timer = setTimeout(() => {

                try {

                    // console.log(exportDataObj.jsonData);

                    // generate worksheet
                    const ws: XLSX.WorkSheet = XLSX.utils.json_to_sheet(exportDataObj.jsonData);

                    const wb: XLSX.WorkBook = XLSX.utils.book_new();
                    XLSX.utils.book_append_sheet(wb, ws, 'Transactions');

                    let fileName = exportDataObj.fileName;

                    // append date to filename, if option enabled
                    if (exportDataObj.date) {

                        // generate today's date for naming
                        const date = new Date(
                            Date.now()).toLocaleDateString(
                                navigator.language,
                                {
                                    year: 'numeric',
                                    month: 'short',
                                    day: 'numeric'
                                }
                            ).trim()
                            .replace(/[^a-zA-Z0-9 ]/g, ' ')
                            .replace(/ +/g, '.');

                        const time = new Date(
                            Date.now()).toLocaleTimeString(
                                navigator.language,
                                {
                                    hour: 'numeric',
                                    minute: 'numeric',
                                    second: 'numeric'
                                }
                            ).trim()
                            .replace(/[^a-zA-Z0-9 ]/g, ' ')
                            .replace(/ +/g, '.');

                        const today = `${date}_${time}`;

                        // appending date to filename base
                        fileName = `${fileName}_${today}`;
                    }

                    // append extension
                    fileName = `${fileName}.${extension}`;

                    // save to file
                    XLSX.writeFile(wb, fileName, { bookType: extension });

                    resolve(true);
                } catch (error) {

                    console.error('export failed: ', error);
                    resolve(false);
                }
            }, 1500);
        });
    }

}
