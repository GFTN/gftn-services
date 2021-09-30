// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package api

import (
	"encoding/json"
	"strings"

	"github.com/GFTN/gftn-services/utility/vault/utils"
)

func ListSafes(session utils.Session) {
	url := "/PasswordVault/WebServices/PIMServices.svc/Safes"
	body, _ := session.Get(url, session.CyberArkLogonResult)

	var safeList utils.SafeList
	if err := json.Unmarshal([]byte(body), &safeList); err != nil {
		LOGGER.Errorf(err.Error())
	}
	for i := 0; i < len(safeList.List); i++ {
		LOGGER.Debugf("%v", safeList.List[i].SafeName)
	}

}

func AddSafe(session utils.Session, safeName string) {
	url := "/PasswordVault/WebServices/PIMServices.svc/Safes"

	payload := strings.NewReader("{\n  \"safe\":\n  {\n    \"SafeName\":\"" + safeName + "\",\n    \"Description\":\"Test, Application, CyberArk, API\",\n    \"OLACEnabled\":false,\n    \"ManagingCPM\":\"PasswordManager\",\n    \"NumberofVersionsRetention\":5,\n    \"NumberofDaysRetention\":7\n  }\n}")

	body, _ := session.Post(url, session.CyberArkLogonResult, payload)

	LOGGER.Debugf(string(body))
}

func ListSafeMember(session utils.Session, safeName string) {

	url := "/PasswordVault/WebServices/PIMServices.svc/Safes/" + safeName + "/Members"

	body, _ := session.Get(url, session.CyberArkLogonResult)

	LOGGER.Debugf("All the members in the %v", safeName)
	var memberList utils.SafeMemberList
	if err := json.Unmarshal([]byte(body), &memberList); err != nil {
		LOGGER.Error("%s", err)
	}
	for i := 0; i < len(memberList.Members); i++ {
		LOGGER.Debugf("%v", memberList.Members[i].UserName)
	}
}

func AddSafeMember(session utils.Session, safeName string, username string) {
	url := "/PasswordVault/WebServices/PIMServices.svc/Safes/" + safeName + "/Members"

	payload := strings.NewReader("{\r\n  \"member\":\r\n  {\r\n    \"MemberName\"        :\"" + username + "\",\r\n    \"SearchIn\"          :\"Vault\",\r\n    \"MembershipExpirationDate\"  :\"\",\r\n    \"Permissions\":\r\n    [\r\n      {\"Key\":\"UseAccounts\", \"Value\":true},\r\n      {\"Key\":\"RetrieveAccounts\", \"Value\":true},\r\n      {\"Key\":\"ListAccounts\", \"Value\":true},\r\n      {\"Key\":\"AddAccounts\", \"Value\":false},\r\n      {\"Key\":\"UpdateAccountContent\",\"Value\":false},\r\n      {\"Key\":\"UpdateAccountProperties\",\"Value\":false},\r\n      {\"Key\":\"InitiateCPMAccountManagementOperations\",\"Value\":false},\r\n      {\"Key\":\"SpecifyNextAccountContent\",\"Value\":false},\r\n      {\"Key\":\"RenameAccounts\", \"Value\":false},\r\n      {\"Key\":\"DeleteAccounts\", \"Value\":false},\r\n      {\"Key\":\"UnlockAccounts\", \"Value\":false},\r\n      {\"Key\":\"ManageSafe\", \"Value\":false},\r\n      {\"Key\":\"ManageSafeMembers\", \"Value\":false},\r\n      {\"Key\":\"BackupSafe\", \"Value\":false},\r\n      {\"Key\":\"ViewAuditLog\", \"Value\":true},\r\n      {\"Key\":\"ViewSafeMembers\", \"Value\":true},\r\n      {\"Key\":\"RequestsAuthorizationLevel\",\"Value\":0},\r\n      {\"Key\":\"AccessWithoutConfirmation\",\"Value\":false},\r\n      {\"Key\":\"CreateFolders\", \"Value\":false},\r\n      {\"Key\":\"DeleteFolders\", \"Value\":false},\r\n      {\"Key\":\"MoveAccountsAndFolders\",\"Value\":false}\r\n    ]\r\n  }\r\n}")

	body, _ := session.Post(url, session.CyberArkLogonResult, payload)

	LOGGER.Debugf(string(body))
}

func UpdateSafe(session utils.Session, safeName string) {
	url := "/PasswordVault/WebServices/PIMServices.svc/Safes/?SafeName=" + safeName

	payload := strings.NewReader("{\r\n  \"safe\":\r\n  {\r\n    \"SafeName\":\"P-LIN-ROOT-SSHKEYS\",\r\n    \"Description\":\"Test, Application, CyberArk, API, Update Safe Test\",\r\n    \"OLACEnabled\":false,\r\n    \"ManagingCPM\":\"PasswordManager\",\r\n    \"NumberOfVersionsRetention\":5,\r\n    \"NumberOfDaysRetention\":7\r\n  }\r\n}")

	body, _ := session.Put(url, session.CyberArkLogonResult, payload, "")

	LOGGER.Debugf(string(body))
}
