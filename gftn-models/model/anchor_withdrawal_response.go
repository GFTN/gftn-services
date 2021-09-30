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

// AnchorWithdrawalResponse settlementReceipt
//
// Settlement Receipt
// swagger:model AnchorWithdrawalResponse
type AnchorWithdrawalResponse struct {

	// The fee amount, should be a float64 number
	// Multiple Of: 1e-07
	AmountFee float64 `json:"amount_fee,omitempty" bson:"amount_fee"`

	// The identifier of the asset the anchor issued. For a list of assets, retrieve all World Wire assets from the /assets endpoint.
	AssetCode string `json:"asset_code,omitempty"`

	// A reference hash to refer transaction in anchor's system of record.
	ReferenceHash string `json:"reference_hash,omitempty"`

	// A hash that identifies the transaction on the ledger.
	TransactionID string `json:"transaction_id,omitempty"`

	// An optional way for customers to name a transaction.
	TransactionNote string `json:"transaction_note,omitempty"`
}

// Validate validates this anchor withdrawal response
func (m *AnchorWithdrawalResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAmountFee(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AnchorWithdrawalResponse) validateAmountFee(formats strfmt.Registry) error {

	if swag.IsZero(m.AmountFee) { // not required
		return nil
	}

	if err := validate.MultipleOf("amount_fee", "body", float64(m.AmountFee), 1e-07); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *AnchorWithdrawalResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AnchorWithdrawalResponse) UnmarshalBinary(b []byte) error {
	var res AnchorWithdrawalResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
