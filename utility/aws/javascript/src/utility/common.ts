// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import * as Utility from "./var"

export function getCredentialId(credential: Utility.CredentialInfo) {
  if (!process.env.AWS_ACCESS_KEY_ID || !process.env.AWS_SECRET_ACCESS_KEY || !process.env.AWS_REGION) {
    return new Error("Cannot fetch the correct AWS session config, please check that you have set access key ID/secret key/region correctly");
  }
  if (!credential.environment || !credential.domain || !credential.service || !credential.variable) {
    return new Error('Credential parameters missing');
  }
  let credentialId = `/${credential.environment}/${credential.domain}/${credential.service}/${credential.variable}`;
  return credentialId;
}