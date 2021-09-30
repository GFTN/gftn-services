// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
export interface CredentialInfo {
	environment: string;
	domain: string;
	service: string;
	variable: string;
}

export interface ParameterContent {
	value: string;
	description: string;
}

export interface SecretContent {
	key?: string;
	value?: string;
	description: string;
	filePath?: string;  //absolute path
}

export const AWS_SECRET = "AWS"

export const LOCAL_SECRET = "LOCAL"

export const VAULT_SECRET = "VAULT"