// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package main

import (
	b64 "encoding/base64"
	"fmt"
	"os"

	"github.com/GFTN/gftn-services/utility/common"

	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/payment/utils/test/functions"
)

func setEnvVariables() {
	os.Setenv(global_environment.ENV_KEY_ENVIRONMENT_VERSION, "eksqa")
	os.Setenv(global_environment.ENV_KEY_IBM_TOKEN_DOMAIN_ID, "ww")
	os.Setenv(global_environment.ENV_KEY_SECRET_STORAGE_LOCATION, common.AWS_SECRET)
	os.Setenv(global_environment.ENV_KEY_AWS_ACCESS_KEY_ID, "")
	os.Setenv(global_environment.ENV_KEY_AWS_SECRET_ACCESS_KEY, "")
	os.Setenv(global_environment.ENV_KEY_AWS_REGION, "ap-southeast-1")

}
func main() {

	setEnvVariables()
	data := functions.Sign()
	sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	fmt.Println(sEnc)

	//functions.Transfer()
}
