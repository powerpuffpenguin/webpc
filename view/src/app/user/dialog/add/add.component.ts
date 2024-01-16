import { HttpClient } from '@angular/common/http';
import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ToasterService } from 'src/app/core/toaster.service';
import { finalize, takeUntil } from 'rxjs/operators';
import { Authorization, AuthorizationName, Authorizations, ServerAPI } from 'src/app/core/core/api';
import { KeysService } from 'src/app/core/group/keys.service';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Closed } from 'src/app/core/utils/closed';
import { md5String } from 'src/app/core/utils/md5';
import { TreeSelectComponent } from 'src/app/shared/tree-select/tree-select.component';
import { Data } from '../../query/query';
import { Element } from 'src/app/core/group/tree';
interface Response {
  id: string
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
export class AddComponent implements OnInit {
  disabled = false
  name = ''
  nickname = ''
  password = ''
  readonly authorizations = Authorizations
  set = new Set<number>()
  private parent_ = new GroupData()
  private closed_ = new Closed()
  constructor(@Inject(MAT_DIALOG_DATA) private readonly data_: InjectData,
    private httpClient: HttpClient,
    private toasterService: ToasterService,
    private matDialogRef: MatDialogRef<AddComponent>,
    private i18nService: I18nService,
    private readonly keysService: KeysService,
    private readonly matDialog: MatDialog,
  ) {
    this.set.add(Authorization.Shell)
    this.set.add(Authorization.Read)
    this.set.add(Authorization.Write)
    this.set.add(Authorization.VNC)
  }
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
  authorizationName(authorization: Authorization): string {
    return AuthorizationName(authorization)
  }
  onCheckedChange(authorization: Authorization, evt: MatCheckboxChange) {
    if (evt.checked) {
      this.set.add(authorization)
    } else {
      this.set.delete(authorization)
    }
  }
  onSubmit() {
    if (this.disabled) {
      return
    }
    this.disabled = true
    const authorization = new Array<number>()
    this.set.forEach((k) => {
      authorization.push(k)
    })
    authorization.sort()
    const password = md5String(this.password)
    const pid = this.parent_.id
    const pname = this.parent_.name
    ServerAPI.v1.users.post<Response>(this.httpClient, {
      parent: pid,
      name: this.name,
      password: password,
      nickname: this.nickname,
      authorization: authorization,
    }).pipe(
      takeUntil(this.closed_.observable),
      finalize(() => {
        this.disabled = false
      })
    ).subscribe((response) => {
      this.data_?.onAdded({
        id: response.id,
        name: this.name,
        nickname: this.nickname,
        authorization: authorization,
        parent: pid,
        parentName: pname,
      })
      this.toasterService.pop('success', undefined, this.i18nService.get('add user successed'))
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
}
