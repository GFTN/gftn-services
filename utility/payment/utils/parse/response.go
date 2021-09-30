// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package parse

import (
	"encoding/json"
	"encoding/xml"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	cancelreportstruct "github.com/GFTN/iso20022/camt03000105"
	reportstruct "github.com/GFTN/iso20022/pacs00200109"
	"github.com/GFTN/gftn-services/gftn-models/model"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/payment/constant"
	"github.com/GFTN/gftn-services/utility/payment/environment"
	"github.com/GFTN/gftn-services/utility/payment/utils"
	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"
	"github.com/GFTN/gftn-services/utility/payment/utils/signing"
)

var letterRunes = []rune("0123456789")

func Init() {
	rand.Seed(time.Now().UnixNano())
}

func getReportMax350Text(text string) *reportstruct.Max350Text {
	res := reportstruct.Max350Text(text)
	return &res
}
func getReportMax35Text(text string) *reportstruct.Max35Text {
	res := reportstruct.Max35Text(text)
	return &res
}

func getCancelMax35Text(text string) *cancelreportstruct.Max35Text {
	res := cancelreportstruct.Max35Text(text)
	return &res
}

func CreatePacs002(participantBIC, originalInstrId, target string, statusCode int, originalGroupInf *reportstruct.OriginalGroupInformation29) []byte {
	timeNow, _ := time.Parse("2006-01-02T15:04:05", time.Now().UTC().Format("2006-01-02T15:04:05"))
	t := reportstruct.ISODateTime(timeNow.String())

	var reason = constant.STATUS_CODE_DEFAULT
	var txStatus = constant.PAYMENT_STATUS_RJCT
	var description = ""
	var statusType = 0

	for _, s := range V.Vars.Status {
		if statusCode == s.Code {
			reason = s.Code
			description = s.Description
			txStatus = s.TxStatus
			statusType = s.Type
		}
	}

	LOGGER.Infof("Reason: %d", reason)
	LOGGER.Info("Description: " + description)

	reasonCode := strconv.Itoa(reason)

	var currencyCode string

	dateToday := time.Now().Format("02-01-2006")
	dateToday = strings.Replace(dateToday, "-", "", -1)

	wwBIC := os.Getenv(environment.ENV_KEY_WW_BIC)

	b := make([]rune, 11)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	seqNum := string(b)

	tx := reportstruct.PaymentTransaction91{}

	if statusType == 0 {
		txsts := reportstruct.ExternalPaymentTransactionStatus1Code(txStatus)
		cd := reportstruct.ExternalStatusReason1Code(reasonCode)
		var statusinformation11 []*reportstruct.StatusReasonInformation11

		info := &reportstruct.StatusReasonInformation11{
			Rsn: &reportstruct.StatusReason6Choice{
				Cd:    &cd,
				Prtry: getReportMax35Text(description),
			},
		}
		statusinformation11 = append(statusinformation11, info)

		tx = reportstruct.PaymentTransaction91{
			TxSts:     &txsts,
			StsRsnInf: statusinformation11,
		}

		currencyCode = constant.WWCCY
	} else {
		txsts := reportstruct.ExternalPaymentTransactionStatus1Code(txStatus)
		reasonInformation11 := []*reportstruct.StatusReasonInformation11{}
		cd := reportstruct.ExternalStatusReason1Code(reasonCode)
		info := &reportstruct.StatusReasonInformation11{
			Rsn: &reportstruct.StatusReason6Choice{
				Cd:    &cd,
				Prtry: getReportMax35Text(description),
			},
		}

		reasonInformation11 = append(reasonInformation11, info)

		tx = reportstruct.PaymentTransaction91{
			OrgnlGrpInf:  originalGroupInf,
			OrgnlInstrId: getReportMax35Text(originalInstrId),
			TxSts:        &txsts,
			StsRsnInf:    reasonInformation11,
		}
		if len(originalInstrId) >= 5 {
			currencyCode = originalInstrId[:5]
		}
	}

	bicfi := reportstruct.BICFIIdentifier(os.Getenv(environment.ENV_KEY_WW_BIC))
	wwId := reportstruct.Max35Text(os.Getenv(environment.ENV_KEY_WW_ID))
	agtBicfi := reportstruct.BICFIIdentifier(participantBIC)
	credt := reportstruct.ISONormalisedDateTime(timeNow.String())
	targetId := reportstruct.Max35Text(target)

	report := &reportstruct.Message{
		Body: &reportstruct.FIToFIPaymentStatusReportV09{
			GrpHdr: &reportstruct.GroupHeader53{
				MsgId:   getReportMax35Text(currencyCode + dateToday + wwBIC + seqNum),
				CreDtTm: &t,
				InstgAgt: &reportstruct.BranchAndFinancialInstitutionIdentification5{
					FinInstnId: &reportstruct.FinancialInstitutionIdentification8{
						BICFI: &bicfi,
					},
				},
				InstdAgt: &reportstruct.BranchAndFinancialInstitutionIdentification5{
					FinInstnId: &reportstruct.FinancialInstitutionIdentification8{
						BICFI: &agtBicfi,
					},
				},
			},
		},
		Head: &reportstruct.BusinessApplicationHeaderV01{
			Fr: &reportstruct.Party9Choice{
				FIId: &reportstruct.BranchAndFinancialInstitutionIdentification5{
					FinInstnId: &reportstruct.FinancialInstitutionIdentification8{
						BICFI: &bicfi,
						Othr: &reportstruct.GenericFinancialIdentification1{
							Id: &wwId,
						},
					},
				},
			},
			To: &reportstruct.Party9Choice{
				FIId: &reportstruct.BranchAndFinancialInstitutionIdentification5{
					FinInstnId: &reportstruct.FinancialInstitutionIdentification8{
						BICFI: &agtBicfi,
						Othr: &reportstruct.GenericFinancialIdentification1{
							Id: &targetId,
						},
					},
				},
			},
			BizMsgIdr: getReportMax35Text(currencyCode + dateToday + wwBIC + seqNum),
			MsgDefIdr: getReportMax35Text(constant.PACS002),
			CreDt:     &credt,
		},
	}

	report.Body.TxInfAndSts = append(report.Body.TxInfAndSts, &tx)

	msg, _ := xml.MarshalIndent(report, "", "\t")

	header := `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	cbMsg := []byte(header + string(msg))

	statusMsgType := constant.ISO20022 + ":" + constant.PACS002

	/*
		Signing message with IBM master account
	*/
	var signedMessage []byte
	var signErr error
	if os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME) == strings.ToLower(os.Getenv(environment.ENV_KEY_WW_ID)) {
		LOGGER.Infof("Signing with utility function")
		var signedRes string
		signedRes, signErr = signing.SignPayloadByMasterAccount(string(cbMsg))
		signedMessage = []byte(signedRes)
	} else {
		LOGGER.Infof("Signing with crypto-service")
		signOperation := signing.InitiateSignOperations(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL))
		signedMessage, signErr = signOperation.SignPayloadByMasterAccount(cbMsg)
	}
	if signErr != nil {
		LOGGER.Errorf("Failed to sign payload: %v", signErr.Error())
		return nil
	}

	gatewayMsg := EncodeBase64(signedMessage)

	reportClient := model.SendPacs{
		MessageType: &statusMsgType,
		Message:     &gatewayMsg,
	}

	sendResult, _ := json.Marshal(reportClient)

	LOGGER.Info("Generate pacs.002 message")

	return sendResult
}

func CreateCamt030(participantBIC, originalInstrId, xmlType, target string, statusCode int) []byte {
	timeNow, _ := time.Parse("2006-01-02T15:04:05", time.Now().UTC().Format("2006-01-02T15:04:05"))
	t := cancelreportstruct.ISODateTime(timeNow.String())

	var reason = constant.STATUS_CODE_DEFAULT
	var txStatus = constant.PAYMENT_STATUS_RJCT
	var description = ""
	var statusType = 0

	for _, s := range V.Vars.Status {
		if statusCode == s.Code {
			reason = s.Code
			description = s.Description
			txStatus = s.TxStatus
			statusType = s.Type
		}
	}

	LOGGER.Infof("Reason: %d", reason)
	LOGGER.Info("Description: " + description)

	var currencyCode, justificationCode string

	dateToday := time.Now().Format("02-01-2006")
	dateToday = strings.Replace(dateToday, "-", "", -1)

	wwBIC := os.Getenv(environment.ENV_KEY_WW_BIC)

	b := make([]rune, 11)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	seqNum := string(b)

	justificationCode = strconv.Itoa(statusCode)
	if statusType == 0 {
		currencyCode = constant.WWCCY
	} else {
		currencyCode = originalInstrId[:5]
	}

	bicfi := cancelreportstruct.BICFIDec2014Identifier(wwBIC)
	pbicfi := cancelreportstruct.BICFIDec2014Identifier(participantBIC)
	justfn := cancelreportstruct.CaseForwardingNotification3Code(justificationCode)
	nt := &cancelreportstruct.NotificationOfCaseAssignmentV05{
		Hdr: &cancelreportstruct.ReportHeader5{
			Id: getCancelMax35Text(currencyCode + dateToday + wwBIC + seqNum),
			Fr: &cancelreportstruct.Party40Choice{
				Agt: &cancelreportstruct.BranchAndFinancialInstitutionIdentification6{
					FinInstnId: &cancelreportstruct.FinancialInstitutionIdentification18{
						BICFI: &bicfi,
					},
				},
			},
			To: &cancelreportstruct.Party40Choice{
				Agt: &cancelreportstruct.BranchAndFinancialInstitutionIdentification6{
					FinInstnId: &cancelreportstruct.FinancialInstitutionIdentification18{
						BICFI: &pbicfi,
					},
				},
			},
			CreDtTm: &t,
		},
		Case: &cancelreportstruct.Case5{
			Id: getCancelMax35Text(currencyCode + dateToday + wwBIC + seqNum),
			Cretr: &cancelreportstruct.Party40Choice{
				Agt: &cancelreportstruct.BranchAndFinancialInstitutionIdentification6{
					FinInstnId: &cancelreportstruct.FinancialInstitutionIdentification18{
						BICFI: &bicfi,
					},
				},
			},
		},
		Assgnmt: &cancelreportstruct.CaseAssignment5{
			Id: getCancelMax35Text(currencyCode + dateToday + wwBIC + seqNum),
			Assgnr: &cancelreportstruct.Party40Choice{
				Agt: &cancelreportstruct.BranchAndFinancialInstitutionIdentification6{
					FinInstnId: &cancelreportstruct.FinancialInstitutionIdentification18{
						BICFI: &bicfi,
					},
				},
			},
			Assgne: &cancelreportstruct.Party40Choice{
				Agt: &cancelreportstruct.BranchAndFinancialInstitutionIdentification6{
					FinInstnId: &cancelreportstruct.FinancialInstitutionIdentification18{
						BICFI: &pbicfi,
					},
				},
			},
			CreDtTm: &t,
		},
		Ntfctn: &cancelreportstruct.CaseForwardingNotification3{
			Justfn: &justfn,
		},
	}

	LOGGER.Infof("status: %s", txStatus)

	report := &cancelreportstruct.Message{}
	report.Body = nt

	headWWBIC := cancelreportstruct.BICFIIdentifier(os.Getenv(environment.ENV_KEY_WW_BIC))
	wwId := cancelreportstruct.Max35Text(os.Getenv(environment.ENV_KEY_WW_ID))
	agtBicfi := cancelreportstruct.BICFIIdentifier(participantBIC)
	targetId := cancelreportstruct.Max35Text(target)
	credt := cancelreportstruct.ISONormalisedDateTime(timeNow.String())

	report.Head = &cancelreportstruct.BusinessApplicationHeaderV01{
		Fr: &cancelreportstruct.Party9Choice{
			FIId: &cancelreportstruct.BranchAndFinancialInstitutionIdentification5{
				FinInstnId: &cancelreportstruct.FinancialInstitutionIdentification8{
					BICFI: &headWWBIC,
					Othr: &cancelreportstruct.GenericFinancialIdentification1{
						Id: &wwId,
					},
				},
			},
		},
		To: &cancelreportstruct.Party9Choice{
			FIId: &cancelreportstruct.BranchAndFinancialInstitutionIdentification5{
				FinInstnId: &cancelreportstruct.FinancialInstitutionIdentification8{
					BICFI: &agtBicfi,
					Othr: &cancelreportstruct.GenericFinancialIdentification1{
						Id: &targetId,
					},
				},
			},
		},
		BizMsgIdr: getCancelMax35Text(currencyCode + dateToday + wwBIC + seqNum),
		MsgDefIdr: getCancelMax35Text(constant.CAMT030),
		CreDt:     &credt,
	}

	msg, _ := xml.MarshalIndent(report, "", "\t")

	header := `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	cbMsg := []byte(header + string(msg))

	/*
		Signing message with IBM master account
	*/
	var signedMessage []byte
	var signErr error
	if os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME) == strings.ToLower(os.Getenv(environment.ENV_KEY_WW_ID)) {
		LOGGER.Infof("Signing with utility function")
		var signedRes string
		signedRes, signErr = signing.SignPayloadByMasterAccount(string(cbMsg))
		signedMessage = []byte(signedRes)
	} else {
		LOGGER.Infof("Signing with crypto-service")
		signOperation := signing.InitiateSignOperations(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL))
		signedMessage, signErr = signOperation.SignPayloadByMasterAccount(cbMsg)
	}
	if signErr != nil {
		LOGGER.Errorf("Failed to sign payload: %v", signErr.Error())
		return nil
	}
	gatewayMsg := EncodeBase64(signedMessage)

	statusMsgType := constant.ISO20022 + ":" + constant.CAMT030

	reportClient := model.SendPacs{
		MessageType: &statusMsgType,
		Message:     &gatewayMsg,
	}

	sendResult, _ := json.Marshal(reportClient)

	LOGGER.Info("Generate camt.030 message")

	return sendResult
}

