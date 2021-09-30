// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Component, OnInit, HostBinding, OnDestroy, NgZone, isDevMode } from '@angular/core';
import { Router } from '@angular/router';
import { INodeAutomation } from '../../../shared/models/node.interface';
import { AngularFireDatabase } from '@angular/fire/database';
import { MatDialog } from '@angular/material/dialog';
import { EditNodeComponent } from './edit-node/edit-node.component';
import { Observable, Observer } from 'rxjs';
import { filter } from 'lodash';
import { SessionService } from '../../../shared/services/session.service';
import { ModalService } from 'carbon-components-angular';
import { DeleteNodeModalComponent } from './delete-node-modal/delete-node-modal.component';
import { RegistrationModalComponent } from './registration-modal/registration-modal.component';
import { NodeService } from '../shared/node.service';

@Component({
  selector: 'app-nodes',
  templateUrl: './nodes.component.html',
  styleUrls: ['./nodes.component.scss']
})
export class NodesComponent implements OnInit, OnDestroy {

  $nodes: Observable<INodeAutomation[]>;
  // institutionId: string;
  // name: string;
  loaded: boolean;

  dbRef: firebase.database.Reference;

  showDeleted = false;

  constructor(
    public dialog: MatDialog,
    private router: Router,
    public sessionService: SessionService,
    private modalService: ModalService,
    public nodeService: NodeService,
    private db: AngularFireDatabase,
    private ngZone: NgZone
  ) {
    this.loaded = false;
  }

  @HostBinding('attr.class') cls = 'flex-fill';

  ngOnInit() {

    // get participant id from currently signed in user
    // this.institutionId = this.sessionService.institution.institutionId;
    // this.name = this.sessionService.institution.name;


    // set reference to firebase db for subscribing/unsubscribing to
    this.dbRef = this.db.database.ref(`/participants/${this.sessionService.institution.info.institutionId}/nodes`);

    this.$nodes = this.getNodes();
  }

  ngOnDestroy() {
    // unsubscribe from firebase db ref
    this.dbRef.off();
  }

  /**
   * Opens up modal for account registration
   * (currently Anchor/issuer ONLY)
   *
   * @param {INodeAutomation} node
   * @memberof NodesComponent
   */
  public openRegistrationDialog(node: INodeAutomation) {
    // creates and opens the modal for deleting node
    this.modalService.create({
      component: RegistrationModalComponent,
      inputs: {
        MODAL_DATA: {
          node: node
        }
      }
    });
  }

  /**
   * Opens dialog for editing a node
   * @param node
   */
  public openEditDialog(node: INodeAutomation) {

    this.dialog.open(EditNodeComponent, {
      data: {
        node: node,
        institution: this.sessionService.institution
      }
    });
  }

  /**
   * Opens confirmation dialog for deleting a node
   * @param node
   */
  public openDeleteDialog(node: INodeAutomation) {

    // creates and opens the modal for deleting node
    this.modalService.create({
      component: DeleteNodeModalComponent,
      inputs: {
        MODAL_DATA: {
          node: node
        }
      }
    });
  }

  public getNodes(): Observable<INodeAutomation[]> {

    let init = true;

    const source = new Observable((observer: Observer<INodeAutomation[]>) => {
      this.dbRef.on('value', (data: any) => {

        // convert data.val() from object to array of nodes
        let nodes: INodeAutomation[] = data.val() ? Object.values(data.val()) : [];

        // map data if nodes exist
        if (nodes) {

          // exclude deleted nodes. these remain in firebase to keep track
          // of nodes that have already been created for this institution/participant
          // so another duplicate node cannot get created with the same participantId
          if (!this.showDeleted) {
            nodes = filter(nodes, (node: INodeAutomation) => {
              return node.status[0] !== 'deleted';
            });
          }
        }

        // since calling from external source
        // need to put result into angular zone
        this.ngZone.run(() => {

          this.loaded = true;

          // redirect if no nodes are returned
          if (nodes.length === 0 && init) {
            init = false;
            // 0.5 sec timeout before immediate redirect
            setTimeout(() => {
              this.router.navigate([`/office/account/${this.sessionService.institution.info.slug}/nodes/add`]);
            }, 500);
          }

          // update observer value
          observer.next(
            nodes
          );

        });

      }, (err) => {
        console.log(err);
      });
    });

    return source;
  }
  /**
   * Toggle showing of deleted nodes
   */
  toggleShowDeleted() {
    this.showDeleted = !this.showDeleted;

    this.$nodes = this.getNodes();
  }
}
