// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { CheckboxOption } from '../../../../shared/models/checkbox-option.model';
import { CheckboxGroupFilter } from '../../models/filter.model';

@Component({
  selector: 'app-quick-filter',
  templateUrl: './quick-filter.component.html',
  styleUrls: ['./quick-filter.component.scss']
})
export class QuickFilterComponent implements OnInit {

  @Input() currentData: any;

  @Input() allData: any;

  @Input() filters: CheckboxGroupFilter;

  // optional layout for 'sm' breakpoint
  // used for narrower filter columnns
  @Input() smLayout: 'column' | 'row' | 'row wrap' = 'column';

  smLayoutStyle: 'row' | 'column' = 'row';

  // notify parent component of filter change event
  @Output() changed = new EventEmitter<CheckboxGroupFilter>();

  constructor() { }

  ngOnInit() {
    this.smLayoutStyle = this.smLayout.includes('row') ? 'column' : 'row';
  }

  /**
   *  Uncheck all options of the checkbox
   *
   * @param {CheckboxOption[]} options
   * @memberof QuickFilterComponent
   */
  clearFilterOptions(options: CheckboxOption[]) {
    for (const option of options) {
      option.checked = false;
    }

    this.changed.emit(this.filters);
  }

  /**
   * Get number of checked options
   *
   * @param {CheckboxOption[]} options
   * @returns {number}
   * @memberof QuickFilterComponent
   */
  getCheckedOptions(options: CheckboxOption[]): number {
    let checkedOptions = 0;
    for (const option of options) {
      if (option.checked) {
        checkedOptions++;
      }
    }

    return checkedOptions;
  }

  /**
   * Resets all filters
   *
   * @memberof QuickFilterComponent
   */
  resetFilters() {
    for (const filterVal of Object.values(this.filters)) {
      this.clearFilterOptions(filterVal.options);
    }

    this.changed.emit(this.filters);
  }
}
