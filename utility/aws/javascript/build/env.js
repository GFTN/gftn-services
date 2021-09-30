// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
"use strict";
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
var SM = require("./awsSecret");
var Var = require("./utility/var");
function CheckVariable() {
    if (!process.env.HOME_DOMAIN_NAME || !process.env.SERVICE_NAME || !process.env.ENV_VERSION || !process.env.SECRET_STORAGE_LOCATION) {
        throw new Error(("Initializing failed, Require the following environment variables to start up the service. HOME_DOMAIN_NAME, SERVICE_NAME, ENV_VERSION, SECRET_STORAGE_LOCATION"));
    }
}
exports.CheckVariable = CheckVariable;
function GetEnv(credential) {
    var _this = this;
    return new Promise(function (resolve, rej) { return __awaiter(_this, void 0, void 0, function () {
        var res, e_1, obj, keys, i;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    if (!(process.env[credential.variable] && credential.variable !== "initialize")) return [3 /*break*/, 1];
                    console.log("env already been initialized");
                    return [3 /*break*/, 6];
                case 1:
                    res = void 0;
                    _a.label = 2;
                case 2:
                    _a.trys.push([2, 4, , 5]);
                    return [4 /*yield*/, SM.getSecret(credential)];
                case 3:
                    res = _a.sent();
                    return [3 /*break*/, 5];
                case 4:
                    e_1 = _a.sent();
                    rej(e_1);
                    return [2 /*return*/];
                case 5:
                    obj = JSON.parse(res);
                    keys = Object.keys(obj);
                    for (i = 0; i < keys.length; i++) {
                        process.env[keys[i]] = obj[keys[i]];
                    }
                    process.env[credential.variable] = "true";
                    _a.label = 6;
                case 6:
                    resolve("success");
                    return [2 /*return*/];
            }
        });
    }); });
}
function InitEnv() {
    var _this = this;
    return new Promise(function (res, rej) { return __awaiter(_this, void 0, void 0, function () {
        var domainId, svcName, envVersion, participant_credential, e_2, service_credential, e_3;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    if (!(process.env["SECRET_STORAGE_LOCATION"] === Var.AWS_SECRET)) return [3 /*break*/, 8];
                    console.log("Initializing service with AWS");
                    //participant
                    console.log("Initializing participant-specific secrets from AWS");
                    domainId = process.env.HOME_DOMAIN_NAME, svcName = process.env.SERVICE_NAME, envVersion = process.env.ENV_VERSION;
                    participant_credential = {
                        environment: envVersion,
                        domain: domainId,
                        service: "participant",
                        variable: "initialize"
                    };
                    _a.label = 1;
                case 1:
                    _a.trys.push([1, 3, , 4]);
                    return [4 /*yield*/, GetEnv(participant_credential)];
                case 2:
                    _a.sent();
                    return [3 /*break*/, 4];
                case 3:
                    e_2 = _a.sent();
                    console.log(e_2);
                    rej(e_2);
                    return [3 /*break*/, 4];
                case 4:
                    //service
                    console.log("Initializing service-specific secrets from AWS");
                    service_credential = {
                        environment: envVersion,
                        domain: domainId,
                        service: svcName,
                        variable: "initialize"
                    };
                    _a.label = 5;
                case 5:
                    _a.trys.push([5, 7, , 8]);
                    return [4 /*yield*/, GetEnv(service_credential)];
                case 6:
                    _a.sent();
                    return [3 /*break*/, 8];
                case 7:
                    e_3 = _a.sent();
                    console.log(e_3);
                    rej(e_3);
                    return [3 /*break*/, 8];
                case 8:
                    res("success");
                    return [2 /*return*/];
            }
        });
    }); });
}
exports.InitEnv = InitEnv;
//# sourceMappingURL=env.js.map