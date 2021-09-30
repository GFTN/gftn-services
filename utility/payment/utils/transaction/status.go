// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package transaction

import (
	"time"

	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/utility/payment"
	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"
)

type Payment struct {
	PaymentStatuses []model.TransactionReceipt
}

func InitiatePaymentLogOperation() (op Payment) {
	LOGGER.Info("Init payment information operation")
	status := Payment{
		PaymentStatuses: []model.TransactionReceipt{},
	}

	return status
}

func CreateDataForTransactionMemo(data *sendmodel.StatusData) model.Send {
	creditorLocationAddress := &model.Address{
		BuildingNumber: checkIsNilString(data.CreditorBuildingNo),
		City:           checkIsNilString(data.CityCdtr),
		Country:        checkIsNilString(data.CountryCdtr),
		PostalCode:     checkIsNilString(data.CreditorPostalCode),
		Street:         checkIsNilString(data.CreditorStreet),
	}

	customerLocationAddress := &model.Address{
		BuildingNumber: checkIsNilString(data.CustomerBuildingNo),
		City:           checkIsNilString(data.CityDbtr),
		Country:        checkIsNilString(data.CountryDbtr),
		PostalCode:     checkIsNilString(data.CustomerCountry),
		Street:         checkIsNilString(data.CustomerStreet),
	}

	if checkIsNilFloat(data.ExchangeRate) != nil && data.ExchangeRate != 1.0 {
		crtyCcy := data.CrtyCcy
		data.AssetCodeBeneficiary = crtyCcy
	}

	if len(data.IdCdtr) == 0 {
		data.IdCdtr = ""
	}

	tr := model.Send{
		AccountNameSend: data.AccountNameSend,
		Creditor: &model.PaymentActor{
			Customer: &model.PaymentAddress{
				Name:            checkIsNilString(data.NameCdtr),
				ID:              checkIsNilString(data.IdCdtr),
				LocationAddress: creditorLocationAddress,
			},
		},
		Debtor: &model.PaymentActor{
			Customer: &model.PaymentAddress{
				Name:            checkIsNilString(data.NameDbtr),
				ID:              checkIsNilString(data.IdDbtr),
				LocationAddress: customerLocationAddress,
			},
		},
		EndToEndID:    data.EndToEndID,
		ExchangeRate:  data.ExchangeRate,
		InstructionID: data.InstructionID,
		TransactionDetails: &model.TransactionDetails{
			AmountBeneficiary:    checkIsNilFloat(data.AmountBeneficiary),
			AssetCodeBeneficiary: checkIsNilString(data.AssetCodeBeneficiary),
			AmountSettlement:     checkIsNilFloat(data.AmountSettlement),
			Assetsettlement: &model.Asset{
				AssetCode: checkIsNilString(data.CurrencyCode),
				AssetType: checkIsNilString(data.AssetType),
				IssuerID:  data.IssuerID,
			},
			Feecreditor: &model.Fee{
				Cost: checkIsNilFloat(data.FeeCost),
				Costasset: &model.Asset{
					AssetCode: checkIsNilString(data.FeeCurrencyCode),
					AssetType: checkIsNilString(data.FeeAssetType),
					IssuerID:  data.IssuerID,
				},
			},
			OfiID:            checkIsNilString(data.IdDbtr),
			RfiID:            checkIsNilString(data.IdCdtr),
			SettlementMethod: checkIsNilString(data.AssetType),
		},
	}

	return tr
}

func (op *Payment) BuildTXMemo(opType string, data *sendmodel.StatusData, stellarTxnId, orgnlMsgId, orgnlInstrId, ofiId, messageType, messageName string) model.FitoFICCTMemoData {
	tr := CreateDataForTransactionMemo(data)
	memo := payment.BuildFiToFiCCTTxnMemo(opType, tr, stellarTxnId, orgnlMsgId, orgnlInstrId, messageType, messageName)
	memo.OfiID = &ofiId

	return memo
}

func (op *Payment) RecordPaymentStatus(data ...string) {
	if len(data) == 0 {
		LOGGER.Errorf("No parameter passed in")
		return
	}
	status := data[0]
	txHash := "N/A"
	if len(data) == 2 {
		txHash = data[1]
		LOGGER.Debugf("Current tx hash in payment status: %v", txHash)
	}
	LOGGER.Debugf("Current payment status: %v", status)
	timeStamp := time.Now().Unix()
	//payment.WritePaymentStatusToLog(memo, status, &timeStamp)
	op.PaymentStatuses = append(op.PaymentStatuses, model.TransactionReceipt{Transactionstatus: &status, Timestamp: &timeStamp, Transactionid: &txHash})
}

func checkIsNilString(str string) *string {
	if len(str) == 0 {
		return nil
	} else {
		return &str
	}
}

func checkIsNilFloat(f float64) *float64 {
	if f > 0 {
		return &f
	} else {
		return nil
	}
}
