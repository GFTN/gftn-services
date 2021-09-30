// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package participant

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/GFTN/gftn-services/automation-service/automate/participant/remove"
	"github.com/GFTN/gftn-services/automation-service/constant"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/op/go-logging"
	"github.com/GFTN/gftn-services/automation-service/environment"
	"github.com/GFTN/gftn-services/automation-service/internal_model"
	"github.com/GFTN/gftn-services/automation-service/model/model"
	"github.com/GFTN/gftn-services/automation-service/utility"
	gftn_model "github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/participant"
	"github.com/GFTN/gftn-services/utility/response"
)

var LOGGER = logging.MustGetLogger("participant")

type DeploymentOperations struct {
	dockerRegistryURL  string
	kafkaKeyPW         string
	clusterEnvironment string
	dockerImageVersion string
	awsSecretTemplate  string
}

func InitiateDeploymentOperations() (DeploymentOperations, error) {
	op := DeploymentOperations{}

	env := os.Getenv(global_environment.ENV_KEY_ENVIRONMENT_VERSION)
	imageVersion := os.Getenv(environment.ENV_DOCKER_IMAGE_VERSION)
	dockerRegistryURL := os.Getenv(environment.ENV_KEY_DOCKER_REGISTRY_URL)
	kafkaKeyPW := os.Getenv(environment.ENV_KEY_KAFKA_KEY_PW)

	// Check if all necessary variables were setup correctly
	LOGGER.Infof("Check environment variable")
	if len(dockerRegistryURL) == 0 || len(kafkaKeyPW) == 0 {
		LOGGER.Error("Missing environment variable: please check if dockerRegistryURL, kafkaKeyPW was correctly setup")
		return op, errors.New("")
	}

	// Check if the cluster environment match the exact environment
	LOGGER.Infof("Check if cluster environment is correct")
	envCheck := utility.ValidateEnv(env)
	LOGGER.Infof("Cluster env check result: %t", envCheck)
	if !envCheck {
		LOGGER.Errorf("Unsupported cluster environment: %s", env)
		return op, errors.New("unsupported cluster environment")
	}

	LOGGER.Infof("Check if the Docker image version exist")
	tagCheck := utility.ValidateDockerTag(imageVersion)
	LOGGER.Infof("Docker tag check result: %t", tagCheck)
	if !tagCheck {
		LOGGER.Errorf("Unsupported Docker image tag: %s", imageVersion)
		return op, errors.New("unsupported Docker image tag")
	}

	// Get the secret template for specific version from AWS secret manager
	LOGGER.Infof("Check if AWS secret template exist")
	template, getErr := utility.GetSecretTemplate(env, imageVersion)
	if getErr != nil {
		LOGGER.Errorf("Error getting secret template: %s", getErr.Error())
		return op, errors.New("can not found secret template")
	} else if len(template) == 0 {
		LOGGER.Error("AWS secret template file is empty")
		return op, errors.New("aws secret template file is empty")
	}

	// Switch the kubeconfig to the correct environment
	script := constant.K8sBasePath + "/script/change_cluster.sh"
	err := utility.RunBashCmd(script, env)
	if err != nil {
		LOGGER.Errorf("Error: %s", err.Error())
		return op, errors.New("failed to switch cluster environment")
	}

	op.dockerRegistryURL = dockerRegistryURL
	op.kafkaKeyPW = kafkaKeyPW
	op.clusterEnvironment = env
	op.dockerImageVersion = imageVersion
	op.awsSecretTemplate = template

	return op, nil
}

