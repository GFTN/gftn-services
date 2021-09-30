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

// TransactionDetails transactionDetails
//
// Transaction Details
// swagger:model transactionDetails
type TransactionDetails struct {

	// The amount the beneficiary should receive in beneficiary currency
	// Required: true
	// Multiple Of: 1e-07
	AmountBeneficiary *float64 `json:"amount_beneficiary" bson:"amount_beneficiary"`

	// The amount of the settlement.
	// Required: true
	// Multiple Of: 1e-07
	AmountSettlement *float64 `json:"amount_settlement" bson:"amount_settlement"`

	// The asset code for the beneficiary
	// Required: true
	AssetCodeBeneficiary *string `json:"asset_code_beneficiary" bson:"asset_code_beneficiary"`

	// assetsettlement
	// Required: true
	Assetsettlement *Asset `json:"assetsettlement"`

	// feecreditor
	// Required: true
	Feecreditor *Fee `json:"feecreditor"`

	// The ID that identifies the OFI Participant on the WorldWire network (i.e. uk.yourbankintheUK.payments.ibm.com).
	// Required: true
	// Max Length: 32
	// Min Length: 5
	// Pattern: ^[a-zA-Z0-9-]{5,32}$
	OfiID *string `json:"ofi_id" bson:"ofi_id"`

	// The ID that identifies the RFI Participant on the WorldWire network (i.e. uk.yourbankintheUK.payments.ibm.com).
	// Required: true
	// Max Length: 32
	// Min Length: 5
	// Pattern: ^[a-zA-Z0-9-]{5,32}$
	RfiID *string `json:"rfi_id" bson:"rfi_id"`

	// The preferred settlement method for this payment request (DA, DO, or XLM)
	// Required: true
	SettlementMethod *string `json:"settlement_method" bson:"settlement_method"`
}

// Validate validates this transaction details
func (m *TransactionDetails) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAmountBeneficiary(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAmountSettlement(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAssetCodeBeneficiary(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAssetsettlement(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFeecreditor(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOfiID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRfiID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSettlementMethod(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TransactionDetails) validateAmountBeneficiary(formats strfmt.Registry) error {

	if err := validate.Required("amount_beneficiary", "body", m.AmountBeneficiary); err != nil {
		return err
	}

	if err := validate.MultipleOf("amount_beneficiary", "body", float64(*m.AmountBeneficiary), 1e-07); err != nil {
		return err
	}

	return nil
}

func (m *TransactionDetails) validateAmountSettlement(formats strfmt.Registry) error {

	if err := validate.Required("amount_settlement", "body", m.AmountSettlement); err != nil {
		return err
	}

	if err := validate.MultipleOf("amount_settlement", "body", float64(*m.AmountSettlement), 1e-07); err != nil {
		return err
	}

	return nil
}

func (m *TransactionDetails) validateAssetCodeBeneficiary(formats strfmt.Registry) error {

	if err := validate.Required("asset_code_beneficiary", "body", m.AssetCodeBeneficiary); err != nil {
		return err
	}

	return nil
}

func (m *TransactionDetails) validateAssetsettlement(formats strfmt.Registry) error {

	if err := validate.Required("assetsettlement", "body", m.Assetsettlement); err != nil {
		return err
	}

	if m.Assetsettlement != nil {
		if err := m.Assetsettlement.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("assetsettlement")
			}
			return err
		}
	}

	return nil
}

func (m *TransactionDetails) validateFeecreditor(formats strfmt.Registry) error {

	if err := validate.Required("feecreditor", "body", m.Feecreditor); err != nil {
		return err
	}

	if m.Feecreditor != nil {
		if err := m.Feecreditor.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("feecreditor")
			}
			return err
		}
	}

	return nil
}

func (m *TransactionDetails) validateOfiID(formats strfmt.Registry) error {

	if err := validate.Required("ofi_id", "body", m.OfiID); err != nil {
		return err
	}

	if err := validate.MinLength("ofi_id", "body", string(*m.OfiID), 5); err != nil {
		return err
	}

	if err := validate.MaxLength("ofi_id", "body", string(*m.OfiID), 32); err != nil {
		return err
	}

	if err := validate.Pattern("ofi_id", "body", string(*m.OfiID), `^[a-zA-Z0-9-]{5,32}$`); err != nil {
		return err
	}

	return nil
}

func (m *TransactionDetails) validateRfiID(formats strfmt.Registry) error {

	if err := validate.Required("rfi_id", "body", m.RfiID); err != nil {
		return err
	}

	if err := validate.MinLength("rfi_id", "body", string(*m.RfiID), 5); err != nil {
		return err
	}

	if err := validate.MaxLength("rfi_id", "body", string(*m.RfiID), 32); err != nil {
		return err
	}

	if err := validate.Pattern("rfi_id", "body", string(*m.RfiID), `^[a-zA-Z0-9-]{5,32}$`); err != nil {
		return err
	}

	return nil
}

func (m *TransactionDetails) validateSettlementMethod(formats strfmt.Registry) error {

	if err := validate.Required("settlement_method", "body", m.SettlementMethod); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TransactionDetails) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TransactionDetails) UnmarshalBinary(b []byte) error {
	var res TransactionDetails
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
