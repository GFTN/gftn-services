// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
/* tslint:disable */
import { Controller, ValidationService, FieldErrors, ValidateError, TsoaRoute } from 'tsoa';
import { TOTPTwoFactorAuthController } from './controllers/totp-two-factor.controller';
import { HelperController } from './controllers/helper.controller';
import { IBMIdController } from './controllers/ibmid.controller';
import { PermissionsController } from './controllers/permissions.controller';
import { InstitutionController } from './controllers/institution.controller';
import { JwtController } from './controllers/jwt.controller';
import { expressAuthentication } from './middleware/authentication';
import * as express from 'express';

const models: TsoaRoute.Models = {
    "TOTPQRcode": {
        "properties": {
            "qrcodeURI": { "dataType": "string", "required": true },
            "accountName": { "dataType": "string", "required": true },
        },
    },
    "TOTPResponse": {
        "properties": {
            "success": { "dataType": "boolean", "required": true },
            "registered": { "dataType": "boolean" },
            "msg": { "dataType": "string" },
            "data": { "ref": "TOTPQRcode" },
        },
    },
    "TokenBody": {
        "properties": {
            "token": { "dataType": "string", "required": true },
        },
    },
    "IEncodeFirebaseCred": {
        "properties": {
            "type": { "dataType": "string", "required": true },
            "project_id": { "dataType": "string", "required": true },
            "private_key_id": { "dataType": "string", "required": true },
            "private_key": { "dataType": "string", "required": true },
            "client_email": { "dataType": "string", "required": true },
            "client_id": { "dataType": "string", "required": true },
            "auth_uri": { "dataType": "string", "required": true },
            "token_uri": { "dataType": "string", "required": true },
            "auth_provider_x509_cert_url": { "dataType": "string", "required": true },
            "client_x509_cert_url": { "dataType": "string", "required": true },
        },
    },
    "IDecodeResult": {
        "properties": {
            "encodedText": { "dataType": "string", "required": true },
            "getJson": { "dataType": "boolean", "required": true },
        },
    },
    "IVerifyCompare": {
        "properties": {
            "endpoint": { "dataType": "string", "required": true },
            "ip": { "dataType": "string", "required": true },
            "account": { "dataType": "string" },
        },
    },
    "IJWTTokenInfoPublic": {
        "properties": {
            "acc": { "dataType": "array", "array": { "dataType": "string" }, "required": true },
            "ver": { "dataType": "string", "required": true },
            "ips": { "dataType": "array", "array": { "dataType": "string" }, "required": true },
            "env": { "dataType": "string", "required": true },
            "enp": { "dataType": "array", "array": { "dataType": "string" }, "required": true },
            "jti": { "dataType": "string", "required": true },
            "aud": { "dataType": "string", "required": true },
            "description": { "dataType": "string", "required": true },
        },
    },
};
const validationService = new ValidationService(models);

