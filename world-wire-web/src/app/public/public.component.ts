// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit } from '@angular/core';
import { Modal } from 'carbon-components';
import { AuthService } from '../shared/services/auth.service';

@Component({
  templateUrl: './public.component.html',
  styleUrls: ['./public.component.scss']
})
export class PublicComponent implements OnInit {

  modal: any;

  constructor(
    public authService: AuthService
  ) { }

  ngOnInit() {
    const modalElement = document.getElementById('modal-side-nav');
    this.modal = Modal.create(modalElement);
  }

  openModal() {
    this.modal.show();
  }

}
