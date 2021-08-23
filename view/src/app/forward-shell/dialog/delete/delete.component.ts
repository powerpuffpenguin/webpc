import { HttpClient } from '@angular/common/http';
import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ToasterService } from 'angular2-toaster';
import { finalize, takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Closed } from 'src/app/core/utils/closed';
import { ListResult } from '../../list/load_state';
export interface Data {
  result: ListResult
  id: string
}
@Component({
  selector: 'app-delete',
  templateUrl: './delete.component.html',
  styleUrls: ['./delete.component.scss']
})
export class DeleteComponent implements OnInit, OnDestroy {
  disabled = false
  private closed_ = new Closed()
  constructor(@Inject(MAT_DIALOG_DATA) private readonly data_: Data,
    private httpClient: HttpClient,
    private toasterService: ToasterService,
    private matDialogRef: MatDialogRef<DeleteComponent>,
    private i18nService: I18nService,
  ) { }
  get data(): ListResult {
    return this.data_.result
  }
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
    const data = this.data_
    ServerAPI.forward.v1.shell.child(data.result.id).delete(this.httpClient, {
      params: {
        slave_id: data.id,
      }
    }).pipe(
      takeUntil(this.closed_.observable),
      finalize(() => {
        this.disabled = false
      })
    ).subscribe(() => {
      this.toasterService.pop('success', undefined, this.i18nService.get('terminal has been deleted'))
      this.matDialogRef.close(true)
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
}
