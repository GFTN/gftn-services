// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// Code generated by go-swagger; DO NOT EDIT.

package model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SweepReceipt SweepReceipt
//
// Sweep Receipt
// swagger:model SweepReceipt
type SweepReceipt struct {

	// Source account balances, after sweeping.
	BalanceResult []*Sweep `json:"balance_result"`

	// Timestamp when the exchange occurred.
	TimeExecuted int64 `json:"time_executed,omitempty"`

	// Transacted hash.
	// Required: true
	TransactionHash *string `json:"transaction_hash"`
}

// Validate validates this sweep receipt
func (m *SweepReceipt) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBalanceResult(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTransactionHash(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SweepReceipt) validateBalanceResult(formats strfmt.Registry) error {

	if swag.IsZero(m.BalanceResult) { // not required
		return nil
	}

	for i := 0; i < len(m.BalanceResult); i++ {
		if swag.IsZero(m.BalanceResult[i]) { // not required
			continue
		}

		if m.BalanceResult[i] != nil {
			if err := m.BalanceResult[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("balance_result" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *SweepReceipt) validateTransactionHash(formats strfmt.Registry) error {

	if err := validate.Required("transaction_hash", "body", m.TransactionHash); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SweepReceipt) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SweepReceipt) UnmarshalBinary(b []byte) error {
	var res SweepReceipt
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
