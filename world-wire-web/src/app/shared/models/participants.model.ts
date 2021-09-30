// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { IInstitution } from './participant.interface';
import { Injectable, NgZone } from '@angular/core';
import { AngularFireDatabase, DatabaseSnapshot } from '@angular/fire/database';
import { Observable, Observer } from 'rxjs';
import { FI_TYPES, ROLES, ENVIRONMENT, STATUS } from '../constants/general.constants';
import * as _ from 'lodash';
import * as firebase from 'firebase';

// export interface IInstitutionAll {
//     // array of institutions
//     arr: IInstitution[];
//     // used to lookup by institution by slug (both should be unique)
//     ids: { [slug: string]: /*institutionId*/ string }[];
// }

@Injectable()
export class ParticipantsModel {

    /**
     * Primarily used to get details about a participant
     * including a list of associated users with permissions.
     * Security: NOT a secure node for evaluating permissions
     * @memberof ParticipantModel
     */
    route = 'participants/{institutionId}'; // firebase ref

    model: IInstitution;

    types = FI_TYPES;
    roles = ROLES;
    status = STATUS;
    searchedParticipants: IInstitution[];

    private allRef: firebase.database.Reference;

    private db: firebase.database.Database;

    constructor(
        private afDb?: AngularFireDatabase,
        // Note: Do not use ngZone here (only use in component view) so that this class can be easily tested.
        // private ngZone: NgZone
    ) {

        // set db ref
        if (this.afDb) {
            this.setDb(afDb.database);
        }

        this.model = {
            info: {
                institutionId: '',
                name: '',
                // short: '',
                geo_lat: 0,
                geo_lon: 0,
                country: '',
                address1: '',
                address2: '',
                city: '',
                state: '',
                zip: '',
                logo_url: '',
                site_url: '',
                kind: 'Money Transfer Operator',
                slug: '',
                status: 'active'
            },
            users: {},
        };

    }

    /**
     * Sets db ref for CRUD operations pertaining
     * to a participant. This method
     * allows for unit testing with a mock db.
     * @param dbRef
     */
    setDb(dbRef: firebase.database.Database) {
        this.db = dbRef;
        this.allRef = this.db.ref('participants');
    }

    /**
     * Sets the slug to a kebab case formated string
     * @param text
     * @memberof ParticipantsModel
     */
    genSlug(text: string) {
        // with type info
        this.model.info.slug = _.kebabCase(text);
    }

    /**
    * Returns all participants listed in firebase database
    * (not participant registry). Firebase database is
    * kept separately in sync with PR so that if anything
    * malicious were to be done to the firebase database
    * of participants it would not affect the WW network
    *
    * @returns {Promise<{ [id: string]: IInstitution }>}
    * @memberof ParticipantsModel
    */
    get(institutionId: string): Promise<IInstitution> {


        return new Promise((resolve, reject) => {

            // watch changes to values affecting firebase ref
            this.allRef.child(institutionId)
                .once('value', (data: DatabaseSnapshot<IInstitution>) => {

                    resolve(data.val());

                }, (error: any) => {
                    console.log('Unable to get participant: ', error);
                    reject('Unable to get participant: ' + error);
                });

        });

    }

    /**
     * Returns all participants listed in firebase database
     * (not participant registry). Firebase database is
     * kept separately in sync with PR so that if anything
     * malicious were to be done to the firebase database
     * of participants it would not affect the WW network
     *
     * @returns {Promise<{ [id: string]: IInstitution }>}
     * @memberof ParticipantsModel
     */
    allPromise(): Promise<IInstitution[]> {

        return new Promise((resolve, reject) => {

            // watch changes to values affecting firebase ref
            this.allRef.once('value', (data: DatabaseSnapshot<{ [id: string]: IInstitution }>) => {

                // sort by 'name' field
                const participants = _.sortBy(
                    // return in array format
                    _.toArray(data.val()),
                    (o: IInstitution) => {
                        // sort case insensitive
                        return o.info.name.toLowerCase();
                    }
                );

                // //  map slug to institution ids for easy lookup of a institution id
                // const mapIds = _.map(data.val(), (institution: IInstitution, institutionId: string) => {
                //     return { [institution.slug]: institutionId };
                // });

                // // create result obj
                // const res: IInstitutionAll = {
                //     arr: participants,
                //     ids: mapIds
                // };

                // resolve(res);

                resolve(participants);

            }, (error: any) => {
                console.log('Unable to get all participants: ', error);
                reject('Unable to get all participants: ' + error);
            });

        });

    }

