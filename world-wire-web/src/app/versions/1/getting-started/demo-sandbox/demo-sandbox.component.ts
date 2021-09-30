// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Input, OnChanges, ViewChild, ElementRef } from '@angular/core';
import { HttpClient, HttpRequest, HttpHeaders } from '@angular/common/http';
import * as highlightjs from 'highlight.js';

@Component({
  selector: 'app-demo-sandbox',
  templateUrl: './demo-sandbox.component.html',
  styleUrls: ['./demo-sandbox.component.scss']
})
export class DemoSandboxComponent implements OnInit, OnChanges {
  @Input() requestRoute: string;
  @Input() requestBody: string;
  @Input() operationType: string;
  @Input() requestLines: number;
  @Input() responseLines: number;

  // TODO: get request response syntax highlighting to work
  // @ViewChild('requestResultSelector') requestResultSelector: ElementRef;

  bodyString = '';
  requestResult: string;
  constructor(private http: HttpClient) { }

  ngOnInit() {
    // this.requestBody = JSON.stringify(this.requestBody);
    this.requestResult = '';
    this.bodyString = JSON.stringify(this.requestBody, undefined, 2);
    // initialize syntax highlighting onload of page
    this.bodyString = highlightjs.highlight('json', this.bodyString).value;
  }

  ngOnChanges() {
    // this.requestBody = JSON.stringify(this.requestBody);
  }

  countLines(input: string): number {
    return JSON.stringify(input, undefined, 2).split(':').length + 4;
  }

  executeRequest() {
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type':  'application/json'
      })
    };

    switch (this.operationType) {
      case 'GET': {
        this.http.get(
          this.requestRoute
        ).subscribe(
          (data: any) => {
            // console.log('Performing GET request on ' +
            //         this.requestScheme + '://' + this.subDomain + '.' + this.baseDomain + '.' + this.domainExtension +
            //         this.apiPort + this.basePath + this.requestUrl + this.baseQuery);
            // console.log('Response data:\n', JSON.stringify(data) );
            this.requestResult = JSON.stringify(data, undefined, 2);
            // highlight syntax for JSON result
            this.requestResult = highlightjs.highlight('json', this.requestResult).value;
          }, (err: any) => {
            console.log(err);
            this.requestResult = JSON.stringify(err, undefined, 2);
          }
        );
        break;
      }
      case 'POST': {
        this.http.post(
          this.requestRoute, this.requestBody, httpOptions
          // 'http://localhost:9080' + this.basePath + this.requestUrl + this.baseQuery,
          // 'http://localhost:9080/v1/helloworldwire',
          // 'https://jsonplaceholder.typicode.com/posts',
          // this.requestPayload, httpOptions
          // JSON.parse('{"account_name":"default","asset_code":"NZD","amount":"300"}'), httpOptions
        )
          .subscribe(
          (data: any) => {
            // console.log('Performing ' + this.operationType + ' request on ' +
            // this.requestScheme + '://' + this.subDomain + '.' + this.baseDomain + '.' + this.domainExtension +
            // this.apiPort + this.basePath + this.requestUrl + this.baseQuery,
            // '\n with headers: \n' + JSON.stringify(httpOptions.headers.keys()),
            // '\n and payload: ' + this.requestPayload);
            console.log(data);
            this.requestResult = JSON.stringify(data, undefined, 2);
            // highlight syntax for JSON result
            this.requestResult = highlightjs.highlight('json', this.requestResult).value;
          }, (err: any) => {
            console.log(err);
            this.requestResult = JSON.stringify(err, undefined, 2);
          }
        );
        break;
      }
    } // end switch

    // TODO: get request response syntax highlighting to work
    // highlightjs.initHighlightingOnLoad();
    // highlightjs.highlightBlock();
  }

}
