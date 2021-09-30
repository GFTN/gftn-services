// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Input, AfterViewInit, ViewChild, ElementRef, forwardRef, NgZone, OnChanges } from '@angular/core';
import { Dropdown } from 'carbon-components';
import { IDropdownOption as DropdownOption } from '../../models/dropdown-option.interface';
import { ControlValueAccessor, NG_VALUE_ACCESSOR, NG_VALIDATORS, AbstractControl, Validator } from '@angular/forms';

@Component({
  selector: 'app-dropdown',
  templateUrl: './dropdown.component.html',
  styleUrls: ['./dropdown.component.scss'],
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => DropdownComponent),
      multi: true
    },
    {
      provide: NG_VALIDATORS,
      useExisting: forwardRef(() => DropdownComponent),
      multi: true,
    }
  ]
})
export class DropdownComponent implements OnInit, AfterViewInit, OnChanges, ControlValueAccessor, Validator {

  dropdownInstance: any;

  @Input() label: string;
  @Input() placeholder: string;
  @Input() options: DropdownOption[];
  @Input() defaultValue: string;
  @Input() disabled: boolean;

  @Input() model: string;

  @ViewChild('dropDownElm') dropDownElm: ElementRef;
  @ViewChild('dropdownText') dropdownText: ElementRef;
  @ViewChild('dropdownInput') dropdownInput: ElementRef;

  // valid: boolean = false;
  error: boolean;

  constructor(private ngZone: NgZone) {}

  ngOnInit() {
  }

  ngAfterViewInit() {
    // initialize dropdown
    this.dropdownInstance = Dropdown.create(this.dropDownElm.nativeElement);

    // selecting default value
    // from dropdown if set
    if (this.defaultValue != null) {
      this.selectOption(this.defaultValue);
    }

    // Notify Angular and update model upon value change
    // from Carbon Design's vanilla JS dropdown
    document.addEventListener('dropdown-beingselected', () => {
      Promise.resolve(null).then(() => {
        this.ngZone.run(() => {
          this.writeValue(this.getValue());
        });
      });
    });
  }

  ngOnChanges() {

    // selecting default value
    // from dropdown if set
    if (this.dropdownInstance && this.defaultValue != null) {
      this.selectOption(this.defaultValue);
    }
  }

  /**
   * GETS value of dropdown
   */
  public getValue() {
    return this.dropDownElm.nativeElement.dataset.value;
  }

  /**
   *
   * @param id
   * SETS selected dropdown option manually
   */
  public selectOption(value) {

    let selectedOption;

    // get from list of options
    // with matching value
    for (let i = 0; i < this.options.length; i++) {
      const option = this.options[i];
      if (option.value === value) {
        selectedOption = document.getElementById(`${option.name}-${i}`);
        break;
      }
    }

    // wrapping setting default in promise
    // to solve async issues with getValue()
    // being called by parent component
    // before selectOption()
    Promise.resolve(null).then(() => {
      this.dropdownInstance.select(selectedOption);
      this.writeValue(this.getValue());
    });
  }

  // Neccessary implementations to notify Angular of model changes
  onChange(value: string) {
    if (value !== '') {
      this.error = true;
    }
    this.error = false;
  }
  onTouched = () => { };

  // NEEDS TO BE IMPLEMENTED
  // tells Angular there is a change to the model
  registerOnChange(fn: (value: string) => void): void {
    this.onChange = fn;
  }

  // NEEDS TO BE IMPLEMENTED
  registerOnTouched(fn: () => void): void {
    this.onTouched = fn;
  }


  // NEEDS TO BE IMPLEMENTED
  // Updates the model value
  writeValue(value: string): void {

    this.onChange(value);
  }

  validate(c: AbstractControl) {
    // if(this.getValue() != '') {
    //   this.valid = true;
    // }
    // return this.valid;
    return (!this.error) ? null : {
      jsonParseError: {
        valid: false,
      },
    };
  }
}
