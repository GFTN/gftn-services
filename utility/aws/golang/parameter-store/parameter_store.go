// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
/*
Parameter Naming Constraints:
	* Parameter names are case sensitive.
	* A parameter name must be unique within an AWS Region
	* A parameter name can't be prefixed with "aws" or "ssm" (case-insensitive).
	* Parameter names can include only the following symbols and letters: a-zA-Z0-9_.-/
	* A parameter name can't include spaces.
	* Parameter hierarchies are limited to a maximum depth of fifteen levels.
*/
package parameter_store

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/GFTN/gftn-services/utility/aws/golang/utility"
)

func errorHandler(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case ssm.ErrCodeInternalServerError:
			LOGGER.Errorf("%s %s", ssm.ErrCodeInternalServerError, aerr.Error())
		case ssm.ErrCodeInvalidKeyId:
			LOGGER.Errorf("%s %s", ssm.ErrCodeInvalidKeyId, aerr.Error())
		case ssm.ErrCodeParameterNotFound:
			LOGGER.Errorf("%s %s", ssm.ErrCodeParameterNotFound, aerr.Error())
		case ssm.ErrCodeParameterVersionNotFound:
			LOGGER.Errorf("%s %s", ssm.ErrCodeParameterVersionNotFound, aerr.Error())
		case ssm.ErrCodeParameterLimitExceeded:
			LOGGER.Errorf("%s %s", ssm.ErrCodeParameterLimitExceeded, aerr.Error())
		case ssm.ErrCodeTooManyUpdates:
			LOGGER.Errorf("%s %s", ssm.ErrCodeTooManyUpdates, aerr.Error())
		case ssm.ErrCodeParameterAlreadyExists:
			LOGGER.Errorf("%s %s", ssm.ErrCodeParameterAlreadyExists, aerr.Error())
		case ssm.ErrCodeHierarchyLevelLimitExceededException:
			LOGGER.Errorf("%s %s", ssm.ErrCodeHierarchyLevelLimitExceededException, aerr.Error())
		case ssm.ErrCodeHierarchyTypeMismatchException:
			LOGGER.Errorf("%s %s", ssm.ErrCodeHierarchyTypeMismatchException, aerr.Error())
		case ssm.ErrCodeInvalidAllowedPatternException:
			LOGGER.Errorf("%s %s", ssm.ErrCodeInvalidAllowedPatternException, aerr.Error())
		case ssm.ErrCodeParameterMaxVersionLimitExceeded:
			LOGGER.Errorf("%s %s", ssm.ErrCodeParameterMaxVersionLimitExceeded, aerr.Error())
		case ssm.ErrCodeParameterPatternMismatchException:
			LOGGER.Errorf("%s %s", ssm.ErrCodeParameterPatternMismatchException, aerr.Error())
		case ssm.ErrCodeUnsupportedParameterType:
			LOGGER.Errorf("%s %s", ssm.ErrCodeUnsupportedParameterType, aerr.Error())
		default:
			LOGGER.Errorf("%s", aerr.Error())
		}
	} else {
		// Print the error, cast err to awserr.Error to get the Code and Message from an error.
		LOGGER.Errorf("%s", err.Error())
	}
}

func GetParameter(credentialInfo utility.CredentialInfo) (string, error) {
	parameterId, err := utility.GetCredentialId(credentialInfo)
	if err != nil {
		return "", err
	}
	LOGGER.Infof("getParameter: %s", parameterId)

	//Create a Parameters Manager client
	svc := ssm.New(session.New())

	// by default, we set decryption always for true
	decryption := true
	input := &ssm.GetParameterInput{
		Name:           aws.String(parameterId),
		WithDecryption: &decryption, // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetParameter(input)
	if err != nil {
		errorHandler(err)
		return "", err
	}
	LOGGER.Infof("getParameter: Success!")
	return *result.Parameter.Value, nil
}

func CreateParameter(credentialInfo utility.CredentialInfo, parameterContent utility.ParameterContent) error {
	return putParameter(credentialInfo, parameterContent, false)
}

func UpdateParameter(credentialInfo utility.CredentialInfo, updatedParameter utility.ParameterContent) error {
	return putParameter(credentialInfo, updatedParameter, true)
}

func DeleteParameter(credentialInfo utility.CredentialInfo) error {
	parameterId, err := utility.GetCredentialId(credentialInfo)
	if err != nil {
		return err
	}
	LOGGER.Infof("deleteParameter: %s", parameterId)

	svc := ssm.New(session.New())
	input := &ssm.DeleteParameterInput{
		Name: aws.String(parameterId),
	}

	_, err = svc.DeleteParameter(input)
	if err != nil {
		errorHandler(err)
		return err
	}
	LOGGER.Infof("deleteParameter: Success!")
	return nil
}

func putParameter(credentialInfo utility.CredentialInfo, parameterContent utility.ParameterContent, overwrite bool) error {
	parameterId, err := utility.GetCredentialId(credentialInfo)
	if err != nil {
		return err
	}

	if overwrite {
		LOGGER.Infof("updateParameter: %s", parameterId)
	} else {
		LOGGER.Infof("createParameter: %s", parameterId)
	}
	//Create a Parameters Manager client
	svc := ssm.New(session.New())

	input := &ssm.PutParameterInput{
		Description: aws.String(parameterContent.Description),
		Name:        aws.String(parameterId),
		Value:       aws.String(parameterContent.Value),
		Type:        aws.String("SecureString"),
		Overwrite:   &overwrite,
	}

	_, err = svc.PutParameter(input)
	if err != nil {
		errorHandler(err)
		return err
	}
	LOGGER.Infof("putParameter: Success!")
	return nil
}
