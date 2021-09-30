// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package helper

import (
	"encoding/json"
	"io/ioutil"
	"os"

	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
)

// Secrets : see /.credentials/{env}/secret.json
type Secrets struct {
	Fb_admin              string
	Send_in_blue_api_key  string
	Passport_secret       string
	Ibmid_client_id       string
	Ibmid_client_secret   string
	Ww_jwt_pepper_encoded string
}

// Envs : see /.credentials/{env}/env.json
type Envs struct {
	Totp_label              string
	Build                   string
	Build_for               string
	Api_port                string
	Site_root               string
	App_root                string
	Gae_service             string
	Ibmid_authorization_url string
	Ibmid_token_url         string
	Ibmid_issuer_id         string
	Enable_2fa              string
	Fb_database_url         string
	Refresh_mins            string
	Initial_mins 			string
}

// SetCustomEnvs : sets environment variables from ./credentials/*
func SetCustomEnvs(credentialsDir string) {

	secretsByte, err := ioutil.ReadFile("../" + credentialsDir + "/secret.json")
	var secrets Secrets
	if err != nil {
		panic(err)
	}
	json.Unmarshal(secretsByte, &secrets)
	// fmt.Print(string(secretsByte))

	envByte, err := ioutil.ReadFile("../" + credentialsDir + "/env.json")
	var envs Envs
	if err != nil {
		panic(err)
	}
	json.Unmarshal(envByte, &envs)
	// fmt.Print(string(envByte))

	os.Setenv(global_environment.ENV_KEY_FIREBASE_DB_URL, envs.Fb_database_url)
	os.Setenv(global_environment.ENV_KEY_WW_JWT_PEPPER_OBJ, secrets.Ww_jwt_pepper_encoded)
	os.Setenv(global_environment.ENV_KEY_FIREBASE_CREDENTIALS, secrets.Fb_admin)

}
