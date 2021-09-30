// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package participant

import (
	"github.com/GFTN/gftn-services/automation-service/constant"
	"github.com/GFTN/gftn-services/utility/wwfirebase"
	"os"
	"testing"
)

var (
	env = "dev"
	iid = "-LmxG1giKEfq3bvbIA5m"
	participantID = "fgdsfasdfadev"
	version = "latest"
	replica = "1"
	callBackURL = ""
	rdoClientURL = ""
	dockerRegistryURL = "ip-team-worldwire-docker-local.artifactory.swg-devops.com"
)

//func Test_Deploy(t *testing.T) {
//	deployParticipantServiceChan := make(chan error)
//	go deployParticipantServices(deployParticipantServiceChan, env, participantID, version, replica, dockerRegistryURL)
//	err := <- deployParticipantServiceChan
//	if err != nil {
//		LOGGER.Errorf("%s", err.Error())
//		return
//	}
//
//	LOGGER.Debugf("Deploy %s services to %s cluster. (%s version)", participantID, env, version)
//	return
//}
//
//func Test_GetSecretTemplate(t *testing.T) {
//	os.Setenv("AWS_REGION", "ap-southeast-1")
//	template, getErr := getSecretTemplate(env)
//	if getErr != nil {
//		LOGGER.Error(getErr.Error())
//		return
//	}
//
//	decodedByte, _ := base64.StdEncoding.DecodeString(template)
//	//LOGGER.Debugf("%s", string(decodedByte))
//
//	LOGGER.Infof("Create AWS secret for participant")
//	setSecretErr := utility.SetAWSSecret(env, participantID, callBackURL, rdoClientURL, decodedByte)
//	if setSecretErr != nil {
//		LOGGER.Error(setSecretErr.Error())
//		return
//	}
//
//	LOGGER.Infof("Participant services secret successfully generated")
//	return
//}

//func Test_Kafka(t *testing.T) {
//	LOGGER.Infof("Create Kafka cert and topics for participant")
//	kcertScript := "./kafka_test.sh"
//	kcertErr := utility.RunBashCmd(kcertScript, env)
//	if kcertErr != nil {
//		LOGGER.Errorf("Failed to deploy participant Kafka certificate and topics")
//		return
//	}
//
//	out := exec.Command("/bin/bash", "-c", "source", "test.txt")
//	out.Start()
//	cmdErr := out.Wait()
//
//	if cmdErr != nil {
//		LOGGER.Errorf("Error running bash command: %s", cmdErr.Error())
//		return
//	}
//
//	brokers := os.Getenv("BROKERS_URL")
//	LOGGER.Debugf("%s", brokers)
//	LOGGER.Debugf("Participant Kafka certificate and topics successfully generated")
//
//	return
//}

func Test_FB(t *testing.T) {
	os.Setenv("FIREBASE_CREDENTIALS", "")
	os.Setenv("FIREBASE_DB_URL", "https://dev-2-c8774.firebaseio.com")
	wwfirebase.FbClient, _, _ = wwfirebase.AuthenticateWithAdminPrivileges()
	wwfirebase.FbRef = wwfirebase.GetRootRef()

	// Get the original data from fireBase
	var participantInfo map[string]interface{}
	ref := wwfirebase.FbClient.NewRef("participants/" + iid + "/nodes/").Child(participantID)
	ref.Get(wwfirebase.AppContext, &participantInfo)

	//var status []string

	if participantInfo == nil {
		LOGGER.Warning("No participant record found")
		//WriteParticipantToFireBase(input)
		//return nil
	} else {
		LOGGER.Infof("%+v", participantInfo["status"])
		for _, s := range participantInfo["status"].([]interface{}) {
			LOGGER.Infof(s.(string))
		}
	}

	status := "complete"
	result := "done"
	var statuses []string
	for key := range participantInfo {
		if key == "status" {
			count := 0
			LOGGER.Info(participantInfo[key].([]interface{})[0])
			// find if the status was already in the firebase
			for _, s := range participantInfo[key].([]interface{}) {
				if s.(string) == status {
					if result == "failed" {
						// status exist in the firebase
						count += 1
					}
				} else {
					if s.(string) == constant.StatusPending || s.(string) == constant.StatusConfiguring || s.(string) == constant.StatusComplete {
						// skip "pending" and "configuring" status
						continue
					} else {
						// get status that's already in firebase
						statuses = append(statuses, s.(string))
					}
				}
			}

			if count == 0 {
				// A new status and it is not "complete"
				if result == "failed" {
					statuses = append(statuses, status)
				} else if len(statuses) == 0 {
					if result == "resolve" {
						statuses = append(statuses, constant.StatusConfiguring)
					} else {
						statuses = append(statuses, constant.StatusComplete)
					}
				}
				participantInfo[key] = statuses
			}
		} else if key == "initialized" {
			participantInfo[key] = true
		}
	}

	//LOGGER.Debug(participantInfo)
	// Update the log into FireBase
	ref.Update(wwfirebase.AppContext, participantInfo)
}