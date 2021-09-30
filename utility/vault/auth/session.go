// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package auth

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/op/go-logging"
	"github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/vault/api"
	"github.com/GFTN/gftn-services/utility/vault/utils"
)

var LOGGER = logging.MustGetLogger("vault")

//singleton vault session object
var VaultSession utils.Session

func GetSession() (utils.Session, error) {

	session := utils.Session{}
	if VaultSession == session {
		//Create session once, should only be calling during init of a service

		session.BaseURL = os.Getenv(global_environment.ENV_KEY_VAULT_BASE_URL)
		session.CertPath = os.Getenv(global_environment.ENV_KEY_VAULT_CERT)
		session.KeyPath = os.Getenv(global_environment.ENV_KEY_VAULT_CERT_PRIVATE_KEY)

		var username, password string
		var ch1 = make(chan error)
		var ch2 = make(chan error)
		go func(ch chan<- error) {
			var err error
			username, err = api.GetPasswordWithAim(session, "IBM", "ADMIN_ACCOUNT")
			ch <- err
		}(ch1)

		go func(ch chan<- error) {
			var err error
			password, err = api.GetPasswordWithAim(session, "IBM", "ADMIN_PASSWORD")
			ch <- err
		}(ch2)

		if <-ch1 != nil || <-ch2 != nil {
			return utils.Session{}, errors.New("Error while fetching admin account")
		}

		url := session.BaseURL + "/PasswordVault/WebServices/auth/Cyberark/CyberArkAuthenticationService.svc/Logon"
		var jsonStr = []byte(`{"username":"` + username + `", "password":"` + password + `", "connectionNumber":2}`)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		if err != nil {
			LOGGER.Errorf("Something goes wrong when creating a new session token request: %s", err)
			return utils.Session{}, err
		}

		req.Header.Set("Content-Type", "application/json")

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{Transport: tr}

		res, err := client.Do(req)

		if err != nil {
			LOGGER.Errorf("Something goes wrong when fetching session token from the vault: %s", err)
			return utils.Session{}, err
		}

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		s := utils.Session{}
		err = json.Unmarshal([]byte(body), &s)
		if err != nil {
			LOGGER.Errorf("Something goes wrong when parsing session token: %s", err)
			return utils.Session{}, err
		}

		s.CertPath = session.CertPath
		s.KeyPath = session.KeyPath
		s.BaseURL = session.BaseURL
		LOGGER.Debugf("vault session: %v, %v, %v, %v", s.BaseURL, s.KeyPath, s.CertPath, s.CyberArkLogonResult)

		VaultSession = s
		return VaultSession, nil
	}
	//If service is initialized return singleton
	return VaultSession, nil
}
