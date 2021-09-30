// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Directive, Input } from '@angular/core';
import { AsyncValidator, AbstractControl, ValidationErrors, NG_ASYNC_VALIDATORS } from '@angular/forms';
import { Observable } from 'rxjs';
import { AngularFireDatabase } from '@angular/fire/database';

/**
 * Checks against firebase if Institution slug already exists
 *
 * @export
 * @class InstitutionIdValidator
 * @implements {AsyncValidator}
 */
@Directive({
    selector: '[appInstitutionIdValid]',
    providers: [{
        provide: NG_ASYNC_VALIDATORS,
        useExisting: InstitutionIdValidator,
        multi: true
    }]
})
export class InstitutionIdValidator implements AsyncValidator {

    @Input('institutionSlug') institutionSlug: string;

    @Input('institutionId') institutionId: string;

    private timeout;

    constructor(
        private db: AngularFireDatabase
    ) { }

    validate(
        ctrl: AbstractControl
    ): Promise<ValidationErrors | null> | Observable<ValidationErrors | null> {

        // prevents submission multiple async requests
        clearTimeout(this.timeout);

        return new Promise((resolve) => {

            // delay validation to allow request to come back
            this.timeout = setTimeout(() => {
                this.db.database.ref('slugs')
                    .child(this.institutionSlug)
                    .once('value', (exists: firebase.database.DataSnapshot) => {

                        // slug exists and the mapped institution is not the one being created
                        const returnVal = (exists.val() && exists.val() !== this.institutionId) ? { 'slugExists': true } : null;

                        resolve(returnVal);
                    });
            }, 400);
        });
    }

}
