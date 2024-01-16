import { HttpClient } from '@angular/common/http';
import { Component, Inject, OnDestroy, OnInit } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ToasterService } from 'src/app/core/toaster.service';
import { finalize, takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { GroupService } from 'src/app/core/group/group.service';
import { Keys, NestedNode, Element } from 'src/app/core/group/tree';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Closed } from 'src/app/core/utils/closed';
import { Data } from '../../query/query';

@Component({
  selector: 'app-group',
  templateUrl: './group.component.html',
  styleUrls: ['./group.component.scss']
})
export class GroupComponent implements OnInit, OnDestroy {
  err: any
  items: Array<NestedNode> = []
  checked: Element | undefined
  disabled = false
  private closed_ = new Closed()
  constructor(@Inject(MAT_DIALOG_DATA) public readonly data: Data,
    private matDialogRef: MatDialogRef<GroupComponent>,
    private readonly groupService: GroupService,
    private readonly httpClient: HttpClient,
    private readonly toasterService: ToasterService,
    private readonly i18nService: I18nService,
  ) {
    this.checked = new Element(data.parent, '')
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
  ngOnDestroy() {
    this.closed_.close()
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
  get isNotChanged(): boolean {
    return this.checked?.id == this.data.parent
  }
  onSubmit() {
    if (this.disabled) {
      return
    }
    const parent = this.checked?.id
    if (!parent) {
      return
    } else if (parent == this.data.parent) {
      this.matDialogRef.close()
      return
    }
    const name = this.checked?.name
    this.disabled = true
    ServerAPI.v1.slaves.child('group', this.data.id).post(this.httpClient, {
      parent: parent,
    }).pipe(
      takeUntil(this.closed_.observable),
      finalize(() => {
        this.disabled = false
      })
    ).subscribe(() => {
      this.data.parent = parent
      this.data.parentName = name

      this.toasterService.pop('success', undefined, this.i18nService.get('change group successed'))
      this.matDialogRef.close()
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
}
