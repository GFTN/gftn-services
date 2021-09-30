// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { SwaggerService } from '../../shared/services/swagger.service';
import { Tooltip } from 'carbon-components';

@Component({
  selector: 'app-getting-started',
  templateUrl: './getting-started.component.html',
  styleUrls: ['./getting-started.component.scss']
})
export class GettingStartedComponent implements OnInit {

  paths: any;

  requestResult: string;

  constructor(
    private http: HttpClient,
    public swaggerService: SwaggerService
  ) { }

  ngOnInit() {

    this.requestResult = '';

    this.createPaths();

  }

  createPaths() {
    this.paths = {
      1: {
        'title': 'Select your Sandbox role',
      },
      2: {
        'title': 'Check if any Participants are already on your Whitelist',
        'url': 'https://one-nz.worldwire.io/api/v1/participant-client-api/participants/whitelist',
        'method': 'GET',
        'body': {},
        'link': '../api/client-api'
      },
      3: {
        'title': 'Find eligible Participants to exchange with',
        'url': 'https://one-nz.worldwire.io/api/v1/participant-client-api/participants?' +
          'country_code=FJ&issuer_id=is-tn-one-nz',
        'method': 'GET',
        'body': {},
        'link': '../api/client-api'
      },
      4: {
        'title': 'Whitelist the Participant',
        'url': 'https://one-nz.worldwire.io/api/v1/participant-client-api/participants/whitelist',
        'method': 'POST',
        'body': {
          'participant_id': 'is.tn.onenz.worldwire.io'
        },
        'link': '../api/client-api'
      },
      5: {
        'title': 'Setup trusted assets',
        'url': 'https://one-nz.worldwire.io/api/v1/participant-client-api/trust?asset_code=FJDDA&permissions=allow',
        'method': 'POST',
        'body': {
          'participant_id': 'pi-tn-one.fiji.worldwire.io',
          'limit': '3000',
          'account_name': 'default'
        },
        'link': '../api/client-api'
      },
      6: {
        'title': 'Submit a request for quotes',
        'url': 'https://one-nz.worldwire.io/api/v1/participant-client-api/quotes/requests',
        'method': 'POST',
        'body': {
          'limit_max': '2000',
          'limit_min': '1',
          'ofi_id': 'is.tn.onenz.worldwire.io',
          'source_asset': {
            'asset_code': 'NZDDA',
            'asset_type': 'DA',
            'issuer_id': 'is.tn.onenz.worldwire.io'
          },
          'target_asset': {
            'asset_code': 'FJDDA',
            'asset_type': 'DA',
            'issuer_id': 'pn.tn.onefj.worldwire.io'
          },
          'time_expire': '1644516034'
        },
        'link': '../api/client-api'
      },
      7: {
        'title': 'Retrieve quotes',
        'url': 'https://one-nz.worldwire.io/api/v1/participant-client-api/quotes/requests/{request_id}?' +
          'request_id=2377229340202',
        'method': 'GET',
        'body': {},
        'link': '../api/client-api'
      },
      8: {
        'title': 'Act on a desired quote',
        'url': 'https://one-nz.worldwire.io/api/v1/participant-client-api/exchange',
        'method': 'POST',
        'body': {
          'exchange': 'LSfkljahSJLHSJKE4',
          'signature': 'asjdflkjasdlkfjadsfj'
        },
        'link': '../api/client-api'
      },
      9: {
        'title': 'Calculate fees to send payment',
        'url': 'https://one-nz.worldwire.io/api/v1/participant-client-api/fees/participants/{participant_id}?' +
          'participant_id=pn.tn.oneFJ.worldwire.io',
        'method': 'GET',
        'body': {
          'participant_id': 'pn.tn.oneNZ.worldwire.io',
          'fee_id': '2342LASF',
          'settlement_amount': '4000.0000001',
          'settlement_asset': 'NZDDA',
          'payout_point': {
            'context': 'http://schema.org',
            'type': 'FinancialInstitute',
            'name': 'CIMB Bank',
            'currencies_accepted': ['USD'],
            'additional_type': 'eligibleTransactionVolume',
            'image': 'https://www.cimbbank.com.sg/content/dam/cimbsingapore/logo/cimblogo.jpg',
            'url': '',
            'telephone': '07-418 6258 / 6276',
            'member_of': ['BankA'],
            'address': {
              'type': 'PostalAddress',
              'street_address': '39A s Rahmat',
              'address_locality': 'changi',
              'address_region': 'North East',
              'postal_code': '83000',
              'address_country': 'SG'
            },
            'geo': {
              'type': 'GeoCoordinates',
              'latitude': '1.8482097',
              'longitude': '102.93248110000002'
            },
            'offer': {
              'type': 'OfferCatalog',
              'name': 'Cash cash_pickup',
              'category': [
                {
                  'type': 'OfferCatalog',
                  'name': 'Cash cash_pickup',
                  'list': [
                    {
                      'type': 'Offer',
                      'detail': {
                        'type': 'delivery',
                        'name': 'Deliver PHP to home address',
                        'terms_of_service': 'Limit of 52,000 PHP'
                      }
                    },
                    {
                      'type': 'Offer',
                      'detail': {
                        'type': 'delivery',
                        'name': 'Deliver USD to home address',
                        'terms_of_service': 'Limit of 1,000 USD'
                      }
                    }
                  ]
                },
                {
                  'type': 'OfferCatalog',
                  'name': 'Cash pick-up',
                  'list': [
                    {
                      'type': 'Offer',
                      'detail': {
                        'type': 'agency_pickup',
                        'name': 'Pick-up PHP from kiosk',
                        'terms_of_service': 'Limit of 52,000 PHP per day per identity verified'
                      }
                    },
                    {
                      'type': 'Offer',
                      'detail': {
                        'type': 'agency_pickup',
                        'name': 'Pick-up USD from kiosk',
                        'terms_of_service': 'Limit of 1,000 USD per day per identity verified'
                      }
                    }
                  ]
                }
              ]
            },
            'opening_hours': [{
              'type': 'OpeningHoursSpecification',
              'day_of_week': [
                'Monday',
                'Tuesday',
                'Wednesday',
                'Thursday'
              ],
              'opens': '09:15',
              'closes': '16:30'
            },
            {
              'type': 'OpeningHoursSpecification',
              'day_of_week': ['Friday'],
              'opens': '09:15',
              'closes': '16:00'
            }]
          },
          'payout_amount': '4000.0000001',
          'payout_asset': 'NZDDA'
        },
        'link': '../api/client-api'
      },
      10: {
        'title': 'Send a payment',
        'url': 'https://one-nz-client.worldwire.io/api/v1/participant-client-api/transactions/send',
        'method': 'POST',
        'body': { 'end_to_end_id': 'ABC123KLM2917' },
        'link': '../api/client-api'
      }
    };
  }

  countLines(input: string): number {
    return JSON.stringify(input, undefined, 2).split(':').length + 4;
  }

  executeRequest(url: string, method: string, body: string) {
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json'
      })
    };

    switch (method) {
      case 'GET': {
        this.http.request(method,
          url
        ).subscribe(
          (data: any) => {
            this.requestResult = JSON.stringify(data, undefined, 2);
          }, (err: any) => {
            console.log(err);
            this.requestResult = JSON.stringify(err, undefined, 2);
          }
        );
        break;
      }
      case 'POST': {
        this.http.request(method,
          url, { body: body }
        )
          .subscribe(
            (data: any) => {
              console.log(data);
              this.requestResult = JSON.stringify(data, undefined, 2);
            }, (err: any) => {
              console.log(err);
              this.requestResult = JSON.stringify(err, undefined, 2);
            }
          );
        break;
      }

    }
  }
}
