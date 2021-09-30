// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, Input, OnChanges, OnDestroy } from '@angular/core';
import { SwaggerService, INavPathDef } from '../../../shared/services/swagger.service';
import { MatDialog } from '@angular/material/dialog';

@Component({
  selector: 'app-operation',
  templateUrl: './operation.component.html',
  styleUrls: ['./operation.component.scss']
})
export class OperationComponent implements OnChanges, OnDestroy {

  @Input() navPath: INavPathDef;
  responseArr: string[];

  private markdown: string;

  constructor(
    public swaggerService: SwaggerService,
    public dialog: MatDialog
  ) { }

  ngOnChanges() {
    this.responseArr = Object.keys(this.navPath.operation.responses);
    this.markdown = this.swaggerService.toMarkdown(this.navPath.operation.description);
  }

  ngOnDestroy() {
    this.markdown = '';
  }

}