func (op *DeploymentOperations) DeployParticipantServicesAndConfigs(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("X-XSS-Protection", "1")
	var input model.Automation
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		LOGGER.Errorf("Error decoding the payload: %s", err.Error())
		response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1001", err)
		return
	}

	err = input.Validate(strfmt.Default)
	if err != nil {
		LOGGER.Errorf("Error while validating Payload: %s", err.Error())
		response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1002", err)
		return
	}

	institutionID := *input.InstitutionID
	participantID := *input.ParticipantID
	env := op.clusterEnvironment
	imageVersion := op.dockerImageVersion
	awsSecretTemplate := op.awsSecretTemplate
	dockerRegistryURL := op.dockerRegistryURL
	kafkaKeyPW := op.kafkaKeyPW
	orgId := os.Getenv(environment.ENV_KEY_AWS_ORG_ID)
	secretManagerRegion := os.Getenv(global_environment.ENV_KEY_AWS_REGION)
	participantRole := *input.Role

	LOGGER.Infof("Participant: %s, Cluster env: %s, Docker Image Tag: %s", participantID, env, imageVersion)

	// Check if the participant was recorded in the firebase
	// if not, Write the participant information into FireBase: configuring
	// if yes, check which process was failed and redo it again
	statuses := utility.CheckParticipantRecordInFireBase(input)
	if statuses == nil || len(statuses) == 0 {
		LOGGER.Infof("Create participant entry in the participant-registry: %s", participantID)
		createEntryErr := createParticipantEntry(input)
		if createEntryErr != nil {
			// Update the status in FireBase: create_participant_entry_failed
			go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreatePREntryFailed, "failed")
			response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", createEntryErr)
			return
		}

		// If the role of the participant is an issuer and it is Stronghold, only deploy the api-gateway for it
		if utility.IsStrongHold(input) {
			LOGGER.Info("Participant is Stronghold, only need to deploy api-gateway")
			err := createAWSAPIGateway(env, institutionID, participantID, orgId, "create_stronghold_aws_settings.sh")
			if err == nil {
				go func() {
					cdnErr := createAWSCustomDomainName(env, institutionID, participantID, orgId, "create_stronghold_aws_settings.sh")
					route53Err := createAWSRoute53(env, institutionID, participantID, orgId, "create_stronghold_aws_settings.sh")

					if cdnErr == nil && route53Err == nil {
						utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusComplete, "done")
					}
				}()
			}
		} else {
			LOGGER.Info("Deploy participant backend micro services")
			deployErr := deployNodesProcess(participantRole, env, orgId, kafkaKeyPW, dockerRegistryURL, secretManagerRegion, imageVersion, awsSecretTemplate, input)
			if deployErr != nil {
				response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", deployErr)
				return
			}

			// If the role of the participant is a market maker, create the issuing and operating account for it
			if participantRole == constant.MarketMakerParticipant {
				LOGGER.Info("Create participant accounts")
				go createAccountsProcess(institutionID, participantID)
			} else {
				go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusComplete, "done")
			}
		}

		response.NotifySuccess(w, req, "Successfully deployed")
		return
	} else if statuses[0] == constant.StatusComplete {
		response.NotifySuccess(w, req, "Successfully deployed")
		return
	} else {
		for _, s := range statuses {
			LOGGER.Debugf("Doing retry, status is %s", s)
			switch s {
			case constant.StatusCreatePREntryFailed:
				LOGGER.Debug("**Retry**: Create Participant entry")
				err := createParticipantEntry(input)
				if err != nil {
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreatePREntryFailed, "failed")
					response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", err)
					return
				} else {
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreatePREntryFailed, "resolve")
				}

				if utility.IsStrongHold(input) {
					LOGGER.Debug("**Retry**: Deploy api-gateway for Stronghold")
					err := createAWSAPIGateway(env, institutionID, participantID, orgId, "create_stronghold_aws_settings.sh")
					if err == nil {
						go func() {
							cdnErr := createAWSCustomDomainName(env, institutionID, participantID, orgId, "create_stronghold_aws_settings.sh")
							route53Err := createAWSRoute53(env, institutionID, participantID, orgId, "create_stronghold_aws_settings.sh")

							if cdnErr == nil && route53Err == nil {
								utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusComplete, "done")
							}
						}()
					}
				} else {
					LOGGER.Debug("**Retry**: Deploy participant backend micro services")
					deployErr := deployNodesProcess(participantRole, env, orgId, kafkaKeyPW, dockerRegistryURL, secretManagerRegion, imageVersion, awsSecretTemplate, input)
					if deployErr != nil {
						go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateMicroServicesFailed, "failed")
						response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", deployErr)
						return
					}

					if participantRole == constant.MarketMakerParticipant {
						LOGGER.Debug("**Retry**: Create participant accounts")
						go createAccountsProcess(institutionID, participantID)
					} else {
						go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusComplete, "done")
					}
				}
			case constant.StatusCreateIAMPolicyFailed:
				remove.IAMPolicy(participantID, env, orgId)
				createIAMPolicyChan := make(chan error)
				LOGGER.Debug("**Retry**: Create AWS IAM policy")
				go createIAMPolicy(createIAMPolicyChan, participantRole, env, institutionID, participantID, orgId, secretManagerRegion)
				createIAMPolicyErr := <-createIAMPolicyChan
				if createIAMPolicyErr != nil {
					// Update the status in FireBase: create_iam_policy_failed
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateIAMPolicyFailed, "failed")
					response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", createIAMPolicyErr)
					return
				} else {
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateIAMPolicyFailed, "resolve")
				}

				LOGGER.Debug("**Retry**: Deploy participant backend micro services")
				deployParticipantServiceError := deployParticipantServices(participantRole, env, participantID, imageVersion, input.Replica, dockerRegistryURL, orgId)
				if deployParticipantServiceError != nil {
					// Update the status in FireBase: create_micro_services_failed
					go utility.UpdateStatusInFireBase(*input.InstitutionID, *input.ParticipantID, constant.StatusCreateMicroServicesFailed, "failed")
					response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", deployParticipantServiceError)
					return
				}

				if participantRole == constant.MarketMakerParticipant {
					LOGGER.Debug("**Retry**: Create participant accounts")
					go createAccountsProcess(institutionID, participantID)
				} else {
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusComplete, "done")
				}
			case constant.StatusCreateKafkaTopicFailed:
				remove.KafkaCert(participantID, env, participantRole)
				LOGGER.Debug("**Retry**: Create Kafka certificate and topics")
				kcertErr := createKafkaCertAndTopics(participantRole, env, institutionID, participantID, kafkaKeyPW, dockerRegistryURL, orgId, imageVersion)
				if kcertErr != nil {
					// Update the status in FireBase: create_kafka_topic_failed
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateKafkaTopicFailed, "failed")
					response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", kcertErr)
					return
				} else {
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateKafkaTopicFailed, "done")
				}
			case constant.StatusCreateAWSSecretFailed:
				remove.AWSSecret(participantID, env)
				createSecretChan := make(chan error)
				LOGGER.Debug("**Retry**: Create AWS secret")
				go createAWSSecret(createSecretChan, participantRole, env, institutionID, participantID, awsSecretTemplate)
				createSecretError := <-createSecretChan
				if createSecretError != nil {
					// Update the status in FireBase: create_aws_secret_failed
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSSecretFailed, "failed")
					response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", createSecretError)
					return
				} else {
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSSecretFailed, "resolve")
				}

				LOGGER.Debug("**Retry**: Deploy participant backend micro services")
				deployParticipantServiceError := deployParticipantServices(participantRole, env, participantID, imageVersion, input.Replica, dockerRegistryURL, orgId)
				if deployParticipantServiceError != nil {
					// Update the status in FireBase: create_micro_services_failed
					go utility.UpdateStatusInFireBase(*input.InstitutionID, *input.ParticipantID, constant.StatusCreateMicroServicesFailed, "failed")
					response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", deployParticipantServiceError)
					return
				}

				if participantRole == constant.MarketMakerParticipant {
					LOGGER.Debug("**Retry**: Create participant accounts")
					go createAccountsProcess(institutionID, participantID)
				} else {
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusComplete, "done")
				}
			case constant.StatusCreateAWSAPIGatewayFailed:
				remove.AWSResources(participantID, env, "apigateway")
				LOGGER.Debug("**Retry**: Create AWS resources (api-gateway, custom domain name, route 53 domain)")
				switch participantRole {
				case constant.MarketMakerParticipant:
					gwErr := createAWSAPIGateway(env, institutionID, participantID, orgId, "create_aws_settings.sh")
					if gwErr != nil {
						go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSAPIGatewayFailed, "failed")
						response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", gwErr)
						return
					} else {
						utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSAPIGatewayFailed, "resolve")
						go createAWSCustomDomainName(env, institutionID, participantID, orgId, "create_aws_settings.sh")
						go createAWSRoute53(env, institutionID, participantID, orgId, "create_aws_settings.sh")
					}
				case constant.AssetIssuerParticipant:
					if utility.IsStrongHold(input) {
						gwErr := createAWSAPIGateway(env, institutionID, participantID, orgId, "create_stronghold_aws_settings.sh")
						if gwErr != nil {
							go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSAPIGatewayFailed, "failed")
							response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", gwErr)
							return
						} else {
							utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSAPIGatewayFailed, "resolve")
							go createAWSCustomDomainName(env, institutionID, participantID, orgId, "create_stronghold_aws_settings.sh")
							go createAWSRoute53(env, institutionID, participantID, orgId, "create_stronghold_aws_settings.sh")
						}
					} else {
						gwErr := createAWSAPIGateway(env, institutionID, participantID, orgId, "create_global_only_aws_settings.sh")
						if gwErr != nil {
							go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSAPIGatewayFailed, "failed")
							response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", gwErr)
							return
						} else {
							utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSAPIGatewayFailed, "resolve")
							go createAWSCustomDomainName(env, institutionID, participantID, orgId, "create_global_only_aws_settings.sh")
							go createAWSRoute53(env, institutionID, participantID, orgId, "create_global_only_aws_settings.sh")
						}
					}
				default:
					LOGGER.Error("Unknown participant role")
					response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", errors.New("unknown participant role"))
					return
				}

				go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusComplete, "done")
			case constant.StatusCreateAWSCustomDomainNameFailed:
				remove.AWSResources(participantID, env, "customdomainname")
				LOGGER.Debug("**Retry**: Create AWS Custom Domain Name")
				switch participantRole {
				case constant.MarketMakerParticipant:
					cdnErr := createAWSCustomDomainName(env, institutionID, participantID, orgId, "create_aws_settings.sh")
					if cdnErr != nil {
						go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSCustomDomainNameFailed, "failed")
						response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", cdnErr)
						return
					} else {
						go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSCustomDomainNameFailed, "done")
					}
				case constant.AssetIssuerParticipant:
					if utility.IsStrongHold(input) {
						cdnErr := createAWSCustomDomainName(env, institutionID, participantID, orgId, "create_stronghold_aws_settings.sh")
						if cdnErr != nil {
							go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSCustomDomainNameFailed, "failed")
							response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", cdnErr)
							return
						} else {
							go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSCustomDomainNameFailed, "done")
						}
					} else {
						cdnErr := createAWSCustomDomainName(env, institutionID, participantID, orgId, "create_global_only_aws_settings.sh")
						if cdnErr != nil {
							go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSCustomDomainNameFailed, "failed")
							response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", cdnErr)
							return
						} else {
							go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSCustomDomainNameFailed, "done")
						}
					}

				default:
					LOGGER.Error("Unknown participant role")
					response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", errors.New("unknown participant role"))
					return
				}
			case constant.StatusCreateAWSRoute53DomainFailed:
				remove.AWSResources(participantID, env, "route53domain")
				LOGGER.Debug("**Retry**: Create AWS Route 53 Domain")
				switch participantRole {
				case constant.MarketMakerParticipant:
					r53Err := createAWSRoute53(env, institutionID, participantID, orgId, "create_aws_settings.sh")
					if r53Err != nil {
						go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSRoute53DomainFailed, "failed")
						response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", r53Err)
						return
					} else {
						go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSRoute53DomainFailed, "done")
					}
				case constant.AssetIssuerParticipant:
					if utility.IsStrongHold(input) {
						r53Err := createAWSRoute53(env, institutionID, participantID, orgId, "create_stronghold_aws_settings.sh")
						if r53Err != nil {
							go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSRoute53DomainFailed, "failed")
							response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", r53Err)
							return
						} else {
							go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSRoute53DomainFailed, "done")
						}
					} else {
						r53Err := createAWSRoute53(env, institutionID, participantID, orgId, "create_global_only_aws_settings.sh")
						if r53Err != nil {
							go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSRoute53DomainFailed, "failed")
							response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", r53Err)
							return
						} else {
							go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSRoute53DomainFailed, "done")
						}
					}

				default:
					LOGGER.Error("Unknown participant role")
					response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", errors.New("unknown participant role"))
					return
				}
			case constant.StatusCreateAWSDynamoDBFailed:
				remove.AWSResources(participantID, env, "dynamodb")
				dbErr := createAWSDynamoDB(env, institutionID, participantID, orgId, "create_aws_settings.sh")
				if dbErr != nil {
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSDynamoDBFailed, "failed")
					response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", dbErr)
					return
				} else {
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSDynamoDBFailed, "done")
				}
			case constant.StatusCreateMicroServicesFailed:
				LOGGER.Debug("**Retry**: Deploy participant backend micro services")
				deployParticipantServiceError := deployParticipantServices(participantRole, env, participantID, imageVersion, input.Replica, dockerRegistryURL, orgId)
				if deployParticipantServiceError != nil {
					// Update the status in FireBase: create_micro_services_failed
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateMicroServicesFailed, "failed")
					response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", deployParticipantServiceError)
					return
				} else {
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateMicroServicesFailed, "resolve")
				}

				if participantRole == constant.MarketMakerParticipant {
					LOGGER.Debug("**Retry**: Create participant accounts")
					go createAccountsProcess(institutionID, participantID)
				} else {
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusComplete, "done")
				}
			case constant.StatusCreateIssuingAccountFailed:
				LOGGER.Debug("**Retry**: Create issuing account")
				apiSvcClient, getErr := getAPISVCClient(participantID)
				if getErr != nil {
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateIssuingAccountFailed, "failed")
					response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", getErr)
					return
				}

				createIssuingAccountErr := createIssuingAccount(apiSvcClient)
				if createIssuingAccountErr != nil {
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateIssuingAccountFailed, "failed")
					response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", createIssuingAccountErr)
					return
				} else {
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateIssuingAccountFailed, "done")
				}
			case constant.StatusCreateOperatingAccountFailed:
				LOGGER.Debug("**Retry**: Create operating account")
				apiSvcClient, getErr := getAPISVCClient(participantID)
				if getErr != nil {
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateOperatingAccountFailed, "failed")
					response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", getErr)
					return
				}

				createOperatingAccountErr := createOperatingAccount(apiSvcClient)
				if createOperatingAccountErr != nil {
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateOperatingAccountFailed, "failed")
					response.NotifyWWError(w, req, http.StatusBadRequest, "DEPLOYMENT-1100", createOperatingAccountErr)
					return
				} else {
					go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateOperatingAccountFailed, "done")
				}
			default:
				LOGGER.Warningf("Unknown status: %s", s)
			}
		}
		response.NotifySuccess(w, req, "Successfully deployed")
		return
	}
}

