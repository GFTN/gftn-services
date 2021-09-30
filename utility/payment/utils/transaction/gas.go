// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package transaction

import (
	"errors"
	"github.com/go-resty/resty"
)

func (op *CreateFundingOpereations) getSequenceAndIBMAccount() ([]byte, error){
	getResponse, restErr := resty.R().Get(op.gasServiceURL + "/lockaccount")
	if restErr != nil {
		LOGGER.Errorf("Error while getting the response from gas service: %v", restErr.Error())
		return nil, restErr
	}

	account := getResponse.Body()
	if len(account) == 0 {
		LOGGER.Error("Response from gas service is empty")
		return nil, errors.New("response from gas service is empty")
	}

	return account, nil
}
