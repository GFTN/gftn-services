// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { TOTPTwoFactorPrivate } from './totp-two-factor.controller';
// unit testing:
// install jest: $ npm i jest @types/jest ts-jest -D
// to kick off the test: $ npx jest

test('randstr generator len 32', () => {
    const obj = new TOTPTwoFactorPrivate();
    const str = obj.randStr(32, false);
    expect(str.length).toBe(32);
});

test('randstr generator len 64', () => {
    const obj = new TOTPTwoFactorPrivate();
    const str = obj.randStr(64, false);
    expect(str.length).toBe(64);
});

// note: this is just a basic check but not any sort of testing that checks the reliability/randomness of the random number generator
test('randstr generator sanity check for non repetition', () => {
    const obj = new TOTPTwoFactorPrivate();
    const strHashmap:{[id :string]: number} = {};
    for (let i = 0; i < 10; i++){
        const str = obj.randStr(32, false);
        expect(!(str in strHashmap)).toBe(true);
        strHashmap[str] = 1; 
    }
});