func deployNodesProcess(role, env, orgId, kafkaKeyPW, dockerRegistryURL, secretManagerRegion, imageVersion, secretTemplate string, input model.Automation) error {
	institutionID := *input.InstitutionID
	participantID := *input.ParticipantID
	replica := input.Replica

	createIAMPolicyChan := make(chan error)
	createSecretChan := make(chan error)

	LOGGER.Debugf("Deploying participant: %s", role)
	// Create IAM policy for each service
	go createIAMPolicy(createIAMPolicyChan, role, env, institutionID, participantID, orgId, secretManagerRegion)

	// Create Kafka certificate and topics
	go createKafkaCertAndTopics(role, env, institutionID, participantID, kafkaKeyPW, dockerRegistryURL, orgId, imageVersion)

	// Create AWS secret
	go createAWSSecret(createSecretChan, role, env, institutionID, participantID, secretTemplate)

	// Create AWS settings (api-gateway, custom domain name, route53 domain, dynamoDB)
	go createAWSSettings(role, env, institutionID, participantID, orgId)

	LOGGER.Debug("---------- Waiting for all the pre-process to be finished ----------")

	createSecretError := <-createSecretChan
	if createSecretError != nil {
		return createSecretError
	}

	createIAMPolicyErr := <-createIAMPolicyChan
	if createIAMPolicyErr != nil {
		return createIAMPolicyErr
	}

	LOGGER.Debug("---------- All the pre-process were successfully finished ----------")

	deployParticipantServiceError := deployParticipantServices(role, env, participantID, imageVersion, replica, dockerRegistryURL, orgId)
	if deployParticipantServiceError != nil {
		// Update the status in FireBase: create_micro_services_failed
		go utility.UpdateStatusInFireBase(*input.InstitutionID, *input.ParticipantID, constant.StatusCreateMicroServicesFailed, "failed")
		return deployParticipantServiceError
	}

	LOGGER.Info("===> Successfully deployed participant services")

	return nil
}

