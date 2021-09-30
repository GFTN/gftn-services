// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, HostListener } from '@angular/core';

@Component({
  selector: 'app-back-to-top',
  templateUrl: './back-to-top.component.html',
  styleUrls: ['./back-to-top.component.scss']
})
export class BackToTopComponent implements OnInit {

  showBackToTop = false;

  constructor() { }

  ngOnInit() {
  }

  @HostListener('window:scroll', [])
  onScroll(): void {
    if (window.scrollY === 0) {
      // Scroll position is at the top of the page
      this.showBackToTop = false;
    } else {
      // Scroll position is not at the top of the page
      this.showBackToTop = true;
    }
  }

  scrollToTop() {

    // animates scroll to top using window.setInterval
    const scrollToTop = window.setInterval(function () {
      const pos = window.pageYOffset;
      if (pos > 0) {
        window.scrollTo(0, pos - 10); // how far to scroll on each step
      } else {
        window.clearInterval(scrollToTop);
      }
    }, 15);
  }
}
