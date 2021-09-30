// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package constant

var (
	ACCEPT_STRING                        = "ACCEPT"
	REJECT_STRING                        = "REJECT"
	DENIED_STRING                        = "denied"
	FALSE_STRING                         = "false"
	ERROR_STRING                         = "res_err"
	SETTLEMENT_METHOD_DIGITAL_ASSET      = "DA"
	SETTLEMENT_METHOD_DIGITAL_OBLIGATION = "DO"
	SETTLEMENT_METHOD_XLM                = "XLM"
	KAFKA_ENABLE_SSL_TRUE                = "true"
	KAFKA_ENABLE_SSL_FALSE               = "false"
	KAFKA_SSL                            = "ssl"
	KAFKA_INITIAL_ERROR                  = "kafka consumer initial error"
	KAFKA_INITIAL_SUCCESS                = "kafka consumer initial success"
	KAFKA_CONSUMER_RECONNECT             = "kafka consumer reconnect to the broker"
	KAFKA_PRODUCE_ERROR                  = "encounter error while producing message to Kafka"
	SERVICE_LOG_TOPIC                    = "service_log"
	SERVICE_LOG_ACTION_SEND_OFI          = "Send-OFI"
	SERVICE_LOG_ACTION_SEND_RFI          = "Send-RFI"
	SERVICE_LOG_ACTION_REPLY_OFI         = "Reply-OFI"
	SERVICE_LOG_ACTION_REPLY_RFI         = "Reply-RFI"
	ENV_DEV_STRING                       = "dev"
	OK_STRING                            = "OK"
	TX_DONE                              = "DONE"
	EMPTY_STRING                         = "EMPTY"
	BIC_STRING                           = "BIC"
	UNKNOWN_STRING                       = "UNKNOWN"
	XLM_STRING                           = "XLM"
	DEFAULT_STRING                       = "default"
	LOW_THRESHOLD                        = uint32(1)
	MEDIUM_THRESHOLD                     = uint32(2)
	HIGH_THRESHOLD                       = uint32(3)
	MASTER_WEIGHT                        = uint32(2)
	WW_ADMIN_WEIGHT                      = uint32(1)
	SHA_WEIGHT                           = uint32(2)
	ISO20022                             = "iso20022"
	ISO8385                              = "iso8385"
	JSON                                 = "json"
	MT                                   = "mt"
	PACS008                              = "pacs.008.001.07"
	PACS004                              = "pacs.004.001.09"
	CAMT056                              = "camt.056.001.08"
	CAMT026                              = "camt.026.001.07"
	CAMT029                              = "camt.029.001.09"
	IBWF001                              = "ibwf.001.001.01"
	IBWF002                              = "ibwf.002.001.01"
	PACS002                              = "pacs.002.001.09"
	CAMT030                              = "camt.030.001.05"
	PACS009                              = "pacs.009.001.08"
	CAMT087                              = "camt.087.001.06"
	SUPPORT_MESSAGE_TYPES                = []string{PACS008, PACS004, CAMT056, CAMT029, IBWF001, IBWF002, PACS002, CAMT030, PACS009, PACS002, CAMT026, CAMT087}
	SUPPORT_CAMT_MESSAGES                = []string{CAMT056, CAMT029, CAMT026, CAMT087}
	REQUEST                              = "REQUEST"
	RESPONSE                             = "RESPONSE"
	REDEEM                               = "REDEEM"
	WWBIC                                = "WORLDWIRE00"
	WWID                                 = "WW"
	REASON_CODE_PAYMENT_CANCELLATION     = 1000
	REASON_CODE_RDO                      = 1001
	DO_SETTLEMENT                        = "WWDO"
	DA_SETTLEMENT                        = "WWDA"
	XLM_SETTLEMENT                       = "XLM"
	WWCCY                                = "XXXXX"
)