func createAccountsProcess(institutionID, participantID string) {
	apiSvcClient, getErr := getAPISVCClient(participantID)
	if getErr != nil {
		LOGGER.Errorf("Failed to create api-service http client: %s", getErr.Error())
		go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateIssuingAccountFailed, "failed")
		return
	}

	// Wait for the api-service to start up
	LOGGER.Debug("Create issuing account")
	time.Sleep(time.Second * 120)

	createIssuingAccountErr := createIssuingAccount(apiSvcClient)
	if createIssuingAccountErr != nil {
		LOGGER.Errorf("Failed to create issuing account: %s", createIssuingAccountErr.Error())
		// Update the status in FireBase: create_issuing_account_failed
		go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateIssuingAccountFailed, "failed")
	}

	// Wait for the api-service to start up
	LOGGER.Debug("Create operating account")
	time.Sleep(time.Second * 120)

	createOperatingAccountErr := createOperatingAccount(apiSvcClient)
	if createOperatingAccountErr != nil {
		LOGGER.Errorf("Failed to create operating account: %s", createOperatingAccountErr.Error())
		// Update the status in FireBase: create_operating_account_failed
		go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateOperatingAccountFailed, "failed")
		return
	}

	LOGGER.Info("===> Successfully create participant accounts")
	go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusComplete, "done")
	return
}

