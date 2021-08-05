import { HttpClient } from '@angular/common/http';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { ToasterService } from 'angular2-toaster';
import { takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { Closed } from 'src/app/core/utils/closed';
import { FlatTreeControl } from '@angular/cdk/tree';
import { MatTreeFlatDataSource, MatTreeFlattener } from '@angular/material/tree';
import { Tree, NetElement, Helper, FlatNode, NestedNode } from 'src/app/shared/tree/tree';
import { NodeEvent } from './node/node.component';

interface ListResponse {
  items: Array<NetElement>
}
@Component({
  selector: 'app-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.scss'],
  providers: [Helper],
})
export class ListComponent implements OnInit, OnDestroy {
  private closed_ = new Closed()
  err: any
  ready = false
  tree: Tree | undefined
  private keys_ = new Map<string, NestedNode>()
  private flatNodeMap_ = new Map<FlatNode, NestedNode>()
  private nestedNodeMap_ = new Map<NestedNode, FlatNode>()
  private _transformer = (node: NestedNode, level: number) => {
    this.keys_.set(node.data.id, node)
    const nestedNodeMap = this.nestedNodeMap_
    const flatNodeMap = this.flatNodeMap_
    const existingNode = nestedNodeMap.get(node)
    if (existingNode) {
      if (existingNode.update) {
        flatNodeMap.delete(existingNode)
      } else {
        return existingNode
      }
    }
    const flatNode = new FlatNode(node.data)
    flatNode.level = level
    flatNode.expandable = !!node.children && node.children.length > 0
    if (flatNode.expandable) {
      this.treeControl.expand(flatNode)
    }
    flatNodeMap.set(flatNode, node)
    nestedNodeMap.set(node, flatNode)

    return flatNode
  }
  treeControl = new FlatTreeControl<FlatNode>(node => node.level, node => node.expandable)
  private treeFlattener_ = new MatTreeFlattener(this._transformer, node => node.level, node => node.expandable, node => node.children)
  dataSource = new MatTreeFlatDataSource(this.treeControl, this.treeFlattener_)
  constructor(
    private readonly helper: Helper,
    private readonly httpClient: HttpClient,
    private readonly toasterService: ToasterService,
  ) {
  }

  ngOnInit(): void {
    this.helper.dataChange.pipe(
      takeUntil(this.closed_.observable)
    ).subscribe(
      (data) => {
        this.dataSource.data = data
      }
    )
    this.onClickRefresh()
  }
  ngOnDestroy() {
    this.closed_.close()
  }
  onClickRefresh() {
    this.err = undefined
    ServerAPI.v1.groups.get<ListResponse>(this.httpClient).pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((resp) => {
      try {
        const tree = new Tree(this.helper, resp.items)
        this.tree = tree
        this.ready = true
      } catch (e) {
        this.err = e
      }
    }, (e) => {
      this.err = e
    })
  }
  hasChild = (_: number, node: FlatNode) => {
    return node.expandable
  }
  onNodeChanged(evt: NodeEvent) {
    if (this.closed_.isClosed) {
      return
    }
    const tree = this.tree
    if (!tree) {
      return
    }
    const node = this.flatNodeMap_.get(evt.node)
    if (!node) {
      return
    }
    try {
      if (evt.what === 'add') {
        const data = evt.data
        if (node && data) {
          evt.node.update = true
          tree.add(node, data)
          this.treeControl.expand(evt.node)
        }
      } else if (evt.what === 'changed') {
        const parent = node.parent
        if (parent) {
          parent.children.sort(NestedNode.compareFn)
          this.helper.updateView()
        }
      } else if (evt.what === 'delete') {
        const parent = node.parent
        if (!parent) {
          throw new Error(`parent not exists: ${node.data.parent?.id}`)
        }
        const index = parent.children.indexOf(node)
        parent.children.splice(index, 1)
        this._remove(node)
        const pf = this.nestedNodeMap_.get(parent)
        if (pf) {
          pf.update = true
        }
        this.helper.updateView()
      } else if (evt.what === 'move') {
        const parent = evt.data
        if (!parent) {
          throw new Error(`parent null`)
        }
        this.tree?.move(node.data, parent)
        const p = this.keys_.get(parent.id)
        if (p && this._move(node, p)) {
          this._setUpdate(evt.node)
          this.helper.updateView()
        }
      }
    } catch (e) {
      this.toasterService.pop('error', undefined, e)
    }
  }
  private _setUpdate(flat: FlatNode) {
    flat.update = true
    const node = this.flatNodeMap_.get(flat)
    node?.children.forEach((node) => {
      const flat = this.nestedNodeMap_.get(node)
      if (flat) {
        this._setUpdate(flat)
      }
    })
  }
  private _move(node: NestedNode, parent: NestedNode): boolean {
    const po = node.parent
    if (!po) {
      return false
    }
    if (po == parent) {
      return false
    }

    // rm from old
    let children = po.children
    children.splice(children.indexOf(node), 1)
    let fp = this.nestedNodeMap_.get(po)
    if (fp) {
      fp.update = true
    }
    // add to new
    fp = this.nestedNodeMap_.get(parent)
    if (fp) {
      fp.update = true
    }
    node.parent = parent
    children = parent.children
    children.push(node)
    children.sort(NestedNode.compareFn)
    return true
  }
  private _remove(node: NestedNode) {
    this.keys_.delete(node.data.id)
    const flat = this.nestedNodeMap_.get(node)
    if (flat) {
      this.flatNodeMap_.delete(flat)
    }
    this.nestedNodeMap_.delete(node)
    node.children.forEach((node) => {
      this._remove(node)
    })
  }
}
