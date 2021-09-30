// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Pipe } from '@angular/core';
import { PipeTransform } from '@angular/core/src/change_detection/pipe_transform';

@Pipe({ name: 'phone' })
export class PhonePipe implements PipeTransform {

    transform(tel, format?: 'euro' | 'us') {
        const value = tel.toString().trim().replace(/^\+/, '');

        // console.log(format);

        if (value.match(/[^0-9]/)) {
            return tel;
        }

        let country, city, number;

        switch (value.length) {
            case 10: // +1PPP####### -> C (PPP) ###-####
                country = 1;
                city = value.slice(0, 3);
                number = value.slice(3);
                break;

            case 11: // +CPPP####### -> CCC (PP) ###-####
                country = '+' + value[0];
                city = value.slice(1, 4);
                number = value.slice(4);
                break;

            case 12: // +CCCPP####### -> CCC (PP) ###-####
                country = '+' + value.slice(0, 3);
                city = value.slice(3, 5);
                number = value.slice(5);
                break;

            default:
                return tel;
        }

        if (country === 1) {
            country = '';
        }

        if (format === 'euro') {
            // euro number format
            number = number.slice(0, 3) + '.' + number.slice(3);
            return (country + '.' + city + '.' + number).trim();
        } else {
            // classic phone number format
            number = number.slice(0, 3) + '-' + number.slice(3);
            return (country + ' (' + city + ') ' + number).trim();
        }

    }
}
