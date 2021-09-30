// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// Purpose: these endpoints are common utilities that can be 
// exposed to do certain common tasks. For example, certain 
// environment variables need to be converted to base64 when 
// importing into kubernetes. Use of these endpoints ensures 
// consistency when encoding and decoding. 

import { Route, Controller, Post, Body, OperationId
    // , Security 
} from 'tsoa';
import { IDecodeResult, IEncodeFirebaseCred } from '../models/auth.model';

/**
 * Utility endpoints to assist in development
 *
 * @export
 * @class HelperController
 * @extends {Controller}
 */
@Route('util')
export class HelperController extends Controller {

    /**
     * encode firebase credentials to sanitized base64 string
     *
     * @param {{ text: string; }} body
     * @returns
     * @memberof HelperController
     */
    @Post('base64/encode-fb-cred')
    @OperationId('base64EncodeFbCred')
    // @Security('api_key')
    public async encodeFbCred(
        @Body() body: IEncodeFirebaseCred
    ) {

        // set to escaped json string
        const sanitizedText = JSON.stringify({
            type: body.type,
            project_id: body.project_id,
            private_key_id: body.private_key_id,
            private_key: body.private_key,
            client_email: body.client_email,
            client_id: body.client_id,
            auth_uri: body.auth_uri,
            token_uri: body.token_uri,
            auth_provider_x509_cert_url: body.auth_provider_x509_cert_url,
            client_x509_cert_url: body.client_x509_cert_url,
        });

        // encode to base 64
        const encodedStr = Buffer.from(sanitizedText).toString('base64');
        return encodedStr;
    }

    /**
     * alias for base64/decode-text - can be used to decode firebase credential
     *
     * @param {IDecodeResult} body
     * @returns
     * @memberof HelperController
     */
    @Post('base64/decode-fb-cred')
    @OperationId('base64DecodeFbCred')
    // @Security('api_key')
    public async decodeFbCred(
        @Body() body: IDecodeResult
    ) {
        return this.decodeBase64(body);
    }


    /**
     * encode to base64 string
     * obj = json object in body, or text = someLongString
     *
     * @param {{ obj?: {}; text?: string  }} body
     * @returns
     * @memberof HelperController
     */
    @Post('base64/encode')
    @OperationId('base64Encode')
    // @Security('api_key')
    public async encodeBase64(
        @Body() body: { obj?: {}; text?: string }
    ) {

        // encode to base 64
        if (body.text) {
            const encodedStr = Buffer.from(body.text).toString('base64');
            return encodedStr;
        } else {
            const encodedStr = Buffer.from(JSON.stringify(body.obj)).toString('base64');
            return encodedStr;
        }

    }

    /**
     * decode to base64 string or obj
     *
     * @param {IDecodeResult} body
     * @returns
     * @memberof HelperController
     */
    @Post('base64/decode')
    @OperationId('base64Decode')
    // @Security('api_key')
    public async decodeBase64(
        @Body() body: IDecodeResult
    ) {
        const decodedStr = Buffer.from(body.encodedText, 'base64').toString('utf8');
        if (body.getJson) {
            return JSON.parse(decodedStr);
        } else {
            return decodedStr;
        }
    }

    /**
     * Creates a JSON escaped string 
     *
     * @param {{}} body a json object
     * @memberof HelperController
     */
    @Post('json/escape')
    @OperationId('escapeJson')
    // @Security('api_key')
    public async escapeJson(
        @Body() body: {}
    ) {

        const jsonStr = JSON.stringify(body);
        return jsonStr.replace(/\\n/g, "\\n")
            .replace(/\\'/g, "\\'")
            .replace(/\\"/g, '\\"')
            .replace(/\\&/g, "\\&")
            .replace(/\\r/g, "\\r")
            .replace(/\\t/g, "\\t")
            .replace(/\\b/g, "\\b")
            .replace(/\\f/g, "\\f");

    }

}
