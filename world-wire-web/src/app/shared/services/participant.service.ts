// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable } from '@angular/core';
import { HttpHeaders, HttpClient } from '@angular/common/http';
import { ENVIRONMENT } from '../constants/general.constants';

import { Participant } from '../models/participant.interface';

import { SessionService } from './session.service';
import { AuthService } from './auth.service';


@Injectable()
export class ParticipantService {

    constructor(
        private http: HttpClient,
        private authService: AuthService,
        private sessionService: SessionService
    ) { }

    /**
     * Returns true if the Participant ID already exists in the database
     * TODO: Need to sync this officially with the PR of each env for nodes made outside of the portal
     *
     * @param {string} id
     * @param {string} env
     * @returns {Promise<boolean>}
     * @memberof ParticipantService
     */
    participantIdExists(id: string): Promise<boolean> {

        return new Promise(async (resolve) => {

            // searches through ALL nodes, active and 'deleted'
            const iid = await this.getParticipantInstitution(id);

            return iid ? resolve(true) : resolve(false);
        });
    }

    /**
     * Get Institution associated with this participant id
     * (ID of Participant node)
     *
     * @param {string} participantId
     * @param {string} [env]
     * @returns {Promise<string>}
     * @memberof ParticipantService
     */
    async getParticipantInstitution(participantId: string): Promise<string> {

        // const pid = await this.db.database.ref(`/nodes/${participantId}`).once('value');

        // only make request if value is provided,
        // since request requires a parameter
        if (participantId) {
            const pRequest = `https://admin.${ENVIRONMENT.envGlobalRoot}/pr/v1/admin/pr`;

            const h: HttpHeaders = await this.authService.getFirebaseIdToken(this.sessionService.institution.info.institutionId);

            const options = {
                headers: h
            };

            let participantFound: Participant;

            // check against official PR from /admin/pr
            try {
                const participants = await this.http.get(
                    pRequest,
                    options
                ).toPromise() as Participant[];

                participantFound = participants.find((participant: Participant) => {
                    return participant.id === participantId;
                });
            } catch (err) {
                return null;
            }

            return participantFound ? participantFound.id : null;
        }

        return null;
    }
}
