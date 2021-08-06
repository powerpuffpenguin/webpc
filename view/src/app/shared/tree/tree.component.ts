import { FlatTreeControl } from '@angular/cdk/tree';
import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { MatTreeFlatDataSource, MatTreeFlattener } from '@angular/material/tree';
import { takeUntil } from 'rxjs/operators';
import { Closed } from 'src/app/core/utils/closed';
import { environment } from 'src/environments/environment';
import { NestedNode, Helper, FlatNode, Element } from '../../core/group/tree';

@Component({
  selector: 'shared-tree',
  templateUrl: './tree.component.html',
  styleUrls: ['./tree.component.scss'],
  providers: [Helper],
})
export class TreeComponent implements OnInit {
  @Input()
  disabled = false
  @Input("data")
  set data(source: Array<NestedNode>) {
    this.helper.dataChange.next(source)
  }
  @Output() valChange = new EventEmitter<Element>();

  private checked_ = ''
  @Input()
  set checked(id: string) {
    if (this.checked_ == id) {
      return
    }
    const node = this.keys_.get(id)
    this.checked_ = id
    if (node) {
      this.valChange.next(node.data)
    }
  }
  get checked(): string {
    return this.checked_
  }
  private closed_ = new Closed()
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
  hasChild = (_: number, node: FlatNode) => {
    return node.expandable
  }
  constructor(
    private readonly helper: Helper,
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
  }
  ngOnDestroy() {
    this.closed_.close()
  }
  getName(data: Element): string {
    if (environment.production) {
      return data.name
    }
    return `${data.name} -> ${data.id}`
  }
}
