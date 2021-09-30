// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Input, TemplateRef } from '@angular/core';

/**
 * NotificationTile
 *
 * This combines the tile and Inline Notification components
 * into a stylized form of tile used for displaying
 * different states of a data object.
 * This provides more control over the notification object
 * by allow for actions on the actual tile, in-content actions
 * and other custom tile stylings not present
 * in the default Carbon Design Notification.
 *
 * @export
 * @class NotificationTileComponent
 * @implements {OnInit}
 */
@Component({
  selector: 'app-notification-tile',
  templateUrl: './notification-tile.component.html',
  styleUrls: ['./notification-tile.component.scss']
})


export class NotificationTileComponent implements OnInit {

  // extending the Carbon styles for the Notification Tile
  @Input() type: 'info' | 'success' | 'warning' | 'error' = 'info';

  // main title of Notification Tile
  @Input() title = '';

  // main content of the tile.
  // Can pass in a template for custom HTML or template
  @Input() content: TemplateRef<any> | string;

  // content can be a template
  hasContentTemplate = false;

  // OPTIONAL: link for the notification,
  // to provide a corresponding action to the tile
  @Input() link = '';

  @Input() hasAction = false;

  constructor() { }

  ngOnInit() {
    this.hasContentTemplate = this.content instanceof TemplateRef;
  }

}
