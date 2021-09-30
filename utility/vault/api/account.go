// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package api

import (
	"encoding/json"
	"fmt"
	"strings"

	logging "github.com/op/go-logging"
	"github.com/GFTN/gftn-services/utility/vault/utils"
)

var LOGGER = logging.MustGetLogger("vault")

func AddAccount(session utils.Session, safeName string, accountName string, password string, groupName string) error {
	accountName = strings.ToUpper(accountName)
	url := "/PasswordVault/WebServices/PIMServices.svc/Account"
	LOGGER.Infof("Storing " + accountName + " to the vault...")
	payload := strings.NewReader("{\r\n  \"account\" : {\r\n    \"safe\":\"" + safeName + "\",\r\n    \"platformID\":\"WinDesktopLocal\",\r\n    \"address\":\"https://worldwire.io\",\r\n    \"accountName\":\"" + accountName + "\",\r\n    \"password\":\"" + password + "\",\r\n    \"username\":\"" + accountName + "\",\r\n    \"disableAutoMgmt\":\"false\",\r\n    \"disableAutoMgmtReason\":\"N/A\",\r\n    \"groupName\":\"\",\r\n    \"" + groupName + "\":\"\",\r\n}")

	_, err := session.Post(url, session.CyberArkLogonResult, payload)
	if err != nil {
		LOGGER.Errorf("Error while adding account into the vault")
	}
	return err

}

func GetAccount(session utils.Session, safeName string, keyword string) (string, error) {
	keyword = strings.ToUpper(keyword)
	url := "/PasswordVault/WebServices/PIMServices.svc/Accounts?Safe=" + safeName + "&Keywords=" + keyword

	body, err := session.Get(url, session.CyberArkLogonResult)
	var acc utils.Accounts
	if err := json.Unmarshal([]byte(body), &acc); err != nil {
		LOGGER.Errorf("Error while parsing account data from vault: %s", err)
		return "", err
	}

	if len(acc.Accounts) <= 0 {
		LOGGER.Warningf("No matching keyword:%s found in safe %s", keyword, safeName)
		return "null", nil
	}
	LOGGER.Infof("Successfully retrieve %s in safe: %s", keyword, safeName)
	return acc.Accounts[0].AccountID, err

}

func GetAccountGroup(session utils.Session, safeName string) {
	url := "/PasswordVault/API/AccountGroups?Safe=" + safeName

	body, _ := session.Get(url, session.CyberArkLogonResult)
	fmt.Println(string(body))

}
