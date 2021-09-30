// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package parse

import (
	"encoding/json"
	"strconv"

	pacs008struct "github.com/GFTN/iso20022/pacs00800107"
	pacs009struct "github.com/GFTN/iso20022/pacs00900108"

	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/utility/payment/utils/database"

	"github.com/GFTN/gftn-services/utility/payment/constant"
	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"
)

func CreatePacs008DbData(d *pacs008struct.FIToFICustomerCreditTransferV07) *sendmodel.DBData {
	// Data that will be store into the DynamoDB for later use
	exchangeRate, _ := strconv.ParseFloat(string(*d.CdtTrfTxInf[0].XchgRate), 64)
	dbData := &sendmodel.DBData{
		MessageId:             string(*d.GrpHdr.MsgId),
		CreateDateTime:        string(*d.GrpHdr.CreDtTm),
		InstrId:               string(*d.CdtTrfTxInf[0].PmtId.InstrId),
		EndToEndId:            string(*d.CdtTrfTxInf[0].PmtId.EndToEndId),
		TxId:                  string(*d.CdtTrfTxInf[0].PmtId.TxId),
		SettlementAmount:      d.CdtTrfTxInf[0].IntrBkSttlmAmt.Value,
		SettlementCurrency:    d.CdtTrfTxInf[0].IntrBkSttlmAmt.Currency,
		SettlementAccountName: string(*d.GrpHdr.SttlmInf.SttlmAcct.Nm),
		SettlementMethod:      string(*d.GrpHdr.SttlmInf.SttlmMtd),
		SettlementParticipant: string(*d.GrpHdr.SttlmInf.SttlmAcct.Id.Othr.Id),
		AssetIssuer:           string(*d.GrpHdr.PmtTpInf.SvcLvl.Prtry),
		InstructedAgentBIC:    string(*d.GrpHdr.InstdAgt.FinInstnId.BICFI),
		InstructedAgentId:     string(*d.GrpHdr.InstdAgt.FinInstnId.Othr.Id),
		InstructingAgentBIC:   string(*d.GrpHdr.InstgAgt.FinInstnId.BICFI),
		InstructingAgentId:    string(*d.GrpHdr.InstgAgt.FinInstnId.Othr.Id),
		SettlementDate:        string(*d.CdtTrfTxInf[0].IntrBkSttlmDt),
		ExchangeRate:          exchangeRate,
		ChargeBear:            string(*d.CdtTrfTxInf[0].ChrgBr),
		ChargeAmount:          d.CdtTrfTxInf[0].ChrgsInf[0].Amt.Value,
		ChargeCurrency:        d.CdtTrfTxInf[0].ChrgsInf[0].Amt.Currency,
		ChargeAgentBIC:        string(*d.CdtTrfTxInf[0].ChrgsInf[0].Agt.FinInstnId.BICFI),
		ChargeAgentId:         string(*d.CdtTrfTxInf[0].ChrgsInf[0].Agt.FinInstnId.Othr.Id),
		InstructedAmount:      d.CdtTrfTxInf[0].InstdAmt.Value,
		InstructedCurrency:    d.CdtTrfTxInf[0].InstdAmt.Value,
	}

	return dbData
}

func CreatePacs009DbData(d *pacs009struct.FinancialInstitutionCreditTransferV08) *sendmodel.DBData {
	// Data that will be store into the DynamoDB for later use
	dbData := &sendmodel.DBData{
		MessageId:             string(*d.GrpHdr.MsgId),
		CreateDateTime:        string(*d.GrpHdr.CreDtTm),
		InstrId:               string(*d.CdtTrfTxInf[0].PmtId.InstrId),
		EndToEndId:            string(*d.CdtTrfTxInf[0].PmtId.EndToEndId),
		TxId:                  string(*d.CdtTrfTxInf[0].PmtId.TxId),
		SettlementAmount:      d.CdtTrfTxInf[0].IntrBkSttlmAmt.Value,
		SettlementCurrency:    d.CdtTrfTxInf[0].IntrBkSttlmAmt.Currency,
		SettlementAccountName: string(*d.GrpHdr.SttlmInf.SttlmAcct.Nm),
		SettlementMethod:      string(*d.GrpHdr.SttlmInf.SttlmMtd),
		SettlementParticipant: string(*d.GrpHdr.SttlmInf.SttlmAcct.Id.Othr.Id),
		InstructedAgentBIC:    string(*d.GrpHdr.InstdAgt.FinInstnId.BICFI),
		InstructedAgentId:     string(*d.GrpHdr.InstdAgt.FinInstnId.Othr.Id),
		InstructingAgentBIC:   string(*d.GrpHdr.InstgAgt.FinInstnId.BICFI),
		InstructingAgentId:    string(*d.GrpHdr.InstgAgt.FinInstnId.Othr.Id),
		SettlementDate:        string(*d.CdtTrfTxInf[0].IntrBkSttlmDt),
	}

	return dbData
}

func GetDBData(msgId string) (interface{}, []model.TransactionReceipt) {
	// Get transaction data and payment status information from DynamoDB
	dbData, txStatus, _, dbPaymentInfo, dbErr := database.DC.GetTransactionData(msgId)
	if dbErr != nil || *txStatus == "" {
		LOGGER.Error("Error getting transaction data")
		return nil, nil
	}

	var data sendmodel.DBData

	if *dbData != constant.EMPTY_STRING {
		pbDBData, _ := DecodeBase64(*dbData)
		json.Unmarshal(pbDBData, &data)
	}

	paymentInfo, _ := DecodeBase64(*dbPaymentInfo)

	var info []model.TransactionReceipt
	json.Unmarshal(paymentInfo, &info)

	return &data, info
}
