// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package client

import (
	"github.com/go-errors/errors"
	"github.com/go-resty/resty"
	"github.com/GFTN/gftn-services/utility/global-environment"
	participant_util "github.com/GFTN/gftn-services/utility/participant"
	"os"
)

type RestPaymentListenerClient struct {
	PaymentListenerURL string
}

func CreateRestPaymentListenerClient() (RestPaymentListenerClient, error) {

	client := RestPaymentListenerClient{}
	url, _ := participant_util.GetServiceUrl(os.Getenv(global_environment.ENV_KEY_PAYMENT_SVC_URL), os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME))
	if url == "" {
		LOGGER.Errorf("You MUST set the %v environment variable to point to this participant's payment Listener",
			global_environment.ENV_KEY_PAYMENT_SVC_URL)
		os.Exit(1)
	}
	client.PaymentListenerURL = url

	return client, nil

}

func (client RestPaymentListenerClient) SubscribePayments(distAccountName string) (err error) {
	// API Service =>> Payment Listener

	if distAccountName == "" {
		LOGGER.Debugf("SubscribePayments called with empty Dist account")
		return nil
	}
	///client/accounts/{account_name}/{cursor}
	url := client.PaymentListenerURL + "/internal/accounts/" + distAccountName + "/now"

	response, err := resty.R().Post(url)
	LOGGER.Infof("Start Payment listener for Operating account:  %v", distAccountName)

	if err != nil {
		LOGGER.Errorf("Error Starting Payment listener for Operating account:  %v, error: %v", distAccountName, err.Error())
		return err
	}

	if response.StatusCode() != 200 {
		LOGGER.Errorf("Error Starting Payment listener for Operating account:  %v", distAccountName)
		return errors.New("Error Starting Payment listener for Operating")
	}

	return nil
}
