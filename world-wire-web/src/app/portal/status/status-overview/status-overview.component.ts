// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, AfterViewInit, ElementRef, HostBinding } from '@angular/core';
import { NotificationService, Notification, ModalService } from 'carbon-components-angular';
import { NotificationContent } from 'carbon-components-angular/notification/notification-content.interface';
import { slideAnimation } from '../../../shared/animations';

import { StatusDetailsModalComponent } from '../status-details-modal/status-details-modal.component';
import { StatusByDate, Service } from '../../../shared/models/log.interface';
import { MediaObserver, MediaChange } from '@angular/flex-layout';
import { LogService } from '../../shared/services/log.service';
import { SessionService } from '../../../shared/services/session.service';


export class StatusMessage {
  name?: string;
  icon?: string;
  message?: string;
  notification: NotificationContent;
}

@Component({
  selector: 'app-status-overview',
  templateUrl: './status-overview.component.html',
  styleUrls: ['./status-overview.component.scss'],
  animations: [
    slideAnimation
  ],
})

export class StatusOverviewComponent implements OnInit, AfterViewInit {

  // stores current date of today
  today: Date;

  // stores static messages viewed for the individual status by date
  statusMessages: StatusMessage[];

  // stores reference to the current overall Notification status
  currentRef: Notification;

  // variables for individual error details by date.
  // detailsPaneOpened - to toggle open/close of details pane display
  detailsPaneOpened = false;

  // stores reference to the current errors being view by date
  openedLinkRef: ElementRef;

  detailsObject: StatusByDate;

  state = '';

  viewInit = false;

  iconPath = '/assets/icons/ibm/carbon-icons.svg#';

  constructor(
    private notificationService: NotificationService,
    protected modalService: ModalService,
    public sessionService: SessionService,
    public media: MediaObserver,
    public logService: LogService,
  ) {
    this.today = new Date();

    this.statusMessages = [];

    this.setStatusMessages();

    media.media$
      .subscribe((change: MediaChange) => {
        this.state = change ? `'${change.mqAlias}' = (${change.mediaQuery})` : '';
      });
  }

  @HostBinding('attr.class') cls = 'flex-fill';

  ngOnInit() {
    // initializing date onload of page
    this.logService.today = this.logService.today ? this.logService.today : new Date();

    // initialize data only once per page load
    if (!this.logService.services) {

      this.updateServicesOverview();
    }
  }

  ngAfterViewInit() {

    this.viewInit = true;

    // called once after view inits to account for initial services init
    this.getCurrentOverallStatus();

  }

  /**
   * Update view of services overview to reflect the currently selected node
   */
  updateServicesOverview() {

    // get version of current node in session
    const currentVersion = this.logService.getApiVersionForNode();

    // get associated services for the current node's version of the API
    this.logService.getServicesForApi(currentVersion);

    if (this.viewInit) {
      // this can only run after view inits since the service
      // from carbon-components-angular uses a ElementRef.
      this.getCurrentOverallStatus();
    }
  }

  getCurrentOverallStatus() {

    let downServices = 0;

    // check current status of all services
    // to get overall current statuses
    if (this.logService.services) {
      for (const service of this.logService.services) {
        if (this.getOverallStatusByDate(service).name !== 'success') {
          downServices++;
        }
      }

      if (downServices === 0) {
        this.replaceNotification(0);
      }

      if (downServices > 0 && downServices < this.logService.services.length) {
        this.replaceNotification(1);
      }

      if (downServices > 0 && downServices === this.logService.services.length) {
        this.replaceNotification(2);
      }
    } else {
      this.sessionService.currentNode ? this.replaceNotification(3) : this.replaceNotification(4);
    }
  }

  /**
   * Closes existing notification
   * before displaying a new one so that
   * only one is visible at a time
   * @param errorRef
   */
  replaceNotification(errorRef: number) {

    // close existing notification
    // if one exists on the screen
    this.closeNotification().then(() => {
      if (errorRef >= 0) {
        this.showNotification(errorRef);
      }
    });
  }