func CreateSuccessPacs002(participantBIC, target string, statusCode int, xmlData *sendmodel.XMLData) []byte {
	timeNow, _ := time.Parse("2006-01-02T15:04:05", time.Now().UTC().Format("2006-01-02T15:04:05"))
	t := reportstruct.ISODateTime(timeNow.String())

	var reason = constant.STATUS_CODE_DEFAULT
	var txStatus = constant.PAYMENT_STATUS_RJCT
	var description = ""
	var statusType = 0

	for _, s := range V.Vars.Status {
		if statusCode == s.Code {
			reason = s.Code
			description = s.Description
			txStatus = s.TxStatus
			statusType = s.Type
		}
	}

	LOGGER.Infof("Reason: %d", reason)
	LOGGER.Info("Description: " + description)

	reasonCode := strconv.Itoa(reason)

	var currencyCode string

	dateToday := time.Now().Format("02-01-2006")
	dateToday = strings.Replace(dateToday, "-", "", -1)

	wwBIC := os.Getenv(environment.ENV_KEY_WW_BIC)

	b := make([]rune, 11)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	seqNum := string(b)

	tx := reportstruct.PaymentTransaction91{}

	if statusType == 0 {
		txsts := reportstruct.ExternalPaymentTransactionStatus1Code(txStatus)
		cd := reportstruct.ExternalStatusReason1Code(reasonCode)
		var statusinformation11 []*reportstruct.StatusReasonInformation11

		info := &reportstruct.StatusReasonInformation11{
			Rsn: &reportstruct.StatusReason6Choice{
				Cd:    &cd,
				Prtry: getReportMax35Text(description),
			},
		}
		statusinformation11 = append(statusinformation11, info)

		tx = reportstruct.PaymentTransaction91{
			TxSts:     &txsts,
			StsRsnInf: statusinformation11,
		}

		currencyCode = constant.WWCCY
	} else {
		txsts := reportstruct.ExternalPaymentTransactionStatus1Code(txStatus)
		reasonInformation11 := []*reportstruct.StatusReasonInformation11{}
		cd := reportstruct.ExternalStatusReason1Code(reasonCode)
		info := &reportstruct.StatusReasonInformation11{
			Rsn: &reportstruct.StatusReason6Choice{
				Cd:    &cd,
				Prtry: getReportMax35Text(description),
			},
		}

		reasonInformation11 = append(reasonInformation11, info)

		tx = reportstruct.PaymentTransaction91{
			OrgnlEndToEndId: getReportMax35Text(xmlData.OriginalEndtoEndId),
			OrgnlInstrId:    getReportMax35Text(xmlData.OriginalInstructionId),
			TxSts:           &txsts,
			StsRsnInf:       reasonInformation11,
		}
		currencyCode = xmlData.CurrencyCode

	}

	bicfi := reportstruct.BICFIIdentifier(os.Getenv(environment.ENV_KEY_WW_BIC))
	wwId := reportstruct.Max35Text(os.Getenv(environment.ENV_KEY_WW_ID))
	agtBicfi := reportstruct.BICFIIdentifier(participantBIC)
	credt := reportstruct.ISONormalisedDateTime(timeNow.String())
	targetId := reportstruct.Max35Text(target)

	report := &reportstruct.Message{
		Body: &reportstruct.FIToFIPaymentStatusReportV09{
			GrpHdr: &reportstruct.GroupHeader53{
				MsgId:   getReportMax35Text(currencyCode + dateToday + wwBIC + seqNum),
				CreDtTm: &t,
				InstgAgt: &reportstruct.BranchAndFinancialInstitutionIdentification5{
					FinInstnId: &reportstruct.FinancialInstitutionIdentification8{
						BICFI: &bicfi,
					},
				},
				InstdAgt: &reportstruct.BranchAndFinancialInstitutionIdentification5{
					FinInstnId: &reportstruct.FinancialInstitutionIdentification8{
						BICFI: &agtBicfi,
					},
				},
			},
		},
		Head: &reportstruct.BusinessApplicationHeaderV01{
			Fr: &reportstruct.Party9Choice{
				FIId: &reportstruct.BranchAndFinancialInstitutionIdentification5{
					FinInstnId: &reportstruct.FinancialInstitutionIdentification8{
						BICFI: &bicfi,
						Othr: &reportstruct.GenericFinancialIdentification1{
							Id: &wwId,
						},
					},
				},
			},
			To: &reportstruct.Party9Choice{
				FIId: &reportstruct.BranchAndFinancialInstitutionIdentification5{
					FinInstnId: &reportstruct.FinancialInstitutionIdentification8{
						BICFI: &agtBicfi,
						Othr: &reportstruct.GenericFinancialIdentification1{
							Id: &targetId,
						},
					},
				},
			},
			BizMsgIdr: getReportMax35Text(currencyCode + dateToday + wwBIC + seqNum),
			MsgDefIdr: getReportMax35Text(constant.PACS002),
			CreDt:     &credt,
		},
	}

	report.Body.TxInfAndSts = append(report.Body.TxInfAndSts, &tx)
	for index, splmtryData := range xmlData.SupplementaryData {
		id := reportstruct.Max34Text(splmtryData)
		supplementaryData1 := &reportstruct.SupplementaryData1{
			PlcAndNm: getReportMax350Text(index),
			Envlp: &reportstruct.SupplementaryDataEnvelope1{
				Id: &id,
			},
		}
		report.Body.TxInfAndSts[0].SplmtryData = append(report.Body.TxInfAndSts[0].SplmtryData, supplementaryData1)
	}

	msg, _ := xml.MarshalIndent(report, "", "\t")

	header := `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	cbMsg := []byte(header + string(msg))

	statusMsgType := constant.ISO20022 + ":" + constant.PACS002

	/*
		Signing message with IBM master account
	*/
	var signedMessage []byte
	var signErr error
	if utils.StringsEqual(os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME), os.Getenv(environment.ENV_KEY_WW_ID)) {
		LOGGER.Infof("Signing with utility function")
		var signedRes string
		signedRes, signErr = signing.SignPayloadByMasterAccount(string(cbMsg))
		signedMessage = []byte(signedRes)
	} else {
		LOGGER.Infof("Signing with crypto-service")
		signOperation := signing.InitiateSignOperations(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL))
		signedMessage, signErr = signOperation.SignPayloadByMasterAccount(cbMsg)
	}
	if signErr != nil {
		LOGGER.Errorf("Failed to sign payload: %v", signErr.Error())
		return nil
	}

	gatewayMsg := EncodeBase64(signedMessage)

	reportClient := model.SendPacs{
		MessageType: &statusMsgType,
		Message:     &gatewayMsg,
	}

	sendResult, _ := json.Marshal(reportClient)

	LOGGER.Info("Generate pacs.002 message")

	return sendResult
}
