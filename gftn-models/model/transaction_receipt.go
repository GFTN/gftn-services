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

// TransactionReceipt transactionReceipt
//
// Transaction Receipt
// swagger:model transactionReceipt
type TransactionReceipt struct {

	// The timestamp of the transaction.
	// Required: true
	Timestamp *int64 `json:"timestamp"`

	// A unique transaction identifier generated by the ledger.
	// Required: true
	Transactionid *string `json:"transactionid"`

	// This would capture the new status of a transaction while transaction travel through payment flow.
	// Required: true
	Transactionstatus *string `json:"transactionstatus"`
}

// Validate validates this transaction receipt
func (m *TransactionReceipt) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTimestamp(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTransactionid(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTransactionstatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TransactionReceipt) validateTimestamp(formats strfmt.Registry) error {

	if err := validate.Required("timestamp", "body", m.Timestamp); err != nil {
		return err
	}

	return nil
}

func (m *TransactionReceipt) validateTransactionid(formats strfmt.Registry) error {

	if err := validate.Required("transactionid", "body", m.Transactionid); err != nil {
		return err
	}

	return nil
}

func (m *TransactionReceipt) validateTransactionstatus(formats strfmt.Registry) error {

	if err := validate.Required("transactionstatus", "body", m.Transactionstatus); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TransactionReceipt) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TransactionReceipt) UnmarshalBinary(b []byte) error {
	var res TransactionReceipt
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
