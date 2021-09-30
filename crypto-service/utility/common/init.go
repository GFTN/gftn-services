// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package common

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/GFTN/gftn-services/crypto-service/utility/constant"

	"github.com/miekg/pkcs11"

	"github.com/GFTN/gftn-services/crypto-service/environment"
	"github.com/GFTN/gftn-services/crypto-service/utility/hsm"
)

func ReinitializeHSM(c *hsm.Crypto, session pkcs11.SessionHandle) (*hsm.Crypto, pkcs11.SessionHandle, error) {
	var mutex = &sync.Mutex{}
	defer mutex.Unlock()
	if strings.ToUpper(os.Getenv(environment.ENV_KEY_ACCOUNT_SOURCE)) == constant.ACCOUNT_FROM_HSM {
		LOGGER.Infof("Reinitializing HSM connection")
		mutex.Lock()
		err := LogoutHSM(c, session)
		if err != nil {
			LOGGER.Warningf("Encounter error while logging out of the HSM")
		} else {
			LOGGER.Infof("Successfully logged out the HSM")
		}
		c, _, session = InitiateHSM()
		LOGGER.Infof("HSM connection reinitialized")
		return c, session, nil
	}
	return nil, 0, errors.New("Account source is not HSM")
}

// Open a HSM session with the specific slot then login with the pin code.
func InitiateHSM() (*hsm.Crypto, uint, pkcs11.SessionHandle) {
	LOGGER.Infof("Initializing HSM connection...")
	lib := os.Getenv(environment.ENV_KEY_PKCS11_LIB)
	rawSlot := os.Getenv(environment.ENV_KEY_PKCS11_SLOT)
	slot64, err := strconv.ParseUint(rawSlot, 10, 64)
	if err != nil {
		LOGGER.Errorf("Error while opening a login-session: %v\n", err)
	}
	slot := uint(slot64)
	pin := os.Getenv(environment.ENV_KEY_PKCS11_PIN)
	c := hsm.NewCrypto(lib)

	session, err := c.OpenSession(slot)
	if err != nil {
		LOGGER.Errorf("Error while opening a session: %v", err)
	}
	err = c.Login(session, pin)
	if err != nil {
		LOGGER.Errorf("Error while login to an opened login-session: %v\n", err)
	}

	if err == nil {
		LOGGER.Infof("HSM connection successfully initialized!")
	}
	//Cannot defer here as we have to maintain the login session
	//defer HsmInstance.C.Logout(session)
	//defer HsmInstance.C.CloseSession(session)
	return c, slot, session
}

func LogoutHSM(c *hsm.Crypto, session pkcs11.SessionHandle) error {
	LOGGER.Infof("Logging out of HSM...")
	err := c.Logout(session)
	if err != nil {
		LOGGER.Warningf("Error while logout the HSM session: %v", err)
	}
	err = c.FinalizeCrypto()
	if err != nil {
		LOGGER.Warningf("Error while finalizing crypto object: %v", err)
	}
	return err
}
