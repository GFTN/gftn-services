// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package internal_model

import "net/http"

type DeploymentGlobal struct {
	Environment string `json:"environment"`
	ImageTag    string `json:"image_tag"`
	Replica     string `json:"replica"`
}

type DeploymentParticipant struct {
	Participants []Participant `json:"participants"`
	Environment  string        `json:"environment"`
	ImageTag     string        `json:"image_tag"`
}

type Participant struct {
	ID           string `json:"id"`
	CallBackURL  string `json:"call_back_url"`
	RDOClientURL string `json:"rdo_client_url"`
	Replica      string `json:"replica"`
}

type DeploymentKafka struct {
	Environment  string `json:"environment"`
	AWSRegion    string `json:"aws_region"`
}

type TearDownGlobal struct {
	Environment    string   `json:"environment"`
}

type TearDownParticipant struct {
	Participants []string `json:"participants"`
	Environment  string   `json:"environment"`
}

type TearDownKafka struct {
	Environment  string   `json:"environment"`
}

type RestartKafka struct {
	Environment string `json:"environment"`
	AWSRegion   string `json:"aws_region"`
}

type PasswordRotate struct {
	Environment  string   `json:"environment"`
}

type Client struct {
	HTTPClient *http.Client
	URL        string
}

type Maintain struct {
	ParticipantID []string `json:"participant_id"`
	Global        bool     `json:"global"`
}