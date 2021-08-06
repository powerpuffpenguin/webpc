import { Component, Inject, OnInit } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { GroupService } from 'src/app/core/group/group.service';
import { Element, Keys, NestedNode } from 'src/app/core/group/tree';

@Component({
  selector: 'shared-tree-select',
  templateUrl: './tree-select.component.html',
  styleUrls: ['./tree-select.component.scss']
})
export class TreeSelectComponent implements OnInit {
  err: any
  items: Array<NestedNode> = []
  checked: Element | undefined
  constructor(
    @Inject(MAT_DIALOG_DATA) public readonly data: Element,
    private matDialogRef: MatDialogRef<TreeSelectComponent>,
    private readonly groupService: GroupService,
  ) {
    this.checked = data
  }

  ngOnInit(): void {
    this.groupService.promise.then((items) => {
      const keys = new Keys(items)
      const root = keys.createNested()
      this.items = [root]
    }).catch((e) => {
      this.err = e
    })
  }
  onClose() {
    this.matDialogRef.close()
  }
  onChanged(parent: Element) {
    this.checked = parent
  }
  get checkedID(): string {
    return this.checked?.id ?? ''
  }
  onSubmit() {
    if (this.checked?.id == this.data.id) {
      this.matDialogRef.close()
    } else {
      this.matDialogRef.close(this.checked)
    }
  }
}
