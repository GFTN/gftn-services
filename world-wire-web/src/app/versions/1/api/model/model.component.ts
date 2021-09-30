// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, Input, OnInit } from '@angular/core';
import * as _ from 'lodash';
import { RegexPipe } from '../../../../shared/pipes/regex.pipe';
import { Schema } from 'swagger-schema-official';
import { SwaggerService } from '../../../shared/services/swagger.service';

@Component({
  selector: 'app-model',
  templateUrl: './model.component.html',
  styleUrls: ['./model.component.scss'],
  providers: [RegexPipe]
})
export class ModelComponent implements OnInit {

  @Input() model: {
    // title
    key: string,
    // model definition
    value: Schema
  };

  constructor(
    public swaggerService: SwaggerService
  ) { }

  ngOnInit() {
    // console.log(this.model);
  }

}
