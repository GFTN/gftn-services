// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package message_converter

import (
	"encoding/xml"
	"errors"
	"os"
	"strconv"

	ibwf "github.com/GFTN/iso20022/ibwf00100101"
	pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/ibwf00100101"
	"github.com/GFTN/gftn-services/utility/payment/constant"
	"github.com/GFTN/gftn-services/utility/payment/environment"
	"github.com/GFTN/gftn-services/utility/payment/utils"
	"github.com/GFTN/gftn-services/utility/payment/utils/parse"
	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"
)

type Ibwf001 struct {
	Message     *ibwf.Message
	SendPayload pbstruct.SendPayload
	Raw         []byte
}

func (msg *Ibwf001) SanityCheck(bic, participantId string) error {

	//check if BIC code is defined & the destination is world wire
	if msg.Message.Head.To == nil ||
		msg.Message.Head.To.FIId == nil ||
		msg.Message.Head.To.FIId.FinInstnId == nil ||
		msg.Message.Head.To.FIId.FinInstnId.Othr == nil ||
		!utils.StringsEqual(string(*msg.Message.Head.To.FIId.FinInstnId.BICFI), os.Getenv(environment.ENV_KEY_WW_BIC)) ||
		!utils.StringsEqual(string(*msg.Message.Head.To.FIId.FinInstnId.Othr.Id), os.Getenv(environment.ENV_KEY_WW_ID)) {
		return errors.New("Destination BIC or ID is undefined or incorrect")
	}

	//check if BIC code is defined & the source is OFI
	if msg.Message.Head.Fr == nil ||
		msg.Message.Head.Fr.FIId == nil ||
		msg.Message.Head.Fr.FIId.FinInstnId == nil ||
		msg.Message.Head.Fr.FIId.FinInstnId.Othr == nil ||
		!utils.StringsEqual(string(*msg.Message.Head.Fr.FIId.FinInstnId.BICFI), bic) ||
		!utils.StringsEqual(string(*msg.Message.Head.Fr.FIId.FinInstnId.Othr.Id), participantId) {
		return errors.New("Source BIC or ID is undefined or incorrect")
	}

	//check the format of both message identifier & business message identifier in the header
	if err := parse.HeaderIdentifierCheck(string(*msg.Message.Head.BizMsgIdr), string(*msg.Message.Head.MsgDefIdr), constant.IBWF001); err != nil {
		return err
	}

	//check the format & value of original instruction id
	for _, res := range msg.Message.Body.ResBody {
		if string(res.FedRes.PmtId.InstrId) != "" {
			if err := parse.InstructionIdCheck(string(res.FedRes.PmtId.InstrId)); err != nil {
				return err
			}
		} else {
			return errors.New("Original instruction ID is missing in the payload")
		}

		if string(res.CmpRes.PmtId.InstrId) != "" {
			if err := parse.InstructionIdCheck(string(res.CmpRes.PmtId.InstrId)); err != nil {
				return err
			}
		} else {
			return errors.New("Original instruction ID is missing in the payload")
		}
	}

	//check the format & value of instruction id
	if err := parse.InstructionIdCheck(string(msg.Message.Body.ResBody[0].Id)); err != nil {
		return errors.New("Instruction ID format is incorrect")
	}

	//check the BIC code and ID of both OFI & RFI is defined
	if msg.Message.Body.GrpHdr.InstdAgt.FinInstnId.Othr == nil ||
		msg.Message.Body.GrpHdr.InstdAgt.FinInstnId.Othr.Id == nil ||
		msg.Message.Body.GrpHdr.InstgAgt.FinInstnId.Othr == nil ||
		msg.Message.Body.GrpHdr.InstgAgt.FinInstnId.Othr.Id == nil {
		return errors.New("Instructing/Instructed agent ID is missing")
	}

	return nil
}

func (msg *Ibwf001) RequestToStruct() error {
	LOGGER.Infof("Constructing request data to go struct...")
	err := xml.Unmarshal(msg.Raw, &msg.Message)
	if err != nil {
		return err
	}

	LOGGER.Infof("Go struct constructed successfully")
	return nil
}

func (msg *Ibwf001) StructToProto() error {
	LOGGER.Infof("Constructing to protobuffer...")

	//putting raw xml message into the kafka send payload
	msg.SendPayload.Message = msg.Raw
	msg.SendPayload.MsgType = constant.ISO20022 + ":" + constant.IBWF001
	msg.SendPayload.OfiId = string(*msg.Message.Body.GrpHdr.InstdAgt.FinInstnId.Othr.Id)
	msg.SendPayload.RfiId = string(*msg.Message.Body.GrpHdr.InstgAgt.FinInstnId.Othr.Id)
	msg.SendPayload.InstructionId = string(msg.Message.Body.ResBody[0].Id)
	msg.SendPayload.OriginalInstructionId = string(msg.Message.Body.ResBody[0].FedRes.PmtId.InstrId)
	msg.SendPayload.OriginalMsgId = string(msg.Message.Body.ResBody[0].FedRes.PmtId.InstrId)
	msg.SendPayload.MsgId = string(msg.Message.Body.GrpHdr.MsgId)

	LOGGER.Infof("Protobuffer successfully constructed")
	return nil
}

