// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// Code generated by go-swagger; DO NOT EDIT.

package model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Blocklist blocklist
//
// A blocklist that records all the currencies/countries/particpants that is forbidden to transact with
// swagger:model Blocklist
type Blocklist struct {

	// The id of the block type
	ID string `json:"id,omitempty"`

	// The name of the block type
	Name string `json:"name,omitempty"`

	// The type of the blocklist element
	// Required: true
	// Enum: [CURRENCY COUNTRY INSTITUTION]
	Type *string `json:"type"`

	// The value of the block type
	// Required: true
	Value []string `json:"value"`
}

// Validate validates this blocklist
func (m *Blocklist) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateValue(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var blocklistTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["CURRENCY","COUNTRY","INSTITUTION"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		blocklistTypeTypePropEnum = append(blocklistTypeTypePropEnum, v)
	}
}

const (

	// BlocklistTypeCURRENCY captures enum value "CURRENCY"
	BlocklistTypeCURRENCY string = "CURRENCY"

	// BlocklistTypeCOUNTRY captures enum value "COUNTRY"
	BlocklistTypeCOUNTRY string = "COUNTRY"

	// BlocklistTypeINSTITUTION captures enum value "INSTITUTION"
	BlocklistTypeINSTITUTION string = "INSTITUTION"
)

// prop value enum
func (m *Blocklist) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, blocklistTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Blocklist) validateType(formats strfmt.Registry) error {

	if err := validate.Required("type", "body", m.Type); err != nil {
		return err
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", *m.Type); err != nil {
		return err
	}

	return nil
}

func (m *Blocklist) validateValue(formats strfmt.Registry) error {

	if err := validate.Required("value", "body", m.Value); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Blocklist) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Blocklist) UnmarshalBinary(b []byte) error {
	var res Blocklist
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
