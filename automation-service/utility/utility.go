// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package utility

import (
	"encoding/json"
	"errors"
	"github.com/GFTN/gftn-services/automation-service/constant"
	"github.com/GFTN/gftn-services/utility/global-environment"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/op/go-logging"
	"github.com/GFTN/gftn-services/automation-service/environment"
	"github.com/GFTN/gftn-services/automation-service/model/model"
	"github.com/GFTN/gftn-services/utility/wwfirebase"
)

var LOGGER = logging.MustGetLogger("utility")

type FBLog struct {
	callbackUrl   string
	countryCode   string
	env           string
	initialized   bool
	institutionId string
	participantId string
	status        string
	version       string
}

type ClusterInfo struct {
	Clusters []Cluster `json:"clusters"`
}

type Cluster struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

type Tags struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func RunKubeCmd(action, target, arg string) error {
	out := exec.Command("kubectl", action, target, arg)
	out.Start()
	cmdErr := out.Wait()

	if cmdErr != nil {
		LOGGER.Errorf("Error: %s", cmdErr.Error())
		return cmdErr
	}

	return nil
}

func RunBashCmd(args ...string) error {
	out := exec.Command("/bin/sh", args...)
	out.Start()
	cmdErr := out.Wait()

	if cmdErr != nil {
		LOGGER.Errorf("Error running bash command: %s", cmdErr.Error())
		return cmdErr
	}

	return nil
}

func ValidateEnv(env string) bool {
	for _, e := range constant.AvailableEnvs {
		if env == e {
			return true
		}
	}

	return false
}

func FetchDockerImageTag() []string {
	var dockerTags []string
	dockerRegistryURL := os.Getenv(environment.ENV_KEY_DOCKER_REGISTRY_URL)

	url := "https://" + dockerRegistryURL + "/artifactory/api/docker/ip-team-worldwire-docker-local/v2/gftn%2Fapi-service/tags/list"

	req, _ := http.NewRequest("GET", url, nil)

	userName := os.Getenv(environment.ENV_KEY_DOCKER_USER)
	pw := os.Getenv(environment.ENV_KEY_DOCKER_PW)

	req.SetBasicAuth(userName, pw)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	tags := &Tags{}
	json.Unmarshal(body, &tags)

	for _, t := range tags.Tags {
		dockerTags = append(dockerTags, t)
	}

	LOGGER.Debugf("%+v", dockerTags)
	return dockerTags
}

func ValidateDockerTag(tag string) bool {
	dockerTags := FetchDockerImageTag()

	for _, t := range dockerTags {
		if tag == t {
			return true
		}
	}

	return false
}

func CheckAWSSecretExist(env, domain, variable string) error {
	os.Setenv("AWS_REGION", "ap-southeast-1")
	LOGGER.Debugf("AWS_REGION: %s", os.Getenv("AWS_REGION"))

	for _, svc := range constant.ParticipantServices {
		var credential CredentialInfo

		if len(variable) == 0 {
			credential = CredentialInfo{
				Environment: env,
				Domain:      domain,
				Service:     svc,
			}
		} else {
			credential = CredentialInfo{
				Environment: env,
				Domain:      domain,
				Service:     svc,
				Variable:    variable,
			}
		}

		secretStr, getSecretErr := GetSecret(credential)
		if getSecretErr != nil {
			LOGGER.Warningf("Can not find secret in secret manager: %s", getSecretErr)
			return getSecretErr
		} else if len(secretStr) == 0 {
			LOGGER.Warningf("Secret is empty: secret=%s", secretStr)
			return errors.New("secret is empty")
		}

		os.Setenv(credential.Variable, "true")
	}

	return nil
}

func UpdateAWSSecret(environment, category, domain, version string, byteData []byte) error {
	// Secret file structure:
	// environment-variable{service-name{secrets}...}

	// Convert json into golang map
	c := make(map[string]interface{})
	json.Unmarshal(byteData, &c)

	// Iterate through environment variable
	for env, value := range c {
		// Choose the corresponding environment base on request payload
		if environment == env {
			secretType := make(map[string]interface{})

			jsonValue, _ := json.Marshal(value)
			json.Unmarshal(jsonValue, &secretType)

			for t, v := range secretType {
				if t == category {
					update(v, environment, domain, version)
					break
				}
			}
		}
	}

	return nil
}

