// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package utility

import (
	"encoding/base64"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"io/ioutil"
	"os"
)

type CredentialInfo struct {
	Environment string
	Domain      string
	Service     string
	Variable    string
}

type SecretContent struct {
	Entry       []SecretEntry
	Description string
	FilePath    string
	RawJson     []byte
}

type SecretEntry struct {
	Key   string
	Value string
}

func errorHandler(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case secretsmanager.ErrCodeDecryptionFailure:
			LOGGER.Warningf("%s %s", secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
		case secretsmanager.ErrCodeInternalServiceError:
			LOGGER.Warningf("%s %s", secretsmanager.ErrCodeInternalServiceError, aerr.Error())
		case secretsmanager.ErrCodeInvalidParameterException:
			LOGGER.Warningf("%s %s", secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
		case secretsmanager.ErrCodeInvalidRequestException:
			LOGGER.Warningf("%s %s", secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
		case secretsmanager.ErrCodeResourceNotFoundException:
			LOGGER.Warningf("%s %s", secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
		case secretsmanager.ErrCodeLimitExceededException:
			LOGGER.Warningf("%s %s", secretsmanager.ErrCodeLimitExceededException, aerr.Error())
		case secretsmanager.ErrCodeEncryptionFailure:
			LOGGER.Warningf("%s %s", secretsmanager.ErrCodeEncryptionFailure, aerr.Error())
		case secretsmanager.ErrCodeResourceExistsException:
			LOGGER.Warningf("%s %s", secretsmanager.ErrCodeResourceExistsException, aerr.Error())
		case secretsmanager.ErrCodeMalformedPolicyDocumentException:
			LOGGER.Warningf("%s %s", secretsmanager.ErrCodeMalformedPolicyDocumentException, aerr.Error())
		case secretsmanager.ErrCodePreconditionNotMetException:
			LOGGER.Warningf("%s %s", secretsmanager.ErrCodePreconditionNotMetException, aerr.Error())
		default:
			LOGGER.Errorf("%s", aerr.Error())
		}
	} else {
		// Print the error, cast err to awserr.Error to get the Code and Message from an error.
		LOGGER.Errorf("%s", err.Error())
	}
}

func GetCredentialId(credentialInfo CredentialInfo) (string, error) {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" || os.Getenv("AWS_REGION") == "" {
		LOGGER.Warningf("Cannot fetch the correct AWS session config, please check that you have set access key ID/secret key/region correctly")
		return "", errors.New("Cannot fetch the correct AWS session config, please check that you have set access key ID/secret key/region correctly")
	}
	if credentialInfo.Environment == "" || credentialInfo.Domain == "" || credentialInfo.Service == "" || credentialInfo.Variable == "" {
		LOGGER.Errorf("Some parameters are missing in the credential info")
		return "", errors.New("Some parameters are missing in the credential info")
	}
	return "/" + credentialInfo.Environment + "/" + credentialInfo.Domain + "/" + credentialInfo.Service + "/" + credentialInfo.Variable, nil
}

func GetSecret(credentialInfo CredentialInfo) (string, error) {
	secretId, err := GetCredentialId(credentialInfo)
	if err != nil {
		return "", err
	}
	LOGGER.Infof("getSecret: %s", secretId)

	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New())

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretId),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		errorHandler(err)
		return "", err
	}

	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString string
	if result.SecretString != nil {
		secretString = *result.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			LOGGER.Errorf("Base64 Decode Error: %s", err)
			return "", err
		}
		secretString = string(decodedBinarySecretBytes[:len])
	}
	LOGGER.Infof("getSecret: Success!")
	return secretString, nil
}

func CreateSecret(credentialInfo CredentialInfo, secretContent SecretContent) error {

	secretId, err := GetCredentialId(credentialInfo)
	if err != nil {
		return err
	}
	LOGGER.Infof("createSecret: %s", secretId)

	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New())

	secretString, err := getSecretString(secretContent)
	if err != nil {
		return err
	}

	input := &secretsmanager.CreateSecretInput{
		// ClientRequestToken: aws.String("EXAMPLE1-90ab-cdef-fedc-ba987SECRET1"),
		Description:  aws.String(secretContent.Description),
		Name:         aws.String(secretId),
		SecretString: aws.String(secretString),
	}

	_, err = svc.CreateSecret(input)
	if err != nil {
		errorHandler(err)
		return err
	}
	LOGGER.Infof("createSecret: Success!")
	return nil
}

func UpdateSecret(credentialInfo CredentialInfo, secretContent SecretContent) error {

	secretId, err := GetCredentialId(credentialInfo)
	if err != nil {
		return err
	}
	LOGGER.Infof("updateSecret: %s", secretId)

	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New())

	secretString, err := getSecretString(secretContent)
	if err != nil {
		return err
	}

	input := &secretsmanager.UpdateSecretInput{
		Description:  aws.String(secretContent.Description),
		SecretId:         aws.String(secretId),
		SecretString: aws.String(secretString),
	}

	_, err = svc.UpdateSecret(input)
	if err != nil {
		errorHandler(err)
		return err
	}
	LOGGER.Infof("updateSecret: Success!")
	return nil
}

func getSecretString(secretContent SecretContent) (string, error) {
	if secretContent.FilePath != "" {
		temp, err := ioutil.ReadFile(secretContent.FilePath)
		if err != nil {
			LOGGER.Errorf("Error encountered when reading the file: %s", err)
			return "", err
		}
		return string(temp), nil
	} else if len(secretContent.RawJson) > 0 {
		return string(secretContent.RawJson), nil
	} else if len(secretContent.Entry) > 0 {
		var payload string
		payload += "{"
		for key, subEntry := range secretContent.Entry {
			postfix := ","
			if key == len(secretContent.Entry)-1 {
				postfix = ""
			}
			payload += "\"" + subEntry.Key + "\":\"" + subEntry.Value + "\"" + postfix
		}
		payload += "}"
		return payload, nil
	} else {
		LOGGER.Errorf("No secret value specified")
		return "", errors.New("No secret value specified")
	}
	return "", errors.New("Error when parsing secret value")
}