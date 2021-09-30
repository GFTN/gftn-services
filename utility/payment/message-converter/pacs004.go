// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package message_converter

import (
	"encoding/xml"
	"errors"
	"os"

	"github.com/GFTN/gftn-services/utility/payment/constant"
	"github.com/GFTN/gftn-services/utility/payment/utils"

	pacs "github.com/GFTN/iso20022/pacs00400109"
	pbstruct "github.com/GFTN/iso20022/proto/github.com/GFTN/iso20022/pacs00400109"
	"github.com/GFTN/gftn-services/utility/payment/environment"
	"github.com/GFTN/gftn-services/utility/payment/utils/parse"
	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"
)

type Pacs004 struct {
	Message     *pacs.Message
	SendPayload pbstruct.SendPayload
	Raw         []byte
}

func (msg *Pacs004) SanityCheck(bic, participantId string) error {

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
	if err := parse.HeaderIdentifierCheck(string(*msg.Message.Head.BizMsgIdr), string(*msg.Message.Head.MsgDefIdr), constant.PACS004); err != nil {
		return err
	}

	for _, txInfo := range msg.Message.Body.TxInf {
		//check the format & value of original instruction id
		if txInfo.OrgnlInstrId != nil {
			if err := parse.InstructionIdCheck(string(*txInfo.OrgnlInstrId)); err != nil {
				return err
			}
		} else {
			return errors.New("Original instruction ID is missing in the payload")
		}

		//check the format & value of instruction id
		if txInfo.RtrId != nil {
			if err := parse.InstructionIdCheck(string(*txInfo.RtrId)); err != nil {
				return err
			}
		} else {
			return errors.New("Instruction ID is missing in the payload")
		}
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

func (msg *Pacs004) RequestToStruct() error {
	LOGGER.Infof("Constructing request data to go struct...")
	err := xml.Unmarshal(msg.Raw, &msg.Message)
	if err != nil {
		return err
	}

	LOGGER.Infof("Go struct constructed successfully")
	return nil
}

func (msg *Pacs004) StructToProto() error {
	return nil
}

func (msg *Pacs004) ProtobuftoStruct() (*sendmodel.XMLData, error) {
	return nil, nil
}
