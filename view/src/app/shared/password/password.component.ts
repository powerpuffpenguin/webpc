import { Component, OnInit, OnDestroy } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { MatDialogRef } from '@angular/material/dialog';
import { ToasterService } from 'src/app/core/toaster.service';
import { Closed } from 'src/app/core/utils/closed';
import { ServerAPI } from 'src/app/core/core/api';
import { finalize, takeUntil } from 'rxjs/operators';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Manager } from 'src/app/core/session/session';
import { md5String } from 'src/app/core/utils/md5';
interface Response {
  changed: boolean
}
@Component({
  selector: 'app-password',
  templateUrl: './password.component.html',
  styleUrls: ['./password.component.scss']
})
export class PasswordComponent implements OnInit, OnDestroy {
  disabled = false
  old = ''
  val = ''
  private closed_ = new Closed()
  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private matDialogRef: MatDialogRef<PasswordComponent>,
    private i18nService: I18nService,
  ) { }
  ngOnInit(): void {
  }
  ngOnDestroy() {
    this.closed_.close()
  }
  onSave() {
    this.disabled = true
    ServerAPI.v1.sessions.child('password').post<Response>(this.httpClient,
      {
        'old': md5String(this.old),
        'password': md5String(this.val),
      },
    ).pipe(
      takeUntil(this.closed_.observable),
      finalize(() => {
        this.disabled = false
      })
    ).subscribe((resp) => {
      if (resp.changed) {
        this.toasterService.pop('success', undefined, this.i18nService.get('password changed'))
        const session = Manager.instance.session
        if (session) {
          Manager.instance.clear(session)
        }
        this.matDialogRef.close()
      } else {
        this.toasterService.pop('error', undefined, 'not matched or not changed')
      }
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
  onClose() {
    this.matDialogRef.close()
  }
}
