// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package utils

//var BaseURL = "https://3.0.15.221"

//var certPath = "/vagrant/go/src/github.com/GFTN/gftn-services/utility/vault/certs/certificate.crt"
//var keyPath = "/vagrant/go/src/github.com/GFTN/gftn-services/utility/vault/certs/privateKey.key"

type Session struct{
	CyberArkLogonResult string
	BaseURL string
	CertPath string
	KeyPath string
}

type Accounts struct{
	Count uint
	Accounts []Account
}

type Account struct{
	AccountID string
	InternalProperties []InternalProperties
	Properties []Properties
}

type Properties struct{
	Key string
	Value string
}

type InternalProperties struct{
	Key string
	Value string
}

type SafeList struct{
	List []Safe `json:"GetSafesResult"`
}

type Safe struct{
	Description string
	ManagingCPM string
	NumberOfDaysRetention uint
	NumberOfVersionsRetention uint
	OLACEnabled bool
	SafeName string
}

type SafeMemberList struct{
	Members []Member
}

type Member struct{
	Permissions interface{}
	UserName string
}

type Secret struct{
	Content string
	PolicyID string
	CreationMethod string
	Folder string
	Address string
	Name string
	Safe string
	DeviceType string
	UserName string
	PasswordChangeInProcess string
}
