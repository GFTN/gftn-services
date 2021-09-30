// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
"use strict";
/*
Parameter Naming Constraints:
    * Parameter names are case sensitive.
    * A parameter name must be unique within an AWS Region
    * A parameter name can't be prefixed with "aws" or "ssm" (case-insensitive).
    * Parameter names can include only the following symbols and letters: a-zA-Z0-9_.-/
    * A parameter name can't include spaces.
    * Parameter hierarchies are limited to a maximum depth of fifteen levels.
*/
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
Object.defineProperty(exports, "__esModule", { value: true });
var Common = require("./utility/common");
var aws_sdk_1 = require("aws-sdk");
function getParameter(credentialInfo) {
    return new Promise(function (res, rej) {
        var credentialId;
        credentialId = Common.getCredentialId(credentialInfo);
        if (credentialId instanceof Error) {
            rej(credentialId);
        }
        console.log("getParameter: " + credentialId);
        var ssm = new aws_sdk_1.SSM();
        var params = {
            Name: credentialId,
            WithDecryption: true
        };
        ssm.getParameter(params, function (err, data) {
            if (err) {
                console.error("Error getting parameter");
                rej(err);
            }
            else {
                console.log(credentialId + " successfully retrieved");
                var resString = JSON.stringify(data);
                var resJson = JSON.parse(resString);
                res(resJson.Parameter.Value);
            }
        });
    });
}
exports.getParameter = getParameter;
function createParameter(credentialInfo, parameterContent) {
    return __awaiter(this, void 0, void 0, function () {
        return __generator(this, function (_a) {
            return [2 /*return*/, putParameter(credentialInfo, parameterContent, false)];
        });
    });
}
exports.createParameter = createParameter;
function updateParameter(credentialInfo, parameterContent) {
    return __awaiter(this, void 0, void 0, function () {
        return __generator(this, function (_a) {
            return [2 /*return*/, putParameter(credentialInfo, parameterContent, true)];
        });
    });
}
exports.updateParameter = updateParameter;
function removeParameter(credentialInfo) {
    return new Promise(function (res, rej) {
        var credentialId;
        credentialId = Common.getCredentialId(credentialInfo);
        if (credentialId instanceof Error) {
            rej(credentialId);
        }
        console.log("removeParameter: " + credentialId);
        var ssm = new aws_sdk_1.SSM();
        var params = {
            Name: credentialId
        };
        ssm.deleteParameter(params, function (err, data) {
            if (err) {
                console.error("Error deleting parameter");
                rej(err);
            }
            else {
                console.log(credentialId + " successfully removed");
                res("success");
            }
        });
    });
}
exports.removeParameter = removeParameter;
function putParameter(credentialInfo, parameterContent, overwrite) {
    return new Promise(function (res, rej) {
        var credentialId;
        credentialId = Common.getCredentialId(credentialInfo);
        if (credentialId instanceof Error) {
            rej(credentialId);
        }
        if (overwrite) {
            console.info("updateParameter: " + credentialId);
        }
        else {
            console.info("createParameter: " + credentialId);
        }
        var ssm = new aws_sdk_1.SSM();
        var params = {
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
            }
            else {
                console.log(credentialId + " successfully added");
                var resString = JSON.stringify(data);
                var resJson = JSON.parse(resString);
                res(resJson);
            }
        });
    });
}
//# sourceMappingURL=awsParameter.js.map