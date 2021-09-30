// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package message_converter

import (
	"github.com/op/go-logging"
	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"
)

var LOGGER = logging.MustGetLogger("message-converter")

type MessageInterface interface {
	//unmarshaling http request to go struct
	RequestToStruct() error
	//marshaling go struct to proto buffer
	StructToProto() error
	//restoring proto buffer back to go struct
	ProtobuftoStruct() (*sendmodel.XMLData, error)
	//XML msg payload format & value check
	SanityCheck(string, string) error
}
