// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package whitelist_handler

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/op/go-logging"
	"github.com/GFTN/gftn-services/utility/payment/constant"

	"github.com/GFTN/gftn-services/gftn-models/model"
	"github.com/GFTN/gftn-services/global-whitelist-service/whitelistclient"
	comn "github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
)

var LOGGER = logging.MustGetLogger("whitlist-handler")

type ParticipantWhiteList struct {
	whiteListServiceURL string
	whiteListClient     whitelistclient.Client
}

func CreateWhiteListServiceOperations() (op ParticipantWhiteList) {
	op.whiteListServiceURL = os.Getenv(global_environment.ENV_KEY_WL_SVC_URL)

	wlc := whitelistclient.Client{
		HTTPClient: &http.Client{Timeout: time.Second * 30},
		WLURL:      op.whiteListServiceURL,
	}
	op.whiteListClient = wlc

	return op
}

func (op *ParticipantWhiteList) CheckWhiteListParticipant(inspector, inspected, settlementAccountName string) (string, error) {
	mutualWhiteList, whiteListErr := op.whiteListClient.GetMutualWhiteListParticipants(inspector)
	if whiteListErr != nil {
		return "", whiteListErr
	}

	if len(mutualWhiteList) == 0 {
		return "", nil
	} else if settlementAccountName == constant.EMPTY_STRING {
		LOGGER.Infof("Check if %s was whitelisted", inspected)
		result := findInspected(inspected, mutualWhiteList)
		if !result {
			return "", errors.New("no mutual whitelisted between OFI and RFI")
		}

		return "null", nil
	} else {
		LOGGER.Infof("Check if %s was whitelisted and get the settlement account %v", inspected, settlementAccountName)
		pkey, findErr := findInspectedInPR(inspected, settlementAccountName, mutualWhiteList)
		if !findErr {
			return "", errors.New("participant not register or is inactive")
		}

		return pkey, nil
	}

	return "", nil
}

func findInspected(inspected string, inspectorWhiteList []model.Participant) bool {
	found := false

	for _, p := range inspectorWhiteList {
		if *p.ID == inspected {
			if strings.ToLower(p.Status) == "active" {
				found = true
				return found
			}
		}
	}

	return found
}

func findInspectedInPR(rfiDomain, settlementAccountName string, prObject []model.Participant) (string, bool) {
	rfiInPR := false
	var accountAddress = ""

	for _, p := range prObject {
		if rfiDomain == *p.ID {
			if p.Status == "active" {
				if settlementAccountName == comn.ISSUING {
					accountAddress = p.IssuingAccount
					rfiInPR = true
					return accountAddress, rfiInPR
				} else {
					for _, oa := range p.OperatingAccounts {
						if oa.Name == settlementAccountName {
							accountAddress = *oa.Address
							rfiInPR = true
							return accountAddress, rfiInPR
						}
					}
				}
			}
		}
	}

	return accountAddress, rfiInPR
}