  /**
   * Calls Carbon notification service
   * to handle showing of notification
   * @param errorRef
   */
  showNotification(errorRef: number) {

    // append current time to notification
    const msg = this.statusMessages[errorRef].notification;
    msg.message = msg.message ? msg.message : 'as of ' + new Date().toLocaleString();

    const ref = this.notificationService.showNotification(msg);
    this.currentRef = ref;

    // subscribe to close button
    // to dereference current notification
    this.currentRef.close.subscribe(() => {
      this.currentRef = null;
    });
  }

  /**
   * Allows for automatic/programmatic
   * closing of notification
   */
  closeNotification(): Promise<void> {
    return new Promise((resolve) => {

      // quick resolve if no existing notification to close
      if (!this.currentRef) {
        return resolve();
      }

      this.currentRef.onClose();
      setTimeout(() => {
        this.currentRef = null;
        return resolve();
      }, 201);

    });
  }

  getOverallStatusByDate(service: Service): StatusMessage {

    // default to 'success' state
    let status = this.statusMessages[0];

    // errors < 3 means service is in partially degraded state
    if (service.errorHistory[0].errors.length > 0 && service.errorHistory[0].errors.length < 3) {
      status = this.statusMessages[1];
    }

    // errors >= 3 means service is has been fully degraded.
    // this can be changed if we want to adjust the error rate
    // or add other factors in which to determine full degradation
    if (service.errorHistory[0].errors.length >= 3) {
      status = this.statusMessages[2];
    }

    return status;
  }

  /**
  * returns message associated with status
  * @param status
  */
  getStatusName(status: number): string {

    return this.statusMessages[status].name;
  }

  /**
  * returns message associated with status
  * @param status
  */
  getStatusMessage(status: number): string {

    return this.statusMessages[status].message;
  }

  /**
   * initialize status messages for current overall status
   */
  setStatusMessages() {
    this.statusMessages.push(
      {
        name: 'success',
        icon: 'icon--checkmark--solid',
        message: 'No downtime',
        notification: {
          type: 'success',
          title: 'All systems are currently online.',
          message: '',
          target: '#notification'
        }
      },
      {
        name: 'warning',
        icon: 'icon--warning--solid',
        message: 'Partial downtime occured',
        notification: {
          type: 'warning',
          title: 'Some systems are currently down or degraded.',
          message: '',
          target: '#notification'
        }
      },
      {
        name: 'error',
        icon: 'icon--close--solid',
        message: 'Downtime',
        notification: {
          type: 'error',
          title: 'All systems are currently down or degraded.',
          message: '',
          target: '#notification'
        },
      },
      {
        name: 'error',
        icon: 'icon--close--solid',
        message: 'Downtime',
        notification: {
          type: 'error',
          title: 'Unsupported API Version',
          message: 'This version of your API is in pre-production and is currently not supported for system status checks.',
          target: '#notification'
        },
      },
      {
        name: 'info',
        icon: 'icon--info--solid',
        message: 'Downtime',
        notification: {
          type: 'info',
          title: 'Participant Node Not Set',
          message: 'There is currently no participant node set for this environment.',
          target: '#notification'
        },
      });
  }

  toggleDetailsPane(linkRef: ElementRef, currStatusDate: StatusByDate) {

    if (linkRef == null || this.openedLinkRef === linkRef) {
      // toggle opened/close
      this.detailsPaneOpened = !this.detailsPaneOpened;
    } else {
      this.detailsPaneOpened = true;
      // TODO: get details for clicked link
    }

    // set current clicked link
    this.openedLinkRef = linkRef;

    // set currently viewed error details by date
    this.detailsObject = currStatusDate;
  }

  openModal(linkRef: ElementRef, currStatusDate: StatusByDate) {

    // set current clicked link
    this.openedLinkRef = linkRef;

    // set currently viewed error details by date
    this.detailsObject = currStatusDate;

    // creates and opens the modal
    this.modalService.create({
      component: StatusDetailsModalComponent,
      inputs: {
        MODAL_DATA: this.detailsObject
      }
    });
  }
}
