// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable, OnInit } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class WwHttpService implements OnInit {
  requestUrl: string;
  requestScheme = 'http';
  subDomain = 'nightly-nz';
  baseDomain = 'worldwire-ibm';
  domainExtension = 'com';
  apiPort = ':8080';
  basePath = '/v1';
  baseQuery = '';
  requestPath = JSON.stringify(this.requestUrl);
  // jsonSpec: Spec;
  requestPayload: JSON = {} as JSON;
  payloadString: string;
  defaultPayloads: JSON;

  constructor(private http: HttpClient) { }

  ngOnInit() {
    this.defaultPayloads = this.executeRequest('GET', './assets/open-api/1/quick-start.json', '{}');
  }

  setRequestPayload(payload: JSON) {
    this.requestPayload = payload;
  }

  getRequestPayload(): JSON {
    return this.requestPayload;
  }

  returnDefaults() {
    return this.defaultPayloads;
  }

  executeRequest(method: string, url: string, body?: string): any {
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
            return JSON.stringify(data, undefined, 2);
          }, (err: any) => {
            console.log(err);
            return JSON.stringify(err, undefined, 2);
          }
        );
        break;
      }
      default: {
        this.http.request(method,
          url, {body: body}
        )
          .subscribe(
            (data: any) => {
              console.log(data);
              return JSON.stringify(data, undefined, 2);
            }, (err: any) => {
              console.log(err);
              return JSON.stringify(err, undefined, 2);
            }
          );
        break;
      }

    }
  }
}
