// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package parse

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/GFTN/gftn-services/utility/payment/environment"

	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/types"
	"github.com/lestrrat-go/libxml2/xsd"
	"github.com/GFTN/gftn-services/utility/payment/constant"
	"github.com/GFTN/gftn-services/utility/payment/utils/sendmodel"
)

var V EndPointInput
var schema *xsd.Schema

type validate struct {
	MessageSchema *xsd.Schema
	XMLFile       types.Document
}

type EndPointInput struct {
	SchemaFile []byte
	Vars       sendmodel.SendVariables
	v          validate
}

//Basic struct used to test XML format
type Document struct {
	Format string `xml:"xmlns,attr"`
}

//Defines an ISO Message
//Can either be ISO 8583 or 20022
type ISOMessage interface {
	String() (result string, ok bool)
}

func SchemaInitiate(path []string) ([]*xsd.Schema, error) {
	var xsdSchemas []*xsd.Schema
	if len(path) == 0 {
		return nil, errors.New("Unable to read XSD path")
	}
	schema, readErr := V.SetSchema(path[0])
	if readErr != nil {
		LOGGER.Errorf("Failed to set schema: %v", readErr.Error())
		schema.Free()
		return nil, readErr
	}

	LOGGER.Infof("Read and set XSD.")
	xsdSchemas = append(xsdSchemas, schema)
	return xsdSchemas, nil
}

func init() {
	LOGGER.Infof("Initializing the XSD schema files")
	if os.Getenv(environment.ENV_KEY_SERVICE_FILE) == "" {
		LOGGER.Info("XSD Initialization failed as environment variable SERVICE_FILE not set")
	} else {
		LOGGER.Infof("Opening configuration file to get XSD path:%s", os.Getenv(environment.ENV_KEY_SERVICE_FILE))
		jsonFile, openErr := os.Open(os.Getenv(environment.ENV_KEY_SERVICE_FILE))
		if openErr != nil {
			LOGGER.Error("Opening of the configuration file failed")
			panic(openErr)
		}

		//UnMarshall the configuration file to struct
		byteData, readErr := ioutil.ReadAll(jsonFile)
		if readErr != nil {
			LOGGER.Errorf("Error while reading json configuration file", readErr)
			panic(readErr)
		}
		configVar := sendmodel.SendVariables{}
		json.Unmarshal(byteData, &configVar)
		s, err := xsd.ParseFromFile(configVar.XSDPath[0])
		if err != nil {
			LOGGER.Error("Parsing of the XSD schema file failed")
			panic(err)
		}

		//Set the schema variable, the variable will not be closed to ensure that we dont need to open the schema file every time
		schema = s
	}
}

func ValidateSchema(xml string) error {
	doc, err := libxml2.ParseString(xml)
	if err != nil {
		LOGGER.Warningf("Error while parsing the XML document", err)
		return err
	}
	if schema == nil {
		LOGGER.Error("schema object is null")
		return errors.New("schema validation failed because schema isn't initialized")
	}
	err = schema.Validate(doc)
	doc.Free()
	if err != nil {
		errs, ok := err.(xsd.SchemaValidationError)
		if ok {
			for _, e := range errs.Errors() {
				LOGGER.Error(e)
			}
		}
		LOGGER.Warningf("Error while validating against schema", err)
		return err
	}

	return nil
}

func (i *EndPointInput) SetSchema(path string) (*xsd.Schema, error) {
	byteValueXSD, err := ioutil.ReadFile(path)
	if err != nil {
		LOGGER.Errorf("Read Schema from path: %s failed: %v", path, err.Error())
		return nil, err
	}

	i.SchemaFile = byteValueXSD

	s, err := xsd.Parse(i.SchemaFile)
	if err != nil {
		LOGGER.Errorf("Failed to parse schema: %v", err.Error())
		return nil, err
	}

	i.v.MessageSchema = s

	return i.v.MessageSchema, nil
}

// BizMsgIdr Format : ByyyymmddbbbbbbbbbbbXAAnnnnnnn
func HeaderIdentifierCheck(bizMsgIdr string, msgDefIdr string, msgType string) error {
	msgDefIdr = strings.TrimSpace(strings.ToUpper(msgDefIdr))
	bizMsgIdr = strings.TrimSpace(strings.ToUpper(bizMsgIdr))
	msgType = strings.TrimSpace(strings.ToUpper(msgType))

	if msgDefIdr != msgType {
		return errors.New("MsgDefIdr does not match the actual message type")
	}

	r := regexp.MustCompile(`^B([12]\d{3}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01]))[A-Z0-9]{11}[B,H,G]{1}[A-Z]{2}[0-9]{7}$`)
	if !r.MatchString(bizMsgIdr) {
		return errors.New("BizMsgIdr value format is incorrect")
	}

	for _, n := range constant.SUPPORT_MESSAGE_TYPES {
		if msgDefIdr == strings.ToUpper(n) {
			return nil
		}
	}
	return errors.New("MsgDefIdr value format is incorrect")
}

func InstructionIdCheck(instrId string) error {
	instrId = strings.ToUpper(instrId)
	r := regexp.MustCompile(`^[A-Z]{3}(DO|DA|XX)([12]\d{3}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01]))[A-Z0-9]{11}[B,H,G]{1}[0-9]{10}$`)
	if !r.MatchString(instrId) {
		return errors.New("Instruction ID format is incorrect")
	}

	return nil

}
