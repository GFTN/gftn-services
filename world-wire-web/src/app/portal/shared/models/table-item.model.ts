// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { TableItem, TableHeaderItem } from 'carbon-components-angular';

/**
 * Custom Table Item class for use in our data table.
 * This will be provided in the "expanded" portion of a row
 * of our data table
 */
export class CustomTableItem {
    name: string;
    value: any;
    tooltip?: string;
    delimiter?: string;
    class?: string; // optional classes for table item
    containerClass?: string;
}

/**
 * Custom Table Header Item class used in our data table
 */
export class CustomTableHeaderItem extends TableHeaderItem {

    // override compare function to support manual sort based on custom template
    compare(one: TableItem, two: TableItem): number {

        const metaData = this.metadata ? this.metadata : {};

        // previous: one.data vs. two.data
        // now: one.data.value vs. two.data.value
        const firstValue = metaData.date ? new Date(one.data.value) : one.data.value;
        const secondValue = metaData.date ? new Date(two.data.value) : two.data.value;

        if (firstValue < secondValue) {
            return -1;
        } else if (firstValue > secondValue) {
            return 1;
        } else {
            return 0;
        }
    }

}
