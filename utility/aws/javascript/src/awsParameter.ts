// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
/*
Parameter Naming Constraints:
	* Parameter names are case sensitive.
	* A parameter name must be unique within an AWS Region
	* A parameter name can't be prefixed with "aws" or "ssm" (case-insensitive).
	* Parameter names can include only the following symbols and letters: a-zA-Z0-9_.-/
	* A parameter name can't include spaces.
	* Parameter hierarchies are limited to a maximum depth of fifteen levels.
*/

import * as Var from './utility/var';
import * as Common from './utility/common';
import {SSM} from 'aws-sdk';

export function getParameter(credentialInfo: Var.CredentialInfo) {
  return new Promise((res, rej) => {
    let credentialId: any;
    credentialId = Common.getCredentialId(credentialInfo);
    if (credentialId instanceof Error) {
      rej(credentialId);
    }
    console.log("getParameter: " + credentialId);
    const ssm = new SSM();
    const params = {
      Name: credentialId,
      WithDecryption: true
    };

    ssm.getParameter(params, function (err, data) {
      if (err) {
        console.error("Error getting parameter");
        rej(err);
      } else {
        console.log(`${credentialId} successfully retrieved`);
        const resString = JSON.stringify(data);
        const resJson = JSON.parse(resString);
        res(resJson.Parameter.Value);
      }
    });
  });
}

export async function createParameter(credentialInfo: Var.CredentialInfo, parameterContent: Var.ParameterContent) {
  return putParameter(credentialInfo, parameterContent, false);
}

export async function updateParameter(credentialInfo: Var.CredentialInfo, parameterContent: Var.ParameterContent) {
  return putParameter(credentialInfo, parameterContent, true);
}

export function removeParameter(credentialInfo: Var.CredentialInfo) {
  return new Promise((res, rej) => {
    let credentialId: any;
    credentialId = Common.getCredentialId(credentialInfo);
    if (credentialId instanceof Error) {
      rej(credentialId);
    }
    console.log(`removeParameter: ${credentialId}`);
    const ssm = new SSM();
    const params = {
      Name: credentialId
    };

    ssm.deleteParameter(params, function (err, data) {
      if (err) {
        console.error("Error deleting parameter");
        rej(err);
      } else {
        console.log(`${credentialId} successfully removed`);
        res("success");
      }
    });
  });
}


function putParameter(credentialInfo: Var.CredentialInfo, parameterContent: Var.ParameterContent, overwrite: boolean) {
  return new Promise((res, rej) => {
    let credentialId: any;
    credentialId = Common.getCredentialId(credentialInfo);
    if (credentialId instanceof Error) {
      rej(credentialId);
    }
    if (overwrite) {
      console.info(`updateParameter: ${credentialId}`);
    } else {
      console.info(`createParameter: ${credentialId}`);
    }
    const ssm = new SSM();
    const params = {
      Name: credentialId,
      Type: "SecureString",
      Value: parameterContent.value,
      Description: parameterContent.description,
      Overwrite: overwrite
    };

    ssm.putParameter(params, function (err, data) {
      if (err) {
        console.error("Error creating/updating parameter");
        rej(err);
      } else {
        console.log(`${credentialId} successfully added`);
        const resString = JSON.stringify(data);
        const resJson = JSON.parse(resString);
        res(resJson);
      }
    });
  });

}