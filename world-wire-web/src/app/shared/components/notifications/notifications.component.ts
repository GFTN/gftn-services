// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component } from '@angular/core';
import * as _ from 'lodash';
import { NotificationService } from '../../services/notification.service';
import { simpleFadeAnimation } from '../../animations';

@Component({
  selector: 'app-notification',
  templateUrl: './notifications.component.html',
  styleUrls: ['./notifications.component.scss'],
  animations: [
    simpleFadeAnimation
  ]
})
export class NotificationComponent {

  constructor(
    public n: NotificationService
  ) { }

}
