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

// FinancialInstitutionDefinition institution
//
// Institution
// swagger:model FinancialInstitutionDefinition
type FinancialInstitutionDefinition struct {

	// The name of the Institution.
	// Required: true
	Name *string `json:"name"`

	// The ID that identifies a Participant on the WorldWire network (i.e. uk.yourbankintheUK.payments.ibm.com).
	// Required: true
	// Max Length: 32
	// Min Length: 5
	// Pattern: ^[a-zA-Z0-9-]{5,32}$
	ParticipantID *string `json:"participant_id"`
}

// Validate validates this financial institution definition
func (m *FinancialInstitutionDefinition) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateParticipantID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *FinancialInstitutionDefinition) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *FinancialInstitutionDefinition) validateParticipantID(formats strfmt.Registry) error {

	if err := validate.Required("participant_id", "body", m.ParticipantID); err != nil {
		return err
	}

	if err := validate.MinLength("participant_id", "body", string(*m.ParticipantID), 5); err != nil {
		return err
	}

	if err := validate.MaxLength("participant_id", "body", string(*m.ParticipantID), 32); err != nil {
		return err
	}

	if err := validate.Pattern("participant_id", "body", string(*m.ParticipantID), `^[a-zA-Z0-9-]{5,32}$`); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *FinancialInstitutionDefinition) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *FinancialInstitutionDefinition) UnmarshalBinary(b []byte) error {
	var res FinancialInstitutionDefinition
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
