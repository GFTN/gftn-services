// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package constant

var (
	K8sBasePath = "/var/k8s"
	MSKBasePath = "/var/msk"
	KafkaBasePath = "/var/kafka-cluster"

	AvailableEnvs = []string{"dev", "qa", "st", "tn", "pen", "prod", "eksdev", "eksqa"}
	ParticipantServices = []string{"api-service", "crypto-service", "payment-service", "send-service", "ww-gateway", "participant"}
	GlobalDomain = "ww"
	Participant = "participant"
	Deployment = "deployment"
	Global = "global"

	StatusPending = "pending"
	StatusConfiguring = "configuring"
	StatusConfigurationFailed = "configuration_failed"
	StatusCreatePREntryFailed = "create_participant_entry_failed"
	StatusCreateIAMPolicyFailed = "create_iam_policy_failed"
	StatusCreateKafkaTopicFailed = "create_kafka_topic_failed"
	StatusCreateAWSSecretFailed = "create_aws_secret_failed"
	StatusCreateAWSAPIGatewayFailed = "create_aws_api_gateway_failed"
	StatusCreateAWSCustomDomainNameFailed = "create_aws_domain_custom_domain_name_failed"
	StatusCreateAWSRoute53DomainFailed = "create_aws_route53_domain_failed"
	StatusCreateAWSDynamoDBFailed = "create_aws_dynamodb_failed"
	StatusCreateMicroServicesFailed = "create_micro_services_failed"
	StatusCreateIssuingAccountFailed = "create_issuing_account_failed"
	StatusCreateOperatingAccountFailed = "create_operating_account_failed"
	StatusComplete = "complete"

	MarketMakerParticipant = "MM"
	AssetIssuerParticipant = "IS"

	StrongholdID = "stronghold"
	StrongholdBICCode = "USASTGHD101"
)