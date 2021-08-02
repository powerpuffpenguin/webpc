import { HttpClient } from '@angular/common/http';
import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ToasterService } from 'angular2-toaster';
import { finalize, takeUntil } from 'rxjs/operators';
import { Authorization, AuthorizationName, Authorizations, ServerAPI } from 'src/app/core/core/api';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Closed } from 'src/app/core/utils/closed';
import { md5String } from 'src/app/core/utils/md5';
import { Data } from '../../query/query';
interface Response {
  id: string
}
interface InjectData {
  onAdded(data: Data): void
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
  private closed_ = new Closed()
  constructor(@Inject(MAT_DIALOG_DATA) private readonly data_: InjectData,
    private httpClient: HttpClient,
    private toasterService: ToasterService,
    private matDialogRef: MatDialogRef<AddComponent>,
    private i18nService: I18nService,
  ) { }
  ngOnInit(): void {
  }
  ngOnDestroy() {
    this.closed_.close()
  }
  onClose() {
    this.matDialogRef.close()
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
    ServerAPI.v1.users.post<Response>(this.httpClient, {
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
      })
      this.toasterService.pop('success', undefined, this.i18nService.get('add user successed'))
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
}
