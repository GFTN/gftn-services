// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable } from '@angular/core';

@Injectable()

/**
 * Stores data for a file export
 */
export class ExportData {

    // object to store json data for export
    // will typically be Array<Array<any>>
    // according to spec for XLSX
    // https://github.com/SheetJS/js-xlsx
    jsonData: Array<Array<any>>;

    fileName: string;

    // optional: toggle date true/false
    // to append current date to filename
    date?: boolean;
}
