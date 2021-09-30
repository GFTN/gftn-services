// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit } from '@angular/core';
import { SessionService } from '../../../../shared/services/session.service';

@Component({
  selector: 'app-nodeconfig',
  templateUrl: './nodeconfig.component.html',
  styleUrls: ['./nodeconfig.component.scss']
})
export class NodeConfigComponent implements OnInit {

  // Information (labels, ids)
  // for Configuration Steps
  public configSteps: any[];

  constructor(
    public sessionService: SessionService
  ) { }

  ngOnInit() {
    // configuration information
    this.configSteps = [
      // {
      //   text: `Set Region`,
      //   state: `current`,
      //   title: `Set Container Region for Node`,
      //   description: `Select an EC2 region in which to spin up a new participant node.`
      // },
      {
        text: `Set Participant Details`,
        state: `incomplete`,
        title: `Set Participant Details`,
        description: `Set a unique ID and other details about your participant for your participant's node.`
      },
      {
        text: `Configure Variables`,
        state: `incomplete`,
        title: `Configure Environment Variables`
      },
      {
        text: `Review`,
        state: `incomplete`,
        title: `Review Node Configuration`
      }
    ];
  }

}
