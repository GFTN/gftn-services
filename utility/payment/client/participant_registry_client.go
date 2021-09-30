// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package client

import (
	pr_client "github.com/GFTN/gftn-services/participant-registry-client/pr-client"
	"github.com/GFTN/gftn-services/utility/common"
	"github.com/GFTN/gftn-services/utility/payment/constant"
)

func GetParticipantAccount(prServiceURL, homeDomain, queryStr string) *string {
	prc, prcErr := pr_client.CreateRestPRServiceClient(prServiceURL)
	if prcErr != nil {
		LOGGER.Error("Can not create connection to PR client service, please check if PR service is running")
		return nil
	}

	pr, prcGetErr := prc.GetParticipantForDomain(homeDomain)
	if prcGetErr != nil {
		LOGGER.Error("Could not found participant from PR service")
		return nil
	}

	if pr.Status == "active" {
		if queryStr == common.ISSUING {
			return &pr.IssuingAccount
		} else if queryStr == constant.BIC_STRING {
			return pr.Bic
		} else {
			for _, oa := range pr.OperatingAccounts {
				if oa.Name == queryStr {
					return oa.Address
				}
			}
		}
	} else {
		LOGGER.Errorf("Participant status is inactive")
		return nil
	}

	return nil
}

func GetParticipantRole(prServiceURL, homeDomain string) *string {
	prc, prcErr := pr_client.CreateRestPRServiceClient(prServiceURL)
	if prcErr != nil {
		LOGGER.Error("Can not create connection to PR client service, please check if PR service is running")
		return nil
	}

	pr, prcGetErr := prc.GetParticipantForDomain(homeDomain)
	if prcGetErr != nil {
		LOGGER.Error("Could not found participant from PR service")
		return nil
	}

	if pr.Status == "active" {
		return pr.Role
	} else {
		LOGGER.Errorf("Participant status is inactive")
		return nil
	}

	return nil
}
