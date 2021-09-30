// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { EscapeHtmlPipe } from './pipes/keep-html.pipe';

@NgModule({
  declarations: [
    EscapeHtmlPipe
  ],
  imports: [
    CommonModule
  ],
  exports: [
    EscapeHtmlPipe
  ]
})
export class PipeModule { }
