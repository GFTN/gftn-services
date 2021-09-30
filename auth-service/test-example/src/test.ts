// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { hello, another } from './hello';
import { expect } from 'chai';
import 'mocha';

describe('Hello function', () => {
  
  it('should return hello world', () => {
    const result = hello();
    expect(result).to.equal('Hello World!');
  });

  it('should also return hello world', () => {
    const result = new another().hello();
    expect(result).to.equal('Hello World!');
  });

});