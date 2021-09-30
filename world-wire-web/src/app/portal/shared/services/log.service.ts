// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable } from '@angular/core';
import { Service, StatusByDate } from '../../../shared/models/log.interface';
import { SessionService } from '../../../shared/services/session.service';
import { VERSION_DETAILS, IWWApi } from '../../../shared/constants/versions.constant';
import { startCase, toLower, find, toArray } from 'lodash';

@Injectable()
export class LogService {

    // stores current date of today
    today: Date;

    dateOptions = { year: 'numeric', month: 'long', day: 'numeric' };

    // number of days to grab most recent error history for
    // can be changed later if we want to grab more
    public numDays = 30;

    // stores all world wire-hosted services for selected API
    public services: Service[];

    constructor(
        public sessionService: SessionService
    ) { }


    /**
     * Get API version associated with
     * the currently viewed node
     */
    getApiVersionForNode(): string {

        const version = this.sessionService.currentNode ? this.sessionService.currentNode.version : '';

        return version;
    }

    /**
     * Get API Configuration/Services by version
     * @param version
     */
    getApiWithVersion(version: string): IWWApi {
        return find(VERSION_DETAILS, (details: IWWApi) => {
            return details.releaseTag === version;
        });
    }

    getServicesForApi(currentVersion: string) {

        // resetting services every time the API changes
        this.services = null;

        const apis = this.getApiWithVersion(currentVersion) ? toArray(this.getApiWithVersion(currentVersion).config) : null;

        if (apis) {

            // initializing services array
            this.services = [];

            // get statuses for shared services
            this.getService('participant-registry');

            // get statues of services
            // for current version of the node's API
            for (let i = 0; i < apis.length; i++) {

                const api = apis[i];

                // skip if not world wire-hosted service
                if (api.wwHosted === false) {
                    continue;
                }


                // TODO: check if participant is anchor. This is pseudocode for now.
                // if (this.getParticipantRole() !== 'anchor') {
                // // skip getting service details if not anchor
                //   continue;
                // }

                // get service + errors associated with it
                this.getService(api.name);
            }
        }
    }

    /**
     * Gets API Service by name
     * @param serviceName
     */
    getService(serviceName: string) {

        // strip out hypens and convert to
        // start case (first letter capitalized)
        // for human-readable display name
        let displayName = serviceName.replace('/-/g', ' ');
        displayName = startCase(toLower(displayName));

        // convert Api to API
        displayName = displayName.replace('Api', 'API');

        const service = {
            name: serviceName,
            displayName: displayName,
            status: 0,
            url: this.getUrlForService(serviceName),
            errorHistory: this.getRecentHistory()
        };

        // get errors for service
        service.errorHistory = this.getErrorsForService(service);

        this.services.push(service);
    }

    getUrlForService(serviceName: string): string {
        const homeDomain = '';
        const sharedDomain = '';
        let url = '';

        switch (serviceName) {
            case 'participant-registry':
                url = sharedDomain + '/pr';
                break;
        }

        return url;
    }

    getErrorsForService(service: Service) {

        // check for error history initialization
        if (!service.errorHistory) {
            service.errorHistory = [];
        }

        // TODO: query last 30 days of logs
        // for service
        const errorLogs = [];
        // const errorLogs = [{
        //   code: 'WW-001',
        //   details: 'either source amount or beneficiary amount is required:',
        //   message: 'either source amount or beneficiary amount is required',
        //   participant_id: 'nz.one.payments.gftn.io',
        //   time_stamp: '2018-12-12T12:35:23+00:00',
        //   type: 'NotifyWWError',
        //   url: '/v1/client/fees?source_asset=NZD&target_asset=FJD&price=1.000&beneficiary_domain=nz.one.worldwire.io'
        // },
        // {
        //   code: 'WW-001',
        //   details: 'Generic Error Message: Cannot fetch correct env variables for GenericGetAccount function',
        //   message: 'Generic Error',
        //   participant_id: 'nz.one.payments.gftn.io',
        //   time_stamp: '2018-12-12T01:35:23+00:00',
        //   type: 'NotifyWWError',
        //   url: '/v1/crypto/internal/sign'
        // },
        // {
        //   code: 'WW-001',
        //   details: 'Generic Error Message: Cannot fetch correct env variables for GenericGetAccount function',
        //   message: 'G 01eneric Error',
        //   participant_id: 'nz.one.payments.gftn.io',
        //   time_stamp: '2018-12-12T15:35:23+00:00',
        //   type: 'NotifyWWError',
        //   url: '/v1/crypto/internal/sign'
        // }];

        // append logs to date
        for (const log of errorLogs) {

            const date: string = new Date(log.time_stamp).toLocaleDateString(window.navigator.language, this.dateOptions);

            const errorHistory: StatusByDate = find(service.errorHistory, (statusDate: StatusByDate) => {
                return statusDate.date === date;
            });

            // append to error history for date
            errorHistory.errors.push(log);

            // sort by timestamp from most to least recent
            errorHistory.errors.sort((a, b) => {
                const date1 = +new Date(a.time_stamp);
                const date2 = +new Date(b.time_stamp);

                return date1 - date2;
            });

            // order from most to least recent
            errorHistory.errors.reverse();
        }

        for (const dateStatus of service.errorHistory) {
            if (dateStatus.errors.length > 0 && dateStatus.errors.length < 3) {
                dateStatus.status = 1;
            }
            if (dateStatus.errors.length >= 3) {
                dateStatus.status = 2;
            }
        }

        return service.errorHistory;
    }

    /**
     * Initialize the most recent date
     * history to gather system statuses for
     */
    getRecentHistory(): StatusByDate[] {

        const errorHistory: StatusByDate[] = [];

        // gets the last number of days to grab errors
        for (let i = 0; i < this.numDays; i++) {

            const newDate = new Date().setDate(this.today.getDate() - i);

            const newDateString: string = new Date(newDate).toLocaleDateString(window.navigator.language, this.dateOptions);

            const statusByDate: StatusByDate = {
                status: 0,
                date: newDateString,
                errors: []
            };

            errorHistory.push(statusByDate);
        }

        // reorder from latest to earliest
        return errorHistory.reverse();
    }

}