func createIAMPolicy(createIAMPolicyChan chan<- error, role, env, institutionID, participantID, orgId, secretManagerRegion string) {
	switch role {
	case constant.MarketMakerParticipant:
		for _, svcName := range constant.ParticipantServices {
			if svcName == "participant" {
				continue
			}

			createErr := utility.CreateIAMPolicy(participantID, svcName, env, orgId, secretManagerRegion)
			if createErr != nil {
				LOGGER.Errorf("Error create IAM policy for service %s: %s", svcName, createErr.Error())
				// Update the status in FireBase: create_iam_policy_failed
				go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateIAMPolicyFailed, "failed")
				createIAMPolicyChan <- createErr
				return
			}
		}

		LOGGER.Info("Successfully create IAM policy for each service")
		createIAMPolicyChan <- nil
	case constant.AssetIssuerParticipant:
		svcName := "ww-gateway"
		createErr := utility.CreateIAMPolicy(participantID, svcName, env, orgId, secretManagerRegion)
		if createErr != nil {
			LOGGER.Errorf("Error create IAM policy for service %s: %s", svcName, createErr.Error())
			// Update the status in FireBase: create_iam_policy_failed
			go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateIAMPolicyFailed, "failed")
			createIAMPolicyChan <- createErr
			return
		}

		LOGGER.Info("Successfully create IAM policy for ww-gateway")
		createIAMPolicyChan <- nil
	}

	return
}

func createKafkaCertAndTopics(role, env, institutionID, participantID, kafkaKeyPW, dockerRegistryURL, orgId, imageVersion string) error {
	switch role {
	case constant.MarketMakerParticipant:
		// Create Kafka certificate for each participant
		LOGGER.Infof("Create Kafka cert and topics for participant")
		kcertScript := constant.MSKBasePath + "/create_cert_and_topic.sh"
		kcertErr := utility.RunBashCmd(kcertScript, participantID, kafkaKeyPW, env, constant.MSKBasePath, dockerRegistryURL, orgId, imageVersion)
		if kcertErr != nil {
			LOGGER.Errorf("Failed to deploy participant Kafka certificate and topics")
			// Update the status in FireBase: create_kafka_topic_failed
			go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateKafkaTopicFailed, "failed")
			return kcertErr
		}

		LOGGER.Debugf("Participant Kafka certificate and topics successfully generated")
	case constant.AssetIssuerParticipant:
		// Create Kafka certificate for anchor
		LOGGER.Infof("Create Kafka cert and topics for anchor")
		kcertScript := constant.MSKBasePath + "/create_is_cert_and_topic.sh"
		kcertErr := utility.RunBashCmd(kcertScript, participantID, kafkaKeyPW, env, constant.MSKBasePath, dockerRegistryURL, orgId, imageVersion)
		if kcertErr != nil {
			LOGGER.Errorf("Failed to deploy anchor Kafka certificate and topics")
			// Update the status in FireBase: create_kafka_topic_failed
			go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateKafkaTopicFailed, "failed")
			return kcertErr
		}

		LOGGER.Debugf("Anchor Kafka certificate and topics successfully generated")
	}
	return nil
}

