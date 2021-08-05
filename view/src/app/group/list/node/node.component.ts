import { Component, OnInit, Input, Output, EventEmitter, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Closed } from 'src/app/core/utils/closed';
import { EditComponent } from 'src/app/group/dialog/edit/edit.component';
import { AddComponent } from 'src/app/group/dialog/add/add.component';
import { FlatNode, Element } from '../../../shared/tree/tree';
import { DeleteComponent } from '../../dialog/delete/delete.component';
import { takeUntil } from 'rxjs/operators';
import { SelectComponent } from '../../dialog/select/select.component';
export interface NodeEvent {
  what: 'add' | 'delete' | 'changed' | 'move'
  node: FlatNode
  data?: Element
}
@Component({
  selector: 'group-list-node',
  templateUrl: './node.component.html',
  styleUrls: ['./node.component.scss']
})
export class NodeComponent implements OnInit, OnDestroy {
  private closed_ = new Closed()
  constructor(private readonly matDialog: MatDialog,
  ) { }
  @Input()
  node: FlatNode | undefined
  @Output() valChange = new EventEmitter<NodeEvent>();
  ngOnInit(): void {
  }
  ngOnDestroy() {
    this.closed_.close()
  }

  onClickAdd(node: FlatNode) {
    this.matDialog.open(AddComponent, {
      data: {
        parent: node.data,
        onAdded: (data: Element) => {
          if (this.closed_.isNotClosed) {
            this.valChange.next({
              what: 'add',
              node: node,
              data: data,
            })
          }
        },
      },
      disableClose: true,
    })
  }
  onClickEdit(node: FlatNode) {
    this.matDialog.open(EditComponent, {
      data: node.data,
      disableClose: true,
    }).afterClosed().pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((ok) => {
      if (typeof ok === "boolean" && ok) {
        this.valChange.next({
          what: 'changed',
          node: node,
        })
      }
    })
  }
  onClickMove(node: FlatNode) {
    this.matDialog.open(SelectComponent, {
      data: node.data,
      disableClose: false,
    }).afterClosed().pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((data: Element) => {
      if (data instanceof Element) {
        this.valChange.next({
          what: 'move',
          node: node,
          data: data
        })
      }
    })
  }
  onClickDelete(node: FlatNode) {
    this.matDialog.open(DeleteComponent, {
      data: node.data,
      disableClose: true,
    }).afterClosed().pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((ok) => {
      if (typeof ok === "boolean" && ok) {
        this.valChange.next({
          what: 'delete',
          node: node,
        })
      }
    })
  }
}
