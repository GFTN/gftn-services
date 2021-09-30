// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
/*
Secret Naming Reminder:
	The secret name must be ASCII letters, digits, or the following characters: /_+=.@-
	Don't end your secret name with a hyphen followed by six characters. If you
	do so, you risk confusion and unexpected results when searching for a secret
	by partial ARN. This is because Secrets Manager automatically adds a hyphen
	and six random characters at the end of the ARN.
*/

package secret_manager

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/GFTN/gftn-services/utility/aws/golang/utility"
)

func errorHandler(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case secretsmanager.ErrCodeDecryptionFailure:
			LOGGER.Errorf("%s %s", secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
		case secretsmanager.ErrCodeInternalServiceError:
			LOGGER.Errorf("%s %s", secretsmanager.ErrCodeInternalServiceError, aerr.Error())
		case secretsmanager.ErrCodeInvalidParameterException:
			LOGGER.Errorf("%s %s", secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
		case secretsmanager.ErrCodeInvalidRequestException:
			LOGGER.Errorf("%s %s", secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
		case secretsmanager.ErrCodeResourceNotFoundException:
			LOGGER.Errorf("%s %s", secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
		case secretsmanager.ErrCodeLimitExceededException:
			LOGGER.Errorf("%s %s", secretsmanager.ErrCodeLimitExceededException, aerr.Error())
		case secretsmanager.ErrCodeEncryptionFailure:
			LOGGER.Errorf("%s %s", secretsmanager.ErrCodeEncryptionFailure, aerr.Error())
		case secretsmanager.ErrCodeResourceExistsException:
			LOGGER.Warningf("%s %s", secretsmanager.ErrCodeResourceExistsException, aerr.Error())
		case secretsmanager.ErrCodeMalformedPolicyDocumentException:
			LOGGER.Errorf("%s %s", secretsmanager.ErrCodeMalformedPolicyDocumentException, aerr.Error())
		case secretsmanager.ErrCodePreconditionNotMetException:
			LOGGER.Errorf("%s %s", secretsmanager.ErrCodePreconditionNotMetException, aerr.Error())
		default:
			LOGGER.Errorf("%s", aerr.Error())
		}
	} else {
		// Print the error, cast err to awserr.Error to get the Code and Message from an error.
		LOGGER.Errorf("%s", err.Error())
	}
}

/*
	IAM user required permission to call GetSecret function:
		* secretsmanager:GetSecretValue

		* kms:Decrypt - required only if you use a customer-managed AWS KMS key
		to encrypt the secret. You do not need this permission to use the account's
		default AWS managed CMK for Secrets Manager.
*/

func GetSecret(credentialInfo utility.CredentialInfo) (string, error) {
	secretId, err := utility.GetCredentialId(credentialInfo)
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

func GetSingleSecretEntry(credential utility.CredentialInfo, entryName string) (string, error) {
	res, err := GetSecret(credential)
	if err != nil {
		LOGGER.Errorf("Cannot get the specified environment variable: %s", err)
		return "", err
	}

	// unmarshal to map so that we can retrieve the value of the object
	secretResult := map[string]string{}
	err = json.Unmarshal([]byte(res), &secretResult)
	if err != nil {
		errorMsg := "Error parsing secret object format from AWS"
		LOGGER.Errorf(errorMsg)
		return "", errors.New(errorMsg)
	}

	// if cannot find killswitch string from aws, error out
	if _, ok := secretResult[entryName]; !ok {
		errMsg := "Cannot find corresponding killswitch string"
		LOGGER.Errorf(errMsg)
		return "", errors.New(errMsg)
	}
	return secretResult[entryName], nil
}

/*
	IAM user required permission to call CreateSecret function:
		* secretsmanager:CreateSecret

		* kms:GenerateDataKey - needed only if you use a customer-managed AWS
		KMS key to encrypt the secret. You do not need this permission to use
		the account's default AWS managed CMK for Secrets Manager.

		* kms:Decrypt - needed only if you use a customer-managed AWS KMS key
		to encrypt the secret. You do not need this permission to use the account's
		default AWS managed CMK for Secrets Manager.

		* secretsmanager:TagResource - needed only if you include the Tags parameter.
*/

func CreateSecret(credentialInfo utility.CredentialInfo, secretContent utility.SecretContent) error {

	secretId, err := utility.GetCredentialId(credentialInfo)
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

/*
	IAM user required permission to call DeleteSecret function:
		* secretsmanager:DeleteSecret
	Note:
		recoveryDays should be 7 days at minimum
*/
func DeleteSecret(credentialInfo utility.CredentialInfo, recoveryDays int64) error {
	secretId, err := utility.GetCredentialId(credentialInfo)
	if err != nil {
		return err
	}
	LOGGER.Infof("deleteSecret: %s", secretId)

	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New())
	input := &secretsmanager.DeleteSecretInput{
		RecoveryWindowInDays: aws.Int64(recoveryDays),
		SecretId:             aws.String(secretId),
	}

	_, err = svc.DeleteSecret(input)
	if err != nil {
		errorHandler(err)
		return err
	}
	LOGGER.Infof("deleteSecret: Success!")
	return nil
}

func DeleteSingleSecretEntry(credentialInfo utility.CredentialInfo, entryName string) error {

	LOGGER.Infof("DeleteSingleSecretEntry: %s, entry name: %s", credentialInfo, entryName)

	res, err := GetSecret(credentialInfo)
	if err != nil {
		LOGGER.Errorf("Cannot get the specified environment variable: %s", err)
		return err
	}

	// unmarshal to map so that we can retrieve the value of the object
	secretResult := map[string]string{}
	err = json.Unmarshal([]byte(res), &secretResult)
	if err != nil {
		errorMsg := "Error parsing secret object format from AWS"
		LOGGER.Errorf(errorMsg)
		return errors.New(errorMsg)
	}

	var sc = utility.SecretContent{}
	for key, val := range secretResult {
		if key == entryName {
			continue
		}
		sc.Entry = append(sc.Entry, utility.SecretEntry{Key: key, Value: val})
	}

	err = UpdateSecret(credentialInfo, sc)
	if err != nil {
		errMsg := "Error while updating secret in the deleteSingeSecretEntry process"
		LOGGER.Errorf(errMsg)
		return errors.New(errMsg)
	}
	LOGGER.Infof("DeleteSingleSecretEntry: Success!")
	return nil
}

/*
	IAM user required permission to call UpdateSecret function:
		* secretsmanager:UpdateSecret

		* kms:GenerateDataKey - needed only if you use a custom AWS KMS key to
		encrypt the secret. You do not need this permission to use the account's
		AWS managed CMK for Secrets Manager.

		* kms:Decrypt - needed only if you use a custom AWS KMS key to encrypt
		the secret. You do not need this permission to use the account's AWS managed
		CMK for Secrets Manager.
*/

func UpdateSecret(credentialInfo utility.CredentialInfo, updatedSecret utility.SecretContent) error {
	secretId, err := utility.GetCredentialId(credentialInfo)
	if err != nil {
		return err
	}
	LOGGER.Infof("updateSecret: %s", secretId)

	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New())

	secretString, err := getSecretString(updatedSecret)
	if err != nil {
		return err
	}
	input := &secretsmanager.UpdateSecretInput{
		Description:  aws.String(updatedSecret.Description),
		SecretId:     aws.String(secretId),
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

/*
	IAM user required permission to call AppendSecret function:
		* secretsmanager:PutSecretValue

		* kms:GenerateDataKey - needed only if you use a customer-managed AWS
		KMS key to encrypt the secret. You do not need this permission to use
		the account's default AWS managed CMK for Secrets Manager.
*/

func AppendSecret(credentialInfo utility.CredentialInfo, updatedSecret utility.SecretContent) error {
	secretId, err := utility.GetCredentialId(credentialInfo)
	if err != nil {
		return err
	}
	LOGGER.Infof("appendSecret: %s", secretId)

	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New())

	secretString, err := getSecretString(updatedSecret)
	if err != nil {
		return err
	}

	res, err := GetSecret(credentialInfo)
	if err != nil {
		return err
	}

	appendData := map[string]string{}
	err = json.Unmarshal([]byte(secretString), &appendData)
	if err != nil {
		errMsg := errors.New("Error parsing secret object format from AWS")
		LOGGER.Errorf("%s", errMsg)
		return errMsg
	}

	remoteData := map[string]string{}
	err = json.Unmarshal([]byte(res), &remoteData)
	if err != nil {
		errMsg := errors.New("Error parsing secret object format from AWS")
		LOGGER.Errorf("%s", errMsg)
		return errMsg
	}

	for key, val := range appendData {
		remoteData[key] = val
	}

	finalResult, _ := json.Marshal(remoteData)

	input := &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(secretId),
		SecretString: aws.String(string(finalResult)),
	}

	_, err = svc.PutSecretValue(input)
	if err != nil {
		errorHandler(err)
		return err
	}
	LOGGER.Infof("appendSecret: Success!")
	return nil
}

func getSecretString(secretContent utility.SecretContent) (string, error) {
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