func createAWSSecret(createSecretChan chan<- error, role, env, institutionID, participantID, secretTemplate string) {
	// Check if the secrets for participant services already exist in the AWS secret manager
	checkResult := utility.CheckAWSSecretExist(env, participantID, "initialize")
	if checkResult != nil {
		// Currently all the secret store in region ap-southeast-1 (Singapore).
		// So switch the region to ap-southeast-1 in order to access the secret.
		os.Setenv(global_environment.ENV_KEY_AWS_REGION, "ap-southeast-1")
		LOGGER.Infof("Switch AWS region to %s for accessing the secret manager", os.Getenv(global_environment.ENV_KEY_AWS_REGION))

		decodedByte, _ := base64.StdEncoding.DecodeString(secretTemplate)

		// Create secret in AWS secret manager for participant
		LOGGER.Infof("Create AWS secret for participant")
		setSecretErr := utility.SetAWSSecret(role, env, participantID, decodedByte)
		if setSecretErr != nil {
			// Update the status in FireBase: create_aws_secret_failed
			go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSSecretFailed, "failed")
			createSecretChan <- errors.New("failed to create participant services secret")
			return
		}
		LOGGER.Infof("Participant services secret successfully generated")
		createSecretChan <- nil
		return
	}
	createSecretChan <- nil
	return
}

func createParticipantEntry(input model.Automation) error {
	// Construct a participant record
	prServiceURL := os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL)
	participantModel := &gftn_model.Participant{
		Bic:         input.Bic,
		CountryCode: input.CountryCode,
		ID:          input.ParticipantID,
		Role:        input.Role,
	}

	participantPayload, _ := json.Marshal(participantModel)

	prClient := &internal_model.Client{
		HTTPClient: &http.Client{Timeout: time.Second * 30},
		URL:        prServiceURL,
	}

	// call the /internal/pr/domain/input.ParticipantID "GET" endpoint to check if the participant already exist in the participant registry
	rGet, err := http.NewRequest(http.MethodGet, prClient.URL+"/internal/pr/domain/"+*input.ParticipantID, bytes.NewBuffer(participantPayload))
	if err != nil {
		LOGGER.Errorf("Error creating HTTP request: %s", err.Error())
		return err
	}

	resGet, restGetErr := http.DefaultClient.Do(rGet)
	if restGetErr != nil {
		LOGGER.Errorf("Unable to get participant from pr-serivce: %s", restGetErr.Error())
		return restGetErr
	} else if resGet.StatusCode != http.StatusOK {
		if resGet.StatusCode == http.StatusBadRequest {
			for i := 1; i < 6; i++ {
				LOGGER.Warningf("Service not ready, retry it %d", i)
				time.Sleep(10 * time.Second)
				rGet, err := http.NewRequest(http.MethodGet, prClient.URL+"/internal/pr/domain/"+*input.ParticipantID, bytes.NewBuffer(participantPayload))
				if err != nil {
					LOGGER.Errorf("Error creating HTTP request: %s", err.Error())
					return err
				}

				resGet, restErr := http.DefaultClient.Do(rGet)
				if restErr != nil {
					return restErr
				} else if resGet.StatusCode != http.StatusOK {
					if resGet.StatusCode == http.StatusBadRequest {
						if i == 5 {
							return errors.New("unable to get participant from pr-serivce")
						} else {
							continue
						}
					} else if resGet.StatusCode == http.StatusNotFound {
						LOGGER.Warning("Participant did not exist")
						break
					}
				} else {
					LOGGER.Error("Participant already exist")
					return errors.New("participant already exist")
				}
			}
		} else if resGet.StatusCode == http.StatusNotFound {
			LOGGER.Warning("Participant did not exist")
		}
	} else {
		LOGGER.Error("Participant already exist")
		return errors.New("participant already exist")
	}

	LOGGER.Debug("Create participant entry")
	// call the /internal/pr "POST" endpoint to create a participant entry into participant registry
	r, reqErr := http.NewRequest(http.MethodPost, prClient.URL+"/internal/pr", bytes.NewBuffer(participantPayload))
	if reqErr != nil {
		LOGGER.Errorf("Error creating HTTP request: %s", reqErr.Error())
		return reqErr
	}

	res, restErr := http.DefaultClient.Do(r)
	if restErr != nil {
		LOGGER.Errorf("Unable to create participant into pr-serivce: %s", restErr.Error())
		return restErr
	} else if res.StatusCode != http.StatusOK {
		responseBody, responseErr := ioutil.ReadAll(resGet.Body)
		if responseErr != nil {
			return errors.New("error trying to read the response from pr-service: " + responseErr.Error())
		}
		LOGGER.Errorf("%s", string(responseBody))
		LOGGER.Errorf("Unable to create participant into pr-serivce: %d", res.StatusCode)
		return errors.New("unable to create participant entry, http response status code from pr-service endpoint is not 200")
	}

	// Activate the participant using pr-service endpoint
	participantStatus := "active"
	participantStatusModel := &gftn_model.ParticipantStatus{
		Status: &participantStatus,
	}

	participantStatusPayload, _ := json.Marshal(participantStatusModel)

	LOGGER.Debug("Activate participant")
	// call the internal/pr/{participantID}/status endpoint to update the status to `active`
	statusReq, reqErr := http.NewRequest(http.MethodPut, prClient.URL+"/internal/pr/"+*input.ParticipantID+"/status", bytes.NewBuffer(participantStatusPayload))
	if reqErr != nil {
		LOGGER.Errorf("Error creating HTTP request: %s", reqErr.Error())
		return reqErr
	}

	statusRes, restErr := http.DefaultClient.Do(statusReq)
	if restErr != nil {
		LOGGER.Errorf("Unable to update the status: %s", restErr.Error())
		return restErr
	} else if statusRes.StatusCode != http.StatusOK {
		responseBody, responseErr := ioutil.ReadAll(resGet.Body)
		if responseErr != nil {
			return errors.New("error trying to read the response from pr-service: " + responseErr.Error())
		}
		LOGGER.Errorf("%s", string(responseBody))
		LOGGER.Errorf("Unable to update the status: %d", statusRes.StatusCode)
		return errors.New("unable to update the participant status, http response status code from pr-service endpoint is not 200")
	}

	LOGGER.Info("===> Successfully create participant entry in participant registry")

	return nil
}

