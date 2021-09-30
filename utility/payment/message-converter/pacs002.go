// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package message_converter

import (
	"encoding/xml"
	"errors"
	"os"

	"github.com/beevik/etree"

	pacs "github.com/GFTN/iso20022/pacs00200109"

	pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/pacs00200109"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/payment/client"
	"github.com/GFTN/gftn-services/utility/payment/constant"
	"github.com/GFTN/gftn-services/utility/payment/environment"
	"github.com/GFTN/gftn-services/utility/payment/utils"
	"github.com/GFTN/gftn-services/utility/payment/utils/parse"
	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"
)

type Pacs002 struct {
	Message     *pacs.Message
	SendPayload pbstruct.SendPayload
	Raw         []byte
}

func (msg *Pacs002) SanityCheck(bic, participantId string) error {

	//check if BIC code is defined & the destination is world wire
	if msg.Message.Head.To == nil ||
		msg.Message.Head.To.FIId == nil ||
		msg.Message.Head.To.FIId.FinInstnId == nil ||
		msg.Message.Head.To.FIId.FinInstnId.Othr == nil ||
		!utils.StringsEqual(string(*msg.Message.Head.To.FIId.FinInstnId.BICFI), os.Getenv(environment.ENV_KEY_WW_BIC)) ||
		!utils.StringsEqual(string(*msg.Message.Head.To.FIId.FinInstnId.Othr.Id), os.Getenv(environment.ENV_KEY_WW_ID)) {
		return errors.New("Destination BIC or ID is undefined or incorrect")
	}

	//check if BIC code is defined & the source is RFI
	if msg.Message.Head.Fr == nil ||
		msg.Message.Head.Fr.FIId == nil ||
		msg.Message.Head.Fr.FIId.FinInstnId == nil ||
		msg.Message.Head.Fr.FIId.FinInstnId.Othr == nil ||
		!utils.StringsEqual(string(*msg.Message.Head.Fr.FIId.FinInstnId.BICFI), bic) ||
		!utils.StringsEqual(string(*msg.Message.Head.Fr.FIId.FinInstnId.Othr.Id), participantId) {
		return errors.New("Source BIC or ID is undefined or incorrect")
	}

	//check the format of both message identifier & business message identifier in the header
	if err := parse.HeaderIdentifierCheck(string(*msg.Message.Head.BizMsgIdr), string(*msg.Message.Head.MsgDefIdr), constant.PACS002); err != nil {
		return err
	}

	//check the format & value of instruction id
	if err := parse.InstructionIdCheck(string(*msg.Message.Body.GrpHdr.MsgId)); err != nil {
		return err
	}

	//check the BIC code and ID of both OFI & RFI is defined
	if msg.Message.Body.GrpHdr.InstdAgt.FinInstnId.Othr == nil ||
		msg.Message.Body.GrpHdr.InstdAgt.FinInstnId.Othr.Id == nil ||
		msg.Message.Body.GrpHdr.InstgAgt.FinInstnId.Othr == nil ||
		msg.Message.Body.GrpHdr.InstgAgt.FinInstnId.Othr.Id == nil {
		return errors.New("Instructing/Instructed agent ID is missing")
	}

	//check if the supplementary data are correct
	splmtryDataSet := make(map[string]string)

	for _, splmtryData := range msg.Message.Body.TxInfAndSts[0].SplmtryData {
		index := string(*splmtryData.PlcAndNm)
		if splmtryData.Envlp != nil && splmtryData.Envlp.Id != nil {
			splmtryDataSet[index] = string(*splmtryData.Envlp.Id)
		}
	}

	if splmtryDataSet[constant.PACS002_SUPPLEMENTARY_DATA_ISSUER] == "" ||
		splmtryDataSet[constant.PACS002_SUPPLEMENTARY_DATA_ACCOUNT_ADDRESS] == "" ||
		splmtryDataSet[constant.PACS002_SUPPLEMENTARY_DATA_PAY_REFERENCE] == "" {
		return errors.New("Supplementary Data is missing")
	}

	return nil
}

func GetElementByName(element *etree.Element, tag string) *etree.Element {
	for _, child := range element.ChildElements() {
		if child.Tag == tag {
			return child
		}
	}

	return nil
}

func (msg *Pacs002) RequestToStruct() error {
	LOGGER.Infof("Constructing request data to go struct...")
	err := xml.Unmarshal(msg.Raw, &msg.Message)
	if err != nil {
		return err
	}

	LOGGER.Infof("Go struct constructed successfully")
	return nil
}

func (msg *Pacs002) StructToProto() error {
	LOGGER.Infof("Constructing to protobuffer...")

	//putting raw xml message into the kafka send payload
	msg.SendPayload.Message = msg.Raw
	msg.SendPayload.MsgType = constant.ISO20022 + ":" + constant.PACS002
	msg.SendPayload.OfiId = string(*msg.Message.Body.GrpHdr.InstdAgt.FinInstnId.Othr.Id)
	msg.SendPayload.RfiId = string(*msg.Message.Body.GrpHdr.InstgAgt.FinInstnId.Othr.Id)
	msg.SendPayload.InstructionId = string(*msg.Message.Body.GrpHdr.MsgId)
	msg.SendPayload.OriginalInstructionId = string(*msg.Message.Body.TxInfAndSts[0].OrgnlInstrId)
	msg.SendPayload.MsgId = string(*msg.Message.Body.GrpHdr.MsgId)

	LOGGER.Infof("Protobuffer successfully constructed")
	return nil
}

