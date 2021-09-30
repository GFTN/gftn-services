// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// Code generated by go-swagger; DO NOT EDIT.

package model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// DiscoverParticipantResponse addressLedger
//
// Address Ledger
// swagger:model DiscoverParticipantResponse
type DiscoverParticipantResponse struct {

	// Can be either 'issuing' or the Participants operating account's name.
	// Required: true
	AccountName *string `json:"account_name"`

	// The ledger address which is expected to be the recipient for this transaction, once compliance checks are complete.
	// Required: true
	Address *string `json:"address"`
}

// Validate validates this discover participant response
func (m *DiscoverParticipantResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAccountName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAddress(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DiscoverParticipantResponse) validateAccountName(formats strfmt.Registry) error {

	if err := validate.Required("account_name", "body", m.AccountName); err != nil {
		return err
	}

	return nil
}

func (m *DiscoverParticipantResponse) validateAddress(formats strfmt.Registry) error {

	if err := validate.Required("address", "body", m.Address); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DiscoverParticipantResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DiscoverParticipantResponse) UnmarshalBinary(b []byte) error {
	var res DiscoverParticipantResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
