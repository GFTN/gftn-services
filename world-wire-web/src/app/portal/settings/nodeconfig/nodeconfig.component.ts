// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit } from '@angular/core';
import { ENVIRONMENT } from '../../../shared/constants/general.constants';
import { INodeAutomation, NodeConfigData } from '../../../shared/models/node.interface';
import { IInstitution } from '../../../shared/models/participant.interface';
import { SessionService } from '../../../shared/services/session.service';

@Component({
  selector: 'app-nodeconfig',
  templateUrl: './nodeconfig.component.html',
  styleUrls: ['./nodeconfig.component.scss']
})
export class NodeConfigComponent implements OnInit {

  public institutions: IInstitution[] = [];

  // stores get node configurations
  public nodeConfigs: Map<string, NodeConfigData> = new Map();

  public currInstitution: IInstitution;

  constructor(
    public sessionService: SessionService,
  ) { }

  ngOnInit() {

    this.currInstitution = this.sessionService.institution;
  }

  /**
   * returns current environment name (short)
   * from configuration
   *
   * @returns
   * @memberof NodeConfigComponent
   */
  getShortEnvironmentName() {
    return ENVIRONMENT.name;
  }

  /**
   * returns current environment name (full)
   * from configuration
   *
   * @returns
   * @memberof NodeConfigComponent
   */
  getFullEnvironmentName() {
    return ENVIRONMENT.text;
  }
  /**
   * Gets the configuration data for
   * the selected node
   * @param node
   * returns NodeConfigData
   */
  getConfigData(node: INodeAutomation): NodeConfigData {

    const data: NodeConfigData = node;

    return data;
  }

}
