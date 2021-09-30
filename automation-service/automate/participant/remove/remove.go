// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package remove

import (
	"github.com/op/go-logging"
	"github.com/GFTN/gftn-services/automation-service/constant"
	"github.com/GFTN/gftn-services/automation-service/utility"
)

var LOGGER = logging.MustGetLogger("remove")

func IAMPolicy(participantID, env, orgId string) {
	LOGGER.Infof("Remove IAM policies")

	rIAMScript := constant.K8sBasePath + "/script/delete_iam_policy.sh"
	utility.RunBashCmd(rIAMScript, participantID, env, orgId)

	LOGGER.Debugf("IAM policies successfully removed")
	return
}

func KafkaCert(participantID, env, role string) {
	LOGGER.Info("Remove k8s secret and job")
	switch role {
	case constant.AssetIssuerParticipant:
		rJobScript := constant.MSKBasePath + "/delete_is_kafka_secret_job.sh"
		utility.RunBashCmd(rJobScript, participantID, env)
	case constant.MarketMakerParticipant:
		rJobScript := constant.MSKBasePath + "/delete_kafka_secret_job.sh"
		utility.RunBashCmd(rJobScript, participantID, env)
	}

	LOGGER.Info("Remove certificate in AWS ACM")
	rACMScript := constant.MSKBasePath + "/delete_acm.sh"
	utility.RunBashCmd(rACMScript, participantID, env)

	LOGGER.Debugf("Kafka certificate successfully removed")
	return
}

func AWSSecret(participantID, env string) {
	LOGGER.Info("Remove AWS secret")
	rSecretScript := constant.K8sBasePath + "/script/delete_secret.sh"
	utility.RunBashCmd(rSecretScript, participantID, env)

	LOGGER.Debugf("AWS secret successfully removed")
	return
}

func AWSResources(participantID, env, resource string) {
	LOGGER.Info("Remove AWS resources")
	rResourcesScript := constant.K8sBasePath + "/script/delete_aws_resources.sh"
	utility.RunBashCmd(rResourcesScript, participantID, env, resource)

	LOGGER.Debugf("AWS resources successfully removed")
	return
}