// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { CheckboxOption } from '../../../shared/models/checkbox-option.model';

export class Filter {
    key?: string;
    logic = '';
    value: string | number;
}

export interface CheckboxGroup {
    name: string;
    options: CheckboxOption[];
}

export interface CheckboxGroupFilter {
    [name: string]: CheckboxGroup;
}
