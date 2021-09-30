// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Directive, Input } from '@angular/core';
import { AsyncValidator, AbstractControl, ValidationErrors, NG_ASYNC_VALIDATORS } from '@angular/forms';
import { Observable } from 'rxjs';

/**
 * Validates an input against an original value for confirmation.
 * Can be used for passwords, usernames, etc.
 *
 * @export
 * @class ConfirmInputValidator
 * @implements {AsyncValidator}
 */
@Directive({
    selector: '[appConfirmInputValid]',
    providers: [{
        provide: NG_ASYNC_VALIDATORS,
        useExisting: ConfirmInputValidator,
        multi: true
    }]
})
export class ConfirmInputValidator implements AsyncValidator {

    @Input('initInput') initInput: string;

    private timeout;

    validate(
        ctrl: AbstractControl
    ): Promise<ValidationErrors | null> | Observable<ValidationErrors | null> {

        // prevents submission multiple async requests
        clearTimeout(this.timeout);

        return new Promise((resolve) => {

            // delay validation to allow input to finish
            this.timeout = setTimeout(() => {
                const returnVal = (this.initInput !== ctrl.value) ? { 'unconfirmed': true } : null;
                resolve(returnVal);
            }, 400);
        });
    }

}
