import { HttpClient } from '@angular/common/http';
import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ToasterService } from 'src/app/core/toaster.service';
import { finalize, takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Closed } from 'src/app/core/utils/closed';
import { Data } from '../../query/query';
interface Response {
  changed: boolean
  code: string
}
@Component({
  selector: 'app-code',
  templateUrl: './code.component.html',
  styleUrls: ['./code.component.scss']
})
export class CodeComponent implements OnInit, OnDestroy {
  disabled = false
  private closed_ = new Closed()
  constructor(@Inject(MAT_DIALOG_DATA) public readonly data: Data,
    private httpClient: HttpClient,
    private toasterService: ToasterService,
    private matDialogRef: MatDialogRef<CodeComponent>,
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
    ServerAPI.v1.slaves.child('code', this.data.id).post<Response>(this.httpClient,
      undefined,
    ).pipe(
      takeUntil(this.closed_.observable),
      finalize(() => {
        this.disabled = false
      })
    ).subscribe((resp) => {
      if (resp.changed) {
        this.data.code = resp.code
      }
      this.toasterService.pop('success', undefined, this.i18nService.get('device code changed'))
      this.matDialogRef.close()
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
}
