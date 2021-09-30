// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package participant

import (
	"crypto/sha256"
	"math/rand"
	"time"

	"github.com/stellar/go/strkey"
	"github.com/GFTN/gftn-services/utility/common"
)

const separator = ";"

/*
 This function would accept one string value as argument and generate SHA256 value
 key with Stellar's HashX encoding. This would be used as signer of an account.
*/
func GenerateSHA256Hash(key string) string {
	hasher := sha256.New()
	hasher.Write([]byte(key))

	actual, err := strkey.Encode(strkey.VersionByteHashX, hasher.Sum(nil))
	if err != nil {
		LOGGER.Fatal(err)
		return ""
	}
	return actual
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func GetSecretPhrase() string {
	n := common.KillswitchStringLength
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r1.Intn(len(letterBytes))]
	}
	return string(b)
}
