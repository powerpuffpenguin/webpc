import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'src/app/core/toaster.service';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ListResult } from '../../list/load_state';
import { Closed } from 'src/app/core/utils/closed';
import { finalize, takeUntil } from 'rxjs/operators';
export interface Data {
  result: ListResult
  id: string
}
@Component({
  selector: 'app-edit',
  templateUrl: './edit.component.html',
  styleUrls: ['./edit.component.scss']
})
export class EditComponent implements OnInit, OnDestroy {
  private closed_ = new Closed()
  constructor(private readonly httpClient: HttpClient,
    private readonly toasterService: ToasterService,
    private readonly i18nService: I18nService,
    private readonly matDialogRef: MatDialogRef<EditComponent>,
    @Inject(MAT_DIALOG_DATA) private readonly data_: Data,
  ) {
  }
  disabled = false
  name = ''
  get data(): ListResult {
    return this.data_.result
  }
  ngOnInit(): void {
    this.name = this.data.name
  }
  ngOnDestroy() {
    this.closed_.close()
  }
  get isNotChanged(): boolean {
    return this.name.trim() == this.data.name.trim()
  }
  onSubmit() {
    const name = this.name.trim()
    const data = this.data_
    const result = data.result
    ServerAPI.forward.v1.shell.child('rename').post(this.httpClient, {
      id: result.id,
      name: name,
    }, {
      params: {
        slave_id: data.id,
      }
    }).pipe(
      takeUntil(this.closed_.observable),
      finalize(() => {
        this.disabled = false
      })
    ).subscribe(() => {
      result.name = name
      this.toasterService.pop('success', undefined, this.i18nService.get('change terminal completed'))
      this.matDialogRef.close(true)
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
  onClose() {
    this.matDialogRef.close()
  }
}
