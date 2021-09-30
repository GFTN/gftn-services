// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Pipe, PipeTransform } from '@angular/core';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';

@Pipe({
    name: 'keepHtml', pure: false
})

export class EscapeHtmlPipe implements PipeTransform {

    constructor(private sanitizer: DomSanitizer) {
    }

    transform(content): SafeHtml {
        return this.sanitizer.bypassSecurityTrustHtml(content);
    }
}
