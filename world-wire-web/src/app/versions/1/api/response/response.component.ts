// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Response } from 'swagger-schema-official';
import { Component, OnInit, Input } from '@angular/core';
import { SwaggerService } from '../../../shared/services/swagger.service';

@Component({
  selector: '[app-response]',
  templateUrl: './response.component.html',
  styleUrls: ['./response.component.scss']
})
export class ResponseComponent implements OnInit {

  @Input() responseCode: string;
  @Input() responseDef: Response;

  constructor(public swaggerService: SwaggerService) { }

  ngOnInit() {
    // console.log({ code: this.responseCode, def: this.responseDef });
  }

}
