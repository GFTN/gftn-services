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

// Coordinate coordinate
//
// Geographic coordinates for a location. Based on https://schema.org/geo
// swagger:model Coordinate
type Coordinate struct {

	// The latitude of the geo coordinates
	// Required: true
	Lat *float64 `json:"lat"`

	// The longitude of the geo coordinates
	// Required: true
	Long *float64 `json:"long"`
}

// Validate validates this coordinate
func (m *Coordinate) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLat(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLong(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Coordinate) validateLat(formats strfmt.Registry) error {

	if err := validate.Required("lat", "body", m.Lat); err != nil {
		return err
	}

	return nil
}

func (m *Coordinate) validateLong(formats strfmt.Registry) error {

	if err := validate.Required("long", "body", m.Long); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Coordinate) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Coordinate) UnmarshalBinary(b []byte) error {
	var res Coordinate
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}