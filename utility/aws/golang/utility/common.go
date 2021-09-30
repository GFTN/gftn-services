// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package utility

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/GFTN/gftn-services/utility/global-environment"
)

func GetCredentialId(credentialInfo CredentialInfo) (string, error) {
	if os.Getenv(global_environment.ENV_KEY_AWS_ACCESS_KEY_ID) == "" || os.Getenv(global_environment.ENV_KEY_AWS_SECRET_ACCESS_KEY) == "" || os.Getenv(global_environment.ENV_KEY_AWS_REGION) == "" {
		LOGGER.Errorf("Cannot fetch the correct AWS session config, please check that you have set access key ID/secret key/region correctly")
		return "", errors.New("Cannot fetch the correct AWS session config, please check that you have set access key ID/secret key/region correctly")
	}
	if credentialInfo.Environment == "" || credentialInfo.Domain == "" || credentialInfo.Service == "" || credentialInfo.Variable == "" {
		LOGGER.Errorf("Some parameters are missing in the credential info")
		return "", errors.New("Some parameters are missing in the credential info")
	}
	return "/" + credentialInfo.Environment + "/" + credentialInfo.Domain + "/" + credentialInfo.Service + "/" + credentialInfo.Variable, nil
}

func StoreSecret(filePath string, payload string) error {
	err := ioutil.WriteFile(filePath, []byte(payload), 0440)
	if err != nil {
		LOGGER.Errorf("%s", err)
	}
	LOGGER.Infof("Storing secret file as %s", filePath)
	return err
}
