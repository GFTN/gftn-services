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

// Fee fee
//
// Fee
// swagger:model Fee
type Fee struct {

	// The fee amount, should be a float64 number
	// Required: true
	// Multiple Of: 1e-07
	Cost *float64 `json:"cost"`

	// costasset
	// Required: true
	Costasset *Asset `json:"costasset"`
}

// Validate validates this fee
func (m *Fee) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCost(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCostasset(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Fee) validateCost(formats strfmt.Registry) error {

	if err := validate.Required("cost", "body", m.Cost); err != nil {
		return err
	}

	if err := validate.MultipleOf("cost", "body", float64(*m.Cost), 1e-07); err != nil {
		return err
	}

	return nil
}

func (m *Fee) validateCostasset(formats strfmt.Registry) error {

	if err := validate.Required("costasset", "body", m.Costasset); err != nil {
		return err
	}

	if m.Costasset != nil {
		if err := m.Costasset.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("costasset")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Fee) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Fee) UnmarshalBinary(b []byte) error {
	var res Fee
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
