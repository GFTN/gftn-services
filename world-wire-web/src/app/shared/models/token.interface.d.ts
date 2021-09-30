// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
//
export interface IJWTTokens {

    // requires "auth-service to write"
    // should NEVER be exposed via firebase client side sdk
    // jwt info is used to track jwt information
    // status of a request
    jwt_info: {
        [participantId: string]: IJWTPublic
    }

    // jwt_secure requires "triggers to  write" to
    // should NEVER be exposed via firebase client side sdk
    // and is used by the worldwire application to determine
    // if a jwt token is active and not revoked
    jwt_secure: {

        // jti claim (firebase assigned id)
        [jwtId: string]: IWWJWTSecure

    }

}

export interface IDeveloperJWTSecure {

    // NOTE: keyid will be used to find the decryptions secret
    // in the database and will equate the the firebase push
    // id for new record

    // s = the decryptions secret
    // usage: find this secret in firebase db by keyid to decrypt the token
    s: string;
    // i = the jti
    // usage: match the decrypted tokens jti to the stored jti
    i: string;
}

export interface IWWJWTSecure extends IDeveloperJWTSecure {

    // high watermark, a number that is incremented by 1 every
    // time the session is refreshed, (count of sessions) that
    // prevents replay and makes it more difficult to guess
    // the encrypted fields
    n: number;

    // // last 8 characters of the current token,
    // // Trying to produce a token that passes would require the the output encrypted
    // // text to match the last 8 characters of this string.
    // // used as extra validation, this way if someone tries
    // // to reproduce the token by guessing the passphrase
    // // they also have to comprise a combination of fields that when encrypted
    // // the token string matches **i** thereby making it the difficulty level
    // // hard to forge a token. **w** would be not as difficult to guess and check.
    // i: string;
}

export interface IJWTPublic extends IJWTTokenInfoGeneratedPublic, IJWTTokenClaimsAndPayloadPublic { }

export interface IJWTTokenInfoGeneratedPublic {
    stage:
    // created a request for approval
    'review' |
    // 2nd person to approve
    'approved' |
    // ready for a system to make it's first request
    // this is when the timer will start
    'ready' |
    // initialized and is in use 'This is when the token becomes active'
    'initialized' |
    // token was not refreshed in time and was allowed to expire
    // the only thing the user can do at this point is create a
    // new token to restore the session
    'revoked'

    // if the token is currently in use
    active: boolean;

    // who created the token and when
    createdAt: number; // ie: moment.utc().valueOf()
    createdBy: string;
    approvedAt: number;
    approvedBy: string;

    // used have log on who revoked tokens
    // can be used in the view to toggle revoked tokens
    revokedAt?: number;
    revokedBy?: string;

    // when it was last refreshed
    refreshedAt: number; // ie: moment.utc().valueOf()

}

export interface IJWTTokenInfoPublic extends IJWTTokenClaimsAndPayloadPublic {
    description: string;
}


/**
 * JWT Token Claim Values that need to be set
 * (meaning not exposed to the front-end and will
 *  be automatically set by back-end)
 *
 * @export
 * @interface IJWTTokenClaims
 */
export interface IJWTTokenClaimsAndPayloadPublic extends IJWTTokenPayloadPublic, IJWTTokenClaimsPublic { }


/**
 * JWT Token Claims and Payload that includes fields set by the back-end
 *
 * @export
 * @interface IJWTTokenClaimsAndPayloadSecure
 * @extends {IJWTTokenPayloadSecure}
 * @extends {IJWTTokenClaimsSecure}
 */
export interface IJWTTokenClaimsAndPayloadSecure extends IJWTTokenPayloadSecure, IJWTTokenClaimsSecure { }

export interface IJWTTokenClaimsPublic {

    // decrypted jwt token id
    // (firebase id to prevent collision)
    jti: string;

    // aud claim will be used to store **Participant's ID in the Participant Registry**
    // NOTE 1: Participant Id per PR (different than firebase participant id)
    // NOTE 2: was originally using AUD for endpoints but this
    // does not allow for string[] of endpoints
    aud: string; // JSON.stringify

}

/**
 * These are secure (meaning not exposed to the front-end)
 *
 * @export
 * @interface IJWTTokenClaimsSecure
 */
export interface IJWTTokenClaimsSecure extends IJWTTokenClaimsPublic {

    // to be set by back-end
    // issued at
    iat?: number;

    // expired time
    exp?: number;
    nbf?: number;

    // sub will be used to store the institution id per firebase user
    // Note:  **AUD** will be used to store the participant Id per the PR
    // this will be set via the header in the post request "x-iid"
    sub: string;

    // claim included in token header to lookup secret
    kid?: string;

    // not being used
    // iss?: string;
}

/**
 * JWT Token Payload that needs to be set
 * characters kept to a minimum for faster decryption
 *
 * @export
 * @interface IJWTTokenPayload
 */
export interface IJWTTokenPayloadPublic {

    // allowable 'human readable' stellar accounts
    acc: string[]; // JSON.stringify

    // versions
    ver: string;

    // allowable ips
    ips: string[];

    // environment this token can be used for.
    // generally pulled from environment variable
    env: string;

    // endpoints: use enp in middleware to check endpoints
    // instead aud claim so that access info is consistent
    // with where other access rights info is stored in payload,
    // AUD would require JSON.stringify to convert to a string
    // since audience is defined in spec as a string and NOT string[]
    enp: string[];

}

/**
 * This is secure (meaning not exposed to the front-end in jwt_info)
 *
 * @export
 * @interface IJWTTokenPayloadSecure
 * @extends {IJWTTokenPayloadPublic}
 */
