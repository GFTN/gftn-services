// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package crypto_handler

import (
	"errors"
	"os"
	"strings"

	"github.com/op/go-logging"
	"github.com/GFTN/gftn-services/crypto-service/environment"
	"github.com/GFTN/gftn-services/crypto-service/utility/common"
	"github.com/GFTN/gftn-services/crypto-service/utility/constant"
	util "github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	vauth "github.com/GFTN/gftn-services/utility/vault/auth"
	"github.com/GFTN/gftn-services/utility/vault/utils"
)

var LOGGER = logging.MustGetLogger("crypto-handler")

type CryptoOperations struct {
	HSMInstance  common.HsmObject
	VaultSession utils.Session
}

//Global handler variable used for clean up session in the end
var CYPTO_OPERATIONS = CryptoOperations{}

func CreateCryptoOperations() (op CryptoOperations, err error) {

	//op.HSMInstance = common.HsmInstance
	op.VaultSession = utils.Session{}

	if strings.ToUpper(os.Getenv(global_environment.ENV_KEY_SECRET_STORAGE_LOCATION)) == util.VAULT_SECRET {
		//Vault location
		op.VaultSession, err = vauth.GetSession()
		if err != nil {
			LOGGER.Errorf("Error reading account source environment settings")
			return op, err
		}
	}

	if strings.ToUpper(os.Getenv(environment.ENV_KEY_ACCOUNT_SOURCE)) == constant.ACCOUNT_FROM_HSM {
		LOGGER.Infof("Using HSM as account source")
		if os.Getenv(environment.ENV_KEY_PKCS11_PIN) == "" || os.Getenv(environment.ENV_KEY_PKCS11_SLOT) == "" {
			LOGGER.Errorf("Error reading PKCS11_PIN && PKCS11_SLOT environment settings")
			return op, errors.New("Error reading PKCS11_PIN && PKCS11_SLOT environment settings")
		}
		if op.HSMInstance.C == nil {
			LOGGER.Infof("Initializing HSM client")
			op.HSMInstance.C, op.HSMInstance.Slot, op.HSMInstance.Session = common.InitiateHSM()
		} else {
			LOGGER.Infof("HSM client already initialized. Skipped")
		}
	} else if strings.ToUpper(os.Getenv(environment.ENV_KEY_ACCOUNT_SOURCE)) == constant.ACCOUNT_FROM_STELLAR {
		LOGGER.Infof("Using HSM as account source")
	} else {
		LOGGER.Errorf("Error reading account source environment settings")
	}
	CYPTO_OPERATIONS = op
	return op, nil
}
