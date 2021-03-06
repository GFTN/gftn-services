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

// PayloadWithSignature payloadWithSignature
//
// PayloadWithSignature
// swagger:model PayloadWithSignature
type PayloadWithSignature struct {

	// Signed ISO 20022 message.
	// Required: true
	// Format: byte
	PayloadWithSignature *strfmt.Base64 `json:"payload_with_signature"`
}

// Validate validates this payload with signature
func (m *PayloadWithSignature) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePayloadWithSignature(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PayloadWithSignature) validatePayloadWithSignature(formats strfmt.Registry) error {

	if err := validate.Required("payload_with_signature", "body", m.PayloadWithSignature); err != nil {
		return err
	}

	// Format "byte" (base64 string) is already validated when unmarshalled

	return nil
}

// MarshalBinary interface implementation
func (m *PayloadWithSignature) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PayloadWithSignature) UnmarshalBinary(b []byte) error {
	var res PayloadWithSignature
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
