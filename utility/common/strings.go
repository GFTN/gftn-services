// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package common

import (
	"bytes"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/stellar/go/support/errors"
)

const (
	//True and False : string constant for query parameters
	True    string = "true"
	False   string = "false"
	ISSUING string = "issuing"
	DEFAULT string = "default"
)

// Separator seperates the name and domain portions of an address
const Separator = "*"

var (
	// ErrInvalidAddress is the error returned when an address is invalid in
	// such a way that we do not know if the name or domain portion is at fault.
	ErrInvalidAddress = errors.New("invalid address")

	// ErrInvalidName is the error returned when an address's name portion is
	// invalid.
	ErrInvalidName = errors.New("name part of address is invalid")

	// ErrInvalidDomain is the error returned when an address's domain portion
	// is invalid.
	ErrInvalidDomain = errors.New("domain part of address is invalid")
)

/*
	Concatenates given number of string arguments
*/
func Cat(arguments ...string) string {
	var b bytes.Buffer
	for _, s := range arguments {
		b.WriteString(s)
	}
	return b.String()
}

// Split takes an address, of the form "name*domain" and provides the
// constituent elements.
func Split(address string) (name, domain string, err error) {
	parts := strings.Split(address, Separator)

	if len(parts) != 2 {
		err = ErrInvalidAddress
		return
	}

	name = parts[0]
	domain = parts[1]

	if name == "" {
		err = ErrInvalidName
		return
	}

	if !govalidator.IsDNSName(domain) {
		err = ErrInvalidDomain
		return
	}

	return
}
