// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package parse

import (
	"encoding/base64"
)

func EncodeBase64(content []byte) (string){
	LOGGER.Infof("Encode to base 64 format")

	encoded := base64.StdEncoding.EncodeToString(content)

	return encoded
}

func DecodeBase64(str string) ([]byte, error){
	LOGGER.Infof("Decode base 64 format to byte array")

	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		LOGGER.Error(err.Error())
		return nil, err
	}

	return data, nil
}
