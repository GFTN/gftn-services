// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Pipe } from '@angular/core';
import { PipeTransform } from '@angular/core/src/change_detection/pipe_transform';

@Pipe({ name: 'regex' })
export class RegexPipe implements PipeTransform {

    transform(text: string, regex: string, flags: string, replaceText: string) {

        // > this code will remove all special characters
        // https://stackoverflow.com/questions/4374822/remove-all-special-characters-with-regexp
        // return text.replace(/[^\w\s]/gi, '_');

        // > How to use 'new RegEx()'
        // https://stackoverflow.com/questions/874709/converting-user-input-string-to-regular-expression
        // var re = new RegExp("a|b", "i");
        // // same as
        // var re = /a|b/i;

        // > Escaping \ in string for 'new RegEx()'
        // https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/RegExp

        // > Test a regEx expression using this Generator
        // http://scriptular.com/

        return text.replace(new RegExp(regex, flags), replaceText);
    }
}
