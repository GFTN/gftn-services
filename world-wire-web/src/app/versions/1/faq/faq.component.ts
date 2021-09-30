// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { HttpClient } from '@angular/common/http';
import { Component, OnInit, AfterViewInit } from '@angular/core';
import * as _ from 'lodash';
import { VersionService } from '../../shared/services/version.service';

@Component({
  selector: 'app-faq',
  templateUrl: './faq.component.html',
  styleUrls: ['./faq.component.scss']
})
export class FaqComponent implements OnInit {
  questions: any[];
  categories: string[];
  category_questions: Map<string, any[]> = new Map();
  accordionInitialized = false;

  constructor(
    public http: HttpClient,
    private versionService: VersionService
  ) { }

  ngOnInit() {
    this.http.request('GET', './assets/open-api/' + this.versionService.current.module + '/faq.json')
      .toPromise().then((data: any) => {
        this.questions = _.sortBy(data, ['question']);

        // transform all questions into one object each based on category
        for (const question of this.questions) {
          if (!this.category_questions.has(question.category)) {
            this.category_questions.set(question.category, []);
          }

          const questionArray = this.category_questions.get(question.category);
          questionArray.push(question);
        }
      });
  }
}