func (msg *Pacs002) ProtobuftoStruct() (*sendmodel.XMLData, error) {
	LOGGER.Infof("Restoring protobuffer to go struct...")

	err := xml.Unmarshal(msg.SendPayload.Message, &msg.Message)
	if err != nil {
		return nil, errors.New("Encounter error while unmarshaling xml message")
	}

	msgId := string(*msg.Message.Body.GrpHdr.MsgId)
	ofiId := string(*msg.Message.Body.GrpHdr.InstgAgt.FinInstnId.Othr.Id)
	rfiId := string(*msg.Message.Body.GrpHdr.InstdAgt.FinInstnId.Othr.Id)
	bizMsgIdr := msg.Message.Head.BizMsgIdr
	msgDefIdr := msg.Message.Head.MsgDefIdr
	creDt := msg.Message.Head.CreDt
	instructionId := string(*msg.Message.Body.GrpHdr.MsgId)
	status := string(*msg.Message.Body.TxInfAndSts[0].TxSts)
	currencyCode := msg.Message.Body.TxInfAndSts[0].OrgnlTxRef.IntrBkSttlmAmt.Currency
	settlementAmount := msg.Message.Body.TxInfAndSts[0].OrgnlTxRef.IntrBkSttlmAmt.Value
	feeAmount := msg.Message.Body.TxInfAndSts[0].ChrgsInf[0].Amt.Value
	splmtryDataSet := make(map[string]string)

	for _, splmtryData := range msg.Message.Body.TxInfAndSts[0].SplmtryData {
		index := string(*splmtryData.PlcAndNm)
		if splmtryData.Envlp != nil && splmtryData.Envlp.Id != nil {
			splmtryDataSet[index] = string(*splmtryData.Envlp.Id)
		}
	}

	if splmtryDataSet[constant.PACS002_SUPPLEMENTARY_DATA_ISSUER] == "" ||
		splmtryDataSet[constant.PACS002_SUPPLEMENTARY_DATA_ACCOUNT_ADDRESS] == "" ||
		splmtryDataSet[constant.PACS002_SUPPLEMENTARY_DATA_PAY_REFERENCE] == "" {
		return nil, errors.New("Supplementary Data is missing")
	}

	settlementAccountName := splmtryDataSet[constant.PACS002_SUPPLEMENTARY_DATA_ACCOUNT_ADDRESS]
	assetIssuer := splmtryDataSet[constant.PACS002_SUPPLEMENTARY_DATA_ISSUER]

	LOGGER.Infof("Retrieving rfi BIC code from participant registry")
	rfiBic := client.GetParticipantAccount(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL), rfiId, constant.BIC_STRING)
	wwBicfi := pacs.BICFIIdentifier(os.Getenv(environment.ENV_KEY_WW_BIC))
	wwId := pacs.Max35Text(os.Getenv(environment.ENV_KEY_WW_ID))
	rfiBicfi := pacs.BICFIIdentifier(*rfiBic)
	id := pacs.Max35Text(rfiId)

	LOGGER.Infof("Updating Business Application Header at RFI side")
	msg.Message.Head = &pacs.BusinessApplicationHeaderV01{}
	msg.Message.Head.Fr = &pacs.Party9Choice{}
	msg.Message.Head.Fr.FIId = &pacs.BranchAndFinancialInstitutionIdentification5{}
	msg.Message.Head.Fr.FIId.FinInstnId = &pacs.FinancialInstitutionIdentification8{}
	msg.Message.Head.Fr.FIId.FinInstnId.Othr = &pacs.GenericFinancialIdentification1{}

	msg.Message.Head.To = &pacs.Party9Choice{}
	msg.Message.Head.To.FIId = &pacs.BranchAndFinancialInstitutionIdentification5{}
	msg.Message.Head.To.FIId.FinInstnId = &pacs.FinancialInstitutionIdentification8{}
	msg.Message.Head.To.FIId.FinInstnId.Othr = &pacs.GenericFinancialIdentification1{}

	msg.Message.Head.Fr.FIId.FinInstnId.BICFI = &wwBicfi
	msg.Message.Head.Fr.FIId.FinInstnId.Othr.Id = &wwId
	msg.Message.Head.To.FIId.FinInstnId.BICFI = &rfiBicfi
	msg.Message.Head.To.FIId.FinInstnId.Othr.Id = &id

	msg.Message.Head.BizMsgIdr = bizMsgIdr
	msg.Message.Head.MsgDefIdr = msgDefIdr
	msg.Message.Head.CreDt = creDt

	//marshaling go struct back to raw xml message
	xmlMsg, err := xml.MarshalIndent(msg.Message, "", "\t")
	if err != nil {
		LOGGER.Warningf("Error while marshaling the xml")
		return nil, err
	}

	xmlData := &sendmodel.XMLData{
		RequestXMLMsg:            xmlMsg,
		MessageId:                msgId,
		RequestMsgType:           msg.SendPayload.MsgType,
		OFIId:                    ofiId,
		RFIId:                    rfiId,
		RFISettlementAccountName: settlementAccountName,
		InstructionId:            instructionId,
		OriginalInstructionId:    instructionId,
		StatusCode:               status,
		AssetIssuer:              assetIssuer,
		CurrencyCode:             currencyCode,
		SettlementAmount:         settlementAmount,
		FeeAmount:                feeAmount,
		SupplementaryData:        splmtryDataSet,
	}
	LOGGER.Infof("Restoring protobuffer successfully")
	return xmlData, nil

}
