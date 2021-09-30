// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, Inject } from '@angular/core';
import { INodeAutomation, NodeConfigData } from '../../../../shared/models/node.interface';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { IInstitution } from '../../../../shared/models/participant.interface';
import { NodeService } from '../../shared/node.service';

@Component({
  selector: 'app-edit-node',
  templateUrl: './edit-node.component.html',
  styleUrls: ['./edit-node.component.scss']
})
export class EditNodeComponent implements OnInit {

  node: INodeAutomation;
  configData: NodeConfigData;
  institution: IInstitution;

  loaded: boolean;

  subtitle = 'View/Edit';

  constructor(
    public dialogRef: MatDialogRef<EditNodeComponent>,
    private nodeService: NodeService,
    @Inject(MAT_DIALOG_DATA) public data,
  ) {

    this.node = data.node;
    this.subtitle = this.node.update ? 'Review' : this.subtitle;
    this.institution = data.institution;

    this.loaded = false;

  }

  ngOnInit() {

    // get node details
    if (this.node) {

      const nodeStatus = this.nodeService.getStatus(this.node.status);

      // allow retries if configuration failed for some reason or another
      if (nodeStatus.includes('failed') && !this.node.update) {
        this.node.update = this.node;
      }

      // get configData
      this.configData = this.node;

      this.loaded = true;
    }
  }

  closeModal() {
    this.dialogRef.close();
  }
}
