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

// Compliance compliance
//
// Compliance
// swagger:model compliance
type Compliance struct {

	// The hash of the entire one-way Send (payment) bundle that is stored in the txn memo field on the ledger.
	// Required: true
	Compliance *string `json:"compliance"`

	// pending send
	// Required: true
	PendingSend *Send `json:"pending_send"`
}

// Validate validates this compliance
func (m *Compliance) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCompliance(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePendingSend(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Compliance) validateCompliance(formats strfmt.Registry) error {

	if err := validate.Required("compliance", "body", m.Compliance); err != nil {
		return err
	}

	return nil
}

func (m *Compliance) validatePendingSend(formats strfmt.Registry) error {

	if err := validate.Required("pending_send", "body", m.PendingSend); err != nil {
		return err
	}

	if m.PendingSend != nil {
		if err := m.PendingSend.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("pending_send")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Compliance) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Compliance) UnmarshalBinary(b []byte) error {
	var res Compliance
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
