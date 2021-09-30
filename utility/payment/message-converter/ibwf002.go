// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package message_converter

import (
	"encoding/xml"
	"errors"
	"os"

	ibwf "github.com/GFTN/iso20022/ibwf00200101"
	pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/ibwf00200101"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/payment/client"
	"github.com/GFTN/gftn-services/utility/payment/constant"
	"github.com/GFTN/gftn-services/utility/payment/environment"
	"github.com/GFTN/gftn-services/utility/payment/utils"
	"github.com/GFTN/gftn-services/utility/payment/utils/parse"
	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"
)

type Ibwf002 struct {
	Message     *ibwf.Message
	SendPayload pbstruct.SendPayload
	Raw         []byte
}

func (msg *Ibwf002) SanityCheck(bic, participantId string) error {

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
	if err := parse.HeaderIdentifierCheck(string(*msg.Message.Head.BizMsgIdr), string(*msg.Message.Head.MsgDefIdr), constant.IBWF002); err != nil {
		return err
	}

	for _, sttlOblInf := range msg.Message.DigOblSetNotif.SttlOblInf {
		for _, txInfo := range sttlOblInf.TxInf {
			//check the format & value of original instruction id
			if txInfo.OrgnlInstrId != nil {
				if err := parse.InstructionIdCheck(string(*txInfo.OrgnlInstrId)); err != nil {
					return err
				}
			} else {
				return errors.New("Original instruction ID is missing in the payload")
			}

			//check the format & value of instruction id
			if txInfo.NtfId != nil {
				if err := parse.InstructionIdCheck(string(*txInfo.NtfId)); err != nil {
					return err
				}
			} else {
				return errors.New("Instruction ID is missing in the payload")
			}
		}
	}

	//check the BIC code and ID of both OFI & RFI is defined
	if msg.Message.DigOblSetNotif.GrpHdr.InstdAgt.FinInstnId.Othr == nil ||
		msg.Message.DigOblSetNotif.GrpHdr.InstdAgt.FinInstnId.Othr.Id == nil ||
		msg.Message.DigOblSetNotif.GrpHdr.InstgAgt.FinInstnId.Othr == nil ||
		msg.Message.DigOblSetNotif.GrpHdr.InstgAgt.FinInstnId.Othr.Id == nil {
		return errors.New("Instructing/Instructed agent ID is missing")
	}

	return nil
}

func (msg *Ibwf002) RequestToStruct() error {
	LOGGER.Infof("Constructing request data to go struct...")
	err := xml.Unmarshal(msg.Raw, &msg.Message)
	if err != nil {
		return err
	}

	LOGGER.Infof("Go struct constructed successfully")
	return nil
}

func (msg *Ibwf002) StructToProto() error {
	LOGGER.Infof("Constructing to protobuffer..")

	//putting raw xml message into the kafka send payload
	msg.SendPayload.Message = msg.Raw
	msg.SendPayload.MsgType = constant.ISO20022 + ":" + constant.IBWF002
	msg.SendPayload.OfiId = string(*msg.Message.DigOblSetNotif.GrpHdr.InstgAgt.FinInstnId.Othr.Id)
	msg.SendPayload.RfiId = string(*msg.Message.DigOblSetNotif.GrpHdr.InstdAgt.FinInstnId.Othr.Id)
	msg.SendPayload.InstructionId = string(*msg.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].NtfId)
	msg.SendPayload.OriginalInstructionId = string(*msg.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlInstrId)
	msg.SendPayload.OriginalMsgId = string(*msg.Message.DigOblSetNotif.SttlOblInf[0].OrgnlGrpInf.OrgnlMsgId)
	msg.SendPayload.MsgId = string(*msg.Message.DigOblSetNotif.GrpHdr.MsgId)

	LOGGER.Infof("Protobuffer successfully constructed")
	return nil
}

func (msg *Ibwf002) ProtobuftoStruct() (*sendmodel.XMLData, error) {
	LOGGER.Infof("Restoring protobuffer to go struct...")

	err := xml.Unmarshal(msg.SendPayload.Message, &msg.Message)
	if err != nil {
		return nil, errors.New("Encounter error while unmarshaling xml message")
	}

	msgId := string(*msg.Message.DigOblSetNotif.GrpHdr.MsgId)
	ofiId := string(*msg.Message.DigOblSetNotif.GrpHdr.InstgAgt.FinInstnId.Othr.Id)
	rfiId := string(*msg.Message.DigOblSetNotif.GrpHdr.InstdAgt.FinInstnId.Othr.Id)
	bizMsgIdr := msg.Message.Head.BizMsgIdr
	msgDefIdr := msg.Message.Head.MsgDefIdr
	creDt := msg.Message.Head.CreDt

	if rfiId != os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME) {
		LOGGER.Error("Wrong participant id for rfi")
		return nil, errors.New("Wrong participant id for rfi")
	}

	settlementAccountName := string(*msg.Message.DigOblSetNotif.GrpHdr.SttlmInf.SttlmAcct.Nm)

	LOGGER.Infof("Retrieving rfi BIC code from participant registry")
	rfiBic := client.GetParticipantAccount(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL), rfiId, constant.BIC_STRING)
	wwBicfi := ibwf.BICFIIdentifier(os.Getenv(environment.ENV_KEY_WW_BIC))
	wwId := ibwf.Max35Text(os.Getenv(environment.ENV_KEY_WW_ID))
	rfiBicfi := ibwf.BICFIIdentifier(*rfiBic)
	id := ibwf.Max35Text(rfiId)

	LOGGER.Infof("Updating Business Application Header at RFI side")
	msg.Message.Head = &ibwf.BusinessApplicationHeaderV01{}
	msg.Message.Head.Fr = &ibwf.Party9Choice{}
	msg.Message.Head.Fr.FIId = &ibwf.BranchAndFinancialInstitutionIdentification5{}
	msg.Message.Head.Fr.FIId.FinInstnId = &ibwf.FinancialInstitutionIdentification8{}
	msg.Message.Head.Fr.FIId.FinInstnId.Othr = &ibwf.GenericFinancialIdentification1{}

	msg.Message.Head.To = &ibwf.Party9Choice{}
	msg.Message.Head.To.FIId = &ibwf.BranchAndFinancialInstitutionIdentification5{}
	msg.Message.Head.To.FIId.FinInstnId = &ibwf.FinancialInstitutionIdentification8{}
	msg.Message.Head.To.FIId.FinInstnId.Othr = &ibwf.GenericFinancialIdentification1{}

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
		RequestMsgType:           msg.SendPayload.MsgType,
		MessageId:                msgId,
		OFIId:                    ofiId,
		OFISettlementAccountName: settlementAccountName,
		OriginalEndtoEndId:       string(*msg.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlEndToEndId),
		OriginalMsgId:            string(*msg.Message.DigOblSetNotif.SttlOblInf[0].OrgnlGrpInf.OrgnlMsgId),
		OriginalInstructionId:    string(*msg.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].OrgnlInstrId),
		InstructionId:            string(*msg.Message.DigOblSetNotif.SttlOblInf[0].TxInf[0].NtfId),
	}

	LOGGER.Infof("Restoring protobuffer successfully")
	return xmlData, nil

}