func createAWSSettings(role, env, institutionID, participantID, orgId string) {
	switch role {
	case constant.MarketMakerParticipant:
		err := createAWSAPIGateway(env, institutionID, participantID, orgId, "create_aws_settings.sh")
		if err == nil {
			go createAWSCustomDomainName(env, institutionID, participantID, orgId, "create_aws_settings.sh")
			go createAWSRoute53(env, institutionID, participantID, orgId, "create_aws_settings.sh")
		}
		createAWSDynamoDB(env, institutionID, participantID, orgId, "create_aws_settings.sh")
	case constant.AssetIssuerParticipant:
		err := createAWSAPIGateway(env, institutionID, participantID, orgId, "create_global_only_aws_settings.sh")
		if err == nil {
			go createAWSCustomDomainName(env, institutionID, participantID, orgId, "create_global_only_aws_settings.sh")
			go createAWSRoute53(env, institutionID, participantID, orgId, "create_global_only_aws_settings.sh")
		}
	}
	return
}

func createAWSAPIGateway(env, institutionID, participantID, orgId, scriptName string) error {
	// Deploy AWS settings
	LOGGER.Infof("Setting up AWS API-Gateway for participant")
	script := constant.K8sBasePath + "/script/" + scriptName
	gwErr := utility.RunBashCmd(script, participantID, env, "apigateway", orgId)
	if gwErr != nil {
		LOGGER.Errorf("Failed to set up participant AWS API-Gateway")
		// Update the status in FireBase: create_aws_api_gateway_failed
		go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSAPIGatewayFailed, "failed")
		return gwErr
	}
	LOGGER.Debugf("Participant AWS API-Gateway successfully generated")
	return nil
}

func createAWSCustomDomainName(env, institutionID, participantID, orgId, scriptName string) error {
	// Deploy AWS settings
	LOGGER.Infof("Setting up AWS Custom Domain Name for participant")
	script := constant.K8sBasePath + "/script/" + scriptName
	cdnErr := utility.RunBashCmd(script, participantID, env, "customdomainname", orgId)
	if cdnErr != nil {
		LOGGER.Errorf("Failed to set up participant AWS Custom Domain Name")
		// Update the status in FireBase: create_aws_domain_custom_domain_name_failed
		go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSCustomDomainNameFailed, "failed")
		return cdnErr
	}
	LOGGER.Debugf("Participant AWS Custom Domain Name successfully generated")
	return nil
}

func createAWSRoute53(env, institutionID, participantID, orgId, scriptName string) error {
	// Deploy AWS settings
	LOGGER.Infof("Setting up AWS Route 53 domain for participant")
	script := constant.K8sBasePath + "/script/" + scriptName
	r53Err := utility.RunBashCmd(script, participantID, env, "route53domain", orgId)
	if r53Err != nil {
		LOGGER.Errorf("Failed to set up participant AWS Route 53 domain")
		// Update the status in FireBase: create_aws_route53_domain_failed
		go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSRoute53DomainFailed, "failed")
		return r53Err
	}
	LOGGER.Debugf("Participant AWS Route 53 domain successfully generated")
	return nil
}

func createAWSDynamoDB(env, institutionID, participantID, orgId, scriptName string) error {
	// Deploy AWS settings
	LOGGER.Infof("Setting up AWS DynamoDB for participant")
	script := constant.K8sBasePath + "/script/" + scriptName
	dbErr := utility.RunBashCmd(script, participantID, env, "dynamodb", orgId)
	if dbErr != nil {
		LOGGER.Errorf("Failed to set up participant AWS DynamoDB")
		// Update the status in FireBase: create_aws_dynamodb_failed
		go utility.UpdateStatusInFireBase(institutionID, participantID, constant.StatusCreateAWSDynamoDBFailed, "failed")
		return dbErr
	}
	LOGGER.Debugf("Participant AWS DynamoDB successfully generated")
	return nil
}

