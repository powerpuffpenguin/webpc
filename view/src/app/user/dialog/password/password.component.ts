import { HttpClient } from '@angular/common/http';
import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ToasterService } from 'angular2-toaster';
import { finalize, takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Closed } from 'src/app/core/utils/closed';
import { Data } from '../../query/query';
import { md5String } from 'src/app/core/utils/md5';

@Component({
  selector: 'app-password',
  templateUrl: './password.component.html',
  styleUrls: ['./password.component.scss']
})
export class PasswordComponent implements OnInit, OnDestroy {
  disabled = false
  val = ''
  private closed_ = new Closed()
  constructor(@Inject(MAT_DIALOG_DATA) public readonly data: Data,
    private httpClient: HttpClient,
    private toasterService: ToasterService,
    private matDialogRef: MatDialogRef<PasswordComponent>,
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
  onSubmit() {
    if (this.disabled) {
      return
    }
    this.disabled = true
    ServerAPI.v1.users.child('password', this.data.id).post(this.httpClient,
      {
        'value': md5String(this.val),
      },
    ).pipe(
      takeUntil(this.closed_.observable),
      finalize(() => {
        this.disabled = false
      })
    ).subscribe(() => {
      this.toasterService.pop('success', undefined, this.i18nService.get('password changed'))
      this.matDialogRef.close()
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
}
