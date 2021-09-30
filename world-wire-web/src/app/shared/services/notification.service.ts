// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable } from '@angular/core';

@Injectable()
export class NotificationService {

    // types of notifications:
    // success: boolean;
    // error: boolean;
    // info: boolean;
    // warning: boolean;

    // shows the message
    msg: string;

    // shows the notification
    _show: boolean;

    timer: any;

    title: string;

    locationClass: string;

    kindClass: string;

    constructor(
    ) {
        // setTimeout(() => {
        //     this._show = true;
        // }, 1000);
        this._show = false;
        this.locationClass = 'notify-bottom';
        this.kindClass = 'bx--toast-notification--info';
        this.timer = null;
    }

    show(
        kind: 'success' | 'error' | 'info' | 'warning',
        msg: string,
        duration?: number,
        title?: string,
        location?: 'top' | 'bottom',
    ) {

        this._show = true;
        this.msg = msg;

        if (location) {
            this.locationClass = 'notify-' + location;
        } else {
            this.locationClass = 'notify-bottom';
        }

        if (this.timer) {
            clearTimeout(this.timer);
        }

        // set default duration
        if (!duration) {
            duration = 5000;
        }

        // set default title
        if (title) {
            this.title = title;
        } else {
            this.title = kind;
        }

        // make notification visible
        this.kindClass = 'bx--toast-notification--' + kind;

        // hide notification after duration
        this.timer = setTimeout(() => {
            // this[kind] = false;
            this._show = false;
        }, duration);

    }

}
