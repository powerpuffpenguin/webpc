import { HttpClient } from '@angular/common/http';
import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ToasterService } from 'src/app/core/toaster.service';
import { finalize, takeUntil } from 'rxjs/operators';
import { Authorizations, ServerAPI, Authorization, AuthorizationName } from 'src/app/core/core/api';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Closed } from 'src/app/core/utils/closed';
import { Data } from '../../query/query';

@Component({
  selector: 'app-edit',
  templateUrl: './edit.component.html',
  styleUrls: ['./edit.component.scss']
})
export class EditComponent implements OnInit, OnDestroy {
  disabled = false
  private closed_ = new Closed()
  nickname = ''
  readonly authorizations = Authorizations
  private set_ = new Set<number>()
  set = new Set<number>()
  constructor(@Inject(MAT_DIALOG_DATA) public readonly data: Data,
    private httpClient: HttpClient,
    private toasterService: ToasterService,
    private matDialogRef: MatDialogRef<EditComponent>,
    private i18nService: I18nService,
  ) {
    this.nickname = data.nickname
    this.data.authorization?.forEach((v) => {
      if (!this.set.has(v)) {
        this.set.add(v)
        this.set_.add(v)
      }
    })
  }
  ngOnInit(): void {
  }
  ngOnDestroy() {
    this.closed_.close()
  }
  onClose() {
    this.matDialogRef.close()
  }
  private _equal(l: Set<number>, r: Set<number>): boolean {
    if (l.size != r.size) {
      return false
    }
    const itear = l.values()
    while (true) {
      const result = itear.next()
      if (result.done) {
        break
      }
      if (!r.has(result.value)) {
        return false
      }
    }
    return true
  }
  get isNotChanged(): boolean {
    return this.nickname == (this.data.nickname ?? '') &&
      this._equal(this.set, this.set_)
  }
  name(authorization: Authorization): string {
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
    if (this.disabled || this.isNotChanged) {
      return
    }
    this.disabled = true
    const authorization = new Array<number>()
    this.set.forEach((k) => {
      authorization.push(k)
    })
    authorization.sort()
    ServerAPI.v1.users.child('change', this.data.id).post(this.httpClient, {
      nickname: this.nickname,
      authorization: authorization,
    }).pipe(
      takeUntil(this.closed_.observable),
      finalize(() => {
        this.disabled = false
      })
    ).subscribe(() => {
      this.toasterService.pop('success', undefined, this.i18nService.get('user properties changed'))
      this.data.nickname = this.nickname
      this.data.authorization = authorization
      this.matDialogRef.close(true)
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
}
