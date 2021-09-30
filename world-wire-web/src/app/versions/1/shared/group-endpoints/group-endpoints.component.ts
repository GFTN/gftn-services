// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, AfterViewInit, ViewChild, ElementRef, Input } from '@angular/core';
import { Accordion } from 'carbon-components';
import { SwaggerService, INavPathDef } from '../../../shared/services/swagger.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-group-endpoints',
  templateUrl: './group-endpoints.component.html',
  styleUrls: ['./group-endpoints.component.scss']
})
export class GroupEndpointsComponent implements OnInit, AfterViewInit {

  @Input() group: INavPathDef;
  @ViewChild('accordion') accordion: ElementRef;

  private accordionInitialized = false;

  constructor(
    public swaggerService: SwaggerService,
    public router: Router
  ) { }

  ngOnInit() {
    // console.log(this.group);
  }

  ngAfterViewInit() {
    // initialize accordion only once
    if (this.accordionInitialized === false) {
      Accordion.create(this.accordion.nativeElement);
      this.accordionInitialized = true;
    }
  }

}
