// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package main

import (
	"fmt"

	"github.com/GFTN/gftn-services/utility/aws/golang/parameter-store"
	"github.com/GFTN/gftn-services/utility/aws/golang/secret-manager"
	"github.com/GFTN/gftn-services/utility/aws/golang/utility"
	"github.com/GFTN/gftn-services/utility/aws/golang/utility/environment"
)

func main() {
	env_test()
	//secret_test()
	//parameter_test()
	/*
		var newTest2 = utility.CredentialInfo{
			Environment: "dev",
			Domain:      "worldwire.io",
			Service:     "IBM",
			Variable:    "secret12",
		}

		var newTestContent = utility.SecretContent{
			Key:         "ding",
			Value:       "dong",
			Description: "wahaha",
			//FilePath:    "./test.json",
		}
		_ = secret_manager.CreateSecret(newTest2, newTestContent)

		res2, _ := secret_manager.GetSecret(newTest2)
		fmt.Printf("result: %s\n", res2)
	*/
}

func env_test() {
	var test = utility.CredentialInfo{
		Environment: "dev",
		Domain:      "IBM",
		Service:     "COMMON",
		Variable:    "TOKEN",
	}

	var entries = []utility.SecretEntry{
		utility.SecretEntry{
			Key:   "node_address",
			Value: "GDZXEYRX3L3KQO4CBGWXZVWJ24ATHUULWD6XTJ7YL3YYTTV4KUWIYTTO",
		},
		utility.SecretEntry{
			Key:   "node_seed",
			Value: "SCAKW7YYACCOUHQJX5JWK7YGIW3BLEE2LVE2WY2FSNE24RTR3FE2TDZU",
		},
		utility.SecretEntry{
			Key:   "node_passcpde",
			Value: "",
		},
		utility.SecretEntry{
			Key:   "public_label",
			Value: "",
		},
		utility.SecretEntry{
			Key:   "private_label",
			Value: "",
		},
	}

	var accountSecretContent = utility.SecretContent{
		Entry:       entries,
		Description: "Token account of IBM",
	}
	err := environment.Setenv(test, accountSecretContent)
	if err != nil {
		fmt.Printf("Encounter error while storing secret to AWS: %s", err)
	}

	ress, _ := environment.GetAccount(test)

	fmt.Println(ress)

	var test2 = utility.CredentialInfo{
		Environment: "dev",
		Domain:      "p1.worldwire.io",
		Service:     "crypto-service",
		Variable:    "STELLAR_NETWORK",
	}
	entries = []utility.SecretEntry{
		utility.SecretEntry{
			Key:   "STELLAR_NETWORK",
			Value: "Test SDF Network ; September 2015",
		},
	}

	accountSecretContent = utility.SecretContent{
		Entry:       entries,
		Description: "STELLAR_NETWORK",
	}

	res2s, _ := environment.GetAccount(test2)

	fmt.Println(res2s)
	/*
		entries = []utility.SecretEntry{
			utility.SecretEntry{
				// the key can be set to environment.ENV_KEY_NODE_ADDRESS to make sure naming consistency
				Key:   "keysw",
				Value: "address5d",
			},
		}
		newTestContent = utility.SecretContent{
			Entry:       entries,
			Description: "newone",
		}
		_ = environment.Updateenv(test, newTestContent)
		resss, _ := environment.Getenv(test)

		fmt.Println(resss)
	*/
}

func secret_test() {
	//secret manager
	var test = utility.CredentialInfo{
		Environment: "dev",
		Domain:      "worldwire.io",
		Service:     "IBM",
		Variable:    "secret",
	}
	res, _ := secret_manager.GetSecret(test)
	fmt.Printf("result: %s\n", res)

	var newTest = utility.CredentialInfo{
		Environment: "dev",
		Domain:      "worldwire.io",
		Service:     "IBM",
		Variable:    "dweifowehfowe",
	}

	var entries = []utility.SecretEntry{
		utility.SecretEntry{
			Key:   "account1",
			Value: "account1value",
		},
		utility.SecretEntry{
			Key:   "account2",
			Value: "account2value",
		},
	}

	var newTestContent = utility.SecretContent{
		Entry:       entries,
		Description: "newone",
	}

	_ = secret_manager.CreateSecret(newTest, newTestContent)
	res2, _ := secret_manager.GetSecret(newTest)
	fmt.Printf("result: %s\n", res2)
	/*
		var updateTestContent = utility.SecretContent{
			FilePath:    "./test.json",
			Description: "the newest one",
		}
		_ = secret_manager.UpdateSecret(newTest, updateTestContent)
		res2, _ = secret_manager.GetSecret(newTest)
		fmt.Printf("result: %s\n", res2)
	*/
	//_ = secret_manager.DeleteSecret(newTest, 7)
}

func parameter_test() {
	var test2 = utility.CredentialInfo{
		Environment: "dev",
		Domain:      "worldwire.io",
		Service:     "IBM",
		Variable:    "token_secret_key",
	}
	res3, _ := parameter_store.GetParameter(test2)

	fmt.Printf("result: %s\n", res3)

	var newParam = utility.CredentialInfo{
		Environment: "dev",
		Domain:      "worldwire.io",
		Service:     "IBM",
		Variable:    "secretNewe",
	}

	var newParameterContent = utility.ParameterContent{
		Value:       "dong",
		Description: "newone",
	}
	_ = parameter_store.CreateParameter(newParam, newParameterContent)
	res2, _ := parameter_store.GetParameter(newParam)
	fmt.Printf("result: %s\n", res2)

	var updateParameterContent = utility.ParameterContent{
		Value:       "dongdong",
		Description: "the newest one",
	}
	_ = parameter_store.UpdateParameter(newParam, updateParameterContent)
	res2, _ = parameter_store.GetParameter(newParam)
	fmt.Printf("result: %s\n", res2)
	_ = parameter_store.DeleteParameter(newParam)

}
