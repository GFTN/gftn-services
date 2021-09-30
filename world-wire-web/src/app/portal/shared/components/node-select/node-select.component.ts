// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Input, Output, EventEmitter, OnDestroy } from '@angular/core';
import { SessionService } from '../../../../shared/services/session.service';
import { ListItem } from 'carbon-components-angular';
import { INodeAutomation } from '../../../../shared/models/node.interface';
import { VERSION_DETAILS, IWWApi } from '../../../../shared/constants/versions.constant';
import { find } from 'lodash';

@Component({
  selector: 'app-node-select',
  templateUrl: './node-select.component.html',
  styleUrls: ['./node-select.component.scss']
})
export class NodeSelectComponent implements OnInit {

  @Input() showVersion = false;

  init = false;

  currentNodeId: string;

  currentNodes: ListItem[];

  constructor(
    public sessionService: SessionService
  ) { }

  ngOnInit() {
    this.currentNodes = this.setCurrentNodes();

    this.currentNodeId = this.sessionService.currentNode ? this.sessionService.currentNode.participantId : null;

  }

  /**
   * Sets the dropdown list of available nodes
   * for the current environment
   */
  setCurrentNodes(): ListItem[] {
    const list: ListItem[] = [];

    // setting drodown init to false since dropdown list is being reset
    this.init = false;

    if (this.sessionService.institutionNodes) {

      for (const node of this.sessionService.institutionNodes) {
        const item = {
          value: node,
          content: node.participantId,
          selected: false,
        };

        if (node === this.sessionService.currentNode) {
          item.selected = true;
        }
        list.push(item);
      }
    }
    return list;
  }

  /**
   * Get value from selected node to set in service
   * and notify parent component of a change.
   * @param event
   * is of type any, following the ibm-dropdown spec
   */
  selectNode(event: any) {

    const item = event.item;

    // use init to prevent selection event from firing
    // upon initialization of the dropdown
    if (item && item.content && this.init) {

      const node = item.value;
      this.sessionService.setCurrentNode(node);

      // this.currentNodeId = item.content;
      this.sessionService.propogateNodeChange();
    }

    if (!this.init) {
      this.init = true;
    }
  }

  /**
   * Get Human-readable name for version
   * as opposed to release tag of the API.
   * @param releaseTag
   */
  getVersionName(releaseTag: string): string {
    const versionDetail = find(VERSION_DETAILS, (details: IWWApi) => {
      return details.releaseTag === releaseTag;
    });

    // default to text noting that current version of API is in pre-production
    // but DON'T release the actual release tag name in case this changes
    return versionDetail ? versionDetail.version : '(pre-production)';
  }
}
