// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { ParticipantRole } from "./participant.interface";

export interface INodeAutomation {

    // id of firebase participant record per '/participants' in firebase
    institutionId: string;

    // Release tag of the API version of this node
    version?: string;

    // id per '/nodes' in firebase (aka: homeDomain or homeId per client-api env vars)
    participantId: string;

    // country code of this participant
    countryCode: string;

    // Role of the registered participant on the network
    // 'MM' = Market Maker/regular participant, issues only DOs
    // 'IS' = Issuer of real world DAs/stable coins in addition to DOs
    role: ParticipantRole;


    // OPTIONAL: BIC code, for banks
    bic?: string;

    // see details below:
    status?: string[]

    // fields per world wire infrastructure (status codes in the firebase db):
    // pending
    // configuring
    // configuration_failed
    // complete
    // create_participant_entry_failed
    // create_iam_policy_failed
    // create_kafka_topic_failed
    // create_aws_secret_failed
    // create_aws_api_gateway_failed
    // create_aws_domain_custom_domain_name_failed
    // create_aws_route53_domain_failed
    // create_aws_dynamodb_failed
    // create_micro_services_failed
    // create_issuing_account_failed
    // create_operating_account_failed
    // deleted

    // initialized is true if participant already added to registry and false if not
    initialized: boolean;

    // list of approval ids for any requested actions (create, update, delete)
    approvalIds?: string[];

    accountApprovalId?: string;

    // address of issuing account. Single account for anchors
    issuingAccount?: string;

    // information requested for update for any requested action (create, update, delete)
    update?: INodeAutomation;
}

export interface NodeConfigData extends INodeAutomation {
    participantIdBase?: string;
}
