// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable } from '@angular/core';
import { lowerCase } from 'lodash';

@Injectable()
export class NodeService {

  // user-friendly mappings for node statuses
  // shared among components that need and use these
  statusMap: Map<string, string> = new Map([
    ['pending', 'pending'],
    ['configuring', 'pending'],
    ['complete', 'complete exists'],
    ['configuration_failed', 'failed'],
    ['create_participant_entry_failed', 'failed'],
    ['create_iam_policy_failed', 'failed'],
    ['create_kafka_topic_failed', 'failed'],
    ['create_aws_secret_failed', 'failed'],
    ['create_aws_api_gateway_failed', 'failed'],
    ['create_aws_domain_custom_domain_name_failed', 'failed'],
    ['create_aws_route53_domain_failed', 'failed'],
    ['create_aws_dynamodb_failed', 'failed'],
    ['create_micro_services_failed', 'failed'],
    ['create_issuing_account_failed', 'failed'],
    ['create_operating_account_failed', 'failed'],
    ['deleted', 'deleted'],
  ]);

  constructor() { }

  /**
   * Get status from node mapping
   *
   * @param {string[]} status
   * @returns {string}
   * @memberof NodeService
   */
  public getStatus(status: string[]): string {

    if (this.statusMap.get(status[0])) {
      return this.statusMap.get(status[0]);
    }
    return status[0];

  }

  public toHumanReadable(status: string): string {
    return lowerCase(status);
  }

}
