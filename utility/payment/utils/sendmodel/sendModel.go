// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package sendmodel

import (
	"github.com/GFTN/gftn-services/gftn-models/model"
)

var TxMemo model.FitoFICCTMemoData

type DBData struct {
	MessageId      string `json:"message_id"`
	CreateDateTime string `json:"create_date_time"`
	InstrId        string `json:"instr_id" `
	EndToEndId     string `json:"end_to_end_id" `
	TxId           string `json:"tx_id" `

	// Method used to settle the (batch of) payment instructions.
	SettlementMethod string `json:"sttlm_mtd" `

	// A specific purpose account used to post debit and credit entries as a result of the transaction.
	SettlementAccountName string `json:"sttlm_acct" `
	SettlementParticipant string `json:"sttlm_participant" `
	SettlementAmount      string `json:"sttlm_amt" `
	SettlementCurrency    string `json:"sttlm_ccy" `

	AssetIssuer string `json:"issr" `

	InstructedAgentBIC string `json:"instructed_agent_bic"`
	InstructedAgentId  string `json:"instructed_agent_id"`

	InstructingAgentBIC string `json:"instructing_agent_bic"`
	InstructingAgentId  string `json:"instructing_agent_id"`

	SettlementDate string  `json:"settlement_date"`
	ExchangeRate   float64 `json:"exchange_rate"`
	ChargeBear     string  `json:"charge_bear"`
	ChargeAmount   string  `json:"charge_amount"`
	ChargeCurrency string  `json:"charge_currency"`

	ChargeAgentBIC string `json:"charge_agent_bic"`
	ChargeAgentId  string `json:"charge_agent_id"`

	InstructedAmount   string `json:"payout_amount"`
	InstructedCurrency string `json:"payout_currency"`
}

type ErrorModel struct {
	MessageID   string `json:"end_to_end_id"`
	MessageType string `json:"message_type"`
	ErrorType   string `json:"error_msg"`
}

type SendVariables struct {
	PaymentStatus    []string `json:"payment_status"`
	WWBIC            string   `json:"wwbic"`
	WWCCY            string   `json:"wwccy"`
	XSDType          []string `json:"xsd_type"`
	MSGType          []string `json:"msg_type"`
	XSDPath          []string `json:"xsd_path"`
	Status           []Code   `json:"status"`
	SettlementMethod []string `json:"settlement_method"`
	DatabaseStatus   []string `json:"database_status"`
}

type Code struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	TxStatus    string `json:"tx_status"`
	Type        int    `json:"type"`
}

type SignData struct {
	OFIId                 string
	SettlementAccountName string
	SettlementAmount      string
	AssetIssuerId         string
	CurrencyCode          string
}

type RDOData struct {
	SettlementAccountName string
	SettlementAmount      float64
	FeeAmount             float64
	MessageId             string
	OFIId                 string
	RFIId                 string
	SettlementCcy         string
	IssuerId              string
}

type StatusData struct {
	CityCdtr              string
	CountryCdtr           string
	NameCdtr              string
	IdCdtr                string
	BICCdtr               string
	CityDbtr              string
	CountryDbtr           string
	NameDbtr              string
	IdDbtr                string
	BICDbtr               string
	CurrencyCode          string
	AssetType             string
	FeeCost               float64
	FeeCurrencyCode       string
	FeeAssetType          string
	CreditorStreet        string
	CreditorBuildingNo    string
	CreditorPostalCode    string
	CustomerStreet        string
	CustomerBuildingNo    string
	CustomerCountry       string
	AccountNameSend       string
	EndToEndID            string
	InstructionID         string
	OriginalInstructionID string
	IssuerID              string
	AmountBeneficiary     float64
	AssetCodeBeneficiary  string
	CrtyCcy               string
	AmountSettlement      float64
	ExchangeRate          float64
	SupplementaryData     map[string]string
}

type XMLData struct {
	ProtoBufData             interface{}
	RequestXMLMsg            []byte
	Signature                []byte
	MessageId                string
	OriginalMsgId            string
	OriginalEndtoEndId       string
	OriginalInstructionId    string
	OperationType            string
	OFIId                    string
	RFIId                    string
	RequestMsgType           string
	RFIAccount               string
	OFISettlementAccountName string
	RFISettlementAccountName string
	SettlementAmount         string
	FeeAmount                string
	StatusCode               string
	AssetIssuer              string
	OFIBIC                   string
	ErrorCode                int
	InstructionId            string
	CurrencyCode             string
	OfiCountry               string
	RfiCountry               string
	SupplementaryData        map[string]string
}
