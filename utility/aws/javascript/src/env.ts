// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import * as SM from './awsSecret';
import * as Var from './utility/var';

export function CheckVariable() {
    if (!process.env.HOME_DOMAIN_NAME || !process.env.SERVICE_NAME || !process.env.ENV_VERSION || !process.env.SECRET_STORAGE_LOCATION) {
        throw new Error(("Initializing failed, Require the following environment variables to start up the service. HOME_DOMAIN_NAME, SERVICE_NAME, ENV_VERSION, SECRET_STORAGE_LOCATION"));
    }
}

function GetEnv(credential: Var.CredentialInfo) {
    return new Promise(async (resolve, rej) => {

        if (process.env[credential.variable] && credential.variable !== "initialize") {
            console.log("env already been initialized");
        } else {
            let res: any;
            try {
                res = await SM.getSecret(credential);
            } catch (e) {
                rej(e);
                return;
            }
            const obj = JSON.parse(res);
            const keys = Object.keys(obj);
            for (let i = 0; i < keys.length; i++) {
                process.env[keys[i]] = obj[keys[i]];
            }
            process.env[credential.variable] = "true";
        }

        resolve("success");

    });
}

export function InitEnv() {
    return new Promise(async (res, rej) => {
        if (process.env["SECRET_STORAGE_LOCATION"] === Var.AWS_SECRET) {
            console.log("Initializing service with AWS");

            //participant
            console.log("Initializing participant-specific secrets from AWS");
            const domainId = process.env.HOME_DOMAIN_NAME,
                svcName = process.env.SERVICE_NAME,
                envVersion = process.env.ENV_VERSION;

            const participant_credential: Var.CredentialInfo = {
                environment: envVersion,
                domain: domainId,
                service: "participant",
                variable: "initialize"
            };
            try {
                await GetEnv(participant_credential);
            } catch (e) {
                console.log(e);
                rej(e);
            }

            //service
            console.log("Initializing service-specific secrets from AWS");
            const service_credential: Var.CredentialInfo = {
                environment: envVersion,
                domain: domainId,
                service: svcName,
                variable: "initialize"
            };
            try {
                await GetEnv(service_credential);
            } catch (e) {
                console.log(e);
                rej(e);
            }
        }
        res("success");
    });


}
