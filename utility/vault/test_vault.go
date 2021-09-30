// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package vault

import (
	"github.com/GFTN/gftn-services/utility/vault/api"
	"github.com/GFTN/gftn-services/utility/vault/auth"
	// "encoding/json"
	// "github.com/GFTN/gftn-services/utility/vault/utils"
	// "fmt"
	// "net/http"
	// "log"
	//"fmt"
)

func main() {

	var appId = "SSLcert"
	var safeName = "ApiSafe"
	//var newCredential = "test123"

	s := auth.GetSession()

	//create application and enable certificateSerialNumber
	api.AddApplication(s, appId)
	api.ListAuthentication(s, appId)
	api.AddAuthentication(s, appId)

	//create safe
	api.AddSafe(s, safeName)
	api.ListSafeMember(s, safeName)

	// for AIM operation
	// api.AddSafeMember(s, safeName, appId)
	// api.AddSafeMember(s, safeName, "Prov_EC2AMAZ-TQ8PDII")
	// api.AddSafeMember(s, safeName, "AIMWebService")
	// api.ListSafeMember(s, safeName)

	//add account in the safe
	api.AddAccount(s, safeName, "IBM_Token_ACC_PUBLIC", "ie.one.payments.worldwire.io", "ibm-account-public-key", "")
	api.AddAccount(s, safeName, "IBM_Token_ACC_PRIVATE", "ie.one.payments.worldwire.io", "ibm-account-private-key", "")

	//get account id inside the safe
	//accountId := api.GetAccount(s, safeName)
	//api.GetAccountGroup(s, safeName)

	//update password
	// auth.GetPasswordValue(s, accountId)
	// auth.RandomCredential(s, accountId)
	// auth.GetPasswordValue(s, accountId)
	// auth.SetCredential(s, accountId, newCredential)
	// auth.GetPasswordValue(s, accountId)

	//AIM
	// body := auth.GetPassword(appId, safeName, "IBM_Token_ACC_PRIVATE")
	// var secret utils.Secret
	// if err := json.Unmarshal([]byte(body), &secret); err != nil{
	// 	panic(err)
	// }

	//this will not be prininted since the output will be catched by eval() in env.sh
	// fmt.Println("IBM_PRIVATE_LABEL=", secret.Content)
	// auth.SetEnv()
	// utils.GetEnv()
	// fmt.Println(os.Getenv("test"))

}
