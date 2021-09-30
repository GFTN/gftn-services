// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package services

import (
	"os"
	"strings"

	"github.com/GFTN/gftn-services/utility/common"

	secret_manager "github.com/GFTN/gftn-services/utility/aws/golang/secret-manager"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
)

func InitEnv() {
	if strings.ToUpper(os.Getenv(global_environment.ENV_KEY_SECRET_STORAGE_LOCATION)) == common.AWS_SECRET {
		secret_manager.InitEnv()
	}
}
