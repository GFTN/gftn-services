// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package message

import (
	"net/http"
	"os"
	"time"

	global_environment "github.com/GFTN/gftn-services/utility/global-environment"

	bs "github.com/BurntSushi/toml"
	"github.com/op/go-logging"
	"github.com/GFTN/gftn-services/gftn-models/model"
)

var (
	ERROR_MESSAGE_CONFIG ErrorMessageConfig
)

var LOGGER = logging.MustGetLogger("errorMessageconfig")

type ErrorMessageConfig struct {
	WWError map[string]WWErrorMessage `toml:"WW_ERROR"`
}
type WWErrorMessage struct {
	ShortMessage string `toml:"SHORT_MESSAGE"`
	Details      string `toml:"DETAILS"`
}

func LoadMessage(locale, errorCode string) model.WorldWireError {

	//We can handle internationalization in future
	//tErrorCode := common.Cat(errorCode,".", locale)
	errorMessage := ERROR_MESSAGE_CONFIG.WWError[errorCode]
	if (errorMessage != WWErrorMessage{}) {
		return model.WorldWireError{Message: &errorMessage.ShortMessage, URL: "", Code: errorCode, Details: &errorMessage.Details}
	}
	return model.WorldWireError{}
}

// Each service will Need to load error config and init and pass to WWError function each time
func LoadErrorConfig(fileName string) error {
	configFile := fileName
	if _, err := bs.DecodeFile(configFile, &ERROR_MESSAGE_CONFIG); err != nil {
		LOGGER.Errorf("errordecoding ErrorMessageConfig: %v", err)
		return err
	}
	return nil
}

func Translate(r *http.Request, errorCode string, err error) model.WorldWireError {

	localeStr := r.Header.Get("Accept-Language")
	//Hardcoding it to english for now, will need to follow  internationalization as :
	//https://phraseapp.com/blog/posts/internationalization-i18n-go/
	localeStr = "EN-US"
	errorMessage := LoadMessage(localeStr, errorCode)

	// static fields set for every error
	errorMessage.URL = r.RequestURI
	timeNow := time.Now().Unix()
	errorMessage.TimeStamp = &timeNow
	errorMessage.ParticipantID = os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)
	// Kind is set in the location of the code where thrown
	errorMessage.Type = "Generic"

	ww001 := "WW-001"
	genError := "Generic Error"
	genErrDetail := "Generic Error Message: "

	// set generic error  by default
	if errorMessage.Code == "" {
		errorMessage.Code = ww001
		errorMessage.Message = &genError
		errorMessage.Details = &genErrDetail
	}
	if err != nil {
		LOGGER.Debug(*errorMessage.Details, err.Error())
	}
	return errorMessage
}
