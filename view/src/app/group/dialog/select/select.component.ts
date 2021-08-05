import { HttpClient } from '@angular/common/http';
import { Component, OnInit, Inject, OnDestroy } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ToasterService } from 'angular2-toaster';
import { finalize, takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Closed } from 'src/app/core/utils/closed';
import { Element, NestedNode } from '../../../shared/tree/tree';

@Component({
  selector: 'app-select',
  templateUrl: './select.component.html',
  styleUrls: ['./select.component.scss']
})
export class SelectComponent implements OnInit, OnDestroy {
  private closed_ = new Closed()
  disabled = false
  parent: Element | undefined
  root: Array<NestedNode> = []
  constructor(@Inject(MAT_DIALOG_DATA) public readonly data: Element,
    private matDialogRef: MatDialogRef<SelectComponent>,
    private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
  ) {
  }

  ngOnInit(): void {
    let node = this.data
    while (node.parent) {
      node = node.parent
    }
    const root = new NestedNode(node)
    this._appendChildren(root)

    this.root = [root]
  }
  private _appendChildren(node: NestedNode) {
    const ele = node.data
    ele.children.forEach((child) => {
      if (child == this.data) {
        return
      }
      const nodeChild = new NestedNode(child)
      nodeChild.parent = node
      node.children.push(nodeChild)
      this._appendChildren(nodeChild)
    })
    node.children.sort(NestedNode.compareFn)
  }
  ngOnDestroy() {
    this.closed_.close()
  }
  onClose() {
    this.matDialogRef.close()
  }
  get isNotChanged(): boolean {
    return !this.parent || this.parent == this.data.parent
  }
  onChanged(parent: Element) {
    this.parent = parent
  }
  onSubmit() {
    if (this.disabled || this.isNotChanged) {
      return
    }
    this.disabled = true
    ServerAPI.v1.groups.child('move', this.data.id).post(this.httpClient, {
      parent: this.parent?.id,
    }).pipe(
      takeUntil(this.closed_.observable),
      finalize(() => {
        this.disabled = false
      })
    ).subscribe(() => {
      this.toasterService.pop('success', undefined, this.i18nService.get('group move changed'))

      this.matDialogRef.close(this.parent)
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
}
