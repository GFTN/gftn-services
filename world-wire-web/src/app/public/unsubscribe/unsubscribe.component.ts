// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// The Transaction component is used for sending single
// api transactions for which the response is a simple message.
// ie: unsubscribing from emails

import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ActivatedRoute } from '@angular/router';
import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-transaction',
  templateUrl: './unsubscribe.component.html',
  styleUrls: ['./unsubscribe.component.scss']
})
export class UnsubscribeComponent implements OnInit {

  msg: string;

  constructor(
    private http: HttpClient,
    private route: ActivatedRoute
  ) {  }

  ngOnInit() {
    this.msg = 'Please wait while we process your request.';
    this.sendRequest();
  }

  sendRequest() {

    // get data from route params
    const data = this.route.snapshot.params as { emailHash: string, mailingList: any, all: any };

    this.http.post(
      environment.apiRootUrl + '/unsubscribe',
      {
        emailHash: data.emailHash,
        // converts csv of mailing lists names into an array of mailing list names
        mailingList: data.mailingList.split(','),
        // if 0 then false, if 1 true
        all: Boolean(Number(data.all))
      }
    ).subscribe((res) => {
      console.log(res);
      this.msg = 'We have unsubscribed you from our mailing list.';
    });


  }

}