export function RegisterRoutes(app: express.Express) {
    app.get('/totp/:accountName',
        function(request: any, response: any, next: any) {
            const args = {
                accountName: { "in": "path", "name": "accountName", "required": true, "dataType": "string" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new TOTPTwoFactorAuthController();


            const promise = controller.createTOTP.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/totp/:accountName/confirm',
        function(request: any, response: any, next: any) {
            const args = {
                accountName: { "in": "path", "name": "accountName", "required": true, "dataType": "string" },
                body: { "in": "body", "name": "body", "required": true, "ref": "TokenBody" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new TOTPTwoFactorAuthController();


            const promise = controller.confirmTOTP.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/totp/check',
        function(request: any, response: any, next: any) {
            const args = {
                req: { "in": "request", "name": "req", "required": true, "dataType": "object" },
                body: { "in": "body", "name": "body", "required": true, "ref": "TokenBody" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new TOTPTwoFactorAuthController();


            const promise = controller.checkTOTP.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/util/base64/encode-fb-cred',
        function(request: any, response: any, next: any) {
            const args = {
                body: { "in": "body", "name": "body", "required": true, "ref": "IEncodeFirebaseCred" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new HelperController();


            const promise = controller.encodeFbCred.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/util/base64/decode-fb-cred',
        function(request: any, response: any, next: any) {
            const args = {
                body: { "in": "body", "name": "body", "required": true, "ref": "IDecodeResult" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new HelperController();


            const promise = controller.decodeFbCred.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/util/base64/encode',
        function(request: any, response: any, next: any) {
            const args = {
                body: { "in": "body", "name": "body", "required": true, "dataType": "any" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new HelperController();


            const promise = controller.encodeBase64.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/util/base64/decode',
        function(request: any, response: any, next: any) {
            const args = {
                body: { "in": "body", "name": "body", "required": true, "ref": "IDecodeResult" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new HelperController();


            const promise = controller.decodeBase64.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/util/json/escape',
        function(request: any, response: any, next: any) {
            const args = {
                body: { "in": "body", "name": "body", "required": true, "dataType": "any" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new HelperController();


            const promise = controller.escapeJson.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.get('/sso/login',
        function(request: any, response: any, next: any) {
            const args = {
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new IBMIdController();


            const promise = controller.login.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.get('/sso/token',
        function(request: any, response: any, next: any) {
            const args = {
                req: { "in": "request", "name": "req", "required": true, "dataType": "object" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new IBMIdController();


            const promise = controller.token.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.get('/sso/portal-login',
        function(request: any, response: any, next: any) {
            const args = {
                req: { "in": "request", "name": "req", "required": true, "dataType": "object" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new IBMIdController();


            const promise = controller.portalToken.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/sso/portal-login-totp',
        function(request: any, response: any, next: any) {
            const args = {
                req: { "in": "request", "name": "req", "required": true, "dataType": "object" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new IBMIdController();


            const promise = controller.portalTokenTOTP.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.get('/sso/callback',
        function(request: any, response: any, next: any) {
            const args = {
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new IBMIdController();


            const promise = controller.callback.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/sso/logout',
        function(request: any, response: any, next: any) {
            const args = {
                req: { "in": "request", "name": "req", "required": true, "dataType": "object" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new IBMIdController();


            const promise = controller.logout.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/sso/failure',
        function(request: any, response: any, next: any) {
            const args = {
                req: { "in": "request", "name": "req", "required": true, "dataType": "object" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new IBMIdController();


            const promise = controller.failure.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/permissions/participant',
        function(request: any, response: any, next: any) {
            const args = {
                req: { "in": "request", "name": "req", "required": true, "dataType": "object" },
                body: { "in": "body", "name": "body", "required": true, "dataType": "any" },
                fid: { "in": "header", "name": "x-fid", "required": true, "dataType": "string" },
                iid: { "in": "header", "name": "x-iid", "required": true, "dataType": "string" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new PermissionsController();


            const promise = controller.updateParticipantUserPermissions.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/permissions/super',
        function(request: any, response: any, next: any) {
            const args = {
                req: { "in": "request", "name": "req", "required": true, "dataType": "object" },
                body: { "in": "body", "name": "body", "required": true, "dataType": "any" },
                fid: { "in": "header", "name": "x-fid", "required": true, "dataType": "string" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new PermissionsController();


            const promise = controller.updateSuperPermissions.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/institution/create',
        authenticateMiddleware([{ "api_key": [] }]),
        function(request: any, response: any, next: any) {
            const args = {
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new InstitutionController();


            const promise = controller.create.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/jwt/refresh',
        function(request: any, response: any, next: any) {
            const args = {
                req: { "in": "request", "name": "req", "required": true, "dataType": "object" },
                authorization: { "in": "header", "name": "Authorization", "required": true, "dataType": "string" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new JwtController();


            const promise = controller.refresh.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/jwt/verify',
        function(request: any, response: any, next: any) {
            const args = {
                req: { "in": "request", "name": "req", "required": true, "dataType": "object" },
                body: { "in": "body", "name": "body", "required": true, "ref": "IVerifyCompare" },
                authorization: { "in": "header", "name": "Authorization", "required": true, "dataType": "string" },
                fid: { "in": "header", "name": "x-fid", "required": true, "dataType": "string" },
                iid: { "in": "header", "name": "x-iid", "required": true, "dataType": "string" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new JwtController();


            const promise = controller.verify.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/jwt/request',
        function(request: any, response: any, next: any) {
            const args = {
                req: { "in": "request", "name": "req", "required": true, "dataType": "object" },
                body: { "in": "body", "name": "body", "required": true, "ref": "IJWTTokenInfoPublic" },
                fid: { "in": "header", "name": "x-fid", "required": true, "dataType": "string" },
                iid: { "in": "header", "name": "x-iid", "required": true, "dataType": "string" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new JwtController();


            const promise = controller.request.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/jwt/approve',
        function(request: any, response: any, next: any) {
            const args = {
                req: { "in": "request", "name": "req", "required": true, "dataType": "object" },
                body: { "in": "body", "name": "body", "required": true, "dataType": "any" },
                fid: { "in": "header", "name": "x-fid", "required": true, "dataType": "string" },
                iid: { "in": "header", "name": "x-iid", "required": true, "dataType": "string" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new JwtController();


            const promise = controller.approve.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/jwt/reject',
        function(request: any, response: any, next: any) {
            const args = {
                req: { "in": "request", "name": "req", "required": true, "dataType": "object" },
                body: { "in": "body", "name": "body", "required": true, "dataType": "any" },
                fid: { "in": "header", "name": "x-fid", "required": true, "dataType": "string" },
                iid: { "in": "header", "name": "x-iid", "required": true, "dataType": "string" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new JwtController();


            const promise = controller.reject.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/jwt/revoke',
        function(request: any, response: any, next: any) {
            const args = {
                req: { "in": "request", "name": "req", "required": true, "dataType": "object" },
                body: { "in": "body", "name": "body", "required": true, "dataType": "any" },
                fid: { "in": "header", "name": "x-fid", "required": true, "dataType": "string" },
                iid: { "in": "header", "name": "x-iid", "required": true, "dataType": "string" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new JwtController();


            const promise = controller.revoke.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.post('/jwt/generate',
        function(request: any, response: any, next: any) {
            const args = {
                req: { "in": "request", "name": "req", "required": true, "dataType": "object" },
                body: { "in": "body", "name": "body", "required": true, "dataType": "any" },
                fid: { "in": "header", "name": "x-fid", "required": true, "dataType": "string" },
                iid: { "in": "header", "name": "x-iid", "required": true, "dataType": "string" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new JwtController();


            const promise = controller.generate.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });
    app.get('/jwt/rotate-pepper',
        function(request: any, response: any, next: any) {
            const args = {
                fid: { "in": "header", "name": "x-fid", "required": true, "dataType": "string" },
            };

            let validatedArgs: any[] = [];
            try {
                validatedArgs = getValidatedArgs(args, request);
            } catch (err) {
                return next(err);
            }

            const controller = new JwtController();


            const promise = controller.rotatePepper.apply(controller, validatedArgs as any);
            promiseHandler(controller, promise, response, next);
        });

    function authenticateMiddleware(security: TsoaRoute.Security[] = []) {
        return (request: any, _response: any, next: any) => {
            let responded = 0;
            let success = false;

            const succeed = function(user: any) {
                if (!success) {
                    success = true;
                    responded++;
                    request['user'] = user;
                    next();
                }
            }

            const fail = function(error: any) {
                responded++;
                if (responded == security.length && !success) {
                    error.status = 401;
                    next(error)
                }
            }

            for (const secMethod of security) {
                if (Object.keys(secMethod).length > 1) {
                    let promises: Promise<any>[] = [];

                    for (const name in secMethod) {
                        promises.push(expressAuthentication(request, name, secMethod[name]));
                    }

                    Promise.all(promises)
                        .then((users) => { succeed(users[0]); })
                        .catch(fail);
                } else {
                    for (const name in secMethod) {
                        expressAuthentication(request, name, secMethod[name])
                            .then(succeed)
                            .catch(fail);
                    }
                }
            }
        }
    }

    function isController(object: any): object is Controller {
        return 'getHeaders' in object && 'getStatus' in object && 'setStatus' in object;
    }

    function promiseHandler(controllerObj: any, promise: any, response: any, next: any) {
        return Promise.resolve(promise)
            .then((data: any) => {
                let statusCode;
                if (isController(controllerObj)) {
                    const headers = controllerObj.getHeaders();
                    Object.keys(headers).forEach((name: string) => {
                        response.set(name, headers[name]);
                    });

                    statusCode = controllerObj.getStatus();
                }

                if (data || data === false) { // === false allows boolean result
                    response.status(statusCode || 200).json(data);
                } else {
                    response.status(statusCode || 204).end();
                }
            })
            .catch((error: any) => next(error));
    }

    function getValidatedArgs(args: any, request: any): any[] {
        const fieldErrors: FieldErrors = {};
        const values = Object.keys(args).map((key) => {
            const name = args[key].name;
            switch (args[key].in) {
                case 'request':
                    return request;
                case 'query':
                    return validationService.ValidateParam(args[key], request.query[name], name, fieldErrors);
                case 'path':
                    return validationService.ValidateParam(args[key], request.params[name], name, fieldErrors);
                case 'header':
                    return validationService.ValidateParam(args[key], request.header(name), name, fieldErrors);
                case 'body':
                    return validationService.ValidateParam(args[key], request.body, name, fieldErrors, name + '.');
                case 'body-prop':
                    return validationService.ValidateParam(args[key], request.body[name], name, fieldErrors, 'body.');
            }
        });
        if (Object.keys(fieldErrors).length > 0) {
            throw new ValidateError(fieldErrors, '');
        }
        return values;
    }
}