export interface IJWTTokenPayloadSecure extends IJWTTokenPayloadPublic {
    // increment count watermark
    // increments by 1 each time the token is refreshed
    // if 0 then the token has never been initialized
    // maximum expiration time on a un-initialize token should be
    // less than 24hrs for increased security purposes, then the token should
    // be refreshed every 15 minutes thereafter
    n: number
}

// this.options = {
//     algorithm: 'RS256',
//     // The "kid" (key ID) Header Parameter is a hint indicating which key
//     // was used to secure the JWS.  This parameter allows originators to
//     // explicitly signal a change of key to recipients.  The structure of
//     // the "kid" value is unspecified.  Its value MUST be a case-sensitive
//     // string.  Use of this Header Parameter is OPTIONAL.
//     // When used with a JWK, the "kid" value is used to match a JWK "kid"
//     // parameter value.
//     // JWK description here: https://auth0.com/docs/jwks
//     // Since WorldWire is both the **issuer and consumer** of a JWT token
//     // symmetrical encryption should suffice for our purposes and as such
//     // we will not need to use a key id in the header.
//     // see additional discussion at https://github.com/GFTN/gftn-services/issues/9
//     // The only useful scenario for using keyId in the header is in the
//     // event we use different secret keys for each participant node
//     // in which case we would need to lookup by "keyid" the appropriate
//     // secret key to use (as stored in a secure database, HSM and/or vault)
//     // keyid: '',
//     // The "jti" (JWT ID) claim provides a unique identifier for the JWT.
//     // The identifier value MUST be assigned in a manner that ensures that
//     // there is a negligible probability that the same value will be
//     // accidentally assigned to a different data object; if the application
//     // uses multiple issuers, collisions MUST be prevented among values
//     // produced by different issuers as well.  The "jti" claim can be used
//     // to prevent the JWT from being replayed.  The "jti" value is a case-
//     // sensitive string.  Use of this claim is OPTIONAL.
//     // Rather than hashing which presents a risk of possible collision,
//     // the firebase id assigned to the new record should suffice for
//     // quick lookup of the validity of the token as it was originally issued
//     // see discussion at https://github.com/GFTN/gftn-services/issues/78
//     jwtid: this.id,
//     // The "exp" (expiration time) claim identifies the expiration time on
//     // or after which the JWT MUST NOT be accepted for processing.  The
//     // processing of the "exp" claim requires that the current date/time
//     // MUST be before the expiration date/time listed in the "exp" claim.
//     // expiresIn: '15m',Implementers MAY provide for some small leeway, usually no more than
//     // a few minutes, to account for clock skew.  Its value MUST be a number
//     // containing a NumericDate value.  Use of this claim is OPTIONAL.
//     // Eg: 60, "2 days", "10h", "7d". A numeric value is interpreted as a seconds count.
//     // If you use a string be sure you provide the time units (days, hours, etc),
//     // otherwise milliseconds unit is used by default ("120" is equal to "120ms").
//     expiresIn: '15m',
//     // The "nbf" (not before) claim identifies the time before which the JWT
//     // MUST NOT be accepted for processing.  The processing of the "nbf"
//     // claim requires that the current date/time MUST be after or equal to
//     // the not-before date/time listed in the "nbf" claim.  Implementers MAY
//     // provide for some small leeway, usually no more than a few minutes, to
//     // account for clock skew.  Its value MUST be a number containing a
//     // NumericDate value.  Use of this claim is OPTIONAL.
//     // Eg: 60, "2 days", "10h", "7d". A numeric value is interpreted as a seconds count.
//     // If you use a string be sure you provide the time units (days, hours, etc),
//     // otherwise milliseconds unit is used by default ("120" is equal to "120ms").
//     // set to 4 seconds (~4000ms) since this should only apply to future blocks
//     // since stellar refreshes ever 5 seconds
//     notBefore: '4000',
//     // The "aud" (audience) claim identifies the recipients that the JWT is
//     // intended for.  Each principal intended to process the JWT MUST
//     // identify itself with a value in the audience claim.  If the principal
//     // processing the claim does not identify itself with a value in the
//     // "aud" claim when this claim is present, then the JWT MUST be
//     // rejected.  In the general case, the "aud" value is an array of case-
//     // sensitive strings, each containing a StringOrURI value.  In the
//     // special case when the JWT has one audience, the "aud" value MAY be a
//     // single case-sensitive string containing a StringOrURI value.  The
//     // interpretation of audience values is generally application specific.
//     // Use of this claim is OPTIONAL.
//     // The audience value is a string -- typically, the base address of the
//     // resource being accessed, such as "https://contoso.com".
//     // ie: allowable world wire endpoints
//     // aka: the endpoints for which this token is to be used
//     // eg: '*' for all endpoints or '/NAME_OF_SPECIFIC_ENDPOINT'
//     audience: ['/test'],
//     // audience: '/test',
//     // The "sub" (subject) claim identifies the principal that is the
//     // subject of the JWT.  The claims in a JWT are normally statements
//     // about the subject.  The subject value MUST either be scoped to be
//     // locally unique in the context of the issuer or be globally unique.
//     // The processing of this claim is generally application specific.  The
//     // "sub" value is a case-sensitive string containing a StringOrURI
//     // value.  Use of this claim is OPTIONAL.
//     // aka: world wire participantId
//     subject: '',
//     // The "iss" (issuer) claim identifies the principal that issued the
//     // JWT.  The processing of this claim is generally application specific.
//     // The "iss" value is a case-sensitive string containing a StringOrURI
//     // value.  Use of this claim is OPTIONAL.
//     issuer: this.env.apiRoot,
//     // noTimestamp: true, // Generated JWTs will include an iat (issued at) claim by default unless noTimestamp is specified.
//     // header: {},
//     // encoding: ''
// };
