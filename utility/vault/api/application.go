// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package api

import (
	"fmt"
	"strings"

	"github.com/GFTN/gftn-services/utility/vault/utils"
)

func AddApplication(session utils.Session, appId string) {
	url := "/PasswordVault/WebServices/PIMServices.svc/Applications/"

	payload := strings.NewReader("{\r\n  \"application\":\r\n  {\r\n    \"AppID\":\"RESTExamples2\",\r\n    \"Description\":\"Testing DevOps Deployments with CyberArk\",\r\n    \"Location\":\"/Applications\",\r\n    \"AccessPermittedFrom\":0,\r\n    \"AccessPermittedTo\":23,\r\n    \"ExpirationDate\":\"\",\r\n    \"Disabled\":\"No\",\r\n    \"BusinessOwnerFName\":\"John\",\r\n    \"BusinessOwnerLName\":\"Doe\",\r\n    \"BusinessOwnerEmail\":\"John.Doe@CyberArk.com\",\r\n    \"BusinessOwnerPhone\":\"555-555-1212\"\r\n  }\r\n}")

	body, _ := session.Post(url, session.CyberArkLogonResult, payload)

	fmt.Println(string(body))
}

func ListAuthentication(session utils.Session, appId string) {
	url := "/PasswordVault/WebServices/PIMServices.svc/Applications/" + appId + "/Authentications"
	body, _ := session.Get(url, session.CyberArkLogonResult)
	fmt.Println(string(body))
}

func AddAuthentication(session utils.Session, appId string) {
	url := "/PasswordVault/WebServices/PIMServices.svc/Applications/" + appId + "/Authentications"
	//only path/ hostuser / hash?
	payload := strings.NewReader("{\r\n  \"authentication\":\r\n  {\r\n    \"AuthType\":\"certificateSerialNumber\",\r\n    \"AuthValue\":\"1234555\",\r\n}")
	body, _ := session.Post(url, session.CyberArkLogonResult, payload)
	fmt.Println(string(body))
}