func update(value interface{}, environment, domain, version string) error {
	svc := make(map[string]interface{})

	// Convert json into golang map
	jsonValue, _ := json.Marshal(value)
	json.Unmarshal(jsonValue, &svc)

	// Iterate through the service name
	for svcName, v := range svc {
		credential := CredentialInfo{
			Environment: environment,
			Domain:      domain,
			Service:     svcName,
			Variable:    "initialize",
		}

		secrets, _ := json.Marshal(v)

		secretContent := SecretContent{
			Entry:       nil,
			Description: "World Wire Service Secret: " + version,
			FilePath:    "",
			RawJson:     secrets,
		}

		err := UpdateSecret(credential, secretContent)
		if err != nil {
			LOGGER.Errorf("Encounter error while updating secret to AWS: %s", err.Error())
			return err
		}
	}

	return nil
}

func SetAWSSecret(role, environment, participantID string, byteData []byte) error {
	// Create secret: /environment/domain/key/initialize

	c := make(map[string]interface{})

	// Open the file that contained the secret for each services
	// Secret file structure: environment-variable{service-name{secrets}...}
	// Convert json into golang map
	err := json.Unmarshal(byteData, &c)
	if err != nil {
		LOGGER.Error(err.Error())
		return err
	}

	// Iterate through environment variable
	for env, value := range c {
		// Choose the corresponding environment base on request payload
		if environment == env {
			domain := make(map[string]interface{})

			// Convert json into golang map
			jsonValue, _ := json.Marshal(value)
			json.Unmarshal(jsonValue, &domain)

			for d, svcs := range domain {
				if d == constant.Participant {
					svc := make(map[string]interface{})

					// Convert json into golang map
					jsonValue, _ := json.Marshal(svcs)
					json.Unmarshal(jsonValue, &svc)

					// Iterate through the service name
					for svcName, v := range svc {
						if role == constant.AssetIssuerParticipant {
							if svcName == "ww-gateway" || svcName == "participant" {
								credential := CredentialInfo{
									Environment: environment,
									Domain:      participantID,
									Service:     svcName,
									Variable:    "initialize",
								}

								secrets, _ := json.Marshal(v)

								var finalSecretStr string
								var finalSecret []byte

								brokersURL := os.Getenv("BROKERS_URL")

								if d != constant.GlobalDomain {
									finalSecretStr = strings.Replace(string(secrets), "kafka_broker_url", brokersURL, -1)
									finalSecret = []byte(finalSecretStr)
								} else {
									finalSecret = secrets
								}

								secretContent := SecretContent{
									Entry:       nil,
									Description: "World Wire Service Secret",
									FilePath:    "",
									RawJson:     finalSecret,
								}

								err := CreateSecret(credential, secretContent)
								if err != nil {
									LOGGER.Errorf("Encounter error while storing secret to AWS: %s", err.Error())
									return err
								}
							}
						} else {
							credential := CredentialInfo{
								Environment: environment,
								Domain:      participantID,
								Service:     svcName,
								Variable:    "initialize",
							}

							secrets, _ := json.Marshal(v)

							var finalSecretStr string
							var finalSecret []byte

							brokersURL := os.Getenv("BROKERS_URL")

							if d != constant.GlobalDomain {
								finalSecretStr = strings.Replace(string(secrets), "kafka_broker_url", brokersURL, -1)
								finalSecret = []byte(finalSecretStr)
							} else {
								finalSecret = secrets
							}

							secretContent := SecretContent{
								Entry:       nil,
								Description: "World Wire Service Secret",
								FilePath:    "",
								RawJson:     finalSecret,
							}

							err := CreateSecret(credential, secretContent)
							if err != nil {
								LOGGER.Errorf("Encounter error while storing secret to AWS: %s", err.Error())
								return err
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func CreateIAMPolicy(participantName, serviceName, envName, orgId, secretManagerRegion string) error {
	dbRegion := os.Getenv(environment.ENV_KEY_DYNAMO_DB_REGION)
	policyBoundary := `/` + envName + `/` + participantName + `/` + serviceName + `/*`
	participantBoundary := `/` + envName + `/` + participantName + `/participant/*`
	username := envName + `_` + participantName + `_` + serviceName
	localServiceGroupName := envName + `_local_service`
	commonServiceGroupName := envName + `_service`
	tokenSecret := "/" + envName + "/ww/account/token*"
	accountSecret := "/" + envName + "/" + participantName + "/account/*"
	killSwitchSecret := "/" + envName + "/" + participantName + "/killswitch/accounts*"
	adminKillSwitchSecret := "/" + envName + "/*/killswitch/accounts*"

	var dat []byte
	var err error

	switch serviceName {
	case "api-service":
		dat, err = ioutil.ReadFile(constant.K8sBasePath + "/iam-policy/api-service.json")
	case "crypto-service":
		dat, err = ioutil.ReadFile(constant.K8sBasePath + "/iam-policy/crypto-service.json")
	case "send-service":
		dat, err = ioutil.ReadFile(constant.K8sBasePath + "/iam-policy/send-service.json")
	case "payment-service":
		dat, err = ioutil.ReadFile(constant.K8sBasePath + "/iam-policy/payment-service.json")
	case "ww-gateway":
		dat, err = ioutil.ReadFile(constant.K8sBasePath + "/iam-policy/ww-gateway.json")
	case "gas-service":
		dat, err = ioutil.ReadFile(constant.K8sBasePath + "/iam-policy/gas-service.json")
	case "admin-service":
		dat, err = ioutil.ReadFile(constant.K8sBasePath + "/iam-policy/admin-service.json")
	case "pr-service":
		dat, err = ioutil.ReadFile(constant.K8sBasePath + "/iam-policy/pr-service.json")
	case "payout-service":
		dat, err = ioutil.ReadFile(constant.K8sBasePath + "/iam-policy/payout-service.json")
	case "fee-service":
		dat, err = ioutil.ReadFile(constant.K8sBasePath + "/iam-policy/fee-service.json")
	case "anchor-service":
		dat, err = ioutil.ReadFile(constant.K8sBasePath + "/iam-policy/anchor-service.json")
	case "whitelist-service":
		dat, err = ioutil.ReadFile(constant.K8sBasePath + "/iam-policy/whitelist-service.json")
	case "quotes-service":
		dat, err = ioutil.ReadFile(constant.K8sBasePath + "/iam-policy/quotes-service.json")
	default:
		return errors.New("unknown service type")
	}

	if err != nil {
		LOGGER.Error(err.Error())
		return err
	}

	file := string(dat)
	file = strings.Replace(file, "<env>", envName, -1)
	file = strings.Replace(file, "<db_region>", dbRegion, -1)
	file = strings.Replace(file, "<region>", secretManagerRegion, -1)
	file = strings.Replace(file, "<accountId>", orgId, -1)
	file = strings.Replace(file, "<service_secret>", policyBoundary, -1)
	file = strings.Replace(file, "<participant_secret>", participantBoundary, -1)
	file = strings.Replace(file, "<account_secret>", accountSecret, -1)
	file = strings.Replace(file, "<token_secret>", tokenSecret, -1)
	file = strings.Replace(file, "<killswitch_secret>", killSwitchSecret, -1)
	file = strings.Replace(file, "<admin_killswitch_secret>", adminKillSwitchSecret, -1)
	file = strings.Replace(file, "<participantId>", participantName, -1)

	input := &iam.CreatePolicyInput{
		Description:    aws.String("policy for /" + envName + "/" + participantName + "/" + serviceName),
		PolicyDocument: aws.String(file),
		PolicyName:     aws.String(username + "_policy"),
	}

	svc := iam.New(session.New())
	output, err := svc.CreatePolicy(input)
	err = iamErrorHandler(err)
	if err != nil {
		LOGGER.Error("Failed to create IAM policy")
		return err
	}

	LOGGER.Info(output)
	LOGGER.Info("Policy created: " + policyBoundary)
	//========
	userInput := &iam.CreateUserInput{
		UserName: aws.String(username),
		//PermissionsBoundary: output.Policy.Arn,
	}

	result, err := svc.CreateUser(userInput)
	err = iamErrorHandler(err)
	if err != nil {
		LOGGER.Error("Failed to create IAM user")
		return err
	}

	LOGGER.Info(result)
	LOGGER.Info("Created user " + username)
	//========
	localServiceGroupInput := &iam.AddUserToGroupInput{
		GroupName: aws.String(localServiceGroupName),
		UserName:  aws.String(username),
	}

	addToLocalGroup, err := svc.AddUserToGroup(localServiceGroupInput)
	err = iamErrorHandler(err)
	if err != nil {
		LOGGER.Error("Failed to add IAM user to local service group")
		return err
	}

	LOGGER.Info(addToLocalGroup)
	LOGGER.Info("Add user to local service group " + localServiceGroupName)
	//========
	commonServiceGroupInput := &iam.AddUserToGroupInput{
		GroupName: aws.String(commonServiceGroupName),
		UserName:  aws.String(username),
	}

	addToCommonGroup, err := svc.AddUserToGroup(commonServiceGroupInput)
	err = iamErrorHandler(err)
	if err != nil {
		LOGGER.Error("Failed to add IAM user to common service group")
		return err
	}

	LOGGER.Info(addToCommonGroup)
	LOGGER.Info("Add user to common service group " + commonServiceGroupName)

	_, attachErr := svc.AttachUserPolicy(&iam.AttachUserPolicyInput{
		PolicyArn: output.Policy.Arn,
		UserName:  aws.String(username),
	})
	err = iamErrorHandler(attachErr)
	if err != nil {
		LOGGER.Error("Failed to create IAM user policy")
		return err
	}

	//======
	keyInput := &iam.CreateAccessKeyInput{
		UserName: aws.String(username),
	}

	accessKeyRes, err := svc.CreateAccessKey(keyInput)
	err = iamErrorHandler(err)
	if err != nil {
		LOGGER.Error("Failed to create access key")
		return err
	}

	script := constant.K8sBasePath + "/script/create_aws_iam_key.sh"
	err = RunBashCmd(script, participantName, serviceName, *accessKeyRes.AccessKey.AccessKeyId, *accessKeyRes.AccessKey.SecretAccessKey, envName, orgId)
	if err != nil {
		LOGGER.Errorf("Failed to create AWS IAM secret")
		return err
	}

	return nil
}

func iamErrorHandler(err error) error {
	if err == nil {
		return nil
	}
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case iam.ErrCodeLimitExceededException:
			LOGGER.Errorf("%s %s", iam.ErrCodeLimitExceededException, aerr.Error())
		case iam.ErrCodeEntityAlreadyExistsException:
			LOGGER.Errorf("%s %s", iam.ErrCodeEntityAlreadyExistsException, aerr.Error())
		case iam.ErrCodeNoSuchEntityException:
			LOGGER.Errorf("%s %s", iam.ErrCodeNoSuchEntityException, aerr.Error())
		case iam.ErrCodeInvalidInputException:
			LOGGER.Errorf("%s %s", iam.ErrCodeInvalidInputException, aerr.Error())
		case iam.ErrCodeConcurrentModificationException:
			LOGGER.Errorf("%s %s", iam.ErrCodeConcurrentModificationException, aerr.Error())
		case iam.ErrCodeServiceFailureException:
			LOGGER.Errorf("%s %s", iam.ErrCodeServiceFailureException, aerr.Error())
		default:
			LOGGER.Errorf("%s", aerr.Error())
		}

		return err
	} else {
		// Print the error, cast err to awserr.Error to get the Code and Message from an error.
		LOGGER.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func WriteParticipantToFireBase(input model.Automation) {
	LOGGER.Info("Writing participant information to FireBase")

	input.Status = []string{constant.StatusConfiguring}
	isInitialized := true
	input.Initialized = &isInitialized
	nodes := map[string]interface{}{*input.ParticipantID: input}
	// write to firebase the details of the node
	updateInstitutionErr := wwfirebase.FbRef.Child("participants/"+*input.InstitutionID+"/nodes/").Update(wwfirebase.AppContext, nodes) //.Push(wwfirebase.AppContext, &input)
	if updateInstitutionErr != nil {
		LOGGER.Errorf("Failed to update the institution in FireBase: %s", updateInstitutionErr.Error())
		return
	}

	institution := map[string]interface{}{*input.ParticipantID: *input.InstitutionID}
	// duplicate the write to firebase so that "institutionId" can be looked-up by
	// "participantId" and "env" to check within firebase security rules
	updateNodeErr := wwfirebase.FbRef.Child("nodes/").Update(wwfirebase.AppContext, institution) //.Push(wwfirebase.AppContext, *input.InstitutionID)
	if updateNodeErr != nil {
		LOGGER.Errorf("Failed to update the nodes in FireBase: %s", updateNodeErr.Error())
		return
	}

	return
}

func CheckParticipantRecordInFireBase(input model.Automation) ([]string) {
	wwfirebase.FbClient, _, _ = wwfirebase.AuthenticateWithAdminPrivileges()
	wwfirebase.FbRef = wwfirebase.GetRootRef()

	// Get the original data from fireBase
	var participantInfo map[string]interface{}
	ref := wwfirebase.FbClient.NewRef("participants/" + *input.InstitutionID + "/nodes/").Child(*input.ParticipantID)
	ref.Get(wwfirebase.AppContext, &participantInfo)

	var status []string
	if participantInfo == nil {
		LOGGER.Warning("No participant record found")
		WriteParticipantToFireBase(input)
		return nil
	} else {
		if len(participantInfo["status"].([]interface{})) == 1 {
			LOGGER.Debug("Status length: 1")
			if participantInfo["status"].([]interface{})[0] == constant.StatusConfiguring {
				LOGGER.Debug("status is: configuring")
				WriteParticipantToFireBase(input)
				return nil
			} else if participantInfo["status"].([]interface{})[0] == constant.StatusComplete {
				LOGGER.Debug("status is: complete")
				status = append(status, constant.StatusComplete)
				return status
			} else if participantInfo["status"].([]interface{})[0] == constant.StatusPending {
				LOGGER.Debug("status is: pending")
				WriteParticipantToFireBase(input)
				return nil
			} else {
				LOGGER.Debugf("status is: %s", participantInfo["status"].([]interface{})[0])
				UpdateStatusInFireBase(*input.InstitutionID, *input.ParticipantID, constant.StatusConfiguring, "configure")
				s := participantInfo["status"].([]interface{})[0]
				status = append(status, s.(string))
				return status
			}
		} else {
			LOGGER.Debug("More then one status")
			for _, s := range participantInfo["status"].([]interface{}) {
				if s.(string) == constant.StatusConfiguring {
					LOGGER.Debug("status is: configuring")
					continue
				} else if s.(string) == constant.StatusPending {
					LOGGER.Debug("status is: pending")
					UpdateStatusInFireBase(*input.InstitutionID, *input.ParticipantID, constant.StatusConfiguring, "configure")
				} else {
					status = append(status, s.(string))
				}
			}
		}
		return status
	}
}

func UpdateStatusInFireBase(iid, pid, status, result string) {
	LOGGER.Debugf("Update participant in FireBase")

	// Get the original data from fireBase
	var participantInfo map[string]interface{}
	ref := wwfirebase.FbClient.NewRef("participants/" + iid + "/nodes/").Child(pid)
	ref.Get(wwfirebase.AppContext, &participantInfo)

	var statuses []string
	for key := range participantInfo {
		if key == "status" {
			count := 0
			// find if the status was already in the firebase
			for _, s := range participantInfo[key].([]interface{}) {
				if s.(string) == status {
					if result == "failed" {
						// status exist in the firebase
						count += 1
					}
				} else {
					if s.(string) == constant.StatusPending ||
						s.(string) == constant.StatusConfiguring ||
						s.(string) == constant.StatusComplete ||
						s.(string) == constant.StatusConfigurationFailed {
						// skip "pending" and "configuring" and "complete" and "configuration_failed" status
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
				} else if result == "configure" {
					// Prepend "configuring"
					statuses = append([]string{status}, statuses...)
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

	// Update the log into FireBase
	ref.Update(wwfirebase.AppContext, participantInfo)
}

func GetSecretTemplate(env, version string) (string, error) {
	os.Setenv(global_environment.ENV_KEY_AWS_REGION, "ap-southeast-1")

	credential := CredentialInfo{
		Environment: env,
		Domain:      "aws",
		Service:     "secret",
		Variable:    version,
	}

	res, err := GetSecret(credential)
	if err != nil {
		LOGGER.Errorf("Cannot get the secret template for version %s: %s", version, err)
		return "", err
	}

	secretResult := map[string]string{}
	err = json.Unmarshal([]byte(res), &secretResult)
	if err != nil {
		errMsg := errors.New("error parsing secret object format from AWS")
		LOGGER.Errorf("%s", errMsg)
		return "", errMsg
	}

	for key := range secretResult {
		if key == "file" {
			return secretResult[key], nil
		}
	}

	return "", nil
}

func IsStrongHold(input model.Automation) bool {
	pId := *input.ParticipantID
	bicCode := *input.Bic
	lowerCaseParticipantId := strings.ToLower(pId)

	return bicCode == constant.StrongholdBICCode && strings.Contains(lowerCaseParticipantId, constant.StrongholdID)
}