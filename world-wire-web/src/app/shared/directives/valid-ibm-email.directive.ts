// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Directive, forwardRef, Input } from '@angular/core';
import { NG_VALIDATORS, AbstractControl,  Validator, FormControl, ValidatorFn } from '@angular/forms';


@Directive({
  selector: '[appValidEmail]',
  providers: [
    { provide: NG_VALIDATORS,
      useExisting: EmailValidator,
      multi: true }
  ]
})

export class EmailValidator implements Validator {

  validate(control: AbstractControl): {[key: string]: any} | null {

    let isEmailValid = false;
    const input = control.value;
    // The value in the input box is null
    if (!control.value) {
      return null;
    }

    // the email entered in the input box is not acceptable. It has to be (ibm.com, us.ibm.com, sg.ibm.com)
    // checking to see if the email has an ibm domain and has a length greater than 0.
    if (input.substr(-8) === '@ibm.com' && input.substring(0, input.length - 8).length > 0) {
      isEmailValid = true;
    } else if ( (input.substr(-11) === '@us.ibm.com' && input.substring(0, input.length - 11).length > 0) || (input.substr(-11) === '@sg.ibm.com' && input.substring(0, input.length - 11).length > 0) ) {
      isEmailValid = true;
    } else {
      isEmailValid = false;
    }

    if (!isEmailValid) {
      return { email: 'The email entered is not an IBM email address' };
    }

    return null;

  }
}
