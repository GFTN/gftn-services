// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Directive } from '@angular/core';
import { AsyncValidator, NG_ASYNC_VALIDATORS, AbstractControl, ValidationErrors } from '@angular/forms';
import * as cc from 'currency-codes';
import { Observable } from 'rxjs';

@Directive({
  selector: '[appValidAsset]',
  providers: [
    {
      provide: NG_ASYNC_VALIDATORS,
      useExisting: AssetValidator,
      multi: true
    }
  ]
})
export class AssetValidator implements AsyncValidator {

  private timeout;

  invalidText: ValidationErrors = { 'invalidAsset': true };

  validate(
    control: AbstractControl
  ): Promise<ValidationErrors | null> | Observable<ValidationErrors | null> {

    // prevents submission multiple async requests
    clearTimeout(this.timeout);

    return new Promise((resolve) => {

      // delay validation to allow input to finish
      this.timeout = setTimeout(() => {

        if (!control.value) {
          resolve(null);
        }

        // asset code should be exactly 3 chars.
        if (control.value.length < 3) {
          resolve(this.invalidText);
        }

        const code = control.value.substring(0, 3);

        const returnVal = !cc.code(code) ? this.invalidText : null;

        resolve(returnVal);
      }, 400);
    });
  }
}
