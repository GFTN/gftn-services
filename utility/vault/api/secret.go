// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/GFTN/gftn-services/utility/vault/utils"
)

func GetPasswordWithAim(session utils.Session, safeName string, objectName string) (string, error) {
	objectName = strings.ToUpper(objectName)
	url := "/AIMWebService/api/Accounts?AppID=SSLcert&Safe=" + safeName + "&Folder=Root&Object=" + objectName
	resp, err := session.Get(url, "")
	var secret utils.Secret
	if err := json.Unmarshal([]byte(resp), &secret); err != nil {
		LOGGER.Error("%s", err)
	}
	return secret.Content, err
}

func GetPasswordValue(session utils.Session, safeName string, keyword string) (string, error) {
	accountId, err := GetAccount(session, safeName, keyword)

	if accountId == "null" {
		LOGGER.Warningf("Data " + keyword + " not exists in safe " + safeName)
		return "null", errors.New("Data not exists on vault")
	}

	url := "/PasswordVault/API/Accounts/" + accountId + "/Password/Retrieve"
	payload := strings.NewReader("{\n\t\"Reason\":\"Automatically retrieved password by " + safeName + "\",\n\t}")
	body, err := session.Post(url, session.CyberArkLogonResult, payload)
	if err != nil {
		LOGGER.Errorf("Something goes wrong when retrieving the credential in the vault: %s", err)
	} else {
		LOGGER.Infof("%s safe credential: %s successfully retrieved!", safeName, keyword)
	}
	result := strings.Replace(string(body), "\"", "", -1)
	return result, err
}

func RandomCredential(session utils.Session, accountId string) {
	url := "/PasswordVault/WebServices/PIMServices.svc/Accounts/" + accountId + "/ChangeCredentials"

	payload := strings.NewReader("{\n  \"ChangeCredsForGroup\":\"No\"\n}")

	body, _ := session.Put(url, session.CyberArkLogonResult, payload, "ImmediateChangeByCPM")

	fmt.Println(string(body))
}

func SetCredential(session utils.Session, safeName string, accountName string, newCredential string) error {
	accountId, err := GetAccount(session, safeName, accountName)
	url := "/PasswordVault/API/Accounts/" + accountId + "/Password/Update"
	payload := strings.NewReader("{\n\t\"ChangeCredsForGroup\":\"true\",\n\t\"NewCredentials\":\"" + newCredential + "\",\n\t\"AutoGenerate\":\"false\"}")
	_, err = session.Post(url, session.CyberArkLogonResult, payload)
	if err != nil {
		LOGGER.Errorf("Something goes wrong when updating the credential in the vault: %s", err)
	} else {
		LOGGER.Infof("%s safe credential: %s successfully updated!", safeName, accountName)
	}
	return err
}

func StoreSecret(session utils.Session, safeName string, objectName string, filePath string) error {
	body, err := GetPasswordWithAim(session, safeName, objectName)
	if err != nil {
		LOGGER.Errorf("%s", err)
	}
	err = ioutil.WriteFile(filePath, []byte(body), 0440)
	if err != nil {
		LOGGER.Errorf("%s", err)
	}
	LOGGER.Infof("Storing secret file as %s", filePath)
	return err
}