func deployParticipantServices(role, env, participantID, imageVersion, replica, dockerRegistryURL, orgId string) error {
	switch role {
	case constant.MarketMakerParticipant:
		// Deploy participant services
		LOGGER.Infof("Deploy participant micro services on Kubernetes cluster")
		pScript := constant.K8sBasePath + "/script/create_participant_msvc.sh"
		pErr := utility.RunBashCmd(pScript, participantID, env, imageVersion, replica, dockerRegistryURL, orgId)
		if pErr != nil {
			LOGGER.Errorf("Failed to deploy participant services")
			return errors.New("failed to deploy participant services")
		}

		LOGGER.Debugf("Participant services successfully generated")
	case constant.AssetIssuerParticipant:
		// Deploy anchor's local services
		LOGGER.Infof("Deploy anchor local micro services on Kubernetes cluster")
		pScript := constant.K8sBasePath + "/script/create_anchor_msvc.sh"
		pErr := utility.RunBashCmd(pScript, participantID, env, imageVersion, replica, dockerRegistryURL, orgId)
		if pErr != nil {
			LOGGER.Errorf("Failed to deploy anchor local micro services")
			return errors.New("failed to deploy anchor local micro services")
		}

		LOGGER.Debugf("Anchor local micro services successfully generated")
	}

	return nil
}

func getAPISVCClient(participantID string) (*internal_model.Client, error) {
	// Call the api-service internal endpoint
	url := os.Getenv(global_environment.ENV_KEY_API_SVC_URL)
	apiServiceURL, convertErr := participant.GetServiceUrl(url, participantID)
	if convertErr != nil {
		LOGGER.Error(convertErr.Error())
		return nil, convertErr
	}

	apiSvcClient := &internal_model.Client{
		HTTPClient: &http.Client{Timeout: time.Second * 80},
		URL:        apiServiceURL,
	}

	return apiSvcClient, nil
}

func createIssuingAccount(apiSvcClient *internal_model.Client) error {
	// call the /internal/accounts/issuing endpoint to create issuing account
	issuingR, err := http.NewRequest(http.MethodPost, apiSvcClient.URL+"/internal/accounts/issuing", nil)
	if err != nil {
		LOGGER.Errorf("Error creating HTTP request: %s", err.Error())
		return err
	}

	resIssuing, restErr := http.DefaultClient.Do(issuingR)
	if restErr != nil || resIssuing.StatusCode != http.StatusOK {
		LOGGER.Errorf("Unable to create participant issuing account")
		for i := 1; i < 6; i++ {
			LOGGER.Warningf("Service not ready, retry it %d", i)
			time.Sleep(15 * time.Second)
			issuingR, err := http.NewRequest(http.MethodPost, apiSvcClient.URL+"/internal/accounts/issuing", nil)
			if err != nil {
				LOGGER.Errorf("Error creating HTTP request: %s", err.Error())
				return err
			}

			resIssuing, restErr := http.DefaultClient.Do(issuingR)
			if restErr != nil || resIssuing.StatusCode != http.StatusOK {
				LOGGER.Errorf("Unable to create participant issuing account")
				if i == 5 {
					LOGGER.Errorf("Unable to create participant issuing account")
					return errors.New("unable to create participant issuing account")
				} else {
					continue
				}
			} else {
				break
			}
		}
	}

	return nil
}

func createOperatingAccount(apiSvcClient *internal_model.Client) error {
	// call the /internal/accounts/default endpoint to create issuing account
	defaultR, err := http.NewRequest(http.MethodPost, apiSvcClient.URL+"/internal/accounts/default", nil)
	if err != nil {
		LOGGER.Errorf("Error creating HTTP request: %s", err.Error())
		return err
	}

	resDefault, restErr := http.DefaultClient.Do(defaultR)
	if restErr != nil || resDefault.StatusCode != http.StatusOK {
		LOGGER.Errorf("Unable to create participant operating account")
		for i := 1; i < 6; i++ {
			LOGGER.Warningf("Service not ready, retry it %d", i)
			time.Sleep(15 * time.Second)
			defaultR, err := http.NewRequest(http.MethodPost, apiSvcClient.URL+"/internal/accounts/default", nil)
			if err != nil {
				LOGGER.Errorf("Error creating HTTP request: %s", err.Error())
				return err
			}

			resDefault, restErr := http.DefaultClient.Do(defaultR)
			if restErr != nil || resDefault.StatusCode != http.StatusOK {
				LOGGER.Errorf("Unable to create participant operating account")
				if i == 5 {
					LOGGER.Errorf("Unable to create participant operating account")
					return errors.New("unable to create participant operating account")
				} else {
					continue
				}
			} else {
				break
			}
		}
	}

	return nil
}