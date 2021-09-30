// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package message_converter

import (
	"encoding/xml"
	"errors"
	"os"

	camt "github.com/GFTN/iso20022/camt08700106"
	pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/camt08700106"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/payment/client"
	"github.com/GFTN/gftn-services/utility/payment/constant"
	"github.com/GFTN/gftn-services/utility/payment/environment"
	"github.com/GFTN/gftn-services/utility/payment/utils"
	"github.com/GFTN/gftn-services/utility/payment/utils/parse"
	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"
)

type Camt087 struct {
	Message     *camt.Message
	SendPayload pbstruct.SendPayload
	Raw         []byte
}

func (msg *Camt087) SanityCheck(bic, participantId string) error {

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
	if err := parse.HeaderIdentifierCheck(string(*msg.Message.Head.BizMsgIdr), string(*msg.Message.Head.MsgDefIdr), constant.CAMT087); err != nil {
		return err
	}

	//check the format & value of original instruction id
	if msg.Message.Body.Undrlyg.Initn != nil && msg.Message.Body.Undrlyg.Initn.OrgnlInstrId != nil {
		if err := parse.InstructionIdCheck(string(*msg.Message.Body.Undrlyg.Initn.OrgnlInstrId)); err != nil {
			return err
		}
	} else {
		return errors.New("Original instruction ID is missing in the payload")
	}

	//check the format & value of instruction id
	if msg.Message.Body != nil && msg.Message.Body.Assgnmt != nil && msg.Message.Body.Assgnmt.Id != nil {
		if err := parse.InstructionIdCheck(string(*msg.Message.Body.Assgnmt.Id)); err != nil {
			return err
		}
	} else {
		return errors.New("Instruction ID is missing in the payload")
	}

	//check the BIC code and ID of both OFI & RFI is defined
	if msg.Message.Body.Assgnmt.Assgnr.Agt.FinInstnId.Othr == nil ||
		msg.Message.Body.Assgnmt.Assgnr.Agt.FinInstnId.Othr.Id == nil ||
		msg.Message.Body.Assgnmt.Assgne.Agt.FinInstnId.Othr == nil ||
		msg.Message.Body.Assgnmt.Assgne.Agt.FinInstnId.Othr.Id == nil {
		return errors.New("Assigner/Assignee ID is missing")
	}

	return nil
}

func (msg *Camt087) RequestToStruct() error {
	LOGGER.Infof("Constructing request data to go struct...")
	err := xml.Unmarshal(msg.Raw, &msg.Message)
	if err != nil {
		return err
	}

	LOGGER.Infof("Go struct constructed successfully")
	return nil
}

func (msg *Camt087) StructToProto() error {
	LOGGER.Infof("Constructing to protobuffer...")

	//putting raw xml message into the kafka send payload
	msg.SendPayload.Message = msg.Raw
	msg.SendPayload.MsgType = constant.ISO20022 + ":" + constant.CAMT087
	msg.SendPayload.OfiId = string(*msg.Message.Body.Assgnmt.Assgnr.Agt.FinInstnId.Othr.Id)
	msg.SendPayload.RfiId = string(*msg.Message.Body.Assgnmt.Assgne.Agt.FinInstnId.Othr.Id)
	msg.SendPayload.InstructionId = string(*msg.Message.Body.Assgnmt.Id)
	msg.SendPayload.OriginalInstructionId = string(*msg.Message.Body.Undrlyg.Initn.OrgnlInstrId)
	msg.SendPayload.OriginalMsgId = string(*msg.Message.Body.Undrlyg.Initn.OrgnlGrpInf.OrgnlMsgId)
	msg.SendPayload.MsgId = string(*msg.Message.Body.Assgnmt.Id)

	LOGGER.Infof("Protobuffer successfully constructed")
	return nil

}

func (msg *Camt087) ProtobuftoStruct() (*sendmodel.XMLData, error) {
	LOGGER.Infof("Restoring protobuffer to go struct...")

	err := xml.Unmarshal(msg.SendPayload.Message, &msg.Message)
	if err != nil {
		return nil, errors.New("Encounter error while unmarshaling xml message")
	}

	msgId := string(*msg.Message.Body.Assgnmt.Id)
	originalMsgId := string(*msg.Message.Body.Undrlyg.Initn.OrgnlGrpInf.OrgnlMsgId)
	oringalInstructionId := string(*msg.Message.Body.Undrlyg.Initn.OrgnlInstrId)
	ofiId := string(*msg.Message.Body.Assgnmt.Assgnr.Agt.FinInstnId.Othr.Id)
	rfiId := string(*msg.Message.Body.Assgnmt.Assgne.Agt.FinInstnId.Othr.Id)
	bizMsgIdr := msg.Message.Head.BizMsgIdr
	msgDefIdr := msg.Message.Head.MsgDefIdr
	creDt := msg.Message.Head.CreDt

	if rfiId != os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME) {
		LOGGER.Error("Wrong participant id for rfi")
		return nil, errors.New("Wrong participant id for rfi")
	}

	LOGGER.Infof("Retrieving rfi BIC code from participant registry")
	rfiBic := client.GetParticipantAccount(os.Getenv(global_environment.ENV_KEY_PARTICIPANT_REGISTRY_URL), rfiId, constant.BIC_STRING)
	wwBicfi := camt.BICFIIdentifier(os.Getenv(environment.ENV_KEY_WW_BIC))
	wwId := camt.Max35Text(os.Getenv(environment.ENV_KEY_WW_ID))
	rfiBicfi := camt.BICFIIdentifier(*rfiBic)
	id := camt.Max35Text(rfiId)

	LOGGER.Infof("Updating Business Application Header at RFI side")
	msg.Message.Head = &camt.BusinessApplicationHeaderV01{}
	msg.Message.Head.Fr = &camt.Party9Choice{}
	msg.Message.Head.Fr.FIId = &camt.BranchAndFinancialInstitutionIdentification5{}
	msg.Message.Head.Fr.FIId.FinInstnId = &camt.FinancialInstitutionIdentification8{}
	msg.Message.Head.Fr.FIId.FinInstnId.Othr = &camt.GenericFinancialIdentification1{}

	msg.Message.Head.To = &camt.Party9Choice{}
	msg.Message.Head.To.FIId = &camt.BranchAndFinancialInstitutionIdentification5{}
	msg.Message.Head.To.FIId.FinInstnId = &camt.FinancialInstitutionIdentification8{}
	msg.Message.Head.To.FIId.FinInstnId.Othr = &camt.GenericFinancialIdentification1{}

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
		RequestXMLMsg:         xmlMsg,
		MessageId:             msgId,
		OriginalMsgId:         originalMsgId,
		InstructionId:         msgId,
		OriginalInstructionId: oringalInstructionId,
		RequestMsgType:        msg.SendPayload.MsgType,
		OFIId:                 ofiId,
		RFIId:                 rfiId,
	}
	LOGGER.Infof("Restoring protobuffer successfully")

	return xmlData, nil
}
