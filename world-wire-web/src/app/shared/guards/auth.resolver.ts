// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable } from '@angular/core';
import { Resolve } from '@angular/router';
import { Observable } from 'rxjs';
import { map, take, tap } from 'rxjs/operators';
import { AngularFireAuth } from '@angular/fire/auth';
import { Component } from '@angular/compiler/src/core';

@Injectable()
export class AuthResolver implements Resolve<Component> {

    constructor(
        private auth: AngularFireAuth
    ) { }

    resolve(): Observable<any> | Promise<any> | any {

        return this.auth.authState.pipe(
            take(1),
            map((user: firebase.User) => user),
            tap(
                (user: firebase.User) => {

                    return new Promise((resolve) => {
                        if (user) {
                            resolve(user);
                        } else {
                            resolve(null);
                        }
                    });
                }
            ));
    }
}