    /**
     * Returns all participants listed in firebase database
     * (not participant registry). Firebase database is
     * kept separately in sync with PR so that if anything
     * malicious were to be done to the firebase database
     * of participants it would not affect the WW network
     *
     * @returns {Observable<{ [id: string]: IInstitution }>}
     * @memberof ParticipantsModel
     */
    allObservable(): Observable<IInstitution[]> {

        const source = new Observable((observer: Observer<IInstitution[]>) => {

            // watch changes to values affecting firebase ref
            this.allRef.on('value',
                (participants: DatabaseSnapshot<{ [id: string]: IInstitution }>) => {

                    // sort by 'name' field
                    const data = _.sortBy(
                        // return in array format
                        _.toArray(participants.val()),
                        (o: IInstitution) => {
                            // sort case insensitive
                            return o.info.name.toLowerCase();
                        }
                    );

                    // since calling from external source
                    // need to put result into angular zone
                    // this.ngZone.run(() => {
                    // update observer value
                    observer.next(
                        data
                    );

                    // });

                }, (error: any) => {
                    console.log(error);

                }
            );

        });

        return source;

    }

    /**
     * This detaches a callback so that callback listeners
     * are appropriately removed when no longer used.
     * This should be used with 'ngOnDestroy' since leaving a
     * angular component usually means that the data associated
     * with this callback is no longer needed
     *
     * @memberof ParticipantsModel
     */
    unsubscribeObservable() {

        // console.log('off');

        // stop listener on
        // https://firebase.google.com/docs/reference/js/firebase.database.Reference#off
        this.allRef.off();
    }

    /**
    * Creates a new participant
    *
    * @memberof ParticipantModel
    */
    create(): Promise<boolean> {

        return new Promise((resolve, reject) => {

            // create new key
            this.allRef.push()
                .then((data) => {
                    // console.log(data.key);
                    // set new record id
                    this.model.info.institutionId = data.key;
                    // call update to add new record to firebase
                    this.update().then(() => {
                        resolve(true);
                    }, () => {
                        reject(false);
                    });
                }, () => {
                    reject(false);
                });
        });

    }

    /**
     * Updates an existing participant's info
     *
     * @returns
     * @memberof ParticipantModel
     */
    update(): Promise<boolean> {

        console.log(this.model);

        // return promise since result is a single success or failure
        return new Promise((resolve, reject) => {

            if (this.model.info.institutionId) {

                // update record in firebase
                this.allRef.child(this.model.info.institutionId).child('info')
                    .update(this.model.info).then(() => {
                        resolve(true);
                    }, (err: any) => {
                        console.log('Error: Unable to save participant details.' + err);
                        reject(false);
                    });

            } else {
                console.log('No institutionId provided.');
                reject(false);
            }

        });

    }

    delete(institutionId: string) {
        return new Promise((resolve, reject) => {

            if (institutionId) {

                // update record in firebase
                this.allRef.child(this.model.info.institutionId)
                    .remove()
                    .then(() => {
                        resolve(true);
                    }, () => {
                        console.log('Error: Unable to delete participant details.');
                        reject(false);
                    });

            } else {
                console.log('No institutionId provided.');
                reject(false);
            }

        });
    }

}
