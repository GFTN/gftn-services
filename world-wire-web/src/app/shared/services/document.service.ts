// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable } from '@angular/core';

@Injectable()
export class DocumentService {

    // constructor(){
    //     console.log('init window service');
    // }

    get documentRef(): Document {
        return document;
    }

}
