// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
"use strict";
/*
Secret Naming Reminder:
    The secret name must be ASCII letters, digits, or the following characters: /_+=.@-
    Don't end your secret name with a hyphen followed by six characters. If you
    do so, you risk confusion and unexpected results when searching for a secret
    by partial ARN. This is because Secrets Manager automatically adds a hyphen
    and six random characters at the end of the ARN.
*/
Object.defineProperty(exports, "__esModule", { value: true });
var Common = require("./utility/common");
var fs = require("fs");
var aws_sdk_1 = require("aws-sdk");
/*
    IAM user required permission to call GetSecret function:
        * secretsmanager:GetSecretValue

        * kms:Decrypt - required only if you use a customer-managed AWS KMS key
        to encrypt the secret. You do not need this permission to use the account's
        default AWS managed CMK for Secrets Manager.
*/
function getSecret(credentialInfo) {
    return new Promise(function (res, rej) {
        var credentialId;
        credentialId = Common.getCredentialId(credentialInfo);
        if (credentialId instanceof Error) {
            rej(credentialId);
        }
        console.info("retrieving secret: " + credentialId);
        var client = new aws_sdk_1.SecretsManager();
        client.getSecretValue({ SecretId: credentialId }, function (err, data) {
            if (err) {
                console.error("Error getting secret " + credentialId);
                rej(err.message);
            }
            else {
                console.info("retrieving secret: Success!");
                // Decrypts secret using the associated KMS CMK.
                // Depending on whether the secret is a string or binary, one of these fields will be populated.
                if ('SecretString' in data) {
                    var secret = data.SecretString;
                    res(secret);
                }
                else {
                    var buff = new Buffer(data.SecretBinary, 'base64');
                    var decodedBinarySecret = buff.toString('ascii');
                    res(decodedBinarySecret);
                }
            }
        });
    });
}
exports.getSecret = getSecret;
/*
    IAM user required permission to call UpdateSecret function:
        * secretsmanager:UpdateSecret

        * kms:GenerateDataKey - needed only if you use a custom AWS KMS key to
        encrypt the secret. You do not need this permission to use the account's
        AWS managed CMK for Secrets Manager.

        * kms:Decrypt - needed only if you use a custom AWS KMS key to encrypt
        the secret. You do not need this permission to use the account's AWS managed
        CMK for Secrets Manager.
*/
function updateSecret(credentialInfo, secretContent) {
    return new Promise(function (res, rej) {
        var credentialId;
        credentialId = Common.getCredentialId(credentialInfo);
        if (credentialId instanceof Error) {
            rej(credentialId);
        }
        console.info("updating secret: " + credentialId);
        var secretString = getSecretString(secretContent);
        if (secretString instanceof Error) {
            rej(secretString);
        }
        var params = {
            SecretId: credentialId,
            SecretString: secretString
        };
        var client = new aws_sdk_1.SecretsManager();
        client.putSecretValue(params, function (err, data) {
            if (err) {
                console.error("Error updating secret");
                rej(err.message);
            }
            else {
                console.info("updating secret: Success!");
                res(data);
            }
        });
    });
}
exports.updateSecret = updateSecret;
/*
    IAM user required permission to call CreateSecret function:
        * secretsmanager:CreateSecret

        * kms:GenerateDataKey - needed only if you use a customer-managed AWS
        KMS key to encrypt the secret. You do not need this permission to use
        the account's default AWS managed CMK for Secrets Manager.

        * kms:Decrypt - needed only if you use a customer-managed AWS KMS key
        to encrypt the secret. You do not need this permission to use the account's
        default AWS managed CMK for Secrets Manager.

        * secretsmanager:TagResource - needed only if you include the Tags parameter.
*/
function createSecret(credentialInfo, secretFilePath, secretJSONString) {
    return new Promise(function (res, rej) {
        var credentialId;
        credentialId = Common.getCredentialId(credentialInfo);
        if (credentialId instanceof Error) {
            rej(credentialId);
        }
        console.info("adding secret: " + credentialId);
        // default to secret string passed in
        var _secretString = secretJSONString;
        // if path provided the override the secretJSONString
        var description = '-'; // default description to '-' dash so that it doesn't error when secretFilePath is null
        if (secretFilePath) {
            description = secretFilePath.description;
            _secretString = getSecretString(secretFilePath);
            if (_secretString instanceof Error) {
                rej(_secretString);
            }
        }
        var params = {
            Description: description,
            Name: credentialId,
            SecretString: _secretString
        };
        var client = new aws_sdk_1.SecretsManager();
        client.createSecret(params, function (err, data) {
            if (err) {
                console.error("Error creating secret");
                rej(err.message);
            }
            else {
                console.info("adding secret: Success!");
                res(data);
            }
        });
    });
}
exports.createSecret = createSecret;
/*
    IAM user required permission to call DeleteSecret function:
        * secretsmanager:DeleteSecret
    Note:
        recoveryDays should be 7 days at minimum
*/
function removeSecret(credentialInfo, forceDeleteWithoutRecovery) {
    return new Promise(function (res, rej) {
        var credentialId;
        credentialId = Common.getCredentialId(credentialInfo);
        if (credentialId instanceof Error) {
            rej(credentialId);
        }
        console.info("removing secret: " + credentialId);
        // sectet info needed to delete
        var params = {
            SecretId: credentialId
        };
        // by default wait 7 days to delete a secret 
        // unless explictly decideing to force delete
        if (forceDeleteWithoutRecovery === true) {
            params['ForceDeleteWithoutRecovery'] = true;
        }
        else {
            // recoverable for 7 days before permanent deletion
            params['RecoveryWindowInDays'] = 7;
        }
        // see https://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/SecretsManager.html#deleteSecret-property
        var client = new aws_sdk_1.SecretsManager();
        client.deleteSecret(params, function (err, data) {
            if (err) {
                console.error("Error deleting secret");
                rej(err.message);
            }
            else {
                console.info("removing secret: Success!");
                res(data);
            }
        });
    });
}
exports.removeSecret = removeSecret;
function getSecretString(secretContent) {
    if (secretContent.filePath) {
        return fs.readFileSync(secretContent.filePath, 'utf8');
    }
    else if (secretContent.key && secretContent.value) {
        return "{\"" + secretContent.key + "\":\"" + secretContent.value + "\"}";
    }
    else {
        return new Error("Error parsing secret content, please specify the correct key/value or file path");
    }
}
//# sourceMappingURL=awsSecret.js.map