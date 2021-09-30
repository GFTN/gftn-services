// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
/**
 * TOTP profile for firebase
 *
 * @export
 * @interface TOTPProfile
 */
export interface TOTPProfile {
    key: string;
    registered: boolean;
}
/**
 * TOTP QR code object containing QR code URL containing the TOPT seed and user profile.
 *
 * @export
 * @interface TOTPQRcode
 */
export interface TOTPQRcode {
    // The base64 encoded representation the the QR Code.
    qrcodeURI: string;
    // A user friendly name for the registration
    accountName: string;
}
/**
 * TOTP response for TOTP API endpoints
 *
 * @export
 * @interface TOTPResponse
 */
export interface TOTPResponse {
    // response status for the api call; true for success
    success: boolean;
    // status of user's TOTP registration; true for registered
    registered?: boolean;
    msg?: string;
    data?: TOTPQRcode;
}
/**
 * TOTP token as body for TOTP API endpoints
 *
 * @export
 * @interface TokenBody
 */
export interface TokenBody {
    token: string;
}

export interface TOTPRegistrationData {
    email: string;
}
