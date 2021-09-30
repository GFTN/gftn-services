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

// Participant participant
//
// Participant
// swagger:model Participant
type Participant struct {

	// The business identifier code of each participant
	// Required: true
	// Max Length: 11
	// Min Length: 11
	// Pattern: ^[A-Z]{3}[A-Z]{3}[A-Z2-9]{1}[A-NP-Z0-9]{1}[A-Z0-9]{3}$
	Bic *string `json:"bic"`

	// Callback url of the finiancial institute's backend system.
	// Required: true
	CallbackURL *string `json:"callbackUrl"`

	// Participant's country of residence
	// Required: true
	CountryCode *string `json:"countryCode"`

	// The participant domain for the participant
	// Required: true
	// Max Length: 32
	// Min Length: 5
	// Pattern: ^[a-zA-Z0-9-]{5,32}$
	ParticipantID *string `json:"participantId"`

	// RDO client url of the finiancial institute's backend system.
	// Required: true
	RdoClientURL *string `json:"rdoClientUrl"`

	// The Role of this registered participant, it can be MM for Market Maker and IS for Issuer or anchor
	// Required: true
	// Max Length: 2
	// Min Length: 2
	// Enum: [MM IS]
	Role *string `json:"role"`
}

// Validate validates this participant
func (m *Participant) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBic(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCallbackURL(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCountryCode(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateParticipantID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRdoClientURL(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRole(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Participant) validateBic(formats strfmt.Registry) error {

	if err := validate.Required("bic", "body", m.Bic); err != nil {
		return err
	}

	if err := validate.MinLength("bic", "body", string(*m.Bic), 11); err != nil {
		return err
	}

	if err := validate.MaxLength("bic", "body", string(*m.Bic), 11); err != nil {
		return err
	}

	if err := validate.Pattern("bic", "body", string(*m.Bic), `^[A-Z]{3}[A-Z]{3}[A-Z2-9]{1}[A-NP-Z0-9]{1}[A-Z0-9]{3}$`); err != nil {
		return err
	}

	return nil
}

func (m *Participant) validateCallbackURL(formats strfmt.Registry) error {

	if err := validate.Required("callbackUrl", "body", m.CallbackURL); err != nil {
		return err
	}

	return nil
}

func (m *Participant) validateCountryCode(formats strfmt.Registry) error {

	if err := validate.Required("countryCode", "body", m.CountryCode); err != nil {
		return err
	}

	return nil
}

func (m *Participant) validateParticipantID(formats strfmt.Registry) error {

	if err := validate.Required("participantId", "body", m.ParticipantID); err != nil {
		return err
	}

	if err := validate.MinLength("participantId", "body", string(*m.ParticipantID), 5); err != nil {
		return err
	}

	if err := validate.MaxLength("participantId", "body", string(*m.ParticipantID), 32); err != nil {
		return err
	}

	if err := validate.Pattern("participantId", "body", string(*m.ParticipantID), `^[a-zA-Z0-9-]{5,32}$`); err != nil {
		return err
	}

	return nil
}

func (m *Participant) validateRdoClientURL(formats strfmt.Registry) error {

	if err := validate.Required("rdoClientUrl", "body", m.RdoClientURL); err != nil {
		return err
	}

	return nil
}

var participantTypeRolePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["MM","IS"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		participantTypeRolePropEnum = append(participantTypeRolePropEnum, v)
	}
}

const (

	// ParticipantRoleMM captures enum value "MM"
	ParticipantRoleMM string = "MM"

	// ParticipantRoleIS captures enum value "IS"
	ParticipantRoleIS string = "IS"
)

// prop value enum
func (m *Participant) validateRoleEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, participantTypeRolePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Participant) validateRole(formats strfmt.Registry) error {

	if err := validate.Required("role", "body", m.Role); err != nil {
		return err
	}

	if err := validate.MinLength("role", "body", string(*m.Role), 2); err != nil {
		return err
	}

	if err := validate.MaxLength("role", "body", string(*m.Role), 2); err != nil {
		return err
	}

	// value enum
	if err := m.validateRoleEnum("role", "body", *m.Role); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Participant) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Participant) UnmarshalBinary(b []byte) error {
	var res Participant
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
