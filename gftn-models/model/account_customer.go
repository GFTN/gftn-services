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

// AccountCustomer accountCustomer
//
// Account customer
// swagger:model AccountCustomer
type AccountCustomer struct {

	// Identifier for the customer account
	// Required: true
	AccountNumber *string `json:"account_number" bson:"account_number"`

	// Account type for customer account
	// Required: true
	// Enum: [checking savings]
	AccountType *string `json:"account_type" bson:"account_type"`

	// A routing number to an institution
	// Required: true
	RoutingNumber *string `json:"routing_number" bson:"routing_number"`
}

// Validate validates this account customer
func (m *AccountCustomer) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAccountNumber(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAccountType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRoutingNumber(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AccountCustomer) validateAccountNumber(formats strfmt.Registry) error {

	if err := validate.Required("account_number", "body", m.AccountNumber); err != nil {
		return err
	}

	return nil
}

var accountCustomerTypeAccountTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["checking","savings"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		accountCustomerTypeAccountTypePropEnum = append(accountCustomerTypeAccountTypePropEnum, v)
	}
}

const (

	// AccountCustomerAccountTypeChecking captures enum value "checking"
	AccountCustomerAccountTypeChecking string = "checking"

	// AccountCustomerAccountTypeSavings captures enum value "savings"
	AccountCustomerAccountTypeSavings string = "savings"
)

// prop value enum
func (m *AccountCustomer) validateAccountTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, accountCustomerTypeAccountTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *AccountCustomer) validateAccountType(formats strfmt.Registry) error {

	if err := validate.Required("account_type", "body", m.AccountType); err != nil {
		return err
	}

	// value enum
	if err := m.validateAccountTypeEnum("account_type", "body", *m.AccountType); err != nil {
		return err
	}

	return nil
}

func (m *AccountCustomer) validateRoutingNumber(formats strfmt.Registry) error {

	if err := validate.Required("routing_number", "body", m.RoutingNumber); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *AccountCustomer) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AccountCustomer) UnmarshalBinary(b []byte) error {
	var res AccountCustomer
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
