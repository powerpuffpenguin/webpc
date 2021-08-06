import { HttpClient } from '@angular/common/http';
import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ToasterService } from 'angular2-toaster';
import { finalize, takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { KeysService } from 'src/app/core/group/keys.service';
import { Element } from 'src/app/core/group/tree';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Closed } from 'src/app/core/utils/closed';
import { TreeSelectComponent } from 'src/app/shared/tree-select/tree-select.component';
import { Data } from '../../query/query';
interface Response {
  id: string
  code: string
}
interface InjectData {
  onAdded(data: Data): void
}
class GroupData {
  id: string = '1'
  name: string = ''
}
@Component({
  selector: 'app-add',
  templateUrl: './add.component.html',
  styleUrls: ['./add.component.scss']
})
export class AddComponent implements OnInit, OnDestroy {
  disabled = false
  name = 'new'
  description = ''
  private parent_ = new GroupData()
  private closed_ = new Closed()
  constructor(@Inject(MAT_DIALOG_DATA) private readonly data_: InjectData,
    private readonly httpClient: HttpClient,
    private readonly toasterService: ToasterService,
    private readonly matDialogRef: MatDialogRef<AddComponent>,
    private readonly i18nService: I18nService,
    private readonly keysService: KeysService,
    private readonly matDialog: MatDialog,
  ) { }

  ngOnInit(): void {
    this.keysService.promise.then((keys) => {
      if (this.closed_.isClosed) {
        return
      }
      const ele = keys.get(this.parent_.id)
      if (ele) {
        this.parent_.name = ele.name
      }
    }).catch((e) => {
      if (this.closed_.isClosed) {
        return
      }
      console.warn(e)
    })
  }
  ngOnDestroy() {
    this.closed_.close()
  }
  onClose() {
    this.matDialogRef.close()
  }
  get parentName(): string {
    return this.keysService.parentName(this.parent_.id, this.parent_.name)
  }
  onClickSelect() {
    this.matDialog.open(TreeSelectComponent, {
      data: new Element(this.parent_.id, this.parent_.name)
    }).afterClosed().pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((v) => {
      if (v) {
        this.parent_.id = v.id
        this.parent_.name = v.name
      }
    })
  }
  onSubmit() {
    if (this.disabled) {
      return
    }
    this.disabled = true
    const name = this.name.trim()
    const description = this.description.trim()
    const pid = this.parent_.id
    const pname = this.parent_.name
    ServerAPI.v1.slaves.post<Response>(this.httpClient, {
      parent: pid,
      name: name,
      description: description,
    }).pipe(
      takeUntil(this.closed_.observable),
      finalize(() => {
        this.disabled = false
      })
    ).subscribe((response) => {
      this.data_?.onAdded({
        id: response.id,
        name: name,
        description: description,
        code: response.code,
        parent: pid,
        parentName: pname,
      })
      this.toasterService.pop('success', undefined, this.i18nService.get('add device successed'))
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
}
