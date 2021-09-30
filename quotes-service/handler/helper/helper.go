// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	uuid "github.com/satori/go.uuid"
	middlewares "github.com/GFTN/gftn-services/auth-service/authorization/middleware"
	"github.com/GFTN/gftn-services/utility/wwfirebase"
)

func GetIdentity(req *http.Request) (string, error) {
	sessionContext, err := middlewares.GetSessionContext(req)
	if err != nil {
		return "", err
	}
	identity := sessionContext.ParticipantID
	LOGGER.Info("Caller Identity: ", identity)
	return identity, nil

	// ori
	// token, _ := request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
	// 	return nil, nil
	// })
	// claim, _ := token.Claims.(jwt.MapClaims)
	// if str, ok := claim["domain"].(string); ok {
	// 	LOGGER.Info("Caller Identity: ", str)
	// 	return str, nil
	// }
	// return "", errors.New("Get identity from jwt token failed")
}

func CheckIdentity(req *http.Request, callerExpected string) error {
	caller, err := GetIdentity(req)
	if err != nil {
		return err
	}
	if caller != callerExpected {
		return errors.New("Caller Identity does not match with jwt token")
	}
	return nil
}

func ExtractJwt(req *http.Request) (string, error) {
	token, err := request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	if err != nil {
		return "", err
	}
	if token == nil {
		return "", errors.New("Error extracting jwt token")
	}
	tokenString := token.Raw
	return tokenString, nil
}

func CanonicaliseJson(jsonb []byte) []byte {
	temp := make(map[string]interface{})
	json.Unmarshal(jsonb, &temp)
	canonicaliseJson, _ := json.Marshal(temp)
	return canonicaliseJson
}

// SendLogToFirebase sends a log to Firebase based on the childRef provided
// `obj`: data to be posted
// `childRefs`: location to post to on Firebase (variadic)
func SendLogToFirebase(obj interface{}, childRefs ...string) error {
	if len(childRefs) == 0 {
		return errors.New("Forbiddened: childRefs = 0")
	}
	// generating uuid for obj
	uuid := uuid.Must(uuid.NewV4()).String()
	LOGGER.Debug("uuid for exchange log:", uuid)
	updates := make(map[string]interface{})
	for _, childRef := range childRefs {
		path := childRef + "/" + uuid
		updates[path] = obj
	}
	err := wwfirebase.FbRef.Update(wwfirebase.AppContext, updates)
	if err != nil {
		LOGGER.Error("Error sending log to Firebase: %s", err.Error())
		return err
	}

	return nil
}

func SendQuoteRequestToRFI(URLCallback string, quoteRequestToRFIJsonb []byte, HTTP *http.Client) error {
	resp, err := HTTP.Post(URLCallback+"/quotes/request", "application/json", bytes.NewBuffer(quoteRequestToRFIJsonb))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		err = errors.New("Error sending quote request for URL: " + URLCallback + " with status code: " + strconv.Itoa(resp.StatusCode))
		return err
	}
	return nil
}
