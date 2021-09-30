// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable } from '@angular/core';
import { DocumentService } from '../services/document.service';
import * as moment from 'moment-timezone';

@Injectable()
export class UtilsService {

    timezoneAbbr: string;

    constructor(
        private documentService: DocumentService
    ) {
        // get timezone for user
        const timezone = moment.tz.guess();

        this.timezoneAbbr = moment.tz(timezone).zoneAbbr();
    }

    /**
     * Helper function to convert string number
     * to decimal/float number in the view template
     *
     * @param {string} text
     * @returns {number}
     * @memberof UtilsService
     */
    strToFloat(text: string): number {
        return parseFloat(text);
    }

    /**
     * Converts JS object to array of object values
     *
     * @param {object} obj
     * @returns
     * @memberof UtilsService
     */
    toArray(obj: object) {
        return Object.values(obj);
    }

    /**
     * Sets a new cookie with name and expiry time
     *
     * @param {string} name
     * @param {string} dateString
     * @memberof UtilsService
     */
    setCookie(name: string, dateString: string, value: string = 'true') {

        this.documentService.documentRef.cookie = name + '=' + value + '; expires=' +
            dateString +
            ';path=/';
    }

    /**
     * Get cookie by cookie's name
     *
     * @param {string} name
     * @returns
     */
    getCookie(name: string) {
        const value = '; ' + this.documentService.documentRef.cookie;
        const parts = value.split('; ' + name + '=');
        if (parts.length === 2) {
            return parts.pop().split(';').shift();
        }
    }

    /**
     * Capitalize the first letter in a string
     *
     * @param {string} text
     * @returns {string}
     * @memberof UtilsService
     */
    capitalizeFirstLetter(text: string): string {
        return text.charAt(0).toUpperCase() + text.slice(1);
    }

    /**
     * Converts basic action verbs to past tense for user text
     *
     * @param {string} verb
     * @returns {string}
     * @memberof UtilsService
     */
    convertVerbToPastTense(verb: string): string {
        const lastLetter = verb.slice(-1);

        return (lastLetter === 'e') ? verb + 'd' : verb + 'ed';
    }

    /**
     * Converts Unix timestamp to date and time in user's local timezone.
     * Also appends the timezone abbreviation for this user.
     *
     * @param {number} timestamp
     * @returns {string}
     * @memberof UtilsService
     */
    toLocaleDateTime(timestamp: number): string {
        return new Date(timestamp * 1000).toLocaleString() + ' ' + this.timezoneAbbr;
    }

    /**
     * Generates a Unix timestamp in seconds.
     * Default for Date.now() is in millisecs
     * so we need to do some Math to round it up
     *
     * @param {number} [timestampInMillisecs]
     * @returns {number}
     * @memberof UtilsService
     */
    getTimestampInSecs(timestampInMillisecs?: number): number {
        timestampInMillisecs = timestampInMillisecs ? timestampInMillisecs : Date.now();
        return Math.floor(timestampInMillisecs / 1000);
    }

    /**
     * Checks if string is a properly formatted ip address
     *
     * @param {string} text
     * @returns
     * @memberof UtilsService
     */
    isIpv4(text: string) {
        if (/^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/.test(text)) {
            return true;
        } else {
            return false;
        }
    }

    /**
     * used by a view to convert a commas separated string
     * into a set of text displayed on newlines
     *
     * @param {string} text
     * @returns
     * @memberof UtilsService
     */
    convertCSVtoNewLine(text: string) {
        return (text.toString()).replace(/,/g, '\n');
    }

}
