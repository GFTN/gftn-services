// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package services

import (
	"os"

	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
)

func VariableCheck() {
	domainId := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	svcName := os.Getenv(global_environment.ENV_KEY_SERVICE_NAME)
	envVersion := os.Getenv(global_environment.ENV_KEY_ENVIRONMENT_VERSION)
	secretLocation := os.Getenv(global_environment.ENV_KEY_SECRET_STORAGE_LOCATION)

	if domainId == "" || svcName == "" || envVersion == "" || secretLocation == "" {
		panic("Initializing failed, Require the following environment variables to start up the service.\nHOME_DOMAIN_NAME: " + domainId + "\nSERVICE_NAME: " + svcName + "\nENVIRONMENT_VERSION: " + envVersion + "\nSECRET_STORAGE_LOCATION: " + secretLocation)
	}
}
