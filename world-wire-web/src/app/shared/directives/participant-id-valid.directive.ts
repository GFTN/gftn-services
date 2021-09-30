// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Directive, Input } from '@angular/core';
import { AsyncValidator, AbstractControl, ValidationErrors, NG_ASYNC_VALIDATORS } from '@angular/forms';
import { ParticipantService } from '../services/participant.service';
import { Observable } from 'rxjs';

/**
 * Checks against firebase if a specific participant stack id already exists
 *
 * @export
 * @class ParticipantIdValidator
 * @implements {AsyncValidator}
 */
@Directive({
    selector: '[appParticipantIdValid]',
    providers: [{
        provide: NG_ASYNC_VALIDATORS,
        useExisting: ParticipantIdValidator,
        multi: true
    }]
})
export class ParticipantIdValidator implements AsyncValidator {

    @Input('participantIdVal') participantIdVal: string;

    private timeout;

    constructor(private participantService: ParticipantService) { }

    validate(
        ctrl: AbstractControl
    ): Promise<ValidationErrors | null> | Observable<ValidationErrors | null> {

        // prevents submission multiple async requests
        clearTimeout(this.timeout);

        return new Promise((resolve) => {

            // delay validation to allow request to come back
            this.timeout = setTimeout(() => {
                this.participantService.participantIdExists(this.participantIdVal).then(
                    (idExists: boolean) => {
                        const returnVal = idExists ? { 'participantIdExists': true } : null;
                        resolve(returnVal);
                    });
            }, 400);
        });
    }

}