func (msg *Ibwf001) ProtobuftoStruct() (*sendmodel.XMLData, error) {
	LOGGER.Infof("Restoring protobuffer to go struct...")

	err := xml.Unmarshal(msg.SendPayload.Message, &msg.Message)
	if err != nil {
		return nil, errors.New("Encounter error while unmarshaling xml message")
	}

	msgId := string(msg.Message.Body.GrpHdr.MsgId)
	instrId := string(msg.Message.Body.ResBody[0].Id)
	rfiAccountName := string(msg.Message.Body.GrpHdr.SttlmInf.SttlmAcct.Nm)
	rfiId := string(*msg.Message.Body.GrpHdr.InstgAgt.FinInstnId.Othr.Id)
	ofiId := string(*msg.Message.Body.GrpHdr.InstdAgt.FinInstnId.Othr.Id)
	ofiBIC := string(*msg.Message.Body.GrpHdr.InstdAgt.FinInstnId.BICFI)

	bizMsgIdr := msg.Message.Head.BizMsgIdr
	msgDefIdr := msg.Message.Head.MsgDefIdr
	creDt := msg.Message.Head.CreDt

	var federationStatus string
	var complianceInfoStatus string
	var complianceTransactionStatus string
	var rfiAccountPkey string

	for _, res := range msg.Message.Body.ResBody {
		complianceInfoStatus = string(res.CmpRes.InfSts)
		complianceTransactionStatus = string(res.CmpRes.TxnSts)
		rfiAccountPkey = string(res.FedRes.AccId)
		federationStatus = string(res.FedRes.FedSts)
	}

	pacs008InstructionId := string(msg.Message.Body.ResBody[0].FedRes.PmtId.InstrId)

	status, checkErr := checkFedAndCompResult(federationStatus, complianceInfoStatus, complianceTransactionStatus)
	if checkErr != nil && status == constant.REJECT_STRING {
		LOGGER.Errorf("RFI reject the transaction request")
		statusCode := strconv.Itoa(constant.STATUS_CODE_FED_COMP_RJCT)
		xmlData := &sendmodel.XMLData{}
		xmlData.OriginalMsgId = pacs008InstructionId
		xmlData.InstructionId = msgId
		xmlData.StatusCode = statusCode
		xmlData.OFIBIC = ofiBIC
		xmlData.RFIAccount = rfiAccountPkey

		return xmlData, checkErr
	}

	LOGGER.Infof("Restoring ibwf document")
	wwBicfi := msg.Message.Head.To.FIId.FinInstnId.BICFI
	wwId := msg.Message.Head.To.FIId.FinInstnId.Othr.Id
	ofiBicfi := msg.Message.Head.Fr.FIId.FinInstnId.BICFI
	id := msg.Message.Head.Fr.FIId.FinInstnId.Othr.Id

	msg.Message.Head = &ibwf.BusinessApplicationHeaderV01{}
	msg.Message.Head.Fr = &ibwf.Party9Choice{}
	msg.Message.Head.Fr.FIId = &ibwf.BranchAndFinancialInstitutionIdentification5{}
	msg.Message.Head.Fr.FIId.FinInstnId = &ibwf.FinancialInstitutionIdentification8{}
	msg.Message.Head.Fr.FIId.FinInstnId.Othr = &ibwf.GenericFinancialIdentification1{}

	msg.Message.Head.To = &ibwf.Party9Choice{}
	msg.Message.Head.To.FIId = &ibwf.BranchAndFinancialInstitutionIdentification5{}
	msg.Message.Head.To.FIId.FinInstnId = &ibwf.FinancialInstitutionIdentification8{}
	msg.Message.Head.To.FIId.FinInstnId.Othr = &ibwf.GenericFinancialIdentification1{}

	msg.Message.Head.Fr.FIId.FinInstnId.BICFI = ofiBicfi
	msg.Message.Head.Fr.FIId.FinInstnId.Othr.Id = id
	msg.Message.Head.To.FIId.FinInstnId.BICFI = wwBicfi
	msg.Message.Head.To.FIId.FinInstnId.Othr.Id = wwId

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
		RequestMsgType:           msg.SendPayload.MsgType,
		MessageId:                msgId,
		OFIId:                    ofiId,
		RFIId:                    rfiId,
		OriginalMsgId:            pacs008InstructionId,
		OriginalInstructionId:    pacs008InstructionId,
		InstructionId:            instrId,
		StatusCode:               strconv.Itoa(constant.STATUS_CODE_DEFAULT),
		RFIAccount:               rfiAccountPkey,
		RFISettlementAccountName: rfiAccountName,
		OFIBIC:                   ofiBIC,
	}

	LOGGER.Infof("Restoring protobuffer successfully")
	return xmlData, nil

}

func checkFedAndCompResult(fs, cis, cts string) (string, error) {
	if fs == constant.PAYMENT_STATUS_RJCT {
		return constant.REJECT_STRING, errors.New("federation not accepted")
	} else if cis == constant.PAYMENT_STATUS_RJCT {
		return constant.REJECT_STRING, errors.New("compliance not accepted")
	} else if cts == constant.PAYMENT_STATUS_RJCT {
		return constant.REJECT_STRING, errors.New("compliance not accepted")
	} else if fs == constant.PAYMENT_STATUS_ACTC && cis == constant.PAYMENT_STATUS_ACTC && cts == constant.PAYMENT_STATUS_ACTC {
		LOGGER.Infof("Federation and Compliance all accepted")
		return constant.ACCEPT_STRING, nil
	} else {
		return constant.REJECT_STRING, errors.New("wrong federation and compliance status")
	}
}
