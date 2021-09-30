// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
export const CUSTOM_REGEXES: RegexMap = {
    // general regex for all universal text. Should prevent basic unknown Unicode/non-unicode characters from being entered
    text: {
        pattern: new RegExp(/^[0-9A-Za-z!@#$%&*()_\-+={[}\]|\:;"'<,>.?\/\\~`]+[0-9A-Za-z!@#$%&*()_\-+={[}\]|\:;"'<,>.?\/\\~`\s]*$/),
        validationText: 'Invalid characters were entered.'
    },
    bic: {
        pattern: '^[A-Z]{3}[A-Z]{3}[A-Z2-9]{1}[A-NP-Z0-9]{1}[A-Z0-9]{3}$',
        validationText: `BIC Code must follow the format of CCCXXXXX000, where:
        • CCC - ISO country code (A-Z): 3 letter code. E.g, for Singapore SGP
        • XXXXXX - first 5 characters of the participant's name. E.g, for MatchMove, MATCH
        • 000 - a unique representation number in WW (0-9): 3-letter code`,
    },

    assetDO: {
        pattern: '^([a-zA-Z]){3}DO$',
        validationText: `Asset Code for a Digital Obligation must end with "DO"`,
    },

    assetDA: {
        pattern: '^([a-zA-Z]){3}$',
        validationText: `Asset Code must be exactly 3 alphabetic characters.`
    },
    ipV4: {
        pattern: '(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)',
        validationText: `Invalid IPv4 address entered.`
    },
    email: {
        pattern: '^([\\w-]+(?:\\.[\\w-]+)*)@((?:[\\w-]+\\.)*\\w[\\w-]{0,66})\\.([a-z]{2,6}(?:\\.[a-z]{2})?)$',
        validationText: `Provided email address must be valid.`,
    },
};

export interface RegexMap {
    [key: string]: CustomRegex;
}

export interface CustomRegex {
    pattern: string | RegExp;
    validationText: string;
}
