// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package kafka

const (
	KAFKA_SSL                          = "ssl"
	ENV_KEY_KAFKA_ENABLE_SSL           = "KAFKA_ENABLE_SSL"
	ENV_KEY_KAFKA_CA_LOCATION          = "KAFKA_CA_LOCATION"
	ENV_KEY_KAFKA_CERTIFICATE_LOCATION = "KAFKA_CERTIFICATE_LOCATION"
	ENV_KEY_KAFKA_KEY_LOCATION         = "KAFKA_KEY_LOCATION"
	ENV_KEY_KAFKA_KEY_PASSWORD         = "KAFKA_KEY_PASSWORD"
	ENV_KEY_KAFKA_BROKER_URL           = "KAFKA_BROKER_URL"
	PAYMENT_TOPIC                      = "PAYMENT"
	QUOTES_TOPIC                       = "QUOTES"
	FEE_TOPIC                          = "FEE"
	TRANSACTION_TOPIC                  = "TRANSACTIONS"
	REQUEST_TOPIC                      = "_req"
	RESPONSE_TOPIC                     = "_res"
	ANCHOR_REDEMPTION_TOPIC            = "ANCHOR_REDEMPTION" + REQUEST_TOPIC
)

var SUPPORT_MESSAGE_TYPES = []string{PAYMENT_TOPIC, QUOTES_TOPIC, FEE_TOPIC, TRANSACTION_TOPIC}
